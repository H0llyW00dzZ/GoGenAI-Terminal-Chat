// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import genai "github.com/google/generative-ai-go/genai"

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
//  *ChatConfig: A pointer to a ChatConfig instance populated with default settings.
//
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

// Option is a function type that takes a pointer to a GenerativeModel instance and applies a specific configuration to it.
type Option func(m *genai.GenerativeModel)

// WithTemperature sets the temperature parameter of the model.
func WithTemperature(temperature float32) Option {
	return func(m *genai.GenerativeModel) {
		m.SetTemperature(float32(temperature))
	}
}

// WithTopP sets the top_p parameter of the model.
func WithTopP(topP float32) Option {
	return func(m *genai.GenerativeModel) {
		m.SetTopP(float32(topP))
	}
}

// WithTopK sets the top_k parameter of the model.
func WithTopK(topK int32) Option {
	return func(m *genai.GenerativeModel) {
		m.SetTopK(int32(topK))
	}
}

// WithMaxOutputTokens sets the max_output_tokens parameter of the model.
func WithMaxOutputTokens(maxOutputTokens int32) Option {
	return func(m *genai.GenerativeModel) {
		m.SetMaxOutputTokens(int32(maxOutputTokens))
	}
}

// ApplyOptions applies the provided options to the given generative AI model.
func ApplyOptions(m *genai.GenerativeModel, options ...Option) {
	for _, option := range options {
		option(m)
	}
}
