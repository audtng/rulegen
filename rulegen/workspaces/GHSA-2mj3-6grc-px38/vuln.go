package main

	errEmptyKey                  = errors.New("empty key")
	errInvalidDatatype           = errors.New("invalid data type")
	errMissingDatatype           = errors.New("missing data type")
)
			},
			Fail: false,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			d, err := New(test.Tok)
			if !assert.NoError(t, err) {
				return
			}

			if test.Fail {
				_, err := d.DissectConvert(test.Msg)
				assert.Error(t, err)
				return
			}

			r, err := d.DissectConvert(test.Msg)
			if !assert.NoError(t, err) {
		return newSkipField(id), nil
	}

	key, dataType, ordinal, length, greedy := extractKeyParts(rawKey)

	// rawKey will have | as suffix when data type is missing
	if strings.HasSuffix(rawKey, dataTypeIndicator) {
	}
}

func extractKeyParts(rawKey string) (key string, dataType string, ordinal int, length int, greedy bool) {
	m := suffixRE.FindAllStringSubmatch(rawKey, -1)

	if m[0][3] != "" {
		ordinal, _ = strconv.Atoi(m[0][3])
	}

	dataType = m[0][8]

	return m[0][1], dataType, ordinal, length, greedy
}
