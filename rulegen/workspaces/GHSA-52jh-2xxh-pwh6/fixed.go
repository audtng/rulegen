package main

	}

	acc := c.acc
	// Guard against LS+ arriving before CONNECT has been processed, which
	// can happen when compression is enabled.
	if acc == nil {
		c.mu.Unlock()
		c.sendErr("Authorization Violation")
		c.closeConnection(ProtocolViolation)
		return nil
	}
	// Check if we have a loop.
	ldsPrefix := bytes.HasPrefix(sub.subject, []byte(leafNodeLoopDetectionSubjectPrefix))

	// Indicate any activity, so pub and sub or unsubs.
	c.in.subs++

	srv := c.srv

	c.mu.Lock()
		return nil
	}

	acc := c.acc
	// Guard against LS- arriving before CONNECT has been processed.
	if acc == nil {
		c.mu.Unlock()
		c.sendErr("Authorization Violation")
		c.closeConnection(ProtocolViolation)
		return nil
	}

	spoke := c.isSpokeLeafNode()
	// We store local subs by account and subject and optionally queue name.
	// LS- will have the arg exactly as the key.
