package main

	"github.com/prebid/prebid-server/v4/macros"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
	"github.com/prebid/prebid-server/v4/util/jsonutil"
	"github.com/prebid/prebid-server/v4/util/urlutil"
)

const TAPPX_BIDDER_VERSION = "1.6"
	}}, []error{}
}

// Builds endpoint url based on adapter-specific pub settings from imp.ext.
// NOTE: The deprecated "host" param (params.Host) is intentionally unused here.
// If it is ever reintroduced, it MUST be validated with urlutil.IsSafeHost to prevent SSRF.
func (a *TappxAdapter) buildEndpointURL(params *openrtb_ext.ExtImpTappx, test int) (string, error) {

	if params.Endpoint == "" {
			Message: "Tappx key undefined",
		}
	}
	isNewEndpoint, err := regexp.Match(`^(zz|vz)[0-9]{3,}([a-z]{2,3}|test)$`, []byte(params.Endpoint))
	if err != nil {
		return "", &errortypes.BadInput{
	}
	var tappxHost string
	if isNewEndpoint {
		// Defense-in-depth: the runtime regex above is stricter than IsSafeHost,
		// so this check is currently unreachable. It guards against future relaxation
		// of the regex inadvertently allowing SSRF via subdomain injection.
		if !urlutil.IsSafeHost(params.Endpoint) {
			return "", &errortypes.BadInput{
				Message: "Invalid Tappx endpoint",
			}
		}
		tappxHost = params.Endpoint + ".pub.tappx.com/rtb/"
	} else {
		// endpoint is used as a path segment on fixed domain: ssp.api.tappx.com/rtb/v2/{endpoint}
		tappxHost = "ssp.api.tappx.com/rtb/v2/"
	}

package urlutil

import "regexp"

var safeHostPattern = regexp.MustCompile(`^[a-zA-Z0-9.-]+(:[0-9]+)?$`)

// IsSafeHost returns true for bare hostnames with an optional port.
// It intentionally rejects URL control characters such as '/', '?', '#', and '@'
// so user-supplied host values cannot rewrite the outbound request URL.
func IsSafeHost(host string) bool {
	return safeHostPattern.MatchString(host)
}
	"github.com/prebid/prebid-server/v4/macros"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
	"github.com/prebid/prebid-server/v4/util/jsonutil"
	"github.com/prebid/prebid-server/v4/util/urlutil"
)

type adapter struct {
	modifiedRequest := *request
	modifiedRequest.Imp[0].Ext = modifiedExt

	if !urlutil.IsSafeHost(params.Account) {
		return nil, []error{&errortypes.BadInput{Message: "Invalid account"}}
	}

	// create a map of macros to resolve
	endpointParams := macros.EndpointTemplateParams{AccountID: params.Account}

