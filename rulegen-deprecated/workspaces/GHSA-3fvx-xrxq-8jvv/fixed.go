package main


				ip := net.ParseIP(host)
				if ip == nil {
					ips, err := net.LookupIP(host) //nolint
					if err != nil {
						return nil, fmt.Errorf("DNS resolution failed for host %s: %v", host, err)
					}
					if len(ips) == 0 {
						return nil, fmt.Errorf("no IP addresses found for host: %s", host)
					}
					ip = ips[0] // Use the first resolved IP address
				}
				if isPrivateOrInternal(ip) {
					return nil, fmt.Errorf("%w", ErrPrivateIP)
