// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"strings"
)

// Define the ASCII patterns for the 'slant' font for the characters
var asciiPatterns = map[rune][]string{
	// Figlet in a compiled language, not an interpreted language.
	// This literally header in your machine lmao.
	'G': {
		"   ______      ______           ___    ____  ",
		"  / ____/___  / ____/__  ____  /   |  /  _/  ",
		" / / __/ __ \\/ / __/ _ \\/ __ \\/ /| |  / /    ",
		"/ /_/ / /_/ / /_/ /  __/ / / / ___ |_/ /     ",
		"\\____/\\____/\\____/\\___/_/ /_/_/  |_/___/     ",
	},
	'V': {
		"",
		"",
		"",
		"Current Version: " + CurrentVersion,
	},
}

// Define a map for character colors
var asciiColors = map[rune]string{
	'G': BoldText + colors.ColorHex95b806,
	'V': BoldText + colors.ColorCyan24Bit,
}

// applyColor applies a color to a given line if the hacker color exists.
func applyColor(char rune, line string) string {
	color, colorOK := asciiColors[char]
	if !colorOK {
		return line // No color to apply
	}
	return color + line + colors.ColorReset
}

// appendPatternLine appends a pattern line to the output with the proper color.
func appendPatternLine(output []string, char rune, pattern []string) []string {
	for i, line := range pattern {
		coloredLine := applyColor(char, line)
		output[i] += coloredLine
	}
	return output
}

// Convert a string to ASCII art using the slant font and colorize it hacker colors
func toASCIIArt(input string) string {
	output := make([]string, len(asciiPatterns['G']))

	for _, char := range input {
		pattern, ok := asciiPatterns[char]
		if !ok {
			continue // Skip characters without a pattern
		}

		output = appendPatternLine(output, char, pattern)
	}

	return strings.Join(output, "\n")
}
