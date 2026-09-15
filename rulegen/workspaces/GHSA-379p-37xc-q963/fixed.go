package main

}

func (me *LogoutProvider) DoCommand(c *Context, channelId string, message string) *model.CommandResponse {
	FAIL := &model.CommandResponse{ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command_logout.fail_message")}
	SUCCESS := &model.CommandResponse{GotoLocation: "/", ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command_logout.success_message")}

	// We can't actually remove the user's cookie from here so we just dump their session and let the browser figure it out
	if c.Session.Id != "" {
		RevokeSessionById(c, c.Session.Id)
		if c.Err != nil {
			return FAIL
		}
		return SUCCESS
	}
	return FAIL
}
	c.LogAudit("")
	c.RemoveSessionCookie(w, r)
	if c.Session.Id != "" {
		RevokeSessionById(c, c.Session.Id)
	}
}

