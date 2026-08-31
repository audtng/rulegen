package main

	return false
}

var localHosts = [...]string{"127.0.0.1", "::1"}

// IsLocalHost will return true if address is a localhost address.
func (*Ctx) isLocalHost(address string) bool {
	for _, h := range localHosts {
		if address == h {
			return true
		}
	}

// IsFromLocal will return true if request came from local.
func (c *Ctx) IsFromLocal() bool {
	return c.isLocalHost(c.fasthttp.RemoteIP().String())
}
