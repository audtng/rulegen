package rules

import (
	"encoding/json"
	"os/exec"
	"strconv"
	"strings"
)

type Config struct {
	ExtraArgs []string `json:"extra_args"`
	Command   string   `json:"command"`
}

type ArgProvider interface {
	GetArgs() []string
}

func (c Config) GetArgs() []string {
	return c.ExtraArgs
}

// 1. Direct stdlib: unmarshaled config passed directly to exec.Command
func testDirectStdlib(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	// ruleid: untrusted-config-command-injection
	exec.Command("virtiofsd", cfg.ExtraArgs...)
}

// 2. Proper patch: avoid using untrusted input from config, use static safe arguments
func testProperPatch(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	// ok: untrusted-config-command-injection
	exec.Command("virtiofsd", "--cache=always", "--sandbox=chroot")
}

// 3. Cross-function taint: taint passed to a helper function before sink
func prepareDaemonArgs(raw []string) []string {
	return append([]string{"--syslog"}, raw...)
}

func testCrossFunctionTaint(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	args := prepareDaemonArgs(cfg.ExtraArgs)
	// ruleid: untrusted-config-command-injection
	exec.Command("virtiofsd", args...)
}

// 4. Interface bypass: taint propagated through interface method
func testInterfaceBypass(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	var provider ArgProvider = cfg
	args := provider.GetArgs()
	// ruleid: untrusted-config-command-injection
	exec.Command("virtiofsd", args...)
}

// 5. Fake sanitizer: trimming whitespace or string manipulation doesn't prevent flag/argument injection
func testFakeSanitizer(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	if len(cfg.ExtraArgs) > 0 {
		cleaned := strings.TrimSpace(cfg.ExtraArgs[0])
		// ruleid: untrusted-config-command-injection
		exec.Command("virtiofsd", cleaned)
	}
}

// 6. Real sanitizer: escaping/quoting or allowlist validation
func testRealSanitizer(data []byte) {
	var cfg Config
	json.Unmarshal(data, &cfg)
	if len(cfg.ExtraArgs) > 0 {
		quoted := strconv.Quote(cfg.ExtraArgs[0])
		// ok: untrusted-config-command-injection
		exec.Command("virtiofsd", quoted)
	}
}
