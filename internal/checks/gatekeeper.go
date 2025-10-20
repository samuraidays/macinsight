package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/pkg/types"
)

// Gatekeeper が有効かを spctl --status で確認
// 重みは 20 点（pass=20, fail=0, unknown=10）
func Gatekeeper(ctx context.Context) types.CheckResult {
	const weight = 20

    res := runCommand(ctx, 3*time.Second, "/usr/sbin/spctl", "--status")
	ev := map[string]string{"spctl_status": strings.TrimSpace(res.Stdout)}

	cr := types.CheckResult{
		ID:       "gatekeeper",
		Title:    "Gatekeeper enabled",
		Evidence: ev,
	}

	// 実行エラー時は unknown
	if res.Err != nil {
		cr.Status = "unknown"
		cr.Score = weight / 2
		cr.Recommendation = "spctl の実行権限/パスやOSバージョン差を確認"
		return cr
	}

	// 出力に "assessments enabled" を含むかで判定
	if strings.Contains(res.Stdout, "assessments enabled") {
		cr.Status = "pass"
		cr.Score = weight
	} else {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "システム設定 > プライバシーとセキュリティ > App のダウンロード元 を制限"
	}

	return cr
}
