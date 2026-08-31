package main

	channelID := postActionIntegrationRequest.Context[channelIDForContext].(string)
	rootID := postActionIntegrationRequest.Context[rootIDForContext].(string)

	userIDFromHeader := r.Header.Get("Mattermost-User-Id")
	if userIDFromHeader != userID {
		p.API.LogWarn("User ID mismatch", "header_user_id", userIDFromHeader, "context_user_id", userID)
		http.Error(w, "user ID mismatch", http.StatusBadRequest)
		return
	}

	if action != usePersonalMeetingID && action != useAUniqueMeetingID {
		p.API.LogWarn("Invalid meeting action", "action", action)
		http.Error(w, "invalid meeting action", http.StatusBadRequest)
		return
	}

	// Attempt to get ephemeral post should return an error.
	// Validate bot ownership if not an ephemeral post.
	oldPost, appErr := p.client.Post.GetPost(rootID)
	if appErr == nil && oldPost.UserId != p.botUserID {
		p.API.LogWarn("Post not created by bot", "post_id", rootID, "user_id", oldPost.UserId)
		http.Error(w, "cannot update post created by non-bot user", http.StatusForbidden)
		return
	}

	slackAttachment := model.SlackAttachment{
		Text: fmt.Sprintf("You have selected `%s` to start the meeting.", action),
	}

		if err := p.storeUserPreference(userID, val); err != nil {
			p.API.LogWarn("failed to update preferences for the user", "Error", err.Error())
			http.Error(w, "failed to update preferences for the user", http.StatusInternalServerError)
			return
		}

	p.API.UpdateEphemeralPost(userID, post)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(post); err != nil {
		p.API.LogError("failed to write response", "Error", err.Error())
	}
