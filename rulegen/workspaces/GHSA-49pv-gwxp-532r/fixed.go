package main

		pathname := ctx.R.URL.Path

		// ban malicious requests
		if strings.HasSuffix(pathname, ".env") || strings.HasSuffix(pathname, ".php") || strings.Contains(pathname, "/.") {
			return rex.Status(404, "not found")
		}

