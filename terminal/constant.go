// Copyright (c) 2024 H0llyW00dzZ

package terminal

import "time"

// Defined constants for the terminal package
const (
	SignalMessage = "\nReceived an interrupt, shutting down gracefully..."
	StripChars    = "---"
	NewLineChars  = '\n'
	// this animated chars is magic, it used to show the user that the AI is typing just like human would type
	AnimatedChars = "%c"
	// this model is subject to changed in future
	ModelAi     = "gemini-pro"
	TypingDelay = 100 * time.Millisecond
)

// Defined constants for language
const (
	YouNerd               = "🤓 You: "
	AiNerd                = "🤖 AI: "
	ContextPrompt         = "Hello! How can I assist you today?"
	ShutdownMessage       = "Shutting down gracefully..."
	UnknownCommand        = "Unknown command."
	ContextPromptShutdown = "The user has issued a quit command. Please provide a shutdown message as you are Assistant."
)

// Defined constants for commands
// Note: will add more in future based on the need
const (
	QuitCommand = ":quit"
	PrefixChar  = ":"
)

// Defined List error message
const (
	ErrorGettingShutdownMessage = "Error getting shutdown message from AI:"
)
