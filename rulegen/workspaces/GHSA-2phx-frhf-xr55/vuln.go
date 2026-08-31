package main

	channelID := postActionIntegrationRequest.Context[channelIDForContext].(string)
	rootID := postActionIntegrationRequest.Context[rootIDForContext].(string)

	slackAttachment := model.SlackAttachment{
		Text: fmt.Sprintf("You have selected `%s` to start the meeting.", action),
	}

		if err := p.storeUserPreference(userID, val); err != nil {
			p.API.LogWarn("failed to update preferences for the user", "Error", err.Error())
			return
		}

	p.API.UpdateEphemeralPost(userID, post)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(post); err != nil {
		p.API.LogError("failed to write response", "Error", err.Error())
	}
