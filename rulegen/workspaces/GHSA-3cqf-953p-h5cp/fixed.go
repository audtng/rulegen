package main


// Get returns a cluster from a query
func (s *Server) Get(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	c, err := s.getClusterAndVerifyAccess(ctx, q, rbacpolicy.ActionGet)
	if err != nil {
		return nil, err
	}

	return s.toAPIResponse(c), nil
}

	return repo, nil
}

func (s *Server) getClusterAndVerifyAccess(ctx context.Context, q *cluster.ClusterQuery, action string) (*appv1.Cluster, error) {
	c, err := s.getClusterWith403IfNotExist(ctx, q)
	if err != nil {
		return nil, err
	}

	// verify that user can do the specified action inside project where cluster is located
	if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, action, CreateClusterRBACObject(c.Project, c.Server)) {
		log.WithField("cluster", q.Server).Warnf("encountered permissions issue while processing request: %v", err)
		return nil, common.PermissionDeniedAPIError
	}

	return c, nil
}

func (s *Server) getCluster(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	if q.Id != nil {
		q.Server = ""

// Update updates a cluster
func (s *Server) Update(ctx context.Context, q *cluster.ClusterUpdateRequest) (*appv1.Cluster, error) {
	c, err := s.getClusterAndVerifyAccess(ctx, &cluster.ClusterQuery{
		Server: q.Cluster.Server,
		Name:   q.Cluster.Name,
		Id:     q.Id,
	}, rbacpolicy.ActionUpdate)

	if err != nil {
		return nil, err
	}

	if len(q.UpdatedFields) == 0 || sets.NewString(q.UpdatedFields...).Has("project") {
		// verify that user can do update inside project where cluster will be located
		if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(q.Cluster.Project, c.Server)) {
	if q.Name != "" {
		servers, err := s.db.GetClusterServersByName(ctx, q.Name)
		if err != nil {
			log.WithField("cluster", q.Name).Warnf("failed to get cluster servers by name: %v", err)
			return nil, common.PermissionDeniedAPIError
		}
		for _, server := range servers {
			if err := enforceAndDelete(s, ctx, server, c.Project); err != nil {

func enforceAndDelete(s *Server, ctx context.Context, server, project string) error {
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionDelete, CreateClusterRBACObject(project, server)); err != nil {
		log.WithField("cluster", server).Warnf("encountered permissions issue while processing request: %v", err)
		return common.PermissionDeniedAPIError
	}
	if err := s.db.DeleteCluster(ctx, server); err != nil {
		return err
	if q.Name != "" {
		servers, err = s.db.GetClusterServersByName(ctx, q.Name)
		if err != nil {
			log.WithField("cluster", q.Name).Warnf("failed to get cluster servers by name: %v", err)
			return nil, common.PermissionDeniedAPIError
		}
		for _, server := range servers {
			if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(clust.Project, server)); err != nil {
				log.WithField("cluster", server).Warnf("encountered permissions issue while processing request: %v", err)
				return nil, common.PermissionDeniedAPIError
			}
		}
	} else {
		if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(clust.Project, q.Server)); err != nil {
			log.WithField("cluster", q.Server).Warnf("encountered permissions issue while processing request: %v", err)
			return nil, common.PermissionDeniedAPIError
		}
		servers = append(servers, q.Server)
	}

// InvalidateCache invalidates cluster cache
func (s *Server) InvalidateCache(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	cls, err := s.getClusterAndVerifyAccess(ctx, q, rbacpolicy.ActionUpdate)
	if err != nil {
		return nil, err
	}
	now := v1.Now()
	cls.RefreshRequestedAt = &now
	cls, err = s.db.UpdateCluster(ctx, cls)
	"context"
	"encoding/json"
	"fmt"
	"github.com/argoproj/argo-cd/v2/server/rbacpolicy"
	"github.com/argoproj/argo-cd/v2/util/assets"
	"github.com/golang-jwt/jwt/v4"
	"reflect"
	"testing"
	"time"
	return enf
}

func newEnforcer() *rbac.Enforcer {
	enforcer := rbac.NewEnforcer(fake.NewSimpleClientset(test.NewFakeConfigMap()), test.FakeArgoCDNamespace, common.ArgoCDRBACConfigMapName, nil)
	_ = enforcer.SetBuiltinPolicy(assets.BuiltinPolicyCSV)
	enforcer.SetDefaultRole("role:test")
	enforcer.SetClaimsEnforcerFunc(func(claims jwt.Claims, rvals ...interface{}) bool {
		return true
	})
	return enforcer
}

func TestUpdateCluster_RejectInvalidParams(t *testing.T) {
	testCases := []struct {
		name    string
		})
	}
}

func TestGetClusterAndVerifyAccess(t *testing.T) {
	t.Run("GetClusterAndVerifyAccess - No Cluster", func(t *testing.T) {
		db := &dbmocks.ArgoDB{}

		mockCluster := v1alpha1.Cluster{
			Name:       "test/ing",
			Server:     "https://127.0.0.1",
			Namespaces: []string{"default", "kube-system"},
		}
		mockClusterList := v1alpha1.ClusterList{
			ListMeta: v1.ListMeta{},
			Items: []v1alpha1.Cluster{
				mockCluster,
			},
		}

		db.On("ListClusters", mock.Anything).Return(&mockClusterList, nil)

		server := NewServer(db, newNoopEnforcer(), newServerInMemoryCache(), &kubetest.MockKubectlCmd{})
		cluster, err := server.getClusterAndVerifyAccess(context.Background(), &clusterapi.ClusterQuery{
			Name: "test/not-exists",
		}, rbacpolicy.ActionGet)

		assert.Nil(t, cluster)
		assert.ErrorIs(t, err, common.PermissionDeniedAPIError)
	})

	t.Run("GetClusterAndVerifyAccess - Permissions Denied", func(t *testing.T) {
		db := &dbmocks.ArgoDB{}

		mockCluster := v1alpha1.Cluster{
			Name:       "test/ing",
			Server:     "https://127.0.0.1",
			Namespaces: []string{"default", "kube-system"},
		}
		mockClusterList := v1alpha1.ClusterList{
			ListMeta: v1.ListMeta{},
			Items: []v1alpha1.Cluster{
				mockCluster,
			},
		}

		db.On("ListClusters", mock.Anything).Return(&mockClusterList, nil)

		server := NewServer(db, newEnforcer(), newServerInMemoryCache(), &kubetest.MockKubectlCmd{})
		cluster, err := server.getClusterAndVerifyAccess(context.Background(), &clusterapi.ClusterQuery{
			Name: "test/ing",
		}, rbacpolicy.ActionGet)

		assert.Nil(t, cluster)
		assert.ErrorIs(t, err, common.PermissionDeniedAPIError)
	})
}

func TestNoClusterEnumeration(t *testing.T) {
	db := &dbmocks.ArgoDB{}

	mockCluster := v1alpha1.Cluster{
		Name:       "test/ing",
		Server:     "https://127.0.0.1",
		Namespaces: []string{"default", "kube-system"},
	}
	mockClusterList := v1alpha1.ClusterList{
		ListMeta: v1.ListMeta{},
		Items: []v1alpha1.Cluster{
			mockCluster,
		},
	}

	db.On("ListClusters", mock.Anything).Return(&mockClusterList, nil)
	db.On("GetCluster", mock.Anything, mock.Anything).Return(&mockCluster, nil)

	server := NewServer(db, newEnforcer(), newServerInMemoryCache(), &kubetest.MockKubectlCmd{})

	t.Run("Get", func(t *testing.T) {
		_, err := server.Get(context.Background(), &clusterapi.ClusterQuery{
			Name: "cluster-not-exists",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")

		_, err = server.Get(context.Background(), &clusterapi.ClusterQuery{
			Name: "test/ing",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")
	})

	t.Run("Update", func(t *testing.T) {
		_, err := server.Update(context.Background(), &clusterapi.ClusterUpdateRequest{
			Cluster: &v1alpha1.Cluster{
				Name: "cluster-not-exists",
			},
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")

		_, err = server.Update(context.Background(), &clusterapi.ClusterUpdateRequest{
			Cluster: &v1alpha1.Cluster{
				Name: "test/ing",
			},
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")
	})

	t.Run("Delete", func(t *testing.T) {
		_, err := server.Delete(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.2",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")

		_, err = server.Delete(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.1",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")
	})

	t.Run("RotateAuth", func(t *testing.T) {
		_, err := server.RotateAuth(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.2",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")

		_, err = server.RotateAuth(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.1",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")
	})

	t.Run("InvalidateCache", func(t *testing.T) {
		_, err := server.InvalidateCache(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.2",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")

		_, err = server.InvalidateCache(context.Background(), &clusterapi.ClusterQuery{
			Server: "https://127.0.0.1",
		})
		assert.Error(t, err)
		assert.Equal(t, common.PermissionDeniedAPIError.Error(), err.Error(), "error message must be _only_ the permission error, to avoid leaking information about cluster existence")
	})
}
