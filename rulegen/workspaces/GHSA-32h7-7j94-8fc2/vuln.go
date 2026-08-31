package main

	})

	t.Run("save-second-reaction", func(t *testing.T) {
		reaction.EmojiName = "sad"

		rr, _, err := client.SaveReaction(context.Background(), reaction)
		require.NoError(t, err)
	r2 := &model.Reaction{
		UserId:    userId,
		PostId:    postId,
		EmojiName: "smile-",
	}

	r3 := &model.Reaction{
	r4 := &model.Reaction{
		UserId:    user2Id,
		PostId:    postId,
		EmojiName: "smile_",
	}

	// Check the appropriate permissions are enforced.
	reactionObject := model.Reaction{
		UserId:    th.BasicUser.Id,
		PostId:    post.Id,
		EmojiName: "emoji",
		CreateAt:  model.GetMillis(),
	}
	reactionObjectDeleted := model.Reaction{
		UserId:    th.BasicUser2.Id,
		PostId:    post.Id,
		EmojiName: "emoji",
		CreateAt:  model.GetMillis(),
	}

	th.App.SaveReactionForPost(th.Context, &reactionObject)
	th.App.SaveReactionForPost(th.Context, &reactionObjectDeleted)
	reactionsOfPost, err := th.App.BuildPostReactions(th.Context, post.Id)
	require.Nil(t, err)

		return nil, err
	}

	channel, err := a.GetChannel(c, post.ChannelId)
	if err != nil {
		return nil, err
	"github.com/mattermost/mattermost/server/v8/channels/testlib"
)

func TestSharedChannelSyncForReactionActions(t *testing.T) {
	t.Run("adding a reaction in a shared channel performs a content sync when sync service is running on that node", func(t *testing.T) {
		th := Setup(t).InitBasic()
		assert.NotNil(t, err)
	})
}
	return err
}

func (s *OpenTracingLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.GetForPost")
	return result, err
}

func (s *OpenTracingLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	origCtx := s.Root.Store.Context()
	span, newCtx := tracing.StartSpanWithParentByContext(s.Root.Store.Context(), "ReactionStore.PermanentDeleteBatch")

}

func (s *RetryLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {

	tries := 0

}

func (s *RetryLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {

	tries := 0
package sqlstore

import (
	"time"

	sq "github.com/mattermost/squirrel"
	return reactions, nil
}

// GetForPostSince returns all reactions associated with `postId` updated after `since`.
func (s *SqlReactionStore) GetForPostSince(postId string, since int64, excludeRemoteId string, inclDeleted bool) ([]*model.Reaction, error) {
	query := s.getQueryBuilder().
	return reactions, nil
}

func (s *SqlReactionStore) BulkGetForPosts(postIds []string) ([]*model.Reaction, error) {
	placeholder, values := constructArrayArgs(postIds)
	var reactions []*model.Reaction
	Delete(reaction *model.Reaction) (*model.Reaction, error)
	GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error)
	GetForPostSince(postId string, since int64, excludeRemoteId string, inclDeleted bool) ([]*model.Reaction, error)
	DeleteAllWithEmojiName(emojiName string) error
	BulkGetForPosts(postIds []string) ([]*model.Reaction, error)
	DeleteOrphanedRowsByIds(r *model.RetentionIdsForDeletion) error
	return r0
}

// GetForPost provides a mock function with given fields: postID, allowFromCache
func (_m *ReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	ret := _m.Called(postID, allowFromCache)
	return r0, r1
}

// PermanentDeleteBatch provides a mock function with given fields: endTime, limit
func (_m *ReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	ret := _m.Called(endTime, limit)
	t.Run("PermanentDeleteBatch", func(t *testing.T) { testReactionStorePermanentDeleteBatch(t, ss) })
	t.Run("ReactionBulkGetForPosts", func(t *testing.T) { testReactionBulkGetForPosts(t, ss) })
	t.Run("ReactionDeadlock", func(t *testing.T) { testReactionDeadlock(t, ss) })
}

func testReactionSave(t *testing.T, ss store.Store) {
	}()
	wg.Wait()
}
	return err
}

func (s *TimerLayerReactionStore) GetForPost(postID string, allowFromCache bool) ([]*model.Reaction, error) {
	start := time.Now()

	return result, err
}

func (s *TimerLayerReactionStore) PermanentDeleteBatch(endTime int64, limit int64) (int64, error) {
	start := time.Now()

	props["PersistentNotificationMaxRecipients"] = strconv.FormatInt(int64(*c.ServiceSettings.PersistentNotificationMaxRecipients), 10)
	props["AllowSyncedDrafts"] = strconv.FormatBool(*c.ServiceSettings.AllowSyncedDrafts)
	props["DelayChannelAutocomplete"] = strconv.FormatBool(*c.ExperimentalSettings.DelayChannelAutocomplete)

	if license != nil {
		props["ExperimentalEnableAuthenticationTransfer"] = strconv.FormatBool(*c.ServiceSettings.ExperimentalEnableAuthenticationTransfer)

	SitenameMaxLength = 30

	ServiceSettingsDefaultSiteURL          = "http://localhost:8065"
	ServiceSettingsDefaultTLSCertFile      = ""
	ServiceSettingsDefaultTLSKeyFile       = ""
	ServiceSettingsDefaultReadTimeout      = 300
	ServiceSettingsDefaultWriteTimeout     = 300
	ServiceSettingsDefaultIdleTimeout      = 60
	ServiceSettingsDefaultMaxLoginAttempts = 10
	ServiceSettingsDefaultAllowCorsFrom    = ""
	ServiceSettingsDefaultListenAndAddress = ":8065"
	ServiceSettingsDefaultGfycatAPIKey     = "2_KtH_W5"
	ServiceSettingsDefaultGfycatAPISecret  = "3wLVZPiswc3DnaiaFoLkDvB4X0IV6CpMkj4tf2inJRsBY6-FnkT08zGmppWFgeof"
	ServiceSettingsDefaultGiphySdkKeyTest  = "s0glxvzVg9azvPipKxcPLpXV0q1x1fVP"
	ServiceSettingsDefaultDeveloperFlags   = ""

	TeamSettingsDefaultSiteName              = "Mattermost"
	TeamSettingsDefaultMaxUsersPerTeam       = 50
	EnableCustomGroups                                *bool   `access:"site_users_and_teams"`
	SelfHostedPurchase                                *bool   `access:"write_restrictable,cloud_restrictable"`
	AllowSyncedDrafts                                 *bool   `access:"site_posts"`
}

var MattermostGiphySdkKey string
	if s.SelfHostedPurchase == nil {
		s.SelfHostedPurchase = NewBool(true)
	}
}

type ClusterSettings struct {
