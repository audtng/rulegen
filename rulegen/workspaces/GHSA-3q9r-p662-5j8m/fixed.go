package main

	}

	if _, ok := forwardReq.Header[forwardedheaders.XForwardedPort]; !ok {
		forwardReq.Header.Set(forwardedheaders.XForwardedPort, forwardedPort(req, forwardReq))
	}

	if _, ok := forwardReq.Header[forwardedheaders.XForwardedHost]; !ok {
	return filteredHeaders
}

func forwardedPort(req, forwardReq *http.Request) string {
	if req == nil {
		return ""
	}
		return port
	}

	if forwardReq.Header.Get(forwardedheaders.XForwardedProto) == "https" || forwardReq.Header.Get(forwardedheaders.XForwardedProto) == "wss" {
		return "443"
	}

