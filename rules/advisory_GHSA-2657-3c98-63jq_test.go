package rules

import (
	"archive/tar"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

// 1. Direct stdlib: extracting from tar.Header without validation
func DirectStdlib(hdr *tar.Header, dest string) error {
	target := filepath.Join(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 2. Proper patch: validate path locality with filepath.IsLocal
func ProperPatch(hdr *tar.Header, dest string) error {
	name := hdr.Name
	if !filepath.IsLocal(name) {
		return errors.New("invalid path: attempts directory traversal")
	}
	target := filepath.Join(dest, name)
	// ok: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

func helperBuildPath(baseDir, entryName string) string {
	return filepath.Join(baseDir, entryName)
}

// 3. Cross-function taint: taint propagates through helper function
func CrossFunctionTaint(hdr *tar.Header, dest string) error {
	target := helperBuildPath(dest, hdr.Name)
	// ruleid: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 4. Interface bypass: io.Reader wrapped in tar.Reader
func InterfaceBypass(r io.Reader, dest string) error {
	tr := tar.NewReader(r)
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

// 5. Fake sanitizer: using path.Clean or filepath.Clean does not prevent traversal
func FakeSanitizer(hdr *tar.Header, dest string) error {
	cleaned := path.Clean(hdr.Name)
	target := filepath.Join(dest, cleaned)
	// ruleid: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// 6. Real sanitizer: filepath.Base extracts only filename, preventing traversal
func RealSanitizer(hdr *tar.Header, dest string) error {
	safeName := filepath.Base(hdr.Name)
	target := filepath.Join(dest, safeName)
	// ok: archive-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}
