// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import "strings"

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
	// Combine the parts after the command keyword to match the ClearChatHistoryArgs
	args := strings.Join(parts[1:], " ")
	return len(parts) > 1 && args == ClearChatHistoryArgs
}

// handleSafetyCommand is the command to adjust safety settings.
type handleSafetyCommand struct {
	SafetySettings *SafetySettings
}

// IsValid checks if the safety command is valid based on the input parts.
func (cmd *handleSafetyCommand) IsValid(parts []string) bool {
	return len(parts) == 2 && (parts[1] == Low || parts[1] == High || parts[1] == Default)
}

// setSafetyLevel updates the safety settings based on the command argument.
func (cmd *handleSafetyCommand) setSafetyLevel(level string) {
	switch level {
	case Low:
		cmd.SafetySettings.SetLowSafety()
	case High:
		cmd.SafetySettings.SetHighSafety()
	case Default:
		*cmd.SafetySettings = *DefaultSafetySettings()
	}
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
	return parts[languageFlagIndex] == AITranslateCommandsArg
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
