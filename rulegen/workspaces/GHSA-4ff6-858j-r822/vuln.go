package main

	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/matrix-org/gomatrix"
}

type clientOptions struct {
	transport    http.RoundTripper
	dnsCache     *DNSCache
	timeout      time.Duration
	skipVerify   bool
	keepAlives   bool
	wellKnownSRV bool
	userAgent    string
}

// ClientOption are supplied to NewClient or NewFederationClient.
			clientOpts.dnsCache,
			clientOpts.keepAlives,
			clientOpts.wellKnownSRV,
		)
	}
	client := &Client{
	}
}

const destinationTripperLifetime = time.Minute * 5 // how long to keep an entry
const destinationTripperReapInterval = time.Minute // how often to check for dead entries

	dnsCache        *DNSCache
	keepAlives      bool
	wellKnownSRV    bool
}

func newDestinationTripper(skipVerify bool, dnsCache *DNSCache, keepAlives, wellKnownSRV bool) *destinationTripper {
	tripper := &destinationTripper{
		transports:   make(map[string]*destinationTripperTransport),
		skipVerify:   skipVerify,
		dnsCache:     dnsCache,
		keepAlives:   keepAlives,
		wellKnownSRV: wellKnownSRV,
	}
	time.AfterFunc(destinationTripperReapInterval, tripper.reaper)
	return tripper
	time.AfterFunc(destinationTripperReapInterval, f.reaper)
}

// destinationTripperDialer enforces dial timeouts on the federation requests. If
// the TCP connection doesn't complete within 5 seconds, it's probably just not
// going to.
var destinationTripperDialer = &net.Dialer{
	Timeout: time.Second * 5,
}

type destinationTripperTransport struct {
// We need to use one transport per TLS server name (instead of giving our round
// tripper a single transport) because there is no way to specify the TLS
// ServerName on a per-connection basis.
func (f *destinationTripper) getTransport(tlsServerName string) http.RoundTripper {
	f.transportsMutex.Lock()
	defer f.transportsMutex.Unlock()

					InsecureSkipVerify: f.skipVerify,
					ClientSessionCache: tls.NewLRUClientSessionCache(0), // 0 = use default
				},
				Dial:              destinationTripperDialer.Dial, // nolint: staticcheck
				DialContext:       destinationTripperDialer.DialContext,
				Proxy:             http.ProxyFromEnvironment,
				ForceAttemptHTTP2: true, // if we can multiplex requests over HTTP/2, we should
			},
		u := makeHTTPSURL(r.URL, result.Destination)
		r.URL = &u
		r.Host = string(result.Host)
		resp, err = f.getTransport(result.TLSServerName).RoundTrip(r)
		if err == nil {
			return resp, nil
		}
	size     int
	duration time.Duration
	entries  map[string]*dnsCacheEntry
}

func NewDNSCache(size int, duration time.Duration) *DNSCache {
	return &DNSCache{
		resolver: net.DefaultResolver,
		size:     size,
		duration: duration,
		entries:  make(map[string]*dnsCacheEntry),
	}
}

	// retried set to true. This stops us from recursing more than
	// once.
	retried := false
	dialer := net.Dialer{}

retryLookup:
	// Consult the cache for the hostname. This will cause the OS to
	// Try each address in the cached entry. If we successfully connect
	// to one of those addresses then return the conn and stop there.
	for _, addr := range entry.addrs {
		conn, err := dialer.DialContext(ctx, "tcp", addr.String()+":"+port)
		if err != nil {
			continue
		}
