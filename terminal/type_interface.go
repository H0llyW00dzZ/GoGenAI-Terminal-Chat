// Copyright (c) 2024 H0llyW00dzZ
//
// License: MIT License

package terminal

import (
	"context"
)

// Worker defines the interface for a background worker in the terminal application.
type Worker interface {
	Start(ctx context.Context) error
	Stop() error
}
