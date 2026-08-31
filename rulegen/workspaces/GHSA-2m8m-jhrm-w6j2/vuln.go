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
}

func TestShellEscapePowerShell(t *testing.T) {
	for i, test := range []struct {
		unescaped, escaped string
	}{
		{"c:/test&notepad", "'c:/test&notepad'"},
		{"c:/test\"&\"notepad", "'c:/test\"&\"notepad'"},
		{"c:/test'&'notepad", "'c:/test''&''notepad'"},
	} {
		got, err := quoteOrEscapeShellPath("powershell", test.unescaped)
		assert.NoError(t, err)
		assert.Equal(t, test.escaped, got, fmt.Sprintf("Test %d unescaped = %q", i, test.unescaped))
	}
}

