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
	defer s.cleanup()

	s.setupSignalHandling()

	// Simulate AI starting the conversation by Gopher Nerd
	// This A prompt Context as starting point for AI to start the conversation
	fmt.Print(AiNerd)
	PrintTypingChat(ContextPrompt, TypingDelay)
	fmt.Println()

	// Add AI's initial message to chat history
	s.ChatHistory.AddMessage(AiNerd, ContextPrompt)

	// Main loop for processing user input
	for {
		if done := s.processInput(); done {
			break
		}
	}
}

// setupSignalHandling configures the handling of interrupt signals to ensure graceful
// shutdown of the session. It listens for SIGINT and SIGTERM signals.
func (s *Session) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Goroutine Officer (Known as Gopher Officer) to handle graceful shutdown
	go func() {
		<-sigChan // Block until a signal is received
		fmt.Println(SignalMessage)
		s.cleanup()
		os.Exit(0)
	}()
}

// processInput reads user input from the terminal. It returns true if the session
// should end, either due to a command or an error.
func (s *Session) processInput() bool {
	fmt.Print(YouNerd)
	userInput, err := bufio.NewReader(os.Stdin).ReadString(NewLineChars)
	if err != nil {
		log.Fatal(err)
	}

	userInput = strings.TrimSpace(userInput)
	if done := s.handleUserInput(userInput); done {
		return true
	}

	return false
}

// handleUserInput processes the user's input. If the input is a command, it is handled
// accordingly. Otherwise, the input is sent to the AI for a response. It returns true
// if the session should end.
func (s *Session) handleUserInput(input string) bool {
	s.ChatHistory.AddMessage(YouNerd, input)
	fmt.Println()

	if isCommand, err := HandleCommand(input, s); isCommand {
		if err != nil {
			log.Fatal(err)
		}
		return true
	}

	aiResponse, err := SendMessage(s.Ctx, s.Client, s.ChatHistory.GetHistory())
	if err != nil {
		log.Fatal(err)
	}

	s.ChatHistory.AddMessage(AiNerd, aiResponse)
	return false
}

// cleanup releases resources used by the session. It cancels the context and closes
// the AI client connection.
func (s *Session) cleanup() {
	s.Cancel()
	s.Client.Close()
}
