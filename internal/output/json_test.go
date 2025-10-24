package output

import (
	"bytes"
	"encoding/json"
	"testing"

	"github.com/samuraidays/macinsight/pkg/types"
)

func TestWriteJSON_EncodesReport(t *testing.T) {
	rep := types.Report{
		Version: "vtest",
		Host:    types.HostInfo{Hostname: "host", OS: types.OSInfo{Product: "macOS", Version: "26.0.1", Build: "23A344"}},
		Score:   42,
		Checks:  []types.CheckResult{{ID: "sip", Title: "System Integrity Protection enabled", Status: "pass", Score: 20}},
	}

	var buf bytes.Buffer
	if err := WriteJSON(&buf, rep); err != nil {
		t.Fatalf("WriteJSON error: %v", err)
	}

	var got types.Report
	if err := json.Unmarshal(buf.Bytes(), &got); err != nil {
		t.Fatalf("json.Unmarshal error: %v", err)
	}

	if got.Version != rep.Version || got.Score != rep.Score {
		t.Fatalf("decoded report mismatch: got=%+v", got)
	}

	if len(got.Checks) != 1 || got.Checks[0].ID != "sip" {
		t.Fatalf("decoded checks mismatch: %+v", got.Checks)
	}
}
