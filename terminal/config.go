// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"

	genai "github.com/google/generative-ai-go/genai"
)

// DefaultChatConfig constructs a new ChatConfig with pre-defined default values.
// These defaults are chosen to balance the need for context awareness by the AI
// and efficient memory usage. The function is useful for initializing chat sessions
// with standard behavior without requiring manual configuration.
//
// Returns:
//
//	*ChatConfig: A pointer to a ChatConfig instance populated with default settings.
func DefaultChatConfig() *ChatConfig {
	return &ChatConfig{
		// Note: This history size is stable. It is automatically handled by the garbage collector.
		// Ref:
		// - https://tip.golang.org/doc/gc-guide
		// - https://pkg.go.dev/builtin
		HistorySize: 10, // Default to retaining the last 10 messages
		// Note: HistorySendToAI currently is unimplemented, will implemented it later when I am free
		HistorySendToAI: 10, // Default to sending the last 10 messages for AI context
	}
}

// ConfigureModel applies a series of configuration options to a GenerativeModel.
// This function is variadic, meaning it can accept multiple configuration options
// that will be applied in the order they are provided.
//
// Parameters:
//
//	model *genai.GenerativeModel: The generative AI model to configure.
//	opts ...ModelConfig: A variadic number of configuration options.
func ConfigureModel(model *genai.GenerativeModel, opts ...ModelConfig) {
	for _, opt := range opts {
		opt(model)
	}
}

// WithTemperature creates a ModelConfig function to set the temperature
// of a GenerativeModel. Temperature controls the randomness of the AI's
// responses, with higher values leading to more varied output.
//
// Parameters:
//
//	temperature float32: The temperature value to set.
//
// Returns:
//
//	ModelConfig: A function that sets the temperature when applied to a model.
func WithTemperature(temperature float32) ModelConfig {
	return func(m *genai.GenerativeModel) {
		m.SetTemperature(temperature)
	}
}

// WithTopP creates a ModelConfig function to set the top_p parameter
// of a GenerativeModel. Top_p controls the nucleus sampling strategy, where
// a smaller value leads to less randomness in token selection.
//
// Parameters:
//
//	topP float32: The top_p value to set.
//
// Returns:
//
//	ModelConfig: A function that sets the top_p value when applied to a model.
func WithTopP(topP float32) ModelConfig {
	return func(m *genai.GenerativeModel) {
		m.SetTopP(topP)
	}
}

// WithTopK creates a ModelConfig function to set the top_k parameter
// of a GenerativeModel. Top_k restricts the sampling pool to the k most likely
// tokens, where a lower value increases the likelihood of high-probability tokens.
//
// Parameters:
//
//	topK int32: The top_k value to set.
//
// Returns:
//
//	ModelConfig: A function that sets the top_k value when applied to a model.
func WithTopK(topK int32) ModelConfig {
	return func(m *genai.GenerativeModel) {
		m.SetTopK(topK)
	}
}

// WithMaxOutputTokens creates a ModelConfig function to set the maximum number
// of output tokens for a GenerativeModel. This parameter limits the length of
// the AI's responses.
//
// Parameters:
//
//	maxOutputTokens int32: The maximum number of tokens to set.
//
// Returns:
//
//	ModelConfig: A function that sets the maximum number of output tokens when applied to a model.
//	error: An error if maxOutputTokens is below 50.
func WithMaxOutputTokens(maxOutputTokens int32) (ModelConfig, error) {
	if maxOutputTokens < 50 {
		return nil, fmt.Errorf(ErrorMaxOutputTokenMustbe, maxOutputTokens)
	}
	return func(m *genai.GenerativeModel) {
		m.SetMaxOutputTokens(maxOutputTokens)
	}, nil
}

// WithSafetyOptions creates a ModelConfig function to set the safety options
// for a GenerativeModel.
//
// Parameters:
//
//	safety *SafetySettings: The safety settings to apply.
//
// Returns:
//
//	ModelConfig: A function that, when applied to a model, sets the safety options.
//
// Note: This is currently marked as TODO, as it's not used anywhere in the code. However, it would be advantageous to implement this feature in the future.
// For instance, it could be used with Vertex AI.
func WithSafetyOptions(safety *SafetySettings, modelName string) ModelConfig {
	// Note: This advanced idiomatic Go makes use of pointers hahaha.
	return func(m *genai.GenerativeModel) {
		safety.ApplyToModel(m, modelName)
	}
}

// ApplyOptions is a convenience function that applies a series of configuration
// options to a GenerativeModel. This allows for flexible and dynamic model
// configuration at runtime.
//
// Parameters:
//
//	m *genai.GenerativeModel: The generative AI model to configure.
//	configs ...ModelConfig: A variadic number of configuration options.
//
// Returns:
//
//	bool: A boolean indicating whether the options were applied successfully.
//	error: An error if any of the configuration options are nil.
func ApplyOptions(m *genai.GenerativeModel, configs ...ModelConfig) (bool, error) {
	for _, option := range configs {
		if option == nil {
			return false, fmt.Errorf(ErrorGenAiReceiveNil)
		}
		option(m)
	}
	return true, nil
}
