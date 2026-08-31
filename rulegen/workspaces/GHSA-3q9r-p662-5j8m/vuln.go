package main

	}

	if _, ok := forwardReq.Header[forwardedheaders.XForwardedPort]; !ok {
		forwardReq.Header.Set(forwardedheaders.XForwardedPort, forwardedPort(req))
	}

	if _, ok := forwardReq.Header[forwardedheaders.XForwardedHost]; !ok {
	return filteredHeaders
}

func forwardedPort(req *http.Request) string {
	if req == nil {
		return ""
	}
		return port
	}

	if req.Header.Get(forwardedheaders.XForwardedProto) == "https" || req.Header.Get(forwardedheaders.XForwardedProto) == "wss" {
		return "443"
	}

