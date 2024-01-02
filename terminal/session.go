package terminal

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

type Session struct {
	Client      *genai.Client
	ChatHistory ChatHistory
	Ctx         context.Context
	Cancel      context.CancelFunc
}

func NewSession(apiKey string) (*Session, error) {
	ctx, cancel := context.WithCancel(context.Background())

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		cancel()
		return nil, err
	}

	return &Session{
		Client:      client,
		ChatHistory: ChatHistory{},
		Ctx:         ctx,
		Cancel:      cancel,
	}, nil
}

func (s *Session) Start() {
	defer s.Cancel()
	defer s.Client.Close()

	reader := bufio.NewReader(os.Stdin)

	// Print the chat history before starting the conversation
	s.ChatHistory.PrintHistory()

	// Set up channel to listen for interrupt signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine to handle graceful shutdown
	go func() {
		<-sigChan // Block until a signal is received
		fmt.Println(SignalMessage)
		s.Cancel() // Cancel the context to cleanup resources
		s.Client.Close()
		os.Exit(0)
	}()

	for {
		fmt.Print(YouNerd)
		userInput, err := reader.ReadString(NewLineChars)
		if err != nil {
			log.Fatal(err)
		}

		userInput = strings.TrimSpace(userInput)
		s.ChatHistory.AddMessage(YouNerd, userInput)

		// Pass the entire chat history as context for the AI's response
		chatContext := s.ChatHistory.GetHistory()
		aiResponse, err := SendMessage(s.Ctx, s.Client, chatContext)
		if err != nil {
			log.Fatal(err)
		}

		s.ChatHistory.AddMessage(AiNerd, aiResponse)
	}
}
