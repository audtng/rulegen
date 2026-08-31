package main

	})

	t.Run("save-second-reaction", func(t *testing.T) {
		reaction.EmojiName = "cry"

		rr, _, err := client.SaveReaction(context.Background(), reaction)
		require.NoError(t, err)
	r2 := &model.Reaction{
		UserId:    userId,
		PostId:    postId,
		EmojiName: "cry",
	}

	r3 := &model.Reaction{
	r4 := &model.Reaction{
		UserId:    user2Id,
		PostId:    postId,
		EmojiName: "grin",
	}

	// Check the appropriate permissions are enforced.
	reactionObject := model.Reaction{
		UserId:    th.BasicUser.Id,
		PostId:    post.Id,
		EmojiName: "smile",
		CreateAt:  model.GetMillis(),
	}
	reactionObjectDeleted := model.Reaction{
		UserId:    th.BasicUser2.Id,
		PostId:    post.Id,
		EmojiName: "smile",
		CreateAt:  model.GetMillis(),
	}

	_, err := th.App.SaveReactionForPost(th.Context, &reactionObject)
	require.Nil(t, err)
	_, err = th.App.SaveReactionForPost(th.Context, &reactionObjectDeleted)
	require.Nil(t, err)
	reactionsOfPost, err := th.App.BuildPostReactions(th.Context, post.Id)
	require.Nil(t, err)

		return nil, err
	}

	// Check whether this is a valid emoji
	if _, ok := model.GetSystemEmojiId(reaction.EmojiName); !ok {
		if _, emojiErr := a.GetEmojiByName(c, reaction.EmojiName); emojiErr != nil {
			return nil, emojiErr
		}
	}

	existing, dErr := a.Srv().Store().Reaction().ExistsOnPost(reaction.PostId, reaction.EmojiName)
	if dErr != nil {
		return nil, model.NewAppError("SaveReactionForPost", "app.reaction.save.save.app_error", nil, "", http.StatusInternalServerError).Wrap(dErr)
	}

	// If it exists already, we don't need to check for the limit
	if !existing {
		count, dErr := a.Srv().Store().Reaction().GetUniqueCountForPost(reaction.PostId)
		if dErr != nil {
			return nil, model.NewAppError("SaveReactionForPost", "app.reaction.save.save.app_error", nil, "", http.StatusInternalServerError).Wrap(dErr)
		}

		if count >= *a.Config().ServiceSettings.UniqueEmojiReactionLimitPerPost {
			return nil, model.NewAppError("SaveReactionForPost", "app.reaction.save.save.too_many_reactions", nil, "", http.StatusBadRequest)
		}
	}

	channel, err := a.GetChannel(c, post.ChannelId)
	if err != nil {
		return nil, err
	"github.com/mattermost/mattermost/server/v8/channels/testlib"
)

func TestSaveReactionForPost(t *testing.T) {
	th := Setup(t).InitBasic()

	post := th.CreatePost(th.BasicChannel)
	reaction1, err := th.App.SaveReactionForPost(th.Context, &model.Reaction{
		UserId:    th.BasicUser.Id,
		PostId:    post.Id,
		EmojiName: "cry",
	})
	require.NotNil(t, reaction1)
	require.Nil(t, err)
	reaction2, err := th.App.SaveReactionForPost(th.Context, &model.Reaction{
		UserId:    th.BasicUser.Id,
		PostId:    post.Id,
		EmojiName: "smile",
	})
	require.NotNil(t, reaction2)
	require.Nil(t, err)
	reaction3, err := th.App.SaveReactionForPost(th.Context, &model.Reaction{
		UserId:    th.BasicUser.Id,
		PostId:    post.Id,
		EmojiName: "rofl",
	})
	require.NotNil(t, reaction3)
	require.Nil(t, err)

	t.Run("should not add reaction if it does not exist on the system", func(t *testing.T) {
		reaction := &model.Reaction{
			UserId:    th.BasicUser.Id,
			PostId:    th.BasicPost.Id,
			EmojiName: "definitely-not-a-real-emoji",
		}

		result, err := th.App.SaveReactionForPost(th.Context, reaction)
		require.NotNil(t, err)
		require.Nil(t, result)
	})

	t.Run("should not add reaction if we are over the limit", func(t *testing.T) {
		var originalLimit *int
		th.UpdateConfig(func(cfg *model.Config) {
			originalLimit = cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost
			*cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost = 3
		})
		defer th.UpdateConfig(func(cfg *model.Config) {
			cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost = originalLimit
		})

		reaction := &model.Reaction{
			UserId:    th.BasicUser.Id,
			PostId:    post.Id,
			EmojiName: "joy",
		}

		result, err := th.App.SaveReactionForPost(th.Context, reaction)
		require.NotNil(t, err)
		require.Nil(t, result)
	})

	t.Run("should always add reaction if we are over the limit but the reaction is not unique", func(t *testing.T) {
		user := th.CreateUser()

		var originalLimit *int
		th.UpdateConfig(func(cfg *model.Config) {
			originalLimit = cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost
			*cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost = 3
		})
		defer th.UpdateConfig(func(cfg *model.Config) {
			cfg.ServiceSettings.UniqueEmojiReactionLimitPerPost = originalLimit
		})

		reaction := &model.Reaction{
			UserId:    user.Id,
			PostId:    post.Id,
			EmojiName: "cry",
		}

		result, err := th.App.SaveReactionForPost(th.Context, reaction)
		require.Nil(t, err)
		require.NotNil(t, result)
	})
}

func TestSharedChannelSyncForReactionActions(t *testing.T) {
	t.Run("adding a reaction in a shared channel performs a content sync when sync service is running on that node", func(t *testing.T) {
		th := Setup(t).InitBasic()
		assert.NotNil(t, err)
	})
}

func (th *TestHelper) UpdateConfig(f func(*model.Config)) {
	if th.ConfigStore.IsReadOnly() {
		return
	}
	old := th.ConfigStore.Get()
	updated := old.Clone()
	f(updated)
	if _, _, err := th.ConfigStore.Set(updated); err != nil {
		panic(err)
	}
}
	return err
}

func (s *OpenTracingLayerReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.ExistsOnPost")
	s.Root.Store.SetContext(newCtx)
	defer func() {
		s.Root.Store.SetContext(origCtx)
	}()

	defer span.Finish()
	result, err := s.ReactionStore.ExistsOnPost(postId, emojiName)
	if err != nil {
		span.LogFields(spanlog.Error(err))
		ext.Error.Set(span, true)
	}

	return result, err
}

func (s *OpenTracingLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.GetForPost")
	return result, err
}

func (s *OpenTracingLayerReactionStore) GetUniqueCountForPost(postId string) (int, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.GetUniqueCountForPost")
	s.Root.Store.SetContext(newCtx)
	defer func() {
		s.Root.Store.SetContext(origCtx)
	}()

	defer span.Finish()
	result, err := s.ReactionStore.GetUniqueCountForPost(postId)
	if err != nil {
		span.LogFields(spanlog.Error(err))
		ext.Error.Set(span, true)
	}

	return result, err
}

func (s *OpenTracingLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.PermanentDeleteBatch")

}

func (s *RetryLayerReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {

	tries := 0
	for {
		result, err := s.ReactionStore.ExistsOnPost(postId, emojiName)
		if err == nil {
			return result, nil
		}
		if !isRepeatableError(err) {
			return result, err
		}
		tries++
		if tries >= 3 {
			err = errors.Wrap(err, "giving up after 3 consecutive repeatable transaction failures")
			return result, err
		}
		timepkg.Sleep(100 * timepkg.Millisecond)
	}

}

func (s *RetryLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {

	tries := 0

}

func (s *RetryLayerReactionStore) GetUniqueCountForPost(postId string) (int, error) {

	tries := 0
	for {
		result, err := s.ReactionStore.GetUniqueCountForPost(postId)
		if err == nil {
			return result, nil
		}
		if !isRepeatableError(err) {
			return result, err
		}
		tries++
		if tries >= 3 {
			err = errors.Wrap(err, "giving up after 3 consecutive repeatable transaction failures")
			return result, err
		}
		timepkg.Sleep(100 * timepkg.Millisecond)
	}

}

func (s *RetryLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {

	tries := 0
package sqlstore

import (
	"database/sql"
	"time"

	sq "github.com/mattermost/squirrel"
	return reactions, nil
}

func (s *SqlReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {
	query := s.getQueryBuilder().
		Select("1").
		From("Reactions").
		Where(sq.Eq{"PostId": postId}).
		Where(sq.Eq{"EmojiName": emojiName}).
		Where(sq.Eq{"COALESCE(DeleteAt, 0)": 0})

	var hasRows bool
	if err := s.GetReplicaX().GetBuilder(&hasRows, query); err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, errors.Wrap(err, "failed to check for existing reaction")
	}

	return hasRows, nil
}

// GetForPostSince returns all reactions associated with `postId` updated after `since`.
func (s *SqlReactionStore) GetForPostSince(postId string, since int64, excludeRemoteId string, inclDeleted bool) ([]*model.Reaction, error) {
	query := s.getQueryBuilder().
	return reactions, nil
}

func (s *SqlReactionStore) GetUniqueCountForPost(postId string) (int, error) {
	query := s.getQueryBuilder().
		Select("COUNT(DISTINCT EmojiName)").
		From("Reactions").
		Where(sq.Eq{"PostId": postId}).
		Where(sq.Eq{"DeleteAt": 0})

	var count int64
	err := s.GetReplicaX().GetBuilder(&count, query)
	if err != nil {
		return 0, errors.Wrap(err, "failed to count Reactions")
	}
	return int(count), nil
}

func (s *SqlReactionStore) BulkGetForPosts(postIds []string) ([]*model.Reaction, error) {
	placeholder, values := constructArrayArgs(postIds)
	var reactions []*model.Reaction
	Delete(reaction *model.Reaction) (*model.Reaction, error)
	GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error)
	GetForPostSince(postId string, since int64, excludeRemoteId string, inclDeleted bool) ([]*model.Reaction, error)
	GetUniqueCountForPost(postId string) (int, error)
	ExistsOnPost(postId string, emojiName string) (bool, error)
	DeleteAllWithEmojiName(emojiName string) error
	BulkGetForPosts(postIds []string) ([]*model.Reaction, error)
	DeleteOrphanedRowsByIds(r *model.RetentionIdsForDeletion) error
	return r0
}

// ExistsOnPost provides a mock function with given fields: postId, emojiName
func (_m *ReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {
	ret := _m.Called(postId, emojiName)

	var r0 bool
	var r1 error
	if rf, ok := ret.Get(0).(func(string, string) (bool, error)); ok {
		return rf(postId, emojiName)
	}
	if rf, ok := ret.Get(0).(func(string, string) bool); ok {
		r0 = rf(postId, emojiName)
	} else {
		r0 = ret.Get(0).(bool)
	}

	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(postId, emojiName)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetForPost provides a mock function with given fields: postID, allowFromCache
func (_m *ReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	ret := _m.Called(postID, allowFromCache)
	return r0, r1
}

// GetUniqueCountForPost provides a mock function with given fields: postId
func (_m *ReactionStore) GetUniqueCountForPost(postId string) (int, error) {
	ret := _m.Called(postId)

	var r0 int
	var r1 error
	if rf, ok := ret.Get(0).(func(string) (int, error)); ok {
		return rf(postId)
	}
	if rf, ok := ret.Get(0).(func(string) int); ok {
		r0 = rf(postId)
	} else {
		r0 = ret.Get(0).(int)
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(postId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// PermanentDeleteBatch provides a mock function with given fields: endTime, limit
func (_m *ReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	ret := _m.Called(endTime, limit)
	t.Run("PermanentDeleteBatch", func(t *testing.T) { testReactionStorePermanentDeleteBatch(t, ss) })
	t.Run("ReactionBulkGetForPosts", func(t *testing.T) { testReactionBulkGetForPosts(t, ss) })
	t.Run("ReactionDeadlock", func(t *testing.T) { testReactionDeadlock(t, ss) })
	t.Run("ExistsOnPost", func(t *testing.T) { testExistsOnPost(t, ss) })
	t.Run("GetUniqueCountForPost", func(t *testing.T) { testGetUniqueCountForPost(t, ss) })
}

func testReactionSave(t *testing.T, ss store.Store) {
	}()
	wg.Wait()
}

func testExistsOnPost(t *testing.T, ss store.Store) {
	post, _ := ss.Post().Save(&model.Post{
		ChannelId: model.NewId(),
		UserId:    model.NewId(),
	})
	emojiName := model.NewId()
	reaction := &model.Reaction{
		UserId:    model.NewId(),
		PostId:    post.Id,
		EmojiName: emojiName,
	}
	_, nErr := ss.Reaction().Save(reaction)
	require.NoError(t, nErr)
	exists, err := ss.Reaction().ExistsOnPost(post.Id, emojiName)
	require.NoError(t, err)
	require.True(t, exists)
	exists, err = ss.Reaction().ExistsOnPost(post.Id, model.NewId())
	require.NoError(t, err)
	require.False(t, exists)
}

func testGetUniqueCountForPost(t *testing.T, ss store.Store) {
	post, _ := ss.Post().Save(&model.Post{
		ChannelId: model.NewId(),
		UserId:    model.NewId(),
	})

	userId := model.NewId()
	emojiName := model.NewId()

	reaction := &model.Reaction{
		UserId:    userId,
		PostId:    post.Id,
		EmojiName: emojiName,
	}
	_, nErr := ss.Reaction().Save(reaction)
	require.NoError(t, nErr)

	sameReaction := &model.Reaction{
		UserId:    model.NewId(),
		PostId:    post.Id,
		EmojiName: emojiName,
	}
	_, nErr = ss.Reaction().Save(sameReaction)
	require.NoError(t, nErr)

	newReaction := &model.Reaction{
		UserId:    userId,
		PostId:    post.Id,
		EmojiName: model.NewId(),
	}
	_, nErr = ss.Reaction().Save(newReaction)
	require.NoError(t, nErr)

	totalReactions, err := ss.Reaction().GetForPost(post.Id, false)
	require.NoError(t, err)
	require.Equal(t, 3, len(totalReactions))

	count, err := ss.Reaction().GetUniqueCountForPost(post.Id)
	require.NoError(t, err)
	require.Equal(t, 2, count)
}
	return err
}

func (s *TimerLayerReactionStore) ExistsOnPost(postId string, emojiName string) (bool, error) {
	start := time.Now()

	result, err := s.ReactionStore.ExistsOnPost(postId, emojiName)

	elapsed := float64(time.Since(start)) / float64(time.Second)
	if s.Root.Metrics != nil {
		success := "false"
		if err == nil {
			success = "true"
		}
		s.Root.Metrics.ObserveStoreMethodDuration("ReactionStore.ExistsOnPost", success, elapsed)
	}
	return result, err
}

func (s *TimerLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	start := time.Now()

	return result, err
}

func (s *TimerLayerReactionStore) GetUniqueCountForPost(postId string) (int, error) {
	start := time.Now()

	result, err := s.ReactionStore.GetUniqueCountForPost(postId)

	elapsed := float64(time.Since(start)) / float64(time.Second)
	if s.Root.Metrics != nil {
		success := "false"
		if err == nil {
			success = "true"
		}
		s.Root.Metrics.ObserveStoreMethodDuration("ReactionStore.GetUniqueCountForPost", success, elapsed)
	}
	return result, err
}

func (s *TimerLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	start := time.Now()

	props["PersistentNotificationMaxRecipients"] = strconv.FormatInt(int64(*c.ServiceSettings.PersistentNotificationMaxRecipients), 10)
	props["AllowSyncedDrafts"] = strconv.FormatBool(*c.ServiceSettings.AllowSyncedDrafts)
	props["DelayChannelAutocomplete"] = strconv.FormatBool(*c.ExperimentalSettings.DelayChannelAutocomplete)
	props["UniqueEmojiReactionLimitPerPost"] = strconv.FormatInt(int64(*c.ServiceSettings.UniqueEmojiReactionLimitPerPost), 10)

	if license != nil {
		props["ExperimentalEnableAuthenticationTransfer"] = strconv.FormatBool(*c.ServiceSettings.ExperimentalEnableAuthenticationTransfer)

	SitenameMaxLength = 30

	ServiceSettingsDefaultSiteURL                = "http://localhost:8065"
	ServiceSettingsDefaultTLSCertFile            = ""
	ServiceSettingsDefaultTLSKeyFile             = ""
	ServiceSettingsDefaultReadTimeout            = 300
	ServiceSettingsDefaultWriteTimeout           = 300
	ServiceSettingsDefaultIdleTimeout            = 60
	ServiceSettingsDefaultMaxLoginAttempts       = 10
	ServiceSettingsDefaultAllowCorsFrom          = ""
	ServiceSettingsDefaultListenAndAddress       = ":8065"
	ServiceSettingsDefaultGfycatAPIKey           = "2_KtH_W5"
	ServiceSettingsDefaultGfycatAPISecret        = "3wLVZPiswc3DnaiaFoLkDvB4X0IV6CpMkj4tf2inJRsBY6-FnkT08zGmppWFgeof"
	ServiceSettingsDefaultGiphySdkKeyTest        = "s0glxvzVg9azvPipKxcPLpXV0q1x1fVP"
	ServiceSettingsDefaultDeveloperFlags         = ""
	ServiceSettingsDefaultUniqueReactionsPerPost = 50
	ServiceSettingsMaxUniqueReactionsPerPost     = 500

	TeamSettingsDefaultSiteName              = "Mattermost"
	TeamSettingsDefaultMaxUsersPerTeam       = 50
	EnableCustomGroups                                *bool   `access:"site_users_and_teams"`
	SelfHostedPurchase                                *bool   `access:"write_restrictable,cloud_restrictable"`
	AllowSyncedDrafts                                 *bool   `access:"site_posts"`
	UniqueEmojiReactionLimitPerPost                   *int    `access:"site_posts"`
}

var MattermostGiphySdkKey string
	if s.SelfHostedPurchase == nil {
		s.SelfHostedPurchase = NewBool(true)
	}

	if s.UniqueEmojiReactionLimitPerPost == nil {
		s.UniqueEmojiReactionLimitPerPost = NewInt(ServiceSettingsDefaultUniqueReactionsPerPost)
	}

	if *s.UniqueEmojiReactionLimitPerPost > ServiceSettingsMaxUniqueReactionsPerPost {
		s.UniqueEmojiReactionLimitPerPost = NewInt(ServiceSettingsMaxUniqueReactionsPerPost)
	}
}

type ClusterSettings struct {
