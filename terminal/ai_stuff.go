// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License
//
// Note: By simplify like this, list function here can be reusable

package terminal

import (
	"fmt"
	"os"
	"strings"

	genai "github.com/google/generative-ai-go/genai"
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

// handleTokenCount processes multiple file paths to count the number of tokens for each file.
// It uses the generative AI model to calculate token counts for text and image files.
// Only files with supported extensions are processed. Unsupported files are logged and skipped.
// After processing, it aggregates the token counts from all valid files and logs the total.
// If errors occur during file processing, they are logged, and the file is excluded from the total count.
//
// Parameters:
//
//	apiKey string: The API key used for authenticating requests to the AI service.
//	filePaths []string: A slice of file paths to be processed for token counting.
//	session *Session: The current chat session containing state and context.
//
// Returns:
//
//	bool: Always returns false to indicate the session should continue.
//	error: Returns nil as errors are handled internally and do not require external handling.
//
// Note: This approach simplifies maintenance and improvements by abstracting logic in this manner,
// in contrast to less optimal practices where functions are made overly complex (e.g, stupid human) with excessive conditional statements.
func (cmd *handleTokeCountingCommand) handleTokenCount(apiKey string, filePaths []string, session *Session) (bool, error) {
	var validFilePaths []string
	totalTokenCount := 0
	// Note: This functionality may only be compatible with Go version 1.22 and onwards hahahaha.
	// Additionally, while it may seem complex due to the 'if statement' error handling, it's not actually that complex.
	for _, filePath := range filePaths {
		// Prepare the parameters for token counting based on the file type.
		params, err := cmd.prepareTokenCountParams(apiKey, filePath)
		if err != nil {
			// Log the error directly without formatting.
			logger.Error("%s", err)
			continue
		}

		// Count the tokens using the prepared parameters.
		tokenCount, err := params.CountTokens()
		if err != nil {
			// Log the error with the file path and error details.
			logger.Error(ErrorFailedToCountTokens, filePath, err)
			continue
		}

		// Add the valid file path to the list and accumulate the token count.
		validFilePaths = append(validFilePaths, filePath)
		totalTokenCount += tokenCount
	}

	if len(validFilePaths) > 0 {
		// Log the total token count for all valid files if any valid files were processed.
		logger.Any(InfoTokenCountFile, strings.Join(validFilePaths, dotStringComma), totalTokenCount)
	}

	return false, nil // The session continues regardless of token counting results.
}

// prepareTokenCountParams prepares the necessary parameters for counting tokens based on file type.
// It differentiates between text and image files and sets up the parameters accordingly.
//
// Parameters:
//
//	apiKey string: The API key used for authenticating requests to the AI service.
//	filePath string: The path to the file for which tokens should be counted.
//
// Returns:
//
//	TokenCountParams: The prepared parameters for token counting.
//	error: An error if the file cannot be read or if the file extension is not supported.
//
// This method is responsible for determining the file type and reading the file's content.
// For image files, it also determines the image format, which is necessary for the token counting operation.
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

// readImageFile reads the content of an image file and updates the provided TokenCountParams
// with the image data and format.
//
// Parameters:
//
//	filePath string: The path to the image file.
//	params *TokenCountParams: A pointer to the TokenCountParams struct to be updated with the image data.
//
// Returns:
//
//	error: An error if the file cannot be read.
//
// This method reads the specified image file and sets the ImageData and ImageFormat fields
// in the TokenCountParams struct. It is used in conjunction with prepareTokenCountParams to
// handle image files for token counting.
func (cmd *handleTokeCountingCommand) readImageFile(filePath string, params *TokenCountParams) error {
	imageData, err := os.ReadFile(filePath)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		return fmt.Errorf(ErrorFailedToReadFile, filePath, err) // low level in 2024
	}
	// Note: Avoid attempting to inspect "imageData" using fmt.Println(imageData) unless you are professional/master of go programming
	// as it will literally print the binary data of the image.
	params.ImageData = [][]byte{imageData} // Wrap imageData in a slice of byte slices
	params.ImageFormat = getImageFormat(filePath)
	params.ModelName = GeminiProVision
	return nil
}

// readTextFile reads the content of a text file and updates the provided TokenCountParams
// with the text content.
//
// Parameters:
//
//	filePath string: The path to the text file.
//	params *TokenCountParams: A pointer to the TokenCountParams struct to be updated with the text content.
//
// Returns:
//
//	error: An error if the file cannot be read or if the file extension is not supported.
//
// This method reads the specified text file and sets the Input field in the TokenCountParams struct.
// It is used in conjunction with prepareTokenCountParams to handle text files for token counting.
func (cmd *handleTokeCountingCommand) readTextFile(filePath string, params *TokenCountParams) error {
	if err := verifyFileExtension(filePath); err != nil {
		return err
	}
	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		// Magic FMT, unlike stupid hard coding
		return fmt.Errorf(ErrorFailedToReadFile, filePath, err) // low level in 2024
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

// DisplayModelInfo formats and logs the information about a generative AI model.
// It compiles the model's details into a single string and logs it using the
// logger.Any method for a consistent logging experience.
//
// Parameters:
//
//	info *genai.ModelInfo: A pointer to the ModelInfo struct containing the model's metadata.
func DisplayModelInfo(info *genai.ModelInfo) {
	// Compile the model information into a single formatted string.
	modelInfo := fmt.Sprintf(
		ModelFormat,
		info.DisplayName,
		info.BaseModelID,
		info.Version,
		info.Description,
		info.InputTokenLimit,
		info.OutputTokenLimit,
	)

	// Log the compiled model information.
	logger.Any(modelInfo)
}
