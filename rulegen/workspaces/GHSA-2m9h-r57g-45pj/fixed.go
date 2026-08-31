package main

		if len(wantPatterns) != 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}

		if !filepathDescendsFrom(destDir, opts.DestinationDir) {
			return fmt.Errorf("error downloading %s: would result in path traversal", a.Name)
		}

		err := opts.Platform.Download(a.DownloadURL, destDir)
		if err != nil {
			return fmt.Errorf("error downloading %s: %w", a.Name, err)
					})
			},
		},
		{
			name: "given artifact name contains `..`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "..",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `..`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "..",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "../etc/passwd",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "../etc/passwd",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "../../etc/passwd",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ../../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "../../etc/passwd",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "error downloading ../../etc/passwd: would result in path traversal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
func filepathDescendsFrom(p, dir string) bool {
	p = filepath.Clean(p)
	dir = filepath.Clean(dir)
	if dir == "." && p == ".." {
		return false
	}
	if dir == "." && !filepath.IsAbs(p) {
		return !strings.HasPrefix(p, ".."+string(filepath.Separator))
	}
pkg/cmd/run/download/download_test.go | 12 ++++++------
1 file changed, 6 insertions(+), 6 deletions(-)
			},
		},
		{
			name: "given artifact name contains `..` and the DestinationDir is `.`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `..` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd` and the DestinationDir is `.`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd` and the DestinationDir is `.`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ../../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
pkg/cmd/run/download/download.go      |   1 +
pkg/cmd/run/download/download_test.go | 524 +++++++++++++++-----------
pkg/cmd/run/download/zip.go           |  14 +-
3 files changed, 309 insertions(+), 230 deletions(-)
			}
		}
		destDir := opts.DestinationDir
		// Why do we only include the artifact name in the destination directory if there are multiple?
		if len(wantPatterns) != 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

	}
}

type testArtifact struct {
	artifact shared.Artifact
	files    []string
}

type fakePlatform struct {
	runArtifacts map[string][]testArtifact
}

func (f *fakePlatform) List(runID string) ([]shared.Artifact, error) {
	var runIds []string
	if runID != "" {
		runIds = []string{runID}
	} else {
		for k := range f.runArtifacts {
			runIds = append(runIds, k)
		}
	}

	var artifacts []shared.Artifact
	for _, id := range runIds {
		for _, a := range f.runArtifacts[id] {
			artifacts = append(artifacts, a.artifact)
		}
	}

	return artifacts, nil
}

func (f *fakePlatform) Download(url string, dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// Now to be consistent, we find the artifact with the provided URL.
	// It's a bit janky to iterate the runs, to find the right artifact
	// rather than keying directly to it, but it allows the setup of the
	// fake platform to be declarative rather than imperative.
	// Think fakePlatform { artifacts: ... } rather than fakePlatform.makeArtifactAvailable()
	for _, testArtifacts := range f.runArtifacts {
		for _, testArtifact := range testArtifacts {
			if testArtifact.artifact.DownloadURL == url {
				for _, file := range testArtifact.files {
					path := filepath.Join(dir, file)
					return os.WriteFile(path, []byte{}, 0600)
				}
			}
		}
	}

	return errors.New("no artifact matches the provided URL")
}

func Test_runDownload(t *testing.T) {
	tests := []struct {
		name          string
		opts          DownloadOptions
		platform      *fakePlatform
		promptStubs   func(*prompter.MockPrompter)
		expectedFiles []string
		wantErr       string
	}{
		{
			name: "download non-expired",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "./tmp",
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "expired-artifact",
								DownloadURL: "http://download.com/expired.zip",
								Expired:     true,
							},
							files: []string{
								"expired",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("artifact-1", "artifact-1-file"),
				filepath.Join("artifact-2", "artifact-2-file"),
			},
		},
		{
			name: "all artifacts are expired",
			opts: DownloadOptions{
				RunID: "2345",
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     true,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     true,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{},
			wantErr:       "no valid artifacts found to download",
		},
		{
			name: "no name matches",
			opts: DownloadOptions{
				RunID: "2345",
				Names: []string{"artifact-3"},
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{},
			wantErr:       "no artifact matches any of the names or patterns provided",
		},
		{
			name: "no pattern matches",
			opts: DownloadOptions{
				RunID:        "2345",
				FilePatterns: []string{"artifiction-*"},
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{},
			wantErr:       "no artifact matches any of the names or patterns provided",
		},
		{
			name: "avoid redownloading files of the same name",
			opts: DownloadOptions{
				RunID: "2345",
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("artifact-1", "artifact-1-file"),
			},
		},
		{
			name: "prompt to select artifact",
			opts: DownloadOptions{
				RunID:    "",
				DoPrompt: true,
				Names:    []string(nil),
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "expired-artifact",
								DownloadURL: "http://download.com/expired.zip",
								Expired:     true,
							},
							files: []string{
								"expired",
							},
						},
					},
					"6789": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			promptStubs: func(pm *prompter.MockPrompter) {
				pm.RegisterMultiSelect("Select artifacts to download:", nil, []string{"artifact-1", "artifact-2"},
					func(_ string, _, opts []string) ([]int, error) {
						for i, o := range opts {
							if o == "artifact-2" {
								return []int{i}, nil
							}
						}
						return nil, fmt.Errorf("no artifact-2 found in %v", opts)
					})
			},
			expectedFiles: []string{
				filepath.Join("artifact-2-file"),
			},
		},
		{
			name: "handling artifact name with path traversal exploit",
			opts: DownloadOptions{
				RunID: "2345",
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "..",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"etc/passwd",
							},
						},
					},
				},
			},
			expectedFiles: []string{},
			wantErr:       "error downloading ..: would result in path traversal",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := &tt.opts
			if opts.DestinationDir == "" {
				opts.DestinationDir = t.TempDir()
			} else {
				opts.DestinationDir = filepath.Join(t.TempDir(), opts.DestinationDir)
			}

			ios, _, stdout, stderr := iostreams.Test()
			opts.IO = ios
			opts.Platform = tt.platform

			pm := prompter.NewMockPrompter(t)
			opts.Prompter = pm
				require.NoError(t, err)
			}

			// Check that the exact number of files exist
			require.Equal(t, len(tt.expectedFiles), countFilesInDirRecursively(t, opts.DestinationDir))

			// Then check that the exact files are correct
			for _, name := range tt.expectedFiles {
				require.FileExists(t, filepath.Join(opts.DestinationDir, name))
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())
		})
	}
}

func countFilesInDirRecursively(t *testing.T, dir string) int {
	t.Helper()

	count := 0
	require.NoError(t, filepath.Walk(dir, func(_ string, info os.FileInfo, err error) error {
		require.NoError(t, err)
		if !info.IsDir() {
			count++
		}
		return nil
	}))

	return count
}
}

func filepathDescendsFrom(p, dir string) bool {
	relativePath, _ := filepath.Rel(dir, p)
	return !strings.HasPrefix(relativePath, "..")
}
pkg/cmd/run/download/download.go      |  11 +-
pkg/cmd/run/download/download_test.go | 189 +++++++++++++++++++++++++-
pkg/cmd/run/download/zip.go           |  21 ++-
pkg/cmd/run/download/zip_test.go      |  80 +++++++++++
4 files changed, 297 insertions(+), 4 deletions(-)
			}
		}
		destDir := opts.DestinationDir

		// Isolate the downloaded artifact file to avoid potential conflicts from other downloaded artifacts when:
		//
		// 1. len(wantPatterns) > 0: Any pattern can result in 2+ artifacts
		// 2. len(wantNames) == 0: User wants all artifacts regardless what they are named
		// 3. len(wantNames) > 1: User wants multiple, specific artifacts
		//
		// Otherwise if a single artifact is wanted, then the protective subdirectory is an unnecessary inconvenience.
		if len(wantPatterns) > 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}

		wantErr       string
	}{
		{
			name: "download non-expired to relative directory",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "./tmp",
				filepath.Join("artifact-2", "artifact-2-file"),
			},
		},
		{
			name: "download non-expired to absolute directory",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "/tmp",
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "expired-artifact",
								DownloadURL: "http://download.com/expired.zip",
								Expired:     true,
							},
							files: []string{
								"expired",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-2",
								DownloadURL: "http://download.com/artifact2.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-2-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("artifact-1", "artifact-1-file"),
				filepath.Join("artifact-2", "artifact-2-file"),
			},
		},
		{
			name: "all artifacts are expired",
			opts: DownloadOptions{
			expectedFiles: []string{},
			wantErr:       "no artifact matches any of the names or patterns provided",
		},
		{
			name: "pattern matches",
			opts: DownloadOptions{
				RunID:        "2345",
				FilePatterns: []string{"artifact-*"},
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "non-artifact-2",
								DownloadURL: "http://download.com/non-artifact-2.zip",
								Expired:     false,
							},
							files: []string{
								"non-artifact-2-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-3",
								DownloadURL: "http://download.com/artifact3.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-3-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("artifact-1", "artifact-1-file"),
				filepath.Join("artifact-3", "artifact-3-file"),
			},
		},
		{
			name: "no pattern matches",
			opts: DownloadOptions{
			expectedFiles: []string{},
			wantErr:       "no artifact matches any of the names or patterns provided",
		},
		{
			name: "want specific single artifact",
			opts: DownloadOptions{
				RunID: "2345",
				Names: []string{"non-artifact-2"},
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "non-artifact-2",
								DownloadURL: "http://download.com/non-artifact-2.zip",
								Expired:     false,
							},
							files: []string{
								"non-artifact-2-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-3",
								DownloadURL: "http://download.com/artifact3.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-3-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("non-artifact-2-file"),
			},
		},
		{
			name: "want specific multiple artifacts",
			opts: DownloadOptions{
				RunID: "2345",
				Names: []string{"artifact-1", "artifact-3"},
			},
			platform: &fakePlatform{
				runArtifacts: map[string][]testArtifact{
					"2345": {
						{
							artifact: shared.Artifact{
								Name:        "artifact-1",
								DownloadURL: "http://download.com/artifact1.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-1-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "non-artifact-2",
								DownloadURL: "http://download.com/non-artifact-2.zip",
								Expired:     false,
							},
							files: []string{
								"non-artifact-2-file",
							},
						},
						{
							artifact: shared.Artifact{
								Name:        "artifact-3",
								DownloadURL: "http://download.com/artifact3.zip",
								Expired:     false,
							},
							files: []string{
								"artifact-3-file",
							},
						},
					},
				},
			},
			expectedFiles: []string{
				filepath.Join("artifact-1", "artifact-1-file"),
				filepath.Join("artifact-3", "artifact-3-file"),
			},
		},
		{
			name: "avoid redownloading files of the same name",
			opts: DownloadOptions{
}

func filepathDescendsFrom(p, dir string) bool {
	// Regardless of the logic below, `p` is never allowed to be current directory `.` or parent directory `..`
	// however we check explicitly here before filepath.Rel() which doesn't cover all cases.
	p = filepath.Clean(p)

	if p == "." || p == ".." {
		return false
	}

	// filepathDescendsFrom() takes advantage of filepath.Rel() to determine if `p` is descended from `dir`:
	//
	// 1. filepath.Rel() calculates a path to traversal from fictious `dir` to `p`.
	// 2. filepath.Rel() errors in a handful of cases where absolute and relative paths are compared as well as certain traversal edge cases
	//    For more information, https://github.com/golang/go/blob/00709919d09904b17cfe3bfeb35521cbd3fb04f8/src/path/filepath/path_test.go#L1510-L1515
	// 3. If the path to traverse `dir` to `p` requires `..`, then we know it is not descend from / contained in `dir`
	//
	// As-is, this function requires the caller to ensure `p` and `dir` are either 1) both relative or 2) both absolute.
	relativePath, err := filepath.Rel(dir, p)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(relativePath, "..")
}
			},
			want: false,
		},
		{
			name: "deny parent directory filename (`..`) escaping absolute directory",
			args: args{
				p:   filepath.FromSlash(".."),
				dir: filepath.FromSlash("/var/logs/"),
			},
			want: false,
		},
		{
			name: "deny parent directory filename (`..`) escaping current directory",
			args: args{
				p:   filepath.FromSlash(".."),
				dir: filepath.FromSlash("."),
			},
			want: false,
		},
		{
			name: "deny parent directory filename (`..`) escaping parent directory",
			args: args{
				p:   filepath.FromSlash(".."),
				dir: filepath.FromSlash(".."),
			},
			want: false,
		},
		{
			name: "deny parent directory filename (`..`) escaping relative directory",
			args: args{
				p:   filepath.FromSlash(".."),
				dir: filepath.FromSlash("relative-dir"),
			},
			want: false,
		},
		{
			name: "deny current directory filename (`.`) in absolute directory",
			args: args{
				p:   filepath.FromSlash("."),
				dir: filepath.FromSlash("/var/logs/"),
			},
			want: false,
		},
		{
			name: "deny current directory filename (`.`) in current directory",
			args: args{
				p:   filepath.FromSlash("."),
				dir: filepath.FromSlash("."),
			},
			want: false,
		},
		{
			name: "deny current directory filename (`.`) in parent directory",
			args: args{
				p:   filepath.FromSlash("."),
				dir: filepath.FromSlash(".."),
			},
			want: false,
		},
		{
			name: "deny current directory filename (`.`) in relative directory",
			args: args{
				p:   filepath.FromSlash("."),
				dir: filepath.FromSlash("relative-dir"),
			},
			want: false,
		},
		{
			name: "relative path, absolute dir",
			args: args{
				p:   filepath.FromSlash("whatever"),
				dir: filepath.FromSlash("/a/b/c"),
			},
			want: false,
		},
		{
			name: "absolute path, relative dir",
			args: args{
				p:   filepath.FromSlash("/a/b/c"),
				dir: filepath.FromSlash("whatever"),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
pkg/cmd/run/download/download.go | 34 ++++++++++++++++++++++----------
1 file changed, 24 insertions(+), 10 deletions(-)
	opts.IO.StartProgressIndicator()
	defer opts.IO.StopProgressIndicator()

	// track downloaded artifacts and avoid re-downloading any of the same name, isolate if multiple artifacts
	downloaded := set.NewStringSet()
	isolateArtifacts := isolateArtifacts(wantNames, wantPatterns)

	for _, a := range artifacts {
		if a.Expired {
			continue
				continue
			}
		}

		destDir := opts.DestinationDir
		if isolateArtifacts {
			destDir = filepath.Join(destDir, a.Name)
		}

	return nil
}

func isolateArtifacts(wantNames []string, wantPatterns []string) bool {
	if len(wantPatterns) > 0 {
		// Patterns can match multiple artifacts
		return true
	}

	if len(wantNames) == 0 {
		// All artifacts wanted regardless what they are named
		return true
	}

	if len(wantNames) > 1 {
		// Multiple, specific artifacts wanted
		return true
	}

	return false
}

func matchAnyName(names []string, name string) bool {
	for _, n := range names {
		if name == n {
