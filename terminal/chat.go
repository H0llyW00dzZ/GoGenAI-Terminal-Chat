// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

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
	// Forgot to remove "if" previously, an "if" statement here was causing an google internal error when
	// Shift+Enter was pressed, due to unintended handling of input. This has been corrected.
	// Additionally, the issue of duplicate "🤖 AI:" prefixes, which resulted in the AI
	// misinterpreting its own output, has been resolved.
	// Now, the message is simply appended to the history with the correct user prefix.
	//
	// Note: the duplicate of "🤖 AI:" it's by AI It's self, not causing of the code.
	// Now It fixed by sanitize the message before append it to the history.
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
