package main

	UpdatePendingResource(applicationID, pendingID, userID string, res charmresource.Resource, r io.Reader) (resources.Resource, error)
}

// ResourcesUploadHandler is the HTTP handler for client
// uploads of resources.
type ResourcesUploadHandler struct {
	ChangeAllowedFunc func(*http.Request) error
	StateFunc         func(*http.Request) (ResourcesBackend, state.PoolHelper, names.Tag, error)
}

// ResourcesDownloadHandler is the HTTP handler for client
// downloads of resources.
type ResourcesDownloadHandler struct {
	StateAuthFunc func(*http.Request, ...string) (ResourcesBackend, state.PoolHelper, error)
}

// ServeHTTP implements http.Handler.
func (h *ResourcesDownloadHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	backend, poolhelper, err := h.StateAuthFunc(req, names.UserTagKind, names.MachineTagKind, names.ControllerAgentTagKind, names.ApplicationTagKind)
	if err != nil {
		if err := sendError(resp, err); err != nil {
			logger.Errorf("%v", err)
		if _, err := io.Copy(resp, reader); err != nil {
			logger.Errorf("resource download failed: %v", err)
		}
	default:
		if err := sendError(resp, errors.MethodNotAllowedf("unsupported method: %q", req.Method)); err != nil {
			logger.Errorf("%v", err)
		}
	}
}

func (h *ResourcesDownloadHandler) download(backend ResourcesBackend, req *http.Request) (io.ReadCloser, int64, error) {
	defer req.Body.Close()

	query := req.URL.Query()
	application := query.Get(":application")
	name := query.Get(":resource")

	resource, reader, err := backend.OpenResource(application, name)
	return reader, resource.Size, errors.Trace(err)
}

// ServeHTTP implements http.Handler.
func (h *ResourcesUploadHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	backend, closer, tag, err := h.StateFunc(req)
	if err != nil {
		if err := sendError(resp, err); err != nil {
			logger.Errorf("%v", err)
		}
		return
	}
	defer closer.Release()

	switch req.Method {
	case "PUT":
		if err := h.ChangeAllowedFunc(req); err != nil {
			if err := sendError(resp, err); err != nil {
	}
}

func (h *ResourcesUploadHandler) upload(backend ResourcesBackend, req *http.Request, username string) (*params.UploadResult, error) {
	defer req.Body.Close()

	uploaded, err := h.readResource(backend, req)
}

// readResource extracts the relevant info from the request.
func (h *ResourcesUploadHandler) readResource(backend ResourcesBackend, req *http.Request) (*uploadedResource, error) {
	uReq, err := extractUploadRequest(req)
	if err != nil {
		return nil, errors.Trace(err)
		},
	}
	modelToolsDownloadHandler := newToolsDownloadHandler(httpCtxt)
	var resourcesUploadAuthorizer httpcontext.CompositeAuthorizer = []httpcontext.Authorizer{
		controllerAdminAuthorizer{
			st: systemState,
		},
		modelPermissionAuthorizer{
			userAccess: systemState.UserPermission,
			perm:       permission.WriteAccess,
		},
	}
	resourceUploadHandler := &ResourcesUploadHandler{
		StateFunc: func(req *http.Request) (ResourcesBackend, state.PoolHelper, names.Tag, error) {
			st, entity, err := httpCtxt.stateForRequestAuthenticated(req)
			if err != nil {
				return nil, nil, nil, errors.Trace(err)
			}
			return nil
		},
	}
	resourceDownloadHandler := &ResourcesDownloadHandler{
		StateAuthFunc: func(req *http.Request, tagKinds ...string) (ResourcesBackend, state.PoolHelper, error) {
			st, _, err := httpCtxt.stateForRequestAuthenticatedTag(req, tagKinds...)
			if err != nil {
				return nil, nil, errors.Trace(err)
			}
			rst := st.Resources()
			return rst, st, nil
		},
	}
	unitResourcesHandler := &UnitResourcesHandler{
		NewOpener: func(req *http.Request, tagKinds ...string) (resources.Opener, state.PoolHelper, error) {
			st, _, err := httpCtxt.stateForRequestAuthenticatedTag(req, tagKinds...)
		unauthenticated: true,
	}, {
		pattern: modelRoutePrefix + "/applications/:application/resources/:resource",
		methods: []string{"GET"},
		handler: resourceDownloadHandler,
	}, {
		pattern:    modelRoutePrefix + "/applications/:application/resources/:resource",
		methods:    []string{"PUT"},
		handler:    resourceUploadHandler,
		authorizer: resourcesUploadAuthorizer,
	}, {
		pattern: modelRoutePrefix + "/units/:unit/resources/:resource",
		handler: unitResourcesHandler,
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/go-macaroon-bakery/macaroon-bakery/v3/bakery"
	"github.com/go-macaroon-bakery/macaroon-bakery/v3/httpbakery"
func unmarshalHTTPErrorResponse(resp *http.Response) error {
	var body json.RawMessage
	if err := httprequest.UnmarshalJSONResponse(resp, &body); err != nil {
		// Auth errors can come back as plain text as the http handler
		// authentication function does not write its error response as json.
		// We want to return a suitable error code so the caller can
		// handle that case. We don't want to propagate the message
		// "unexpected content type text/plain; want application/json".
		var decodeErr *httprequest.DecodeResponseError
		if errors.As(err, &decodeErr) {
			var code string
			switch resp.StatusCode {
			case http.StatusForbidden:
				code = params.CodeForbidden
			case http.StatusUnauthorized:
				code = params.CodeUnauthorized
			}
			var errMsg string
			msg, readErr := io.ReadAll(decodeErr.Response.Body)
			if readErr == nil {
				errMsg = strings.Trim(string(msg), "\n")
			} else {
				// Should never happen.
				errMsg = err.Error()
			}
			return params.Error{
				Message: errMsg,
				Code:    code,
			}
		}
		return errors.Trace(err)
	}
	// genericErrorResponse defines a struct that is compatible with all the
