package main

import "C"

import (
	"fmt"
	"unsafe"
)

// Uncompress with a known output size. len(out) should be equal to
// the length of the uncompressed out.
func Uncompress(in, out []byte) (err error) {
	read := int(C.LZ4_uncompress(p(in), p(out), clen(out)))

	if read != len(in) {
		err = fmt.Errorf("uncompress read %d bytes should have read %d",
			read, len(in))
	}
	return
}

// CompressBound calculates the size of the output buffer needed by
