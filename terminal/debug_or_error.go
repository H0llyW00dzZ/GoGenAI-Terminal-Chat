package terminal

import (
	"fmt"
	"log"
	"os"
)

type DebugOrErrorLogger struct {
	logger *log.Logger
}

func NewDebugOrErrorLogger() *DebugOrErrorLogger {
	return &DebugOrErrorLogger{
		logger: log.New(os.Stderr, "", log.LstdFlags),
	}
}

func (l *DebugOrErrorLogger) Debug(format string, v ...interface{}) {
	// Check the environment variable to determine if the application is in debug mode
	if os.Getenv(DEBUG_MODE) == "true" {
		l.logger.Printf(format, v...)
	}
}

func (l *DebugOrErrorLogger) Error(format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	colorizedMsg := ColorRed + msg + ColorReset // Apply red color to the entire message
	l.logger.Println(colorizedMsg)
}
