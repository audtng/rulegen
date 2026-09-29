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
