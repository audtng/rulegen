package main

	if !c.Repo.GitRepo.HasBranch(branchName) {
		return
	}
	if len(commitID) > 0 {
		branchCommitID, err := c.Repo.GitRepo.BranchCommitID(branchName)
		if err != nil {
