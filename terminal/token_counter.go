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

// CountTokens uses a generative AI model to count the number of tokens in the provided text input or image data.
// It returns the token count and any error encountered in the process. A new client is created and closed within the function.
func CountTokens(params TokenCountParams) (int, error) {
	ctx := context.Background()
	return countTokensWithClient(ctx, params)
}

// countTokensWithClient handles the token counting process, including client initialization, request execution,
// and error handling with retry logic. It abstracts away the complexities of interacting with the AI service.
func countTokensWithClient(ctx context.Context, params TokenCountParams) (int, error) {
	var tokenCount int

	success, err := performTokenCount(ctx, params, &tokenCount)
	if err != nil {
		return 0, err
	}
	if !success {
		return 0, fmt.Errorf(ErrorLowLevelFailedToCountTokensAfterRetries)
	}

	return tokenCount, nil
}

// performTokenCount manages retry logic for token counting, using a retry function and an error handler
// to determine if errors are transient. It delegates the actual token counting to makeTokenCountRequest.
func performTokenCount(ctx context.Context, params TokenCountParams, tokenCount *int) (bool, error) {
	retryFunc := func() (bool, error) {
		return makeTokenCountRequest(ctx, params, tokenCount)
	}

	return retryWithExponentialBackoff(retryFunc, standardAPIErrorHandler)
}

// makeTokenCountRequest attempts to count tokens by initializing a new AI client and sending the token counting request.
// It updates the token count based on the response and indicates if the operation was successful.
func makeTokenCountRequest(ctx context.Context, params TokenCountParams, tokenCount *int) (bool, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(params.APIKey))
	if err != nil {
		return false, err
	}
	defer client.Close()

	model := client.GenerativeModel(params.ModelName)

	resp, err := prepareAndCountTokens(ctx, model, params)
	if err != nil {
		return false, err
	}

	*tokenCount = int(resp.TotalTokens)
	return true, nil
}

// prepareAndCountTokens creates and executes a token counting request using the provided model,
// based on the input text or image data.
func prepareAndCountTokens(ctx context.Context, model *genai.GenerativeModel, params TokenCountParams) (*genai.CountTokensResponse, error) {
	if len(params.ImageData) > 0 && len(params.Input) > 0 {
		return model.CountTokens(ctx,
			genai.Text(params.Input),
			genai.ImageData(
				params.ImageFormat,
				params.ImageData))
	} else if len(params.ImageData) > 0 {
		return model.CountTokens(ctx,
			genai.ImageData(params.ImageFormat,
				params.ImageData))
	} else {
		return model.CountTokens(ctx,
			genai.Text(params.Input))
	}
}
