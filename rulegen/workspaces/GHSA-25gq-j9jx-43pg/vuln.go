package main

	UpdatedBefore    timeutil.TimeStamp
	ConcurrencyGroup string
	OrderBy          db.SearchOrderBy
}

var JobOrderByMap = map[string]map[string]db.SearchOrderBy{
		}
		cond = cond.And(builder.Eq{"`action_run_job`.concurrency_group": opts.ConcurrencyGroup})
	}
	return cond
}

	Status           []Status
	ConcurrencyGroup string
	CommitSHA        string
}

func (opts FindRunOptions) ToConds() builder.Cond {
	if opts.CommitSHA != "" {
		cond = cond.And(builder.Eq{"`action_run`.commit_sha": opts.CommitSHA})
	}
	return cond
}


const ssh2keyStart = "---- BEGIN SSH2 PUBLIC KEY ----"

func extractTypeFromBase64Key(key string) (string, error) {
	b, err := base64.StdEncoding.DecodeString(key)
	if err != nil || len(b) < 4 {
		return "", fmt.Errorf("invalid key format: %w", err)

// parseKeyString parses any key string in OpenSSH or SSH2 format to clean OpenSSH string (RFC4253).
func parseKeyString(content string) (string, error) {
	// remove whitespace at start and end
	content = strings.TrimSpace(content)

		// Transform all legal line endings to a single "\n".
		content = strings.NewReplacer("\r\n", "\n", "\r", "\n").Replace(content)

		lines := strings.Split(content, "\n")
		continuationLine := false

			if continuationLine || strings.ContainsAny(line, ":-") {
				continuationLine = strings.HasSuffix(line, "\\")
			} else {
				keyContent += line
			}
		}

		t, err := extractTypeFromBase64Key(keyContent)
		if err != nil {
	}
}

func Test_PublicKeysAreExternallyManaged(t *testing.T) {
	key1 := unittest.AssertExistsAndLoadBean(t, &PublicKey{ID: 1})
	externals, err := PublicKeysAreExternallyManaged(t.Context(), []*PublicKey{key1})
	assert.NoError(t, err)
	assert.Len(t, externals, 1)
	assert.False(t, externals[0])
}
	return bitmap.hasScope(AccessTokenScopePublicOnly)
}

// HasScope returns true if the string has the given scope
func (s AccessTokenScope) HasScope(scopes ...AccessTokenScope) (bool, error) {
	bitmap, err := s.parse()
		})
	}
}
		return nil
	}
	return db.WithTx(ctx, func(ctx context.Context) error {
		attachments, err := repo_model.GetAttachmentsByUUIDs(ctx, uuids)
		if err != nil {
			return fmt.Errorf("getAttachmentsByUUIDs [uuids: %v]: %w", uuids, err)
		}
		for i := range attachments {
			attachments[i].IssueID = c.IssueID
			attachments[i].CommentID = c.ID
			if err := repo_model.UpdateAttachment(ctx, attachments[i]); err != nil {
	assert.NoError(t, unittest.PrepareTestDatabase())

	comment := unittest.AssertExistsAndLoadBean(t, &issues_model.Comment{ID: 1})
	attachment := repo_model.Attachment{
		Name: "test.txt",
	}
	assert.NoError(t, db.Insert(t.Context(), &attachment))

	return err
}

// UpdateIssueAttachments update attachments by UUIDs for the issue
func UpdateIssueAttachments(ctx context.Context, issueID int64, uuids []string) (err error) {
	return db.WithTx(ctx, func(ctx context.Context) error {
		attachments, err := repo_model.GetAttachmentsByUUIDs(ctx, uuids)
		if err != nil {
			return fmt.Errorf("getAttachmentsByUUIDs [uuids: %v]: %w", uuids, err)
		}
		for i := range attachments {
			attachments[i].IssueID = issueID
			if err := repo_model.UpdateAttachment(ctx, attachments[i]); err != nil {
				return fmt.Errorf("update attachment [id: %d]: %w", attachments[i].ID, err)

import (
	"context"
	"fmt"
	"slices"
	"strconv"
	RepoID  int64
}

// IsErrRepoLabelNotExist checks if an error is a RepoErrLabelNotExist.
func IsErrRepoLabelNotExist(err error) bool {
	_, ok := err.(ErrRepoLabelNotExist)
	return ok
}

func (err ErrRepoLabelNotExist) Error() string {
	return fmt.Sprintf("label does not exist [label_id: %d, repo_id: %d]", err.LabelID, err.RepoID)
}
	return l, nil
}

// GetLabelInRepoByID returns a label by ID in given repository.
func GetLabelInRepoByID(ctx context.Context, repoID, labelID int64) (*Label, error) {
	if labelID <= 0 || repoID <= 0 {
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/timeutil"

	"github.com/stretchr/testify/assert"
)
	assert.Equal(t, "label1", label.Name)

	_, err = issues_model.GetLabelInRepoByName(t.Context(), 1, "")
	assert.True(t, issues_model.IsErrRepoLabelNotExist(err))

	_, err = issues_model.GetLabelInRepoByName(t.Context(), unittest.NonexistentID, "nonexistent")
	assert.True(t, issues_model.IsErrRepoLabelNotExist(err))
}

func TestGetLabelInRepoByNames(t *testing.T) {
	assert.EqualValues(t, 1, label.ID)

	_, err = issues_model.GetLabelInRepoByID(t.Context(), 1, -1)
	assert.True(t, issues_model.IsErrRepoLabelNotExist(err))

	_, err = issues_model.GetLabelInRepoByID(t.Context(), unittest.NonexistentID, unittest.NonexistentID)
	assert.True(t, issues_model.IsErrRepoLabelNotExist(err))
}

func TestGetLabelsInRepoByIDs(t *testing.T) {
}

// userOrgTeamUnitRepoBuilder returns repo ids where user's teams can access the special unit.
func userOrgTeamUnitRepoBuilder(userID int64, unitType unit.Type) *builder.Builder {
	return userOrgTeamRepoBuilder(userID).
		Join("INNER", "team_unit", "`team_unit`.team_id = `team_repo`.team_id").
		Where(builder.Eq{"`team_unit`.`type`": unitType}).
		And(builder.Gt{"`team_unit`.`access_mode`": int(perm.AccessModeNone)})
}

// userOrgTeamUnitRepoCond returns a condition to select repo ids where user's teams can access the special unit.
func UserOrgUnitRepoCond(idStr string, userID, orgID int64, unitType unit.Type) builder.Cond {
	return builder.In(idStr,
		userOrgTeamUnitRepoBuilder(userID, unitType).
			And(builder.Eq{"`team_unit`.org_id": orgID}),
	)
}

	))
}

// GetUserRepositories returns a list of repositories of given user.
func GetUserRepositories(ctx context.Context, opts SearchRepoOptions) (RepositoryList, int64, error) {
	if len(opts.OrderBy) == 0 {

	"gitea.dev/models/db"
	repo_model "gitea.dev/models/repo"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/optional"
		})
	}
}
		// only include private repos the actor can still access, so metadata does not leak after access revocation
		cond = cond.And(AccessibleRepositoryCondition(opts.Actor, unit.TypeInvalid))
	} else {
		cond = cond.And(builder.Eq{
			"repository.is_private": false,
		})
	}
	return cond
}
		// only include private repos the actor can still access, so metadata does not leak after access revocation
		cond = cond.And(AccessibleRepositoryCondition(opts.Actor, unit.TypeInvalid))
	} else {
		cond = cond.And(builder.Eq{
			"repository.is_private": false,
		})
	}
	return cond.And(builder.Neq{
		"watch.mode": WatchModeDont,
	require.Len(t, users, 1)
	assert.Equal(t, "user2", users[0].Name)
}
package openid

import (
	"time"

	"github.com/yohcop/openid-go"
)

var (
	nonceStore     = openid.NewSimpleNonceStore()
	discoveryCache = newTimedDiscoveryCache(24 * time.Hour)
)

// Verify handles response from OpenID provider
func Verify(fullURL string) (id string, err error) {
	return openid.Verify(fullURL, discoveryCache, nonceStore)
}

// Normalize normalizes an OpenID URI

// RedirectURL redirects browser
func RedirectURL(id, callbackURL, realm string) (string, error) {
	return openid.RedirectURL(id, callbackURL, realm)
}

import (
	"context"
	"fmt"
	"net"
	"net/url"
	"syscall"
	"time"
		return dialer.DialContext(ctx, network, addrOrHost)
	}
}
	"strconv"
	"strings"

	"gitea.dev/modules/util"
	"gitea.dev/modules/validation"

	namePattern = regexp.MustCompile(`\A[a-zA-Z0-9@._+-]+\z`)
	// (epoch:pkgver-pkgrel)
	versionPattern = regexp.MustCompile(`\A(?:\d:)?[\w.+~]+(?:-[-\w.+~]+)?\z`)
)

type Package struct {
	}

	var p *Package
	files := make([]string, 0, 10)

	tr := tar.NewReader(inner)
	for {
				return nil, err
			}
		} else if !strings.HasPrefix(filename, ".") {
			files = append(files, hd.Name)
		}
	}

		return nil, ErrMissingPKGINFOFile
	}

	p.FileMetadata.Files = files
	p.FileCompressionExtension = compressionType

	return p, nil
	"io"
	"testing"

	"github.com/klauspost/compress/zstd"
	"github.com/stretchr/testify/assert"
	"github.com/ulikunitz/xz"
		assert.ElementsMatch(t, []string{"usr/bin/paket1"}, p.FileMetadata.Backup)
	})
}
				return nil, GlobalVars().ErrUnsupportedCompression
			}

			tr := tar.NewReader(inner)
			for {
				hd, err := tr.Next()
				if err == io.EOF {
	key := ""
	var depends strings.Builder
	var control strings.Builder

	// https://www.debian.org/doc/debian-policy/ch-controlfields.html#syntax-of-control-files
	s := bufio.NewScanner(r)
		control.WriteString(line)
		control.WriteByte('\n')

		if line[0] == ' ' || line[0] == '\t' {
			switch key {
			case "Description":
				p.Metadata.Description += line
			case "Depends":
				depends.WriteString(trimmed)
			}
					p.Metadata.Maintainer = a.Name
				}
			case "Description":
				p.Metadata.Description = value
			case "Depends":
				depends.WriteString(value)
			case "Homepage":
		return nil, GlobalVars().ErrInvalidArchitecture
	}

	dependencies := strings.Split(depends.String(), ",")
	for i := range dependencies {
		dependencies[i] = strings.TrimSpace(dependencies[i])
		assert.True(t, IsValidDistributionOrComponent(name), "good=%q", name)
	}
}
	"crypto/tls"
	"net"
	"net/http"
	"os"
	"strings"
	"sync"
	return conn, nil
}

var internalAPITransport = sync.OnceValue(func() http.RoundTripper {
	return &http.Transport{
		DialContext: dialContextInternalAPI,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
			ServerName:         setting.Domain,
		},
	}
})
	XContentTypeOptions string

	ContentSecurityPolicyGeneral string // it only supports empty (default policy) or "unset", maybe it can support more in the future
}{
	XFrameOptions:       "SAMEORIGIN",
	XContentTypeOptions: "nosniff",
}

var (
)

func TestLoadSecurityFrom(t *testing.T) {
	cfg, err := NewConfigProviderFromData(`[security]
X_FRAME_OPTIONS = DENY
X_CONTENT_TYPE_OPTIONS = unset
CONTENT_SECURITY_POLICY_GENERAL = "script-src *; foo"`)
	assert.NoError(t, err)
	loadSecurityFrom(cfg)
	assert.Equal(t, "DENY", Security.XFrameOptions)
	assert.Equal(t, "unset", Security.XContentTypeOptions)
	assert.Equal(t, `"script-src *`, Security.ContentSecurityPolicyGeneral) // holy shit ini package bug
}
	Webhook.QueueLength = sec.Key("QUEUE_LENGTH").MustInt(1000)
	Webhook.DeliverTimeout = sec.Key("DELIVER_TIMEOUT").MustInt(5)
	Webhook.SkipTLSVerify = sec.Key("SKIP_TLS_VERIFY").MustBool()
	Webhook.AllowedHostList = sec.Key("ALLOWED_HOST_LIST").MustString("")
	Webhook.Types = []string{"gitea", "gogs", "slack", "discord", "dingtalk", "telegram", "msteams", "feishu", "matrix", "wechatwork", "packagist"}
	Webhook.PagingNum = sec.Key("PAGING_NUM").MustInt(10)
	Webhook.ProxyURL = sec.Key("PROXY_URL").MustString("")

// Open open a local file or a remote file
func Open(uriStr string) (io.ReadCloser, error) {
	u, err := url.Parse(uriStr)
	if err != nil {
		return nil, err
	}
	switch strings.ToLower(u.Scheme) {
	case "http", "https":
		f, err := http.Get(uriStr)
		if err != nil {
			return nil, err
		}
package uri

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadURI(t *testing.T) {
	assert.NoError(t, err)
	defer f.Close()
}
				}

				if publicOnly {
					if ctx.Package != nil && ctx.Package.Owner.Visibility.IsPrivate() {
						ctx.HTTPError(http.StatusForbidden, "reqToken", "token scope is limited to public packages")
						return
					}
func AddPackageTag(ctx *context.Context) {
	packageName := packageNameFromParams(ctx)

	body, err := io.ReadAll(ctx.Req.Body)
	if err != nil {
		apiError(ctx, http.StatusInternalServerError, err)
		return
	}
	version := strings.Trim(string(body), "\"") // is as "version" in the body

	pv, err := packages_model.GetVersionByNameAndVersion(ctx, ctx.Package.Owner.ID, packages_model.TypeNpm, packageName, version)
					return
				}
			case auth_model.AccessTokenScopeCategoryPackage:
				if ctx.Package != nil && ctx.Package.Owner.Visibility.IsPrivate() {
					ctx.APIError(http.StatusForbidden, "token scope is limited to public packages")
					return
				}
		return
	}

	label, err := issues_model.GetLabelByID(ctx, ctx.PathParamInt64("id"))
	if err != nil {
		if issues_model.IsErrLabelNotExist(err) {
			ctx.APIError(http.StatusUnprocessableEntity, err.Error())
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

			// allow only RepoAdmin, Admin and User to add time
			user, err = user_model.GetUserByName(ctx, form.User)
			if err != nil {
				ctx.APIErrorInternal(err)
			}
		}
	}
		l, err = issues_model.GetLabelInRepoByID(ctx, ctx.Repo.Repository.ID, intID)
	}
	if err != nil {
		if issues_model.IsErrRepoLabelNotExist(err) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

	form := web.GetForm(ctx).(*api.EditLabelOption)
	l, err := issues_model.GetLabelInRepoByID(ctx, ctx.Repo.Repository.ID, ctx.PathParamInt64("id"))
	if err != nil {
		if issues_model.IsErrRepoLabelNotExist(err) {
			ctx.APIErrorNotFound()
		} else {
			ctx.APIErrorInternal(err)
		}
		return
	}

		return
	}

	// keep API back-compat: when no style is given, default to "merge" rather than the repo's DefaultUpdateStyle,
	// so existing API clients keep getting a merge update.
	rebase := repo_model.UpdateStyle(ctx.FormString("style", string(repo_model.UpdateStyleMerge))) == repo_model.UpdateStyleRebase
	"gitea.dev/routers/api/v1/utils"
	"gitea.dev/services/context"
	"gitea.dev/services/convert"
)

// ListJobs lists jobs for api route validated ownerID and repoID
// ownerID == 0 and repoID == 0 means all jobs
// ownerID == 0 and repoID != 0 means all jobs for the given repo
		opts.Statuses = append(opts.Statuses, values...)
	}

	jobs, total, err := db.FindAndCount[actions_model.ActionRunJob](ctx, opts)
	if err != nil {
		ctx.APIErrorInternal(err)
	}
	excludePullRequests := ctx.FormBool("exclude_pull_requests")

	runs, total, err := db.FindAndCount[actions_model.ActionRun](ctx, opts)
	if err != nil {
		ctx.APIErrorInternal(err)
	}
	t.Scope = scope

	if err := auth_model.NewAccessToken(ctx, t); err != nil {
		ctx.APIErrorInternal(err)
		return
	"net/http"

	user_model "gitea.dev/models/user"
	api "gitea.dev/modules/structs"
	"gitea.dev/modules/web"
	"gitea.dev/services/context"
	//   "422":
	//     "$ref": "#/responses/validationError"

	form := web.GetForm(ctx).(*api.CreateEmailOption)
	if len(form.Emails) == 0 {
		ctx.APIError(http.StatusUnprocessableEntity, "Email list empty")
	//   "404":
	//     "$ref": "#/responses/notFound"

	form := web.GetForm(ctx).(*api.DeleteEmailOption)
	if len(form.Emails) == 0 {
		ctx.Status(http.StatusNoContent)
	user_model "gitea.dev/models/user"
	auth_module "gitea.dev/modules/auth"
	"gitea.dev/modules/container"
	"gitea.dev/modules/httplib"
	"gitea.dev/modules/log"
	"gitea.dev/modules/optional"
	"gitea.dev/modules/session"
	"gitea.dev/modules/setting"
	source_service "gitea.dev/services/auth/source"
	ctx.Redirect(setting.AppSubURL + "/user/link_account")
}

var oauth2AvatarHTTPClient = &http.Client{Timeout: 30 * time.Second}

func oauth2UpdateAvatarIfNeed(ctx *context.Context, avatarURL string, u *user_model.User) {
	if !setting.OAuth2Client.UpdateAvatar || len(avatarURL) == 0 {
	// Some hosts (e.g. Wikimedia) reject Go's default User-Agent.
	req.Header.Set("User-Agent", "Gitea "+setting.AppVer)

	resp, err := oauth2AvatarHTTPClient.Do(req)
	if err != nil {
		log.Warn("fetch %q failed: %v", avatarURL, err)
		return
			ctx.ServerError("GetExternalLogin", err)
			return
		}
		isDisabledByAutoSync := hasExt && extLogin.RefreshToken == ""
		if isDisabledByAutoSync {
			opts.IsActive = optional.Some(true)
		}
		return
	}

	response := &userInfoResponse{
		Sub:               strconv.FormatInt(ctx.Doer.ID, 10),
		Name:              ctx.Doer.DisplayName(),
package auth

import (
	"testing"

	"gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/services/oauth2_provider"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func createAndParseToken(t *testing.T, grant *auth.OAuth2Grant) *oauth2_provider.OIDCToken {
	assert.Equal(t, user.Email, oidcToken.Email)
	assert.Equal(t, user.IsActive, oidcToken.EmailVerified)
}
		includePrivate = isOrgMember
	}

	actions, _, err := feed_service.GetFeeds(ctx, activities_model.GetFeedsOptions{
		RequestedUser:   ctx.ContextUser,
		Actor:           ctx.Doer,
	"path"
	"strings"

	repo_model "gitea.dev/models/repo"
	"gitea.dev/modules/setting"
	"gitea.dev/modules/util"
	"gitea.dev/services/context"
		return
	}
	branchName := setting.Repository.DefaultBranch

	repo, err := repo_model.GetRepositoryByOwnerAndName(ctx, ownerName, repoName)
	if err == nil && len(repo.DefaultBranch) > 0 {
		branchName = repo.DefaultBranch
	}
	prefix := setting.AppURL + path.Join(url.PathEscape(ownerName), url.PathEscape(repoName), "src", "branch", util.PathEscapeSegments(branchName))

	ctx.RespHeader().Set("Content-Type", "text/html")
	_, _ = ctx.Write([]byte(res))
}
			}
		}
	case "attach", "detach", "toggle", "toggle-alt":
		label, err := issues_model.GetLabelByID(ctx, ctx.FormInt64("id"))
		if err != nil {
			if issues_model.IsErrRepoLabelNotExist(err) {
				ctx.HTTPError(http.StatusNotFound, "GetLabelByID")
			} else {
				ctx.ServerError("GetLabelByID", err)
	"strings"
	"time"

	"gitea.dev/models/db"
	git_model "gitea.dev/models/git"
	repo_model "gitea.dev/models/repo"
		return
	}

	// Check whether the repo is viewable: not in migration, and the code unit should be enabled
	// Ideally the "feed" logic should be after this, but old code did so, so keep it as-is.
	checkHomeCodeViewable(ctx)
		date := ctx.FormString("date")
		pagingNum = setting.UI.FeedPagingNum
		showPrivate := ctx.IsSigned && (ctx.Doer.IsAdmin || ctx.Doer.ID == ctx.ContextUser.ID)
		items, feedCount, err := feed_service.GetFeedsForDashboard(ctx, activities_model.GetFeedsOptions{
			RequestedUser:   ctx.ContextUser,
			Actor:           ctx.Doer,
		return
	}

	if err := auth_model.NewAccessToken(ctx, t); err != nil {
		ctx.ServerError("NewAccessToken", err)
		return
	"errors"
	"fmt"
	"sync"

	runnerv1 "gitea.dev/actions-proto-go/runner/v1"
	actions_model "gitea.dev/models/actions"
	"gitea.dev/models/db"
	secret_model "gitea.dev/models/secret"
	"gitea.dev/modules/log"
	"gitea.dev/modules/setting"

	return task, ok, false, err
}

func PickTask(ctx context.Context, runner *actions_model.ActionRunner) (*runnerv1.Task, bool, error) {
	var (
		task       *runnerv1.Task
		// The job was already claimed but assembling its payload failed; release the
		// claim so the job returns to the waiting queue instead of being stranded in
		// running state with no runner ever executing it.
		if relErr := actions_model.ReleaseTaskForRunner(ctx, t); relErr != nil {
			log.Error("ReleaseTaskForRunner [task_id: %d]: %v", t.ID, relErr)
		}
		return nil, false, err
	}
	actionTask = t
	// The job is claimed and its payload assembled, but if the request context was cancelled meanwhile, response can no longer reach the runner.
	// Release the claim so another runner can pick the job up.
	if err := ctx.Err(); err != nil {
		if relErr := actions_model.ReleaseTaskForRunner(ctx, t); relErr != nil {
			log.Error("ReleaseTaskForRunner [task_id: %d]: %v", t.ID, relErr)
		}
		return nil, false, err
	}

	"testing"

	actions_model "gitea.dev/models/actions"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	assert.False(t, ok)
	assert.True(t, throttled)
}
	if subtle.ConstantTimeCompare([]byte(t.TokenHash), []byte(hex.EncodeToString(hashedToken[:]))) == 0 {
		// If an attacker steals a token and uses the token to create a new session the hash gets updated.
		// When the victim uses the old token the hashes don't match anymore and the victim should be notified about the compromised token.
		return nil, ErrAuthTokenInvalidHash
	}

	auth_model "gitea.dev/models/auth"
	"gitea.dev/models/unittest"
	"gitea.dev/modules/timeutil"

	"github.com/stretchr/testify/assert"
)
		assert.ErrorIs(t, err, ErrAuthTokenInvalidHash)
		assert.Nil(t, at2)

		assert.NoError(t, auth_model.DeleteAuthTokenByID(t.Context(), at.ID))
	})

	t.Run("Valid", func(t *testing.T) {
}

// TokenCanAccessRepo reports whether the current API token is allowed to access the repository.
// A public-only token cannot reach a private repo; any other token is unrestricted by this check.
func (ctx *APIContext) TokenCanAccessRepo(repo *repo_model.Repository) bool {
	return repo == nil || !ctx.PublicOnly || !repo.IsPrivate
}

func init() {
package context

import (
	"net/http"
	"slices"

	"gitea.dev/models/unit"
)

// CheckTokenScopes checks whether the authenticated API token contains any of the given scopes.
func CheckTokenScopes(ctx *Context, repo *repo_model.Repository, scopes ...auth_model.AccessTokenScope) {
	if ctx.Data["IsApiToken"] != true {
		return
	}

	if publicOnly && repo != nil && repo.IsPrivate {
		ctx.HTTPError(http.StatusForbidden)
		return
	}
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
						}
						defer rc.Close()
					} else {
						resp, err := http.Get(*asset.DownloadURL)
						if err != nil {
							return err
						}
		}

		// SECURITY: We will assume that the pr.PatchURL has been checked
		// pr.PatchURL maybe a local file - but note EnsureSafe should be asserting that this safe
		resp, err := http.Get(u) // TODO: This probably needs to use the downloader as there may be rate limiting issues here
		if err != nil {
			return err
		}
		baseURL,
		gitea_sdk.SetToken(token),
		gitea_sdk.SetBasicAuth(username, password),
		gitea_sdk.SetHTTPClient(NewMigrationHTTPClient()),
	)
	if err != nil {
		log.Error(fmt.Sprintf("Failed to create NewGiteaDownloader for: %s. Error: %v", baseURL, err))
		Created:         rel.CreatedAt,
	}

	httpClient := NewMigrationHTTPClient()

	for _, asset := range rel.Attachments {
		assetID := asset.ID // Don't optimize this, for closure we need a local variable
						return err
					}
				} else if asset.DownloadURL != nil {
					rc, err = uri.Open(*asset.DownloadURL)
					if err != nil {
						return err
					}
		}

		// SECURITY: We will assume that the pr.PatchURL has been checked
		// pr.PatchURL maybe a local file - but note EnsureSafe should be asserting that this safe
		ret, err := uri.Open(pr.PatchURL) // TODO: This probably needs to use the downloader as there may be rate limiting issues here
		if err != nil {
			return err
		}
		r.Published = rel.PublishedAt.Time
	}

	httpClient := NewMigrationHTTPClient()

	for _, asset := range rel.Assets {
		assetID := asset.GetID() // Don't optimize this, for closure we need a local variable TODO: no need to do so in new Golang
//	Use either a username/password, personal token entered into the username field, or anonymous/public access
//	Note: Public access only allows very basic access
func NewGitlabDownloader(ctx context.Context, baseURL, repoPath, token string) (*GitlabDownloader, error) {
	gitlabClient, err := gitlab.NewClient(token, gitlab.WithBaseURL(baseURL), gitlab.WithHTTPClient(NewMigrationHTTPClient()))
	if err != nil {
		log.Trace("Error logging into gitlab: %v", err)
		return nil, err
		PublisherName:   rel.Author.Username,
	}

	httpClient := NewMigrationHTTPClient()

	for _, asset := range rel.Assets.Links {
		assetID := asset.ID // Don't optimize this, for closure we need a local variable
	"gitea.dev/modules/hostmatcher"
	"gitea.dev/modules/proxy"
	"gitea.dev/modules/setting"
)

// NewMigrationHTTPClient returns a HTTP client for migration
func NewMigrationHTTPClient() *http.Client {
	return &http.Client{
		Transport: NewMigrationHTTPTransport(),
	}
}

// NewMigrationHTTPTransport returns a HTTP transport for migration
func NewMigrationHTTPTransport() *http.Transport {
	return &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: setting.Migrations.SkipTLSVerify},
		Proxy:           proxy.Proxy(),
		DialContext:     hostmatcher.NewDialContext("migration", allowList, blockList, setting.Proxy.ProxyURLFixed),
	}
}
		blockList.AppendBuiltin(hostmatcher.MatchBuiltinLoopback)
	}

	return nil
}
		log.Error("SyncMirrors [repo: %-v]: GetRemoteURL Error %v", m.Repo, remoteErr)
		return nil, false
	}
	envs := proxy.EnvWithProxy(remoteURL.URL)
	timeout := time.Duration(setting.Git.Timeout.Mirror) * time.Second

			return err
		}

		if err := repo_model.ClearRepoStars(ctx, repo.ID); err != nil {
			return err
		}
	}

	// Create/Remove git-daemon-export-ok for git-daemon...
		require.NoError(t, ChangeOrganizationVisibility(t.Context(), org, structs.VisibleTypePrivate))
		unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: org.ID, Visibility: structs.VisibleTypePrivate})
	})
}
	user2 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 2})
	org3 := unittest.AssertExistsAndLoadBean(t, &org_model.Organization{ID: 3})
	user4 := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: 4})

	t.Run("User projects", func(t *testing.T) {
		pi1 := project_model.ProjectIssue{
		})

		t.Run("Authenticated user with no permission to the private repo", func(t *testing.T) {
			columnIssues, err := LoadIssuesFromProject(t.Context(), projects[0], &issues_model.IssuesOptions{
				Owner: org3.AsUser(),
				Doer:  user2,
			})
			assert.NoError(t, err)
			assert.Len(t, columnIssues, 1)
			assert.Len(t, columnIssues[defaultColumn.ID], 1) // user2 can only visit public repo issues
		})
	})

			if err = repo_model.ClearRepoStars(ctx, repo.ID); err != nil {
				return err
			}
		}

		// Create/Remove git-daemon-export-ok for git-daemon...
	assert.True(t, updatedRepo.IsPrivate)
	assert.Zero(t, updatedRepo.NumWatches)
}
// Init starts the hooks delivery thread
func Init() error {
	timeout := time.Duration(setting.Webhook.DeliverTimeout) * time.Second

	allowedHostListValue := setting.Webhook.AllowedHostList
	if allowedHostListValue == "" {
		allowedHostListValue = hostmatcher.MatchBuiltinExternal
	}
	allowedHostMatcher := hostmatcher.ParseHostMatchList("webhook.ALLOWED_HOST_LIST", allowedHostListValue)

	webhookHTTPClient = &http.Client{
		Timeout: timeout,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: setting.Webhook.SkipTLSVerify},
			Proxy:           webhookProxy(allowedHostMatcher),
			DialContext:     hostmatcher.NewDialContext("webhook", allowedHostMatcher, nil, setting.Webhook.ProxyURLFixed),
		},
	}

	hookQueue = queue.CreateUniqueQueue(graceful.GetManager().ShutdownContext(), "webhook_sender", handler)
		config["color"] = s.Color
	}

	authorizationHeader, err := w.HeaderAuthorization()
	if err != nil {
		return nil, err
	}

	return &api.Hook{
		ID:                  w.ID,
		Name:                w.Name,
		Type:                w.Type,
		URL:                 fmt.Sprintf("%s/settings/hooks/%d", repoLink, w.ID),
		Active:              w.IsActive,
		Config:              config,
		Events:              w.EventsArray(),
		AuthorizationHeader: authorizationHeader,
		Updated:             w.UpdatedUnix.AsTime(),
		Created:             w.CreatedUnix.AsTime(),
		BranchFilter:        w.BranchFilter,
)

func TestMain(m *testing.M) {
	// for tests, allow only loopback IPs
	setting.Webhook.AllowedHostList = hostmatcher.MatchBuiltinLoopback
	unittest.MainTest(m, &unittest.TestOptions{
		SetUp: func() error {
			setting.LoadQueueSettings()
			return Init()
		},
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/migration"
	api "gitea.dev/modules/structs"
	mirror_service "gitea.dev/services/mirror"
	repo_service "gitea.dev/services/repository"
	files_service "gitea.dev/services/repository/files"

	"github.com/stretchr/testify/assert"
)

func TestScheduleUpdate(t *testing.T) {

func testScheduleUpdateMirrorSync(t *testing.T) {
	doTestScheduleUpdate(t, func(t *testing.T, u *url.URL, testContext APITestContext, user *user_model.User, repo *repo_model.Repository) (commitID, expectedSpec string) {
		// create mirror repo
		opts := migration.MigrateOptions{
			RepoName:    "actions-schedule-mirror",
		assert.NotNil(t, run.TriggerActor, "trigger_actor should be populated")
	}
}
	unittest.AssertExistsAndLoadBean(t, &issues_model.IssueLabel{IssueID: issue.ID, LabelID: 2})
}

func TestAPIAddIssueLabelsWithLabelNames(t *testing.T) {
	assert.NoError(t, unittest.LoadFixtures())

	assert.EqualValues(t, 33, apiNewTime.Time)
	assert.Equal(t, user2.ID, apiNewTime.UserID)
	assert.EqualValues(t, 947688818, apiNewTime.Created.Unix())
}
	neturl "net/url"
	"testing"

	"gitea.dev/models/packages"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
		})
	})
}
		test(t, http.StatusNotFound, packageTag2, "1.2")
		test(t, http.StatusOK, packageTag, packageVersion)
		test(t, http.StatusOK, packageTag2, packageVersion)
	})

	t.Run("ListTags", func(t *testing.T) {
		}
	}
}

	apiHook := DecodeJSON(t, resp, &api.Hook{})
	assert.Equal(t, "http://example.com/", apiHook.Config["url"])
	assert.Equal(t, "Bearer s3cr3t", apiHook.AuthorizationHeader)
	assert.Equal(t, "CI notifications", apiHook.Name)

	newName := "Deploy hook"
	patchReq := NewRequestWithJSON(t, "PATCH", fmt.Sprintf("/api/v1/repos/%s/%s/hooks/%d", owner.Name, repo.Name, apiHook.ID), api.EditHookOption{
		Name: &newName,
	deleteAPIAccessToken(t, newAccessToken, user)
}

// TestAPIDeleteMissingToken ensures that error is thrown when token not found
func TestAPIDeleteMissingToken(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	"testing"

	auth_model "gitea.dev/models/auth"
	api "gitea.dev/modules/structs"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

func TestAPIListEmails(t *testing.T) {
	defer tests.PrepareTestEnv(t)()

	}
	require.NoError(t, user_model.LinkExternalToUser(t.Context(), u, extLink))

	prepareUserExternalLink := func(t *testing.T, refreshToken string) {
		err := user_model.UpdateUserCols(t.Context(), &user_model.User{ID: u.ID, IsActive: false}, "is_active")
		require.NoError(t, err)
		_, err = db.GetEngine(t.Context()).Where(builder.Eq{"user_id": u.ID}).Cols("refresh_token").
			Update(&user_model.ExternalLoginUser{RefreshToken: refreshToken})
		require.NoError(t, err)
		require.False(t, unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID}).IsActive)
	}

	t.Run("admin-disabled user is not reactivated", func(t *testing.T) {
		prepareUserExternalLink(t, "non-empty-refresh-token")
		doOIDCSignIn(t, authSource.Name)
		after := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID})
		assert.False(t, after.IsActive, "OAuth callback must not re-enable an administrator-disabled account")
	})

	t.Run("auto-sync-disabled user is reactivated", func(t *testing.T) {
		prepareUserExternalLink(t, "" /* empty refresh token */)
		doOIDCSignIn(t, authSource.Name)
		after := unittest.AssertExistsAndLoadBean(t, &user_model.User{ID: u.ID})
		assert.True(t, after.IsActive, "OAuth callback must reactivate a sync-disabled account on successful login")
	"net/http"
	"testing"

	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
		})
	})
}
	"net/http"
	"testing"

	"gitea.dev/modules/setting"
	"gitea.dev/modules/test"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

func TestGoGet(t *testing.T) {

	assert.Equal(t, expected, resp.Body.String())
}
package integration

import (
	"slices"
	"testing"

	"gitea.dev/models/db"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/gitrepo"
	"gitea.dev/modules/migration"
	mirror_service "gitea.dev/services/mirror"
	release_service "gitea.dev/services/release"
	repo_service "gitea.dev/services/repository"
	mirror = unittest.AssertExistsAndLoadBean(t, &repo_model.Mirror{RepoID: mirrorRepo.ID})
	assert.Equal(t, lastMirrorSync, mirror.LastSyncUnix)
}
func TestOAuth2AvatarFromPicture(t *testing.T) {
	defer tests.PrepareTestEnv(t)()
	defer test.MockVariableValue(&setting.OAuth2Client.UpdateAvatar, true)()

	mockServer := createOAuth2MockProvider()
	defer mockServer.Close()
	testOAuth2(t, "/user/oauth2/test%2Bplus", http.StatusTemporaryRedirect)
	testOAuth2(t, "/user/oauth2/test%20plus", http.StatusNotFound)
}
	"path"
	"strings"
	"testing"

	"gitea.dev/models/db"
	issues_model "gitea.dev/models/issues"
	"gitea.dev/models/unittest"
	user_model "gitea.dev/models/user"
	"gitea.dev/modules/git"
	"gitea.dev/modules/test"
	issue_service "gitea.dev/services/issue"
	repo_service "gitea.dev/services/repository"
	"gitea.dev/tests"

	"github.com/stretchr/testify/assert"
)

func TestPullView_ReviewerMissed(t *testing.T) {
			unittest.AssertExistsAndLoadBean(t, &issues_model.Review{IssueID: pr.IssueID, Type: issues_model.ReviewTypeRequest, ReviewerID: 5})
			assert.NoError(t, pr.LoadIssue(t.Context()))

			// update the file on the pr branch
			_, err = files_service.ChangeRepoFiles(t.Context(), repo, user2, &files_service.ChangeRepoFilesOptions{
				OldBranch: "codeowner-basebranch",
			})
			assert.NoError(t, err)

			reviewNotifiers, err := issue_service.PullRequestCodeOwnersReview(t.Context(), pr)
			assert.NoError(t, err)
			assert.Len(t, reviewNotifiers, 1)
			assert.EqualValues(t, 8, reviewNotifiers[0].Reviewer.ID)

			err = issue_service.ChangeTitle(t.Context(), pr.Issue, user2, "[WIP] Test Pull Request")
	})
}

func updateRepoPullRequestConfig(t *testing.T, repoID int64, update func(*repo_model.PullRequestsConfig)) {
	t.Helper()

