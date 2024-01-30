// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"context"

	genai "github.com/google/generative-ai-go/genai"
)

// GetEmbedding computes the numerical embedding for a given piece of text using
// the specified generative AI model. Embeddings are useful for a variety of machine
// learning tasks, such as semantic search, where they can represent the meaning of
// text in a form that can be processed by algorithms.
//
// Parameters:
//
//	ctx     context.Context: The context for controlling the lifetime of the request. It allows
//	                         the function to be canceled or to time out, and it carries request-scoped values.
//	client  *genai.Client:   The client used to interact with the generative AI service. It should be
//	                         already initialized and authenticated before calling this function.
//	modelID string:          The identifier for the embedding model to be used. This specifies which
//	                         AI model will generate the embeddings.
//	text    string:          The input text to be converted into an embedding.
//
// Returns:
//
//	[]float32: An array of floating-point numbers representing the embedding of the input text.
//	error:     An error that may occur during the embedding process. If the operation is successful,
//	           the error is nil.
//
// The function delegates the embedding task to the genai client's EmbeddingModel method and
// retrieves the embedding values from the response. It is the caller's responsibility to manage
// the lifecycle of the genai.Client, including its creation and closure.
//
// Note: This function marked as TODO for now, since it is not used in the main because,
// a current version of chat system it's consider fully stable with better logic.
func GetEmbedding(ctx context.Context, client *genai.Client, modelID, text string) ([]float32, error) {
	em := client.EmbeddingModel(modelID)
	res, err := em.EmbedContent(ctx, genai.Text(text))
	if err != nil {
		return nil, err
	}
	return res.Embedding.Values, nil
}
