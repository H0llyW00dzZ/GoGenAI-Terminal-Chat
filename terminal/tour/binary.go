// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package tour

import (
	"fmt"
)

// ConvertToBinary converts a string to its binary representation.
// Each character of the input string is represented by its binary
// equivalent, separated by spaces in the resulting string.
//
// For example:
//
//	ConvertToBinary("Go") returns "1000111 1101111"
func ConvertToBinary(input string) string {
	var binaryRepresentation string
	for _, runeValue := range input {
		// Convert each rune (Unicode code point) in the input string to a binary string.
		binaryRepresentation += fmt.Sprintf(BinaryCharHighLevel, runeValue)
	}
	// Trim the trailing space for a cleaner output.
	return binaryRepresentation[:len(binaryRepresentation)-1]
}
