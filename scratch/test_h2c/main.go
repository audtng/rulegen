package rules

import (
	"net/http"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func Dummy() {
	_ = &http.Server{}
	_ = &http2.Server{}
	_ = h2c.NewHandler(http.DefaultServeMux, &http2.Server{})
}
