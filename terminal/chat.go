// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
)

// DetermineMessageType analyzes the content of a message to classify its type.
// It returns the MessageType based on predefined criteria for identifying user, AI, and system messages.
func DetermineMessageType(message string) MessageType {
	if isSysMessage(message) {
		return SystemMessage
	} else if isAIMessage(message) {
		return AIMessage
	}
	return UserMessage
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

// AddMessage appends a new message to the chat history.
// It takes the username and the text of the message as inputs and formats
// them before adding to the Messages slice.
//
// Parameters:
//
//	user string: The username of the individual sending the message.
//	text string: The content of the message to be added to the history.
//	config *ChatConfig: Configuration parameters for the chat session, including history size.
//
// This method does not return any value or error. It assumes that all input
// is valid and safe to add to the chat history.
func (h *ChatHistory) AddMessage(user string, text string, config *ChatConfig) {
	// Warning!!! Explicit 🤪
	h.mu.Lock()         // Lock for writing
	defer h.mu.Unlock() // Unlock when the function returns

	// Sanitize and format the message before adding it to the history of RAM's labyrinth.
	sanitizedText := h.SanitizeMessage(text)
	message := fmt.Sprintf(ObjectHighLevelStringWithNewLine, user, sanitizedText) // Add newlines around the message
	hashValue := h.hashMessage(sanitizedText)
	messageType := DetermineMessageType(sanitizedText)

	// Delegate message handling based on type.
	// Note: This becomes easier to maintain by Go routines.
	switch messageType {
	case SystemMessage:
		h.handleSystemMessage(sanitizedText, message, hashValue)
	case AIMessage:
		h.handleAIMessage(message, hashValue)
	default:
		h.handleUserMessage(user, message, hashValue)
	}
	// Check if the message hash already exists to prevent duplicates
	h.manageHistorySize(config)
}

// handleSystemMessage checks and replaces an existing system message.
// Returns true if a system message was handled (either replaced or identified for addition).
func (h *ChatHistory) handleSystemMessage(sanitizedText, message, hashValue string) bool {
	if !isSysMessage(sanitizedText) {
		return false // Not a system message, no action required.
	}
	// Warning!!! Explicit 🤪
	h.mu.Lock()         // Lock for writing
	defer h.mu.Unlock() // Ensure unlocking
	// Check if there is an existing system message
	if h.isExistingSysMessage(hashValue) {
		h.replaceExistingSysMessage(message, hashValue)
	} else {
		h.addNewSysMessage(message, hashValue)
	}
	h.cleanupOldSysMessages()
	return true // Indicate a system message was handled.
}

// handleAIMessage processes an AI message.
func (h *ChatHistory) handleAIMessage(message, hashValue string) {

	if _, exists := h.Hashes[hashValue]; !exists {
		h.addMessageToHistory(message, hashValue)
		h.AIMessageCount++
	}
}

// handleUserMessage processes a user message.
func (h *ChatHistory) handleUserMessage(user, message, hashValue string) {
	if _, exists := h.Hashes[hashValue]; exists {
		return
	}

	h.addMessageToHistory(message, hashValue)
	h.updateMessageCounts(user)
}

// updateMessageCounts updates the message counts based on the user.
func (h *ChatHistory) updateMessageCounts(user string) {

	if user == SYSTEMPREFIX {
		h.updateSystemMessageCount()
	} else if user == AiNerd {
		h.AIMessageCount++
	} else if user == YouNerd {
		h.UserMessageCount++
	}
}

// updateSystemMessageCount increments the system message count if it's currently zero.
func (h *ChatHistory) updateSystemMessageCount() {
	if h.SystemMessageCount == 0 {
		h.SystemMessageCount = 1
	}
}

// addMessageToHistory adds a message to the history.
func (h *ChatHistory) addMessageToHistory(message, hashValue string) {
	// Note: this remove the oldest message are automated handle by Garbage Collector.
	// For example, free memory to avoid memory leak.
	h.Messages = append(h.Messages, message)  // Add the new message
	h.Hashes[hashValue] = len(h.Messages) - 1 // Map the hash to the new message index
}

// isExistingSysMessage checks if there is an existing system message with the same hash.
func (h *ChatHistory) isExistingSysMessage(hashValue string) bool {
	_, exists := h.Hashes[hashValue]
	return exists && isSysMessage(h.Messages[h.Hashes[hashValue]])
}

// replaceExistingSysMessage replaces the existing system message with the new one.
func (h *ChatHistory) replaceExistingSysMessage(message, hashValue string) {
	existingIndex := h.Hashes[hashValue]
	h.Messages[existingIndex] = message
}

// addNewSysMessage adds a new system message to the history.
func (h *ChatHistory) addNewSysMessage(message, hashValue string) {
	h.Messages = append(h.Messages, message)
	h.Hashes[hashValue] = len(h.Messages) - 1
}

// cleanupOldSysMessages removes older system messages from the history.
func (h *ChatHistory) cleanupOldSysMessages() {
	for i := len(h.Messages) - 2; i >= 0; i-- {
		if isSysMessage(h.Messages[i]) {
			h.Messages = append(h.Messages[:i], h.Messages[i+1:]...)
			break // Assuming only one system message exists at a time.
		}
	}
}

// manageHistorySize manages the size of the chat history based on the ChatConfig.
func (h *ChatHistory) manageHistorySize(config *ChatConfig) {
	// Remove the oldest two messages (one user and one AI) to maintain a fixed history size in RAM's labyrinth.
	// Note: The fixed history size might be increased in the future. Currently, the application's memory usage is minimal, consuming only 16 MB (Average).
	// then keep a maximum of 10 history entries for transmission to Google AI.
	for len(h.Messages) > config.HistorySize*2 {
		oldestUserHash := h.hashMessage(h.Messages[0])
		oldestAIHash := h.hashMessage(h.Messages[1])
		delete(h.Hashes, oldestUserHash) // Remove the hash of the oldest user message
		delete(h.Hashes, oldestAIHash)   // Remove the hash of the oldest AI message
		h.Messages = h.Messages[2:]      // Remove the oldest two messages
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
func (h *ChatHistory) GetHistory(config *ChatConfig) string {
	// Note: if you're still wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
	// Define the prefixes to be removed
	// Additional Note: If issues still arise due to ANSI color codes in AI responses, it's not because of the 'this' or 'Colorize' function in Genai.go.
	// The issue lies with the AI's attempt to apply formatting, which fails due to incorrect ANSI sequences, reminiscent of issues one might encounter with "PYTHON" or Your Machine is bad LMAO.
	h.mu.RLock()         // Lock for reading
	defer h.mu.RUnlock() // Unlock when the function returns

	// Determine the starting index based on the number of messages to include.
	// Additional Note: This required go1.21.0 ~ latest
	// Ref: https://pkg.go.dev/builtin#max
	startIndex := max(0, len(h.Messages)-config.HistorySize)
	historySubset := h.Messages[startIndex:]

	return h.buildHistoryString(historySubset)
}

// buildHistoryString builds the chat history string from a subset of messages.
func (h *ChatHistory) buildHistoryString(historySubset []string) string {
	// Use a strings.Builder to build the chat history string efficiently.
	builder := strings.Builder{}

	// Check for system messages and prepend them to the history.
	sysMsgs, chatMsgs := h.separateSystemMessages(historySubset)
	h.appendSystemMessages(&builder, sysMsgs)
	h.appendChatMessages(&builder, chatMsgs)

	return builder.String()
}

// appendSystemMessages appends system messages to the StringBuilder.
func (h *ChatHistory) appendSystemMessages(builder *strings.Builder, sysMsgs []string) {
	// Keep track of the latest system message index
	latestSysMsgIndex := -1

	for i, sysMsg := range sysMsgs {
		if !isSysMessage(sysMsg) {
			continue // Skip non-system messages
		}

		// Check if this system message is the latest one
		if latestSysMsgIndex == -1 {
			latestSysMsgIndex = i
		} else {
			// Remove the older system message
			builder.Reset()
			latestSysMsgIndex = i
		}

		builder.WriteString(sysMsg)
		builder.WriteRune(nl.NewLineChars)
		builder.WriteString(StripChars)    // Append the separator
		builder.WriteRune(nl.NewLineChars) // Append a newline character after the separator
		builder.WriteRune(nl.NewLineChars) // Append an extra newline character after the system message
	}
}

// appendChatMessages appends chat messages to the StringBuilder, adding separators as needed.
func (h *ChatHistory) appendChatMessages(builder *strings.Builder, chatMsgs []string) {
	totalMessages := len(chatMsgs)
	for i, msg := range chatMsgs {
		sanitizedMsg := h.SanitizeMessage(msg) // Sanitize each message
		builder.WriteString(sanitizedMsg)      // Append the sanitized message to the builder
		builder.WriteRune(nl.NewLineChars)     // Append a newline character after each message
		// After printing an AI message and if it's not the last message, add a separator
		// Note: This a better way instead of structuring it then stored in RAM's labyrinth.
		// For example how it work it's like this
		//
		// ⚙️  SYSTEM: Discussion Summary:
		//
		// ---
		//
		// 🤓 You: :checkversion
		//
		// 🤖 AI: You are using the latest version, v0.5.0 of GoGenAI Terminal Chat. There is no need to update at the moment. Is there anything else I can help you with today?
		//
		// ---
		//
		// Add a separator after an AI message if the next message is a user message
		// Determine if a separator should be appended after the current message.
		if i < totalMessages-1 && h.shouldAppendSeparator(sanitizedMsg, chatMsgs[i+1], i, totalMessages) {
			builder.WriteString(StripChars)    // Append the separator
			builder.WriteRune(nl.NewLineChars) // Append a newline character after the separator
			builder.WriteRune(nl.NewLineChars) // Append a newline character after the separator for the user
		}
	}
}

// shouldAppendSeparator determines if a separator should be added between messages.
func (h *ChatHistory) shouldAppendSeparator(currentMessage, nextMessage string, currentIndex, totalMessages int) bool {
	// Check if the current message is from the AI and the next message is from the user.
	return isAIMessage(currentMessage) &&
		currentIndex < totalMessages-1 &&
		isUserMessage(nextMessage)
}

// separateSystemMessages separates system messages from chat messages.
func (h *ChatHistory) separateSystemMessages(messages []string) (sysMsgs, chatMsgs []string) {
	for _, msg := range messages {
		if isSysMessage(msg) {
			sysMsgs = append(sysMsgs, msg)
		} else {
			chatMsgs = append(chatMsgs, msg)
		}
	}
	return sysMsgs, chatMsgs
}

// isUserMessage checks if the message is from the user
func isUserMessage(message string) bool {
	// Assuming YouNerd is the prefix for user messages
	return strings.HasPrefix(message, YouNerd)
}

// isAIMessage checks if the message is from the AI
func isAIMessage(message string) bool {
	return strings.HasPrefix(message, AiNerd)
}

// isSysMessage checks if the message is from the System
func isSysMessage(message string) bool {
	return strings.HasPrefix(message, SYSTEMPREFIX)
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
	h.AIMessageCount = 0
	h.SystemMessageCount = 0
	h.UserMessageCount = 0
}

// Note: This a different way unlike "Clear"
func (h *ChatHistory) cleanup() {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.Messages = nil
	h.Hashes = nil
	h.AIMessageCount = 0
	h.SystemMessageCount = 0
	h.UserMessageCount = 0
}

// ClearAllSystemMessages removes all system messages from the chat history.
func (h *ChatHistory) ClearAllSystemMessages() {
	h.mu.Lock()         // Lock for writing
	defer h.mu.Unlock() // Ensure unlocking

	var newMessages []string
	var newHashes = make(map[string]int)
	h.SystemMessageCount = 0 // Reset the system message count

	for _, message := range h.Messages {
		if !isSysMessage(message) {
			// This is not a system message; keep it.
			newMessages = append(newMessages, message)
			hashValue := h.hashMessage(message)
			newHashes[hashValue] = len(newMessages) - 1
		}
	}

	// Replace the old Messages and Hashes with the new ones that exclude system messages.
	h.Messages = newMessages
	h.Hashes = newHashes
}

// GetMessageStats safely retrieves the message counts from the ChatHistory instance.
// It returns a MessageStats struct containing the counts of user, AI, and system messages.
// Access to the ChatHistory's message counts is read-locked to ensure thread safety.
func (h *ChatHistory) GetMessageStats() MessageStats {
	h.mu.RLock()         // Lock for reading
	defer h.mu.RUnlock() // Unlock when the function exits

	// Return a new instance of MessageStats with the current counts.
	return MessageStats{
		UserMessages:   h.UserMessageCount,
		AIMessages:     h.AIMessageCount,
		SystemMessages: h.SystemMessageCount,
	}
}

// HasSystemMessages checks if there are any system messages in the chat history.
func (h *ChatHistory) HasSystemMessages() bool {
	h.mu.RLock()         // Lock for reading
	defer h.mu.RUnlock() // Ensure unlocking

	// Check if the SystemMessageCount is greater than 0.
	return h.SystemMessageCount > 0
}
