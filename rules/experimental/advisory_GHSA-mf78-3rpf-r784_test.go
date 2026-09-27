package rules

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SQLiteConn struct{}
type VTab interface{}

// Edge Case 1: Direct unvalidated file read in Connect callback (vulnerable)
type DirectModule struct{}

func (m *DirectModule) Connect(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	// ruleid: sqlite-vtab-path-traversal
	f, err := os.Open(args[0])
	return f, err
}

// Edge Case 2: Path traversal via filepath.Join in Create callback (vulnerable)
type JoinedModule struct{}

func (m *JoinedModule) Create(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	targetPath := filepath.Join("/var/app/data", args[0])
	// ruleid: sqlite-vtab-path-traversal
	data, err := os.ReadFile(targetPath)
	return data, err
}

// Edge Case 3: Sprintf path injection into os.OpenFile (vulnerable)
type FormattedModule struct{}

func (m *FormattedModule) Connect(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	targetPath := fmt.Sprintf("/var/data/%s.csv", args[0])
	// ruleid: sqlite-vtab-path-traversal
	f, err := os.OpenFile(targetPath, os.O_RDONLY, 0)
	return f, err
}

// Edge Case 4: Safe access with filepath.IsLocal validation (safe)
type LocalModule struct{}

func (m *LocalModule) Connect(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	filePath := args[0]
	if !filepath.IsLocal(filePath) {
		return nil, errors.New("path must be local and not traverse upwards")
	}
	// ok: sqlite-vtab-path-traversal
	f, err := os.Open(filePath)
	return f, err
}

// Edge Case 5: Safe access by stripping directory traversal via filepath.Base (safe)
type BaseModule struct{}

func (m *BaseModule) Create(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	safeName := filepath.Base(args[0])
	safePath := filepath.Join("/var/app/data", safeName)
	// ok: sqlite-vtab-path-traversal
	data, err := os.ReadFile(safePath)
	return data, err
}

// Edge Case 6: Safe access with containment verification via prefix check (safe)
type PrefixModule struct{}

func (m *PrefixModule) Connect(c *SQLiteConn, args []string) (VTab, error) {
	if len(args) < 1 {
		return nil, errors.New("missing path argument")
	}
	clean := filepath.Clean(args[0])
	if !strings.HasPrefix(clean, "/var/app/data/") {
		return nil, errors.New("path outside allowed directory")
	}
	// ok: sqlite-vtab-path-traversal
	f, err := os.Open(clean)
	return f, err
}
