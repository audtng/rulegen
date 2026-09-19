package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// 1. Direct stdlib: Unsanitized extraction directly from tar.Reader Next()
func directStdlib(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(destDir, hdr.Name)
		// ruleid: archive-path-traversal
		f, err := os.Create(target)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// 2. Proper patch: Validate that destination path remains within destDir
func properPatch(hdr *tar.Header, destDir string) error {
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	target := filepath.Join(destDir, hdr.Name)
	cleanTarget := filepath.Clean(target)
	if !strings.HasPrefix(cleanTarget, cleanDest) {
		return errors.New("illegal file path: directory traversal attempt")
	}
	// ok: archive-path-traversal
	f, err := os.OpenFile(cleanTarget, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// Helper for cross-function taint
func extractTarEntry(hdr *tar.Header, destDir string) error {
	target := filepath.Join(destDir, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// 3. Cross-function taint: Taint flows across function boundaries
func crossFunctionTaint(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if err := extractTarEntry(hdr, destDir); err != nil {
			return err
		}
	}
	return nil
}

// 4. Interface bypass: Taint passed through an interface / type assertion
func interfaceBypass(entry any, destDir string) error {
	if f, ok := entry.(*zip.File); ok {
		target := filepath.Join(destDir, f.Name)
		// ruleid: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// 5. Fake sanitizer: filepath.Clean does not prevent traversal out of destination directory
func fakeSanitizer(hdr *tar.Header, destDir string) error {
	cleaned := filepath.Clean(hdr.Name)
	target := filepath.Join(destDir, cleaned)
	// ruleid: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// 6. Real sanitizer: filepath.Base safely strips directory traversal
func realSanitizer(hdr *tar.Header, destDir string) error {
	baseName := filepath.Base(hdr.Name)
	target := filepath.Join(destDir, baseName)
	// ok: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}
