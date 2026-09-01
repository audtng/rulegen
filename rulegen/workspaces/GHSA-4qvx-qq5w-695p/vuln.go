package main

	enableRemoteScriptChecks := b.boolVal(c.EnableScriptChecks)
	enableLocalScriptChecks := b.boolValWithDefault(c.EnableLocalScriptChecks, enableRemoteScriptChecks)

	// ----------------------------------------------------------------
	// build runtime config
	//
		VerifyIncoming:                          b.boolVal(c.VerifyIncoming),
		VerifyIncomingHTTPS:                     b.boolVal(c.VerifyIncomingHTTPS),
		VerifyIncomingRPC:                       b.boolVal(c.VerifyIncomingRPC),
		VerifyOutgoing:                          b.boolVal(c.VerifyOutgoing),
		VerifyServerHostname:                    b.boolVal(c.VerifyServerHostname),
		Watches:                                 c.Watches,
	}

// requests. It will return a nil config if this configuration should
// not use TLS for outgoing connections.
func (c *Config) OutgoingTLSConfig() (*tls.Config, error) {
	// If VerifyServerHostname is true, that implies VerifyOutgoing
	if c.VerifyServerHostname {
		c.VerifyOutgoing = true
	}
	if !c.UseTLS && !c.VerifyOutgoing {
		return nil, nil
	}
