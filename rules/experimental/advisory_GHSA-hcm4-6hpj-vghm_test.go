package rules

import (
	"archive/tar"
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
)

// Case 1: Direct tar symlink target extraction without path validation (Vulnerable)
func extractTarDirect(hdr *tar.Header, targetPath string) error {
	// ruleid: archive-symlink-target-traversal
	return os.Symlink(hdr.Linkname, targetPath)
}

// Case 2: filepath.Join does not prevent traversal for symlink target (Vulnerable)
func extractTarJoined(hdr *tar.Header, destDir, targetPath string) error {
	resolvedTarget := filepath.Join(destDir, hdr.Linkname)
	// ruleid: archive-symlink-target-traversal
	return os.Symlink(resolvedTarget, targetPath)
}

// Case 3: Hardlink traversal via archive file header without validation (Vulnerable)
func extractZipHardlink(fh *zip.FileHeader, targetPath string) error {
	// ruleid: archive-symlink-target-traversal
	return os.Link(fh.Name, targetPath)
}

// Case 4: Sanitized symlink target via filepath.Base (Safe)
func extractTarSafeBase(hdr *tar.Header, targetPath string) error {
	safeTarget := filepath.Base(hdr.Linkname)
	// ok: archive-symlink-target-traversal
	return os.Symlink(safeTarget, targetPath)
}

// Case 5: Validated symlink target via filepath.IsLocal (Safe)
func extractTarSafeLocal(hdr *tar.Header, targetPath string) error {
	target := hdr.Linkname
	if !filepath.IsLocal(target) {
		return errors.New("insecure symlink target")
	}
	// ok: archive-symlink-target-traversal
	return os.Symlink(target, targetPath)
}

// Case 6: Hardcoded safe symlink target (Safe)
func extractStaticSymlink(targetPath string) error {
	// ok: archive-symlink-target-traversal
	return os.Symlink("relative/internal/target", targetPath)
}
