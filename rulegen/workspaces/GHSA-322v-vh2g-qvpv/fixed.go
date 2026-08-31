package main

	_, err = th.SystemAdminClient.SetDefaultProfileImage(context.Background(), user.Id)
	require.NoError(t, err)

	// Check that a system admin can set the default profile image for another system admin
	anotherAdmin := th.CreateUser()
	_, appErr := th.App.UpdateUserRoles(th.Context, anotherAdmin.Id, model.SystemAdminRoleId+" "+model.SystemUserRoleId, false)
	require.Nil(t, appErr)

	_, err = th.SystemAdminClient.SetDefaultProfileImage(context.Background(), anotherAdmin.Id)
	require.NoError(t, err)

	ruser, appErr := th.App.GetUser(user.Id)
	require.Nil(t, appErr)
	assert.Less(t, ruser.LastPictureUpdate, iuser.LastPictureUpdate, "LastPictureUpdate should be updated to a lower negative number")
	if userID == "" {
		return false
	}
	if session.IsUnrestricted() || a.SessionHasPermissionTo(session, model.PermissionManageSystem) {
		return true
	}

		return true
	}

	if !a.SessionHasPermissionTo(session, model.PermissionEditOtherUsers) {
		return false
	}

	user, err := a.GetUser(userID)
	if err != nil {
		return false
	}

	if user.IsSystemAdmin() {
		return false
	}

	return true
}

func (a *App) SessionHasPermissionToUserOrBot(rctx request.CTX, session model.Session, userID string) bool {

		th.AddPermissionToRole(model.PermissionEditOtherUsers.Id, model.SystemUserManagerRoleId)
		assert.True(t, th.App.SessionHasPermissionToUser(session, th.BasicUser2.Id))
		assert.False(t, th.App.SessionHasPermissionToUser(session, th.SystemAdminUser.Id))
		th.RemovePermissionFromRole(model.PermissionEditOtherUsers.Id, model.SystemUserManagerRoleId)

		bot, err := th.App.CreateBot(th.Context, &model.Bot{
