package main


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
