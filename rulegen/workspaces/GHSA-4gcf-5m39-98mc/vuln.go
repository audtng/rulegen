package main

	_store := store.FromContext(c)
	forge := server.Config.Services.Forge

	tmpRepo, tmpBuild, err := forge.Hook(c, c.Request)
	if err != nil {
		if errors.Is(err, &types.ErrIgnoreEvent{}) {
		return
	}

	// skip the tmpBuild if any case-insensitive combination of the words "skip" and "ci"
	// wrapped in square brackets appear in the commit message
	skipMatch := skipRe.FindString(tmpBuild.Message)
		return
	}

	repo, err := _store.GetRepoNameFallback(tmpRepo.ForgeRemoteID, tmpRepo.FullName)
	if err != nil {
		msg := fmt.Sprintf("failure to get repo %s from store", tmpRepo.FullName)
		c.String(http.StatusNoContent, msg)
		return
	}

	oldFullName := repo.FullName
	if oldFullName != tmpRepo.FullName {
		// create a redirection
		err = _store.CreateRedirection(&model.Redirection{RepoID: repo.ID, FullName: repo.FullName})
		if err != nil {
			_ = c.AbortWithError(http.StatusInternalServerError, err)
			return
		}
	}

	repo.Update(tmpRepo)
	err = _store.UpdateRepo(repo)
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	// get the token and verify the hook is authorized
	parsed, err := token.ParseRequest(c.Request, func(_ *token.Token) (string, error) {
		return repo.Hash, nil
		return
	}

	if repo.UserID == 0 {
		msg := fmt.Sprintf("ignoring hook. repo %s has no owner.", repo.FullName)
		log.Warn().Msg(msg)
		c.String(http.StatusNoContent, msg)
		return
	}

	if tmpBuild.Event == model.EventPull && !repo.AllowPull {
		msg := "ignoring hook: pull requests are disabled for this repo in woodpecker"
		log.Debug().Str("repo", repo.FullName).Msg(msg)
		return
	}

	pl, err := pipeline.Create(c, _store, repo, tmpBuild)
	if err != nil {
		handlePipelineErr(c, err)
