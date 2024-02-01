// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import genai "github.com/google/generative-ai-go/genai"

// ErrorHandlerFunc is a type that represents a function that handles an error and
// decides whether the operation should be retried.
type ErrorHandlerFunc func(error) bool

// ModelConfig defines a function type for configuring a GenerativeModel.
// Functions of this type take a pointer to a GenerativeModel and apply
// specific settings to it.
type ModelConfig func(m *genai.GenerativeModel)

// RetryableFunc is a type that represents a function that can be retried.
type RetryableFunc func() (bool, error)
