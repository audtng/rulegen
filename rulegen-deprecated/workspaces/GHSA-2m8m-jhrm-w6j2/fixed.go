package main

	return hashString, nil
}

// powerShellQuoteEscaper doubles every character PowerShell accepts as a
// single-quote string delimiter. As well as the ASCII apostrophe, PowerShell
// treats the Unicode smart quotes U+2018, U+2019, U+201A and U+201B as single
// quotes, so a path wrapped in apostrophes must double all of them or an
// attacker controlled filename could close the literal and inject a statement.
// Doubling a delimiter is PowerShell's escape for a literal occurrence of it,
// and preserves the exact character.
var powerShellQuoteEscaper = strings.NewReplacer(
	"'", "''",
	"‘", "‘‘",
	"’", "’’",
	"‚", "‚‚",
	"‛", "‛‛",
)

// quoteOrEscapeShellPath makes path a valid string argument in configured shell
// and also ensures it cannot cause unintended behavior.
func quoteOrEscapeShellPath(shellType string, shellPath string) (string, error) {
	// PowerShell
	if shellType == "powershell" {
		return "'" + powerShellQuoteEscaper.Replace(shellPath) + "'", nil
	}
	// Windows Command Prompt
	if shellType == "cmd" {
