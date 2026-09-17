package main

	})
}

func (api *ApiManagerCtx) UpdateProfile(w http.ResponseWriter, r *http.Request) error {
	session, _ := auth.GetSession(r)

	data := session.Profile()
	if err := utils.HttpJsonRequest(w, r, &data); err != nil {
		return err
	}

	err := api.sessions.Update(session.ID(), data)
	if err != nil {
		if errors.Is(err, types.ErrSessionNotFound) {
			return utils.HttpBadRequest("session does not exist")
