// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"context"
	"log"
	"sync"
	"time"

	genai "github.com/google/generative-ai-go/genai"
)

// ASCIIArtChar represents a styled ASCII character with its pattern and color.
// The Pattern field is a slice of strings, with each string representing a line
// of the character's ASCII representation. The Color field specifies the color
// to be used when displaying the character.
type ASCIIArtChar struct {
	Pattern []string // Lines of the ASCII representation of the character.
	Color   string   // Color code or label for the character's color.
}

// ANSIColorCodes defines a struct for holding ANSI color escape sequences.
type ANSIColorCodes struct {
	ColorRed         string
	ColorGreen       string
	ColorYellow      string
	ColorBlue        string
	ColorPurple      string
	ColorCyan        string
	ColorHex95b806   string // 24-bit color
	ColorCyan24Bit   string // 24-bit color
	ColorPurple24Bit string // 24-bit color
	ColorReset       string
}

// BinaryAnsiChars is a struct that contains the ANSI characters used to print the typing effect.
type BinaryAnsiChars struct {
	BinaryAnsiChar          rune
	BinaryAnsiSquenseChar   rune
	BinaryAnsiSquenseString string
	BinaryLeftSquareBracket rune
}

// ChatHistory manages the state of chat messages exchanged during a session.
// It tracks the messages, their unique hashes, and counts of different types of messages (user, AI, system).
// This struct also ensures concurrent access safety using a read-write mutex.
type ChatHistory struct {
	Messages           []string       // Messages contains all the chat messages in chronological order.
	Hashes             map[string]int // Hashes maps the SHA-256 hash of each message to its index in Messages.
	UserMessageCount   int            // UserMessageCount holds the total number of user messages.
	AIMessageCount     int            // AIMessageCount holds the total number of AI messages.
	SystemMessageCount int            // SystemMessageCount holds the total number of system messages.
	mu                 sync.RWMutex   // Explicit 🤪
}

// ChatConfig encapsulates settings that affect the management of chat history
// during a session with the generative AI. It determines the amount of chat history
// retained in memory and the portion of that history used to provide context to the AI.
type ChatConfig struct {
	// HistorySize specifies the total number of chat messages to retain in the session's history.
	// This helps in limiting the memory footprint and ensures that only recent interactions
	// are considered for maintaining context.
	HistorySize int

	// HistorySendToAI indicates the number of recent messages from the history to be included
	// when sending context to the AI. This allows the AI to generate responses that are
	// relevant to the current conversation flow without being overwhelmed by too much history.
	HistorySendToAI int
}

// ChatWorker is responsible for handling background tasks related to chat sessions.
type ChatWorker struct {
	session *Session
	ticker  *time.Ticker
	done    chan bool
}

// ColorizationOptions encapsulates the settings necessary for the Colorize function to apply color to text.
// It includes the text to be colorized, pairs of delimiters and their corresponding colors, a map to determine
// if delimiters should be retained, and a map defining additional ANSI formatting codes.
type ColorizationOptions struct {
	Text           string            // Text is the original string that will be colorized.
	ColorPairs     []string          // ColorPairs is a slice where each pair consists of a delimiter and its associated ANSI color code.
	KeepDelimiters map[string]bool   // KeepDelimiters dictates whether each delimiter should be kept in the final output.
	Formatting     map[string]string // Formatting is a map associating delimiters with their ANSI formatting codes for additional text styling.
}

// ColorizationPartOptions holds the options for colorizing parts of a text.
// It contains the text to be colorized, the delimiter that marks the text to be colorized,
// the color to apply, and maps that determine whether to keep delimiters and how to format the text.
type ColorizationPartOptions struct {
	Text           string            // Text is the string that will be processed for colorization.
	Delimiter      string            // Delimiter is the string that marks the beginning and end of text to be colorized.
	Color          string            // Color is the ANSI color code that will be applied to the text within delimiters.
	KeepDelimiters map[string]bool   // KeepDelimiters is a map indicating whether to keep each delimiter in the output text.
	Formatting     map[string]string // Formatting is a map of delimiters to their corresponding ANSI formatting codes.
}

// DebugOrErrorLogger provides a simple logger with support for debug and error logging.
// It encapsulates a standard log.Logger and adds functionality for conditional debug
// logging and colorized error output.
type DebugOrErrorLogger struct {
	logger          *log.Logger
	debugMode       bool
	PrintTypingChat func(string, time.Duration)
}

// FormattingOptions encapsulates the settings required to apply formatting to a text segment.
// It includes the text to be formatted, the delimiter used to identify the text segment,
// the color to apply, and a map specifying the ANSI formatting codes associated with different delimiters.
type FormattingOptions struct {
	Text       string            // Text is the string to which formatting will be applied.
	Delimiter  string            // Delimiter is the string used to identify segments of text for formatting.
	Color      string            // Color is the ANSI color code that will be applied to the text within the delimiter.
	Formatting map[string]string // Formatting is a map that associates delimiters with their respective ANSI formatting codes.
}

// TypingChars is a struct that contains the Animated Chars used to print the typing effect.
type TypingChars struct {
	AnimatedChars string
}

// GitHubRelease represents the metadata of a software release from GitHub.
// It includes information such as the tag name, release name, and a description body,
// typically containing the changelog or release notes.
type GitHubRelease struct {
	TagName string `json:"tag_name"`     // The tag associated with the release, e.g., "v1.2.3"
	Name    string `json:"name"`         // The official name of the release
	Body    string `json:"body"`         // Detailed description or changelog for the release
	Date    string `json:"published_at"` // Published Date
}

// MessageStats encapsulates the counts of different types of messages in the chat history.
// It holds separate counts for user messages, AI messages, and system messages.
type MessageStats struct {
	UserMessages   int // UserMessages is the count of messages sent by users.
	AIMessages     int // AIMessages is the count of messages sent by the AI.
	SystemMessages int // SystemMessages is the count of system-generated messages.
}

// Session encapsulates the state and functionality for a chat session with a generative AI model.
// It holds the AI client, chat history, and context for managing the session lifecycle.
type Session struct {
	Client         *genai.Client      // Client is the generative AI client used to communicate with the AI model.
	ChatHistory    *ChatHistory       // ChatHistory stores the history of the chat session.
	ChatConfig     *ChatConfig        // ChatConfig contains the settings for managing the chat history size.
	Ctx            context.Context    // Ctx is the context governing the session, used for cancellation.
	Cancel         context.CancelFunc // Cancel is a function to cancel the context, used for cleanup.
	Ended          bool               // Ended indicates whether the session has ended.
	SafetySettings *SafetySettings    // Holds the current safety settings for the session.
	// mu protects the concurrent access to session's state, ensuring thread safety.
	// It should be locked when accessing or modifying the session's state.
	mu sync.Mutex
	// this reference pretty useful, which can handle runtime 24/7, unlike original ai chat session systems.
	// for example, if session is ended not cause of client, then it will be renew with previous chat history.
	lastInput string // Stores the last user input for reference

}

// SafetyOption is a function type that takes a pointer to a SafetySettings
// instance and applies a specific safety configuration to it. It is used
// to abstract the different safety level settings (e.g., low, high, default)
// and allows for a flexible and scalable way to manage safety configurations
// through function mapping.
type SafetyOption struct {
	Setter func(s *SafetySettings)
	Valid  bool
}

// SafetySettings encapsulates the content safety configuration for the AI model.
// It defines thresholds for various categories of potentially harmful content,
// allowing users to set the desired level of content filtering based on the
// application's requirements and user preferences.
type SafetySettings struct {
	// DangerousContentThreshold defines the threshold for filtering dangerous content.
	DangerousContentThreshold genai.HarmBlockThreshold
	// HarassmentContentThreshold defines the threshold for filtering harassment-related content.
	HarassmentContentThreshold genai.HarmBlockThreshold
	// SexuallyExplicitContentThreshold defines the threshold for filtering sexually explicit content.
	SexuallyExplicitContentThreshold genai.HarmBlockThreshold
	// MedicalThreshold defines the threshold for filtering medical-related content.
	MedicalThreshold genai.HarmBlockThreshold
	// ViolenceThreshold defines the threshold for filtering violent content.
	ViolenceThreshold genai.HarmBlockThreshold
	// HateSpeechThreshold defines the threshold for filtering hate speech.
	HateSpeechThreshold genai.HarmBlockThreshold
	// ToxicityThreshold defines the threshold for filtering toxic content.
	ToxicityThreshold genai.HarmBlockThreshold
	// DerogatoryThershold defines the threshold for filtering derogatory content.
	DerogatoryThershold genai.HarmBlockThreshold
}

// TokenCountParams encapsulates the parameters needed for counting tokens using a generative AI model.
type TokenCountParams struct {
	// Authentication key for the AI service.
	APIKey string
	// Name of the AI model to use.
	ModelName string
	// Text input for token counting.
	Input string
	// Format of the image if provided (e.g., "png", "jpeg").
	ImageFormat string
	// Image data as a byte slice.
	ImageData []byte
}

// NewLineChar is a struct that containt Rune for New Line Character
type NewLineChar struct {
	NewLineChars rune
}
