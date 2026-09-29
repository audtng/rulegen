package rules

import (
	"errors"
	"net"

	"golang.org/x/crypto/ssh"
)

func testVulnClientIgnoreHostKey() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "admin",
		// ruleid: ssh-insecure-authentication
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
}

func testVulnClientIgnoreHostKeyAssign(cfg *ssh.ClientConfig) {
	// ruleid: ssh-insecure-authentication
	cfg.HostKeyCallback = ssh.InsecureIgnoreHostKey()
}

func testVulnClientCallbackNilReturn() *ssh.ClientConfig {
	return &ssh.ClientConfig{
		User: "admin",
		// ruleid: ssh-insecure-authentication
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}
}

func testVulnServerNoClientAuth() *ssh.ServerConfig {
	// ruleid: ssh-insecure-authentication
	return &ssh.ServerConfig{
		NoClientAuth: true,
	}
}

func testVulnServerNoClientAuthAssign(cfg *ssh.ServerConfig) {
	// ruleid: ssh-insecure-authentication
	cfg.NoClientAuth = true
}

func testVulnServerPasswordAcceptAll() *ssh.ServerConfig {
	return &ssh.ServerConfig{
		// ruleid: ssh-insecure-authentication
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			return nil, nil
		},
	}
}

func testVulnServerPublicKeyAcceptAll() *ssh.ServerConfig {
	return &ssh.ServerConfig{
		// ruleid: ssh-insecure-authentication
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			return nil, nil
		},
	}
}

func testSafeClientValidCallback(trustedKey ssh.PublicKey) *ssh.ClientConfig {
	// ok: ssh-insecure-authentication
	return &ssh.ClientConfig{
		User: "admin",
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			if string(key.Marshal()) == string(trustedKey.Marshal()) {
				return nil
			}
			return errors.New("host key mismatch")
		},
	}
}

func testSafeServerConfigAuth() *ssh.ServerConfig {
	// ok: ssh-insecure-authentication
	return &ssh.ServerConfig{
		NoClientAuth: false,
		PasswordCallback: func(conn ssh.ConnMetadata, password []byte) (*ssh.Permissions, error) {
			if string(password) == "secret" {
				return &ssh.Permissions{}, nil
			}
			return nil, errors.New("invalid password")
		},
	}
}

func testSafeServerPublicKeyAuth(expectedKey ssh.PublicKey) *ssh.ServerConfig {
	// ok: ssh-insecure-authentication
	return &ssh.ServerConfig{
		PublicKeyCallback: func(conn ssh.ConnMetadata, key ssh.PublicKey) (*ssh.Permissions, error) {
			if string(key.Marshal()) == string(expectedKey.Marshal()) {
				return &ssh.Permissions{}, nil
			}
			return nil, errors.New("unauthorized public key")
		},
	}
}
