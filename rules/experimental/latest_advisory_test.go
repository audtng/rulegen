package rules

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type yamlMock struct{}

func (yamlMock) Unmarshal(data []byte, v any) error {
	return nil
}

var yaml yamlMock

type VolumeConfig struct {
	Items map[string][]byte `json:"items"`
}

type FileConfig struct {
	FilePath string `json:"filePath"`
	Content  string `json:"content"`
}

type CleanupConfig struct {
	TargetDir string `json:"targetDir"`
}

// 1. Vulnerable: YAML unmarshaling of volume items written to disk without path sanitization
func vulnerableKubeVolumeConfig(data []byte, mountPoint string) error {
	var cfg VolumeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for k, v := range cfg.Items {
		dataPath := filepath.Join(mountPoint, k)
		// ruleid: untrusted-config-path-traversal
		f, err := os.Create(dataPath)
		if err != nil {
			return err
		}
		if _, err := f.Write(v); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

// 2. Safe: YAML unmarshaling of volume items sanitized with filepath.Base
func safeKubeVolumeConfigBase(data []byte, mountPoint string) error {
	var cfg VolumeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for k, v := range cfg.Items {
		baseName := filepath.Base(k)
		dataPath := filepath.Join(mountPoint, baseName)
		// ok: untrusted-config-path-traversal
		f, err := os.Create(dataPath)
		if err != nil {
			return err
		}
		if _, err := f.Write(v); err != nil {
			f.Close()
			return err
		}
		f.Close()
	}
	return nil
}

// 3. Vulnerable: JSON unmarshaling into struct, path joined with root and written via os.WriteFile
func vulnerableJSONConfigFile(data []byte, rootDir string) error {
	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	destPath := filepath.Join(rootDir, cfg.FilePath)
	// ruleid: untrusted-config-path-traversal
	return os.WriteFile(destPath, []byte(cfg.Content), 0600)
}

// 4. Safe: JSON unmarshaling into struct, validated with filepath.IsLocal
func safeJSONConfigIsLocal(data []byte, rootDir string) error {
	var cfg FileConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.FilePath) {
		return errors.New("unsafe path detected")
	}
	destPath := filepath.Join(rootDir, cfg.FilePath)
	// ok: untrusted-config-path-traversal
	return os.WriteFile(destPath, []byte(cfg.Content), 0600)
}

// 5. Vulnerable: Unmarshaled config path cleaned and formatted, used in os.OpenFile
func vulnerableConfigOpenFile(data []byte, baseDir string) error {
	var cfg FileConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	clean := filepath.Clean(cfg.FilePath)
	target := fmt.Sprintf("%s/%s", baseDir, clean)
	// ruleid: untrusted-config-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// 6. Safe: Unmarshaled config path verified with filepath.IsLocal before os.RemoveAll
func safeConfigRemoveAll(data []byte, baseDir string) error {
	var cfg CleanupConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.TargetDir) {
		return errors.New("directory is not local")
	}
	target := filepath.Join(baseDir, cfg.TargetDir)
	// ok: untrusted-config-path-traversal
	return os.RemoveAll(target)
}
