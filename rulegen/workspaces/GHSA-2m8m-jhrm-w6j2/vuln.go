package main

	return hashString, nil
}

// quoteOrEscapeShellPath makes path a valid string argument in configured shell
// and also ensures it cannot cause unintended behavior.
func quoteOrEscapeShellPath(shellType string, shellPath string) (string, error) {
	// PowerShell
	if shellType == "powershell" {
		return "'" + strings.ReplaceAll(shellPath, "'", "''") + "'", nil
	}
	// Windows Command Prompt
	if shellType == "cmd" {
