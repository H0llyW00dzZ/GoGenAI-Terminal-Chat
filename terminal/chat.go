// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"sync"
)

// Note: This is subject to change (for example, it can be customized in commands). For now, it's stable. Additionally, a token is inexpensive since, with Google AI's Gemini-Pro model, the maximum is 32K tokens.
const MaxChatHistory = 5 // Maximum number of messages to keep in history

// ChatHistory holds the chat messages exchanged during a session.
// It provides methods to add new messages to the history and to retrieve
// the current state of the conversation.
type ChatHistory struct {
	Messages []string
	Hashes   map[string]int // Maps hash values (as hex strings) to indices in the Messages slice
	mu       sync.RWMutex   // Explicit 🤪

}

// NewChatHistory creates and initializes a new ChatHistory struct.
// It sets up an empty slice for storing messages and initializes the hash map
// used to track unique messages. A new, random seed is generated for hashing
// to ensure the uniqueness of hash values across different instances.
//
// Returns:
//
//	*ChatHistory: A pointer to the newly created ChatHistory struct ready for use.
func NewChatHistory() *ChatHistory {
	return &ChatHistory{
		Messages: make([]string, 0),
		Hashes:   make(map[string]int),
	}
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
func (h *ChatHistory) AddMessage(user string, text string) {
	h.mu.Lock()         // Lock for writing
	defer h.mu.Unlock() // Unlock when the function returns

	// Sanitize and format the message before adding it to the history of RAM's labyrinth.
	sanitizedText := h.SanitizeMessage(text)
	message := fmt.Sprintf(ObjectHighLevelString, user, sanitizedText)
	hashValue := h.hashMessage(sanitizedText)

	// Check if the message hash already exists to prevent duplicates
	if _, exists := h.Hashes[hashValue]; !exists {
		// Remove the oldest message to maintain a fixed history size in RAM's labyrinth.
		// Note: The fixed history size might be increased in the future. Currently, the application's memory usage is minimal, consuming only 16 MB (Average).
		// then keep a maximum of 5 history entries for transmission to Google AI.
		if len(h.Messages) >= MaxChatHistory {
			// Remove the oldest message and its hash
			oldestHash := h.hashMessage(h.Messages[0])
			delete(h.Hashes, oldestHash) // Remove the hash of the oldest message
			h.Messages = h.Messages[1:]  // Remove the oldest message

			// Update the indices of the remaining hashes
			for hash, index := range h.Hashes {
				if index > 0 {
					h.Hashes[hash] = index - 1
				}
			}
		}
		// Note: this remove the oldest message are automated handle by Garbage Collector.
		// For example, free memory to avoid memory leak.
		h.Messages = append(h.Messages, message)  // Add the new message
		h.Hashes[hashValue] = len(h.Messages) - 1 // Map the hash to the new message index
	}
}

// SanitizeMessage removes ANSI color codes and other non-content prefixes from a message.
//
// Parameters:
//
//	message string: The message to be sanitized.
//
// Returns:
//
//	string: The sanitized message.
func (h *ChatHistory) SanitizeMessage(message string) string {
	// This better way to sanitize message instead of struct again.
	// It fix truncated message about color codes.
	// Note: now more Simplicity and yet powerful.
	// Remove all ANSI color codes from the message.
	return ansiRegex.ReplaceAllString(message, "")
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
	h.mu.RLock()         // Lock for reading
	defer h.mu.RUnlock() // Unlock when the function returns
	// fix concurrency issue
	var builder strings.Builder // Create a new builder for this method call

	for i, msg := range h.Messages {
		sanitizedMsg := h.SanitizeMessage(msg) // Sanitize each message
		builder.WriteString(sanitizedMsg)      // Append the sanitized message to the builder
		builder.WriteRune(nl.NewLineChars)     // Append a newline character after each message
		// After printing an AI message and if it's not the last message, add a separator
		// Note: This a better way instead of structuring it then stored in RAM's labyrinth.
		// For example how it work it's like this
		//
		// 🤓 You: :checkversion
		//
		// 🤖 AI: You are using the latest version, v0.5.0 of GoGenAI Terminal Chat. There is no need to update at the moment. Is there anything else I can help you with today?
		//
		// ---
		if i%2 == 1 && i < len(h.Messages)-1 {
			builder.WriteString(StripChars)    // Insert a separator
			builder.WriteRune(nl.NewLineChars) // Append a newline character after the separator
		}
	}

	return builder.String() // Return the complete, concatenated chat history
}

// hashMessage generates a SHA-256 hash for a given message.
func (h *ChatHistory) hashMessage(message string) string {
	hasher := sha256.New()
	hasher.Write([]byte(message))
	return hex.EncodeToString(hasher.Sum(nil))
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
// Also it used to be maintain the RAM's labyrinth hahaha and automated handle by Garbage Collector.
func (h *ChatHistory) RemoveMessages(numMessages int, messageContent string) {
	// Note: This simple and yet powerful unlike shitty complex code Hahaha.
	h.mu.Lock()
	defer h.mu.Unlock()

	if messageContent != "" {
		h.removeMessagesByContent(messageContent)
	} else {
		h.removeRecentMessages(numMessages)
	}
}

// removeMessagesByContent removes all messages that contain the specified content.
func (h *ChatHistory) removeMessagesByContent(content string) {
	// Filter out messages that do not contain the content.
	var newMessages []string
	for _, message := range h.Messages {
		if !strings.Contains(message, content) {
			newMessages = append(newMessages, message)
		} else {
			// Remove the hash of the message being removed.
			delete(h.Hashes, h.hashMessage(message))
		}
	}
	h.Messages = newMessages
}

// removeRecentMessages removes the specified number of most recent messages.
func (h *ChatHistory) removeRecentMessages(num int) {
	numToRemove := min(num, len(h.Messages))
	if numToRemove == 0 {
		return
	}
	// Remove hashes of messages being removed.
	for _, message := range h.Messages[len(h.Messages)-numToRemove:] {
		delete(h.Hashes, h.hashMessage(message))
	}
	h.Messages = h.Messages[:len(h.Messages)-numToRemove]
}

// min returns the smaller of x or y.
func min(x, y int) int {
	if x < y {
		return x
	}
	return y
}

// FilterMessages returns a slice of messages that match the predicate function.
//
// Parameters:
//
//	predicate func(string) bool: A function that returns true for messages that should be included.
//
// Returns:
//
//	[]string: A slice of messages that match the predicate.
//
// TODO: Filtering Messages
func (h *ChatHistory) FilterMessages(predicate func(string) bool) []string {
	h.mu.RLock() // Lock for reading
	defer h.mu.RUnlock()

	filtered := []string{}
	for _, msg := range h.Messages {
		if predicate(msg) {
			filtered = append(filtered, msg)
		}
	}
	return filtered
}

// Clear removes all messages from the chat history, effectively resetting it.
func (h *ChatHistory) Clear() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Messages = []string{}
	h.Hashes = make(map[string]int)
}
