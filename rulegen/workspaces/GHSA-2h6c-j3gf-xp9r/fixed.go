package main

// NOTE: Don't bother replacing the divisions/modulo with shifts/ands, go is smart.

import (
	"fmt"
	"math/bits"
)

// NewBitfield creates a new fixed-sized Bitfield (allocated up-front).
func NewBitfield(size int) (Bitfield, error) {
	if size < 0 {
		return nil, fmt.Errorf("bitfield size must be positive; got %d", size)
	}
	if size%8 != 0 {
		return nil, fmt.Errorf("bitfield size must be a multiple of 8; got %d", size)
	}
	return make([]byte, size/8), nil
}

// FromBytes constructs a new bitfield from a serialized bitfield.
func FromBytes(size int, bits []byte) (Bitfield, error) {
	bf, err := NewBitfield(size)
	if err != nil {
		return nil, err
	}
	start := len(bf) - len(bits)
	if start < 0 {
		return nil, fmt.Errorf("bitfield too small: got %d; need %d", size, len(bits)*8)
	}
	copy(bf[start:], bits)
	return bf, nil
}

func (bf Bitfield) offset(i int) (uint, uint8) {
)

func TestExhaustive24(t *testing.T) {
	bf, err := NewBitfield(24)
	assertNoError(t, err)
	max := 1 << 24

	bint := new(big.Int)
}

func TestBitfield(t *testing.T) {
	bf, err := NewBitfield(128)
	assertNoError(t, err)
	if bf.OnesBefore(20) != 0 {
		t.Fatal("expected no bits set")
	}
	}
}

func TestBadSizeFails(t *testing.T) {
	for _, size := range [...]int{-8, 2, 1337, -3} {
		_, err := NewBitfield(size)
		if err == nil {
			t.Fatalf("missing error for %d sized bitfield", size)
		}
	}
}

var benchmarkSize = 512

func BenchmarkBitfield(t *testing.B) {
	bf, err := NewBitfield(benchmarkSize)
	assertNoError(t, err)
	t.ResetTimer()
	for i := 0; i < t.N; i++ {
		if bf.Bit(i % benchmarkSize) {
	}
}

func BenchmarkOnes(b *testing.B) {
	bf, err := NewBitfield(benchmarkSize)
	assertNoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j*4 < benchmarkSize; j++ {
			if bf.Ones() != j {
				b.Fatal("bad", i)
			}
			bf.SetBit(j * 4)
		}
	}
}

func BenchmarkBytes(b *testing.B) {
	bfa, err := NewBitfield(216)
	assertNoError(b, err)
	bfb, err := NewBitfield(216)
	assertNoError(b, err)
	for j := 0; j*4 < 216; j++ {
		bfa.SetBit(j * 4)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		bfb.SetBytes(bfa.Bytes())
	}
}
		}
	}
}

func FuzzFromBytes(f *testing.F) {
	f.Fuzz(func(_ *testing.T, size int, bytes []byte) {
		if size > 1<<20 { // We relly on consumers for limit checks, hopefully they understand that a New... factory allocates memory.
			return
		}
		FromBytes(size, bytes)
	})
}

func assertNoError(t testing.TB, e error) {
	if e != nil {
		t.Fatal(e)
	}
}
