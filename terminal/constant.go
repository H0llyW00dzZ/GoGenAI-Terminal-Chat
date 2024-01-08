// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import "time"

// Defined List of GitHub API
const (
	// GitHubAPIURL is the endpoint for the latest release information of the application.
	GitHubAPIURL = "https://api.github.com/repos/H0llyW00dzZ/GoGenAI-Terminal-Chat/releases/latest"
	// CurrentVersion represents the current version of the application.
	CurrentVersion = "v0.2.1"
)

// Defined constants for the terminal package
const (
	SignalMessage = "\nReceived an interrupt, shutting down gracefully..."
	StripChars    = "---"
	NewLineChars  = '\n'
	// this animated chars is magic, it used to show the user that the AI is typing just like human would type
	AnimatedChars = "%c"
	// this model is subject to changed in future
	ModelAi = "gemini-pro"
	// this may subject to changed in future for example can customize the delay
	TypingDelay = 60 * time.Millisecond
)

// Defined constants for language
const (
	YouNerd                = "🤓 You: "
	AiNerd                 = "🤖 AI: "
	ContextPrompt          = "Hello! How can I assist you today?"
	ShutdownMessage        = "Shutting down gracefully..."
	UnknownCommand         = "Unknown command."
	ContextPromptShutdown  = "The user has issued a quit command. Please provide a shutdown message as you are Assistant."
	ContextCancel          = "Context canceled, shutting down..." // sending a messages to gopher officer
	ANewVersionIsAvailable = "A newer version is available: %s\n"
	YouAreusingLatest      = "You are using the latest version of GoGenAI Terminal Chat."
)

// Defined constants for commands
//
// Note: will add more in future based on the need,
// for example, to change the model, or to change the delay, another thing is syncing ai with goroutine (known as gopher)
const (
	QuitCommand    = ":quit"
	VersionCommand = ":version"
	PrefixChar     = ":"
)

// Defined List error message
const (
	ErrorGettingShutdownMessage         = "Error getting shutdown message from AI: %v"
	ErrorHandlingCommand                = "Error handling command: %v"
	ErrorCountingTokens                 = "Error counting tokens: %v"
	ErrorSendingMessage                 = "Error sending message to AI: %v"
	ErrorReadingUserInput               = "Error reading user input: %v"
	ErrorFailedToFetchReleaseInfo       = "Failed to fetch the latest release info: %v"
	ErrorReceivedNon200StatusCode       = "Received non-200 status code: %v"
	ErrorFailedToReadTheResponseBody    = "Failed to read the response body: %v"
	ErrorFaileduUnmarshalTheReleaseData = "Failed to unmarshal the release data: %v"
)

// Defined List of characters
const (
	SingleAsterisk = "*"
	DoubleAsterisk = "**"
	SingleBacktick = "`"
	StringNewLine  = "\n"
)

// Defined List of Environment variables
const (
	DEBUG_MODE = "DEBUG_MODE"
)
