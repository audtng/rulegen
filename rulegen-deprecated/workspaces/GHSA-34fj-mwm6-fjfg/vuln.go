package main


	ret.UserData = MaskedUserData
	ret.MCPOAuth = ""
	if "" != ret.AccessAuthCode {
		ret.AccessAuthCode = MaskedAccessAuthCode
	}
func HideConfSecret(c *AppConf) {
	c.AI = &conf.AI{}
	c.MCPOAuth = ""
	c.Api = &conf.API{}
	c.Flashcard = &conf.Flashcard{}
	c.ServerAddrs = []string{}
