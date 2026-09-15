package main

	"path/filepath"
	"slices"
	"strings"
	"sync"
)

// HostMatchList is used to check if a host or IP is in a list.
	ipNets []*net.IPNet
}

// MatchBuiltinExternal A valid global-unicast IP that is neither private (see MatchBuiltinPrivate)
// nor a reserved special-purpose range (see reservedIPNets); i.e. a routable host on the public internet.
const MatchBuiltinExternal = "external"

// reservedIPNets are special-purpose ranges that net.IP.IsPrivate omits but that must not be
// treated as public/external destinations (CGNAT, cloud metadata, IPv6 transition, etc.). We layer
// these on top of net.IP.IsPrivate (RFC 1918 / RFC 4193) so future additions to Go's IsPrivate are
// picked up automatically, while still covering the ranges it leaves out; otherwise the default
// allow-list would let authenticated users reach cloud metadata, internal, and IPv6 transition
// endpoints (SSRF), and a "private" block-list would fail to catch them.
var reservedIPNets = sync.OnceValue(func() []*net.IPNet {
	var nets []*net.IPNet
	for _, cidr := range []string{
		// IPv4
		"100.64.0.0/10",    // RFC 6598 Carrier-Grade NAT
		"168.63.129.16/32", // Azure WireServer metadata endpoint
		"192.0.0.0/24",     // RFC 6890 IETF protocol assignments
		"192.0.2.0/24",     // RFC 5737 TEST-NET-1
		"192.88.99.0/24",   // RFC 7526 6to4 relay anycast (deprecated)
		"198.18.0.0/15",    // RFC 2544 benchmarking
		"198.51.100.0/24",  // RFC 5737 TEST-NET-2
		"203.0.113.0/24",   // RFC 5737 TEST-NET-3
		// IPv6
		"100::/64",       // RFC 6666 discard-only
		"64:ff9b::/96",   // RFC 6052 NAT64 (can embed IPv4 such as 169.254.169.254)
		"64:ff9b:1::/48", // RFC 8215 local-use NAT64
		"2001::/32",      // RFC 4380 Teredo tunneling (embeds IPv4)
		"2001:10::/28",   // RFC 4843 ORCHID (deprecated)
		"2001:20::/28",   // RFC 7343 ORCHIDv2
		"2001:db8::/32",  // RFC 3849 documentation
		"2002::/16",      // RFC 3056 6to4 (embeds IPv4)
	} {
		_, ipNet, err := net.ParseCIDR(cidr)
		if err != nil {
			panic("hostmatcher: invalid reserved CIDR " + cidr + ": " + err.Error())
		}
		nets = append(nets, ipNet)
	}
	return nets
})

// isPrivateIP reports whether ip falls in a private (net.IP.IsPrivate) or reserved special-purpose
// range (see reservedIPNets) that must not be considered a public/external destination.
func isPrivateIP(ip net.IP) bool {
	if ip.IsPrivate() {
		return true
	}
	for _, ipNet := range reservedIPNets() {
		if ipNet.Contains(ip) {
			return true
		}
	}
	return false
}

// MatchBuiltinPrivate RFC 1918 (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) and RFC 4193 (FC00::/7),
// plus the reserved special-purpose ranges in reservedIPNets (CGNAT, NAT64, cloud metadata, etc.).
// Also called LAN/Intranet.
const MatchBuiltinPrivate = "private"

// MatchBuiltinLoopback 127.0.0.0/8 for IPv4 and ::1/128 for IPv6, localhost is included.
	for _, builtin := range hl.builtins {
		switch builtin {
		case MatchBuiltinExternal:
			if ip.IsGlobalUnicast() && !isPrivateIP(ip) {
				return true
			}
		case MatchBuiltinPrivate:
			if isPrivateIP(ip) {
				return true
			}
		case MatchBuiltinLoopback:
