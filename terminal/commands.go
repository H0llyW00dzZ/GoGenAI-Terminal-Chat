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
			fmt.Println() // Add a newline before the command is executed
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
	// Request a shutdown message from the AI but don't print it
	_, err := SendMessage(session.Ctx, session.AiChatSession, ContextPromptShutdown)
	if err != nil {
		// Log the error if there's an issue getting the message
		logger.Error(ErrorGettingShutdownMessage, err)
	}
	// Print only the shutdown message
	fmt.Println() // A better newline instead of hardcoding "\n"
	fmt.Println(StripChars)
	fmt.Println(ShutdownMessage) // Print the shutdown message
	session.EndSession()         // Use EndSession to perform cleanup and signal that the session has ended

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

// handleCheckVersionCommand checks if the user is running the latest version of the application.
// If the user is not on the latest version, it fetches and displays the release information for
// the newer version. This function is intended to be used as a command handler within a chat session
// that allows users to check for software updates.
//
// The function performs the following actions:
//   - Calls CheckLatestVersion with the current application version to determine if an update is available.
//   - If the current version is the latest, it informs the user accordingly.
//   - If there is a newer version, it fetches the full release information and prepares a prompt for the AI
//     to explain the release notes to the user.
//
// Parameters:
//
//	session *Session: A pointer to the current Session object, which contains the chat session state and context.
//
// Returns:
//
//	bool: A boolean flag indicating whether the command has been fully handled (always false as the session should continue after checking the version).
//	error: An error that may occur during the version check or while fetching the release information.
//
// Note: The function sends messages to the AI using SendMessage and assumes that the session and its context
// are properly initialized and active. It does not return any AI-generated messages directly to the user but
// assumes that the AI response is handled elsewhere in the chat session flow.
func handleCheckVersionCommand(session *Session) (bool, error) {
	// Check if the current version is the latest.
	isLatest, latestVersion, err := CheckLatestVersion(CurrentVersion)
	if err != nil {
		return false, err
	}

	if isLatest {
		aiPrompt = fmt.Sprintf(YouAreusingLatest, CurrentVersion)
	} else {
		// Fetch the release information for the latest version.
		releaseInfo, err := GetFullReleaseInfo(latestVersion)
		if err != nil {
			return false, err
		}

		aiPrompt = fmt.Sprintf(ReleaseNotesPrompt,
			CurrentVersion,
			releaseInfo.TagName,
			releaseInfo.Name,
			releaseInfo.Body)
	}

	// Send the constructed message to the AI and get the response.
	_, err = SendMessage(session.Ctx, session.AiChatSession, aiPrompt)
	if err != nil {
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		return false, err
	}

	fmt.Println() // A better newline instead of hardcoding "\n"
	fmt.Println(StripChars)

	return false, nil
}
