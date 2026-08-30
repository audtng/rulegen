package main

	UpdatedBefore    timeutil.TimeStamp
	ConcurrencyGroup string
	OrderBy          db.SearchOrderBy
	// AccessibleRepoIDsSubQuery, when non-nil, restricts results to the repo IDs selected by the
	// subquery (the caller's accessible repos). A nil value means no restriction. Using a subquery
	// instead of a materialized ID slice avoids exceeding DB parameter limits for large owners.
	AccessibleRepoIDsSubQuery *builder.Builder
}

var JobOrderByMap = map[string]map[string]db.SearchOrderBy{
		}
		cond = cond.And(builder.Eq{"`action_run_job`.concurrency_group": opts.ConcurrencyGroup})
	}
	if opts.AccessibleRepoIDsSubQuery != nil {
		cond = cond.And(builder.In("`action_run_job`.repo_id", opts.AccessibleRepoIDsSubQuery))
	}
	return cond
}

	Status           []Status
	ConcurrencyGroup string
	CommitSHA        string
	// AccessibleRepoIDsSubQuery, when non-nil, restricts results to the repo IDs selected by the
	// subquery (the caller's accessible repos). A nil value means no restriction. Using a subquery
	// instead of a materialized ID slice avoids exceeding DB parameter limits for large owners.
	AccessibleRepoIDsSubQuery *builder.Builder
}

func (opts FindRunOptions) ToConds() builder.Cond {
	if opts.CommitSHA != "" {
		cond = cond.And(builder.Eq{"`action_run`.commit_sha": opts.CommitSHA})
	}
	if opts.AccessibleRepoIDsSubQuery != nil {
		cond = cond.And(builder.In("`action_run`.repo_id", opts.AccessibleRepoIDsSubQuery))
	}
	return cond
}


const ssh2keyStart = "---- BEGIN SSH2 PUBLIC KEY ----"

const (
	// the longest RSA key ssh-keygen allows to generate is 16384 bits (2048 bytes), we still relax the limit a little here
	maxKeyBinaryBytes        = 4096
	maxKeyContentBase64Bytes = maxKeyBinaryBytes * 4 / 3
	maxKeyContentExtraBytes  = 4 * 1024 // header, footer, comment
	maxKeyContentBytes       = maxKeyContentBase64Bytes + maxKeyContentExtraBytes
)

func extractTypeFromBase64Key(key string) (string, error) {
	if len(key) > maxKeyContentBase64Bytes {
		return "", util.NewInvalidArgumentErrorf("SSH public key base64 is too long")
	}
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(b) < 4 {
		return "", fmt.Errorf("invalid key format: %w", err)

// parseKeyString parses any key string in OpenSSH or SSH2 format to clean OpenSSH string (RFC4253).
func parseKeyString(content string) (string, error) {
	if len(content) > maxKeyContentBytes {
		return "", util.NewInvalidArgumentErrorf("SSH public key content is too long")
	}

	// remove whitespace at start and end
	content = strings.TrimSpace(content)

		// Transform all legal line endings to a single "\n".
		content = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(content)

		var b strings.Builder
		b.Grow(len(content))
		lines := strings.Split(content, "\n")
		continuationLine := false

			if continuationLine || strings.ContainsAny(line, ":-") {
				continuationLine = strings.HasSuffix(line, "\\")
			} else {
				b.WriteString(line)
			}
		}
		keyContent = b.String()

		t, err := extractTypeFromBase64Key(keyContent)
		if err != nil {
	}
}

func TestPublicKeysAreExternallyManaged(t *testing.T) {
	key1 := unittest.AssertExistsAndLoadBean(t, &PublicKey{ID: 1})
	externals, err := PublicKeysAreExternallyManaged(t.Context(), []*PublicKey{key1})
	assert.NoError(t, err)
	assert.Len(t, externals, 1)
	assert.False(t, externals[0])
}

// TestCheckPublicKeyStringOversized tests if oversized SSH2 public key strings are rejected before triggering costly operations.
func TestCheckPublicKeyStringOversized(t *testing.T) {
	_, err := parseKeyString(strings.Repeat("a", maxKeyContentBytes+1))
	assert.ErrorContains(t, err, "SSH public key content is too long")

	content := "---- BEGIN SSH2 PUBLIC KEY ----\n" + strings.Repeat("a", maxKeyContentBase64Bytes+1) + "\n--- END SSH2 PUBLIC KEY ----"
	_, err = parseKeyString(content)
	assert.ErrorContains(t, err, "SSH public key base64 is too long")
}
	return bitmap.hasScope(AccessTokenScopePublicOnly)
}

// CanCreateChildScope reports whether a request authenticated by this (parent) scope may mint a token
// carrying the child scope. It rejects any grantable scope the parent does not hold, closing the
// scope-escalation path. public-only is a restriction rather than a grantable permission, so it is
// ignored here (a child may always be public-only); EnforcePublicOnlyFrom handles carrying it down.
func (s AccessTokenScope) CanCreateChildScope(child AccessTokenScope) (bool, error) {
	requested := child.StringSlice()
	scopes := make([]AccessTokenScope, 0, len(requested))
	for _, sc := range requested {
		childScope := AccessTokenScope(sc)
		if childScope == AccessTokenScopePublicOnly {
			continue
		}
		scopes = append(scopes, childScope)
	}
	return s.HasScope(scopes...)
}

// EnforcePublicOnlyFrom adds the public-only restriction to s when the authorizing parent scope is
// public-only, so a public-only token cannot mint a child token that drops the restriction.
func (s AccessTokenScope) EnforcePublicOnlyFrom(parent AccessTokenScope) (AccessTokenScope, error) {
	publicOnly, err := parent.PublicOnly()
	if err != nil {
		return "", err
	}
	if !publicOnly {
		return s, nil
	}
	return AccessTokenScope(string(s) + "," + string(AccessTokenScopePublicOnly)).Normalize()
}

// HasScope returns true if the string has the given scope
func (s AccessTokenScope) HasScope(scopes ...AccessTokenScope) (bool, error) {
	bitmap, err := s.parse()
		})
	}
}

func TestAccessTokenScope_EnforcePublicOnlyFrom(t *testing.T) {
	tests := []struct {
		in     AccessTokenScope
		parent AccessTokenScope
		out    AccessTokenScope
	}{
		// public-only parent forces the restriction onto the minted scope
		{"write:user", "write:user,public-only", "public-only,write:user"},
		// already public-only stays public-only
		{"public-only,read:user", "public-only", "public-only,read:user"},
		// non-public-only parent leaves the scope untouched
		{"write:user", "write:user", "write:user"},
		{"all", "all", "all"},
	}
	for _, test := range tests {
		t.Run(string(test.parent)+"->"+string(test.in), func(t *testing.T) {
			got, err := test.in.EnforcePublicOnlyFrom(test.parent)
			assert.NoError(t, err)
			assert.Equal(t, test.out, got)
		})
	}
}
		return nil
	}
	return db.WithTx(ctx, func(ctx context.Context) error {
		issue, err := GetIssueByID(ctx, c.IssueID)
		if err != nil {
			return err
		}
		attachments, err := repo_model.GetAttachmentsByUUIDs(ctx, uuids)
		if err != nil {
			return fmt.Errorf("getAttachmentsByUUIDs [uuids: %v]: %w", uuids, err)
		}
		for i := range attachments {
			if err := validateAttachmentForIssue(ctx, issue, attachments[i]); err != nil {
				return err
			}
			attachments[i].IssueID = c.IssueID
			attachments[i].CommentID = c.ID
			if err := repo_model.UpdateAttachment(ctx, attachments[i]); err != nil {
	assert.NoError(t, unittest.PrepareTestDatabase())

	comment := unittest.AssertExistsAndLoadBean(t, &issues_model.Comment{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: comment.IssueID})
	attachment := repo_model.Attachment{
		RepoID: issue.RepoID, // must match the comment's repo, else the cross-repo guard rejects it
		Name:   "test.txt",
	}
	assert.NoError(t, db.Insert(t.Context(), &attachment))

	return err
}

// validateAttachmentForIssue rejects a foreign or already-linked attachment before it is linked to
// issue: a known UUID could otherwise re-link (and expose) another repo's private attachment. A
// legacy attachment predating repo_id-on-upload is adopted into the issue's repo.
func validateAttachmentForIssue(ctx context.Context, issue *Issue, attachment *repo_model.Attachment) error {
	if attachment.RepoID == 0 && attachment.CreatedUnix < repo_model.LegacyAttachmentMissingRepoIDCutoff {
		attachment.RepoID = issue.RepoID
		if err := repo_model.UpdateAttachmentByUUID(ctx, attachment, "repo_id"); err != nil {
			return fmt.Errorf("update attachment repo_id [id: %d]: %w", attachment.ID, err)
		}
	}
	if attachment.RepoID != issue.RepoID {
		return util.NewPermissionDeniedErrorf("attachment belongs to a different repository")
	}
	if attachment.IssueID != 0 && attachment.IssueID != issue.ID {
		return util.NewPermissionDeniedErrorf("attachment is already linked to another issue")
	}
	if attachment.ReleaseID != 0 {
		return util.NewPermissionDeniedErrorf("attachment is already linked to a release")
	}
	return nil
}

// UpdateIssueAttachments update attachments by UUIDs for the issue
func UpdateIssueAttachments(ctx context.Context, issueID int64, uuids []string) (err error) {
	if len(uuids) == 0 {
		return nil
	}
	return db.WithTx(ctx, func(ctx context.Context) error {
		issue, err := GetIssueByID(ctx, issueID)
		if err != nil {
			return err
		}
		attachments, err := repo_model.GetAttachmentsByUUIDs(ctx, uuids)
		if err != nil {
			return fmt.Errorf("getAttachmentsByUUIDs [uuids: %v]: %w", uuids, err)
		}
		for i := range attachments {
			if err := validateAttachmentForIssue(ctx, issue, attachments[i]); err != nil {
				return err
			}
			attachments[i].IssueID = issueID
			if err := repo_model.UpdateAttachment(ctx, attachments[i]); err != nil {
				return fmt.Errorf("update attachment [id: %d]: %w", attachments[i].ID, err)
// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package issues_test

import (
	"testing"

	issues_model "gitea.dev/models/issues"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unittest"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpdateIssueAttachmentsCrossRepo(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	// attachment id 2 belongs to repo 2 / issue 4; issue 1 lives in repo 1
	issue1 := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{ID: 1})
	foreign := unittest.AssertExistsAndLoadBean(t, &repo_model.Attachment{ID: 2})
	require.NotEqual(t, issue1.RepoID, foreign.RepoID)

	// re-linking a foreign repo's attachment by UUID must be rejected
	err := issues_model.UpdateIssueAttachments(t.Context(), issue1.ID, []string{foreign.UUID})
	assert.ErrorIs(t, err, util.ErrPermissionDenied)

	// the foreign attachment must be left untouched
	reloaded := unittest.AssertExistsAndLoadBean(t, &repo_model.Attachment{ID: 2})
	assert.Equal(t, foreign.IssueID, reloaded.IssueID)
}

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strconv"
	RepoID  int64
}

func (err ErrRepoLabelNotExist) Error() string {
	return fmt.Sprintf("label does not exist [label_id: %d, repo_id: %d]", err.LabelID, err.RepoID)
}
	return l, nil
}

// GetLabelInRepoOrOrgByID returns the label with labelID scoped to the repo, falling back to the
// repo's owning organization when ownerIsOrg is set. It returns ErrRepoLabelNotExist /
// ErrOrgLabelNotExist when the label is in neither scope, so a foreign-but-existing label ID is
// indistinguishable from a nonexistent one (no cross-repo enumeration oracle).
func GetLabelInRepoOrOrgByID(ctx context.Context, repoID, ownerID int64, ownerIsOrg bool, labelID int64) (*Label, error) {
	label, err := GetLabelInRepoByID(ctx, repoID, labelID)
	if err != nil && errors.Is(err, util.ErrNotExist) && ownerIsOrg {
		return GetLabelInOrgByID(ctx, ownerID, labelID)
	}
	return label, err
}

// GetLabelInRepoByID returns a label by ID in given repository.
func GetLabelInRepoByID(ctx context.Context, repoID, labelID int64) (*Label, error) {
	if labelID <= 0 || repoID <= 0 {
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/timeutil"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
)
	assert.Equal(t, "label1", label.Name)

	_, err = issues_model.GetLabelInRepoByName(t.Context(), 1, "")
	assert.ErrorIs(t, err, util.ErrNotExist)

	_, err = issues_model.GetLabelInRepoByName(t.Context(), unittest.NonexistentID, "nonexistent")
	assert.ErrorIs(t, err, util.ErrNotExist)
}

func TestGetLabelInRepoByNames(t *testing.T) {
	assert.EqualValues(t, 1, label.ID)

	_, err = issues_model.GetLabelInRepoByID(t.Context(), 1, -1)
	assert.ErrorIs(t, err, util.ErrNotExist)

	_, err = issues_model.GetLabelInRepoByID(t.Context(), unittest.NonexistentID, unittest.NonexistentID)
	assert.ErrorIs(t, err, util.ErrNotExist)
}

func TestGetLabelsInRepoByIDs(t *testing.T) {
}

// userOrgTeamUnitRepoBuilder returns repo ids where user's teams can access the special unit.
// A team grants the unit either through an explicit team_unit row (access_mode > none) or by being an
// admin/owner team (team.authorize >= admin), which grants every unit regardless of team_unit rows —
// mirroring the HasAdminAccess() short-circuit in access.GetIndividualUserRepoPermission.
func userOrgTeamUnitRepoBuilder(userID int64, unitType unit.Type) *builder.Builder {
	return userOrgTeamRepoBuilder(userID).
		Join("INNER", "team", "`team`.id = `team_repo`.team_id").
		Join("LEFT", "team_unit", builder.Expr("`team_unit`.team_id = `team_repo`.team_id AND `team_unit`.`type` = ?", unitType)).
		Where(builder.Or(
			builder.Gte{"`team`.authorize": int(perm.AccessModeAdmin)},
			builder.Gt{"`team_unit`.`access_mode`": int(perm.AccessModeNone)},
		))
}

// userOrgTeamUnitRepoCond returns a condition to select repo ids where user's teams can access the special unit.
func UserOrgUnitRepoCond(idStr string, userID, orgID int64, unitType unit.Type) builder.Cond {
	return builder.In(idStr,
		userOrgTeamUnitRepoBuilder(userID, unitType).
			And(builder.Eq{"`team`.org_id": orgID}),
	)
}

	))
}

// PublicRepoUnderPublicOwnerCond restricts to public repos whose owner is publicly visible: the
// "genuinely public" set a public-only token or an anonymous caller may see (a public repo under a
// limited/private owner is not publicly reachable and must be excluded).
func PublicRepoUnderPublicOwnerCond() builder.Cond {
	return builder.And(
		builder.Eq{"`repository`.is_private": false},
		builder.In("`repository`.owner_id", builder.Select("id").From("`user`").Where(builder.Eq{"visibility": structs.VisibleTypePublic})),
	)
}

// UserActionsAccessibleOwnerRepoCond selects the repos owned by ownerID whose Actions `user` may read.
// It is used to list an org/user's Actions runs and jobs (see the callers in routers/api/v1/shared).
//   - owner_id = ownerID: only that owner's repos.
//   - AccessibleRepositoryCondition(user, TypeActions): only repos whose Actions the user can read
//     (admin/owner teams are handled inside it; a site admin is not, callers must skip the filter for one).
//   - publicOnly (a public-only token): additionally limit to public repos under a public owner.
func UserActionsAccessibleOwnerRepoCond(ownerID int64, user *user_model.User, publicOnly bool) builder.Cond {
	cond := builder.NewCond().And(
		builder.Eq{"`repository`.owner_id": ownerID},
		AccessibleRepositoryCondition(user, unit.TypeActions),
	)
	if publicOnly {
		cond = cond.And(PublicRepoUnderPublicOwnerCond())
	}
	return cond
}

// FindUserActionsAccessibleOwnerRepoIDsSubQuery returns a subquery selecting the repository IDs the user
// can see for the given owner. Callers embed it in an `IN (...)` condition so that a large owner does not
// materialize every repo ID into the SQL statement, which could exceed database parameter limits.
func FindUserActionsAccessibleOwnerRepoIDsSubQuery(ownerID int64, user *user_model.User, publicOnly bool) *builder.Builder {
	return builder.Select("id").From("repository").Where(UserActionsAccessibleOwnerRepoCond(ownerID, user, publicOnly))
}

// GetUserRepositories returns a list of repositories of given user.
func GetUserRepositories(ctx context.Context, opts SearchRepoOptions) (RepositoryList, int64, error) {
	if len(opts.OrderBy) == 0 {

	"gitea.dev/models/db"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unit"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/optional"
		})
	}
}

func TestFindUserActionsAccessibleOwnerRepoIDs(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())
	// user2 is on org3's owner team, so it can access org3's private repo3 (which has the actions unit)
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// org3 is a public org owning repo3 (private) and repo32 (public), both with the actions unit
	const orgID = 3

	all, err := repo_model.SearchRepositoryIDsByCondition(t.Context(), repo_model.UserActionsAccessibleOwnerRepoCond(orgID, user, false))
	require.NoError(t, err)
	assert.Contains(t, all, int64(3), "without public-only the private repo's actions are listed")

	publicOnly, err := repo_model.SearchRepositoryIDsByCondition(t.Context(), repo_model.UserActionsAccessibleOwnerRepoCond(orgID, user, true))
	require.NoError(t, err)
	assert.NotContains(t, publicOnly, int64(3), "a public-only token must not list a private repo's actions")
	assert.Contains(t, publicOnly, int64(32), "a public repo under a public owner stays listed")
}

// TestUserOrgUnitRepoCondTeamAuthorize pins the team.authorize behavior of userOrgTeamUnitRepoBuilder
// (exercised through UserOrgUnitRepoCond): an admin/owner team grants every unit even without an explicit
// team_unit row, while a non-admin team only grants a unit it has an explicit row for. This guards both
// directions — hiding repos from admin-team members, and over-broadening a plain team's access.
func TestUserOrgUnitRepoCondTeamAuthorize(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	accessibleRepoIDs := func(userID, orgID int64, unitType unit.Type) []int64 {
		ids, err := repo_model.SearchRepositoryIDsByCondition(t.Context(),
			repo_model.UserOrgUnitRepoCond("`repository`.id", userID, orgID, unitType))
		require.NoError(t, err)
		return ids
	}

	// Case A: user18 is only on org17's owner team (team5, authorize=owner), linked to the private repo24
	// but with no Actions team_unit row. The owner authorize must still grant it, mirroring the runtime
	// HasAdminAccess() short-circuit in access.GetIndividualUserRepoPermission.
	assert.Contains(t, accessibleRepoIDs(18, 17, unit.TypeActions), int64(24),
		"an owner team grants a unit it has no explicit team_unit row for")

	// Cases B and C share one subject so the team_unit row is the only difference: user4 is only on org3's
	// write team (team2, authorize=write, non-admin), linked to the private repo3. team2 has an explicit
	// Projects row but none for Actions.
	assert.Contains(t, accessibleRepoIDs(4, 3, unit.TypeProjects), int64(3),
		"a non-admin team grants a unit it has an explicit team_unit row for")
	assert.NotContains(t, accessibleRepoIDs(4, 3, unit.TypeActions), int64(3),
		"a non-admin team must NOT grant a unit it has no team_unit row for")
}
		// only include private repos the actor can still access, so metadata does not leak after access revocation
		cond = cond.And(AccessibleRepositoryCondition(opts.Actor, unit.TypeInvalid))
	} else {
		// a public repo under a limited/private owner is not publicly reachable, so exclude it too
		cond = cond.And(PublicRepoUnderPublicOwnerCond())
	}
	return cond
}
		// only include private repos the actor can still access, so metadata does not leak after access revocation
		cond = cond.And(AccessibleRepositoryCondition(opts.Actor, unit.TypeInvalid))
	} else {
		// a public repo under a limited/private owner is not publicly reachable, so exclude it too
		cond = cond.And(PublicRepoUnderPublicOwnerCond())
	}
	return cond.And(builder.Neq{
		"watch.mode": WatchModeDont,
	require.Len(t, users, 1)
	assert.Equal(t, "user2", users[0].Name)
}

func TestStarredWatchedReposExcludeNonPublicOwners(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	const viewerID = 2
	// repo1: public repo under a public owner; repo38: public repo under a limited org (not publicly reachable)
	const publicOwnerRepo, limitedOwnerRepo = 1, 38

	require.NoError(t, db.Insert(t.Context(), &repo_model.Star{UID: viewerID, RepoID: publicOwnerRepo}))
	require.NoError(t, db.Insert(t.Context(), &repo_model.Star{UID: viewerID, RepoID: limitedOwnerRepo}))
	require.NoError(t, db.Insert(t.Context(), &repo_model.Watch{UserID: viewerID, RepoID: publicOwnerRepo, Mode: repo_model.WatchModeNormal}))
	require.NoError(t, db.Insert(t.Context(), &repo_model.Watch{UserID: viewerID, RepoID: limitedOwnerRepo, Mode: repo_model.WatchModeNormal}))

	listOpts := db.ListOptions{Page: 1, PageSize: 50}

	starred, err := repo_model.GetStarredRepos(t.Context(), &repo_model.StarredReposOptions{
		ListOptions: listOpts, StarrerID: viewerID, IncludePrivate: false,
	})
	require.NoError(t, err)
	assert.NotContains(t, repoIDs(starred), int64(limitedOwnerRepo), "a public repo under a limited owner must be hidden from a public star listing")
	assert.Contains(t, repoIDs(starred), int64(publicOwnerRepo), "a public repo under a public owner stays visible")

	watched, _, err := repo_model.GetWatchedRepos(t.Context(), &repo_model.WatchedReposOptions{
		ListOptions: listOpts, WatcherID: viewerID, IncludePrivate: false,
	})
	require.NoError(t, err)
	assert.NotContains(t, repoIDs(watched), int64(limitedOwnerRepo), "a public repo under a limited owner must be hidden from a public watch listing")
	assert.Contains(t, repoIDs(watched), int64(publicOwnerRepo), "a public repo under a public owner stays visible")
}

func repoIDs(repos []*repo_model.Repository) []int64 {
	ids := make([]int64, len(repos))
	for i, r := range repos {
		ids[i] = r.ID
	}
	return ids
}
package openid

import (
	"net/http"
	"sync"
	"time"

	"gitea.dev/modules/hostmatcher"
	"gitea.dev/modules/proxy"
	"gitea.dev/modules/setting"

	"github.com/yohcop/openid-go"
)

var (
	nonceStore     = openid.NewSimpleNonceStore()
	discoveryCache = newTimedDiscoveryCache(24 * time.Hour)

	// openIDInstance does discovery/verification via an SSRF-protected client, so a user-supplied
	// OpenID identifier can't reach internal/loopback/reserved addresses. It honors the operator's
	// [security] ALLOWED_HOST_LIST (empty defaults to "external"), matching the avatar/webhook/migration
	// clients, and validates the proxy path too. Lazy: reads proxy/settings once.
	openIDInstance = sync.OnceValue(func() *openid.OpenID {
		allowList := hostmatcher.ParseHostMatchList("security.ALLOWED_HOST_LIST", setting.Security.AllowedHostList)
		return openid.NewOpenID(&http.Client{
			Timeout:   30 * time.Second,
			Transport: hostmatcher.NewHTTPTransport("openid", allowList, nil, proxy.Proxy(), setting.Proxy.ProxyURLFixed, nil),
		})
	})
)

// Verify handles response from OpenID provider
func Verify(fullURL string) (id string, err error) {
	return openIDInstance().Verify(fullURL, discoveryCache, nonceStore)
}

// Normalize normalizes an OpenID URI

// RedirectURL redirects browser
func RedirectURL(id, callbackURL, realm string) (string, error) {
	return openIDInstance().RedirectURL(id, callbackURL, realm)
}
// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package openid

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOpenIDDiscoveryBlocksInternalHost(t *testing.T) {
	var reached atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Store(true)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// RedirectURL performs server-side discovery of the identifier URL; a loopback URL
	// must be refused at dial time instead of reaching the internal server
	_, err := RedirectURL(srv.URL, "http://example.com/callback", "http://example.com/")
	require.Error(t, err)
	assert.False(t, reached.Load(), "OpenID discovery must not reach an internal/loopback host")
}

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"net/url"
	"syscall"
	"time"
		return dialer.DialContext(ctx, network, addrOrHost)
	}
}

// NewHTTPTransport builds an http.Transport that validates the request target against the allow/block
// lists on the direct-dial path (DialContext). When an HTTP proxy is configured the proxy resolves and
// dials the target itself, so restricting the proxied target is the proxy server's responsibility, not
// Gitea's. proxyFunc selects the proxy URL per request (the http.Transport.Proxy selector, e.g.
// proxy.Proxy()); proxyURLFixed is the fixed proxy address the dialer must always permit; tlsConfig may
// be nil. blockList may be nil for callers that only maintain an allow-list.
func NewHTTPTransport(usage string, allowList, blockList *HostMatchList, proxyFunc func(*http.Request) (*url.URL, error), proxyURLFixed *url.URL, tlsConfig *tls.Config) *http.Transport {
	return &http.Transport{
		TLSClientConfig: tlsConfig,
		Proxy:           proxyFunc,
		DialContext:     NewDialContext(usage, allowList, blockList, proxyURLFixed),
	}
}
	"strconv"
	"strings"

	"gitea.dev/modules/packages"
	"gitea.dev/modules/util"
	"gitea.dev/modules/validation"

	namePattern = regexp.MustCompile(`\A[a-zA-Z0-9@._+-]+\z`)
	// (epoch:pkgver-pkgrel)
	versionPattern = regexp.MustCompile(`\A(?:\d:)?[\w.+~]+(?:-[-\w.+~]+)?\z`)

	// caps on the accumulated package file list (vars so tests can lower them); far above
	// any legitimate package, but low enough to stop metadata amplification
	maxFileEntries   = 100000
	maxFileNameBytes = 16 * 1024 * 1024
)

type Package struct {
	}

	var p *Package
	files := packages.NewBoundedFileList(maxFileEntries, maxFileNameBytes)

	tr := tar.NewReader(inner)
	for {
				return nil, err
			}
		} else if !strings.HasPrefix(filename, ".") {
			if err := files.Add(hd.Name); err != nil {
				return nil, err
			}
		}
	}

		return nil, ErrMissingPKGINFOFile
	}

	p.FileMetadata.Files = files.Files()
	p.FileCompressionExtension = compressionType

	return p, nil
	"io"
	"testing"

	"gitea.dev/modules/test"
	"gitea.dev/modules/util"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/ulikunitz/xz"
		assert.ElementsMatch(t, []string{"usr/bin/paket1"}, p.FileMetadata.Backup)
	})
}

// TestParsePackageTooManyFiles ensures the accumulated file list is bounded to prevent
// metadata amplification from a package with a huge number of (tiny) file entries.
func TestParsePackageTooManyFiles(t *testing.T) {
	defer test.MockVariableValue(&maxFileEntries, 3)()
	buf := test.WriteTarCompression(func(w io.Writer) io.WriteCloser { return gzip.NewWriter(w) }, map[string]string{
		"file1":    "content1",
		".PKGINFO": string(createPKGINFOContent(packageName, packageVersion)),
	})
	_, err := ParsePackage(buf)
	assert.NoError(t, err)

	buf = test.WriteTarCompression(func(w io.Writer) io.WriteCloser { return gzip.NewWriter(w) }, map[string]string{
		"file1":    "content1",
		"file2":    "content2",
		"file3":    "content3",
		"file4":    "content4",
		".PKGINFO": string(createPKGINFOContent(packageName, packageVersion)),
	})
	_, err = ParsePackage(buf)
	assert.ErrorIs(t, err, util.ErrInvalidArgument)
}
				return nil, GlobalVars().ErrUnsupportedCompression
			}

			// bound the decompressed control archive: it holds only the small control file
			// and maintainer scripts, so a much larger stream is a decompression bomb
			const maxControlTarSize = 32 * 1024 * 1024
			tr := tar.NewReader(io.LimitReader(inner, maxControlTarSize))
			for {
				hd, err := tr.Next()
				if err == io.EOF {
	key := ""
	var depends strings.Builder
	var control strings.Builder
	var description strings.Builder

	// https://www.debian.org/doc/debian-policy/ch-controlfields.html#syntax-of-control-files
	s := bufio.NewScanner(r)
		control.WriteString(line)
		control.WriteByte('\n')

		// a leading space or tab marks a folded continuation line that belongs to the previous field
		// (identified by key), not a new "Key: value" pair; only the multi-line fields append here.
		// Continuation lines may themselves contain a colon, so they must not be re-split on ":".
		if line[0] == ' ' || line[0] == '\t' {
			switch key {
			case "Description":
				description.WriteString(line)
			case "Depends":
				depends.WriteString(trimmed)
			}
					p.Metadata.Maintainer = a.Name
				}
			case "Description":
				description.Reset()
				description.WriteString(value)
			case "Depends":
				depends.WriteString(value)
			case "Homepage":
		return nil, GlobalVars().ErrInvalidArchitecture
	}

	p.Metadata.Description = description.String()

	dependencies := strings.Split(depends.String(), ",")
	for i := range dependencies {
		dependencies[i] = strings.TrimSpace(dependencies[i])
		assert.True(t, IsValidDistributionOrComponent(name), "good=%q", name)
	}
}

// TestParseControlFileMultilineDescription verifies a multi-line Description is assembled in order
// (the parser accumulates it in a strings.Builder); it guards the assembled value, not its timing.
func TestParseControlFileMultilineDescription(t *testing.T) {
	var buf bytes.Buffer
	buf.WriteString("Package: testpkg\nVersion: 1.0\nArchitecture: amd64\nDescription: short summary\n more details\n even more\n")

	p, err := ParseControlFile(&buf)
	assert.NoError(t, err)
	assert.NotNil(t, p)
	assert.Equal(t, "short summary more details even more", p.Metadata.Description)
}
// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package packages

import "gitea.dev/modules/util"

// BoundedFileList accumulates file names from a package archive while enforcing caps on the number of
// entries and their total name length, returning an error once either cap would be exceeded.
type BoundedFileList struct {
	files     []string
	nameBytes int
	maxFiles  int
	maxBytes  int
}

// NewBoundedFileList creates a BoundedFileList with the given caps; a non-positive cap falls back to the
// corresponding default.
func NewBoundedFileList(maxFiles, maxNameBytes int) *BoundedFileList {
	return &BoundedFileList{maxFiles: maxFiles, maxBytes: maxNameBytes}
}

// Add appends name, returning util.ErrInvalidArgument once the entry count or accumulated byte length
// would exceed the configured cap.
func (b *BoundedFileList) Add(name string) error {
	if len(b.files) >= b.maxFiles || b.nameBytes+len(name) > b.maxBytes {
		return util.NewInvalidArgumentErrorf("package contains too many file entries")
	}
	b.nameBytes += len(name)
	b.files = append(b.files, name)
	return nil
}

// Files returns the accumulated file names.
func (b *BoundedFileList) Files() []string {
	return b.files
}
	"crypto/tls"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	return conn, nil
}

// internalAPIConnectionIsLocal reports whether the internal API transport connects to a local target,
// where the self-signed local certificate cannot be verified so skipping verification is safe. It mirrors
// what dialContextInternalAPI actually dials: a unix socket whenever Protocol is HTTPUnix (always local,
// whatever LOCAL_ROOT_URL says), otherwise the LOCAL_ROOT_URL host directly. A non-loopback LOCAL_ROOT_URL
// is a real network hop, so its certificate must be verified, else the internal token can be MITM'd. An
// unparseable LOCAL_ROOT_URL is a hard misconfiguration and fails closed (verify).
func internalAPIConnectionIsLocal(protocol setting.Scheme, localURL string) bool {
	if protocol == setting.HTTPUnix {
		return true
	}
	u, err := url.Parse(localURL)
	if err != nil {
		return false
	}
	host := u.Hostname()
	if host == "localhost" {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}

var internalAPITransport = sync.OnceValue(func() http.RoundTripper {
	return &http.Transport{
		DialContext: dialContextInternalAPI,
		TLSClientConfig: &tls.Config{
			// Skip verification only for a local target (unix socket, or a loopback LOCAL_ROOT_URL), where the
			// self-signed local cert can't be verified anyway; a non-loopback LOCAL_ROOT_URL is a real network
			// hop and must be verified so the internal token can't be MITM'd. When verifying, Go's default
			// ServerName (the dialed LOCAL_ROOT_URL host) is already correct, so it is not overridden.
			InsecureSkipVerify: internalAPIConnectionIsLocal(setting.Protocol, setting.LocalURL),
		},
	}
})
// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package private

import (
	"testing"

	"gitea.dev/modules/setting"

	"github.com/stretchr/testify/assert"
)

func TestInternalAPIConnectionIsLocal(t *testing.T) {
	cases := []struct {
		name     string
		protocol setting.Scheme
		localURL string
		want     bool
	}{
		// HTTPUnix always dials the unix socket (a local target), whatever LOCAL_ROOT_URL says
		{"unix socket", setting.HTTPUnix, "https://gitea.example.com/", true},
		{"localhost", setting.HTTP, "http://localhost:3000/", true},
		{"loopback ipv4", setting.HTTPS, "https://127.0.0.1:3000/", true},
		{"loopback ipv6", setting.HTTPS, "https://[::1]:3000/", true},
		// a non-loopback LOCAL_ROOT_URL is a real network hop and must be verified
		{"remote host", setting.HTTPS, "https://gitea.internal:443/", false},
		{"remote ip", setting.HTTPS, "https://10.0.0.5:3000/", false},
		// an unparseable LOCAL_ROOT_URL is a hard misconfiguration; fail closed to verification
		{"invalid url", setting.HTTPS, "://bad", false},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, internalAPIConnectionIsLocal(c.protocol, c.localURL))
		})
	}
}
	XContentTypeOptions string

	ContentSecurityPolicyGeneral string // it only supports empty (default policy) or "unset", maybe it can support more in the future
	AllowedHostList              string
}{
	XFrameOptions:       "SAMEORIGIN",
	XContentTypeOptions: "nosniff",
	AllowedHostList:     "external",
}

var (
)

func TestLoadSecurityFrom(t *testing.T) {
	assert.Equal(t, "SAMEORIGIN", Security.XFrameOptions)
	assert.Equal(t, "nosniff", Security.XContentTypeOptions)
	assert.Equal(t, "external", Security.AllowedHostList)

	cfg, err := NewConfigProviderFromData(`[security]
X_FRAME_OPTIONS = DENY
X_CONTENT_TYPE_OPTIONS = unset
ALLOWED_HOST_LIST = foo
CONTENT_SECURITY_POLICY_GENERAL = "script-src *; foo"
`)
	assert.NoError(t, err)
	loadSecurityFrom(cfg)
	assert.Equal(t, "DENY", Security.XFrameOptions)
	assert.Equal(t, "unset", Security.XContentTypeOptions)
	assert.Equal(t, "foo", Security.AllowedHostList)
	assert.Equal(t, `"script-src *`, Security.ContentSecurityPolicyGeneral) // holy shit ini package bug
}
	Webhook.QueueLength = sec.Key("QUEUE_LENGTH").MustInt(1000)
	Webhook.DeliverTimeout = sec.Key("DELIVER_TIMEOUT").MustInt(5)
	Webhook.SkipTLSVerify = sec.Key("SKIP_TLS_VERIFY").MustBool()

	deprecatedSetting(rootCfg, "webhook", "ALLOWED_HOST_LIST", "security", "ALLOWED_HOST_LIST", "v28.0.0")
	Webhook.AllowedHostList = sec.Key("ALLOWED_HOST_LIST").MustString(Security.AllowedHostList)

	Webhook.Types = []string{"gitea", "gogs", "slack", "discord", "dingtalk", "telegram", "msteams", "feishu", "matrix", "wechatwork", "packagist"}
	Webhook.PagingNum = sec.Key("PAGING_NUM").MustInt(10)
	Webhook.ProxyURL = sec.Key("PROXY_URL").MustString("")

// Open open a local file or a remote file
func Open(uriStr string) (io.ReadCloser, error) {
	return OpenWithClient(uriStr, http.DefaultClient)
}

// OpenWithClient opens a local file or a remote file, using the given (non-nil) HTTP client
// for http/https URLs. Callers that must confine remote access (e.g. to defeat SSRF via
// redirects) should pass a client whose transport validates the peer at dial time; Open
// passes http.DefaultClient.
func OpenWithClient(uriStr string, client *http.Client) (io.ReadCloser, error) {
	u, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		f, err := client.Get(uriStr)
		if err != nil {
			return nil, err
		}
package uri

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadURI(t *testing.T) {
	assert.NoError(t, err)
	defer f.Close()
}

// TestOpenWithClientValidatesRedirectTarget verifies OpenWithClient routes the
// whole request chain (including redirects) through the provided client, so a
// client whose transport refuses to dial an internal target blocks a redirect to
// it — whereas the default client (old Open behavior) follows it.
func TestOpenWithClientValidatesRedirectTarget(t *testing.T) {
	var internalHit atomic.Bool
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		internalHit.Store(true)
		_, _ = w.Write([]byte("secret"))
	}))
	defer internal.Close()

	front := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, internal.URL, http.StatusFound)
	}))
	defer front.Close()

	internalAddr := strings.TrimPrefix(internal.URL, "http://")

	// a client that refuses to dial the internal target, mimicking the migration
	// hostmatcher dialer that re-validates every hop
	blockingClient := &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == internalAddr {
				return nil, errors.New("blocked internal address")
			}
			return (&net.Dialer{}).DialContext(ctx, network, addr)
		},
	}}

	_, err := OpenWithClient(front.URL, blockingClient)
	require.Error(t, err)
	assert.False(t, internalHit.Load(), "the redirect target must not be reached through the validating client")

	// the default client (the previous behavior) follows the redirect to the internal target
	internalHit.Store(false)
	rc, err := Open(front.URL)
	require.NoError(t, err)
	_ = rc.Close()
	assert.True(t, internalHit.Load(), "sanity check: the default client follows the redirect")
}
				}

				if publicOnly {
					// a public-only token must not reach limited-visibility owners either,
					// matching how orgs/users are enforced elsewhere in this file
					if ctx.Package != nil && !ctx.Package.Owner.Visibility.IsPublic() {
						ctx.HTTPError(http.StatusForbidden, "reqToken", "token scope is limited to public packages")
						return
					}
func AddPackageTag(ctx *context.Context) {
	packageName := packageNameFromParams(ctx)

	// the dist-tag body is only a quoted version string; bound it to avoid an unbounded
	// read that could exhaust memory
	const maxDistTagBodySize = 4 * 1024
	body, err := io.ReadAll(io.LimitReader(ctx.Req.Body, maxDistTagBodySize+1))
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	if len(body) > maxDistTagBodySize {
		apiError(ctx, http.StatusRequestEntityTooLarge, errors.New("request body too large"))
		return
	}
	version := strings.Trim(string(body), "\"") // is as "version" in the body

	pv, err := packages_model.GetVersionByNameAndVersion(ctx, ctx.Package.Owner.ID, packages_model.TypeNpm, packageName, version)
					return
				}
			case auth_model.AccessTokenScopeCategoryPackage:
				// a public-only token must not reach limited-visibility owners either,
				// matching the org/user public-only enforcement above
				if ctx.Package != nil && !ctx.Package.Owner.Visibility.IsPublic() {
					ctx.APIError(http.StatusForbidden, "token scope is limited to public packages")
					return
				}
		return
	}

	// the label must belong to this repo (or its owning org); otherwise a foreign label ID
	// is rejected the same way as a nonexistent one, closing a cross-repo enumeration oracle
	labelID := ctx.PathParamInt64("id")
	label, err := issues_model.GetLabelInRepoOrOrgByID(ctx, ctx.Repo.Repository.ID, ctx.Repo.Owner.ID, ctx.Repo.Owner.IsOrganization(), labelID)
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}

			// allow only RepoAdmin, Admin and User to add time
			user, err = user_model.GetUserByName(ctx, form.User)
			if err != nil {
				ctx.APIErrorAuto(err)
				return
			}
		}
	}
		l, err = issues_model.GetLabelInRepoByID(ctx, ctx.Repo.Repository.ID, intID)
	}
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}

	form := web.GetForm(ctx).(*api.EditLabelOption)
	l, err := issues_model.GetLabelInRepoByID(ctx, ctx.Repo.Repository.ID, ctx.PathParamInt64("id"))
	if err != nil {
		ctx.APIErrorAuto(err)
		return
	}

		return
	}

	// a public-only token must not update (push into) a private head repo,
	// even when the base repo named in the route is public
	if !ctx.TokenCanAccessRepo(pr.HeadRepo) {
		ctx.APIErrorNotFound()
		return
	}

	// keep API back-compat: when no style is given, default to "merge" rather than the repo's DefaultUpdateStyle,
	// so existing API clients keep getting a merge update.
	rebase := repo_model.UpdateStyle(ctx.FormString("style", string(repo_model.UpdateStyleMerge))) == repo_model.UpdateStyleRebase
	"gitea.dev/routers/api/v1/utils"
	"gitea.dev/services/context"
	"gitea.dev/services/convert"

	"xorm.io/builder"
)

// actionsOwnerAccessibleRepoIDsSubQuery returns the sub-query restricting an owner-scoped actions
// listing to the repos whose actions the caller can read, or nil when no restriction applies. A bare
// org member must not be able to enumerate runs/jobs of repos they have no access to. A site admin may
// skip the access filter, but a public-only token must stay confined to public repos even for an admin.
func actionsOwnerAccessibleRepoIDsSubQuery(ctx *context.APIContext, ownerID int64) *builder.Builder {
	if ownerID > 0 && (ctx.Doer == nil || !ctx.Doer.IsAdmin || ctx.PublicOnly) {
		return repo_model.FindUserActionsAccessibleOwnerRepoIDsSubQuery(ownerID, ctx.Doer, ctx.PublicOnly)
	}
	return nil
}

// ListJobs lists jobs for api route validated ownerID and repoID
// ownerID == 0 and repoID == 0 means all jobs
// ownerID == 0 and repoID != 0 means all jobs for the given repo
		opts.Statuses = append(opts.Statuses, values...)
	}

	opts.AccessibleRepoIDsSubQuery = actionsOwnerAccessibleRepoIDsSubQuery(ctx, opts.OwnerID)

	jobs, total, err := db.FindAndCount[actions_model.ActionRunJob](ctx, opts)
	if err != nil {
		ctx.APIErrorInternal(err)
	}
	excludePullRequests := ctx.FormBool("exclude_pull_requests")

	opts.AccessibleRepoIDsSubQuery = actionsOwnerAccessibleRepoIDsSubQuery(ctx, opts.OwnerID)

	runs, total, err := db.FindAndCount[actions_model.ActionRun](ctx, opts)
	if err != nil {
		ctx.APIErrorInternal(err)
	}
	t.Scope = scope

	// a token-authenticated request must not mint a token with a broader scope than its own
	if ctx.Data["IsApiToken"] == true {
		apiTokenScope, ok := ctx.Data["ApiTokenScope"].(auth_model.AccessTokenScope)
		if !ok {
			ctx.APIError(http.StatusForbidden, "the authenticating token has no scope")
			return
		}
		hasScope, err := apiTokenScope.CanCreateChildScope(scope)
		if err != nil {
			ctx.APIErrorInternal(err)
			return
		}
		if !hasScope {
			ctx.APIError(http.StatusForbidden, "cannot create an access token with a broader scope than the authenticating token")
			return
		}
		// a public-only token must not mint a token that drops the public-only restriction
		if t.Scope, err = t.Scope.EnforcePublicOnlyFrom(apiTokenScope); err != nil {
			ctx.APIErrorInternal(err)
			return
		}
	}

	if err := auth_model.NewAccessToken(ctx, t); err != nil {
		ctx.APIErrorInternal(err)
		return
	"net/http"

	user_model "gitea.dev/models/user"
	"gitea.dev/modules/setting"
	api "gitea.dev/modules/structs"
	"gitea.dev/modules/web"
	"gitea.dev/services/context"
	//   "422":
	//     "$ref": "#/responses/validationError"

	if user_model.IsFeatureDisabledWithLoginType(ctx.Doer, setting.UserFeatureManageCredentials) {
		ctx.APIErrorNotFound("emails are not allowed to be changed")
		return
	}

	form := web.GetForm(ctx).(*api.CreateEmailOption)
	if len(form.Emails) == 0 {
		ctx.APIError(http.StatusUnprocessableEntity, "Email list empty")
	//   "404":
	//     "$ref": "#/responses/notFound"

	if user_model.IsFeatureDisabledWithLoginType(ctx.Doer, setting.UserFeatureManageCredentials) {
		ctx.APIErrorNotFound("emails are not allowed to be changed")
		return
	}

	form := web.GetForm(ctx).(*api.DeleteEmailOption)
	if len(form.Emails) == 0 {
		ctx.Status(http.StatusNoContent)
	user_model "gitea.dev/models/user"
	auth_module "gitea.dev/modules/auth"
	"gitea.dev/modules/container"
	"gitea.dev/modules/hostmatcher"
	"gitea.dev/modules/httplib"
	"gitea.dev/modules/log"
	"gitea.dev/modules/optional"
	"gitea.dev/modules/proxy"
	"gitea.dev/modules/session"
	"gitea.dev/modules/setting"
	source_service "gitea.dev/services/auth/source"
	ctx.Redirect(setting.AppSubURL + "/user/link_account")
}

// oauth2AvatarAllowList parses the host allow-list applied to avatar fetches from the global
// [security] ALLOWED_HOST_LIST, defaulting an empty setting to the built-in "external" set. An empty
// host-match list would otherwise disable the allow-list check entirely and permit any host, including
// loopback/private addresses (SSRF).
func oauth2AvatarAllowList() *hostmatcher.HostMatchList {
	return hostmatcher.ParseHostMatchList("security.ALLOWED_HOST_LIST", setting.Security.AllowedHostList)
}

// oauth2AvatarHTTPClient builds the SSRF-protected client for avatar fetches. It is constructed per call
// so a changed allowlist takes effect (avatar fetches are infrequent, so this is not a hot path).
func oauth2AvatarHTTPClient() *http.Client {
	allowList := oauth2AvatarAllowList()
	return &http.Client{
		Timeout:   30 * time.Second,
		Transport: hostmatcher.NewHTTPTransport("oauth2-avatar", allowList, nil, proxy.Proxy(), setting.Proxy.ProxyURLFixed, nil),
	}
}

func oauth2UpdateAvatarIfNeed(ctx *context.Context, avatarURL string, u *user_model.User) {
	if !setting.OAuth2Client.UpdateAvatar || len(avatarURL) == 0 {
	// Some hosts (e.g. Wikimedia) reject Go's default User-Agent.
	req.Header.Set("User-Agent", "Gitea "+setting.AppVer)

	resp, err := oauth2AvatarHTTPClient().Do(req)
	if err != nil {
		log.Warn("fetch %q failed: %v", avatarURL, err)
		return
			ctx.ServerError("GetExternalLogin", err)
			return
		}
		// the cron clears all three token fields when it disables a user, so require the
		// full signature; a RefreshToken alone is empty for many normal logins (e.g. GitHub
		// or OIDC without offline_access), which would otherwise reactivate admin-disabled users
		isDisabledByAutoSync := hasExt && extLogin.AccessToken == "" && extLogin.RefreshToken == "" && extLogin.ExpiresAt.IsZero()
		if isDisabledByAutoSync {
			opts.IsActive = optional.Some(true)
		}
		return
	}

	// enforce the same user scope the REST API requires before returning identity
	// claims; OIDC access tokens map to the "all" scope, so standard OIDC clients
	// are unaffected and only explicitly-restricted tokens are rejected
	tokenScope, _ := ctx.Data["ApiTokenScope"].(auth.AccessTokenScope)
	if allowed, err := tokenScope.HasScope(auth.AccessTokenScopeReadUser); err != nil {
		ctx.ServerError("HasScope", err)
		return
	} else if !allowed {
		ctx.Resp.Header().Set("WWW-Authenticate", `Bearer realm="Gitea OAuth2"`)
		ctx.PlainText(http.StatusForbidden, "token does not have required scope: read:user")
		return
	}

	response := &userInfoResponse{
		Sub:               strconv.FormatInt(ctx.Doer.ID, 10),
		Name:              ctx.Doer.DisplayName(),
package auth

import (
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/hostmatcher"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/test"
	"gitea.dev/services/oauth2_provider"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func createAndParseToken(t *testing.T, grant *auth.OAuth2Grant) *oauth2_provider.OIDCToken {
	assert.Equal(t, user.Email, oidcToken.Email)
	assert.Equal(t, user.IsActive, oidcToken.EmailVerified)
}

func TestOAuth2AvatarClientBlocksLoopback(t *testing.T) {
	var hit atomic.Bool
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hit.Store(true)
		_, _ = w.Write([]byte("img"))
	}))
	defer srv.Close()

	// the httptest server binds a loopback address, which the SSRF-protected dialer must refuse
	resp, err := oauth2AvatarHTTPClient().Get(srv.URL)
	if resp != nil {
		_ = resp.Body.Close()
	}
	require.Error(t, err)
	assert.False(t, hit.Load(), "avatar client must refuse to dial a loopback address")
}

func TestOAuth2AvatarAllowListRestricts(t *testing.T) {
	defer test.MockVariableValue(&setting.Security.AllowedHostList, "avatars.example.com")()
	allowList := oauth2AvatarAllowList()
	assert.True(t, allowList.MatchHostName("avatars.example.com"), "the configured host must be allowed")
	assert.False(t, allowList.MatchHostName("8.8.8.8"), "an unrelated external host must be rejected")

	// the default `external` allow-list still permits external hosts
	setting.Security.AllowedHostList = hostmatcher.MatchBuiltinExternal
	assert.True(t, oauth2AvatarAllowList().MatchHostName("8.8.8.8"), "default allow-list permits external hosts")
}

func TestOAuth2AvatarClientBlocksCloudMetadata(t *testing.T) {
	// external-only allow-list must reject link-local cloud metadata (169.254.169.254) at dial time
	resp, err := oauth2AvatarHTTPClient().Get("http://169.254.169.254/latest/meta-data/")
	if resp != nil {
		_ = resp.Body.Close()
	}
	require.Error(t, err)
	assert.ErrorContains(t, err, "can only call allowed HTTP servers",
		"avatar client must refuse a link-local cloud-metadata address")
}
		includePrivate = isOrgMember
	}

	// a public-only API token must not surface private activity, even for its own owner
	if includePrivate && context.TokenIsPublicOnly(ctx) {
		includePrivate = false
	}

	actions, _, err := feed_service.GetFeeds(ctx, activities_model.GetFeedsOptions{
		RequestedUser:   ctx.ContextUser,
		Actor:           ctx.Doer,
	"path"
	"strings"

	auth_model "gitea.dev/models/auth"
	access_model "gitea.dev/models/perm/access"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unit"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/util"
	"gitea.dev/services/context"
		return
	}
	branchName := setting.Repository.DefaultBranch
	if repo, err := repo_model.GetRepositoryByOwnerAndName(ctx, ownerName, repoName); err == nil {
		branchName = goGetDefaultBranch(ctx, repo)
	}
	prefix := setting.AppURL + path.Join(url.PathEscape(ownerName), url.PathEscape(repoName), "src", "branch", util.PathEscapeSegments(branchName))

	ctx.RespHeader().Set("Content-Type", "text/html")
	_, _ = ctx.Write([]byte(res))
}

// goGetDefaultBranch returns the repository's real default branch only when the caller may genuinely
// reach it, otherwise the neutral instance default, so the meta response does not disclose the branch
// name (or the repo's existence) to callers who cannot see the repository.
func goGetDefaultBranch(ctx *context.Context, repo *repo_model.Repository) string {
	def := setting.Repository.DefaultBranch
	if len(repo.DefaultBranch) == 0 {
		return def
	}
	if err := repo.LoadOwner(ctx); err != nil || repo.Owner == nil {
		return def
	}
	// a token that was not granted repository read scope must not learn repository details, even when the
	// account behind it could read the repo through the web UI
	if !goGetTokenCanReadRepo(ctx) {
		return def
	}
	// a public-only token may only reach genuinely public resources (a public repo under a public owner)
	if context.TokenIsPublicOnly(ctx) && (repo.IsPrivate || !repo.Owner.Visibility.IsPublic()) {
		return def
	}
	// the caller must be able to read the code and see the owner: a limited/private owner hides its repos
	// from anonymous/non-member callers even when the repo itself is public
	perm, err := access_model.GetDoerRepoPermission(ctx, repo, ctx.Doer)
	if err != nil || !perm.CanRead(unit.TypeCode) || !user_model.IsUserVisibleToViewer(ctx, repo.Owner, ctx.Doer) {
		return def
	}
	return repo.DefaultBranch
}

// goGetTokenCanReadRepo reports whether the request may learn repository details. A non-token request
// always may; a token request may only when its scope grants repository read, so a PAT that was never
// scoped for repositories cannot disclose the branch even if its owner can read the repo.
func goGetTokenCanReadRepo(ctx *context.Context) bool {
	if ctx.Data["IsApiToken"] != true {
		return true
	}
	scope, ok := ctx.Data["ApiTokenScope"].(auth_model.AccessTokenScope)
	if !ok {
		return false
	}
	has, err := scope.HasScope(auth_model.AccessTokenScopeReadRepository)
	return err == nil && has
}
			}
		}
	case "attach", "detach", "toggle", "toggle-alt":
		// scope the label to this repo (or its org) so a foreign label ID is 404, not an oracle
		labelID := ctx.FormInt64("id")
		label, err := issues_model.GetLabelInRepoOrOrgByID(ctx, ctx.Repo.Repository.ID, ctx.Repo.Owner.ID, ctx.Repo.Owner.IsOrganization(), labelID)
		if err != nil {
			if errors.Is(err, util.ErrNotExist) {
				ctx.HTTPError(http.StatusNotFound, "GetLabelByID")
			} else {
				ctx.ServerError("GetLabelByID", err)
	"strings"
	"time"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/db"
	git_model "gitea.dev/models/git"
	repo_model "gitea.dev/models/repo"
		return
	}

	// a scoped or public-only API token authenticating this web request must still satisfy
	// the repository read scope before private repo content is served
	context.CheckRepoScopedToken(ctx, ctx.Repo.Repository, auth_model.Read)
	if ctx.Written() {
		return
	}

	// Check whether the repo is viewable: not in migration, and the code unit should be enabled
	// Ideally the "feed" logic should be after this, but old code did so, so keep it as-is.
	checkHomeCodeViewable(ctx)
		date := ctx.FormString("date")
		pagingNum = setting.UI.FeedPagingNum
		showPrivate := ctx.IsSigned && (ctx.Doer.IsAdmin || ctx.Doer.ID == ctx.ContextUser.ID)
		// a public-only API token must not surface private activity, even for its own owner
		if showPrivate && context.TokenIsPublicOnly(ctx) {
			showPrivate = false
		}
		items, feedCount, err := feed_service.GetFeedsForDashboard(ctx, activities_model.GetFeedsOptions{
			RequestedUser:   ctx.ContextUser,
			Actor:           ctx.Doer,
		return
	}

	// a token-authenticated request must not mint a token with a broader scope than its own, nor
	// drop the public-only restriction. Web routes accept basic-auth PATs/OAuth tokens too, so this
	// must mirror the REST API guard in routers/api/v1/user/app.go.
	if ctx.Data["IsApiToken"] == true {
		apiTokenScope, ok := ctx.Data["ApiTokenScope"].(auth_model.AccessTokenScope)
		if !ok {
			ctx.HTTPError(http.StatusForbidden, "the authenticating token has no scope")
			return
		}
		hasScope, err := apiTokenScope.CanCreateChildScope(t.Scope)
		if err != nil {
			ctx.ServerError("CanCreateChildScope", err)
			return
		}
		if !hasScope {
			ctx.HTTPError(http.StatusForbidden, "cannot create an access token with a broader scope than the authenticating token")
			return
		}
		if t.Scope, err = t.Scope.EnforcePublicOnlyFrom(apiTokenScope); err != nil {
			ctx.ServerError("EnforcePublicOnlyFrom", err)
			return
		}
	}

	if err := auth_model.NewAccessToken(ctx, t); err != nil {
		ctx.ServerError("NewAccessToken", err)
		return
	"errors"
	"fmt"
	"sync"
	"time"

	runnerv1 "gitea.dev/actions-proto-go/runner/v1"
	actions_model "gitea.dev/models/actions"
	"gitea.dev/models/db"
	secret_model "gitea.dev/models/secret"
	"gitea.dev/modules/graceful"
	"gitea.dev/modules/log"
	"gitea.dev/modules/setting"

	return task, ok, false, err
}

// releaseTaskForRunnerCleanup releases a claimed task using a fresh, bounded context. The request context
// is typically already canceled when we reach the release paths below, and a DB transaction on a canceled
// context fails immediately, which would strand the claimed job in running state.
func releaseTaskForRunnerCleanup(t *actions_model.ActionTask) {
	ctx, cancel := context.WithTimeout(graceful.GetManager().ShutdownContext(), 10*time.Second)
	defer cancel()
	if relErr := actions_model.ReleaseTaskForRunner(ctx, t); relErr != nil {
		log.Error("ReleaseTaskForRunner [task_id: %d]: %v", t.ID, relErr)
	}
}

func PickTask(ctx context.Context, runner *actions_model.ActionRunner) (*runnerv1.Task, bool, error) {
	var (
		task       *runnerv1.Task
		// The job was already claimed but assembling its payload failed; release the
		// claim so the job returns to the waiting queue instead of being stranded in
		// running state with no runner ever executing it.
		releaseTaskForRunnerCleanup(t)
		return nil, false, err
	}
	actionTask = t
	// The job is claimed and its payload assembled, but if the request context was cancelled meanwhile, response can no longer reach the runner.
	// Release the claim so another runner can pick the job up.
	if err := ctx.Err(); err != nil {
		releaseTaskForRunnerCleanup(t)
		return nil, false, err
	}

	"testing"

	actions_model "gitea.dev/models/actions"
	"gitea.dev/models/db"
	"gitea.dev/models/unittest"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	assert.False(t, ok)
	assert.True(t, throttled)
}

// TestReleaseTaskForRunnerCleanup verifies the cleanup used by PickTask releases a claimed task through a
// fresh context. PickTask reaches this path when the request context is already canceled, and on
// PostgreSQL/MySQL a DB transaction on a canceled context fails immediately; reusing it would strand the
// claimed job in running state, so the cleanup must not use the request context.
func TestReleaseTaskForRunnerCleanup(t *testing.T) {
	require.NoError(t, unittest.PrepareTestDatabase())

	run := &actions_model.ActionRun{
		Title: "cleanup-run", RepoID: 1, OwnerID: 2, WorkflowID: "test.yaml",
		TriggerUserID: 2, Ref: "refs/heads/main",
		CommitSHA: "c2d72f548424103f01ee1dc02889c1e2bff816b0", Event: "push", TriggerEvent: "push",
		Status: actions_model.StatusWaiting,
	}
	require.NoError(t, db.Insert(t.Context(), run))
	job := &actions_model.ActionRunJob{
		RunID: run.ID, RepoID: run.RepoID, OwnerID: run.OwnerID, CommitSHA: run.CommitSHA,
		Name: "cleanup-job", Attempt: 1, JobID: "cleanup-job", Status: actions_model.StatusWaiting,
		RunsOn:          []string{"ubuntu-latest"},
		WorkflowPayload: []byte("on: push\njobs:\n  cleanup-job:\n    runs-on: ubuntu-latest\n    steps:\n      - run: echo hi\n"),
	}
	require.NoError(t, db.Insert(t.Context(), job))
	runner := &actions_model.ActionRunner{Name: "cleanup-runner", AgentLabels: []string{"ubuntu-latest"}}
	runner.GenerateAndFillToken()
	require.NoError(t, db.Insert(t.Context(), runner))

	task, ok, err := actions_model.CreateTaskForRunner(t.Context(), runner)
	require.NoError(t, err)
	require.True(t, ok)
	require.Equal(t, actions_model.StatusRunning, unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunJob{ID: job.ID}).Status)

	// the cleanup helper uses its own context, so the claimed job is returned to the waiting queue
	releaseTaskForRunnerCleanup(task)
	released := unittest.AssertExistsAndLoadBean(t, &actions_model.ActionRunJob{ID: job.ID})
	assert.Equal(t, actions_model.StatusWaiting, released.Status)
	assert.Zero(t, released.TaskID)
	unittest.AssertNotExistsBean(t, &actions_model.ActionTask{ID: task.ID})
}
	if subtle.ConstantTimeCompare([]byte(t.TokenHash), []byte(hex.EncodeToString(hashedToken[:]))) == 0 {
		// If an attacker steals a token and uses the token to create a new session the hash gets updated.
		// When the victim uses the old token the hashes don't match anymore and the victim should be notified about the compromised token.
		// Revoke the token so the attacker's rotated token (which shares this ID) can no longer be used.
		if err := auth_model.DeleteAuthTokenByID(ctx, t.ID); err != nil {
			return nil, err
		}
		return nil, ErrAuthTokenInvalidHash
	}

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	"gitea.dev/modules/timeutil"
	"gitea.dev/modules/util"

	"github.com/stretchr/testify/assert"
)
		assert.ErrorIs(t, err, ErrAuthTokenInvalidHash)
		assert.Nil(t, at2)

		// a hash mismatch signals a compromised token, which must be revoked
		_, err = auth_model.GetAuthTokenByID(t.Context(), at.ID)
		assert.ErrorIs(t, err, util.ErrNotExist)
	})

	t.Run("Valid", func(t *testing.T) {
}

// TokenCanAccessRepo reports whether the current API token is allowed to access the repository.
// A public-only token cannot reach a private repo or a repo owned by a non-public (limited or
// private) owner; any other token is unrestricted by this check.
func (ctx *APIContext) TokenCanAccessRepo(repo *repo_model.Repository) bool {
	return !ctx.PublicOnly || !publicOnlyTokenDeniedRepo(ctx, repo)
}

func init() {
package context

import (
	"context"
	"net/http"
	"slices"

	"gitea.dev/models/unit"
)

// isOwnerHidden reports whether repo's owner is not publicly visible (a limited or private owner), so
// the owner's repositories must be hidden from callers that may only reach genuinely public resources.
func isOwnerHidden(ctx context.Context, repo *repo_model.Repository) bool {
	if err := repo.LoadOwner(ctx); err != nil || repo.Owner == nil {
		return true // fail closed if the owner visibility can't be determined
	}
	return !repo.Owner.Visibility.IsPublic()
}

// publicOnlyTokenDeniedRepo reports whether a public-only API token must be denied access to
// repo. A public-only token may only reach genuinely public resources, so it is denied for
// private repos and for repos owned by a non-public (limited or private) owner.
func publicOnlyTokenDeniedRepo(ctx context.Context, repo *repo_model.Repository) bool {
	if repo == nil {
		return false
	}
	return repo.IsPrivate || isOwnerHidden(ctx, repo)
}

// TokenIsPublicOnly reports whether the request is authenticated by a public-only API token. A
// non-token request, or a token with no recorded scope, is not public-only.
func TokenIsPublicOnly(ctx *Context) bool {
	if ctx.Data["IsApiToken"] != true {
		return false
	}
	scope, ok := ctx.Data["ApiTokenScope"].(auth_model.AccessTokenScope)
	if !ok {
		return false
	}
	publicOnly, _ := scope.PublicOnly()
	return publicOnly
}

// CheckTokenScopes checks whether the authenticated API token contains any of the given scopes.
func CheckTokenScopes(ctx *Context, repo *repo_model.Repository, scopes ...auth_model.AccessTokenScope) {
	if ctx.Data["IsApiToken"] != true {
		return
	}

	if publicOnly && publicOnlyTokenDeniedRepo(ctx, repo) {
		ctx.HTTPError(http.StatusForbidden)
		return
	}
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
						}
						defer rc.Close()
					} else {
						// use the migration client so the fetch (including any redirect) is
						// validated against the migration host allow/block list
						resp, err := getMigrationHTTPClient().Get(*asset.DownloadURL)
						if err != nil {
							return err
						}
		}

		// SECURITY: We will assume that the pr.PatchURL has been checked
		// pr.PatchURL maybe a local file - but note EnsureSafe should be asserting that this safe.
		// Use the migration client so an http(s) PatchURL (and any redirect it follows) is
		// validated against the migration host allow/block list at dial time.
		resp, err := getMigrationHTTPClient().Get(u)
		if err != nil {
			return err
		}
		baseURL,
		gitea_sdk.SetToken(token),
		gitea_sdk.SetBasicAuth(username, password),
		gitea_sdk.SetHTTPClient(newMigrationHTTPClient()),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create NewGiteaDownloader for: %s. Error: %v", baseURL, err))
		Created:         rel.CreatedAt,
	}

	httpClient := newMigrationHTTPClient()

	for _, asset := range rel.Attachments {
		assetID := asset.ID // Don't optimize this, for closure we need a local variable
						return err
					}
				} else if asset.DownloadURL != nil {
					// use the migration client so the fetch (including any redirect) is
					// validated against the migration host allow/block list
					rc, err = uri.OpenWithClient(*asset.DownloadURL, getMigrationHTTPClient())
					if err != nil {
						return err
					}
		}

		// SECURITY: We will assume that the pr.PatchURL has been checked
		// pr.PatchURL maybe a local file - but note EnsureSafe should be asserting that this safe.
		// Use the migration client so an http(s) PatchURL (and any redirect it follows) is
		// validated against the migration host allow/block list at dial time.
		ret, err := uri.OpenWithClient(pr.PatchURL, getMigrationHTTPClient())
		if err != nil {
			return err
		}
		r.Published = rel.PublishedAt.Time
	}

	httpClient := newMigrationHTTPClient()

	for _, asset := range rel.Assets {
		assetID := asset.GetID() // Don't optimize this, for closure we need a local variable TODO: no need to do so in new Golang
//	Use either a username/password, personal token entered into the username field, or anonymous/public access
//	Note: Public access only allows very basic access
func NewGitlabDownloader(ctx context.Context, baseURL, repoPath, token string) (*GitlabDownloader, error) {
	gitlabClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL), gitlab.WithHTTPClient(newMigrationHTTPClient()))
	if err != nil {
		log.Trace("Error logging into gitlab: %v", err)
		return nil, err
		PublisherName:   rel.Author.Username,
	}

	httpClient := newMigrationHTTPClient()

	for _, asset := range rel.Assets.Links {
		assetID := asset.ID // Don't optimize this, for closure we need a local variable
	"gitea.dev/modules/hostmatcher"
	"gitea.dev/modules/proxy"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/util"
)

// migrationHTTPClient is the shared migration client. Callers that would otherwise build a client per
// request use it (via getMigrationHTTPClient) so a single connection pool is reused across downloads —
// e.g. many release assets from the same host — instead of a fresh pool and TLS handshake each time. It
// is built lazily on first use and reset by Init whenever the allow/block lists change; OnceValue keeps
// concurrent callers sharing a single client instead of racing to create their own.
var migrationHTTPClient = util.OnceValue[*http.Client]{Func: newMigrationHTTPClient}

// newMigrationHTTPClient returns a HTTP client for migration
func newMigrationHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewMigrationHTTPTransport(),
	}
}

// getMigrationHTTPClient returns the shared migration client, building it on first use so no request
// escapes the SSRF-validated transport even before Init has run.
func getMigrationHTTPClient() *http.Client {
	return migrationHTTPClient.Value()
}

// NewMigrationHTTPTransport returns a HTTP transport for migration. The target is validated against the
// allow/block lists on both the direct-dial and proxy paths, so a configured proxy cannot be used to
// reach an otherwise-forbidden target (SSRF).
func NewMigrationHTTPTransport() *http.Transport {
	return hostmatcher.NewHTTPTransport("migration", allowList, blockList, proxy.Proxy(), setting.Proxy.ProxyURLFixed,
		&tls.Config{InsecureSkipVerify: setting.Migrations.SkipTLSVerify})
}
		blockList.AppendBuiltin(hostmatcher.MatchBuiltinLoopback)
	}

	// reset the shared client so it is rebuilt from the freshly parsed lists on next use; download paths
	// then reuse one connection pool instead of creating a client (and pool) per request
	migrationHTTPClient.Reset()

	return nil
}
		log.Error("SyncMirrors [repo: %-v]: GetRemoteURL Error %v", m.Repo, remoteErr)
		return nil, false
	}
	// re-validate on every sync: the host may now resolve to an internal IP (rebinding) or the
	// allow/block list may have changed. ssh/file are skipped (not an HTTP SSRF vector).
	switch remoteURL.URL.Scheme {
	case "http", "https", "git":
		if allowErr := migrations.IsMigrateURLAllowed(remoteURL.String(), m.Repo.MustOwner(ctx)); allowErr != nil {
			log.Error("SyncMirrors [repo: %-v]: remote URL is not allowed: %v", m.Repo, allowErr)
			return nil, false
		}
	}
	envs := proxy.EnvWithProxy(remoteURL.URL)
	timeout := time.Duration(setting.Git.Timeout.Mirror) * time.Second

			return err
		}

		// the repo is no longer publicly visible, so drop stars and watches from users who can no longer
		// see it, matching the direct repository-private transition (see services/repository)
		if err := repo_model.ClearRepoStars(ctx, repo.ID); err != nil {
			return err
		}
		if err := repo_model.ClearRepoWatches(ctx, repo.ID); err != nil {
			return err
		}
	}

	// Create/Remove git-daemon-export-ok for git-daemon...
		require.NoError(t, ChangeOrganizationVisibility(t.Context(), org, structs.VisibleTypePrivate))
		unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: org.ID, Visibility: structs.VisibleTypePrivate})
	})

	t.Run("ChangeVisibilityClearsWatchesAndStars", func(t *testing.T) {
		// org3 is a public organization owning the public repo32
		org := unittest.AssertExistsAndLoadBean(t, &organization.Organization{ID: 3})
		require.Equal(t, structs.VisibleTypePublic, org.Visibility)
		repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 32, OwnerID: org.ID})

		// an outside user watches and stars the repo while the org is still visible
		watcher := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
		require.NoError(t, repo_model.WatchRepo(t.Context(), watcher, repo, true))
		require.NoError(t, repo_model.StarRepo(t.Context(), watcher, repo, true))
		unittest.AssertExistsAndLoadBean(t, &repo_model.Watch{UserID: watcher.ID, RepoID: repo.ID})

		require.NoError(t, ChangeOrganizationVisibility(t.Context(), org, structs.VisibleTypePrivate))

		// making the org private must drop watches, not only stars, from users who can no longer see it
		unittest.AssertNotExistsBean(t, &repo_model.Watch{UserID: watcher.ID, RepoID: repo.ID})
		unittest.AssertNotExistsBean(t, &repo_model.Star{UID: watcher.ID, RepoID: repo.ID})
	})
}
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	org3 := unittest.AssertExistsAndLoadBean(t, &org_model.Organization{ID: 3})
	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})
	// user15 is on org3's team7 (write access to the public repo32 only), so it can see org3 public repos
	// but has no access to the private repo3 — a genuine "no permission to the private repo" org member.
	user15 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 15})

	t.Run("User projects", func(t *testing.T) {
		pi1 := project_model.ProjectIssue{
		})

		t.Run("Authenticated user with no permission to the private repo", func(t *testing.T) {
			// user2 is on org3's Owners team and has owner access to the private repo3, so it is not a
			// valid "no permission" subject; user15 has no access to repo3 but can see the public repo32.
			columnIssues, err := LoadIssuesFromProject(t.Context(), projects[0], &issues_model.IssuesOptions{
				Owner: org3.AsUser(),
				Doer:  user15,
			})
			assert.NoError(t, err)
			assert.Len(t, columnIssues, 1)
			assert.Len(t, columnIssues[defaultColumn.ID], 1) // user15 can only visit public repo issues
		})

		t.Run("Org owner team member", func(t *testing.T) {
			// user2 is on org3's Owners team, so it has access to the private repo3 and must see both the
			// public and the private issue — the owner-team access that team.authorize grants at runtime.
			columnIssues, err := LoadIssuesFromProject(t.Context(), projects[0], &issues_model.IssuesOptions{
				Owner: org3.AsUser(),
				Doer:  user2,
			})
			assert.NoError(t, err)
			assert.Len(t, columnIssues, 1)
			assert.Len(t, columnIssues[defaultColumn.ID], 2) // owner-team member visits both public and private issues
		})
	})

			if err = repo_model.ClearRepoStars(ctx, repo.ID); err != nil {
				return err
			}

			// watchers who lost access must not keep watching the now-private repo
			if err = repo_model.ClearRepoWatches(ctx, repo.ID); err != nil {
				return err
			}
		}

		// Create/Remove git-daemon-export-ok for git-daemon...
	assert.True(t, updatedRepo.IsPrivate)
	assert.Zero(t, updatedRepo.NumWatches)
}

// TestUpdateRepositoryClearsWatchesOnVisibilityChange ensures the shared updateRepository
// helper (used by the API EditRepo path) also clears watches when a repo goes private.
func TestUpdateRepositoryClearsWatchesOnVisibilityChange(t *testing.T) {
	assert.NoError(t, unittest.PrepareTestDatabase())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	assert.False(t, repo.IsPrivate)

	watchers, err := repo_model.GetRepoWatchersIDs(t.Context(), repo.ID)
	require.NoError(t, err)
	require.NotEmpty(t, watchers)

	repo.IsPrivate = true
	require.NoError(t, updateRepository(t.Context(), repo, true))

	watchers, err = repo_model.GetRepoWatchersIDs(t.Context(), repo.ID)
	assert.NoError(t, err)
	assert.Empty(t, watchers)

	updatedRepo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: repo.ID})
	assert.Zero(t, updatedRepo.NumWatches)
}
// Init starts the hooks delivery thread
func Init() error {
	timeout := time.Duration(setting.Webhook.DeliverTimeout) * time.Second
	allowedHostMatcher := hostmatcher.ParseHostMatchList("security.ALLOWED_HOST_LIST", setting.Webhook.AllowedHostList)

	// NewHTTPTransport enforces the allow-list on direct connections; when webhookProxy routes a request
	// through a configured proxy, restricting the proxied target is the proxy server's responsibility.
	webhookHTTPClient = &http.Client{
		Timeout: timeout,
		Transport: hostmatcher.NewHTTPTransport("webhook", allowedHostMatcher, nil, webhookProxy(allowedHostMatcher), setting.Webhook.ProxyURLFixed,
			&tls.Config{InsecureSkipVerify: setting.Webhook.SkipTLSVerify}),
	}

	hookQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "webhook_sender", handler)
		config["color"] = s.Color
	}

	return &api.Hook{
		ID:     w.ID,
		Name:   w.Name,
		Type:   w.Type,
		URL:    fmt.Sprintf("%s/settings/hooks/%d", repoLink, w.ID),
		Active: w.IsActive,
		Config: config,
		Events: w.EventsArray(),
		// the stored authorization header is a secret and must never be returned by the API,
		// consistent with the webhook secret which is also omitted from the response
		AuthorizationHeader: "",
		Updated:             w.UpdatedUnix.AsTime(),
		Created:             w.CreatedUnix.AsTime(),
		BranchFilter:        w.BranchFilter,
)

func TestMain(m *testing.M) {
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			// for tests, allow only loopback IPs. This must run after the test config is loaded (which
			// resets the shared Security.AllowedHostList) and before Init() builds the delivery client.
			setting.Security.AllowedHostList = hostmatcher.MatchBuiltinLoopback
			setting.LoadQueueSettings()
			return Init()
		},
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/migration"
	"gitea.dev/modules/setting"
	api "gitea.dev/modules/structs"
	"gitea.dev/modules/test"
	migrations_service "gitea.dev/services/migrations"
	mirror_service "gitea.dev/services/mirror"
	repo_service "gitea.dev/services/repository"
	files_service "gitea.dev/services/repository/files"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestScheduleUpdate(t *testing.T) {

func testScheduleUpdateMirrorSync(t *testing.T) {
	doTestScheduleUpdate(t, func(t *testing.T, u *url.URL, testContext APITestContext, user *user_model.User, repo *repo_model.Repository) (commitID, expectedSpec string) {
		// the mirror sync re-validates the remote URL, which rejects the local test server unless local
		// networks are allowed; migrations.Init rebuilds the host allow-list from the setting
		defer test.MockVariableValue(&setting.Migrations.AllowLocalNetworks, true)()
		require.NoError(t, migrations_service.Init())

		// create mirror repo
		opts := migration.MigrateOptions{
			RepoName:    "actions-schedule-mirror",
		assert.NotNil(t, run.TriggerActor, "trigger_actor should be populated")
	}
}

// TestAPIOrgActionsRunsAccessControl ensures the org-level Actions run/job listing does not
// leak runs/jobs from repos the caller cannot access.
func TestAPIOrgActionsRunsAccessControl(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// org3 has action run 802 (and its jobs) in the private repo5; user28 is an org3 member
	// (teams 12/13) with no access to repo5.
	token := getUserToken(t, "user28", auth_model.AccessTokenScopeReadOrganization)

	req := NewRequest(t, "GET", "/api/v1/orgs/org3/actions/runs").AddTokenAuth(token)
	resp := MakeRequest(t, req, http.StatusOK)
	runs := DecodeJSON(t, resp, &api.ActionWorkflowRunsResponse{})
	for _, r := range runs.Entries {
		assert.NotEqual(t, int64(802), r.ID, "must not leak a run from an inaccessible repo")
	}

	req = NewRequest(t, "GET", "/api/v1/orgs/org3/actions/jobs").AddTokenAuth(token)
	resp = MakeRequest(t, req, http.StatusOK)
	jobs := DecodeJSON(t, resp, &api.ActionWorkflowJobsResponse{})
	for _, j := range jobs.Entries {
		assert.NotEqual(t, int64(802), j.RunID, "must not leak a job from an inaccessible repo run")
	}

	// user1 is a site admin: it normally bypasses the per-repo access filter, but a public-only token
	// must stay confined to public repos, so the run/job in the private repo5 must not be listed.
	adminPublicOnly := getUserToken(t, "user1", auth_model.AccessTokenScopeReadOrganization, auth_model.AccessTokenScopePublicOnly)

	req = NewRequest(t, "GET", "/api/v1/orgs/org3/actions/runs").AddTokenAuth(adminPublicOnly)
	resp = MakeRequest(t, req, http.StatusOK)
	adminRuns := DecodeJSON(t, resp, &api.ActionWorkflowRunsResponse{})
	for _, r := range adminRuns.Entries {
		assert.NotEqual(t, int64(802), r.ID, "a public-only admin token must not list a private repo's run")
	}

	req = NewRequest(t, "GET", "/api/v1/orgs/org3/actions/jobs").AddTokenAuth(adminPublicOnly)
	resp = MakeRequest(t, req, http.StatusOK)
	adminJobs := DecodeJSON(t, resp, &api.ActionWorkflowJobsResponse{})
	for _, j := range adminJobs.Entries {
		assert.NotEqual(t, int64(802), j.RunID, "a public-only admin token must not list a private repo's job")
	}
}
	unittest.AssertExistsAndLoadBean(t, &issues_model.IssueLabel{IssueID: issue.ID, LabelID: 2})
}

// TestAPIDeleteIssueLabelCrossRepo ensures DeleteIssueLabel does not act on a label
// belonging to another repository, and that a foreign-but-existing label ID and a
// nonexistent label ID return the same status (no cross-repo enumeration oracle).
func TestAPIDeleteIssueLabelCrossRepo(t *testing.T) {
	assert.NoError(t, unittest.LoadFixtures())

	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	issue := unittest.AssertExistsAndLoadBean(t, &issues_model.Issue{RepoID: repo.ID})
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: repo.OwnerID})
	// label 5 exists but belongs to repo 10, not repo 1
	foreignLabel := unittest.AssertExistsAndLoadBean(t, &issues_model.Label{ID: 5})
	assert.NotEqual(t, repo.ID, foreignLabel.RepoID)

	session := loginUser(t, owner.Name)
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteIssue)
	base := fmt.Sprintf("/api/v1/repos/%s/%s/issues/%d/labels", repo.OwnerName, repo.Name, issue.Index)

	// a foreign-but-existing label ID must not be accepted (was 204 before the fix)
	req := NewRequest(t, "DELETE", fmt.Sprintf("%s/%d", base, foreignLabel.ID)).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)

	// a nonexistent label ID must return the SAME status, so no oracle exists (was 422 before the fix)
	req = NewRequest(t, "DELETE", fmt.Sprintf("%s/%d", base, 9999999)).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)

	// a label that belongs to the repo is still removable
	unittest.AssertExistsAndLoadBean(t, &issues_model.Label{ID: 2, RepoID: repo.ID})
	addReq := NewRequestWithJSON(t, "POST", base, &api.IssueLabelsOption{Labels: []any{2}}).AddTokenAuth(token)
	MakeRequest(t, addReq, http.StatusOK)
	req = NewRequest(t, "DELETE", fmt.Sprintf("%s/%d", base, 2)).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNoContent)
}

func TestAPIAddIssueLabelsWithLabelNames(t *testing.T) {
	assert.NoError(t, unittest.LoadFixtures())

	assert.EqualValues(t, 33, apiNewTime.Time)
	assert.Equal(t, user2.ID, apiNewTime.UserID)
	assert.EqualValues(t, 947688818, apiNewTime.Created.Unix())

	// adding time for a nonexistent user must return 404, not panic with a 500
	req = NewRequestWithJSON(t, "POST", urlStr, &api.AddTimeOption{
		Time: 33,
		User: "nonexistentuser",
	}).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusNotFound)
}
	neturl "net/url"
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/packages"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
		})
	})
}

// TestPackageGenericPublicOnlyTokenLimitedOwner ensures a public-only token cannot
// access packages owned by a limited-visibility owner (only genuinely public owners).
func TestPackageGenericPublicOnlyTokenLimitedOwner(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user33 has limited visibility (visible only to authenticated users, not public)
	owner := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 33})
	base := fmt.Sprintf("/api/packages/%s/generic/pkg/1.0.0", owner.Name)

	// upload a package into the limited owner's namespace
	req := NewRequestWithBody(t, "PUT", base+"/file.bin", bytes.NewReader([]byte{1, 2, 3})).
		AddBasicAuth(owner.Name)
	MakeRequest(t, req, http.StatusCreated)

	// a public-only read:package token (even the owner's own) must be refused
	publicOnlyToken := getUserToken(t, owner.Name, auth_model.AccessTokenScopeReadPackage, auth_model.AccessTokenScopePublicOnly)
	req = NewRequest(t, "GET", base+"/file.bin").AddTokenAuth(publicOnlyToken)
	MakeRequest(t, req, http.StatusForbidden)
	// same via the v1 package API surface (checkTokenPublicOnly)
	req = NewRequest(t, "GET", fmt.Sprintf("/api/v1/packages/%s/generic/pkg/1.0.0", owner.Name)).AddTokenAuth(publicOnlyToken)
	MakeRequest(t, req, http.StatusForbidden)

	// a normal read:package token still works, proving only public-only is restricted
	token := getUserToken(t, owner.Name, auth_model.AccessTokenScopeReadPackage)
	req = NewRequest(t, "GET", base+"/file.bin").AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
}
		test(t, http.StatusNotFound, packageTag2, "1.2")
		test(t, http.StatusOK, packageTag, packageVersion)
		test(t, http.StatusOK, packageTag2, packageVersion)

		// an oversized dist-tag body is rejected instead of being read unbounded
		oversized := strings.Repeat("a", 5*1024)
		req := NewRequestWithBody(t, "PUT", fmt.Sprintf("%s/%s", tagsRoot, packageTag), strings.NewReader(oversized)).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusRequestEntityTooLarge)
	})

	t.Run("ListTags", func(t *testing.T) {
		}
	}
}

// TestAPIRepoLimitedOwnerPublicOnly ensures a public-only token cannot reach a public repo
// owned by a limited-visibility owner (which is not reachable anonymously).
func TestAPIRepoLimitedOwnerPublicOnly(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// limited_org (limited visibility) owns the non-private repo public_repo_on_limited_org
	const url = "/api/v1/repos/limited_org/public_repo_on_limited_org"

	// a public-only token is confined to genuinely public resources
	publicOnlyToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadRepository, auth_model.AccessTokenScopePublicOnly)
	req := NewRequest(t, "GET", url).AddTokenAuth(publicOnlyToken)
	MakeRequest(t, req, http.StatusNotFound)

	// a normal token can still reach the limited owner's public repo
	token := getUserToken(t, "user2", auth_model.AccessTokenScopeReadRepository)
	req = NewRequest(t, "GET", url).AddTokenAuth(token)
	MakeRequest(t, req, http.StatusOK)
}

	apiHook := DecodeJSON(t, resp, &api.Hook{})
	assert.Equal(t, "http://example.com/", apiHook.Config["url"])
	// the stored authorization header is a secret and must never be returned by the API
	assert.Empty(t, apiHook.AuthorizationHeader)
	assert.Equal(t, "CI notifications", apiHook.Name)

	// a read-scoped token must not be able to read back the authorization header
	readToken := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeReadRepository)
	getReq := NewRequest(t, "GET", fmt.Sprintf("/api/v1/repos/%s/%s/hooks/%d", owner.Name, repo.Name, apiHook.ID)).
		AddTokenAuth(readToken)
	getResp := MakeRequest(t, getReq, http.StatusOK)
	assert.NotContains(t, getResp.Body.String(), "s3cr3t")

	newName := "Deploy hook"
	patchReq := NewRequestWithJSON(t, "PATCH", fmt.Sprintf("/api/v1/repos/%s/%s/hooks/%d", owner.Name, repo.Name, apiHook.ID), api.EditHookOption{
		Name: &newName,
	deleteAPIAccessToken(t, newAccessToken, user)
}

// TestAPICreateTokenScopeEscalation ensures a token-authenticated request cannot
// mint a new token with a broader scope than the authenticating token.
func TestAPICreateTokenScopeEscalation(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})

	// a write:user-scoped token authenticates the create requests below
	writeUserToken := getUserToken(t, user.Name, auth_model.AccessTokenScopeWriteUser)

	// requesting a broader scope ("all") than the authenticating token is rejected
	req := NewRequestWithJSON(t, "POST", "/api/v1/users/"+user.LoginName+"/tokens", map[string]any{
		"name":   "escalated",
		"scopes": []string{"all"},
	})
	req.Request.SetBasicAuth(user.Name, writeUserToken)
	MakeRequest(t, req, http.StatusForbidden)

	// requesting a subset scope ("read:user") is allowed
	req = NewRequestWithJSON(t, "POST", "/api/v1/users/"+user.LoginName+"/tokens", map[string]any{
		"name":   "subset",
		"scopes": []string{"read:user"},
	})
	req.Request.SetBasicAuth(user.Name, writeUserToken)
	MakeRequest(t, req, http.StatusCreated)

	// password (non-token) auth may still create a token with any scope
	req = NewRequestWithJSON(t, "POST", "/api/v1/users/"+user.LoginName+"/tokens", map[string]any{
		"name":   "by-password",
		"scopes": []string{"all"},
	}).AddBasicAuth(user.Name)
	MakeRequest(t, req, http.StatusCreated)

	// a public-only token must not mint a token that drops the public-only restriction
	publicOnlyToken := getUserToken(t, user.Name, auth_model.AccessTokenScopeWriteUser, auth_model.AccessTokenScopePublicOnly)
	req = NewRequestWithJSON(t, "POST", "/api/v1/users/"+user.LoginName+"/tokens", map[string]any{
		"name":   "still-public-only",
		"scopes": []string{"write:user"},
	})
	req.Request.SetBasicAuth(user.Name, publicOnlyToken)
	resp := MakeRequest(t, req, http.StatusCreated)
	var createdToken api.AccessToken
	DecodeJSON(t, resp, &createdToken)
	assert.Contains(t, createdToken.Scopes, string(auth_model.AccessTokenScopePublicOnly))

	// an unrestricted parent token may create a narrower public-only child: public-only is a restriction,
	// not a grantable permission, so the subset check must not reject it
	req = NewRequestWithJSON(t, "POST", "/api/v1/users/"+user.LoginName+"/tokens", map[string]any{
		"name":   "narrower-public-only",
		"scopes": []string{"write:user", "public-only"},
	})
	req.Request.SetBasicAuth(user.Name, writeUserToken)
	resp = MakeRequest(t, req, http.StatusCreated)
	DecodeJSON(t, resp, &createdToken)
	assert.Contains(t, createdToken.Scopes, string(auth_model.AccessTokenScopePublicOnly))
}

// TestAPIDeleteMissingToken ensures that error is thrown when token not found
func TestAPIDeleteMissingToken(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/modules/setting"
	api "gitea.dev/modules/structs"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

// TestAPIManageEmailsFeatureDisabled ensures the email management API honors the
// manage_credentials feature restriction, matching the web UI.
func TestAPIManageEmailsFeatureDisabled(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	session := loginUser(t, "user2")
	token := getTokenForLoggedInUser(t, session, auth_model.AccessTokenScopeWriteUser)

	WithDisabledFeatures(t, setting.UserFeatureManageCredentials)

	addReq := NewRequestWithJSON(t, "POST", "/api/v1/user/emails", &api.CreateEmailOption{
		Emails: []string{"user2-3@example.com"},
	}).AddTokenAuth(token)
	MakeRequest(t, addReq, http.StatusNotFound)

	delReq := NewRequestWithJSON(t, "DELETE", "/api/v1/user/emails", &api.DeleteEmailOption{
		Emails: []string{"user2-2@example.com"},
	}).AddTokenAuth(token)
	MakeRequest(t, delReq, http.StatusNotFound)
}

func TestAPIListEmails(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	}
	require.NoError(t, user_model.LinkExternalToUser(t.Context(), u, extLink))

	prepareUserExternalLink := func(t *testing.T, accessToken, refreshToken string, expiresAt time.Time) {
		err := user_model.UpdateUserCols(t.Context(), &user_model.User{ID: u.ID, IsActive: false}, "is_active")
		require.NoError(t, err)
		_, err = db.GetEngine(t.Context()).Where(builder.Eq{"user_id": u.ID}).Cols("access_token", "refresh_token", "expires_at").
			Update(&user_model.ExternalLoginUser{AccessToken: accessToken, RefreshToken: refreshToken, ExpiresAt: expiresAt})
		require.NoError(t, err)
		require.False(t, unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID}).IsActive)
	}

	t.Run("admin-disabled user is not reactivated", func(t *testing.T) {
		prepareUserExternalLink(t, "an-access-token", "non-empty-refresh-token", time.Now().Add(time.Hour))
		doOIDCSignIn(t, authSource.Name)
		after := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID})
		assert.False(t, after.IsActive, "OAuth callback must not re-enable an administrator-disabled account")
	})

	t.Run("admin-disabled user without refresh token is not reactivated", func(t *testing.T) {
		// GitHub / OIDC-without-offline_access sources never store a refresh token, so a
		// stored access token (and no refresh token) is the normal admin-disabled state
		prepareUserExternalLink(t, "an-access-token", "" /* no refresh token */, time.Now().Add(time.Hour))
		doOIDCSignIn(t, authSource.Name)
		after := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID})
		assert.False(t, after.IsActive, "OAuth callback must not re-enable an admin-disabled account on a no-refresh-token source")
	})

	t.Run("auto-sync-disabled user is reactivated", func(t *testing.T) {
		// the auto-sync cron clears all three token fields when it disables a user
		prepareUserExternalLink(t, "", "", time.Time{})
		doOIDCSignIn(t, authSource.Name)
		after := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID})
		assert.True(t, after.IsActive, "OAuth callback must reactivate a sync-disabled account on successful login")
	"net/http"
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
		})
	})
}

// TestFeedUserPublicOnlyToken ensures a public-only API token cannot surface a user's
// private activity through their profile feed, even when authenticated as the owner.
func TestFeedUserPublicOnlyToken(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user2 has activity on the private repo user2/repo2
	const privateMarker = "user2/repo2"

	// a normal read:user token authenticated as the owner sees private activity
	fullToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadUser)
	reqFull := NewRequest(t, "GET", "/user2.rss")
	reqFull.SetBasicAuth("user2", fullToken)
	respFull := MakeRequest(t, reqFull, http.StatusOK)
	assert.Contains(t, respFull.Body.String(), privateMarker)

	// a public-only token must not surface the private activity
	publicOnlyToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadUser, auth_model.AccessTokenScopePublicOnly)
	reqPublicOnly := NewRequest(t, "GET", "/user2.rss")
	reqPublicOnly.SetBasicAuth("user2", publicOnlyToken)
	respPublicOnly := MakeRequest(t, reqPublicOnly, http.StatusOK)
	assert.NotContains(t, respPublicOnly.Body.String(), privateMarker)
}

// TestProfileActivityPublicOnlyToken ensures the HTML profile activity tab does not
// surface a user's private activity to a public-only API token, mirroring the RSS/Atom
// guard. The /{username} route is AllowBasic, so a public-only PAT used as the Basic
// password must still be downgraded even for the feed owner.
func TestProfileActivityPublicOnlyToken(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user2 has activity on the private repo user2/repo2
	const privateMarker = "user2/repo2"

	// a normal read:user token authenticated as the owner sees private activity
	fullToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadUser)
	reqFull := NewRequest(t, "GET", "/user2?tab=activity")
	reqFull.SetBasicAuth("user2", fullToken)
	respFull := MakeRequest(t, reqFull, http.StatusOK)
	assert.Contains(t, respFull.Body.String(), privateMarker)

	// a public-only token must not surface the private activity
	publicOnlyToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadUser, auth_model.AccessTokenScopePublicOnly)
	reqPublicOnly := NewRequest(t, "GET", "/user2?tab=activity")
	reqPublicOnly.SetBasicAuth("user2", publicOnlyToken)
	respPublicOnly := MakeRequest(t, reqPublicOnly, http.StatusOK)
	assert.NotContains(t, respPublicOnly.Body.String(), privateMarker)
}
	"net/http"
	"testing"

	auth_model "gitea.dev/models/auth"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unittest"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/test"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGoGet(t *testing.T) {

	assert.Equal(t, expected, resp.Body.String())
}

// TestGoGetPrivateRepoBranchNotLeaked ensures the go-get meta endpoint does not disclose a
// private repository's default branch name to unauthorized callers.
func TestGoGetPrivateRepoBranchNotLeaked(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user2/repo2 is private; give it a non-default branch name
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	require.True(t, repo.IsPrivate)
	repo.DefaultBranch = "secretbranch"
	require.NoError(t, repo_model.UpdateRepositoryColsNoAutoTime(t.Context(), repo, "default_branch"))

	// an unauthenticated caller must see the neutral instance default, not the real branch
	req := NewRequest(t, "GET", "/user2/repo2?go-get=1")
	resp := MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "/src/branch/master{/dir}")
	assert.NotContains(t, resp.Body.String(), "secretbranch")

	// the owner may still see the real default branch
	session := loginUser(t, "user2")
	req = NewRequest(t, "GET", "/user2/repo2?go-get=1")
	resp = session.MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "secretbranch")
}

// TestGoGetPublicRepoUnderLimitedOwnerBranchNotLeaked ensures the go-get meta endpoint does not disclose
// the default branch of a public repo whose owner is not visible to the caller (here a limited org, which
// is hidden from anonymous callers but visible to authenticated ones).
func TestGoGetPublicRepoUnderLimitedOwnerBranchNotLeaked(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// repo id 38 is a public repo owned by a limited org; give it a non-default branch name
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 38})
	require.False(t, repo.IsPrivate)
	require.NoError(t, repo.LoadOwner(t.Context()))
	require.False(t, repo.Owner.Visibility.IsPublic())
	repo.DefaultBranch = "secretbranch"
	require.NoError(t, repo_model.UpdateRepositoryColsNoAutoTime(t.Context(), repo, "default_branch"))

	url := fmt.Sprintf("/%s/%s?go-get=1", repo.OwnerName, repo.Name)

	// an anonymous caller cannot see the limited owner, so the neutral instance default is returned
	resp := MakeRequest(t, NewRequest(t, "GET", url), http.StatusOK)
	assert.Contains(t, resp.Body.String(), "/src/branch/master{/dir}")
	assert.NotContains(t, resp.Body.String(), "secretbranch")

	// an authenticated caller can see a limited org's public repo, so the real branch is shown
	session := loginUser(t, "user2")
	resp = session.MakeRequest(t, NewRequest(t, "GET", url), http.StatusOK)
	assert.Contains(t, resp.Body.String(), "secretbranch")
}

// TestGoGetPrivateRepoBranchNotLeakedToTokenWithoutRepoScope ensures a PAT that was not granted repository
// read scope cannot learn a private repo's default branch through go-get, even when the account behind the
// token could read the repository.
func TestGoGetPrivateRepoBranchNotLeakedToTokenWithoutRepoScope(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user2/repo2 is private; give it a non-default branch name
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 2})
	require.True(t, repo.IsPrivate)
	repo.DefaultBranch = "secretbranch"
	require.NoError(t, repo_model.UpdateRepositoryColsNoAutoTime(t.Context(), repo, "default_branch"))

	// a token scoped only to read:misc does not grant repository read: the branch must stay hidden.
	// Web routes authenticate a token via basic auth (username + token), which also records the scope.
	miscToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadMisc)
	req := NewRequest(t, "GET", "/user2/repo2?go-get=1")
	req.Request.SetBasicAuth("user2", miscToken)
	resp := MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "/src/branch/master{/dir}")
	assert.NotContains(t, resp.Body.String(), "secretbranch")

	// a token that includes repository read scope may see the real branch
	repoToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadRepository)
	req = NewRequest(t, "GET", "/user2/repo2?go-get=1")
	req.Request.SetBasicAuth("user2", repoToken)
	resp = MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "secretbranch")
}
// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	"gitea.dev/tests"
)

// TestUpdateIssueLabelForeignLabel verifies that the web issue-label action rejects a
// label owned by a different repo/org with 404 (indistinguishable from a nonexistent
// id), closing the cross-repo label enumeration oracle also on the web side.
func TestUpdateIssueLabelForeignLabel(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	sess := loginUser(t, "user2")

	// label 3 belongs to org3 — foreign to user2/repo1 (issue 1); must be 404, not a 500 oracle
	req := NewRequestWithValues(t, "POST", "/user2/repo1/issues/labels?issue_ids=1", map[string]string{
		"action": "attach",
		"id":     "3",
	})
	sess.MakeRequest(t, req, http.StatusNotFound)

	// a label owned by the repo still attaches normally
	req = NewRequestWithValues(t, "POST", "/user2/repo1/issues/labels?issue_ids=1", map[string]string{
		"action": "attach",
		"id":     "2",
	})
	sess.MakeRequest(t, req, http.StatusOK)
}
package integration

import (
	"net/http"
	"net/http/httptest"
	"slices"
	"sync/atomic"
	"testing"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/gitrepo"
	"gitea.dev/modules/migration"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/test"
	migrations "gitea.dev/services/migrations"
	mirror_service "gitea.dev/services/mirror"
	release_service "gitea.dev/services/release"
	repo_service "gitea.dev/services/repository"
	mirror = unittest.AssertExistsAndLoadBean(t, &repo_model.Mirror{RepoID: mirrorRepo.ID})
	assert.Equal(t, lastMirrorSync, mirror.LastSyncUnix)
}

// TestMirrorPullSSRFRevalidation ensures a pull mirror re-validates its remote URL against
// the migration allow/block list on every sync, so a mirror whose (network) remote now
// points at a disallowed internal host is never fetched.
func TestMirrorPullSSRFRevalidation(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	ctx := t.Context()
	user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	repo := unittest.AssertExistsAndLoadBean(t, &repo_model.Repository{ID: 1})
	repoPath := repo_model.RepoPath(user.Name, repo.Name)

	// an "internal" server that records whether it was reached
	var reached atomic.Bool
	internal := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		reached.Store(true)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer internal.Close()

	mirrorRepo, err := repo_service.CreateRepositoryDirectly(ctx, user, user, repo_service.CreateRepoOptions{
		Name:     "ssrf_mirror",
		IsMirror: true,
		Status:   repo_model.RepositoryBeingMigrated,
	}, false)
	require.NoError(t, err)
	_, err = repo_service.MigrateRepositoryGitData(ctx, user, mirrorRepo, migration.MigrateOptions{
		RepoName:  "ssrf_mirror",
		Mirror:    true,
		CloneAddr: repoPath,
	}, nil)
	require.NoError(t, err)

	mirror, err := repo_model.GetMirrorByRepoID(ctx, mirrorRepo.ID)
	require.NoError(t, err)

	// repoint the mirror at the loopback server, which is disallowed once local networks are off
	require.NoError(t, mirror_service.UpdateAddress(ctx, mirror, internal.URL+"/repo.git"))
	defer test.MockVariableValue(&setting.Migrations.AllowLocalNetworks, false)()
	require.NoError(t, migrations.Init())
	t.Cleanup(func() { _ = migrations.Init() })

	assert.False(t, mirror_service.SyncPullMirror(ctx, mirrorRepo.ID))
	assert.False(t, reached.Load(), "the disallowed internal remote must not be reached")
}
func TestOAuth2AvatarFromPicture(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.OAuth2Client.UpdateAvatar, true)()
	// the SSRF-protected avatar client blocks loopback by default; allow the loopback mock server here
	defer test.MockVariableValue(&setting.Security.AllowedHostList, "loopback,external")()

	mockServer := createOAuth2MockProvider()
	defer mockServer.Close()
	testOAuth2(t, "/user/oauth2/test%2Bplus", http.StatusTemporaryRedirect)
	testOAuth2(t, "/user/oauth2/test%20plus", http.StatusNotFound)
}

// TestOAuthUserInfoTokenScope verifies the OIDC userinfo endpoint enforces the
// read:user token scope, so a restrictively-scoped token cannot read identity claims.
func TestOAuthUserInfoTokenScope(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// a token without the user scope must be rejected
	miscToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadMisc)
	req := NewRequest(t, "GET", "/login/oauth/userinfo")
	req.SetHeader("Authorization", "Bearer "+miscToken)
	MakeRequest(t, req, http.StatusForbidden)

	// a token with read:user is allowed and returns the identity claims
	userToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadUser)
	req = NewRequest(t, "GET", "/login/oauth/userinfo")
	req.SetHeader("Authorization", "Bearer "+userToken)
	resp := MakeRequest(t, req, http.StatusOK)
	assert.Contains(t, resp.Body.String(), "user2@example.com")
}
	"path"
	"strings"
	"testing"
	"time"

	"gitea.dev/models/db"
	issues_model "gitea.dev/models/issues"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/git"
	"gitea.dev/modules/gitrepo"
	"gitea.dev/modules/test"
	issue_service "gitea.dev/services/issue"
	repo_service "gitea.dev/services/repository"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPullView_ReviewerMissed(t *testing.T) {
			unittest.AssertExistsAndLoadBean(t, &issues_model.Review{IssueID: pr.IssueID, Type: issues_model.ReviewTypeRequest, ReviewerID: 5})
			assert.NoError(t, pr.LoadIssue(t.Context()))

			// capture the current PR head ref so we can wait for the async
			// refs/pull/N/head sync triggered by the next push to complete
			baseGitRepo, err := gitrepo.OpenRepository(t.Context(), repo)
			require.NoError(t, err)
			defer baseGitRepo.Close()
			headRefBefore, err := baseGitRepo.GetRefCommitID(pr.GetGitHeadRefName())
			require.NoError(t, err)

			// update the file on the pr branch
			_, err = files_service.ChangeRepoFiles(t.Context(), repo, user2, &files_service.ChangeRepoFilesOptions{
				OldBranch: "codeowner-basebranch",
			})
			assert.NoError(t, err)

			// refs/pull/N/head is refreshed asynchronously by the push hook; wait for
			// it before evaluating code owners, otherwise the changed-file set may not
			// yet include user8-file.md and the review request would be missed
			require.Eventually(t, func() bool {
				headRefAfter, err := baseGitRepo.GetRefCommitID(pr.GetGitHeadRefName())
				return err == nil && headRefAfter != headRefBefore
			}, 30*time.Second, 100*time.Millisecond)

			reviewNotifiers, err := issue_service.PullRequestCodeOwnersReview(t.Context(), pr)
			require.NoError(t, err)
			require.Len(t, reviewNotifiers, 1)
			assert.EqualValues(t, 8, reviewNotifiers[0].Reviewer.ID)

			err = issue_service.ChangeTitle(t.Context(), pr.Issue, user2, "[WIP] Test Pull Request")
	})
}

// TestAPIPullUpdatePublicOnlyToken verifies that a public-only API token cannot
// update (push into) a PR whose head repo is private, even when the base repo
// named in the route is public.
func TestAPIPullUpdatePublicOnlyToken(t *testing.T) {
	onGiteaRun(t, func(t *testing.T, giteaURL *url.URL) {
		user := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
		org26 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 26})
		pr := createOutdatedPR(t, user, org26)
		require.NoError(t, pr.LoadBaseRepo(t.Context()))
		require.NoError(t, pr.LoadHeadRepo(t.Context()))
		require.NoError(t, pr.LoadIssue(t.Context()))

		// make the head repo private while the base repo stays public
		require.NoError(t, repo_model.UpdateRepositoryColsNoAutoTime(t.Context(),
			&repo_model.Repository{ID: pr.HeadRepo.ID, IsPrivate: true}, "is_private"))

		// a public-only write token must be refused (404), not perform the push
		publicOnlyToken := getUserToken(t, user.Name, auth_model.AccessTokenScopeWriteRepository, auth_model.AccessTokenScopePublicOnly)
		req := NewRequestf(t, "POST", "/api/v1/repos/%s/%s/pulls/%d/update", pr.BaseRepo.OwnerName, pr.BaseRepo.Name, pr.Issue.Index).
			AddTokenAuth(publicOnlyToken)
		MakeRequest(t, req, http.StatusNotFound)

		// the head branch must still be outdated (no push happened)
		diffCount, err := gitrepo.GetDivergingCommits(t.Context(), pr.BaseRepo, pr.BaseBranch, pr.GetGitHeadRefName())
		require.NoError(t, err)
		assert.Equal(t, 1, diffCount.Behind)

		// a normal write token still works, proving the guard is scope-specific
		token := getUserToken(t, user.Name, auth_model.AccessTokenScopeWriteRepository)
		req = NewRequestf(t, "POST", "/api/v1/repos/%s/%s/pulls/%d/update", pr.BaseRepo.OwnerName, pr.BaseRepo.Name, pr.Issue.Index).
			AddTokenAuth(token)
		MakeRequest(t, req, http.StatusOK)
	})
}

func updateRepoPullRequestConfig(t *testing.T, repoID int64, update func(*repo_model.PullRequestsConfig)) {
	t.Helper()

// Copyright 2026 The Gitea Authors. All rights reserved.
// SPDX-License-Identifier: MIT

package integration

import (
	"net/http"
	"testing"

	auth_model "gitea.dev/models/auth"
	"gitea.dev/tests"
)

// TestRepoHomeContentTokenScopes ensures the web repository home page enforces the
// repository read scope (and public-only confinement) of an API token used via basic
// auth, so a wrongly-scoped token cannot read private repository content.
func TestRepoHomeContentTokenScopes(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	// user2/repo2 is a private repository owned by user2
	const url = "/user2/repo2"

	// a token without repository scope must be denied
	miscToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadMisc)
	reqDenied := NewRequest(t, "GET", url)
	reqDenied.SetBasicAuth("user2", miscToken)
	MakeRequest(t, reqDenied, http.StatusForbidden)

	// a public-only token must be denied on a private repo
	publicOnlyToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadRepository, auth_model.AccessTokenScopePublicOnly)
	reqPublicOnly := NewRequest(t, "GET", url)
	reqPublicOnly.SetBasicAuth("user2", publicOnlyToken)
	MakeRequest(t, reqPublicOnly, http.StatusForbidden)

	// a token with repository read scope is allowed
	ownerReadToken := getUserToken(t, "user2", auth_model.AccessTokenScopeReadRepository)
	reqAllowed := NewRequest(t, "GET", url)
	reqAllowed.SetBasicAuth("user2", ownerReadToken)
	MakeRequest(t, reqAllowed, http.StatusOK)
}
