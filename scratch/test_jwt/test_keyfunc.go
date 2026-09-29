package rules

import "github.com/golang-jwt/jwt/v5"

func k1(s string) {
    _, _ = jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
        return []byte("secret"), nil
    })
}

func k2(s string) {
    _, _ = jwt.Parse(s, func(t *jwt.Token) (any, error) {
        return []byte("secret"), nil
    })
}

func k_safe(s string) {
    _, _ = jwt.Parse(s, func(t *jwt.Token) (interface{}, error) {
        if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
            return nil, nil
        }
        return []byte("secret"), nil
    })
}
