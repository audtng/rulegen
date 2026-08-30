package main


import (
	"io"
	"slices"
	"strings"

	"github.com/juju/charm/v8"
	charmresource "github.com/juju/charm/v8/resource"
	"github.com/juju/errors"
	"github.com/juju/names/v4"
	"gopkg.in/errgo.v1"
	"gopkg.in/macaroon.v2"

	"github.com/juju/juju/api/base"
	return args, nil
}

func translateUploadError(err error) error {
	permissionCodes := []string{params.CodeForbidden, params.CodeUnauthorized}
	var uploadErr *errgo.Err
	if errors.As(err, &uploadErr) {
		if slices.Contains(permissionCodes, params.ErrCode(uploadErr.Cause())) ||
			slices.Contains(permissionCodes, params.ErrCode(uploadErr.Underlying())) {
			return apiservererrors.ErrPerm
		}
	}
	return err
}

// Upload sends the provided resource blob up to Juju.
func (c Client) Upload(application, name, filename string, reader io.ReadSeeker) error {
	uReq, err := NewUploadRequest(application, name, filename, reader)

	var response params.UploadResult // ignored
	if err := c.httpClient.Do(c.facade.RawAPICaller().Context(), req, &response); err != nil {
		return errors.Trace(translateUploadError(err))
	}

	return nil

		var response params.UploadResult // ignored
		if err := c.httpClient.Do(c.facade.RawAPICaller().Context(), req, &response); err != nil {
			return "", translateUploadError(errors.Trace(err))
		}
	}

	"github.com/kr/pretty"
	"go.uber.org/mock/gomock"
	gc "gopkg.in/check.v1"
	"gopkg.in/errgo.v1"

	"github.com/juju/juju/api/client/resources"
	apicharm "github.com/juju/juju/api/common/charm"
	c.Check(err, gc.ErrorMatches, `.*invalid application.*`)
}

func (s *UploadSuite) assertUploadFailed(c *gc.C, uploadErr error, expectedMsg string) {
	defer s.setUpMocks(c).Finish()

	ctx := context.TODO()
	req.Header.Set("Content-Disposition", "form-data; filename=foo.zip")
	req.ContentLength = int64(len(data))

	s.httpClient.EXPECT().Do(ctx, reqMatcher{c, req}, gomock.Any()).Return(uploadErr)

	err = s.client.Upload("a-application", "spam", "foo.zip", strings.NewReader(data))
	c.Assert(err, gc.ErrorMatches, expectedMsg)
}

func (s *UploadSuite) TestUploadFailed(c *gc.C) {
	s.assertUploadFailed(c, errors.New("upload failed"), "upload failed")
}

func (s *UploadSuite) TestUploadAuthError(c *gc.C) {
	authErr := errgo.Mask(params.Error{
		Code:    params.CodeUnauthorized,
		Message: "user unauthorized",
	})
	s.assertUploadFailed(c, authErr, "permission denied")
}

func (s *UploadSuite) TestAddPendingResources(c *gc.C) {
	c.Assert(err, gc.ErrorMatches, `.*invalid application.*`)
}

func (s *UploadSuite) assertUploadPendingResourceFailed(c *gc.C, uploadErr error, expectedMsg string) {
	defer s.setUpMocks(c).Finish()

	res, apiResult := newResourceResult(c, "spam")
	req.URL.RawQuery = "pendingid=" + expected
	req.Header.Set("Content-Disposition", "form-data; filename=file.zip")

	s.httpClient.EXPECT().Do(ctx, reqMatcher{c, req}, gomock.Any()).Return(uploadErr)

	_, err = s.client.UploadPendingResource("a-application", res[0].Resource, "file.zip", strings.NewReader(data))
	c.Assert(err, gc.ErrorMatches, expectedMsg)
}

func (s *UploadSuite) TestUploadPendingResourceFailed(c *gc.C) {
	s.assertUploadPendingResourceFailed(c, errors.New("upload failed"), "upload failed")
}

func (s *UploadSuite) TestUploadPendingResourceAuthError(c *gc.C) {
	authErr := errgo.Mask(params.Error{
		Code:    params.CodeUnauthorized,
		Message: "user unauthorized",
	})
	s.assertUploadPendingResourceFailed(c, authErr, "permission denied")
}
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
}, {
	about:       "non-JSON error response",
	handler:     http.NotFound,
	expectError: `Get http://.*/: 404 page not found`,
}, {
	about: "non-JSON auth error response",
	handler: func(w http.ResponseWriter, req *http.Request) {
		http.Error(w, "some unauth error", http.StatusUnauthorized)
	},
	expectError:     `Get http://.*/: some unauth error`,
	expectErrorCode: "unauthorized access",
}, {
	about: "bad error response",
	handler: func(w http.ResponseWriter, req *http.Request) {
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
// Copyright 2025 Canonical Ltd.
// Licensed under the AGPLv3, see LICENCE file for details.

package apiserver_test

import (
	"encoding/json"
	"fmt"
	"mime"
	"net/http"
	"net/url"
	"strings"

	jc "github.com/juju/testing/checkers"
	"github.com/juju/utils/v3"
	gc "gopkg.in/check.v1"

	apitesting "github.com/juju/juju/apiserver/testing"
	"github.com/juju/juju/core/permission"
	"github.com/juju/juju/core/resources"
	"github.com/juju/juju/rpc/params"
	"github.com/juju/juju/state"
	"github.com/juju/juju/testing/factory"
)

type resourcesAuthSuite struct {
	apiserverBaseSuite
}

func (s *resourcesAuthSuite) resourcesURL(app, res string) *url.URL {
	u := s.URL(fmt.Sprintf("/model/%s/applications/%s/resources/%s", s.Model.UUID(), app, res), nil)
	return u
}

func (s *resourcesAuthSuite) assertJSONErrorResponse(c *gc.C, resp *http.Response, expCode int, expError string) {
	uploadResponse := s.assertResponse(c, resp, expCode)
	c.Check(uploadResponse.Error, gc.NotNil)
	c.Check(uploadResponse.Error.Message, gc.Matches, expError)
}

func (s *resourcesAuthSuite) assertPlainErrorResponse(c *gc.C, resp *http.Response, expCode int, expError string) {
	body := apitesting.AssertResponse(c, resp, expCode, "text/plain; charset=utf-8")
	c.Assert(string(body), gc.Matches, expError+"\n")
}

func (s *resourcesAuthSuite) assertResponse(c *gc.C, resp *http.Response, expStatus int) params.UploadResult {
	body := apitesting.AssertResponse(c, resp, expStatus, params.ContentTypeJSON)
	var uploadResult params.UploadResult
	err := json.Unmarshal(body, &uploadResult)
	c.Assert(err, jc.ErrorIsNil, gc.Commentf("Body: %s", body))
	return uploadResult
}

var _ = gc.Suite(&resourcesAuthSuite{})

func (s *resourcesAuthSuite) TestResourcesUploadedSecurely(c *gc.C) {
	url := s.resourcesURL("tomcat", "jdk")
	url.Scheme = "http"
	resp := apitesting.SendHTTPRequest(c, apitesting.HTTPRequestParams{
		Method:       "PUT",
		URL:          url.String(),
		ExpectStatus: http.StatusBadRequest,
	})
	defer resp.Body.Close()
}

func (s *resourcesAuthSuite) TestRequiresAuth(c *gc.C) {
	resp := apitesting.SendHTTPRequest(c, apitesting.HTTPRequestParams{Method: "GET", URL: s.resourcesURL("tomcat", "jdk").String()})
	defer resp.Body.Close()
	s.assertPlainErrorResponse(c, resp, http.StatusUnauthorized, "authentication failed: no credentials provided")
}

func (s *resourcesAuthSuite) TestAuthRejectsNonsUser(c *gc.C) {
	// Add a machine and try to login.
	machine, err := s.State.AddMachine("quantal", state.JobHostUnits)
	c.Assert(err, jc.ErrorIsNil)
	err = machine.SetProvisioned("foo", "", "fake_nonce", nil)
	c.Assert(err, jc.ErrorIsNil)
	password, err := utils.RandomPassword()
	c.Assert(err, jc.ErrorIsNil)
	err = machine.SetPassword(password)
	c.Assert(err, jc.ErrorIsNil)

	resp := apitesting.SendHTTPRequest(c, apitesting.HTTPRequestParams{
		Tag:      machine.Tag().String(),
		Password: password,
		Method:   "PUT",
		URL:      s.resourcesURL("tomcat", "jdk").String(),
		Nonce:    "fake_nonce",
	})
	s.assertPlainErrorResponse(
		c, resp, http.StatusForbidden,
		"authorization failed: permission denied",
	)
	resp.Body.Close()

	// Now try a user login.
	content, err := resources.GenerateContent(strings.NewReader("resource"))
	c.Assert(err, jc.ErrorIsNil)
	filename := mime.BEncoding.Encode("utf-8", "foo.txt")
	disp := mime.FormatMediaType(
		"form-data",
		map[string]string{"filename": filename},
	)

	resp = s.sendHTTPRequest(c, apitesting.HTTPRequestParams{
		Method:      "PUT",
		URL:         s.resourcesURL("tomcat", "jdk").String(),
		ContentType: "application/octet-stream",
		ExtraHeaders: map[string]string{
			"Content-Sha384":      content.Fingerprint.String(),
			"Content-Length":      fmt.Sprintf("%d", content.Size),
			"Content-Disposition": disp,
		},
		Body: strings.NewReader("fake_nonce"),
	})
	s.assertJSONErrorResponse(c, resp, http.StatusNotFound, `application "tomcat" not found`)
	resp.Body.Close()
}

func (s *resourcesAuthSuite) TestAuthRejectsUserWithoutPermission(c *gc.C) {
	u := s.Factory.MakeUser(c, &factory.UserParams{
		Name:     "oryx",
		Password: "gardener",
		Access:   permission.ReadAccess,
	})

	resp := apitesting.SendHTTPRequest(c, apitesting.HTTPRequestParams{
		Tag:      u.Tag().String(),
		Password: "gardener",
		Method:   "PUT",
		URL:      s.resourcesURL("tomcat", "jdk").String(),
	})
	defer resp.Body.Close()
	s.assertPlainErrorResponse(
		c, resp, http.StatusForbidden,
		"authorization failed: permission denied",
	)
}
type ResourcesHandlerSuite struct {
	testing.IsolationSuite

	stateAuthErr    error
	backend         *fakeBackend
	username        string
	req             *http.Request
	recorder        *httptest.ResponseRecorder
	uploadHandler   *apiserver.ResourcesUploadHandler
	downloadHandler *apiserver.ResourcesDownloadHandler
}

var _ = gc.Suite(&ResourcesHandlerSuite{})
	c.Assert(err, jc.ErrorIsNil)
	s.req = req
	s.recorder = httptest.NewRecorder()
	s.uploadHandler = &apiserver.ResourcesUploadHandler{
		StateFunc:         s.stateUpload,
		ChangeAllowedFunc: func(*http.Request) error { return nil },
	}
	s.downloadHandler = &apiserver.ResourcesDownloadHandler{
		StateAuthFunc: s.authStateDownload,
	}
}

func (s *ResourcesHandlerSuite) authStateDownload(req *http.Request, tagKinds ...string) (
	apiserver.ResourcesBackend, state.PoolHelper, error,
) {
	if s.stateAuthErr != nil {
		return nil, nil, errors.Trace(s.stateAuthErr)
	}

	ph := apiservertesting.StubPoolHelper{StubRelease: func() bool { return false }}
	return s.backend, ph, nil
}

func (s *ResourcesHandlerSuite) stateUpload(req *http.Request) (
	apiserver.ResourcesBackend, state.PoolHelper, names.Tag, error,
) {
	if s.stateAuthErr != nil {
func (s *ResourcesHandlerSuite) TestExpectedAuthTags(c *gc.C) {
	expectedTags := set.NewStrings(names.UserTagKind, names.MachineTagKind, names.ControllerAgentTagKind, names.ApplicationTagKind)

	s.downloadHandler.StateAuthFunc = func(req *http.Request, tagKinds ...string) (apiserver.ResourcesBackend, state.PoolHelper, error) {
		gotTags := set.NewStrings(tagKinds...)
		if gotTags.Difference(expectedTags).Size() != 0 || expectedTags.Difference(gotTags).Size() != 0 {
			c.Fatalf("unexpected tag kinds %v", tagKinds)
			return nil, nil, errors.NotValidf("tag kinds %v", tagKinds)
		}
		ph := apiservertesting.StubPoolHelper{StubRelease: func() bool { return false }}
		return s.backend, ph, nil
	}
	s.req.Method = "GET"
	s.downloadHandler.ServeHTTP(s.recorder, s.req)
	s.checkResp(c, http.StatusOK, "application/octet-stream", resourceBody)
}

	failure, expected := apiFailure("<failure>", "")
	s.stateAuthErr = failure

	s.uploadHandler.ServeHTTP(s.recorder, s.req)

	s.checkResp(c, http.StatusInternalServerError, "application/json", expected)
}

func (s *ResourcesHandlerSuite) TestDownloadUnsupportedMethod(c *gc.C) {
	s.req.Method = "PUT"

	s.downloadHandler.ServeHTTP(s.recorder, s.req)

	_, expected := apiFailure(`unsupported method: "PUT"`, params.CodeMethodNotAllowed)
	s.checkResp(c, http.StatusMethodNotAllowed, "application/json", expected)
}

func (s *ResourcesHandlerSuite) TestUploadUnsupportedMethod(c *gc.C) {
	s.req.Method = "GET"

	s.uploadHandler.ServeHTTP(s.recorder, s.req)

	_, expected := apiFailure(`unsupported method: "GET"`, params.CodeMethodNotAllowed)
	s.checkResp(c, http.StatusMethodNotAllowed, "application/json", expected)
}

func (s *ResourcesHandlerSuite) TestGetSuccess(c *gc.C) {
	s.req.Method = "GET"
	s.downloadHandler.ServeHTTP(s.recorder, s.req)
	s.checkResp(c, http.StatusOK, "application/octet-stream", resourceBody)
}

	s.backend.ReturnSetResource = res

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.uploadHandler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),
	s.backend.ReturnSetResource = res

	expectedError := apiservererrors.OperationBlockedError("test block")
	s.uploadHandler.ChangeAllowedFunc = func(*http.Request) error {
		return expectedError
	}

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.uploadHandler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.ErrorResult{apiservererrors.ServerError(expectedError)})
	s.checkResp(c, http.StatusBadRequest, "application/json", string(expected))
	s.backend.ReturnSetResource = res

	req, _ := newUploadRequest(c, "spam", "a-application", uploadContent)
	s.uploadHandler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	req.Header.Set("Content-Disposition", "form-data; filename=different.ext")
	s.uploadHandler.ServeHTTP(s.recorder, req)

	_, expected := apiFailure(`incorrect extension on resource upload "different.ext", expected ".tgz"`,
		"")

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	req.URL.RawQuery += "&pendingid=some-unique-id"
	s.uploadHandler.ServeHTTP(s.recorder, req)

	expected := mustMarshalJSON(&params.UploadResult{
		Resource: api.Resource2API(res),
	s.backend.SetResourceErr = failure

	req, _ := newUploadRequest(c, "spam", "a-application", content)
	s.uploadHandler.ServeHTTP(s.recorder, req)
	s.checkResp(c, http.StatusInternalServerError, "application/json", expected)
}

