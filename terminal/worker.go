// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"context"
	"time"
)

// NewChatWorker creates a new ChatWorker for a given chat session.
func NewChatWorker(session *Session) *ChatWorker {
	return &ChatWorker{
		session: session,
		ticker:  time.NewTicker(1 * time.Second), // Adjust the ticker interval as needed.
		done:    make(chan bool),
	}
}

// Start begins the background work loop of the ChatWorker.
func (cw *ChatWorker) Start(ctx context.Context) error {
	go func() {
		for {
			select {
			case <-cw.ticker.C:
				// Note: The current implementation is a prototype sample logic. These tasks are placeholders
				// and are not yet implemented. They serve as examples of the kind of periodic work a ChatWorker
				// might perform. Future development will include concrete implementations as per the application's requirements.
			case <-cw.done:
				// Handle cleanup and shutdown of the worker.
				return
			case <-ctx.Done():
				// Context cancellation has been requested, stop the worker.
				return
			}
		}
	}()
	return nil
}

// Stop signals the ChatWorker to stop its work loop.
func (cw *ChatWorker) Stop() error {
	cw.ticker.Stop()
	cw.done <- true
	return nil
}
