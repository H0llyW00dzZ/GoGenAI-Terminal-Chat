// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

import (
	"strings"
)

// ASCII Art
const (
	// NOTE: ' is rune not a string
	G = 'G'
	V = 'V'
	// ASCII slant font
	_G   = "   ______      ______           ___    ____  "
	_O   = "  / ____/___  / ____/__  ____  /   |  /  _/  "
	_GEN = " / / __/ __ \\/ / __/ _ \\/ __ \\/ /| |  / /    "
	A_   = "/ /_/ / /_/ / /_/ /  __/ / / / ___ |_/ /     "
	I_   = "\\____/\\____/\\____/\\___/_/ /_/_/  |_/___/     "
	// Blank Art
	BLANK_ = "                                      "
)

// Text
const (
	Current_Version = "Current Version: " + CurrentVersion
	Copyright       = "Copyright (c) 2024 @H0llyW00dzZ"
)

// Define the ASCII patterns for the 'slant' font for the characters
var asciiPatterns = map[rune][]string{
	// Figlet in a compiled language, not an interpreted language.
	// This literally header in your machine lmao.
	// It so easy implement Header like this in go, also it possible to made it animated drawing/human typing this ascii art
	// unlike "interpreted language" 🤪
	G: {
		_G,
		_O,
		_GEN,
		A_,
		I_,
	},
	V: {
		BLANK_,
		BLANK_,
		BLANK_, // TODO: Implement a notification to display here when a new version is available.
		//			 For checking the version and viewing the change log, implement the command ":checkversion".
		Current_Version,
		Copyright,
	},
}

// Define a map for character colors
var asciiColors = map[rune]string{
	G: BoldText + colors.ColorHex95b806,
	V: BoldText + colors.ColorCyan24Bit,
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

// ToASCIIArt converts a string to ASCII art representation using predefined
// patterns for each character. It applies a slant font style and colorizes
// the output in hacker-style colors. Characters that do not have a corresponding
// pattern in the asciiPatterns map are skipped.
//
// The function iterates over each character in the input string, looks up the
// pattern, and appends it line by line to the output slice. The output slice
// is then joined into a single string with newline characters separating the
// lines of the ASCII art.
//
// Parameters:
//
//	input - The string to be converted into ASCII art.
//
// Returns:
//
//	A string representing the ASCII art of the input text.
//	If the input contains characters without corresponding patterns, those
//	characters are omitted from the ASCII art representation.
//
// Example:
//
//	asciiArt := ToASCIIArt("G")
//	fmt.Println(asciiArt)
//
// This will output the ASCII art representation of "G" using the slant font
// and hacker-style colors.
func ToASCIIArt(input string) string {
	output := make([]string, len(asciiPatterns[G]))

	for _, char := range input {
		pattern, ok := asciiPatterns[char]
		if !ok {
			continue // Skip characters without a pattern
		}

		output = appendPatternLine(output, char, pattern)
	}

	return strings.Join(output, "\n")
}
