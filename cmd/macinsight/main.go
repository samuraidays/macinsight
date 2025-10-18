package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/internal/output"
	"github.com/samuraidays/macinsight/internal/runner"
)

// ldflags で埋め込む用（go build -ldflags "-X main.version=v0.1.0"）
var version = "v0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	// サブコマンド：audit / list-checks / version
	switch os.Args[1] {
	case "audit":
		runAudit(os.Args[2:])
	case "list-checks":
		fmt.Println("sip,gatekeeper,filevault,firewall")
	case "version":
		fmt.Println(version)
	default:
		usage()
	}
}

func usage() {
	fmt.Print(`macinsight - macOS Security Audit CLI

Usage:
  macinsight audit [--json] [--only <checks>] [--exclude <checks>] [--timeout 3s]
  macinsight list-checks
  macinsight version

Examples:
  macinsight audit
  macinsight audit --json --only filevault,gatekeeper
`)
}

func runAudit(args []string) {
	// フラグ定義
	fs := flag.NewFlagSet("audit", flag.ExitOnError)
	var asJSON bool
	var only, exclude string
	var timeout time.Duration
	fs.BoolVar(&asJSON, "json", false, "print JSON")
	fs.StringVar(&only, "only", "", "comma-separated checks to include")
	fs.StringVar(&exclude, "exclude", "", "comma-separated checks to skip")
	fs.DurationVar(&timeout, "timeout", 3*time.Second, "per-check timeout")
	_ = fs.Parse(args)

	// 実行オプションを作成
	opt := runner.Option{
		Only:    toSet(only),
		Exclude: toSet(exclude),
		Timeout: timeout,
	}

	// 監査の実行
	rep := runner.Run(version, opt)

	// 出力モード
	var err error
	if asJSON {
		err = output.WriteJSON(os.Stdout, rep)
	} else {
		err = output.WriteTable(os.Stdout, rep)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

func toSet(csv string) map[string]struct{} {
	m := map[string]struct{}{}
	if strings.TrimSpace(csv) == "" {
		return m
	}
	for _, v := range strings.Split(csv, ",") {
		v = strings.TrimSpace(v)
		if v != "" {
			m[v] = struct{}{}
		}
	}
	return m
}
