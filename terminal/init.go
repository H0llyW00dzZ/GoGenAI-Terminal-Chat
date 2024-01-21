// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"regexp"
)

// apiKey holds the API key used for authenticating requests to the generative
// AI service. It should be initialized with a valid API key before making any
// requests that require authentication.
//
// Note: Storing API keys in source code is not recommended due to security
// concerns. It is better to use environment variables or secure storage mechanisms
// to handle sensitive information such as API keys.
var apiKey string

// logger is a package-level variable that can be used throughout the terminal package.
var logger *DebugOrErrorLogger

// this is a package-level variable that holds the command registry.
// Caution is advised: if you're not familiar with these practices, improper handling in this "CommandRegistry" could lead to frequent panics 24/7 🤪.
var registry *CommandRegistry

// checkVersion is a package-level variable that holds the latest release information
// fetched from the GitHub API. It is used to cache the details of the latest release
// to avoid multiple API calls when checking for updates within the application's
// lifecycle. This variable should be updated only through the CheckLatestVersion function.
var checkVersion GitHubRelease

var aiPrompt string

// colors holds the ANSI color codes and is accessible throughout the package.
var colors = ANSIColorCodes{
	ColorRed:         ColorRed,
	ColorGreen:       ColorGreen,
	ColorYellow:      ColorYellow,
	ColorBlue:        ColorBlue,
	ColorPurple:      ColorPurple,
	ColorCyan:        ColorCyan,
	ColorHex95b806:   ColorHex95b806,
	ColorCyan24Bit:   ColorCyan24Bit,
	ColorPurple24Bit: ColorPurple24Bit,
	ColorReset:       ColorReset,
}

// ansichar
var ansichar = BinaryAnsiChars{
	BinaryAnsiChar:          BinaryAnsiChar,
	BinaryAnsiSquenseChar:   BinaryAnsiSquenseChar,
	BinaryAnsiSquenseString: BinaryAnsiSquenseString,
	BinaryLeftSquareBracket: BinaryLeftSquareBracket,
}

// animatedchars
var humantyping = TypingChars{
	AnimatedChars: AnimatedChars,
}

// newLineChar
var nl = NewLineChar{
	NewLineChars: NewLineChars,
}

// totalTokenCount is a package-level variable that holds the total number of tokens
var totalTokenCount int = 0

// ansiRegex is a compiled regular expression that matches ANSI color codes.
// It is compiled once when the package is initialized.
// Note: Removing Struct now, this a `Go` not a `Rust`
var ansiRegex *regexp.Regexp

var tripleBacktickColor string

// scalable safetyOptions maps safety level strings to their corresponding setter functions and validity.
var safetyOptions = map[string]SafetyOption{
	Low: {
		Setter: func(s *SafetySettings) { s.SetLowSafety() },
		Valid:  true,
	},
	High: {
		Setter: func(s *SafetySettings) { s.SetHighSafety() },
		Valid:  true,
	},
	Default: {
		Setter: func(s *SafetySettings) { *s = *DefaultSafetySettings() },
		Valid:  true,
	},
}

// scalable a global variable for the ASCII style.
var slantStyle = NewASCIIArtStyle()
var stripStyle = NewASCIIArtStyle()
var newLine = NewASCIIArtStyle()

func init() {
	// Initialize the logger when the package is imported.
	logger = NewDebugOrErrorLogger()
	// Compile the ANSI color code regular expression pattern.
	ansiRegex = regexp.MustCompile(BinaryRegexAnsi)

	// Initialize the command registry.
	// Note: This NewCommandRegistry offers excellent scalability. For Example: You can easily add numerous commands without impacting
	// the AI's performance or synchronization ai, such as `:quit` or `:checkversion`.
	registry = NewCommandRegistry()
	registry.Register(QuitCommand, &handleQuitCommand{})
	registry.Register(ShortQuitCommand, &handleQuitCommand{})
	registry.Register(VersionCommand, &handleCheckVersionCommand{})
	registry.Register(HelpCommand, &handleHelpCommand{})
	registry.Register(ShortHelpCommand, &handleHelpCommand{})
	registry.Register(ClearCommand, &handleClearCommand{})
	registry.Register(SafetyCommand, &handleSafetyCommand{})
	registry.Register(AITranslateCommand, &handleAITranslateCommand{})
	registry.Register(CryptoRandCommand, &handleCryptoRandCommand{})
	//TODO: Will add more commands here, example: :help, :about, :credits, :k8s, syncing AI With Go Routines (Known as Gopher hahaha) etc.
	// Note: In python, I don't think so it's possible hahaahaha, also I am using prefix ":" instead of "/" is respect to git and command line, fuck prefix "/" which is confusing for command line

	// ASCII Scalable mode 🤪
	// Additional Note: let me know if scalable logic like this are possible in other language that can avoid complexity "overhead" for human and the machine hahaahahahaa
	slantStyle.AddChar(G, []string{
		_G,
		_O,
		_GEN,
		A_,
		I_,
	}, BoldText+colors.ColorHex95b806)
	// Initialize the text style patterns for 'V'.
	slantStyle.AddChar(V, []string{
		BLANK_,
		BLANK_,
		BLANK_, // TODO: Implement a notification to be displayed here when a new version is available.
		Current_Version,
		Copyright,
		// Note: This utilizes a struct for color definitions to ensure consistency. This is important for compatibility with operating systems that may not handle ANSI colors properly.
	}, BoldText+colors.ColorCyan24Bit)
	// Initialize the text style patterns for 'V'.
	// Note: Although it may appear that patterns are duplicated, they are, in fact, distinct. This structure ensures scalability that powered by Go 🤪.
	stripStyle.AddChar(V, []string{
		eMpty, // this better unlike hardcoded "\n" lmao.
		StripChars,
		// Note: This utilizes a struct for color definitions to ensure consistency. This is important for compatibility with operating systems that may not handle ANSI colors properly.
	}, BoldText+colors.ColorCyan24Bit)
	newLine.AddChar(N, []string{
		eMpty, // this better unlike hardcoded "\n" lmao.
	}, colors.ColorReset)
}
