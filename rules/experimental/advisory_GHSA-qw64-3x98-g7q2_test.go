package rules

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"
)

type VulnerableFS struct {
	baseDir string
}

type ValidatingFS struct {
	baseDir string
}

type LocalFS struct {
	baseDir string
}

type BaseFS struct {
	baseDir string
}

// Edge case 1: Vulnerable Open without boundary validation
func (v *VulnerableFS) Open(name string) (fs.File, error) {
	target := filepath.Join(v.baseDir, name)
	// ruleid: generic-fs-path-traversal
	return os.Open(target)
}

// Edge case 2: Safe Open with fs.ValidPath boundary validation
func (s *ValidatingFS) Open(name string) (fs.File, error) {
	if !fs.ValidPath(name) {
		return nil, errors.New("invalid path")
	}
	target := filepath.Join(s.baseDir, name)
	// ok: generic-fs-path-traversal
	return os.Open(target)
}

// Edge case 3: Vulnerable Create without boundary check
func (v *VulnerableFS) Create(name string) (*os.File, error) {
	target := filepath.Join(v.baseDir, name)
	// ruleid: generic-fs-path-traversal
	return os.Create(target)
}

// Edge case 4: Safe Create with filepath.IsLocal sanitization
func (s *LocalFS) Create(name string) (*os.File, error) {
	if !filepath.IsLocal(name) {
		return nil, errors.New("path escapes root")
	}
	target := filepath.Join(s.baseDir, name)
	// ok: generic-fs-path-traversal
	return os.Create(target)
}

// Edge case 5: Vulnerable RemoveAll with flawed filepath.Clean sanitization
func (v *VulnerableFS) RemoveAll(name string) error {
	cleaned := filepath.Clean(name)
	target := filepath.Join(v.baseDir, cleaned)
	// ruleid: generic-fs-path-traversal
	return os.RemoveAll(target)
}

// Edge case 6: Safe OpenFile with filepath.Base filename extraction
func (s *BaseFS) OpenFile(name string, flag int, perm fs.FileMode) (*os.File, error) {
	baseName := filepath.Base(name)
	target := filepath.Join(s.baseDir, baseName)
	// ok: generic-fs-path-traversal
	return os.OpenFile(target, flag, perm)
}
