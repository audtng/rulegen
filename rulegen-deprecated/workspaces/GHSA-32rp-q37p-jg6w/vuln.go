package main


func TestGroupStore(t *testing.T, ss store.Store) {
	t.Run("Create", func(t *testing.T) { testGroupStoreCreate(t, ss) })
	t.Run("Get", func(t *testing.T) { testGroupStoreGet(t, ss) })
	t.Run("GetByName", func(t *testing.T) { testGroupStoreGetByName(t, ss) })
	t.Run("GetByIDs", func(t *testing.T) { testGroupStoreGetByIDs(t, ss) })
	t.Run("GetMemberUsersNotInChannel", func(t *testing.T) { testGroupGetMemberUsersNotInChannel(t, ss) })

	t.Run("UpsertMember", func(t *testing.T) { testUpsertMember(t, ss) })
	t.Run("DeleteMember", func(t *testing.T) { testGroupDeleteMember(t, ss) })
	t.Run("PermanentDeleteMembersByUser", func(t *testing.T) { testGroupPermanentDeleteMembersByUser(t, ss) })

	t.Run("CreateGroupSyncable", func(t *testing.T) { testCreateGroupSyncable(t, ss) })
	t.Run("GroupMemberCount", func(t *testing.T) { groupTestGroupMemberCount(t, ss) })
	t.Run("DistinctGroupMemberCount", func(t *testing.T) { groupTestDistinctGroupMemberCount(t, ss) })
	t.Run("GroupCountWithAllowReference", func(t *testing.T) { groupTestGroupCountWithAllowReference(t, ss) })
}

func testGroupStoreCreate(t *testing.T, ss store.Store) {
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}

	// Happy path
		Name:        model.NewString(model.NewId()),
		DisplayName: "",
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	data, err := ss.Group().Create(g2)
	require.Nil(t, data)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	_, err = ss.Group().Create(g4)
	require.NoError(t, err)
		Name:        g4.Name,
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	data, err = ss.Group().Create(g4b)
	require.Nil(t, data)
		DisplayName: strings.Repeat("x", model.GroupDisplayNameMaxLength),
		Description: strings.Repeat("x", model.GroupDescriptionMaxLength),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	require.Nil(t, g5.IsValidForCreate())

		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSource("fake"),
		RemoteId:    model.NewId(),
	}
	require.Equal(t, g6.IsValidForCreate().Id, "model.group.source.app_error")

		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	require.Equal(t, g7.IsValidForCreate().Id, "model.group.name.invalid_chars.app_error")
}

func testGroupStoreGet(t *testing.T, ss store.Store) {
	// Create a group
	g1 := &model.Group{
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	d1, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	g1Opts := model.GroupSearchOpts{
		FilterAllowReference: false,
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewId(),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	d1, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.Len(t, d1.Id, 26)

	// Get the group
	d2, err := ss.Group().GetByRemoteID(d1.RemoteId, model.GroupSourceLdap)
	require.NoError(t, err)
	require.Equal(t, d1.Id, d2.Id)
	require.Equal(t, *d1.Name, *d2.Name)
			DisplayName: model.NewId(),
			Description: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewId(),
		}
		groups = append(groups, g)
		_, err := ss.Group().Create(g)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	g1, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	g2, err = ss.Group().Create(g2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}

	// Create a group
	g1Update.Name = model.NewString(model.NewId())
	g1Update.DisplayName = model.NewId()
	g1Update.Description = model.NewId()
	g1Update.RemoteId = model.NewId()

	ud1, err := ss.Group().Update(g1Update)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: "",
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.Nil(t, data)
	require.Error(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	d2, err := ss.Group().Create(g2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	})
	require.Error(t, err)
	require.Contains(t, err.Error(), fmt.Sprintf("Group with name %s already exists", *g1Update.Name))

	// Cannot update CreateAt
	someVal := model.GetMillis()
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}

	d1, err := ss.Group().Create(g1)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.Equal(t, beforeRestoreCount+1, afterRestoreCount)
}

func testGroupDeleteMember(t *testing.T, ss store.Store) {
	// Create group
	g1 := &model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	require.True(t, errors.As(err, &nfErr))
}

func testGroupPermanentDeleteMembersByUser(t *testing.T, ss store.Store) {
	var g *model.Group
	var groups []*model.Group
			Name:        model.NewString(model.NewId()),
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			RemoteId:    model.NewId(),
		}
		group, err := ss.Group().Create(g)
		groups = append(groups, group)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Description: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(g1)
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group2, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "ChannelMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group2, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "TeamMembersToAdd Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: "Pending[Channel|Team]MemberRemovals Test Group",
		RemoteId:    model.NewId(),
		Source:      model.GroupSourceLdap,
	})
	require.NoError(t, err)
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-2",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-3",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group1, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId()),
		DisplayName:    "group-1",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	group2, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-2"),
		DisplayName:    "group-2",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
	})
	deletedGroup, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-deleted"),
		DisplayName:    "group-deleted",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: false,
		DeleteAt:       1,
	group3, err := ss.Group().Create(&model.Group{
		Name:           model.NewString(model.NewId() + "-group-3"),
		DisplayName:    "group-3",
		RemoteId:       model.NewId(),
		Source:         model.GroupSourceLdap,
		AllowReference: true,
	})
	_, err = ss.Group().UpsertMember(group1.Id, user2.Id)
	require.NoError(t, err)

	_, err = ss.Group().UpsertMember(deletedGroup.Id, user1.Id)
	require.NoError(t, err)

				return len(groups) > 0
			},
		},
	}

	for _, tc := range testCases {
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			Description: model.NewId(),
			RemoteId:    model.NewId(),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
			DisplayName: model.NewId(),
			Source:      model.GroupSourceLdap,
			Description: model.NewId(),
			RemoteId:    model.NewId(),
		}
		group, err := ss.Group().Create(group)
		require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group, err := ss.Group().Create(group)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group1, err = ss.Group().Create(group1)
	require.NoError(t, err)
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		Description: model.NewId(),
		RemoteId:    model.NewId(),
	}
	group2, err = ss.Group().Create(group2)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)
}

func groupTestGroupMemberCount(t *testing.T, ss store.Store) {
	group, err := ss.Group().Create(&model.Group{
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group.Id)

	member1, err := ss.Group().UpsertMember(group.Id, model.NewId())
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group.Id, member1.UserId)

	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))

	member2, err := ss.Group().UpsertMember(group.Id, model.NewId())
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group.Id, member2.UserId)

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group2.Id)

	member1, err := ss.Group().UpsertMember(group1.Id, model.NewId())
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group1.Id, member1.UserId)

	require.NoError(t, err)
	require.GreaterOrEqual(t, count, int64(1))

	member2, err := ss.Group().UpsertMember(group1.Id, model.NewId())
	require.NoError(t, err)
	defer ss.Group().DeleteMember(group1.Id, member2.UserId)

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	})
	require.NoError(t, err)
	defer ss.Group().Delete(group1.Id)
		Name:           model.NewString(model.NewId()),
		DisplayName:    model.NewId(),
		Source:         model.GroupSourceLdap,
		RemoteId:       model.NewId(),
		AllowReference: true,
	})
	require.NoError(t, err)
	require.NoError(t, err)
	require.Greater(t, countAfter, count)
}
	Purpose     string `json:"purpose"`
}

// channelInternal is a struct without the db:"-" tags
// which does not work with sqlx. This would be removed once we
// move to the new migration system.
type channelInternal struct {
	Id                string
	CreateAt          int64
	UpdateAt          int64
	DeleteAt          int64
	TeamId            string
	Type              model.ChannelType
	DisplayName       string
	Name              string
	Header            string
	Purpose           string
	LastPostAt        int64
	TotalMsgCount     int64
	ExtraUpdateAt     int64
	CreatorId         string
	SchemeId          *string
	Props             map[string]interface{}
	GroupConstrained  *bool
	Shared            *bool
	TotalMsgCountRoot int64
	PolicyId          *string
	LastRootPostAt    int64
}

func (ci *channelInternal) ToModel() *model.Channel {
	return &model.Channel{
		Id:                ci.Id,
		CreateAt:          ci.CreateAt,
		UpdateAt:          ci.UpdateAt,
		DeleteAt:          ci.DeleteAt,
		TeamId:            ci.TeamId,
		Type:              ci.Type,
		DisplayName:       ci.DisplayName,
		Name:              ci.Name,
		Header:            ci.Header,
		Purpose:           ci.Purpose,
		LastPostAt:        ci.LastPostAt,
		TotalMsgCount:     ci.TotalMsgCount,
		ExtraUpdateAt:     ci.ExtraUpdateAt,
		CreatorId:         ci.CreatorId,
		SchemeId:          ci.SchemeId,
		Props:             ci.Props,
		GroupConstrained:  ci.GroupConstrained,
		Shared:            ci.Shared,
		TotalMsgCountRoot: ci.TotalMsgCountRoot,
		PolicyID:          ci.PolicyId,
		LastRootPostAt:    ci.LastRootPostAt,
	}
}

type channelWithTeamDataInternal struct {
	channelInternal
	TeamDisplayName string
	TeamName        string
	TeamUpdateAt    int64
}

func (ctd *channelWithTeamDataInternal) ToModel() *model.ChannelWithTeamData {
	res := &model.ChannelWithTeamData{
		TeamDisplayName: ctd.TeamDisplayName,
		TeamName:        ctd.TeamName,
		TeamUpdateAt:    ctd.TeamUpdateAt,
	}
	res.Channel = *ctd.channelInternal.ToModel()
	return res
}

func channelWithTeamDataSliceToModel(channels []*channelWithTeamDataInternal) model.ChannelListWithTeamData {
	res := make(model.ChannelListWithTeamData, 0, len(channels))
	for _, ch := range channels {
		res = append(res, ch.ToModel())
	}
	return res
}

var allChannelMembersForUserCache = cache.NewLRU(cache.LRUOptions{
	Size: AllChannelMembersForUserCacheSize,
})
}

func newSqlChannelStore(sqlStore *SqlStore, metrics einterfaces.MetricsInterface) store.ChannelStore {
	s := &SqlChannelStore{
		SqlStore: sqlStore,
		metrics:  metrics,
	}

	for _, db := range sqlStore.GetAllConns() {
		table := db.AddTableWithName(model.Channel{}, "Channels").SetKeys(false, "Id")
		table.ColMap("Id").SetMaxSize(26)
		table.ColMap("TeamId").SetMaxSize(26)
		table.ColMap("Type").SetMaxSize(1)
		table.ColMap("DisplayName").SetMaxSize(64)
		table.ColMap("Name").SetMaxSize(64)
		table.SetUniqueTogether("Name", "TeamId")
		table.ColMap("Header").SetMaxSize(1024)
		table.ColMap("Purpose").SetMaxSize(250)
		table.ColMap("CreatorId").SetMaxSize(26)
		table.ColMap("SchemeId").SetMaxSize(26)
		table.ColMap("LastRootPostAt").SetDefaultConstraint(model.NewString("0"))

		tablem := db.AddTableWithName(channelMember{}, "ChannelMembers").SetKeys(false, "ChannelId", "UserId")
		tablem.ColMap("ChannelId").SetMaxSize(26)
		tablem.ColMap("UserId").SetMaxSize(26)
		tablem.ColMap("Roles").SetMaxSize(model.UserRolesMaxLength)
		tablem.ColMap("NotifyProps").SetDataType(sqlStore.jsonDataType())

		tablePublicChannels := db.AddTableWithName(publicChannel{}, "PublicChannels").SetKeys(false, "Id")
		tablePublicChannels.ColMap("Id").SetMaxSize(26)
		tablePublicChannels.ColMap("TeamId").SetMaxSize(26)
		tablePublicChannels.ColMap("DisplayName").SetMaxSize(64)
		tablePublicChannels.ColMap("Name").SetMaxSize(64)
		tablePublicChannels.SetUniqueTogether("Name", "TeamId")
		tablePublicChannels.ColMap("Header").SetMaxSize(1024)
		tablePublicChannels.ColMap("Purpose").SetMaxSize(250)

		tableSidebarCategories := db.AddTableWithName(model.SidebarCategory{}, "SidebarCategories").SetKeys(false, "Id")
		tableSidebarCategories.ColMap("Id").SetMaxSize(128)
		tableSidebarCategories.ColMap("UserId").SetMaxSize(26)
		tableSidebarCategories.ColMap("TeamId").SetMaxSize(26)
		tableSidebarCategories.ColMap("Sorting").SetMaxSize(64)
		tableSidebarCategories.ColMap("Type").SetMaxSize(64)
		tableSidebarCategories.ColMap("DisplayName").SetMaxSize(64)

		tableSidebarChannels := db.AddTableWithName(model.SidebarChannel{}, "SidebarChannels").SetKeys(false, "ChannelId", "UserId", "CategoryId")
		tableSidebarChannels.ColMap("ChannelId").SetMaxSize(26)
		tableSidebarChannels.ColMap("UserId").SetMaxSize(26)
		tableSidebarChannels.ColMap("CategoryId").SetMaxSize(128)
	}

	return s
}

func (s SqlChannelStore) upsertPublicChannelT(transaction *sqlxTxWrapper, channel *model.Channel) error {
			    PublicChannels(Id, DeleteAt, TeamId, DisplayName, Name, Header, Purpose)
			VALUES
			    (:id, :deleteat, :teamid, :displayname, :name, :header, :purpose)
			ON DUPLICATE KEY UPDATE
			    DeleteAt = :deleteat,
			    TeamId = :teamid,
			    DisplayName = :displayname,
			    Name = :name,
			    Header = :header,
			    Purpose = :purpose;
		`, vals)
	} else {
		_, err = transaction.NamedExec(`
			INSERT INTO

	if channel.Type != model.ChannelTypeDirect && channel.Type != model.ChannelTypeGroup && maxChannelsPerTeam >= 0 {
		var count int64
		if err := transaction.Get(&count, "SELECT COUNT(0) FROM Channels WHERE TeamId = ? AND DeleteAt = 0 AND (Type = 'O' OR Type = 'P')", channel.TeamId); err != nil {
			return nil, errors.Wrapf(err, "save_channel_count: teamId=%s", channel.TeamId)
		} else if count >= maxChannelsPerTeam {
			return nil, store.NewErrLimitExceeded("channels_per_team", int(count), "teamId="+channel.TeamId)
func (s SqlChannelStore) GetPinnedPosts(channelId string) (*model.PostList, error) {
	pl := model.NewPostList()

	posts := []*postInternal{}
	if err := s.GetReplicaX().Select(&posts, "SELECT *, (SELECT count(Posts.Id) FROM Posts WHERE Posts.RootId = (CASE WHEN p.RootId = '' THEN p.Id ELSE p.RootId END) AND Posts.DeleteAt = 0) as ReplyCount  FROM Posts p WHERE IsPinned = true AND ChannelId = ? AND DeleteAt = 0 ORDER BY CreateAt ASC", channelId); err != nil {
		return nil, errors.Wrap(err, "failed to find Posts")
	}
	for _, post := range posts {
		pl.AddPost(post.ToModel())
		pl.AddOrder(post.Id)
	}
	return pl, nil
	return nil
}

func (s SqlChannelStore) GetChannels(teamId string, userId string, includeDeleted bool, lastDeleteAt int) (model.ChannelList, error) {
	query := s.getQueryBuilder().
		Select("Channels.*").
		From("Channels, ChannelMembers").
		Where(
			sq.And{
				sq.Expr("Id = ChannelId"),
				sq.Eq{"UserId": userId},
			},
		).
		OrderBy("DisplayName")

	if teamId != "" {
		query = query.Where(sq.Or{
			sq.Eq{"TeamId": teamId},
			sq.Eq{"TeamId": ""},
		})
	}

	if includeDeleted {
		if lastDeleteAt != 0 {
			// We filter by non-archived, and archived >= a timestamp.
			query = query.Where(sq.Or{
				sq.Eq{"DeleteAt": 0},
				sq.GtOrEq{"DeleteAt": lastDeleteAt},
			})
		}
		// If lastDeleteAt is not set, we include everything. That means no filter is needed.
	} else {
		// Don't include archived channels.
		query = query.Where(sq.Eq{"DeleteAt": 0})
	}

	channels := model.ChannelList{}
		return nil, errors.Wrap(err, "failed to create query")
	}

	data := []*channelWithTeamDataInternal{}
	err = s.GetReplicaX().Select(&data, queryString, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get all channels")
	}

	return channelWithTeamDataSliceToModel(data), nil
}

func (s SqlChannelStore) GetAllChannelsCount(opts store.ChannelSearchOpts) (int64, error) {
	} else {
		selectStr = "c.*, Teams.DisplayName AS TeamDisplayName, Teams.Name AS TeamName, Teams.UpdateAt AS TeamUpdateAt"
		if opts.IncludePolicyID {
			selectStr += ", RetentionPoliciesChannels.PolicyId"
		}
	}


func (s SqlChannelStore) GetTeamChannels(teamId string) (model.ChannelList, error) {
	data := model.ChannelList{}
	err := s.GetReplicaX().Select(&data, "SELECT * FROM Channels WHERE TeamId = ? And Type != 'D' ORDER BY DisplayName", teamId)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Channels with teamId=%s", teamId)
	}
		SELECT * FROM Channels
		WHERE (TeamId = ? OR TeamId = '')
		AND DeleteAt != 0
		AND Type != 'P'
		UNION
			SELECT * FROM Channels
			WHERE (TeamId = ? OR TeamId = '')
			AND DeleteAt != 0
			AND Type = 'P'
			AND Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = ?)
		ORDER BY DisplayName LIMIT ? OFFSET ?
	`

	if err := s.GetReplicaX().Select(&channels, query, teamId, teamId, userId, limit, offset); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.NewErrNotFound("Channel", fmt.Sprintf("TeamId=%s,UserId=%s", teamId, userId))
		}
	query := s.getQueryBuilder().
		Select(selectStr).
		From("ChannelMembers").
		Join("GroupMembers ON GroupMembers.UserId = ChannelMembers.UserId")

	if includeTimezones {
		query = query.Join("Users ON Users.Id = GroupMembers.UserId")
	return nil
}

// TODO: convert to squirrel (https://github.com/mattermost/mattermost-server/issues/19332)
func (s SqlChannelStore) UpdateLastViewedAt(channelIds []string, userId string, updateThreads bool) (map[string]int64, error) {
	var threadsToUpdate []string
	now := model.GetMillis()
	if updateThreads {
		var err error
		threadsToUpdate, err = s.Thread().CollectThreadsWithNewerReplies(userId, channelIds, now)
		if err != nil {
			return nil, err
		}
	}

	keys, props := MapStringsToQueryParams(channelIds, "Channel")
	props["UserId"] = userId

	var lastPostAtTimes []struct {
		Id                string
		LastPostAt        int64
		TotalMsgCount     int64
		TotalMsgCountRoot int64
	}

	query := `SELECT Id, LastPostAt, TotalMsgCount, TotalMsgCountRoot FROM Channels WHERE Id IN ` + keys
	// TODO: use a CTE for mysql too when version 8 becomes the minimum supported version.
	if s.DriverName() == model.DatabaseDriverPostgres {
		query = `WITH c AS ( ` + query + `),
	updated AS (
	UPDATE
		ChannelMembers cm
	SET
		MentionCount = 0,
		MentionCountRoot = 0,
		MsgCount = greatest(cm.MsgCount, c.TotalMsgCount),
		MsgCountRoot = greatest(cm.MsgCountRoot, c.TotalMsgCountRoot),
		LastViewedAt = greatest(cm.LastViewedAt, c.LastPostAt),
		LastUpdateAt = greatest(cm.LastViewedAt, c.LastPostAt)
	FROM c
		WHERE cm.UserId = :UserId
		AND c.Id=cm.ChannelId
)
	SELECT Id, LastPostAt FROM c`
	}

	_, err := s.GetMaster().Select(&lastPostAtTimes, query, props)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to find ChannelMembers data with userId=%s and channelId in %v", userId, channelIds)
	}
		for _, t := range lastPostAtTimes {
			times[t.Id] = t.LastPostAt
		}
		if updateThreads {
			s.Thread().UpdateUnreadsByChannel(userId, threadsToUpdate, now, true)
		}
		return times, nil
	}

	msgCountQuery := ""
	msgCountQueryRoot := ""
	lastViewedQuery := ""

	for index, t := range lastPostAtTimes {
		times[t.Id] = t.LastPostAt

		props["msgCount"+strconv.Itoa(index)] = t.TotalMsgCount
		msgCountQuery += fmt.Sprintf("WHEN :channelId%d THEN GREATEST(MsgCount, :msgCount%d) ", index, index)

		props["msgCountRoot"+strconv.Itoa(index)] = t.TotalMsgCountRoot
		msgCountQueryRoot += fmt.Sprintf("WHEN :channelId%d THEN GREATEST(MsgCountRoot, :msgCountRoot%d) ", index, index)

		props["lastViewed"+strconv.Itoa(index)] = t.LastPostAt
		lastViewedQuery += fmt.Sprintf("WHEN :channelId%d THEN GREATEST(LastViewedAt, :lastViewed%d) ", index, index)

		props["channelId"+strconv.Itoa(index)] = t.Id
	}

	updateQuery := `UPDATE
			ChannelMembers
		SET
			MentionCount = 0,
			MentionCountRoot = 0,
			MsgCount = CASE ChannelId ` + msgCountQuery + ` END,
			MsgCountRoot = CASE ChannelId ` + msgCountQueryRoot + ` END,
			LastViewedAt = CASE ChannelId ` + lastViewedQuery + ` END,
			LastUpdateAt = LastViewedAt
		WHERE
				UserId = :UserId
				AND ChannelId IN ` + keys

	if _, err := s.GetMaster().Exec(updateQuery, props); err != nil {
		return nil, errors.Wrapf(err, "failed to update ChannelMembers with userId=%s and channelId in %v", userId, channelIds)
	}

	if updateThreads {
		s.Thread().UpdateUnreadsByChannel(userId, threadsToUpdate, now, true)
	}
	return times, nil
}

// UpdateLastViewedAtPost updates a ChannelMember as if the user last read the channel at the time of the given post.
// If the provided mentionCount is -1, the given post and all posts after it are considered to be mentions. Returns
// an updated model.ChannelUnreadAt that can be returned to the client.
func (s SqlChannelStore) UpdateLastViewedAtPost(unreadPost *model.Post, userID string, mentionCount, mentionCountRoot int, updateThreads bool, setUnreadCountRoot bool) (*model.ChannelUnreadAt, error) {
	var threadsToUpdate []string
	unreadDate := unreadPost.CreateAt - 1
	if updateThreads {
		var err error
		threadsToUpdate, err = s.Thread().CollectThreadsWithNewerReplies(userID, []string{unreadPost.ChannelId}, unreadDate)
		if err != nil {
			return nil, err
		}
	}

	unread, unreadRoot, err := s.CountPostsAfter(unreadPost.ChannelId, unreadDate, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get ChannelMember with channelId=%s", unreadPost.ChannelId)
	}

	if updateThreads {
		s.Thread().UpdateUnreadsByChannel(userID, threadsToUpdate, unreadDate, false)
	}
	return result, nil
}

func (s SqlChannelStore) IncrementMentionCount(channelId string, userId string, updateThreads, isRoot bool) error {
	now := model.GetMillis()
	var threadsToUpdate []string
	if updateThreads {
		var err error
		threadsToUpdate, err = s.Thread().CollectThreadsWithNewerReplies(userId, []string{channelId}, now)
		if err != nil {
			return err
		}
	}
	rootInc := 0
	if isRoot {
		rootInc = 1
	if err != nil {
		return errors.Wrapf(err, "failed to Update ChannelMembers with channelId=%s and userId=%s", channelId, userId)
	}
	if updateThreads {
		s.Thread().UpdateUnreadsByChannel(userId, threadsToUpdate, now, false)
	}
	return nil
}

func (s SqlChannelStore) GetAll(teamId string) ([]*model.Channel, error) {
	data := []*model.Channel{}
	err := s.GetReplicaX().Select(&data, "SELECT * FROM Channels WHERE TeamId = ? AND Type != 'D' ORDER BY Name", teamId)

	if err != nil {
		return nil, errors.Wrapf(err, "failed to find Channels with teamId=%s", teamId)
	return value, nil
}

func (s SqlChannelStore) AnalyticsDeletedTypeCount(teamId string, channelType string) (int64, error) {
	query := s.getQueryBuilder().
		Select("COUNT(Id) AS Value").
		From("Channels").
	return dbMembers.ToModel(), nil
}

func (s SqlChannelStore) GetMembersForUserWithPagination(userId string, page, perPage int) (model.ChannelMembersWithTeamData, error) {
	dbMembers := channelMemberWithTeamWithSchemeRolesList{}
	offset := page * perPage
		`+deleteFilter+`
		SEARCH_CLAUSE
		AND (
			c.Type != 'P'
			OR (
				c.Type = 'P'
				AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
			)
		)
	ORDER BY c.DisplayName
	`, term, map[string]interface{}{
		"UserId": userID,
	})
}

		`+deleteFilter+`
		SEARCH_CLAUSE
		AND (
			c.Type != 'P'
			OR (
				c.Type = 'P'
				AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
			)
		)
	ORDER BY c.DisplayName
	LIMIT :Limit
	`, term, map[string]interface{}{
		"TeamId": teamID,
		"UserId": userID,
		"Limit":  model.ChannelSearchDefaultLimit,
	})
}

		JOIN
			ChannelMembers AS CM ON CM.ChannelId = C.Id
		WHERE
			(C.TeamId = :TeamId OR (C.TeamId = '' AND C.Type = 'G'))
			AND CM.UserId = :UserId
			` + deleteFilter + `
			%v
	var channels model.ChannelList

	if likeClause, likeTerm := s.buildLIKEClause(term, "Name, DisplayName, Purpose"); likeClause == "" {
		if _, err := s.GetReplica().Select(&channels, fmt.Sprintf(queryFormat, ""), map[string]interface{}{"TeamId": teamId, "UserId": userId}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	} else {
		fulltextQuery := fmt.Sprintf(queryFormat, "AND "+fulltextClause)
		query := fmt.Sprintf("(%v) UNION (%v) LIMIT 50", likeQuery, fulltextQuery)

		if _, err := s.GetReplica().Select(&channels, query, map[string]interface{}{"TeamId": teamId, "UserId": userId, "LikeTerm": likeTerm, "FulltextTerm": fulltextTerm}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	}
					%v
				) AS OtherUsers ON OtherUsers.ChannelId = C.Id
			WHERE
			    C.Type = 'D'
				AND CM.UserId = :UserId
			LIMIT 50`

	var channels model.ChannelList

	if likeClause, likeTerm := s.buildLIKEClause(term, "IU.Username, IU.Nickname"); likeClause == "" {
		if _, err := s.GetReplica().Select(&channels, fmt.Sprintf(queryFormat, ""), map[string]interface{}{"UserId": userId}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	} else {
		query := fmt.Sprintf(queryFormat, "AND "+likeClause)

		if _, err := s.GetReplica().Select(&channels, query, map[string]interface{}{"UserId": userId, "LikeTerm": likeTerm}); err != nil {
			return nil, errors.Wrapf(err, "failed to find Channels with term='%s'", term)
		}
	}
			c.TeamId = :TeamId
			SEARCH_CLAUSE
			AND c.DeleteAt != 0
			AND c.Type != 'P'
		ORDER BY c.DisplayName
		LIMIT 100
		`, term, map[string]interface{}{
		"TeamId": teamId,
		"UserId": userId,
	})

	privateChannels, privateErr := s.performSearch(`
			c.TeamId = :TeamId
			SEARCH_CLAUSE
			AND c.DeleteAt != 0
			AND c.Type = 'P'
			AND c.Id IN (SELECT ChannelId FROM ChannelMembers WHERE UserId = :UserId)
		ORDER BY c.DisplayName
		LIMIT 100
		`, term, map[string]interface{}{
		"TeamId": teamId,
		"UserId": userId,
	})

	outputErr := publicErr
			selectStr += ", t.DisplayName AS TeamDisplayName, t.Name AS TeamName, t.UpdateAt as TeamUpdateAt"
		}
		if opts.IncludePolicyID {
			selectStr += ", RetentionPoliciesChannels.PolicyId"
		}
	}

	if err != nil {
		return nil, 0, errors.Wrap(err, "channel_tosql")
	}
	channels := []*channelWithTeamDataInternal{}
	if err2 := s.GetReplicaX().Select(&channels, queryString, args...); err2 != nil {
		return nil, 0, errors.Wrapf(err2, "failed to find Channels with term='%s'", term)
	}
		totalCount = int64(len(channels))
	}

	return channelWithTeamDataSliceToModel(channels), totalCount, nil
}

// TODO: rewrite in squrrel
                        JOIN
                            Users u on u.Id = cm.UserId
                        WHERE
                            c.Type = 'G'
                        AND
                            u.Id = :UserId
                        GROUP BY
                JOIN
                    Users u on u.Id = cm.UserId
                WHERE
                    c.Type = 'G'
                AND
                    u.Id = :UserId
                GROUP BY
	}

	var likeClauses []string
	args := map[string]interface{}{"UserId": userId}
	terms := strings.Split(strings.ToLower(strings.Trim(term, " ")), " ")

	for idx, term := range terms {
	return dbMembers.ToModel(), nil
}

func (s SqlChannelStore) GetChannelsByScheme(schemeId string, offset int, limit int) (model.ChannelList, error) {
	channels := model.ChannelList{}
	err := s.GetReplicaX().Select(&channels, "SELECT * FROM Channels WHERE SchemeId = ? ORDER BY DisplayName LIMIT ? OFFSET ?", schemeId, limit, offset)
			Schemes ON Channels.SchemeId = Schemes.Id
		WHERE
			Channels.Id > ?
			AND Channels.Type IN ('O', 'P')
		ORDER BY
			Id
		LIMIT ?`,
		afterId, limit); err != nil {
		return nil, errors.Wrap(err, "failed to find Channels for export")
	}

		Where(sq.And{
			sq.Gt{"Channels.Id": afterId},
			sq.Eq{"Channels.DeleteAt": int(0)},
			sq.Eq{"Channels.Type": []string{"D", "G"}},
		}).
		OrderBy("Channels.Id").
		Limit(uint64(limit))

import (
	"context"
	"encoding/json"
	"errors"
	"sort"
	"testing"
	"time"

	"github.com/mattermost/gorp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

)

type SqlStore interface {
	GetMaster() *gorp.DbMap
	DriverName() string
}

func cleanupChannels(t *testing.T, ss store.Store) {
	list, err := ss.Channel().GetAllChannels(0, 100000, store.ChannelSearchOpts{IncludeDeleted: true})
	require.NoError(t, err, "error cleaning all channels", err)
	t.Run("RemoveMembers", func(t *testing.T) { testChannelRemoveMembers(t, ss) })
	t.Run("ChannelDeleteMemberStore", func(t *testing.T) { testChannelDeleteMemberStore(t, ss) })
	t.Run("GetChannels", func(t *testing.T) { testChannelStoreGetChannels(t, ss) })
	t.Run("GetChannelsByUser", func(t *testing.T) { testChannelStoreGetChannelsByUser(t, ss) })
	t.Run("GetAllChannels", func(t *testing.T) { testChannelStoreGetAllChannels(t, ss, s) })
	t.Run("GetMoreChannels", func(t *testing.T) { testChannelStoreGetMoreChannels(t, ss) })
	t.Run("GetPublicChannelsByIdsForTeam", func(t *testing.T) { testChannelStoreGetPublicChannelsByIdsForTeam(t, ss) })
	t.Run("GetChannelCounts", func(t *testing.T) { testChannelStoreGetChannelCounts(t, ss) })
	t.Run("GetMembersForUser", func(t *testing.T) { testChannelStoreGetMembersForUser(t, ss) })
	t.Run("GetMembersForUserWithPagination", func(t *testing.T) { testChannelStoreGetMembersForUserWithPagination(t, ss) })
	t.Run("CountPostsAfter", func(t *testing.T) { testCountPostsAfter(t, ss) })
	t.Run("UpdateLastViewedAt", func(t *testing.T) { testChannelStoreUpdateLastViewedAt(t, ss) })
	t.Run("SearchAllChannels", func(t *testing.T) { testChannelStoreSearchAllChannels(t, ss) })
	t.Run("GetMembersByIds", func(t *testing.T) { testChannelStoreGetMembersByIds(t, ss) })
	t.Run("GetMembersByChannelIds", func(t *testing.T) { testChannelStoreGetMembersByChannelIds(t, ss) })
	t.Run("SearchGroupChannels", func(t *testing.T) { testChannelStoreSearchGroupChannels(t, ss) })
	t.Run("AnalyticsDeletedTypeCount", func(t *testing.T) { testChannelStoreAnalyticsDeletedTypeCount(t, ss) })
	t.Run("GetPinnedPosts", func(t *testing.T) { testChannelStoreGetPinnedPosts(t, ss) })
	require.Len(t, members, 1, "should have saved just 1 member")

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreCreateDirectChannel(t *testing.T, ss store.Store) {
	require.True(t, errors.As(err, &nfErr))

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsByIds(t *testing.T, ss store.Store) {
	nErr = ss.Channel().Delete(o3.Id, model.GetMillis())
	require.NoError(t, nErr, nErr)

	list, nErr := ss.Channel().GetChannels(o1.TeamId, m1.UserId, false, 0)
	require.NoError(t, nErr)
	require.Len(t, list, 1, "invalid number of channels")

	cresult := ss.Channel().PermanentDelete(o2.Id)
	require.NoError(t, cresult)

	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, false, 0)
	if assert.Error(t, nErr) {
		var nfErr *store.ErrNotFound
		require.True(t, errors.As(nErr, &nfErr))

func testChannelStoreGetChannels(t *testing.T, ss store.Store) {
	team := model.NewId()
	o1 := model.Channel{}
	o1.TeamId = team
	o1.DisplayName = "Channel1"
	o1.Name = NewTestId()
	o1.Type = model.ChannelTypeOpen
	_, nErr := ss.Channel().Save(&o1, -1)
	require.NoError(t, nErr)

	o2 := model.Channel{}
	_, err = ss.Channel().SaveMember(&m4)
	require.NoError(t, err)

	list, nErr := ss.Channel().GetChannels(o1.TeamId, m1.UserId, false, 0)
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	_, ok = ids4[o1.Id]
	require.True(t, ok, "missing channel")

	nErr = ss.Channel().Delete(o2.Id, 10)
	require.NoError(t, nErr)

	require.NoError(t, nErr)

	// should return 1
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, false, 0)
	require.NoError(t, nErr)
	require.Len(t, list, 1)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")

	// Should return all
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, true, 0)
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	require.Equal(t, o3.Id, list[2].Id, "missing channel")

	// Should still return all
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, true, 10)
	require.NoError(t, nErr)
	require.Len(t, list, 3)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	require.Equal(t, o3.Id, list[2].Id, "missing channel")

	// Should return 2
	list, nErr = ss.Channel().GetChannels(o1.TeamId, m1.UserId, true, 20)
	require.NoError(t, nErr)
	require.Len(t, list, 2)
	require.Equal(t, o1.Id, list[0].Id, "missing channel")
	ss.Channel().InvalidateAllChannelMembersForUser(m1.UserId)
}

func testChannelStoreGetChannelsByUser(t *testing.T, ss store.Store) {
	team := model.NewId()
	team2 := model.NewId()
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	_, err = ss.Group().Create(group)
	require.NoError(t, err)
	assert.Equal(t, *list[0].PolicyID, policy.ID)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetMoreChannels(t *testing.T, ss store.Store) {
	})
}

func testChannelStoreGetMembersForUserWithPagination(t *testing.T, ss store.Store) {
	t1 := model.Team{
		DisplayName: "team1",
	require.NoError(t, err)

	var times map[string]int64
	times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, m1.UserId, false)
	require.NoError(t, err, "failed to update ", err)
	require.Equal(t, o1.LastPostAt, times[o1.Id], "last viewed at time incorrect")

	times, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId, m2.ChannelId}, m1.UserId, false)
	require.NoError(t, err, "failed to update ", err)
	require.Equal(t, o2.LastPostAt, times[o2.Id], "last viewed at time incorrect")

	assert.Equal(t, o2.LastPostAt, rm2.LastUpdateAt)
	assert.Equal(t, o2.TotalMsgCount, rm2.MsgCount)

	_, err = ss.Channel().UpdateLastViewedAt([]string{m1.ChannelId}, "missing id", false)
	require.NoError(t, err, "failed to update")
}

	_, err := ss.Channel().SaveMember(&m1)
	require.NoError(t, err)

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, m1.UserId, false, false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount(m1.ChannelId, "missing id", false, false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount("missing id", m1.UserId, false, false)
	require.NoError(t, err, "failed to update")

	err = ss.Channel().IncrementMentionCount("missing id", "missing id", false, false)
	require.NoError(t, err, "failed to update")
}

		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	_, err := ss.Group().Create(g1)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	_, err = ss.Group().Create(g2)
	require.NoError(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}

	_, err = ss.Group().Create(g3)

	t.Run("error", func(t *testing.T) {
		// trigger a SQL error
		s.GetMaster().Exec("ALTER TABLE Channels RENAME TO Channels_renamed")
		defer s.GetMaster().Exec("ALTER TABLE Channels_renamed RENAME TO Channels")

		list, err := ss.Channel().SearchArchivedInTeam(teamId, "term", userId)
		require.Error(t, err)
		Name:        model.NewString(model.NewId()),
		DisplayName: model.NewId(),
		Source:      model.GroupSourceLdap,
		RemoteId:    model.NewId(),
	}
	_, err = ss.Group().Create(group)
	require.NoError(t, err)
	})
}

func testChannelStoreSearchGroupChannels(t *testing.T, ss store.Store) {
	// Users
	u1 := &model.User{}
	}()

	var openStartCount int64
	openStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "O")
	require.NoError(t, nErr, nErr)

	var privateStartCount int64
	privateStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "P")
	require.NoError(t, nErr, nErr)

	var directStartCount int64
	directStartCount, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "D")
	require.NoError(t, nErr, nErr)

	nErr = ss.Channel().Delete(o1.Id, model.GetMillis())

	var count int64

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "O")
	require.NoError(t, err, nErr)
	assert.Equal(t, openStartCount+2, count, "Wrong open channel deleted count.")

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "P")
	require.NoError(t, nErr, nErr)
	assert.Equal(t, privateStartCount+1, count, "Wrong private channel deleted count.")

	count, nErr = ss.Channel().AnalyticsDeletedTypeCount("", "D")
	require.NoError(t, nErr, nErr)
	assert.Equal(t, directStartCount+1, count, "Wrong direct channel deleted count.")
}
		Type:        model.ChannelTypeOpen,
	}

	_, execerr := s.GetMaster().ExecNoTimeout(`
		INSERT INTO
		    PublicChannels(Id, DeleteAt, TeamId, DisplayName, Name, Header, Purpose)
		VALUES
		    (:Id, :DeleteAt, :TeamId, :DisplayName, :Name, :Header, :Purpose);
	`, map[string]interface{}{
		"Id":          o3.Id,
		"DeleteAt":    o3.DeleteAt,
		"TeamId":      o3.TeamId,
		"DisplayName": o3.DisplayName,
		"Name":        o3.Name,
		"Header":      o3.Header,
		"Purpose":     o3.Purpose,
	})
	require.NoError(t, execerr)

	o3.DisplayName = "Open Channel 3 - Modified"

	_, execerr = s.GetMaster().ExecNoTimeout(`
		INSERT INTO
		    Channels(Id, CreateAt, UpdateAt, DeleteAt, TeamId, Type, DisplayName, Name, Header, Purpose, LastPostAt, LastRootPostAt, TotalMsgCount, ExtraUpdateAt, CreatorId, TotalMsgCountRoot)
		VALUES
				(:Id, :CreateAt, :UpdateAt, :DeleteAt, :TeamId, :Type, :DisplayName, :Name, :Header, :Purpose, :LastPostAt, :LastRootPostAt, :TotalMsgCount, :ExtraUpdateAt, :CreatorId, 0);
	`, map[string]interface{}{
		"Id":             o3.Id,
		"CreateAt":       o3.CreateAt,
		"UpdateAt":       o3.UpdateAt,
		"DeleteAt":       o3.DeleteAt,
		"TeamId":         o3.TeamId,
		"Type":           o3.Type,
		"DisplayName":    o3.DisplayName,
		"Name":           o3.Name,
		"Header":         o3.Header,
		"Purpose":        o3.Purpose,
		"LastPostAt":     o3.LastPostAt,
		"LastRootPostAt": o3.LastRootPostAt,
		"TotalMsgCount":  o3.TotalMsgCount,
		"ExtraUpdateAt":  o3.ExtraUpdateAt,
		"CreatorId":      o3.CreatorId,
	})
	require.NoError(t, execerr)

	_, nErr = ss.Channel().Save(&o4, -1)
	require.NoError(t, nErr)

	_, execerr = s.GetMaster().ExecNoTimeout(`
		DELETE FROM
		    PublicChannels
		WHERE
		    Id = :Id
	`, map[string]interface{}{
		"Id": o4.Id,
	})
	require.NoError(t, execerr)

	o4.DisplayName += " - Modified"
	assert.Equal(t, u3.Id, d2[0].UserId)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannels(t *testing.T, ss store.Store, s SqlStore) {
	assert.ElementsMatch(t, []string{o1.DisplayName, o2.DisplayName}, []string{d1[0].DisplayName, d1[1].DisplayName})

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsExcludePrivateAndPublic(t *testing.T, ss store.Store, s SqlStore) {
	assert.Equal(t, o1.DisplayName, d1[0].DisplayName)

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreExportAllDirectChannelsDeletedChannel(t *testing.T, ss store.Store, s SqlStore) {
	assert.Equal(t, 0, len(d1))

	// Manually truncate Channels table until testlib can handle cleanups
	s.GetMaster().Exec("TRUNCATE Channels")
}

func testChannelStoreGetChannelsBatchForIndexing(t *testing.T, ss store.Store) {
