package main

package util

import (
	"math/rand"
	"strings"
	"time"
)

func SubString(sourceString string, begin, end int) string {
	runs := seed.Runes()
	result := ""
	for i := 0; i < length; i++ {
		rand.Seed(time.Now().UnixNano())
		randNumber := rand.Intn(len(runs))
		result += string(runs[randNumber])
	}
	return result
}
				}
			}
			// try append write, get response
			log.LogDebugf("action[streamer.write] doAppendWrite req %v FileOffset %v size %v", req.ExtentKey, req.FileOffset, req.Size)
			var status int32
			// First, attempt sequential writes using neighboring extent keys. If the last extent has a different version,
			// it indicates that the extent may have been fully utilized by the previous version.
		m = "OpLcNodeScan"
	case OpLcNodeSnapshotVerDel:
		m = "OpLcNodeSnapshotVerDel"
	default:
		m = fmt.Sprintf("op:%v not found", p.Opcode)
	}
