package main

	protocol, address, dialer := makeHTTPDialer(remoteAddr)
	return protocol + "://" + address, &http.Client{
		Transport: &http.Transport{
			Dial: dialer,
		},
	}
}
