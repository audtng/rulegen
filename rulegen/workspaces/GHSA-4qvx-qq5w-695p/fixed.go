package main

	enableRemoteScriptChecks := b.boolVal(c.EnableScriptChecks)
	enableLocalScriptChecks := b.boolValWithDefault(c.EnableLocalScriptChecks, enableRemoteScriptChecks)

	// VerifyServerHostname implies VerifyOutgoing
	verifyServerName := b.boolVal(c.VerifyServerHostname)
	verifyOutgoing := b.boolVal(c.VerifyOutgoing)
	if verifyServerName {
		// Setting only verify_server_hostname is documented to imply
		// verify_outgoing. If it doesn't then we risk sending communication over TCP
		// when we documented it as forcing TLS for RPCs. Enforce this here rather
		// than in several different places through the code that need to reason
		// about it. (See CVE-2018-19653)
		verifyOutgoing = true
	}

	// ----------------------------------------------------------------
	// build runtime config
	//
		VerifyIncoming:                          b.boolVal(c.VerifyIncoming),
		VerifyIncomingHTTPS:                     b.boolVal(c.VerifyIncomingHTTPS),
		VerifyIncomingRPC:                       b.boolVal(c.VerifyIncomingRPC),
		VerifyOutgoing:                          verifyOutgoing,
		VerifyServerHostname:                    verifyServerName,
		Watches:                                 c.Watches,
	}

// requests. It will return a nil config if this configuration should
// not use TLS for outgoing connections.
func (c *Config) OutgoingTLSConfig() (*tls.Config, error) {
	if !c.UseTLS && !c.VerifyOutgoing {
		return nil, nil
	}
