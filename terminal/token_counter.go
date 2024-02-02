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
// and counts the number of tokens in the given input string or image data.
// This function is useful for understanding the token usage of inputs in the context
// of generative AI, which can help manage API usage and costs.
//
// Parameters:
//
//	apiKey     string         : The API key used to authenticate with the generative AI service.
//	input      string         : The text input for which the number of tokens will be counted.
//	imageFormat string  : The format of the image (e.g., "png", "jpeg"), if image data is provided.
//	imageData  []byte   : The byte slice containing the image data.
//
// Returns:
//
//	int   : The number of tokens that the input contains.
//	error : An error that occurred while creating the client, connecting to the service,
//	        or counting the tokens. If the operation is successful, the error is nil.
func CountTokens(apiKey, input, imageFormat string, imageData []byte) (int, error) {
	ctx := context.Background()
	return countTokensWithClient(ctx, apiKey, input, imageFormat, imageData)
}

// countTokensWithClient orchestrates the process of counting the number of tokens
// in a given input string and/or image data using a generative AI model. This function
// is designed to handle the complexities of interacting with the AI service, including
// client initialization, request execution, and error handling with retry logic.
//
// Parameters:
//
//	ctx         : The context for controlling the lifetime of the request.
//	apiKey      : The API key used to authenticate with the generative AI service.
//	input       : The text input for which the number of tokens will be counted.
//	imageFormat : The format of the image (e.g., "png", "jpeg"), if image data is provided.
//	imageData   : The byte slice containing the image data.
//
// Returns:
//
//	tokenCount : The number of tokens in the input string and/or image data.
//	err        : An error encountered during the token counting process.
//
// Note: This function leverages performTokenCount to manage retries and error handling,
// abstracting the retry logic away from the core token counting operation.
func countTokensWithClient(ctx context.Context, apiKey, input, imageFormat string, imageData []byte) (int, error) {
	var tokenCount int

	success, err := performTokenCount(ctx, apiKey, input, imageFormat, imageData, &tokenCount)
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
//	imageFormat : The format of the image (e.g., "png", "jpeg"), if image data is provided.
//	imageData   : The byte slice containing the image data.
//	tokenCount  : A pointer to an integer that will hold the token count result.
//
// Returns:
//
//	success : A boolean indicating whether the token counting operation succeeded.
//	err     : An error encountered during the token counting process.
//
// Note: This function delegates the actual token counting to makeTokenCountRequest
// and is responsible for invoking the retry logic.
func performTokenCount(ctx context.Context, apiKey, input, imageFormat string, imageData []byte, tokenCount *int) (bool, error) {
	retryFunc := func() (bool, error) {
		return makeTokenCountRequest(ctx, apiKey, input, imageFormat, imageData, tokenCount)
	}

	return retryWithExponentialBackoff(retryFunc, standardAPIErrorHandler)
}

// makeTokenCountRequest performs a single attempt to count tokens using the AI model.
// It initializes a new AI client, sends the token counting request with either text,
// image data, or both, and updates the token count based on the response.
//
// Parameters:
//
//	ctx         : The context for controlling the lifetime of the request.
//	apiKey      : The API key used to authenticate with the generative AI service.
//	input       : The text input for which the number of tokens will be counted.
//	imageFormat : The format of the image (e.g., "png", "jpeg"), if image data is provided.
//	imageData   : The byte slice containing the image data.
//	tokenCount  : A pointer to an integer that will hold the token count result.
//
// Returns:
//
//	success : A boolean indicating whether the token counting operation succeeded.
//	err     : An error encountered during the token counting process.
//
// Note: This function is called within the retry logic of performTokenCount and
// handles the direct interaction with the AI service for counting tokens.
func makeTokenCountRequest(ctx context.Context, apiKey, input, imageFormat string, imageData []byte, tokenCount *int) (bool, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return false, err
	}
	defer client.Close()

	model := client.GenerativeModel(ModelAi)

	resp, err := prepareAndCountTokens(ctx, model, input, imageFormat, imageData)
	if err != nil {
		return false, err
	}

	*tokenCount = int(resp.TotalTokens)
	return true, nil
}

// prepareAndCountTokens prepares the token counting request based on the input and image data
// and executes the request using the provided model.
func prepareAndCountTokens(ctx context.Context, model *genai.GenerativeModel, input, imageFormat string, imageData []byte) (*genai.CountTokensResponse, error) {
	if len(imageData) > 0 && len(input) > 0 {
		return model.CountTokens(ctx, genai.Text(input), genai.ImageData(imageFormat, imageData))
	} else if len(imageData) > 0 {
		return model.CountTokens(ctx, genai.ImageData(imageFormat, imageData))
	} else {
		return model.CountTokens(ctx, genai.Text(input))
	}
}
