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

// Execute gracefully terminates the chat session. It sends a shutdown message to the AI,
// prints a farewell message to the user, and signals that the session should end. This method
// is the designated handler for the ":quit" command.
//
// Parameters:
//
//	session *Session: The current chat session, which provides context and state for the operation.
//
// Returns:
//
//	bool: Always returns true to indicate that the session should be terminated.
//	error: Returns an error if one occurs during the shutdown message transmission; otherwise, nil.
//
// The method sends a formatted shutdown message to the AI, which includes the entire chat history
// for context. If an error occurs during message transmission, it is logged. The method then prints
// a predefined shutdown message and invokes a session cleanup function.
//
// Note: The function assumes the presence of constants for the shutdown message format (ContextPromptShutdown)
// and a predefined shutdown message (ShutdownMessage). It relies on the session's endSession method to perform
// any necessary cleanup. The method's return value of true indicates to the calling code that the session loop
// should exit and the application should terminate.
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

// Execute processes the ":help" command within a chat session. It constructs a help prompt
// that includes a list of available commands and sends it to the generative AI model for a response.
// The AI's response, which contains information on how to use the commands, is then logged.
//
// This method provides the AI with the session's current chat history for context, ensuring
// the help message is relevant to the state of the conversation. If an error occurs during
// message transmission, it is logged.
//
// The method assumes the presence of a HelpCommandPrompt constant that contains the format
// string for the AI's help prompt, as well as constants for the various commands (e.g.,
// QuitCommand, VersionCommand, HelpCommand).
//
// Parameters:
//
//	session *Session: the current chat session, which contains state information such as the chat history
//	          and the generative AI client.
//
// Returns:
//
//	bool: Indicates whether the command was successfully handled. It returns false to continue the session.
//	error: Any error that occurs during the version check or message sending process.
//
// Note: The method does not add the AI's response to the chat history to avoid potential
// loops in the AI's behavior.
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
	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}

// Execute checks if the current version of the software is the latest and informs the user accordingly.
// If the current version is not the latest, it retrieves and provides release notes for the latest version.
// This method uses the session's chat history for context and sends an appropriate message to the generative
// AI model for a response.
//
// Parameters:
//
//	session *Session: The current session containing the chat history and other relevant context.
//
// Returns:
//
//	bool: Indicates whether the command was successfully handled. It returns false to continue the session.
//	error: Any error that occurs during the version check or message sending process.
//
// Note: This method does not terminate the session. It is designed to be used with `RenewSession` if needed,
// to ensure that the session state is correctly maintained. The method assumes the presence of constants
// for formatting messages to the AI (YouAreUsingLatest and ReleaseNotesPrompt) and relies on external
// functions (CheckLatestVersion and GetFullReleaseInfo) to determine version information and fetch release details.
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
	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}
