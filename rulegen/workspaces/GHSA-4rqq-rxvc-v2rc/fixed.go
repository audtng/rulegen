package main

			continue
		}

		// Unfortunately browsers consider a redirect Location with preceding "//" and "/\" as meaning redirect to "http(s)://REST_OF_PATH"
		// Therefore we should ignore these redirect locations to prevent open redirects
		if len(loc) > 1 && loc[0] == '/' && (loc[1] == '/' || loc[1] == '\\') {
			continue
		}

		u, err := url.Parse(loc)
		if err != nil || ((u.Scheme != "" || u.Host != "") && !strings.HasPrefix(strings.ToLower(loc), strings.ToLower(setting.AppURL))) {
			continue
