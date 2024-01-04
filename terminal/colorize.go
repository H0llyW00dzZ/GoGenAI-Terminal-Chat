// Copyright (c) 2024 H0llyW00dzZ

package terminal

import (
	"strings"
)

// ANSI color codes
const (
	ColorRed    = "\033[31m"
	ColorGreen  = "\033[32m"
	ColorYellow = "\033[33m"
	ColorBlue   = "\033[34m"
	ColorPurple = "\033[35m"
	ColorCyan   = "\033[36m"
	ColorReset  = "\033[0m"
)

// Colorize applies ANSI color codes to the text surrounded by specified delimiters.
// It can process multiple delimiters, each with a corresponding color. The function
// can also conditionally retain or remove the delimiters in the final output.
//
// Parameters:
//
//	text          string: The text to be colorized.
//	colorPairs    []string: A slice where each pair of elements represents a delimiter and its color.
//	keepDelimiters map[string]bool: A map to indicate whether to keep the delimiter in the output.
//
// Returns:
//
//	string: The colorized text.
//
// Note: This function may not work as expected in Windows Command Prompt due to its limited
// support for ANSI color codes. It is designed for terminals that support ANSI, such as those
// in Linux/Unix environments.
func Colorize(text string, colorPairs []string, keepDelimiters map[string]bool) string {
	for i := 0; i < len(colorPairs); i += 2 {
		delimiter := colorPairs[i]
		color := colorPairs[i+1]
		parts := strings.Split(text, delimiter)
		for j := 1; j < len(parts); j += 2 {
			if keep, exists := keepDelimiters[delimiter]; exists && keep {
				parts[j] = color + delimiter + parts[j] + delimiter + ColorReset
			} else {
				parts[j] = color + parts[j] + ColorReset
			}
		}
		text = strings.Join(parts, "")
	}
	return text
}
