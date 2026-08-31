package main

		id := normalizeLicenseID(license)

		if cp.LicensePath != "" {
			// Clean and localize the path
			cleanPath := filepath.Clean(cp.LicensePath)
			localPath, err := filepath.Localize(cleanPath)
			if err != nil {
				return nil, fmt.Errorf("invalid license-path %q: %w", cp.LicensePath, err)
			}
			fullPath := filepath.Join(workspaceDir, localPath)

			content, err := os.ReadFile(fullPath) // #nosec G304 - Reading license file from build workspace
			if err != nil {
				return nil, fmt.Errorf("failed to read licensepath %q: %w", cp.LicensePath, err)
			}
		})
	}
}

// TestLicensingInfosPathTraversal tests that LicensingInfos properly validates
// license-path
func TestLicensingInfosPathTraversal(t *testing.T) {
	ctx := slogtest.Context(t)

	// Create a temporary directory structure simulating a workspace with files outside it
	tmpRoot := t.TempDir()
	workspaceDir := filepath.Join(tmpRoot, "workspace")
	require.NoError(t, os.MkdirAll(workspaceDir, 0o755))

	// Create a legitimate license file inside workspace
	legitimateLicense := "This is a legitimate license file"
	require.NoError(t, os.WriteFile(filepath.Join(workspaceDir, "LICENSE"), []byte(legitimateLicense), 0o644))

	// Create sensitive files
	sensitiveDir := filepath.Join(tmpRoot, "sensitive")
	require.NoError(t, os.MkdirAll(sensitiveDir, 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(sensitiveDir, "secrets.txt"), []byte("SENSITIVE"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(tmpRoot, "credentials.env"), []byte("CREDS"), 0o644))

	tests := []struct {
		name        string
		copyright   []Copyright
		expectError bool
		description string
	}{
		{
			name: "legitimate relative path within workspace",
			copyright: []Copyright{
				{License: "MIT", LicensePath: "LICENSE"},
			},
			expectError: false,
			description: "A normal license file path should work",
		},
		{
			name: "simple parent directory traversal",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: "../sensitive/secrets.txt"},
			},
			expectError: true,
			description: "Simple ../ should be blocked to prevent reading files outside workspace",
		},
		{
			name: "multiple parent directory traversal",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: "../../credentials.env"},
			},
			expectError: true,
			description: "Multiple ../ levels should be blocked",
		},
		{
			name: "traversal with valid path prefix",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: "subdir/../../sensitive/secrets.txt"},
			},
			expectError: true,
			description: "Paths that go up and then outside should be blocked",
		},
		{
			name: "absolute path",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: "/etc/passwd"},
			},
			expectError: true,
			description: "Absolute paths should be rejected",
		},
		{
			name: "absolute path to sensitive file",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: filepath.Join(sensitiveDir, "secrets.txt")},
			},
			expectError: true,
			description: "Absolute paths to any file should be rejected",
		},
		{
			name: "dot-dot in middle of path escaping",
			copyright: []Copyright{
				{License: "CustomLicense", LicensePath: "foo/../../../sensitive/secrets.txt"},
			},
			expectError: true,
			description: "Complex paths with .. that escape should be blocked",
		},
		{
			name: "relative path with ./ prefix",
			copyright: []Copyright{
				{License: "MIT", LicensePath: "./LICENSE"},
			},
			expectError: false,
			description: "Relative paths with ./ prefix within workspace should work",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pkg := Package{
				Copyright: tt.copyright,
			}

			result, err := pkg.LicensingInfos(ctx, workspaceDir)

			if tt.expectError {
				require.Error(t, err, "Expected an error but got none for: %s", tt.description)
			} else {
				require.NoError(t, err, "Expected no error but got: %v for: %s", err, tt.description)

				// For successful cases, verify the content if a file was read
				if tt.copyright[0].LicensePath != "" {
					id := normalizeLicenseID(tt.copyright[0].License)
					content, exists := result[id]
					require.True(t, exists, "Expected license info to be present")
					require.Equal(t, legitimateLicense, content, "Content should match legitimate license")
				}
			}
		})
	}
}

// TestLicensingInfosMultipleLicenses tests that one malicious path fails the entire operation
func TestLicensingInfosMultipleLicenses(t *testing.T) {
	ctx := slogtest.Context(t)

	tmpRoot := t.TempDir()
	workspaceDir := filepath.Join(tmpRoot, "workspace")
	require.NoError(t, os.MkdirAll(workspaceDir, 0o755))

	goodLicense := "This is a good license"
	require.NoError(t, os.WriteFile(filepath.Join(workspaceDir, "LICENSE"), []byte(goodLicense), 0o644))

	// Create sensitive file outside workspace
	require.NoError(t, os.WriteFile(filepath.Join(tmpRoot, "secret"), []byte("SENSITIVE"), 0o644))

	t.Run("one bad path in multiple licenses should fail entire operation", func(t *testing.T) {
		pkg := Package{
			Copyright: []Copyright{
				{License: "MIT", LicensePath: "LICENSE"},
				{License: "Apache-2.0", LicensePath: "../../secret"},
				{License: "BSD-3-Clause", LicensePath: "LICENSE"},
			},
		}

		_, err := pkg.LicensingInfos(ctx, workspaceDir)
		require.Error(t, err, "Should fail when any license path is malicious")
	})

	t.Run("all good paths should succeed", func(t *testing.T) {
		pkg := Package{
			Copyright: []Copyright{
				{License: "MIT", LicensePath: "LICENSE"},
				{License: "Apache-2.0", LicensePath: "LICENSE"},
				{License: "BSD-3-Clause", LicensePath: ""},
			},
		}

		result, err := pkg.LicensingInfos(ctx, workspaceDir)
		require.NoError(t, err)
		require.Len(t, result, 2)
		require.Equal(t, goodLicense, result["MIT"])
		require.Equal(t, goodLicense, result["Apache-2.0"])
	})
}
