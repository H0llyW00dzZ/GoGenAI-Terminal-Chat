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

// IsANSISequence checks if the current index in the rune slice is the start of an ANSI sequence.
//
// Deprecated: This method is no longer used, and was replaced by SanitizeMessage.
// It used to be filter for fix truncated message.
func IsANSISequence(runes []rune, index int) bool {
	return index+1 < len(runes) && runes[index] == ansichar.BinaryAnsiChar && runes[index+1] == ansichar.BinaryLeftSquareBracket
}

// PrintANSISequence prints the full ANSI sequence without delay and returns the new index.
//
// Deprecated: This method is no longer used, and was replaced by SanitizeMessage.
// It used to be filter for fix truncated message.
func PrintANSISequence(runes []rune, index int) int {
	// Print the full ANSI sequence without delay.
	for index < len(runes) && runes[index] != ansichar.BinaryAnsiSquenseChar {
		fmt.Printf(humantyping.AnimatedChars, runes[index])
		index++ // Move past the current character.
	}
	if index < len(runes) {
		fmt.Printf(humantyping.AnimatedChars, runes[index]) // Print the 'm' character to complete the ANSI sequence.
	}
	return index // Return the new index position.
}

// ApplyBold applies bold formatting to the provided text if the delimiter indicates bold.
//
// Deprecated: This method is no longer used, and was replaced by ApplyFormatting.
// It used to be used for formatting text.
func ApplyBold(text string, delimiter string, color string) string {
	if delimiter == DoubleAsterisk {
		return color + BoldText + text + ResetText + ColorReset
	}
	return color + text + ColorReset
}
