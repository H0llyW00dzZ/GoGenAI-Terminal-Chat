package terminal

// handleK8sCommand would be a handler function for a hypothetical ":k8s" command.
// Note: This command is secret of H0llyW00dzZ (Original Author) would be used to interact with a Kubernetes cluster.
func (c *handleK8sCommand) Execute(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// storageCommand would be a handler function for a hypothetical ":storage" command.
func (c *storageCommand) Execute(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// savehistorytostorageCommand would be a handler function for a hypothetical ":save history" command.
func (c *savehistorytostorageCommand) Execute(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// loadhistoryfromstorageCommand would be a handler function for a hypothetical ":load history" command.
func (c *loadhistoryfromstorageCommand) Execute(session *Session) (bool, error) {
	// currently unimplemented
	return true, nil
}

// reportshitFunctionthatTooComplexCommand would be a handler function for a hypothetical ":report" command.
func (c *reportshitFunctionthatTooComplexCommand) Execute(session *Session) (bool, error) {
	// currently unimplemented
	// this would be used to report when goroutines panic in other side not because of this terminal LOL
	return true, nil
}