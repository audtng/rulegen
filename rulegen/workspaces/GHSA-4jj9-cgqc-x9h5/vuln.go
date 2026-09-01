package main

	}
}
func createDefaultServiceMeshMonitor() {
	acc := access.NewReaderAccessControl()
	cfg, rev := clusHelper.GetSystemConfigRev(acc)
	if !cfg.TapProxymesh {
		cfg.TapProxymesh = true
		_ = clusHelper.PutSystemConfigRev(cfg, rev)
	}
}
