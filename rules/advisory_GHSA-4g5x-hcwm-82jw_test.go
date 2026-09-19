package rules

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// PathResolver defines an interface for retrieving path strings.
type PathResolver interface {
	Resolve() string
}

// UserPathResolver implements PathResolver containing user-controlled input.
type UserPathResolver struct {
	rawPath string
}

func (u UserPathResolver) Resolve() string {
	return u.rawPath
}

func buildPath(base, untrusted string) string {
	return filepath.Join(base, untrusted)
}

// 1. Direct stdlib: Direct flow from standard library HTTP request to file reading sink.
func DirectStdlib(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	// ruleid: go-path-traversal-arbitrary-file-read
	data, err := os.ReadFile(filePath)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 2. Proper patch: Verification with filepath.IsLocal prevents directory traversal.
func ProperPatch(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	if !filepath.IsLocal(filePath) {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	// ok: go-path-traversal-arbitrary-file-read
	data, err := os.ReadFile(filepath.Join("/safe/base/dir", filePath))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 3. Cross-function taint: Untrusted path propagated through a helper function into file sink.
func CrossFunctionTaint(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	target := buildPath("/safe/base/dir", filePath)
	// ruleid: go-path-traversal-arbitrary-file-read
	data, err := os.ReadFile(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 4. Interface bypass: Untrusted path stored inside a struct and accessed via interface method.
func InterfaceBypass(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	var resolver PathResolver = UserPathResolver{rawPath: filePath}
	// ruleid: go-path-traversal-arbitrary-file-read
	f, err := os.Open(resolver.Resolve())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer f.Close()
}

// 5. Fake sanitizer: Incomplete sanitization (naive string replacement) bypassable with "../".
func FakeSanitizer(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	fakeSanitized := strings.ReplaceAll(filePath, "../", "")
	target := filepath.Join("/safe/base/dir", fakeSanitized)
	// ruleid: go-path-traversal-arbitrary-file-read
	data, err := os.ReadFile(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}

// 6. Real sanitizer: Stripping all directory components using filepath.Base.
func RealSanitizer(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	safeName := filepath.Base(filePath)
	target := filepath.Join("/safe/base/dir", safeName)
	// ok: go-path-traversal-arbitrary-file-read
	data, err := os.ReadFile(target)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(data)
}
