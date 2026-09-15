package main

	"strings"
)

func dropIPv6zone(address string) string {
	i := strings.IndexByte(address, '%')
	if i != -1 {
		address = address[:i]
	}
	return address
}

// FindClientIP returns client real IP address.
func FindClientIP(r *http.Request) string {
	headers := []string{"X-Forwarded-For", "X-Real-Ip"}
	for _, header := range headers {
	}

	// Fallback to TCP/IP source IP address.
	remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		remoteIP = r.RemoteAddr

	return remoteIP
}

				// Returns a 404 if the client is not authorized to access the metrics endpoint.
				if route.GetName() == "metrics" && !isAllowedToAccessMetricsEndpoint(r) {
					logger.Error(`[Metrics] Client not allowed: %s`, request.ClientIP(r))
					http.NotFound(w, r)
					return
				}
}

func isAllowedToAccessMetricsEndpoint(r *http.Request) bool {
	clientIP := request.ClientIP(r)

	if config.Opts.MetricsUsername() != "" && config.Opts.MetricsPassword() != "" {
		username, password, authOK := r.BasicAuth()
		if !authOK {
			logger.Info("[Metrics] [ClientIP=%s] No authentication header sent", clientIP)
			logger.Fatal(`[Metrics] Unable to parse CIDR %v`, err)
		}

		if network.Contains(net.ParseIP(clientIP)) {
			return true
		}
	}
