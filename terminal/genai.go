// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"context"
	"fmt"
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
// which may include a portion of the previous chat history determined by the session's ChatConfig,
// to the AI model for generating a response.
//
// Parameters:
//
//	ctx context.Context: The context for controlling the cancellation of the request.
//	client *genai.Client: The client instance used to create a generative model session and send messages to the AI model.
//	chatContext string: The chat context or message to be sent to the AI model.
//	session *Session: The current chat session containing the ChatHistory and ChatConfig. The ChatConfig determines
//	                  how much of the chat history is sent to the AI for context.
//
// Returns:
//
//	string: The AI's response as a string, which includes the AI's message with a simulated typing effect.
//	error: An error message if the message sending or response retrieval fails. If the operation is successful,
//	       the error is nil.
//
// The function initializes a new chat session and sends the chat context, along with the portion of chat history
// specified by the session's ChatConfig, to the generative AI model. It then calls `printResponse` to process
// and print the AI's response. The final AI response is returned as a concatenated string of all parts from the AI response.
func SendMessage(ctx context.Context, client *genai.Client, chatContext string, session *Session) (string, error) {
	model := client.GenerativeModel(ModelAi)
	cs := model.StartChat()
	// Note: This is a good balance between safety and readability.
	// It allows for a wider range of content to be generated while still maintaining a reasonable level of safety.
	// Additional Note: This method unlike static "model.SafetySettings = []*genai.SafetySetting" in official genai docs lmao.
	safetySettings := DefaultSafetySettings()
	safetySettings.ApplyToModel(session.Client.GenerativeModel(ModelAi))
	tempOption := WithTemperature(0.9) // Improved output, especially when using :summarize
	ApplyOptions(model, tempOption)
	// Retrieve the relevant chat history using ChatConfig
	chatHistory := session.ChatHistory.GetHistory(session.ChatConfig)

	fullContext := chatContext
	if len(chatHistory) > 0 {
		// Append the new message to the chat history to form the full context
		fullContext = chatHistory + StringNewLine + chatContext
	}

	resp, err := cs.SendMessage(ctx, genai.Text(fullContext))
	if err != nil {
		return "", err
	}

	return printResponse(resp), nil
}

// SendDummyMessage verifies the validity of the API key by sending a dummy message.
//
// Parameters:
//
//	client *genai.Client: The AI client used to send the message.
//
// Returns:
//
//	A boolean indicating the validity of the API key.
//	An error if sending the dummy message fails.
func SendDummyMessage(client *genai.Client) (bool, error) {
	// Initialize a dummy chat session or use an appropriate lightweight method.
	model := client.GenerativeModel(ModelAi)
	// Configure the model with options.
	// Apply the configurations to the model.
	// Note: This a testing in live production by sending a Dummy messages lmao
	tempOption := WithTemperature(0.9)
	topPOption := WithTopP(0.5)
	topKOption := WithTopK(20)
	// Exercise caution: setting the max output tokens below 50 may cause a panic.
	// This could be a bug in official genai package or an unintended issue from Google's side.
	maxOutputTokensOption, err := WithMaxOutputTokens(50)
	if err != nil {
		return handleGenAIError(err)
	}

	success, err := ApplyOptions(model, tempOption, topPOption, topKOption, maxOutputTokensOption)
	if !success {
		return false, fmt.Errorf(ErrorFailedToApplyModelConfiguration)
	}

	cs := model.StartChat()

	// Attempt to send a dummy message.
	resp, err := cs.SendMessage(context.Background(), genai.Text(DummyMessages))
	if err != nil {
		return handleGenAIError(err)
	}

	// A non-nil response indicates a valid API key.
	return resp != nil, nil
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
				// Filter out the language identifier from code blocks before any other processing
				filteredContent := FilterLanguageFromCodeBlock(content)
				colorized := colorizeResponse(filteredContent)
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

// concatenateChatHistory concatenates the AI's response with the chat history.
func concatenateChatHistory(aiResponse string, chatHistory ...string) string {
	if len(chatHistory) > 0 {
		return strings.Join(chatHistory, StringNewLine) + StringNewLine + aiResponse
	}
	return aiResponse
}

// sanitizeAIResponse removes unwanted separators from the AI's response.
func sanitizeAIResponse(response string) string {
	// Split the response by the separator
	parts := strings.Split(response, SanitizeTextAIResponse)

	// Rejoin the parts without the separator to get the sanitized response
	sanitizedResponse := strings.Join(parts, StringNewLine)

	return sanitizedResponse
}
