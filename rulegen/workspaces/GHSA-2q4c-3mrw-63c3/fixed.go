package main

// Package insecureserverbind validates listen addresses for insecure, unauthenticated Kopia servers.
package insecureserverbind

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
)

// AllowDangerousUnauthenticatedNetworkFlag is the CLI flag that disables bind restrictions.
const AllowDangerousUnauthenticatedNetworkFlag = "allow-extremely-dangerous-unauthenticated-server-on-the-network"

// AllowDangerousUnauthenticatedNetworkFlagHelp is the kingpin description for that flag.
const AllowDangerousUnauthenticatedNetworkFlagHelp = "Allow unauthenticated server to listen on non-loopback addresses; " +
	"exposes full repository and control API to the network without authentication which allows any external attacker to take full control of the server host (extremely dangerous)"

// ErrDisallowedPublicBind is returned when the address would expose an unauthenticated server beyond loopback.
var ErrDisallowedPublicBind = errors.New("refusing to expose unauthenticated server on non-loopback network bind")

// RestrictionApplies reports whether insecure unauthenticated bind checks must run.
func RestrictionApplies(insecure, withoutPassword, allowDangerousNetwork bool) bool {
	return insecure && withoutPassword && !allowDangerousNetwork
}

// ValidateListenAddressIfRestricted runs [ValidateListenAddressFlag] only when [RestrictionApplies] is true.
func ValidateListenAddressIfRestricted(insecure, withoutPassword, allowDangerousNetwork bool, address string) error {
	if !RestrictionApplies(insecure, withoutPassword, allowDangerousNetwork) {
		return nil
	}

	return ValidateListenAddressFlag(address)
}

// ValidateListenerAddrIfRestricted runs [ValidateListenerAddr] only when [RestrictionApplies] is true.
func ValidateListenerAddrIfRestricted(insecure, withoutPassword, allowDangerousNetwork bool, addr net.Addr) error {
	if !RestrictionApplies(insecure, withoutPassword, allowDangerousNetwork) {
		return nil
	}

	return ValidateListenerAddr(addr)
}

func stripProtocol(addr string) string {
	return strings.TrimPrefix(strings.TrimPrefix(addr, "https://"), "http://")
}

// ParseListenHost extracts the host part of a server listen address flag value.
// If isUnix is true, host is empty and the address refers to a Unix domain socket.
//
// Unix detection runs after stripping a leading http:// or https:// (same as the server’s
// stripProtocol). Any form that becomes unix:… is treated as a Unix socket, including:
//   - unix:/path/to/socket
//   - http://unix:/path/to/socket
//   - https://unix:/path/to/socket
func ParseListenHost(address string) (host string, isUnix bool, err error) {
	stripped := stripProtocol(address)
	if strings.HasPrefix(stripped, "unix:") {
		return "", true, nil
	}

	s := stripped
	if !strings.Contains(s, "://") {
		s = "http://" + s
	}

	u, err := url.Parse(s)
	if err != nil {
		return "", false, fmt.Errorf("parsing listen address: %w", err)
	}

	return u.Hostname(), false, nil
}

// ValidateListenAddressFlag checks that --address is safe for an insecure server without a UI password.
func ValidateListenAddressFlag(address string) error {
	host, isUnix, err := ParseListenHost(address)
	if err != nil {
		return err
	}

	if isUnix {
		return nil
	}

	if host == "" {
		return fmt.Errorf("%w: missing host in listen address %q binds all interfaces; use loopback, a unix socket, or pass --%s (extremely dangerous)",
			ErrDisallowedPublicBind, address, AllowDangerousUnauthenticatedNetworkFlag)
	}

	if strings.EqualFold(host, "localhost") {
		return nil
	}

	if ip := net.ParseIP(host); ip != nil {
		if ip.IsLoopback() {
			return nil
		}

		return fmt.Errorf("%w: %q is not a loopback address; pass --%s only in isolated lab environments (extremely dangerous)",
			ErrDisallowedPublicBind, host, AllowDangerousUnauthenticatedNetworkFlag)
	}

	return fmt.Errorf("%w: hostname %q is not localhost; pass --%s only in isolated lab environments (extremely dangerous)",
		ErrDisallowedPublicBind, host, AllowDangerousUnauthenticatedNetworkFlag)
}

// ValidateListenerAddr checks the bound listener address after Listen (covers socket activation).
func ValidateListenerAddr(addr net.Addr) error {
	switch a := addr.(type) {
	case *net.UnixAddr:
		return nil
	case *net.TCPAddr:
		if a.IP != nil && a.IP.IsLoopback() {
			return nil
		}

		return fmt.Errorf("%w: listener %v is not loopback; pass --%s only in isolated lab environments (extremely dangerous)",
			ErrDisallowedPublicBind, addr, AllowDangerousUnauthenticatedNetworkFlag)
	default:
		if addr.Network() == "unix" {
			return nil
		}

		return fmt.Errorf("%w: cannot validate listener type %T %v; pass --%s only if you accept the risk",
			ErrDisallowedPublicBind, addr, addr, AllowDangerousUnauthenticatedNetworkFlag)
	}
}
	htpasswd "github.com/tg123/go-htpasswd"

	"github.com/kopia/kopia/internal/auth"
	"github.com/kopia/kopia/internal/insecureserverbind"
	"github.com/kopia/kopia/internal/server"
	"github.com/kopia/kopia/notification"
	"github.com/kopia/kopia/notification/sender/jsonsender"

	logServerRequests bool

	serverStartAllowDangerousUnauthenticatedNetwork bool

	disableCSRFTokenChecks bool // disable CSRF token checks - used for development/debugging only

	sf  serverFlags
	cmd.Flag("insecure", "Allow insecure configurations (do not use in production)").Hidden().BoolVar(&c.serverStartInsecure)
	cmd.Flag("max-concurrency", "Maximum number of server goroutines").Default("0").IntVar(&c.serverStartMaxConcurrency)

	cmd.Flag(insecureserverbind.AllowDangerousUnauthenticatedNetworkFlag, insecureserverbind.AllowDangerousUnauthenticatedNetworkFlagHelp).
		Hidden().
		BoolVar(&c.serverStartAllowDangerousUnauthenticatedNetwork)

	cmd.Flag("without-password", "Start the server without a password").Hidden().BoolVar(&c.serverStartWithoutPassword)
	cmd.Flag("random-password", "Generate random password and print to stderr").Hidden().BoolVar(&c.serverStartRandomPassword)
	cmd.Flag("htpasswd-file", "Path to htpasswd file that contains allowed user@hostname entries").Hidden().ExistingFileVar(&c.serverStartHtpasswdFile)
}

func (c *commandServerStart) run(ctx context.Context) (reterr error) {
	if err := insecureserverbind.ValidateListenAddressIfRestricted(
		c.serverStartInsecure,
		c.serverStartWithoutPassword,
		c.serverStartAllowDangerousUnauthenticatedNetwork,
		c.sf.serverAddress,
	); err != nil {
		return errors.Wrap(err, "listen address not allowed for insecure server without password")
	}

	opts, err := c.serverStartOptions(ctx)
	if err != nil {
		return err
	"github.com/coreos/go-systemd/v22/activation"
	"github.com/pkg/errors"

	"github.com/kopia/kopia/internal/insecureserverbind"
	"github.com/kopia/kopia/internal/tlsutil"
)

		return errors.Errorf("Too many activated sockets found.  Expected 1, got %v", len(listeners))
	}

	if err := insecureserverbind.ValidateListenerAddrIfRestricted(
		c.serverStartInsecure,
		c.serverStartWithoutPassword,
		c.serverStartAllowDangerousUnauthenticatedNetwork,
		l.Addr(),
	); err != nil {
		l.Close() //nolint:errcheck

		return errors.Wrap(err, "insecure server bind validation")
	}

	defer l.Close() //nolint:errcheck

	httpServer.Addr = l.Addr().String()
