package main

	}

	// The entry is a directory
	dir, err := gitRepo.LsTree(entry.ID().String())
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get tree")
		return
	}

	entries, err := dir.Entries()
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "list entries")
		return
	}

	sha := c.Params(":sha")
	tree, err := gitRepo.LsTree(sha)
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get tree")
		return
	}

	entries, err := tree.Entries()
	if err != nil {
		c.Error(err, "list entries")
		return

	// Get page list.
	if isViewPage {
		entries, err := commit.Entries()
		if err != nil {
			c.Error(err, "list entries")
			return nil, ""
		return
	}

	entries, err := commit.Entries()
	if err != nil {
		c.Error(err, "list entries")
		return
