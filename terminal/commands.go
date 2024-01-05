// Copyright (c) 2024 H0llyW00dzZ

package terminal

import (
	"fmt"
	"os"
	"strings"
)

// HandleCommand checks if the user input is a command and executes it.
// It supports the :quit command to gracefully shut down the application.
//
// Parameters:
//
//	input string: The user input to be checked for commands.
//	session *Session: The current chat session for context.
//
// Returns:
//
//	bool: A boolean indicating if the input was a command.
//	error: An error that may occur while handling the command.
func HandleCommand(input string, session *Session) (bool, error) {
	if strings.HasPrefix(input, PrefixChar) {
		switch input {
		case QuitCommand:
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
			session.Client.Close()
			os.Exit(0) // Exit the application immediately
			return true, nil
		default:
			fmt.Println(UnknownCommand)
			return true, nil
		}
	}
	return false, nil
}
