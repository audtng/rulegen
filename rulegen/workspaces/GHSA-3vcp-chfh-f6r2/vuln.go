package main


import (
	"net"

	"github.com/emicklei/go-restful/v3"

	"github.com/kumahq/kuma/v2/pkg/core"
	rest_errors "github.com/kumahq/kuma/v2/pkg/core/rest/errors"
	"github.com/kumahq/kuma/v2/pkg/core/user"
)

var log = core.Log.WithName("api-server").WithName("authn")

func LocalhostAuthenticator(request *restful.Request, response *restful.Response, chain *restful.FilterChain) {
	host, _, err := net.SplitHostPort(request.Request.RemoteAddr)
	if err != nil {
		rest_errors.HandleError(request.Request.Context(), response, err, "Could not parse Remote Address from the Request")
		return
	}
	if host == "127.0.0.1" || host == "::1" {
		log.V(1).Info("authenticated as admin because requests originates from the same machine")
		request.Request = request.Request.WithContext(user.Ctx(request.Request.Context(), user.Admin.Authenticated()))
	}
	chain.ProcessFilter(request, response)
}
	}
	container.Filter(rt.APIServerAuthenticator())

	cors := restful.CrossOriginResourceSharing{
		ExposeHeaders:  []string{restful.HEADER_AccessControlAllowOrigin},
		AllowedDomains: serverConfig.CorsAllowedDomains,
		Container:      container,
	}
	container.Filter(cors.Filter)

	// We create a WebService and set up resources endpoints and index endpoint instead of creating WebService
	// for every resource like /meshes/{mesh}/traffic-permissions, /meshes/{mesh}/traffic-log etc.
func DefaultApiServerConfig() *ApiServerConfig {
	return &ApiServerConfig{
		ReadOnly:           false,
		CorsAllowedDomains: []string{".*"},
		BasePath:           "/",
		HTTP: ApiServerHTTPConfig{
			Enabled:   true,
