package main

)

const (
	mapEntryCost  = 8
	listEntryCost = 4
)

// This file should be identical to the general feature in the parent package,
	// entire block to be part of the CBOR object and will error if there is
	// extraneous data after the end of the object.
	DontParseBeyondEnd bool

	// AllowBudget sets the maximum budget for the decoder. The budget is
	// decremented as the decoder allocates resources (nodes, map entries, list
	// elements, string/bytes content). If the budget is exhausted, the decoder
	// returns ErrAllocationBudgetExceeded.
	//
	// When zero, a default budget is used which is generous for typical IPLD
	// block sizes.
	AllocationBudget int64

	// MaxCollectionPrealloc sets the maximum size hint passed to
	// BeginMap/BeginList. CBOR headers declare collection sizes upfront;
	// this caps the initial allocation while collections grow dynamically
	// as entries are decoded.
	//
	// When zero, a default of 1024 is used.
	MaxCollectionPrealloc int64
}

const (
	defaultAllocationBudget      int64 = 1048576 * 10
	defaultMaxCollectionPrealloc int64 = 1024
)

// Decode deserializes data from the given io.Reader and feeds it into the given datamodel.NodeAssembler.
// Decode fits the codec.Decoder function interface.
//
// Unmarshal is a deprecated function.
// Please consider switching to DecodeOptions.Decode instead.
func Unmarshal(na datamodel.NodeAssembler, tokSrc shared.TokenSource, options DecodeOptions) error {
	budget := options.AllocationBudget
	if budget == 0 {
		budget = defaultAllocationBudget
	}
	return unmarshal1(na, tokSrc, &budget, options)
}

func (cfg DecodeOptions) maxPrealloc() int64 {
	if cfg.MaxCollectionPrealloc > 0 {
		return cfg.MaxCollectionPrealloc
	}
	return defaultMaxCollectionPrealloc
}

func unmarshal1(na datamodel.NodeAssembler, tokSrc shared.TokenSource, budget *int64, options DecodeOptions) error {
	var tk tok.Token
	done, err := tokSrc.Step(&tk)
	if err == io.EOF {
	if done && !tk.Type.IsValue() && tk.Type != tok.TNull {
		return fmt.Errorf("unexpected eof")
	}
	return unmarshal2(na, tokSrc, &tk, budget, options)
}

// starts with the first token already primed.  Necessary to get recursion
//
//	to flow right without a peek+unpeek system.
func unmarshal2(na datamodel.NodeAssembler, tokSrc shared.TokenSource, tk *tok.Token, budget *int64, options DecodeOptions) error {
	// FUTURE: check for schema.TypedNodeBuilder that's going to parse a Link (they can slurp any token kind they want).
	switch tk.Type {
	case tok.TMapOpen:
			expectLen = math.MaxInt64
			allocLen = 0
		} else {
			*budget -= allocLen
			if *budget < 0 {
				return ErrAllocationBudgetExceeded
			}
			if allocLen > options.maxPrealloc() {
				allocLen = options.maxPrealloc()
			}
		}
		ma, err := na.BeginMap(allocLen)
		if err != nil {
				}
				return ma.Finish()
			case tok.TString:
				*budget -= int64(len(tk.Str) + mapEntryCost)
				if *budget < 0 {
					return ErrAllocationBudgetExceeded
				}
				// continue
			if err != nil { // return in error if the key was rejected
				return err
			}
			err = unmarshal1(mva, tokSrc, budget, options)
			if err != nil { // return in error if some part of the recursion errored
				return err
			}
			expectLen = math.MaxInt64
			allocLen = 0
		} else {
			*budget -= allocLen
			if *budget < 0 {
				return ErrAllocationBudgetExceeded
			}
			if allocLen > options.maxPrealloc() {
				allocLen = options.maxPrealloc()
			}
		}
		la, err := na.BeginList(allocLen)
		if err != nil {
				}
				return la.Finish()
			default:
				*budget -= listEntryCost
				if *budget < 0 {
					return ErrAllocationBudgetExceeded
				}
				observedLen++
				if observedLen > expectLen {
					return fmt.Errorf("unexpected continuation of array elements beyond declared length")
				}
				err := unmarshal2(la.AssembleValue(), tokSrc, tk, budget, options)
				if err != nil { // return in error if some part of the recursion errored
					return err
				}
	case tok.TNull:
		return na.AssignNull()
	case tok.TString:
		*budget -= int64(len(tk.Str))
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignString(tk.Str)
	case tok.TBytes:
		*budget -= int64(len(tk.Bytes))
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		if !tk.Tagged {
			return fmt.Errorf("unhandled cbor tag %d", tk.Tag)
		}
	case tok.TBool:
		*budget -= 1
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignBool(tk.Bool)
	case tok.TInt:
		*budget -= 1
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignInt(tk.Int)
	case tok.TUint:
		*budget -= 1
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		// note that this pushes any overflow errors up the stack when AsInt() may
		}
		return na.AssignInt(int64(tk.Uint))
	case tok.TFloat64:
		*budget -= 1
		if *budget < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignFloat(tk.Float64)
package dagcbor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"runtime"
	"strings"
	"testing"
		qt.Assert(t, err, qt.Equals, ErrTrailingBytes)
	})
}

func cborMapHeader(length uint32) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0xBA)
	binary.Write(&buf, binary.BigEndian, length)
	return buf.Bytes()
}

func cborArrayHeader(length uint32) []byte {
	var buf bytes.Buffer
	buf.WriteByte(0x9A)
	binary.Write(&buf, binary.BigEndian, length)
	return buf.Bytes()
}

// TestDecodeOptions_AllocationBudget verifies that the configurable allocation
// budget is respected, both with defaults and custom values.
func TestDecodeOptions_AllocationBudget(t *testing.T) {
	t.Run("default budget rejects oversized structure", func(t *testing.T) {
		// A map header claiming more entries than the default budget allows
		payload := cborMapHeader(20_000_000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})

	t.Run("custom budget accepts within limit", func(t *testing.T) {
		// Build a small valid map: {"a": 1}
		var buf bytes.Buffer
		buf.Write(cborMapHeader(1))
		buf.WriteByte(0x61) // text(1)
		buf.WriteByte('a')
		buf.WriteByte(0x01) // uint(1)

		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 100}.Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)
		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_Map)
	})

	t.Run("custom budget rejects when exhausted", func(t *testing.T) {
		// A map claiming 50 entries with a budget of only 10 should be rejected
		payload := cborMapHeader(50)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 10}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})

	t.Run("budget accounts for declared collection sizes", func(t *testing.T) {
		// A list claiming 1000 entries consumes budget even if no entries follow
		payload := cborArrayHeader(1000)
		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, AllocationBudget: 500}.Decode(nb, bytes.NewReader(payload))
		qt.Assert(t, err, qt.Equals, ErrAllocationBudgetExceeded)
	})
}

// TestDecodeOptions_MaxCollectionPrealloc verifies that the preallocation cap
// is respected and that large collections still decode correctly.
func TestDecodeOptions_MaxCollectionPrealloc(t *testing.T) {
	t.Run("large map decodes correctly with default cap", func(t *testing.T) {
		const numEntries = 5_000
		var buf bytes.Buffer
		buf.Write(cborMapHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			key := fmt.Sprintf("k%05d", i)
			buf.WriteByte(0x66) // text(6)
			buf.WriteString(key)
			if i < 24 {
				buf.WriteByte(byte(i))
			} else if i < 256 {
				buf.WriteByte(0x18)
				buf.WriteByte(byte(i))
			} else {
				buf.WriteByte(0x19)
				binary.Write(&buf, binary.BigEndian, uint16(i))
			}
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_Map)
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))

		v, err := node.LookupByString("k00000")
		qt.Assert(t, err, qt.IsNil)
		vi, err := v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(0))

		v, err = node.LookupByString("k04999")
		qt.Assert(t, err, qt.IsNil)
		vi, err = v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(4999))
	})

	t.Run("large list decodes correctly with default cap", func(t *testing.T) {
		const numEntries = 5_000
		var buf bytes.Buffer
		buf.Write(cborArrayHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			if i < 24 {
				buf.WriteByte(byte(i))
			} else if i < 256 {
				buf.WriteByte(0x18)
				buf.WriteByte(byte(i))
			} else {
				buf.WriteByte(0x19)
				binary.Write(&buf, binary.BigEndian, uint16(i))
			}
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Kind(), qt.Equals, datamodel.Kind_List)
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))

		v, err := node.LookupByIndex(0)
		qt.Assert(t, err, qt.IsNil)
		vi, err := v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(0))

		v, err = node.LookupByIndex(numEntries - 1)
		qt.Assert(t, err, qt.IsNil)
		vi, err = v.AsInt()
		qt.Assert(t, err, qt.IsNil)
		qt.Assert(t, vi, qt.Equals, int64(numEntries-1))
	})

	t.Run("custom prealloc cap with valid data", func(t *testing.T) {
		// 100-entry list with a prealloc cap of 10 should still decode fine
		const numEntries = 100
		var buf bytes.Buffer
		buf.Write(cborArrayHeader(numEntries))
		for i := 0; i < numEntries; i++ {
			buf.WriteByte(byte(i % 24))
		}

		nb := basicnode.Prototype.Any.NewBuilder()
		err := DecodeOptions{AllowLinks: true, MaxCollectionPrealloc: 10}.Decode(nb, &buf)
		qt.Assert(t, err, qt.IsNil)

		node := nb.Build()
		qt.Assert(t, node.Length(), qt.Equals, int64(numEntries))
	})
}
