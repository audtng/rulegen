package main

}

func (me *LogoutProvider) DoCommand(c *Context, channelId string, message string) *model.CommandResponse {

	return &model.CommandResponse{GotoLocation: c.GetTeamURL() + "/logout", ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command_logout.success_message")}
}
	c.LogAudit("")
	c.RemoveSessionCookie(w, r)
	if c.Session.Id != "" {
		if result := <-Srv.Store.Session().Remove(c.Session.Id); result.Err != nil {
			c.Err = result.Err
			return
		}
	}
}

