// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"math"
	"time"
)

// RetrypolicyFunc is a type that represents a function that can be retried.
type RetryableFunc func() (bool, error)

// retryWithExponentialBackoff attempts to execute a RetryableFunc with a retry policy.
func retryWithExponentialBackoff(retryFunc RetryableFunc) (bool, error) {
	const maxRetries = 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		success, err := retryFunc()
		if err == nil {
			return success, nil
		}

		if logger.HandleGoogleAPIError(err) {
			delay := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
			time.Sleep(delay)
			continue // Retry the request
		} else {
			// Non-retryable error or max retries exceeded
			return false, err
		}
	}

	return false, fmt.Errorf(ErrorLowLevelMaximumRetries)
}
