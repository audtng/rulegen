package main

	// Pre-scan query string to detect excessive parameters before expensive parsing.
	// This prevents DoS via url.ParseQuery allocating large maps/slices.
	if len(query) > maxQueryBufferSize {
		return boundKeySegment(query)
	}

	// Fast path: single key=value pair needs no parsing or sorting
	if strings.IndexByte(query, '&') < 0 {
		return boundKeySegment(query)
	}

	// Quick count of potential parameters (ampersands + 1)
			paramCount++
			if paramCount > maxQueryParams {
				// Too many parameters detected, hash without parsing
				return boundKeySegment(query)
			}
		}
	}

	parsed, err := url.ParseQuery(query)
	if err != nil {
		return boundKeySegment(query)
	}

	// Double-check actual parameter count after parsing
	for _, values := range parsed {
		actualCount += len(values)
		if actualCount > maxQueryParams {
			return boundKeySegment(query)
		}
	}

					*bufPtr = buf
					keyBufferPool.Put(bufPtr)
				}
				return boundKeySegment(query)
			}

			buf = append(buf, escapedKey...)
		{"at start", "no-cache, max-age=0", "no-cache", true},
		{"at end", "public, no-cache", "no-cache", true},
		{"not present", "public, max-age=0", "no-cache", false},
		{"partial match (truncated)", "no-cach", "no-cache", false}, // cspell:disable-line -- intentionally truncated directive
		{"substring of longer token", "no-cache-extended", "no-cache", false},

		// Trailing whitespace (#4143)
		cfg.Methods = ConfigDefault.Methods
	} else {
		// Normalize method names to uppercase (HTTP methods are case-sensitive
		// and c.Method() returns uppercase, e.g. "GET" not "get")
		for i, m := range cfg.Methods {
			cfg.Methods[i] = utilsstrings.ToUpper(m)
		}
	}
	if cfg.KeyGenerator == nil {
		cfg.KeyGenerator = func(c fiber.Ctx) string {
