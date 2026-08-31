package main

	// DefaultMaxFileBytes is the default maximum size of files that will be unpacked. Larger files are ignored.
	// The max is large because some files are hundreds of megabytes.
	DefaultMaxFileBytes = 1024 * 1024 * 1024 // 1GB
	// DefaultMaxSymlinkDepth is the default maximum symlink depth. This should be no more than 8,
	// since that is the maximum number of symlinks the os.Root API will handle. From the os.Root API,
	// "8 is __POSIX_SYMLOOP_MAX (the minimum allowed value for SYMLOOP_MAX), and a common limit".
	DefaultMaxSymlinkDepth = 6
)

	chainLayers    []*chainLayer
	config         *Config
	size           int64
	root           *os.Root
	ExtractDir     string
	BaseImageIndex int
}
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	// OpenRoot assumes that the provided directory is trusted. In this case, we created the
	// imageExtractionPath directory, so it is indeed trusted.
	root, err := os.OpenRoot(imageExtractionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open root directory: %w", err)
	}
	// Close the root directory at the end of the function, since no more files will be unpacked
	// afterward.
	defer root.Close()

	baseImageIndex, err := findBaseImageIndex(history)
	if err != nil {
		baseImageIndex = -1
	outputImage := &Image{
		chainLayers:    chainLayers,
		config:         config,
		root:           root,
		ExtractDir:     imageExtractionPath,
		BaseImageIndex: baseImageIndex,
	}
	// Add the root directory to each chain layer. If this is not done, then the virtual paths won't
	// be rooted, and traversal in the virtual filesystem will be broken.
	if err := addRootDirectoryToChainLayers(outputImage.chainLayers, imageExtractionPath); err != nil {
		return nil, handleImageError(outputImage, fmt.Errorf("failed to add root directory to chain layers: %w", err))
	}

	// Since the layers are in reverse order, the v1LayerIndex starts at the last layer and works
		layerDir := layerDirectory(i)

		// Create the chain layer directory if it doesn't exist.
		if err := root.Mkdir(layerDir, dirPermission); err != nil && !errors.Is(err, fs.ErrExist) {
			return nil, handleImageError(outputImage, fmt.Errorf("failed to create chain layer directory: %w", err))
		}

		if v1LayerIndex < 0 {
			return nil, handleImageError(outputImage, fmt.Errorf("mismatch between v1 layers and chain layers, on v1 layer index %d, but only %d v1 layers", v1LayerIndex, len(v1Layers)))
		}

		chainLayersToFill := chainLayers[i:]
		v1Layer := v1Layers[v1LayerIndex]
		layerReader, err := v1Layer.Uncompressed()
		if err != nil {
			return nil, handleImageError(outputImage, err)
		}
		v1LayerIndex--

			defer layerReader.Close()

			tarReader := tar.NewReader(layerReader)
			layerSize, err := fillChainLayersWithFilesFromTar(outputImage, tarReader, layerDir, chainLayersToFill)
			if err != nil {
				return fmt.Errorf("failed to fill chain layer with v1 layer tar: %w", err)
			}
		}()

		if err != nil {
			return nil, handleImageError(outputImage, err)
		}
	}


// handleImageError cleans up the image and returns the provided error. The image is cleaned up
// regardless of the error, as the image is in an invalid state if an error is returned.
func handleImageError(image *Image, err error) error {
	if image != nil {
		_ = image.CleanUp()
	}
	return err
}

// validateHistory makes sure that the number of v1 layers matches the number of non-empty history
// fillChainLayersWithFilåesFromTar fills the chain layers with the files found in the tar. The
// chainLayersToFill are the chain layers that will be filled with the files via the virtual
// filesystem.
func fillChainLayersWithFilesFromTar(img *Image, tarReader *tar.Reader, layerDir string, chainLayersToFill []*chainLayer) (int64, error) {
	if len(chainLayersToFill) == 0 {
		return 0, errors.New("no chain layers provided, this should not happen")
	}
			virtualPath = "/" + path.Join(dirname, basename)
		}

		// If the file already exists in the current chain layer, then skip it. This is done because
		// the tar file could be read multiple times to handle required symlinks.
		if currentChainLayer.fileNodeTree.Get(virtualPath) != nil {
		var newNode *fileNode
		switch header.Typeflag {
		case tar.TypeDir:
			newNode, err = img.handleDir(virtualPath, layerDir, header, isWhiteout)
		case tar.TypeReg:
			newNode, err = img.handleFile(virtualPath, layerDir, tarReader, header, isWhiteout)
		case tar.TypeSymlink, tar.TypeLink:
			newNode, err = img.handleSymlink(virtualPath, layerDir, header, isWhiteout)
		default:

		// If the virtual path has any directories and those directories have not been populated, then
		// populate them with file nodes.
		populateEmptyDirectoryNodes(virtualPath, layerDir, img.ExtractDir, chainLayersToFill)

		// In each outer loop, a layer is added to each relevant output chainLayer slice. Because the
		// outer loop is looping backwards (latest layer first), we ignore any files that are already in
}

// handleDir creates the directory specified by path, if it doesn't exist.
func (img *Image) handleDir(virtualPath, layerDir string, header *tar.Header, isWhiteout bool) (*fileNode, error) {
	realFilePath := filepath.Join(img.ExtractDir, layerDir, filepath.FromSlash(virtualPath))
	if _, err := img.root.Stat(filepath.Join(layerDir, filepath.FromSlash(virtualPath))); err != nil {
		if err := os.MkdirAll(realFilePath, dirPermission); err != nil {
			return nil, fmt.Errorf("failed to create directory with realFilePath %s: %w", realFilePath, err)
		}

// handleFile creates the file specified by path, and then copies the contents of the tarReader into
// the file.
func (img *Image) handleFile(virtualPath, layerDir string, tarReader *tar.Reader, header *tar.Header, isWhiteout bool) (*fileNode, error) {
	realFilePath := filepath.Join(img.ExtractDir, layerDir, filepath.FromSlash(virtualPath))
	parentDirectory := filepath.Dir(realFilePath)
	if err := os.MkdirAll(parentDirectory, dirPermission); err != nil {
		return nil, fmt.Errorf("failed to create parent directory %s: %w", parentDirectory, err)
	}

	// Write all files as read/writable by the current user, inaccessible by anyone else. Actual
	// permission bits are stored in FileNode.
	f, err := img.root.OpenFile(filepath.Join(layerDir, filepath.FromSlash(virtualPath)), os.O_CREATE|os.O_RDWR, filePermission)
	if err != nil {
		return nil, err
	}
import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
//  4. Devise a pathtree that will return an error when inserting a path. Make sure that Load()
//     returns an error.
func TestFromV1Image(t *testing.T) {
	tests := []struct {
		name                  string
		v1Image               v1.Image
		config                *Config
		wantChainLayerEntries []chainLayerEntries
		wantErr               bool
		wantNonZeroSize       bool
				},
				errorOnConfigFile: true,
			},
			config: DefaultConfig(),
			wantChainLayerEntries: []chainLayerEntries{
				chainLayerEntries{
					filepathContentPairs: []filepathContentPair{},
				},
			},
		},
		{
			name: "image with error on layers",
				},
				errorOnLayers: true,
			},
			config:  DefaultConfig(),
			wantErr: true,
		},
		{
			name: "image with single package",
			v1Image: constructImageWithTarEntries(t, []*tarEntry{
				{
					Header: &tar.Header{
						Name: "etc/os-release",
						Mode: 0777,
						Size: int64(len(osContents)),
					},
					Data: bytes.NewBufferString(osContents),
				},
				{
					Header: &tar.Header{
						Name: "var/lib/dpkg/status",
						Mode: 0777,
						Size: int64(len("Package: fake-package-name\nVersion: 1.0\nStatus: install ok installed")),
					},
					Data: bytes.NewBufferString("Package: fake-package-name\nVersion: 1.0\nStatus: install ok installed"),
				},
			}),
			wantChainLayerEntries: []chainLayerEntries{
				{
					filepathContentPairs: []filepathContentPair{
						{
							filepath: "/etc/os-release",
							content:  osContents,
						},
						{
							filepath: "/var/lib/dpkg/status",
							content:  "Package: fake-package-name\nVersion: 1.0\nStatus: install ok installed",
						},
					},
				},
			},
			config:          DefaultConfig(),
			wantNonZeroSize: true,
		},
		{
					},
				},
			},
			config:  DefaultConfig(),
			wantErr: true,
		},
		{
			name: "image attempting trampoline path traversal attack",
			v1Image: constructImageWithTarEntries(t, []*tarEntry{
				{
					Header: &tar.Header{
						Name: "escape/poc.txt",
						Mode: 0777,
						Size: int64(len("👻")),
					},
					Data: bytes.NewBufferString("👻"),
				},
				{
					Header: &tar.Header{
						Name:     "usr/share/doc/a/copyright",
						Typeflag: tar.TypeSymlink,
						Linkname: "/trampoline",
						Mode:     0777,
					},
				},
				{
					Header: &tar.Header{
						Name:     "trampoline/",
						Typeflag: tar.TypeSymlink,
						Linkname: ".",
						Mode:     0777,
					},
				},
				{
					Header: &tar.Header{
						Name:     "usr/share/doc/b/copyright",
						Typeflag: tar.TypeSymlink,
						Linkname: "/escape",
						Mode:     0777,
					},
				},
				{
					Header: &tar.Header{
						Name:     "escape/",
						Typeflag: tar.TypeSymlink,
						Linkname: "trampoline/trampoline/trampoline/trampoline/trampoline/../../../../tmp",
						Mode:     0777,
					},
				},
				{
					Header: &tar.Header{
						Name:     "usr/share/doc/c/copyright",
						Typeflag: tar.TypeSymlink,
						Linkname: "/escape/poc.txt",
						Mode:     0777,
					},
				},
			}),
			config: &Config{
				MaxFileBytes:    DefaultMaxFileBytes,
				MaxSymlinkDepth: DefaultMaxSymlinkDepth,
				Requirer: require.NewFileRequirerPaths([]string{
					"/usr/share/doc/a/copyright",
					"/usr/share/doc/b/copyright",
					"/usr/share/doc/c/copyright",
				}),
			},
			wantChainLayerEntries: []chainLayerEntries{
				{
					filepathContentPairs: []filepathContentPair{
						{
							filepath: "/escape/poc.txt",
							content:  "👻",
						},
					},
				},
			},
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// some. This is needed to compare the files found after the extractor runs.
			filesInTmpWant := scalibrFilesInTmp(t)

			gotImage, gotErr := FromV1Image(tc.v1Image, tc.config)

			if tc.wantErr {
				if gotErr == nil {
					t.Fatalf("FromV1Image() returned nil error, but want non-nil error")
				}
				return
			}

			if gotErr != nil {
				t.Fatalf("FromV1Image() returned error: %v", gotErr)
			}

			if tc.wantNonZeroSize && gotImage.Size() == 0 {
				t.Errorf("got image with size 0, but want non-zero size")
			}

			// Make sure the expected files are in the chain layers.
			chainLayers, err := gotImage.ChainLayers()
			if err != nil {
				t.Fatalf("ChainLayers() returned error: %v", err)
			}

			// If the number of chain layers does not match the number of expected chain layer entries,
			// then there is no point in continuing the test.
			if len(chainLayers) != len(tc.wantChainLayerEntries) {
				t.Fatalf("ChainLayers() returned incorrect number of chain layers: got %d chain layers, want %d chain layers", len(chainLayers), len(tc.wantChainLayerEntries))
			}

			for i := range chainLayers {
				chainLayer := chainLayers[i]
				wantChainLayerEntries := tc.wantChainLayerEntries[i]

				if wantChainLayerEntries.ignore {
					continue
				}

				compareChainLayerEntries(t, chainLayer, wantChainLayerEntries, nil)
			}

			if gotImage != nil {
				if err := gotImage.CleanUp(); err != nil {
					t.Fatalf("CleanUp() returned error: %v", err)
	}
}

// tarEntry represents a single entry in a tarball. It contains the header and data for the entry.
// If the data is nil, the entry will be written without any content.
type tarEntry struct {
	Header *tar.Header
	Data   io.Reader
}

func constructImageWithTarEntries(t *testing.T, tarEntries []*tarEntry) v1.Image {
	t.Helper()

	var buf bytes.Buffer
	w := tar.NewWriter(&buf)

	// Put them in a single tarball to make a single layer and put that layer in an empty image to
	// make the minimal image that will work.
	for _, entry := range tarEntries {
		if err := w.WriteHeader(entry.Header); err != nil {
			t.Fatalf("couldn't write header for %s: %v", entry.Header.Name, err)
		}
		if entry.Data != nil {
			if _, err := io.Copy(w, entry.Data); err != nil {
				t.Fatalf("writing content for %s: %v", entry.Header.Name, err)
			}
		}
	}
	w.Close()
		return io.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
	})
	if err != nil {
		t.Fatalf("unable to create layer: %v", err)
	}

	image, err := mutate.AppendLayers(empty.Image, layer)
	if err != nil {
		t.Fatalf("unable append layer to image: %v", err)
	}

	return image
}
	return image
}

// filesInTmp returns the list of filenames in tmpDir.
func filesInTmp(t *testing.T, tmpDir string) []string {
	t.Helper()

	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		filenames = append(filenames, f.Name())
	}
	return filenames
