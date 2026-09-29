package rules

import (
	"crypto/tls"
	"crypto/x509"
)

func testVulnRequestClientCert() *tls.Config {
	// ruleid: tls-insecure-client-auth
	return &tls.Config{
		ClientAuth: tls.RequestClientCert,
		MinVersion: tls.VersionTLS12,
	}
}

func testVulnRequireAnyClientCert() *tls.Config {
	// ruleid: tls-insecure-client-auth
	cfg := tls.Config{
		ClientAuth: tls.RequireAnyClientCert,
	}
	return &cfg
}

func testVulnInsecureSkipVerify() *tls.Config {
	// ruleid: tls-insecure-client-auth
	return &tls.Config{
		InsecureSkipVerify: true,
	}
}

func testVulnAssignClientAuth(cfg *tls.Config) {
	// ruleid: tls-insecure-client-auth
	cfg.ClientAuth = tls.RequestClientCert
}

func testVulnAssignInsecureSkipVerify(cfg *tls.Config) {
	// ruleid: tls-insecure-client-auth
	cfg.InsecureSkipVerify = true
}

func testSafeRequireAndVerify() *tls.Config {
	caPool := x509.NewCertPool()
	// ok: tls-insecure-client-auth
	return &tls.Config{
		ClientAuth: tls.RequireAndVerifyClientCert,
		ClientCAs:  caPool,
		MinVersion: tls.VersionTLS13,
	}
}

func testSafeVerifyPeerCertificate() *tls.Config {
	// ok: tls-insecure-client-auth
	return &tls.Config{
		ClientAuth: tls.RequireAnyClientCert,
		VerifyPeerCertificate: func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
			return nil
		},
	}
}

func testSafeVerifyConnection() *tls.Config {
	// ok: tls-insecure-client-auth
	return &tls.Config{
		InsecureSkipVerify: true,
		VerifyConnection: func(cs tls.ConnectionState) error {
			return nil
		},
	}
}

func testSafeAssignWithVerifyPeerCert(cfg *tls.Config) {
	// ok: tls-insecure-client-auth
	cfg.InsecureSkipVerify = true
	cfg.VerifyPeerCertificate = func(rawCerts [][]byte, verifiedChains [][]*x509.Certificate) error {
		return nil
	}
}

func testSafeAssignWithVerifyConnection(cfg *tls.Config) {
	// ok: tls-insecure-client-auth
	cfg.ClientAuth = tls.RequireAnyClientCert
	cfg.VerifyConnection = func(cs tls.ConnectionState) error {
		return nil
	}
}

func testSafeDefaultConfig() *tls.Config {
	// ok: tls-insecure-client-auth
	return &tls.Config{
		MinVersion: tls.VersionTLS13,
	}
}
