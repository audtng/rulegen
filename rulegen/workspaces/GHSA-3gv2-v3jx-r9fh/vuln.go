package main

	"io"
	"os"
	"path/filepath"

	errs "github.com/ctfer-io/chall-manager/pkg/errors"
	"github.com/pkg/errors"
		return cd, errors.Wrap(err, "base64 decoded invalid zip archive")
	}
	for _, f := range r.File {
		filePath := filepath.Join(cd, f.Name)
		if f.FileInfo().IsDir() {
			continue
		}

		// Save output directory i.e. the directory containing the Pulumi.yaml file,
		// the scenario entrypoint.
		}

		// Create and write the file
		if err := copy(filePath, f); err != nil {
			return cd, &errs.ErrInternal{Sub: err}
		}
	}
	return outDir, Validate(ctx, outDir)
}

func copy(filePath string, f *zip.File) error {
	outFile, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE, 0777)
	if err != nil {
		return err
