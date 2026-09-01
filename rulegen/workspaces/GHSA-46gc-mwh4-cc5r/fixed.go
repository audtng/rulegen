package main

	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	}
}

// isAllowedOrigin validates that the origin is from localhost.
// Returns true if the origin's hostname is "localhost" or "127.0.0.1" (any port allowed).
func isAllowedOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
		return false // Invalid URL format
	}

	// Only allow http or https schemes
	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}

	// Extract hostname (without port)
	host := u.Hostname()

	// Only allow localhost or 127.0.0.1
	return host == "localhost" || host == "127.0.0.1"
}

// originSecurityHandler validates Origin header to prevent DNS rebinding attacks.
// This implements the security requirement from the MCP specification:
// https://modelcontextprotocol.io/specification/2024-11-05/basic/transports#security-warning
		// This handles:
		// - Non-browser clients (curl, SDKs) - no Origin header sent
		// - Same-origin requests - browsers don't send Origin for same-origin
		if origin != "" && !isAllowedOrigin(origin) {
			http.Error(w, "Forbidden: Invalid Origin header", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	"io"
	"net"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	sseHandler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return g.mcpServer
	}, nil)
	mux.Handle("/sse", originSecurityHandler(sseHandler))
	httpServer := &http.Server{
		Handler: mux,
	}
	streamHandler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return g.mcpServer
	}, nil)
	mux.Handle("/mcp", originSecurityHandler(streamHandler))
	httpServer := &http.Server{
		Handler: mux,
	}
		}
	}
}

// originSecurityHandler validates Origin header to prevent DNS rebinding attacks.
// This implements the security requirement from the MCP specification:
// https://modelcontextprotocol.io/specification/2024-11-05/basic/transports#security-warning
func originSecurityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow requests with no Origin header
		// This handles:
		// - Non-browser clients (curl, SDKs) - no Origin header sent
		// - Same-origin requests - browsers don't send Origin for same-origin
		if origin != "" {
			// For cross-origin requests (browser-based), only allow localhost origins
			// This prevents DNS rebinding attacks using 0.0.0.0 or malicious domains
			allowed := origin == "http://localhost" ||
				origin == "https://localhost" ||
				origin == "http://127.0.0.1" ||
				origin == "https://127.0.0.1" ||
				strings.HasPrefix(origin, "http://localhost:") ||
				strings.HasPrefix(origin, "https://localhost:") ||
				strings.HasPrefix(origin, "http://127.0.0.1:") ||
				strings.HasPrefix(origin, "https://127.0.0.1:")

			if !allowed {
				http.Error(w, "Forbidden: Invalid Origin header", http.StatusForbidden)
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

import (
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/modelcontextprotocol/go-sdk/mcp"

}

// isAllowedOrigin validates that the origin is from localhost.
// Returns true if the origin's hostname is "localhost", "127.0.0.1", or "::1" (IPv6 localhost).
func isAllowedOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
	// Extract hostname (without port)
	host := u.Hostname()

	// Only allow localhost, IPv4 loopback, or IPv6 loopback
	return host == "localhost" || host == "127.0.0.1" || host == "::1"
}

// originSecurityHandler validates Origin header to prevent DNS rebinding attacks.
// This implements the security requirement from the MCP specification:
// https://modelcontextprotocol.io/specification/2024-11-05/basic/transports#security-warning
//
// Note: Origin validation is NOT skipped in container mode because:
// - Container services use HTTP clients which don't send Origin headers (validation allows them)
// - If container port is exposed to host, Origin validation still protects against browsers
func originSecurityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")

		// Allow requests with no Origin header
		// This handles:
		// - Non-browser clients (curl, SDKs) - no Origin header sent
		// - Same-origin requests - browsers don't send Origin for same-origin
		// - Container-to-container requests (HTTP clients don't send Origin)
		if origin != "" && !isAllowedOrigin(origin) {
			msg := fmt.Sprintf("Forbidden: Origin must be localhost, 127.0.0.1, or ::1, got: %s", origin)
			http.Error(w, msg, http.StatusForbidden)
			return
		}

