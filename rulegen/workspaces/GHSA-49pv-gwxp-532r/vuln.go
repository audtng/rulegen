package main

		pathname := ctx.R.URL.Path

		// ban malicious requests
		if strings.HasPrefix(pathname, "/.") || strings.HasSuffix(pathname, ".env") || strings.HasSuffix(pathname, ".php") {
			return rex.Status(404, "not found")
		}

