// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License
//
// Note: By simplify like this, list function here can be reusable

package terminal

import (
	"fmt"
	"os"
)

// addMessageWithContext adds a message to the chat history with context.
func addMessageWithContext(session *Session, sender, message string) {
	session.ChatHistory.AddMessage(sender, message, session.ChatConfig)
}

// executeCommand is a generic function to execute a command.
func executeCommand(session *Session, command string, constructPrompt func(string) string) (bool, error) {
	// Assuming command is the user input that triggered the AI.
	// Note: The command execution process is now more dynamic.
	addMessageWithContext(session, YouNerd, command)
	success, err := sendCommandToAI(session, command, constructPrompt)
	if err != nil {
		logger.Error(ErrorFailedToSendCommandToAI, err)
		return false, err
	}
	return !success, nil // Return false to continue the session if successful.
}

// sendCommandToAI sends a command to the AI after sanitizing and applying retry logic.
func sendCommandToAI(session *Session, command string, constructPrompt func(string) string) (bool, error) {
	// Construct the AI prompt using the provided function.
	aiPrompt := constructPrompt(command)

	// Define a retryable operation with a function that sends a message to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// This function is called repeatedly by retryWithExponentialBackoff if it fails.
			return sendMessageToAI(session, aiPrompt)
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	return operation.retryWithExponentialBackoff(standardAPIErrorHandler)
}

// sendMessageToAI sends a message to the AI and handles the response.
func sendMessageToAI(session *Session, message string) (bool, error) {
	// Fix Duplicated by using Magic "_" Identifier
	_, err := session.SendMessage(session.Ctx, session.Client, message)
	return err == nil, err
}

// sendShutdownMessage sends a formatted shutdown message to the AI and logs it to the chat history.
func sendShutdownMessage(session *Session) error {
	// Clear the chat history in preparation for shutdown.
	session.ChatHistory.Clear()
	// Add context and quit command messages to the chat history.
	addMessageWithContext(session, AiNerd, ContextPrompt)
	addMessageWithContext(session, YouNerd, QuitCommand)

	// Construct the AI prompt for shutdown.
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, QuitCommand, ApplicationName)

	// Define a retryable operation with a function that sends the shutdown message to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// This function is called repeatedly by retryWithExponentialBackoff if it fails.
			return sendMessageToAI(session, aiPrompt)
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	_, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)
	return err
}

// constructSummarizePrompt constructs the prompt to be sent to the AI for summarization.
func (h *handleSummarizeCommand) constructSummarizePrompt() string {
	return fmt.Sprintf(SummarizePrompt)
}

// sendSummarizePrompt sends the summarize prompt to the AI and handles the response.
func (h *handleSummarizeCommand) sendSummarizePrompt(session *Session, sanitizedMessage string) (bool, error) {
	// Define a retryable operation with a function that sends the summarize prompt to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// Note: This is subject to change, for example,
			// to implement another functionality without displaying AI response in the terminal,
			// but only adding it to the chat history.
			// This function is called repeatedly by retryWithExponentialBackoff if it fails.
			aiResponse, err := session.SendMessage(session.Ctx, session.Client, sanitizedMessage)
			if err != nil {
				return false, err
			}
			// Handle the AI's response after a successful send.
			h.handleAIResponse(session, sanitizedMessage, aiResponse, SummaryPrefix)
			return true, nil
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	return operation.retryWithExponentialBackoff(standardAPIErrorHandler)
}

// handleAIResponse processes the AI's response to the summarize command.
func (h *handleSummarizeCommand) handleAIResponse(session *Session, sanitizedMessage, aiResponse, customText string) {
	// Concatenate custom text with the AI response.
	fullResponse := SummaryPrefix + aiResponse

	// Instead of directly adding, check if a system message already exists and replace it.
	formattedResponse := fmt.Sprintf(ObjectHighLevelString, SYSTEMPREFIX, fullResponse)
	if !session.ChatHistory.handleSystemMessage(sanitizedMessage, formattedResponse, session.ChatHistory.hashMessage(fullResponse)) {
		// If it was not a system message or no existing system message was found to replace,
		// add the new system message to the chat history.
		session.ChatHistory.AddMessage(SYSTEMPREFIX, fullResponse, session.ChatConfig)
	}
}

// Note: This approach simplifies maintenance and improvements by abstracting logic in this manner,
// in contrast to less optimal practices where functions are made overly complex (e.g, stupid human) with excessive conditional statements.
func (cmd *handleTokeCountingCommand) handleTokenCount(apiKey, filePath string, session *Session) (bool, error) {
	// Verify the file extension before reading the file.
	params, err := cmd.prepareTokenCountParams(apiKey, filePath)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		logger.Error(ErrorFailedToReadFile, filePath, err) // Using logger.Error with formatting directive.
		return false, nil
	}

	tokenCount, err := params.CountTokens()
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		logger.Error(ErrorFailedToCountTokens, filePath, err) // Using logger.Error with formatting directive.
		return false, nil
	}
	logger.Any(InfoTokenCountFile, filePath, tokenCount)
	return false, nil // Continue the session after displaying the token count.
}

func (cmd *handleTokeCountingCommand) prepareTokenCountParams(apiKey, filePath string) (TokenCountParams, error) {
	var params TokenCountParams
	params.APIKey = apiKey

	if isImage := verifyImageFileExtension(filePath) == nil; isImage {
		if err := cmd.readImageFile(filePath, &params); err != nil {
			return TokenCountParams{}, err
		}
	} else {
		if err := cmd.readTextFile(filePath, &params); err != nil {
			return TokenCountParams{}, err
		}
	}

	return params, nil
}

func (cmd *handleTokeCountingCommand) readImageFile(filePath string, params *TokenCountParams) error {
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		return fmt.Errorf(ObjectHighLevelFMT, ErrorFailedToReadFile, err) // low level in 2024
	}
	// Note: Avoid attempting to inspect "imageData" using fmt.Println(imageData) unless you are professional/master of go programming
	// as it will literally print the binary data of the image.
	params.ImageData = imageData
	params.ImageFormat = getImageFormat(filePath)
	params.ModelName = GeminiProVision
	return nil
}

func (cmd *handleTokeCountingCommand) readTextFile(filePath string, params *TokenCountParams) error {
	if err := verifyFileExtension(filePath); err != nil {
		return err
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		return fmt.Errorf(ObjectHighLevelFMT, ErrorFailedToReadFile, err)
	}
	params.Input = string(fileContent)
	params.ModelName = GeminiPro
	return nil
}

// constructAITranslatePrompt constructs the AI translation prompt.
//
// Note: This is currently unstable and will be fixed later. The issue lies only with the prompt.
func constructAITranslatePrompt(applicationName, command, text, targetLanguage string) string {
	return fmt.Sprintf(AITranslateCommandPrompt,
		applicationName,
		command,
		text,
		targetLanguage)
}

// handleAIInteraction handles sending messages to the AI and processing responses.
func handleAIInteraction(session *Session, aiPrompt string, postProcess func(session *Session, aiResponse string) error) error {
	// Sanitize the AI prompt to ensure it is safe to send.
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	// Define a retryable operation with a function that sends messages to the AI.
	operation := RetryableOperation{
		retryFunc: func() (bool, error) {
			// This function is called repeatedly by retryWithExponentialBackoff if it fails.
			aiResponse, err := session.SendMessage(session.Ctx, session.Client, sanitizedMessage)
			if err != nil {
				return false, err
			}
			// Call the provided post-processing function on the AI's response.
			return true, postProcess(session, aiResponse)
		},
	}

	// Execute the retryable operation with an exponential backoff strategy.
	success, err := operation.retryWithExponentialBackoff(standardAPIErrorHandler)

	if err != nil {
		// If an error occurs that is not recoverable by retries, return the error.
		return err
	}

	if !success {
		// If the operation was not successful after retries, return an error.
		return fmt.Errorf(ErrorFailedToSendTranslationMessageAfterRetries)
	}

	return nil
}

// postProcessAITranslate handles the post-processing of AI's response for translation.
func postProcessAITranslate(session *Session, aiResponse string) error {
	// Sanitize AI's response to remove any separators
	aiResponse = sanitizeAIResponse(aiResponse)
	// Add the sanitized AI's response to the chat history
	session.ChatHistory.AddMessage(AiNerd, aiResponse, session.ChatConfig)
	return nil
}
