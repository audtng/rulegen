package main

import (
	"unsafe"
)

type ByteSliceView struct {
	is_nil bool
	ptr    *byte
	len    uint
}

func makeView(s []byte) ByteSliceView {
	if s == nil {
		return ByteSliceView{is_nil: true, ptr: nil, len: 0}
	}
	if len(s) == 0 {
		return ByteSliceView{is_nil: false, ptr: nil, len: 0}
	}
	return ByteSliceView{
		is_nil: false,
		ptr:    (*byte)(unsafe.Pointer(&s[0])),
		len:    uint(len(s)),
	}
}
