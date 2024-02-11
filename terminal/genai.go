// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"
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
	// Note: Improve a human typing effect.
	writer := bufio.NewWriter(os.Stdout) // Create a buffered writer

	for _, char := range message {
		// Additional Note: This improvement eliminates the use of fmt + animated characters, enhancing smoothness, especially with 100+ messages.
		// Also, ignore Go routines in pprof debugger that frequently switch (e.g., from 50 to 100 Go routines) and are waiting in "I/O Wait".
		writer.WriteString(string(char)) // Write to the buffer
		writer.Flush()                   // Flush the buffer to print the character group
		time.Sleep(delay)                // Sleep for the desired delay
	}

	printnewlineAscii()
	writer.Flush() // Make sure to flush any remaining output
}

// ConfigureModelForSession prepares and configures a generative AI model for use in a chat session.
// It applies safety settings from the session to the model and sets additional configuration options
// such as temperature. This function is essential for ensuring that the AI model behaves according
// to the desired safety guidelines and operational parameters before engaging in a chat session.
//
// Parameters:
//
//	ctx context.Context: A context.Context that carries deadlines, cancellation signals, and other request-scoped
//	       values across API boundaries and between processes.
//	client *genai.Client: A pointer to a genai.Client, which provides the functionality to interact with the
//	          generative AI service.
//	session *Session: A pointer to a Session struct that contains the current chat session's state,
//	           including safety settings and chat history.
//	modelName string: The identifier for the generative AI model to be used in the session. This name
//	                   is used to apply specific configurations and safety settings tailored to the
//	                   particular AI model.
//
// Returns:
//
//	*genai.GenerativeModel: A pointer to a generative AI model that is configured and ready for
//	                          initiating a chat session.
//
// Note: The function assumes that the client has been properly initialized and that the session
// contains valid safety settings. If no safety settings are present in the session, default
// safety settings are applied. The modelName parameter allows for model-specific configuration,
// enabling more granular control over the behavior and safety of different AI models.
func ConfigureModelForSession(ctx context.Context, client *genai.Client, session *Session, modelName string) *genai.GenerativeModel {
	// Initialize the model with the specific AI model identifier.
	model := client.GenerativeModel(modelName)

	// Apply safety settings from the session or use default settings if none are provided.
	if session.SafetySettings == nil {
		session.SafetySettings = DefaultSafetySettings()
	}
	session.SafetySettings.ApplyToModel(model, modelName)

	// Set additional configuration options, such as the temperature, to control the creativity
	// and randomness of the AI's responses.
	tempOption := WithTemperature(0.9)
	ApplyOptions(model, tempOption)

	return model
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
func (s *Session) SendMessage(ctx context.Context, client *genai.Client, chatContext string) (string, error) {
	// Get the generative model from the client
	model := ConfigureModelForSession(ctx, client, s, GeminiPro) // Simplify 🤪

	// Retrieve the relevant chat history using ChatConfig
	chatHistory := s.ChatHistory.GetHistory(s.ChatConfig)

	// Form the full context by appending the new message to the chat history
	fullContext := chatContext
	if len(chatHistory) > 0 {
		// Append the new message to the chat history to form the full context
		fullContext = chatHistory + StringNewLine + chatContext
	}

	// Start a new chat session with the model
	cs := model.StartChat()

	// Send the full context to the AI and get the response
	resp, err := cs.SendMessage(ctx, genai.Text(fullContext))
	if err != nil {
		logger.Error(ErrorFailedTosendmessagesToAI, err)
		return "", err
	}

	// Process the AI's response using the Session's method
	return s.printResponse(resp), nil
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
	model := client.GenerativeModel(GeminiPro)
	// Configure the model with options.
	// Apply the configurations to the model.
	// Note: This a testing in live production by sending a Dummy messages lmao
	tempOption := WithTemperature(0.9)
	topPOption := WithTopP(0.5)
	topKOption := WithTopK(20)
	// Exercise caution: setting the max output tokens below 20 may cause a panic.
	// This could be a bug in official genai package or an unintended issue from Google's side.
	maxOutputTokensOption, err := WithMaxOutputTokens(20)
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
func (s *Session) printResponse(resp *genai.GenerateContentResponse) string {
	aiResponse := ""
	// Note: this method are better instead of resp.Candidates[0] because it's more efficient and faster.
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content := fmt.Sprint(part)

				// Note: The function removeAIPrefix is invoked here to prevent the occurrence of duplicate AIPrefix entries in ChatHistory (known as RAM's labyrinth),
				// which could lead to confusion.
				// Remove the AI prefix from the content.
				content = removeAIPrefix(content)

				// Store the original AI response in the chat history (known as RAM's labyrinth)
				s.ChatHistory.AddMessage(AiNerd, content, s.ChatConfig)

				// Process the AI response for display
				// Filter out the language identifier from code blocks before any other processing
				filteredContent := FilterLanguageFromCodeBlock(content)
				colorized := colorizeResponse(filteredContent)
				colorized = handleSingleAsterisks(colorized)
				colorized = handleSingleMinusSign(colorized)

				// Display the processed AI response
				// Note: "false" indicate that AI Prefix not System Prefix
				// This how I like Go, unlike other language that sometimes not accurate about boolean lmao
				printAIResponse(colorized, false)
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

// sendToAIWithoutDisplay sends a message to the AI, processes the response, and updates the chat history without displaying the response.
//
// Note: This function is currently unused, but it will be employed for automated summarization in the future.
func (s *Session) sendToAIWithoutDisplay(ctx context.Context, client *genai.Client, chatContext string) error {
	model := ConfigureModelForSession(ctx, client, s, GeminiPro)

	// Retrieve the relevant chat history using ChatConfig
	chatHistory := s.ChatHistory.GetHistory(s.ChatConfig)

	fullContext := chatContext
	if len(chatHistory) > 0 {
		// Append the new message to the chat history to form the full context
		fullContext = chatHistory + StringNewLine + chatContext
	}

	return s.sendMessageAndProcessResponse(ctx, model, fullContext)
}

// sendMessageAndProcessResponse handles the full communication cycle with the generative AI model.
// It sends the provided context to the model, processes the response, and updates the chat history.
func (s *Session) sendMessageAndProcessResponse(ctx context.Context, model *genai.GenerativeModel, fullContext string) error {
	// Send the message to the AI
	resp, err := model.StartChat().SendMessage(ctx, genai.Text(fullContext))
	if err != nil {
		return err
	}

	// Process the AI's response and add it to the chat history
	aiResponse := processAIResponse(resp)
	sanitizedMessage := s.ChatHistory.SanitizeMessage(aiPrompt)
	formattedResponse := fmt.Sprintf(ObjectHighLevelString, SYSTEMPREFIX, aiResponse)
	if !s.ChatHistory.handleSystemMessage(sanitizedMessage, formattedResponse, s.ChatHistory.hashMessage(aiResponse)) {
		// If it was not a system message or no existing system message was found to replace,
		// add the new system message to the chat history.
		s.ChatHistory.AddMessage(SYSTEMPREFIX, aiResponse, s.ChatConfig)
	}

	return nil
}

// processAIResponse processes the AI's response and returns it as a string.
func processAIResponse(resp *genai.GenerateContentResponse) string {
	var aiResponse strings.Builder
	for _, cand := range resp.Candidates {
		if cand.Content != nil {
			for _, part := range cand.Content.Parts {
				content := fmt.Sprint(part)
				aiResponse.WriteString(content)
			}
		}
	}
	return aiResponse.String()
}
