package api

import (
	"fmt"
	"math/rand"
	"time"
)

type RetryableFunc func() error

func RetryWithBackoff(attempts int, baseDelay time.Duration, maxDelay time.Duration, retryableFunc RetryableFunc) error {
	var err error
	for i := 0; i < attempts; i++ {
		err = retryableFunc()
		if err == nil {
			return nil // Success, no need to retry anymore
		}

		// Calculate the backoff delay using exponential growth, with jitter.
		// delay = baseDelay * (2^i) + random jitter up to baseDelay
		backoff := baseDelay * (1 << i)
		if backoff > maxDelay {
			backoff = maxDelay
		}
		// Add jitter to the delay to avoid thundering herd problem
		jitter := time.Duration(rand.Int63n(int64(baseDelay)))
		backoff = backoff + jitter

		fmt.Printf("Attempt %d failed. Retrying in %v...\n", i+1, backoff)
		time.Sleep(backoff)
	}
	return fmt.Errorf("operation failed after %d attempts: %v", attempts, err)
}
