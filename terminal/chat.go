// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// Note: This is subject to change (for example, it can be customized in commands). For now, it's stable. Additionally, a token is inexpensive since, with Google AI's Gemini-Pro model, the maximum is 32K tokens.
const MaxChatHistory = 5 // Maximum number of messages to keep in history

// ChatHistory holds the chat messages exchanged during a session.
// It provides methods to add new messages to the history and to retrieve
// the current state of the conversation.
type ChatHistory struct {
	Messages []string
}

// AddMessage appends a new message to the chat history.
// It takes the username and the text of the message as inputs and formats
// them before adding to the Messages slice.
//
// Parameters:
//
//	user string: The username of the individual sending the message.
//	text string: The content of the message to be added to the history.
//
// This method does not return any value or error. It assumes that all input
// is valid and safe to add to the chat history.
func (h *ChatHistory) AddMessage(user, text string) {
	message := fmt.Sprintf("%s %s", user, text)
	h.Messages = append(h.Messages, message)
	if len(h.Messages) > MaxChatHistory {
		// Remove the oldest message to maintain a fixed history size
		h.Messages = h.Messages[1:]
	}
}

// GetHistory concatenates all messages in the chat history into a single
// string, with each message separated by a newline character. This provides
// a simple way to view the entire chat history as a single text block.
//
// Returns:
//
//	string: A newline-separated string of all messages in the chat history.
func (h *ChatHistory) GetHistory() string {
	// Note: if you're still wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
	// Define the prefixes to be removed
	prefixesToRemove := []string{YouNerd, AiNerd}

	for _, msg := range h.Messages {
		sanitizedMsg := msg
		// Remove each prefix from the start of the message
		for _, prefix := range prefixesToRemove {
			if strings.HasPrefix(sanitizedMsg, prefix) {
				sanitizedMsg = strings.TrimPrefix(sanitizedMsg, prefix)
				break // Assume only one prefix will match and then break the loop
			}
		}
		// Optimized to use Builder.WriteString() for better performance and to avoid memory allocation overhead.
		builder.WriteString(sanitizedMsg)
		builder.WriteRune(NewLineChars) // Append a newline character after each message.
	}

	// The builder.String() method returns the complete, concatenated chat history.
	return builder.String()

}
