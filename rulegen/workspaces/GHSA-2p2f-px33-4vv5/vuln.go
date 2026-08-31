package main

	if err != nil {
		return nil, fmt.Errorf("resolve CA: %w", err)
	}

	// Get network.
	network, err := s.GetNetwork(ctx, host.NetworkID)
