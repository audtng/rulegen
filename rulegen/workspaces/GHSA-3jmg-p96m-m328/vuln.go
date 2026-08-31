package main

	// FileLoading settings
	u.FileLoading = d.FileLoading

	// Permissions
	u.Permissions.Api = d.Account.Permissions.Api
	u.Permissions.Admin = d.Account.Permissions.Admin
		u.LoginMethod = users.LoginMethod(d.Account.LoginMethod)
	}

	if len(u.Scopes) == 0 && u.Username != "anonymous" {
		for _, source := range Config.Server.Sources {
			if source.Config.DefaultEnabled {
				u.Scopes = append(u.Scopes, users.SourceScope{
