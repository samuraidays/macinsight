package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
	"github.com/samuraidays/macinsight/pkg/types"
)

// SIP（System Integrity Protection）状態
// 重みは 15 点
func SIP(ctx context.Context) types.CheckResult {
	const weight = 15

	res := executil.Run(ctx, 3*time.Second, "/usr/bin/csrutil", "status")
	ev := map[string]string{"csrutil": strings.TrimSpace(res.Stdout)}

	cr := types.CheckResult{
		ID:       "sip",
		Title:    "System Integrity Protection enabled",
		Evidence: ev,
	}

	if res.Err != nil {
		cr.Status = "unknown"
		cr.Score = weight / 2
		cr.Recommendation = "csrutil の場所/実行可否やOSバージョン差を確認"
		return cr
	}

	// "enabled" を含んでいれば pass とする（言語/表記差を吸収）
	if strings.Contains(strings.ToLower(res.Stdout), "enabled") {
		cr.Status = "pass"
		cr.Score = weight
	} else {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "SIP 有効化にはリカバリモードでの操作が必要"
	}

	return cr
}
