package main

		if p2c <= 0 {
			return nil, fmt.Errorf("go-jose/go-jose: invalid P2C: must be a positive integer")
		}

		// salt is UTF8(Alg) || 0x00 || Salt Input
		alg := headers.getAlgorithm()
