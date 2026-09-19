package rules

import (
	"errors"
	"net/http"
	"path/filepath"
)

// 1. Direct stdlib: user input from HTTP request flows into filepath.Join and http.ServeFile
func DirectStdlib(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	target := filepath.Join("/var/www/uploads", userInput)
	// ruleid: http-path-traversal
	http.ServeFile(w, r, target)
}

// 2. Proper patch: user input is validated with filepath.IsLocal before use
func ProperPatch(w http.ResponseWriter, r *http.Request) error {
	userInput := r.URL.Query().Get("file")
	if !filepath.IsLocal(userInput) {
		return errors.New("path escapes root directory")
	}
	target := filepath.Join("/var/www/uploads", userInput)
	// ok: http-path-traversal
	http.ServeFile(w, r, target)
	return nil
}

// 3. Cross-function taint: taint flows across a helper function
func formatPath(p string) string {
	return p
}

func CrossFunctionTaint(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	formatted := formatPath(userInput)
	target := filepath.Join("/var/www/uploads", formatted)
	// ruleid: http-path-traversal
	http.ServeFile(w, r, target)
}

// 4. Interface bypass: tainted input flows through an interface{} / any type assertion
func InterfaceBypass(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	var val any = userInput
	target := filepath.Join("/var/www/uploads", val.(string))
	// ruleid: http-path-traversal
	http.ServeFile(w, r, target)
}

// 5. Fake sanitizer: filepath.Clean does not prevent traversal outside root
func FakeSanitizer(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	cleaned := filepath.Clean(userInput)
	target := filepath.Join("/var/www/uploads", cleaned)
	// ruleid: http-path-traversal
	http.ServeFile(w, r, target)
}

// 6. Real sanitizer: filepath.Base extracts only filename, preventing traversal
func RealSanitizer(w http.ResponseWriter, r *http.Request) {
	userInput := r.URL.Query().Get("file")
	base := filepath.Base(userInput)
	target := filepath.Join("/var/www/uploads", base)
	// ok: http-path-traversal
	http.ServeFile(w, r, target)
}
