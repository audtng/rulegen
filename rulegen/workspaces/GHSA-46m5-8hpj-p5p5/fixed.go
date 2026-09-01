package main


// buildReceiverIntegrations builds a list of integration notifiers off of a receiver config.
func (am *alertmanager) buildReceiverIntegrations(receiver *alertingNotify.APIReceiver, tmpl *alertingTemplates.Template) ([]*alertingNotify.Integration, error) {
	err := patchNewSecureFields(context.Background(), receiver, am.decryptFn)
	if err != nil {
		return nil, err
	}
	receiverCfg, err := alertingNotify.BuildReceiverConfiguration(context.Background(), receiver, am.decryptFn)
	if err != nil {
		return nil, err
	return integrations, nil
}

func patchNewSecureFields(ctx context.Context, api *alertingNotify.APIReceiver, decrypt alertingNotify.GetDecryptedValueFn) error {
	for _, integration := range api.Integrations {
		switch integration.Type {
		case "dingding":
			err := patchSettingsFromSecureSettings(ctx, integration, "url", decrypt)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func patchSettingsFromSecureSettings(ctx context.Context, integration *alertingNotify.GrafanaIntegrationConfig, key string, decrypt alertingNotify.GetDecryptedValueFn) error {
	var encrypted string
	var ok bool
	if encrypted, ok = integration.SecureSettings[key]; !ok {
		return nil
	}
	decoded, err := decode(encrypted)
	if err != nil {
		return err
	}
	settings := map[string]any{}
	err = json.Unmarshal(integration.Settings, &settings)
	if err != nil {
		return err
	}
	currentValue, ok := settings[key]
	currentString := ""
	if ok {
		currentString, _ = currentValue.(string)
	}
	secretValue := decrypt(ctx, map[string][]byte{key: decoded}, key, currentString)
	if secretValue == currentString {
		return nil
	}
	settings[key] = secretValue
	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}
	integration.Settings = data
	return nil
}

// PutAlerts receives the alerts and then sends them through the corresponding route based on whenever the alert has a receiver embedded or not
func (am *alertmanager) PutAlerts(_ context.Context, postableAlerts apimodels.PostableAlerts) error {
	alerts := make(alertingNotify.PostableAlerts, 0, len(postableAlerts.PostableAlerts))
				SecureSettings:        gr.SecureSettings,
			})
		}
		recv := &alertingNotify.APIReceiver{
			ConfigReceiver: r.Receiver,
			GrafanaIntegrations: alertingNotify.GrafanaIntegrations{
				Integrations: integrations,
			},
		}
		err := patchNewSecureFields(ctx, recv, am.decryptFn)
		if err != nil {
			return nil, err
		}
		receivers = append(receivers, recv)
	}
	var alert *alertingNotify.TestReceiversConfigAlertParams
	if c.Alert != nil {
					Placeholder:  "https://oapi.dingtalk.com/robot/send?access_token=xxxxxxxxx",
					PropertyName: "url",
					Required:     true,
					Secure:       true,
				},
				{
					Label:        "Message Type",
