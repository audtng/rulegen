package main

		}
	}

	if err := validatePath(subDir); err != nil {
		return nil, fmt.Errorf("invalid subDir %q: %v", subDir, err)
	}
	if err := validatePath(baseDir); err != nil {
		return nil, fmt.Errorf("invalid baseDir %q: %v", baseDir, err)
	}

	return &nfsVolume{
		id:       id,
		server:   server,
	}, nil
}

func validatePath(path string) error {
	if strings.Contains(path, "..") {
		return fmt.Errorf("path contains directory traversal sequence")
	}
	return nil
}

// Given a CSI snapshot ID, return a nfsSnapshot
// sample snapshot ID:
//
	}
	return nil
}

func validatePath(path string) error {
	for _, segment := range strings.Split(path, "/") {
		if segment == ".." {
			return fmt.Errorf("path contains directory traversal sequence")
		}
	}
	return nil
}
	}, nil
}

// Given a CSI snapshot ID, return a nfsSnapshot
// sample snapshot ID:
//
