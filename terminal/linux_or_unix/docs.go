// Package linux_or_unix provides a utility to manage input history for terminal-based
// applications running on Linux or Unix systems. It includes functionality to store,
// retrieve, and navigate through a history of user inputs, which enhances the user
// experience by allowing easy access to previously entered commands or messages.

// The InputHistoryManager struct is the core component of this package, offering
// methods to add new inputs to the history, navigate to previous and next entries,
// and update the current input being edited. This functionality is particularly
// useful for command-line interfaces or chat applications where users may need to
// recall or edit prior inputs.

// The package is designed with portability in mind, targeting non-Windows platforms
// where terminal behavior for input history is consistent with Unix-like conventions.
// It employs conditional compilation to ensure that it is only built on compatible
// systems, thereby avoiding runtime errors on unsupported platforms.

// Usage:
// To use the InputHistoryManager, create a new instance using the
// NewInputHistoryManager function. As users enter inputs, add each new input to
// the manager using the Add method. When the user wishes to navigate through their
// input history, use the Previous and Next methods to retrieve past inputs.

// Example:
//  manager := linux_or_unix.NewInputHistoryManager()
//  manager.Add("first command")
//  manager.Add("second command")
//  previousInput := manager.Previous() // retrieves "second command"
//  nextInput := manager.Next()         // retrieves ""

// Note that the package is built only on non-Windows platforms. If you need similar
// functionality for Windows, you should create a corresponding package that handles
// input history according to Windows terminal behavior.

package linux_or_unix
