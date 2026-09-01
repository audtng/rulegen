package main

	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/api/norman"
	"github.com/rancher/rancher/pkg/clusterrouter"
	"github.com/rancher/rancher/pkg/features"
	"github.com/rancher/rancher/pkg/types/config"
	"github.com/rancher/rancher/pkg/utils"
	"github.com/rancher/rancher/pkg/wrangler"
	steveauth "github.com/rancher/steve/pkg/auth"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apiserver/pkg/endpoints/request"
)


	root := mux.NewRouter()
	root.UseEncodedPath()

	apiLimit, err := quantityAsInt64(getEnvWithDefault("CATTLE_AUTH_API_BODY_LIMIT", "1Mi"), 1024*1024)
	if err != nil {
		return nil, err
	}
	logrus.Infof("Configuring auth server API body limit to %v bytes", apiLimit)

	limitingHandler := utils.APIBodyLimitingHandler(apiLimit)
	root.PathPrefix("/v3-public").Handler(limitingHandler(publicAPI))
	root.PathPrefix("/v1-saml").Handler(limitingHandler(saml))
	root.NotFoundHandler = privateAPI

	return func(next http.Handler) http.Handler {
		next.ServeHTTP(rw, req)
	})
}

func quantityAsInt64(s string, d int64) (int64, error) {
	i, err := resource.ParseQuantity(s)
	if err != nil {
		return 0, fmt.Errorf("parsing setting: %w", err)
	}

	q, ok := i.AsInt64()
	if ok {
		return q, nil
	}

	return d, nil
}

func getEnvWithDefault(key, defaultValue string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return defaultValue
}
	fleetconst "github.com/rancher/rancher/pkg/fleet"
	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
)

const (
	SQLCacheGCKeepCount = NewSetting("sql-cache-gc-keep-count", "1000")

	SCCOperatorImage = NewSetting("scc-operator-image", buildconfig.DefaultSccOperatorImage)

	// This is the limit for request bodies sent to /v3-public/* endpoints in
	// bytes.
	// The default = 1MiB
	APIBodyLimit = NewSetting("public-api-body-limit", "1Mi")
)

// FullShellImage returns the full private registry name of the rancher shell image.
	return i
}

// GetQuantityAsInt64 will return the currently stored value of the setting as an int64
// parsed from a Kubernetes Quantity format string.
//
// See https://pkg.go.dev/k8s.io/apimachinery/pkg/api/resource#ParseQuantity for
// format details.
//
// If the quantity cannot be expressed as an int64 d will be returned.
func (s Setting) GetQuantityAsInt64(d int64) (int64, error) {
	v := s.Get()
	i, err := resource.ParseQuantity(v)
	if err != nil {
		return 0, fmt.Errorf("parsing setting: %w", err)
	}

	q, ok := i.AsInt64()
	if ok {
		return q, nil
	}

	return d, nil
}

// SetProvider will set the given provider as the global provider for all settings.
func SetProvider(p Provider) error {
	if err := p.SetAll(settings); err != nil {

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rancher/rancher/pkg/metrics"
	"github.com/rancher/rancher/pkg/multiclustermanager/whitelist"
	"github.com/rancher/rancher/pkg/rbac"
	"github.com/rancher/rancher/pkg/settings"
	"github.com/rancher/rancher/pkg/tunnelserver/mcmauthorizer"
	"github.com/rancher/rancher/pkg/types/config"
	"github.com/rancher/rancher/pkg/utils"
	"github.com/rancher/rancher/pkg/version"
	"github.com/rancher/steve/pkg/auth"
	"github.com/sirupsen/logrus"
)

func router(ctx context.Context, localClusterEnabled bool, tunnelAuthorizer *mcmauthorizer.Authorizer, scaledContext *config.ScaledContext, clusterManager *clustermanager.Manager) (func(http.Handler) http.Handler, error) {
	unauthed := mux.NewRouter()
	unauthed.UseEncodedPath()

	publicLimit, err := settings.APIBodyLimit.GetQuantityAsInt64(1024 * 1024)
	if err != nil {
		return nil, fmt.Errorf("parsing the public API body limit: %w", err)
	}
	logrus.Infof("Configuring public API body limit to %v bytes", publicLimit)
	limitingHandler := utils.APIBodyLimitingHandler(publicLimit)

	unauthed.Path("/").MatcherFunc(parse.MatchNotBrowser).Handler(managementAPI)
	unauthed.Handle("/v3/connect", connectHandler)
	unauthed.Handle("/v3/connect/register", connectHandler)

	return func(next http.Handler) http.Handler {
		metricsAuthed.NotFoundHandler = next
		return limitingHandler(unauthed)
	}, nil
}

