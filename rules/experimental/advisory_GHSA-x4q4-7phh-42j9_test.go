package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
)

// Test 1: Vulnerable tar extraction loop with tr.Next()
func vulnTarExtract(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, hdr.Name)
		// ruleid: archive-path-traversal
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Test 2: Vulnerable zip extraction loop with path.Join
func vulnZipExtract(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := path.Join(dest, f.Name)
		// ruleid: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Test 3: Vulnerable tar extraction where filepath.Clean is incorrectly assumed to sanitize
func vulnTarClean(hdr *tar.Header, dest string) error {
	cleaned := filepath.Clean(hdr.Name)
	target := filepath.Join(dest, cleaned)
	// ruleid: archive-path-traversal
	return os.MkdirAll(target, 0755)
}

// Test 4: Safe tar extraction validated with filepath.IsLocal
func safeTarIsLocal(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("invalid path: %s", name)
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}

// Test 5: Safe zip extraction sanitized with filepath.Base
func safeZipBase(f *zip.File, dest string) error {
	safeName := filepath.Base(f.Name)
	target := filepath.Join(dest, safeName)
	// ok: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}

// Test 6: Safe zip extraction sanitized with path.Base
func safeZipPathBase(f *zip.File, dest string) error {
	safeName := path.Base(f.Name)
	target := path.Join(dest, safeName)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}
