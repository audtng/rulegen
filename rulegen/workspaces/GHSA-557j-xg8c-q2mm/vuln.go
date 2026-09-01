package main

		lockfileName = "requirements.lock"
	}
	dest := filepath.Join(chartpath, lockfileName)
	return os.WriteFile(dest, data, 0644)
}

