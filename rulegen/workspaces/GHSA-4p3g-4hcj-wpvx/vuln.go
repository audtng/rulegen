package main

	"github.com/prebid/prebid-server/v4/macros"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
	"github.com/prebid/prebid-server/v4/util/jsonutil"
)

const TAPPX_BIDDER_VERSION = "1.6"
	}}, []error{}
}

// Builds enpoint url based on adapter-specific pub settings from imp.ext
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
		tappxHost = params.Endpoint + ".pub.tappx.com/rtb/"
	} else {
		tappxHost = "ssp.api.tappx.com/rtb/v2/"
	}

	"github.com/prebid/prebid-server/v4/macros"
	"github.com/prebid/prebid-server/v4/openrtb_ext"
	"github.com/prebid/prebid-server/v4/util/jsonutil"
)

type adapter struct {
	modifiedRequest := *request
	modifiedRequest.Imp[0].Ext = modifiedExt

	// create a map of macros to resolve
	endpointParams := macros.EndpointTemplateParams{AccountID: params.Account}

