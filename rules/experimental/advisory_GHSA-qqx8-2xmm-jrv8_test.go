package rules

import (
	"encoding/json"
	"os"
	"path"
	"path/filepath"
)

type ChallengePayload struct {
	Token    string `json:"token"`
	Filename string `json:"filename"`
}

func testVulnerableLegoChallenge(data []byte, webroot string) error {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	challengePath := path.Join(".well-known", "acme-challenge", payload.Token)
	dest := filepath.Join(webroot, challengePath)
	// ruleid: path-traversal-file-write
	return os.WriteFile(dest, []byte("keyAuth"), 0644)
}

func testVulnerableRemoveFile(data []byte, webroot string) error {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	dest := filepath.Join(webroot, payload.Token)
	// ruleid: path-traversal-file-write
	return os.Remove(dest)
}

func testVulnerableOpenFile(data []byte, webroot string) (*os.File, error) {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, err
	}
	dest := filepath.Join(webroot, payload.Filename)
	// ruleid: path-traversal-file-write
	return os.OpenFile(dest, os.O_CREATE|os.O_WRONLY, 0644)
}

func testSafeBaseSanitizer(data []byte, webroot string) error {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	safeName := filepath.Base(payload.Filename)
	dest := filepath.Join(webroot, safeName)
	// ok: path-traversal-file-write
	return os.WriteFile(dest, []byte("safe"), 0644)
}

func testSafeIsLocalSanitizer(data []byte, webroot string) error {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	if !filepath.IsLocal(payload.Filename) {
		return os.ErrInvalid
	}
	dest := filepath.Join(webroot, payload.Filename)
	// ok: path-traversal-file-write
	return os.WriteFile(dest, []byte("safe"), 0644)
}

func testSafeConstantPath(data []byte, webroot string) error {
	var payload ChallengePayload
	if err := json.Unmarshal(data, &payload); err != nil {
		return err
	}
	dest := filepath.Join(webroot, "fixed_challenge.txt")
	// ok: path-traversal-file-write
	return os.WriteFile(dest, []byte("fixed"), 0644)
}
