package main

		router.Methods(http.MethodHead).HandlerFunc(
			collectAPIStats("headbucket", maxClients(gz(httpTraceAll(api.HeadBucketHandler)))))
		// PostPolicy
		router.Methods(http.MethodPost).MatcherFunc(func(r *http.Request, _ *mux.RouteMatch) bool {
			return isRequestPostPolicySignatureV4(r)
		}).HandlerFunc(collectAPIStats("postpolicybucket", maxClients(gz(httpTraceHdrs(api.PostPolicyBucketHandler)))))
		// DeleteMultipleObjects
		router.Methods(http.MethodPost).HandlerFunc(
			collectAPIStats("deletemultipleobjects", maxClients(gz(httpTraceAll(api.DeleteMultipleObjectsHandler))))).Queries("delete", "")
	"encoding/hex"
	"errors"
	"io"
	"mime"
	"net/http"
	"net/url"
	"strconv"

// Verify if request has AWS Post policy Signature Version '4'.
func isRequestPostPolicySignatureV4(r *http.Request) bool {
	mediaType, _, err := mime.ParseMediaType(r.Header.Get(xhttp.ContentType))
	if err != nil {
		return false
	}
	return mediaType == "multipart/form-data" && r.Method == http.MethodPost
}

// Verify if the request has AWS Streaming Signature Version '4'. This is only valid for 'PUT' operation.
