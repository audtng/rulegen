package rules

import "github.com/golang-jwt/jwt/v5"

func f(s string) {
    t, _ := jwt.Parse(s, nil)
    _ = t
}
