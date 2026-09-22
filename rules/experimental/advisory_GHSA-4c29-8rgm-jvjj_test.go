package rules

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type ContainerSpec struct {
	ID      string `json:"id"`
	WorkDir string `json:"work_dir"`
}

// Edge Case 1: Unmarshaled container ID joined with storage root and written via os.OpenFile
func HandleContainerCreateVulnerable(data []byte, storageRoot string) error {
	var spec ContainerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return err
	}
	targetPath := filepath.Join(storageRoot, spec.ID)
	// ruleid: go-untrusted-unmarshal-path-traversal
	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Edge Case 2: Unmarshaled nested path propagated via filepath.Clean into os.MkdirAll
func HandleStateDirVulnerable(data []byte, storageRoot string) error {
	var spec ContainerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return err
	}
	targetDir := filepath.Clean(filepath.Join(storageRoot, spec.WorkDir))
	// ruleid: go-untrusted-unmarshal-path-traversal
	return os.MkdirAll(targetDir, 0755)
}

// Edge Case 3: Unmarshaled map formatted via fmt.Sprintf into os.RemoveAll
func HandlePurgeVulnerable(data []byte, storageRoot string) error {
	var params map[string]string
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}
	target := fmt.Sprintf("%s/%s", storageRoot, params["id"])
	// ruleid: go-untrusted-unmarshal-path-traversal
	return os.RemoveAll(target)
}

// Edge Case 4: Sanitized using filepath.Base before file creation
func HandleContainerCreateSafe(data []byte, storageRoot string) error {
	var spec ContainerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return err
	}
	safeID := filepath.Base(spec.ID)
	targetPath := filepath.Join(storageRoot, safeID)
	// ok: go-untrusted-unmarshal-path-traversal
	f, err := os.OpenFile(targetPath, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()
	return nil
}

// Edge Case 5: Sanitized using filepath.IsLocal validation before directory creation
func HandleStateDirSafe(data []byte, storageRoot string) error {
	var spec ContainerSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return err
	}
	if !filepath.IsLocal(spec.WorkDir) {
		return os.ErrInvalid
	}
	targetDir := filepath.Join(storageRoot, spec.WorkDir)
	// ok: go-untrusted-unmarshal-path-traversal
	return os.MkdirAll(targetDir, 0755)
}

// Edge Case 6: Safe fixed path within storage root independent of unmarshaled message
func HandlePurgeSafe(data []byte, storageRoot string) error {
	var params map[string]string
	if err := json.Unmarshal(data, &params); err != nil {
		return err
	}
	target := filepath.Join(storageRoot, "default_state")
	// ok: go-untrusted-unmarshal-path-traversal
	return os.RemoveAll(target)
}
