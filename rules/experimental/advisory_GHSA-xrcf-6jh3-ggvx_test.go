package rules

import (
	"bufio"
	"context"
	"database/sql"
	"fmt"
	"io"
	"net"
	"os"
	"path/filepath"
	"testing"
)

// Case 1: Vulnerable - Server connection input formatted into ATTACH DATABASE query executed on sql.DB
func TestVulnerableAttachDatabase(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	dbPath, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	defer db.Close()

	query := fmt.Sprintf("ATTACH DATABASE '%s' AS attached_db", dbPath)
	// ruleid: go-server-arbitrary-file-write
	_, _ = db.Exec(query)
}

// Case 2: Safe - Server connection input sanitized with filepath.Base before ATTACH query
func TestSafeAttachDatabaseWithBase(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()

	reader := bufio.NewReader(conn)
	rawPath, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	safePath := filepath.Base(rawPath)
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	defer db.Close()

	query := fmt.Sprintf("ATTACH DATABASE '%s' AS safe_db", safePath)
	// ok: go-server-arbitrary-file-write
	_, _ = db.Exec(query)
}

// Case 3: Vulnerable - Scanner reading server stream writing arbitrary file via os.WriteFile
func TestVulnerableFileWriteFromScanner(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	filename := scanner.Text()

	targetPath := filepath.Join("/var/data", filename)
	// ruleid: go-server-arbitrary-file-write
	_ = os.WriteFile(targetPath, []byte("payload"), 0644)
}

// Case 4: Safe - Scanner input validated with filepath.IsLocal before writing file
func TestSafeFileWriteWithIsLocal(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()

	scanner := bufio.NewScanner(conn)
	if !scanner.Scan() {
		return
	}
	filename := scanner.Text()

	if !filepath.IsLocal(filename) {
		return
	}

	targetPath := filepath.Join("/var/data", filename)
	// ok: go-server-arbitrary-file-write
	_ = os.WriteFile(targetPath, []byte("safe content"), 0644)
}

// Case 5: Vulnerable - Reading bytes from net.Conn into sql.Open creating arbitrary DB file
func TestVulnerableSQLOpenFromConn(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return
	}
	defer conn.Close()

	data, err := io.ReadAll(conn)
	if err != nil {
		return
	}

	targetFile := string(data)
	// ruleid: go-server-arbitrary-file-write
	db, _ := sql.Open("sqlite3", targetFile)
	if db != nil {
		defer db.Close()
	}
}

// Case 6: Safe - In-memory and constant database operations / files
func TestSafeConstantOperations(t *testing.T) {
	ctx := context.Background()
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return
	}
	defer db.Close()

	// ok: go-server-arbitrary-file-write
	_, _ = db.ExecContext(ctx, "ATTACH DATABASE ':memory:' AS memdb")

	// ok: go-server-arbitrary-file-write
	f, err := os.Create("/var/log/fixed_audit.log")
	if err == nil {
		f.Close()
	}
}
