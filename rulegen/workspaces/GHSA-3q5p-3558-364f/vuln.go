package main

	return false
}

var localHosts = [...]string{"127.0.0.1", "0.0.0.0", "::1"}

// IsLocalHost will return true if address is a localhost address.
func (*Ctx) isLocalHost(address string) bool {
	for _, h := range localHosts {
		if strings.Contains(address, h) {
			return true
		}
	}

// IsFromLocal will return true if request came from local.
func (c *Ctx) IsFromLocal() bool {
	ips := c.IPs()
	if len(ips) == 0 {
		ips = append(ips, c.IP())
	}
	return c.isLocalHost(ips[0])
}
