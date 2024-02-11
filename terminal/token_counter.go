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
func (p *TokenCountParams) CountTokens() (int, error) {
	ctx := context.Background()
	return p.countTokensWithClient(ctx)
}

// countTokensWithClient handles the token counting process, including client initialization, request execution,
// and error handling with retry logic. It abstracts away the complexities of interacting with the AI service.
func (p *TokenCountParams) countTokensWithClient(ctx context.Context) (int, error) {
	var tokenCount int

	success, err := p.performTokenCount(ctx, &tokenCount)
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
func (p *TokenCountParams) performTokenCount(ctx context.Context, tokenCount *int) (bool, error) {
	retryFunc := func() (bool, error) {
		return p.makeTokenCountRequest(ctx, tokenCount)
	}

	// Assuming retryWithExponentialBackoff and standardAPIErrorHandler are defined elsewhere
	return retryWithExponentialBackoff(retryFunc, standardAPIErrorHandler)
}

// makeTokenCountRequest attempts to count tokens by initializing a new AI client and sending the token counting request.
// It updates the token count based on the response and indicates if the operation was successful.
func (p *TokenCountParams) makeTokenCountRequest(ctx context.Context, tokenCount *int) (bool, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey(p.APIKey))
	if err != nil {
		return false, err
	}
	defer client.Close()

	model := client.GenerativeModel(p.ModelName)

	resp, err := p.prepareAndCountTokens(ctx, model)
	if err != nil {
		return false, err
	}

	*tokenCount = int(resp.TotalTokens)
	return true, nil
}

// prepareAndCountTokens creates and executes a token counting request using the provided model,
// based on the input text or image data.
func (p *TokenCountParams) prepareAndCountTokens(ctx context.Context, model *genai.GenerativeModel) (*genai.CountTokensResponse, error) {
	if len(p.ImageData) > 0 && len(p.Input) > 0 {
		return model.CountTokens(ctx,
			genai.Text(p.Input),
			genai.ImageData(
				p.ImageFormat,
				p.ImageData))
	} else if len(p.ImageData) <= 0 {
		return model.CountTokens(ctx,
			genai.Text(p.Input))
	}

	return model.CountTokens(ctx,
		genai.ImageData(
			p.ImageFormat,
			p.ImageData))
}
