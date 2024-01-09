// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

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

// RenewSession attempts to renew the chat session with the AI service.
func (s *Session) RenewSession() error {
	s.mutex.Lock()         // Lock the mutex before accessing shared resources
	defer s.mutex.Unlock() // Ensure the mutex is unlocked at the end of the method

	// Close the current session if it exists
	if s.AiChatSession != nil {
		s.endSession()
		s.AiChatSession = nil // Explicitly set to nil to avoid using a closed session
	}

	// Create a new session
	newSession, err := startChatSession(s.Client)
	if err != nil {
		return err
	}

	// Replace the old session with the new one
	s.AiChatSession = newSession
	return nil
}
