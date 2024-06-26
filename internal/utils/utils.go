// Package utils contains various utilities.
package utils

import (
	"context"
	"time"

	"github.com/avast/retry-go"
)

// RetryFunc re-calls the specified method if errors occur.
func RetryFunc(ctx context.Context, retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,
		retry.Attempts(3),
		retry.Delay(2*time.Second),
		retry.Context(ctx),
	)
}
