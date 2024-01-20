// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package tools

import (
	"crypto/rand"
	"errors"
)

func GenerateRandomString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New(errorLengthMustbePositiveInteger)
	}

	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	for i, b := range bytes {
		bytes[i] = charset[b%byte(len(charset))]
	}
	return string(bytes), nil
}
