package main

package proxy

import (
	"bytes"
	"context"
	"fmt"
	}
	defer backendConn.Close()

	// Once the connection is hijacked, the ErrorResponder will no longer work, so
	// hijacking should be the last step in the upgrade.
	requestHijacker, ok := w.(http.Hijacker)
		}
	}

	// Proxy the connection. This is bidirectional, so we need a goroutine
	// to copy in each direction. Once one side of the connection exits, we
	// exit the function which performs cleanup and in the process closes
	return dial(updatedReq, h.UpgradeTransport)
}

// dial dials the backend at req.URL and writes req to it.
func dial(req *http.Request, transport http.RoundTripper) (net.Conn, error) {
	conn, err := DialURL(req.Context(), req.URL, transport)
