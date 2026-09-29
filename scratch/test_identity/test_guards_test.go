package rules
import (
    "context"
    "net/http"
)
func isTrusted(a string) bool { return true }

func vuln1(r *http.Request) string {
    // ruleid: test-guards
    return r.Header.Get("X-Remote-User")
}

func safe1(r *http.Request) string {
    // ok: test-guards
    if !isTrusted(r.RemoteAddr) {
        return ""
    }
    return r.Header.Get("X-Remote-User")
}

func safe2(r *http.Request) string {
    // ok: test-guards
    if r.TLS == nil {
        return ""
    }
    return r.Header.Get("X-Remote-User")
}

func safe3(r *http.Request) string {
    // ok: test-guards
    if isTrusted(r.RemoteAddr) {
        return r.Header.Get("X-Remote-User")
    }
    return ""
}
