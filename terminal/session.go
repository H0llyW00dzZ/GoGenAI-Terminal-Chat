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

// Session encapsulates the state and functionality for a chat session with a generative AI model.
// It holds the AI client, chat history, and context for managing the session lifecycle.
type Session struct {
	Client      *genai.Client      // Client is the generative AI client used to communicate with the AI model.
	ChatHistory ChatHistory        // ChatHistory stores the history of the chat session.
	Ctx         context.Context    // Ctx is the context governing the session, used for cancellation.
	Cancel      context.CancelFunc // Cancel is a function to cancel the context, used for cleanup.
}

// NewSession creates a new chat session with the provided API key for authentication.
// It initializes the generative AI client and sets up a context for managing the session.
//
// Parameters:
//
//	apiKey string: The API key used for authenticating requests to the AI service.
//
// Returns:
//
//	*Session: A pointer to the newly created Session object.
//	error: An error object if initialization fails.
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

// Start begins the chat session, managing user input and AI responses.
// It sets up a signal listener for graceful shutdown and enters a loop to
// read user input and fetch AI responses indefinitely until an interrupt signal is received.
//
// This method handles user input errors and AI communication errors by logging them and exiting.
// It ensures resources are cleaned up properly on exit by deferring the cancellation of the session's context
// and the closure of the AI client.
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
