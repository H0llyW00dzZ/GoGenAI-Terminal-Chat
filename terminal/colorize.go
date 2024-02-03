// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"strings"
)

// Colorize applies ANSI color codes to the text surrounded by specified delimiters.
// It can process multiple delimiters, each with a corresponding color. The function
// can also conditionally retain or remove the delimiters in the final output.
//
// Parameters:
//
//	options ColorizationOptions: A struct containing all the necessary options, including:
//		- Text: The text to be colorized.
//		- ColorPairs: A slice where each pair of elements represents a delimiter and its color.
//		- KeepDelimiters: A map to indicate whether to keep the delimiter in the output.
//		- Formatting: A map of delimiters to their corresponding ANSI formatting codes.
//
// Returns:
//
//	string: The colorized text.
//
// Note: This function may not work as expected in Windows Command Prompt due to its limited
// support for ANSI color codes. It is designed for terminals that support ANSI, such as those
// in Linux/Unix environments.
func Colorize(options ColorizationOptions) string {
	text := strings.ReplaceAll(options.Text, TripleBacktick, ObjectTripleHighLevelString)

	var result strings.Builder
	result.Grow(len(text) * 2) // Preallocate with an estimated size

	// Assume tripleBacktickColor is defined elsewhere or add it to ColorizationOptions if needed
	var tripleBacktickColor string

	// Process each color pair separately
	for i := 0; i < len(options.ColorPairs); i += 2 {
		delimiter := options.ColorPairs[i]
		color := options.ColorPairs[i+1]
		if delimiter == TripleBacktick {
			// Set the color for triple backticks
			tripleBacktickColor = color
		}
		colorizationPartOptions := ColorizationPartOptions{
			Text:           text,
			Delimiter:      delimiter,
			Color:          color,
			KeepDelimiters: options.KeepDelimiters,
			Formatting:     options.Formatting,
		}
		text = applyColorToDelimitedText(colorizationPartOptions)

		// Create a FormattingOptions struct for the processDelimiters call
		formattingOptions := FormattingOptions{
			Text:       text,
			Delimiter:  delimiter,
			Color:      color,
			Formatting: options.Formatting,
		}
		text = processDelimiters(formattingOptions, options.KeepDelimiters)
	}

	result.WriteString(text)
	processedText := result.String()

	// Replace the placeholder with the colorized triple backtick sequence
	// Note: This can refactor easily, for example changing color inside a triple backtick
	if tripleBacktickColor != "" {
		colorizedTripleBacktick := tripleBacktickColor + TripleBacktick + ColorReset
		processedText = strings.ReplaceAll(processedText,
			ObjectTripleHighLevelString,
			colorizedTripleBacktick)
	}

	return processedText
}

// applyColorToDelimitedText applies the specified color to delimited sections of the given text.
func applyColorToDelimitedText(options ColorizationPartOptions) string {
	var result strings.Builder
	parts := strings.Split(options.Text, options.Delimiter)
	partsLen := len(parts) // Get the length of parts once and pass it to processPart

	// Process parts with a consistent pattern to avoid complex conditionals
	for i, part := range parts {
		processPart(&result, i, partsLen, part, options)
	}
	return result.String()
}

// processPart processes an individual part of the text, applying color if necessary.
func processPart(result *strings.Builder, index, partsLen int, part string, options ColorizationPartOptions) {
	if index%2 == 0 { // Even index, regular text
		result.WriteString(part)
	} else { // Odd index, colorized text
		colorizePart(result, part, options)
	}
	appendDelimiterIfNeeded(result, index, partsLen, options)
}

// colorizePart applies color and formatting to a part of the text.
func colorizePart(result *strings.Builder, part string, options ColorizationPartOptions) {
	// Apply any formatting (bold, italic, etc.) before the color
	if format, hasFormat := options.Formatting[options.Delimiter]; hasFormat {
		result.WriteString(format)
	}
	// Apply the color
	result.WriteString(options.Color)
	// Append the actual text
	result.WriteString(part)
	// Reset the color first
	result.WriteString(ColorReset)
	// Reset any formatting (bold, italic, etc.) if it was applied
	if _, hasFormat := options.Formatting[options.Delimiter]; hasFormat {
		result.WriteString(ResetBoldText)
		result.WriteString(ResetItalicText)
	}
}

// appendDelimiterIfNeeded appends the delimiter to the result if the conditions are met.
func appendDelimiterIfNeeded(result *strings.Builder, index, partsLen int, options ColorizationPartOptions) {
	if shouldKeepDelimiter(options.Delimiter,
		options.KeepDelimiters) &&
		index < partsLen-1 {
		result.WriteString(options.Delimiter)
	}
}

// shouldKeepDelimiter checks if a delimiter should be kept in the final result.
func shouldKeepDelimiter(delimiter string, keepDelimiters map[string]bool) bool {
	keep, exists := keepDelimiters[delimiter]
	return exists && keep
}

// ApplyFormatting applies text formatting based on the provided FormattingOptions.
// If the delimiter is recognized, it applies the appropriate ANSI formatting codes.
//
// Parameters:
//
//	options FormattingOptions: The struct that contains the formatting options.
//
// Returns:
//
//	string: The formatted text.
func ApplyFormatting(options FormattingOptions) string {
	if formatCode, ok := options.Formatting[options.Delimiter]; ok {
		return options.Color + formatCode +
			options.Text + ResetBoldText +
			ResetItalicText + ColorReset
	}
	return options.Color + options.Text + ColorReset
}

// processDelimiters processes the delimiters in the text and applies the corresponding color and formatting.
// It takes a FormattingOptions struct containing the text to process and formatting details,
// and a map that dictates whether to keep or remove each delimiter after processing.
//
// Parameters:
//
//	options FormattingOptions: The struct that contains the text and formatting details.
//	keepDelimiters map[string]bool: A map indicating whether to keep each delimiter in the output.
//
// Returns:
//
//	string: The text with delimiters processed and formatting applied.
func processDelimiters(options FormattingOptions, keepDelimiters map[string]bool) string {
	parts := strings.Split(options.Text, options.Delimiter)
	for j := 1; j < len(parts); j += 2 {
		if keep, exists := keepDelimiters[options.Delimiter]; exists && keep {
			parts[j] = options.Color + options.Delimiter +
				parts[j] + options.Delimiter + ColorReset
		} else {
			formattingOptions := FormattingOptions{
				Text:       parts[j],
				Delimiter:  options.Delimiter,
				Color:      options.Color,
				Formatting: options.Formatting,
			}
			parts[j] = ApplyFormatting(formattingOptions)
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
