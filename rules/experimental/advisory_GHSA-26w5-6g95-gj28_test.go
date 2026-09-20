package rules

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path/filepath"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{
	Unmarshal: json.Unmarshal,
}

type Config struct {
	Workspace string `json:"workspace" xml:"workspace"`
	Dir       string `json:"dir" xml:"dir"`
	Path      string `json:"path" xml:"path"`
}

func VulnerableYamlWorkspaceDelete(data []byte) error {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetPath := filepath.Join("/var/workspaces", cfg.Workspace)
	// ruleid: unmarshaled-config-path-traversal
	return os.RemoveAll(targetPath)
}

func SafeYamlWorkspaceIsLocal(data []byte) error {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.Workspace) {
		return errors.New("workspace must be a local path")
	}
	targetPath := filepath.Join("/var/workspaces", cfg.Workspace)
	// ok: unmarshaled-config-path-traversal
	return os.RemoveAll(targetPath)
}

func VulnerableJsonDirCreation(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	targetDir := filepath.Join("/var/data", cfg.Dir)
	// ruleid: unmarshaled-config-path-traversal
	return os.MkdirAll(targetDir, 0750)
}

func SafeJsonPathBase(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	safeName := filepath.Base(cfg.Path)
	targetFile := filepath.Join("/var/data", safeName)
	// ok: unmarshaled-config-path-traversal
	return os.WriteFile(targetFile, []byte("content"), 0600)
}

func VulnerableDecoderOpenFile(r io.Reader) (*os.File, error) {
	var cfg Config
	dec := json.NewDecoder(r)
	if err := dec.Decode(&cfg); err != nil {
		return nil, err
	}
	targetPath := filepath.Join("/var/storage", cfg.Path)
	// ruleid: unmarshaled-config-path-traversal
	return os.OpenFile(targetPath, os.O_RDWR|os.O_CREATE, 0600)
}

func SafeXmlStaticPath(data []byte) error {
	var cfg Config
	if err := xml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	staticPath := filepath.Join("/var/workspaces", "default", "build.log")
	// ok: unmarshaled-config-path-traversal
	return os.Remove(staticPath)
}
