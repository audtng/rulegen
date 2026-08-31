package main

import (
	"fmt"
	"io"
	"regexp"
	"strings"
	"unicode"

	"github.com/taskcluster/shell"
)

var exportCommandPattern = regexp.MustCompile(`^export [a-zA-Z_][a-zA-Z0-9_-]*=([A-Za-z0-9/:=-]+|'.*')$`)
var quietCommandPattern = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*=([A-Za-z0-9/:=-]+|'.*')$`)

// ValidateExportCommand checks if a string matches the format: export some_username='anycharacters'
func ValidateExportCommand(cmd string) bool {
	return exportCommandPattern.MatchString(cmd)
}

// ValidateQuietCommand checks if a string matches the format: some_username='anycharacters'
func ValidateQuietCommand(cmd string) bool {
	return quietCommandPattern.MatchString(cmd)
}

// ExportEnv writes the passed environment values to the passed
// io.Writer.
func ExportEnv(w io.Writer, values map[string]string) {
	for key, value := range values {
		cmd := fmt.Sprintf("export %s=%s", key, escape(value))
		if ValidateExportCommand(cmd) {
			fmt.Fprintln(w, cmd)
		}
	}
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format.
func ExportQuiet(w io.Writer, values map[string]string) {
	for key, value := range values {
		cmd := fmt.Sprintf("%s=%s", key, escape(value))
		if ValidateQuietCommand(cmd) {
			fmt.Fprintln(w, cmd)
		}
	}
}

func TrimLeadingUnderscoreExportWrapper(exportfunc ExportFunction) ExportFunction {
	}
}

func escape(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, v)
	return shell.Escape(printable)
}
	"github.com/Shopify/ejson2env/v2"
)

func TestExportEnv(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		env      map[string]string
		expected string
	}{
		"empty": {
			env:      map[string]string{},
			env: map[string]string{
				"key": "value",
			},
			expected: "export key=value\n",
		},
		"attempt command injection in key": {
			env: map[string]string{
				"key; touch pwned.txt": "value",
			},
			expected: "",
		},
		"attempt command injection in value": {
			env: map[string]string{
				"key": "value; touch pwned.txt",
			},
			expected: "export key='value; touch pwned.txt'\n",
		},
		"attempt command injection via control characters": {
			env: map[string]string{
				"key": "\bvalue; touch pwned.txt",
			},
			expected: "export key='value; touch pwned.txt'\n",
		},
	}

		t.Run(label, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			ejson2env.ExportEnv(&buf, tc.env)
			t.Log(buf.String())

			if buf.String() != tc.expected {
				t.Errorf("expected %q, got %q", tc.expected, buf.String())
			}
		})
	}
}

	ExportEnv(&buf, testValues)

	expectedOutput := "export test='test value'\\''; echo dangerous; echo '\\''done'\n"

	if expectedOutput != buf.String() {
		t.Fatal(formatInvalid(buf.String(), expectedOutput))
	if !strings.Contains(output, "export key1='value 1'") {
		t.Errorf("output missing 'export key1='value 1''")
	}
	if !strings.Contains(output, "export key2='value '\\'' with some \" quotes and emoji 🐈'") {
		t.Errorf("output missing key2 with proper escaping")
	}
}
exportfunctions.go      |  2 +-
exportfunctions_test.go | 11 +++++++++--
2 files changed, 10 insertions(+), 3 deletions(-)

func filteredValue(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
			env: map[string]string{
				"key": "value\nnewline",
			},
			expected:      "export key=valuenewline\n",
			expectedQuiet: "key=valuenewline\n",
		},
	}

CHANGELOG.md | 4 ++++
VERSION      | 2 +-
2 files changed, 5 insertions(+), 1 deletion(-)
