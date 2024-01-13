// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import "strings"

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

// commandHandlers maps command strings to their corresponding handler functions.
// This allows for a scalable and maintainable way to manage chat commands.
var commandHandlers = map[string]CommandHandler{
	// Note: This map offers excellent scalability. For Example: You can easily add numerous commands without impacting
	// the AI's performance or synchronization ai, such as `:quit` or `:checkversion`.
	QuitCommand:    &handleQuitCommand{},
	VersionCommand: &handleCheckVersionCommand{},
	HelpCommand:    &handleHelpCommand{},
	//TODO: Will add more commands here, example: :help, :about, :credits, :k8s, syncing AI With Go Routines (Known as Gopher hahaha) etc.
	//Note: In python, I don't think so it's possible hahaahaha, also I am using prefix ":" instead of "/" is respect to git and command line, fuck prefix "/" which is confusing for command line
}

// checkVersion is a package-level variable that holds the latest release information
// fetched from the GitHub API. It is used to cache the details of the latest release
// to avoid multiple API calls when checking for updates within the application's
// lifecycle. This variable should be updated only through the CheckLatestVersion function.
var checkVersion GitHubRelease

// This consider stable to avoid memory allocation overhead.
var buildeR strings.Builder

var aiPrompt string

// colors holds the ANSI color codes and is accessible throughout the package.
var colors = ANSIColorCodes{
	ColorRed:       "\x1b[31m",
	ColorGreen:     "\x1b[32m",
	ColorYellow:    "\x1b[33m",
	ColorBlue:      "\x1b[34m",
	ColorPurple:    "\x1b[35m",
	ColorCyan:      "\x1b[36m",
	ColorHex95b806: "\x1b[38;2;149;184;6m",
	ColorCyan24Bit: "\x1b[38;2;17;240;247m",
	ColorReset:     "\x1b[0m",
}

func init() {
	// Initialize the logger when the package is imported.
	logger = NewDebugOrErrorLogger()

}
