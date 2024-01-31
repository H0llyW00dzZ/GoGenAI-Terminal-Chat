// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// addMessageWithContext adds a message to the chat history with context.
func addMessageWithContext(session *Session, sender, message string) {
	session.ChatHistory.AddMessage(sender, message, session.ChatConfig)
}

// executeCommand is a generic function to execute a command.
func executeCommand(session *Session, command string, constructPrompt func(string) string) (bool, error) {
	success, err := sendCommandToAI(session, command, constructPrompt)
	if err != nil {
		logger.Error(ErrorFailedToSendCommandToAI, err)
		return false, err
	}
	return !success, nil // Return false to continue the session if successful.
}

// sendCommandToAI sends a command to the AI after sanitizing and applying retry logic.
func sendCommandToAI(session *Session, command string, constructPrompt func(string) string) (bool, error) {
	sanitizedCommand := session.ChatHistory.SanitizeMessage(command)
	aiPrompt := constructPrompt(sanitizedCommand)

	apiErrorHandler := func(err error) bool {
		return strings.Contains(err.Error(), Error500GoogleApi)
	}

	return retryWithExponentialBackoff(func() (bool, error) {
		return sendMessageToAI(session, aiPrompt)
	}, apiErrorHandler)
}

// sendMessageToAI sends a message to the AI and handles the response.
func sendMessageToAI(session *Session, message string) (bool, error) {
	_, err := SendMessage(session.Ctx, session.Client, message, session)
	return err == nil, err
}

// sendShutdownMessage sends a formatted shutdown message to the AI and logs it to the chat history.
func sendShutdownMessage(session *Session) error {
	session.ChatHistory.Clear()
	session.ChatHistory.AddMessage(AiNerd,
		ContextPrompt,
		session.ChatConfig)
	// Assuming QuitCommand is the user input that triggered the shutdown.
	addMessageWithContext(session, StringNewLine+YouNerd, QuitCommand)

	// Sanitize the message before sending it to the AI.
	sanitizedMessage := session.ChatHistory.SanitizeMessage(QuitCommand)

	// Send a shutdown message to the AI including the chat history with the context prompt
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, sanitizedMessage, ApplicationName)

	// Retry logic for sending the shutdown message to the AI.
	apiErrorHandler := func(err error) bool {
		// Error 500 Google Api
		return strings.Contains(err.Error(), Error500GoogleApi)
	}

	// Attempt to send the shutdown message to the AI with retry logic
	_, err := retryWithExponentialBackoff(func() (bool, error) {
		return sendMessageToAI(session, aiPrompt)
	}, apiErrorHandler)

	return err
}
