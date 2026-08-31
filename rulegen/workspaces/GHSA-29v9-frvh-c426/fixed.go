package main

	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/controller"
	"github.com/monetr/monetr/server/database"
	"github.com/monetr/monetr/server/internal/myownsanity"
	"github.com/monetr/monetr/server/internal/source"
	"github.com/monetr/monetr/server/jobs"
	"github.com/monetr/monetr/server/logging"
	"github.com/monetr/monetr/server/stripe_helper"
	"github.com/monetr/monetr/server/ui"
	"github.com/monetr/monetr/server/zoneinfo"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)
				log.Info("config file loaded", "config", configFileName)
			}

			// As soon as we load the config try to validate it. In the future, other
			// validation functions can be added here in order to prevent
			// misconfiguration. This is done AFTER the logger is loaded so that if
			// there are validation functions in the future that require a logger in
			// order to present warnings then that can be done at the configuration
			// code level.
			if err := myownsanity.FirstError(
				configuration.LunchFlow.ValidateConfig(),
			); err != nil {
				return errors.Wrap(err, "there are configuration problems")
			}

			// TODO Move this to a configuration validation function
			if configuration.ReCAPTCHA.Enabled {
				log.Warn("DEPRECATION WARNING: ReCAPTCHA will be removed in a future release. If you are currently using it then please comment on the issue on GitHub. It is recommended to instead rate limit monetr authentication endpoints instead of using a captcha at this time.",
					"issueUrl", "https://github.com/monetr/monetr/issues/2979",
	v.SetDefault("Logging.Level", LogLevel) // Info
	// Lunch Flow is enabled by default for self-hosted deployments!
	v.SetDefault("LunchFlow.Enabled", true)
	v.SetDefault("LunchFlow.AllowedApiUrls", []string{DefaultLunchFlowAPIURL})
	v.SetDefault("KeyManagement.Provider", "plaintext")
	v.SetDefault("Plaid.Enabled", true)
	v.SetDefault("Plaid.CountryCodes", []plaid.CountryCode{plaid.COUNTRYCODE_US})
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
package config_test

import (
	"testing"

	"github.com/monetr/monetr/server/config"
	"github.com/stretchr/testify/assert"
)

func TestLunchFlow_ValidateConfig(t *testing.T) {
	t.Run("disabled is a no-op", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled:        false,
			AllowedApiUrls: []string{"not a valid url"},
		}
		assert.NoError(t, configuration.ValidateConfig())
	})

	t.Run("no allowed urls is a no-op", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled: true,
		}
		assert.NoError(t, configuration.ValidateConfig())
	})

	t.Run("valid urls are accepted", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled: true,
			AllowedApiUrls: []string{
				"https://lunchflow.app/api/v1",
				"http://lunchflow.app/api/v1",
			},
		}
		assert.NoError(t, configuration.ValidateConfig())
	})

	t.Run("unparseable url is rejected", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled:        true,
			AllowedApiUrls: []string{"https://lunchflow.app/%zz"},
		}
		assert.EqualError(t, configuration.ValidateConfig(), `configured Lunch Flow url (https://lunchflow.app/%zz) is not valid: parse "https://lunchflow.app/%zz": invalid URL escape "%zz"`)
	})

	t.Run("url with query parameters is rejected", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled:        true,
			AllowedApiUrls: []string{"https://lunchflow.app/api/v1?token=secret"},
		}
		assert.EqualError(t, configuration.ValidateConfig(), "Lunch Flow url (https://lunchflow.app/api/v1?token=secret) cannot contain query parameters")
	})

	t.Run("invalid url in config", func(t *testing.T) {
		configuration := config.LunchFlow{
			Enabled:        true,
			AllowedApiUrls: []string{"example.com"},
		}
		assert.EqualError(t, configuration.ValidateConfig(), "Lunch Flow url (example.com) must use an http or https scheme")
	})
}

func TestLunchFlow_IsAllowedApiUrl(t *testing.T) {
	t.Run("exact match is allowed", func(t *testing.T) {
		configuration := config.LunchFlow{
			AllowedApiUrls: []string{"https://lunchflow.app/api/v1"},
		}
		assert.True(t, configuration.IsAllowedApiUrl("https://lunchflow.app/api/v1"))
	})

	t.Run("mismatch is rejected", func(t *testing.T) {
		configuration := config.LunchFlow{
			AllowedApiUrls: []string{"https://lunchflow.app/api/v1"},
		}
		assert.False(t, configuration.IsAllowedApiUrl("http://169.254.169.254/latest/meta-data"))
		assert.False(t, configuration.IsAllowedApiUrl("http://127.0.0.1"))
		assert.False(t, configuration.IsAllowedApiUrl("https://lunchflow.app/api/v2"))
	})

	t.Run("empty list rejects everything", func(t *testing.T) {
		configuration := config.LunchFlow{}
		assert.False(t, configuration.IsAllowedApiUrl("https://lunchflow.app/api/v1"))
		assert.False(t, configuration.IsAllowedApiUrl(""))
	})

	t.Run("multiple entries match any of them", func(t *testing.T) {
		configuration := config.LunchFlow{
			AllowedApiUrls: []string{
				"https://lunchflow.app/api/v1",
				"https://lunchflow.compatible.app/api/v1",
			},
		}
		assert.True(t, configuration.IsAllowedApiUrl("https://lunchflow.app/api/v1"))
		assert.True(t, configuration.IsAllowedApiUrl("https://lunchflow.compatible.app/api/v1"))
		assert.False(t, configuration.IsAllowedApiUrl("https://other.lunchflow.app/api/v1"))
	})
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

					return true
				}, "Lunch Flow API URL must be a full valid URL"),
				validation.NewStringRule(
					c.Configuration.LunchFlow.IsAllowedApiUrl,
					"Lunch Flow API URL is not valid or is not in the configured allowlist",
				),
			).Required(validators.Require),
			validation.Key(
				"apiKey",
		log,
		link.ApiUrl,
		secret.Value,
		c.Configuration.LunchFlow,
	)
	if err != nil {
		return c.wrapAndReturnError(
		response.JSON().Path("$.problems.lunchFlowURL").String().IsEqual("Lunch Flow API URL must be a full valid URL")
	})

	t.Run("URL not in allowlist", func(t *testing.T) {
		for _, disallowed := range []string{
			"http://169.254.169.254/latest/meta-data",
			"http://127.0.0.1",
			"http://localhost",
			"https://attacker.example.com/api/v1",
		} {
			_, e := NewTestApplication(t)
			token := GivenIHaveToken(t, e)
			response := e.POST("/api/lunch_flow/link").
				WithCookie(TestCookieName, token).
				WithJSON(map[string]any{
					"name":         "Not Allowed",
					"lunchFlowURL": disallowed,
					"apiKey":       "foobar",
				}).
				Expect()

			response.Status(http.StatusBadRequest)
			response.JSON().Path("$.error").String().IsEqual("Invalid request")
			response.JSON().Path("$.problems.lunchFlowURL").String().IsEqual("Lunch Flow API URL is not valid or is not in the configured allowlist")
		}
	})

	t.Run("allowlist with multiple entries accepts any", func(t *testing.T) {
		config := NewTestApplicationConfig(t)
		config.LunchFlow.AllowedApiUrls = []string{
			"https://lunchflow.app/api/v1",
			"https://lunchflow.compatible.app/api/v1",
		}
		_, e := NewTestApplicationWithConfig(t, config)
		token := GivenIHaveToken(t, e)

		response := e.POST("/api/lunch_flow/link").
			WithCookie(TestCookieName, token).
			WithJSON(map[string]any{
				"name":         "Staging",
				"lunchFlowURL": "https://lunchflow.compatible.app/api/v1",
				"apiKey":       "foobar",
			}).
			Expect()
		response.Status(http.StatusOK)
	})

	t.Run("invalid api key", func(t *testing.T) {
		_, e := NewTestApplication(t)
		token := GivenIHaveToken(t, e)
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

	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/datasources/lunch_flow"
	"github.com/monetr/monetr/server/internal/mock_lunch_flow"
	"github.com/monetr/monetr/server/internal/testutils"
func TestLunchFlowClient_GetAccounts(t *testing.T) {
	t.Run("happy path, retrieve a few accounts", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		accountOne := lunch_flow.Account{
			Id:              "1234",

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	t.Run("fail to retrieve accounts", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mock_lunch_flow.MockFetchAccountsError(t)


		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")
	})
}

func TestLunchFlowClient_Constructor(t *testing.T) {
	t.Run("rejects URL that is not in the allowlist", func(t *testing.T) {
		log := testutils.GetLog(t)

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			"http://169.254.169.254/latest/meta-data",
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.EqualError(t, err, "Lunch Flow API URL is not in the configured allowlist")
		assert.Nil(t, client, "client must be nil on rejection")
	})

	t.Run("rejects when allowlist is empty", func(t *testing.T) {
		log := testutils.GetLog(t)

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{},
			},
		)
		assert.EqualError(t, err, "Lunch Flow API URL is not in the configured allowlist")
		assert.Nil(t, client)
	})

	t.Run("accepts URL that matches an allowlist entry", func(t *testing.T) {
		log := testutils.GetLog(t)

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.NoError(t, err, "must accept an allowlisted URL")
		assert.NotNil(t, client, "client must be created")
	})
}

func TestLunchFlowClient_GetBalance(t *testing.T) {
	t.Run("happy path read balance", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		expectedBalance := lunch_flow.Balance{
			Amount:   "1234.56",

		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	t.Run("fails to read balance", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		mock_lunch_flow.MockFetchBalanceError(t, "1234")


		client, err := lunch_flow.NewLunchFlowClient(
			log,
			config.DefaultLunchFlowAPIURL,
			"bogus-token",
			config.LunchFlow{
				Enabled:        true,
				AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
			},
		)
		assert.NoError(t, err, "must not return an error creating the client")
		assert.NotNil(t, client, "client must have a value")

	"github.com/benbjohnson/clock"
	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/datasources/lunch_flow/lunch_flow_jobs"
	"github.com/monetr/monetr/server/internal/fixtures"
	"github.com/monetr/monetr/server/internal/mockgen"
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    config.DefaultLunchFlowAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    config.DefaultLunchFlowAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    config.DefaultLunchFlowAPIURL,
			Status:    models.LunchFlowLinkStatusPending,
			CreatedBy: user.UserId,
		}
			AccountId: user.AccountId,
			SecretId:  secret.SecretId,
			Name:      "Test Lunch Flow Link",
			ApiUrl:    config.DefaultLunchFlowAPIURL,
			Status:    models.LunchFlowLinkStatusActive,
			CreatedBy: user.UserId,
		}
			s.log,
			link.LunchFlowLink.ApiUrl,
			secret.Value,
			ctx.Configuration().LunchFlow,
		)
		if err != nil {
			return errors.Wrap(err, "failed to create Lunch Flow API client")
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled:        true,
						AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled:        true,
						AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled:        true,
						AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled:        true,
						AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()
				context.EXPECT().Clock().Return(clock).AnyTimes()
				context.EXPECT().Configuration().Return(config.Configuration{
					LunchFlow: config.LunchFlow{
						Enabled:        true,
						AllowedApiUrls: []string{config.DefaultLunchFlowAPIURL},
					},
				}).AnyTimes()
				context.EXPECT().KMS().Return(kms).AnyTimes()

	"github.com/benbjohnson/clock"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/monetr/monetr/server/config"
	"github.com/monetr/monetr/server/consts"
	"github.com/monetr/monetr/server/internal/testutils"
	. "github.com/monetr/monetr/server/models"
	"github.com/monetr/monetr/server/repository"
		AccountId:            user.AccountId,
		SecretId:             secret.SecretId,
		Name:                 fmt.Sprintf("Lunch Flow Budget %s", gofakeit.City()),
		ApiUrl:               config.DefaultLunchFlowAPIURL,
		Status:               LunchFlowLinkStatusActive,
		LastManualSync:       nil,
		LastSuccessfulUpdate: nil,
	"testing"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/monetr/monetr/server/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Path(t *testing.T, relative string) string {
	require.NotEmpty(t, relative, "relative url cannot be empty")
	parsed, err := url.Parse(config.DefaultLunchFlowAPIURL)
	require.NoError(t, err, "must be able to parse lunch flow's default base URL")
	parsed.Path = relative
	return parsed.String()
		LunchFlow: config.LunchFlow{
			// By default lunch flow is enabled in tests, disable it to simulate
			// alternate behaviors.
			Enabled:        true,
			AllowedApiUrls: []string{"https://lunchflow.app/api/v1"},
		},
		Plaid: config.Plaid{
			Enabled:      true,
