package checks

import (
	"context"
	"time"

	"github.com/samuraidays/macinsight/internal/executil"
)

// runCommand is an indirection over executil.Run to allow tests to mock command execution.
var runCommand = func(ctx context.Context, timeout time.Duration, name string, args ...string) executil.Result {
	return executil.Run(ctx, timeout, name, args...)
}
