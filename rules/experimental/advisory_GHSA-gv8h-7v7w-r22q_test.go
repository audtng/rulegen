package rules

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Filename string `json:"filename"`
}

type OCIArtifact struct {
	Annotations map[string]string `json:"annotations"`
}

type XMLDocument struct {
	XMLName xml.Name `xml:"document"`
	Path    string   `xml:"path"`
}

// Case 1: Vulnerable JSON unmarshaling where untrusted path is joined without validation
func VulnJSONUnmarshal(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	target := filepath.Join(baseDir, cfg.Filename)
	// ruleid: unmarshaled-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 2: Safe JSON unmarshaling where path is sanitized with filepath.Base
func SafeJSONUnmarshalBase(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	safeName := filepath.Base(cfg.Filename)
	target := filepath.Join(baseDir, safeName)
	// ok: unmarshaled-path-traversal
	f, err := os.Create(target)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 3: Vulnerable JSON streaming decoder with map annotations (OCI layer metadata pattern)
func VulnJSONDecoderMap(r io.Reader, baseDir string) error {
	var artifact OCIArtifact
	if err := json.NewDecoder(r).Decode(&artifact); err != nil {
		return err
	}
	filePath := artifact.Annotations["com.docker.compose.file"]
	dest := filepath.Join(baseDir, filePath)
	// ruleid: unmarshaled-path-traversal
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	return f.Close()
}

// Case 4: Safe JSON streaming decoder validated with filepath.IsLocal
func SafeJSONDecoderIsLocal(r io.Reader, baseDir string) error {
	var artifact OCIArtifact
	if err := json.NewDecoder(r).Decode(&artifact); err != nil {
		return err
	}
	filePath := artifact.Annotations["com.docker.compose.envfile"]
	if !filepath.IsLocal(filePath) {
		return fmt.Errorf("insecure path: %s", filePath)
	}
	dest := filepath.Join(baseDir, filePath)
	// ok: unmarshaled-path-traversal
	return os.WriteFile(dest, []byte("ENV_VAR=1"), 0600)
}

// Case 5: Vulnerable XML unmarshaling leading to path traversal in os.RemoveAll
func VulnXMLUnmarshal(data []byte, baseDir string) error {
	var doc XMLDocument
	if err := xml.Unmarshal(data, &doc); err != nil {
		return err
	}
	targetDir := filepath.Join(baseDir, doc.Path)
	// ruleid: unmarshaled-path-traversal
	return os.RemoveAll(targetDir)
}

// Case 6: Safe path verification using strings.HasPrefix containment check
func SafePathPrefixCheck(data []byte, baseDir string) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	cleanBase := filepath.Clean(baseDir) + string(filepath.Separator)
	target := filepath.Join(baseDir, cfg.Filename)
	cleanTarget := filepath.Clean(target)
	if !strings.HasPrefix(cleanTarget, cleanBase) {
		return fmt.Errorf("path traversal attempt: %s", cfg.Filename)
	}
	// ok: unmarshaled-path-traversal
	return os.WriteFile(cleanTarget, []byte("data"), 0644)
}
