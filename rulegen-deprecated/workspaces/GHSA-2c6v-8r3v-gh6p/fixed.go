package main

	if !c.Repo.GitRepo.HasBranch(branchName) {
		return
	}
	if branchName == c.Repo.Repository.DefaultBranch {
		c.Flash.Error(c.Tr("repo.branches.default_deletion_not_allowed"))
		return
	}

	protectBranch, err := database.GetProtectBranchOfRepoByName(c.Repo.Repository.ID, branchName)
	if err != nil && !database.IsErrBranchNotExist(err) {
		log.Error("Failed to get protected branch %q: %v", branchName, err)
		return
	}
	if protectBranch != nil && protectBranch.Protected {
		c.Flash.Error(c.Tr("repo.branches.protected_deletion_not_allowed"))
		return
	}

	if len(commitID) > 0 {
		branchCommitID, err := c.Repo.GitRepo.BranchCommitID(branchName)
		if err != nil {
