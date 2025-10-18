package types

// 各チェックの結果を表す構造体
type CheckResult struct {
	ID             string            `json:"id"`                       // 例: "gatekeeper"
	Title          string            `json:"title"`                    // 例: "Gatekeeper enabled"
	Status         string            `json:"status"`                   // "pass" | "fail" | "unknown"
	Score          int               `json:"score"`                    // このチェックに対して付与された点数
	Evidence       map[string]string `json:"evidence,omitempty"`       // コマンド出力などの証跡
	Recommendation string            `json:"recommendation,omitempty"` // 改善提案（v0.1は任意）
}

// ホスト情報（OSなど）
type HostInfo struct {
	Hostname string `json:"hostname"`
	OS       OSInfo `json:"os"`
}

type OSInfo struct {
	Product string `json:"product"`
	Version string `json:"version"`
	Build   string `json:"build"`
}

// 監査レポートの全体構造
type Report struct {
	Version string        `json:"version"` // macinsight のバージョン
	Host    HostInfo      `json:"host"`
	Score   int           `json:"score"`   // 0〜100
	Checks  []CheckResult `json:"checks"`
}
