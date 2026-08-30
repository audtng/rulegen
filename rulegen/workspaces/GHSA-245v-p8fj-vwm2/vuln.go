package main


import (
	"io"
	"strings"

	"github.com/juju/charm/v8"
	charmresource "github.com/juju/charm/v8/resource"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"gopkg.in/macaroon.v2"

	"github.com/juju/juju/api/base"
	return args, nil
}

// Upload sends the provided resource blob up to Juju.
func (c Client) Upload(application, name, filename string, reader io.ReadSeeker) error {
	uReq, err := NewUploadRequest(application, name, filename, reader)

	var response params.UploadResult // ignored
	if err := c.httpClient.Do(c.facade.RawAPICaller().Context(), req, &response); err != nil {
		return errors.Trace(err)
	}

	return nil

		var response params.UploadResult // ignored
		if err := c.httpClient.Do(c.facade.RawAPICaller().Context(), req, &response); err != nil {
			return "", errors.Trace(err)
		}
	}

	"github.com/kr/pretty"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"

	"github.com/juju/juju/api/client/resources"
	apicharm "github.com/juju/juju/api/common/charm"
	c.Check(err, gc.ErrorMatches, `.*invalid application.*`)
}

func (s *UploadSuite) TestUploadFailed(c *gc.C) {
	defer s.setUpMocks(c).Finish()

	ctx := context.TODO()
	req.Header.Set("Content-Disposition", "form-data; filename=foo.zip")
	req.ContentLength = int64(len(data))

	s.httpClient.EXPECT().Do(ctx, reqMatcher{c, req}, gomock.Any()).Return(errors.New("boom"))

	err = s.client.Upload("a-application", "spam", "foo.zip", strings.NewReader(data))
	c.Assert(err, gc.ErrorMatches, "boom")
}

func (s *UploadSuite) TestAddPendingResources(c *gc.C) {
	c.Assert(err, gc.ErrorMatches, `.*invalid application.*`)
}

func (s *UploadSuite) TestUploadPendingResourceFailed(c *gc.C) {
	defer s.setUpMocks(c).Finish()

	res, apiResult := newResourceResult(c, "spam")
	req.URL.RawQuery = "pendingid=" + expected
	req.Header.Set("Content-Disposition", "form-data; filename=file.zip")

	s.httpClient.EXPECT().Do(ctx, reqMatcher{c, req}, gomock.Any()).Return(errors.New("boom"))

	_, err = s.client.UploadPendingResource("a-application", res[0].Resource, "file.zip", strings.NewReader(data))
	c.Assert(err, gc.ErrorMatches, "boom")
}
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
}, {
	about:       "non-JSON error response",
	handler:     http.NotFound,
	expectError: `Get http://.*/: unexpected content type text/plain; want application/json; content: 404 page not found`,
}, {
	about: "bad error response",
	handler: func(w http.ResponseWriter, req *http.Request) {
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
type ResourcesHandlerSuite struct {
	testing.IsolationSuite

	stateAuthErr error
	backend      *fakeBackend
	username     string
	req          *http.Request
	recorder     *httptest.ResponseRecorder
	handler      *apiserver.ResourcesHandler
}

var _ = gc.Suite(&ResourcesHandlerSuite{})
	c.Assert(err, jc.ErrorIsNil)
	s.req = req
	s.recorder = httptest.NewRecorder()
	s.handler = &apiserver.ResourcesHandler{
		StateAuthFunc:     s.authState,
		ChangeAllowedFunc: func(*http.Request) error { return nil },
	}
}

func (s *ResourcesHandlerSuite) authState(req *http.Request, tagKinds ...string) (
	apiserver.ResourcesBackend, state.PoolHelper, names.Tag, error,
) {
	if s.stateAuthErr != nil {
func (s *ResourcesHandlerSuite) TestExpectedAuthTags(c *gc.C) {
	expectedTags := set.NewStrings(names.UserTagKind, names.MachineTagKind, names.ControllerAgentTagKind, names.ApplicationTagKind)

	s.handler.StateAuthFunc = func(req *http.Request, tagKinds ...string) (apiserver.ResourcesBackend, state.PoolHelper, names.Tag, error) {
		gotTags := set.NewStrings(tagKinds...)
		if gotTags.Difference(expectedTags).Size() != 0 || expectedTags.Difference(gotTags).Size() != 0 {
			c.Fatalf("unexpected tag kinds %v", tagKinds)
			return nil, nil, nil, errors.NotValidf("tag kinds %v", tagKinds)
		}
		ph := apiservertesting.StubPoolHelper{StubRelease: func() bool { return false }}
		tag := names.NewUserTag(s.username)
		return s.backend, ph, tag, nil
	}
	s.req.Method = "GET"
	s.handler.ServeHTTP(s.recorder, s.req)
	s.checkResp(c, http.StatusOK, "application/octet-stream", resourceBody)
}

	failure, expected := apiFailure("<failure>", "")
	s.stateAuthErr = failure

	s.handler.ServeHTTP(s.recorder, s.req)

	s.checkResp(c, http.StatusInternalServerError, "application/json", expected)
}

func (s *ResourcesHandlerSuite) TestUnsupportedMethod(c *gc.C) {
	s.req.Method = "POST"

	s.handler.ServeHTTP(s.recorder, s.req)

	_, expected := apiFailure(`unsupported method: "POST"`, params.CodeMethodNotAllowed)
	s.checkResp(c, http.StatusMethodNotAllowed, "application/json", expected)
}

func (s *ResourcesHandlerSuite) TestGetSuccess(c *gc.C) {
	s.req.Method = "GET"
	s.handler.ServeHTTP(s.recorder, s.req)
	s.checkResp(c, http.StatusOK, "application/octet-stream", resourceBody)
}

	s.backend.ReturnSetResource = res

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.handler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),
	s.backend.ReturnSetResource = res

	expectedError := apiservererrors.OperationBlockedError("test block")
	s.handler.ChangeAllowedFunc = func(*http.Request) error {
		return expectedError
	}

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.handler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.ErrorResult{apiservererrors.ServerError(expectedError)})
	s.checkResp(c, http.StatusBadRequest, "application/json", string(expected))
	s.backend.ReturnSetResource = res

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.handler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	req.Header.Set("Content-Disposition", "form-data; filename=different.ext")
	s.handler.ServeHTTP(s.recorder, req)

	_, expected := apiFailure(`incorrect extension on resource upload "different.ext", expected ".tgz"`,
		"")

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	req.URL.RawQuery += "&pendingid=some-unique-id"
	s.handler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),
	s.backend.SetResourceErr = failure

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	s.handler.ServeHTTP(s.recorder, req)
	s.checkResp(c, http.StatusInternalServerError, "application/json", expected)
}

