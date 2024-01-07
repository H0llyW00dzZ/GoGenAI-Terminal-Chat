// Copyright (c) 2024 H0llyW00dzZ

package terminal

import (
	"fmt"
	"log"
	"os"
)

// DebugOrErrorLogger provides a simple logger with support for debug and error logging.
// It encapsulates a standard log.Logger and adds functionality for conditional debug
// logging and colorized error output.
type DebugOrErrorLogger struct {
	logger *log.Logger
}

// NewDebugOrErrorLogger initializes a new DebugOrErrorLogger with a logger that writes
// to os.Stderr with the standard log flags.
//
// Returns:
//
//	*DebugOrErrorLogger: A pointer to a newly created DebugOrErrorLogger.
func NewDebugOrErrorLogger() *DebugOrErrorLogger {
	return &DebugOrErrorLogger{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

// Debug logs a formatted debug message if the DEBUG_MODE environment variable is set to "true".
// It behaves like Printf and allows for formatted messages.
//
// Parameters:
//
//	format string: The format string for the debug message.
//	v ...interface{}: The values to be formatted according to the format string.
//
// TODO: Add a DEBUG_MODE constant to the terminal package and use it here.
func (l *DebugOrErrorLogger) Debug(format string, v ...interface{}) {
	// Check the environment variable to determine if the application is in debug mode
	if os.Getenv(DEBUG_MODE) == "true" {
		l.logger.Printf(format, v...)
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
	msg := fmt.Sprintf(format, v...)
	colorizedMsg := ColorRed + msg + ColorReset // Apply red color to the entire message
	l.logger.Println(colorizedMsg)
}
