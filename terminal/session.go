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
	Ended       bool               // Ended indicates whether the session has ended.
	mutex       sync.Mutex         // Mutex is a mutex to ensure thread-safe access to the session's state.
	lastInput   string             // Stores the last user input for reference

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
	// Note: This doesn't use a storage system like a database or file system to keep the chat history, nor does it use a JSON structure (as a front-end might) for sending request to Google AI.
	// So if you're wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
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
	// Note: This is securely managed by the Gopher Officer, which handles the session and is linked to the `processInput` function.
	// Additionally, the Gopher Officer may occasionally sleep during the session's lifecycle and will wake up when needed.
	defer s.cleanup()
	// This Automated Spawn another Goroutine Officer (Known as Gopher Officer) to handle signal
	s.setupSignalHandling()

	// Simulate AI starting the conversation by Gopher Nerd
	// This is a prompt context as the starting point for AI to start the conversation
	fmt.Print(AiNerd)
	PrintTypingChat(ContextPrompt, TypingDelay)
	fmt.Println() // Ensure there's a newline after the AI's initial message

	// Add AI's initial message to chat history
	s.ChatHistory.AddMessage(AiNerd, ContextPrompt)

	// Main loop for processing user input
	for {
		select {
		case <-s.Ctx.Done():
			fmt.Println(ContextCancel)
			return
		default:
			done := s.processInput()
			if done {
				return
			}
			// Add a newline after handling input, but only if it's not a command
			// This helps separate the AI's response from the user's next prompt
			if !strings.HasPrefix(strings.TrimSpace(s.lastInput), PrefixChar) {
				fmt.Println()
			}
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
		logger.Error(ErrorReadingUserInput, err)
		return false // Continue the loop, hoping for a successful read next time
	}

	userInput = strings.TrimSpace(userInput)
	s.lastInput = userInput // Store the last input

	if isCommand(userInput) {
		return s.handleCommand(userInput)
	}
	return s.handleUserInput(userInput)
}

// isCommand checks if the input is a command based on the prefix.
func isCommand(input string) bool {
	return strings.HasPrefix(input, PrefixChar)
}

// handleCommand processes the input as a command and returns true if the session should end.
func (s *Session) handleCommand(input string) bool {
	if isCommand, err := HandleCommand(input, s); isCommand {
		if err != nil {
			logger.Error(ErrorHandlingCommand, err)
		}
		// If it's a command, whether it's handled successfully or not, we continue the session
		return false
	}
	return false
}

// handleUserInput processes the user's input. If the input is a command, it is handled
// accordingly. Otherwise, the input is sent to the AI for a response. It returns true
// if the session should end.
func (s *Session) handleUserInput(input string) bool {
	// Check if the session is still valid
	if s.Client == nil {
		// Attempt to renew the session
		if err := s.RenewSession(apiKey); err != nil {
			// Handle the error, possibly by logging and returning true to signal the session should end
			logger.Error(ErrorFailedToRenewSession, err)
			return true // Signal to end the session
		}
	}

	s.ChatHistory.AddMessage(YouNerd, input)
	fmt.Println() // Ensure a newline after the user's input

	aiResponse, err := SendMessage(s.Ctx, s.Client, input)
	if err != nil {
		logger.Error(ErrorSendingMessage, err)
		return true // End the session if there's an error sending the message
	}

	// Add the AI's response to the chat history
	s.ChatHistory.AddMessage(AiNerd, aiResponse)
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

// RenewSession attempts to renew the client session with the AI service.
func (s *Session) RenewSession(apiKey string) error {
	s.mutex.Lock()         // Lock the mutex before accessing shared resources
	defer s.mutex.Unlock() // Ensure the mutex is unlocked at the end of the method

	// Close the current session if it exists
	if s.Client != nil {
		s.Client.Close() // Assuming Close is the method to properly shutdown the client
		s.Client = nil   // Set the client to nil after closing
	}

	// Create a new client for the session
	var err error
	s.Client, err = genai.NewClient(s.Ctx, option.WithAPIKey(apiKey))
	if err != nil {
		// this low level error not possible to use logger.Error
		return fmt.Errorf(ErrorLowLevelFailedtoStartAiChatSession, err)
	}

	return nil
}
