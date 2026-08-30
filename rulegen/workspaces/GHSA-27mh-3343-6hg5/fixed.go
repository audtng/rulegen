package main

		// "data" + size (4 bytes each)
		b = b[8:]

		if len(b) < 4 {
			return fmt.Errorf("invalid encoding: expected at least %d bytes, for class, got %d", 4, len(b))
		}
		class := getInt(b[1:4])
		var ok bool
