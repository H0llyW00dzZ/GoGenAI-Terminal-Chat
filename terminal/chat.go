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

// NewLineChar is a struct that containt Rune for New Line Character
type NewLineChar struct {
	NewLineChars rune
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
	message := fmt.Sprintf(ObjectHighLevelString, user, text)
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
	// Additional Note: If issues still arise due to ANSI color codes in AI responses, it's not because of the 'this' or 'Colorize' function in Genai.go.
	// The issue lies with the AI's attempt to apply formatting, which fails due to incorrect ANSI sequences, reminiscent of issues one might encounter with "PYTHON" or Your Machine is bad LMAO.
	prefixesToRemove := []string{youNerd, aiNerd, ColorGreen, ColorYellow, ColorBlue, ColorPurple, ColorCyan, ColorReset} // Fix all Issues

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
		buildeR.WriteString(sanitizedMsg)
		buildeR.WriteRune(nl.NewLineChars) // Append a newline character after each message.
	}

	// The buildeR.String() method returns the complete, concatenated chat history.
	return buildeR.String()

}

// RemoveMessages removes messages from the chat history. If a specific message is provided,
// it removes messages that contain that text; otherwise, it removes the specified number of
// most recent messages.
//
// Parameters:
//
//	numMessages int: The number of most recent messages to remove. If set to 0 and a specific
//	                 message is provided, all instances of that message are removed.
//	messageContent string: The specific content of messages to remove. If empty, it removes
//	                       the number of most recent messages specified by numMessages.
//
// This method does not return any value. It updates the chat history in place.
//
// Note: This currently marked as TODO since it's not used anywhere in the code. It's a good idea to add this feature in the future.
func (h *ChatHistory) RemoveMessages(numMessages int, messageContent string) {
	// Note: This simple and yet powerful unlike shitty complex code Hahaha.
	if messageContent != "" {
		h.removeMessagesByContent(messageContent)
		return
	}

	// If numMessages is provided, remove the most recent messages.
	if numMessages > 0 {
		h.removeRecentMessages(numMessages)
	}
}

// removeMessagesByContent removes all messages that contain the specified content.
func (h *ChatHistory) removeMessagesByContent(content string) {
	filteredMessages := h.filterMessages(func(msg string) bool {
		return !strings.Contains(msg, content)
	})
	h.Messages = filteredMessages
}

// removeRecentMessages removes the specified number of most recent messages.
func (h *ChatHistory) removeRecentMessages(count int) {
	if count <= len(h.Messages) {
		h.Messages = h.Messages[:len(h.Messages)-count]
	}
}

// filterMessages filters the messages using the provided predicate function.
func (h *ChatHistory) filterMessages(predicate func(string) bool) []string {
	filtered := []string{}
	for _, msg := range h.Messages {
		if predicate(msg) {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}
