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
type CommandHandler interface {
	Execute(*Session) (bool, error)
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
	trimmedInput := strings.TrimSpace(input)
	if !strings.HasPrefix(trimmedInput, PrefixChar) {
		return false, nil
	}

	parts := strings.Fields(trimmedInput)
	if len(parts) == 0 {
		logger.Error(UnknownCommand)
		return true, nil
	}

	command, exists := commandHandlers[parts[0]]
	if !exists {
		// If the command is not recognized, inform the AI about the unrecognized command (Free Error Messages hahaha).
		// Note: This cheap since Google AI's Gemini-Pro model, the maximum is 32K tokens
		aiPrompt := fmt.Sprintf(ErrorUserAttemptUnrecognizedCommandPrompt, ApplicationName, parts[0])

		// Get the entire chat history as a string
		chatHistory := session.ChatHistory.GetHistory()

		// Send the constructed message to the AI and get the response.
		_, err := SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
		if err != nil {
			errMsg := fmt.Sprintf(ErrorFailedtoSendUnrecognizedCommandToAI, err)
			logger.Error(errMsg)
			return false, fmt.Errorf(errMsg)
		}
		return false, nil
	}

	return command.Execute(session)
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
func (q *handleQuitCommand) Execute(session *Session) (bool, error) {
	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()

	// Send a shutdown message to the AI including the chat history
	// this method better instead of hardcode LOL
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, ApplicationName, QuitCommand)
	_, err := SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
	if err != nil {
		// If there's an error sending the message, log it
		logger.Error(ErrorGettingShutdownMessage, err)
	}

	// Proceed with shutdown
	fmt.Println(ShutdownMessage)

	// End the session and perform cleanup
	session.endSession()

	// Signal to the main loop that it's time to exit
	return true, nil // Return true to end the session.
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
func (h *handleHelpCommand) Execute(session *Session) (bool, error) {
	// Define the help prompt to be sent to the AI, including the list of available commands.
	aiPrompt := fmt.Sprintf(HelpCommandPrompt, ApplicationName, QuitCommand, VersionCommand, HelpCommand)

	// Get the entire chat history as a string.
	chatHistory := session.ChatHistory.GetHistory()

	// Send the constructed message to the AI and get the response.
	_, err := SendMessage(session.Ctx, session.Client, aiPrompt, chatHistory)
	if err != nil {
		logger.Error(ErrorSendingMessage, err)
		return false, err
	}
	// Let gopher Remove last chatHistory of the message from the chat history
	// This for protect loop of the AI hahah
	//session.ChatHistory.RemoveMessages(2, "") // Remove 2 line Ai and user message
	// return false to indicate the command was handled, now it doesn't looping ai hahaha
	return false, nil
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
func (c *handleCheckVersionCommand) Execute(session *Session) (bool, error) {
	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()
	// Check if the current version is the latest.
	isLatest, latestVersion, err := CheckLatestVersion(CurrentVersion)
	if err != nil {
		return false, err
	}

	if isLatest {
		aiPrompt = fmt.Sprintf(YouAreusingLatest, ApplicationName, CurrentVersion)
	} else {
		// Fetch the release information for the latest version.
		releaseInfo, err := GetFullReleaseInfo(latestVersion)
		if err != nil {
			return false, err
		}

		aiPrompt = fmt.Sprintf(ReleaseNotesPrompt, ApplicationName,
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
	// Let gopher Remove last chatHistory of the message from the chat history
	// This for protect loop of the AI hahah
	//session.ChatHistory.RemoveMessages(2, "") // Remove 2 line Ai and user message
	// return false to indicate the command was handled now it doesn't looping ai hahaha
	return false, nil
}
