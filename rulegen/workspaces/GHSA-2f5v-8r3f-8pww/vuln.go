package main

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/labels"
	apiruntime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
				},
			})
		} else {
			err := ctrl.stateCache.IterateHierarchy(a.Spec.Destination.Server, kube.GetResourceKey(live), func(child appv1.ResourceNode, appName string) {
				nodes = append(nodes, child)
			})
			if err != nil {
				return nil, err
	orphanedNodes := make([]appv1.ResourceNode, 0)
	for k := range orphanedNodesMap {
		if k.Namespace != "" && proj.IsGroupKindPermitted(k.GroupKind(), true) && !isKnownOrphanedResourceExclusion(k, proj) {
			err := ctrl.stateCache.IterateHierarchy(a.Spec.Destination.Server, k, func(child appv1.ResourceNode, appName string) {
				belongToAnotherApp := false
				if appName != "" {
					if _, exists, err := ctrl.appInformer.GetIndexer().GetByKey(ctrl.namespace + "/" + appName); exists && err == nil {
						belongToAnotherApp = true
					}
				}
				if !belongToAnotherApp {
					orphanedNodes = append(orphanedNodes, child)
				}
			})
			if err != nil {
				return nil, err
		app.Status.Sync.Status = appv1.SyncStatusCodeUnknown
		app.Status.Health.Status = health.HealthStatusUnknown
		ctrl.persistAppStatus(origApp, &app.Status)
		return
	}

	mockStateCache.On("GetClusterCache", mock.Anything).Return(&clusterCacheMock, nil)
	mockStateCache.On("IterateHierarchy", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		key := args[1].(kube.ResourceKey)
		action := args[2].(func(child argoappv1.ResourceNode, appName string))
		appName := ""
		if res, ok := data.namespacedResources[key]; ok {
			appName = res.AppName
		}
		action(argoappv1.ResourceNode{ResourceRef: argoappv1.ResourceRef{Kind: key.Kind, Group: key.Group, Namespace: key.Namespace, Name: key.Name}}, appName)
	}).Return(nil)
	return ctrl
}
	// Returns synced cluster cache
	GetClusterCache(server string) (clustercache.ClusterCache, error)
	// Executes give callback against resource specified by the key and all its children
	IterateHierarchy(server string, key kube.ResourceKey, action func(child appv1.ResourceNode, appName string)) error
	// Returns state of live nodes which correspond for target nodes of specified application.
	GetManagedLiveObjs(a *appv1.Application, targetObjs []*unstructured.Unstructured) (map[kube.ResourceKey]*unstructured.Unstructured, error)
	// IterateResources iterates all resource stored in cache
	return clusterInfo.IsNamespaced(gk)
}

func (c *liveStateCache) IterateHierarchy(server string, key kube.ResourceKey, action func(child appv1.ResourceNode, appName string)) error {
	clusterInfo, err := c.getSyncedCluster(server)
	if err != nil {
		return err
	}
	clusterInfo.IterateHierarchy(key, func(resource *clustercache.Resource, namespaceResources map[kube.ResourceKey]*clustercache.Resource) {
		action(asResourceNode(resource), getApp(resource, namespaceResources))
	})
	return nil
}
}

// IterateHierarchy provides a mock function with given fields: server, key, action
func (_m *LiveStateCache) IterateHierarchy(server string, key kube.ResourceKey, action func(v1alpha1.ResourceNode, string)) error {
	ret := _m.Called(server, key, action)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, kube.ResourceKey, func(v1alpha1.ResourceNode, string)) error); ok {
		r0 = rf(server, key, action)
	} else {
		r0 = ret.Error(0)

// IsLiveResourcePermitted returns whether a live resource found in the cluster is permitted by an AppProject
func (proj AppProject) IsLiveResourcePermitted(un *unstructured.Unstructured, server string, name string) bool {
	if !proj.IsGroupKindPermitted(un.GroupVersionKind().GroupKind(), un.GetNamespace() != "") {
		return false
	}
	if un.GetNamespace() != "" {
		return proj.IsDestinationPermitted(ApplicationDestination{Server: server, Namespace: un.GetNamespace(), Name: name})
	}
	return true
}
			"involvedObject.namespace": a.Namespace,
		}).String()
	} else {
		namespace = q.ResourceNamespace
		var config *rest.Config
		config, err = s.getApplicationClusterConfig(ctx, a)
	return &tree, err
}

func (s *Server) getAppResource(ctx context.Context, action string, q *application.ApplicationResourceRequest) (*appv1.ResourceNode, *rest.Config, *appv1.Application, error) {
	a, err := s.appLister.Get(*q.Name)
	if err != nil {
		return nil, nil, nil, err
	}

	found := tree.FindNode(q.Group, q.Kind, q.Namespace, q.ResourceName)
	if found == nil {
		return nil, nil, nil, status.Errorf(codes.InvalidArgument, "%s %s %s not found as part of application %s", q.Kind, q.Group, q.ResourceName, *q.Name)
	}
	config, err := s.getApplicationClusterConfig(ctx, a)
}

func (s *Server) GetResource(ctx context.Context, q *application.ApplicationResourceRequest) (*application.ApplicationResourceResponse, error) {
	res, config, _, err := s.getAppResource(ctx, rbacpolicy.ActionGet, q)
	if err != nil {
		return nil, err
	}
		Version:      q.Version,
		Group:        q.Group,
	}
	res, config, a, err := s.getAppResource(ctx, rbacpolicy.ActionUpdate, resourceRequest)
	if err != nil {
		return nil, err
	}
		Version:      q.Version,
		Group:        q.Group,
	}
	res, config, a, err := s.getAppResource(ctx, rbacpolicy.ActionDelete, resourceRequest)
	if err != nil {
		return nil, err
	}
	var pods []appv1.ResourceNode
	isTheOneMap := make(map[string]bool)
	for _, treeNode := range treeNodes {
		if treeNode.Kind == kube.PodKind && treeNode.Group == "" {
			if isTheSelectedOne(&treeNode, q, treeNodes, isTheOneMap) {
				pods = append(pods, treeNode)
			}
}

func (s *Server) ListResourceActions(ctx context.Context, q *application.ApplicationResourceRequest) (*application.ResourceActionsListResponse, error) {
	res, config, _, err := s.getAppResource(ctx, rbacpolicy.ActionGet, q)
	if err != nil {
		return nil, err
	}
		Group:        q.Group,
	}
	actionRequest := fmt.Sprintf("%s/%s/%s/%s", rbacpolicy.ActionAction, q.Group, q.Kind, q.Action)
	res, config, a, err := s.getAppResource(ctx, actionRequest, resourceRequest)
	if err != nil {
		return nil, err
	}
		Expect(SyncStatusIs(SyncStatusCodeSynced))
}

func TestPermissions(t *testing.T) {
	EnsureCleanState(t)
	appName := Name()
	_, err := RunCli("proj", "create", "test")
	assert.NoError(t, err)

	// make sure app cannot be created without permissions in project
	_, err = RunCli("app", "create", appName, "--repo", RepoURL(RepoURLTypeFile),
		"--path", guestbookPath, "--project", "test", "--dest-server", KubernetesInternalAPIServerAddr, "--dest-namespace", DeploymentNamespace())
	assert.Error(t, err)
	sourceError := fmt.Sprintf("application repo %s is not permitted in project 'test'", RepoURL(RepoURLTypeFile))
	destinationError := fmt.Sprintf("application destination {%s %s} is not permitted in project 'test'", KubernetesInternalAPIServerAddr, DeploymentNamespace())

	assert.Contains(t, err.Error(), sourceError)
	assert.Contains(t, err.Error(), destinationError)

	proj, err := AppClientset.ArgoprojV1alpha1().AppProjects(ArgoCDNamespace).Get(context.Background(), "test", metav1.GetOptions{})
	assert.NoError(t, err)

	proj.Spec.Destinations = []ApplicationDestination{{Server: "*", Namespace: "*"}}
	proj.Spec.SourceRepos = []string{"*"}
	proj, err = AppClientset.ArgoprojV1alpha1().AppProjects(ArgoCDNamespace).Update(context.Background(), proj, metav1.UpdateOptions{})
	assert.NoError(t, err)

	// make sure controller report permissions issues in conditions
	_, err = RunCli("app", "create", appName, "--repo", RepoURL(RepoURLTypeFile),
		"--path", guestbookPath, "--project", "test", "--dest-server", KubernetesInternalAPIServerAddr, "--dest-namespace", DeploymentNamespace())
	assert.NoError(t, err)
	defer func() {
		err = AppClientset.ArgoprojV1alpha1().Applications(ArgoCDNamespace).Delete(context.Background(), appName, metav1.DeleteOptions{})
		assert.NoError(t, err)
	}()

	proj.Spec.Destinations = []ApplicationDestination{}
	proj.Spec.SourceRepos = []string{}
	_, err = AppClientset.ArgoprojV1alpha1().AppProjects(ArgoCDNamespace).Update(context.Background(), proj, metav1.UpdateOptions{})
	assert.NoError(t, err)
	time.Sleep(1 * time.Second)
	closer, client, err := ArgoCDClientset.NewApplicationClient()
	assert.NoError(t, err)
	defer io.Close(closer)

	refresh := string(RefreshTypeNormal)
	app, err := client.Get(context.Background(), &applicationpkg.ApplicationQuery{Name: &appName, Refresh: &refresh})
	assert.NoError(t, err)

	destinationErrorExist := false
	sourceErrorExist := false
	for i := range app.Status.Conditions {
		if strings.Contains(app.Status.Conditions[i].Message, destinationError) {
			destinationErrorExist = true
		}
		if strings.Contains(app.Status.Conditions[i].Message, sourceError) {
			sourceErrorExist = true
		}
	}
	assert.True(t, destinationErrorExist)
	assert.True(t, sourceErrorExist)
}

func TestPermissionWithScopedRepo(t *testing.T) {
package project

import (
	"github.com/argoproj/argo-cd/v2/test/e2e/fixture"
)

// this implements the "when" part of given/when/then
	return a
}

func (a *Actions) Name(name string) *Actions {
	a.context.name = name
	return a
func (a *Actions) runCli(args ...string) {
	a.context.t.Helper()
	a.lastOutput, a.lastError = fixture.RunCli(args...)
}
