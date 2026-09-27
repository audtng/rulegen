package rules

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
)

// Vuln 1: Map key from JSON Unmarshal flows through filepath.Join to os.WriteFile
func vulnWriteFileFromJSONMap(data []byte, baseDir string) error {
	var credentials map[string][]byte
	json.Unmarshal(data, &credentials)
	for k, v := range credentials {
		credPath := filepath.Join(baseDir, k)
		// ruleid: go-path-traversal-config-file-write
		os.WriteFile(credPath, v, 0o400)
	}
	return nil
}

// Vuln 2: Map key from JSON Unmarshal flows through filepath.Join to os.Chown
func vulnChownFromJSONMap(data []byte, baseDir string) error {
	var files map[string]string
	json.Unmarshal(data, &files)
	for name := range files {
		fpath := filepath.Join(baseDir, name)
		// ruleid: go-path-traversal-config-file-write
		os.Chown(fpath, 0, 0)
	}
	return nil
}

// Vuln 3: Struct field from JSON NewDecoder flows through filepath.Join to os.Create
func vulnCreateFromDecoder(r io.Reader, baseDir string) error {
	var config struct {
		Filename string `json:"filename"`
	}
	json.NewDecoder(r).Decode(&config)
	outPath := filepath.Join(baseDir, config.Filename)
	// ruleid: go-path-traversal-config-file-write
	os.Create(outPath)
	return nil
}

// Safe 1: filepath.Base strips directory traversal components
func safeWithBase(data []byte, baseDir string) error {
	var credentials map[string][]byte
	json.Unmarshal(data, &credentials)
	for k, v := range credentials {
		safeName := filepath.Base(k)
		credPath := filepath.Join(baseDir, safeName)
		// ok: go-path-traversal-config-file-write
		os.WriteFile(credPath, v, 0o400)
	}
	return nil
}

// Safe 2: filepath.IsLocal validates path is confined before use
func safeWithIsLocal(data []byte, baseDir string) error {
	var credentials map[string][]byte
	json.Unmarshal(data, &credentials)
	for k, v := range credentials {
		if !filepath.IsLocal(k) {
			continue
		}
		credPath := filepath.Join(baseDir, k)
		// ok: go-path-traversal-config-file-write
		os.WriteFile(credPath, v, 0o400)
	}
	return nil
}

// Safe 3: Hardcoded filename — no tainted data reaches the sink
func safeHardcodedPath(baseDir string) error {
	credPath := filepath.Join(baseDir, "default.conf")
	// ok: go-path-traversal-config-file-write
	os.WriteFile(credPath, []byte("safe"), 0o644)
	return nil
}
