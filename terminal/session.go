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
	// mu protects the concurrent access to session's state, ensuring thread safety.
	// It should be locked when accessing or modifying the session's state.
	mu sync.Mutex
	// this reference pretty useful, which can handle runtime 24/7, unlike original ai chat session systems.
	// for example, if session is ended not cause of client, then it will be renew with previous chat history.
	lastInput string // Stores the last user input for reference

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
func NewSession(apiKey string) *Session {
	ctx, cancel := context.WithCancel(context.Background())

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		cancel()
		logger.Error(ErrorFailedToCreateNewAiClient, err)
		logger.HandleGoogleAPIError(err)
		return nil
	}

	// Perform a simple request to validate the API key.
	valid, err := SendDummyMessage(client)
	if err != nil || !valid {
		cancel()
		logger.Error(ErrorInvalidApiKey, err)
		logger.HandleGoogleAPIError(err)
		return nil
	}
	// Note: This doesn't use a storage system like a database or file system to keep the chat history, nor does it use a JSON structure (as a front-end might) for sending request to Google AI.
	// So if you're wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
	return &Session{
		Client:      client,
		ChatHistory: ChatHistory{},
		Ctx:         ctx,
		Cancel:      cancel,
	}
}

// SendDummyMessage verifies the validity of the API key by sending a dummy message.
//
// Parameters:
//
//	client *genai.Client: The AI client used to send the message.
//
// Returns:
//
//	A boolean indicating the validity of the API key.
//	An error if sending the dummy message fails.
func SendDummyMessage(client *genai.Client) (bool, error) {
	// Initialize a dummy chat session or use an appropriate lightweight method.
	model := client.GenerativeModel(ModelAi)
	cs := model.StartChat()

	// Attempt to send a dummy message.
	resp, err := cs.SendMessage(context.Background(), genai.Text(DummyMessages))
	if err != nil {
		return false, err
	}

	// A non-nil response indicates a valid API key.
	return resp != nil, nil
}

// Start begins the chat session, managing user input and AI responses.
// It sets up a signal listener for graceful shutdown and enters a loop to
// read user input and fetch AI responses indefinitely until an interrupt signal is received.
//
// This method handles user input errors and AI communication errors by logging them and exiting.
// It ensures resources are cleaned up properly on exit by deferring the cancellation of the session's context
// and the closure of the AI client.
func (s *Session) Start() {
	text := "G V"
	asciiArt := ToASCIIArt(text)
	fmt.Println(asciiArt)
	fmt.Println()
	// Note: This is securely managed by the Gopher Officer, which handles the session and is linked to the `processInput` function.
	// Additionally, the Gopher Officer may occasionally sleep during the session's lifecycle and will wake up when needed.
	defer s.cleanup()
	// This Automated Spawn another Goroutine Officer (Known as Gopher Officer) to handle signal
	s.setupSignalHandling()

	// Simulate AI starting the conversation by Gopher Nerd
	// This is a prompt context as the starting point for AI to start the conversation
	PrintPrefixWithTimeStamp(AiNerd)
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
		}
	}
}

// setupSignalHandling configures the handling of interrupt signals to ensure graceful
// shutdown of the session. It listens for SIGINT and SIGTERM signals.
func (s *Session) setupSignalHandling() {
	sigChan := make(chan os.Signal, 1)
	// Note: by refactoring a logic like this, it easier monitoring other signal in linux/unix or windows, also it easier catch other signal.
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Gopher Officer to handle graceful shutdown, and monitoring other signal in linux/unix or windows.
	go func() {
		for {
			sig := <-sigChan // Block until a signal is received
			switch sig {
			case syscall.SIGINT, syscall.SIGTERM:
				// Perform cleanup and exit only on SIGINT and SIGTERM.
				fmt.Println(SignalMessage)
				s.cleanup()
				os.Exit(0)
			default:
				fmt.Printf(MonitoringSignal, sig)
			}
		}
	}()
}

// processInput reads user input from the terminal. It returns true if the session
// should end, either due to a command or an error.
func (s *Session) processInput() bool {
	PrintPrefixWithTimeStamp(YouNerd)
	userInput, err := bufio.NewReader(os.Stdin).ReadString(byte(nl.NewLineChars))
	if err != nil {
		logger.Error(ErrorReadingUserInput, err)
		logger.HandleGoogleAPIError(err)
		return false // Continue the loop, hoping for a successful read next time
	}

	userInput = strings.TrimSpace(userInput)
	s.lastInput = userInput // Store the last input

	if isCommand(userInput) {
		return s.handleCommand(userInput)
	}
	return s.handleUserInput(userInput)
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
			logger.HandleGoogleAPIError(err)
			return true // Signal to end the session
		}
	}

	// Add the user's input to the chat history
	s.ChatHistory.AddMessage(YouNerd, input)

	// Get the entire chat history as a string
	chatHistory := s.ChatHistory.GetHistory()

	// Pass Better LLm's Send the user input along with the chat history to the AI
	// Note: This not using json struct of candidate (e.g, User and model in json struct), lmao but it more better.
	aiResponse, err := SendMessage(s.Ctx, s.Client, input, chatHistory)
	if err != nil {
		logger.Error(ErrorSendingMessage, err)
		logger.HandleGoogleAPIError(err)
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
//
// TODO: Utilize this in multiple goroutines, such as for task queues, terminal control, etc.
func (s *Session) HasEnded() (ended bool) {
	return s.Ended
}

// RenewSession attempts to renew the client session with the AI service by reinitializing
// the genai.Client with the provided API key. This method is useful when the existing
// client session has expired or is no longer valid and a new session needs to be established
// to continue communication with the AI service.
//
// The method ensures thread-safe access by using a mutex lock during the client reinitialization
// process. If a client session already exists, it is properly closed and a new client is created.
//
// Parameters:
//
//	apiKey string: The API key used for authenticating requests to the AI service.
//
// Returns:
//
//	error: An error object if reinitializing the client fails. If the operation is successful,
//	       the error is nil.
//
// Upon successful completion, the Session's Client field is updated to reference the new
// genai.Client instance. In case of failure, an error is returned and the Client field is set to nil.
func (s *Session) RenewSession(apiKey string) error {
	s.mu.Lock()         // Lock the mutex before accessing shared resources
	defer s.mu.Unlock() // Ensure the mutex is unlocked at the end of the method

	// Close the current session if it exists
	if s.Client != nil {
		//s.Client.Close() // Assuming Close is the method to properly shutdown the client
		// This just ensure that not looping the chat history after RenewSession.
		defer s.cleanup()
		s.Client = nil // Set the client to nil after closing
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
