package rules

import "github.com/golang-jwt/jwt/v5"

func dup(s string) {
    // ruleid: test-dup
    t, _ := jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })
    _ = t
}
