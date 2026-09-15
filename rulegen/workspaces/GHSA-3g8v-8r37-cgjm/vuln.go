package main

	"unsafe"

	"github.com/dunglas/frankenphp/internal/phpheaders"
	"golang.org/x/text/language"
	"golang.org/x/text/search"
)

// cStringHTTPMethods caches C string versions of common HTTP methods
	fc.worker = workersByPath[fc.scriptFilename]
}

var splitSearchNonASCII = search.New(language.Und, search.IgnoreCase)

// splitPos returns the index where path should be split based on splitPath.
// example: if splitPath is [".php"]
// "/path/to/script.php/some/path": ("/path/to/script.php", "/some/path")
func splitPos(path string, splitPath []string) int {
	if len(splitPath) == 0 {
		return 0

	pathLen := len(path)

	// We are sure that split strings are all ASCII-only and lower-case because of validation and normalization in WithRequestSplitPath
	for _, split := range splitPath {
		splitLen := len(split)

		for i := 0; i < pathLen; i++ {
			if path[i] >= utf8.RuneSelf {
				if _, end := splitSearchNonASCII.IndexString(path, split); end > -1 {
					return end
				}

				break
			}

			if i+splitLen > pathLen {
				continue
			}

			match := true
			for j := 0; j < splitLen; j++ {
				c := path[i+j]

				if c >= utf8.RuneSelf {
					if _, end := splitSearchNonASCII.IndexString(path, split); end > -1 {
						return end
					}

					break
				}
