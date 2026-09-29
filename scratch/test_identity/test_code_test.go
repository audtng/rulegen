package rules
import (
    "net/http"
)
type Session struct { User string }

func vulnTLS(w http.ResponseWriter, r *http.Request, s *Session) {
    // ruleid: test-rule
    user := r.Header.Get("X-Remote-User")
    s.User = user
}

func safeTLSGuard(w http.ResponseWriter, r *http.Request, s *Session) {
    // ok: test-rule
    if r.TLS == nil {
        return
    }
    user := r.Header.Get("X-Remote-User")
    s.User = user
}

func safeTLSBlock(w http.ResponseWriter, r *http.Request, s *Session) {
    // ok: test-rule
    if r.TLS != nil {
        user := r.Header.Get("X-Remote-User")
        s.User = user
    }
}
