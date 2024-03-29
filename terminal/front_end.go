// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"os"
	"strings"
	"time"
	"unicode/utf8"

	genai "github.com/google/generative-ai-go/genai"
)

// FilterLanguageFromCodeBlock searches for Markdown code block delimiters with a language identifier
// (e.g., "```go") and removes the language identifier, leaving just the code block delimiters.
// This function is useful when the language identifier is not required, such as when rendering
// plain text or when the syntax highlighting is not supported.
//
// The function uses a precompiled regular expression `filterCodeBlock` that matches the pattern
// of triple backticks followed by any word characters (representing the language identifier).
// It replaces this pattern with just the triple backticks, effectively stripping the language
// identifier from the code block.
//
// Parameters:
//
//	text (string): The input text containing Markdown code blocks with language identifiers.
//
// Returns:
//
//	string: The modified text with language identifiers removed from all code blocks.
//
// Example:
//
//	input := "Here is some Go code:\n```go\nfmt.Println(\"Hello, World!\")\n```"
//	output := FilterLanguageFromCodeBlock(input)
//	// output now contains "Here is some Go code:\n```\nfmt.Println(\"Hello, World!\")\n```"
func FilterLanguageFromCodeBlock(text string) string {
	return filterCodeBlock.ReplaceAllString(text, TripleBacktick)
}

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
func PrintPrefixWithTimeStamp(prefix, message string) {
	currentTime := time.Now().Format(TimeFormat)
	// Check if the first character is potentially an emoji or wide character.
	if isFirstCharacterWide(prefix) {
		// Add an extra space after the prefix to ensure separation in terminals that might not handle wide characters well.
		fmt.Printf(ObjectHighLevelTripleString, currentTime, prefix, message)
	}
}

// isFirstCharacterWide checks if the first character of the string is wider than a standard character.
// This is a simple heuristic based on the assumption that most emojis and wide characters
// have a UTF-8 encoded length greater than 1. This won't cover all cases but works for many emojis.
func isFirstCharacterWide(s string) bool {
	if len(s) == 0 {
		return false
	}
	_, size := utf8.DecodeRuneInString(s)
	// Assuming characters with a UTF-8 size greater than 1 are wide (e.g., most emojis).
	return size > 1
}

// printPromptFeedback formats and prints the prompt feedback received from the AI.
func printPromptFeedback(feedback *genai.PromptFeedback) {
	printnewlineASCII() // this better new line instead of "\n" for front end hahaha
	if feedback == nil {
		return
	}
	// Iterate over safety ratings and print them.
	for i, rating := range feedback.SafetyRatings {
		safetyPrefix := ShieldEmoji
		PrintPrefixWithTimeStamp(safetyPrefix, "")
		promptFeedback := fmt.Sprintf(PROMPTFEEDBACK, rating.Category.String(), rating.Probability.String())
		if i < len(feedback.SafetyRatings)-1 {
			promptFeedback += StringNewLine
		}
		humanTyping := NewTypingPrinter()
		humanTyping.Print(promptFeedback, TypingDelay)
	}
	// fix front end lmao
	printVisualSeparator()
}

// printTokenCount prints the number of tokens used in the AI's response, including the chat history.
// It now supports image data for the GeminiProVision model.
func (s *Session) printTokenCount(apiKey, aiResponse string, chatHistory ...string) {
	// Concatenate chat history and AI response for token counting
	fullText := concatenateChatHistory(aiResponse, chatHistory...)
	// Use the default model name
	modelName := s.DefaultModelName
	if s.CurrentModelName != "" {
		modelName = s.CurrentModelName // Override with the current model name if set
	}
	var imageData [][]byte
	var imageFormat string

	// Initialize TokenCountParams with the text input
	params := TokenCountParams{
		APIKey:    apiKey,
		ModelName: modelName,
		Input:     fullText,
	}

	// Check if the response contains image data for GeminiProVision
	//
	// Note: Not Ready Yet, it won't work hahaha
	if err := verifyImageFileExtension(aiResponse); err == nil {
		modelName = GeminiProVision
		var singleImageData []byte
		singleImageData, imageFormat = readImageFile(aiResponse) // Adjusted to match the returned values
		imageData = append(imageData, singleImageData)           // Append the single image data to the slice
		params.ImageData = imageData                             // Set the ImageData field
		params.ImageFormat = imageFormat                         // Set the ImageFormat field
		params.ModelName = modelName                             // Set the ModelName field
	}

	tokenCount, err := params.CountTokens() // Adjusted to use the CountTokens function directly
	printnewlineASCII()                     // a better one, instead of "\n"
	if err != nil {
		handleTokenCountError(err)
		return
	}
	// Print the current token count
	printCurrentTokenCount(tokenCount)
	// Update and print the total token count
	updateAndPrintTotalTokenCount(tokenCount)

	// Visual separator for clarity in the output
	printVisualSeparator()
}

// handleTokenCountError handles errors that occur while counting tokens.
func handleTokenCountError(err error) {
	logger.Error(ErrorCountingTokens, err)
}

// printCurrentTokenCount prints the number of tokens used in the AI's response.
func printCurrentTokenCount(tokenCount int) {
	tokenPrefix := TokenEmoji
	tokenMSG := fmt.Sprintf(TokenCount, tokenCount)
	humanTyping := NewTypingPrinter()
	PrintPrefixWithTimeStamp(tokenPrefix, "")
	humanTyping.Print(tokenMSG, TypingDelay)
}

// updateAndPrintTotalTokenCount updates the total token count for the session and prints it.
func updateAndPrintTotalTokenCount(tokenCount int) {
	totalTokenCount += tokenCount // Assuming totalTokenCount is a global or package-level variable
	tokenUsageMSG := fmt.Sprintf(TotalTokenCount, totalTokenCount)
	humanTyping := NewTypingPrinter()
	PrintPrefixWithTimeStamp(StatisticsEmoji, "")
	humanTyping.Print(tokenUsageMSG, TypingDelay)
}

// printVisualSeparator prints a visual separator to the standard output.
func printVisualSeparator() {
	text := "V"
	asciiArt, _ := ToASCIIArt(text, stripStyle)
	fmt.Println(asciiArt)
}

// printnewlineASCII prints a newline character as an ASCII art visual separator to the standard output.
func printnewlineASCII() {
	text := "N"
	asciiArt, _ := ToASCIIArt(text, newLine)
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

	// Create an instance of ColorizationOptions with the necessary information
	options := ColorizationOptions{
		Text:           content,
		ColorPairs:     colorPairs,
		KeepDelimiters: keepDelimiters,
		Formatting:     formatting,
	}

	// Call the Colorize function with the options struct
	return Colorize(options)
}

// handleSingleAsterisks applies color to text surrounded by single-character delimiters.
func handleSingleAsterisks(content string) string {
	return SingleCharColorize(content, SingleAsterisk, colors.ColorCyan24Bit)
}

// handleSingleMinusSign applies color to text surrounded by single-character delimiters.
func handleSingleMinusSign(content string) string {
	return SingleCharColorize(content, SingleMinusSign, colors.ColorCyan24Bit)
}

// printAIResponse prints the AI's response with a typing effect.
//
// Added a new parameter `isSystemMessage` to distinguish between AI and system messages.
//
// Note: By implementing it this way, the function becomes easier to maintain and more reusable.
func printAIResponse(response string, isSystemMessage bool) {
	var prefix string
	if isSystemMessage {
		prefix = SYSTEMPREFIX // Use system prefix for system messages
	} else {
		prefix = AiNerd // Use AI prefix for AI messages
	}
	humanTyping := NewTypingPrinter()
	PrintPrefixWithTimeStamp(prefix, "")
	humanTyping.Print(response, TypingDelay)
}

// printResponseFooter prints the footer after the AI response and includes prompt feedback and token count if enabled.
//
// Note: this functionality are powerful, it won't break a current session of conversation hahaha.
func (s *Session) printResponseFooter(resp *genai.GenerateContentResponse, aiResponse string) {
	showPromptFeedback := os.Getenv(ShowPromptFeedBack) == "true"
	showTokenCount := os.Getenv(ShowTokenCount) == "true"

	// Print the footer separator
	printVisualSeparator()

	// Print prompt feedback if enabled
	if showPromptFeedback && resp.PromptFeedback != nil {
		printPromptFeedback(resp.PromptFeedback)
	}

	// Print token count if enabled
	if showTokenCount {
		apiKey := os.Getenv(APIKey) // Retrieve the API_KEY from the environment
		s.printTokenCount(apiKey, aiResponse)
	}

	// Print the closing footer separator
	printnewlineASCII() // fix front end issue lmao
}

// NewTypingPrinter creates a new instance of TypingPrinter with a default print function.
// This constructor function is idiomatic in Go, providing a way to set up the struct
// with default values, which in this case is the PrintTypingChat function.
func NewTypingPrinter() *TypingPrinter {
	return &TypingPrinter{
		PrintFunc: PrintTypingChat, // Set the default typing effect function.
	}
}

// Print executes the typing effect for the given message using the configured PrintFunc.
// This method offers the flexibility to change the typing effect logic at runtime by
// assigning a different function to PrintFunc. The delay parameter controls the speed
// of the typing effect, making it adaptable to various use cases.
func (tp *TypingPrinter) Print(message string, delay time.Duration) {
	tp.PrintFunc(message, delay) // Delegate the print operation to the configured function.
}
