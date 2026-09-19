package rules

import (
	"bufio"
	"io"
	"net"
	"os"
	"path/filepath"
)

// Direct stdlib
func DirectStdlib(conn net.Conn) {
	reader := bufio.NewReader(conn)
	path, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	// ruleid: network-path-traversal
	os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
}

// Proper patch
func ProperPatch(conn net.Conn) {
	reader := bufio.NewReader(conn)
	path, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	if !filepath.IsLocal(path) {
		return
	}
	// ok: network-path-traversal
	os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
}

// Cross-function taint
func crossFunctionHelper(p string) string {
	return filepath.Join("/base", p)
}

func CrossFunctionTaint(conn net.Conn) {
	reader := bufio.NewReader(conn)
	path, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	target := crossFunctionHelper(path)
	// ruleid: network-path-traversal
	os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0600)
}

// Interface bypass
func InterfaceBypass(r io.Reader) {
	data, err := io.ReadAll(r)
	if err != nil {
		return
	}
	path := string(data)
	// ruleid: network-path-traversal
	os.OpenFile(path, os.O_CREATE|os.O_WRONLY, 0600)
}

// Fake sanitizer
func FakeSanitizer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	path, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	cleaned := filepath.Clean(path)
	// ruleid: network-path-traversal
	os.OpenFile(cleaned, os.O_CREATE|os.O_WRONLY, 0600)
}

// Real sanitizer
func RealSanitizer(conn net.Conn) {
	reader := bufio.NewReader(conn)
	path, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	base := filepath.Base(path)
	// ok: network-path-traversal
	os.OpenFile(base, os.O_CREATE|os.O_WRONLY, 0600)
}
