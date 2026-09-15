package main

package cluster

import (
	"net/url"
	"time"

	"context"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	servercache "github.com/argoproj/argo-cd/v2/server/cache"
func (s *Server) getClusterWith403IfNotExist(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	repo, err := s.getCluster(ctx, q)
	if err != nil || repo == nil {
		return nil, status.Error(codes.PermissionDenied, "permission denied")
	}
	return repo, nil
}
	}

	// verify that user can do update inside project where cluster is located
	if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, createRBACObject(c.Project, q.Cluster.Server)); err != nil {
		return nil, err
	}

	if len(q.UpdatedFields) == 0 || sets.NewString(q.UpdatedFields...).Has("project") {
		// verify that user can do update inside project where cluster will be located
		if err := s.enf.EnforceErr(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, createRBACObject(q.Cluster.Project, q.Cluster.Server)); err != nil {
			return nil, err
		}
	}

	"time"

	"github.com/sirupsen/logrus"
)

// Default service addresses and URLS of Argo CD internal services
const TokenVerificationError = "failed to verify the token"

var TokenVerificationErr = errors.New(TokenVerificationError)
