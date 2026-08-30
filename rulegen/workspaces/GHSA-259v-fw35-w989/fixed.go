package main

import (
	"debug/elf"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var (
	pclineTempDir    string
	pclinetestBinary string
)

func dotest() bool {
	// For now, only works on ELF platforms.
	if pclinetestBinary != "" {
		return true
	}
	var err error
	pclineTempDir, err = ioutil.TempDir("", "pclinetest")
	if err != nil {
		panic(err)
	}
	if strings.Contains(pclineTempDir, " ") {
		panic("unexpected space in tempdir")
	}
	// This command builds pclinetest from pclinetest.asm;
	// the resulting binary looks like it was built from pclinetest.s,
	// but we have renamed it to keep it away from the go tool.
	pclinetestBinary = filepath.Join(pclineTempDir, "pclinetest")
	command := fmt.Sprintf("go tool 6a -o %s.6 pclinetest.asm && go tool 6l -E main -o %s %s.6",
		pclinetestBinary, pclinetestBinary, pclinetestBinary)
	cmd := exec.Command("sh", "-c", command)
	if !dotest() {
		return
	}
	defer os.RemoveAll(pclineTempDir)

	f, tab := crack(pclinetestBinary, t)
	text := f.Section(".text")
