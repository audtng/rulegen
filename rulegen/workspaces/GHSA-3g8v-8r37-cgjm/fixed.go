package main

	"unsafe"

	"github.com/dunglas/frankenphp/internal/phpheaders"
)

// cStringHTTPMethods caches C string versions of common HTTP methods
	fc.worker = workersByPath[fc.scriptFilename]
}

// splitPos returns the index where path should be split based on splitPath.
// example: if splitPath is [".php"]
// "/path/to/script.php/some/path": ("/path/to/script.php", "/some/path")
//
// Matching is strictly ASCII case-insensitive. Bytes >= utf8.RuneSelf in path
// never match any split entry: split strings are validated ASCII-only and
// lower-cased in WithRequestSplitPath, so any Unicode equivalence (e.g.
// fullwidth or mathematical letters folding to ASCII) would let an attacker
// upload a file whose name contains such code points and have it served as
// PHP. See GHSA-3g8v-8r37-cgjm and GHSA-v4h7-cj44-8fc8.
func splitPos(path string, splitPath []string) int {
	if len(splitPath) == 0 {
		return 0

	pathLen := len(path)

	for _, split := range splitPath {
		splitLen := len(split)
		if splitLen == 0 || splitLen > pathLen {
			continue
		}

		for i := 0; i <= pathLen-splitLen; i++ {
			match := true
			for j := 0; j < splitLen; j++ {
				c := path[i+j]
				if c >= utf8.RuneSelf {
					match = false

					break
				}
			splitPath: []string{".php"},
			wantPos:   9,
		},
		// Regression tests for GHSA-3g8v-8r37-cgjm: an inner non-ASCII byte
		// caused the loop to break without resetting match=false, so a path
		// such as "/PoC-match-unset.¡.txt" was reported as ".php" matched.
		{
			name:      "non-ascii byte after dot must not match",
			path:      "/PoC-match-unset.¡.txt",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "non-ascii byte mid-extension must not match",
			path:      "/script.p\xc2\xa1p",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		// Regression tests for GHSA-v4h7-cj44-8fc8: search.IgnoreCase folded
		// Unicode equivalents (fullwidth, mathematical, circled letters,
		// fullwidth/small full-stop) onto ASCII ".php".
		{
			name:      "small full stop ﹒ in extension must not match",
			path:      "/shell﹒php",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "fullwidth full stop ． in extension must not match",
			path:      "/shell．php",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "fullwidth p in extension must not match",
			path:      "/shell.ｐhp",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "circled php must not match",
			path:      "/shell.ⓟⓗⓟ",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "mathematical sans-serif bold php must not match",
			path:      "/shell.\U0001D5FD\U0001D5F5\U0001D5FD",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "mathematical script php must not match",
			path:      "/shell.\U0001D4C5\U0001D4BD\U0001D4C5",
			splitPath: []string{".php"},
			wantPos:   -1,
		},
		{
			name:      "circled php with later real php still picks the real one",
			path:      "/shell.ⓟⓗⓟ.anything-after-payload.php",
			splitPath: []string{".php"},
			// "/shell." (7) + "ⓟⓗⓟ" (3*3 bytes) + ".anything-after-payload.php" (27) = 43
			wantPos: 43,
		},
	}

	for _, tt := range tests {
		assert.Equal(t, ".txt.php", pathInfo, "path info should be the remainder after first .php")
	}
}

// TestSplitPosSecurityRegressionUnicodeBypass guards against
// GHSA-3g8v-8r37-cgjm (uninitialized match flag on inner non-ASCII byte) and
// GHSA-v4h7-cj44-8fc8 (Unicode equivalence via search.IgnoreCase letting
// non-PHP files be picked up as the script). Every payload below produced a
// false positive in the vulnerable implementation; none must match here.
func TestSplitPosSecurityRegressionUnicodeBypass(t *testing.T) {
	t.Parallel()

	split := []string{".php"}
	payloads := []string{
		"/PoC-match-unset.¡.txt",                // GHSA-3g8v: match left set after IndexString fallback returned -1
		"/shell﹒php",                            // U+FE52 small full stop
		"/shell．php",                            // U+FF0E fullwidth full stop
		"/shell.ｐhp",                            // U+FF50 fullwidth p
		"/shell.pｈp",                            // U+FF48 fullwidth h
		"/shell.phｐ",                            // U+FF50 fullwidth p (trailing)
		"/shell.\U0001D5C1\U0001D5B5\U0001D5C1", // mathematical sans-serif p/h
		"/shell.\U0001D5FD\U0001D5F5\U0001D5FD", // mathematical sans-serif bold p/h
		"/shell.\U0001D4C5\U0001D4BD\U0001D4C5", // mathematical script p/h
		"/shell.ⓟⓗⓟ",                            // circled latin small
	}

	for _, p := range payloads {
		assert.Equalf(t, -1, splitPos(p, split), "payload %q must not be detected as .php", p)
	}
}
