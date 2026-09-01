package main

		lockfileName = "requirements.lock"
	}
	dest := filepath.Join(chartpath, lockfileName)

	info, err := os.Lstat(dest)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error getting info for %q: %w", dest, err)
	} else if err == nil {
		if info.Mode()&os.ModeSymlink != 0 {
			link, err := os.Readlink(dest)
			if err != nil {
				return fmt.Errorf("error reading symlink for %q: %w", dest, err)
			}
			return fmt.Errorf("the %s file is a symlink to %q", lockfileName, link)
		}
	}

	return os.WriteFile(dest, data, 0644)
}

