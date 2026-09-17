package main


	DevGuardBotUserID          int    // the user id of the devguard bot user, used to create issues
	DevGuardBotUserAccessToken string // the access token of the devguard bot user, used to create issues
}

func (c *GitlabOauth2Config) GetProviderID() string {
	appID              string
	appSecret          string
	scopes             string
	botUserID          int    // the user id of the devguard bot user, used to create issues
	botUserAccessToken string // the access token of the devguard bot user, used to create issues
}

type gitlabOauth2Client struct {
				conf.botUserID = intValue
			case "botuseraccesstoken":
				conf.botUserAccessToken = value
			}

			urls[name] = conf
		if conf.botUserAccessToken == "" {
			slog.Warn(fmt.Sprintf("GITLAB_%s_BOTUSERACCESSTOKEN is not set", strings.ToUpper(name)))
		}
	}

	return urls
}

func NewGitLabOauth2Config(id, gitlabBaseURL, gitlabOauth2ClientID, gitlabOauth2ClientSecret, gitlabOauth2Scopes string, botUserID int, botUserAccessToken string, gitlabOauth2TokenRepository shared.GitLabOauth2TokenRepository) *GitlabOauth2Config {

	frontendURL := os.Getenv("FRONTEND_URL")
	if frontendURL == "" {
		ProviderID:                 id,
		DevGuardBotUserID:          botUserID,
		DevGuardBotUserAccessToken: botUserAccessToken,
		Oauth2Conf: &oauth2.Config{
			ClientID:     gitlabOauth2ClientID,
			ClientSecret: gitlabOauth2ClientSecret,
	envs := parseGitlabEnvs()
	gitlabIntegrations := make(map[string]*GitlabOauth2Config)
	for id, env := range envs {
		gitlabIntegration := NewGitLabOauth2Config(id, env.baseURL, env.appID, env.appSecret, env.scopes, env.botUserID, env.botUserAccessToken, gitlabOauth2TokenRepository)
		gitlabIntegrations[id] = gitlabIntegration
		slog.Info("gitlab oauth2 integration created", "id", id, "baseURL", env.baseURL, "appID", env.appID)
	}
type externalEntityProviderRBAC struct {
	thirdPartyIntegration    shared.IntegrationAggregate
	externalEntityProviderID string
	ctx                      shared.Context

	rootAccessControl shared.AccessControl

var _ shared.AccessControl = (*externalEntityProviderRBAC)(nil)

func NewExternalEntityProviderRBAC(ctx shared.Context, rootAccessControl shared.AccessControl, thirdPartyIntegration shared.IntegrationAggregate, externalEntityProviderID string) *externalEntityProviderRBAC {
	return &externalEntityProviderRBAC{
		thirdPartyIntegration:    thirdPartyIntegration,
		externalEntityProviderID: externalEntityProviderID,
		ctx:                      ctx,
		rootAccessControl:        rootAccessControl,
	}
}

func (e *externalEntityProviderRBAC) HasAccess(ctx context.Context, userID string) (bool, error) {
	return e.thirdPartyIntegration.HasAccessToExternalEntityProvider(e.ctx, e.externalEntityProviderID)
}

}

func (e *externalEntityProviderRBAC) IsAllowed(ctx context.Context, userID string, object shared.Object, action shared.Action) (bool, error) {
	// ALLOW ORG read access for all users - this is pretty much the same as HasAccess.
	if object == shared.ObjectOrganization && action == shared.ActionRead {
		return true, nil
	if project.ExternalEntityProviderID == nil || project.ExternalEntityID == nil {
		return false, nil
	}
	return e.rootAccessControl.IsAllowedInProject(ctx, project, user, object, action)
}

	if asset.ExternalEntityProviderID == nil || asset.ExternalEntityID == nil {
		return false, nil
	}
	return e.rootAccessControl.IsAllowedInAsset(ctx, asset, user, object, action)
}

			domainRBAC := rbacProvider.GetDomainRBAC(org.ID.String())
			if org.IsExternalEntity() {
				// check if there is an admin token defined
				domainRBAC = accesscontrol.NewExternalEntityProviderRBAC(ctx, rbacProvider.GetDomainRBAC(org.ID.String()), shared.GetThirdPartyIntegration(ctx), *org.ExternalEntityProviderID)
			}

			// check if the user is allowed to access the organization
