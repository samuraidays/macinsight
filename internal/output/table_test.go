package output

import (
	"bytes"
	"strings"
	"testing"

	"github.com/samuraidays/macinsight/pkg/types"
)

func TestWriteTable_RendersRowsAndTotal(t *testing.T) {
	rep := types.Report{
		Version: "vtest",
		Host:    types.HostInfo{Hostname: "host"},
		Score:   30,
		Checks: []types.CheckResult{
			{ID: "a", Title: "AAA", Status: "pass", Score: 10, Evidence: map[string]string{"version": "1"}},
			{ID: "b", Title: "BBB", Status: "fail", Score: 0, Evidence: map[string]string{"x": "y"}},
			{ID: "c", Title: "CCC", Status: "warn", Score: 20, Evidence: map[string]string{"k": "v"}},
		},
	}

	var buf bytes.Buffer
	if err := WriteTable(&buf, rep); err != nil {
		t.Fatalf("WriteTable error: %v", err)
	}
	out := buf.String()
	outUpper := strings.ToUpper(out)

	// Contains header (case-insensitive because library uppercases headers)
	for _, h := range []string{"CHECK", "STATUS", "SCORE", "EVIDENCE"} {
		if !strings.Contains(outUpper, h) {
			t.Fatalf("missing header %q in output: %s", h, out)
		}
	}

	// Contains titles and total score
	for _, want := range []string{"AAA", "BBB", "CCC", "TOTAL", "30"} {
		if !strings.Contains(out, want) {
			t.Fatalf("missing %q in output: %s", want, out)
		}
	}
}
