package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

// Case 1: Vulnerable - Direct unmarshaling and path join flowing to os.WriteFile
func testDirectJoinWrite(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetPath := filepath.Join("/var/lib/vz/snippets", cfg.Filename)
	// ruleid: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte(cfg.Content), 0644)
}

// Case 2: Safe - Filename sanitized using filepath.Base
func testSafeFilepathBase(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	safeName := filepath.Base(cfg.Filename)
	targetPath := filepath.Join("/var/lib/vz/snippets", safeName)
	// ok: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte(cfg.Content), 0644)
}

// Case 3: Safe - Path validated using filepath.IsLocal guard
func testSafeFilepathIsLocal(data []byte) (*os.File, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(cfg.Filename) {
		return nil, fmt.Errorf("invalid local path: %s", cfg.Filename)
	}
	targetPath := filepath.Join("/var/lib/vz/snippets", cfg.Filename)
	// ok: path-traversal-unmarshaled-data
	return os.Create(targetPath)
}

// Case 4: Vulnerable - filepath.Clean does not prevent path traversal escapes
func testVulnerableFilepathClean(data []byte) (*os.File, error) {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	cleanedPath := filepath.Clean(filepath.Join("/var/lib/vz/snippets", cfg.Filename))
	// ruleid: path-traversal-unmarshaled-data
	return os.OpenFile(cleanedPath, os.O_CREATE|os.O_WRONLY, 0644)
}

// Case 5: Safe - Untrusted unmarshaled content written to a static destination path
func testSafeStaticPath(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetPath := "/var/lib/vz/snippets/default.conf"
	// ok: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte(cfg.Content), 0644)
}

// Case 6: Vulnerable - Path formatted via fmt.Sprintf flowing to os.RemoveAll
func testVulnerableSprintfRemove(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetPath := fmt.Sprintf("/var/lib/vz/snippets/%s", cfg.Filename)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(targetPath)
}
