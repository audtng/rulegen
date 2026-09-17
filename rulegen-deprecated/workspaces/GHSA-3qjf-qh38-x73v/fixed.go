package main

	"strings"
)

// FindClientIP returns the client real IP address based on trusted Reverse-Proxy HTTP headers.
func FindClientIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-Ip"}
	for _, header := range headers {
	}

	// Fallback to TCP/IP source IP address.
	return FindRemoteIP(r)
}

// FindRemoteIP returns remote client IP address.
func FindRemoteIP(r *http.Request) string {
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr

	return remoteIP
}

func dropIPv6zone(address string) string {
	i := strings.IndexByte(address, '%')
	if i != -1 {
		address = address[:i]
	}
	return address
}

				// Returns a 404 if the client is not authorized to access the metrics endpoint.
				if route.GetName() == "metrics" && !isAllowedToAccessMetricsEndpoint(r) {
					logger.Error(`[Metrics] [ClientIP=%s] Client not allowed (%s)`, request.ClientIP(r), r.RemoteAddr)
					http.NotFound(w, r)
					return
				}
}

func isAllowedToAccessMetricsEndpoint(r *http.Request) bool {
	if config.Opts.MetricsUsername() != "" && config.Opts.MetricsPassword() != "" {
		clientIP := request.ClientIP(r)
		username, password, authOK := r.BasicAuth()
		if !authOK {
			logger.Info("[Metrics] [ClientIP=%s] No authentication header sent", clientIP)
			logger.Fatal(`[Metrics] Unable to parse CIDR %v`, err)
		}

		// We use r.RemoteAddr in this case because HTTP headers like X-Forwarded-For can be easily spoofed.
		// The recommendation is to use HTTP Basic authentication.
		if network.Contains(net.ParseIP(request.FindRemoteIP(r))) {
			return true
		}
	}
