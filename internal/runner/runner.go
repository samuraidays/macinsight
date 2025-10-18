package runner

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/samuraidays/macinsight/internal/checks"
	"github.com/samuraidays/macinsight/pkg/types"
)

// CLIから渡す実行オプション
type Option struct {
	Only    map[string]struct{} // 実行するチェックを限定（空なら全件）
	Exclude map[string]struct{} // 実行しないチェック
	Timeout time.Duration       // チェックごとのタイムアウト
}

// 監査実行（並列に各チェックを走らせ、スコア集計して返す）
func Run(version string, opt Option) types.Report {
	host := types.HostInfo{
		Hostname: hostname(),
		OS:       osinfo(),
	}

	// v0.1: 代表4チェックのみ登録
	registry := map[string]func(context.Context) types.CheckResult{
		"gatekeeper": checks.Gatekeeper,
		"filevault":  checks.FileVault,
		"sip":        checks.SIP,
		"firewall":   checks.Firewall,
	}

	results := make([]types.CheckResult, 0, len(registry))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for id, fn := range registry {
		// --only が指定されたらその集合にあるものだけ
		if len(opt.Only) > 0 {
			if _, ok := opt.Only[id]; !ok {
				continue
			}
		}
		// --exclude は除外
		if _, skip := opt.Exclude[id]; skip {
			continue
		}

		wg.Add(1)
		go func(id string, fn func(context.Context) types.CheckResult) {
			defer wg.Done()
			// 各チェックに個別タイムアウトを適用
			cctx, cancel := context.WithTimeout(context.Background(), opt.Timeout)
			defer cancel()
			cr := fn(cctx)
			mu.Lock()
			results = append(results, cr)
			mu.Unlock()
		}(id, fn)
	}

	wg.Wait()

	// 合計スコア（上限100）
	total := 0
	for _, r := range results {
		total += r.Score
	}
	if total > 100 {
		total = 100
	}

	return types.Report{
		Version: version,
		Host:    host,
		Score:   total,
		Checks:  results,
	}
}

func hostname() string {
	h, _ := os.Hostname()
	return h
}

// sw_vers を呼んで OS 情報を得る（失敗は空値で返す）
func osinfo() types.OSInfo {
	out, _ := exec.Command("/usr/bin/sw_vers").Output()
	s := string(out)
	return types.OSInfo{
		Product: "macOS",
		Version: lineValue(s, "ProductVersion:"),
		Build:   lineValue(s, "BuildVersion:"),
	}
}

func lineValue(s, key string) string {
	for _, l := range strings.Split(s, "\n") {
		lt := strings.TrimSpace(l)
		if strings.HasPrefix(lt, key) {
			f := strings.Fields(lt)
			if len(f) > 1 {
				return f[len(f)-1]
			}
		}
	}
	return ""
}
