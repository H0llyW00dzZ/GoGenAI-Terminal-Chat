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
	QuitCommand: handleQuitCommand,
	//Note: This Version Command still in development, since it's not sync with AI yet for Prompt such as greeting messages, but it's already working when checking version, so it's not ready for releasing yet.
	VersionCommand: handleCheckVersionCommand,
	//TODO: Will add more commands here, example: :help, :about, :credits, :k8s, syncing AI With Go Routines (Known as Gopher hahaha) etc.
	//Note: In python, I don't think so it's possible hahaahaha, also I am using prefix ":" instead of "/" is respect to git and command line, fuck prefix "/" which is confusing for command line
}

var release GitHubRelease

// This consider stable to avoid memory allocation overhead.
var builder strings.Builder

func init() {
	// Initialize the logger when the package is imported.
	logger = NewDebugOrErrorLogger()

}
