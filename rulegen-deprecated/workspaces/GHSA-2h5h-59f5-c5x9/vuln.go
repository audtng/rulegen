package main

	"github.com/go-openapi/swag"
	jarutils "github.com/sassoftware/relic/lib/signjar"
	"github.com/sigstore/rekor/pkg/generated/models"
)

const (
		return nil, nil, types.ValidationError(err)
	}

	// this ensures that the JAR is signed and the signature verifies, as
	// well as checks that the hashes in the signed manifest are all valid
	jarObjs, err := jarutils.Verify(zipReader, false)

	"github.com/sigstore/sigstore/pkg/signature"
	"github.com/sigstore/sigstore/pkg/signature/options"
	"gopkg.in/ini.v1"
)

		}

		if strings.HasPrefix(header.Name, ".SIGN") && pkg.Signature == nil {
			sigBytes := make([]byte, header.Size)
			if _, err = sigReader.Read(sigBytes); err != nil && err != io.EOF {
				return fmt.Errorf("reading signature: %w", err)
		}

		if header.Name == ".PKGINFO" {
			pkginfoContent := make([]byte, header.Size)
			if _, err = ctlReader.Read(pkginfoContent); err != nil && err != io.EOF {
				return fmt.Errorf("reading .PKGINFO: %w", err)
	rootCmd.PersistentFlags().StringSlice("enabled_api_endpoints", operationIds, "list of API endpoints to enable using operationId from openapi.yaml")

	rootCmd.PersistentFlags().Uint64("max_request_body_size", 0, "maximum size for HTTP request body, in bytes; set to 0 for unlimited")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		log.Logger.Fatal(err)
