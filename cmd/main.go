// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package main

import (
	"os"

	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal"
)

const (
	api_Key  = "API_KEY" // Fixed the typo here
	logFatal = "API_KEY environment variable is not set"
)

// why this so simple ? hahahaha
func main() {
	logger := terminal.NewDebugOrErrorLogger() // Assuming NewDebugOrErrorLogger is exported from the terminal package
	// this goroutines logger panic are not because code of function "terminal package" causing panic, goroutines will tell if there's a panic in other side, indicate that other side system are bad (e.g, too complex).
	defer logger.RecoverFromPanic() // Assuming RecoverFromPanic is exported from the terminal package
	apiKey := os.Getenv(api_Key)
	if apiKey == "" {
		logger.Error(logFatal)
		return // Exit the main function if there's no API key
	}

	session, err := terminal.NewSession(apiKey)
	if err != nil {
		logger.Error("Failed to start session: %v", err)
		return // Exit the main function if session creation fails
	}

	session.Start()
}
