// Copyright (c) 2024 H0llyW00dzZ

package terminal

import (
	"fmt"
	"strings"
)

// CommandHandler defines the function signature for handling chat commands.
// Each command handler function must conform to this signature.
type CommandHandler func(session *Session) (bool, error)

// commandHandlers maps command strings to their corresponding handler functions.
// This allows for a scalable and maintainable way to manage chat commands.
var commandHandlers = map[string]CommandHandler{
	QuitCommand: handleQuitCommand,
	//TODO: Will add more commands here, example: :help, :about, :credits, :k8s, syncing AI With Go Routines (Known as Gopher hahaha) etc.
	//Note: In python, I don't think so it's possible hahaahaha, also I am using prefix ":" instead of "/" is respect to git and command line, fuck prefix "/" which is confusing for command line
}

// HandleCommand interprets the user input as a command and executes the associated action.
// It uses a map of command strings to their corresponding handler functions to manage
// different commands and their execution. If the command is recognized, the respective
// handler is called; otherwise, an unknown command message is displayed.
//
// Parameters:
//
//	input     string: The user input to be checked for commands.
//	session *Session: The current chat session for context.
//
// Returns:
//
//	bool: A boolean indicating if the input was a command and was handled.
//	error: An error that may occur while handling the command.
func HandleCommand(input string, session *Session) (bool, error) {
	if strings.HasPrefix(input, PrefixChar) {
		command := strings.TrimSpace(input)
		if handler, exists := commandHandlers[command]; exists {
			// Call the handler function for the command
			return handler(session)
		} else {
			fmt.Println(UnknownCommand)
			return true, nil
		}
	}
	return false, nil
}

// handleQuitCommand gracefully terminates the chat session by sending a shutdown
// message, printing a farewell message, and signaling the main loop to exit.
// It is the handler function for the ":quit" command.
//
// Parameters:
//
//	session *Session: The current chat session for context.
//
// Returns:
//
//	bool: Always returns true to indicate the session should end.
//	error: Returns nil if no error occurs; otherwise, returns an error object.
func handleQuitCommand(session *Session) (bool, error) {
	fmt.Println() // Add a newline right after the quit command is entered

	// Send a message to the AI asking for a shutdown message
	aiShutdownMessage, err := SendMessage(session.Ctx, session.Client, ContextPromptShutdown)
	if err != nil {
		// If there's an error sending the message, log it and continue with shutdown
		logger.Error(ErrorGettingShutdownMessage, err)
	} else {
		// If AI provides a shutdown message, print it
		fmt.Println(aiShutdownMessage)
	}
	// Proceed with shutdown
	fmt.Println(ShutdownMessage)
	session.Cancel() // Cancel the context to cleanup resources
	return true, nil // Signal to the main loop that it's time to exit
}

// handleHelpCommand would be a handler function for a hypothetical ":help" command.
func handleHelpCommand(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// handleK8sCommand would be a handler function for a hypothetical ":k8s" command.
// Note: This command is secret of H0llyW00dzZ (Original Author) would be used to interact with a Kubernetes cluster.
func k8sCommand(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}
