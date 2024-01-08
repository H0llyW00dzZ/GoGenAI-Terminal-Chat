// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// CommandHandler defines the function signature for handling chat commands.
// Each command handler function must conform to this signature.
type CommandHandler func(session *Session) (bool, error)

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
	fmt.Println() // Add a newline right after the HandleCommand is entered
	// Trim the input and check if it starts with the command prefix.
	trimmedInput := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmedInput, PrefixChar) {
		// If the input doesn't start with the command prefix, it's not a command.
		return false, nil
	}

	// Split the input into command and potential arguments.
	parts := strings.Fields(trimmedInput)
	if len(parts) == 0 {
		logger.Error(UnknownCommand) // Use logger to log the unknown command error
		return true, nil
	}

	// Retrieve the command and check if it exists in the commandHandlers map.
	command := parts[0]
	if handler, exists := commandHandlers[command]; exists {
		// Call the handler function for the command if no extra arguments are provided.
		if len(parts) == 1 {
			return handler(session)
		} else {
			logger.Error(UnknownCommand) // Use logger to log the unknown command error
			return true, nil
		}
	} else {
		logger.Error(UnknownCommand) // Use logger to log the unknown command error
		return true, nil
	}
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
	// Send a message to the AI asking for a shutdown message
	aiShutdownMessage, err := SendMessage(session.Ctx, session.AiChatSession, ContextPromptShutdown)
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

func handleCheckVersionCommand(session *Session) (bool, error) {
	fmt.Println() // Add a newline right after the check version command is entered
	isLatest, latestVersion, err := checkLatestVersion(CurrentVersion)
	if err != nil {
		return false, err
	}

	if isLatest {
		fmt.Println(YouAreusingLatest)
		fmt.Println()
	} else {
		releaseInfo, err := getFullReleaseInfo(latestVersion)
		if err != nil {
			return false, err
		}

		// Define the color pairs and delimiters to keep or remove
		colorPairs := []string{
			DoubleAsterisk, ColorGreen, // Apply green color and remove ** delimiter
		}
		keepDelimiters := map[string]bool{
			DoubleAsterisk: false, // Remove ** delimiter
		}

		// Colorize and format the release information
		newVersionMessage := fmt.Sprintf(ANewVersionIsAvailable, ColorGreen+releaseInfo.TagName+ColorReset)
		releaseNameMessage := fmt.Sprintf(ReleaseName, ColorGreen+releaseInfo.Name+ColorReset)

		// Colorize content that is surrounded by double asterisks or backticks
		colorizedNewVersionMessage := Colorize(newVersionMessage, colorPairs, keepDelimiters)
		colorizedReleaseNameMessage := Colorize(releaseNameMessage, colorPairs, keepDelimiters)
		colorizedChangelogMessage := Colorize(ColorYellow+releaseInfo.Body+ColorReset+StringNewLine, colorPairs, keepDelimiters)

		// Print the colorized and formatted messages, each followed by a new line
		fmt.Println(colorizedNewVersionMessage)
		fmt.Println(colorizedReleaseNameMessage)
		fmt.Println(colorizedChangelogMessage)
	}

	return false, nil
}
