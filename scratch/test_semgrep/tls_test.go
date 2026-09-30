package rules

import (
	"context"
	"crypto/tls"
	"fmt"
)

func testVulnDiscardBlank(conn *tls.Conn) {
	// ruleid: tls-unverified-handshake
	_ = conn.Handshake()
}

func testVulnNakedHandshake(conn *tls.Conn) {
	// ruleid: tls-unverified-handshake
	conn.Handshake()
}

func testVulnUncheckedError(ctx context.Context, conn *tls.Conn) {
	// ruleid: tls-unverified-handshake
	err := conn.HandshakeContext(ctx)
	fmt.Println("continuing despite unverified handshake", err)
}

func testVulnVerifyHostnameIgnored(conn *tls.Conn, host string) {
	// ruleid: tls-unverified-handshake
	_ = conn.VerifyHostname(host)
}

func testSafeHandshakeChecked(conn *tls.Conn) error {
	// ok: tls-unverified-handshake
	if err := conn.Handshake(); err != nil {
		return err
	}
	return nil
}

func testSafeHandshakeContextChecked(ctx context.Context, conn *tls.Conn) error {
	// ok: tls-unverified-handshake
	err := conn.HandshakeContext(ctx)
	if err != nil {
		return err
	}
	return nil
}

func testSafeReturnHandshake(conn *tls.Conn) error {
	// ok: tls-unverified-handshake
	return conn.Handshake()
}

func testSafeVerifyHostname(conn *tls.Conn, host string) error {
	// ok: tls-unverified-handshake
	if err := conn.VerifyHostname(host); err != nil {
		return err
	}
	return nil
}
