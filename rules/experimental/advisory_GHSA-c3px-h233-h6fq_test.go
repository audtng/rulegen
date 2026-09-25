package rules

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
)

type dummyYaml struct{}

func (dummyYaml) Unmarshal(data []byte, v any) error {
	return nil
}

var yaml dummyYaml

// Edge case 1: Direct struct field from json.Unmarshal into os.ReadFile (vulnerable)
func testJSONDirectFileRead(data []byte) ([]byte, error) {
	var cfg struct {
		IncludePath string `json:"include"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// ruleid: go-unmarshaled-path-traversal
	return os.ReadFile(cfg.IncludePath)
}

// Edge case 2: Map value from yaml.Unmarshal joined with base directory into os.Open (vulnerable)
func testYAMLJoinedPathFileOpen(data []byte, baseDir string) (*os.File, error) {
	var composeData map[string]string
	if err := yaml.Unmarshal(data, &composeData); err != nil {
		return nil, err
	}
	targetPath := filepath.Join(baseDir, composeData["path"])
	// ruleid: go-unmarshaled-path-traversal
	return os.Open(targetPath)
}

// Edge case 3: XML struct field cleaned with filepath.Clean into os.Create (vulnerable)
func testXMLCleanedPathCreate(data []byte) (*os.File, error) {
	var doc struct {
		OutputFile string `xml:"output_file"`
	}
	if err := xml.Unmarshal(data, &doc); err != nil {
		return nil, err
	}
	cleaned := filepath.Clean(doc.OutputFile)
	// ruleid: go-unmarshaled-path-traversal
	return os.Create(cleaned)
}

// Edge case 4: Sanitized using filepath.Base before file read (safe)
func testSanitizedWithBase(data []byte, baseDir string) ([]byte, error) {
	var cfg struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	fileName := filepath.Base(cfg.Path)
	safePath := filepath.Join(baseDir, fileName)
	// ok: go-unmarshaled-path-traversal
	return os.ReadFile(safePath)
}

// Edge case 5: Validated with filepath.IsLocal guard before file read (safe)
func testSanitizedWithIsLocal(data []byte) ([]byte, error) {
	var cfg struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if !filepath.IsLocal(cfg.Path) {
		return nil, os.ErrInvalid
	}
	// ok: go-unmarshaled-path-traversal
	return os.ReadFile(cfg.Path)
}

// Edge case 6: Static hardcoded path unaffected by unmarshaled data (safe)
func testStaticSafePath(data []byte) (*os.File, error) {
	var cfg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	// ok: go-unmarshaled-path-traversal
	return os.Open("/etc/app/default.conf")
}
