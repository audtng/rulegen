package main

type externalEntityProviderRBAC struct {
	thirdPartyIntegration    shared.IntegrationAggregate
	externalEntityProviderID string
	adminToken               *string
	ctx                      shared.Context

	rootAccessControl shared.AccessControl

var _ shared.AccessControl = (*externalEntityProviderRBAC)(nil)

func NewExternalEntityProviderRBAC(ctx shared.Context, rootAccessControl shared.AccessControl, thirdPartyIntegration shared.IntegrationAggregate, externalEntityProviderID string, adminToken *string) *externalEntityProviderRBAC {
	return &externalEntityProviderRBAC{
		thirdPartyIntegration:    thirdPartyIntegration,
		externalEntityProviderID: externalEntityProviderID,
		adminToken:               adminToken,
		ctx:                      ctx,
		rootAccessControl:        rootAccessControl,
	}
}

func (e *externalEntityProviderRBAC) HasAccess(ctx context.Context, userID string) (bool, error) {
	if e.adminToken != nil && userID == *e.adminToken {
		return true, nil
	}
	return e.thirdPartyIntegration.HasAccessToExternalEntityProvider(e.ctx, e.externalEntityProviderID)
}

}

func (e *externalEntityProviderRBAC) IsAllowed(ctx context.Context, userID string, object shared.Object, action shared.Action) (bool, error) {
	if e.adminToken != nil && userID == *e.adminToken {
		if action == shared.ActionRead {
			return true, nil
		}
		return false, nil
	}

	// ALLOW ORG read access for all users - this is pretty much the same as HasAccess.
	if object == shared.ObjectOrganization && action == shared.ActionRead {
		return true, nil
	if project.ExternalEntityProviderID == nil || project.ExternalEntityID == nil {
		return false, nil
	}
	if e.adminToken != nil && user == *e.adminToken && action == shared.ActionRead {
		return true, nil
	}
	return e.rootAccessControl.IsAllowedInProject(ctx, project, user, object, action)
}

	if asset.ExternalEntityProviderID == nil || asset.ExternalEntityID == nil {
		return false, nil
	}
	if e.adminToken != nil && user == *e.adminToken && action == shared.ActionRead {
		return true, nil
	}
	return e.rootAccessControl.IsAllowedInAsset(ctx, asset, user, object, action)
}


	"github.com/l3montree-dev/devguard/mocks"
	"github.com/l3montree-dev/devguard/shared"
	"github.com/l3montree-dev/devguard/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIsAllowed(t *testing.T) {

	type testCase struct {
		name           string
		userID         string
		object         shared.Object
		action         shared.Action
		adminToken     *string
		mockResult     bool
		mockErr        error
		expectedResult bool
			userID:         "admin-token",
			object:         shared.ObjectProject,
			action:         shared.ActionRead,
			adminToken:     utils.Ptr("admin-token"),
			expectedResult: true,
		},
		{
			userID:         "user1",
			object:         shared.ObjectOrganization,
			action:         shared.ActionRead,
			adminToken:     utils.Ptr("admin-token"),
			expectedResult: true,
		},

		{
			name:       "error from rootAccessControl",
			userID:     "user5",
			object:     shared.ObjectProject,
			action:     shared.ActionRead,
			adminToken: utils.Ptr("admin-token"),
			mockErr:    errors.New("some error"),
			expectErr:  true,
		},
		{
			name:           "admin token can not create",
			userID:         "admin-token",
			object:         shared.ObjectProject,
			action:         shared.ActionCreate,
			adminToken:     utils.Ptr("admin-token"),
			expectedResult: false,
		},
		{
			userID:         "admin-token",
			object:         shared.ObjectProject,
			action:         shared.ActionDelete,
			adminToken:     utils.Ptr("admin-token"),
			expectedResult: false,
		},
	}
				rootAccessControl,
				thirdpartyIntegrationMock,
				"external-entity-provider-id",
				tc.adminToken,
			)

			result, err := rbac.IsAllowed(context.Background(), tc.userID, tc.object, tc.action)
			rootAccessControl,
			thirdpartyIntegrationMock,
			"external-entity-provider-id",
			utils.Ptr("admin-token"),
		)

		hasAccess, err := rbac.HasAccess(context.Background(), "admin-token")
			rootAccessControl,
			thirdpartyIntegrationMock,
			"external-entity-provider-id",
			nil,
		)

		hasAccess, err := rbac.HasAccess(context.Background(), "user1")
		assert.NoError(t, err)
		assert.True(t, hasAccess)
	})

	t.Run("if no admin token is provided, the third party integration should be called (false)", func(t *testing.T) {
		ctx := mocks.NewContext(t)
		rootAccessControl := mocks.NewAccessControl(t)
		thirdpartyIntegrationMock := mocks.NewIntegrationAggregate(t)
		thirdpartyIntegrationMock.On("HasAccessToExternalEntityProvider", ctx, "external-entity-provider-id").Return(false, nil)

		rbac := NewExternalEntityProviderRBAC(
			ctx,
			rootAccessControl,
			thirdpartyIntegrationMock,
			"external-entity-provider-id",
			nil,
		)

		hasAccess, err := rbac.HasAccess(context.Background(), "user1")
		assert.NoError(t, err)
		assert.False(t, hasAccess)
	})
}

	DevGuardBotUserID          int    // the user id of the devguard bot user, used to create issues
	DevGuardBotUserAccessToken string // the access token of the devguard bot user, used to create issues
	AdminToken                 *string
}

func (c *GitlabOauth2Config) GetProviderID() string {
	appID              string
	appSecret          string
	scopes             string
	botUserID          int     // the user id of the devguard bot user, used to create issues
	botUserAccessToken string  // the access token of the devguard bot user, used to create issues
	adminToken         *string // the admin token for the gitlab instance, used to create issues
}

type gitlabOauth2Client struct {
				conf.botUserID = intValue
			case "botuseraccesstoken":
				conf.botUserAccessToken = value
			case "admintoken":
				if value == "" {
					conf.adminToken = nil
				} else {
					conf.adminToken = &value
				}
			}

			urls[name] = conf
		if conf.botUserAccessToken == "" {
			slog.Warn(fmt.Sprintf("GITLAB_%s_BOTUSERACCESSTOKEN is not set", strings.ToUpper(name)))
		}
		if conf.adminToken == nil {
			slog.Warn(fmt.Sprintf("GITLAB_%s_ADMINTOKEN is not set", strings.ToUpper(name)))
		} else {
			slog.Info(fmt.Sprintf("GITLAB_%s_ADMINTOKEN is set", strings.ToUpper(name)))
		}
	}

	return urls
}

func NewGitLabOauth2Config(id, gitlabBaseURL, gitlabOauth2ClientID, gitlabOauth2ClientSecret, gitlabOauth2Scopes string, botUserID int, botUserAccessToken string, adminToken *string, gitlabOauth2TokenRepository shared.GitLabOauth2TokenRepository) *GitlabOauth2Config {

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		ProviderID:                 id,
		DevGuardBotUserID:          botUserID,
		DevGuardBotUserAccessToken: botUserAccessToken,
		AdminToken:                 adminToken,
		Oauth2Conf: &oauth2.Config{
			ClientID:     gitlabOauth2ClientID,
			ClientSecret: gitlabOauth2ClientSecret,
	envs := parseGitlabEnvs()
	gitlabIntegrations := make(map[string]*GitlabOauth2Config)
	for id, env := range envs {
		gitlabIntegration := NewGitLabOauth2Config(id, env.baseURL, env.appID, env.appSecret, env.scopes, env.botUserID, env.botUserAccessToken, env.adminToken, gitlabOauth2TokenRepository)
		gitlabIntegrations[id] = gitlabIntegration
		slog.Info("gitlab oauth2 integration created", "id", id, "baseURL", env.baseURL, "appID", env.appID)
	}
			domainRBAC := rbacProvider.GetDomainRBAC(org.ID.String())
			if org.IsExternalEntity() {
				// check if there is an admin token defined
				conf, ok := oauth2Config[*org.ExternalEntityProviderID]
				if !ok {
					slog.Error("no oauth2 config found for external entity provider", "provider", *org.ExternalEntityProviderID)
					return ctx.JSON(500, map[string]string{"error": "no oauth2 config found for external entity provider"})
				}

				domainRBAC = accesscontrol.NewExternalEntityProviderRBAC(ctx, rbacProvider.GetDomainRBAC(org.ID.String()), shared.GetThirdPartyIntegration(ctx), *org.ExternalEntityProviderID, conf.AdminToken)
			}

			// check if the user is allowed to access the organization
			var scopes string
			var err error

			adminTokenHeader := ctx.Request().Header.Get("X-Admin-Token")

			if oryKratosSessionCookie != nil {
				userID, err = cookieAuth(ctx.Request().Context(), oryAPIClient, oryKratosSessionCookie.String())
				if err != nil {
				scopesArray := strings.Fields(scopes)
				ctx.Set("session", accesscontrol.NewSession(userID, scopesArray))
				return next(ctx)
			} else if adminTokenHeader != "" {
				slog.Warn("admin token header is set, using it to create session")
				ctx.Set("session", accesscontrol.NewSession(adminTokenHeader, []string{}))
				return next(ctx)
			} else {
				userID, scopes, err = verifier.VerifyRequestSignature(ctx.Request().Context(), ctx.Request())
				if err != nil {
		_ = handler(c)
		assert.True(t, called)
	})

	t.Run("should set the session using admin token header", func(t *testing.T) {
		e := echo.New()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("X-Admin-Token", "admin_token_value")
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		mw := SessionMiddleware(nil, nil)

		var called bool
		handler := mw(func(ctx echo.Context) error {
			called = true
			sess := shared.GetSession(ctx)

			assert.Equal(t, "admin_token_value", sess.GetUserID())
			assert.ElementsMatch(t, []string{}, sess.GetScopes())
			return nil
		})

		_ = handler(c)
		assert.True(t, called)
	})
}
