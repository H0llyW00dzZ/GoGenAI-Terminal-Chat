// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"log"
	"strings"
)

// PrintHistory outputs all messages in the chat history to the standard output,
// one message per line. This method is useful for displaying the chat history
// directly to the terminal.
//
// Each message is printed in the order it was added, preserving the conversation
// flow. This method does not return any value or error.
//
// Deprecated: This method is deprecated and was replaced by GetHistory.
// See GetHistory for current functionality.
func (h *ChatHistory) PrintHistory() {
	for _, msg := range h.Messages {
		fmt.Println(msg)
	}
}

// RecoverFromPanic returns a deferred function that recovers from panics within a goroutine
// or function, preventing the panic from propagating and potentially causing the program to crash.
// Instead, it logs the panic information using the standard logger, allowing for post-mortem analysis
// without interrupting the program's execution flow.
//
// Usage:
//
//	defer terminal.RecoverFromPanic()()
//
// The function returned by RecoverFromPanic should be called by deferring it at the start of
// a goroutine or function. When a panic occurs, the deferred function will handle the panic
// by logging its message and stack trace, as provided by the recover built-in function.
//
// Deprecated: This method is deprecated and was replaced by logger.RecoverFromPanic.
// Refer to logger.RecoverFromPanic for the updated implementation.
func RecoverFromPanic() func() {
	return func() {
		if r := recover(); r != nil {
			// Log the panic with additional context if desired
			log.Printf("Recovered from panic: %+v\n", r)
		}
	}
}

// IsANSISequence checks if the current index in the rune slice is the start of an ANSI sequence.
//
// Deprecated: This method is no longer used, and was replaced by SanitizeMessage.
// Consult SanitizeMessage for handling ANSI sequences in messages.
func IsANSISequence(runes []rune, index int) bool {
	return index+1 < len(runes) && runes[index] == ansichar.BinaryAnsiChar && runes[index+1] == ansichar.BinaryLeftSquareBracket
}

// PrintANSISequence prints the full ANSI sequence without delay and returns the new index.
//
// Deprecated: This method is no longer used, and was replaced by SanitizeMessage.
// See SanitizeMessage for current message processing.
func PrintANSISequence(runes []rune, index int) int {
	// Print the full ANSI sequence without delay.
	for index < len(runes) && runes[index] != ansichar.BinaryAnsiSquenseChar {
		fmt.Printf(humantyping.AnimatedChars, runes[index])
		index++ // Move past the current character.
	}
	if index < len(runes) {
		fmt.Printf(humantyping.AnimatedChars, runes[index]) // Print the 'm' character to complete the ANSI sequence.
	}
	return index // Return the new index position.
}

// ApplyBold applies bold formatting to the provided text if the delimiter indicates bold.
//
// Deprecated: This method is no longer used, and was replaced by ApplyFormatting.
// Refer to ApplyFormatting for text formatting options.
func ApplyBold(text string, delimiter string, color string) string {
	if delimiter == DoubleAsterisk {
		return color + BoldText + text + ResetBoldText + ColorReset
	}
	return color + text + ColorReset
}

// HandleUnrecognizedCommand takes an unrecognized command and the current session,
// constructs a prompt to inform the AI about the unrecognized command, and sends
// this information to the AI service. This function is typically called when a user
// input is detected as a command but does not match any of the known command handlers.
//
// Parameters:
//
//	command string: The unrecognized command input by the user.
//	session *Session: The current chat session containing state and context, including the AI client.
//
// Returns:
//
//	bool: Always returns false as this function does not result in a command execution.
//	error: Returns an error if sending the message to the AI fails; otherwise, nil.
//
// The function constructs an error prompt using the application's name and the unrecognized command,
// retrieves the current chat history, and sends this information to the AI service. If an error occurs
// while sending the message, the function logs the error and returns an error to the caller.
//
// Deprecated: This method is no longer used, and was replaced by CommandRegistry.
// See CommandRegistry for handling unrecognized commands.
func HandleUnrecognizedCommand(command string, session *Session, parts []string) (bool, error) {
	// Debug
	logger.Debug(DEBUGEXECUTINGCMD, command, parts)
	// Pass ContextPrompt
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	// If the command is not recognized, inform the AI about the unrecognized command.
	aiPrompt := fmt.Sprintf(ErrorUserAttemptUnrecognizedCommandPrompt, ApplicationName, command)
	// Sanitize the message before sending it to the AI
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Send the constructed message to the AI and get the response.
	_, err := session.SendMessage(session.Ctx, session.Client, sanitizedMessage)
	if err != nil {
		errMsg := fmt.Sprintf(ErrorFailedtoSendUnrecognizedCommandToAI, err)
		logger.Error(errMsg)
		return false, fmt.Errorf(errMsg)
	}
	return false, nil
}

// SafetySetter defines a function signature for functions that modify the safety settings
// of a generative AI model. It takes a pointer to a SafetySettings struct and applies
// specific safety configurations to it. This type is used to abstract the modification
// of safety settings so that different configurations can be applied without changing
// the struct directly.
//
// Deprecated: This type is deprecated and has been replaced by the SafetyOption struct.
// Refer to SafetyOption for safety settings configurations.
type SafetySetter func(*SafetySettings)

// ASCIIPatterns Define the ASCII patterns for the 'slant' font for the characters
//
// Deprecated: This variable is no longer used, and was replaced by NewASCIIArtStyle().
// Use NewASCIIArtStyle() for ASCII art styles and colors.
var ASCIIPatterns = map[rune][]string{
	// Figlet in a compiled language, not an interpreted language.
	// This literally header in your machine lmao.
	// It so easy implement Header like this in go, also it possible to made it animated drawing/human typing this ascii art
	// unlike "interpreted language" 🤪
	G: {
		_G,
		_O,
		_GEN,
		A,
		I,
	},
	V: {
		Blank,
		Blank,
		Blank, // TODO: Implement a notification to display here when a new version is available.
		//			 For checking the version and viewing the change log, implement the command ":checkversion".
		CurrentVersioN,
		Copyright,
	},
}

// ASCIIColors Define a map for character Ascii colors
//
// Deprecated: This variable is no longer used, and was replaced by NewASCIIArtStyle().
// Consult NewASCIIArtStyle() for ASCII color configurations.
var ASCIIColors = map[rune]string{
	G: BoldText + colors.ColorHex95b806,
	V: BoldText + colors.ColorCyan24Bit,
}

// PrintAnotherVisualSeparator prints a visual separator to the standard output.
//
// Deprecated: This method is no longer used, and was replaced by NewASCIIArtStyle().
// NewASCIIArtStyle() should be used for visual separators.
func PrintAnotherVisualSeparator() {
	fmt.Println(colors.ColorCyan24Bit + StripChars + colors.ColorReset)
}

// ReplaceTripleBackticks replaces all occurrences of triple backticks with a placeholder.
//
// Deprecated: This method is no longer used, and was replaced by Colorize().
// Refer to Colorize() for replacing triple backticks in text.
func ReplaceTripleBackticks(text, placeholder string) string {
	for {
		index := strings.Index(text, TripleBacktick)
		if index == -1 {
			break
		}
		text = strings.Replace(text, TripleBacktick, placeholder, 1)
	}
	return text
}

// IsLastMessage checks if the current index is the last message in the slice
//
// Deprecated: This method is no longer used, and was replaced by separateSystemMessages.
// See separateSystemMessages for message handling logic.
func IsLastMessage(index int, messages []string) bool {
	return index == len(messages)-1
}

// IncrementMessageTypeCount updates the count of messages for the given type.
// It specifically flags if a SystemMessage is encountered, as it may require special handling.
// Returns true if the incremented message type is a system message.
//
// Deprecated: This method is no longer used, and was replaced by ChatHistory.AddMessage.
// Consult ChatHistory.AddMessage for message type counting and handling.
func (h *ChatHistory) IncrementMessageTypeCount(messageType MessageType) bool {
	switch messageType {
	case UserMessage:
		h.UserMessageCount++
	case AIMessage:
		h.AIMessageCount++
	case SystemMessage:
		h.SystemMessageCount++
		return true // Indicates that a system message has been processed.
	}
	return false
}
