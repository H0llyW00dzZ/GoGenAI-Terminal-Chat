// Copyright (c) 2024 H0llyW00dzZ
// License: MIT License

package terminal

// Note: This list of commands has already been implemented.
// It is now located here for ease of maintenance and to avoid unnecessary complexity.
// This approach questions why many developers write Go code in an overly complex manner (that I don't fucking understand),
// which often leads to problems.
type handleQuitCommand struct{}
type handleHelpCommand struct{}
type handleCheckVersionCommand struct{}

// Note: this unimplemented
// Now even it's unimplemented, it wont detected in deadcode indicate that "unreachable func"
type handleK8sCommand struct{}
type storageCommand struct{}
type savehistorytostorageCommand struct{}
type loadhistoryfromstorageCommand struct{}
type reportshitFunctionthatTooComplexCommand struct{}
