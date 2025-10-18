package checks

import (
	"context"
	"strings"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
	"github.com/samuraidays/macinsight/pkg/types"
)

// OSUpdate（OS更新状況）チェック
// 重みは 20 点
func OSUpdate(ctx context.Context) types.CheckResult {
	const weight = 20

	// 現在のOSバージョン（参考情報としてevidenceに載せる）
	swVersRes := executil.Run(ctx, 3*time.Second, "/usr/bin/sw_vers", "-productVersion")
	currentVersion := strings.TrimSpace(swVersRes.Stdout)
	if currentVersion == "" {
		currentVersion = "unknown"
	}

	// 更新確認（非0終了 & stderr出力でも「更新なし」が含まれるケースがある）
	updateRes := executil.Run(ctx, 25*time.Second, "/usr/sbin/softwareupdate", "-l")

	// 判定は stdout と stderr を結合して内容で行う
	rawOut := strings.TrimSpace(updateRes.Stdout)
	rawErr := strings.TrimSpace(updateRes.Stderr)
	combined := strings.TrimSpace(rawOut + "\n" + rawErr)
	lc := strings.ToLower(combined)

	ev := map[string]string{
		"currentVersion": currentVersion,
		"softwareupdate": rawOut,
	}
	if rawErr != "" {
		ev["softwareupdate_stderr"] = rawErr
	}
	if updateRes.Err != nil {
		ev["softwareupdate_error"] = updateRes.Err.Error()
	}

	cr := types.CheckResult{
		ID:       "osupdate",
		Title:    "OS updates current",
		Evidence: ev,
	}

	// 1) 「更新なし」を示す代表的な文言（英/日・表記揺れ）を先に評価
	//    - 例: "No new software available", "No updates available",
	//          "Your Mac is up to date", "ソフトウェアの更新はありません", "最新の状態です"
	noUpdateMarkers := []string{
		"no new software available",
		"no updates available",
		"your mac is up to date",
		"up to date", // 他文脈との誤検知は稀。先に利用可能判定をするなら外してもOK
		"ソフトウェアの更新はありません",
		"最新の状態です",
	}

	if containsAny(lc, noUpdateMarkers) || combined == "" {
		cr.Status = "pass"
		cr.Score = weight
		return cr
	}

	// 2) 「利用可能な更新」を示す文言
	//    - 例: "recommended updates", lines starting with "*", "available"
	availableMarkers := []string{
		"recommended updates",
		"available",
		"updates are available",
		"アップデートが利用できます",
		"利用可能",
		"*", // softwareupdate -l の項目は行頭に "*" が付く
	}

	updatesAvailable := containsAny(lc, availableMarkers)

	// 3) セキュリティ更新の有無を抽出（タイトルや説明に含まれるキーワード）
	securityUpdates := extractSecurityUpdates(combined)

	if updatesAvailable {
		if len(securityUpdates) > 0 {
			// セキュリティ更新が含まれる → fail（最優先）
			cr.Status = "fail"
			cr.Score = 0
			cr.Recommendation = "セキュリティ更新が利用可能です。システム設定 > ソフトウェアアップデート から更新してください。"
			cr.Evidence["securityUpdates"] = strings.Join(securityUpdates, "; ")
			return cr
		}
		// セキュリティ以外の更新 → warn
		cr.Status = "warn"
		cr.Score = weight / 2
		cr.Recommendation = "OS更新が利用可能です。セキュリティ更新が含まれる場合は優先して適用してください。"
		return cr
	}

	// 4) ここまででどちらにも当てはまらない場合：
	//    - 実機で文言が変わることがあるため、unknownではなく warn に寄せると運用が安定
	cr.Status = "warn"
	cr.Score = weight / 2
	cr.Recommendation = "更新状況が判別できませんでした。GUIの「システム設定 > ソフトウェアアップデート」を確認してください。"
	return cr
}

// セキュリティ更新を抽出する関数
// softwareupdate -l の出力は、項目行が "*" で始まり、タイトルや説明が続く形式。
// stdout/stderr両方から渡された combined テキストを対象に、代表キーワードで判定する。
func extractSecurityUpdates(combined string) []string {
	var securityUpdates []string
	lines := strings.Split(combined, "\n")

	// セキュリティ関連の代表キーワード（英/日）
	securityKeywords := []string{
		"security response", // Rapid Security Response
		"security update",
		"セキュリティアップデート",
		"critical",
		"重要",
	}

	for _, line := range lines {
		l := strings.TrimSpace(line)
		if l == "" {
			continue
		}
		// 項目行（先頭が "*"）を主に対象にするが、説明行に含まれる場合も拾う
		if strings.HasPrefix(l, "*") || true {
			low := strings.ToLower(l)
			for _, kw := range securityKeywords {
				if strings.Contains(low, strings.ToLower(kw)) {
					// 目視しやすいよう "*" は外して保存
					title := strings.TrimSpace(strings.TrimPrefix(l, "*"))
					if title == "" {
						title = l
					}
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
