package main

package cluster

import (
	"context"
	"net/url"
	"time"

	"github.com/argoproj/gitops-engine/pkg/utils/kube"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/client-go/kubernetes"

	"github.com/argoproj/argo-cd/v2/common"
	"github.com/argoproj/argo-cd/v2/pkg/apiclient/cluster"
	appv1 "github.com/argoproj/argo-cd/v2/pkg/apis/application/v1alpha1"
	servercache "github.com/argoproj/argo-cd/v2/server/cache"
func (s *Server) getClusterWith403IfNotExist(ctx context.Context, q *cluster.ClusterQuery) (*appv1.Cluster, error) {
	repo, err := s.getCluster(ctx, q)
	if err != nil || repo == nil {
		return nil, common.PermissionDeniedAPIError
	}
	return repo, nil
}
	}

	// verify that user can do update inside project where cluster is located
	if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, createRBACObject(c.Project, c.Server)) {
		return nil, common.PermissionDeniedAPIError
	}

	if len(q.UpdatedFields) == 0 || sets.NewString(q.UpdatedFields...).Has("project") {
		// verify that user can do update inside project where cluster will be located
		if !s.enf.Enforce(ctx.Value("claims"), rbacpolicy.ResourceClusters, rbacpolicy.ActionUpdate, createRBACObject(q.Cluster.Project, c.Server)) {
			return nil, common.PermissionDeniedAPIError
		}
	}

	"time"

	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// Default service addresses and URLS of Argo CD internal services
const TokenVerificationError = "failed to verify the token"

var TokenVerificationErr = errors.New(TokenVerificationError)

var PermissionDeniedAPIError = status.Error(codes.PermissionDenied, "permission denied")
