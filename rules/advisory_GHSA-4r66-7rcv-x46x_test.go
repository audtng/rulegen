package rules

import (
	"archive/zip"
	"errors"
	"os"
	"path/filepath"
)

// DirectStdlib tests direct extraction from standard library archive without validation.
func DirectStdlib(f *zip.File, dest string) error {
	path := filepath.Join(dest, f.Name)
	// ruleid: archive-zip-slip
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// ProperPatch tests extraction using filepath.IsLocal path verification.
func ProperPatch(f *zip.File, dest string) error {
	name := f.Name
	if !filepath.IsLocal(name) {
		return errors.New("insecure file path: escapes destination")
	}
	path := filepath.Join(dest, name)
	// ok: archive-zip-slip
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

func resolveArchivePath(dest, name string) string {
	return filepath.Join(dest, name)
}

// CrossFunctionTaint tests taint propagation through helper functions across call boundaries.
func CrossFunctionTaint(f *zip.File, dest string) error {
	path := resolveArchivePath(dest, f.Name)
	// ruleid: archive-zip-slip
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// InterfaceBypass tests taint persistence through interface{} boxing and type assertion.
func InterfaceBypass(f *zip.File, dest string) error {
	var val any = f.Name
	pathStr := val.(string)
	target := filepath.Join(dest, pathStr)
	// ruleid: archive-zip-slip
	out, err := os.Create(target)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// FakeSanitizer tests flawed sanitization with filepath.Clean which does not prevent directory traversal.
func FakeSanitizer(f *zip.File, dest string) error {
	cleaned := filepath.Clean(f.Name)
	path := filepath.Join(dest, cleaned)
	// ruleid: archive-zip-slip
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}

// RealSanitizer tests proper path stripping using filepath.Base.
func RealSanitizer(f *zip.File, dest string) error {
	safeName := filepath.Base(f.Name)
	path := filepath.Join(dest, safeName)
	// ok: archive-zip-slip
	out, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer out.Close()
	return nil
}
