// Copyright (c) 2024 H0llyW00dzZ

package terminal

import (
	"fmt"
	"strings"
)

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
	// Check if the last character of the user string is not a colon
	if !strings.HasSuffix(user, PrefixChar) {
		// If it is, don't add another colon
		h.Messages = append(h.Messages, fmt.Sprintf("%s %s", user, text))
	}
	// Append the message with the user string (which now has a guaranteed single colon)
	h.Messages = append(h.Messages, fmt.Sprintf("%s %s", user, text))
}

// GetHistory concatenates all messages in the chat history into a single
// string, with each message separated by a newline character. This provides
// a simple way to view the entire chat history as a single text block.
//
// Returns:
//
//	string: A newline-separated string of all messages in the chat history.
func (h *ChatHistory) GetHistory() string {
	// Create a new slice to hold messages without emoji prefixes
	sanitizedMessages := make([]string, 0, len(h.Messages))

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
		sanitizedMessages = append(sanitizedMessages, sanitizedMsg)
	}

	// Join the sanitized messages with a newline character
	return strings.Join(sanitizedMessages, "\n")
}
