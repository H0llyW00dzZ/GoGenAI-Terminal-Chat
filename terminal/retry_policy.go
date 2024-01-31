// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"math"
	"strings"
	"time"
)

// RetryableFunc is a type that represents a function that can be retried.
type RetryableFunc func() (bool, error)

// ErrorHandlerFunc is a type that represents a function that handles an error and
// decides whether the operation should be retried.
type ErrorHandlerFunc func(error) bool

// retryWithExponentialBackoff attempts to execute a RetryableFunc with a retry policy.
// It applies exponential backoff between retries and logs an error if the maximum number of retries is reached.
//
// Note: this a powerful retry policy, unlike that shitty complex go codes
func retryWithExponentialBackoff(retryFunc RetryableFunc, handleError ErrorHandlerFunc) (bool, error) {
	const maxRetries = 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		success, err := retryFunc()
		if err == nil {
			return success, nil
		}
		// Log debug information
		logger.Debug(DEBUGRETRYPOLICY, attempt+1, err)

		// Use the provided error handler to check if we should retry.
		if handleError(err) {
			delay := baseDelay * time.Duration(math.Pow(2, float64(attempt)))
			time.Sleep(delay)
			continue // Retry the request
		} else {
			// Non-retryable error or max retries exceeded
			logger.Error(ErrorNonretryableerror, maxRetries, err)
			return false, err
		}
	}

	// If this point is reached, retries have been exhausted without success.
	err := fmt.Errorf(ErrorLowLevelMaximumRetries)
	logger.Error(err.Error())
	return false, err
}

// standardAPIErrorHandler is the standard error handling strategy for API errors.
func standardAPIErrorHandler(err error) bool {
	// Error 500 Google Api
	return strings.Contains(err.Error(), Error500GoogleApi)
}
