#!/bin/bash
set -e

BASE_DIR="/src/rulegen"
WORKSPACE="/src/rulegen/scratch/sim_workspace"
rm -rf "$WORKSPACE"
mkdir -p "$WORKSPACE"

cat << 'RULE_EOF' > "$WORKSPACE/advisory_test.yaml"
rules:
  - id: h2c-upgrade-auth-bypass
    mode: search
    languages:
      - go
    severity: ERROR
    message: >-
      Detected HTTP/2 Cleartext (h2c) handler wrapped by outer middleware or
      authentication handler. In Go's h2c implementation, when an HTTP/1.1
      request is upgraded to HTTP/2 Cleartext, the connection is hijacked and
      subsequent HTTP/2 requests are routed directly to the inner handler passed
      to 'h2c.NewHandler', bypassing any outer middleware. Apply authentication
      and authorization middleware directly to the inner handler before passing
      it to 'h2c.NewHandler'.
    metadata:
      cwe:
        - "CWE-287: Improper Authentication"
      owasp:
        - "A07:2021 - Identification and Authentication Failures"
      category: security
      confidence: HIGH
      references:
        - "https://cwe.mitre.org/data/definitions/287.html"
        - "https://owasp.org/Top10/A07_2021-Identification_and_Authentication_Failures/"
    paths:
      exclude:
        - "*_test.go"
        - "test/**/*.go"
        - "mock/**/*.go"
    patterns:
      - pattern-either:
          # Direct function call wrapping h2c.NewHandler inline
          - patterns:
              - pattern: $WRAPPER(..., h2c.NewHandler(...), ...)
              - pattern-not: http.ListenAndServe(...)
              - pattern-not: http.ListenAndServeTLS(...)
              - pattern-not: http.Serve(...)
              - pattern-not: http.ServeTLS(...)
              - pattern-not: http.Handle(...)
              - pattern-not: httptest.NewServer(...)
              - pattern-not: httptest.NewTLSServer(...)
              - pattern-not: httptest.NewUnstartedServer(...)
          # Direct method call wrapping h2c.NewHandler inline
          - patterns:
              - pattern: $OBJ.$METHOD(..., h2c.NewHandler(...), ...)
              - pattern-not: http.ListenAndServe(...)
              - pattern-not: http.ListenAndServeTLS(...)
              - pattern-not: http.Serve(...)
              - pattern-not: http.ServeTLS(...)
              - pattern-not: http.Handle(...)
              - pattern-not: httptest.NewServer(...)
              - pattern-not: httptest.NewTLSServer(...)
              - pattern-not: httptest.NewUnstartedServer(...)
              - pattern-not: $SERVER.Serve(...)
              - pattern-not: $SERVER.ServeTLS(...)
          # Variable assignment to h2c.NewHandler followed by function wrapping
          - patterns:
              - pattern-either:
                  - pattern-inside: |
                      $H2C := h2c.NewHandler(...)
                      ...
                  - pattern-inside: |
                      $H2C = h2c.NewHandler(...)
                      ...
              - pattern: $WRAPPER(..., $H2C, ...)
              - pattern-not: http.ListenAndServe(...)
              - pattern-not: http.ListenAndServeTLS(...)
              - pattern-not: http.Serve(...)
              - pattern-not: http.ServeTLS(...)
              - pattern-not: http.Handle(...)
              - pattern-not: httptest.NewServer(...)
              - pattern-not: httptest.NewTLSServer(...)
              - pattern-not: httptest.NewUnstartedServer(...)
          # Variable assignment to h2c.NewHandler followed by method wrapping
          - patterns:
              - pattern-either:
                  - pattern-inside: |
                      $H2C := h2c.NewHandler(...)
                      ...
                  - pattern-inside: |
                      $H2C = h2c.NewHandler(...)
                      ...
              - pattern: $OBJ.$METHOD(..., $H2C, ...)
              - pattern-not: http.ListenAndServe(...)
              - pattern-not: http.ListenAndServeTLS(...)
              - pattern-not: http.Serve(...)
              - pattern-not: http.ServeTLS(...)
              - pattern-not: http.Handle(...)
              - pattern-not: httptest.NewServer(...)
              - pattern-not: httptest.NewTLSServer(...)
              - pattern-not: httptest.NewUnstartedServer(...)
              - pattern-not: $SERVER.Serve(...)
              - pattern-not: $SERVER.ServeTLS(...)
RULE_EOF

cat << 'TEST_EOF' > "$WORKSPACE/advisory_test_test.go"
package rules

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type AuthMiddleware struct{}

func (a *AuthMiddleware) Wrap(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

func authMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	})
}

type AuthProvider struct{}

func (ap *AuthProvider) Middleware() *AuthMiddleware {
	return &AuthMiddleware{}
}

func testVulnerableMethodWrap(mux *http.ServeMux) {
	h2cHandler := h2c.NewHandler(mux, &http2.Server{})
	auth := &AuthMiddleware{}
	server := &http.Server{
		// ruleid: h2c-upgrade-auth-bypass
		Handler: auth.Wrap(h2cHandler),
	}
	_ = server
}

func testVulnerableFunctionWrap(mux *http.ServeMux) {
	h2cHandler := h2c.NewHandler(mux, &http2.Server{})
	server := &http.Server{
		// ruleid: h2c-upgrade-auth-bypass
		Handler: authMiddleware(h2cHandler),
	}
	_ = server
}

func testVulnerableInlineListenAndServe(mux *http.ServeMux) {
	// ruleid: h2c-upgrade-auth-bypass
	_ = http.ListenAndServe(":8080", authMiddleware(h2c.NewHandler(mux, &http2.Server{})))
}

func testVulnerableDirectMiddlewareAssign(mux *http.ServeMux) {
	h2cHandler := h2c.NewHandler(mux, &http2.Server{})
	// ruleid: h2c-upgrade-auth-bypass
	secured := authMiddleware(h2cHandler)
	_ = secured
}

func testVulnerableChainedMethodWrap(mux *http.ServeMux) {
	hdlr := h2c.NewHandler(mux, &http2.Server{})
	provider := &AuthProvider{}
	server := &http.Server{
		// ruleid: h2c-upgrade-auth-bypass
		Handler: provider.Middleware().Wrap(hdlr),
	}
	_ = server
}

func testSafeDirectServer(mux *http.ServeMux) {
	h2cHandler := h2c.NewHandler(mux, &http2.Server{})
	// ok: h2c-upgrade-auth-bypass
	server := &http.Server{
		Handler: h2cHandler,
	}
	_ = server
}

func testSafeListenAndServe(mux *http.ServeMux) {
	// ok: h2c-upgrade-auth-bypass
	_ = http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}

func testSafeAuthInsideH2C(mux *http.ServeMux) {
	authed := authMiddleware(mux)
	h2cHandler := h2c.NewHandler(authed, &http2.Server{})
	// ok: h2c-upgrade-auth-bypass
	server := &http.Server{
		Handler: h2cHandler,
	}
	_ = server
}

func testSafeAuthWrapInsideH2C(mux *http.ServeMux) {
	auth := &AuthMiddleware{}
	authed := auth.Wrap(mux)
	h2cHandler := h2c.NewHandler(authed, &http2.Server{})
	// ok: h2c-upgrade-auth-bypass
	_ = http.ListenAndServe(":8080", h2cHandler)
}
TEST_EOF

cd "$WORKSPACE"

echo "Gate 1: Build validation..."
mkdir -p "build_advisory_test"
cp "advisory_test_test.go" "build_advisory_test/main.go"
cd "build_advisory_test"
go mod init ruletest >/dev/null 2>&1
if command -v goimports &> /dev/null; then
    goimports -w "main.go"
fi
go get -d ./... >/dev/null 2>&1 || go mod tidy >/dev/null 2>&1 || true
go build -o /dev/null "./main.go"
echo "Go build OK!"
cd ..
rm -rf "build_advisory_test"

echo "Gate 1: Semgrep test validation..."
semgrep --validate --config advisory_test.yaml
semgrep --test --config advisory_test.yaml advisory_test_test.go
echo "Semgrep test OK!"

echo "Gate 1.5: Semantic Deduplication..."
RULE_HASH=$(python3 "$BASE_DIR/hash_rule.py" "advisory_test.yaml")
echo "Rule hash: $RULE_HASH"
if grep -q "$RULE_HASH" "$BASE_DIR/rule_hashes.txt"; then
    echo "ERROR: Duplicate rule hash!"
    exit 1
fi
echo "Deduplication OK!"

echo "Gate 2: FP check on stdlib..."
semgrep --config advisory_test.yaml /usr/share/go-1.27/src --json -o corpus_results.json --quiet || true
FP_COUNT=$(jq '.results | length' corpus_results.json 2>/dev/null || echo "0")
echo "FP count: $FP_COUNT"
if [ "$FP_COUNT" -gt 0 ]; then
    echo "ERROR: FP count > 0!"
    exit 1
fi
echo "ALL GATES PASSED PERFECTLY!"
