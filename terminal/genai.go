// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"context"
	"fmt"
	"time"

	genai "github.com/google/generative-ai-go/genai"
)

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

				// Colorize the response
				colorPairs := []string{
					DoubleAsterisk, ColorHex95b806,
					SingleBacktick, ColorYellow,
				}

				keepDelimiters := map[string]bool{
					DoubleAsterisk: false, // Remove double asterisks from the output
					SingleBacktick: true,  // Keep single backticks in the output
				}

				// Colorize content that is surrounded by double asterisks or backticks
				colorized := Colorize(content, colorPairs, keepDelimiters)

				// Handle single asterisks separately
				// Pass Colorize content that is surrounded by single-character delimiters
				colorized = SingleCharColorize(colorized, SingleAsterisk, ColorCyan24Bit)

				// Print "AI:" prefix directly without typing effect
				PrintPrefixWithTimeStamp(AiNerd)

				// Assuming 'part' can be printed directly and is of type string or has a String() method
				// Use the typing banner effect only for the part content
				// Colorized string is printed character by character with a delay between each character
				PrintTypingChat(colorized, TypingDelay)
				// Collect AI response
				aiResponse += colorized
			}
		}
	}
	fmt.Println(StringNewLine + StripChars)
	fmt.Print(StringNewLine)
	return aiResponse
}
