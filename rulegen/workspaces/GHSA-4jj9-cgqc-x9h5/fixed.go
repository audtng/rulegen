package main

	}
}
func createDefaultServiceMeshMonitor() {
	modified := false
	acc := access.NewReaderAccessControl()
	cfg, rev := clusHelper.GetSystemConfigRev(acc)
	if !cfg.TapProxymesh {
		cfg.TapProxymesh = true
		modified = true
	}

	// enable TLS verification for communications to LDAP/OIDC/webhook/registry servers by default for new deployment
	if !cfg.EnableTLSVerification {
		cfg.EnableTLSVerification = true
		modified = true
	}
	if modified {
		_ = clusHelper.PutSystemConfigRev(cfg, rev)
	}
}
