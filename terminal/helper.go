package terminal

import "log"

// RecoverFromPanic returns a deferred function that recovers from panics within a goroutine
// or function, preventing the panic from propagating and potentially causing the program to crash.
// Instead, it logs the panic information using the standard logger, allowing for post-mortem analysis
// without interrupting the program's execution flow.
//
// Usage:
//   defer terminal.RecoverFromPanic()()
//
// The function returned by RecoverFromPanic should be called by deferring it at the start of
// a goroutine or function. When a panic occurs, the deferred function will handle the panic
// by logging its message and stack trace, as provided by the recover built-in function.
func RecoverFromPanic() func() {
	return func() {
		if r := recover(); r != nil {
			// Log the panic with additional context if desired
			log.Printf(RecoverGopher, r)
		}
	}
}
