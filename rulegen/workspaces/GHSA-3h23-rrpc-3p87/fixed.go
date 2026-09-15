package main

	if m.serveGitignore(w, r) {
		return nil
	}

	clientIP, err := clientIPFromRequest(r)
	if err != nil {
		m.log.Error("Invalid client IP", zap.String("remote_addr", r.RemoteAddr), zap.Error(err))
		return caddyhttp.Error(http.StatusForbidden, err)
	}
	m.log.Debug("Ranges", zap.Strings("ranges", m.Ranges))

	// Request should be blocked
	return m.responder.ServeHTTP(w, r, next)
}

func clientIPFromRequest(r *http.Request) (net.IP, error) {
	if clientIP, ok := caddyhttp.GetVar(r.Context(), caddyhttp.ClientIPVarKey).(string); ok && clientIP != "" {
		return parseClientIP(clientIP)
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return nil, fmt.Errorf("invalid client IP format")
	}
	return parseClientIP(host)
}

func parseClientIP(rawIP string) (net.IP, error) {
	clientIP := net.ParseIP(rawIP)
	if clientIP == nil {
		return nil, fmt.Errorf("invalid client IP")
	}
	return clientIP, nil
}
