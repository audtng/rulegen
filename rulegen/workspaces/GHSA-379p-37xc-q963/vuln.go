package main

}

func (me *LogoutProvider) DoCommand(c *Context, channelId string, message string) *model.CommandResponse {

	return &model.CommandResponse{GotoLocation: c.GetTeamURL() + "/logout", ResponseType: model.COMMAND_RESPONSE_TYPE_EPHEMERAL, Text: c.T("api.command_logout.success_message")}
}
package api

import (
	"strings"
	"testing"

	"github.com/mattermost/platform/model"
)

func TestLogoutTestCommand(t *testing.T) {
	th := Setup().InitBasic()

	rs1 := th.BasicClient.Must(th.BasicClient.Command(th.BasicChannel.Id, "/logout", false)).Data.(*model.CommandResponse)
	if !strings.HasSuffix(rs1.GotoLocation, "logout") {
		t.Fatal("failed to logout")
	}
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

