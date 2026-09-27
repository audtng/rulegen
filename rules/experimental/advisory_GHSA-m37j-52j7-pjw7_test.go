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

// Case 1: Vulnerable Tar extraction using os.OpenFile directly
func ExtractTarVulnerable(tr *tar.Reader, targetDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(targetDir, hdr.Name)
		// ruleid: go-archive-path-traversal
		f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 2: Safe Tar extraction validated with filepath.IsLocal
func ExtractTarSafeIsLocal(tr *tar.Reader, targetDir string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		name := hdr.Name
		if !filepath.IsLocal(name) {
			continue
		}
		target := filepath.Join(targetDir, name)
		// ok: go-archive-path-traversal
		f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 3: Vulnerable Zip extraction using os.Create
func ExtractZipVulnerable(zr *zip.Reader, targetDir string) error {
	for _, f := range zr.File {
		target := filepath.Join(targetDir, f.Name)
		// ruleid: go-archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Case 4: Safe Zip extraction sanitized with filepath.Base
func ExtractZipSafeBase(zr *zip.Reader, targetDir string) error {
	for _, f := range zr.File {
		cleanName := filepath.Base(f.Name)
		target := filepath.Join(targetDir, cleanName)
		// ok: go-archive-path-traversal
		out, err := os.Create(target)
		if err != nil {
			return err
		}
		out.Close()
	}
	return nil
}

// Case 5: Vulnerable directory creation in Tar extraction
func ExtractTarDirVulnerable(hdr *tar.Header, targetDir string) error {
	target := filepath.Join(targetDir, hdr.Name)
	// ruleid: go-archive-path-traversal
	return os.MkdirAll(target, 0755)
}

// Case 6: Safe Tar extraction validated with prefix boundary check
func ExtractTarSafePrefix(tr *tar.Reader, targetDir string) error {
	cleanBase := filepath.Clean(targetDir) + string(os.PathSeparator)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(targetDir, hdr.Name)
		if !strings.HasPrefix(target, cleanBase) {
			return fmt.Errorf("path traversal attempt: %s", target)
		}
		// ok: go-archive-path-traversal
		f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}
