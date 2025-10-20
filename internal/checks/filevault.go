package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/pkg/types"
)

// FileVault（フルディスク暗号化）の有効/無効
// 重みは 20 点（重要度高）
func FileVault(ctx context.Context) types.CheckResult {
	const weight = 20

    res := runCommand(ctx, 3*time.Second, "/usr/bin/fdesetup", "status")
	ev := map[string]string{"fdesetup": strings.TrimSpace(res.Stdout)}

	cr := types.CheckResult{
		ID:       "filevault",
		Title:    "FileVault enabled",
		Evidence: ev,
	}

	if res.Err != nil {
		cr.Status = "unknown"
		cr.Score = weight / 2
		cr.Recommendation = "管理者権限が必要な場合があります"
		return cr
	}

	if strings.Contains(res.Stdout, "FileVault is On") {
		cr.Status = "pass"
		cr.Score = weight
	} else {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "FileVault 有効化を検討（システム設定 > プライバシーとセキュリティ > FileVault）"
	}

	return cr
}
