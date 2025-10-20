package checks

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

func TestAutoLogin_PassWhenKeyMissing(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Err: errors.New("The domain/default pair does not exist")}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := AutoLogin(context.Background())
	if cr.Status != "pass" {
		t.Fatalf("AutoLogin should pass when key missing, got %s", cr.Status)
	}
}

func TestAutoLogin_FailWhenUserSet(t *testing.T) {
	orig := runCommand
	runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
		return executil.Result{Stdout: "someuser\n"}
	}
	t.Cleanup(func() { runCommand = orig })

	cr := AutoLogin(context.Background())
	if cr.Status != "fail" {
		t.Fatalf("AutoLogin should fail when user is set, got %s", cr.Status)
	}
}
