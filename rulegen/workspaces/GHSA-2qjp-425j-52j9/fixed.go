package main

package remotecommand

import (
	gocontext "context"
	"encoding/json"
	"errors"
	"fmt"

	if ctx.resizeStream != nil {
		ctx.resizeChan = make(chan remotecommand.TerminalSize)
		go handleResizeEvents(req.Context(), ctx.resizeStream, ctx.resizeChan)
	}

	return ctx, true
// supportsTerminalResizing returns false because v1ProtocolHandler doesn't support it.
func (*v1ProtocolHandler) supportsTerminalResizing() bool { return false }

func handleResizeEvents(ctx gocontext.Context, stream io.Reader, channel chan<- remotecommand.TerminalSize) {
	defer runtime.HandleCrash()
	defer close(channel)

		if err := decoder.Decode(&size); err != nil {
			break
		}

		select {
		case channel <- size:
		case <-ctx.Done():
			// To avoid leaking this routine, exit if the http request finishes. This path
			// would generally be hit if starting the process fails and nothing is started to
			// ingest these resize events.
			return
		}
	}
}

