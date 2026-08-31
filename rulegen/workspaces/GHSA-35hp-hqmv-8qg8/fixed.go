package main

	// Pre-scan query string to detect excessive parameters before expensive parsing.
	// This prevents DoS via url.ParseQuery allocating large maps/slices.
	if len(query) > maxQueryBufferSize {
		return boundKeySegment(escapeKeyDelimiters(query))
	}

	// Fast path: single key=value pair needs no parsing or sorting
	if strings.IndexByte(query, '&') < 0 {
		return boundKeySegment(escapeKeyDelimiters(query))
	}

	// Quick count of potential parameters (ampersands + 1)
			paramCount++
			if paramCount > maxQueryParams {
				// Too many parameters detected, hash without parsing
				return boundKeySegment(escapeKeyDelimiters(query))
			}
		}
	}

	parsed, err := url.ParseQuery(query)
	if err != nil {
		return boundKeySegment(escapeKeyDelimiters(query))
	}

	// Double-check actual parameter count after parsing
	for _, values := range parsed {
		actualCount += len(values)
		if actualCount > maxQueryParams {
			return boundKeySegment(escapeKeyDelimiters(query))
		}
	}

					*bufPtr = buf
					keyBufferPool.Put(bufPtr)
				}
				return boundKeySegment(escapeKeyDelimiters(query))
			}

			buf = append(buf, escapedKey...)
		{"at start", "no-cache, max-age=0", "no-cache", true},
		{"at end", "public, no-cache", "no-cache", true},
		{"not present", "public, max-age=0", "no-cache", false},
		{"shorter token does not match", "no-catch", "no-cache", false},
		{"substring of longer token", "no-cache-extended", "no-cache", false},

		// Trailing whitespace (#4143)
		cfg.Methods = ConfigDefault.Methods
	} else {
		// Normalize method names to uppercase (HTTP methods are case-sensitive
		// and c.Method() returns uppercase, e.g. "GET" not "get").
		// Copy first to avoid mutating the caller's slice.
		normalized := make([]string, len(cfg.Methods))
		for i, m := range cfg.Methods {
			normalized[i] = utilsstrings.ToUpper(m)
		}
		cfg.Methods = normalized
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = func(c fiber.Ctx) string {
