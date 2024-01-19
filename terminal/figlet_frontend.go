package terminal

import (
	"strings"
)

// Define the ASCII patterns for the 'slant' font for the characters
var asciiPatterns = map[rune][]string{
	// Figlet in a compiled language, not an interpreted language.
	'G': {
		"   ______      ______           ___    ____  ",
		"  / ____/___  / ____/__  ____  /   |  /  _/  ",
		" / / __/ __ \\/ / __/ _ \\/ __ \\/ /| |  / /    ",
		"/ /_/ / /_/ / /_/ /  __/ / / / ___ |_/ /     ",
		"\\____/\\____/\\____/\\___/_/ /_/_/  |_/___/     ",
	},
}

// Convert a string to ASCII art using the slant font and colorize it hacker colors
func toASCIIArt(input string) string {
	// Prepare a slice of strings to hold each line of the ASCII art
	output := make([]string, len(asciiPatterns['G']))
	// Iterate over each character in the input text
	for _, char := range input {
		// Get the pattern for the current character
		pattern, ok := asciiPatterns[char]
		if !ok {
			continue // Skip characters we don't have a pattern for
		}
		// Append the pattern to the output line by line
		for i, line := range pattern {
			output[i] += line
		}
	}
	// Join the lines into a single string and colorize it
	return colors.ColorHex95b806 + strings.Join(output, "\n") + colors.ColorReset
}
