package main

			})

			It("extracts the ZIP's files, generating directories, and honoring file permissions and symlinks", extractionTest)

			Context("with a bad zip archive", func() {
				BeforeEach(func() {
					test_helper.CreateZipArchive(extractionSrc, []test_helper.ArchiveFile{
						{
							Name: "../some-file",
							Body: "file-in-bad-dir-contents",
						},
					})
				})

				It("returns an error", func() {
					subdir := filepath.Join(extractionDest, "subdir")
					Expect(os.Mkdir(subdir, 0777)).To(Succeed())
					err := extractor.Extract(extractionSrc, subdir)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when 'unzip' is not in the PATH", func() {
			})

			It("extracts the ZIP's files, generating directories, and honoring file permissions and symlinks", extractionTest)

			Context("with a bad zip archive", func() {
				BeforeEach(func() {
					test_helper.CreateZipArchive(extractionSrc, []test_helper.ArchiveFile{
						{
							Name: "../some-file",
							Body: "file-in-bad-dir-contents",
						},
					})
				})

				It("does not insecurely extract the file outside of the provided destination", func() {
					subdir := filepath.Join(extractionDest, "subdir")
					Expect(os.Mkdir(subdir, 0777)).To(Succeed())
					err := extractor.Extract(extractionSrc, subdir)
					Expect(err).NotTo(HaveOccurred())

					Expect(filepath.Join(extractionDest, "some-file")).NotTo(BeAnExistingFile())
					Expect(filepath.Join(subdir, "some-file")).To(BeAnExistingFile())
				})
			})
		})
	})

			})

			It("extracts the TGZ's files, generating directories, and honoring file permissions and symlinks", extractionTest)

			Context("with a bad tgz archive", func() {
				BeforeEach(func() {
					test_helper.CreateTarGZArchive(extractionSrc, []test_helper.ArchiveFile{
						{
							Name: "../some-file",
							Body: "file-in-bad-dir-contents",
						},
					})
				})

				It("returns an error", func() {
					subdir := filepath.Join(extractionDest, "subdir")
					Expect(os.Mkdir(subdir, 0777)).To(Succeed())
					err := extractor.Extract(extractionSrc, subdir)
					Expect(err).To(HaveOccurred())
				})
			})
		})

		Context("when 'tar' is not in the PATH", func() {
			})

			It("extracts the TGZ's files, generating directories, and honoring file permissions and symlinks", extractionTest)

			Context("with a bad tgz archive", func() {
				BeforeEach(func() {
					test_helper.CreateTarGZArchive(extractionSrc, []test_helper.ArchiveFile{
						{
							Name: "../some-file",
							Body: "file-in-bad-dir-contents",
						},
					})
				})

				It("does not insecurely extract the file outside of the provided destination", func() {
					subdir := filepath.Join(extractionDest, "subdir")
					Expect(os.Mkdir(subdir, 0777)).To(Succeed())
					err := extractor.Extract(extractionSrc, subdir)
					Expect(err).NotTo(HaveOccurred())
					Expect(filepath.Join(extractionDest, "some-file")).NotTo(BeAnExistingFile())
					Expect(filepath.Join(subdir, "some-file")).To(BeAnExistingFile())
				})
			})
		})
	})

		})

		It("extracts the TAR's files, generating directories, and honoring file permissions and symlinks", extractionTest)

		Context("with a bad tar archive", func() {
			BeforeEach(func() {
				test_helper.CreateTarArchive(extractionSrc, []test_helper.ArchiveFile{
					{
						Name: "../some-file",
						Body: "file-in-bad-dir-contents",
					},
				})
			})

			It("does not insecurely extract the file outside of the provided destination", func() {
				subdir := filepath.Join(extractionDest, "subdir")
				Expect(os.Mkdir(subdir, 0777)).To(Succeed())
				err := extractor.Extract(extractionSrc, subdir)
				Expect(err).NotTo(HaveOccurred())
				Expect(filepath.Join(extractionDest, "some-file")).NotTo(BeAnExistingFile())
				Expect(filepath.Join(subdir, "some-file")).To(BeAnExistingFile())
			})
		})
	})
})
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
