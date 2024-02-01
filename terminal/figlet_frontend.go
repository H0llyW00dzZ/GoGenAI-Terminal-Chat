// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"fmt"
	"strings"
)

// NewASCIIArtStyle creates and returns a new ASCIIArtStyle map. It initializes
// an empty map that can be populated with ASCII art characters using the AddChar
// method or by direct assignment.
func NewASCIIArtStyle() ASCIIArtStyle {
	return make(ASCIIArtStyle)
}

// AddChar adds a new character with its ASCII art representation to the style.
// If the character already exists in the style, its pattern and color are updated
// with the new values provided.
//
// Parameters:
//
//	char    - The rune representing the character to add or update.
//	pattern - A slice of strings representing the ASCII art pattern for the character.
//	color   - A string representing the color for the character.
func (style ASCIIArtStyle) AddChar(char rune, pattern []string, color string) {
	style[char] = ASCIIArtChar{Pattern: pattern, Color: color}
}

// applyColor applies a color to a given line if the color exists.
func applyColor(artChar ASCIIArtChar, line string) string {
	if artChar.Color == "" {
		return line // No color to apply
	}
	return artChar.Color + line + ColorReset
}

// ToASCIIArt converts a string to its ASCII art representation using a given style.
// Each character in the input string is mapped to an ASCIIArtChar based on the provided style,
// and the resulting ASCII art is constructed line by line. If a character in the input string
// does not have an ASCII art representation in the style, it will be omitted from the output.
//
// The function returns the complete ASCII art as a string, where each line of the art is
// separated by a newline character. If the style is empty or a character is not found in the
// style, an error is returned.
//
// Example usage:
//
//	art, err := ToASCIIArt("G", slantStyle)
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Println(art)
//
// Parameters:
//
//	input - The string to be converted into ASCII art.
//	style - The ASCIIArtStyle used to style the ASCII art representation.
//
// Returns:
//
//	A string representing the ASCII art of the input text, and an error if the style is empty
//	or a character is not found in the style.
func ToASCIIArt(input string, style ASCIIArtStyle) (string, error) {
	if err := checkStyle(style); err != nil {
		logger.Error(ErrorToASCIIArtcheckstyle, err) // Use the package-level logger
		return "", err
	}

	maxHeight := maxPatternHeight(style)
	output := make([]string, maxHeight)

	for _, char := range input {
		if err := buildOutput(&output, char, style); err != nil {
			logger.Error(ErrorToASCIIArtbuildOutput, err) // Log the error
			// Handle the error as needed, e.g., continue, return, etc.
		}
	}

	return strings.Join(output, StringNewLine), nil
}

// checkStyle verifies that the provided ASCIIArtStyle is not empty.
// An ASCIIArtStyle is considered valid if it contains at least one character pattern.
//
// Parameters:
//
//	style - The ASCIIArtStyle to validate.
//
// Returns:
//
//	An error if the style is empty; otherwise, nil.
func checkStyle(style ASCIIArtStyle) error {
	if len(style) == 0 {
		return fmt.Errorf(ErrorStyleIsEmpty)
	}
	return nil
}

// maxPatternHeight calculates the maximum height of the patterns in the given ASCIIArtStyle.
// The height of a pattern is determined by the number of strings in the Pattern slice
// of an ASCIIArtChar. The function iterates over all characters in the style and returns
// the height of the tallest pattern.
//
// Parameters:
//
//	style - The ASCIIArtStyle from which to calculate the maximum pattern height.
//
// Returns:
//
//	The maximum height (in lines) of the patterns in the ASCIIArtStyle.
func maxPatternHeight(style ASCIIArtStyle) int {
	maxHeight := 0
	for _, art := range style {
		if h := len(art.Pattern); h > maxHeight {
			maxHeight = h
		}
	}
	return maxHeight
}

// buildOutput appends the ASCII art representation of a character to a slice of strings,
// each representing a line of the output. The ASCII art pattern for the character is retrieved
// from the provided style, and each line of the pattern is colored if a color is specified.
// If the character does not exist in the style, an error is returned.
//
// Parameters:
//
//	output - A pointer to a slice of strings, each representing a line of the ASCII art output.
//	char   - The character to be converted into ASCII art.
//	style  - The ASCIIArtStyle that contains patterns and colors for ASCII art characters.
//
// Returns:
//
//	An error if the character is not found in the style; otherwise, nil.
func buildOutput(output *[]string, char rune, style ASCIIArtStyle) error {
	art, exists := style[char]
	if !exists {
		return fmt.Errorf(ErrorCharacterNotFoundinStyle, char)
	}
	for i := range *output {
		if i < len(art.Pattern) {
			(*output)[i] += applyColor(art, art.Pattern[i])
		} else {
			(*output)[i] += " " // Add spaces if the pattern is shorter than the max height
		}
	}
	return nil
}

// MergeStyles combines multiple ASCIIArtStyle objects into a single style.
// In case of overlapping characters, the patterns from the styles appearing later
// in the arguments will overwrite those from earlier ones. This function is useful
// when you want to create a composite style that includes patterns from multiple
// sources or to override specific characters' styles in a base style.
//
// Parameters:
//
//	styles ...ASCIIArtStyle - A variadic slice of ASCIIArtStyle objects to be merged.
//
// Returns:
//
//	A new ASCIIArtStyle object that is the union of all provided styles.
func MergeStyles(styles ...ASCIIArtStyle) ASCIIArtStyle {
	mergedStyle := NewASCIIArtStyle()
	for _, style := range styles {
		for char, artChar := range style {
			mergedStyle[char] = artChar
		}
	}
	return mergedStyle
}
