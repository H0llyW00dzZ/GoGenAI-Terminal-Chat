// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"context"
	"fmt"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// CountTokens connects to a generative AI model using the provided API key
// and counts the number of tokens in the given input string. This function
// is useful for understanding the token usage of text inputs in the context
// of generative AI, which can help manage API usage and costs.
//
// Parameters:
//
//	apiKey string: The API key used to authenticate with the generative AI service.
//	input  string: The text input for which the number of tokens will be counted.
//
// Returns:
//
//	int:   The number of tokens that the input string contains.
//	error: An error that occurred while creating the client, connecting to the service,
//	       or counting the tokens. If the operation is successful, the error is nil.
//
// The function creates a new client for each call, which is then closed before
// returning. It is designed to be a self-contained operation that does not require
// the caller to manage the lifecycle of the generative AI client.
func CountTokens(apiKey, input string) (int, error) {
	ctx := context.Background()
	return countTokensWithClient(ctx, apiKey, input)
}

// countTokensWithClient orchestrates the process of counting the number of tokens
// in a given input string using a generative AI model. This function is designed to
// handle the complexities of interacting with the AI service, including client
// initialization, request execution, and error handling with retry logic.
//
// Parameters:
//
//	ctx     : The context for controlling the lifetime of the request.
//	apiKey  : The API key used to authenticate with the generative AI service.
//	input   : The text input for which the number of tokens will be counted.
//
// Returns:
//
//	tokenCount : The number of tokens in the input string.
//	err        : An error encountered during the token counting process.
//
// Note: This function leverages performTokenCount to manage retries and error handling,
// abstracting the retry logic away from the core token counting operation.
func countTokensWithClient(ctx context.Context, apiKey, input string) (int, error) {
	var tokenCount int

	success, err := performTokenCount(ctx, apiKey, input, &tokenCount)
	if err != nil {
		return 0, err
	}
	if !success {
		return 0, fmt.Errorf(ErrorLowLevelFailedToCountTokensAfterRetries)
	}

	return tokenCount, nil
}

// performTokenCount manages the retry logic for the token counting operation.
// It uses a retry function that attempts to count tokens and an error handler
// that determines whether errors are transient and warrant a retry.
//
// Parameters:
//
//	ctx         : The context for controlling the lifetime of the request.
//	apiKey      : The API key used to authenticate with the generative AI service.
//	input       : The text input for which the number of tokens will be counted.
//	tokenCount  : A pointer to an integer that will hold the token count result.
//
// Returns:
//
//	success : A boolean indicating whether the token counting operation succeeded.
//	err     : An error encountered during the token counting process.
//
// Note: This function delegates the actual token counting to makeTokenCountRequest
// and is responsible for invoking the retry logic.
func performTokenCount(ctx context.Context, apiKey, input string, tokenCount *int) (bool, error) {
	retryFunc := func() (bool, error) {
		return makeTokenCountRequest(ctx, apiKey, input, tokenCount)
	}

	return retryWithExponentialBackoff(retryFunc, standardAPIErrorHandler)
}

// makeTokenCountRequest performs a single attempt to count tokens using the AI model.
// It initializes a new AI client, sends the token counting request, and updates the
// token count based on the response.
//
// Parameters:
//
//	ctx         : The context for controlling the lifetime of the request.
//	apiKey      : The API key used to authenticate with the generative AI service.
//	input       : The text input for which the number of tokens will be counted.
//	tokenCount  : A pointer to an integer that will hold the token count result.
//
// Returns:
//
//	success : A boolean indicating whether the token counting operation succeeded.
//	err     : An error encountered during the token counting process.
//
// Note: This function is called within the retry logic of performTokenCount and
// handles the direct interaction with the AI service for counting tokens.
func makeTokenCountRequest(ctx context.Context, apiKey, input string, tokenCount *int) (bool, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return false, err
	}
	defer client.Close()

	model := client.GenerativeModel(ModelAi)
	resp, err := model.CountTokens(ctx, genai.Text(input))
	if err != nil {
		return false, err
	}

	*tokenCount = int(resp.TotalTokens)
	return true, nil
}
