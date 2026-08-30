package main

import (
	"debug/elf"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
)

var pclinetestBinary string

func dotest() bool {
	// For now, only works on ELF platforms.
	if pclinetestBinary != "" {
		return true
	}
	// This command builds pclinetest from pclinetest.asm;
	// the resulting binary looks like it was built from pclinetest.s,
	// but we have renamed it to keep it away from the go tool.
	pclinetestBinary = os.TempDir() + "/pclinetest"
	command := fmt.Sprintf("go tool 6a -o %s.6 pclinetest.asm && go tool 6l -E main -o %s %s.6",
		pclinetestBinary, pclinetestBinary, pclinetestBinary)
	cmd := exec.Command("sh", "-c", command)
	if !dotest() {
		return
	}

	f, tab := crack(pclinetestBinary, t)
	text := f.Section(".text")
