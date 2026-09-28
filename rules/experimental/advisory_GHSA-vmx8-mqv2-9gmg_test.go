package rules

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// Mock third-party yaml parser to allow offline Go compilation without external modules
var yaml = struct {
	Unmarshal func([]byte, any) error
}{
	Unmarshal: func(data []byte, v any) error { return nil },
}

type PluginMetadata struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

type Config struct {
	FilePath string `json:"file_path"`
}

// Case 1: Vulnerable YAML unmarshal into struct, joined with base dir, leading to arbitrary file write.
func testYAMLUnmarshalPathTraversal_Vulnerable(data []byte) error {
	var meta PluginMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return err
	}
	targetPath := filepath.Join("/opt/helm/plugins", meta.Version)
	// ruleid: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte("payload"), 0644)
}

// Case 2: Safe YAML unmarshal with filepath.Base sanitization.
func testYAMLUnmarshalPathTraversal_Safe_FilepathBase(data []byte) error {
	var meta PluginMetadata
	if err := yaml.Unmarshal(data, &meta); err != nil {
		return err
	}
	safeVersion := filepath.Base(meta.Version)
	targetPath := filepath.Join("/opt/helm/plugins", safeVersion)
	// ok: path-traversal-unmarshaled-data
	return os.WriteFile(targetPath, []byte("payload"), 0644)
}

// Case 3: Vulnerable JSON unmarshal into pointer struct, cleaned with filepath.Clean (insufficient), leading to open.
func testJSONUnmarshalOpenFile_Vulnerable(data []byte) (*os.File, error) {
	cfg := new(Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	cleaned := filepath.Clean(cfg.FilePath)
	target := filepath.Join("/var/data", cleaned)
	// ruleid: path-traversal-unmarshaled-data
	return os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

// Case 4: Safe JSON unmarshal validated with filepath.IsLocal.
func testJSONUnmarshalOpenFile_Safe_IsLocal(data []byte) (*os.File, error) {
	cfg := new(Config)
	if err := json.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(cfg.FilePath) {
		return nil, fmt.Errorf("invalid path: %s", cfg.FilePath)
	}
	target := filepath.Join("/var/data", cfg.FilePath)
	// ok: path-traversal-unmarshaled-data
	return os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
}

// Case 5: Vulnerable streaming JSON decode with fmt.Sprintf, leading to arbitrary directory removal.
func testDecoderRemoveAll_Vulnerable(r io.Reader) error {
	var meta PluginMetadata
	dec := json.NewDecoder(r)
	if err := dec.Decode(&meta); err != nil {
		return err
	}
	fullPath := fmt.Sprintf("/var/plugins/%s", meta.Version)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(fullPath)
}

// Case 6: Safe streaming JSON decode validated with filepath.IsLocal, leading to mkdir.
func testDecoderMkdirAll_Safe_IsLocal(r io.Reader) error {
	var meta PluginMetadata
	dec := json.NewDecoder(r)
	if err := dec.Decode(&meta); err != nil {
		return err
	}
	if !filepath.IsLocal(meta.Version) {
		return fmt.Errorf("traversal detected")
	}
	safeDir := filepath.Join("/var/plugins", meta.Version)
	// ok: path-traversal-unmarshaled-data
	return os.MkdirAll(safeDir, 0755)
}
