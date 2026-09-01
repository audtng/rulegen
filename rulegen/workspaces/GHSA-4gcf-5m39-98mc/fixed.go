package main

	_store := store.FromContext(c)
	forge := server.Config.Services.Forge

	//
	// 1. Parse webhook
	//

	tmpRepo, tmpBuild, err := forge.Hook(c, c.Request)
	if err != nil {
		if errors.Is(err, &types.ErrIgnoreEvent{}) {
		return
	}

	//
	// Skip if commit message contains skip-ci
	// TODO: move into global pipeline conditions logic
	//

	// skip the tmpBuild if any case-insensitive combination of the words "skip" and "ci"
	// wrapped in square brackets appear in the commit message
	skipMatch := skipRe.FindString(tmpBuild.Message)
		return
	}

	//
	// 2. Get related repo from store and take repo renaming into account
	//

	repo, err := _store.GetRepoNameFallback(tmpRepo.ForgeRemoteID, tmpRepo.FullName)
	if err != nil {
		msg := fmt.Sprintf("failure to get repo %s from store", tmpRepo.FullName)
		c.String(http.StatusNoContent, msg)
		return
	}
	oldFullName := repo.FullName

	if repo.UserID == 0 {
		msg := fmt.Sprintf("ignoring hook. repo %s has no owner.", repo.FullName)
		log.Warn().Msg(msg)
		c.String(http.StatusNoContent, msg)
		return
	}

	//
	// 3. Check if the webhook is a valid and authorized one
	//

	// get the token and verify the hook is authorized
	parsed, err := token.ParseRequest(c.Request, func(_ *token.Token) (string, error) {
		return repo.Hash, nil
		return
	}

	//
	// 4. Update repo
	//

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

	//
	// 5. Check if pull requests are allowed for this repo
	//

	if tmpBuild.Event == model.EventPull && !repo.AllowPull {
		msg := "ignoring hook: pull requests are disabled for this repo in woodpecker"
		log.Debug().Str("repo", repo.FullName).Msg(msg)
		return
	}

	//
	// 6. Finally create a pipeline
	//

	pl, err := pipeline.Create(c, _store, repo, tmpBuild)
	if err != nil {
		handlePipelineErr(c, err)
