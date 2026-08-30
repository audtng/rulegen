package main

	if err := validateNonFlagArgument(src.Directory, "directory"); err != nil {
		return err
	}
	return nil
}

			},
			isExpectedFailure: true,
		},
	}

	for _, scenario := range scenarios {
