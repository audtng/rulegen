package rules

import (
    "fmt"
    "github.com/golang-jwt/jwt/v5"
)

func vuln1(s string) {
    t, _ := jwt.Parse(s, nil)
    _ = t
}

func vuln2(s string) {
    t, err := jwt.Parse(s, nil)
    fmt.Println(err)
    _ = t
}

func safe1(s string) {
    t, err := jwt.Parse(s, nil)
    if err != nil {
        return
    }
    _ = t
}

func safe2(s string) {
    if _, err := jwt.Parse(s, nil); err != nil {
        return
    }
}
