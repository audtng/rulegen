package main

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
				},
			})
		} else {
			err := ctrl.stateCache.IterateHierarchy(a.Spec.Destination.Server, kube.GetResourceKey(live), func(child appv1.ResourceNode, appName string) bool {
				if !proj.IsResourcePermitted(schema.GroupKind{Group: child.ResourceRef.Group, Kind: child.ResourceRef.Kind}, child.Namespace, a.Spec.Destination) {
					return false
				}
				nodes = append(nodes, child)
				return true
			})
			if err != nil {
				return nil, err
	orphanedNodes := make([]appv1.ResourceNode, 0)
	for k := range orphanedNodesMap {
		if k.Namespace != "" && proj.IsGroupKindPermitted(k.GroupKind(), true) && !isKnownOrphanedResourceExclusion(k, proj) {
			err := ctrl.stateCache.IterateHierarchy(a.Spec.Destination.Server, k, func(child appv1.ResourceNode, appName string) bool {
				belongToAnotherApp := false
				if appName != "" {
					if _, exists, err := ctrl.appInformer.GetIndexer().GetByKey(ctrl.namespace + "/" + appName); exists && err == nil {
						belongToAnotherApp = true
					}
				}
				if belongToAnotherApp || !proj.IsResourcePermitted(schema.GroupKind{Group: child.ResourceRef.Group, Kind: child.ResourceRef.Kind}, child.Namespace, a.Spec.Destination) {
					return false
				}
				orphanedNodes = append(orphanedNodes, child)
				return true
			})
			if err != nil {
				return nil, err
		app.Status.Sync.Status = appv1.SyncStatusCodeUnknown
		app.Status.Health.Status = health.HealthStatusUnknown
		ctrl.persistAppStatus(origApp, &app.Status)

		if err := ctrl.cache.SetAppResourcesTree(app.Name, &appv1.ApplicationTree{}); err != nil {
			log.Warnf("failed to set app resource tree: %v", err)
		}
		if err := ctrl.cache.SetAppManagedResources(app.Name, nil); err != nil {
			log.Warnf("failed to set app managed resources tree: %v", err)
		}
		return
	}

	mockStateCache.On("GetClusterCache", mock.Anything).Return(&clusterCacheMock, nil)
	mockStateCache.On("IterateHierarchy", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		key := args[1].(kube.ResourceKey)
		action := args[2].(func(child argoappv1.ResourceNode, appName string) bool)
		appName := ""
		if res, ok := data.namespacedResources[key]; ok {
			appName = res.AppName
		}
		_ = action(argoappv1.ResourceNode{ResourceRef: argoappv1.ResourceRef{Kind: key.Kind, Group: key.Group, Namespace: key.Namespace, Name: key.Name}}, appName)
	}).Return(nil)
	return ctrl
}
	// Returns synced cluster cache
	GetClusterCache(server string) (clustercache.ClusterCache, error)
	// Executes give callback against resource specified by the key and all its children
	IterateHierarchy(server string, key kube.ResourceKey, action func(child appv1.ResourceNode, appName string) bool) error
	// Returns state of live nodes which correspond for target nodes of specified application.
	GetManagedLiveObjs(a *appv1.Application, targetObjs []*unstructured.Unstructured) (map[kube.ResourceKey]*unstructured.Unstructured, error)
	// IterateResources iterates all resource stored in cache
	return clusterInfo.IsNamespaced(gk)
}

func (c *liveStateCache) IterateHierarchy(server string, key kube.ResourceKey, action func(child appv1.ResourceNode, appName string) bool) error {
	clusterInfo, err := c.getSyncedCluster(server)
	if err != nil {
		return err
	}
	clusterInfo.IterateHierarchy(key, func(resource *clustercache.Resource, namespaceResources map[kube.ResourceKey]*clustercache.Resource) bool {
		return action(asResourceNode(resource), getApp(resource, namespaceResources))
	})
	return nil
}
}

// IterateHierarchy provides a mock function with given fields: server, key, action
func (_m *LiveStateCache) IterateHierarchy(server string, key kube.ResourceKey, action func(v1alpha1.ResourceNode, string) bool) error {
	ret := _m.Called(server, key, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, kube.ResourceKey, func(v1alpha1.ResourceNode, string) bool) error); ok {
		r0 = rf(server, key, action)
	} else {
		r0 = ret.Error(0)

// IsLiveResourcePermitted returns whether a live resource found in the cluster is permitted by an AppProject
func (proj AppProject) IsLiveResourcePermitted(un *unstructured.Unstructured, server string, name string) bool {
	return proj.IsResourcePermitted(un.GroupVersionKind().GroupKind(), un.GetNamespace(), ApplicationDestination{Server: server, Name: name})
}

func (proj AppProject) IsResourcePermitted(groupKind schema.GroupKind, namespace string, dest ApplicationDestination) bool {
	if !proj.IsGroupKindPermitted(groupKind, namespace != "") {
		return false
	}
	if namespace != "" {
		return proj.IsDestinationPermitted(ApplicationDestination{Server: dest.Server, Name: dest.Name, Namespace: namespace})
	}
	return true
}
			"involvedObject.namespace": a.Namespace,
		}).String()
	} else {
		tree, err := s.getAppResources(ctx, a)
		if err != nil {
			return nil, err
		}
		found := false
		for _, n := range append(tree.Nodes, tree.OrphanedNodes...) {
			if n.ResourceRef.UID == q.ResourceUID && n.ResourceRef.Name == q.ResourceName && n.ResourceRef.Namespace == q.ResourceNamespace {
				found = true
				break
			}
		}
		if !found {
			return nil, status.Errorf(codes.InvalidArgument, "%s not found as part of application %s", q.ResourceName, *q.Name)
		}

		namespace = q.ResourceNamespace
		var config *rest.Config
		config, err = s.getApplicationClusterConfig(ctx, a)
	return &tree, err
}

func (s *Server) getAppLiveResource(ctx context.Context, action string, q *application.ApplicationResourceRequest) (*appv1.ResourceNode, *rest.Config, *appv1.Application, error) {
	a, err := s.appLister.Get(*q.Name)
	if err != nil {
		return nil, nil, nil, err
	}

	found := tree.FindNode(q.Group, q.Kind, q.Namespace, q.ResourceName)
	if found == nil || found.ResourceRef.UID == "" {
		return nil, nil, nil, status.Errorf(codes.InvalidArgument, "%s %s %s not found as part of application %s", q.Kind, q.Group, q.ResourceName, *q.Name)
	}
	config, err := s.getApplicationClusterConfig(ctx, a)
}

func (s *Server) GetResource(ctx context.Context, q *application.ApplicationResourceRequest) (*application.ApplicationResourceResponse, error) {
	res, config, _, err := s.getAppLiveResource(ctx, rbacpolicy.ActionGet, q)
	if err != nil {
		return nil, err
	}
		Version:      q.Version,
		Group:        q.Group,
	}
	res, config, a, err := s.getAppLiveResource(ctx, rbacpolicy.ActionUpdate, resourceRequest)
	if err != nil {
		return nil, err
	}
		Version:      q.Version,
		Group:        q.Group,
	}
	res, config, a, err := s.getAppLiveResource(ctx, rbacpolicy.ActionDelete, resourceRequest)
	if err != nil {
		return nil, err
	}
	var pods []appv1.ResourceNode
	isTheOneMap := make(map[string]bool)
	for _, treeNode := range treeNodes {
		if treeNode.Kind == kube.PodKind && treeNode.Group == "" && treeNode.UID != "" {
			if isTheSelectedOne(&treeNode, q, treeNodes, isTheOneMap) {
				pods = append(pods, treeNode)
			}
}

func (s *Server) ListResourceActions(ctx context.Context, q *application.ApplicationResourceRequest) (*application.ResourceActionsListResponse, error) {
	res, config, _, err := s.getAppLiveResource(ctx, rbacpolicy.ActionGet, q)
	if err != nil {
		return nil, err
	}
		Group:        q.Group,
	}
	actionRequest := fmt.Sprintf("%s/%s/%s/%s", rbacpolicy.ActionAction, q.Group, q.Kind, q.Action)
	res, config, a, err := s.getAppLiveResource(ctx, actionRequest, resourceRequest)
	if err != nil {
		return nil, err
	}
		Expect(SyncStatusIs(SyncStatusCodeSynced))
}

// assertResourceActions verifies if view/modify resource actions are successful/failing for given application
func assertResourceActions(t *testing.T, appName string, successful bool) {
	assertError := func(err error, message string) {
		if successful {
			assert.NoError(t, err)
		} else {
			if assert.Error(t, err) {
				assert.Contains(t, err.Error(), message)
			}
		}
	}

	closer, cdClient := ArgoCDClientset.NewApplicationClientOrDie()
	defer io.Close(closer)

	deploymentResource, err := KubeClientset.AppsV1().Deployments(DeploymentNamespace()).Get(context.Background(), "guestbook-ui", metav1.GetOptions{})
	require.NoError(t, err)

	logs, err := cdClient.PodLogs(context.Background(), &applicationpkg.ApplicationPodLogsQuery{
		Group: pointer.String("apps"), Kind: pointer.String("Deployment"), Name: &appName, Namespace: DeploymentNamespace(),
	})
	require.NoError(t, err)
	_, err = logs.Recv()
	assertError(err, "EOF")

	expectedError := fmt.Sprintf("Deployment apps guestbook-ui not found as part of application %s", appName)

	_, err = cdClient.ListResourceEvents(context.Background(), &applicationpkg.ApplicationResourceEventsQuery{
		Name: &appName, ResourceName: "guestbook-ui", ResourceNamespace: DeploymentNamespace(), ResourceUID: string(deploymentResource.UID)})
	assertError(err, fmt.Sprintf("%s not found as part of application %s", "guestbook-ui", appName))

	_, err = cdClient.GetResource(context.Background(), &applicationpkg.ApplicationResourceRequest{
		Name: &appName, ResourceName: "guestbook-ui", Namespace: DeploymentNamespace(), Version: "v1", Group: "apps", Kind: "Deployment"})
	assertError(err, expectedError)

	_, err = cdClient.DeleteResource(context.Background(), &applicationpkg.ApplicationResourceDeleteRequest{
		Name: &appName, ResourceName: "guestbook-ui", Namespace: DeploymentNamespace(), Version: "v1", Group: "apps", Kind: "Deployment",
	})
	assertError(err, expectedError)

	_, err = cdClient.RunResourceAction(context.Background(), &applicationpkg.ResourceActionRunRequest{
		Name: &appName, ResourceName: "guestbook-ui", Namespace: DeploymentNamespace(), Version: "v1", Group: "apps", Kind: "Deployment", Action: "restart",
	})
	assertError(err, expectedError)
}

func TestPermissions(t *testing.T) {
	appCtx := Given(t)
	projName := "argo-project"
	projActions := projectFixture.
		Given(t).
		Name(projName).
		When().
		Create()

	sourceError := fmt.Sprintf("application repo %s is not permitted in project 'argo-project'", RepoURL(RepoURLTypeFile))
	destinationError := fmt.Sprintf("application destination {%s %s} is not permitted in project 'argo-project'", KubernetesInternalAPIServerAddr, DeploymentNamespace())

	appCtx.
		Path("guestbook-logs").
		Project(projName).
		When().
		IgnoreErrors().
		// ensure app is not created if project permissions are missing
		CreateApp().
		Then().
		Expect(Error("", sourceError)).
		Expect(Error("", destinationError)).
		When().
		DoNotIgnoreErrors().
		// add missing permissions, create and sync app
		And(func() {
			projActions.AddDestination("*", "*")
			projActions.AddSource("*")
		}).
		CreateApp().
		Sync().
		Then().
		// make sure application resource actiions are successful
		And(func(app *Application) {
			assertResourceActions(t, app.Name, true)
		}).
		When().
		// remove projet permissions and "refresh" app
		And(func() {
			projActions.UpdateProject(func(proj *AppProject) {
				proj.Spec.Destinations = nil
				proj.Spec.SourceRepos = nil
			})
		}).
		Refresh(RefreshTypeNormal).
		Then().
		// ensure app resource tree is empty when source/destination permissions are missing
		Expect(Condition(ApplicationConditionInvalidSpecError, destinationError)).
		Expect(Condition(ApplicationConditionInvalidSpecError, sourceError)).
		And(func(app *Application) {
			closer, cdClient := ArgoCDClientset.NewApplicationClientOrDie()
			defer io.Close(closer)
			tree, err := cdClient.ResourceTree(context.Background(), &applicationpkg.ResourcesQuery{ApplicationName: &app.Name})
			require.NoError(t, err)
			assert.Len(t, tree.Nodes, 0)
			assert.Len(t, tree.OrphanedNodes, 0)
		}).
		When().
		// add missing permissions but deny management of Deployment kind
		And(func() {
			projActions.
				AddDestination("*", "*").
				AddSource("*").
				UpdateProject(func(proj *AppProject) {
					proj.Spec.NamespaceResourceBlacklist = []metav1.GroupKind{{Group: "*", Kind: "Deployment"}}
				})
		}).
		Refresh(RefreshTypeNormal).
		Then().
		// make sure application resource actiions are failing
		And(func(app *Application) {
			assertResourceActions(t, "test-permissions", false)
		})
}

func TestPermissionWithScopedRepo(t *testing.T) {
package project

import (
	"context"

	"github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// this implements the "when" part of given/when/then
	return a
}

func (a *Actions) AddDestination(cluster string, namespace string) *Actions {
	a.runCli("proj", "add-destination", a.context.name, cluster, namespace)
	return a
}

func (a *Actions) AddSource(repo string) *Actions {
	a.runCli("proj", "add-source", a.context.name, repo)
	return a
}

func (a *Actions) UpdateProject(updater func(project *v1alpha1.AppProject)) *Actions {
	proj, err := fixture.AppClientset.ArgoprojV1alpha1().AppProjects(fixture.ArgoCDNamespace).Get(context.TODO(), a.context.name, v1.GetOptions{})
	require.NoError(a.context.t, err)
	updater(proj)
	_, err = fixture.AppClientset.ArgoprojV1alpha1().AppProjects(fixture.ArgoCDNamespace).Update(context.TODO(), proj, v1.UpdateOptions{})
	require.NoError(a.context.t, err)
	return a
}

func (a *Actions) Name(name string) *Actions {
	a.context.name = name
	return a
func (a *Actions) runCli(args ...string) {
	a.context.t.Helper()
	a.lastOutput, a.lastError = fixture.RunCli(args...)
	if !a.ignoreErrors {
		require.Empty(a.context.t, a.lastError)
	}
}
