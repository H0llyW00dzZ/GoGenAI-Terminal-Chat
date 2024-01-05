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

// SingleCharColorize applies ANSI color codes to text surrounded by single-character delimiters.
// It is particularly useful when dealing with text that contains list items or other elements
// that should be highlighted, and it ensures that the colorization is only applied to the
// specified delimiter at the beginning of a line.
//
// Parameters:
//
//	text      string: The text containing elements to be colorized.
//	delimiter string: The single-character delimiter indicating the start of a colorizable element.
//	color     string: The ANSI color code to be applied to the elements starting with the delimiter.
//
// Returns:
//
//	string: The resulting string with colorized elements as specified by the delimiter.
//
// This function handles each line separately and checks for the presence of the delimiter
// at the beginning after trimming whitespace. If the delimiter is found, it colorizes the
// delimiter and the following character (typically a space). The rest of the line remains
// unaltered. If the delimiter is not at the beginning of a line, the line is added to the
// result without colorization.
//
// Note: As with the Colorize function, SingleCharColorize may not function correctly in
// Windows Command Prompt or other environments that do not support ANSI color codes.
// It is best used in terminals that support these codes, such as most Linux/Unix terminals.
func SingleCharColorize(text string, delimiter string, color string) string {
	var result strings.Builder
	lines := strings.Split(text, "\n")
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, string(delimiter)) {
			// Colorize the delimiter and the following space if it's a list item
			result.WriteString(color)
			result.WriteString(string(delimiter))
			result.WriteString(ColorReset)
			result.WriteString(trimmedLine[1:])
		} else {
			// No coloring needed
			result.WriteString(trimmedLine)
		}
		result.WriteString("\n")
	}
	return strings.TrimRight(result.String(), "\n")
}
