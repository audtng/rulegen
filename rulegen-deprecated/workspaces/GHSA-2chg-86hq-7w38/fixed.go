package main

	maxWitnessItemsPerInput = 500000

	// maxWitnessItemSize is the maximum allowed size for an item within
	// an input's witness data. This value is bounded by the largest
	// possible block size, post segwit v1 (taproot).
	maxWitnessItemSize = 4_000_000
)

// TxFlagMarker is the first byte of the FLAG field in a bitcoin tx
// transaction from one that would require a different parsing logic.
//
// Position of FLAG in a bitcoin tx message:
//
//	┌─────────┬────────────────────┬─────────────┬─────┐
//	│ VERSION │ FLAG               │ TX-IN-COUNT │ ... │
//	│ 4 bytes │ 2 bytes (optional) │ varint      │     │
//	└─────────┴────────────────────┴─────────────┴─────┘
//
// Zooming into the FLAG field:
//
//	┌── FLAG ─────────────┬────────┐
//	│ TxFlagMarker (0x00) │ TxFlag │
//	│ 1 byte              │ 1 byte │
//	└─────────────────────┴────────┘
const TxFlagMarker = 0x00

// TxFlag is the second byte of the FLAG field in a bitcoin tx message.
			// item itself.
			txin.Witness = make([][]byte, witCount)
			for j := uint64(0); j < witCount; j++ {
				txin.Witness[j], err = readScript(
					r, pver, maxWitnessItemSize, "script witness item",
				)
				if err != nil {
					returnScriptBuffers()
					return err
import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
)

const (
