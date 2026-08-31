package main

		}
	}

	return &nfsVolume{
		id:       id,
		server:   server,
	}, nil
}

// Given a CSI snapshot ID, return a nfsSnapshot
// sample snapshot ID:
//
	}
	return nil
}
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
