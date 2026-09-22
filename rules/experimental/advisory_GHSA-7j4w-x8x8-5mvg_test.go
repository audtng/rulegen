package rules

import (
	"archive/tar"
	"archive/zip"
	"os"
	"path/filepath"
)

// Case 1: Vulnerable tar extraction using os.OpenFile
func VulnTarExtract(hdr *tar.Header, destDir string) error {
	target := filepath.Join(destDir, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 2: Vulnerable zip extraction using os.Create
func VulnZipExtract(f *zip.File, destDir string) error {
	target := filepath.Join(destDir, f.Name)
	// ruleid: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 3: Vulnerable tar extraction using os.WriteFile
func VulnTarWriteFile(hdr *tar.Header, destDir string, data []byte) error {
	target := filepath.Join(destDir, hdr.Name)
	// ruleid: archive-path-traversal
	return os.WriteFile(target, data, 0600)
}

// Case 4: Safe tar extraction sanitized with filepath.Base
func SafeTarBase(hdr *tar.Header, destDir string) error {
	safeName := filepath.Base(hdr.Name)
	target := filepath.Join(destDir, safeName)
	// ok: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 5: Safe zip extraction validated with filepath.IsLocal
func SafeZipLocal(f *zip.File, destDir string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return nil
	}
	target := filepath.Join(destDir, name)
	// ok: archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// Case 6: Safe directory creation sanitized with filepath.Base
func SafeTarMkdir(hdr *tar.Header, destDir string) error {
	cleanName := filepath.Base(hdr.Name)
	dirPath := filepath.Join(destDir, cleanName)
	// ok: archive-path-traversal
	return os.MkdirAll(dirPath, 0750)
}
