package main

		}

		if strings.HasPrefix(link, "socket:[") {
			inode, err := strconv.ParseInt(link[8:len(link)-1], 10, 64)
			if err != nil {
				logp.Debug("procs", "%s", err.Error())
				continue
			}

			inodes = append(inodes, uint64(inode))
		}
	}

			continue
		}

		uid, _ := strconv.Atoi(string(words[7]))
		sock.uid = uint32(uid)
		inode, _ := strconv.Atoi(string(words[9]))
		sock.inode = uint64(inode)

		sockets = append(sockets, &sock)
	}
		return nil, 0, err
	}

	port, err := strconv.ParseInt(string(words[1]), 16, 32)
	if err != nil {
		return nil, 0, err
	}

func hexToIpv6(word string) (net.IP, error) {
	p := make(net.IP, net.IPv6len)
	for i := 0; i < 4; i++ {
		part, err := strconv.ParseUint(word[i*8:(i+1)*8], 16, 32)
		if err != nil {
			return nil, err
		}
package procs

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/elastic/elastic-agent-libs/logp"
)

func TestFindSocketsOfPid(t *testing.T) {
	logp.TestingSetup()

	proc := []testProcFile{
		{path: "/proc/766/fd/0", isLink: true, contents: "/dev/null"},
	}

	// Create fake proc file system
	pathPrefix, err := ioutil.TempDir("", "find-sockets")
	if err != nil {
		t.Error("TempDir failed:", err)
		return
	assertUint64ArraysAreEqual(t, []uint64{7619, 7620}, inodes)
}

func TestParse_Proc_Net_Tcp(t *testing.T) {
	socketInfo, err := socketsFromProc("../tests/files/proc_net_tcp.txt", false)
	if err != nil {
		}

		if !file.isLink {
			err = ioutil.WriteFile(filepath.Join(prefix, file.path),
				[]byte(file.contents), 0o644)
			if err != nil {
				return err
