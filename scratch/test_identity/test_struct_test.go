package rules
import "net/http"
type AuthContext struct {
    Authorized bool
    User       string
    Owner      string
}
func vulnStruct(r *http.Request) *AuthContext {
    // ruleid: test-struct
    return &AuthContext{
        Authorized: true,
        User:       r.Header.Get("X-User-Id"),
    }
}
func safeStruct(r *http.Request, defaultOwner string) *AuthContext {
    // ok: test-struct
    return &AuthContext{
        Authorized: true,
        Owner:      defaultOwner,
    }
}
