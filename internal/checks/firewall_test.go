package checks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

func TestFirewall_Pass(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Stdout: "State = 1\n"}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := Firewall(context.Background())
	if cr.Status != "pass" || cr.Score == 0 {
		t.Fatalf("Firewall pass expected, got status=%s score=%d", cr.Status, cr.Score)
	}
}

func TestFirewall_UnknownOnError(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Err: errors.New("exec error")}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := Firewall(context.Background())
	if cr.Status != "unknown" {
		t.Fatalf("Firewall unknown expected on error, got %s", cr.Status)
	}
}
