package main


import (
	"context"
	"errors"
	"slices"
	"strings"
	"time"

const unknownUserID = "UNKNOWN"

var (
	ErrUserNotFound = func(err error) error {
		return zerrors.ThrowNotFound(err, "EVENT-hodc6", "Errors.User.NotFound")
	}
)

type AuthRequestRepo struct {
	Command      *command.Commands
	Query        *query.Queries
	ProjectProvider           projectProvider
	ApplicationProvider       applicationProvider
	CustomTextProvider        customTextProvider
	PasswordChecker           passwordChecker

	IdGenerator id.Generator
}
}

type userViewProvider interface {
	UserByID(context.Context, string, string) (*user_view_model.UserView, error)
}

type loginPolicyViewProvider interface {
	CustomTextListByTemplate(ctx context.Context, aggregateID string, text string, withOwnerRemoved bool) (texts *query.CustomTexts, err error)
}

type passwordChecker interface {
	HumanCheckPassword(ctx context.Context, resourceOwner, userID, password string, authReq *domain.AuthRequest) error
}

func (repo *AuthRequestRepo) Health(ctx context.Context) error {
	return repo.AuthRequests.Health(ctx)
}
	request, err := repo.getAuthRequestEnsureUser(ctx, authReqID, userAgentID, userID)
	if err != nil {
		if isIgnoreUserNotFoundError(err, request) {
			// use the same errorID as below (otherwise it would expose the error reason)
			return zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
		}
		return err
	}
	err = repo.PasswordChecker.HumanCheckPassword(ctx, resourceOwner, userID, password, request.WithCurrentInfo(info))
	if isIgnoreUserInvalidPasswordError(err, request) {
		// use the same errorID as above (otherwise it would expose the error reason)
		return zerrors.ThrowInvalidArgument(nil, "EVENT-SDe2f", "Errors.User.UsernameOrPassword.Invalid")
	}
	return err
}

func isIgnoreUserNotFoundError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.Is(err, ErrUserNotFound(nil))
}

func isIgnoreUserInvalidPasswordError(err error, request *domain.AuthRequest) bool {
	return request != nil && request.LoginPolicy != nil && request.LoginPolicy.IgnoreUnknownUsernames && errors.Is(err, command.ErrPasswordInvalid(nil))
}

func lockoutPolicyToDomain(policy *query.LockoutPolicy) *domain.LockoutPolicy {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	user, viewErr := viewProvider.UserByID(ctx, userID, authz.GetInstance(ctx).InstanceID())
	if viewErr != nil && !zerrors.IsNotFound(viewErr) {
		return nil, viewErr
	} else if user == nil {
	}
	if len(events) == 0 {
		if viewErr != nil {
			// We already returned all errors apart from not found, but need to make sure that can be checked in case IgnoreUnknownUsernames option is active.
			return nil, ErrUserNotFound(viewErr)
		}
		return user_view_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	ErrPasswordInvalid = func(err error) error {
		return zerrors.ThrowInvalidArgument(err, "COMMAND-3M0fs", "Errors.User.Password.Invalid")
	}
	ErrPasswordUnchanged = func(err error) error {
		return zerrors.ThrowPreconditionFailed(err, "COMMAND-Aesh5", "Errors.User.Password.NotChanged")
	}
)

func (c *Commands) SetPassword(ctx context.Context, orgID, userID, password string, oneTime bool) (objectDetails *domain.ObjectDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()
		return nil
	}
	if errors.Is(err, passwap.ErrPasswordMismatch) {
		return ErrPasswordInvalid(err)
	}
	if errors.Is(err, passwap.ErrPasswordNoChange) {
		return ErrPasswordUnchanged(err)
	}
	return zerrors.ThrowInternal(err, "COMMAND-CahN2", "Errors.Internal")
}
	userTable = "auth.users3"
)

func (v *View) UserByID(ctx context.Context, userID, instanceID string) (*model.UserView, error) {
	return view.UserByID(ctx, v.Db, userID, instanceID)
}

func (v *View) UserByLoginName(ctx context.Context, loginName, instanceID string) (*model.UserView, error) {
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	return view.UserByID(ctx, v.Db, queriedUser.ID, instanceID)
}

func (v *View) UserByLoginNameAndResourceOwner(ctx context.Context, loginName, resourceOwner, instanceID string) (*model.UserView, error) {
	}

	//nolint: contextcheck // no lint was added because refactor would change too much code
	user, err := view.UserByID(ctx, v.Db, queriedUser.ID, instanceID)
	if err != nil {
		return nil, err
	}
		OnError(err).
		Errorf("could not get current sequence for userByID")

	user, err := view.UserByID(ctx, v.Db, queriedUser.ID, instanceID)
	if err != nil && !zerrors.IsNotFound(err) {
		return nil, err
	}
