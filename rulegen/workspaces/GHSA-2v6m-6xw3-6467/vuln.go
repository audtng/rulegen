package main


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
