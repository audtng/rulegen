package main

)

const (
	mapEntryGasScore  = 8
	listEntryGasScore = 4
)

// This file should be identical to the general feature in the parent package,
	// entire block to be part of the CBOR object and will error if there is
	// extraneous data after the end of the object.
	DontParseBeyondEnd bool
}

// Decode deserializes data from the given io.Reader and feeds it into the given datamodel.NodeAssembler.
// Decode fits the codec.Decoder function interface.
//
// Unmarshal is a deprecated function.
// Please consider switching to DecodeOptions.Decode instead.
func Unmarshal(na datamodel.NodeAssembler, tokSrc shared.TokenSource, options DecodeOptions) error {
	// Have a gas budget, which will be decremented as we allocate memory, and an error returned when exceeded (or about to be exceeded).
	//  This is a DoS defense mechanism.
	//  It's *roughly* in units of bytes (but only very, VERY roughly) -- it also treats words as 1 in many cases.
	// FUTURE: this ought be configurable somehow.  (How, and at what granularity though?)
	var gas int64 = 1048576 * 10
	return unmarshal1(na, tokSrc, &gas, options)
}

func unmarshal1(na datamodel.NodeAssembler, tokSrc shared.TokenSource, gas *int64, options DecodeOptions) error {
	var tk tok.Token
	done, err := tokSrc.Step(&tk)
	if err == io.EOF {
	if done && !tk.Type.IsValue() && tk.Type != tok.TNull {
		return fmt.Errorf("unexpected eof")
	}
	return unmarshal2(na, tokSrc, &tk, gas, options)
}

// starts with the first token already primed.  Necessary to get recursion
//
//	to flow right without a peek+unpeek system.
func unmarshal2(na datamodel.NodeAssembler, tokSrc shared.TokenSource, tk *tok.Token, gas *int64, options DecodeOptions) error {
	// FUTURE: check for schema.TypedNodeBuilder that's going to parse a Link (they can slurp any token kind they want).
	switch tk.Type {
	case tok.TMapOpen:
			expectLen = math.MaxInt64
			allocLen = 0
		} else {
			if *gas-allocLen < 0 { // halt early if this will clearly demand too many resources
				return ErrAllocationBudgetExceeded
			}
		}
		ma, err := na.BeginMap(allocLen)
		if err != nil {
				}
				return ma.Finish()
			case tok.TString:
				*gas -= int64(len(tk.Str) + mapEntryGasScore)
				if *gas < 0 {
					return ErrAllocationBudgetExceeded
				}
				// continue
			if err != nil { // return in error if the key was rejected
				return err
			}
			err = unmarshal1(mva, tokSrc, gas, options)
			if err != nil { // return in error if some part of the recursion errored
				return err
			}
			expectLen = math.MaxInt64
			allocLen = 0
		} else {
			if *gas-allocLen < 0 { // halt early if this will clearly demand too many resources
				return ErrAllocationBudgetExceeded
			}
		}
		la, err := na.BeginList(allocLen)
		if err != nil {
				}
				return la.Finish()
			default:
				*gas -= listEntryGasScore
				if *gas < 0 {
					return ErrAllocationBudgetExceeded
				}
				observedLen++
				if observedLen > expectLen {
					return fmt.Errorf("unexpected continuation of array elements beyond declared length")
				}
				err := unmarshal2(la.AssembleValue(), tokSrc, tk, gas, options)
				if err != nil { // return in error if some part of the recursion errored
					return err
				}
	case tok.TNull:
		return na.AssignNull()
	case tok.TString:
		*gas -= int64(len(tk.Str))
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignString(tk.Str)
	case tok.TBytes:
		*gas -= int64(len(tk.Bytes))
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		if !tk.Tagged {
			return fmt.Errorf("unhandled cbor tag %d", tk.Tag)
		}
	case tok.TBool:
		*gas -= 1
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignBool(tk.Bool)
	case tok.TInt:
		*gas -= 1
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignInt(tk.Int)
	case tok.TUint:
		*gas -= 1
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		// note that this pushes any overflow errors up the stack when AsInt() may
		}
		return na.AssignInt(int64(tk.Uint))
	case tok.TFloat64:
		*gas -= 1
		if *gas < 0 {
			return ErrAllocationBudgetExceeded
		}
		return na.AssignFloat(tk.Float64)
package dagcbor

import (
	"runtime"
	"strings"
	"testing"
		qt.Assert(t, err, qt.Equals, ErrTrailingBytes)
	})
}
