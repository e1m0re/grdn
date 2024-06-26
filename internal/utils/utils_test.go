package utils

import (
	"context"
	"fmt"
)

func ExampleRetryFunc() {
	err := RetryFunc(context.Background(), func() error {
		// ... do something
		return nil
	})
	if err != nil {
		fmt.Printf("error")
	}
}
