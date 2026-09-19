package rules

import (
	"net/http"
	"os"
	"path/filepath"
)

// 1. Direct stdlib: Untrusted input from HTTP query parameters flows directly into filepath.Join and os.Open
func DirectStdlib(r *http.Request) {
	filename := r.URL.Query().Get("file")
	targetPath := filepath.Join("/var/data", filename)
	// ruleid: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}

// 2. Proper patch: Input is verified with filepath.IsLocal to guarantee containment within base path
func ProperPatch(r *http.Request) {
	filename := r.URL.Query().Get("file")
	if !filepath.IsLocal(filename) {
		return
	}
	targetPath := filepath.Join("/var/data", filename)
	// ok: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}

// 3. Cross-function taint: Tainted input passes through an auxiliary path-constructing function
func buildTargetPath(base, userPath string) string {
	return filepath.Join(base, userPath)
}

func CrossFunctionTaint(r *http.Request) {
	filename := r.FormValue("file")
	targetPath := buildTargetPath("/var/data", filename)
	// ruleid: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}

// 4. Interface bypass: Tainted input is wrapped inside a struct and accessed via an interface method
type PathProvider interface {
	Path() string
}

type UserInput struct {
	val string
}

func (u UserInput) Path() string {
	return u.val
}

func InterfaceBypass(r *http.Request) {
	var provider PathProvider = UserInput{val: r.URL.Query().Get("file")}
	targetPath := filepath.Join("/var/data", provider.Path())
	// ruleid: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}

// 5. Fake sanitizer: filepath.Clean does NOT prevent traversal when relative components escape base dir
func FakeSanitizer(r *http.Request) {
	filename := r.URL.Query().Get("file")
	cleaned := filepath.Clean(filename)
	targetPath := filepath.Join("/var/data", cleaned)
	// ruleid: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}

// 6. Real sanitizer: filepath.Base safely strips any directory components, retaining only the base filename
func RealSanitizer(r *http.Request) {
	filename := r.URL.Query().Get("file")
	baseName := filepath.Base(filename)
	targetPath := filepath.Join("/var/data", baseName)
	// ok: go-web-path-traversal
	f, _ := os.Open(targetPath)
	if f != nil {
		_ = f.Close()
	}
}
