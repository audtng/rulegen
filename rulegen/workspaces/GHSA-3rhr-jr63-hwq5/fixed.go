package main

		return model.NewAppError("DeletePost", "app.post.get.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}

	postHasFiles := len(post.FileIds) > 0

	// If the post is a burn-on-read post, we should get the original post contents
	if post.Type == model.PostTypeBurnOnRead {
		revealedPost, appErr := a.getBurnOnReadPost(rctx, post)
		if appErr != nil {
			rctx.Logger().Warn("Failed to get burn-on-read post", mlog.Err(appErr))
		}
		if revealedPost != nil {
			postHasFiles = len(revealedPost.FileIds) > 0
		}
	}

	if postHasFiles {
		appErr := a.PermanentDeleteFilesByPost(rctx, post.Id)
		if appErr != nil {
			return appErr
