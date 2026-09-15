package main

package config

import (
	"net/url"
	"strings"

	"github.com/pkg/errors"
)

const DefaultLunchFlowAPIURL = "https://lunchflow.app/api/v1"

type LunchFlow struct {
	// Enabled just determines whether or not Lunch Flow will be an option to
	// configure in the UI. This defaults to true as it requires no additional
	// configuration here for self-hosted users.
	Enabled bool `yaml:"enabled"`
	// AllowedApiUrls is the set of Lunch Flow API URLs this deployment is
	// permitted to contact. Comparison is exact-string.
	AllowedApiUrls []string `yaml:"allowedApiUrls"`
}

// ValidateConfig can be called at startup in order to catch problems with the
// configuration early on. If lunch flow is not enabled then this is a no-op, if
// lunch flow is enabled and there are allowed URLs specified; then this will
// validate that those URLs specified are all valid.
func (l LunchFlow) ValidateConfig() error {
	if !l.IsEnabled() {
		return nil
	}
	for _, allowed := range l.AllowedApiUrls {
		parsed, err := url.Parse(allowed)
		if err != nil {
			return errors.Wrapf(err, "configured Lunch Flow url (%s) is not valid", allowed)
		}

		// Do not allow query parameters in the URL as these will be removed when
		// requests are made!
		if len(parsed.Query()) > 0 {
			return errors.Errorf("Lunch Flow url (%s) cannot contain query parameters", allowed)
		}

		// Require a scheme to be specified
		switch strings.ToLower(parsed.Scheme) {
		case "http", "https":
			// These are considered valid!
		default:
			// Any other scheme is not considered valid here!
			return errors.Errorf("Lunch Flow url (%s) must use an http or https scheme", allowed)
		}
	}

	return nil
}

// IsEnabled only returns true if the lunch flow integration is enabled AND when
// there is at least one allowed API URLs configured.
func (l LunchFlow) IsEnabled() bool {
	return l.Enabled && len(l.AllowedApiUrls) > 0
}

// IsAllowedApiUrl returns true when the provided URL matches one of the
// configured allowed Lunch Flow API URLs. Matching is done after both the input
// URL and the allowed URLs are parsed via [url.Parse] in order to ensure
// correctness.
func (l LunchFlow) IsAllowedApiUrl(input string) bool {
	inputUrl, err := url.Parse(input)
	if err != nil {
		return false
	}

	for _, allowed := range l.AllowedApiUrls {
		allowedUrl, err := url.Parse(allowed)
		if err != nil {
			// If an allowed URL in the configuration is not even considered a valid
			// url then discard it. It will not be considered valid in the http client
			// anyway.
			continue
		}
		// Urls must be equal AFTER parsing, the [url.Parse] function does some
		// transformations here that are considered reasonable. Such as converting
		// the scheme to be lowercase.
		if allowedUrl.String() == inputUrl.String() {
			return true
		}
	}
	return false
}

	"github.com/labstack/echo/v4"
	"github.com/monetr/monetr/server/build"
	"github.com/monetr/monetr/server/icons"
)

		Price int64 `json:"price"`
	}
	var configuration struct {
		RequireLegalName        bool         `json:"requireLegalName"`
		RequirePhoneNumber      bool         `json:"requirePhoneNumber"`
		VerifyLogin             bool         `json:"verifyLogin"`
		VerifyRegister          bool         `json:"verifyRegister"`
		VerifyEmailAddress      bool         `json:"verifyEmailAddress"`
		VerifyForgotPassword    bool         `json:"verifyForgotPassword"`
		ReCAPTCHAKey            string       `json:"ReCAPTCHAKey,omitempty"`
		AllowSignUp             bool         `json:"allowSignUp"`
		AllowForgotPassword     bool         `json:"allowForgotPassword"`
		LongPollPlaidSetup      bool         `json:"longPollPlaidSetup"`
		RequireBetaCode         bool         `json:"requireBetaCode"`
		InitialPlan             *InitialPlan `json:"initialPlan"`
		BillingEnabled          bool         `json:"billingEnabled"`
		IconsEnabled            bool         `json:"iconsEnabled"`
		PlaidEnabled            bool         `json:"plaidEnabled"`
		LunchFlowEnabled        bool         `json:"lunchFlowEnabled"`
		LunchFlowAllowedAPIURLs []string     `json:"lunchFlowAllowedAPIURLs"`
		ManualEnabled           bool         `json:"manualEnabled"`
		UploadsEnabled          bool         `json:"uploadsEnabled"`
		Release                 string       `json:"release"`
		Revision                string       `json:"revision"`
		BuildType               string       `json:"buildType"`
		BuildTime               string       `json:"buildTime"`
	}

	configuration.Release = build.Release
	configuration.ManualEnabled = true
	configuration.UploadsEnabled = c.Configuration.Storage.Enabled

	configuration.LunchFlowEnabled = c.Configuration.LunchFlow.IsEnabled()
	configuration.LunchFlowAllowedAPIURLs = c.Configuration.LunchFlow.AllowedApiUrls

	return ctx.JSON(http.StatusOK, configuration)
}
	"path"
	"time"

	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/crumbs"
	"github.com/monetr/monetr/server/round"
	"github.com/pkg/errors"
)

const DateFormat = "2006-01-02"

// maxResponseBodySize caps how much of an upstream response body we will read.
// This is defense against a hostile or compromised upstream streaming an
// unbounded response to exhaust memory. 10mb is generous for realistic account
// and transaction payloads while bounding worst case allocation.
const maxResponseBodySize = 10 * 1024 * 1024

type LunchFlowAccountId = json.Number

type Account struct {
	log *slog.Logger,
	apiUrl string,
	accessToken string,
	configuration config.LunchFlow,
) (LunchFlowClient, error) {
	if !configuration.Enabled {
		log.Error("lunch flow is not enabled on this server but the client is being instantiated!",
			"bug", true,
		)
		return nil, errors.New("Lunch Flow is not enabled on this server")
	}

	if !configuration.IsAllowedApiUrl(apiUrl) {
		log.Warn("rejected Lunch Flow API URL that is not in the configured allowlist, please update your configuration if this url is valid!",
			"apiUrl", apiUrl,
		)
		return nil, errors.New("Lunch Flow API URL is not in the configured allowlist")
	}

	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	body := io.LimitReader(response.Body, maxResponseBodySize)
	if response.StatusCode != http.StatusOK {
		bodyStr, _ := io.ReadAll(body)
		return errors.Errorf("Lunch Flow request failed %s [%d]: %s", requestUrl, response.StatusCode, string(bodyStr))
	}

	if err := json.NewDecoder(body).Decode(result); err != nil {
		return errors.Wrapf(err, "failed to decode response for request %s [%d]", requestUrl, response.StatusCode)
	}

