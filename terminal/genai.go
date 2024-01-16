// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	genai "github.com/google/generative-ai-go/genai"
)

// BinaryAnsiChars is a struct that contains the ANSI characters used to print the typing effect.
type BinaryAnsiChars struct {
	BinaryAnsiChar          rune
	BinaryAnsiSquenseChar   rune
	BinaryAnsiSquenseString string
	BinaryLeftSquareBracket rune
}

// TypingChars is a struct that contains the Animated Chars used to print the typing effect.
type TypingChars struct {
	AnimatedChars string
}

// PrintTypingChat simulates the visual effect of typing out a message character by character.
// It prints each character of a message to the standard output with a delay between each character
// to give the appearance of real-time typing.
//
// Parameters:
//
//	message string: The message to be displayed with the typing effect.
//	delay time.Duration: The duration to wait between printing each character.
//
// This function does not return any value. It directly prints to the standard output.
//
// Note: This is particularly useful for simulating the Gopher's lifecycle (Known as Goroutines) events in a user-friendly manner.
// For instance, when a Gopher completes a task or job and transitions to a resting state,
// this function can print a message with a typing effect to visually represent the Gopher's "sleeping" activities.
func PrintTypingChat(message string, delay time.Duration) {
	// Note: This now more like a human typing effect.
	// It's safe now alongside with sanitizing message.
	for _, char := range message {
		fmt.Printf(AnimatedChars, char)
		time.Sleep(delay)
	}
	fmt.Println()
}

// SendMessage sends a chat message to the generative AI model and retrieves the response.
// It constructs a chat session using the provided `genai.Client`, which is used to communicate
// with the AI service. The function simulates a chat interaction by sending the chat context,
// optionally preceded by previous chat history, to the AI model.
//
// Parameters:
//
//	ctx context.Context: The context for controlling the cancellation of the request.
//	client *genai.Client: The client instance used to create a generative model session and send messages to the AI model.
//	chatContext string: The chat context or message to be sent to the AI model.
//	chatHistory ...string: An optional slice of strings representing previous chat history. If provided,
//	                       it is prepended to the chatContext, separated by a newline, to provide context to the AI.
//
// Returns:
//
//	string: The AI's response as a string, which includes the AI's message with a simulated typing effect.
//	error: An error message if the message sending or response retrieval fails. If the operation is successful,
//	       the error is nil.
//
// The function initializes a new chat session and sends the chat context, along with any provided chat history,
// to the generative AI model. It then calls `printResponse` to process and print the AI's response. The final
// AI response is returned as a concatenated string of all parts from the AI response.
func SendMessage(ctx context.Context, client *genai.Client, chatContext string, chatHistory ...string) (string, error) {
	model := client.GenerativeModel(ModelAi)
	cs := model.StartChat()

	fullContext := chatContext
	if len(chatHistory) > 0 {
		fullContext = chatHistory[0] + StringNewLine + chatContext
	}

	resp, err := cs.SendMessage(ctx, genai.Text(fullContext))
	if err != nil {
		return "", err
	}

	return printResponse(resp), nil
}

// printResponse processes and prints the AI's response from the generative AI model.
// It iterates over the response candidates, prints each part with a typing effect, and
// aggregates the parts into a single response string.
//
// Parameters:
//
//	resp *genai.GenerateContentResponse: The response object from the AI model.
//
// Returns:
//
//	string: A concatenated string of all parts from the AI response.
//
// This function is unexported and is intended for internal use within the package.
func printResponse(resp *genai.GenerateContentResponse) string {
	aiResponse := ""
	// Note: this method are better instead of resp.Candidates[0] because it's more efficient and faster.
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content := fmt.Sprint(part)
				content = removeAIPrefix(content)
				colorized := colorizeResponse(content)
				colorized = handleSingleAsterisks(colorized)
				printAIResponse(colorized)
				aiResponse += colorized
			}
		}
	}
	// print the prompt feedback if it's present
	// why this so simple ? because it's more efficient and faster.
	// get good get "Go" hahaha
	printResponseFooter(resp, aiResponse)
	return aiResponse
}

// printPromptFeedback formats and prints the prompt feedback received from the AI.
func printPromptFeedback(feedback *genai.PromptFeedback) {
	fmt.Print(StringNewLine)
	if feedback == nil {
		return
	}
	// Iterate over safety ratings and print them.
	for _, rating := range feedback.SafetyRatings {
		fmt.Printf(PROMPTFEEDBACK, rating.Category.String(), rating.Probability.String())
	}
}

// printTokenCount prints the number of tokens used in the AI's response, including the chat history.
// It also updates and prints the total token count for the session.
func printTokenCount(apiKey, aiResponse string, chatHistory ...string) {
	// Concatenate chat history and AI response for token counting
	fullText := aiResponse
	if len(chatHistory) > 0 {
		fullText = chatHistory[0] + StringNewLine + aiResponse
	}

	tokenCount, err := CountTokens(apiKey, fullText)
	fmt.Print(StringNewLine)
	if err != nil {
		// Handle the error appropriately
		logger.Error(ErrorCountingTokens, err)
	} else {
		// Print the current and total token count
		tokenPrefix := TokenEmoji
		tokenMSG := fmt.Sprintf(TokenCount, tokenCount)
		PrintPrefixWithTimeStamp(tokenPrefix + " ")
		// Simulate typing the debug message
		PrintTypingChat(tokenMSG, TypingDelay)
		// Update the total token count
		totalTokenCount += tokenCount
		tokenusageMSG := fmt.Sprintf(TotalTokenCount, totalTokenCount)
		PrintPrefixWithTimeStamp(StatisticsEmoji + " ")
		PrintTypingChat(tokenusageMSG, TypingDelay)
	}
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
// Note: this functionality are powerful, it won't break session of conversation hahaha.
func printResponseFooter(resp *genai.GenerateContentResponse, aiResponse string) {
	showPromptFeedback := os.Getenv(SHOW_PROMPT_FEEDBACK) == "true"
	showTokenCount := os.Getenv(SHOW_TOKEN_COUNT) == "true"

	// Print the footer separator
	fmt.Println(StringNewLine + colors.ColorCyan24Bit + StripChars + colors.ColorReset)

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
	fmt.Println(StringNewLine + colors.ColorCyan24Bit + StripChars + colors.ColorReset)
	fmt.Print(StringNewLine)
}
