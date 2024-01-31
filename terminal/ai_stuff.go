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
	success, err := sendCommandToAI(session, command, constructPrompt)
	if err != nil {
		logger.Error(ErrorFailedToSendCommandToAI, err)
		return false, err
	}
	return !success, nil // Return false to continue the session if successful.
}

// sendCommandToAI sends a command to the AI after sanitizing and applying retry logic.
func sendCommandToAI(session *Session, command string, constructPrompt func(string) string) (bool, error) {
	sanitizedCommand := session.ChatHistory.SanitizeMessage(command)
	aiPrompt := constructPrompt(sanitizedCommand)

	return retryWithExponentialBackoff(func() (bool, error) {
		return sendMessageToAI(session, aiPrompt)
	}, standardAPIErrorHandler)
}

// sendMessageToAI sends a message to the AI and handles the response.
func sendMessageToAI(session *Session, message string) (bool, error) {
	_, err := SendMessage(session.Ctx, session.Client, message, session)
	return err == nil, err
}

// sendShutdownMessage sends a formatted shutdown message to the AI and logs it to the chat history.
func sendShutdownMessage(session *Session) error {
	session.ChatHistory.Clear()
	session.ChatHistory.AddMessage(AiNerd,
		ContextPrompt,
		session.ChatConfig)
	// Assuming QuitCommand is the user input that triggered the shutdown.
	addMessageWithContext(session, StringNewLine+YouNerd, QuitCommand)

	// Sanitize the message before sending it to the AI.
	sanitizedMessage := session.ChatHistory.SanitizeMessage(QuitCommand)

	// Send a shutdown message to the AI including the chat history with the context prompt
	aiPrompt := fmt.Sprintf(ContextPromptShutdown, sanitizedMessage, ApplicationName)

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

		h.handleAIResponse(session, sanitizedMessage, aiResponse)
		return true, nil
	}, standardAPIErrorHandler)
}

// handleAIResponse processes the AI's response to the summarize command.
func (h *handleSummarizeCommand) handleAIResponse(session *Session, sanitizedMessage, aiResponse string) {
	// Instead of directly adding, check if a system message already exists and replace it.
	formattedResponse := fmt.Sprintf(ObjectHighLevelString, SYSTEMPREFIX, aiResponse)
	if !session.ChatHistory.handleSystemMessage(sanitizedMessage, formattedResponse, session.ChatHistory.hashMessage(aiResponse)) {
		// If it was not a system message or no existing system message was found to replace,
		// add the new system message to the chat history.
		session.ChatHistory.AddMessage(SYSTEMPREFIX, aiResponse, session.ChatConfig)
	}
}

func (cmd *handleTokeCountingCommand) handleTokenCount(apiKey, filePath string, session *Session) (bool, error) {
	// Verify the file extension before reading the file.
	if err := verifyFileExtension(filePath); err != nil {
		logger.Error(ErrorInvalidFileExtension, err)
		return false, nil
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		logger.Error(ErrorFailedToReadFile, err)
		return false, nil
	}

	text := string(fileContent)
	sanitizedMessage := session.ChatHistory.SanitizeMessage(text)
	tokenCount, err := CountTokens(apiKey, sanitizedMessage)
	if err != nil {
		logger.Error(ErrorFailedToCountTokens, err)
		return false, nil
	}
	logger.Any(InfoTokenCountFile, filePath, tokenCount)
	return false, nil // Continue the session after displaying the token count.
}

// constructAITranslatePrompt constructs the AI translation prompt.
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
