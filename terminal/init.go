// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
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

// filterCodeBlock is a compiled regular expression that is used to identify and
// remove language identifiers from Markdown code blocks. A Markdown code block is
// typically indicated by triple backticks (```) followed by an optional language
// identifier (e.g., ```go). This regular expression matches the pattern of triple
// backticks followed by any sequence of word characters, which represents the
// language identifier. It is used to transform code blocks to a neutral format
// without language hints, which may be desirable for output that does not support
// syntax highlighting or in scenarios where the language identifier is not needed.
//
// The regular expression is compiled once at package initialization for efficiency,
// allowing it to be reused throughout the application without the overhead of
// recompiling it with each use.
var filterCodeBlock *regexp.Regexp

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

var dynamicErrorFileTypeNotSupported = ErrorFileTypeNotSupported

// helper function
//
// verifyFileExtension checks if the file has an allowed extension.
func verifyFileExtension(filePath string) error {
	allowedExtensions := map[string]bool{
		// Note: Feel free to submit a pull request or issues if you want to add support for other file types
		dotMD:   true,
		dotTxt:  true,
		dotPng:  true,
		dotJpg:  true,
		dotJpeg: true,
		dotHeic: true,
		dotHeif: true,
		dotWebp: true,
	}

	// Extract the file extension and check if it's allowed.
	fileExt := strings.ToLower(filepath.Ext(filePath))
	if _, allowed := allowedExtensions[fileExt]; !allowed {
		// Create a slice to hold the allowed extensions for the error message.
		allowedExts := []string{}
		for ext, isAllowed := range allowedExtensions {
			if isAllowed {
				// Add the extension without the dot for a cleaner error message.
				allowedExts = append(allowedExts, strings.TrimPrefix(ext, dotString))
			}
		}
		// Join the allowed extensions with commas and an "or" before the last one.
		allowedExtsStr := strings.Join(allowedExts[:len(allowedExts)-1], dotStringComma) + oRString + allowedExts[len(allowedExts)-1]
		return fmt.Errorf(dynamicErrorFileTypeNotSupported, allowedExtsStr)
	}

	return nil
}

// Dynamic ErrorImageFileTypeNotSupported is a format string for the error message when an unsupported file type is encountered.
var dynamicErrorImageFileTypeNotSupported = ErrorVariableImageFileTypeNotSupported

// helper function
//
// verifyImageFileExtension checks if the image file has an allowed extension.
//
// Note: This marked as todo since currently it unused
func verifyImageFileExtension(filePath string) error {
	allowedExtensions := map[string]bool{
		dotPng:  true,
		dotJpg:  true,
		dotJpeg: true,
		dotHeic: true,
		dotHeif: true,
		dotWebp: true,
	}

	// Extract the file extension and check if it's allowed.
	fileExt := strings.ToLower(filepath.Ext(filePath))
	if _, allowed := allowedExtensions[fileExt]; !allowed {
		// Create a slice to hold the allowed extensions for the error message.
		allowedExts := []string{}
		for ext, isAllowed := range allowedExtensions {
			if isAllowed {
				// Add the extension without the dot for a cleaner error message.
				allowedExts = append(allowedExts, strings.TrimPrefix(ext, dotString))
			}
		}
		// Join the allowed extensions with commas and an "or" before the last one.
		allowedExtsStr := strings.Join(allowedExts[:len(allowedExts)-1], dotStringComma) + oRString + allowedExts[len(allowedExts)-1]
		return fmt.Errorf(dynamicErrorImageFileTypeNotSupported, allowedExtsStr)
	}

	return nil
}

// getImageFormat returns the image format based on the file extension.
func getImageFormat(filePath string) string {
	extToFormat := map[string]string{
		dotJpg:  "jpeg",
		dotJpeg: "jpeg",
		dotPng:  "png",
		dotHeic: "heic",
		dotHeif: "heif",
		dotWebp: "webp",
	}

	// Extract the file extension in lowercase.
	ext := strings.ToLower(filepath.Ext(filePath))
	// Lookup the image format based on file extension.
	if format, ok := extToFormat[ext]; ok {
		return format
	}
	return ""
}

func init() {
	// Initialize the logger when the package is imported.
	logger = NewDebugOrErrorLogger()
	// Compile the ANSI color code regular expression pattern.
	ansiRegex = regexp.MustCompile(BinaryRegexAnsi)
	filterCodeBlock = regexp.MustCompile(CodeBlockRegex)

	// Initialize the command registry.
	// Note: This NewCommandRegistry offers excellent scalability. For Example: You can easily add numerous commands without impacting
	// the AI's performance or synchronization ai, such as `:quit` or `:checkversion`.
	// Additional Note: The scalability of this setup allows the codebase to support a large volume of code, potentially billions of lines.
	// Additionally, this structure simplifies maintenance by reducing the complexity often associated with individual functions, such as numerous 'if', 'for', 'case' statements, and '&&' or '||' operators.
	registry = NewCommandRegistry()
	registry.Register(QuitCommand, &handleQuitCommand{})
	registry.Register(ShortQuitCommand, &handleQuitCommand{})
	registry.Register(VersionCommand, &handleCheckVersionCommand{})
	registry.Register(HelpCommand, &handleHelpCommand{})
	registry.Register(ShortHelpCommand, &handleHelpCommand{})
	registry.Register(AITranslateCommand, &handleAITranslateCommand{})
	registry.Register(SummarizeCommands, &handleSummarizeCommand{})
	// Assume handleClearCommand is capable of handling subcommands for ":clear"
	clearCommandHandler := &handleClearCommand{}
	registry.Register(ClearCommand, &handleClearCommand{})
	// Register subcommands for ":clear"
	// Note: These subcommands are as scalable as the `NewCommandRegistry`.
	registry.RegisterSubcommand(ClearCommand, ChatCommands, clearCommandHandler)
	registry.RegisterSubcommand(ClearCommand, SummarizeCommands, clearCommandHandler)
	// Assume handleStatsCommand is capable of handling subcommands for ":stats"
	statsCommandHandler := &handleStatsCommand{}
	registry.Register(StatsCommand, &handleStatsCommand{})
	registry.RegisterSubcommand(StatsCommand, ChatCommands, statsCommandHandler)
	// Assume handleCryptoRandCommand is capable of handling subcommands for ":cryptorand"
	cryptoRandCommandHandler := &handleCryptoRandCommand{}
	registry.Register(CryptoRandCommand, &handleCryptoRandCommand{})
	registry.RegisterSubcommand(CryptoRandCommand, LengthArgs, cryptoRandCommandHandler)
	// Assume safetySettingsCommandHandler is capable of handling subcommands for ":safety"
	safetySettingsCommandHandler := &handleSafetyCommand{}
	registry.Register(SafetyCommand, &handleSafetyCommand{})
	registry.RegisterSubcommand(SafetyCommand, Low, safetySettingsCommandHandler)
	registry.RegisterSubcommand(SafetyCommand, Default, safetySettingsCommandHandler)
	registry.RegisterSubcommand(SafetyCommand, High, safetySettingsCommandHandler)
	// Assume showChatCommandHandler is capable of handling subcommands for ":chat"
	showChatCommandHandler := &handleShowChatCommand{}
	registry.Register(ChatCommands, &handleShowChatCommand{})
	registry.RegisterSubcommand(ChatCommands, ShowCommands, showChatCommandHandler)
	// Register the token count command and its handler.
	tokenCountCommandHandler := &handleTokeCountingCommand{}
	registry.Register(TokenCountCommands, tokenCountCommandHandler)
	registry.RegisterSubcommand(TokenCountCommands, FileCommands, tokenCountCommandHandler)

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
		TIP,
		BLANK_,
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
