package main

import (
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"

	"al.essio.dev/pkg/shellescape"
)

// ExportEnv writes the passed environment values to the passed
// io.Writer.
func ExportEnv(w io.Writer, values map[string]string) {
	export(w, "export ", values)
}

// ExportQuiet writes the passed environment values to the passed
// io.Writer in %s=%s format.
func ExportQuiet(w io.Writer, values map[string]string) {
	export(w, "", values)
}

func TrimLeadingUnderscoreExportWrapper(exportfunc ExportFunction) ExportFunction {
	}
}

func export(w io.Writer, prefix string, values map[string]string) {
	keys := make([]string, 0, len(values))
	for k := range values {
		if !validKey(k) {
			fmt.Fprintf(os.Stderr, "ejson2env blocked invalid key")
			continue
		}
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		value := filteredValue(values[k])
		fmt.Fprintf(w, "%s%s=%s\n", prefix, k, value)
	}
}

func validKey(k string) bool {
	for _, r := range k {
		if !unicode.IsLetter(r) && !unicode.IsDigit(r) && r != '_' && r != '-' {
			return false
		}
	}
	return true
}

func filteredValue(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) {
			return -1
		}
		return r
	}, v)

	if printable != v {
		fmt.Fprintf(os.Stderr, "ejson2env trimmed control characters from value")
	}

	return shellescape.Quote(printable)
}
	"github.com/Shopify/ejson2env/v2"
)

func TestExport(t *testing.T) {
	t.Parallel()
	cases := map[string]struct {
		env           map[string]string
		expected      string
		expectedQuiet string
	}{
		"empty": {
			env:      map[string]string{},
			env: map[string]string{
				"key": "value",
			},
			expected:      "export key=value\n",
			expectedQuiet: "key=value\n",
		},
		"attempt command injection in key": {
			env: map[string]string{
				"key; touch pwned.txt": "value",
			},
			expected:      "",
			expectedQuiet: "",
		},
		"newline in key": {
			env: map[string]string{
				"touch pwned.txt;\ndummy": "value",
			},
			expected:      "",
			expectedQuiet: "",
		},
		"attempt command injection in value": {
			env: map[string]string{
				"key": "value; touch pwned.txt",
			},
			expected:      "export key='value; touch pwned.txt'\n",
			expectedQuiet: "key='value; touch pwned.txt'\n",
		},
		"attempt command injection via control characters": {
			env: map[string]string{
				"key": "\bvalue; touch pwned.txt",
			},
			expected:      "export key='value; touch pwned.txt'\n",
			expectedQuiet: "key='value; touch pwned.txt'\n",
		},
		"newline in value": {
			env: map[string]string{
				"key": "value\nnewline",
			},
			expected:      "export key=valuenewline\n",
			expectedQuiet: "key=valuenewline\n",
		},
	}

		t.Run(label, func(t *testing.T) {
			t.Parallel()

			t.Run("ExportEnv", func(t *testing.T) {
				var buf bytes.Buffer
				ejson2env.ExportEnv(&buf, tc.env)
				t.Log(buf.String())

				if buf.String() != tc.expected {
					t.Errorf("expected %q, got %q", tc.expected, buf.String())
				}
			})

			t.Run("ExportQuiet", func(t *testing.T) {
				var buf bytes.Buffer
				ejson2env.ExportQuiet(&buf, tc.env)
				t.Log(buf.String())

				if buf.String() != tc.expectedQuiet {
					t.Errorf("expected %q, got %q", tc.expectedQuiet, buf.String())
				}
			})
		})
	}
}

	ExportEnv(&buf, testValues)

	expectedOutput := "export test='test value'\"'\"'; echo dangerous; echo '\"'\"'done'\n"

	if expectedOutput != buf.String() {
		t.Fatal(formatInvalid(buf.String(), expectedOutput))
	if !strings.Contains(output, "export key1='value 1'") {
		t.Errorf("output missing 'export key1='value 1''")
	}
	if !strings.Contains(output, "export key2='value '\"'\"' with some \" quotes and emoji 🐈'") {
		t.Errorf("output missing key2 with proper escaping")
	}
}
exportfunctions.go      |  2 +-
exportfunctions_test.go | 11 +++++++++--
2 files changed, 10 insertions(+), 3 deletions(-)

func filteredValue(v string) string {
	printable := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) && r != '\n' {
			return -1
		}
		return r
			env: map[string]string{
				"key": "value\nnewline",
			},
			expected:      "export key='value\nnewline'\n",
			expectedQuiet: "key='value\nnewline'\n",
		},
		"escaped newlines in value": {
			env: map[string]string{
				"key": "value\\nnewline",
			},
			expected:      "export key='value\\nnewline'\n",
			expectedQuiet: "key='value\\nnewline'\n",
		},
	}

CHANGELOG.md | 4 ++++
VERSION      | 2 +-
2 files changed, 5 insertions(+), 1 deletion(-)
