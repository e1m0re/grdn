package server

import (
	"context"
	"time"

	"github.com/avast/retry-go"
)

func retryFunc(ctx context.Context, retryableFunc retry.RetryableFunc) error {
	return retry.Do(
		retryableFunc,
		retry.Attempts(3),
		retry.Delay(2*time.Second),
		retry.Context(ctx),
	)
}
