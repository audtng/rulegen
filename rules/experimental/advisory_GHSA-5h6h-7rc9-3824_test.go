package rules

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
)

type Request struct {
	Filepath string
	Target   string
}

// Case 1: Direct unvalidated SFTP client path reaching os.Open
func handleDirectReadFile(r *Request) (*os.File, error) {
	// ruleid: sftp-path-traversal
	return os.Open(r.Filepath)
}

// Case 2: Flawed prefix-based path validation allowing root escape to sibling directories
func handlePrefixValidationBypass(r *Request, root string) (*os.File, error) {
	cleanPath := filepath.Clean("/" + r.Filepath)
	if !strings.HasPrefix(cleanPath, root) {
		return nil, errors.New("access denied: outside of root")
	}
	// ruleid: sftp-path-traversal
	return os.Create(cleanPath)
}

// Case 3: SFTP rename command where Target destination path is unvalidated
func handleRenameTargetVulnerable(r *Request, root string) error {
	src := filepath.Join(root, filepath.Clean(r.Filepath))
	dst := filepath.Join(root, r.Target)
	// ruleid: sftp-path-traversal
	return os.Rename(src, dst)
}

// Case 4: Safe path handling using filepath.Base sanitization
func handleSafeBase(r *Request, root string) (*os.File, error) {
	safeName := filepath.Base(r.Filepath)
	safePath := filepath.Join(root, safeName)
	// ok: sftp-path-traversal
	return os.Open(safePath)
}

// Case 5: Safe path handling using filepath.IsLocal validation
func handleSafeIsLocal(r *Request, root string) (*os.File, error) {
	clean := filepath.Clean(r.Filepath)
	if !filepath.IsLocal(clean) {
		return nil, errors.New("invalid path")
	}
	targetPath := filepath.Join(root, clean)
	// ok: sftp-path-traversal
	return os.Open(targetPath)
}

// Case 6: Safe static path operation unaffected by untrusted input
func handleSafeInternalStatic(r *Request) (*os.File, error) {
	staticLogPath := "/var/log/sftp/audit.log"
	// ok: sftp-path-traversal
	return os.OpenFile(staticLogPath, os.O_APPEND|os.O_WRONLY, 0600)
}
