package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/pkg/types"
)

// OSUpdate（OS更新状況）チェック
// 重みは 20 点
func OSUpdate(ctx context.Context) types.CheckResult {
	const weight = 20

	// 現在のOSバージョン（参考情報としてevidenceに載せる）
	swVersRes := runCommand(ctx, 3*time.Second, "/usr/bin/sw_vers", "-productVersion")
	currentVersion := strings.TrimSpace(swVersRes.Stdout)
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	// 利用可能な更新をチェック（--no-scanでキャッシュを使用、高速化）
	updateRes := runCommand(ctx, 8*time.Second, "/usr/sbin/softwareupdate", "-l", "--no-scan")

	ev := map[string]string{
		"version": currentVersion,
	}

	cr := types.CheckResult{
		ID:       "osupdate",
		Title:    "OS updates current",
		Evidence: ev,
	}

	// softwareupdate -l --no-scan が失敗した場合の代替手段
	if updateRes.Err != nil {
		cr.Status = "warn"
		cr.Score = weight / 2
		cr.Recommendation = "OS更新状況の確認に失敗しました。システム設定 > ソフトウェアアップデート から手動で確認してください"
		return cr
	}

	// 更新が利用可能かチェック
	updateOutput := strings.ToLower(strings.TrimSpace(updateRes.Stdout))

	// 「更新なし」を示す文言をチェック
	noUpdateMarkers := []string{
		"no new software available",
		"no updates available",
		"your mac is up to date",
		"ソフトウェアの更新はありません",
		"最新の状態です",
	}

	if containsAny(updateOutput, noUpdateMarkers) || strings.TrimSpace(updateRes.Stdout) == "" {
		cr.Status = "pass"
		cr.Score = weight
		cr.Recommendation = "OSは最新の状態です。定期的な更新を継続してください"
		return cr
	}

	// 実際に更新項目があるかチェック（* で始まる行があるか）
	hasUpdateItems := strings.Contains(updateRes.Stdout, "*")

	if !hasUpdateItems {
		// 更新項目がない場合は更新なしとして扱う
		cr.Status = "pass"
		cr.Score = weight
		cr.Recommendation = "OSは最新の状態です。定期的な更新を継続してください"
		return cr
	}

	// セキュリティ更新を優先的にチェック
	securityUpdates := extractSecurityUpdates(updateRes.Stdout)
	if len(securityUpdates) > 0 {
		cr.Status = "fail"
		cr.Score = 0
		cr.Recommendation = "セキュリティ更新が利用可能です。システム設定 > ソフトウェアアップデート から更新してください"
		ev["updates"] = strings.Join(securityUpdates, "; ")
	} else {
		// セキュリティ以外の更新のみの場合
		cr.Status = "warn"
		cr.Score = weight / 2
		cr.Recommendation = "OS更新が利用可能です。セキュリティ更新を優先して適用してください"
		ev["updates"] = "一般更新が利用可能"
	}

	return cr
}

// セキュリティ更新を抽出する関数
func extractSecurityUpdates(output string) []string {
	var securityUpdates []string
	lines := strings.Split(output, "\n")

	// セキュリティ関連のキーワード
	securityKeywords := []string{
		"Security Update",
		"セキュリティアップデート",
		"Security",
		"セキュリティ",
		"Critical",
		"重要",
		"Rapid Security Response",
	}

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// タイトル行を探す（通常は * で始まる）
		if strings.HasPrefix(line, "*") {
			title := strings.TrimSpace(strings.TrimPrefix(line, "*"))
			for _, keyword := range securityKeywords {
				if strings.Contains(strings.ToLower(title), strings.ToLower(keyword)) {
					securityUpdates = append(securityUpdates, title)
					break
				}
			}
		}
	}

	return securityUpdates
}

// ヘルパー：いずれかのフレーズを含むか
func containsAny(s string, subs []string) bool {
	for _, sub := range subs {
		if strings.Contains(s, sub) {
			return true
		}
	}
	return false
}
