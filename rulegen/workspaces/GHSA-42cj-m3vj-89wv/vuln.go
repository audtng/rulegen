package main

		return "", nil
	}

	if !p.AllowCrossNamespace && strings.HasSuffix(serversTransportName, providerNamespaceSeparator+providerName) {
		// Since we are not able to know if another namespace is in the name (namespace-name@kubernetescrd),
		// if the provider namespace kubernetescrd is used,
		// we don't allow this format to avoid cross namespace references.
		return "", fmt.Errorf("invalid reference to serversTransport %s: namespace-name@kubernetescrd format is not allowed when crossnamespace is disallowed", serversTransportName)
	}

	if strings.Contains(serversTransportName, providerNamespaceSeparator) {
		return serversTransportName, nil
	}

