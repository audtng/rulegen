package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Edge Case 1: Direct extraction of zip.File without path sanitization (vulnerable)
func extractZipDirectVuln(f *zip.File, destDir string) error {
	target := filepath.Join(destDir, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Edge Case 2: Extraction of zip.File sanitized with filepath.Base (safe)
func extractZipBaseSafe(f *zip.File, destDir string) error {
	cleanName := filepath.Base(f.Name)
	target := filepath.Join(destDir, cleanName)
	// ok: archive-path-traversal
	out, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Edge Case 3: Zip reader iteration writing files without validation (vulnerable)
func extractZipReaderVuln(r *zip.Reader, destDir string) error {
	for _, f := range r.File {
		target := filepath.Join(destDir, f.Name)
		// ruleid: archive-path-traversal
		if err := os.WriteFile(target, []byte("content"), 0644); err != nil {
			return err
		}
	}
	return nil
}

// Edge Case 4: Zip entry validation using filepath.IsLocal (safe)
func extractZipIsLocalSafe(f *zip.File, destDir string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("unsafe file path: %s", name)
	}
	target := filepath.Join(destDir, name)
	// ok: archive-path-traversal
	if err := os.WriteFile(target, []byte("content"), 0644); err != nil {
		return err
	}
	return nil
}

// Edge Case 5: Tar archive extraction with Next() reading unvalidated header (vulnerable)
func extractTarArchiveVuln(tr *tar.Reader, destDir string) error {
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
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Edge Case 6: Tar archive extraction verified using filepath.Rel prefix check (safe)
func extractTarArchiveSafe(tr *tar.Reader, destDir string) error {
	cleanDest := filepath.Clean(destDir)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(cleanDest, hdr.Name)
		rel, err := filepath.Rel(cleanDest, target)
		if err != nil || strings.HasPrefix(rel, "..") {
			return fmt.Errorf("path traversal in tar header: %s", hdr.Name)
		}
		// ok: archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}
