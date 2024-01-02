package terminal

import (
	"fmt"
	"strings"
)

type ChatHistory struct {
	Messages []string
}

func (h *ChatHistory) AddMessage(user, text string) {
	h.Messages = append(h.Messages, fmt.Sprintf("%s: %s", user, text))
}

func (h *ChatHistory) GetHistory() string {
	return strings.Join(h.Messages, "\n")
}

func (h *ChatHistory) PrintHistory() {
	for _, msg := range h.Messages {
		fmt.Println(msg)
	}
}
