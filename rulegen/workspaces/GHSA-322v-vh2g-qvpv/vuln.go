package main

	_, err = th.SystemAdminClient.SetDefaultProfileImage(context.Background(), user.Id)
	require.NoError(t, err)

	ruser, appErr := th.App.GetUser(user.Id)
	require.Nil(t, appErr)
	assert.Less(t, ruser.LastPictureUpdate, iuser.LastPictureUpdate, "LastPictureUpdate should be updated to a lower negative number")
	if userID == "" {
		return false
	}
	if session.IsUnrestricted() {
		return true
	}

		return true
	}

	if a.SessionHasPermissionTo(session, model.PermissionEditOtherUsers) {
		return true
	}

	return false
}

func (a *App) SessionHasPermissionToUserOrBot(rctx request.CTX, session model.Session, userID string) bool {

		th.AddPermissionToRole(model.PermissionEditOtherUsers.Id, model.SystemUserManagerRoleId)
		assert.True(t, th.App.SessionHasPermissionToUser(session, th.BasicUser2.Id))
		th.RemovePermissionFromRole(model.PermissionEditOtherUsers.Id, model.SystemUserManagerRoleId)

		bot, err := th.App.CreateBot(th.Context, &model.Bot{
