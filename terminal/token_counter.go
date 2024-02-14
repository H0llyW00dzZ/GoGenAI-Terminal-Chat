// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License
//
// Note: This CountTokens function supports multi-modal inputs and is designed for stability. It can utilize any model available in Google's AI offerings to count tokens.
// It also supports processing multiple images (e.g., 999999+ images hahaha) simultaneously.

package terminal

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"

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

	// Create an instance of RetryableOperation with the defined retryFunc.
	operation := RetryableOperation{
		retryFunc: retryFunc,
	}

	// Call the retryWithExponentialBackoff method on the RetryableOperation instance.
	return operation.retryWithExponentialBackoff(standardAPIErrorHandler)
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

// prepareAndCountTokens determines the method for token counting based on the provided input.
// If text input is provided, it directly counts tokens for the text. If image data is provided,
// it delegates to countTokensConcurrently for concurrent token counting on image data.
func (p *TokenCountParams) prepareAndCountTokens(ctx context.Context, model *genai.GenerativeModel) (*genai.CountTokensResponse, error) {
	// If there is text input, count tokens for text.
	if len(p.Input) > 0 {
		// Text input is present; prepare a request for concurrent token counting for text.
		request := TokenCountRequest{
			Ctx:   ctx,
			Model: model,
			Texts: []string{p.Input},
		}
		return p.countTokensConcurrently(request)
	} else if len(p.ImageData) > 0 {
		// Image data is present; prepare a request for concurrent token counting for images.
		request := TokenCountRequest{
			Ctx:    ctx,
			Model:  model,
			Images: p.ImageData,
		}
		return p.countTokensConcurrently(request)
	}
	return nil, fmt.Errorf(ErrorNoInputProvideForTokenCounting)
}

// countTokensConcurrently orchestrates concurrent token counting for multiple texts or images
// and aggregates the results into a single response.
func (p *TokenCountParams) countTokensConcurrently(req TokenCountRequest) (*genai.CountTokensResponse, error) {
	// Note: This a cheap in terms of efficiency, especially if the task is I/O-bound.
	var totalTokens int64
	var err error

	// Handle concurrent token counting for text inputs if any.
	if len(req.Texts) > 0 {
		totalTokens, err = p.launchTokenCountGoroutinesForText(req)
		if err != nil {
			// An error occurred during concurrent token counting; return the error.
			return nil, err
		}
	}

	// Handle concurrent token counting for image data if any.
	if len(req.Images) > 0 {
		imageTokens, err := p.launchTokenCountGoroutinesForImage(req)
		if err != nil {
			// An error occurred during concurrent token counting; return the error.
			return nil, err
		}
		totalTokens += imageTokens
	}

	// Return the total token count after successfully counting tokens for all inputs.
	return &genai.CountTokensResponse{TotalTokens: int32(totalTokens)}, nil
}

// launchTokenCountGoroutinesForImage starts a goroutine for each image in the request to count tokens in parallel.
// It waits for all goroutines to complete and returns the accumulated token count and the first error encountered, if any.
func (p *TokenCountParams) launchTokenCountGoroutinesForImage(req TokenCountRequest) (int64, error) {
	var totalTokens int64
	var wg sync.WaitGroup
	errChan := make(chan error, len(req.Images))
	// Note: This functionality may only be compatible with Go version 1.22 and onwards.
	// Ref: Range over integers (https://go.dev/doc/go1.22)
	for i, imageData := range req.Images {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine.
		go func(data []byte, index int) {
			defer wg.Done() // Decrement the counter when the goroutine completes.
			tokens, err := p.countTokensForImage(req.Ctx, req.Model, data)
			if err != nil {
				errChan <- fmt.Errorf(ErrorGopherEncounteredAnError, index, err) // Just incase adding this logger
				return
			}
			// Safely add the tokens from this image to the total count.
			atomic.AddInt64(&totalTokens, tokens)
		}(imageData, i)
	}

	go func() {
		wg.Wait()      // Wait for all goroutines to finish.
		close(errChan) // Close the error channel after all goroutines have finished.
	}()

	return totalTokens, collectErrors(errChan) // Return the total tokens and any error that occurred.
}

// countTokensForImage counts the tokens for a single image using the provided generative AI model.
// It returns the token count for the image and any error encountered during the process.
func (p *TokenCountParams) countTokensForImage(ctx context.Context, model *genai.GenerativeModel, imageData []byte) (int64, error) {
	resp, err := model.CountTokens(ctx, genai.ImageData(p.ImageFormat, imageData))
	if err != nil {
		// An error occurred while counting tokens for this image; return the error.
		return 0, err
	}
	// Token counting for this image was successful; return the count.
	return int64(resp.TotalTokens), nil
}

// launchTokenCountGoroutinesForText starts a goroutine for each text input in the request to count tokens in parallel.
// It waits for all goroutines to complete and returns the accumulated token count and the first error encountered, if any.
func (p *TokenCountParams) launchTokenCountGoroutinesForText(req TokenCountRequest) (int64, error) {
	var totalTokens int64
	var wg sync.WaitGroup
	errChan := make(chan error, len(req.Texts))
	// Note: This functionality may only be compatible with Go version 1.22 and onwards.
	// Ref: Range over integers (https://go.dev/doc/go1.22)
	for i, text := range req.Texts {
		wg.Add(1) // Increment the WaitGroup counter for each goroutine.
		// Note: This a better way, for example how it work it's inputValueString1 handle by goroutine 1, inputValueString2 handle by goroutine 2
		go func(t string, index int) {
			defer wg.Done() // Decrement the counter when the goroutine completes.
			tokens, err := p.countTokensForText(req.Ctx, req.Model, t)
			if err != nil {
				errChan <- fmt.Errorf(ErrorGopherEncounteredAnError, index, err) // Just incase adding this logger
				return
			}
			// Safely add the tokens from this text to the total count.
			atomic.AddInt64(&totalTokens, tokens)
		}(text, i)
	}

	go func() {
		wg.Wait()      // Wait for all goroutines to finish.
		close(errChan) // Close the error channel after all goroutines have finished.
	}()

	return totalTokens, collectErrors(errChan) // Return the total tokens and any error that occurred.
}

// countTokensForText counts the tokens for a single text input using the provided generative AI model.
// It returns the token count for the text and any error encountered during the process.
func (p *TokenCountParams) countTokensForText(ctx context.Context, model *genai.GenerativeModel, text string) (int64, error) {
	resp, err := model.CountTokens(ctx, genai.Text(text))
	if err != nil {
		// An error occurred while counting tokens for this text; return the error.
		return 0, err
	}
	// Token counting for this text was successful; return the count.
	return int64(resp.TotalTokens), nil
}

// collectErrors waits for all errors to be sent on the given channel and returns the first non-nil error.
func collectErrors(errChan <-chan error) error {
	var collectedError error
	for err := range errChan {
		if collectedError == nil && err != nil {
			collectedError = err // Record the first non-nil error encountered.
		}
	}
	return collectedError
}
