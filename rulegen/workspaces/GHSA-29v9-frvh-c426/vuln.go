package main

package config

type LunchFlow struct {
	// Enabled just determines whether or not Lunch Flow will be an option to
	// configure in the UI. This defaults to true as it requires no additional
	// configuration here for self-hosted users.
	Enabled bool `yaml:"enabled"`
}

	"github.com/labstack/echo/v4"
	"github.com/monetr/monetr/server/build"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/monetr/monetr/server/icons"
)

		Price int64 `json:"price"`
	}
	var configuration struct {
		RequireLegalName       bool         `json:"requireLegalName"`
		RequirePhoneNumber     bool         `json:"requirePhoneNumber"`
		VerifyLogin            bool         `json:"verifyLogin"`
		VerifyRegister         bool         `json:"verifyRegister"`
		VerifyEmailAddress     bool         `json:"verifyEmailAddress"`
		VerifyForgotPassword   bool         `json:"verifyForgotPassword"`
		ReCAPTCHAKey           string       `json:"ReCAPTCHAKey,omitempty"`
		AllowSignUp            bool         `json:"allowSignUp"`
		AllowForgotPassword    bool         `json:"allowForgotPassword"`
		LongPollPlaidSetup     bool         `json:"longPollPlaidSetup"`
		RequireBetaCode        bool         `json:"requireBetaCode"`
		InitialPlan            *InitialPlan `json:"initialPlan"`
		BillingEnabled         bool         `json:"billingEnabled"`
		IconsEnabled           bool         `json:"iconsEnabled"`
		PlaidEnabled           bool         `json:"plaidEnabled"`
		LunchFlowEnabled       bool         `json:"lunchFlowEnabled"`
		LunchFlowDefaultAPIURL string       `json:"lunchFlowDefaultAPIURL"`
		ManualEnabled          bool         `json:"manualEnabled"`
		UploadsEnabled         bool         `json:"uploadsEnabled"`
		Release                string       `json:"release"`
		Revision               string       `json:"revision"`
		BuildType              string       `json:"buildType"`
		BuildTime              string       `json:"buildTime"`
	}

	configuration.Release = build.Release
	configuration.ManualEnabled = true
	configuration.UploadsEnabled = c.Configuration.Storage.Enabled

	configuration.LunchFlowEnabled = c.Configuration.LunchFlow.Enabled
	configuration.LunchFlowDefaultAPIURL = lunch_flow.DefaultAPIURL

	return ctx.JSON(http.StatusOK, configuration)
}
	"path"
	"time"

	"github.com/monetr/monetr/server/crumbs"
	"github.com/monetr/monetr/server/round"
	"github.com/pkg/errors"
)

const DefaultAPIURL = "https://lunchflow.app/api/v1"

const DateFormat = "2006-01-02"

type LunchFlowAccountId = json.Number

type Account struct {
	log *slog.Logger,
	apiUrl string,
	accessToken string,
) (LunchFlowClient, error) {
	parsedUrl, err := url.Parse(apiUrl)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyStr, _ := io.ReadAll(response.Body)
		return errors.Errorf("Lunch Flow request failed %s [%d]: %s", requestUrl, response.StatusCode, string(bodyStr))
	}

	if err := json.NewDecoder(response.Body).Decode(result); err != nil {
		return errors.Wrapf(err, "failed to decode response for request %s [%d]", requestUrl, response.StatusCode)
	}

