package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable tar extraction using filepath.Join directly
func extractTarBasic(tr *tar.Reader, dest string) error {
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

// Case 2: Vulnerable tar extraction using naive path.Clean (archive traversal preserved)
func extractTarNaiveClean(h *tar.Header, dest string) error {
	cleaned := path.Clean(h.Name)
	target := path.Join(dest, cleaned)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, []byte("payload"), 0644)
}

// Case 3: Vulnerable zip extraction using zip.Reader and filepath.Join
func extractZipBasic(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := filepath.Join(dest, f.Name)
		// ruleid: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Case 4: Patched tar extraction using filepath.Base
func extractTarPatchedBase(h *tar.Header, dest string) error {
	baseName := filepath.Base(h.Name)
	target := filepath.Join(dest, baseName)
	// ok: archive-path-traversal
	return os.WriteFile(target, []byte("safe"), 0644)
}

// Case 5: Patched tar extraction using filepath.IsLocal (Go 1.20+)
func extractTarPatchedIsLocal(tr *tar.Reader, dest string) error {
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path traversal: %s", name)
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}

// Case 6: Patched zip extraction using strings.HasPrefix validation
func extractZipPatchedHasPrefix(f *zip.File, dest string) error {
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	target := filepath.Join(cleanDest, f.Name)
	cleanTarget := filepath.Clean(target)
	if !strings.HasPrefix(cleanTarget, cleanDest) {
		return fmt.Errorf("insecure path traversal: %s", f.Name)
	}
	// ok: archive-path-traversal
	out, err := os.Create(cleanTarget)
	if err != nil {
		return err
	}
	return out.Close()
}
