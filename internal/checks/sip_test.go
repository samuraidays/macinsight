package checks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

func TestSIP_Pass(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Stdout: "System Integrity Protection status: enabled\n"}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := SIP(context.Background())
	if cr.Status != "pass" || cr.Score == 0 {
		t.Fatalf("SIP pass expected, got status=%s score=%d", cr.Status, cr.Score)
	}
}

func TestSIP_UnknownOnError(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Err: errors.New("exec error")}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := SIP(context.Background())
	if cr.Status != "unknown" {
		t.Fatalf("SIP unknown expected on error, got %s", cr.Status)
	}
}
