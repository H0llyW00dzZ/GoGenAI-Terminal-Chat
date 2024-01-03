package terminal

import (
	"context"

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
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		return 0, err
	}
	defer client.Close()

	model := client.GenerativeModel(ModelAi)

	resp, err := model.CountTokens(ctx, genai.Text(input))
	if err != nil {
		return 0, err
	}

	// Convert int32 to int
	return int(resp.TotalTokens), nil
}
