// Package utils contains various utilities.
package utils

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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

// GetMD5Hash returns md5-sum of string.
func GetMD5Hash(text string) string {
	hash := md5.Sum([]byte(text))
	return hex.EncodeToString(hash[:])
}
