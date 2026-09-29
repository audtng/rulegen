package rules
import (
    "context"
    "net/http"
)
func vulnVar(r *http.Request) {
    // ruleid: test-var
    user := r.Header.Get("X-Forwarded-User")
    _ = user
}
func vulnCtx(r *http.Request) context.Context {
    // ruleid: test-var
    return context.WithValue(r.Context(), "user", r.Header.Get("X-Remote-User"))
}
func vulnRet(r *http.Request) string {
    // ruleid: test-var
    return r.Header.Get("X-User-Id")
}
func safeOther(r *http.Request) string {
    // ok: test-var
    return r.Header.Get("X-Request-Id")
}
