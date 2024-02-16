// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License
//
// Note: This approach addresses the common issue of developers writing Go code in an unnecessarily complex way, particularly for logging operations hahahaha.

package terminal

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

// NewDebugOrErrorLogger initializes a new DebugOrErrorLogger with a logger that writes
// to os.Stderr with the standard log flags.
//
// Returns:
//
//	*DebugOrErrorLogger: A pointer to a newly created DebugOrErrorLogger.
func NewDebugOrErrorLogger() *DebugOrErrorLogger {
	debugMode := os.Getenv(DebugMode) == "true" // Read the environment variable once
	return &DebugOrErrorLogger{
		logger:          log.New(os.Stderr, "", log.LstdFlags),
		debugMode:       debugMode,
		PrintTypingChat: PrintTypingChat,
	}
}

// Debug logs a formatted debug message if the DEBUG_MODE environment variable is set to "true".
// It behaves like Printf and allows for formatted messages.
//
// Parameters:
//
//	format string: The format string for the debug message.
//	v ...interface{}: The values to be formatted according to the format string.
func (l *DebugOrErrorLogger) Debug(format string, v ...interface{}) {
	if l.debugMode {
		// Use strings.Builder for efficient string concatenation
		var builder strings.Builder

		// Format the debug message
		message := fmt.Sprintf(format, v...)
		builder.WriteString(colors.ColorHex95b806)
		builder.WriteString(DEBUGPREFIX)
		builder.WriteString(colors.ColorReset)
		builder.WriteString(" ") // adding back this
		builder.WriteString(message)

		// Simulate typing the debug message
		l.PrintTypingChat(builder.String(), TypingDelay)

		// Print a newline after the message
		printnewlineASCII()
	}
}

// Error logs a formatted error message in red color to signify error conditions.
// It behaves like Println and allows for formatted messages.
//
// Parameters:
//
//	format string: The format string for the error message.
//	v ...interface{}: The values to be formatted according to the format string.
func (l *DebugOrErrorLogger) Error(format string, v ...interface{}) {
	var builder strings.Builder

	// Format the error message
	message := fmt.Sprintf(format, v...)
	builder.WriteString(colors.ColorRed)
	builder.WriteString(message)
	builder.WriteString(colors.ColorReset)

	// Print the error prefix with a timestamp
	PrintPrefixWithTimeStamp(SYSTEMPREFIX, "")

	// Simulate typing the error message
	l.PrintTypingChat(builder.String(), TypingDelay)

	// Print a newline after the message
	printnewlineASCII()
}

// RecoverFromPanic should be deferred at the beginning of a function or goroutine
// to handle any panics that may occur. It logs the panic information with a
// colorized output to distinguish the log message clearly in the terminal.
//
// The message "Recovered from panic:" is displayed in green, followed by the panic
// value in red and the stack trace. This method ensures that the panic does not cause
// the program to crash and provides a clear indication in the logs that a panic was
// caught and handled.
//
// Usage:
//
//	func someFunction() {
//	    logger := terminal.NewDebugOrErrorLogger()
//	    defer logger.RecoverFromPanic()
//
//	    // ... function logic that might panic ...
//	}
//
// It is essential to call this method using defer right after obtaining a logger instance.
// This ensures that it can catch and handle panics from anywhere within the scope of the
// function or goroutine.
func (l *DebugOrErrorLogger) RecoverFromPanic() {
	if r := recover(); r != nil {
		var builder strings.Builder
		combinedStyle := MergeStyles(panicDetected)
		text := "PD"
		asciiArt, _ := ToASCIIArt(text, combinedStyle)
		fmt.Println(asciiArt)
		printnewlineASCII()
		// Format the message for panic
		// Include the application name and version in the panic log
		builder.WriteString(fmt.Sprintf(
			RecoverGopher,
			ApplicationName,
			CurrentVersion,
			colors.ColorHex95b806,
			colors.ColorReset,
			colors.ColorRed,
			r,
			colors.ColorReset))

		// Retrieve the stack trace
		stack := make([]byte, 4096)           // Start with a 4KB buffer
		length := runtime.Stack(stack, false) // Pass 'false' to get only the current goroutine's stack trace

		// Check if the stack trace might be truncated
		// Note: This should not happen for lower complexity. If the complexity is very high, this may occur.
		if length == len(stack) {
			builder.WriteString(fmt.Sprintf(StackPossiblyTruncated))
		}

		builder.WriteString(fmt.Sprintf(StackTracePanic,
			colors.ColorHex95b806,
			colors.ColorReset,
			stack[:length]))

		// Output the message to the logger
		l.PrintTypingChat(builder.String(), TypingDelay)
	}
}

// HandleGoogleAPIError checks for Google API server errors and logs them.
// If it's a server error, it returns true, indicating a retry might be warranted.
//
// Parameters:
//
//	err error: The error returned from a Google API call.
//
// Returns:
//
//	bool: Indicates whether the error is a server error (500).
func (l *DebugOrErrorLogger) HandleGoogleAPIError(err error) bool {
	if err != nil {
		// Check if the error message contains a 500 status code.
		if strings.Contains(err.Error(), Error500GoogleAPI) {
			// Log the Google Internal Error with the error message
			l.Error(ErrorGoogleInternal, err)
			return true // Indicate that this is a server error so bad hahaha
		}
	}
	return false // Not a server error
}

// HandleOtherStupidAPIError checks for non-Google API 500 server errors and logs them.
// If it's a server error, it returns true, indicating a retry might be warranted.
//
// Parameters:
//
//	err error: The error returned from an API call.
//	apiName string: The name of the API for logging purposes.
//
// Returns:
//
//	bool: Indicates whether the error is a bad server error (500).
func (l *DebugOrErrorLogger) HandleOtherStupidAPIError(err error, apiName string) bool {
	if err != nil {
		// Check if the error message contains a 500 status code, but is not from Google API.
		if strings.Contains(err.Error(), Code500) && !strings.Contains(err.Error(), Error500GoogleAPI) {
			// Log the Internal Error with the error message
			l.Error(ErrorOtherAPI, apiName, err)
			return true // Indicate that this is a bad server error
		}
	}
	return false // Not a server error
}

// Info logs a formatted information message. It behaves like Println and allows for formatted messages.
//
// Parameters:
//
//	format string: The format string for the information message.
//	v ...interface{}: The values to be formatted according to the format string.
func (l *DebugOrErrorLogger) Info(format string, v ...interface{}) {
	var builder strings.Builder

	message := fmt.Sprintf(format, v...)
	builder.WriteString(colors.ColorBlue)
	builder.WriteString(message)
	builder.WriteString(colors.ColorReset)

	// Print the message with a timestamp and colored output.
	PrintPrefixWithTimeStamp(SYSTEMPREFIX, "")
	l.PrintTypingChat(builder.String(), TypingDelay)

	// Print a newline after the message
	printnewlineASCII()
}

// Any logs a general message without any colorization. It behaves like Println and allows for formatted messages.
//
// Parameters:
//
//	format string: The format string for the general message.
//	v ...interface{}: The values to be formatted according to the format string.
func (l *DebugOrErrorLogger) Any(format string, v ...interface{}) {
	var builder strings.Builder

	message := fmt.Sprintf(format, v...)
	builder.WriteString(message)

	// Print the message with a timestamp but without any color output.
	PrintPrefixWithTimeStamp(SYSTEMPREFIX, "")
	l.PrintTypingChat(builder.String(), TypingDelay)

	// Print a newline after the message
	printnewlineASCII() // this a modern now instead of fmt hahaha
}
