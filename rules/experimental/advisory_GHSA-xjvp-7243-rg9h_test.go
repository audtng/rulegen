package rules

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
)

func handleSCPStreamWriteVuln(r *bufio.Reader, baseDir string) error {
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	clean := filepath.Clean(line)
	target := filepath.Join(baseDir, clean)
	// ruleid: generic-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func handleSCPStreamWriteSafe(r *bufio.Reader, baseDir string) error {
	line, err := r.ReadString('\n')
	if err != nil {
		return err
	}
	base := filepath.Base(line)
	target := filepath.Join(baseDir, base)
	// ok: generic-path-traversal
	f, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	return f.Close()
}

func handleScannerReadVuln(s *bufio.Scanner, root string) ([]byte, error) {
	s.Scan()
	filename := s.Text()
	target := filepath.Join(root, filename)
	// ruleid: generic-path-traversal
	return os.ReadFile(target)
}

func handleScannerReadSafe(s *bufio.Scanner, root string) ([]byte, error) {
	s.Scan()
	filename := s.Text()
	if !filepath.IsLocal(filename) {
		return nil, os.ErrInvalid
	}
	target := filepath.Join(root, filename)
	// ok: generic-path-traversal
	return os.ReadFile(target)
}

func handleConnRemoveVuln(conn net.Conn, root string) error {
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}
	target := filepath.Join(root, string(buf))
	// ruleid: generic-path-traversal
	return os.Remove(target)
}

func handleConnRemoveSafe(conn net.Conn, root string) error {
	buf := make([]byte, 256)
	_, err := conn.Read(buf)
	if err != nil {
		return err
	}
	path := string(buf)
	if !filepath.IsLocal(path) {
		return os.ErrInvalid
	}
	target := filepath.Join(root, path)
	// ok: generic-path-traversal
	return os.Remove(target)
}
