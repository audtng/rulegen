package main


				// Get app before creating to see if it is being updated or no change
				existing, err := appIf.Get(ctx, &applicationpkg.ApplicationQuery{Name: &app.Name})
				if grpc.UnwrapGRPCStatus(err).Code() != codes.NotFound {
					errors.CheckError(err)
				}

)

var (
	watchAPIBufferSize = env.ParseNumFromEnv(argocommon.EnvWatchAPIBufferSize, 1000, 0, math.MaxInt32)
)

// Server provides an Application service
	appclientset      appclientset.Interface
	appLister         applisters.ApplicationLister
	appInformer       cache.SharedIndexInformer
	appBroadcaster    *broadcasterHandler
	repoClientset     apiclient.Clientset
	kubectl           kube.Kubectl
	db                db.ArgoDB
	appclientset appclientset.Interface,
	appLister applisters.ApplicationLister,
	appInformer cache.SharedIndexInformer,
	repoClientset apiclient.Clientset,
	cache *servercache.Cache,
	kubectl kube.Kubectl,
	projInformer cache.SharedIndexInformer,
	enabledNamespaces []string,
) (application.ApplicationServiceServer, AppResourceTreeFn) {
	appBroadcaster := &broadcasterHandler{}
	appInformer.AddEventHandler(appBroadcaster)
	s := &Server{
		ns:                namespace,
	return s, s.getAppResources
}

// List returns list of applications
func (s *Server) List(ctx context.Context, q *application.ApplicationQuery) (*appv1.ApplicationList, error) {
	selector, err := labels.Parse(q.GetSelector())
	if q.Name == nil || *q.Name == "" {
		return nil, fmt.Errorf("invalid request: application name is missing")
	}
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

		return fmt.Errorf("invalid request: application name is missing")
	}

	appName := query.GetName()
	appNs := s.appNamespaceOrDefault(query.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)

	if err != nil {
		return fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return err
	}

	// We must use a client Get instead of an informer Get, because it's common to call Get immediately
	// following a Watch (which is not yet powered by an informer), and the Get must reflect what was
	// previously seen by the client.
	a, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{
		ResourceVersion: q.GetResourceVersion(),
	})

	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}


// ListResourceEvents returns a list of event resources
func (s *Server) ListResourceEvents(ctx context.Context, q *application.ApplicationResourceEventsQuery) (*v1.EventList, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

	return list, nil
}

func (s *Server) validateAndUpdateApp(ctx context.Context, newApp *appv1.Application, merge bool, validate bool) (*appv1.Application, error) {
	s.projectLock.RLock(newApp.Spec.GetProject())
	defer s.projectLock.RUnlock(newApp.Spec.GetProject())

	app, err := s.appclientset.ArgoprojV1alpha1().Applications(newApp.Namespace).Get(ctx, newApp.Name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}

	err = s.validateAndNormalizeApp(ctx, newApp, validate)
	if q.Validate != nil {
		validate = *q.Validate
	}
	return s.validateAndUpdateApp(ctx, q.Application, false, validate)
}

// UpdateSpec updates an application spec and filters out any invalid parameter overrides
	if q.GetSpec() == nil {
		return nil, fmt.Errorf("error updating application spec: spec is nil in request")
	}
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionUpdate, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

	if q.Validate != nil {
		validate = *q.Validate
	}
	a, err = s.validateAndUpdateApp(ctx, a, false, validate)
	if err != nil {
		return nil, fmt.Errorf("error validating and updating app: %w", err)
	}

// Patch patches an application
func (s *Server) Patch(ctx context.Context, q *application.ApplicationPatchRequest) (*appv1.Application, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	app, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}

	if err = s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionUpdate, app.RBACName(s.ns)); err != nil {
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling patched app: %w", err)
	}
	return s.validateAndUpdateApp(ctx, newApp, false, true)
}

// Delete removes an application and all associated resources
func (s *Server) Delete(ctx context.Context, q *application.ApplicationDeleteRequest) (*application.ApplicationResponse, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}

	s.projectLock.RLock(a.Spec.Project)
	proj, err := argo.GetAppProject(app, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
	if err != nil {
		if apierr.IsNotFound(err) {
			return status.Errorf(codes.InvalidArgument, "application references project %s which does not exist", app.Spec.Project)
		}
		return fmt.Errorf("error getting application's project: %w", err)
	}
		return s.cache.GetAppResourcesTree(a.InstanceName(s.ns), &tree)
	})
	if err != nil {
		return &tree, fmt.Errorf("error getting cached app state: %w", err)
	}
	return &tree, nil
}

func (s *Server) getAppLiveResource(ctx context.Context, action string, q *application.ApplicationResourceRequest) (*appv1.ResourceNode, *rest.Config, *appv1.Application, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting app by name: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, action, a.RBACName(s.ns)); err != nil {
		return nil, nil, nil, err
	}

	tree, err := s.getAppResources(ctx, a)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("error getting app resources: %w", err)
func (s *Server) GetResource(ctx context.Context, q *application.ApplicationResourceRequest) (*application.ApplicationResourceResponse, error) {
	res, config, _, err := s.getAppLiveResource(ctx, rbacpolicy.ActionGet, q)
	if err != nil {
		return nil, fmt.Errorf("error getting app live resource: %w", err)
	}

	// make sure to use specified resource version if provided
	}
	res, config, a, err := s.getAppLiveResource(ctx, rbacpolicy.ActionUpdate, resourceRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting app live resource: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionUpdate, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

		}
		return nil, fmt.Errorf("error patching resource: %w", err)
	}
	manifest, err = replaceSecretValues(manifest)
	if err != nil {
		return nil, fmt.Errorf("error replacing secret values: %w", err)
	}
	res, config, a, err := s.getAppLiveResource(ctx, rbacpolicy.ActionDelete, resourceRequest)
	if err != nil {
		return nil, fmt.Errorf("error getting live resource for delete: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionDelete, a.RBACName(s.ns)); err != nil {
		return nil, err
	}
	var deleteOption metav1.DeleteOptions
}

func (s *Server) ResourceTree(ctx context.Context, q *application.ResourcesQuery) (*appv1.ApplicationTree, error) {
	appName := q.GetApplicationName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, fmt.Errorf("error getting application by name: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

}

func (s *Server) WatchResourceTree(q *application.ResourcesQuery, ws application.ApplicationService_WatchResourceTreeServer) error {
	appName := q.GetApplicationName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return fmt.Errorf("error getting application by name: %w", err)
	}

	if err := s.enf.EnforceErr(ws.Context().Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return err
	}

}

func (s *Server) RevisionMetadata(ctx context.Context, q *application.RevisionMetadataQuery) (*appv1.RevisionMetadata, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, fmt.Errorf("error getting app by name: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

}

func (s *Server) ManagedResources(ctx context.Context, q *application.ResourcesQuery) (*application.ManagedResourcesResponse, error) {
	appName := q.GetApplicationName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return nil, fmt.Errorf("error getting application: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, fmt.Errorf("error verifying rbac: %w", err)
	}

	items := make([]*appv1.ResourceDiff, 0)
		return s.cache.GetAppManagedResources(a.InstanceName(s.ns), &items)
	})
	if err != nil {
		return nil, fmt.Errorf("error getting cached app state: %w", err)
	}
	res := &application.ManagedResourcesResponse{}
	for i := range items {
		}
	}

	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	a, err := s.appLister.Applications(appNs).Get(appName)
	if err != nil {
		return fmt.Errorf("error getting application by name: %w", err)
	}

	if err := s.enf.EnforceErr(ws.Context().Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return err
	}


// Sync syncs an application to its target state
func (s *Server) Sync(ctx context.Context, syncReq *application.ApplicationSyncRequest) (*appv1.Application, error) {
	appName := syncReq.GetName()
	appNs := s.appNamespaceOrDefault(syncReq.GetAppNamespace())
	appIf := s.appclientset.ArgoprojV1alpha1().Applications(appNs)
	a, err := appIf.Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application by name: %w", err)
	}

	proj, err := argo.GetAppProject(a, applisters.NewAppProjectLister(s.projInformer.GetIndexer()), s.ns, s.settingsMgr, s.db, ctx)
		op.Retry = *retry
	}

	a, err = argo.SetAppOperation(appIf, appName, &op)
	if err != nil {
		return nil, fmt.Errorf("error setting app operation: %w", err)
}

func (s *Server) Rollback(ctx context.Context, rollbackReq *application.ApplicationRollbackRequest) (*appv1.Application, error) {
	appName := rollbackReq.GetName()
	appNs := s.appNamespaceOrDefault(rollbackReq.GetAppNamespace())
	appIf := s.appclientset.ArgoprojV1alpha1().Applications(appNs)
	a, err := appIf.Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application by name: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionSync, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

		},
		InitiatedBy: appv1.OperationInitiator{Username: session.Username(ctx)},
	}
	a, err = argo.SetAppOperation(appIf, appName, &op)
	if err != nil {
		return nil, fmt.Errorf("error setting app operation: %w", err)
}

func (s *Server) ListLinks(ctx context.Context, req *application.ListAppLinksRequest) (*application.LinksResponse, error) {
	appName := req.GetName()
	appNs := s.appNamespaceOrDefault(req.GetNamespace())

	a, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		log.WithFields(map[string]interface{}{
			"application": appName,
			"ns":          appNs,
		}).Errorf("failed to get application, error=%v", err.Error())
		return nil, fmt.Errorf("error getting application")
	}

	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		log.WithFields(map[string]interface{}{
			"application": appName,
			"ns":          appNs,
		}).Warnf("unauthorized access to app, error=%v", err.Error())
		return nil, fmt.Errorf("error getting application")
	}

	obj, err := kube.ToUnstructured(a)
func (s *Server) TerminateOperation(ctx context.Context, termOpReq *application.OperationTerminateRequest) (*application.OperationTerminateResponse, error) {
	appName := termOpReq.GetName()
	appNs := s.appNamespaceOrDefault(termOpReq.GetAppNamespace())
	a, err := s.appclientset.ArgoprojV1alpha1().Applications(appNs).Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application by name: %w", err)
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionSync, a.RBACName(s.ns)); err != nil {
		return nil, err
	}


func (s *Server) getUnstructuredLiveResourceOrApp(ctx context.Context, rbacRequest string, q *application.ApplicationResourceRequest) (obj *unstructured.Unstructured, res *appv1.ResourceNode, app *appv1.Application, config *rest.Config, err error) {
	if q.GetKind() == "Application" && q.GetGroup() == "argoproj.io" && q.GetName() == q.GetResourceName() {
		namespace := s.appNamespaceOrDefault(q.GetAppNamespace())
		app, err = s.appLister.Applications(namespace).Get(q.GetName())
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting app by name: %w", err)
		}
		if err = s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacRequest, app.RBACName(s.ns)); err != nil {
			return nil, nil, nil, nil, err
	} else {
		res, config, app, err = s.getAppLiveResource(ctx, rbacRequest, q)
		if err != nil {
			return nil, nil, nil, nil, fmt.Errorf("error getting app live resource: %w", err)
		}
		obj, err = s.kubectl.GetResource(ctx, config, res.GroupKindVersion(), res.Name, res.Namespace)
	}
}

func (s *Server) GetApplicationSyncWindows(ctx context.Context, q *application.ApplicationSyncWindowsQuery) (*application.ApplicationSyncWindowsResponse, error) {
	appName := q.GetName()
	appNs := s.appNamespaceOrDefault(q.GetAppNamespace())
	appIf := s.appclientset.ArgoprojV1alpha1().Applications(appNs)
	a, err := appIf.Get(ctx, appName, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("error getting application by name: %w", err)
	}

	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceApplications, rbacpolicy.ActionGet, a.RBACName(s.ns)); err != nil {
		return nil, err
	}

	"context"
	coreerrors "errors"
	"fmt"
	"sync/atomic"
	"testing"
	"time"

	"github.com/argoproj/gitops-engine/pkg/health"
	synccommon "github.com/argoproj/gitops-engine/pkg/sync/common"
	"github.com/argoproj/gitops-engine/pkg/utils/kube/kubetest"
	"github.com/argoproj/pkg/sync"
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes/fake"
	kubetesting "k8s.io/client-go/testing"
	k8scache "k8s.io/client-go/tools/cache"
	"k8s.io/utils/pointer"
	appinformer "github.com/argoproj/argo-cd/v2/pkg/client/informers/externalversions"
	"github.com/argoproj/argo-cd/v2/reposerver/apiclient"
	"github.com/argoproj/argo-cd/v2/reposerver/apiclient/mocks"
	servercache "github.com/argoproj/argo-cd/v2/server/cache"
	"github.com/argoproj/argo-cd/v2/server/rbacpolicy"
	"github.com/argoproj/argo-cd/v2/test"
	mockRepoServiceClient.On("GenerateManifest", mock.Anything, mock.Anything).Return(&apiclient.ManifestResponse{}, nil)
	mockRepoServiceClient.On("GetAppDetails", mock.Anything, mock.Anything).Return(&apiclient.RepoAppDetailsResponse{}, nil)
	mockRepoServiceClient.On("TestRepository", mock.Anything, mock.Anything).Return(&apiclient.TestRepositoryResponse{}, nil)

	if isHelm {
		mockRepoServiceClient.On("ResolveRevision", mock.Anything, mock.Anything).Return(fakeResolveRevesionResponseHelm(), nil)
}

// return an ApplicationServiceServer which returns fake data
func newTestAppServer(objects ...runtime.Object) *Server {
	f := func(enf *rbac.Enforcer) {
		_ = enf.SetBuiltinPolicy(assets.BuiltinPolicyCSV)
		enf.SetDefaultRole("role:admin")
	}
	return newTestAppServerWithEnforcerConfigure(f, objects...)
}

func newTestAppServerWithEnforcerConfigure(f func(*rbac.Enforcer), objects ...runtime.Object) *Server {
	kubeclientset := fake.NewSimpleClientset(&v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: testNamespace,
		panic("Timed out waiting for caches to sync")
	}

	server, _ := NewServer(
		testNamespace,
		kubeclientset,
		fakeAppsClientset,
		factory.Argoproj().V1alpha1().Applications().Lister(),
		appInformer,
		mockRepoClient,
		nil,
		&kubetest.MockKubectlCmd{},
		db,
		enforcer,
		sync.NewKeyLock(),
	return &app
}

func TestListAppsInNamespaceWithLabels(t *testing.T) {
	appServer := newTestAppServer(newTestApp(func(app *appsv1.Application) {
		app.Name = "App1"
		app.ObjectMeta.Namespace = "test-namespace"
		app.SetLabels(map[string]string{"key1": "value1", "key2": "value1"})
}

func TestListAppsInDefaultNSWithLabels(t *testing.T) {
	appServer := newTestAppServer(newTestApp(func(app *appsv1.Application) {
		app.Name = "App1"
		app.SetLabels(map[string]string{"key1": "value1", "key2": "value1"})
	}), newTestApp(func(app *appsv1.Application) {
}

func TestListAppWithProjects(t *testing.T) {
	appServer := newTestAppServer(newTestApp(func(app *appsv1.Application) {
		app.Name = "App1"
		app.Spec.Project = "test-project1"
	}), newTestApp(func(app *appsv1.Application) {
}

func TestListApps(t *testing.T) {
	appServer := newTestAppServer(newTestApp(func(app *appsv1.Application) {
		app.Name = "bcd"
	}), newTestApp(func(app *appsv1.Application) {
		app.Name = "abc"
`
		_ = enf.SetUserPolicy(policy)
	}
	appServer := newTestAppServerWithEnforcerConfigure(f, objects...)

	res, err := appServer.List(ctx, &application.ApplicationQuery{})


func TestCreateApp(t *testing.T) {
	testApp := newTestApp()
	appServer := newTestAppServer()
	testApp.Spec.Project = ""
	createReq := application.ApplicationCreateRequest{
		Application: testApp,
}

func TestCreateAppWithDestName(t *testing.T) {
	appServer := newTestAppServer()
	testApp := newTestAppWithDestName()
	createReq := application.ApplicationCreateRequest{
		Application: testApp,

func TestUpdateApp(t *testing.T) {
	testApp := newTestApp()
	appServer := newTestAppServer(testApp)
	testApp.Spec.Project = ""
	app, err := appServer.Update(context.Background(), &application.ApplicationUpdateRequest{
		Application: testApp,

func TestUpdateAppSpec(t *testing.T) {
	testApp := newTestApp()
	appServer := newTestAppServer(testApp)
	testApp.Spec.Project = ""
	spec, err := appServer.UpdateSpec(context.Background(), &application.ApplicationUpdateSpecRequest{
		Name: &testApp.Name,

func TestDeleteApp(t *testing.T) {
	ctx := context.Background()
	appServer := newTestAppServer()
	createReq := application.ApplicationCreateRequest{
		Application: newTestApp(),
	}
	})
}

func TestDeleteApp_InvalidName(t *testing.T) {
	appServer := newTestAppServer()
	_, err := appServer.Delete(context.Background(), &application.ApplicationDeleteRequest{
		Name: pointer.StringPtr("foo"),
	})
	if !assert.Error(t, err) {
		return
	}
	assert.True(t, apierrors.IsNotFound(err))
}

func TestSyncAndTerminate(t *testing.T) {
	ctx := context.Background()
	appServer := newTestAppServer()
	testApp := newTestApp()
	testApp.Spec.Source.RepoURL = "https://github.com/argoproj/argo-cd.git"
	createReq := application.ApplicationCreateRequest{

func TestSyncHelm(t *testing.T) {
	ctx := context.Background()
	appServer := newTestAppServer()
	testApp := newTestApp()
	testApp.Spec.Source.RepoURL = "https://argoproj.github.io/argo-helm"
	testApp.Spec.Source.Path = ""

func TestSyncGit(t *testing.T) {
	ctx := context.Background()
	appServer := newTestAppServer()
	testApp := newTestApp()
	testApp.Spec.Source.RepoURL = "https://github.com/org/test"
	testApp.Spec.Source.Path = "deploy"
		Revision: "abc",
		Source:   *testApp.Spec.Source.DeepCopy(),
	}}
	appServer := newTestAppServer(testApp)

	updatedApp, err := appServer.Rollback(context.Background(), &application.ApplicationRollbackRequest{
		Name: &testApp.Name,
	ctx := context.Background()
	// nolint:staticcheck
	ctx = context.WithValue(ctx, "claims", &jwt.StandardClaims{Subject: "admin"})
	appServer := newTestAppServer(testApp)
	appServer.enf.SetDefaultRole("")

	// Verify normal update works (without changing project)
	_ = appServer.enf.SetBuiltinPolicy(`p, admin, applications, update, default/test-app, allow`)
	_, err := appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	assert.NoError(t, err)

	// Verify caller cannot update to another project
	testApp.Spec.Project = "my-proj"
	_, err = appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	assert.Equal(t, status.Code(err), codes.PermissionDenied)

	// Verify inability to change projects without create privileges in new project
	_ = appServer.enf.SetBuiltinPolicy(`
p, admin, applications, update, default/test-app, allow
p, admin, applications, update, my-proj/test-app, allow
`)
	_, err = appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	statusErr := grpc.UnwrapGRPCStatus(err)
	assert.NotNil(t, statusErr)
	assert.Equal(t, codes.PermissionDenied, statusErr.Code())

	// Verify inability to change projects without update privileges in new project
	_ = appServer.enf.SetBuiltinPolicy(`
p, admin, applications, update, default/test-app, allow
p, admin, applications, create, my-proj/test-app, allow
`)
	_, err = appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	assert.Equal(t, status.Code(err), codes.PermissionDenied)

	// Verify inability to change projects without update privileges in old project
	_ = appServer.enf.SetBuiltinPolicy(`
p, admin, applications, create, my-proj/test-app, allow
p, admin, applications, update, my-proj/test-app, allow
`)
	_, err = appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	statusErr = grpc.UnwrapGRPCStatus(err)
	assert.NotNil(t, statusErr)
	assert.Equal(t, codes.PermissionDenied, statusErr.Code())

	// Verify can update project with proper permissions
	_ = appServer.enf.SetBuiltinPolicy(`
p, admin, applications, update, default/test-app, allow
p, admin, applications, create, my-proj/test-app, allow
p, admin, applications, update, my-proj/test-app, allow
`)
	updatedApp, err := appServer.Update(ctx, &application.ApplicationUpdateRequest{Application: testApp})
	assert.NoError(t, err)
	assert.Equal(t, "my-proj", updatedApp.Spec.Project)
}

func TestAppJsonPatch(t *testing.T) {
	ctx := context.Background()
	// nolint:staticcheck
	ctx = context.WithValue(ctx, "claims", &jwt.StandardClaims{Subject: "admin"})
	appServer := newTestAppServer(testApp)
	appServer.enf.SetDefaultRole("")

	app, err := appServer.Patch(ctx, &application.ApplicationPatchRequest{Name: &testApp.Name, Patch: pointer.String("garbage")})
	ctx := context.Background()
	// nolint:staticcheck
	ctx = context.WithValue(ctx, "claims", &jwt.StandardClaims{Subject: "admin"})
	appServer := newTestAppServer(testApp)
	appServer.enf.SetDefaultRole("")

	app, err := appServer.Patch(ctx, &application.ApplicationPatchRequest{
	t.Run("Active", func(t *testing.T) {
		testApp := newTestApp()
		testApp.Spec.Project = "proj-maint"
		appServer := newTestAppServer(testApp)

		active, err := appServer.GetApplicationSyncWindows(context.Background(), &application.ApplicationSyncWindowsQuery{Name: &testApp.Name})
		assert.NoError(t, err)
	t.Run("Inactive", func(t *testing.T) {
		testApp := newTestApp()
		testApp.Spec.Project = "default"
		appServer := newTestAppServer(testApp)

		active, err := appServer.GetApplicationSyncWindows(context.Background(), &application.ApplicationSyncWindowsQuery{Name: &testApp.Name})
		assert.NoError(t, err)
	t.Run("ProjectDoesNotExist", func(t *testing.T) {
		testApp := newTestApp()
		testApp.Spec.Project = "none"
		appServer := newTestAppServer(testApp)

		active, err := appServer.GetApplicationSyncWindows(context.Background(), &application.ApplicationSyncWindowsQuery{Name: &testApp.Name})
		assert.Contains(t, err.Error(), "not found")
			Namespace: testNamespace,
		},
	}
	appServer := newTestAppServer(testApp, testProj)
	fakeClientSet := appServer.appclientset.(*apps.Clientset)
	fakeClientSet.AddReactor("get", "applications", func(action kubetesting.Action) (handled bool, ret runtime.Object, err error) {
		return true, &appsv1.Application{Spec: appsv1.ApplicationSpec{Source: &appsv1.ApplicationSource{}}}, nil
	defer cancel()
	testApp := newTestApp()
	testApp.ObjectMeta.ResourceVersion = "1"
	appServer := newTestAppServer(testApp)

	var patched int32

	defer cancel()
	testApp := newTestApp()
	testApp.ObjectMeta.ResourceVersion = "1"
	appServer := newTestAppServer(testApp)

	var getAppDetailsQuery *apiclient.RepoServerAppDetailsQuery
	mockRepoServiceClient := mocks.RepoServerServiceClient{}
		Name:      "guestbook-stateful",
		Namespace: "default",
	}}
	appServer := newTestAppServer(testApp)
	appStateCache := appstate.NewCache(cacheClient, time.Minute)
	err := appStateCache.SetAppResourcesTree(testApp.Name, &appsv1.ApplicationTree{Nodes: []appsv1.ResourceNode{{
		ResourceRef: appsv1.ResourceRef{
	return true
}

type broadcasterHandler struct {
	lock        sync.Mutex
	subscribers []*subscriber
		a.AppClientset,
		a.appLister,
		a.appInformer,
		a.RepoClientset,
		a.Cache,
		kubectl,
		IgnoreErrors().
		CreateApp().
		Then().
		Expect(Error("", "application references project does-not-exist which does not exist"))
}

func TestNamespacedAppDeletion(t *testing.T) {
		IgnoreErrors().
		CreateApp().
		Then().
		Expect(Error("", "application references project does-not-exist which does not exist"))
}

func TestAppDeletion(t *testing.T) {
