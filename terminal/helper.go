package terminal

import "log"

// RecoverFromPanic returns a function that recovers from panics and logs the error.
// This function can be deferred at the beginning of a goroutine or function to
// handle unexpected panics in a controlled manner.
//
// Usage:
//   defer terminal.RecoverFromPanic()
func RecoverFromPanic() func() {
	return func() {
		if r := recover(); r != nil {
			// Log the panic with additional context if desired
			log.Printf(RecoverGopher, r)
		}
	}
}
