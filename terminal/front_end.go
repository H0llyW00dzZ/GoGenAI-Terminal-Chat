// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"fmt"
	"time"
)

// PrintPrefixWithTimeStamp prints a message to the terminal, prefixed with a formatted timestamp.
// The timestamp is formatted according to the TimeFormat constant.
//
// For example, with TimeFormat set to "2006/01/02 15:04:05" and the prefix "🤓 You: ",
// the output might be "2024/01/10 16:30:00 🤓 You:".
//
// This function is designed for terminal outputs that benefit from a timestamped context,
// providing clarity and temporal reference for the message displayed.
//
// The prefix parameter is appended to the timestamp and can be a log level, a descriptor,
// or any other string that aids in categorizing or highlighting the message.
func PrintPrefixWithTimeStamp(prefix string) {
	currentTime := time.Now().Format(TimeFormat)
	fmt.Printf(ObjectHighLevelString, currentTime, prefix)
}
