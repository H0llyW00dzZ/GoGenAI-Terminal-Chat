// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"os"
	"strings"
	"time"

	genai "github.com/google/generative-ai-go/genai"
)

// PrintPrefixWithTimeStamp prints a message to the terminal, prefixed with a formatted timestamp.
// The timestamp is formatted according to the TimeFormat constant.
//
// For example, with TimeFormat set to "2006/01/02 15:04:05" and the prefix "🤓 You: ",
// the output might be "2024/01/10 16:30:00 🤓 You:".
//
// This function is designed for terminal outputs that benefit from a timestamped context,
// providing clarity and temporal reference for the message displayed.
//
// The prefix parameter is appended to the timestamp and can be a log level, a descriptor,
// or any other string that aids in categorizing or highlighting the message.
func PrintPrefixWithTimeStamp(prefix string) {
	currentTime := time.Now().Format(TimeFormat)
	fmt.Printf(ObjectHighLevelString, currentTime, prefix)
}

// printPromptFeedback formats and prints the prompt feedback received from the AI.
func printPromptFeedback(feedback *genai.PromptFeedback) {
	fmt.Print(StringNewLine)
	if feedback == nil {
		return
	}
	// Iterate over safety ratings and print them.
	for _, rating := range feedback.SafetyRatings {
		safetyPrefix := ShieldEmoji
		PrintPrefixWithTimeStamp(safetyPrefix)
		promptFeedback := fmt.Sprintf(PROMPTFEEDBACK, rating.Category.String(), rating.Probability.String())
		PrintTypingChat(promptFeedback, TypingDelay)
	}
	// fix front end lmao
	printVisualSeparator()
}

// printTokenCount prints the number of tokens used in the AI's response, including the chat history.
// It also updates and prints the total token count for the session.
func printTokenCount(apiKey, aiResponse string, chatHistory ...string) {
	// Concatenate chat history and AI response for token counting
	fullText := concatenateChatHistory(aiResponse, chatHistory...)

	tokenCount, err := CountTokens(apiKey, fullText)
	fmt.Print(StringNewLine)
	if err != nil {
		handleTokenCountError(err)
		return
	}
	// print the current token count
	printCurrentTokenCount(tokenCount)
	// update and print the total token count
	updateAndPrintTotalTokenCount(tokenCount)

	// Visual separator for clarity in the output
	printVisualSeparator()
}

// handleTokenCountError handles errors that occur while counting tokens.
func handleTokenCountError(err error) {
	logger.Error(ErrorCountingTokens, err)
	logger.HandleGoogleAPIError(err)
}

// printCurrentTokenCount prints the number of tokens used in the AI's response.
func printCurrentTokenCount(tokenCount int) {
	tokenPrefix := TokenEmoji
	tokenMSG := fmt.Sprintf(TokenCount, tokenCount)
	PrintPrefixWithTimeStamp(tokenPrefix)
	PrintTypingChat(tokenMSG, TypingDelay)
}

// updateAndPrintTotalTokenCount updates the total token count for the session and prints it.
func updateAndPrintTotalTokenCount(tokenCount int) {
	totalTokenCount += tokenCount // Assuming totalTokenCount is a global or package-level variable
	tokenUsageMSG := fmt.Sprintf(TotalTokenCount, totalTokenCount)
	PrintPrefixWithTimeStamp(StatisticsEmoji)
	PrintTypingChat(tokenUsageMSG, TypingDelay)
}

// printVisualSeparator prints a visual separator to the standard output.
func printVisualSeparator() {
	text := "V"
	asciiArt, _ := ToASCIIArt(text, stripStyle)
	fmt.Println(asciiArt)
}

// removeAIPrefix checks for and removes the AI prefix if it's present in the response.
func removeAIPrefix(content string) string {
	aiPrefix := AiNerd // Define the AI prefix
	if strings.HasPrefix(content, aiPrefix) {
		return strings.TrimPrefix(content, aiPrefix)
	}
	return content
}

// colorizeResponse applies color to the response content.
func colorizeResponse(content string) string {
	// Define color pairs and delimiters for colorization
	colorPairs := []string{
		TripleBacktick, colors.ColorPurple24Bit,
		SingleBacktick, colors.ColorYellow,
		DoubleAsterisk, colors.ColorHex95b806,
	}
	keepDelimiters := map[string]bool{
		TripleBacktick: true,  // Keep triple backticks in the output
		SingleBacktick: false, // Remove single backticks in the output
		DoubleAsterisk: false, // Remove double asterisks from the output
	}
	formatting := map[string]string{
		DoubleAsterisk: BoldText, // Assuming DoubleAsterisk in the output
	}

	return Colorize(content, colorPairs, keepDelimiters, formatting)
}

// handleSingleAsterisks applies color to text surrounded by single-character delimiters.
func handleSingleAsterisks(content string) string {
	return SingleCharColorize(content, SingleAsterisk, colors.ColorCyan24Bit)
}

// printAIResponse prints the AI's response with a typing effect.
func printAIResponse(colorized string) {
	PrintPrefixWithTimeStamp(AiNerd)
	PrintTypingChat(colorized, TypingDelay)
}

// printResponseFooter prints the footer after the AI response and includes prompt feedback and token count if enabled.
//
// Note: this functionality are powerful, it won't break a current session of conversation hahaha.
func printResponseFooter(resp *genai.GenerateContentResponse, aiResponse string) {
	showPromptFeedback := os.Getenv(SHOW_PROMPT_FEEDBACK) == "true"
	showTokenCount := os.Getenv(SHOW_TOKEN_COUNT) == "true"

	// Print the footer separator
	printVisualSeparator()

	// Print prompt feedback if enabled
	if showPromptFeedback && resp.PromptFeedback != nil {
		printPromptFeedback(resp.PromptFeedback)
	}

	// Print token count if enabled
	if showTokenCount {
		apiKey := os.Getenv(API_KEY) // Retrieve the API_KEY from the environment
		printTokenCount(apiKey, aiResponse)
	}

	// Print the closing footer separator
	fmt.Println() // fix front end issue lmao
}
