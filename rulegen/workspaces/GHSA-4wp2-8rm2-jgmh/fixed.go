package main

import "C"

import (
	"errors"
	"fmt"
	"unsafe"
)

// Uncompress with a known output size. len(out) should be equal to
// the length of the uncompressed out.
func Uncompress(in, out []byte) (error) {
	if int(C.LZ4_decompress_safe(p(in), p(out), clen(in), clen(out))) < 0 {
		return errors.New("Malformed compression stream")
	}

	return nil
}

// CompressBound calculates the size of the output buffer needed by
