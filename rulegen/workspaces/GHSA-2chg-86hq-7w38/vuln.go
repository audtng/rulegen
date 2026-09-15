package main

	maxWitnessItemsPerInput = 500000

	// maxWitnessItemSize is the maximum allowed size for an item within
	// an input's witness data. This number is derived from the fact that
	// for script validation, each pushed item onto the stack must be less
	// than 10k bytes.
	maxWitnessItemSize = 11000
)

// TxFlagMarker is the first byte of the FLAG field in a bitcoin tx
// transaction from one that would require a different parsing logic.
//
// Position of FLAG in a bitcoin tx message:
//   ┌─────────┬────────────────────┬─────────────┬─────┐
//   │ VERSION │ FLAG               │ TX-IN-COUNT │ ... │
//   │ 4 bytes │ 2 bytes (optional) │ varint      │     │
//   └─────────┴────────────────────┴─────────────┴─────┘
//
// Zooming into the FLAG field:
//   ┌── FLAG ─────────────┬────────┐
//   │ TxFlagMarker (0x00) │ TxFlag │
//   │ 1 byte              │ 1 byte │
//   └─────────────────────┴────────┘
const TxFlagMarker = 0x00

// TxFlag is the second byte of the FLAG field in a bitcoin tx message.
			// item itself.
			txin.Witness = make([][]byte, witCount)
			for j := uint64(0); j < witCount; j++ {
				txin.Witness[j], err = readScript(r, pver,
					maxWitnessItemSize, "script witness item")
				if err != nil {
					returnScriptBuffers()
					return err
import (
	"fmt"

	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcd/btcutil"
)

const (
