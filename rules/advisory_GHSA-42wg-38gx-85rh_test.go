package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Direct stdlib usage of archive entry name in file creation
func testDirectStdlib(f *zip.File) error {
	target := filepath.Join("/tmp/extract", f.Name)
	// ruleid: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}

// Case 2: Proper patch using filepath.IsLocal validation
func testProperPatch(f *zip.File) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("invalid path: %s", name)
	}
	target := filepath.Join("/tmp/extract", name)
	// ok: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}

// Case 3: Cross-function taint propagation
func buildPath(base, name string) string {
	return filepath.Join(base, name)
}

func testCrossFunctionTaint(f *zip.File) error {
	target := buildPath("/tmp/extract", f.Name)
	// ruleid: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}

// Case 4: Interface bypass
type PathGetter interface {
	GetPath() string
}

type archiveItem struct {
	relPath string
}

func (a archiveItem) GetPath() string {
	return a.relPath
}

func testInterfaceBypass(hdr *tar.Header) error {
	var getter PathGetter = archiveItem{relPath: hdr.Name}
	target := filepath.Join("/tmp/extract", getter.GetPath())
	// ruleid: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}

// Case 5: Fake sanitizer (checks file extension suffix without verifying path traversal components)
func testFakeSanitizer(f *zip.File) error {
	fname := f.Name
	if !strings.HasSuffix(fname, ".json") {
		return fmt.Errorf("expected json file")
	}
	target := filepath.Join("/tmp/extract", fname)
	// ruleid: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}

// Case 6: Real sanitizer using filepath.Base
func testRealSanitizer(f *zip.File) error {
	name := filepath.Base(f.Name)
	target := filepath.Join("/tmp/extract", name)
	// ok: go-archive-path-traversal
	_, err := os.Create(target)
	return err
}
