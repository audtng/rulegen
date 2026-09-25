package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable zip extraction using filepath.Join without validation
func ExtractZipVulnerable(f *zip.File, dest string) error {
	path := filepath.Join(dest, f.Name)
	// ruleid: archive-path-traversal
	outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
	if err != nil {
		return err
	}
	defer outFile.Close()
	return nil
}

// Case 2: Safe zip extraction validated with strings.HasPrefix (patch archetype)
func ExtractZipSafePrefix(f *zip.File, dest string) error {
	path := filepath.Join(dest, f.Name)
	cleanDest := filepath.Clean(dest) + string(os.PathSeparator)
	if !strings.HasPrefix(path, cleanDest) {
		return fmt.Errorf("illegal file path: %s", path)
	}
	// ok: archive-path-traversal
	outFile, err := os.Create(path)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return nil
}

// Case 3: Vulnerable tar extraction using Next() header without validation
func ExtractTarVulnerable(tr *tar.Reader, dest string) error {
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}

// Case 4: Safe zip extraction using filepath.Base
func ExtractZipSafeBase(f *zip.File, dest string) error {
	safePath := filepath.Join(dest, filepath.Base(f.Name))
	// ok: archive-path-traversal
	outFile, err := os.Create(safePath)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return nil
}

// Case 5: Safe zip extraction using filepath.IsLocal
func ExtractZipSafeIsLocal(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path: %s", name)
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	return os.MkdirAll(target, 0755)
}

// Case 6: Safe tar extraction using filepath.Rel validation
func ExtractTarSafeRel(h *tar.Header, dest string) error {
	target := filepath.Join(dest, h.Name)
	rel, err := filepath.Rel(dest, target)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("illegal relative path: %s", target)
	}
	// ok: archive-path-traversal
	outFile, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer outFile.Close()
	return nil
}
