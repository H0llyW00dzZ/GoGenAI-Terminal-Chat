// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"

	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal/fun_stuff"
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
	// Note: The list of command handlers here does not use os.Args; instead, it employs advanced idiomatic Go practices. 🤪
	Execute(session *Session, parts []string) (bool, error) // new method
	IsValid(parts []string) bool                            // new method
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

	// Retrieve the command handler and check if it exists.
	commandHandler, exists := commandHandlers[parts[0]]
	if !exists {
		// Handle unrecognized commands.
		return handleUnrecognizedCommand(parts[0], session, parts)
	}

	// Validate the command arguments.
	if !commandHandler.IsValid(parts) {
		logger.Debug(DEBUGEXECUTINGCMD, parts[0], parts)
		logger.Error(HumanErrorWhileTypingCommandArgs)
		fmt.Println()
		return true, nil
	}

	// Execute the command if it is recognized and valid.
	return commandHandler.Execute(session, parts)
}

// handleUnrecognizedCommand takes an unrecognized command and the current session,
// constructs a prompt to inform the AI about the unrecognized command, and sends
// this information to the AI service. This function is typically called when a user
// input is detected as a command but does not match any of the known command handlers.
//
// Parameters:
// - command string: The unrecognized command input by the user.
// - session *Session: The current chat session containing state and context, including the AI client.
//
// Returns:
// - bool: Always returns false as this function does not result in a command execution.
// - error: Returns an error if sending the message to the AI fails; otherwise, nil.
//
// The function constructs an error prompt using the application's name and the unrecognized command,
// retrieves the current chat history, and sends this information to the AI service. If an error occurs
// while sending the message, the function logs the error and returns an error to the caller.
func handleUnrecognizedCommand(command string, session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, command, parts)
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	// If the command is not recognized, inform the AI about the unrecognized command.
	aiPrompt := fmt.Sprintf(ErrorUserAttemptUnrecognizedCommandPrompt, ApplicationName, command)
	chatHistory := session.ChatHistory.GetHistory()
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Send the constructed message to the AI and get the response.
	_, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, chatHistory)
	if err != nil {
		errMsg := fmt.Sprintf(ErrorFailedtoSendUnrecognizedCommandToAI, err)
		logger.Error(errMsg)
		return false, fmt.Errorf(errMsg)
	}
	return false, nil
}

// Execute gracefully terminates the chat session. It sends a shutdown message to the AI,
// prints a farewell message to the user, and signals that the session should end. This method
// is the designated handler for the ":quit" command.
//
// Parameters:
//
//	session *Session: The current chat session, which provides context and state for the operation.
//	parts []string: The slice containing the command and its arguments.
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
func (q *handleQuitCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, QuitCommand, parts)
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	// Get the entire chat history as a string
	chatHistory := session.ChatHistory.GetHistory()
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(QuitCommand)

	// Send a shutdown message to the AI including the chat history
	// this method better instead of hardcode LOL
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, sanitizedMessage, ApplicationName)
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
//	parts 	[]string: The slice containing the command and its arguments.
//
// Returns:
//
//	bool: Indicates whether the command was successfully handled. It returns false to continue the session.
//	error: Any error that occurs during the version check or message sending process.
//
// Note: The method does not add the AI's response to the chat history to avoid potential
// loops in the AI's behavior.
func (h *handleHelpCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, HelpCommand, parts)
	// Pass ContextPrompt 🤪
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	// Define the help prompt to be sent to the AI, including the list of available commands.
	aiPrompt := fmt.Sprintf(HelpCommandPrompt,
		ApplicationName,
		HelpCommand,
		ShortHelpCommand,
		QuitCommand,
		ShortQuitCommand,
		VersionCommand,
		HelpCommand,
		ShortHelpCommand,
		ClearCommand,
		ClearChatHistoryArgs)
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Get the entire chat history as a string.
	chatHistory := session.ChatHistory.GetHistory()

	// Send the constructed message to the AI and get the response.
	_, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, chatHistory)
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
//	parts 	[]string: The slice containing the command and its arguments.
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
func (c *handleCheckVersionCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug:
	logger.Debug(DEBUGEXECUTINGCMD, VersionCommand, parts)
	// Pass ContextPrompt 🤪
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	// Get the entire chat history as a string
	// adding back this hahaha
	chatHistory := session.ChatHistory.GetHistory()

	// Check if the current version is the latest.
	isLatest, latestVersion, err := CheckLatestVersion(CurrentVersion)
	if err != nil {
		return false, err
	}
	// this avoid duplicate logic, by adding variable in this function
	// unlike quit or help command which is already functionality used variable in `init.go`
	var aiPrompt string
	if isLatest {
		aiPrompt = fmt.Sprintf(YouAreusingLatest, VersionCommand, CurrentVersion, ApplicationName)
	} else {
		// Fetch the release information for the latest version.
		releaseInfo, err := GetFullReleaseInfo(latestVersion)
		if err != nil {
			return false, err
		}
		aiPrompt = fmt.Sprintf(ReleaseNotesPrompt,
			VersionCommand,
			CurrentVersion,
			ApplicationName,
			releaseInfo.TagName,
			releaseInfo.Name,
			releaseInfo.Body)
	}

	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Send the sanitized message to the AI and get the response.
	_, err = SendMessage(session.Ctx, session.Client, sanitizedMessage, chatHistory)
	if err != nil {
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		return false, err
	}
	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}

// Execute performs the ping operation on the provided IP address.
// It uses system utilities to send ICMP packets to the IP and returns the result.
//
// session *Session: The current chat session containing state and context, including the AI client.
// parts   []string: The slice containing the command and its arguments.
//
// Returns true if the ping command was executed, and an error if there was an issue executing the command.
func (cmd *handlepingCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, PingCommand, parts)
	// Note: WIP
	// Validate the command arguments.
	if !cmd.IsValid(parts) {
		logger.Error(HumanErrorWhileTypingCommandArgs)
		fmt.Println()
		return true, nil
	}

	ip := parts[1]
	_, err := fun_stuff.PingIP(ip)
	if err != nil {
		logger.Error(ErrorPingFailed, err)
		fmt.Println()
	}

	return false, nil
}

// Execute clears the chat history if the command is valid.
//
// session *Session: The current chat session containing state and context.
// parts   []string: The slice containing the command and its arguments.
//
// Returns true if the clear command was executed, and an error if there was an issue executing the command.
func (cmd *handleClearCommand) Execute(session *Session, parts []string) (bool, error) {
	// Heads-Up: The current implementation is sleek and storage-agnostic, but beware of the ever-lurking feature creep!
	// Future enhancements might include targeted message purges—think selective user word-bombs or a full-on message-specific snipe hunt.
	// But let's cross that bridge when we get to it. For now, we revel in the simplicity of our logic. Stay tuned, fellow code whisperers! 😜

	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, ClearCommand, parts)
	if cmd.IsValid(parts) {
		session.ChatHistory.Clear()
		PrintPrefixWithTimeStamp(SYSTEMPREFIX)
		PrintTypingChat(ChatHistoryClear, TypingDelay)
		// Added back the context prompt after clearing the chat history
		session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
		fmt.Println()
		return false, nil
	} else {
		// Log the error using the logger instead of returning fmt.Errorf
		errorMessage := HumanErrorWhileTypingCommandArgs
		logger.Error(errorMessage)
		// Return nil for the error since we've already logged it
		return false, nil
	}
}

// Execute processes the ":safety" command within a chat session.
func (cmd *handleSafetyCommand) Execute(session *Session, parts []string) (bool, error) {
	if !cmd.IsValid(parts) {
		logger.Error(HumanErrorWhileTypingCommandArgs)
		fmt.Println()
		return true, nil
	}

	// Ensure SafetySettings is initialized.
	if cmd.SafetySettings == nil {
		cmd.SafetySettings = DefaultSafetySettings()
		// Pass ContextPrompt just incase
		session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	}

	// Set the safety level based on the command argument.
	cmd.setSafetyLevel(parts[1])

	// Apply the updated safety settings and notify the user.
	cmd.SafetySettings.ApplyToModel(session.Client.GenerativeModel(ModelAi))
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt)
	PrintPrefixWithTimeStamp(SYSTEMPREFIX)
	PrintTypingChat(fmt.Sprintf(SystemSafety, parts[1]), TypingDelay)
	fmt.Println()

	return false, nil // Continue the session after setting safety levels
}
