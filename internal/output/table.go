package output

import (
	"fmt"
	"io"
	"sort"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/samuraidays/macinsight/pkg/types"
)

// レポートを表形式で出力
func WriteTable(w io.Writer, r types.Report) error {
	t := table.NewWriter()
	t.SetOutputMirror(w)
	t.AppendHeader(table.Row{"Check", "Status", "Score", "Evidence"})

	// 見やすさのためタイトルでソート
	checks := append([]types.CheckResult(nil), r.Checks...)
	sort.Slice(checks, func(i, j int) bool { return checks[i].Title < checks[j].Title })

	for _, c := range checks {
		ev := ""
		// Evidenceを "k=v " で簡易整形
		keys := make([]string, 0, len(c.Evidence))
		for k := range c.Evidence {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			ev += fmt.Sprintf("%s=%s ", k, c.Evidence[k])
		}
		t.AppendRow(table.Row{c.Title, c.Status, c.Score, ev})
	}

	t.AppendFooter(table.Row{"TOTAL", "", r.Score, ""})
	t.Render()
	return nil
}
