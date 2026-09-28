package rules

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"
)

type ManifestEntry struct {
	Name string `json:"name"`
}

type Manifest struct {
	Entries []ManifestEntry `json:"entries"`
}

// 1. Vulnerable: JSON unmarshaling into struct with relative path joined to base dir and created
func TestVulnerableJSONRestore(t *testing.T) {
	data := []byte(`{"entries": [{"name": "../../etc/cron.d/evil"}]}`)
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}
	baseDir := "/var/data"
	for _, entry := range m.Entries {
		dest := filepath.Join(baseDir, entry.Name)
		// ruleid: unmarshaled-path-traversal
		f, err := os.Create(dest)
		if err != nil {
			continue
		}
		f.Close()
	}
}

// 2. Safe: Same restore logic, but path is sanitized with filepath.Base
func TestSafeJSONRestoreWithBase(t *testing.T) {
	data := []byte(`{"entries": [{"name": "../../etc/cron.d/evil"}]}`)
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return
	}
	baseDir := "/var/data"
	for _, entry := range m.Entries {
		cleanName := filepath.Base(entry.Name)
		dest := filepath.Join(baseDir, cleanName)
		// ok: unmarshaled-path-traversal
		f, err := os.Create(dest)
		if err != nil {
			continue
		}
		f.Close()
	}
}

// 3. Vulnerable: Chained Decoder decode into struct and path.Join into os.WriteFile
func TestVulnerableDecoderPathJoin(t *testing.T) {
	reader := strings.NewReader(`{"name": "../evil.sh"}`)
	var entry ManifestEntry
	if err := json.NewDecoder(reader).Decode(&entry); err != nil {
		return
	}
	targetPath := path.Join("/app/scripts", entry.Name)
	// ruleid: unmarshaled-path-traversal
	_ = os.WriteFile(targetPath, []byte("#!/bin/sh\n"), 0755)
}

// 4. Safe: Validated with filepath.IsLocal before writing
func TestSafeIsLocalValidation(t *testing.T) {
	data := []byte(`{"name": "sub/valid.txt"}`)
	var entry ManifestEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return
	}
	if !filepath.IsLocal(entry.Name) {
		return
	}
	dest := filepath.Join("/var/data", entry.Name)
	// ok: unmarshaled-path-traversal
	_ = os.WriteFile(dest, []byte("valid content"), 0644)
}

// 5. Vulnerable: Unmarshaled path passed to os.RemoveAll
func TestVulnerableRemoveAll(t *testing.T) {
	data := []byte(`{"name": "../../var/log"}`)
	var entry ManifestEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return
	}
	targetDir := filepath.Join("/tmp/backup", entry.Name)
	// ruleid: unmarshaled-path-traversal
	_ = os.RemoveAll(targetDir)
}

// 6. Safe: Hardcoded / constant target path independent of unmarshaled filename
func TestSafeHardcodedPath(t *testing.T) {
	data := []byte(`{"name": "../../evil"}`)
	var entry ManifestEntry
	if err := json.Unmarshal(data, &entry); err != nil {
		return
	}
	safePath := filepath.Join("/var/log", "audit.log")
	// ok: unmarshaled-path-traversal
	_ = os.WriteFile(safePath, []byte(entry.Name), 0644)
}
