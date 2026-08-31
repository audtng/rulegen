package main

	}

	// The entry is a directory
	dir, err := gitRepo.LsTree(entry.ID().String(), git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get tree")
		return
	}

	entries, err := dir.Entries(git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "list entries")
		return
	}

	sha := c.Params(":sha")
	tree, err := gitRepo.LsTree(sha, git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.NotFoundOrError(gitutil.NewError(err), "get tree")
		return
	}

	entries, err := tree.Entries(git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.Error(err, "list entries")
		return

	// Get page list.
	if isViewPage {
		entries, err := commit.Entries(git.LsTreeOptions{Verbatim: true})
		if err != nil {
			c.Error(err, "list entries")
			return nil, ""
		return
	}

	entries, err := commit.Entries(git.LsTreeOptions{Verbatim: true})
	if err != nil {
		c.Error(err, "list entries")
		return
