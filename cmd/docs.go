// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

// Package cmd/main is the entry point of the GoGenAI-Terminal-Chat application.
// It sets up the environment for the chat session and starts the interaction
// between the user and the generative AI model.
//
// The application retrieves the necessary API key from environment variables
// and initializes a new chat session. If the API key is not set or an error
// occurs during the initialization of the session, the application will log
// the fatal error and exit.
//
// Usage:
// To run the application, ensure that the API_KEY environment variable is set
// with your generative AI service provider's API key. Then execute the binary:
//
//	API_KEY=your_api_key_here ./GoGenAI-Terminal-Chat
//
// The application will start, and you can begin chatting with the AI. To stop
// the application, send an interrupt signal (Ctrl+C) from the terminal.
//
// Note: This application requires the 'terminal' package, which provides the
// necessary functionality for handling the chat session and communication with
// the AI model. Make sure to import the 'terminal' package and any of its
// dependencies correctly.
//
// Copyright (c) 2024 H0llyW00dzZ
package main
