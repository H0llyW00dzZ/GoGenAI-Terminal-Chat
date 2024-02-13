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

// retryWithExponentialBackoff attempts to execute the RetryableFunc with a retry policy.
// It applies exponential backoff between retries and logs an error if the maximum number of retries is reached.
//
// Note: this a powerful retry policy, unlike that shitty complex go codes
func (op *RetryableOperation) retryWithExponentialBackoff(handleError ErrorHandlerFunc) (bool, error) {
	const maxRetries = 3
	baseDelay := time.Second

	for attempt := 0; attempt < maxRetries; attempt++ {
		success, err := op.retryFunc()
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
	return strings.Contains(err.Error(), Error500GoogleAPI)
}

// standardOtherAPIErrorHandler is the standard error handling strategy for API errors.
func standardOtherAPIErrorHandler(err error) bool {
	// Error 500 Other Api
	return strings.Contains(err.Error(), Code500)
}
