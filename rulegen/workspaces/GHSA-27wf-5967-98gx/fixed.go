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

