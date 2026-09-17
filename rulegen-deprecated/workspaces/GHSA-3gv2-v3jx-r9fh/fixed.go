package main

	"io"
	"os"
	"path/filepath"
	"strings"

	errs "github.com/ctfer-io/chall-manager/pkg/errors"
	"github.com/pkg/errors"
		return cd, errors.Wrap(err, "base64 decoded invalid zip archive")
	}
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		filePath, err := sanitizeArchivePath(cd, f.Name)
		if err != nil {
			return cd, &errs.ErrInternal{Sub: err}
		}

		// Save output directory i.e. the directory containing the Pulumi.yaml file,
		// the scenario entrypoint.
		}

		// Create and write the file
		if err := copyTo(f, filePath); err != nil {
			return cd, &errs.ErrInternal{Sub: err}
		}
	}
	return outDir, Validate(ctx, outDir)
}

func sanitizeArchivePath(d, t string) (v string, err error) {
	v = filepath.Join(d, t)
	if strings.HasPrefix(v, filepath.Clean(d)) {
		return v, nil
	}
	return "", fmt.Errorf("filepath is tainted: %s", t)
}

func copyTo(f *zip.File, filePath string) error {
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
