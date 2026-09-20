package rules

import (
	"archive/tar"
	"archive/zip"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Case 1: Vulnerable Tar extraction joining header Name directly with destination directory
func VulnTarExtract(tr *tar.Reader, dest string) error {
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		target := filepath.Join(dest, hdr.Name)
		// ruleid: go-archive-path-traversal
		f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 2: Safe Tar extraction using filepath.Base to prevent directory traversal
func SafeTarExtractBase(hdr *tar.Header, dest string) error {
	cleanName := filepath.Base(hdr.Name)
	target := filepath.Join(dest, cleanName)
	// ok: go-archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 3: Vulnerable Tar extraction with flawed sanitization using filepath.Clean
func VulnTarCleanBypass(hdr *tar.Header, dest string) error {
	cleaned := filepath.Clean(hdr.Name)
	target := filepath.Join(dest, cleaned)
	// ruleid: go-archive-path-traversal
	return os.WriteFile(target, []byte("data"), 0644)
}

// Case 4: Safe Tar extraction validating path with filepath.IsLocal
func SafeTarIsLocal(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return os.ErrInvalid
	}
	target := filepath.Join(dest, name)
	// ok: go-archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 5: Vulnerable Zip extraction creating parent directories from untrusted entry path
func VulnZipExtract(zr *zip.Reader, dest string) error {
	for _, f := range zr.File {
		target := filepath.Join(dest, f.Name)
		// ruleid: go-archive-path-traversal
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			return err
		}
	}
	return nil
}

// Case 6: Safe Zip extraction verifying destination prefix to prevent escaping target root
func SafeZipPrefixCheck(f *zip.File, dest string) error {
	target := filepath.Join(dest, f.Name)
	cleanDest := filepath.Clean(dest) + string(filepath.Separator)
	if !strings.HasPrefix(target, cleanDest) {
		return os.ErrInvalid
	}
	// ok: go-archive-path-traversal
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	return out.Close()
}
