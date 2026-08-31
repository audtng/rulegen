package main

	})
}

// TODO: Remove when legacy mode is removed as all sessions must be synced with their providers.
func (api *ApiManagerCtx) UpdateProfile(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)

	profile := session.Profile()
	if !profile.IsAdmin {
		// Name is the only updatable field in the profile for non-admins
		var payload types.MemberProfile
		if err := utils.HttpJsonRequest(w, r, &payload); err != nil {
			return err
		}
		profile.Name = payload.Name
	} else {
		if err := utils.HttpJsonRequest(w, r, &profile); err != nil {
			return err
		}
	}

	err := api.sessions.Update(session.ID(), profile)
	if err != nil {
		if errors.Is(err, types.ErrSessionNotFound) {
			return utils.HttpBadRequest("session does not exist")
