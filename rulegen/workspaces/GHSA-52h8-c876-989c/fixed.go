package main


import (
	"context"
	"fmt"
	"github.com/segmentfault/pacman/log"
	"time"

	"github.com/answerdev/answer/internal/base/constant"
	"github.com/answerdev/answer/pkg/converter"

	"github.com/answerdev/answer/internal/base/pager"
	"github.com/answerdev/answer/internal/service/rank"
	"github.com/answerdev/answer/pkg/obj"

	"xorm.io/builder"

	"github.com/answerdev/answer/internal/base/data"
	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/entity"
	"github.com/answerdev/answer/internal/schema"
	"github.com/answerdev/answer/internal/service"
	"github.com/answerdev/answer/internal/service/activity_common"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/xorm"
)
// VoteRepo activity repository
type VoteRepo struct {
	data                     *data.Data
	activityRepo             activity_common.ActivityRepo
	userRankRepo             rank.UserRankRepo
	notificationQueueService notice_queue.NotificationQueueService
}

// NewVoteRepo new repository
func NewVoteRepo(
	data *data.Data,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
	notificationQueueService notice_queue.NotificationQueueService,
) service.VoteRepo {
	return &VoteRepo{
		data:                     data,
		activityRepo:             activityRepo,
		userRankRepo:             userRankRepo,
		notificationQueueService: notificationQueueService,
	}
}

func (vr *VoteRepo) Vote(ctx context.Context, op *schema.VoteOperationInfo) (err error) {
	noNeedToVote, err := vr.votePreCheck(ctx, op)
	if err != nil {
		return err
	}
	if noNeedToVote {
		return nil
	}

	sendInboxNotification := false
	maxDailyRank, err := vr.userRankRepo.GetMaxDailyRank(ctx)
	if err != nil {
		return err
	}
	var userIDs []string
	for _, activity := range op.Activities {
		userIDs = append(userIDs, activity.ActivityUserID)
	}

	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)

		userInfoMapping, err := vr.acquireUserInfo(session, userIDs)
		if err != nil {
			return nil, err
		}

		err = vr.setActivityRankToZeroIfUserReachLimit(ctx, session, op, maxDailyRank)
		if err != nil {
			return nil, err
		}

		sendInboxNotification, err = vr.saveActivitiesAvailable(session, op)
		if err != nil {
			return nil, err
		}

		err = vr.changeUserRank(ctx, session, op, userInfoMapping)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	for _, activity := range op.Activities {
		if activity.Rank == 0 {
			continue
		}
		vr.sendAchievementNotification(ctx, activity.ActivityUserID, op.ObjectCreatorUserID, op.ObjectID)
	}
	if sendInboxNotification {
		vr.sendVoteInboxNotification(ctx, op.OperatingUserID, op.ObjectCreatorUserID, op.ObjectID, op.VoteUp)
	}
	return nil
}

func (vr *VoteRepo) CancelVote(ctx context.Context, op *schema.VoteOperationInfo) (err error) {
	// Pre-Check
	// 1. check if the activity exist
	// 2. check if the activity is not cancelled
	// 3. if all activities are cancelled, return directly
	activities, err := vr.getExistActivity(ctx, op)
	if err != nil {
		return err
	}
	var userIDs []string
	for _, activity := range activities {
		if activity.Cancelled == entity.ActivityCancelled {
			continue
		}
		userIDs = append(userIDs, activity.UserID)
	}
	if len(userIDs) == 0 {
		return nil
	}

	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)

		userInfoMapping, err := vr.acquireUserInfo(session, userIDs)
		if err != nil {
			return nil, err
		}

		err = vr.cancelActivities(session, activities)
		if err != nil {
			return nil, err
		}

		err = vr.rollbackUserRank(ctx, session, activities, userInfoMapping)
		if err != nil {
			return nil, err
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	for _, activity := range activities {
		if activity.Rank == 0 {
			continue
		}
		vr.sendAchievementNotification(ctx, activity.UserID, op.ObjectCreatorUserID, op.ObjectID)
	}
	return nil
}

func (vr *VoteRepo) GetAndSaveVoteResult(ctx context.Context, objectID, objectType string) (
	up, down int64, err error) {
	up = vr.countVoteUp(ctx, objectID, objectType)
	down = vr.countVoteDown(ctx, objectID, objectType)
	err = vr.updateVotes(ctx, objectID, objectType, int(up-down))
	return
}

func (vr *VoteRepo) ListUserVotes(ctx context.Context, userID string,
	page int, pageSize int, activityTypes []int) (voteList []*entity.Activity, total int64, err error) {
	session := vr.data.DB.Context(ctx)
	cond := builder.
		And(
			builder.Eq{"user_id": userID},
			builder.Eq{"cancelled": 0},
			builder.In("activity_type", activityTypes),
		)

	session.Where(cond).Desc("updated_at")

	total, err = pager.Help(page, pageSize, &voteList, &entity.Activity{}, session)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return
}

func (vr *VoteRepo) votePreCheck(ctx context.Context, op *schema.VoteOperationInfo) (noNeedToVote bool, err error) {
	activities, err := vr.getExistActivity(ctx, op)
	if err != nil {
		return false, err
	}
	done := 0
	for _, activity := range activities {
		if activity.Cancelled == entity.ActivityAvailable {
			done++
		}
	}
	return done == len(op.Activities), nil
}

func (vr *VoteRepo) acquireUserInfo(session *xorm.Session, userIDs []string) (map[string]*entity.User, error) {
	us := make([]*entity.User, 0)
	err := session.In("id", userIDs).ForUpdate().Find(&us)
	if err != nil {
		log.Error(err)
		return nil, err
	}

	users := make(map[string]*entity.User, 0)
	for _, u := range us {
		users[u.ID] = u
	}
	return users, nil
}

func (vr *VoteRepo) setActivityRankToZeroIfUserReachLimit(ctx context.Context, session *xorm.Session,
	op *schema.VoteOperationInfo, maxDailyRank int) (err error) {
	// check if user reach daily rank limit
	for _, activity := range op.Activities {
		reach, err := vr.userRankRepo.CheckReachLimit(ctx, session, activity.ActivityUserID, maxDailyRank)
		if err != nil {
			log.Error(err)
			return err
		}
		if reach {
			activity.Rank = 0
		}
	}
	return nil
}

func (vr *VoteRepo) changeUserRank(ctx context.Context, session *xorm.Session,
	op *schema.VoteOperationInfo,
	userInfoMapping map[string]*entity.User) (err error) {
	for _, activity := range op.Activities {
		if activity.Rank == 0 {
			continue
		}
		user := userInfoMapping[activity.ActivityUserID]
		if user == nil {
			continue
		}
		if err = vr.userRankRepo.ChangeUserRank(ctx, session,
			activity.ActivityUserID, user.Rank, activity.Rank); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (vr *VoteRepo) rollbackUserRank(ctx context.Context, session *xorm.Session,
	activities []*entity.Activity,
	userInfoMapping map[string]*entity.User) (err error) {
	for _, activity := range activities {
		if activity.Rank == 0 {
			continue
		}
		user := userInfoMapping[activity.UserID]
		if user == nil {
			continue
		}
		if err = vr.userRankRepo.ChangeUserRank(ctx, session,
			activity.UserID, user.Rank, -activity.Rank); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

// saveActivitiesAvailable save activities
// If activity not exist it will be created or else will be updated
// If this activity is already exist, set activity rank to 0
// So after this function, the activity rank will be correct for update user rank
func (vr *VoteRepo) saveActivitiesAvailable(session *xorm.Session, op *schema.VoteOperationInfo) (newAct bool, err error) {
	for _, activity := range op.Activities {
		existsActivity := &entity.Activity{}
		exist, err := session.
			Where(builder.Eq{"object_id": op.ObjectID}).
			And(builder.Eq{"user_id": activity.ActivityUserID}).
			And(builder.Eq{"trigger_user_id": activity.TriggerUserID}).
			And(builder.Eq{"activity_type": activity.ActivityType}).
			Get(existsActivity)
		if err != nil {
			return false, err
		}
		if exist && existsActivity.Cancelled == entity.ActivityAvailable {
			activity.Rank = 0
			continue
		}
		if exist {
			if _, err = session.Where("id = ?", existsActivity.ID).Cols("`cancelled`").
				Update(&entity.Activity{Cancelled: entity.ActivityAvailable}); err != nil {
				return false, err
			}
		} else {
			insertActivity := entity.Activity{
				ObjectID:         op.ObjectID,
				OriginalObjectID: op.ObjectID,
				UserID:           activity.ActivityUserID,
				TriggerUserID:    converter.StringToInt64(activity.TriggerUserID),
				ActivityType:     activity.ActivityType,
				Rank:             activity.Rank,
				HasRank:          activity.HasRank(),
				Cancelled:        entity.ActivityAvailable,
			}
			_, err = session.Insert(&insertActivity)
			if err != nil {
				return false, err
			}
			newAct = true
		}
	}
	return newAct, nil
}

// cancelActivities cancel activities
// If this activity is already cancelled, set activity rank to 0
// So after this function, the activity rank will be correct for update user rank
func (vr *VoteRepo) cancelActivities(session *xorm.Session, activities []*entity.Activity) (err error) {
	for _, activity := range activities {
		t := &entity.Activity{}
		exist, err := session.ID(activity.ID).Get(t)
		if err != nil {
			log.Error(err)
			return err
		}
		if !exist {
			log.Error(fmt.Errorf("%s activity not exist", activity.ID))
			return fmt.Errorf("%s activity not exist", activity.ID)
		}
		//  If this activity is already cancelled, set activity rank to 0
		if t.Cancelled == entity.ActivityCancelled {
			activity.Rank = 0
		}
		if _, err = session.ID(activity.ID).Cols("cancelled", "cancelled_at").
			Update(&entity.Activity{
				Cancelled:   entity.ActivityCancelled,
				CancelledAt: time.Now(),
			}); err != nil {
			log.Error(err)
			return err
		}
	}
	return nil
}

func (vr *VoteRepo) getExistActivity(ctx context.Context, op *schema.VoteOperationInfo) ([]*entity.Activity, error) {
	var activities []*entity.Activity
	for _, action := range op.Activities {
		t := &entity.Activity{}
		exist, err := vr.data.DB.Context(ctx).
			Where(builder.Eq{"user_id": action.ActivityUserID}).
			And(builder.Eq{"trigger_user_id": action.TriggerUserID}).
			And(builder.Eq{"activity_type": action.ActivityType}).
			And(builder.Eq{"object_id": op.ObjectID}).
			Get(t)
		if err != nil {
			return nil, errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
		}
		if exist {
			activities = append(activities, t)
		}
	}
	return activities, nil
}

func (vr *VoteRepo) countVoteUp(ctx context.Context, objectID, objectType string) (count int64) {
	count, err := vr.countVote(ctx, objectID, objectType, constant.ActVoteUp)
	if err != nil {
		log.Errorf("get vote up count error: %v", err)
	}
	return count
}

func (vr *VoteRepo) countVoteDown(ctx context.Context, objectID, objectType string) (count int64) {
	count, err := vr.countVote(ctx, objectID, objectType, constant.ActVoteDown)
	if err != nil {
		log.Errorf("get vote down count error: %v", err)
	}
	return count
}

func (vr *VoteRepo) countVote(ctx context.Context, objectID, objectType, action string) (count int64, err error) {
	activity := &entity.Activity{}
	activityType, _ := vr.activityRepo.GetActivityTypeByObjectType(ctx, objectType, action)
	count, err = vr.data.DB.Context(ctx).Where(builder.Eq{"object_id": objectID}).
		And(builder.Eq{"activity_type": activityType}).
		And(builder.Eq{"cancelled": 0}).
		Count(activity)
	if err != nil {
		err = errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	return count, err
}

func (vr *VoteRepo) updateVotes(ctx context.Context, objectID, objectType string, voteCount int) (err error) {
	session := vr.data.DB.Context(ctx)
	switch objectType {
	case constant.QuestionObjectType:
		_, err = session.ID(objectID).Cols("vote_count").Update(&entity.Question{VoteCount: voteCount})
	case constant.AnswerObjectType:
		_, err = session.ID(objectID).Cols("vote_count").Update(&entity.Answer{VoteCount: voteCount})
	case constant.CommentObjectType:
		_, err = session.ID(objectID).Cols("vote_count").Update(&entity.Comment{VoteCount: voteCount})
	}
	if err != nil {
		log.Error(err)
	}
	return
}

func (vr *VoteRepo) sendAchievementNotification(ctx context.Context, activityUserID, objectUserID, objectID string) {
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		return

import (
	"context"
	"github.com/answerdev/answer/internal/service/activity_common"
	"strings"

	"github.com/answerdev/answer/internal/base/constant"
	"github.com/answerdev/answer/internal/base/handler"
	"github.com/answerdev/answer/internal/service/config"
	"github.com/answerdev/answer/internal/service/object_info"
	"github.com/answerdev/answer/pkg/htmltext"
	"github.com/segmentfault/pacman/log"

	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/schema"
	answercommon "github.com/answerdev/answer/internal/service/answer_common"
	questioncommon "github.com/answerdev/answer/internal/service/question_common"
	"github.com/segmentfault/pacman/errors"
)

// VoteRepo activity repository
type VoteRepo interface {
	Vote(ctx context.Context, op *schema.VoteOperationInfo) (err error)
	CancelVote(ctx context.Context, op *schema.VoteOperationInfo) (err error)
	GetAndSaveVoteResult(ctx context.Context, objectID, objectType string) (up, down int64, err error)
	ListUserVotes(ctx context.Context, userID string, page int, pageSize int, activityTypes []int) (
		voteList []*entity.Activity, total int64, err error)
}

// VoteService user service
type VoteService struct {
	voteRepo          VoteRepo
	configService     *config.ConfigService
	questionRepo      questioncommon.QuestionRepo
	answerRepo        answercommon.AnswerRepo
	commentCommonRepo comment_common.CommentCommonRepo
	objectService     *object_info.ObjService
	activityRepo      activity_common.ActivityRepo
}

func NewVoteService(
	voteRepo VoteRepo,
	configService *config.ConfigService,
	questionRepo questioncommon.QuestionRepo,
	answerRepo answercommon.AnswerRepo,
	objectService *object_info.ObjService,
) *VoteService {
	return &VoteService{
		voteRepo:          voteRepo,
		configService:     configService,
		questionRepo:      questionRepo,
		answerRepo:        answerRepo,
}

// VoteUp vote up
func (vs *VoteService) VoteUp(ctx context.Context, req *schema.VoteReq) (resp *schema.VoteResp, err error) {
	objectInfo, err := vs.objectService.GetInfo(ctx, req.ObjectID)
	if err != nil {
		return nil, err
	}
	// make object id must be decoded
	objectInfo.ObjectID = req.ObjectID

	// check user is voting self or not
	if objectInfo.ObjectCreatorUserID == req.UserID {
		return nil, errors.BadRequest(reason.DisallowVoteYourSelf)
	}

	voteUpOperationInfo := vs.createVoteOperationInfo(ctx, req.UserID, true, objectInfo)

	// vote operation
	if req.IsCancel {
		err = vs.voteRepo.CancelVote(ctx, voteUpOperationInfo)
	} else {
		// cancel vote down if exist
		voteOperationInfo := vs.createVoteOperationInfo(ctx, req.UserID, false, objectInfo)
		err = vs.voteRepo.CancelVote(ctx, voteOperationInfo)
		if err != nil {
			return nil, err
		}
		err = vs.voteRepo.Vote(ctx, voteUpOperationInfo)
	}

	resp = &schema.VoteResp{}
	resp.UpVotes, resp.DownVotes, err = vs.voteRepo.GetAndSaveVoteResult(ctx, req.ObjectID, objectInfo.ObjectType)
	if err != nil {
		log.Error(err)
	}
	resp.Votes = resp.UpVotes - resp.DownVotes
	if !req.IsCancel {
		resp.VoteStatus = constant.ActVoteUp
	}
	return resp, nil
}

// VoteDown vote down
func (vs *VoteService) VoteDown(ctx context.Context, req *schema.VoteReq) (resp *schema.VoteResp, err error) {
	objectInfo, err := vs.objectService.GetInfo(ctx, req.ObjectID)
	if err != nil {
		return nil, err
	}
	// make object id must be decoded
	objectInfo.ObjectID = req.ObjectID

	// check user is voting self or not
	if objectInfo.ObjectCreatorUserID == req.UserID {
		return nil, errors.BadRequest(reason.DisallowVoteYourSelf)
	}

	// vote operation
	voteDownOperationInfo := vs.createVoteOperationInfo(ctx, req.UserID, false, objectInfo)
	if req.IsCancel {
		err = vs.voteRepo.CancelVote(ctx, voteDownOperationInfo)
	} else {
		// cancel vote up if exist
		err = vs.voteRepo.CancelVote(ctx, vs.createVoteOperationInfo(ctx, req.UserID, true, objectInfo))
		if err != nil {
			return nil, err
		}
		err = vs.voteRepo.Vote(ctx, voteDownOperationInfo)
	}

	resp = &schema.VoteResp{}
	resp.UpVotes, resp.DownVotes, err = vs.voteRepo.GetAndSaveVoteResult(ctx, req.ObjectID, objectInfo.ObjectType)
	if err != nil {
		log.Error(err)
	}
	resp.Votes = resp.UpVotes - resp.DownVotes
	if !req.IsCancel {
		resp.VoteStatus = constant.ActVoteDown
	}
	return resp, nil
}

// ListUserVotes list user's votes
	}
	return pager.NewPageModel(total, votes), err
}

func (vs *VoteService) createVoteOperationInfo(ctx context.Context,
	userID string, voteUp bool, objectInfo *schema.SimpleObjectInfo) *schema.VoteOperationInfo {
	// warp vote operation
	voteOperationInfo := &schema.VoteOperationInfo{
		ObjectID:            objectInfo.ObjectID,
		ObjectType:          objectInfo.ObjectType,
		ObjectCreatorUserID: objectInfo.ObjectCreatorUserID,
		OperatingUserID:     userID,
		VoteUp:              voteUp,
		VoteDown:            !voteUp,
	}
	voteOperationInfo.Activities = vs.getActivities(ctx, voteOperationInfo)
	return voteOperationInfo
}

func (vs *VoteService) getActivities(ctx context.Context, op *schema.VoteOperationInfo) (
	activities []*schema.VoteActivity) {
	activities = make([]*schema.VoteActivity, 0)

	var actions []string
	switch op.ObjectType {
	case constant.QuestionObjectType:
		if op.VoteUp {
			actions = []string{activity_type.QuestionVoteUp, activity_type.QuestionVotedUp}
		} else {
			actions = []string{activity_type.QuestionVoteDown, activity_type.QuestionVotedDown}
		}
	case constant.AnswerObjectType:
		if op.VoteUp {
			actions = []string{activity_type.AnswerVoteUp, activity_type.AnswerVotedUp}
		} else {
			actions = []string{activity_type.AnswerVoteDown, activity_type.AnswerVotedDown}
		}
	case constant.CommentObjectType:
		actions = []string{activity_type.CommentVoteUp}
	}

	for _, action := range actions {
		t := &schema.VoteActivity{}
		cfg, err := vs.configService.GetConfigByKey(ctx, action)
		if err != nil {
			log.Warnf("get config by key error: %v", err)
			continue
		}
		t.ActivityType, t.Rank = cfg.ID, cfg.GetIntValue()

		if strings.Contains(action, "voted") {
			t.ActivityUserID = op.ObjectCreatorUserID
			t.TriggerUserID = op.OperatingUserID
		} else {
			t.ActivityUserID = op.OperatingUserID
			t.TriggerUserID = "0"
		}
		activities = append(activities, t)
	}
	return activities
}
	"github.com/answerdev/answer/internal/service/rank"
	"github.com/answerdev/answer/pkg/converter"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/xorm"
)

	}
}

// AcceptAnswer accept other answer
func (ar *AnswerActivityRepo) AcceptAnswer(ctx context.Context,
	answerObjID, questionObjID, questionUserID, answerUserID string, isSelf bool,
	}
	return err
}
