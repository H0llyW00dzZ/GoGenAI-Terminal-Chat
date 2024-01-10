// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// isCommand checks if the input is a command based on the prefix.
func isCommand(input string) bool {
	fmt.Println() // Add newline if it's a command or unrecognized command
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
		}
	}

	// If the command is not recognized, inform the AI about the unrecognized command (Free Error Messages hahaha).
	// Note: This cheap since Google AI's Gemini-Pro model, the maximum is 32K tokens
	aiPrompt := fmt.Sprintf(ErrorUserAttemptUnrecognizedCommandPrompt, command)

	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()

	// Send the constructed message to the AI and get the response.
	_, err := SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
	if err != nil {
		errMsg := fmt.Sprintf(ErrorFailedtoSendUnrecognizedCommandToAI, err)
		logger.Error(errMsg)
		return false, fmt.Errorf(errMsg)
	}

	// Since the command was handled (even though unrecognized), return true.
	return true, nil
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
	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()

	// Send a shutdown message to the AI including the chat history
	_, err := SendMessage(session.Ctx, session.Client, ContextPromptShutdown, chatHistory)
	if err != nil {
		// If there's an error sending the message, log it
		logger.Error(ErrorGettingShutdownMessage, err)
	}

	// Proceed with shutdown
	fmt.Println(ShutdownMessage)

	// End the session and perform cleanup
	session.endSession()

	// Signal to the main loop that it's time to exit
	return true, nil
}

// handleHelpCommand processes the ":help" command within a chat session. When a user
// inputs the ":help" command, this function constructs a help prompt that includes a list
// of available commands and sends it to the generative AI model for a response.
//
// The function uses the session's current chat history to provide context for the AI's response,
// ensuring that the help message is relevant to the conversation's state. After sending the
// prompt to the AI, the function retrieves and logs the AI's response, which contains
// information on how to use the commands.
//
// Parameters:
//
//	session *Session: A pointer to the current chat session, which contains state information
//	                  such as the chat history and the generative AI client.
//
// Returns:
//
//	bool: A boolean indicating whether the command was handled. It returns true to signal
//	      that indicate that the command was successfully handled.
//	error: An error object that may occur during the sending of the message to the AI. If the
//	       operation is successful, the error is nil.
//
// The function ensures that the session's context and AI client are utilized to communicate
// with the AI model. It also handles any errors that may occur during the message-sending
// process by logging them appropriately.
//
// Note: The function assumes the presence of a HelpCommandPrompt constant that contains the
// format string for the AI's help prompt, as well as constants for the various commands
// (e.g., QuitCommand, VersionCommand, HelpCommand). It also relies on a logger variable
// to log any errors encountered during the operation.
func handleHelpCommand(session *Session) (bool, error) {
	// Define the help prompt to be sent to the AI, including the list of available commands.
	aiPrompt := fmt.Sprintf(HelpCommandPrompt, QuitCommand, VersionCommand, HelpCommand)

	// Get the entire chat history as a string.
	chatHistory := session.ChatHistory.GetHistory()

	// Send the constructed message to the AI and get the response.
	_, err := SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
	if err != nil {
		logger.Error(ErrorSendingMessage, err)
		return false, err
	}

	// return true to indicate the command was handled, but the session should continue since it's safe to do so alongside with RenewSession
	return true, nil
}

// handleK8sCommand would be a handler function for a hypothetical ":k8s" command.
// Note: This command is secret of H0llyW00dzZ (Original Author) would be used to interact with a Kubernetes cluster.
func k8sCommand(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// handleCheckVersionCommand checks if the current version of the software is the latest.
// It updates the aiPrompt with either a confirmation that the current version is up to date
// or with release notes for the latest version available.
//
// Parameters:
//
//	session *Session: The current session containing the chat history and other context.
//
// Returns:
//
//	bool: Returns true to indicate that the command was successfully handled.
//	error: Returns an error if any occurs during version check or message sending.
//
// Note:
// The function returns `true` to indicate that the command was successfully handled
// and the session should continue. This is safe to perform in conjunction with `RenewSession`.
func handleCheckVersionCommand(session *Session) (bool, error) {
	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()
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
	_, err = SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
	if err != nil {
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		return false, err
	}
	// return true to indicate the command was handled, but the session should continue since it's safe to do so alongside with RenewSession
	return true, nil
}
