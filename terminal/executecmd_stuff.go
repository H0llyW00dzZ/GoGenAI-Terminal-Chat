// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"os"
)

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
	logger.Any(clearMessage) // simplify
	// Added back the context prompt after clearing the chat history
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	return false, nil // Continue the session
}

// clearSummarizeHistory clears the summarized messages from the chat history.
func (cmd *handleClearCommand) clearSummarizeHistory(session *Session) (bool, error) {
	session.ChatHistory.ClearAllSystemMessages()
	logger.Any(ChatSysSummaryMessages) // simplify
	session.ChatHistory.AddMessage(AiNerd, ContextPrompt, session.ChatConfig)
	return false, nil // Continue the session
}

// showChatStats displays the chat statistics in a session. It retrieves the statistics
// from the session's ChatHistory and prints them to the console with a typing effect.
// The typing delay is specified by the TypingDelay constant/variable.
// The function prints a system message prefix with a timestamp before the stats.
// After printing the stats, it continues the session without error.
func (cmd *handleStatsCommand) showChatStats(session *Session) (bool, error) {
	// Retrieve chat statistics from the session's ChatHistory.
	stats := session.ChatHistory.GetMessageStats()

	// Use the logger's Any method to print the statistics without colorization.
	// The SYSTEMPREFIX is included directly in the formatted message.
	logger.Any(ListChatStats,
		stats.UserMessages, stats.AIMessages, stats.SystemMessages)

	return false, nil // Continue the session without error.
}
