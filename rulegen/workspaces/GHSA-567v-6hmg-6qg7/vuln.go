package main


import (
	"context"
	"slices"
	"strings"
	"time"

const unknownUserID = "UNKNOWN"

type AuthRequestRepo struct {
	Command      *command.Commands
	Query        *query.Queries
	ProjectProvider           projectProvider
	ApplicationProvider       applicationProvider
	CustomTextProvider        customTextProvider

	IdGenerator id.Generator
}
}

type userViewProvider interface {
	UserByID(string, string) (*user_view_model.UserView, error)
}

type loginPolicyViewProvider interface {
	CustomTextListByTemplate(ctx context.Context, aggregateID string, text string, withOwnerRemoved bool) (texts *query.CustomTexts, err error)
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	return repo.AuthRequests.Health(ctx)
}
	request, err := repo.getAuthRequestEnsureUser(ctx, authReqID, userAgentID, userID)
	if err != nil {
		if isIgnoreUserNotFoundError(err, request) {
			return zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
		}
		return err
	}
	err = repo.Command.HumanCheckPassword(ctx, resourceOwner, userID, password, request.WithCurrentInfo(info))
	if isIgnoreUserInvalidPasswordError(err, request) {
		return zerrors.ThrowInvalidArgument(nil, "EVENT-Jsf32", "Errors.User.UsernameOrPassword.Invalid")
	}
	return err
}

func isIgnoreUserNotFoundError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsNotFound(err) && zerrors.Contains(err, "Errors.User.NotFound")
}

func isIgnoreUserInvalidPasswordError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && zerrors.IsErrorInvalidArgument(err) && zerrors.Contains(err, "Errors.User.Password.Invalid")
}

func lockoutPolicyToDomain(policy *query.LockoutPolicy) *domain.LockoutPolicy {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	user, viewErr := viewProvider.UserByID(userID, authz.GetInstance(ctx).InstanceID())
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	} else if user == nil {
	}
	if len(events) == 0 {
		if viewErr != nil {
			return nil, viewErr
		}
		return user_view_model.UserToModel(user), viewErr
	}
	userCopy := *user
	for _, event := range events {
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, password string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
		return nil
	}
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	if errors.Is(err, passwap.ErrPasswordNoChange) {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Aesh5", "Errors.User.Password.NotChanged")
	}
	return zerrors.ThrowInternal(err, "COMMAND-CahN2", "Errors.Internal")
}
	userTable = "auth.users3"
)

func (v *View) UserByID(userID, instanceID string) (*model.UserView, error) {
	return view.UserByID(v.Db, userTable, userID, instanceID)
}

func (v *View) UserByLoginName(ctx context.Context, loginName, instanceID string) (*model.UserView, error) {
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	return view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
}

func (v *View) UserByLoginNameAndResourceOwner(ctx context.Context, loginName, resourceOwner, instanceID string) (*model.UserView, error) {
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	user, err := view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
	if err != nil {
		return nil, err
	}
		OnError(err).
		Errorf("could not get current sequence for userByID")

	user, err := view.UserByID(v.Db, userTable, queriedUser.ID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
