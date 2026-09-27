package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type PackageMetadata struct {
	Name string `json:"name"`
}

type PackageConfig struct {
	Metadata PackageMetadata `json:"metadata"`
}

// Case 1: Vulnerable direct join of unmarshaled field
func testVulnerableDirectJoin(data []byte, baseDir string) error {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	target := filepath.Join(baseDir, pkg.Metadata.Name)
	// ruleid: path-traversal-unmarshal
	_, err := os.Create(target)
	return err
}

// Case 2: Patched direct join using filepath.Base
func testSafeDirectJoinBase(data []byte, baseDir string) error {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	target := filepath.Join(baseDir, filepath.Base(pkg.Metadata.Name))
	// ok: path-traversal-unmarshal
	_, err := os.Create(target)
	return err
}

// Case 3: Vulnerable formatted string propagation (variant from CVE diff)
func testVulnerableFormattedString(data []byte, baseDir string) error {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-documentation", pkg.Metadata.Name)
	target := filepath.Join(baseDir, filename)
	// ruleid: path-traversal-unmarshal
	return os.WriteFile(target, []byte("docs"), 0600)
}

// Case 4: Patched formatted string with filepath.Base
func testSafeFormattedStringBase(data []byte, baseDir string) error {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return err
	}
	filename := fmt.Sprintf("%s-documentation", filepath.Base(pkg.Metadata.Name))
	target := filepath.Join(baseDir, filename)
	// ok: path-traversal-unmarshal
	return os.WriteFile(target, []byte("docs"), 0600)
}

// Case 5: Vulnerable usage of filepath.Clean which does not prevent path traversal
func testVulnerableCleanNotSanitizer(data []byte, baseDir string) (*os.File, error) {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	target := filepath.Join(baseDir, filepath.Clean(pkg.Metadata.Name))
	// ruleid: path-traversal-unmarshal
	return os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0600)
}

// Case 6: Patched path validation using filepath.IsLocal
func testSafeIsLocalCheck(data []byte, baseDir string) (*os.File, error) {
	var pkg PackageConfig
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(pkg.Metadata.Name) {
		return nil, fmt.Errorf("invalid path")
	}
	target := filepath.Join(baseDir, pkg.Metadata.Name)
	// ok: path-traversal-unmarshal
	return os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0600)
}
