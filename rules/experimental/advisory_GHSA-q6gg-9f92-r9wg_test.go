package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Direct zip extraction without sanitization (Vulnerable)
func extractZipDirect(f *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, f.Name)
	// ruleid: go-archive-path-traversal
	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 2: Zip extraction sanitized with filepath.Base (Safe)
func extractZipWithBase(f *zip.File, destDir string) error {
	safeName := filepath.Base(f.Name)
	targetPath := filepath.Join(destDir, safeName)
	// ok: go-archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 3: Zip extraction sanitized with filepath.IsLocal (Safe)
func extractZipWithIsLocal(f *zip.File, destDir string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	targetPath := filepath.Join(destDir, name)
	// ok: go-archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0600)
}

// Case 4: Tar extraction in loop without sanitization (Vulnerable)
func extractTarLoop(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		destPath := filepath.Join(destDir, hdr.Name)
		if hdr.Typeflag == tar.TypeDir {
			// ruleid: go-archive-path-traversal
			if err := os.MkdirAll(destPath, 0750); err != nil {
				return err
			}
			continue
		}
	}
	return nil
}

// Case 5: Zip extraction with path splitting / partial clean but vulnerable (Vulnerable)
func extractZipSplitTraverse(f *zip.File, destDir string) error {
	parts := strings.SplitN(f.Name, "/", 2)
	if len(parts) < 2 {
		return os.ErrInvalid
	}
	cleanPart := filepath.Clean(parts[1])
	targetPath := filepath.Join(destDir, cleanPart)
	// ruleid: go-archive-path-traversal
	out, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 6: Archive extraction sanitized with filepath.Rel (Safe)
func extractArchiveSafeRel(hdr *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, hdr.Name)
	rel, err := filepath.Rel(destDir, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return os.ErrInvalid
	}
	// ok: go-archive-path-traversal
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
