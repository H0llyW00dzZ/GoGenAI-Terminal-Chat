// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"strings"
)

// ANSIColorCodes defines a struct for holding ANSI color escape sequences.
type ANSIColorCodes struct {
	ColorRed         string
	ColorGreen       string
	ColorYellow      string
	ColorBlue        string
	ColorPurple      string
	ColorCyan        string
	ColorHex95b806   string // 24-bit color
	ColorCyan24Bit   string // 24-bit color
	ColorPurple24Bit string // 24-bit color
	ColorReset       string
}

// ANSI color codes
const (
	// Note: By replacing the ANSI escape sequence from "\033" to "\x1b", might can avoid a rare bug that sometimes occurs on different machines,
	// although the original code works fine on mine (Author: @H0llyW00dzZ).
	ColorRed    = "\x1b[31m"
	ColorGreen  = "\x1b[32m"
	ColorYellow = "\x1b[33m"
	ColorBlue   = "\x1b[34m"
	ColorPurple = "\x1b[35m"
	ColorCyan   = "\x1b[36m"
	// ColorHex95b806 represents the color #95b806 using an ANSI escape sequence for 24-bit color.
	ColorHex95b806 = "\x1b[38;2;149;184;6m"
	// ColorCyan24Bit represents the color #11F0F7 using an ANSI escape sequence for 24-bit color.
	ColorCyan24Bit   = "\x1b[38;2;17;240;247m"
	ColorPurple24Bit = "\x1b[38;2;255;0;255m"
	ColorReset       = "\x1b[0m"
)

// ANSI Text Formatting.
const (
	// bold text.
	BoldText = "\x1b[1m"
	// reset text formatting.
	ResetText = "\x1b[22m"
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
//	formatting     map[string]string: A map of delimiters to their corresponding ANSI formatting codes.
//
// Returns:
//
//	string: The colorized text.
//
// Note: This function may not work as expected in Windows Command Prompt due to its limited
// support for ANSI color codes. It is designed for terminals that support ANSI, such as those
// in Linux/Unix environments.
func Colorize(text string, colorPairs []string, keepDelimiters map[string]bool, formatting map[string]string) string {
	tripleBacktickPlaceholder := ObjectTripleHighLevelString
	text = replaceTripleBackticks(text, tripleBacktickPlaceholder)

	for i := 0; i < len(colorPairs); i += 2 {
		delimiter := colorPairs[i]
		color := colorPairs[i+1]
		if delimiter == TripleBacktick {
			tripleBacktickColor = color
		}
		text = processDelimiters(text, delimiter, color, keepDelimiters, formatting)
	}

	if tripleBacktickColor != "" {
		// Append the ResetColor after the colorizedTripleBacktick to reset color formatting
		colorizedTripleBacktick := tripleBacktickColor + TripleBacktick + colors.ColorReset
		text = strings.Replace(text, tripleBacktickPlaceholder, colorizedTripleBacktick, -1)
	}

	return text
}

// replaceTripleBackticks replaces all occurrences of triple backticks with a placeholder.
func replaceTripleBackticks(text, placeholder string) string {
	for {
		index := strings.Index(text, TripleBacktick)
		if index == -1 {
			break
		}
		text = strings.Replace(text, TripleBacktick, placeholder, 1)
	}
	return text
}

// ApplyFormatting applies text formatting based on the provided delimiter.
// If the delimiter is recognized, it applies the appropriate ANSI formatting codes.
//
// Parameters:
//
//	text string: The text to format.
//	delimiter string: The delimiter that indicates what kind of formatting to apply.
//	color string: The ANSI color code to apply to the text.
//	formatting map[string]string: A map of delimiters to their corresponding ANSI formatting codes.
//
// Returns:
//
//	string: The formatted text.
func ApplyFormatting(text string, delimiter string, color string, formatting map[string]string) string {
	if formatCode, ok := formatting[delimiter]; ok {
		return color + formatCode + text + ResetText + ColorReset
	}
	return color + text + ColorReset
}

// processDelimiters processes the delimiters in the text and applies the corresponding color and formatting.
func processDelimiters(text string, delimiter, color string, keepDelimiters map[string]bool, formatting map[string]string) string {
	parts := strings.Split(text, delimiter)
	for j := 1; j < len(parts); j += 2 {
		if keep, exists := keepDelimiters[delimiter]; exists && keep {
			parts[j] = color + delimiter + parts[j] + delimiter + ColorReset
		} else {
			parts[j] = ApplyFormatting(parts[j], delimiter, color, formatting)
		}
	}
	return strings.Join(parts, "")
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
	//Note: This variable result are not possible to register it in the init.go, because it's used to be avoid the duplicate, so better keep like this.
	var result strings.Builder
	lines := strings.Split(text, StringNewLine)
	for _, line := range lines {
		trimmedLine := strings.TrimSpace(line)
		if strings.HasPrefix(trimmedLine, delimiter) {
			// Colorize the delimiter and the following space if it's a list item
			result.WriteString(color)
			result.WriteString(trimmedLine[:1])
			result.WriteString(colors.ColorReset)
			result.WriteString(trimmedLine[1:])
		} else {
			// No coloring needed
			result.WriteString(trimmedLine)
		}
		result.WriteRune(nl.NewLineChars)
	}
	return strings.TrimRight(result.String(), StringNewLine)
}
