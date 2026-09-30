package main

import "crypto/tls"

func testNaked(c *tls.Conn) {
    c.Handshake()
}

func testAssigned(c *tls.Conn) error {
    err := c.Handshake()
    if err != nil {
        return err
    }
    return nil
}

func testReturn(c *tls.Conn) error {
    return c.Handshake()
}

func testInlineIf(c *tls.Conn) error {
    if err := c.Handshake(); err != nil {
        return err
    }
    return nil
}
