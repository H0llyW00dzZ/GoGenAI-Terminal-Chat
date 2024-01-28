// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// CommandHandler defines the function signature for handling chat commands.
// Each command handler function must conform to this signature.
type CommandHandler interface {
	// Note: The list of command handlers here does not use os.Args; instead, it employs advanced idiomatic Go practices. 🤪
	Execute(session *Session, parts []string) (bool, error) // new method
	IsValid(parts []string) bool                            // new method
}

// CommandRegistry is a centralized registry to manage chat commands.
// It maps command names to their corresponding CommandHandler implementations.
// This allows for a scalable and maintainable way to manage chat commands
// and their execution within a chat session.
type CommandRegistry struct {
	commands map[string]CommandHandler // commands holds the association of command names to their handlers.
}

// NewCommandRegistry initializes a new instance of CommandRegistry.
// It creates a CommandRegistry with an empty map ready to register command handlers.
//
// Returns:
//
//	*CommandRegistry: A pointer to a newly created CommandRegistry with initialized command map.
func NewCommandRegistry() *CommandRegistry {
	return &CommandRegistry{
		commands: make(map[string]CommandHandler),
	}
}

// Register adds a new command and its associated handler to the registry.
// If a command with the same name is already registered, it will be overwritten.
//
// Parameters:
//
//	name string: The name of the command to register.
//	cmd  CommandHandler: The handler that will be associated with the command.
func (r *CommandRegistry) Register(name string, cmd CommandHandler) {
	r.commands[name] = cmd
}

// ExecuteCommand looks up and executes a command based on its name.
// It first validates the command arguments using the IsValid method of the command handler.
// If the command is valid, it executes the command using the Execute method.
// If the command name is not registered, it logs an error.
//
// Parameters:
//
//	name    string: The name of the command to execute.
//	session *Session: The current chat session, providing context for the command execution.
//	parts   []string: The arguments passed along with the command.
//
// Returns:
//
//	bool: A boolean indicating if the command execution should terminate the session.
//	error: An error if one occurs during command validation or execution. Returns nil if no error occurs.
//
// Note:
//
//	If the command is unrecognized, it logs an error but does not return it,
//	as the error is already handled within the method.
func (r *CommandRegistry) ExecuteCommand(name string, session *Session, parts []string) (bool, error) {
	if cmd, exists := r.commands[name]; exists {
		// First, validate the command arguments.
		if !cmd.IsValid(parts) {
			// If the command is not valid, log the error and return.
			logger.Error(HumanErrorWhileTypingCommandArgs, parts)
			return false, nil
		}
		// If the command is valid, execute it.
		return cmd.Execute(session, parts)
	}
	// If the command does not exist, log the error and return.
	logger.Error(ErrorUnrecognizedCommand, name)
	return false, nil // Return nil error since it's already handled.
}

// isCommand checks if the input is a command based on the prefix.
func isCommand(input string) bool {
	fmt.Println() // Add newline if it's a command or unrecognized command
	return strings.HasPrefix(input, PrefixChar)
}

// handleCommand processes the input as a command and returns true if the session should end.
func (s *Session) handleCommand(input string) bool {
	handled, err := HandleCommand(input, s)
	if err != nil {
		logger.Error(ErrorUnknown, err)
	}
	return handled
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
		// Note: this low-level error and should be not happen, but just in case
		return true, fmt.Errorf(ErrorLowLevelCommand)
	}

	// Validate the command arguments.
	commandName := parts[0]
	// Use Magic identifier "_" to ignore the error element, since it duplicates the error handling.
	handled, _ := registry.ExecuteCommand(commandName, session, parts)
	// if err != nil {
	// 	// Since ExecuteCommand already logs errors,
	// 	// keep like this for now, because this palace are low-level error
	// 	return false, err
	// }
	return handled, nil
}

// Note: This list of commands has already been implemented.
// It is now located here for ease of maintenance and to avoid unnecessary complexity.
// This approach questions why many developers write Go code in an overly complex manner (that I don't fucking understand),
// which often leads to problems.
type handleQuitCommand struct{}

// IsValid checks if the quit command is valid based on the input parts.
// The quit command is valid only if there are no additional arguments, hence
// the length of parts must be exactly 1.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handleQuitCommand) IsValid(parts []string) bool {
	// The quit command should not have any arguments.
	return len(parts) == 1
}

type handleHelpCommand struct{}

// IsValid checks if the help command is valid based on the input parts.
// The help command is valid only if there are no additional arguments, hence
// the length of parts must be exactly 1.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handleHelpCommand) IsValid(parts []string) bool {
	// The help command should not have any arguments.
	return len(parts) == 1
}

type handleCheckVersionCommand struct{}

// IsValid checks if the checkversion command is valid based on the input parts.
// The checkversion command is valid only if there are no additional arguments, hence
// the length of parts must be exactly 1.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handleCheckVersionCommand) IsValid(parts []string) bool {
	// The checkversion command should not have any arguments.
	return len(parts) == 1
}

type handleClearCommand struct{}

// IsValid checks if the clear command is valid based on the input parts.
// The clear command is valid only if there are no additional arguments, hence
// the length of parts must be exactly 2.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handleClearCommand) IsValid(parts []string) bool {
	return len(parts) == 2 && parts[0] == ChatCommands && parts[1] == ClearCommand
}

// handleSafetyCommand is the command to adjust safety settings.
type handleSafetyCommand struct {
	SafetySettings *SafetySettings
}

// IsValid checks if the safety command is valid based on the input parts.
func (cmd *handleSafetyCommand) IsValid(parts []string) bool {
	if len(parts) != 2 {
		return false
	}
	option, exists := safetyOptions[parts[1]]
	return exists && option.Valid
}

// setSafetyLevel updates the safety settings based on the command argument.
func (cmd *handleSafetyCommand) setSafetyLevel(level string) {
	option, exists := safetyOptions[level]
	if !exists || !option.Valid {
		// Handle unknown or invalid level, possibly log it or return an error
		logger.Error(ErrorUnknownSafetyLevel, level) // This low-level error won't be shown since it is already handled by the error `IsValid` method
		return
	}

	// Call the setter function associated with the safety level
	option.Setter(cmd.SafetySettings)
}

// handleAITranslateCommand is the command to translate text using the AI model.
type handleAITranslateCommand struct{}

// IsValid checks if the translate command is valid based on the input parts.
// The translate command is expected to follow the pattern: :aitranslate <text> :lang <targetlanguage>
func (cmd *handleAITranslateCommand) IsValid(parts []string) bool {
	// There should be at least 4 parts: the command itself, the text to translate, the language flag, and the target language.
	// Additionally, check for the presence of the language flag ":lang".
	if len(parts) < 4 {
		return false
	}
	languageFlagIndex := len(parts) - 2
	return parts[languageFlagIndex] == LangArgs
}

// handleCryptoRandCommand is the command to translate text using the AI model.
type handleCryptoRandCommand struct{}

// IsValid checks if the cryptorand command is valid based on the input parts.
// The cryptorand command is expected to follow the pattern: :cryptorand :length <number>
func (cmd *handleCryptoRandCommand) IsValid(parts []string) bool {
	// The cryptorand command should have exactly two parts: the command itself and the length argument.
	return len(parts) == 3 && parts[1] == LengthArgs
}

// handleChatShowCommand is responsible for executing the ":show chat history" command.
type handleShowChatCommand struct{}

// IsValid checks if the chat show command is valid based on the input parts.
// The chat show command is valid only if there are no additional arguments, hence
// the length of parts must be exactly 1.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handleShowChatCommand) IsValid(parts []string) bool {
	// Combine the parts after the command keyword to match the ChatHistoryArgs
	args := strings.Join(parts[1:], " ")
	return len(parts) > 1 && args == ChatHistoryArgs
}

// handleSummarizeCommand executes the ":summarize" command.
type handleSummarizeCommand struct{}

// IsValid checks if the summarize command is valid.
func (cmd *handleSummarizeCommand) IsValid(parts []string) bool {
	// The summarize command should not have any arguments.
	return len(parts) == 1
}

// handleClearAllSystemMessagesCommand handles the ":clear :summarize" command.
type handleClearAllSystemMessagesCommand struct{}

// IsValid checks if the clear system messages command is valid based on the input parts.
// This command is valid only if it exactly matches ":clear :summarize".
func (cmd *handleClearAllSystemMessagesCommand) IsValid(parts []string) bool {
	return len(parts) == 2 && parts[0] == ClearCommand && parts[1] == SummarizeCommands
}

// Note: this unimplemented
// Now even it's unimplemented, it wont detected in deadcode indicate that "unreachable func"
type handleK8sCommand struct{}
type storageCommand struct{}
type savehistorytostorageCommand struct{}
type loadhistoryfromstorageCommand struct{}
type reportshitFunctionthatTooComplexCommand struct{}
type handlepingCommand struct{}

// IsValid checks if the ping command is valid based on the input parts.
// The ping command is valid only if it is followed by exactly one argument, which is the IP address.
//
// parts []string: The slice containing the command and its arguments.
//
// Returns true if the command is valid, otherwise false.
func (cmd *handlepingCommand) IsValid(parts []string) bool {
	// The ping command should have exactly two parts: the command itself and the IP address.
	return len(parts) == 2
}

type translateCommand struct{}

type fixDocsFormattingCommand struct{}

type handleTokecountingCommand struct{}

type handlePromptfileCommand struct{}

type handleKaliDocsCommand struct{}
