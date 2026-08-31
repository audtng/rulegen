package main

		return model.NewAppError("DeletePost", "app.post.get.app_error", nil, "", http.StatusBadRequest).Wrap(err)
	}

	// If the post is a burn-on-read post, we should get the original post contents
	if post.Type == model.PostTypeBurnOnRead {
		tmpPost, appErr := a.getBurnOnReadPost(rctx, post)
		if appErr != nil {
			rctx.Logger().Warn("Failed to get burn-on-read post", mlog.Err(appErr))
		}
		if tmpPost != nil {
			post = tmpPost
		}
	}

	if len(post.FileIds) > 0 {
		appErr := a.PermanentDeleteFilesByPost(rctx, post.Id)
		if appErr != nil {
			return appErr
