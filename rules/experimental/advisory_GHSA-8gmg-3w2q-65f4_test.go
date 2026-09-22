package rules

import (
	"os"
	"path/filepath"
)

// Case 1: Vulnerable os.Getenv without sanitization used in os.OpenFile
func VulnEnvOpenFile(baseDir string) error {
	tmp := os.Getenv("TMPDIR")
	target := filepath.Join(baseDir, tmp, "payload.jar")
	// ruleid: env-path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 2: Vulnerable os.LookupEnv without sanitization used in os.Create
func VulnLookupEnvCreate(baseDir string) error {
	tmp, ok := os.LookupEnv("TMPDIR")
	if !ok {
		tmp = "/tmp"
	}
	target := filepath.Join(baseDir, tmp, "payload.jar")
	// ruleid: env-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 3: Vulnerable os.Getenv with filepath.Clean used in os.WriteFile
func VulnEnvCleanWriteFile(destDir string, data []byte) error {
	subDir := os.Getenv("CUSTOM_DIR")
	cleanSub := filepath.Clean(subDir)
	target := filepath.Join(destDir, cleanSub, "output.txt")
	// ruleid: env-path-traversal
	return os.WriteFile(target, data, 0600)
}

// Case 4: Safe os.Getenv sanitized with filepath.Base
func SafeEnvBase(destDir string) error {
	tmp := os.Getenv("TMPDIR")
	clean := filepath.Base(tmp)
	target := filepath.Join(destDir, clean)
	// ok: env-path-traversal
	f, err := os.OpenFile(target, os.O_WRONLY|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 5: Safe os.LookupEnv validated with filepath.IsLocal
func SafeLookupEnvIsLocal(destDir string) error {
	tmp, ok := os.LookupEnv("TMPDIR")
	if !ok {
		return nil
	}
	if !filepath.IsLocal(tmp) {
		return nil
	}
	target := filepath.Join(destDir, tmp)
	// ok: env-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Case 6: Safe directory creation sanitized with filepath.Base
func SafeEnvBaseMkdir(destDir string) error {
	tmp := os.Getenv("TMPDIR")
	clean := filepath.Base(tmp)
	dirPath := filepath.Join(destDir, clean)
	// ok: env-path-traversal
	return os.MkdirAll(dirPath, 0750)
}
