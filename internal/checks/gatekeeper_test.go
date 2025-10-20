package checks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

func TestGatekeeper_Pass(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Stdout: "assessments enabled\n"}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := Gatekeeper(context.Background())
	if cr.Status != "pass" || cr.Score == 0 {
		t.Fatalf("Gatekeeper pass expected, got status=%s score=%d", cr.Status, cr.Score)
	}
}

func TestGatekeeper_UnknownOnError(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Err: errors.New("exec error")}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := Gatekeeper(context.Background())
	if cr.Status != "unknown" {
		t.Fatalf("Gatekeeper unknown expected on error, got %s", cr.Status)
	}
}
