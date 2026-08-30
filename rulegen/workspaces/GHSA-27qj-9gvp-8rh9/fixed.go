package main

		}

		if strings.HasPrefix(link, "socket:[") {
			inode, err := strconv.ParseUint(link[8:len(link)-1], 10, 64)
			if err != nil {
				logp.Debug("procs", "%s", err.Error())
				continue
			}

			inodes = append(inodes, inode)
		}
	}

			continue
		}

		uid, _ := strconv.ParseUint(string(words[7]), 10, 32)
		sock.uid = uint32(uid)
		inode, _ := strconv.ParseUint(string(words[9]), 10, 64)
		sock.inode = inode

		sockets = append(sockets, &sock)
	}
		return nil, 0, err
	}

	port, err := strconv.ParseUint(string(words[1]), 16, 16)
	if err != nil {
		return nil, 0, err
	}

func hexToIpv6(word string) (net.IP, error) {
	p := make(net.IP, net.IPv6len)
	if len(word) < 32 {
		return nil, fmt.Errorf("got ip6 address of invalid length: %d", len(word))
	}
	for i := 0; i < 4; i++ {
		start := i * 8
		end := (i + 1) * 8
		part, err := strconv.ParseUint(word[start:end], 16, 32)
		if err != nil {
			return nil, err
		}
package procs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFindSocketsOfPid(t *testing.T) {

	proc := []testProcFile{
		{path: "/proc/766/fd/0", isLink: true, contents: "/dev/null"},
	}

	// Create fake proc file system
	pathPrefix, err := os.MkdirTemp("", "find-sockets")
	if err != nil {
		t.Error("TempDir failed:", err)
		return
	assertUint64ArraysAreEqual(t, []uint64{7619, 7620}, inodes)
}

func TestParseInvalidIp6Addr(t *testing.T) {
	_, err := hexToIpv6("B80D012000000000FFFF2301EFCD")
	require.Error(t, err)
}

func TestParse_Proc_Net_Tcp(t *testing.T) {
	socketInfo, err := socketsFromProc("../tests/files/proc_net_tcp.txt", false)
	if err != nil {
		}

		if !file.isLink {
			err = os.WriteFile(filepath.Join(prefix, file.path),
				[]byte(file.contents), 0o644)
			if err != nil {
				return err
