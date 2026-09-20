package rules

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"strconv"
)

type Config struct {
	Binary    string   `json:"binary"`
	ExtraArgs []string `json:"extra_args"`
	Timeout   string   `json:"timeout"`
	Mode      string   `json:"mode"`
}

// 1. Vulnerable: unmarshaled extra args passed directly as command arguments
func testVulnerableArgs(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	// ruleid: go-exec-untrusted-config-injection
	cmd := exec.Command("virtiofsd", cfg.ExtraArgs...)
	return cmd.Run()
}

// 2. Vulnerable: unmarshaled binary path passed to exec.CommandContext
func testVulnerableCommand(ctx context.Context, data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	// ruleid: go-exec-untrusted-config-injection
	cmd := exec.CommandContext(ctx, cfg.Binary, "--daemon")
	return cmd.Run()
}

// 3. Vulnerable: unmarshaled args propagated via fmt.Sprintf and slice append
func testVulnerablePropagated(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	args := []string{"--syslog"}
	for _, extra := range cfg.ExtraArgs {
		formatted := fmt.Sprintf("--extra-opt=%s", extra)
		args = append(args, formatted)
	}
	// ruleid: go-exec-untrusted-config-injection
	cmd := exec.Command("virtiofsd", args...)
	return cmd.Run()
}

// 4. Safe: static hardcoded command and arguments
func testSafeStaticCommand() error {
	// ok: go-exec-untrusted-config-injection
	cmd := exec.Command("virtiofsd", "--syslog", "--cache=always")
	return cmd.Run()
}

// 5. Safe: input sanitized through strconv.Atoi before being used
func testSafeSanitizedNumeric(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	timeout, err := strconv.Atoi(cfg.Timeout)
	if err != nil {
		return err
	}
	// ok: go-exec-untrusted-config-injection
	cmd := exec.Command("virtiofsd", "--timeout", strconv.Itoa(timeout))
	return cmd.Run()
}

// 6. Safe: input validated via switch statement mapping to constant flags
func testSafeMappedConstant(data []byte) error {
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}
	var modeFlag string
	switch cfg.Mode {
	case "always":
		modeFlag = "--cache=always"
	case "never":
		modeFlag = "--cache=never"
	default:
		modeFlag = "--cache=auto"
	}
	// ok: go-exec-untrusted-config-injection
	cmd := exec.Command("virtiofsd", modeFlag)
	return cmd.Run()
}
