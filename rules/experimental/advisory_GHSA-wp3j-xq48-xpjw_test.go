package rules

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
)

var yaml = struct {
	Unmarshal func([]byte, any) error
}{
	Unmarshal: func(data []byte, v any) error { return nil },
}

type VolumeConfig struct {
	Items map[string]string `json:"items" yaml:"items"`
}

type PodSpec struct {
	Name string `xml:"name"`
	Path string `xml:"path"`
}

type AppConfig struct {
	FileName string `json:"fileName"`
	Data     string `json:"data"`
}

// Case 1: Vulnerable - JSON unmarshaling with filepath.Join into os.Create
func testVulnerableJSONJoin(data []byte, mountPoint string) error {
	var cfg VolumeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for k := range cfg.Items {
		dataPath := filepath.Join(mountPoint, k)
		// ruleid: path-traversal-unmarshaled-data
		f, err := os.Create(dataPath)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 2: Vulnerable - YAML unmarshaling with direct path into os.WriteFile
func testVulnerableYAMLWrite(data []byte) error {
	var files map[string]string
	if err := yaml.Unmarshal(data, &files); err != nil {
		return err
	}
	for targetPath, content := range files {
		// ruleid: path-traversal-unmarshaled-data
		if err := os.WriteFile(targetPath, []byte(content), 0644); err != nil {
			return err
		}
	}
	return nil
}

// Case 3: Vulnerable - XML unmarshaling through filepath.Clean into os.OpenFile
func testVulnerableXMLClean(data []byte, baseDir string) error {
	var pod PodSpec
	if err := xml.Unmarshal(data, &pod); err != nil {
		return err
	}
	cleaned := filepath.Clean(pod.Path)
	target := filepath.Join(baseDir, cleaned)
	// ruleid: path-traversal-unmarshaled-data
	f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	f.Close()
	return nil
}

// Case 4: Patched - JSON unmarshaling sanitized with filepath.Base
func testPatchedBaseSanitizer(data []byte, mountPoint string) error {
	var cfg VolumeConfig
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for k := range cfg.Items {
		safeName := filepath.Base(k)
		dataPath := filepath.Join(mountPoint, safeName)
		// ok: path-traversal-unmarshaled-data
		f, err := os.Create(dataPath)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 5: Patched - YAML unmarshaling guarded with filepath.IsLocal
func testPatchedIsLocalValidation(data []byte, baseDir string) error {
	var cfg VolumeConfig
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return err
	}
	for k := range cfg.Items {
		if !filepath.IsLocal(k) {
			continue
		}
		target := filepath.Join(baseDir, k)
		// ok: path-traversal-unmarshaled-data
		f, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
		f.Close()
	}
	return nil
}

// Case 6: Safe - Unmarshaled config data, but file path is hardcoded/static
func testSafeHardcodedPath(data []byte, baseDir string) error {
	var appCfg AppConfig
	if err := json.Unmarshal(data, &appCfg); err != nil {
		return err
	}
	staticPath := filepath.Join(baseDir, "fixed_config.json")
	// ok: path-traversal-unmarshaled-data
	return os.WriteFile(staticPath, []byte(appCfg.Data), 0644)
}
