package rules

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{
	Unmarshal: func(b []byte, v any) error { return nil },
}

type ConfigPayload struct {
	FilePath string `json:"file_path"`
	Content  string `json:"content"`
}

// Case 1: Vulnerable - JSON unmarshaling into struct followed by filepath.Join and os.ReadFile
func testVulnerableJSON(data []byte, baseDir string) ([]byte, error) {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	targetPath := filepath.Join(baseDir, cfg.FilePath)
	// ruleid: path-traversal-unmarshaled-data
	return os.ReadFile(targetPath)
}

// Case 2: Vulnerable - YAML unmarshaling followed by filepath.Join and os.Open
func testVulnerableYAML(data []byte, baseDir string) (*os.File, error) {
	var cfg ConfigPayload
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	targetPath := filepath.Join(baseDir, cfg.FilePath)
	// ruleid: path-traversal-unmarshaled-data
	return os.Open(targetPath)
}

// Case 3: Vulnerable - JSON stream decoder followed by filepath.Join and os.OpenFile
func testVulnerableDecoder(data string, baseDir string) (*os.File, error) {
	var cfg ConfigPayload
	dec := json.NewDecoder(strings.NewReader(data))
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	targetPath := filepath.Join(baseDir, cfg.FilePath)
	// ruleid: path-traversal-unmarshaled-data
	return os.OpenFile(targetPath, os.O_RDONLY, 0)
}

// Case 4: Safe - Path sanitized using filepath.Base
func testSafeBase(data []byte, baseDir string) ([]byte, error) {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	safeName := filepath.Base(cfg.FilePath)
	targetPath := filepath.Join(baseDir, safeName)
	// ok: path-traversal-unmarshaled-data
	return os.ReadFile(targetPath)
}

// Case 5: Safe - Path validated using filepath.IsLocal
func testSafeIsLocal(data []byte, baseDir string) ([]byte, error) {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(cfg.FilePath) {
		return nil, os.ErrInvalid
	}
	targetPath := filepath.Join(baseDir, cfg.FilePath)
	// ok: path-traversal-unmarshaled-data
	return os.ReadFile(targetPath)
}

// Case 6: Safe - Unmarshaled payload data used as content, path is static constant
func testSafeConstantPath(data []byte, baseDir string) error {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetPath := filepath.Join(baseDir, "static_output.txt")
	// ok: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte(cfg.Content), 0o600)
}
