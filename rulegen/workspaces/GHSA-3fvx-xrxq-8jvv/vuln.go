package main


				ip := net.ParseIP(host)
				if ip == nil {
					return nil, fmt.Errorf("unexpected non-IP address in dial: %s", host)
				}
				if isPrivateOrInternal(ip) {
					return nil, fmt.Errorf("%w", ErrPrivateIP)
