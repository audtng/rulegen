package main


// Get returns a cluster from a query
func (s *Server) Get(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	c, err := s.getClusterWith403IfNotExist(ctx, q)
	if err != nil {
		return nil, err
	}

	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionGet, CreateClusterRBACObject(c.Project, q.Server)); err != nil {
		return nil, err
	}

	return s.toAPIResponse(c), nil
}

	return repo, nil
}

func (s *Server) getCluster(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	if q.Id != nil {
		q.Server = ""

// Update updates a cluster
func (s *Server) Update(ctx context.Context, q *cluster.ClusterUpdateRequest) (*appv1.Cluster, error) {
	c, err := s.getClusterWith403IfNotExist(ctx, &cluster.ClusterQuery{
		Server: q.Cluster.Server,
		Name:   q.Cluster.Name,
		Id:     q.Id,
	})
	if err != nil {
		return nil, err
	}

	// verify that user can do update inside project where cluster is located
	if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(c.Project, c.Server)) {
		return nil, common.PermissionDeniedAPIError
	}

	if len(q.UpdatedFields) == 0 || sets.NewString(q.UpdatedFields...).Has("project") {
		// verify that user can do update inside project where cluster will be located
		if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(q.Cluster.Project, c.Server)) {
	if q.Name != "" {
		servers, err := s.db.GetClusterServersByName(ctx, q.Name)
		if err != nil {
			return nil, err
		}
		for _, server := range servers {
			if err := enforceAndDelete(s, ctx, server, c.Project); err != nil {

func enforceAndDelete(s *Server, ctx context.Context, server, project string) error {
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionDelete, CreateClusterRBACObject(project, server)); err != nil {
		return err
	}
	if err := s.db.DeleteCluster(ctx, server); err != nil {
		return err
	if q.Name != "" {
		servers, err = s.db.GetClusterServersByName(ctx, q.Name)
		if err != nil {
			return nil, status.Errorf(codes.NotFound, "failed to get cluster servers by name: %v", err)
		}
		for _, server := range servers {
			if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(clust.Project, server)); err != nil {
				return nil, status.Errorf(codes.PermissionDenied, "encountered permissions issue while processing request: %v", err)
			}
		}
	} else {
		if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(clust.Project, q.Server)); err != nil {
			return nil, status.Errorf(codes.PermissionDenied, "encountered permissions issue while processing request: %v", err)
		}
		servers = append(servers, q.Server)
	}

// InvalidateCache invalidates cluster cache
func (s *Server) InvalidateCache(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	cls, err := s.getClusterWith403IfNotExist(ctx, q)
	if err != nil {
		return nil, err
	}
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, CreateClusterRBACObject(cls.Project, q.Server)); err != nil {
		return nil, err
	}
	now := v1.Now()
	cls.RefreshRequestedAt = &now
	cls, err = s.db.UpdateCluster(ctx, cls)
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
	"time"
	return enf
}

func TestUpdateCluster_RejectInvalidParams(t *testing.T) {
	testCases := []struct {
		name    string
		})
	}
}
