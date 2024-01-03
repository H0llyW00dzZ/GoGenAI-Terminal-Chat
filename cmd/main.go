// Copyright (c) 2024 H0llyW00dzZ

package main

import (
	"log"
	"os"

	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal"
)

const (
	api_Key  = "API_KEY" // Fixed the typo here
	logFatal = "API_KEY environment variable is not set"
)

func main() {
	apiKey := os.Getenv(api_Key)
	if apiKey == "" {
		log.Fatal(logFatal)
	}

	session, err := terminal.NewSession(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	session.Start()
}
