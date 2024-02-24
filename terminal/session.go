// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

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
	// Initialize ChatConfig with default values.
	chatConfig := DefaultChatConfig()
	ctx, cancel := context.WithCancel(context.Background())

	client, err := genai.NewClient(ctx, option.WithAPIKey(apiKey))
	if err != nil {
		cancel()
		logger.Error(ErrorFailedToCreateNewAiClient, err)
		return nil
	}

	// Define a retryable operation for validating the API key.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			return SendDummyMessage(client)
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	valid, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)

	// Handle the result of the retry operation.
	if err != nil || !valid {
		cancel()
		logger.Error(ErrorFailedToStartSession, err)
		return nil
	}
	// Note: This doesn't use a storage system like a database or file system to keep the chat history, nor does it use a JSON structure (as a front-end might) for sending request to Google AI.
	// So if you're wondering where this is all stored, it's in a place you won't find—somewhere in the RAM's labyrinth, hahaha!
	// Initialize the ChatHistory here instead of using an empty struct
	chatHistory := NewChatHistory() // Hash RAM's labyrinth, hahaha!
	return &Session{
		Client:           client,
		ChatHistory:      chatHistory, // Store the pointer to ChatHistory in RAM's labyrinth
		ChatConfig:       chatConfig,  // Initialize ChatConfig
		SafetySettings:   DefaultSafetySettings(),
		DefaultModelName: GeminiPro, // Set the default model name
		Ctx:              ctx,
		Cancel:           cancel,
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
	// Merge styles before using.
	combinedStyle := MergeStyles(slantStyle)
	text := "GV"
	asciiArt, _ := ToASCIIArt(text, combinedStyle)
	fmt.Println(asciiArt)
	// Note: This is securely managed by the Gopher Officer, which handles the session and is linked to the `processInput` function.
	// Additionally, the Gopher Officer may occasionally sleep during the session's lifecycle and will wake up when needed.
	defer s.cleanup()
	// This Automated Spawn another Goroutine Officer (Known as Gopher Officer) to handle signal
	s.setupSignalHandling()

	// Simulate AI starting the conversation by Gopher Nerd
	// This is a prompt context as the starting point for AI to start the conversation
	humanTyping := NewTypingPrinter()
	PrintPrefixWithTimeStamp(AiNerd, "")
	humanTyping.Print(ContextPrompt, TypingDelay)
	printnewlineASCII() // Ensure there's a newline after the AI's initial message

	// Add AI's initial message to chat history
	s.ChatHistory.AddMessage(AiNerd, ContextPrompt, s.ChatConfig)

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
	PrintPrefixWithTimeStamp(YouNerd, "")
	userInput, err := bufio.NewReader(os.Stdin).ReadString(byte(nl.NewLineChars))
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

// handleUserInput processes the user's input. If the input is a command, it is handled
// accordingly. Otherwise, the input is sent to the AI for a response. It returns true
// if the session should end.
func (s *Session) handleUserInput(input string) bool {
	if !s.ensureClientIsValid() {
		return true // End the session if the client is not valid
	}

	s.ChatHistory.AddMessage(YouNerd, input, s.ChatConfig) // Add the user's input to the chat history

	if success := s.sendInputToAI(input); !success {
		s.endSession() // Ensure the session ends with cleanup.
		return true    // End the session if sending input to AI failed
	}

	return false // Continue the session
}

// ensureClientIsValid checks the validity of the current client and renews it if necessary.
// It returns true if the client is valid or has been successfully renewed, otherwise false.
func (s *Session) ensureClientIsValid() bool {
	if s.Client != nil {
		return true // Client is valid, no action needed
	}
	// Attempt to renew the session if the client is not initialized.
	if err := s.RenewSession(apiKey); err != nil {
		logger.Error(ErrorFailedToRenewSession, err)
		return false // Client is not valid and renewal failed
	}
	return true // Client was successfully renewed
}

// sendInputToAI sends the user input to the AI and updates the chat history with the AI's response.
// It returns true if the input was successfully sent and the response was received, otherwise false.
func (s *Session) sendInputToAI(input string) bool {
	// Define a retryable operation for sending input to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Fix Duplicated by using Magic "_" Identifier
			// Send the input message to the AI, discarding the response.
			_, err := s.SendMessage(s.Ctx, s.Client, input)
			// If there's an error, the operation is not successful.
			return err == nil, err
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)

	if err != nil || !success {
		logger.Error(ErrorSendingMessage, err)
		return false // Sending input to AI failed.
	}

	return true // Input was successfully sent to AI.
}

// cleanup releases resources used by the session. It cancels the context and closes
// the AI client connection.
func (s *Session) cleanup() {
	s.ChatHistory.cleanup() // Perform Clean
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

// handleGenAIError is a utility function that handles errors by logging them and returning an error value.
//
// Parameters:
//
//	err error: The error to handle.
//
// Returns:
//
//	bool: A boolean indicating whether the error was handled successfully.
//	error: The original error value.
func handleGenAIError(err error) (bool, error) {
	if err != nil {
		logger.Error(ErrorGenAI, err)
		return false, err
	}
	return true, nil
}
