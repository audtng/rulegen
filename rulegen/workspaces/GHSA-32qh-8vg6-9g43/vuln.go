package main

			})

			It("extracts the ZIP's files, generating directories, and honoring file permissions and symlinks", extractionTest)
		})

		Context("when 'unzip' is not in the PATH", func() {
			})

			It("extracts the ZIP's files, generating directories, and honoring file permissions and symlinks", extractionTest)
		})
	})

			})

			It("extracts the TGZ's files, generating directories, and honoring file permissions and symlinks", extractionTest)
		})

		Context("when 'tar' is not in the PATH", func() {
			})

			It("extracts the TGZ's files, generating directories, and honoring file permissions and symlinks", extractionTest)
		})
	})

		})

		It("extracts the TAR's files, generating directories, and honoring file permissions and symlinks", extractionTest)
	})
})
	"os"
	"os/exec"
	"path/filepath"
)

type tgzExtractor struct{}
}

func extractTarArchiveFile(header *tar.Header, dest string, input io.Reader) error {
	filePath := filepath.Join(dest, header.Name)
	fileInfo := header.FileInfo()

	if fileInfo.IsDir() {
		return os.MkdirAll(filePath, fileInfo.Mode())
	}

	err := os.MkdirAll(filepath.Dir(filePath), 0755)
	if err != nil {
		return err
	}
	"os"
	"os/exec"
	"path/filepath"
)

type zipExtractor struct{}
}

func extractZipArchiveFile(file *zip.File, dest string, input io.Reader) error {
	filePath := filepath.Join(dest, file.Name)
	fileInfo := file.FileInfo()

	if fileInfo.IsDir() {
		err := os.MkdirAll(filePath, fileInfo.Mode())
		if err != nil {
			return err
		}
	} else {
		err := os.MkdirAll(filepath.Dir(filePath), 0755)
		if err != nil {
			return err
		}
