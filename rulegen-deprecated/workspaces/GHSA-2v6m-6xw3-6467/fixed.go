package main


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
