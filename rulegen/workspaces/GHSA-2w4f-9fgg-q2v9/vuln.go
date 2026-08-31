package main

		id := normalizeLicenseID(license)

		if cp.LicensePath != "" {
			// Read license content from file
			content, err := os.ReadFile(filepath.Join(workspaceDir, cp.LicensePath)) // #nosec G304 - Reading license file from build workspace
			if err != nil {
				return nil, fmt.Errorf("failed to read licensepath %q: %w", cp.LicensePath, err)
			}
		})
	}
}
