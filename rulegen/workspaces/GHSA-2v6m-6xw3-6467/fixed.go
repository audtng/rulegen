package main

	assert.Len(t, appliedMacProfiles, 1)
	assert.Len(t, appliedWinProfiles, 1)
	require.Len(t, savedAppConfig.Integrations.GoogleCalendar, 1)
	assert.Equal(t, "service@example.com", savedAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["client_email"])
	assert.True(t, savedAppConfig.ActivityExpirySettings.ActivityExpiryEnabled)
	assert.Equal(t, 60, savedAppConfig.ActivityExpirySettings.ActivityExpiryWindow)
	assert.True(t, savedAppConfig.ServerSettings.AIFeaturesDisabled)
	switch {
	case config.API != nil:
		// Use the provided API.
	case config.IntegrationConfig.ApiKey.Values[fleet.GoogleCalendarEmail] == loadEmail:
		config.API = &GoogleCalendarLoadAPI{Logger: config.Logger}
	case config.IntegrationConfig.ApiKey.Values[fleet.GoogleCalendarEmail] == MockEmail:
		config.API = &GoogleCalendarMockAPI{config.Logger}
	default:
		config.API = &GoogleCalendarLowLevelAPI{logger: config.Logger}
func (c *GoogleCalendar) Configure(userEmail string) error {
	adjustedUserEmail := adjustEmail(userEmail)
	err := c.config.API.Configure(
		c.config.Context, c.config.IntegrationConfig.ApiKey.Values[fleet.GoogleCalendarEmail],
		c.config.IntegrationConfig.ApiKey.Values[fleet.GoogleCalendarPrivateKey], adjustedUserEmail,
		c.config.ServerURL,
	)
	if err != nil {
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				"client_email": loadEmail,
				"private_key":  s.server.URL,
			}},
		},
		Logger: kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout)),
	}
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				"client_email": loadEmail,
				"private_key":  s.server.URL,
			}},
		},
		Logger: kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout)),
	}
	config := &GoogleCalendarConfig{
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				fleet.GoogleCalendarEmail:      baseServiceEmail,
				fleet.GoogleCalendarPrivateKey: basePrivateKey,
			}},
		},
		Logger:    logger,
		API:       mockAPI,
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						}},
					},
				},
			},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						}},
					},
				},
			},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						}},
					},
				},
			},
	integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				"fleet": "test",
			}},
		},
	}

			stats.Organization = lic.GetOrganization()
		}
		stats.AIFeaturesDisabled = appConfig.ServerSettings.AIFeaturesDisabled
		stats.MaintenanceWindowsConfigured = len(appConfig.Integrations.GoogleCalendar) > 0 && appConfig.Integrations.GoogleCalendar[0].Domain != "" && !appConfig.Integrations.GoogleCalendar[0].ApiKey.IsEmpty()

		stats.MaintenanceWindowsEnabled = false
		teams, err := ds.ListTeams(ctx, fleet.TeamFilter{User: &fleet.User{
	for _, zdIntegration := range c.Integrations.Zendesk {
		zdIntegration.APIToken = MaskedPassword
	}
	for _, gcIntegration := range c.Integrations.GoogleCalendar {
		gcIntegration.ApiKey.SetMasked()
	}
	// // TODO(hca): confirm that we're properly masking credentials in the new endpoints
	// if c.Integrations.NDESSCEPProxy.Valid {
	// 	c.Integrations.NDESSCEPProxy.Value.Password = MaskedPassword
		for i, g := range c.Integrations.GoogleCalendar {
			gCal := *g
			clone.Integrations.GoogleCalendar[i] = &gCal
			if len(g.ApiKey.Values) > 0 {
				clone.Integrations.GoogleCalendar[i].ApiKey.Values = make(map[string]string, len(g.ApiKey.Values))
				maps.Copy(clone.Integrations.GoogleCalendar[i].ApiKey.Values, g.ApiKey.Values)
			}
		}
	}
	// // TODO(hca): do we want to cache the new grouped CAs datastore method?

import (
	"encoding/json"
	"fmt"
	"reflect"
	"testing"

		})
	}
}

func TestGoogleCalendarApiKeyMarshalUnmarshal(t *testing.T) {
	t.Run("marshal masked", func(t *testing.T) {
		key := GoogleCalendarApiKey{
			Values: map[string]string{
				"client_email": "test@example.com",
				"private_key":  "secret-key",
			},
		}
		key.SetMasked()

		data, err := json.Marshal(key)
		require.NoError(t, err)
		require.Equal(t, fmt.Sprintf(`"%s"`, MaskedPassword), string(data))
	})

	t.Run("marshal unmasked", func(t *testing.T) {
		key := GoogleCalendarApiKey{
			Values: map[string]string{
				"client_email": "test@example.com",
				"private_key":  "secret-key",
			},
		}

		data, err := json.Marshal(key)
		require.NoError(t, err)
		// Unmarshal to verify it's a valid JSON object
		var parsed map[string]string
		err = json.Unmarshal(data, &parsed)
		require.NoError(t, err)
		require.Equal(t, "test@example.com", parsed["client_email"])
		require.Equal(t, "secret-key", parsed["private_key"])
	})

	t.Run("unmarshal masked string", func(t *testing.T) {
		data := fmt.Appendf(nil, `"%s"`, MaskedPassword)
		var key GoogleCalendarApiKey
		err := json.Unmarshal(data, &key)
		require.NoError(t, err)
		require.True(t, key.IsMasked())
		require.True(t, key.IsEmpty())
	})

	t.Run("unmarshal json object", func(t *testing.T) {
		data := []byte(`{"client_email": "test@example.com", "private_key": "secret-key"}`)
		var key GoogleCalendarApiKey
		err := json.Unmarshal(data, &key)
		require.NoError(t, err)
		require.False(t, key.IsMasked())
		require.False(t, key.IsEmpty())
		require.Equal(t, "test@example.com", key.Values["client_email"])
		require.Equal(t, "secret-key", key.Values["private_key"])
	})

	t.Run("unmarshal invalid string", func(t *testing.T) {
		data := []byte(`"some-invalid-string"`)
		var key GoogleCalendarApiKey
		err := json.Unmarshal(data, &key)
		require.Error(t, err)
	})

	t.Run("unmarshal null", func(t *testing.T) {
		data := []byte(`null`)
		var key GoogleCalendarApiKey
		err := json.Unmarshal(data, &key)
		require.NoError(t, err)
		require.False(t, key.IsMasked())
		require.True(t, key.IsEmpty())
	})

	t.Run("full integration roundtrip", func(t *testing.T) {
		// Test the full struct with the API key
		intg := GoogleCalendarIntegration{
			Domain: "example.com",
			ApiKey: GoogleCalendarApiKey{
				Values: map[string]string{
					"client_email": "test@example.com",
					"private_key":  "secret-key",
				},
			},
		}

		// Marshal with unmasked key
		data, err := json.Marshal(intg)
		require.NoError(t, err)

		// Unmarshal and verify
		var parsed GoogleCalendarIntegration
		err = json.Unmarshal(data, &parsed)
		require.NoError(t, err)
		require.Equal(t, "example.com", parsed.Domain)
		require.Equal(t, "test@example.com", parsed.ApiKey.Values["client_email"])
		require.False(t, parsed.ApiKey.IsMasked())

		// Now mask and marshal again
		intg.ApiKey.SetMasked()
		data, err = json.Marshal(intg)
		require.NoError(t, err)
		require.Contains(t, string(data), `"api_key_json":"********"`)
	})
}

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	GoogleCalendarPrivateKey = "private_key"
)

// GoogleCalendarApiKey is a custom type for the Google Calendar API key JSON.
// It handles JSON marshaling/unmarshaling with support for masking sensitive data.
// When marshaled in masked state, it serializes to just "********".
// When unmarshaled, it accepts either "********" (indicating masked/preserve) or a JSON object.
type GoogleCalendarApiKey struct {
	// Values contains the actual API key fields when not masked
	Values map[string]string
	// masked indicates if this key should be serialized as masked
	masked bool
}

// MarshalJSON implements json.Marshaler. When masked, returns "********".
// Otherwise, returns the JSON object representation of the values.
func (k GoogleCalendarApiKey) MarshalJSON() ([]byte, error) {
	if k.masked {
		return json.Marshal(MaskedPassword)
	}
	return json.Marshal(k.Values)
}

// UnmarshalJSON implements json.Unmarshaler. Accepts either "********" string
// (sets masked=true) or a JSON object (populates Values).
func (k *GoogleCalendarApiKey) UnmarshalJSON(data []byte) error {
	// Handle null - treat as empty (will fail validation if required)
	if string(data) == "null" {
		k.Values = nil
		k.masked = false
		return nil
	}

	// Try to unmarshal as a string first (for masked value)
	var str string
	if err := json.Unmarshal(data, &str); err == nil {
		if str == MaskedPassword {
			k.masked = true
			k.Values = nil
			return nil
		}
		// Some other string value - invalid
		return errors.New("api_key_json must be a JSON object or the masked value")
	}

	// Try to unmarshal as a map
	var values map[string]string
	if err := json.Unmarshal(data, &values); err != nil {
		return fmt.Errorf("api_key_json must be a JSON object: %w", err)
	}

	k.Values = values
	k.masked = false
	return nil
}

// IsMasked returns true if this API key was unmarshaled from a masked value
// or has been explicitly marked as masked.
func (k GoogleCalendarApiKey) IsMasked() bool {
	return k.masked
}

// SetMasked marks this API key as masked for serialization.
func (k *GoogleCalendarApiKey) SetMasked() {
	k.masked = true
}

// IsEmpty returns true if there are no values in the API key.
func (k GoogleCalendarApiKey) IsEmpty() bool {
	return len(k.Values) == 0
}

type GoogleCalendarIntegration struct {
	Domain string               `json:"domain"`
	ApiKey GoogleCalendarApiKey `json:"api_key_json"`
}

// Integrations configures the integrations with external systems.
		invalid.Append("integrations.google_calendar", "integrating with >1 Google Workspace service account is not yet supported.")
	}
	for _, intg := range intgs {
		if email, ok := intg.ApiKey.Values[GoogleCalendarEmail]; !ok {
			invalid.Append(
				fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarEmail),
				fmt.Sprintf("%s is required", GoogleCalendarEmail),
			)
		} else {
			email = strings.TrimSpace(email)
			intg.ApiKey.Values[GoogleCalendarEmail] = email
			if email == "" {
				invalid.Append(
					fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarEmail),
				)
			}
		}
		if privateKey, ok := intg.ApiKey.Values["private_key"]; !ok {
			invalid.Append(
				fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarPrivateKey),
				fmt.Sprintf("%s is required", GoogleCalendarPrivateKey),
			)
		} else {
			privateKey = strings.TrimSpace(privateKey)
			intg.ApiKey.Values[GoogleCalendarPrivateKey] = privateKey
			if privateKey == "" {
				invalid.Append(
					fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarPrivateKey),
	appConfig.MDM.IOSUpdates.UpdateNewHosts = optjson.Bool{}
	appConfig.MDM.IPadOSUpdates.UpdateNewHosts = optjson.Bool{}

	// Handle Google Calendar API key preservation/replacement.
	// The custom GoogleCalendarApiKey type handles unmarshaling "********" as masked.
	if newAppConfig.Integrations.GoogleCalendar != nil {
		for i, newGC := range newAppConfig.Integrations.GoogleCalendar {
			if i < len(appConfig.Integrations.GoogleCalendar) {
				// If api_key_json was omitted (empty) or masked ("********"), preserve the existing value
				if newGC.ApiKey.IsEmpty() || newGC.ApiKey.IsMasked() {
					if len(oldAppConfig.Integrations.GoogleCalendar) > i {
						appConfig.Integrations.GoogleCalendar[i].ApiKey = oldAppConfig.Integrations.GoogleCalendar[i].ApiKey
					}
				} else {
					// api_key_json was provided with real values, use it
					appConfig.Integrations.GoogleCalendar[i].ApiKey = newGC.ApiKey
				}
			}
		}
	}

	// if turning off Windows MDM and Windows Migration is not explicitly set to
	// on in the same update, set it to off (otherwise, if it is explicitly set
	// to true, return an error that it can't be done when MDM is off, this is
			}
		}
	}
	// If google_calendar is null, we keep the existing setting.
	if newAppConfig.Integrations.GoogleCalendar == nil {
		appConfig.Integrations.GoogleCalendar = oldAppConfig.Integrations.GoogleCalendar
	}
					{APIToken: "zendesktoken"},
				},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
						"type":                         "service_account",
						"project_id":                   "test-project-123",
						"private_key_id":               "key-id-456",
						fleet.GoogleCalendarPrivateKey: "-----BEGIN PRIVATE KEY-----\nMIIE...\n-----END PRIVATE KEY-----",
						fleet.GoogleCalendarEmail:      "test@test-project.iam.gserviceaccount.com",
						"client_id":                    "123456789",
						"auth_uri":                     "https://accounts.google.com/o/oauth2/auth",
						"token_uri":                    "https://oauth2.googleapis.com/token",
						"auth_provider_x509_cert_url":  "https://www.googleapis.com/oauth2/v1/certs",
						"client_x509_cert_url":         "https://www.googleapis.com/robot/v1/metadata/x509/test",
						"universe_domain":              "googleapis.com",
					}}},
				},
			},
		}, nil
				require.Equal(t, ac.SMTPSettings.SMTPPassword, fleet.MaskedPassword)
				require.Equal(t, ac.Integrations.Jira[0].APIToken, fleet.MaskedPassword)
				require.Equal(t, ac.Integrations.Zendesk[0].APIToken, fleet.MaskedPassword)
				// Verify Google Calendar API key is masked (will serialize to "********")
				require.True(t, ac.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
			}
		})
	}
		})
	}
}

func TestValidAddress(t *testing.T) {
	testCases := []struct {
		name     string
		hostname string
		expected bool
	}{
		// Empty and basic cases
		{name: "empty string", hostname: "", expected: false},

		// Make sure we don't allow URLs
		{name: "http prefix", hostname: "http://example.com", expected: false},
		{name: "https prefix", hostname: "https://example.com", expected: false},
		{name: "with path", hostname: "example.com/path", expected: false},
		{name: "with query", hostname: "example.com?query=value", expected: false},
		{name: "with fragment", hostname: "example.com#fragment", expected: false},

		// Test ports are allowd
		{name: "with port", hostname: "example.com:9090", expected: true},
		{name: "port without hostname", hostname: ":9090", expected: false},
		{name: "port without hostname", hostname: "   :9090", expected: false},

		// Valid IPv4 addresses
		{name: "IPv4 localhost", hostname: "127.0.0.1", expected: true},
		{name: "IPv4 address", hostname: "192.168.1.1", expected: true},
		{name: "IPv4 all zeros", hostname: "0.0.0.0", expected: false},
		{name: "IPv4 loopback with port", hostname: "127.0.0.1:9090", expected: true},

		// Valid IPv6 addresses
		{name: "IPv6 localhost", hostname: "::1", expected: true},
		{name: "IPv6 full", hostname: "2001:0db8:85a3:0000:0000:8a2e:0370:7334", expected: true},
		{name: "IPv6 compressed", hostname: "2001:db8::1", expected: true},
		{name: "IPv6 all zeros", hostname: "::", expected: false},

		// IPv6 with brackets
		{name: "IPv6 localhost with brackets", hostname: "[::1]", expected: true},
		{name: "IPv6 with brackets", hostname: "[2001:db8::1]", expected: true},
		{name: "brackets only", hostname: "[]", expected: false},
		{name: "empty brackets", hostname: "[", expected: false},
		{name: "IPv6 locahost brackets with port", hostname: "[::1]:8089", expected: true},

		// Valid DNS hostnames
		{name: "localhostname", hostname: "localhost", expected: true},
		{name: "hostname with subdomain", hostname: "api.example.com", expected: true},
		{name: "hostname with multiple subdomains", hostname: "a.b.c.example.com", expected: true},
		{name: "hostname with numbers", hostname: "server1.example.com", expected: true},
		{name: "hostname starting with number", hostname: "1server.example.com", expected: true},
		{name: "all numeric label", hostname: "123.example.com", expected: true},
		{name: "hostname with hyphen", hostname: "my-server.example.com", expected: true},
		{name: "hostname with multiple hyphens", hostname: "my-cool-server.example.com", expected: true},
		{name: "single character label", hostname: "a.b.c", expected: true},
		{name: "single character hostname", hostname: "a", expected: true},

		// Invalid DNS hostnames - hyphen rules
		{name: "label starting with hyphen", hostname: "-example.com", expected: false},
		{name: "label ending with hyphen", hostname: "example-.com", expected: false},
		{name: "label starting and ending with hyphen", hostname: "-example-.com", expected: false},
		{name: "only hyphen label", hostname: "-.com", expected: false},

		// Invalid DNS hostnames - special characters
		{name: "hostname with underscore", hostname: "my_server.example.com", expected: false},
		{name: "hostname with space", hostname: "my server.example.com", expected: false},
		{name: "hostname with at symbol", hostname: "user@example.com", expected: false},
		{name: "hostname with exclamation", hostname: "example!.com", expected: false},

		// Invalid DNS hostnames - empty labels
		{name: "empty label (double dot)", hostname: "example..com", expected: false},
		{name: "leading dot", hostname: ".example.com", expected: false},
		{name: "trailing dot only", hostname: "example.com.", expected: false},

		// Length limits
		{name: "label exactly 63 chars", hostname: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", expected: true},
		{name: "label 64 chars (too long)", hostname: "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa.com", expected: false},

		// Real-world examples
		{name: "fleet server URL", hostname: "fleet.example.com", expected: true},
		{name: "AWS endpoint", hostname: "s3.us-west-2.amazonaws.com", expected: true},
		{name: "internal hostname", hostname: "db-primary-01.internal", expected: true},
		{name: "gibberish", hostname: "asdfasdfasdfashttps://lucas-fleet.ngrok.app", expected: false},
		{name: "gibberish II", hostname: "asdfasdfasdfashttps://lucas-fleet.ngrok.app:9800", expected: false},
		{name: "hostname with port", hostname: "example:8080", expected: true},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := validateAddress(tc.hostname)
			assert.Equal(t, tc.expected, result, "isValidHostnameAndPort(%q) = %v, want %v", tc.hostname, result, tc.expected)
		})
	}
}

// TestModifyAppConfigGoogleCalendarAPIKey tests that Google Calendar API keys
// are preserved when omitted from the request, and replaced (not merged) when provided.
func TestModifyAppConfigGoogleCalendarAPIKey(t *testing.T) {
	ds := new(mock.Store)
	svc, ctx := newTestService(t, ds, nil, nil)

	// Initial config with Google Calendar integration
	dsAppConfig := &fleet.AppConfig{
		OrgInfo: fleet.OrgInfo{
			OrgName: "Test",
		},
		ServerSettings: fleet.ServerSettings{
			ServerURL: "https://example.org",
		},
		Integrations: fleet.Integrations{
			GoogleCalendar: []*fleet.GoogleCalendarIntegration{
				{
					Domain: "example.com",
					ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
						fleet.GoogleCalendarEmail:      "test@example.com",
						fleet.GoogleCalendarPrivateKey: "original-private-key",
						"project_id":                   "original-project",
					}},
				},
			},
		},
	}

	ds.AppConfigFunc = func(ctx context.Context) (*fleet.AppConfig, error) {
		return dsAppConfig.Copy(), nil
	}
	ds.SaveAppConfigFunc = func(ctx context.Context, conf *fleet.AppConfig) error {
		*dsAppConfig = *conf
		return nil
	}
	ds.SaveABMTokenFunc = func(ctx context.Context, tok *fleet.ABMToken) error {
		return nil
	}
	ds.ListVPPTokensFunc = func(ctx context.Context) ([]*fleet.VPPTokenDB, error) {
		return []*fleet.VPPTokenDB{}, nil
	}
	ds.ListABMTokensFunc = func(ctx context.Context) ([]*fleet.ABMToken, error) {
		return []*fleet.ABMToken{}, nil
	}

	admin := &fleet.User{GlobalRole: ptr.String(fleet.RoleAdmin)}
	ctx = viewer.NewContext(ctx, viewer.Viewer{User: admin})

	t.Run("preserve API key when omitted (no changes)", func(t *testing.T) {
		// Reset to original state
		dsAppConfig.Integrations.GoogleCalendar[0].Domain = "example.com"
		dsAppConfig.Integrations.GoogleCalendar[0].ApiKey = fleet.GoogleCalendarApiKey{Values: map[string]string{
			fleet.GoogleCalendarEmail:      "test@example.com",
			fleet.GoogleCalendarPrivateKey: "original-private-key",
			"project_id":                   "original-project",
		}}

		// Update without including api_key_json (simulates frontend sending masked value)
		updateJSON := `{
			"integrations": {
				"google_calendar": [{
					"domain": "example.com"
				}]
			}
		}`

		updatedAppConfig, err := svc.ModifyAppConfig(ctx, []byte(updateJSON), fleet.ApplySpecOptions{})
		require.NoError(t, err)

		// API key should be preserved (check datastore, not returned config which is obfuscated)
		require.Len(t, dsAppConfig.Integrations.GoogleCalendar, 1)
		require.Equal(t, "example.com", dsAppConfig.Integrations.GoogleCalendar[0].Domain)
		require.Equal(t, "test@example.com", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarEmail])
		require.Equal(t, "original-private-key", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarPrivateKey])
		require.Equal(t, "original-project", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["project_id"])

		// Returned config should be obfuscated (masked)
		require.True(t, updatedAppConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	})

	t.Run("preserve API key when updating only domain", func(t *testing.T) {
		// Reset to original state
		dsAppConfig.Integrations.GoogleCalendar[0].Domain = "example.com"
		dsAppConfig.Integrations.GoogleCalendar[0].ApiKey = fleet.GoogleCalendarApiKey{Values: map[string]string{
			fleet.GoogleCalendarEmail:      "test@example.com",
			fleet.GoogleCalendarPrivateKey: "original-private-key",
			"project_id":                   "original-project",
		}}

		// Update only domain, omit api_key_json
		updateJSON := `{
			"integrations": {
				"google_calendar": [{
					"domain": "newdomain.com"
				}]
			}
		}`

		updatedAppConfig, err := svc.ModifyAppConfig(ctx, []byte(updateJSON), fleet.ApplySpecOptions{})
		require.NoError(t, err)

		// Domain should be updated, API key preserved (check datastore)
		require.Len(t, dsAppConfig.Integrations.GoogleCalendar, 1)
		require.Equal(t, "newdomain.com", dsAppConfig.Integrations.GoogleCalendar[0].Domain)
		require.Equal(t, "test@example.com", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarEmail])
		require.Equal(t, "original-private-key", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarPrivateKey])
		require.Equal(t, "original-project", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["project_id"])

		// Returned config should be obfuscated (masked)
		require.True(t, updatedAppConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	})

	t.Run("replace API key when new one provided (not merge)", func(t *testing.T) {
		// Reset to original state
		dsAppConfig.Integrations.GoogleCalendar[0].Domain = "example.com"
		dsAppConfig.Integrations.GoogleCalendar[0].ApiKey = fleet.GoogleCalendarApiKey{Values: map[string]string{
			fleet.GoogleCalendarEmail:      "test@example.com",
			fleet.GoogleCalendarPrivateKey: "original-private-key",
			"project_id":                   "original-project",
		}}

		// Provide new API key with different fields
		updateJSON := `{
			"integrations": {
				"google_calendar": [{
					"domain": "example.com",
					"api_key_json": {
						"client_email": "new@example.com",
						"private_key": "new-private-key",
						"new_field": "new-value"
					}
				}]
			}
		}`

		updatedAppConfig, err := svc.ModifyAppConfig(ctx, []byte(updateJSON), fleet.ApplySpecOptions{})
		require.NoError(t, err)

		// API key should be completely replaced (not merged) - check datastore
		require.Len(t, dsAppConfig.Integrations.GoogleCalendar, 1)
		require.Equal(t, "example.com", dsAppConfig.Integrations.GoogleCalendar[0].Domain)
		require.Equal(t, "new@example.com", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarEmail])
		require.Equal(t, "new-private-key", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarPrivateKey])
		require.Equal(t, "new-value", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["new_field"])
		// Old fields should NOT be present (confirms replacement, not merge)
		_, hasOldProject := dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["project_id"]
		require.False(t, hasOldProject, "old project_id should not be present after replacement")

		// Returned config should be obfuscated (masked)
		require.True(t, updatedAppConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	})

	t.Run("validation passes with preserved API key", func(t *testing.T) {
		// Reset to valid state
		dsAppConfig.Integrations.GoogleCalendar[0].Domain = "example.com"
		dsAppConfig.Integrations.GoogleCalendar[0].ApiKey = fleet.GoogleCalendarApiKey{Values: map[string]string{
			fleet.GoogleCalendarEmail:      "valid@example.com",
			fleet.GoogleCalendarPrivateKey: "-----BEGIN PRIVATE KEY-----\nvalid-key\n-----END PRIVATE KEY-----",
		}}

		// Update without api_key_json (should preserve valid key and pass validation)
		updateJSON := `{
			"integrations": {
				"google_calendar": [{
					"domain": "example.com"
				}]
			}
		}`

		updatedAppConfig, err := svc.ModifyAppConfig(ctx, []byte(updateJSON), fleet.ApplySpecOptions{})
		require.NoError(t, err)

		// Should succeed with preserved API key (check datastore)
		require.Len(t, dsAppConfig.Integrations.GoogleCalendar, 1)
		require.Equal(t, "valid@example.com", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarEmail])
		require.Equal(t, "-----BEGIN PRIVATE KEY-----\nvalid-key\n-----END PRIVATE KEY-----", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarPrivateKey])

		// Returned config should be obfuscated (masked)
		require.True(t, updatedAppConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	})

	t.Run("preserve API key when masked value sent", func(t *testing.T) {
		// Reset to original state
		dsAppConfig.Integrations.GoogleCalendar[0].Domain = "example.com"
		dsAppConfig.Integrations.GoogleCalendar[0].ApiKey = fleet.GoogleCalendarApiKey{Values: map[string]string{
			fleet.GoogleCalendarEmail:      "test@example.com",
			fleet.GoogleCalendarPrivateKey: "original-private-key",
			"project_id":                   "original-project",
		}}

		// Send masked api_key_json (simulates frontend sending back obfuscated value)
		updateJSON := `{
			"integrations": {
				"google_calendar": [{
					"domain": "example.com",
					"api_key_json": "********"
				}]
			}
		}`

		updatedAppConfig, err := svc.ModifyAppConfig(ctx, []byte(updateJSON), fleet.ApplySpecOptions{})
		require.NoError(t, err)

		// API key should be preserved, not overwritten with masked values (check datastore)
		require.Len(t, dsAppConfig.Integrations.GoogleCalendar, 1)
		require.Equal(t, "example.com", dsAppConfig.Integrations.GoogleCalendar[0].Domain)
		require.Equal(t, "test@example.com", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarEmail])
		require.Equal(t, "original-private-key", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values[fleet.GoogleCalendarPrivateKey])
		require.Equal(t, "original-project", dsAppConfig.Integrations.GoogleCalendar[0].ApiKey.Values["project_id"])

		// Returned config should be obfuscated (masked)
		require.True(t, updatedAppConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	})
}

	appConfig := s.getConfig()
	require.Len(t, appConfig.Integrations.GoogleCalendar, 1)
	assert.True(t, appConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	assert.Equal(t, domain, appConfig.Integrations.GoogleCalendar[0].Domain)

	// Add 2nd config -- not allowed at this time
	)
	appConfig = s.getConfig()
	require.Len(t, appConfig.Integrations.GoogleCalendar, 1)
	assert.True(t, appConfig.Integrations.GoogleCalendar[0].ApiKey.IsMasked())
	assert.Equal(t, domain, appConfig.Integrations.GoogleCalendar[0].Domain)

	// Clearing other integrations does not clear Google Calendar integration
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				fleet.GoogleCalendarEmail: "calendar-mock@example.com",
			}},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				fleet.GoogleCalendarEmail: "calendar-mock@example.com",
			}},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				fleet.GoogleCalendarEmail: calendar.MockEmail,
			}},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: fleet.GoogleCalendarApiKey{Values: map[string]string{
				fleet.GoogleCalendarEmail: calendar.MockEmail,
			}},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
