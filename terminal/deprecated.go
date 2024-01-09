// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"log"
)

// PrintHistory outputs all messages in the chat history to the standard output,
// one message per line. This method is useful for displaying the chat history
// directly to the terminal.
//
// Each message is printed in the order it was added, preserving the conversation
// flow. This method does not return any value or error.
//
// Deprecated: This method is deprecated was replaced by GetHistory.
// It used to be used for debugging purposes while made the chat system without storage such as database.
func (h *ChatHistory) PrintHistory() {
	for _, msg := range h.Messages {
		fmt.Println(msg)
	}
}

// RecoverFromPanic returns a deferred function that recovers from panics within a goroutine
// or function, preventing the panic from propagating and potentially causing the program to crash.
// Instead, it logs the panic information using the standard logger, allowing for post-mortem analysis
// without interrupting the program's execution flow.
//
// Usage:
//
//	defer terminal.RecoverFromPanic()()
//
// The function returned by RecoverFromPanic should be called by deferring it at the start of
// a goroutine or function. When a panic occurs, the deferred function will handle the panic
// by logging its message and stack trace, as provided by the recover built-in function.
//
// Deprecated: This method is deprecated was replaced by logger.RecoverFromPanic.
func RecoverFromPanic() func() {
	return func() {
		if r := recover(); r != nil {
			// Log the panic with additional context if desired
			log.Printf("Recovered from panic: %+v\n", r)
		}
	}
}
