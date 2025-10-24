package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/internal/output"
	"github.com/samuraidays/macinsight/internal/runner"
	"github.com/samuraidays/macinsight/internal/schema"
)

// ldflags で埋め込む用（go build -ldflags "-X main.version=v0.1.0"）
var version = "v0.1.0"

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	// サブコマンド：audit / list-checks / version / schema
	switch os.Args[1] {
	case "audit":
		runAudit(os.Args[2:])
	case "list-checks":
		fmt.Println("sip,gatekeeper,filevault,firewall,autologin,osupdate")
	case "version":
		fmt.Println(version)
	case "schema":
		runSchema(os.Args[2:])
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
  macinsight schema [--output <file>]

Examples:
  macinsight audit
  macinsight audit --json --only filevault,gatekeeper
  macinsight schema --output schema.json
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

func runSchema(args []string) {
	// フラグ定義
	fs := flag.NewFlagSet("schema", flag.ExitOnError)
	var outputFile string
	fs.StringVar(&outputFile, "output", "", "output file path (default: stdout)")
	_ = fs.Parse(args)

	// スキーマ生成
	generator := &schema.JSONSchemaGenerator{}

	if outputFile != "" {
		// ファイルに出力
		file, err := os.Create(outputFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating file %s: %v\n", outputFile, err)
			os.Exit(1)
		}
		defer func() {
			if closeErr := file.Close(); closeErr != nil {
				fmt.Fprintf(os.Stderr, "Error closing file %s: %v\n", outputFile, closeErr)
			}
		}()

		if err := generator.WriteSchema(file); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("JSON schema written to %s\n", outputFile)
	} else {
		// 標準出力に出力
		if err := generator.WriteSchema(os.Stdout); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing schema: %v\n", err)
			os.Exit(1)
		}
	}
}
