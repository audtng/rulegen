package rules

import (
	"archive/tar"
	"archive/zip"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// 1. Tar extraction without validation (vulnerable)
func extractTarEntryVulnerable(tr *tar.Reader, destDir string) error {
	hdr, err := tr.Next()
	if err != nil {
		return err
	}
	targetPath := filepath.Join(destDir, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 2. Tar extraction sanitized with filepath.Base (safe)
func extractTarEntrySafeBase(hdr *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, filepath.Base(hdr.Name))
	// ok: archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0644)
}

// 3. Zip extraction directly writing file (vulnerable)
func extractZipFileVulnerable(zf *zip.File, destDir string) error {
	targetPath := filepath.Join(destDir, zf.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0644)
}

// 4. Zip extraction sanitized with filepath.IsLocal (safe)
func extractZipFileSafeIsLocal(zf *zip.File, destDir string) error {
	name := zf.Name
	if !filepath.IsLocal(name) {
		return fmt.Errorf("insecure path: %s", name)
	}
	targetPath := filepath.Join(destDir, name)
	// ok: archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0644)
}

// 5. Tar directory creation with filepath.Clean (vulnerable - filepath.Clean does not prevent traversal)
func extractTarDirTraversalCleanVulnerable(hdr *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, filepath.Clean(hdr.Name))
	// ruleid: archive-path-traversal
	return os.MkdirAll(targetPath, 0755)
}

// 6. Tar extraction validated with filepath.Rel (safe)
func extractTarSafeRelCheck(hdr *tar.Header, destDir string) error {
	targetPath := filepath.Join(destDir, hdr.Name)
	rel, err := filepath.Rel(destDir, targetPath)
	if err != nil || strings.HasPrefix(rel, "..") {
		return fmt.Errorf("path escapes destination directory")
	}
	// ok: archive-path-traversal
	return os.WriteFile(targetPath, []byte("data"), 0644)
}
