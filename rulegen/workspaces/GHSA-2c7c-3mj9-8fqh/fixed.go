package main

		if p2c <= 0 {
			return nil, fmt.Errorf("go-jose/go-jose: invalid P2C: must be a positive integer")
		}
		if p2c > 1000000 {
			// An unauthenticated attacker can set a high P2C value. Set an upper limit to avoid
			// DoS attacks.
			return nil, fmt.Errorf("go-jose/go-jose: invalid P2C: too high")
		}

		// salt is UTF8(Alg) || 0x00 || Salt Input
		alg := headers.getAlgorithm()
