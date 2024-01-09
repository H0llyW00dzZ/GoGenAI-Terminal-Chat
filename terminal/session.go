// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// Session encapsulates the state and functionality for a chat session with a generative AI model.
// It holds the AI client, chat history, and context for managing the session lifecycle.
type Session struct {
	Client        *genai.Client      // Client is the generative AI client used to communicate with the AI model.
	ChatHistory   ChatHistory        // ChatHistory stores the history of the chat session.
	Ctx           context.Context    // Ctx is the context governing the session, used for cancellation.
	Cancel        context.CancelFunc // Cancel is a function to cancel the context, used for cleanup.
	AiChatSession *genai.ChatSession // AiChatSession is the chat session with the generative AI model.
	Ended         bool               // Ended is a flag to indicate if the session has ended.
	mutex         sync.Mutex         // Mutex is a mutex to ensure thread-safe access to the session's state.

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
	client, aiChatSession, ctx, cancel, err := createChatSessionWithRetries(apiKey)
	if err != nil {
		return nil, err
	}

	// Note: This doesn't use a storage system like a database or file system to keep the chat history, nor does it use a JSON structure (as a front-end might) for sending request to Google AI.
	// So if you're wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
	return &Session{
		Client:        client,
		ChatHistory:   ChatHistory{},
		Ctx:           ctx,
		Cancel:        cancel,
		AiChatSession: aiChatSession,
	}, nil
}

func createChatSessionWithRetries(apiKey string) (*genai.Client, *genai.ChatSession, context.Context, context.CancelFunc, error) {
	maxRetries := 3
	for attempt := 0; attempt < maxRetries; attempt++ {
		client, ctx, cancel, err := initializeClient(apiKey)
		if err != nil {
			return nil, nil, nil, nil, err
		}

		aiChatSession, err := startChatSession(client)
		if err == nil {
			return client, aiChatSession, ctx, cancel, nil
		}

		// Cleanup before retrying
		cleanup(client, cancel)
		logger.Error(ErrorFailedToStartAIChatSessionAttempt, attempt+1, maxRetries)
		time.Sleep(1 * time.Second)
	}

	errMsg := fmt.Sprintf(ErrorFailedtoStartAiChatSessionAfter, maxRetries)
	logger.Error(errMsg)
	return nil, nil, nil, nil, fmt.Errorf(errMsg)
}

func initializeClient(apiKey string) (*genai.Client, context.Context, context.CancelFunc, error) {
	ctx, cancel := context.WithCancel(context.Background())
	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		cancel()
		return nil, nil, nil, err
	}
	return client, ctx, cancel, nil
}

func startChatSession(client *genai.Client) (*genai.ChatSession, error) {
	model := client.GenerativeModel(ModelAi)
	aiChatSession := model.StartChat()
	if aiChatSession == nil {
		errMsg := ErrorFailedtoStartAiChatSession
		logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	return aiChatSession, nil
}

func cleanup(client *genai.Client, cancel context.CancelFunc) {
	cancel()
	if client != nil {
		client.Close()
	}
}

// Start begins the chat session, managing user input and AI responses.
// It sets up a signal listener for graceful shutdown and enters a loop to
// read user input and fetch AI responses indefinitely until an interrupt signal is received.
//
// This method handles user input errors and AI communication errors by logging them and exiting.
// It ensures resources are cleaned up properly on exit by deferring the cancellation of the session's context
// and the closure of the AI client.
func (s *Session) Start() {
	// Note: This is securely managed by the Gopher Officer, which handles the session and is linked to the `processInput` function.
	// Additionally, the Gopher Officer may occasionally sleep during the session's lifecycle and will wake up when needed.
	defer s.cleanup()
	// This Automated Spawn another Goroutine Officer (Known as Gopher Officer) to handle signal
	s.setupSignalHandling()

	// Simulate AI starting the conversation by Gopher Nerd
	// This is a prompt context as the starting point for AI to start the conversation
	fmt.Print(AiNerd)
	PrintTypingChat(ContextPrompt, TypingDelay)
	fmt.Println() // A better newline instead of hardcoding "\n"
	fmt.Println() // A better newline instead of hardcoding "\n"

	// Add AI's initial message to chat history
	s.ChatHistory.AddMessage(AiNerd, ContextPrompt)

	// Main loop for processing user input
	for {
		fmt.Print(YouNerd)
		if done := s.processInput(); done {
			break // Exit the loop if processInput signals to stop
		}
		fmt.Println() // A better newline instead of hardcoding "\n"
		fmt.Println() // A better newline instead of hardcoding "\n"
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
	userInput, err := bufio.NewReader(os.Stdin).ReadString(NewLineChars)
	if err != nil {
		logger.Error(ErrorReadingUserInput, err)
		return false // Continue the loop, hoping for a successful read next time
	}

	userInput = strings.TrimSpace(userInput)

	// Check if the input is a command and handle it
	if isCommand, err := HandleCommand(userInput, s); isCommand {
		if err != nil {
			logger.Error(ErrorHandlingCommand, err)
		}
		return true // Return true to indicate that the session should end
	}

	// If the input is not a command, send it to the AI as part of the chat history
	return s.handleUserInput(userInput)
}

// handleUserInput processes the user's input. If the input is a command, it is handled
// accordingly. Otherwise, the input is sent to the AI for a response. It returns true
// if the session should end.
func (s *Session) handleUserInput(input string) bool {
	// Check if the session is still valid
	if s.AiChatSession == nil {
		// Attempt to renew the session
		if err := s.RenewSession(); err != nil {
			// Handle the error, possibly by logging and returning true to signal the session should end
			logger.Error("Failed to renew session: %v", err)
			return true // Signal to end the session
		}
	}

	s.ChatHistory.AddMessage(YouNerd, input)
	fmt.Println()

	// Pass the existing AI chat session to SendMessage
	aiResponse, err := SendMessage(s.Ctx, s.AiChatSession, s.ChatHistory.GetHistory())
	if err != nil {
		logger.Error(ErrorSendingMessage, err)
		return true // Signal to end the session due to an error
	} else {
		s.ChatHistory.AddMessage(AiNerd, aiResponse) // Add AI response to history
	}

	return false // Continue the session
}

// cleanup releases resources used by the session. It cancels the context and closes
// the AI client connection.
func (s *Session) cleanup() {
	s.Cancel()
	s.Client.Close()
}

// endSession terminates the chat session and performs necessary cleanup operations. It should be
// invoked in response to user commands that signify the end of a session, system interrupts, or
// internal errors that require the session to be closed.
//
// After calling endSession, the session's resources are released, and it should not be used further.
func (s *Session) endSession() {
	s.cleanup()    // Perform cleanup operations
	s.Ended = true // Mark the session as ended
}

// HasEnded reports whether the chat session has ended. It can be called at any point to
// check the session's state without altering it.
//
// Return value:
//
//	ended bool: A boolean indicating true if the session has ended, or false if it is still active.
func (s *Session) HasEnded() (ended bool) {
	return s.Ended
}
