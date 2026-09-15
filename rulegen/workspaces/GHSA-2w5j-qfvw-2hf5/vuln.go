package main

		traceSpan.EndWithError(err)
	}

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnProjectID.identifier():  projectID,
	return app, err
}

func (q *Queries) AppByID(ctx context.Context, appID string) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{
		AppColumnID.identifier():         appID,
		AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-immt9", "Errors.Query.SQLStatement")
	return app, err
}

func (q *Queries) AppBySAMLEntityID(ctx context.Context, entityID string) (app *App, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	eq := sq.Eq{
		AppSAMLConfigColumnEntityID.identifier(): entityID,
		AppColumnInstanceID.identifier():         authz.GetInstance(ctx).InstanceID(),
	}
	query, args, err := stmt.Where(eq).ToSql()
	if err != nil {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := prepareAppQuery(ctx, q.client)
	eq := sq.Eq{AppColumnInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()}
	query, args, err := stmt.Where(sq.And{
		eq,
		sq.Or{
	return NewTextQuery(AppColumnProjectID, id, TextEquals)
}

func prepareAppQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Row) (*App, error)) {
	return sq.Select(
			AppColumnID.identifier(),
			AppColumnName.identifier(),
			AppColumnProjectID.identifier(),
			AppColumnCreationDate.identifier(),
			AppColumnChangeDate.identifier(),
			AppColumnResourceOwner.identifier(),
			AppColumnState.identifier(),
			AppColumnSequence.identifier(),

			AppAPIConfigColumnAppID.identifier(),
			AppAPIConfigColumnClientID.identifier(),
			AppAPIConfigColumnAuthMethod.identifier(),

			AppOIDCConfigColumnAppID.identifier(),
			AppOIDCConfigColumnVersion.identifier(),
			AppOIDCConfigColumnClientID.identifier(),
			AppOIDCConfigColumnRedirectUris.identifier(),
			AppOIDCConfigColumnResponseTypes.identifier(),
			AppOIDCConfigColumnGrantTypes.identifier(),
			AppOIDCConfigColumnApplicationType.identifier(),
			AppOIDCConfigColumnAuthMethodType.identifier(),
			AppOIDCConfigColumnPostLogoutRedirectUris.identifier(),
			AppOIDCConfigColumnDevMode.identifier(),
			AppOIDCConfigColumnAccessTokenType.identifier(),
			AppOIDCConfigColumnAccessTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenRoleAssertion.identifier(),
			AppOIDCConfigColumnIDTokenUserinfoAssertion.identifier(),
			AppOIDCConfigColumnClockSkew.identifier(),
			AppOIDCConfigColumnAdditionalOrigins.identifier(),
			AppOIDCConfigColumnSkipNativeAppSuccessPage.identifier(),

			AppSAMLConfigColumnAppID.identifier(),
			AppSAMLConfigColumnEntityID.identifier(),
			AppSAMLConfigColumnMetadata.identifier(),
			AppSAMLConfigColumnMetadataURL.identifier(),
		).From(appsTable.identifier()).
			LeftJoin(join(AppAPIConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppOIDCConfigColumnAppID, AppColumnID)).
			LeftJoin(join(AppSAMLConfigColumnAppID, AppColumnID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*App, error) {
			app := new(App)

			var (
				apiConfig  = sqlAPIConfig{}
				oidcConfig = sqlOIDCConfig{}
				samlConfig = sqlSAMLConfig{}
			)

			err := row.Scan(
				&app.ID,
				&app.Name,
				&app.ProjectID,
				&app.CreationDate,
				&app.ChangeDate,
				&app.ResourceOwner,
				&app.State,
				&app.Sequence,

				&apiConfig.appID,
				&apiConfig.clientID,
				&apiConfig.authMethod,

				&oidcConfig.appID,
				&oidcConfig.version,
				&oidcConfig.clientID,
				&oidcConfig.redirectUris,
				&oidcConfig.responseTypes,
				&oidcConfig.grantTypes,
				&oidcConfig.applicationType,
				&oidcConfig.authMethodType,
				&oidcConfig.postLogoutRedirectUris,
				&oidcConfig.devMode,
				&oidcConfig.accessTokenType,
				&oidcConfig.accessTokenRoleAssertion,
				&oidcConfig.iDTokenRoleAssertion,
				&oidcConfig.iDTokenUserinfoAssertion,
				&oidcConfig.clockSkew,
				&oidcConfig.additionalOrigins,
				&oidcConfig.skipNativeAppSuccessPage,

				&samlConfig.appID,
				&samlConfig.entityID,
				&samlConfig.metadata,
				&samlConfig.metadataURL,
			)

			if err != nil {
				if errors.Is(err, sql.ErrNoRows) {
					return nil, zerrors.ThrowNotFound(err, "QUERY-pCP8P", "Errors.App.NotExisting")
				}
				return nil, zerrors.ThrowInternal(err, "QUERY-4SJlx", "Errors.Internal")
			}

			apiConfig.set(app)
			oidcConfig.set(app)
			samlConfig.set(app)

			return app, nil
		}
}

func prepareOIDCAppQuery() (sq.SelectBuilder, func(*sql.Row) (*App, error)) {
			AppSAMLConfigColumnMetadataURL.identifier(),
		).From(appsTable.identifier()).
			Join(join(AppSAMLConfigColumnAppID, AppColumnID)).
			PlaceholderFormat(sq.Dollar), func(row *sql.Row) (*App, error) {

			app := new(App)
	"github.com/zitadel/zitadel/pkg/grpc/authn"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

func (i *Instance) CreateOIDCClient(ctx context.Context, redirectURI, logoutRedirectURI, projectID string, appType app.OIDCAppType, authMethod app.OIDCAuthMethodType, devMode bool, grantTypes ...app.OIDCGrantType) (*management.AddOIDCAppResponse, error) {
	return client, err
}

func (i *Instance) CreateOIDCImplicitFlowClient(ctx context.Context, redirectURI string) (*management.AddOIDCAppResponse, error) {
	project, err := i.Client.Mgmt.AddProject(ctx, &management.AddProjectRequest{
		Name: fmt.Sprintf("project-%d", time.Now().UnixNano()),
	return machine, name, secret.GetClientId(), secret.GetClientSecret(), nil
}

func (i *Instance) CreateOIDCJWTProfileClient(ctx context.Context) (machine *management.AddMachineUserResponse, name string, keyData []byte, err error) {
	name = gofakeit.Username()
	machine, err = i.Client.Mgmt.AddMachineUser(ctx, &management.AddMachineUserRequest{
		err = oidcError(err)
		span.EndWithError(err)
	}()
	client, err := o.query.GetOIDCClientByID(ctx, id, false)
	if err != nil {
		return nil, err
	}
	if client.State != domain.AppStateActive {
		return nil, zerrors.ThrowPreconditionFailed(nil, "OIDC-sdaGg", "client is not active")
	}
	return ClientFromBusiness(client, o.defaultLoginURL, o.defaultLoginURLV2), nil
}

	if err != nil {
		return err
	}
	var allRoles bool
	roles := make([]string, 0)
	for _, scope := range scopes {
	if projectID != "" {
		roleAudience = append(roleAudience, projectID)
	}
	queries := make([]query.SearchQuery, 0, 2)
	projectQuery, err := query.NewUserGrantProjectIDsSearchQuery(roleAudience)
	if err != nil {
		return nil, nil, err
	}
	queries = append(queries, projectQuery)
	userIDQuery, err := query.NewUserGrantUserIDSearchQuery(userID)
	if err != nil {
		return nil, nil, err
	}
	queries = append(queries, userIDQuery)
	grants, err := o.query.UserGrants(ctx, &query.UserGrantsQueries{
		Queries: queries,
	}, true)
	if err != nil {
		return nil, nil, err
	if err != nil {
		return nil, err
	}
	client, err := s.query.GetOIDCClientByID(ctx, clientID, assertion)
	if zerrors.IsNotFound(err) {
		return nil, oidc.ErrInvalidClient().WithParent(err).WithReturnParentToClient(authz.GetFeatures(ctx).DebugOIDCParentError).WithDescription("client not found")
	}
	if err != nil {
		return nil, err // defaults to server error
	}
	if client.State != domain.AppStateActive {
		return nil, oidc.ErrInvalidClient().WithDescription("client is not active")
	}
	if client.Settings == nil {
		client.Settings = &query.OIDCSettings{
			AccessTokenLifetime: s.defaultAccessTokenLifetime,
