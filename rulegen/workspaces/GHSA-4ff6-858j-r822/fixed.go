package main

	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/matrix-org/gomatrix"
}

type clientOptions struct {
	transport     http.RoundTripper
	dnsCache      *DNSCache
	timeout       time.Duration
	skipVerify    bool
	keepAlives    bool
	wellKnownSRV  bool
	userAgent     string
	allowNetworks []string
	denyNetworks  []string
}

// ClientOption are supplied to NewClient or NewFederationClient.
			clientOpts.dnsCache,
			clientOpts.keepAlives,
			clientOpts.wellKnownSRV,
			clientOpts.allowNetworks,
			clientOpts.denyNetworks,
		)
	}
	client := &Client{
	}
}

// WithAllowDenyNetworks sets the allowed and denied networks for the http client. By default,
// all networks are allowed. The deny list is checked before the allow list.
func WithAllowDenyNetworks(allowCIDRs []string, denyCIDRs []string) ClientOption {
	return func(options *clientOptions) {
		options.allowNetworks = allowCIDRs
		options.denyNetworks = denyCIDRs
	}
}

const destinationTripperLifetime = time.Minute * 5 // how long to keep an entry
const destinationTripperReapInterval = time.Minute // how often to check for dead entries

	dnsCache        *DNSCache
	keepAlives      bool
	wellKnownSRV    bool
	dialer          *net.Dialer
}

func newDestinationTripper(skipVerify bool, dnsCache *DNSCache, keepAlives, wellKnownSRV bool, allowCIDRs []string, denyCIDRs []string) *destinationTripper {
	tripper := &destinationTripper{
		transports:   make(map[string]*destinationTripperTransport),
		skipVerify:   skipVerify,
		dnsCache:     dnsCache,
		keepAlives:   keepAlives,
		wellKnownSRV: wellKnownSRV,
		dialer:       newDestinationTripperDialer(allowCIDRs, denyCIDRs),
	}
	time.AfterFunc(destinationTripperReapInterval, tripper.reaper)
	return tripper
	time.AfterFunc(destinationTripperReapInterval, f.reaper)
}

// newDestinationTripperDialer creates a dialer which enforces dial timeouts on the federation requests. If
// the TCP connection doesn't complete within 5 seconds, it's probably just not
// going to.
// The dialer can also be limited to CIDR ranges, if allow or deny networks is non-empty.
func newDestinationTripperDialer(allowNetworks []string, denyNetworks []string) *net.Dialer {
	if len(allowNetworks) == 0 && len(denyNetworks) == 0 {
		return &net.Dialer{
			Timeout: time.Second * 5,
		}
	}

	return &net.Dialer{
		Timeout:        time.Second * 5,
		ControlContext: allowDenyNetworksControl(allowNetworks, denyNetworks),
	}
}

// allowDenyNetworksControl is used to allow/deny access to certain networks
func allowDenyNetworksControl(allowNetworks, denyNetworks []string) func(_ context.Context, network string, address string, conn syscall.RawConn) error {
	return func(_ context.Context, network string, address string, conn syscall.RawConn) error {
		if network != "tcp4" && network != "tcp6" {
			return fmt.Errorf("%s is not a safe network type", network)
		}

		host, _, err := net.SplitHostPort(address)
		if err != nil {
			return fmt.Errorf("%s is not a valid host/port pair: %s", address, err)
		}

		ipaddress := net.ParseIP(host)
		if ipaddress == nil {
			return fmt.Errorf("%s is not a valid IP address", host)
		}

		if !isAllowed(ipaddress, allowNetworks, denyNetworks) {
			return fmt.Errorf("%s is denied", address)
		}

		return nil // allow connection
	}
}

func isAllowed(ip net.IP, allowCIDRs []string, denyCIDRs []string) bool {
	if inRange(ip, denyCIDRs) {
		return false
	}
	if inRange(ip, allowCIDRs) {
		return true
	}
	return false // "should never happen"
}

func inRange(ip net.IP, CIDRs []string) bool {
	for i := 0; i < len(CIDRs); i++ {
		cidr := CIDRs[i]
		_, network, err := net.ParseCIDR(cidr)
		if err != nil {
			return false
		}
		if network.Contains(ip) {
			return true
		}
	}

	return false
}

type destinationTripperTransport struct {
// We need to use one transport per TLS server name (instead of giving our round
// tripper a single transport) because there is no way to specify the TLS
// ServerName on a per-connection basis.
func (f *destinationTripper) getTransport(tlsServerName string, dialer *net.Dialer) http.RoundTripper {
	f.transportsMutex.Lock()
	defer f.transportsMutex.Unlock()

					InsecureSkipVerify: f.skipVerify,
					ClientSessionCache: tls.NewLRUClientSessionCache(0), // 0 = use default
				},
				Dial:              dialer.Dial, // nolint: staticcheck
				DialContext:       dialer.DialContext,
				Proxy:             http.ProxyFromEnvironment,
				ForceAttemptHTTP2: true, // if we can multiplex requests over HTTP/2, we should
			},
		u := makeHTTPSURL(r.URL, result.Destination)
		r.URL = &u
		r.Host = string(result.Host)
		resp, err = f.getTransport(result.TLSServerName, f.dialer).RoundTrip(r)
		if err == nil {
			return resp, nil
		}
	size     int
	duration time.Duration
	entries  map[string]*dnsCacheEntry
	dialer   net.Dialer
}

func NewDNSCache(size int, duration time.Duration, allowNetworks, denyNetworks []string) *DNSCache {
	return &DNSCache{
		resolver: net.DefaultResolver,
		size:     size,
		duration: duration,
		entries:  make(map[string]*dnsCacheEntry),
		dialer: net.Dialer{
			ControlContext: allowDenyNetworksControl(allowNetworks, denyNetworks),
		},
	}
}

	// retried set to true. This stops us from recursing more than
	// once.
	retried := false

retryLookup:
	// Consult the cache for the hostname. This will cause the OS to
	// Try each address in the cached entry. If we successfully connect
	// to one of those addresses then return the conn and stop there.
	for _, addr := range entry.addrs {
		conn, err := c.dialer.DialContext(ctx, "tcp", addr.String()+":"+port)
		if err != nil {
			continue
		}
