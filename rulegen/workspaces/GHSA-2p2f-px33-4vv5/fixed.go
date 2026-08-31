package main

	if err != nil {
		return nil, fmt.Errorf("resolve CA: %w", err)
	}
	// GHSA-2p2f-px33-4vv5: zeroise the decrypted plaintext CA signing key on
	// every return path. The web mobile-bundle handler passes the real
	// CAResolver straight in and has no Wipe of its own, so wiping here is the
	// shared chokepoint for both callers. Wipe is idempotent and nil-safe, so
	// the API handler's own defer caMgr.Wipe() stays a harmless second pass.
	defer caMgr.Wipe()

	// Get network.
	network, err := s.GetNetwork(ctx, host.NetworkID)
