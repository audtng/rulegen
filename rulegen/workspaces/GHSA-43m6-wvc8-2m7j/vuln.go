package main

type WebConn struct {
	WebSocket               *websocket.Conn
	Send                    chan *model.Message
	SessionId               string
	UserId                  string
	hasPermissionsToChannel map[string]bool
	hasPermissionsToTeam    map[string]bool
}

func NewWebConn(ws *websocket.Conn, userId string, sessionId string) *WebConn {
	go func() {
		achan := Srv.Store.User().UpdateUserAndSessionActivity(userId, sessionId, model.GetMillis())
		pchan := Srv.Store.User().UpdateLastPingAt(userId, model.GetMillis())

		if result := <-achan; result.Err != nil {
			l4g.Error(utils.T("api.web_conn.new_web_conn.last_activity.error"), userId, sessionId, result.Err)
		}

		if result := <-pchan; result.Err != nil {
		Send:                    make(chan *model.Message, 64),
		WebSocket:               ws,
		UserId:                  userId,
		SessionId:               sessionId,
		hasPermissionsToChannel: make(map[string]bool),
		hasPermissionsToTeam:    make(map[string]bool),
	}
func (c *WebConn) HasPermissionsToTeam(teamId string) bool {
	perm, ok := c.hasPermissionsToTeam[teamId]
	if !ok {
		session := GetSession(c.SessionId)
		if session == nil {
			perm = false
			c.hasPermissionsToTeam[teamId] = perm
		} else {
			session = sessionResult.Data.(*model.Session)

			if session.IsExpired() {
				return nil
			} else {
				AddSessionToCache(session)
		return
	}

	wc := NewWebConn(ws, c.Session.UserId, c.Session.Id)
	hub.Register(wc)
	go wc.writePump()
	wc.readPump()
