package main

import (
	"context"
	"crypto/tls"
	"fmt"
)

type handshakeContexter interface {
	HandshakeContext(context.Context) error
}

func serverExample(ctx context.Context, rwc any) {
	var err error
	if handshaker, ok := rwc.(handshakeContexter); ok {
		err = handshaker.HandshakeContext(ctx)
	}
	if err != nil {
		fmt.Println("handshake error", err)
		return
	}
}
