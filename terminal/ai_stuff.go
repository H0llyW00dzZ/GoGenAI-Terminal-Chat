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
	aiPrompt := constructPrompt(command)

	return retryWithExponentialBackoff(func() (bool, error) {
		return sendMessageToAI(session, aiPrompt)
	}, standardAPIErrorHandler)
}

// sendMessageToAI sends a message to the AI and handles the response.
func sendMessageToAI(session *Session, message string) (bool, error) {
	aiResponse, err := SendMessage(session.Ctx, session.Client, message, session)
	addMessageWithContext(session, AiNerd, aiResponse)
	return err == nil, err
}

// sendShutdownMessage sends a formatted shutdown message to the AI and logs it to the chat history.
func sendShutdownMessage(session *Session) error {
	session.ChatHistory.Clear()
	// Context
	addMessageWithContext(session, AiNerd, ContextPrompt)
	// Assuming QuitCommand is the user input that triggered the shutdown.
	addMessageWithContext(session, YouNerd, QuitCommand)

	// Send a shutdown message to the AI including the chat history with the context prompt
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, QuitCommand, ApplicationName)

	// Attempt to send the shutdown message to the AI with retry logic
	_, err := retryWithExponentialBackoff(func() (bool, error) {
		return sendMessageToAI(session, aiPrompt)
	}, standardAPIErrorHandler)

	return err
}

// constructSummarizePrompt constructs the prompt to be sent to the AI for summarization.
func (h *handleSummarizeCommand) constructSummarizePrompt() string {
	return fmt.Sprintf(SummarizePrompt)
}

// sendSummarizePrompt sends the summarize prompt to the AI and handles the response.
func (h *handleSummarizeCommand) sendSummarizePrompt(session *Session, sanitizedMessage string) (bool, error) {
	// Retry logic for sending the summarize prompt to the AI.
	return retryWithExponentialBackoff(func() (bool, error) {
		// Note: This is subject to change, for example,
		// to implement another functionality without displaying AI response in the terminal,
		// but only adding it to the chat history.
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		if err != nil {
			return false, err
		}

		h.handleAIResponse(session, sanitizedMessage, aiResponse, SummaryPrefix)
		return true, nil
	}, standardAPIErrorHandler)
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
		logger.Error("%v", err) // Using logger.Error with formatting directive.
		return false, nil
	}

	tokenCount, err := CountTokens(params)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		logger.Error("%v", err) // Using logger.Error with formatting directive.
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
		return fmt.Errorf(ObjectHighLevelFMT, ErrorFailedToReadFile, err)
	}
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
	sanitizedMessage := session.ChatHistory.SanitizeMessage(aiPrompt)

	success, err := retryWithExponentialBackoff(func() (bool, error) {
		aiResponse, err := SendMessage(session.Ctx, session.Client, sanitizedMessage, session)
		if err != nil {
			return false, err
		}
		return true, postProcess(session, aiResponse)
	}, standardAPIErrorHandler)

	if err != nil {
		return err
	}

	if !success {
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
