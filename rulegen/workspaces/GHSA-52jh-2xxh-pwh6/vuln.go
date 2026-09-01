package main

	}

	acc := c.acc
	// Check if we have a loop.
	ldsPrefix := bytes.HasPrefix(sub.subject, []byte(leafNodeLoopDetectionSubjectPrefix))

	// Indicate any activity, so pub and sub or unsubs.
	c.in.subs++

	acc := c.acc
	srv := c.srv

	c.mu.Lock()
		return nil
	}

	spoke := c.isSpokeLeafNode()
	// We store local subs by account and subject and optionally queue name.
	// LS- will have the arg exactly as the key.
