package rules
import "net/http"
type Session struct {
    User string
    UserAgent string
}
func vulnAssign(r *http.Request, s *Session) {
    // ruleid: test-assign
    s.User = r.Header.Get("X-Remote-User")
}
func safeAssign(r *http.Request, s *Session) {
    // ok: test-assign
    s.UserAgent = r.Header.Get("User-Agent")
}
