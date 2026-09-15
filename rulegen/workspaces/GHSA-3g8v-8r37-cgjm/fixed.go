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
