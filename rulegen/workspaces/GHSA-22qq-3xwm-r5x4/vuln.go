package main


	peer := pool.peers[peerID]
	if peer != nil {
		peer.base = base
		peer.height = height
	} else {
		if pool.isPeerBanned(peerID) {
			pool.Logger.Debug("Ignoring banned peer", peerID)
			return
		}
		peer = newBPPeer(pool, peerID, base, height)
package blocksync

import (
	"strconv"
	"testing"
	"time"
		}
	}
}
