package main

	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/controller"
	"github.com/monetr/monetr/server/database"
	"github.com/monetr/monetr/server/internal/source"
	"github.com/monetr/monetr/server/jobs"
	"github.com/monetr/monetr/server/logging"
	"github.com/monetr/monetr/server/stripe_helper"
	"github.com/monetr/monetr/server/ui"
	"github.com/monetr/monetr/server/zoneinfo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
				log.Info("config file loaded", "config", configFileName)
			}

			if configuration.ReCAPTCHA.Enabled {
				log.Warn("DEPRECATION WARNING: ReCAPTCHA will be removed in a future release. If you are currently using it then please comment on the issue on GitHub. It is recommended to instead rate limit monetr authentication endpoints instead of using a captcha at this time.",
					"issueUrl", "https://github.com/monetr/monetr/issues/2979",
	v.SetDefault("Logging.Level", LogLevel) // Info
	// Lunch Flow is enabled by default for self-hosted deployments!
	v.SetDefault("LunchFlow.Enabled", true)
	v.SetDefault("KeyManagement.Provider", "plaintext")
	v.SetDefault("Plaid.Enabled", true)
	v.SetDefault("Plaid.CountryCodes", []plaid.CountryCode{plaid.COUNTRYCODE_US})
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

					return true
				}, "Lunch Flow API URL must be a full valid URL"),
			).Required(validators.Require),
			validation.Key(
				"apiKey",
		log,
		link.ApiUrl,
		secret.Value,
	)
	if err != nil {
		return c.wrapAndReturnError(
		response.JSON().Path("$.problems.lunchFlowURL").String().IsEqual("Lunch Flow API URL must be a full valid URL")
	})

	t.Run("invalid api key", func(t *testing.T) {
		_, e := NewTestApplication(t)
		token := GivenIHaveToken(t, e)
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

	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/monetr/monetr/server/internal/mock_lunch_flow"
	"github.com/monetr/monetr/server/internal/testutils"
func TestLunchFlowClient_GetAccounts(t *testing.T) {
	t.Run("happy path, retrieve a few accounts", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		accountOne := lunch_flow.Account{
			Id:              "1234",

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			lunch_flow.DefaultAPIURL,
			"bogus-token",
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	t.Run("fail to retrieve accounts", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		mock_lunch_flow.MockFetchAccountsError(t)


		client, err := lunch_flow.NewLunchFlowClient(
			log,
			lunch_flow.DefaultAPIURL,
			"bogus-token",
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")
	})
}

func TestLunchFlowClient_GetBalance(t *testing.T) {
	t.Run("happy path read balance", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		expectedBalance := lunch_flow.Balance{
			Amount:   "1234.56",

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			lunch_flow.DefaultAPIURL,
			"bogus-token",
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	t.Run("fails to read balance", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.Deactivate()

		mock_lunch_flow.MockFetchBalanceError(t, "1234")


		client, err := lunch_flow.NewLunchFlowClient(
			log,
			lunch_flow.DefaultAPIURL,
			"bogus-token",
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	"github.com/benbjohnson/clock"
	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/monetr/monetr/server/datasources/lunch_flow/lunch_flow_jobs"
	"github.com/monetr/monetr/server/internal/fixtures"
	"github.com/monetr/monetr/server/internal/mockgen"
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    lunch_flow.DefaultAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    lunch_flow.DefaultAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    lunch_flow.DefaultAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    lunch_flow.DefaultAPIURL,
			Status:    models.LunchFlowLinkStatusActive,
			CreatedBy: user.UserId,
		}
			s.log,
			link.LunchFlowLink.ApiUrl,
			secret.Value,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create Lunch Flow API client")
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled: true,
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled: true,
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled: true,
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled: true,
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled: true,
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()

	"github.com/benbjohnson/clock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/monetr/monetr/server/consts"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/monetr/monetr/server/internal/testutils"
	. "github.com/monetr/monetr/server/models"
	"github.com/monetr/monetr/server/repository"
		AccountId:            user.AccountId,
		SecretId:             secret.SecretId,
		Name:                 fmt.Sprintf("Lunch Flow Budget %s", gofakeit.City()),
		ApiUrl:               lunch_flow.DefaultAPIURL,
		Status:               LunchFlowLinkStatusActive,
		LastManualSync:       nil,
		LastSuccessfulUpdate: nil,
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Path(t *testing.T, relative string) string {
	require.NotEmpty(t, relative, "relative url cannot be empty")
	parsed, err := url.Parse(lunch_flow.DefaultAPIURL)
	require.NoError(t, err, "must be able to parse lunch flow's default base URL")
	parsed.Path = relative
	return parsed.String()
		LunchFlow: config.LunchFlow{
			// By default lunch flow is enabled in tests, disable it to simulate
			// alternate behaviors.
			Enabled: true,
		},
		Plaid: config.Plaid{
			Enabled:      true,
