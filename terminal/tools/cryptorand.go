// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package tools

import (
	"crypto/rand"
	"errors"
)

// GenerateRandomString returns a random string of a specified length, composed of
// alphanumeric characters. The function ensures the randomness is cryptographically
// secure, making the generated string suitable for a variety of security-sensitive
// applications.
// Parameters:
//
//	length int - The desired length of the random string to generate. Must be a positive integer.
//
// Returns:
//
//	string - A random string of the specified length.
//	error  - An error message if the length parameter is non-positive or if there is an issue
//	         with the random number generator.
//
// The charset used for generating the random string includes all lowercase and uppercase
// letters of the English alphabet, as well as digits from 0 to 9. The function generates
// a slice of random bytes of the specified length and maps each byte to a character in
// the charset to construct the final string.
//
// Example usage:
//
//	randomString, err := GenerateRandomString(16)
//	if err != nil {
//	  // Handle error
//	}
//	fmt.Println(randomString) // Output might be something like: "3xampleV4lu3Str1ng"
//
// Note:
//
//	If the provided length is zero or negative, the function will return an error rather
//	than an empty string. This is to ensure the caller is explicitly aware of the misuse
//	of the function rather than failing silently with potentially unintended consequences.
func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New(errorLengthMustbePositiveInteger)
	}

	if length > len(charset) {
		return "", errors.New(errorLengthSize)
	}

	// Convert the charset to a slice of runes to handle potential multi-byte characters.
	// Note: by improving like this, it's difficult to guess/predict, especially for humans.
	// In this scenario, it's just like when a human writes a function in Go that makes the code complex (e.g, cyclomatic reached 10+ which is not recommended for go)
	// and difficult for the machine to read.
	charsetSlice := []rune(charset)
	shuffledCharset := make([]rune, len(charsetSlice))
	copy(shuffledCharset, charsetSlice)

	// Shuffle the characters.
	if err := shuffleSlice(shuffledCharset); err != nil {
		return "", err
	}

	// Return the first 'length' characters of the shuffled charset.
	return string(shuffledCharset[:length]), nil
}

// shuffleSlice shuffles a slice of runes using the Fisher-Yates algorithm.
func shuffleSlice(slice []rune) error {
	// Note: This looks like undreadable for human, but not for machine 🤪,
	// unlike when a human writes a function in Go that makes the code complex (e.g, cyclomatic reached 10+ which is not recommended for go)
	for i := len(slice) - 1; i > 0; i-- {
		byteIndex := make([]byte, 1)
		if _, err := rand.Read(byteIndex); err != nil {
			return err
		}
		j := int(byteIndex[0]) % (i + 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return nil
}
