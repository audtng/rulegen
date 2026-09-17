package main

	"os"
	"os/exec"
	"path/filepath"

	securejoin "github.com/cyphar/filepath-securejoin"
)

type tgzExtractor struct{}
}

func extractTarArchiveFile(header *tar.Header, dest string, input io.Reader) error {
	filePath, err := securejoin.SecureJoin(dest, header.Name)
	if err != nil {
		return err
	}
	fileInfo := header.FileInfo()

	if fileInfo.IsDir() {
		return os.MkdirAll(filePath, fileInfo.Mode())
	}

	err = os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	"os"
	"os/exec"
	"path/filepath"

	securejoin "github.com/cyphar/filepath-securejoin"
)

type zipExtractor struct{}
}

func extractZipArchiveFile(file *zip.File, dest string, input io.Reader) error {
	filePath, err := securejoin.SecureJoin(dest, file.Name)
	fileInfo := file.FileInfo()

	if fileInfo.IsDir() {
		err = os.MkdirAll(filePath, fileInfo.Mode())
		if err != nil {
			return err
		}
	} else {
		err = os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}
