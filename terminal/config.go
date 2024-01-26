// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	genai "github.com/google/generative-ai-go/genai"
)

// ChatConfig encapsulates settings that affect the management of chat history
// during a session with the generative AI. It determines the amount of chat history
// retained in memory and the portion of that history used to provide context to the AI.
type ChatConfig struct {
	// HistorySize specifies the total number of chat messages to retain in the session's history.
	// This helps in limiting the memory footprint and ensures that only recent interactions
	// are considered for maintaining context.
	HistorySize int

	// HistorySendToAI indicates the number of recent messages from the history to be included
	// when sending context to the AI. This allows the AI to generate responses that are
	// relevant to the current conversation flow without being overwhelmed by too much history.
	HistorySendToAI int
}

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

// ModelConfig defines a function type for configuring a GenerativeModel.
// Functions of this type take a pointer to a GenerativeModel and apply
// specific settings to it.
type ModelConfig func(m *genai.GenerativeModel)

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
func WithMaxOutputTokens(maxOutputTokens int32) ModelConfig {
	return func(m *genai.GenerativeModel) {
		m.SetMaxOutputTokens(maxOutputTokens)
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
func ApplyOptions(m *genai.GenerativeModel, configs ...ModelConfig) {
	for _, option := range configs {
		option(m)
	}
}
