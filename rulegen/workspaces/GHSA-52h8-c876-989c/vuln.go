package main


import (
	"context"
	"strings"
	"time"

	"github.com/answerdev/answer/internal/base/constant"
	"github.com/answerdev/answer/pkg/converter"

	"github.com/answerdev/answer/internal/base/pager"
	"github.com/answerdev/answer/internal/service/config"
	"github.com/answerdev/answer/internal/service/rank"
	"github.com/answerdev/answer/pkg/obj"

	"xorm.io/builder"

	"github.com/answerdev/answer/internal/service/activity_common"
	"github.com/answerdev/answer/internal/service/unique"

	"github.com/answerdev/answer/internal/base/data"
	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/entity"
	"github.com/answerdev/answer/internal/schema"
	"github.com/answerdev/answer/internal/service"
	"github.com/segmentfault/pacman/errors"
	"xorm.io/xorm"
)
// VoteRepo activity repository
type VoteRepo struct {
	data                     *data.Data
	uniqueIDRepo             unique.UniqueIDRepo
	configService            *config.ConfigService
	activityRepo             activity_common.ActivityRepo
	userRankRepo             rank.UserRankRepo
	voteCommon               activity_common.VoteRepo
	notificationQueueService notice_queue.NotificationQueueService
}

// NewVoteRepo new repository
func NewVoteRepo(
	data *data.Data,
	uniqueIDRepo unique.UniqueIDRepo,
	configService *config.ConfigService,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
	voteCommon activity_common.VoteRepo,
	notificationQueueService notice_queue.NotificationQueueService,
) service.VoteRepo {
	return &VoteRepo{
		data:                     data,
		uniqueIDRepo:             uniqueIDRepo,
		configService:            configService,
		activityRepo:             activityRepo,
		userRankRepo:             userRankRepo,
		voteCommon:               voteCommon,
		notificationQueueService: notificationQueueService,
	}
}

var LimitUpActions = map[string][]string{
	"question": {"vote_up", "voted_up"},
	"answer":   {"vote_up", "voted_up"},
	"comment":  {"vote_up"},
}

var LimitDownActions = map[string][]string{
	"question": {"vote_down", "voted_down"},
	"answer":   {"vote_down", "voted_down"},
	"comment":  {"vote_down"},
}

func (vr *VoteRepo) vote(ctx context.Context, objectID string, userID, objectUserID string, actions []string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	achievementNotificationUserIDs := make([]string, 0)
	sendInboxNotification := false
	upVote := false
	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		result = nil
		for _, action := range actions {
			var (
				existsActivity entity.Activity
				insertActivity entity.Activity
				has            bool
				triggerUserID,
				activityUserID string
				activityType, deltaRank, hasRank int
			)

			activityUserID, activityType, deltaRank, hasRank, err = vr.CheckRank(ctx, objectID, objectUserID, userID, action)
			if err != nil {
				return
			}

			triggerUserID = userID
			if userID == activityUserID {
				triggerUserID = "0"
			}

			// check is voted up
			has, _ = session.
				Where(builder.Eq{"object_id": objectID}).
				And(builder.Eq{"user_id": activityUserID}).
				And(builder.Eq{"trigger_user_id": triggerUserID}).
				And(builder.Eq{"activity_type": activityType}).
				Get(&existsActivity)

			// is is voted,return
			if has && existsActivity.Cancelled == entity.ActivityAvailable {
				return
			}

			insertActivity = entity.Activity{
				ObjectID:         objectID,
				OriginalObjectID: objectID,
				UserID:           activityUserID,
				TriggerUserID:    converter.StringToInt64(triggerUserID),
				ActivityType:     activityType,
				Rank:             deltaRank,
				HasRank:          hasRank,
				Cancelled:        entity.ActivityAvailable,
			}

			// trigger user rank and send notification
			if hasRank != 0 {
				var isReachStandard bool
				isReachStandard, err = vr.userRankRepo.TriggerUserRank(ctx, session, activityUserID, deltaRank, activityType)
				if err != nil {
					return nil, err
				}
				if isReachStandard {
					insertActivity.Rank = 0
				}
				achievementNotificationUserIDs = append(achievementNotificationUserIDs, activityUserID)
			}

			if has {
				if _, err = session.Where("id = ?", existsActivity.ID).Cols("`cancelled`").
					Update(&entity.Activity{
						Cancelled: entity.ActivityAvailable,
					}); err != nil {
					return
				}
			} else {
				_, err = session.Insert(&insertActivity)
				if err != nil {
					return nil, err
				}
				sendInboxNotification = true
			}

			// update votes
			if action == constant.ActVoteDown || action == constant.ActVoteUp {
				votes := 1
				if action == constant.ActVoteDown {
					upVote = false
					votes = -1
				} else {
					upVote = true
				}
				err = vr.updateVotes(ctx, session, objectID, votes)
				if err != nil {
					return
				}
			}
		}
		return
	})
	if err != nil {
		return
	}

	resp, err = vr.GetVoteResultByObjectId(ctx, objectID)
	resp.VoteStatus = vr.voteCommon.GetVoteStatus(ctx, objectID, userID)

	for _, activityUserID := range achievementNotificationUserIDs {
		vr.sendNotification(ctx, activityUserID, objectUserID, objectID)
	}
	if sendInboxNotification {
		vr.sendVoteInboxNotification(ctx, userID, objectUserID, objectID, upVote)
	}
	return
}

func (vr *VoteRepo) voteCancel(ctx context.Context, objectID string, userID, objectUserID string, actions []string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	notificationUserIDs := make([]string, 0)
	_, err = vr.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		for _, action := range actions {
			var (
				existsActivity entity.Activity
				has            bool
				triggerUserID,
				activityUserID string
				activityType,
				deltaRank, hasRank int
			)
			result = nil

			activityUserID, activityType, deltaRank, hasRank, err = vr.CheckRank(ctx, objectID, objectUserID, userID, action)
			if err != nil {
				return
			}

			triggerUserID = userID
			if userID == activityUserID {
				triggerUserID = "0"
			}

			has, err = session.
				Where(builder.Eq{"user_id": activityUserID}).
				And(builder.Eq{"trigger_user_id": triggerUserID}).
				And(builder.Eq{"activity_type": activityType}).
				And(builder.Eq{"object_id": objectID}).
				Get(&existsActivity)

			if !has {
				return
			}

			if existsActivity.Cancelled == entity.ActivityCancelled {
				return
			}

			if _, err = session.Where("id = ?", existsActivity.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{
					Cancelled:   entity.ActivityCancelled,
					CancelledAt: time.Now(),
				}); err != nil {
				return
			}

			// trigger user rank and send notification
			if hasRank != 0 && existsActivity.Rank != 0 {
				_, err = vr.userRankRepo.TriggerUserRank(ctx, session, activityUserID, -deltaRank, activityType)
				if err != nil {
					return
				}
				notificationUserIDs = append(notificationUserIDs, activityUserID)
			}

			// update votes
			if action == "vote_down" || action == "vote_up" {
				votes := -1
				if action == "vote_down" {
					votes = 1
				}
				err = vr.updateVotes(ctx, session, objectID, votes)
				if err != nil {
					return
				}
			}
		}

		return
	})
	if err != nil {
		return
	}
	resp, err = vr.GetVoteResultByObjectId(ctx, objectID)
	resp.VoteStatus = vr.voteCommon.GetVoteStatus(ctx, objectID, userID)

	for _, activityUserID := range notificationUserIDs {
		vr.sendNotification(ctx, activityUserID, objectUserID, objectID)
	}
	return
}

func (vr *VoteRepo) VoteUp(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}

	actions, ok := LimitUpActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	_, _ = vr.VoteDownCancel(ctx, objectID, userID, objectUserID)
	return vr.vote(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteDown(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitDownActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	_, _ = vr.VoteUpCancel(ctx, objectID, userID, objectUserID)
	return vr.vote(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteUpCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	var objectType string
	resp = &schema.VoteResp{}

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitUpActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	return vr.voteCancel(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) VoteDownCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error) {
	var objectType string
	resp = &schema.VoteResp{}

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		err = errors.BadRequest(reason.ObjectNotFound)
		return
	}
	actions, ok := LimitDownActions[objectType]
	if !ok {
		err = errors.BadRequest(reason.DisallowVote)
		return
	}

	return vr.voteCancel(ctx, objectID, userID, objectUserID, actions)
}

func (vr *VoteRepo) CheckRank(ctx context.Context, objectID, objectUserID, userID string, action string) (activityUserID string, activityType, rank, hasRank int, err error) {
	activityType, rank, hasRank, err = vr.activityRepo.GetActivityTypeByObjID(ctx, objectID, action)

	if err != nil {
		return
	}

	activityUserID = userID
	if strings.Contains(action, "voted") {
		activityUserID = objectUserID
	}

	return activityUserID, activityType, rank, hasRank, nil
}

func (vr *VoteRepo) GetVoteResultByObjectId(ctx context.Context, objectID string) (resp *schema.VoteResp, err error) {
	resp = &schema.VoteResp{}
	for _, action := range []string{"vote_up", "vote_down"} {
		var (
			activity     entity.Activity
			votes        int64
			activityType int
		)

		activityType, _, _, _ = vr.activityRepo.GetActivityTypeByObjID(ctx, objectID, action)

		votes, err = vr.data.DB.Context(ctx).Where(builder.Eq{"object_id": objectID}).
			And(builder.Eq{"activity_type": activityType}).
			And(builder.Eq{"cancelled": 0}).
			Count(&activity)

		if err != nil {
			return
		}

		if action == "vote_up" {
			resp.UpVotes = int(votes)
		} else {
			resp.DownVotes = int(votes)
		}
	}

	resp.Votes = resp.UpVotes - resp.DownVotes

	return resp, nil
}

func (vr *VoteRepo) ListUserVotes(ctx context.Context, userID string,
	page int, pageSize int, activityTypes []int) (voteList []entity.Activity, total int64, err error) {
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

// updateVotes
// if votes < 0 Decr object vote_count,otherwise Incr object vote_count
func (vr *VoteRepo) updateVotes(ctx context.Context, session *xorm.Session, objectID string, votes int) (err error) {
	var (
		objectType string
		e          error
	)

	objectType, err = obj.GetObjectTypeStrByObjectID(objectID)
	switch objectType {
	case "question":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Question{})
	case "answer":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Answer{})
	case "comment":
		_, err = session.Where("id = ?", objectID).Incr("vote_count", votes).Update(&entity.Comment{})
	default:
		e = errors.BadRequest(reason.DisallowVote)
	}

	if e != nil {
		err = e
	} else if err != nil {
		err = errors.BadRequest(reason.DatabaseError).WithError(err).WithStack()
	}

	return
}

// sendNotification send rank triggered notification
func (vr *VoteRepo) sendNotification(ctx context.Context, activityUserID, objectUserID, objectID string) {
	objectType, err := obj.GetObjectTypeStrByObjectID(objectID)
	if err != nil {
		return

import (
	"context"

	"github.com/answerdev/answer/internal/base/constant"
	"github.com/answerdev/answer/internal/base/handler"
	"github.com/answerdev/answer/internal/service/config"
	"github.com/answerdev/answer/internal/service/object_info"
	"github.com/answerdev/answer/pkg/htmltext"
	"github.com/answerdev/answer/pkg/obj"
	"github.com/segmentfault/pacman/log"

	"github.com/answerdev/answer/internal/base/reason"
	"github.com/answerdev/answer/internal/schema"
	answercommon "github.com/answerdev/answer/internal/service/answer_common"
	questioncommon "github.com/answerdev/answer/internal/service/question_common"
	"github.com/answerdev/answer/internal/service/unique"
	"github.com/segmentfault/pacman/errors"
)

// VoteRepo activity repository
type VoteRepo interface {
	VoteUp(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error)
	VoteDown(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error)
	VoteUpCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error)
	VoteDownCancel(ctx context.Context, objectID string, userID, objectUserID string) (resp *schema.VoteResp, err error)
	GetVoteResultByObjectId(ctx context.Context, objectID string) (resp *schema.VoteResp, err error)
	ListUserVotes(ctx context.Context, userID string, page int, pageSize int, activityTypes []int) (
		voteList []entity.Activity, total int64, err error)
}

// VoteService user service
type VoteService struct {
	voteRepo          VoteRepo
	UniqueIDRepo      unique.UniqueIDRepo
	configService     *config.ConfigService
	questionRepo      questioncommon.QuestionRepo
	answerRepo        answercommon.AnswerRepo
	commentCommonRepo comment_common.CommentCommonRepo
	objectService     *object_info.ObjService
}

func NewVoteService(
	VoteRepo VoteRepo,
	uniqueIDRepo unique.UniqueIDRepo,
	configService *config.ConfigService,
	questionRepo questioncommon.QuestionRepo,
	answerRepo answercommon.AnswerRepo,
	objectService *object_info.ObjService,
) *VoteService {
	return &VoteService{
		voteRepo:          VoteRepo,
		UniqueIDRepo:      uniqueIDRepo,
		configService:     configService,
		questionRepo:      questionRepo,
		answerRepo:        answerRepo,
}

// VoteUp vote up
func (vs *VoteService) VoteUp(ctx context.Context, dto *schema.VoteDTO) (voteResp *schema.VoteResp, err error) {
	voteResp = &schema.VoteResp{}

	var objectUserID string

	objectUserID, err = vs.GetObjectUserID(ctx, dto.ObjectID)
	if err != nil {
		return
	}

	// check user is voting self or not
	if objectUserID == dto.UserID {
		err = errors.BadRequest(reason.DisallowVoteYourSelf)
		return
	}

	if dto.IsCancel {
		return vs.voteRepo.VoteUpCancel(ctx, dto.ObjectID, dto.UserID, objectUserID)
	} else {
		return vs.voteRepo.VoteUp(ctx, dto.ObjectID, dto.UserID, objectUserID)
	}
}

// VoteDown vote down
func (vs *VoteService) VoteDown(ctx context.Context, dto *schema.VoteDTO) (voteResp *schema.VoteResp, err error) {
	voteResp = &schema.VoteResp{}

	var objectUserID string

	objectUserID, err = vs.GetObjectUserID(ctx, dto.ObjectID)
	if err != nil {
		return
	}

	// check user is voting self or not
	if objectUserID == dto.UserID {
		err = errors.BadRequest(reason.DisallowVoteYourSelf)
		return
	}

	if dto.IsCancel {
		return vs.voteRepo.VoteDownCancel(ctx, dto.ObjectID, dto.UserID, objectUserID)
	} else {
		return vs.voteRepo.VoteDown(ctx, dto.ObjectID, dto.UserID, objectUserID)
	}
}

func (vs *VoteService) GetObjectUserID(ctx context.Context, objectID string) (userID string, err error) {
	var objectKey string
	objectKey, err = obj.GetObjectTypeStrByObjectID(objectID)

	if err != nil {
		err = nil
		return
	}

	switch objectKey {
	case "question":
		object, has, e := vs.questionRepo.GetQuestion(ctx, objectID)
		if e != nil || !has {
			err = errors.BadRequest(reason.QuestionNotFound).WithError(e).WithStack()
			return
		}
		userID = object.UserID
	case "answer":
		object, has, e := vs.answerRepo.GetAnswer(ctx, objectID)
		if e != nil || !has {
			err = errors.BadRequest(reason.AnswerNotFound).WithError(e).WithStack()
			return
		}
		userID = object.UserID
	case "comment":
		object, has, e := vs.commentCommonRepo.GetComment(ctx, objectID)
		if e != nil || !has {
			err = errors.BadRequest(reason.CommentNotFound).WithError(e).WithStack()
			return
		}
		userID = object.UserID
	default:
		err = errors.BadRequest(reason.DisallowVote).WithError(err).WithStack()
		return
	}

	return
}

// ListUserVotes list user's votes
	}
	return pager.NewPageModel(total, votes), err
}
	"github.com/answerdev/answer/internal/service/rank"
	"github.com/answerdev/answer/pkg/converter"
	"github.com/segmentfault/pacman/errors"
	"github.com/segmentfault/pacman/log"
	"xorm.io/xorm"
)

	}
}

// NewQuestionActivityRepo new repository
func NewQuestionActivityRepo(
	data *data.Data,
	activityRepo activity_common.ActivityRepo,
	userRankRepo rank.UserRankRepo,
) activity.QuestionActivityRepo {
	return &AnswerActivityRepo{
		data:         data,
		activityRepo: activityRepo,
		userRankRepo: userRankRepo,
	}
}

func (ar *AnswerActivityRepo) DeleteQuestion(ctx context.Context, questionID string) (err error) {
	questionInfo := &entity.Question{}
	exist, err := ar.data.DB.Context(ctx).Where("id = ?", questionID).Get(questionInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil
	}

	// get all this object activity
	activityList := make([]*entity.Activity, 0)
	session := ar.data.DB.Context(ctx).Where("has_rank = 1")
	session.Where("cancelled = ?", entity.ActivityAvailable)
	err = session.Find(&activityList, &entity.Activity{ObjectID: questionID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(activityList) == 0 {
		return nil
	}

	log.Infof("questionInfo %s deleted will rollback activity %d", questionID, len(activityList))

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		for _, act := range activityList {
			log.Infof("user %s rollback rank %d", act.UserID, -act.Rank)
			_, e := ar.userRankRepo.TriggerUserRank(
				ctx, session, act.UserID, -act.Rank, act.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}

			if _, e := session.Where("id = ?", act.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{Cancelled: entity.ActivityCancelled, CancelledAt: time.Now()}); e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}

	// get all answers
	answerList := make([]*entity.Answer, 0)
	err = ar.data.DB.Context(ctx).Find(&answerList, &entity.Answer{QuestionID: questionID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	for _, answerInfo := range answerList {
		err = ar.DeleteAnswer(ctx, answerInfo.ID)
		if err != nil {
			log.Error(err)
		}
	}
	return
}

// AcceptAnswer accept other answer
func (ar *AnswerActivityRepo) AcceptAnswer(ctx context.Context,
	answerObjID, questionObjID, questionUserID, answerUserID string, isSelf bool,
	}
	return err
}

func (ar *AnswerActivityRepo) DeleteAnswer(ctx context.Context, answerID string) (err error) {
	answerInfo := &entity.Answer{}
	exist, err := ar.data.DB.Context(ctx).Where("id = ?", answerID).Get(answerInfo)
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if !exist {
		return nil
	}

	// get all this object activity
	activityList := make([]*entity.Activity, 0)
	session := ar.data.DB.Context(ctx).Where("has_rank = 1")
	session.Where("cancelled = ?", entity.ActivityAvailable)
	err = session.Find(&activityList, &entity.Activity{ObjectID: answerID})
	if err != nil {
		return errors.InternalServer(reason.DatabaseError).WithError(err).WithStack()
	}
	if len(activityList) == 0 {
		return nil
	}

	log.Infof("answerInfo %s deleted will rollback activity %d", answerID, len(activityList))

	_, err = ar.data.DB.Transaction(func(session *xorm.Session) (result any, err error) {
		session = session.Context(ctx)
		for _, act := range activityList {
			log.Infof("user %s rollback rank %d", act.UserID, -act.Rank)
			_, e := ar.userRankRepo.TriggerUserRank(
				ctx, session, act.UserID, -act.Rank, act.ActivityType)
			if e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}

			if _, e := session.Where("id = ?", act.ID).Cols("cancelled", "cancelled_at").
				Update(&entity.Activity{Cancelled: entity.ActivityCancelled, CancelledAt: time.Now()}); e != nil {
				return nil, errors.InternalServer(reason.DatabaseError).WithError(e).WithStack()
			}
		}
		return nil, nil
	})
	if err != nil {
		return err
	}
	return
}
