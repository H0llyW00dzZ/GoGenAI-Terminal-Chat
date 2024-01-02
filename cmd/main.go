package main

import (
	"log"
	"os"

	"github.com/H0llyW00dzZ/GoGenAI-Terminal-Chat/terminal"
)

func main() {
	apiKey := os.Getenv(aPi_Key)
	if apiKey == "" {
		log.Fatal(logFatal)
	}

	session, err := terminal.NewSession(apiKey)
	if err != nil {
		log.Fatal(err)
	}

	session.Start()
}
