package main

	"io"
	"net"
	"net/http"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	}
}

// originSecurityHandler validates Origin header to prevent DNS rebinding attacks.
// This implements the security requirement from the MCP specification:
// https://modelcontextprotocol.io/specification/2024-11-05/basic/transports#security-warning
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
	"io"
	"net"
	"net/http"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	sseHandler := mcp.NewSSEHandler(func(_ *http.Request) *mcp.Server {
		return g.mcpServer
	}, nil)
	mux.Handle("/sse", sseHandler)
	httpServer := &http.Server{
		Handler: mux,
	}
	streamHandler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		return g.mcpServer
	}, nil)
	mux.Handle("/mcp", streamHandler)
	httpServer := &http.Server{
		Handler: mux,
	}
		}
	}
}

import (
	"context"
	"io"
	"net"
	"net/http"
	"net/url"
	"os"

	"github.com/modelcontextprotocol/go-sdk/mcp"

}

// isAllowedOrigin validates that the origin is from localhost.
// Returns true if the origin's hostname is "localhost" or "127.0.0.1" (any port allowed).
func isAllowedOrigin(origin string) bool {
	u, err := url.Parse(origin)
	if err != nil {
	// Extract hostname (without port)
	host := u.Hostname()

	// Only allow localhost or 127.0.0.1
	return host == "localhost" || host == "127.0.0.1"
}

// originSecurityHandler validates Origin header to prevent DNS rebinding attacks.
func originSecurityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Skip origin validation in container environments (compose networking)
		if os.Getenv("DOCKER_MCP_IN_CONTAINER") == "1" {
			next.ServeHTTP(w, r)
			return
		}

		origin := r.Header.Get("Origin")

		// Allow requests with no Origin header
		// This handles:
		// - Non-browser clients (curl, SDKs) - no Origin header sent
		// - Same-origin requests - browsers don't send Origin for same-origin
		if origin != "" && !isAllowedOrigin(origin) {
			http.Error(w, "Forbidden: Invalid Origin header", http.StatusForbidden)
			return
		}

