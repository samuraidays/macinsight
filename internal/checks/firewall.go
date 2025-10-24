package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/pkg/types"
)

// macOSアプリケーションファイアウォールの有効/無効
// 重みは 10 点
func Firewall(ctx context.Context) types.CheckResult {
	const weight = 10

	res := runCommand(ctx, 3*time.Second, "/usr/libexec/ApplicationFirewall/socketfilterfw", "--getglobalstate")
	ev := map[string]string{"socketfilterfw": strings.TrimSpace(res.Stdout)}

	cr := types.CheckResult{
		ID:       "firewall",
		Title:    "Firewall enabled",
		Evidence: ev,
	}

	if res.Err != nil {
		cr.Status = "unknown"
		cr.Score = weight / 2
		cr.Recommendation = "管理者権限が必要な場合があります"
		return cr
	}

	// "State = 1" or "enabled" を含んでいれば pass
	out := strings.ToLower(res.Stdout)
	if strings.Contains(out, "state = 1") || strings.Contains(out, "enabled") {
		cr.Status = "pass"
		cr.Score = weight
	} else {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "システム設定 > ネットワーク > ファイアウォール を有効化"
	}

	return cr
}
