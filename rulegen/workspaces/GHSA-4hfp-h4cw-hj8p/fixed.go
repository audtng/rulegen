package main

	chart "helm.sh/helm/v4/pkg/chart/v2"
)

// MaxDecompressedChartSize is the maximum size of a chart archive that will be
// decompressed. This is the decompressed size of all the files.
// The default value is 100 MiB.
var MaxDecompressedChartSize int64 = 100 * 1024 * 1024 // Default 100 MiB

// MaxDecompressedFileSize is the size of the largest file that Helm will attempt to load.
// The size of the file is the decompressed version of it when it is stored in an archive.
var MaxDecompressedFileSize int64 = 5 * 1024 * 1024 // Default 5 MiB

var drivePathPattern = regexp.MustCompile(`^[a-zA-Z]:/`)

// FileLoader loads a chart from a file

	files := []*BufferedFile{}
	tr := tar.NewReader(unzipped)
	remainingSize := MaxDecompressedChartSize
	for {
		b := bytes.NewBuffer(nil)
		hd, err := tr.Next()
			return nil, errors.New("chart yaml not in base directory")
		}

		if hd.Size > remainingSize {
			return nil, fmt.Errorf("decompressed chart is larger than the maximum size %d", MaxDecompressedChartSize)
		}

		if hd.Size > MaxDecompressedFileSize {
			return nil, fmt.Errorf("decompressed chart file %q is larger than the maximum file size %d", hd.Name, MaxDecompressedFileSize)
		}

		limitedReader := io.LimitReader(tr, remainingSize)

		bytesWritten, err := io.Copy(b, limitedReader)
		if err != nil {
			return nil, err
		}

		remainingSize -= bytesWritten
		// When the bytesWritten are less than the file size it means the limit reader ended
		// copying early. Here we report that error. This is important if the last file extracted
		// is the one that goes over the limit. It assumes the Size stored in the tar header
		// is correct, something many applications do.
		if bytesWritten < hd.Size || remainingSize <= 0 {
			return nil, fmt.Errorf("decompressed chart is larger than the maximum size %d", MaxDecompressedChartSize)
		}

		data := bytes.TrimPrefix(b.Bytes(), utf8bom)

		files = append(files, &BufferedFile{Name: n, Data: data})
			return fmt.Errorf("cannot load irregular file %s as it has file mode type bits set", name)
		}

		if fi.Size() > MaxDecompressedFileSize {
			return fmt.Errorf("chart file %q is larger than the maximum file size %d", fi.Name(), MaxDecompressedFileSize)
		}

		data, err := os.ReadFile(name)
		if err != nil {
			return errors.Wrapf(err, "error reading %s", n)
