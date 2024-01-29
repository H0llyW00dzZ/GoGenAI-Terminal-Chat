// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal/fun_stuff"
	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal/tools"
)

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

	// Context
	// Clear the chat history now that the shutdown message has been sent
	session.ChatHistory.Clear()
	session.ChatHistory.AddMessage(StringNewLine+YouNerd,
		QuitCommand,
		session.ChatConfig) // should be accurate now

	// Sanitize the message before sending it to the AI
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
		_, err := SendMessage(session.Ctx, session.Client, aiPrompt, session)
		return err == nil, err
	}, apiErrorHandler)

	if err != nil {
		// If there's an error sending the message, log it
		logger.Error(ErrorFailedToSendShutdownMessage, err)
	}

	// Proceed with shutdown regardless of the error
	fmt.Println(ShutdownMessage)
	session.endSession() // End the session and perform cleanup

	// Signal to the main loop that it's time to exit
	return true, nil // Always return true to end the session, and nil for error since we handle it above.
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
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	session.ChatHistory.AddMessage(StringNewLine+YouNerd, HelpCommand, session.ChatConfig)
	// Define the help prompt to be sent to the AI, including the list of available commands.
	aiPrompt := fmt.Sprintf(HelpCommandPrompt,
		// Note: This doesn't look complex, as the complex one looks way better than the "hardcoded" one LOL
		ApplicationName,
		HelpCommand,
		QuitCommand,
		ShortQuitCommand,
		HelpCommand,
		ShortHelpCommand,
		VersionCommand,
		SafetyCommand,
		Low,
		Default,
		High,
		AITranslateCommand,
		LangArgs,
		CryptoRandCommand,
		LengthArgs,
		SummarizeCommands,
		ShowCommands,
		ChatHistoryArgs,
		ClearCommand,
		SummarizeCommands,
		ChatCommands,
		ClearCommand,
	)

	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Retry logic for sending the help prompt to the AI.
	apiErrorHandler := func(err error) bool {
		// Error 500 Google Api
		return strings.Contains(err.Error(), Error500GoogleApi)
	}

	success, err := retryWithExponentialBackoff(func() (bool, error) {
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		// Sanitize AI's response to remove any separators
		aiResponse = sanitizeAIResponse(aiResponse)
		// Add the sanitized AI's response to the chat history
		session.ChatHistory.AddMessage(AiNerd, aiResponse, session.ChatConfig)
		return err == nil, err
	}, apiErrorHandler)

	if err != nil {
		logger.Error(ErrorFailedToSendHelpMessage, err)
		return false, err
	}

	if !success {
		return false, fmt.Errorf(ErrorFailedToSendHelpMessagesAfterRetries)
	}

	// Indicate that the command was handled successfully; return false to continue the session.
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
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	session.ChatHistory.AddMessage(StringNewLine+YouNerd, VersionCommand, session.ChatConfig)
	// Get the entire chat history as a string

	// Check if the current version is the latest.
	aiPrompt, err := c.checkVersionAndGetPrompt()
	if err != nil {
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		logger.HandleGoogleAPIError(err)
		return false, err
	}

	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Retry logic for sending the version check prompt to the AI.
	apiErrorHandler := func(err error) bool {
		// Error 500 Google Api
		return strings.Contains(err.Error(), Error500GoogleApi)
	}

	success, err := retryWithExponentialBackoff(func() (bool, error) {
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		// Sanitize AI's response to remove any separators
		aiResponse = sanitizeAIResponse(aiResponse)
		// Add the sanitized AI's response to the chat history
		session.ChatHistory.AddMessage(AiNerd, aiResponse, session.ChatConfig)
		return err == nil, err
	}, apiErrorHandler)

	if err != nil {
		logger.Error(ErrorFailedToSendVersionCheckMessage, err)
		return false, err
	}

	if !success {
		return false, fmt.Errorf(ErrorFailedToSendVersionCheckMessageAfterReties)
	}

	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}

// checkVersionAndGetPrompt checks if the current version of the software is the latest and informs the user accordingly.
func (c *handleCheckVersionCommand) checkVersionAndGetPrompt() (aiPrompt string, err error) {
	// Check if the current version is the latest.
	isLatest, latestVersion, err := checkLatestVersionWithBackoff()
	if err != nil {
		return "", err
	}

	if isLatest {
		aiPrompt = fmt.Sprintf(YouAreusingLatest, VersionCommand, ApplicationName, CurrentVersion)
	} else {
		// Fetch and format the release information for the latest version.
		aiPrompt, err = fetchAndFormatReleaseInfo(latestVersion)
		if err != nil {
			return "", err
		}
	}

	// Return the prompt to the caller.
	return aiPrompt, nil
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
		logger.Error(HumanErrorWhileTypingCommandArgs, parts)
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

	// Log the error if the command is not valid
	logger.Error(HumanErrorWhileTypingCommandArgs, parts)
	return false, nil // Continue the session
}

func (cmd *handleClearCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, ClearCommand, parts)

	// Handle the subcommands of the clear command
	switch subcommand {
	case ChatCommands:
		return cmd.clearChatHistory(session)
	case SummarizeCommands:
		return cmd.clearSummarizeHistory(session)
	default:
		// Handle unrecognized subcommand
		logger.Error(ErrorUnrecognizedCommand, subcommand)
		return false, nil
	}
}

// clearChatHistory clears the chat history.
func (cmd *handleClearCommand) clearChatHistory(session *Session) (bool, error) {
	session.ChatHistory.Clear()
	// Prepare the full message to be printed
	clearMessage := ChatHistoryClear
	showTokenCount := os.Getenv(SHOW_TOKEN_COUNT) == "true"
	// Append token reset message if SHOW_TOKEN_COUNT is true
	if showTokenCount {
		totalTokenCount = 0 // Reset the total token count to zero
		clearMessage += "\n" + ResetTotalTokenUsage
	}
	// Print the message(s) with timestamp and typing effect
	PrintPrefixWithTimeStamp(SYSTEMPREFIX)
	PrintTypingChat(clearMessage, TypingDelay)
	// Added back the context prompt after clearing the chat history
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	fmt.Println()
	return false, nil // Continue the session
}

// clearSummarizeHistory clears the summarized messages from the chat history.
func (cmd *handleClearCommand) clearSummarizeHistory(session *Session) (bool, error) {
	session.ChatHistory.ClearAllSystemMessages()
	PrintPrefixWithTimeStamp(SYSTEMPREFIX)
	PrintTypingChat(ChatSysSummaryMessages, TypingDelay)
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	fmt.Println()
	return false, nil // Continue the session
}

// Execute processes the ":safety" command within a chat session.
//
// Note: The flexibility demonstrated in this function is quite powerful. In many programming languages,
// changing safety settings would typically require constructing and parsing JSON structures for each request.
// However, Go's type system allows us to elegantly manipulate these settings directly through struct methods,
// bypassing the need for repetitive JSON serialization and deserialization hahaha.
func (cmd *handleSafetyCommand) Execute(session *Session, parts []string) (bool, error) {
	// Note: The code in "safety_settings.go" employs advanced idiomatic Go practices. 🤪
	// Caution is advised: if you're not familiar with these practices, improper handling in this "Execute" could lead to frequent panics 24/7 🤪.
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, SafetyCommand, parts)
	if !cmd.IsValid(parts) {
		logger.Error(HumanErrorWhileTypingCommandArgs, parts)
		fmt.Println()
		return true, nil
	}

	// Ensure SafetySettings is initialized.
	if cmd.SafetySettings == nil {
		cmd.SafetySettings = DefaultSafetySettings()
		// Pass ContextPrompt just incase
		session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	}

	// Set the safety level based on the command argument.
	cmd.setSafetyLevel(parts[1])

	// Apply the updated safety settings and notify the user.
	cmd.SafetySettings.ApplyToModel(session.Client.GenerativeModel(ModelAi))
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	PrintPrefixWithTimeStamp(SYSTEMPREFIX)
	PrintTypingChat(fmt.Sprintf(SystemSafety, parts[1]), TypingDelay)
	fmt.Println()     // this correct, fix front end issue
	return false, nil // Continue the session after setting safety levels
}

// Execute processes the ":aitranslate" command within a chat session.
func (cmd *handleAITranslateCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, AITranslateCommand, parts)

	// Find the index of the language flag ":lang" to separate text and target language.
	languageFlagIndex := len(parts) - 2
	textToTranslate := strings.Join(parts[1:languageFlagIndex], " ")
	targetLanguage := parts[languageFlagIndex+1]

	// Define the translation prompt to be sent to the AI.
	aiPrompt := fmt.Sprintf(AITranslateCommandPrompt,
		ApplicationName,
		AITranslateCommand,
		textToTranslate,
		targetLanguage)

	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Retry logic for sending the translation prompt to the AI.
	apiErrorHandler := func(err error) bool {
		// Error 500 Google Api
		return strings.Contains(err.Error(), Error500GoogleApi)
	}

	// Wrap the SendMessage call within retryWithExponentialBackoff
	success, err := retryWithExponentialBackoff(func() (bool, error) {
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		if err != nil {
			return false, err
		}
		// Sanitize AI's response to remove any separators
		aiResponse = sanitizeAIResponse(aiResponse)
		// Add a message to the chat history indicating the translation command was invoked
		translationCommandMessage := fmt.Sprintf(ContextUserInvokeTranslateCommands, targetLanguage, textToTranslate)
		session.ChatHistory.AddMessage(StringNewLine+YouNerd, translationCommandMessage, session.ChatConfig)
		// Add the sanitized AI's response to the chat history
		session.ChatHistory.AddMessage(AiNerd, aiResponse, session.ChatConfig)
		return true, nil
	}, apiErrorHandler)

	if err != nil {
		logger.Error(ErrorFailedToSendTranslationMessage, err)
		return false, err
	}

	if !success {
		return false, fmt.Errorf(ErrorFailedToSendTranslationMessageAfterRetries)
	}

	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}

// Execute processes the ":cryptorand" command within a chat session.
func (cmd *handleCryptoRandCommand) Execute(session *Session, parts []string) (bool, error) {
	lengthStr := parts[2] // The length argument is now the second part of the command
	length, err := strconv.Atoi(lengthStr)
	if err != nil {
		logger.Error(ErrorInvalidLengthArgs, err)
		return false, fmt.Errorf(errorinvalidlengthArgs, err)
	}

	randomString, err := tools.GenerateRandomString(length)
	if err != nil {
		logger.Error(ErrorFailedtoGenerateRandomString, err)
		return false, fmt.Errorf(errorfailedtogeneraterandomstring, err)
	}

	logger.Debug(CryptoRandLength, lengthStr)
	logger.Debug(CryptoRandStringRes, randomString)

	logger.Info(CryptoRandRes, lengthStr, randomString)

	return false, nil
}

// Execute displays the entire chat history.
//
// session *Session: The current chat session containing state and context.
// parts   []string: The slice containing the command and its arguments.
//
// Returns false to indicate the session should continue, and an error if there is an issue.
func (cmd *handleShowChatCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, ShowCommands, parts)
	if !cmd.IsValid(parts) {
		logger.Error(HumanErrorWhileTypingCommandArgs, parts)
		fmt.Println()
		return true, nil
	}

	// Retrieve and log the entire chat history.
	history := session.ChatHistory.GetHistory(session.ChatConfig)
	logger.Info(ShowChatHistory, history)

	return false, nil // Return false to indicate the session should continue.
}

// Execute processes the ":summarize" command within a chat session.
func (h *handleSummarizeCommand) Execute(session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, SummarizeCommands, parts)
	// Define the summarize prompt to be sent to the AI.
	aiPrompt := h.constructSummarizePrompt()
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	success, err := h.sendSummarizePrompt(session, sanitizedMessage)
	if err != nil {
		logger.Error(ErrorFailedToSendHelpMessage, err)
		return false, err
	}

	if !success {
		return false, fmt.Errorf(ErrorFailedToSendHelpMessagesAfterRetries)
	}

	// Indicate that the command was handled successfully; return false to continue the session.
	return false, nil
}

// constructSummarizePrompt constructs the prompt to be sent to the AI for summarization.
func (h *handleSummarizeCommand) constructSummarizePrompt() string {
	return fmt.Sprintf(SummarizePrompt)
}

// sendSummarizePrompt sends the summarize prompt to the AI and handles the response.
func (h *handleSummarizeCommand) sendSummarizePrompt(session *Session, sanitizedMessage string) (bool, error) {
	apiErrorHandler := func(err error) bool {
		// Error 500 Google Api
		return strings.Contains(err.Error(), Error500GoogleApi)
	}
	// Retry logic for sending the summarize prompt to the AI.
	return retryWithExponentialBackoff(func() (bool, error) {
		// Note: This is subject to change, for example,
		// to implement another functionality without displaying AI response in the terminal,
		// but only adding it to the chat history.
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		if err != nil {
			return false, err
		}

		h.handleAIResponse(session, sanitizedMessage, aiResponse)
		return true, nil
	}, apiErrorHandler)
}

// handleAIResponse processes the AI's response to the summarize command.
func (h *handleSummarizeCommand) handleAIResponse(session *Session, sanitizedMessage, aiResponse string) {
	// Instead of directly adding, check if a system message already exists and replace it.
	formattedResponse := fmt.Sprintf(ObjectHighLevelString, SYSTEMPREFIX, aiResponse)
	if !session.ChatHistory.handleSystemMessage(sanitizedMessage, formattedResponse, session.ChatHistory.hashMessage(aiResponse)) {
		// If it was not a system message or no existing system message was found to replace,
		// add the new system message to the chat history.
		session.ChatHistory.AddMessage(SYSTEMPREFIX, aiResponse, session.ChatConfig)
	}
}
