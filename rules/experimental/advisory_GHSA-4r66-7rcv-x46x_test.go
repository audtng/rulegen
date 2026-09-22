package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable zip extraction using *zip.File directly into os.OpenFile
func testVulnerableZipExtract(f *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, f.Name)
	// ruleid: go-archive-path-traversal-zip-slip
	out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 2: Safe zip extraction using filepath.Base to strip path traversal components
func testSafeZipExtractBase(f *zip.File, destDir string) error {
	cleanName := filepath.Base(f.Name)
	targetPath := filepath.Join(destDir, cleanName)
	// ok: go-archive-path-traversal-zip-slip
	out, err := os.OpenFile(targetPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 3: Vulnerable tar extraction reading *tar.Header via reader loop into os.Create
func testVulnerableTarExtract(tr *tar.Reader, destDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		targetPath := filepath.Join(destDir, hdr.Name)
		// ruleid: go-archive-path-traversal-zip-slip
		out, err := os.Create(targetPath)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Case 4: Safe tar extraction validated using filepath.IsLocal
func testSafeTarExtractIsLocal(hdr *tar.Header, destDir string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return nil
	}
	targetPath := filepath.Join(destDir, name)
	// ok: go-archive-path-traversal-zip-slip
	return os.WriteFile(targetPath, []byte("safe"), 0600)
}

// Case 5: Vulnerable zip loop with flawed filepath.Clean sanitization (Clean alone does not prevent directory escape)
func testVulnerableZipLoopFlawedClean(r *zip.Reader, destDir string) error {
	for _, f := range r.File {
		targetPath := filepath.Clean(filepath.Join(destDir, f.Name))
		// ruleid: go-archive-path-traversal-zip-slip
		if err := os.WriteFile(targetPath, []byte("data"), 0600); err != nil {
			return err
		}
	}
	return nil
}

// Case 6: Safe tar extraction using destination prefix boundary verification
func testSafeTarPrefixCheck(hdr *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, hdr.Name)
	cleanDest := filepath.Clean(destDir) + string(filepath.Separator)
	if !strings.HasPrefix(targetPath, cleanDest) {
		return nil
	}
	// ok: go-archive-path-traversal-zip-slip
	out, err := os.Create(targetPath)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
