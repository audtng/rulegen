package main

	"path/filepath"
	"slices"
	"strings"
)

// HostMatchList is used to check if a host or IP is in a list.
	ipNets []*net.IPNet
}

// MatchBuiltinExternal A valid non-private unicast IP, all hosts on public internet are matched
const MatchBuiltinExternal = "external"

// MatchBuiltinPrivate RFC 1918 (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) and RFC 4193 (FC00::/7). Also called LAN/Intranet.
const MatchBuiltinPrivate = "private"

// MatchBuiltinLoopback 127.0.0.0/8 for IPv4 and ::1/128 for IPv6, localhost is included.
	for _, builtin := range hl.builtins {
		switch builtin {
		case MatchBuiltinExternal:
			if ip.IsGlobalUnicast() && !ip.IsPrivate() {
				return true
			}
		case MatchBuiltinPrivate:
			if ip.IsPrivate() {
				return true
			}
		case MatchBuiltinLoopback:
