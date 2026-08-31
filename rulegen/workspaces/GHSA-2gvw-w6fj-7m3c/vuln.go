package main

	"k8s.io/apimachinery/pkg/runtime/schema"
)

// AppProjectList is list of AppProject resources
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AppProjectList struct {
//
// If the user does provide a "project," we can respond more specifically. If the user does not have access to the given
// app name in the given project, we return "permission denied." If the app exists, but the project is different from
func (s *Server) getAppEnforceRBAC(ctx context.Context, action, project, namespace, name string, getApp func() (*appv1.Application, error)) (*appv1.Application, error) {
	user := session.Username(ctx)
	if user == "" {
		user = "Unknown user"
			// but the app is in a different project" response. We don't want the user inferring the existence of the
			// app from response time.
			_, _ = getApp()
			return nil, permissionDeniedErr
		}
	}
	a, err := getApp()
		if apierr.IsNotFound(err) {
			if project != "" {
				// We know that the user was allowed to get the Application, but the Application does not exist. Return 404.
				return nil, status.Errorf(codes.NotFound, apierr.NewNotFound(schema.GroupResource{Group: "argoproj.io", Resource: "applications"}, name).Error())
			}
			// We don't know if the user was allowed to get the Application, and we don't want to leak information about
			// the Application's existence. Return 403.
			logCtx.Warn("application does not exist")
			return nil, permissionDeniedErr
		}
		logCtx.Errorf("failed to get application: %s", err)
		return nil, permissionDeniedErr
	}
	// Even if we performed an initial RBAC check (because the request was fully parameterized), we still need to
	// perform a second RBAC check to ensure that the user has access to the actual Application's project (not just the
			// The user specified a project. We would have returned a 404 if the user had access to the app, but the app
			// did not exist. So we have to return a 404 when the app does exist, but the user does not have access.
			// Otherwise, they could infer that the app exists based on the error code.
			return nil, status.Errorf(codes.NotFound, apierr.NewNotFound(schema.GroupResource{Group: "argoproj.io", Resource: "applications"}, name).Error())
		}
		// The user didn't specify a project. We always return permission denied for both lack of access and lack of
		// existence.
		return nil, permissionDeniedErr
	}
	effectiveProject := "default"
	if a.Spec.Project != "" {
		}).Warnf("user tried to %s application in project %s, but the application is in project %s", action, project, effectiveProject)
		// The user has access to the app, but the app is in a different project. Return 404, meaning "app doesn't
		// exist in that project".
		return nil, status.Errorf(codes.NotFound, apierr.NewNotFound(schema.GroupResource{Group: "argoproj.io", Resource: "applications"}, name).Error())
	}
	return a, nil
}

// getApplicationEnforceRBACInformer uses an informer to get an Application. If the app does not exist, permission is
// denied, or any other error occurs when getting the app, we return a permission denied error to obscure any sensitive
// information.
func (s *Server) getApplicationEnforceRBACInformer(ctx context.Context, action, project, namespace, name string) (*appv1.Application, error) {
	namespaceOrDefault := s.appNamespaceOrDefault(namespace)
	return s.getAppEnforceRBAC(ctx, action, project, namespaceOrDefault, name, func() (*appv1.Application, error) {
		return s.appLister.Applications(namespaceOrDefault).Get(name)
// getApplicationEnforceRBACClient uses a client to get an Application. If the app does not exist, permission is denied,
// or any other error occurs when getting the app, we return a permission denied error to obscure any sensitive
// information.
func (s *Server) getApplicationEnforceRBACClient(ctx context.Context, action, project, namespace, name, resourceVersion string) (*appv1.Application, error) {
	namespaceOrDefault := s.appNamespaceOrDefault(namespace)
	return s.getAppEnforceRBAC(ctx, action, project, namespaceOrDefault, name, func() (*appv1.Application, error) {
		if !s.isNamespaceEnabled(namespaceOrDefault) {
	if q.Validate != nil {
		validate = *q.Validate
	}
	err := s.validateAndNormalizeApp(ctx, a, validate)
	if err != nil {
		return nil, fmt.Errorf("error while validating and normalizing app: %w", err)
	}
	return updated, nil
}

func (s *Server) queryRepoServer(ctx context.Context, a *appv1.Application, action func(
	client apiclient.RepoServerServiceClient,
	repo *appv1.Repository,
	helmRepos []*appv1.Repository,
	if err != nil {
		return fmt.Errorf("error getting kustomize settings options: %w", err)
	}
	proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			return status.Errorf(codes.InvalidArgument, "application references project %s which does not exist", a.Spec.Project)
		}
		return fmt.Errorf("error getting application's project: %w", err)
	}

	helmRepos, err := s.db.ListHelmRepositories(ctx)
	if err != nil {
	if q.Name == nil || *q.Name == "" {
		return nil, fmt.Errorf("invalid request: application name is missing")
	}
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, err
	}
	}

	var manifestInfo *apiclient.ManifestResponse
	err = s.queryRepoServer(ctx, a, func(
		client apiclient.RepoServerServiceClient, repo *appv1.Repository, helmRepos []*appv1.Repository, helmCreds []*appv1.RepoCreds, helmOptions *appv1.HelmOptions, kustomizeOptions *appv1.KustomizeOptions, enableGenerateManifests map[string]bool) error {
		revision := source.TargetRevision
		if q.GetRevision() != "" {
			return fmt.Errorf("error getting API resources: %w", err)
		}

		proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
		if err != nil {
			return fmt.Errorf("error getting app project: %w", err)
		}

		manifestInfo, err = client.GenerateManifest(ctx, &apiclient.ManifestRequest{
			Repo:               repo,
			Revision:           revision,
		return fmt.Errorf("invalid request: application name is missing")
	}

	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, query.GetProject(), query.GetAppNamespace(), query.GetName())
	if err != nil {
		return err
	}

	var manifestInfo *apiclient.ManifestResponse
	err = s.queryRepoServer(ctx, a, func(
		client apiclient.RepoServerServiceClient, repo *appv1.Repository, helmRepos []*appv1.Repository, helmCreds []*appv1.RepoCreds, helmOptions *appv1.HelmOptions, kustomizeOptions *appv1.KustomizeOptions, enableGenerateManifests map[string]bool) error {

		appInstanceLabelKey, err := s.settingsMgr.GetAppInstanceLabelKey()
	// We must use a client Get instead of an informer Get, because it's common to call Get immediately
	// following a Watch (which is not yet powered by an informer), and the Get must reflect what was
	// previously seen by the client.
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, project, appNs, appName, q.GetResourceVersion())
	if err != nil {
		return nil, err
	}

	if refreshType == appv1.RefreshTypeHard {
		// force refresh cached application details
		if err := s.queryRepoServer(ctx, a, func(
			client apiclient.RepoServerServiceClient,
			repo *appv1.Repository,
			helmRepos []*appv1.Repository,

// ListResourceEvents returns a list of event resources
func (s *Server) ListResourceEvents(ctx context.Context, q *application.ApplicationResourceEventsQuery) (*v1.EventList, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, err
	}
	s.projectLock.RLock(newApp.Spec.GetProject())
	defer s.projectLock.RUnlock(newApp.Spec.GetProject())

	app, err := s.getApplicationEnforceRBACClient(ctx, action, currentProject, newApp.Namespace, newApp.Name, "")
	if err != nil {
		return nil, err
	}

	err = s.validateAndNormalizeApp(ctx, newApp, validate)
	if err != nil {
		return nil, fmt.Errorf("error validating and normalizing app: %w", err)
	}
	if q.GetSpec() == nil {
		return nil, fmt.Errorf("error updating application spec: spec is nil in request")
	}
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionUpdate, q.GetProject(), q.GetAppNamespace(), q.GetName(), "")
	if err != nil {
		return nil, err
	}

// Patch patches an application
func (s *Server) Patch(ctx context.Context, q *application.ApplicationPatchRequest) (*appv1.Application, error) {
	app, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName(), "")
	if err != nil {
		return nil, err
	}
	return s.validateAndUpdateApp(ctx, newApp, false, true, rbacpolicy.ActionUpdate, q.GetProject())
}

// Delete removes an application and all associated resources
func (s *Server) Delete(ctx context.Context, q *application.ApplicationDeleteRequest) (*application.ApplicationResponse, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, q.GetProject(), appNs, appName, "")
	if err != nil {
		return nil, err
	}
	}
}

func (s *Server) validateAndNormalizeApp(ctx context.Context, app *appv1.Application, validate bool) error {
	proj, err := argo.GetAppProject(app, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			// Offer no hint that the project does not exist.
			log.Warnf("User attempted to create/update application in non-existent project %q", app.Spec.Project)
			return permissionDeniedErr
		}
		return fmt.Errorf("error getting application's project: %w", err)
	}
	if app.GetName() == "" {
		return fmt.Errorf("resource name may not be empty")
	}
}

func (s *Server) getAppLiveResource(ctx context.Context, action string, q *application.ApplicationResourceRequest) (*appv1.ResourceNode, *rest.Config, *appv1.Application, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, action, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, nil, nil, err
	}
}

func (s *Server) ResourceTree(ctx context.Context, q *application.ResourcesQuery) (*appv1.ApplicationTree, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetApplicationName())
	if err != nil {
		return nil, err
	}
}

func (s *Server) WatchResourceTree(q *application.ResourcesQuery, ws application.ApplicationService_WatchResourceTreeServer) error {
	_, err := s.getApplicationEnforceRBACInformer(ws.Context(), rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetApplicationName())
	if err != nil {
		return err
	}
}

func (s *Server) RevisionMetadata(ctx context.Context, q *application.RevisionMetadataQuery) (*appv1.RevisionMetadata, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, fmt.Errorf("error getting repository by URL: %w", err)
	}
	// We need to get some information with the project associated to the app,
	// so we'll know whether GPG signatures are enforced.
	proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting app project: %w", err)
	}
	conn, repoClient, err := s.repoClientset.NewRepoServerClient()
	if err != nil {
		return nil, fmt.Errorf("error creating repo server client: %w", err)

// RevisionChartDetails returns the helm chart metadata, as fetched from the reposerver
func (s *Server) RevisionChartDetails(ctx context.Context, q *application.RevisionMetadataQuery) (*appv1.ChartDetails, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return nil, err
	}
}

func (s *Server) ManagedResources(ctx context.Context, q *application.ResourcesQuery) (*application.ManagedResourcesResponse, error) {
	a, err := s.getApplicationEnforceRBACInformer(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetApplicationName())
	if err != nil {
		return nil, err
	}
		}
	}

	a, err := s.getApplicationEnforceRBACInformer(ws.Context(), rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName())
	if err != nil {
		return err
	}

// Sync syncs an application to its target state
func (s *Server) Sync(ctx context.Context, syncReq *application.ApplicationSyncRequest) (*appv1.Application, error) {
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, syncReq.GetProject(), syncReq.GetAppNamespace(), syncReq.GetName(), "")
	if err != nil {
		return nil, err
	}

	proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			return a, status.Errorf(codes.InvalidArgument, "application references project %s which does not exist", a.Spec.Project)
		}
		return a, fmt.Errorf("error getting app project: %w", err)
	}

	s.inferResourcesStatusHealth(a)

	if !proj.Spec.SyncWindows.Matches(a).CanSync(true) {
}

func (s *Server) Rollback(ctx context.Context, rollbackReq *application.ApplicationRollbackRequest) (*appv1.Application, error) {
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionSync, rollbackReq.GetProject(), rollbackReq.GetAppNamespace(), rollbackReq.GetName(), "")
	if err != nil {
		return nil, err
	}
}

func (s *Server) ListLinks(ctx context.Context, req *application.ListAppLinksRequest) (*application.LinksResponse, error) {
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, req.GetProject(), req.GetNamespace(), req.GetName(), "")
	if err != nil {
		return nil, err
	}
		return nil, fmt.Errorf("failed to read application deep links from configmap: %w", err)
	}

	clstObj, _, err := s.getObjectsForDeepLinks(ctx, a)
	if err != nil {
		return nil, err
	}
	return finalList, nil
}

func (s *Server) getObjectsForDeepLinks(ctx context.Context, app *appv1.Application) (cluster *unstructured.Unstructured, project *unstructured.Unstructured, err error) {
	proj, err := argo.GetAppProject(app, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("error getting app project: %w", err)
	}

	// sanitize project jwt tokens
	proj.Status = appv1.AppProjectStatus{}

		return nil, err
	}

	clstObj, projObj, err := s.getObjectsForDeepLinks(ctx, app)
	if err != nil {
		return nil, err
	}
func (s *Server) TerminateOperation(ctx context.Context, termOpReq *application.OperationTerminateRequest) (*application.OperationTerminateResponse, error) {
	appName := termOpReq.GetName()
	appNs := s.appNamespaceOrDefault(termOpReq.GetAppNamespace())
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionSync, termOpReq.GetProject(), appNs, appName, "")
	if err != nil {
		return nil, err
	}

func (s *Server) getUnstructuredLiveResourceOrApp(ctx context.Context, rbacRequest string, q *application.ApplicationResourceRequest) (obj *unstructured.Unstructured, res *appv1.ResourceNode, app *appv1.Application, config *rest.Config, err error) {
	if q.GetKind() == applicationType.ApplicationKind && q.GetGroup() == applicationType.Group && q.GetName() == q.GetResourceName() {
		app, err = s.getApplicationEnforceRBACInformer(ctx, rbacRequest, q.GetProject(), q.GetAppNamespace(), q.GetName())
		if err != nil {
			return nil, nil, nil, nil, err
		}
		}
	}

	// First, make sure all the returned resources are permitted, for each operation.
	// Also perform create with dry-runs for all create-operation resources.
	// This is performed separately to reduce the risk of only some of the resources being successfully created later.
	// the dry-run for relevant apply/delete operation would have to be invoked as well.
	for _, impactedResource := range newObjects {
		newObj := impactedResource.UnstructuredObj
		err := s.verifyResourcePermitted(ctx, app, newObj)
		if err != nil {
			return nil, err
		}
	return &application.ApplicationResponse{}, nil
}

func (s *Server) verifyResourcePermitted(ctx context.Context, app *appv1.Application, obj *unstructured.Unstructured) error {
	proj, err := argo.GetAppProject(app, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			return fmt.Errorf("application references project %s which does not exist", app.Spec.Project)
		}
		return fmt.Errorf("failed to get project %s: %w", app.Spec.Project, err)
	}
	permitted, err := proj.IsResourcePermitted(schema.GroupKind{Group: obj.GroupVersionKind().Group, Kind: obj.GroupVersionKind().Kind}, obj.GetNamespace(), app.Spec.Destination, func(project string) ([]*appv1.Cluster, error) {
		clusters, err := s.db.GetProjectClusters(context.TODO(), project)
		if err != nil {
}

func (s *Server) GetApplicationSyncWindows(ctx context.Context, q *application.ApplicationSyncWindowsQuery) (*application.ApplicationSyncWindowsResponse, error) {
	a, err := s.getApplicationEnforceRBACClient(ctx, rbacpolicy.ActionGet, q.GetProject(), q.GetAppNamespace(), q.GetName(), "")
	if err != nil {
		return nil, err
	}

	proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting app project: %w", err)
	}

	windows := proj.Spec.SyncWindows.Matches(a)
	sync := windows.CanSync(true)

		appServer := newTestAppServer(t, testApp)

		active, err := appServer.GetApplicationSyncWindows(context.Background(), &application.ApplicationSyncWindowsQuery{Name: &testApp.Name})
		assert.Contains(t, err.Error(), "not found")
		assert.Nil(t, active)
	})
}
	t.Run("Get application in other namespace when allowed", func(t *testing.T) {
		testApp := newTestApp()
		testApp.Namespace = "argocd-1"
		appServer := newTestAppServer(t, testApp)
		appServer.enabledNamespaces = []string{"argocd-1"}
		app, err := appServer.Get(context.TODO(), &application.ApplicationQuery{
			Name:         pointer.String("test-app"),
		require.Equal(t, "argocd-1", app.Namespace)
		require.Equal(t, "test-app", app.Name)
	})
	t.Run("Create application in other namespace when allowed", func(t *testing.T) {
		testApp := newTestApp()
		testApp.Namespace = "argocd-1"
		})
		require.Error(t, err)
		require.Nil(t, app)
		require.ErrorContains(t, err, "not allowed to use project")
	})

	t.Run("Create application in other namespace when not allowed by configuration", func(t *testing.T) {
		require.Nil(t, app)
		require.ErrorContains(t, err, "namespace 'argocd-1' is not permitted")
	})

}
		return nil, err
	}
	if !proj.IsAppNamespacePermitted(app, ns) {
		return nil, fmt.Errorf("application '%s' in namespace '%s' is not allowed to use project '%s'",
			app.Name, app.Namespace, proj.Name)
	}
	return proj, nil
}
	return expiredCount
}

// Protect admin user from login attempt reset caused by attempts to overflow cache in a brute force attack. Instead remove random non-admin to make room in cache. 
func pickRandomNonAdminLoginFailure(failures map[string]LoginAttempts, username string) *string {
	idx := rand.Intn(len(failures) - 1)
	i := 0
