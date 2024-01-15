//go:build !windows
// +build !windows

package linux_or_unix

// InputHistoryManager manages the history of user inputs during a chat session.
// It allows navigating through the history of commands or messages that have been
// entered previously.
type InputHistoryManager struct {
	History      []string // History stores the list of inputs entered by the user.
	CurrentIndex int      // CurrentIndex represents the position in the input history.
	CurrentInput string   // CurrentInput stores the input currently being edited.
}

// Add appends a new input to the end of the history. After adding the input,
// it resets the current index to -1, which deselects any history entry.
func (m *InputHistoryManager) Add(input string) {
	m.History = append(m.History, input)
	m.CurrentIndex = -1
}

// Previous retrieves the input that precedes the current one in the history.
// If there is a previous input, it updates the current index to point to that
// input and returns it. If the beginning of the history is reached, it returns
// the current input without changing the index.
func (m *InputHistoryManager) Previous() string {
	if m.CurrentIndex < len(m.History)-1 {
		m.CurrentIndex++
		return m.History[len(m.History)-1-m.CurrentIndex]
	}
	return m.CurrentInput
}

// Next retrieves the input that follows the current one in the history.
// If there is a next input, it updates the current index to point to that
// input and returns it. If the end of the history is reached or if there is no
// next input, it returns an empty string and resets the current index.
func (m *InputHistoryManager) Next() string {
	if m.CurrentIndex > 0 {
		m.CurrentIndex--
		return m.History[len(m.History)-1-m.CurrentIndex]
	}
	if m.CurrentIndex == 0 {
		m.CurrentIndex--
	}
	return ""
}

// UpdateCurrentInput sets the CurrentInput field to the provided input value.
// This method is typically called when navigating through the input history
// to update the input being edited by the user.
func (m *InputHistoryManager) UpdateCurrentInput(input string) {
	m.CurrentInput = input
}

// NewInputHistoryManager creates and returns a new instance of InputHistoryManager.
// It initializes the history with an empty slice and sets the current index to -1,
// indicating that no history entry is currently selected.
// TODO: Used it for Arrow Up and Down to retrieve previous and next input.
func NewInputHistoryManager() *InputHistoryManager {
	return &InputHistoryManager{
		History:      make([]string, 0),
		CurrentIndex: -1,
	}
}
