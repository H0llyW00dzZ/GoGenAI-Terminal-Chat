// Copyright (c) 2024 H0llyW00dzZ
//
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
	if err := sendShutdownMessage(session); err != nil {
		logger.Error(ErrorFailedToSendShutdownMessage, err)
	}
	// Proceed with shutdown regardless of the error
	fmt.Println(ShutdownMessage)
	session.endSession() // End the session and perform cleanup
	return true, nil
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
func (cmd *handleHelpCommand) Execute(session *Session, parts []string) (bool, error) {
	return executeCommand(session, HelpCommand, func(cmd string) string {
		// Note: This a better fmt formatting unlike 'C' or 'RUST' hahahaha
		return fmt.Sprintf(HelpCommandPrompt,
			ApplicationName,
			cmd,
			QuitCommand,
			ShortQuitCommand,
			HelpCommand,
			ShortHelpCommand,
			VersionCommand,
			SafetyCommand,
			Low, Default, High, Unspecified, None,
			AITranslateCommand,
			LangArgs,
			CryptoRandCommand,
			LengthArgs,
			SummarizeCommands,
			SwitchModelCommands,
			GeminiPro, GeminiProTuning, GeminiProLatest,
			ChatCommands,
			ShowCommands,
			ChatHistoryArgs,
			StatsCommand,
			ChatCommands,
			ClearCommand,
			SummarizeCommands,
			ClearCommand,
			ChatCommands,
			CheckModelCommands,
			TokenCountCommands,
			FileCommands)
	})
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
	// Pass ContextPrompt 🤪
	// Add messages to the chat history to provide context for the version check.
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	session.ChatHistory.AddMessage(YouNerd, VersionCommand, session.ChatConfig)

	// Check if the current version is the latest and get the prompt for the AI.
	aiPrompt, err := c.checkVersionAndGetPrompt()
	if err != nil {
		// Log errors related to version checking and Google API issues.
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		logger.HandleGoogleAPIError(err)
		return false, err
	}

	// Sanitize the AI prompt to ensure it is safe to send.
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Define a retryable operation with a function that sends the version check message to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Fix Duplicated by using Magic "_" Identifier
			// Send the message to the AI, discarding the response since it's not needed here.
			_, err := session.SendMessage(session.Ctx, session.Client, sanitizedMessage)
			// If there's no error, the operation is successful.
			return err == nil, err
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)

	if err != nil {
		// If an error occurs that is not recoverable by retries, log it and return the error.
		logger.Error(ErrorFailedToSendVersionCheckMessage, err)
		return false, err
	}

	if !success {
		// If the operation was not successful after retries, return an error.
		return false, fmt.Errorf(ErrorFailedToSendVersionCheckMessageAfterReties)
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
func (cmd *handlepingCommand) Execute(session *Session, parts []string, subcommand string) (bool, error) {
	// Note: WIP
	// Validate the command arguments.
	if !cmd.IsValid(parts) {
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return true, nil
	}

	ip := parts[1]
	_, err := fun_stuff.PingIP(ip)
	if err != nil {
		logger.Error(ErrorPingFailed, err)
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

	// Note: This place only, for commands doesn't have any subcommands/args, so it will return error hahaha
	return cmd.HandleSubcommand("", session, parts) // Continue the session
}

func (cmd *handleClearCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	// Handle the subcommands of the clear command
	switch subcommand {
	case ChatCommands:
		return cmd.clearChatHistory(session)
	case SummarizeCommands:
		return cmd.clearSummarizeHistory(session)
	default:
		// Handle unrecognized subcommand
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}
}

// Execute processes the ":safety" command within a chat session.
//
// Note: The flexibility demonstrated in this function is quite powerful. In many programming languages,
// changing safety settings would typically require constructing and parsing JSON structures for each request.
// However, Go's type system allows us to elegantly manipulate these settings directly through struct methods,
// bypassing the need for repetitive JSON serialization and deserialization hahaha.
func (cmd *handleSafetyCommand) Execute(session *Session, parts []string) (bool, error) {
	// Continue the session after setting safety levels
	return cmd.HandleSubcommand("", session, parts) // Continue the session
}

func (cmd *handleSafetyCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	// Note: The code in "safety_settings.go" employs advanced idiomatic Go practices. 🤪
	// Caution is advised: if you're not familiar with these practices, improper handling in this "Execute" could lead to frequent panics 24/7 🤪.
	if !cmd.IsValid(parts) {
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}

	// Set the safety level based on the command argument.
	cmd.setSafetyLevel(session, parts[1])

	// Apply the updated safety settings and notify the user.
	// Note: It should be working now. If it still doesn't work, this may indicate a problem with your machine hahaha.
	session.SafetySettings.ApplyToModel(session.Client.GenerativeModel(GeminiPro), GeminiPro)
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	logger.Any(fmt.Sprintf(SystemSafety, parts[1])) // simplify
	return false, nil
}

// Execute processes the ":aitranslate" command within a chat session.
func (cmd *handleAITranslateCommand) Execute(session *Session, parts []string) (bool, error) {
	// Ensure that the command is valid before proceeding.
	if !cmd.IsValid(parts) {
		// Log the error with the logger instead of returning fmt.Errorf
		logger.Error(ErrorWhileTypingCommandArgs, AITranslateCommand, parts)
		return false, nil // Return nil error because the logger already handled it
	}

	// Find the index of the language flag ":lang" to separate text and target language.
	languageFlagIndex := len(parts) - 2
	textToTranslate := strings.Join(parts[1:languageFlagIndex], " ")
	targetLanguage := parts[languageFlagIndex+1]

	aiPrompt := constructAITranslatePrompt(ApplicationName, AITranslateCommand, textToTranslate, targetLanguage)

	err := handleAIInteraction(session, aiPrompt, func(session *Session, aiResponse string) error {
		// Add a message to the chat history indicating the translation command was invoked
		translationCommandMessage := fmt.Sprintf(ContextUserInvokeTranslateCommands, targetLanguage, textToTranslate)
		session.ChatHistory.AddMessage(YouNerd, translationCommandMessage, session.ChatConfig)
		return postProcessAITranslate(session, aiResponse)
	})

	if err != nil {
		logger.Error(ErrorFailedToSendTranslationMessage, err)
		return false, err
	}

	// Indicate that the command was handled; return false to continue the session.
	return false, nil
}

// Execute processes the ":cryptorand" command within a chat session.
func (cmd *handleCryptoRandCommand) Execute(session *Session, parts []string) (bool, error) {
	// Continue the session without performing any action.
	return cmd.HandleSubcommand("", session, parts) // Continue the session
}

func (cmd *handleCryptoRandCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	// Check if there are enough parts to contain the length argument.
	if len(parts) < 3 {
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}

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
	logger.Any(CryptoRandRes, lengthStr, randomString)
	return false, nil
}

// Execute displays the entire chat history.
//
// session *Session: The current chat session containing state and context.
// parts   []string: The slice containing the command and its arguments.
//
// Returns false to indicate the session should continue, and an error if there is an issue.
func (cmd *handleShowChatCommand) Execute(session *Session, parts []string) (bool, error) {
	// Return false to indicate the session should continue.
	return cmd.HandleSubcommand("", session, parts) // Continue the session
}

func (cmd *handleShowChatCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {

	// Check if there are enough parts to contain the length argument.
	if len(parts) < 3 {
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}

	// Retrieve and log the entire chat history.
	history := session.ChatHistory.GetHistory(session.ChatConfig)
	logger.Info(ShowChatHistory, history)
	return false, nil // Return false to indicate the session should continue.
}

// Execute processes the ":summarize" command within a chat session.
func (h *handleSummarizeCommand) Execute(session *Session, parts []string) (bool, error) {
	// Add a message to the chat history indicating the summarize command was invoked
	session.ChatHistory.AddMessage(YouNerd, SummarizeCommands, session.ChatConfig)
	// Check if there are system messages in the chat history before summarizing.
	if session.ChatHistory.HasSystemMessages() {
		// Remove system messages from the chat history.
		session.ChatHistory.ClearAllSystemMessages()
	}

	// Define the summarize prompt to be sent to the AI.
	aiPrompt := h.constructSummarizePrompt()
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	success, err := h.sendSummarizePrompt(session, sanitizedMessage)
	if err != nil {
		logger.Error(ErrorFailedToSendSummarizeMessage, err)
		return false, err
	}

	if !success {
		return false, fmt.Errorf(ErrorFailedToSendSummarizeMessageAfterRetries)
	}

	// Indicate that the command was handled successfully; return false to continue the session.
	return false, nil
}

// Execute processes the main command for handleStatsCommand. Since handleStatsCommand
// is implemented with subcommands, this method does not perform any action and simply
// returns false and nil to indicate that the session should continue without error.
// The actual command logic is delegated to the HandleSubcommand method.
func (cmd *handleStatsCommand) Execute(session *Session, parts []string) (bool, error) {
	// Continue the session without performing any action.
	return cmd.HandleSubcommand("", session, parts) // Continue the session
}

// HandleSubcommand dispatches the handling of specific subcommands for the stats command.
// It takes a subcommand string, the current session, and the command parts as arguments.
// Based on the subcommand, it calls the appropriate method to handle it.
// If the subcommand is not recognized, it logs an error and continues the session.
func (cmd *handleStatsCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	// Dispatch handling based on the subcommand.
	switch subcommand {
	case ChatCommands:
		// Handle the ':chat' subcommand to show chat statistics.
		return cmd.showChatStats(session)
	default:
		// Log an error for unrecognized subcommands and continue the session.
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}
}

func (cmd *handleTokeCountingCommand) Execute(session *Session, parts []string) (bool, error) {
	// Continue the session
	return cmd.HandleSubcommand("", session, parts)
}

func (cmd *handleTokeCountingCommand) HandleSubcommand(subcommand string, session *Session, parts []string) (bool, error) {
	if !cmd.IsValid(parts) {
		logger.Error(ErrorWhileTypingCommandArgs, subcommand, parts)
		return false, nil
	}

	// The file paths start from index 2
	filePaths := parts[2:]

	apiKey := os.Getenv(APIKey) // Retrieve the API_KEY from the environment
	switch subcommand {
	case FileCommands:
		return cmd.handleTokenCount(apiKey, filePaths)
	default:
		// Log an error for unrecognized subcommands and continue the session.
		logger.Error(ErrorUnrecognizedSubcommandForTokenCount, subcommand)
		return false, nil
	}

}

func (cmd *handleCheckModelCommand) Execute(session *Session, parts []string) (bool, error) {
	// Validate the command arguments.
	if !cmd.IsValid(parts) {
		logger.Error(ErrorWhileTypingCommandArgs, CheckModelCommands, parts)
		return false, nil
	}

	modelName := parts[1] // The model name is the second part.

	// Define a retryable operation for retrieving model info.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Note: this a magic method
			modelInfo, err := session.Client.EmbeddingModel(modelName).Info(session.Ctx)
			if err != nil {
				// Log the error and decide if it's worth retrying based on the error type.
				logger.Error(ErrorFailedToRetriveModelInfo, err)
				return false, err
			}

			// Process and display the model information if retrieval was successful.
			DisplayModelInfo(modelInfo)
			return true, nil
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)

	if err != nil || !success {
		logger.Error("%s", err)
	}

	// Return false to indicate the session should continue.
	return false, nil
}

// isValidModelName checks if the provided model name is supported.
func isValidModelName(modelName string) (bool, error) {
	if valid, exists := supportedModels[modelName]; exists && valid {
		return true, nil
	}
	// If the model name is not found or not valid, return an error.
	return false, fmt.Errorf(ErrorUnsupportedModelName, modelName)
}

// Execute changes the current AI model used in the session to the one specified in the command.
func (cmd *handleSwitchModelCommand) Execute(session *Session, parts []string) (bool, error) {
	// Check if the command is valid.
	if !cmd.IsValid(parts) {
		logger.Error(ErrorWhileTypingCommandArgs, SwitchModelCommands, parts)
		return false, nil
	}

	// Extract the model name from the command parts.
	modelName := parts[1]

	// Validate the model name.
	valid, err := isValidModelName(modelName)
	if !valid {
		// Log the error using logger.Error and return the error message.
		logger.Error(err.Error())
		return false, nil // Continue the session
	}

	// Update the session with the new model name.
	session.CurrentModelName = modelName

	// Log a confirmation message using the logger.
	logger.Any(SwitchedModel, modelName)

	return false, nil // Continue the session.
}
