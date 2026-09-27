package rules

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"os"
	"path"
	"path/filepath"
)

type ConfigPayload struct {
	TargetFile string `json:"target_file"`
}

type XMLMetadata struct {
	FilePath string `xml:"file_path"`
}

type OCIArtifactManifest struct {
	Annotations map[string]string `json:"annotations"`
}

// Edge Case 1: Vulnerable JSON unmarshaling into struct, path joined without validation
func edgeCase1VulnerableJSONUnmarshal(data []byte, baseDir string) error {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	dest := filepath.Join(baseDir, cfg.TargetFile)
	// ruleid: path-traversal-unmarshaled-data
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	return f.Close()
}

// Edge Case 2: Patched JSON unmarshaling using filepath.IsLocal validation
func edgeCase2PatchedJSONUnmarshalIsLocal(data []byte, baseDir string) error {
	var cfg ConfigPayload
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	if !filepath.IsLocal(cfg.TargetFile) {
		return errors.New("invalid path: must be a local relative path")
	}
	dest := filepath.Join(baseDir, cfg.TargetFile)
	// ok: path-traversal-unmarshaled-data
	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	return f.Close()
}

// Edge Case 3: Vulnerable XML unmarshaling into metadata struct, path joined and removed
func edgeCase3VulnerableXMLUnmarshal(data []byte, baseDir string) error {
	var meta XMLMetadata
	if err := xml.Unmarshal(data, &meta); err != nil {
		return err
	}
	targetPath := path.Join(baseDir, meta.FilePath)
	// ruleid: path-traversal-unmarshaled-data
	return os.RemoveAll(targetPath)
}

// Edge Case 4: Patched XML unmarshaling using filepath.Base sanitization
func edgeCase4PatchedXMLUnmarshalFilepathBase(data []byte, baseDir string) error {
	var meta XMLMetadata
	if err := xml.Unmarshal(data, &meta); err != nil {
		return err
	}
	safeName := filepath.Base(meta.FilePath)
	targetPath := filepath.Join(baseDir, safeName)
	// ok: path-traversal-unmarshaled-data
	return os.RemoveAll(targetPath)
}

// Edge Case 5: Vulnerable JSON stream decoding into map annotations (OCI artifact pattern)
func edgeCase5VulnerableDecoderAnnotations(r io.Reader, baseDir string) error {
	var manifest OCIArtifactManifest
	dec := json.NewDecoder(r)
	if err := dec.Decode(&manifest); err != nil {
		return err
	}
	fileName := manifest.Annotations["org.opencontainers.image.title"]
	dest := filepath.Join(baseDir, fileName)
	// ruleid: path-traversal-unmarshaled-data
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

// Edge Case 6: Patched JSON stream decoding validated with filepath.IsLocal
func edgeCase6PatchedDecoderAnnotationsIsLocal(r io.Reader, baseDir string) error {
	var manifest OCIArtifactManifest
	dec := json.NewDecoder(r)
	if err := dec.Decode(&manifest); err != nil {
		return err
	}
	fileName := manifest.Annotations["org.opencontainers.image.title"]
	if !filepath.IsLocal(fileName) {
		return errors.New("invalid filename: must be a local relative path")
	}
	dest := filepath.Join(baseDir, fileName)
	// ok: path-traversal-unmarshaled-data
	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}
