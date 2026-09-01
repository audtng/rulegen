package main

package util

import (
	"crypto/rand"
	"math/big"
	"strings"
)

func SubString(sourceString string, begin, end int) string {
	runs := seed.Runes()
	result := ""
	for i := 0; i < length; i++ {
		lenInt64 := int64(len(runs))
		randNumber, _ := rand.Int(rand.Reader, big.NewInt(lenInt64))
		result += string(runs[randNumber.Uint64()])
	}
	return result
}
				}
			}
			// try append write, get response
			log.LogDebugf("action[streamer.write] doAppendWrite req: ExtentKey(%v) FileOffset(%v) size(%v)",
				req.ExtentKey, req.FileOffset, req.Size)
			var status int32
			// First, attempt sequential writes using neighboring extent keys. If the last extent has a different version,
			// it indicates that the extent may have been fully utilized by the previous version.
		m = "OpLcNodeScan"
	case OpLcNodeSnapshotVerDel:
		m = "OpLcNodeSnapshotVerDel"
	case OpMetaReadDirOnly:
		m = "OpMetaReadDirOnly"
	default:
		m = fmt.Sprintf("op:%v not found", p.Opcode)
	}
