package main

	errEmptyKey                  = errors.New("empty key")
	errInvalidDatatype           = errors.New("invalid data type")
	errMissingDatatype           = errors.New("missing data type")
	errInvalidFieldName          = errors.New("invalid field name")
)
			},
			Fail: false,
		},
		{
			Name:     "Invalid field name should fail gracefully",
			Tok:      "%{\n}",
			Msg:      "test message",
			Expected: map[string]interface{}{},
			Fail:     true,
		},
	}

	for _, test := range tests {
		t.Run(test.Name, func(t *testing.T) {
			d, err := New(test.Tok)
			if test.Fail {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}

			r, err := d.DissectConvert(test.Msg)
			if !assert.NoError(t, err) {
		return newSkipField(id), nil
	}

	key, dataType, ordinal, length, greedy, err := extractKeyParts(rawKey)
	if err != nil {
		return nil, err
	}

	// rawKey will have | as suffix when data type is missing
	if strings.HasSuffix(rawKey, dataTypeIndicator) {
	}
}

func extractKeyParts(rawKey string) (key string, dataType string, ordinal int, length int, greedy bool, err error) {
	m := suffixRE.FindAllStringSubmatch(rawKey, -1)

	// check if we have at least one match otherwise the field is invalid.
	if len(m) == 0 {
		return "", "", 0, 0, false, errInvalidFieldName
	}

	if m[0][3] != "" {
		ordinal, _ = strconv.Atoi(m[0][3])
	}

	dataType = m[0][8]

	return m[0][1], dataType, ordinal, length, greedy, nil
}
