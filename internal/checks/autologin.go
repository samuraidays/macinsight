package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/pkg/types"
)

// AutoLogin（自動ログイン）状態
// 重みは 10 点
func AutoLogin(ctx context.Context) types.CheckResult {
	const weight = 10

	// 自動ログインの設定を確認
	res := runCommand(ctx, 3*time.Second, "/usr/bin/defaults", "read", "/Library/Preferences/com.apple.loginwindow", "autoLoginUser")
	ev := map[string]string{"autoLoginUser": strings.TrimSpace(res.Stdout)}

	cr := types.CheckResult{
		ID:       "autologin",
		Title:    "Auto-login disabled",
		Evidence: ev,
	}

	if res.Err != nil {
		// エラーの場合、設定ファイルが存在しないか、キーが存在しない可能性
		// この場合は自動ログインが無効とみなす
		cr.Status = "pass"
		cr.Score = weight
		cr.Evidence["note"] = "Auto-login setting not found (disabled by default)"
		return cr
	}

	// 空文字列または "()" の場合は自動ログインが無効
	output := strings.TrimSpace(res.Stdout)
	if output == "" || output == "()" || output == "0" {
		cr.Status = "pass"
		cr.Score = weight
	} else {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "自動ログインを無効にしてください: システム設定 > ユーザとグループ > ログインオプション"
	}

	return cr
}
