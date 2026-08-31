package main

		if len(wantPatterns) != 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}
		err := opts.Platform.Download(a.DownloadURL, destDir)
		if err != nil {
			return fmt.Errorf("error downloading %s: %w", a.Name, err)
					})
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
func filepathDescendsFrom(p, dir string) bool {
	p = filepath.Clean(p)
	dir = filepath.Clean(dir)
	if dir == "." && !filepath.IsAbs(p) {
		return !strings.HasPrefix(p, ".."+string(filepath.Separator))
	}
pkg/cmd/run/download/download_test.go | 12 ++++++------
1 file changed, 6 insertions(+), 6 deletions(-)
			},
		},
		{
			name: "given artifact name contains `..`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `..`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			wantErr: "error downloading ..: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "imaginary-dir",
			wantErr: "error downloading ../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd`, verify an error about path traversal is returned",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
			wantErr: "error downloading ../../etc/passwd: would result in path traversal",
		},
		{
			name: "given artifact name contains `../../etc/passwd`, verify an error about path traversal is returned",
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
		if len(wantPatterns) != 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}

import (
	"bytes"
	"io"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

	}
}

func Test_runDownload(t *testing.T) {
	tests := []struct {
		name        string
		opts        DownloadOptions
		mockAPI     func(*mockPlatform)
		promptStubs func(*prompter.MockPrompter)
		wantErr     string
	}{
		{
			name: "download non-expired",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "./tmp",
				Names:          []string(nil),
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "artifact-1",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
					{
						Name:        "expired-artifact",
						DownloadURL: "http://download.com/expired.zip",
						Expired:     true,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.zip",
						Expired:     false,
					},
				}, nil)
				p.On("Download", "http://download.com/artifact1.zip", filepath.FromSlash("tmp/artifact-1")).Return(nil)
				p.On("Download", "http://download.com/artifact2.zip", filepath.FromSlash("tmp/artifact-2")).Return(nil)
			},
		},
		{
			name: "no valid artifacts",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
				Names:          []string(nil),
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "artifact-1",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     true,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.zip",
						Expired:     true,
					},
				}, nil)
			},
			wantErr: "no valid artifacts found to download",
		},
		{
			name: "no name matches",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
				Names:          []string{"artifact-3"},
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "artifact-1",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "no artifact matches any of the names or patterns provided",
		},
		{
			name: "no pattern matches",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: ".",
				FilePatterns:   []string{"artifiction-*"},
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "2345").Return([]shared.Artifact{
					{
						Name:        "artifact-1",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.zip",
						Expired:     false,
					},
				}, nil)
			},
			wantErr: "no artifact matches any of the names or patterns provided",
		},
		{
			name: "prompt to select artifact",
			opts: DownloadOptions{
				RunID:          "",
				DoPrompt:       true,
				DestinationDir: ".",
				Names:          []string(nil),
			},
			mockAPI: func(p *mockPlatform) {
				p.On("List", "").Return([]shared.Artifact{
					{
						Name:        "artifact-1",
						DownloadURL: "http://download.com/artifact1.zip",
						Expired:     false,
					},
					{
						Name:        "expired-artifact",
						DownloadURL: "http://download.com/expired.zip",
						Expired:     true,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.zip",
						Expired:     false,
					},
					{
						Name:        "artifact-2",
						DownloadURL: "http://download.com/artifact2.also.zip",
						Expired:     false,
					},
				}, nil)
				p.On("Download", "http://download.com/artifact2.zip", ".").Return(nil)
			},
			promptStubs: func(pm *prompter.MockPrompter) {
				pm.RegisterMultiSelect("Select artifacts to download:", nil, []string{"artifact-1", "artifact-2"},
					func(_ string, _, opts []string) ([]int, error) {
						return []int{1}, nil
					})
			},
		},
		{
			name: "given artifact name contains `..` and the DestinationDir is `.`, verify an error about path traversal is returned",
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
			name: "given artifact name contains `..` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
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
			name: "given artifact name contains `../etc/passwd` and the DestinationDir is `.`, verify an error about path traversal is returned",
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
			name: "given artifact name contains `../etc/passwd` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
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
			name: "given artifact name contains `../../etc/passwd` and the DestinationDir is `.`, verify an error about path traversal is returned",
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
			name: "given artifact name contains `../../etc/passwd` and the DestinationDir is `imaginary-dir`, verify an error about path traversal is returned",
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
			opts := &tt.opts
			ios, _, stdout, stderr := iostreams.Test()
			opts.IO = ios
			opts.Platform = newMockPlatform(t, tt.mockAPI)

			pm := prompter.NewMockPrompter(t)
			opts.Prompter = pm
				require.NoError(t, err)
			}

			assert.Equal(t, "", stdout.String())
			assert.Equal(t, "", stderr.String())
		})
	}
}

type mockPlatform struct {
	mock.Mock
}

func newMockPlatform(t *testing.T, config func(*mockPlatform)) *mockPlatform {
	m := &mockPlatform{}
	m.Test(t)
	t.Cleanup(func() {
		m.AssertExpectations(t)
	})
	if config != nil {
		config(m)
	}
	return m
}

func (p *mockPlatform) List(runID string) ([]shared.Artifact, error) {
	args := p.Called(runID)
	return args.Get(0).([]shared.Artifact), args.Error(1)
}

func (p *mockPlatform) Download(url string, dir string) error {
	args := p.Called(url, dir)
	return args.Error(0)
}
}

func filepathDescendsFrom(p, dir string) bool {
	p = filepath.Clean(p)
	dir = filepath.Clean(dir)
	if dir == "." && p == ".." {
		return false
	}
	if dir == "." && !filepath.IsAbs(p) {
		return !strings.HasPrefix(p, ".."+string(filepath.Separator))
	}
	if !strings.HasSuffix(dir, string(filepath.Separator)) {
		dir += string(filepath.Separator)
	}
	return strings.HasPrefix(p, dir)
}
pkg/cmd/run/download/download.go      |  11 +-
pkg/cmd/run/download/download_test.go | 189 +++++++++++++++++++++++++-
pkg/cmd/run/download/zip.go           |  21 ++-
pkg/cmd/run/download/zip_test.go      |  80 +++++++++++
4 files changed, 297 insertions(+), 4 deletions(-)
			}
		}
		destDir := opts.DestinationDir
		// Why do we only include the artifact name in the destination directory if there are multiple?
		if len(wantPatterns) != 0 || len(wantNames) != 1 {
			destDir = filepath.Join(destDir, a.Name)
		}

		wantErr       string
	}{
		{
			name: "download non-expired",
			opts: DownloadOptions{
				RunID:          "2345",
				DestinationDir: "./tmp",
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
			name: "no pattern matches",
			opts: DownloadOptions{
			expectedFiles: []string{},
			wantErr:       "no artifact matches any of the names or patterns provided",
		},
		{
			name: "avoid redownloading files of the same name",
			opts: DownloadOptions{
}

func filepathDescendsFrom(p, dir string) bool {
	relativePath, _ := filepath.Rel(dir, p)
	return !strings.HasPrefix(relativePath, "..")
}
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

	// track downloaded artifacts and avoid re-downloading any of the same name
	downloaded := set.NewStringSet()
	for _, a := range artifacts {
		if a.Expired {
			continue
				continue
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

	return nil
}

func matchAnyName(names []string, name string) bool {
	for _, n := range names {
		if name == n {
