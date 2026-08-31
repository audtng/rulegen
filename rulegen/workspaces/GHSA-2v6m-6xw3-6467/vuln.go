package main

	assert.Len(t, appliedMacProfiles, 1)
	assert.Len(t, appliedWinProfiles, 1)
	require.Len(t, savedAppConfig.Integrations.GoogleCalendar, 1)
	assert.Equal(t, "service@example.com", savedAppConfig.Integrations.GoogleCalendar[0].ApiKey["client_email"])
	assert.True(t, savedAppConfig.ActivityExpirySettings.ActivityExpiryEnabled)
	assert.Equal(t, 60, savedAppConfig.ActivityExpirySettings.ActivityExpiryWindow)
	assert.True(t, savedAppConfig.ServerSettings.AIFeaturesDisabled)
	switch {
	case config.API != nil:
		// Use the provided API.
	case config.IntegrationConfig.ApiKey[fleet.GoogleCalendarEmail] == loadEmail:
		config.API = &GoogleCalendarLoadAPI{Logger: config.Logger}
	case config.IntegrationConfig.ApiKey[fleet.GoogleCalendarEmail] == MockEmail:
		config.API = &GoogleCalendarMockAPI{config.Logger}
	default:
		config.API = &GoogleCalendarLowLevelAPI{logger: config.Logger}
func (c *GoogleCalendar) Configure(userEmail string) error {
	adjustedUserEmail := adjustEmail(userEmail)
	err := c.config.API.Configure(
		c.config.Context, c.config.IntegrationConfig.ApiKey[fleet.GoogleCalendarEmail],
		c.config.IntegrationConfig.ApiKey[fleet.GoogleCalendarPrivateKey], adjustedUserEmail,
		c.config.ServerURL,
	)
	if err != nil {
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			Domain: "example.com",
			ApiKey: map[string]string{
				"client_email": loadEmail,
				"private_key":  s.server.URL,
			},
		},
		Logger: kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout)),
	}
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			Domain: "example.com",
			ApiKey: map[string]string{
				"client_email": loadEmail,
				"private_key":  s.server.URL,
			},
		},
		Logger: kitlog.NewLogfmtLogger(kitlog.NewSyncWriter(os.Stdout)),
	}
	config := &GoogleCalendarConfig{
		Context: context.Background(),
		IntegrationConfig: &fleet.GoogleCalendarIntegration{
			ApiKey: map[string]string{
				fleet.GoogleCalendarEmail:      baseServiceEmail,
				fleet.GoogleCalendarPrivateKey: basePrivateKey,
			},
		},
		Logger:    logger,
		API:       mockAPI,
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{
						Domain: "example.com",
						ApiKey: map[string]string{
							fleet.GoogleCalendarEmail: "calendar-mock@example.com",
						},
					},
				},
			},
	integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: map[string]string{
				"fleet": "test",
			},
		},
	}

			stats.Organization = lic.GetOrganization()
		}
		stats.AIFeaturesDisabled = appConfig.ServerSettings.AIFeaturesDisabled
		stats.MaintenanceWindowsConfigured = len(appConfig.Integrations.GoogleCalendar) > 0 && appConfig.Integrations.GoogleCalendar[0].Domain != "" && len(appConfig.Integrations.GoogleCalendar[0].ApiKey) > 0

		stats.MaintenanceWindowsEnabled = false
		teams, err := ds.ListTeams(ctx, fleet.TeamFilter{User: &fleet.User{
	for _, zdIntegration := range c.Integrations.Zendesk {
		zdIntegration.APIToken = MaskedPassword
	}
	// // TODO(hca): confirm that we're properly masking credentials in the new endpoints
	// if c.Integrations.NDESSCEPProxy.Valid {
	// 	c.Integrations.NDESSCEPProxy.Value.Password = MaskedPassword
		for i, g := range c.Integrations.GoogleCalendar {
			gCal := *g
			clone.Integrations.GoogleCalendar[i] = &gCal
			clone.Integrations.GoogleCalendar[i].ApiKey = make(map[string]string, len(g.ApiKey))
			maps.Copy(clone.Integrations.GoogleCalendar[i].ApiKey, g.ApiKey)
		}
	}
	// // TODO(hca): do we want to cache the new grouped CAs datastore method?

import (
	"encoding/json"
	"reflect"
	"testing"

		})
	}
}

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	GoogleCalendarPrivateKey = "private_key"
)

type GoogleCalendarIntegration struct {
	Domain string            `json:"domain"`
	ApiKey map[string]string `json:"api_key_json"`
}

// Integrations configures the integrations with external systems.
		invalid.Append("integrations.google_calendar", "integrating with >1 Google Workspace service account is not yet supported.")
	}
	for _, intg := range intgs {
		if email, ok := intg.ApiKey[GoogleCalendarEmail]; !ok {
			invalid.Append(
				fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarEmail),
				fmt.Sprintf("%s is required", GoogleCalendarEmail),
			)
		} else {
			email = strings.TrimSpace(email)
			intg.ApiKey[GoogleCalendarEmail] = email
			if email == "" {
				invalid.Append(
					fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarEmail),
				)
			}
		}
		if privateKey, ok := intg.ApiKey["private_key"]; !ok {
			invalid.Append(
				fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarPrivateKey),
				fmt.Sprintf("%s is required", GoogleCalendarPrivateKey),
			)
		} else {
			privateKey = strings.TrimSpace(privateKey)
			intg.ApiKey[GoogleCalendarPrivateKey] = privateKey
			if privateKey == "" {
				invalid.Append(
					fmt.Sprintf("integrations.google_calendar.api_key_json.%s", GoogleCalendarPrivateKey),
	appConfig.MDM.IOSUpdates.UpdateNewHosts = optjson.Bool{}
	appConfig.MDM.IPadOSUpdates.UpdateNewHosts = optjson.Bool{}

	// if turning off Windows MDM and Windows Migration is not explicitly set to
	// on in the same update, set it to off (otherwise, if it is explicitly set
	// to true, return an error that it can't be done when MDM is off, this is
			}
		}
	}
	// If google_calendar is null, we keep the existing setting. If it's not null, we update.
	if newAppConfig.Integrations.GoogleCalendar == nil {
		appConfig.Integrations.GoogleCalendar = oldAppConfig.Integrations.GoogleCalendar
	}
					{APIToken: "zendesktoken"},
				},
				GoogleCalendar: []*fleet.GoogleCalendarIntegration{
					{ApiKey: map[string]string{fleet.GoogleCalendarPrivateKey: "google-calendar-private-key"}},
				},
			},
		}, nil
				require.Equal(t, ac.SMTPSettings.SMTPPassword, fleet.MaskedPassword)
				require.Equal(t, ac.Integrations.Jira[0].APIToken, fleet.MaskedPassword)
				require.Equal(t, ac.Integrations.Zendesk[0].APIToken, fleet.MaskedPassword)
				// Google Calendar private key is not obfuscated
				require.Equal(t, ac.Integrations.GoogleCalendar[0].ApiKey[fleet.GoogleCalendarPrivateKey], "google-calendar-private-key")
			}
		})
	}
		})
	}
}

	appConfig := s.getConfig()
	require.Len(t, appConfig.Integrations.GoogleCalendar, 1)
	assert.Equal(t, email, appConfig.Integrations.GoogleCalendar[0].ApiKey[fleet.GoogleCalendarEmail])
	assert.Equal(t, privateKey, appConfig.Integrations.GoogleCalendar[0].ApiKey[fleet.GoogleCalendarPrivateKey])
	assert.Equal(t, domain, appConfig.Integrations.GoogleCalendar[0].Domain)

	// Add 2nd config -- not allowed at this time
	)
	appConfig = s.getConfig()
	require.Len(t, appConfig.Integrations.GoogleCalendar, 1)
	assert.Equal(t, email, appConfig.Integrations.GoogleCalendar[0].ApiKey[fleet.GoogleCalendarEmail])
	assert.Equal(t, privateKey, appConfig.Integrations.GoogleCalendar[0].ApiKey[fleet.GoogleCalendarPrivateKey])
	assert.Equal(t, domain, appConfig.Integrations.GoogleCalendar[0].Domain)

	// Clearing other integrations does not clear Google Calendar integration
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: map[string]string{
				fleet.GoogleCalendarEmail: "calendar-mock@example.com",
			},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: map[string]string{
				fleet.GoogleCalendarEmail: "calendar-mock@example.com",
			},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: map[string]string{
				fleet.GoogleCalendarEmail: calendar.MockEmail,
			},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
	appCfg.Integrations.GoogleCalendar = []*fleet.GoogleCalendarIntegration{
		{
			Domain: "example.com",
			ApiKey: map[string]string{
				fleet.GoogleCalendarEmail: calendar.MockEmail,
			},
		},
	}
	err = s.ds.SaveAppConfig(ctx, appCfg)
