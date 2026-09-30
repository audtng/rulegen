package main

import "crypto/tls"

func test(c *tls.Conn) {
    c.Handshake()
}
