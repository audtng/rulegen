package main

	UpdatePendingResource(applicationID, pendingID, userID string, res charmresource.Resource, r io.Reader) (resources.Resource, error)
}

// ResourcesHandler is the HTTP handler for client downloads and
// uploads of resources.
type ResourcesHandler struct {
	StateAuthFunc     func(*http.Request, ...string) (ResourcesBackend, state.PoolHelper, names.Tag, error)
	ChangeAllowedFunc func(*http.Request) error
}

// ServeHTTP implements http.Handler.
func (h *ResourcesHandler) ServeHTTP(resp http.ResponseWriter, req *http.Request) {
	backend, poolhelper, tag, err := h.StateAuthFunc(req, names.UserTagKind, names.MachineTagKind, names.ControllerAgentTagKind, names.ApplicationTagKind)
	if err != nil {
		if err := sendError(resp, err); err != nil {
			logger.Errorf("%v", err)
		if _, err := io.Copy(resp, reader); err != nil {
			logger.Errorf("resource download failed: %v", err)
		}
	case "PUT":
		if err := h.ChangeAllowedFunc(req); err != nil {
			if err := sendError(resp, err); err != nil {
	}
}

func (h *ResourcesHandler) download(backend ResourcesBackend, req *http.Request) (io.ReadCloser, int64, error) {
	defer req.Body.Close()

	query := req.URL.Query()
	application := query.Get(":application")
	name := query.Get(":resource")

	resource, reader, err := backend.OpenResource(application, name)
	return reader, resource.Size, errors.Trace(err)
}

func (h *ResourcesHandler) upload(backend ResourcesBackend, req *http.Request, username string) (*params.UploadResult, error) {
	defer req.Body.Close()

	uploaded, err := h.readResource(backend, req)
}

// readResource extracts the relevant info from the request.
func (h *ResourcesHandler) readResource(backend ResourcesBackend, req *http.Request) (*uploadedResource, error) {
	uReq, err := extractUploadRequest(req)
	if err != nil {
		return nil, errors.Trace(err)
		},
	}
	modelToolsDownloadHandler := newToolsDownloadHandler(httpCtxt)
	resourcesHandler := &ResourcesHandler{
		StateAuthFunc: func(req *http.Request, tagKinds ...string) (ResourcesBackend, state.PoolHelper, names.Tag, error) {
			st, entity, err := httpCtxt.stateForRequestAuthenticatedTag(req, tagKinds...)
			if err != nil {
				return nil, nil, nil, errors.Trace(err)
			}
			return nil
		},
	}
	unitResourcesHandler := &UnitResourcesHandler{
		NewOpener: func(req *http.Request, tagKinds ...string) (resources.Opener, state.PoolHelper, error) {
			st, _, err := httpCtxt.stateForRequestAuthenticatedTag(req, tagKinds...)
		unauthenticated: true,
	}, {
		pattern: modelRoutePrefix + "/applications/:application/resources/:resource",
		handler: resourcesHandler,
	}, {
		pattern: modelRoutePrefix + "/units/:unit/resources/:resource",
		handler: unitResourcesHandler,
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/go-macaroon-bakery/macaroon-bakery/v3/bakery"
	"github.com/go-macaroon-bakery/macaroon-bakery/v3/httpbakery"
func unmarshalHTTPErrorResponse(resp *http.Response) error {
	var body json.RawMessage
	if err := httprequest.UnmarshalJSONResponse(resp, &body); err != nil {
		return errors.Trace(err)
	}
	// genericErrorResponse defines a struct that is compatible with all the
