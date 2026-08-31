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
}

func TestShellEscapePowerShell(t *testing.T) {
	// PowerShell treats U+2018, U+2019, U+201A and U+201B as single-quote
	// delimiters in addition to the ASCII apostrophe, so all of them are
	// doubled inside the wrapping apostrophes. Doubling is PowerShell's escape
	// for a literal delimiter and preserves the exact character.
	unquote := strings.NewReplacer("''", "'", "‘‘", "‘", "’’", "’", "‚‚", "‚", "‛‛", "‛")
	for i, test := range []struct {
		unescaped, escaped string
	}{
		{"c:/test&notepad", "'c:/test&notepad'"},
		{"c:/test\"&\"notepad", "'c:/test\"&\"notepad'"},
		{"c:/test'&'notepad", "'c:/test''&''notepad'"},
		// injection attempts via the ASCII apostrophe and each smart quote
		{"x';calc;#", "'x'';calc;#'"},
		{"x’;calc;#", "'x’’;calc;#'"},
		{"x‘;calc;#", "'x‘‘;calc;#'"},
		{"x‚;calc;#", "'x‚‚;calc;#'"},
		{"x‛;calc;#", "'x‛‛;calc;#'"},
	} {
		got, err := quoteOrEscapeShellPath("powershell", test.unescaped)
		assert.NoError(t, err)
		assert.Equal(t, test.escaped, got, fmt.Sprintf("Test %d unescaped = %q", i, test.unescaped))
		// Every single-quote delimiter must appear an even number of times so
		// none is left unpaired to close the literal early.
		for _, q := range []string{"'", "‘", "’", "‚", "‛"} {
			assert.Zero(t, strings.Count(got, q)%2, fmt.Sprintf("Test %d odd %q count in %q", i, q, got))
		}
		// Undoubling the quoted body recovers the original path exactly.
		assert.Equal(t, test.unescaped, unquote.Replace(got[1:len(got)-1]), fmt.Sprintf("Test %d round-trip", i))
	}
}

