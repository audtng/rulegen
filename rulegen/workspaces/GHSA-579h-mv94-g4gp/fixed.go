package main

package proxy

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	}
	defer backendConn.Close()

	// determine the http response code from the backend by reading from rawResponse+backendConn
	rawResponseCode, headerBytes, err := getResponseCode(io.MultiReader(bytes.NewReader(rawResponse), backendConn))
	if err != nil {
		klog.V(6).Infof("Proxy connection error: %v", err)
		h.Responder.Error(w, req, err)
		return true
	}
	if len(headerBytes) > len(rawResponse) {
		// we read beyond the bytes stored in rawResponse, update rawResponse to the full set of bytes read from the backend
		rawResponse = headerBytes
	}

	// Once the connection is hijacked, the ErrorResponder will no longer work, so
	// hijacking should be the last step in the upgrade.
	requestHijacker, ok := w.(http.Hijacker)
		}
	}

	if rawResponseCode != http.StatusSwitchingProtocols {
		// If the backend did not upgrade the request, finish echoing the response from the backend to the client and return, closing the connection.
		klog.V(6).Infof("Proxy upgrade error, status code %d", rawResponseCode)
		_, err := io.Copy(requestHijackedConn, backendConn)
		if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
			klog.Errorf("Error proxying data from backend to client: %v", err)
		}
		// Indicate we handled the request
		return true
	}

	// Proxy the connection. This is bidirectional, so we need a goroutine
	// to copy in each direction. Once one side of the connection exits, we
	// exit the function which performs cleanup and in the process closes
	return dial(updatedReq, h.UpgradeTransport)
}

// getResponseCode reads a http response from the given reader, returns the status code,
// the bytes read from the reader, and any error encountered
func getResponseCode(r io.Reader) (int, []byte, error) {
	rawResponse := bytes.NewBuffer(make([]byte, 0, 256))
	// Save the bytes read while reading the response headers into the rawResponse buffer
	resp, err := http.ReadResponse(bufio.NewReader(io.TeeReader(r, rawResponse)), nil)
	if err != nil {
		return 0, nil, err
	}
	// return the http status code and the raw bytes consumed from the reader in the process
	return resp.StatusCode, rawResponse.Bytes(), nil
}

// dial dials the backend at req.URL and writes req to it.
func dial(req *http.Request, transport http.RoundTripper) (net.Conn, error) {
	conn, err := DialURL(req.Context(), req.URL, transport)
