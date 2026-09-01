package main

	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/api/norman"
	"github.com/rancher/rancher/pkg/clusterrouter"
	"github.com/rancher/rancher/pkg/features"
	"github.com/rancher/rancher/pkg/types/config"
	"github.com/rancher/rancher/pkg/wrangler"
	steveauth "github.com/rancher/steve/pkg/auth"
	"github.com/sirupsen/logrus"
	"k8s.io/apiserver/pkg/endpoints/request"
)


	root := mux.NewRouter()
	root.UseEncodedPath()
	root.PathPrefix("/v3-public").Handler(publicAPI)
	root.PathPrefix("/v1-saml").Handler(saml)
	root.NotFoundHandler = privateAPI

	return func(next http.Handler) http.Handler {
		next.ServeHTTP(rw, req)
	})
}
	fleetconst "github.com/rancher/rancher/pkg/fleet"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
)

const (
	SQLCacheGCKeepCount = NewSetting("sql-cache-gc-keep-count", "1000")

	SCCOperatorImage = NewSetting("scc-operator-image", buildconfig.DefaultSccOperatorImage)
)

// FullShellImage returns the full private registry name of the rancher shell image.
	return i
}

// SetProvider will set the given provider as the global provider for all settings.
func SetProvider(p Provider) error {
	if err := p.SetAll(settings); err != nil {

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/metrics"
	"github.com/rancher/rancher/pkg/multiclustermanager/whitelist"
	"github.com/rancher/rancher/pkg/rbac"
	"github.com/rancher/rancher/pkg/tunnelserver/mcmauthorizer"
	"github.com/rancher/rancher/pkg/types/config"
	"github.com/rancher/rancher/pkg/version"
	"github.com/rancher/steve/pkg/auth"
)

func router(ctx context.Context, localClusterEnabled bool, tunnelAuthorizer *mcmauthorizer.Authorizer, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager) (func(http.Handler) http.Handler, error) {
	unauthed := mux.NewRouter()
	unauthed.UseEncodedPath()

	unauthed.Path("/").MatcherFunc(parse.MatchNotBrowser).Handler(managementAPI)
	unauthed.Handle("/v3/connect", connectHandler)
	unauthed.Handle("/v3/connect/register", connectHandler)

	return func(next http.Handler) http.Handler {
		metricsAuthed.NotFoundHandler = next
		return unauthed
	}, nil
}

