package checks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

func TestFileVault_Pass(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Stdout: "FileVault is On.\n"}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := FileVault(context.Background())
	if cr.Status != "pass" || cr.Score == 0 {
		t.Fatalf("FileVault pass expected, got status=%s score=%d", cr.Status, cr.Score)
	}
}

func TestFileVault_UnknownOnError(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Err: errors.New("exec error")}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := FileVault(context.Background())
	if cr.Status != "unknown" {
		t.Fatalf("FileVault unknown expected on error, got %s", cr.Status)
	}
}
