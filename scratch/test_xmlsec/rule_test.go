package rules

import (
	"context"
	"os/exec"
)

// 1. Direct xmlsec1 verify without --enabled-key-data (vulnerable)
func verifyDirectVuln(xmlPath, certPath string) (*exec.Cmd, error) {
	// ruleid: xmlsec-verify-missing-enabled-key-data
	cmd := exec.Command("xmlsec1", "--verify", "--pubkey-cert-pem", certPath, xmlPath)
	return cmd, nil
}

// 2. Direct xmlsec1 verify with --enabled-key-data x509 (safe)
func verifyDirectSafe(xmlPath, certPath string) (*exec.Cmd, error) {
	// ok: xmlsec-verify-missing-enabled-key-data
	cmd := exec.Command("xmlsec1", "--verify", "--enabled-key-data", "x509", "--pubkey-cert-pem", certPath, xmlPath)
	return cmd, nil
}

// 3. CommandContext verify without --enabled-key-data (vulnerable)
func verifyContextVuln(ctx context.Context, xmlPath, certPath string) (*exec.Cmd, error) {
	// ruleid: xmlsec-verify-missing-enabled-key-data
	cmd := exec.CommandContext(ctx, "xmlsec1", "verify", "--pubkey-cert-pem", certPath, xmlPath)
	return cmd, nil
}

// 4. CommandContext verify with --enabled-key-data flag (safe)
func verifyContextSafe(ctx context.Context, xmlPath, certPath string) (*exec.Cmd, error) {
	// ok: xmlsec-verify-missing-enabled-key-data
	cmd := exec.CommandContext(ctx, "xmlsec1", "--verify", "--enabled-key-data", "raw-x509-cert", "--pubkey-cert-pem", certPath, xmlPath)
	return cmd, nil
}

// 5. Binary path /usr/bin/xmlsec1 verify without --enabled-key-data (vulnerable)
func verifyCustomPathVuln(xmlPath string) *exec.Cmd {
	// ruleid: xmlsec-verify-missing-enabled-key-data
	return exec.Command("/usr/bin/xmlsec1", "--verify", xmlPath)
}

// 6. Binary path with --enabled-key-data (safe)
func verifyCustomPathSafe(xmlPath string) *exec.Cmd {
	// ok: xmlsec-verify-missing-enabled-key-data
	return exec.Command("/usr/bin/xmlsec1", "--verify", "--enabled-key-data", "raw-x509-cert", xmlPath)
}

// 7. Non-verify command (e.g. sign) (safe)
func signDoc(xmlPath, keyPath string) *exec.Cmd {
	// ok: xmlsec-verify-missing-enabled-key-data
	return exec.Command("xmlsec1", "--sign", "--privkey-pem", keyPath, xmlPath)
}

// 8. Args slice without --enabled-key-data (vulnerable)
func verifyArgsSliceVuln(xmlPath string) *exec.Cmd {
	args := []string{"--verify", "--pubkey-cert-pem", "cert.pem", xmlPath}
	// ruleid: xmlsec-verify-missing-enabled-key-data
	return exec.Command("xmlsec1", args...)
}

// 9. Args slice with --enabled-key-data (safe)
func verifyArgsSliceSafe(xmlPath string) *exec.Cmd {
	args := []string{"--verify", "--enabled-key-data", "x509", "--pubkey-cert-pem", "cert.pem", xmlPath}
	// ok: xmlsec-verify-missing-enabled-key-data
	return exec.Command("xmlsec1", args...)
}
