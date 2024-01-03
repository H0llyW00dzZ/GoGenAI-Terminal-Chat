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
	YouNerd       = "🤓 You: "
	AiNerd        = "🤖 AI: "
	ContextPrompt = "Hello! How can I assist you today?"
)
