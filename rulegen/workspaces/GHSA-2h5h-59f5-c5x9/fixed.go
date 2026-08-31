package main

	"github.com/go-openapi/swag"
	jarutils "github.com/sassoftware/relic/lib/signjar"
	"github.com/sigstore/rekor/pkg/generated/models"
	"github.com/spf13/viper"
)

const (
		return nil, nil, types.ValidationError(err)
	}

	// Checking that uncompressed metadata files are within acceptable bounds before reading into memory.
	// Checks match those performed by the relic library in the jarutils.Verify method below. For example,
	// the META-INF/MANIFEST.MF is read into memory by the relic lib, but a META-INF/LICENSE file is not.
	for _, f := range zipReader.File {
		dir, name := path.Split(strings.ToUpper(f.Name))
		if dir != "META-INF/" || name == "" || strings.LastIndex(name, ".") < 0 {
			continue
		}
		if f.UncompressedSize64 > viper.GetUint64("max_jar_metadata_size") && viper.GetUint64("max_jar_metadata_size") > 0 {
			return nil, nil, types.ValidationError(
				fmt.Errorf("uncompressed jar metadata of size %d exceeds max allowed size %d", f.UncompressedSize64, viper.GetUint64("max_jar_metadata_size")))
		}
	}

	// this ensures that the JAR is signed and the signature verifies, as
	// well as checks that the hashes in the signed manifest are all valid
	jarObjs, err := jarutils.Verify(zipReader, false)

	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/options"
	"github.com/spf13/viper"
	"gopkg.in/ini.v1"
)

		}

		if strings.HasPrefix(header.Name, ".SIGN") && pkg.Signature == nil {
			if header.Size < 0 {
				return errors.New("negative header size for .SIGN file")
			}
			if uint64(header.Size) > viper.GetUint64("max_apk_metadata_size") && viper.GetUint64("max_apk_metadata_size") > 0 {
				return fmt.Errorf("uncompressed .SIGN file size %d exceeds max allowed size %d", header.Size, viper.GetUint64("max_apk_metadata_size"))
			}
			sigBytes := make([]byte, header.Size)
			if _, err = sigReader.Read(sigBytes); err != nil && err != io.EOF {
				return fmt.Errorf("reading signature: %w", err)
		}

		if header.Name == ".PKGINFO" {
			if header.Size < 0 {
				return errors.New("negative header size for .PKGINFO file")
			}
			if uint64(header.Size) > viper.GetUint64("max_apk_metadata_size") && viper.GetUint64("max_apk_metadata_size") > 0 {
				return fmt.Errorf("uncompressed .PKGINFO file size %d exceeds max allowed size %d", header.Size, viper.GetUint64("max_apk_metadata_size"))
			}
			pkginfoContent := make([]byte, header.Size)
			if _, err = ctlReader.Read(pkginfoContent); err != nil && err != io.EOF {
				return fmt.Errorf("reading .PKGINFO: %w", err)
	rootCmd.PersistentFlags().StringSlice("enabled_api_endpoints", operationIds, "list of API endpoints to enable using operationId from openapi.yaml")

	rootCmd.PersistentFlags().Uint64("max_request_body_size", 0, "maximum size for HTTP request body, in bytes; set to 0 for unlimited")
	rootCmd.PersistentFlags().Uint64("max_jar_metadata_size", 1048576, "maximum permitted size for jar META-INF/ files, in bytes; set to 0 for unlimited")
	rootCmd.PersistentFlags().Uint64("max_apk_metadata_size", 1048576, "maximum permitted size for apk .SIGN and .PKGINFO files, in bytes; set to 0 for unlimited")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Logger.Fatal(err)
