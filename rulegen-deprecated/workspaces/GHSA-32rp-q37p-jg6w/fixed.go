package main


func TestGroupStore(t *testing.T, ss store.Store) {
	t.Run("Create", func(t *testing.T) { testGroupStoreCreate(t, ss) })
	t.Run("CreateWithUserIds", func(t *testing.T) { testGroupCreateWithUserIds(t, ss) })

	t.Run("Get", func(t *testing.T) { testGroupStoreGet(t, ss) })
	t.Run("GetByName", func(t *testing.T) { testGroupStoreGetByName(t, ss) })
	t.Run("GetByIDs", func(t *testing.T) { testGroupStoreGetByIDs(t, ss) })
	t.Run("GetMemberUsersNotInChannel", func(t *testing.T) { testGroupGetMemberUsersNotInChannel(t, ss) })

	t.Run("UpsertMember", func(t *testing.T) { testUpsertMember(t, ss) })
	t.Run("UpsertMembers", func(t *testing.T) { testUpsertMembers(t, ss) })
	t.Run("DeleteMember", func(t *testing.T) { testGroupDeleteMember(t, ss) })
	t.Run("DeleteMembers", func(t *testing.T) { testGroupDeleteMembers(t, ss) })
	t.Run("PermanentDeleteMembersByUser", func(t *testing.T) { testGroupPermanentDeleteMembersByUser(t, ss) })

	t.Run("CreateGroupSyncable", func(t *testing.T) { testCreateGroupSyncable(t, ss) })
	t.Run("GroupMemberCount", func(t *testing.T) { groupTestGroupMemberCount(t, ss) })
	t.Run("DistinctGroupMemberCount", func(t *testing.T) { groupTestDistinctGroupMemberCount(t, ss) })
	t.Run("GroupCountWithAllowReference", func(t *testing.T) { groupTestGroupCountWithAllowReference(t, ss) })

	t.Run("GetMember", func(t *testing.T) { groupTestGetMember(t, ss) })
	t.Run("GetNonMemberUsersPage", func(t *testing.T) { groupTestGetNonMemberUsersPage(t, ss) })
}

func testGroupStoreCreate(t *testing.T, ss store.Store) {
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}

	// Happy path
		Name:        model.NewString(model.NewId()),
		DisplayName: "",
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	data, err := ss.Group().Create(g2)
	require.Nil(t, data)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	_, err = ss.Group().Create(g4)
	require.NoError(t, err)
		Name:        g4.Name,
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	data, err = ss.Group().Create(g4b)
	require.Nil(t, data)
		DisplayName: strings.Repeat("x", model.GroupDisplayNameMaxLength),
		Description: strings.Repeat("x", model.GroupDescriptionMaxLength),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	require.Nil(t, g5.IsValidForCreate())

		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSource("fake"),
		RemoteId:    model.NewString(model.NewId()),
	}
	require.Equal(t, g6.IsValidForCreate().Id, "model.group.source.app_error")

		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	require.Equal(t, g7.IsValidForCreate().Id, "model.group.name.invalid_chars.app_error")
}

func testGroupCreateWithUserIds(t *testing.T, ss store.Store) {
	// Create user 1
	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user1, nErr := ss.User().Save(u1)
	require.NoError(t, nErr)

	// Create user 2
	u2 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user2, nErr := ss.User().Save(u2)
	require.NoError(t, nErr)

	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceCustom,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}

	// Save a new group
	guids1 := &model.GroupWithUserIds{
		Group:   *g1,
		UserIds: []string{user1.Id, user2.Id},
	}

	// Happy path
	d1, err := ss.Group().CreateWithUserIds(guids1)
	require.NoError(t, err)
	require.Len(t, d1.Id, 26)
	require.Equal(t, *guids1.Name, *d1.Name)
	require.Equal(t, guids1.DisplayName, d1.DisplayName)
	require.Equal(t, guids1.Description, d1.Description)
	require.Equal(t, guids1.RemoteId, d1.RemoteId)
	require.NotZero(t, d1.CreateAt)
	require.NotZero(t, d1.UpdateAt)
	require.Zero(t, d1.DeleteAt)
	require.Equal(t, *model.NewInt64(2), int64(*d1.MemberCount))

	// Requires display name

	g2 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "",
		Source:      model.GroupSourceCustom,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}

	guids2 := &model.GroupWithUserIds{
		Group:   *g2,
		UserIds: []string{user1.Id, user2.Id},
	}
	data, err := ss.Group().CreateWithUserIds(guids2)
	require.Nil(t, data)
	require.Error(t, err)
	var appErr *model.AppError
	require.True(t, errors.As(err, &appErr))
	require.Equal(t, appErr.Id, "model.group.display_name.app_error")

	// Won't accept a duplicate name
	g4 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceCustom,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids4 := &model.GroupWithUserIds{
		Group:   *g4,
		UserIds: []string{user1.Id, user2.Id},
	}
	_, err = ss.Group().CreateWithUserIds(guids4)
	require.NoError(t, err)
	g4b := &model.Group{
		Name:        g4.Name,
		DisplayName: model.NewId(),
		Source:      model.GroupSourceCustom,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids4b := &model.GroupWithUserIds{
		Group:   *g4b,
		UserIds: []string{user1.Id},
	}
	data, err = ss.Group().CreateWithUserIds(guids4b)
	require.Nil(t, data)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unique constraint: Name")

	// Fields cannot be greater than max values
	g5 := &model.Group{
		Name:        model.NewString(strings.Repeat("x", model.GroupNameMaxLength)),
		DisplayName: strings.Repeat("x", model.GroupDisplayNameMaxLength),
		Description: strings.Repeat("x", model.GroupDescriptionMaxLength),
		Source:      model.GroupSourceCustom,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids5 := &model.GroupWithUserIds{
		Group: *g5,
	}
	require.Nil(t, guids5.IsValidForCreate())

	guids5.Name = model.NewString(*guids5.Name + "x")
	require.Equal(t, guids5.IsValidForCreate().Id, "model.group.name.invalid_length.app_error")
	guids5.Name = model.NewString(model.NewId())
	require.Nil(t, guids5.IsValidForCreate())

	guids5.DisplayName = guids5.DisplayName + "x"
	require.Equal(t, guids5.IsValidForCreate().Id, "model.group.display_name.app_error")
	guids5.DisplayName = model.NewId()
	require.Nil(t, guids5.IsValidForCreate())

	guids5.Description = guids5.Description + "x"
	require.Equal(t, guids5.IsValidForCreate().Id, "model.group.description.app_error")
	guids5.Description = model.NewId()
	require.Nil(t, guids5.IsValidForCreate())

	// Must use a valid type
	g6 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSource("fake"),
		RemoteId:    model.NewString(model.NewId()),
	}
	guids6 := &model.GroupWithUserIds{
		Group: *g6,
	}
	require.Equal(t, guids6.IsValidForCreate().Id, "model.group.source.app_error")

	//must use valid characters
	g7 := &model.Group{
		Name:        model.NewString("%^#@$$"),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceCustom,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids7 := &model.GroupWithUserIds{
		Group: *g7,
	}
	require.Equal(t, guids7.IsValidForCreate().Id, "model.group.name.invalid_chars.app_error")

	// Invalid user ids
	g8 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceCustom,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids8 := &model.GroupWithUserIds{
		Group:   *g8,
		UserIds: []string{"1234uid"},
	}
	data, err = ss.Group().CreateWithUserIds(guids8)
	require.Nil(t, data)
	require.Error(t, err)
	require.Equal(t, store.NewErrNotFound("User", "1234uid"), err)
}

func testGroupStoreGet(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	d1, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	g1Opts := model.GroupSearchOpts{
		FilterAllowReference: false,
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewString(model.NewId()),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	d1, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.Len(t, d1.Id, 26)

	// Get the group
	d2, err := ss.Group().GetByRemoteID(*d1.RemoteId, model.GroupSourceLdap)
	require.NoError(t, err)
	require.Equal(t, d1.Id, d2.Id)
	require.Equal(t, *d1.Name, *d2.Name)
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewString(model.NewId()),
		}
		groups = append(groups, g)
		_, err := ss.Group().Create(g)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	g1, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	g2, err = ss.Group().Create(g2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}

	// Create a group
	g1Update.Name = model.NewString(model.NewId())
	g1Update.DisplayName = model.NewId()
	g1Update.Description = model.NewId()
	g1Update.RemoteId = model.NewString(model.NewId())

	ud1, err := ss.Group().Update(g1Update)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: "",
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.Nil(t, data)
	require.Error(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	d2, err := ss.Group().Create(g2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), "unique constraint: Name")

	// Cannot update CreateAt
	someVal := model.GetMillis()
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}

	d1, err := ss.Group().Create(g1)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.Equal(t, beforeRestoreCount+1, afterRestoreCount)
}

func testUpsertMembers(t *testing.T, ss store.Store) {
	// Create group
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)

	// Create user
	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user, nErr := ss.User().Save(u1)
	require.NoError(t, nErr)

	// Create user
	u2 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user2, nErr := ss.User().Save(u2)
	require.NoError(t, nErr)

	// Happy path
	m, err := ss.Group().UpsertMembers(group.Id, []string{user.Id, user2.Id})
	require.NoError(t, err)
	require.Equal(t, 2, len(m))

	// Duplicate composite key (GroupId, UserId)
	// Ensure new CreateAt > previous CreateAt for the same (groupId, userId)
	// time.Sleep(1 * time.Millisecond)
	_, err = ss.Group().UpsertMembers(group.Id, []string{user.Id})
	require.NoError(t, err)

	// Invalid GroupId
	_, err = ss.Group().UpsertMembers(model.NewId(), []string{user.Id})
	require.Error(t, err)
	require.Contains(t, err.Error(), "failed to get UserGroup with")

	// Restores a deleted member
	// Ensure new CreateAt > previous CreateAt for the same (groupId, userId)
	time.Sleep(1 * time.Millisecond)
	_, err = ss.Group().UpsertMembers(group.Id, []string{user.Id, user2.Id})
	require.NoError(t, err)

	_, err = ss.Group().DeleteMembers(group.Id, []string{user.Id})
	require.NoError(t, err)

	groupMembers, err := ss.Group().GetMemberUsers(group.Id)
	require.NoError(t, err)
	beforeRestoreCount := len(groupMembers)

	_, err = ss.Group().UpsertMembers(group.Id, []string{user.Id, user2.Id})
	require.NoError(t, err)

	groupMembers, err = ss.Group().GetMemberUsers(group.Id)
	require.NoError(t, err)
	afterRestoreCount := len(groupMembers)

	require.Equal(t, beforeRestoreCount+1, afterRestoreCount)
}

func testGroupDeleteMember(t *testing.T, ss store.Store) {
	// Create group
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.True(t, errors.As(err, &nfErr))
}

func testGroupDeleteMembers(t *testing.T, ss store.Store) {
	// Create user
	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user, nErr := ss.User().Save(u1)
	require.NoError(t, nErr)
	// Create group
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	guids := &model.GroupWithUserIds{
		Group:   *g1,
		UserIds: []string{user.Id},
	}
	group, err := ss.Group().CreateWithUserIds(guids)
	require.NoError(t, err)

	// Happy path
	d2, err := ss.Group().DeleteMembers(group.Id, []string{user.Id})
	require.NoError(t, err)
	require.Equal(t, d2[0].GroupId, group.Id)
	require.Equal(t, d2[0].UserId, user.Id)
	require.NotZero(t, d2[0].DeleteAt)

	// Delete an already deleted member
	_, err = ss.Group().DeleteMembers(group.Id, []string{user.Id})
	var nfErr *store.ErrNotFound
	require.True(t, errors.As(err, &nfErr))

	// Delete with non-existent User
	_, err = ss.Group().DeleteMembers(group.Id, []string{model.NewId()})
	require.True(t, errors.As(err, &nfErr))

	// Delete non-existent Group
	_, err = ss.Group().DeleteMembers(model.NewId(), []string{user.Id})
	require.True(t, errors.As(err, &nfErr))
}

func testGroupPermanentDeleteMembersByUser(t *testing.T, ss store.Store) {
	var g *model.Group
	var groups []*model.Group
			Name:        model.NewString(model.NewId()),
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewString(model.NewId()),
		}
		group, err := ss.Group().Create(g)
		groups = append(groups, group)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group2, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "ChannelMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group2, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "Pending[Channel|Team]MemberRemovals Test Group",
		RemoteId:    model.NewString(model.NewId()),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-2"),
		DisplayName:    "group-2",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-deleted"),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-3"),
		DisplayName:    "group-3",
		RemoteId:       model.NewString(model.NewId()),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	_, err = ss.Group().UpsertMember(group1.Id, user2.Id)
	require.NoError(t, err)

	_, err = ss.Group().UpsertMember(group2.Id, user2.Id)
	require.NoError(t, err)

	_, err = ss.Group().UpsertMember(deletedGroup.Id, user1.Id)
	require.NoError(t, err)

				return len(groups) > 0
			},
		},
		{
			Name:    "Filter by group member",
			Opts:    model.GroupSearchOpts{FilterHasMember: user1.Id},
			Page:    0,
			PerPage: 100,
			Resultf: func(groups []*model.Group) bool {
				return len(groups) == 1 && groups[0].Id == group1.Id
			},
		},
		{
			Name:    "Filter by non-existent group member",
			Opts:    model.GroupSearchOpts{FilterHasMember: model.NewId()},
			Page:    0,
			PerPage: 100,
			Resultf: func(groups []*model.Group) bool {
				return len(groups) == 0
			},
		},
		{
			Name:    "Filter by non-member member",
			Opts:    model.GroupSearchOpts{FilterHasMember: user2.Id},
			Page:    0,
			PerPage: 100,
			Resultf: func(groups []*model.Group) bool {
				return len(groups) == 2
			},
		},
	}

	for _, tc := range testCases {
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			Description: model.NewId(),
			RemoteId:    model.NewString(model.NewId()),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			Description: model.NewId(),
			RemoteId:    model.NewString(model.NewId()),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(group)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewString(model.NewId()),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
}

func groupTestGroupMemberCount(t *testing.T, ss store.Store) {
	user := &model.User{
		Email:    fmt.Sprintf("test.%s@localhost", model.NewId()),
		Username: model.NewId(),
	}
	user, err := ss.User().Save(user)
	require.NoError(t, err)

	user2 := &model.User{
		Email:    fmt.Sprintf("test.%s@localhost", model.NewId()),
		Username: model.NewId(),
	}
	user2, err = ss.User().Save(user2)
	require.NoError(t, err)

	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group.Id)

	member1, err := ss.Group().UpsertMember(group.Id, user.Id)
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group.Id, member1.UserId)

	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))

	member2, err := ss.Group().UpsertMember(group.Id, user2.Id)
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group.Id, member2.UserId)

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)

	user := &model.User{
		Email:    fmt.Sprintf("test.%s@localhost", model.NewId()),
		Username: model.NewId(),
	}
	user, err = ss.User().Save(user)
	require.NoError(t, err)

	user2 := &model.User{
		Email:    fmt.Sprintf("test.%s@localhost", model.NewId()),
		Username: model.NewId(),
	}
	user2, err = ss.User().Save(user2)
	require.NoError(t, err)

	member1, err := ss.Group().UpsertMember(group1.Id, user.Id)
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group1.Id, member1.UserId)

	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))

	member2, err := ss.Group().UpsertMember(group1.Id, user2.Id)
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group1.Id, member2.UserId)

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:           model.NewString(model.NewId()),
		DisplayName:    model.NewId(),
		Source:         model.GroupSourceLdap,
		RemoteId:       model.NewString(model.NewId()),
		AllowReference: true,
	})
	require.NoError(t, err)
	require.NoError(t, err)
	require.Greater(t, countAfter, count)
}

func groupTestGetMember(t *testing.T, ss store.Store) {
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)

	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user1, nErr := ss.User().Save(u1)
	require.NoError(t, nErr)

	u2 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user2, nErr := ss.User().Save(u2)
	require.NoError(t, nErr)

	_, err = ss.Group().UpsertMember(group.Id, user1.Id)
	require.NoError(t, err)

	member, err := ss.Group().GetMember(g1.Id, u1.Id)
	require.NoError(t, err)
	require.NotNil(t, member)

	member, err = ss.Group().GetMember(g1.Id, user2.Id)
	require.Error(t, err)
	require.Nil(t, member)
}

func groupTestGetNonMemberUsersPage(t *testing.T, ss store.Store) {
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)

	u1 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	user1, nErr := ss.User().Save(u1)
	require.NoError(t, nErr)

	u2 := &model.User{
		Email:    MakeEmail(),
		Username: model.NewId(),
	}
	_, nErr = ss.User().Save(u2)
	require.NoError(t, nErr)

	users, err := ss.Group().GetNonMemberUsersPage(group.Id, 0, 1000)
	require.NoError(t, err)

	originalLen := len(users)

	_, err = ss.Group().UpsertMember(group.Id, user1.Id)
	require.NoError(t, err)

	users, err = ss.Group().GetNonMemberUsersPage(group.Id, 0, 1000)
	require.NoError(t, err)
	require.Len(t, users, originalLen-1)

	users, err = ss.Group().GetNonMemberUsersPage(model.NewId(), 0, 1000)
	require.Error(t, err)
	require.Nil(t, users)
}
	Purpose     string `json:"purpose"`
}

var allChannelMembersForUserCache = cache.NewLRU(cache.LRUOptions{
	Size: AllChannelMembersForUserCacheSize,
})
}

func newSqlChannelStore(sqlStore *SqlStore, metrics einterfaces.MetricsInterface) store.ChannelStore {
	return &SqlChannelStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}
}

func (s SqlChannelStore) upsertPublicChannelT(transaction *sqlxTxWrapper, channel *model.Channel) error {
			    PublicChannels(Id, DeleteAt, TeamId, DisplayName, Name, Header, Purpose)
			VALUES
			    (:id, :deleteat, :teamid, :displayname, :name, :header, :purpose)
		`, vals)
		if err != nil && IsUniqueConstraintError(err, []string{"PRIMARY"}) {
			_, err = transaction.NamedExec(`UPDATE PublicChannels
				SET deleteAt = :deleteat,
			    TeamId = :teamid,
			    DisplayName = :displayname,
			    Name = :name,
			    Header = :header,
			    Purpose = :purpose
			    WHERE Id=:id`, vals)
		}
	} else {
		_, err = transaction.NamedExec(`
			INSERT INTO

	if channel.Type != model.ChannelTypeDirect && channel.Type != model.ChannelTypeGroup && maxChannelsPerTeam >= 0 {
		var count int64
		if err := transaction.Get(&count, "SELECT COUNT(0) FROM Channels WHERE TeamId = ? AND DeleteAt = 0 AND (Type = ? OR Type = ?)", channel.TeamId, model.ChannelTypeOpen, model.ChannelTypePrivate); err != nil {
			return nil, errors.Wrapf(err, "save_channel_count: teamId=%s", channel.TeamId)
		} else if count >= maxChannelsPerTeam {
			return nil, store.NewErrLimitExceeded("channels_per_team", int(count), "teamId="+channel.TeamId)
func (s SqlChannelStore) GetPinnedPosts(channelId string) (*model.PostList, error) {
	pl := model.NewPostList()

	posts := []*model.Post{}
	if err := s.GetReplicaX().Select(&posts, "SELECT *, (SELECT count(Posts.Id) FROM Posts WHERE Posts.RootId = (CASE WHEN p.RootId = '' THEN p.Id ELSE p.RootId END) AND Posts.DeleteAt = 0) as ReplyCount  FROM Posts p WHERE IsPinned = true AND ChannelId = ? AND DeleteAt = 0 ORDER BY CreateAt ASC", channelId); err != nil {
		return nil, errors.Wrap(err, "failed to find Posts")
	}
	for _, post := range posts {
		pl.AddPost(post)
		pl.AddOrder(post.Id)
	}
	return pl, nil
	return nil
}

func (s SqlChannelStore) GetChannels(teamId string, userId string, opts *model.ChannelSearchOpts) (model.ChannelList, error) {
	query := s.getQueryBuilder().
		Select("ch.*").
		From("Channels ch, ChannelMembers cm").
		Where(
			sq.And{
				sq.Expr("ch.Id = cm.ChannelId"),
				sq.Eq{"cm.UserId": userId},
			},
		).
		OrderBy("ch.DisplayName")

	if teamId != "" {
		query = query.Where(sq.Or{
			sq.Eq{"ch.TeamId": teamId},
			sq.Eq{"ch.TeamId": ""},
		})
	}

	if opts.IncludeDeleted {
		if opts.LastDeleteAt != 0 {
			// We filter by non-archived, and archived >= a timestamp.
			query = query.Where(sq.Or{
				sq.Eq{"ch.DeleteAt": 0},
				sq.GtOrEq{"ch.DeleteAt": opts.LastDeleteAt},
			})
		}
		// If opts.LastDeleteAt is not set, we include everything. That means no filter is needed.
	} else {
		// Don't include archived channels.
		query = query.Where(sq.Eq{"ch.DeleteAt": 0})
	}

	if opts.LastUpdateAt > 0 {
		query = query.Where(sq.GtOrEq{"ch.UpdateAt": opts.LastUpdateAt})
	}

	channels := model.ChannelList{}
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "getchannels_tosql")
	}

	err = s.GetReplicaX().Select(&channels, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get channels with TeamId=%s and UserId=%s", teamId, userId)
	}

	if len(channels) == 0 {
		return nil, store.NewErrNotFound("Channel", "userId="+userId)
	}

	return channels, nil
}

func (s SqlChannelStore) GetChannelsWithCursor(teamId string, userId string, opts *model.ChannelSearchOpts, afterChannelID string) (model.ChannelList, error) {
	query := s.getQueryBuilder().
		Select("ch.*").
		From("Channels ch, ChannelMembers cm").
		Where(
			sq.And{
				sq.Expr("ch.Id = cm.ChannelId"),
				sq.Eq{"cm.UserId": userId},
			},
		).
		OrderBy("ch.Id")

	if opts.PerPage != nil {
		// The limit is verified at the GraphQL layer.
		query = query.Limit(uint64(*opts.PerPage))
	}

	if afterChannelID != "" {
		query = query.Where(sq.Gt{"ch.Id": afterChannelID})
	}

	if teamId != "" {
		query = query.Where(sq.Or{
			sq.Eq{"ch.TeamId": teamId},
			sq.Eq{"ch.TeamId": ""},
		})
	}

	if opts.IncludeDeleted {
		if opts.LastDeleteAt != 0 {
			// We filter by non-archived, and archived >= a timestamp.
			query = query.Where(sq.Or{
				sq.Eq{"ch.DeleteAt": 0},
				sq.GtOrEq{"ch.DeleteAt": opts.LastDeleteAt},
			})
		}
		// If opts.LastDeleteAt is not set, we include everything. That means no filter is needed.
	} else {
		// Don't include archived channels.
		query = query.Where(sq.Eq{"ch.DeleteAt": 0})
	}

	if opts.LastUpdateAt > 0 {
		query = query.Where(sq.GtOrEq{"ch.UpdateAt": opts.LastUpdateAt})
	}

	channels := model.ChannelList{}
		return nil, errors.Wrap(err, "failed to create query")
	}

	data := model.ChannelListWithTeamData{}
	err = s.GetReplicaX().Select(&data, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all channels")
	}

	return data, nil
}

func (s SqlChannelStore) GetAllChannelsCount(opts store.ChannelSearchOpts) (int64, error) {
	} else {
		selectStr = "c.*, Teams.DisplayName AS TeamDisplayName, Teams.Name AS TeamName, Teams.UpdateAt AS TeamUpdateAt"
		if opts.IncludePolicyID {
			selectStr += ", RetentionPoliciesChannels.PolicyId AS PolicyID"
		}
	}


func (s SqlChannelStore) GetTeamChannels(teamId string) (model.ChannelList, error) {
	data := model.ChannelList{}
	err := s.GetReplicaX().Select(&data, "SELECT * FROM Channels WHERE TeamId = ? And Type != ? ORDER BY DisplayName", teamId, model.ChannelTypeDirect)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Channels with teamId=%s", teamId)
	}
		SELECT * FROM Channels
		WHERE (TeamId = ? OR TeamId = '')
		AND DeleteAt != 0
		AND Type != ?
		UNION
			SELECT * FROM Channels
			WHERE (TeamId = ? OR TeamId = '')
			AND DeleteAt != 0
			AND Type = ?
			AND Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = ?)
		ORDER BY DisplayName LIMIT ? OFFSET ?
	`

	if err := s.GetReplicaX().Select(&channels, query, teamId, model.ChannelTypePrivate, teamId, model.ChannelTypePrivate, userId, limit, offset); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Channel", fmt.Sprintf("TeamId=%s,UserId=%s", teamId, userId))
		}
	query := s.getQueryBuilder().
		Select(selectStr).
		From("ChannelMembers").
		Join("GroupMembers ON GroupMembers.UserId = ChannelMembers.UserId AND GroupMembers.DeleteAt = 0")

	if includeTimezones {
		query = query.Join("Users ON Users.Id = GroupMembers.UserId")
	return nil
}

func (s SqlChannelStore) UpdateLastViewedAt(channelIds []string, userId string) (map[string]int64, error) {
	lastPostAtTimes := []struct {
		Id                string
		LastPostAt        int64
		TotalMsgCount     int64
		TotalMsgCountRoot int64
	}{}

	// We use the question placeholder format for both databases, because
	// we replace that with the dollar format later on.
	// It's needed to support the prefix CTE query. See: https://github.com/Masterminds/squirrel/issues/285.
	query := sq.StatementBuilder.PlaceholderFormat(sq.Question).
		Select("Id, LastPostAt, TotalMsgCount, TotalMsgCountRoot").
		From("Channels").
		Where(sq.Eq{"Id": channelIds})

	// TODO: use a CTE for mysql too when version 8 becomes the minimum supported version.
	if s.DriverName() == model.DatabaseDriverPostgres {
		with := query.Prefix("WITH c AS (").Suffix(") ,")
		update := sq.StatementBuilder.PlaceholderFormat(sq.Question).
			Update("ChannelMembers cm").
			Set("MentionCount", 0).
			Set("MentionCountRoot", 0).
			Set("MsgCount", sq.Expr("greatest(cm.MsgCount, c.TotalMsgCount)")).
			Set("MsgCountRoot", sq.Expr("greatest(cm.MsgCountRoot, c.TotalMsgCountRoot)")).
			Set("LastViewedAt", sq.Expr("greatest(cm.LastViewedAt, c.LastPostAt)")).
			Set("LastUpdateAt", sq.Expr("greatest(cm.LastViewedAt, c.LastPostAt)")).
			SuffixExpr(sq.Expr("FROM c WHERE cm.UserId = ? AND c.Id = cm.ChannelId", userId))
		updateWrap := update.Prefix("updated AS (").Suffix(")")
		query = with.SuffixExpr(updateWrap).Suffix("SELECT Id, LastPostAt FROM c")
	}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "UpdateLastViewedAt_CTE_Tosql")
	}

	if s.DriverName() == model.DatabaseDriverPostgres {
		sql, err = sq.Dollar.ReplacePlaceholders(sql)
		if err != nil {
			return nil, errors.Wrap(err, "UpdateLastViewedAt_ReplacePlaceholders")
		}
	}

	err = s.GetMasterX().Select(&lastPostAtTimes, sql, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ChannelMembers data with userId=%s and channelId in %v", userId, channelIds)
	}
		for _, t := range lastPostAtTimes {
			times[t.Id] = t.LastPostAt
		}
		return times, nil
	}

	var msgCountQuery, msgCountQueryRoot, lastViewedQuery = sq.Case("ChannelId"), sq.Case("ChannelId"), sq.Case("ChannelId")

	for _, t := range lastPostAtTimes {
		times[t.Id] = t.LastPostAt

		msgCountQuery = msgCountQuery.When(
			sq.Expr("?", t.Id),
			sq.Expr("GREATEST(MsgCount, ?)", t.TotalMsgCount))

		msgCountQueryRoot = msgCountQueryRoot.When(
			sq.Expr("?", t.Id),
			sq.Expr("GREATEST(MsgCountRoot, ?)", t.TotalMsgCountRoot))

		lastViewedQuery = lastViewedQuery.When(
			sq.Expr("?", t.Id),
			sq.Expr("GREATEST(LastViewedAt, ?)", t.LastPostAt))
	}

	updateQuery := s.getQueryBuilder().Update("ChannelMembers").
		Set("MentionCount", 0).
		Set("MentionCountRoot", 0).
		Set("MsgCount", msgCountQuery).
		Set("MsgCountRoot", msgCountQueryRoot).
		Set("LastViewedAt", lastViewedQuery).
		Set("LastUpdateAt", sq.Expr("LastViewedAt")).
		Where(sq.Eq{
			"UserId":    userId,
			"ChannelId": channelIds,
		})

	sql, args, err = updateQuery.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "UpdateLastViewedAt_Update_Tosql")
	}

	if _, err := s.GetMasterX().Exec(sql, args...); err != nil {
		return nil, errors.Wrapf(err, "failed to update ChannelMembers with userId=%s and channelId in %v", userId, channelIds)
	}

	return times, nil
}

// UpdateLastViewedAtPost updates a ChannelMember as if the user last read the channel at the time of the given post.
// If the provided mentionCount is -1, the given post and all posts after it are considered to be mentions. Returns
// an updated model.ChannelUnreadAt that can be returned to the client.
func (s SqlChannelStore) UpdateLastViewedAtPost(unreadPost *model.Post, userID string, mentionCount, mentionCountRoot int, setUnreadCountRoot bool) (*model.ChannelUnreadAt, error) {
	unreadDate := unreadPost.CreateAt - 1

	unread, unreadRoot, err := s.CountPostsAfter(unreadPost.ChannelId, unreadDate, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ChannelMember with channelId=%s", unreadPost.ChannelId)
	}

	return result, nil
}

func (s SqlChannelStore) IncrementMentionCount(channelId string, userId string, isRoot bool) error {
	now := model.GetMillis()
	rootInc := 0
	if isRoot {
		rootInc = 1
	if err != nil {
		return errors.Wrapf(err, "failed to Update ChannelMembers with channelId=%s and userId=%s", channelId, userId)
	}
	return nil
}

func (s SqlChannelStore) GetAll(teamId string) ([]*model.Channel, error) {
	data := []*model.Channel{}
	err := s.GetReplicaX().Select(&data, "SELECT * FROM Channels WHERE TeamId = ? AND Type != ? ORDER BY Name", teamId, model.ChannelTypeDirect)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Channels with teamId=%s", teamId)
	return value, nil
}

func (s SqlChannelStore) AnalyticsDeletedTypeCount(teamId string, channelType model.ChannelType) (int64, error) {
	query := s.getQueryBuilder().
		Select("COUNT(Id) AS Value").
		From("Channels").
	return dbMembers.ToModel(), nil
}

func (s SqlChannelStore) GetMembersForUserWithCursor(userID, afterChannel, afterUser string, limit, lastUpdateAt int) (model.ChannelMembers, error) {
	query := s.getQueryBuilder().
		Select("ChannelMembers.*",
			"TeamScheme.DefaultChannelGuestRole TeamSchemeDefaultGuestRole",
			"TeamScheme.DefaultChannelUserRole TeamSchemeDefaultUserRole",
			"TeamScheme.DefaultChannelAdminRole TeamSchemeDefaultAdminRole",
			"ChannelScheme.DefaultChannelGuestRole ChannelSchemeDefaultGuestRole",
			"ChannelScheme.DefaultChannelUserRole ChannelSchemeDefaultUserRole",
			"ChannelScheme.DefaultChannelAdminRole ChannelSchemeDefaultAdminRole").
		From("ChannelMembers").
		InnerJoin("Channels ON ChannelMembers.ChannelId = Channels.Id").
		LeftJoin("Schemes ChannelScheme ON Channels.SchemeId = ChannelScheme.Id").
		LeftJoin("Teams ON Channels.TeamId = Teams.Id").
		LeftJoin("Schemes TeamScheme ON Teams.SchemeId = TeamScheme.Id").
		Where(sq.Eq{
			"ChannelMembers.UserId": userID,
			"Channels.DeleteAt":     0,
		}).
		OrderBy("ChannelId, UserId ASC").
		// The limit is verified at the GraphQL layer.
		Limit(uint64(limit))

	if afterChannel != "" && afterUser != "" {
		query = query.Where(sq.Or{
			sq.Gt{"ChannelMembers.ChannelId": afterChannel},
			sq.And{
				sq.Eq{"ChannelMembers.ChannelId": afterChannel},
				sq.Gt{"ChannelMembers.UserId": afterUser},
			},
		})
	}

	if lastUpdateAt != 0 {
		query = query.Where(sq.GtOrEq{"ChannelMembers.LastUpdateAt": lastUpdateAt})
	}

	queryString, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "getMembersForUserWithCursor_tosql")
	}

	dbMembers := channelMemberWithSchemeRolesList{}
	err = s.GetReplicaX().Select(&dbMembers, queryString, args...)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ChannelMembers data with userId=%s", userID)
	}

	return dbMembers.ToModel(), nil
}

func (s SqlChannelStore) GetMembersForUserWithPagination(userId string, page, perPage int) (model.ChannelMembersWithTeamData, error) {
	dbMembers := channelMemberWithTeamWithSchemeRolesList{}
	offset := page * perPage
		`+deleteFilter+`
		SEARCH_CLAUSE
		AND (
			c.Type != :ChannelType
			OR (
				c.Type = :ChannelType
				AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
			)
		)
	ORDER BY c.DisplayName
	`, term, map[string]interface{}{
		"UserId":      userID,
		"ChannelType": model.ChannelTypePrivate,
	})
}

		`+deleteFilter+`
		SEARCH_CLAUSE
		AND (
			c.Type != :ChannelType
			OR (
				c.Type = :ChannelType
				AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
			)
		)
	ORDER BY c.DisplayName
	LIMIT :Limit
	`, term, map[string]interface{}{
		"TeamId":      teamID,
		"UserId":      userID,
		"Limit":       model.ChannelSearchDefaultLimit,
		"ChannelType": model.ChannelTypePrivate,
	})
}

		JOIN
			ChannelMembers AS CM ON CM.ChannelId = C.Id
		WHERE
			(C.TeamId = :TeamId OR (C.TeamId = '' AND C.Type = :ChannelType))
			AND CM.UserId = :UserId
			` + deleteFilter + `
			%v
	var channels model.ChannelList

	if likeClause, likeTerm := s.buildLIKEClause(term, "Name, DisplayName, Purpose"); likeClause == "" {
		if _, err := s.GetReplica().Select(&channels, fmt.Sprintf(queryFormat, ""), map[string]interface{}{"TeamId": teamId, "UserId": userId, "ChannelType": model.ChannelTypeGroup}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	} else {
		fulltextQuery := fmt.Sprintf(queryFormat, "AND "+fulltextClause)
		query := fmt.Sprintf("(%v) UNION (%v) LIMIT 50", likeQuery, fulltextQuery)

		if _, err := s.GetReplica().Select(&channels, query, map[string]interface{}{"TeamId": teamId, "UserId": userId, "LikeTerm": likeTerm, "FulltextTerm": fulltextTerm, "ChannelType": model.ChannelTypeGroup}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	}
					%v
				) AS OtherUsers ON OtherUsers.ChannelId = C.Id
			WHERE
			    C.Type = :ChannelType
				AND CM.UserId = :UserId
			LIMIT 50`

	var channels model.ChannelList

	if likeClause, likeTerm := s.buildLIKEClause(term, "IU.Username, IU.Nickname"); likeClause == "" {
		if _, err := s.GetReplica().Select(&channels, fmt.Sprintf(queryFormat, ""), map[string]interface{}{"UserId": userId, "ChannelType": model.ChannelTypeDirect}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	} else {
		query := fmt.Sprintf(queryFormat, "AND "+likeClause)

		if _, err := s.GetReplica().Select(&channels, query, map[string]interface{}{"UserId": userId, "LikeTerm": likeTerm, "ChannelType": model.ChannelTypeDirect}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	}
			c.TeamId = :TeamId
			SEARCH_CLAUSE
			AND c.DeleteAt != 0
			AND c.Type != :ChannelType
		ORDER BY c.DisplayName
		LIMIT 100
		`, term, map[string]interface{}{
		"TeamId":      teamId,
		"UserId":      userId,
		"ChannelType": model.ChannelTypePrivate,
	})

	privateChannels, privateErr := s.performSearch(`
			c.TeamId = :TeamId
			SEARCH_CLAUSE
			AND c.DeleteAt != 0
			AND c.Type = :ChannelType
			AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
		ORDER BY c.DisplayName
		LIMIT 100
		`, term, map[string]interface{}{
		"TeamId":      teamId,
		"UserId":      userId,
		"ChannelType": model.ChannelTypePrivate,
	})

	outputErr := publicErr
			selectStr += ", t.DisplayName AS TeamDisplayName, t.Name AS TeamName, t.UpdateAt as TeamUpdateAt"
		}
		if opts.IncludePolicyID {
			selectStr += ", RetentionPoliciesChannels.PolicyId AS PolicyID"
		}
	}

	if err != nil {
		return nil, 0, errors.Wrap(err, "channel_tosql")
	}
	channels := model.ChannelListWithTeamData{}
	if err2 := s.GetReplicaX().Select(&channels, queryString, args...); err2 != nil {
		return nil, 0, errors.Wrapf(err2, "failed to find Channels with term='%s'", term)
	}
		totalCount = int64(len(channels))
	}

	return channels, totalCount, nil
}

// TODO: rewrite in squrrel
                        JOIN
                            Users u on u.Id = cm.UserId
                        WHERE
                            c.Type = :ChannelType
                        AND
                            u.Id = :UserId
                        GROUP BY
                JOIN
                    Users u on u.Id = cm.UserId
                WHERE
                    c.Type = :ChannelType
                AND
                    u.Id = :UserId
                GROUP BY
	}

	var likeClauses []string
	args := map[string]interface{}{"UserId": userId, "ChannelType": model.ChannelTypeGroup}
	terms := strings.Split(strings.ToLower(strings.Trim(term, " ")), " ")

	for idx, term := range terms {
	return dbMembers.ToModel(), nil
}

func (s SqlChannelStore) GetMembersInfoByChannelIds(channelIDs []string) (map[string][]*model.User, error) {
	query := s.getQueryBuilder().
		Select("Channels.Id as ChannelId, Users.Id, Users.FirstName, Users.LastName, Users.Nickname, Users.Username").
		From("ChannelMembers as cm").
		Join("Channels ON cm.ChannelId = Channels.Id").
		Join("Users ON cm.UserId = Users.Id").
		Where(sq.Eq{
			"Channels.Id":       channelIDs,
			"Channels.DeleteAt": 0,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "dm_gm_names_tosql")
	}

	res := []*struct {
		model.User
		ChannelId string
	}{}

	if err := s.GetReplicaX().Select(&res, sql, args...); err != nil {
		return nil, errors.Wrap(err, "failed to find channels display name")
	}

	if len(res) == 0 {
		return nil, store.NewErrNotFound("User", fmt.Sprintf("%v", channelIDs))
	}

	userInfo := make(map[string][]*model.User)
	for _, item := range res {
		userInfo[item.ChannelId] = append(userInfo[item.ChannelId], &item.User)
	}

	return userInfo, nil
}

func (s SqlChannelStore) GetChannelsByScheme(schemeId string, offset int, limit int) (model.ChannelList, error) {
	channels := model.ChannelList{}
	err := s.GetReplicaX().Select(&channels, "SELECT * FROM Channels WHERE SchemeId = ? ORDER BY DisplayName LIMIT ? OFFSET ?", schemeId, limit, offset)
			Schemes ON Channels.SchemeId = Schemes.Id
		WHERE
			Channels.Id > ?
			AND Channels.Type IN (?, ?)
		ORDER BY
			Id
		LIMIT ?`,
		afterId, model.ChannelTypeOpen, model.ChannelTypePrivate, limit); err != nil {
		return nil, errors.Wrap(err, "failed to find Channels for export")
	}

		Where(sq.And{
			sq.Gt{"Channels.Id": afterId},
			sq.Eq{"Channels.DeleteAt": int(0)},
			sq.Eq{"Channels.Type": []model.ChannelType{model.ChannelTypeDirect, model.ChannelTypeGroup}},
		}).
		OrderBy("Channels.Id").
		Limit(uint64(limit))

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

type SqlStore interface {
	GetMasterX() SqlXExecutor
	DriverName() string
}

type SqlXExecutor interface {
	Get(dest interface{}, query string, args ...interface{}) error
	NamedExec(query string, arg interface{}) (sql.Result, error)
	Exec(query string, args ...interface{}) (sql.Result, error)
	ExecRaw(query string, args ...interface{}) (sql.Result, error)
	NamedQuery(query string, arg interface{}) (*sqlx.Rows, error)
	QueryRowX(query string, args ...interface{}) *sqlx.Row
	QueryX(query string, args ...interface{}) (*sqlx.Rows, error)
	Select(dest interface{}, query string, args ...interface{}) error
}

func cleanupChannels(t *testing.T, ss store.Store) {
	list, err := ss.Channel().GetAllChannels(0, 100000, store.ChannelSearchOpts{IncludeDeleted: true})
	require.NoError(t, err, "error cleaning all channels", err)
	t.Run("RemoveMembers", func(t *testing.T) { testChannelRemoveMembers(t, ss) })
	t.Run("ChannelDeleteMemberStore", func(t *testing.T) { testChannelDeleteMemberStore(t, ss) })
	t.Run("GetChannels", func(t *testing.T) { testChannelStoreGetChannels(t, ss) })
	t.Run("GetChannelsWithCursor", func(t *testing.T) { testChannelStoreGetChannelsWithCursor(t, ss) })
	t.Run("GetChannelsByUser", func(t *testing.T) { testChannelStoreGetChannelsByUser(t, ss) })
	t.Run("GetAllChannels", func(t *testing.T) { testChannelStoreGetAllChannels(t, ss, s) })
	t.Run("GetMoreChannels", func(t *testing.T) { testChannelStoreGetMoreChannels(t, ss) })
	t.Run("GetPublicChannelsByIdsForTeam", func(t *testing.T) { testChannelStoreGetPublicChannelsByIdsForTeam(t, ss) })
	t.Run("GetChannelCounts", func(t *testing.T) { testChannelStoreGetChannelCounts(t, ss) })
	t.Run("GetMembersForUser", func(t *testing.T) { testChannelStoreGetMembersForUser(t, ss) })
	t.Run("GetMembersForUserWithCursor", func(t *testing.T) { testChannelStoreGetMembersForUserWithCursor(t, ss) })
	t.Run("GetMembersForUserWithPagination", func(t *testing.T) { testChannelStoreGetMembersForUserWithPagination(t, ss) })
	t.Run("CountPostsAfter", func(t *testing.T) { testCountPostsAfter(t, ss) })
	t.Run("UpdateLastViewedAt", func(t *testing.T) { testChannelStoreUpdateLastViewedAt(t, ss) })
	t.Run("SearchAllChannels", func(t *testing.T) { testChannelStoreSearchAllChannels(t, ss) })
	t.Run("GetMembersByIds", func(t *testing.T) { testChannelStoreGetMembersByIds(t, ss) })
	t.Run("GetMembersByChannelIds", func(t *testing.T) { testChannelStoreGetMembersByChannelIds(t, ss) })
	t.Run("GetMembersInfoByChannelIds", func(t *testing.T) { testChannelStoreGetMembersInfoByChannelIds(t, ss) })
	t.Run("SearchGroupChannels", func(t *testing.T) { testChannelStoreSearchGroupChannels(t, ss) })
	t.Run("AnalyticsDeletedTypeCount", func(t *testing.T) { testChannelStoreAnalyticsDeletedTypeCount(t, ss) })
	t.Run("GetPinnedPosts", func(t *testing.T) { testChannelStoreGetPinnedPosts(t, ss) })
	require.Len(t, members, 1, "should have saved just 1 member")

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreCreateDirectChannel(t *testing.T, ss store.Store) {
	require.True(t, errors.As(err, &nfErr))

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsByIds(t *testing.T, ss store.Store) {
	nErr = ss.Channel().Delete(o3.Id, model.GetMillis())
	require.NoError(t, nErr, nErr)

	list, nErr := ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 1, "invalid number of channels")

	cresult := ss.Channel().PermanentDelete(o2.Id)
	require.NoError(t, cresult)

	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
	})
	if assert.Error(t, nErr) {
		var nfErr *store.ErrNotFound
		require.True(t, errors.As(nErr, &nfErr))

func testChannelStoreGetChannels(t *testing.T, ss store.Store) {
	team := model.NewId()
	o1 := &model.Channel{}
	o1.TeamId = team
	o1.DisplayName = "Channel1"
	o1.Name = NewTestId()
	o1.Type = model.ChannelTypeOpen
	var nErr error
	o1, nErr = ss.Channel().Save(o1, -1)
	require.NoError(t, nErr)

	o2 := model.Channel{}
	_, err = ss.Channel().SaveMember(&m4)
	require.NoError(t, err)

	list, nErr := ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	_, ok = ids4[o1.Id]
	require.True(t, ok, "missing channel")

	// Sleeping to guarantee that the
	// UpdateAt is different.
	// The proper way would be to set UpdateAt during channel creation itself,
	// but the *Channel.PreSave method ignores any existing CreateAt value.
	// TODO: check if using an existing CreateAt breaks anything.
	time.Sleep(time.Millisecond)

	now := model.GetMillis()
	_, nErr = ss.Channel().Update(o1)
	require.NoError(t, nErr)

	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastUpdateAt:   int(now),
	})
	require.NoError(t, nErr)
	// should return 1
	require.Len(t, list, 1)

	nErr = ss.Channel().Delete(o2.Id, 10)
	require.NoError(t, nErr)

	require.NoError(t, nErr)

	// should return 1
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 1)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")

	// Should return all
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   0,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	require.Equal(t, o3.Id, list[2].Id, "missing channel")

	// Should still return all
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   10,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	require.Equal(t, o3.Id, list[2].Id, "missing channel")

	// Should return 2
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   20,
	})
	require.NoError(t, nErr)
	require.Len(t, list, 2)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	ss.Channel().InvalidateAllChannelMembersForUser(m1.UserId)
}

func testChannelStoreGetChannelsWithCursor(t *testing.T, ss store.Store) {
	teamID := model.NewId()
	o1 := &model.Channel{}
	o1.TeamId = teamID
	o1.DisplayName = "Channel1"
	o1.Name = NewTestId()
	o1.Type = model.ChannelTypeOpen
	var nErr error
	o1, nErr = ss.Channel().Save(o1, -1)
	require.NoError(t, nErr)

	o2 := model.Channel{}
	o2.TeamId = teamID
	o2.DisplayName = "Channel2"
	o2.Name = NewTestId()
	o2.Type = model.ChannelTypeOpen
	_, nErr = ss.Channel().Save(&o2, -1)
	require.NoError(t, nErr)

	o3 := model.Channel{}
	o3.TeamId = teamID
	o3.DisplayName = "Channel3"
	o3.Name = NewTestId()
	o3.Type = model.ChannelTypeOpen
	_, nErr = ss.Channel().Save(&o3, -1)
	require.NoError(t, nErr)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err := ss.Channel().SaveMember(&m1)
	require.NoError(t, err)

	m2 := model.ChannelMember{}
	m2.ChannelId = o1.Id
	m2.UserId = model.NewId()
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err = ss.Channel().SaveMember(&m2)
	require.NoError(t, err)

	m3 := model.ChannelMember{}
	m3.ChannelId = o2.Id
	m3.UserId = m1.UserId
	m3.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err = ss.Channel().SaveMember(&m3)
	require.NoError(t, err)

	m4 := model.ChannelMember{}
	m4.ChannelId = o3.Id
	m4.UserId = m1.UserId
	m4.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err = ss.Channel().SaveMember(&m4)
	require.NoError(t, err)

	list, nErr := ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
		PerPage:        model.NewInt(2),
	}, "")
	require.NoError(t, nErr)
	require.Len(t, list, 2)
	require.Equal(t, teamID, list[0].TeamId, "incorrect teamID")
	require.Equal(t, teamID, list[1].TeamId, "incorrect teamID")

	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
		PerPage:        model.NewInt(2),
	}, list[1].Id)
	require.NoError(t, nErr)
	require.Len(t, list, 1)
	require.Equal(t, teamID, list[0].TeamId, "incorrect teamID")

	// Sleeping to guarantee that the
	// UpdateAt is different.
	// The proper way would be to set UpdateAt during channel creation itself,
	// but the *Channel.PreSave method ignores any existing CreateAt value.
	// TODO: check if using an existing CreateAt breaks anything.
	time.Sleep(time.Millisecond)

	now := model.GetMillis()
	_, nErr = ss.Channel().Update(o1)
	require.NoError(t, nErr)

	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastUpdateAt:   int(now),
	}, "")
	require.NoError(t, nErr)
	// should return 1
	require.Len(t, list, 1)

	nErr = ss.Channel().Delete(o2.Id, 10)
	require.NoError(t, nErr)

	nErr = ss.Channel().Delete(o3.Id, 20)
	require.NoError(t, nErr)

	// should return 1
	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: false,
		LastDeleteAt:   0,
	}, "")
	require.NoError(t, nErr)
	require.Len(t, list, 1)

	// Should return all
	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   0,
		PerPage:        model.NewInt(2),
	}, "")
	require.NoError(t, nErr)
	require.Len(t, list, 2)

	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   0,
		PerPage:        model.NewInt(2),
	}, list[1].Id)
	require.NoError(t, nErr)
	require.Len(t, list, 1)

	// Should still return all
	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   10,
	}, "")
	require.NoError(t, nErr)
	require.Len(t, list, 3)

	// Should return 2
	list, nErr = ss.Channel().GetChannelsWithCursor(o1.TeamId, m1.UserId, &model.ChannelSearchOpts{
		IncludeDeleted: true,
		LastDeleteAt:   20,
	}, "")
	require.NoError(t, nErr)
	require.Len(t, list, 2)
}

func testChannelStoreGetChannelsByUser(t *testing.T, ss store.Store) {
	team := model.NewId()
	team2 := model.NewId()
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	_, err = ss.Group().Create(group)
	require.NoError(t, err)
	assert.Equal(t, *list[0].PolicyID, policy.ID)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreGetMoreChannels(t *testing.T, ss store.Store) {
	})
}

func testChannelStoreGetMembersForUserWithCursor(t *testing.T, ss store.Store) {
	t1 := model.Team{}
	t1.DisplayName = "Team1"
	t1.Name = NewTestId()
	t1.Email = MakeEmail()
	t1.Type = model.TeamOpen
	_, err := ss.Team().Save(&t1)
	require.NoError(t, err)

	o1 := model.Channel{}
	o1.TeamId = t1.Id
	o1.DisplayName = "Channel1"
	o1.Name = NewTestId()
	o1.Type = model.ChannelTypeOpen
	_, nErr := ss.Channel().Save(&o1, -1)
	require.NoError(t, nErr)

	o2 := model.Channel{}
	o2.TeamId = o1.TeamId
	o2.DisplayName = "Channel2"
	o2.Name = NewTestId()
	o2.Type = model.ChannelTypeOpen
	_, nErr = ss.Channel().Save(&o2, -1)
	require.NoError(t, nErr)

	m1 := model.ChannelMember{}
	m1.ChannelId = o1.Id
	m1.UserId = model.NewId()
	m1.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err = ss.Channel().SaveMember(&m1)
	require.NoError(t, err)

	m2 := model.ChannelMember{}
	m2.ChannelId = o2.Id
	m2.UserId = m1.UserId
	m2.NotifyProps = model.GetDefaultChannelNotifyProps()
	_, err = ss.Channel().SaveMember(&m2)
	require.NoError(t, err)

	t.Run("with channels", func(t *testing.T) {
		var members model.ChannelMembers
		members, err = ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 1, 0)
		require.NoError(t, err)
		assert.Len(t, members, 1)
		members, err = ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 3, 0)
		require.NoError(t, err)
		assert.Len(t, members, 2)
		members, err = ss.Channel().GetMembersForUserWithCursor(m1.UserId, members[0].ChannelId, m1.UserId, 1, 0)
		require.NoError(t, err)
		assert.Len(t, members, 1)
	})

	t.Run("with channels and direct messages", func(t *testing.T) {
		user := model.User{Id: m1.UserId}
		u1 := model.User{Id: model.NewId()}
		u2 := model.User{Id: model.NewId()}
		u3 := model.User{Id: model.NewId()}
		u4 := model.User{Id: model.NewId()}
		_, nErr = ss.Channel().CreateDirectChannel(&u1, &user)
		require.NoError(t, nErr)
		_, nErr = ss.Channel().CreateDirectChannel(&u2, &user)
		require.NoError(t, nErr)
		// other user direct message
		_, nErr = ss.Channel().CreateDirectChannel(&u3, &u4)
		require.NoError(t, nErr)

		members, err2 := ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 10, 0)
		require.NoError(t, err2)
		assert.Len(t, members, 4)

		members, err2 = ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 2, 0)
		require.NoError(t, err2)
		assert.Len(t, members, 2)

		members, err2 = ss.Channel().GetMembersForUserWithCursor(m1.UserId, members[1].ChannelId, m1.UserId, 2, 0)
		require.NoError(t, err2)
		assert.Len(t, members, 2)
	})

	t.Run("with channels, direct channels and group messages", func(t *testing.T) {
		userIds := []string{model.NewId(), model.NewId(), model.NewId(), m1.UserId}
		group := &model.Channel{
			Name:        model.GetGroupNameFromUserIds(userIds),
			DisplayName: "test",
			Type:        model.ChannelTypeGroup,
		}
		var channel *model.Channel
		channel, nErr = ss.Channel().Save(group, 10000)
		require.NoError(t, nErr)
		for _, userId := range userIds {
			cm := &model.ChannelMember{
				UserId:      userId,
				ChannelId:   channel.Id,
				NotifyProps: model.GetDefaultChannelNotifyProps(),
				SchemeUser:  true,
			}

			_, err = ss.Channel().SaveMember(cm)
			require.NoError(t, err)
		}
		members, err := ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 10, 0)
		require.NoError(t, err)
		assert.Len(t, members, 5)

		members, err = ss.Channel().GetMembersForUserWithCursor(m1.UserId, "", "", 2, 0)
		require.NoError(t, err)
		assert.Len(t, members, 2)

		members, err = ss.Channel().GetMembersForUserWithCursor(m1.UserId, members[1].ChannelId, m1.UserId, 10, 0)
		require.NoError(t, err)
		assert.Len(t, members, 3)
	})
}

func testChannelStoreGetMembersForUserWithPagination(t *testing.T, ss store.Store) {
	t1 := model.Team{
		DisplayName: "team1",
	require.NoError(t, err)

	var times map[string]int64
	times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, m1.UserId)
	require.NoError(t, err, "failed to update ", err)
	require.Equal(t, o1.LastPostAt, times[o1.Id], "last viewed at time incorrect")

	times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId, m2.ChannelId}, m1.UserId)
	require.NoError(t, err, "failed to update ", err)
	require.Equal(t, o2.LastPostAt, times[o2.Id], "last viewed at time incorrect")

	assert.Equal(t, o2.LastPostAt, rm2.LastUpdateAt)
	assert.Equal(t, o2.TotalMsgCount, rm2.MsgCount)

	_, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, "missing id")
	require.NoError(t, err, "failed to update")
}

	_, err := ss.Channel().SaveMember(&m1)
	require.NoError(t, err)

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, m1.UserId, false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, "missing id", false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount("missing id", m1.UserId, false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount("missing id", "missing id", false)
	require.NoError(t, err, "failed to update")
}

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	_, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	_, err = ss.Group().Create(g2)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}

	_, err = ss.Group().Create(g3)

	t.Run("error", func(t *testing.T) {
		// trigger a SQL error
		s.GetMasterX().Exec("ALTER TABLE Channels RENAME TO Channels_renamed")
		defer s.GetMasterX().Exec("ALTER TABLE Channels_renamed RENAME TO Channels")

		list, err := ss.Channel().SearchArchivedInTeam(teamId, "term", userId)
		require.Error(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewString(model.NewId()),
	}
	_, err = ss.Group().Create(group)
	require.NoError(t, err)
	})
}

func testChannelStoreGetMembersInfoByChannelIds(t *testing.T, ss store.Store) {
	u, err := ss.User().Save(&model.User{
		Username: "user.test",
		Email:    MakeEmail(),
		Nickname: model.NewId(),
	})
	require.NoError(t, err)

	// Create a couple channels and add the user to them
	channel1, err := ss.Channel().Save(&model.Channel{
		TeamId:      model.NewId(),
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Type:        model.ChannelTypeOpen,
	}, -1)
	require.NoError(t, err)

	channel2, err := ss.Channel().Save(&model.Channel{
		TeamId:      model.NewId(),
		DisplayName: model.NewId(),
		Name:        model.NewId(),
		Type:        model.ChannelTypeOpen,
	}, -1)
	require.NoError(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   channel1.Id,
		UserId:      u.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.NoError(t, err)

	_, err = ss.Channel().SaveMember(&model.ChannelMember{
		ChannelId:   channel2.Id,
		UserId:      u.Id,
		NotifyProps: model.GetDefaultChannelNotifyProps(),
	})
	require.NoError(t, err)

	t.Run("should return the user's members for the given channels", func(t *testing.T) {
		result, nErr := ss.Channel().GetMembersInfoByChannelIds([]string{channel1.Id, channel2.Id})
		require.NoError(t, nErr)
		assert.Len(t, result, 2)
		for _, item := range result {
			assert.Len(t, item, 1)
			assert.Equal(t, u.Id, item[0].Id)
		}
	})

	t.Run("should not error or return anything for invalid channel IDs", func(t *testing.T) {
		_, err := ss.Channel().GetMembersInfoByChannelIds([]string{model.NewId(), model.NewId()})
		var nfErr *store.ErrNotFound
		require.True(t, errors.As(err, &nfErr))
	})
}

func testChannelStoreSearchGroupChannels(t *testing.T, ss store.Store) {
	// Users
	u1 := &model.User{}
	}()

	var openStartCount int64
	openStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypeOpen)
	require.NoError(t, nErr, nErr)

	var privateStartCount int64
	privateStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypePrivate)
	require.NoError(t, nErr, nErr)

	var directStartCount int64
	directStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypeDirect)
	require.NoError(t, nErr, nErr)

	nErr = ss.Channel().Delete(o1.Id, model.GetMillis())

	var count int64

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypeOpen)
	require.NoError(t, err, nErr)
	assert.Equal(t, openStartCount+2, count, "Wrong open channel deleted count.")

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypePrivate)
	require.NoError(t, nErr, nErr)
	assert.Equal(t, privateStartCount+1, count, "Wrong private channel deleted count.")

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", model.ChannelTypeDirect)
	require.NoError(t, nErr, nErr)
	assert.Equal(t, directStartCount+1, count, "Wrong direct channel deleted count.")
}
		Type:        model.ChannelTypeOpen,
	}

	_, execerr := s.GetMasterX().NamedExec(`
		INSERT INTO
		    PublicChannels(Id, DeleteAt, TeamId, DisplayName, Name, Header, Purpose)
		VALUES
		    (:id, :deleteat, :teamid, :displayname, :name, :header, :purpose);
	`, map[string]interface{}{
		"id":          o3.Id,
		"deleteat":    o3.DeleteAt,
		"teamid":      o3.TeamId,
		"displayname": o3.DisplayName,
		"name":        o3.Name,
		"header":      o3.Header,
		"purpose":     o3.Purpose,
	})
	require.NoError(t, execerr)

	o3.DisplayName = "Open Channel 3 - Modified"

	_, execerr = s.GetMasterX().NamedExec(`
		INSERT INTO
		    Channels(Id, CreateAt, UpdateAt, DeleteAt, TeamId, Type, DisplayName, Name, Header, Purpose, LastPostAt, LastRootPostAt, TotalMsgCount, ExtraUpdateAt, CreatorId, TotalMsgCountRoot)
		VALUES
				(:id, :createat, :updateat, :deleteat, :teamid, :type, :displayname, :name, :header, :purpose, :lastpostat, :lastrootpostat, :totalmsgcount, :extraupdateat, :creatorid, 0);
	`, map[string]interface{}{
		"id":             o3.Id,
		"createat":       o3.CreateAt,
		"updateat":       o3.UpdateAt,
		"deleteat":       o3.DeleteAt,
		"teamid":         o3.TeamId,
		"type":           o3.Type,
		"displayname":    o3.DisplayName,
		"name":           o3.Name,
		"header":         o3.Header,
		"purpose":        o3.Purpose,
		"lastpostat":     o3.LastPostAt,
		"lastrootpostat": o3.LastRootPostAt,
		"totalmsgcount":  o3.TotalMsgCount,
		"extraupdateat":  o3.ExtraUpdateAt,
		"creatorid":      o3.CreatorId,
	})
	require.NoError(t, execerr)

	_, nErr = ss.Channel().Save(&o4, -1)
	require.NoError(t, nErr)

	_, execerr = s.GetMasterX().Exec(`
		DELETE FROM
		    PublicChannels
		WHERE
		    Id = ?
	`, o4.Id)
	require.NoError(t, execerr)

	o4.DisplayName += " - Modified"
	assert.Equal(t, u3.Id, d2[0].UserId)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannels(t *testing.T, ss store.Store, s SqlStore) {
	assert.ElementsMatch(t, []string{o1.DisplayName, o2.DisplayName}, []string{d1[0].DisplayName, d1[1].DisplayName})

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsExcludePrivateAndPublic(t *testing.T, ss store.Store, s SqlStore) {
	assert.Equal(t, o1.DisplayName, d1[0].DisplayName)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsDeletedChannel(t *testing.T, ss store.Store, s SqlStore) {
	assert.Equal(t, 0, len(d1))

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMasterX().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsBatchForIndexing(t *testing.T, ss store.Store) {
