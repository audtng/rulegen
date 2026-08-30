package main

	if err := validateNonFlagArgument(src.Directory, "directory"); err != nil {
		return err
	}
	if (src.Revision != "") && (src.Directory != "") {
		cleanedDir := filepath.Clean(src.Directory)
		if strings.Contains(cleanedDir, "/") || (strings.Contains(cleanedDir, "\\")) {
			return fmt.Errorf("%q is not a valid directory, it must not contain a directory separator", src.Directory)
		}
	}
	return nil
}

			},
			isExpectedFailure: true,
		},
		{
			name: "invalid-revision-directory-combo",
			vol: &v1.Volume{
				Name: "vol1",
				VolumeSource: v1.VolumeSource{
					GitRepo: &v1.GitRepoVolumeSource{
						Repository: gitURL,
						Revision:   "main",
						Directory:  "foo/bar",
					},
				},
			},
			isExpectedFailure: true,
		},
	}

	for _, scenario := range scenarios {
