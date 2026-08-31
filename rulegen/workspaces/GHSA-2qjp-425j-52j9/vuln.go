package main

package remotecommand

import (
	"encoding/json"
	"errors"
	"fmt"

	if ctx.resizeStream != nil {
		ctx.resizeChan = make(chan remotecommand.TerminalSize)
		go handleResizeEvents(ctx.resizeStream, ctx.resizeChan)
	}

	return ctx, true
// supportsTerminalResizing returns false because v1ProtocolHandler doesn't support it.
func (*v1ProtocolHandler) supportsTerminalResizing() bool { return false }

func handleResizeEvents(stream io.Reader, channel chan<- remotecommand.TerminalSize) {
	defer runtime.HandleCrash()
	defer close(channel)

		if err := decoder.Decode(&size); err != nil {
			break
		}
		channel <- size
	}
}

