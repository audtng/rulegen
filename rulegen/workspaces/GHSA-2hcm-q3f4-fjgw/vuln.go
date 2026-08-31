package main

	// DefaultMaxFileBytes is the default maximum size of files that will be unpacked. Larger files are ignored.
	// The max is large because some files are hundreds of megabytes.
	DefaultMaxFileBytes = 1024 * 1024 * 1024 // 1GB
	// DefaultMaxSymlinkDepth is the default maximum symlink depth.
	DefaultMaxSymlinkDepth = 6
)

	chainLayers    []*chainLayer
	config         *Config
	size           int64
	ExtractDir     string
	BaseImageIndex int
}
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}

	baseImageIndex, err := findBaseImageIndex(history)
	if err != nil {
		baseImageIndex = -1
	outputImage := &Image{
		chainLayers:    chainLayers,
		config:         config,
		ExtractDir:     imageExtractionPath,
		BaseImageIndex: baseImageIndex,
	}
	// Add the root directory to each chain layer. If this is not done, then the virtual paths won't
	// be rooted, and traversal in the virtual filesystem will be broken.
	if err := addRootDirectoryToChainLayers(outputImage.chainLayers, imageExtractionPath); err != nil {
		return handleImageError(outputImage, fmt.Errorf("failed to add root directory to chain layers: %w", err))
	}

	// Since the layers are in reverse order, the v1LayerIndex starts at the last layer and works
		layerDir := layerDirectory(i)

		// Create the chain layer directory if it doesn't exist.
		// Use filepath here as it is a path that will be written to disk.
		dirPath := filepath.Join(imageExtractionPath, layerDir)
		if err := os.Mkdir(dirPath, dirPermission); err != nil && !errors.Is(err, fs.ErrExist) {
			return handleImageError(outputImage, fmt.Errorf("failed to create chain layer directory: %w", err))
		}

		if v1LayerIndex < 0 {
			return handleImageError(outputImage, fmt.Errorf("mismatch between v1 layers and chain layers, on v1 layer index %d, but only %d v1 layers", v1LayerIndex, len(v1Layers)))
		}

		chainLayersToFill := chainLayers[i:]
		v1Layer := v1Layers[v1LayerIndex]
		layerReader, err := v1Layer.Uncompressed()
		if err != nil {
			return handleImageError(outputImage, err)
		}
		v1LayerIndex--

			defer layerReader.Close()

			tarReader := tar.NewReader(layerReader)
			layerSize, err := fillChainLayersWithFilesFromTar(outputImage, tarReader, layerDir, dirPath, chainLayersToFill)
			if err != nil {
				return fmt.Errorf("failed to fill chain layer with v1 layer tar: %w", err)
			}
		}()

		if err != nil {
			return handleImageError(outputImage, err)
		}
	}


// handleImageError cleans up the image and returns the provided error. The image is cleaned up
// regardless of the error, as the image is in an invalid state if an error is returned.
func handleImageError(image *Image, err error) (*Image, error) {
	if image != nil {
		_ = image.CleanUp()
	}
	return nil, err
}

// validateHistory makes sure that the number of v1 layers matches the number of non-empty history
// fillChainLayersWithFilåesFromTar fills the chain layers with the files found in the tar. The
// chainLayersToFill are the chain layers that will be filled with the files via the virtual
// filesystem.
func fillChainLayersWithFilesFromTar(img *Image, tarReader *tar.Reader, layerDir string, dirPath string, chainLayersToFill []*chainLayer) (int64, error) {
	if len(chainLayersToFill) == 0 {
		return 0, errors.New("no chain layers provided, this should not happen")
	}
			virtualPath = "/" + path.Join(dirname, basename)
		}

		// realFilePath is where the file will be written to disk. filepath.Join will convert
		// any forward slashes to the appropriate OS specific path separator.
		realFilePath := filepath.Join(dirPath, filepath.FromSlash(cleanedFilePath))

		// If the file already exists in the current chain layer, then skip it. This is done because
		// the tar file could be read multiple times to handle required symlinks.
		if currentChainLayer.fileNodeTree.Get(virtualPath) != nil {
		var newNode *fileNode
		switch header.Typeflag {
		case tar.TypeDir:
			newNode, err = img.handleDir(realFilePath, virtualPath, layerDir, header, isWhiteout)
		case tar.TypeReg:
			newNode, err = img.handleFile(realFilePath, virtualPath, layerDir, tarReader, header, isWhiteout)
		case tar.TypeSymlink, tar.TypeLink:
			newNode, err = img.handleSymlink(virtualPath, layerDir, header, isWhiteout)
		default:

		// If the virtual path has any directories and those directories have not been populated, then
		// populate them with file nodes.
		populateEmptyDirectoryNodes(virtualPath, layerDir, dirPath, chainLayersToFill)

		// In each outer loop, a layer is added to each relevant output chainLayer slice. Because the
		// outer loop is looping backwards (latest layer first), we ignore any files that are already in
}

// handleDir creates the directory specified by path, if it doesn't exist.
func (img *Image) handleDir(realFilePath, virtualPath, layerDir string, header *tar.Header, isWhiteout bool) (*fileNode, error) {
	if _, err := os.Stat(realFilePath); err != nil {
		if err := os.MkdirAll(realFilePath, dirPermission); err != nil {
			return nil, fmt.Errorf("failed to create directory with realFilePath %s: %w", realFilePath, err)
		}

// handleFile creates the file specified by path, and then copies the contents of the tarReader into
// the file.
func (img *Image) handleFile(realFilePath, virtualPath, layerDir string, tarReader *tar.Reader, header *tar.Header, isWhiteout bool) (*fileNode, error) {
	parentDirectory := filepath.Dir(realFilePath)
	if err := os.MkdirAll(parentDirectory, dirPermission); err != nil {
		return nil, fmt.Errorf("failed to create parent directory %s: %w", parentDirectory, err)
	}
	// Write all files as read/writable by the current user, inaccessible by anyone else
	// Actual permission bits are stored in FileNode
	f, err := os.OpenFile(realFilePath, os.O_CREATE|os.O_RDWR, filePermission)

	if err != nil {
		return nil, err
	}
import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
//  4. Devise a pathtree that will return an error when inserting a path. Make sure that Load()
//     returns an error.
func TestFromV1Image(t *testing.T) {
	fakeImage, err := constructImage("1.0", "fake-package-name")
	if err != nil {
		t.Fatalf("Failed to construct image: %v", err)
	}

	tests := []struct {
		name                  string
		v1Image               v1.Image
		wantChainLayerEntries []chainLayerEntries
		wantErr               bool
		wantNonZeroSize       bool
				},
				errorOnConfigFile: true,
			},
		},
		{
			name: "image with error on layers",
				},
				errorOnLayers: true,
			},
			wantErr: true,
		},
		{
			name:    "image with single package",
			v1Image: *fakeImage,
			wantChainLayerEntries: []chainLayerEntries{
				{
					filepathContentPairs: []filepathContentPair{
						{
							filepath: "etc/os-release",
							content:  osContents,
						},
						{
							filepath: "var/lib/dpkg/status",
							content:  "Package: fake-package-name\nVersion: 1.0\nStatus: install ok installed",
						},
					},
				},
			},
			wantNonZeroSize: true,
		},
		{
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// some. This is needed to compare the files found after the extractor runs.
			filesInTmpWant := scalibrFilesInTmp(t)

			gotImage, gotErr := FromV1Image(tc.v1Image, DefaultConfig())

			if tc.wantErr != (gotErr != nil) {
				t.Errorf("FromV1Image() returned error: %v", gotErr)
			}

			if tc.wantNonZeroSize && gotImage.Size() == 0 {
				t.Errorf("got image with size 0, but want non-zero size")
			}

			if gotImage != nil {
				if err := gotImage.CleanUp(); err != nil {
					t.Fatalf("CleanUp() returned error: %v", err)
	}
}

// constructImage constructs a fake v1.Image that contains a single layer with two files:
//   - The file `osContents` defines the OS, which allows the image to be scanned.
//   - The file `statusContents` defines the packages in the image, in this probe it contains a
//     single fake package with a specified version, so it is only affected by the Note created in
//     this execution.
//
// Put them in a single tarball to make a single layer and put that layer in an empty image to
// make the minimal image that will work.
func constructImage(version, fakePackageName string) (*v1.Image, error) {
	// The file containing the fake package version.
	statusContents := fmt.Sprintf("Package: %s\nVersion: %s\nStatus: install ok installed", fakePackageName, version)

	var buf bytes.Buffer
	w := tar.NewWriter(&buf)

	// These files are the minimal set to create an image that will be scanned by Container Analysis.
	// - The file `osContents` defines the OS, which allows the image to be scanned.
	// - The file `statusContents` defines the packages in the image, in this probe it contains a
	//   single fake package with a specified version, so it is only affected by the Note created in
	//   this execution.
	//
	// Put them in a single tarball to make a single layer and put that layer in an empty image to
	// make the minimal image that will work.
	files := []struct {
		name, contents string
	}{
		{"etc/os-release", osContents},
		{"var/lib/dpkg/status", statusContents},
	}
	for _, file := range files {
		hdr := &tar.Header{
			Name:     file.name,
			Mode:     0600,
			Size:     int64(len(file.contents)),
			Typeflag: tar.TypeReg,
		}
		if err := w.WriteHeader(hdr); err != nil {
			return nil, fmt.Errorf("couldn't write header for %s: %w", file.name, err)
		}
		if _, err := w.Write([]byte(file.contents)); err != nil {
			return nil, fmt.Errorf("couldn't write %s: %w", file.name, err)
		}
	}
	w.Close()
		return io.NopCloser(bytes.NewBuffer(buf.Bytes())), nil
	})
	if err != nil {
		return nil, fmt.Errorf("unable to create layer: %w", err)
	}
	image, err := mutate.AppendLayers(empty.Image, layer)
	return &image, err
}
	return image
}

// filesInTmp returns the list of filenames in /tmp.
func filesInTmp(t *testing.T, tmpDir string) []string {
	t.Helper()

	}

	for _, f := range files {
		filenames = append(filenames, f.Name())
	}
	return filenames
