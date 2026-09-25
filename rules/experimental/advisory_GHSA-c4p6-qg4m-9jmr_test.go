package rules

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
)

type VaultCredentialConfig struct {
	ServiceAccountPath string `json:"serviceAccountPath"`
}

type AppConfig struct {
	ConfigFile   string `json:"configFile"`
	AuditLogPath string `xml:"auditLogPath"`
	TemplatePath string `json:"templatePath"`
}

// Case 1: Vulnerable - Direct read from unmarshaled JSON struct field (KEDA archetype)
func testDirectUnmarshalRead(data []byte) ([]byte, error) {
	var cred VaultCredentialConfig
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, err
	}
	// ruleid: arbitrary-file-read-unmarshaled-path
	return os.ReadFile(cred.ServiceAccountPath)
}

// Case 2: Safe - Path sanitized using filepath.Base
func testSanitizedBase(data []byte) ([]byte, error) {
	var cred VaultCredentialConfig
	if err := json.Unmarshal(data, &cred); err != nil {
		return nil, err
	}
	safeName := filepath.Base(cred.ServiceAccountPath)
	safePath := filepath.Join("/safe/tokens", safeName)
	// ok: arbitrary-file-read-unmarshaled-path
	return os.ReadFile(safePath)
}

// Case 3: Vulnerable - Tainted path propagated via fmt.Sprintf and opened with os.Open
func testPropagatedSprintfOpen(data []byte) (*os.File, error) {
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	targetPath := fmt.Sprintf("/var/data/%s", cfg.ConfigFile)
	// ruleid: arbitrary-file-read-unmarshaled-path
	return os.Open(targetPath)
}

// Case 4: Safe - Path validated using filepath.IsLocal
func testSanitizedIsLocal(data []byte) (*os.File, error) {
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(cfg.ConfigFile) {
		return nil, errors.New("path must be local")
	}
	// ok: arbitrary-file-read-unmarshaled-path
	return os.Open(cfg.ConfigFile)
}

// Case 5: Vulnerable - Unmarshaled from XML and opened via os.OpenFile
func testXMLUnmarshalOpenFile(data []byte) (*os.File, error) {
	var cfg AppConfig
	if err := xml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// ruleid: arbitrary-file-read-unmarshaled-path
	return os.OpenFile(cfg.AuditLogPath, os.O_RDONLY, 0)
}

// Case 6: Safe - Path validated using fs.ValidPath
func testSanitizedFSValidPath(data []byte) ([]byte, error) {
	var cfg AppConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if !fs.ValidPath(cfg.TemplatePath) {
		return nil, errors.New("invalid path format")
	}
	// ok: arbitrary-file-read-unmarshaled-path
	return os.ReadFile(cfg.TemplatePath)
}
