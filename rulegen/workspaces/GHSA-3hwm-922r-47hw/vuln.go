package main

package cmd

import (
	"fmt"
	"net/http"
	"os"
	"strings"

	"entgo.io/contrib/entgql"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi/v5"
	_ "github.com/lib/pq"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
			AllowCredentials: true,
			Debug:            os.Getenv("DEBUG") == "true",
		}).Handler)
		router.Use(api.AuthzByPolicyMiddleware)
		router.Use(api.AuthenticationMiddleware)
		if os.Getenv("DEBUG") == "true" {
			router.Use(api.LoggingMiddleware)
		}

		if *playgroudActive {
			router.Handle("/", playground.Handler("GraphQL playground", "/graphql"))
			log.Info().Msgf("connect to http://localhost:%s/ for GraphQL playground", *apiPortFlag)

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"github.com/rs/zerolog/log"
)

func (r *mutationResolver) CreateFriendship(ctx context.Context, userID uuid.UUID) (bool, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return true, nil
}

func (r *mutationResolver) DeleteFriendship(ctx context.Context, userID uuid.UUID) (bool, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return true, nil
}

func (r *mutationResolver) CreateOrUpdateFollowsGroup(ctx context.Context, input typesgen.FollowsGroupInput) (*generated.FollowsGroup, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return fg, nil
}

func (r *mutationResolver) DeleteFollowsGroup(ctx context.Context, id uuid.UUID) (bool, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return true, nil
}

func (r *mutationResolver) AssignFollowsGroupToUser(ctx context.Context, userID uuid.UUID, followsGroupID uuid.UUID, assign bool) (bool, error) {
	// Get the current user from the context
	cu, err := CurrentUserFromContext(ctx)
	return true, nil
}

func (r *mutationResolver) UpdateSettings(ctx context.Context, input typesgen.SettingsInput) (*gotype.Settings, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return &updatedUser.Settings, nil
}

func (r *mutationResolver) InternalCreateUser(ctx context.Context, input typesgen.CreateUserInput) (uuid.UUID, error) {
	campusID, err := r.client.Campus.Query().Where(campus.DuoID(input.CurrentDuoCampusID)).FirstID(ctx)
	if err != nil {
		ID(ctx)
}

func (r *mutationResolver) InternalLinkAccount(ctx context.Context, input typesgen.LinkAccountInput) (*generated.Account, error) {
	id, err := r.client.Account.Create().
		SetProvider(input.Provider.String()).
	return account, nil
}

func (r *mutationResolver) InviteOnDiscord(ctx context.Context) (bool, error) {
	cu, err := CurrentUserFromContext(ctx)
	if err != nil {
	return true, nil
}

func (r *queryResolver) Me(ctx context.Context) (*generated.User, error) {
	return CurrentUserFromContext(ctx)
}

func (r *queryResolver) SearchUser(ctx context.Context, query string, onlyOnline *bool) ([]*generated.User, error) {
	usersID, err := searchengine.NewClient().SearchUser(query, *onlyOnline)
	if err != nil {
	}).All(ctx)
}

func (r *queryResolver) Campus(ctx context.Context, id uuid.UUID) (*generated.Campus, error) {
	return r.client.Campus.Query().Where(campus.ID(id)).Only(ctx)
}

func (r *queryResolver) User(ctx context.Context, id uuid.UUID) (*generated.User, error) {
	return r.client.User.Query().Where(user.ID(id)).Only(ctx)
}

func (r *queryResolver) Location(ctx context.Context, id uuid.UUID) (*generated.Location, error) {
	return r.client.Location.Query().Where(location.ID(id)).Only(ctx)
}

func (r *queryResolver) Locations(ctx context.Context, page typesgen.PageInput, campusID uuid.UUID) (*generated.LocationConnection, error) {
	return r.client.Location.Query().
		Where(location.CampusID(campusID)).
		Paginate(ctx, page.After, &page.First, page.Before, page.Last)
}

func (r *queryResolver) LocationsByCampusName(ctx context.Context, page typesgen.PageInput, campusName string) (*generated.LocationConnection, error) {
	return r.client.Campus.Query().
		Where(campus.Name(campusName)).
		Paginate(ctx, page.After, &page.First, page.Before, page.Last)
}

func (r *queryResolver) LocationsByCluster(ctx context.Context, page typesgen.PageInput, campusName string, identifierPrefix *string) (*generated.LocationConnection, error) {
	return r.client.Campus.Query().
		Where(campus.NameEqualFold(campusName)).
		Paginate(ctx, page.After, &page.First, page.Before, page.Last)
}

func (r *queryResolver) LocationsStatsByPrefixes(ctx context.Context, campusName string, identifierPrefixes []string) ([]*typesgen.LocationStats, error) {
	sqlResults := []*typesgen.LocationStats{}
	prefixes := make([]any, len(identifierPrefixes))
	return finalLocationStats, err
}

func (r *queryResolver) MyFollowings(ctx context.Context, followsGroupID *uuid.UUID, followsGroupSlug *string) ([]*generated.User, error) {
	cu, _ := CurrentUserFromContext(ctx)

		All(ctx)
}

func (r *queryResolver) MyFollowsGroups(ctx context.Context) ([]*generated.FollowsGroup, error) {
	cu, _ := CurrentUserFromContext(ctx)

		All(ctx)
}

func (r *queryResolver) FollowsGroupsForUser(ctx context.Context, userID uuid.UUID) ([]*generated.FollowsGroup, error) {
	// Get the current user from the context
	cu, err := CurrentUserFromContext(ctx)
		All(ctx)
}

func (r *queryResolver) InternalGetUserByAccount(ctx context.Context, provider typesgen.Provider, uid string) (*generated.User, error) {
	return r.client.Account.Query().
		Where(account.Provider(provider.String()), account.ProviderAccountID(uid)).
		Only(ctx)
}

func (r *queryResolver) InternalGetUserByEmail(ctx context.Context, email string) (*generated.User, error) {
	return r.client.User.Query().
		Where(user.Email(email)).
		Only(ctx)
}

func (r *queryResolver) InternalGetUser(ctx context.Context, id uuid.UUID) (*generated.User, error) {
	return r.client.User.Get(ctx, id)
}

func (r *userResolver) IsSwimmer(ctx context.Context, obj *generated.User) (bool, error) {
	if obj.PoolYear == nil || obj.PoolMonth == nil {
		return false, nil
		strings.EqualFold(*obj.PoolMonth, now.Format("January"))), nil
}

func (r *userResolver) IsMe(ctx context.Context, obj *generated.User) (bool, error) {
	cu, _ := CurrentUserFromContext(ctx)

	return cu.ID == obj.ID, nil
}

func (r *userResolver) IsFollowing(ctx context.Context, obj *generated.User) (bool, error) {
	cu, _ := CurrentUserFromContext(ctx)

	return false, nil
}

func (r *userResolver) IsFollower(ctx context.Context, obj *generated.User) (bool, error) {
	cu, _ := CurrentUserFromContext(ctx)

	"net/http"

	sentryhttp "github.com/getsentry/sentry-go/http"
	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
)
