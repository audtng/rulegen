package rules

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

// Request mocks sftp.Request for standalone compilation without external dependencies.
type Request struct {
	Filepath string
	Target   string
}

// Resolver defines an interface for path processing.
type Resolver interface {
	Resolve(p string) string
}

// PassthroughResolver implements Resolver.
type PassthroughResolver struct{}

func (pr PassthroughResolver) Resolve(p string) string {
	return p
}

func helperInsecure(p string) string {
	return filepath.Clean("/" + p)
}

// 1. Direct stdlib flow
func testDirectStdlib(r *Request) error {
	// ruleid: sftp-path-traversal
	f, err := os.Open(r.Filepath)
	if err != nil {
		return err
	}
	return f.Close()
}

// 2. Proper patch (validates local jail boundary)
func testProperPatch(r *Request) error {
	cleanPath := filepath.Clean(r.Filepath)
	if !filepath.IsLocal(cleanPath) {
		return errors.New("insecure path outside jail")
	}
	// ok: sftp-path-traversal
	f, err := os.Open(cleanPath)
	if err != nil {
		return err
	}
	return f.Close()
}

// 3. Cross-function taint propagation
func testCrossFunctionTaint(r *Request) error {
	p := helperInsecure(r.Filepath)
	// ruleid: sftp-path-traversal
	return os.RemoveAll(p)
}

// 4. Interface bypass
func testInterfaceBypass(r *Request, res Resolver) error {
	resolved := res.Resolve(r.Filepath)
	// ruleid: sftp-path-traversal
	_, err := os.Create(resolved)
	return err
}

// 5. Fake sanitizer (prefix-based validation flaw as seen in GHSA-5h6h-7rc9-3824)
func testFakeSanitizer(r *Request, root string) error {
	cleanPath := filepath.Clean("/" + r.Filepath)
	if !strings.HasPrefix(cleanPath, root) {
		return errors.New("access denied: outside root")
	}
	// ruleid: sftp-path-traversal
	f, err := os.Open(cleanPath)
	if err != nil {
		return err
	}
	return f.Close()
}

// 6. Real sanitizer (strips path components)
func testRealSanitizer(r *Request) error {
	base := filepath.Base(r.Filepath)
	// ok: sftp-path-traversal
	f, err := os.Open(base)
	if err != nil {
		return err
	}
	return f.Close()
}
