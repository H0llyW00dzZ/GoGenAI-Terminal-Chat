package terminal

import "fmt"

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
