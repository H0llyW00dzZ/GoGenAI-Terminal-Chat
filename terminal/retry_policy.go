// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"math"
	"time"
)

// RetryableFunc is a type that represents a function that can be retried.
type RetryableFunc func() (bool, error)

// retryWithExponentialBackoff attempts to execute a RetryableFunc with a retry policy.
// It applies exponential backoff between retries and logs an error if the maximum number of retries is reached.
//
// Note: this a powerful retry policy, unlike that shitty complex go codes
func retryWithExponentialBackoff(retryFunc RetryableFunc) (bool, error) {
	const maxRetries = 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		success, err := retryFunc()
		if err == nil {
			return success, nil
		}
		// log debug
		logger.Debug(fmt.Sprintf(DEBUGRETRYPOLICY, attempt+1, err))
		// Log the error
		if logger.HandleGoogleAPIError(err) {
			delay := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
			time.Sleep(delay)
			continue // Retry the request
		} else {
			// Non-retryable error or max retries exceeded
			logger.Error(ErrorFailedToSendMessagesAfterRetryingonInternalServerError, err)
			return false, err
		}
	}

	// If this point is reached, retries have been exhausted without success.
	err := fmt.Errorf(ErrorLowLevelMaximumRetries)
	logger.Error(err.Error())
	return false, err
}
