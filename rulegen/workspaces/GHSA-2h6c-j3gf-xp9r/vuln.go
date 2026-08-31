package main

// NOTE: Don't bother replacing the divisions/modulo with shifts/ands, go is smart.

import (
	"math/bits"
)

// NewBitfield creates a new fixed-sized Bitfield (allocated up-front).
//
// Panics if size is not a multiple of 8.
func NewBitfield(size int) Bitfield {
	if size%8 != 0 {
		panic("Bitfield size must be a multiple of 8")
	}
	return make([]byte, size/8)
}

// FromBytes constructs a new bitfield from a serialized bitfield.
func FromBytes(size int, bits []byte) Bitfield {
	bf := NewBitfield(size)
	start := len(bf) - len(bits)
	if start < 0 {
		panic("bitfield too small")
	}
	copy(bf[start:], bits)
	return bf
}

func (bf Bitfield) offset(i int) (uint, uint8) {
)

func TestExhaustive24(t *testing.T) {
	bf := NewBitfield(24)
	max := 1 << 24

	bint := new(big.Int)
}

func TestBitfield(t *testing.T) {
	bf := NewBitfield(128)
	if bf.OnesBefore(20) != 0 {
		t.Fatal("expected no bits set")
	}
	}
}

var benchmarkSize = 512

func BenchmarkBitfield(t *testing.B) {
	bf := NewBitfield(benchmarkSize)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		if bf.Bit(i % benchmarkSize) {
	}
}

func BenchmarkOnes(t *testing.B) {
	bf := NewBitfield(benchmarkSize)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		for j := 0; j*4 < benchmarkSize; j++ {
			if bf.Ones() != j {
				t.Fatal("bad", i)
			}
			bf.SetBit(j * 4)
		}
	}
}

func BenchmarkBytes(t *testing.B) {
	bfa := NewBitfield(211)
	bfb := NewBitfield(211)
	for j := 0; j*4 < 211; j++ {
		bfa.SetBit(j * 4)
	}
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		bfb.SetBytes(bfa.Bytes())
	}
}
		}
	}
}
