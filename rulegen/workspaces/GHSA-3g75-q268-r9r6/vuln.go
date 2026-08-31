package main

import (
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	// this is purposefully done before attempting to
	// decrypt the materials
	var cekFunc internal.CEKEntry
	if objectMetadata.CEKAlg == internal.AESGCMNoPadding {
		cekFunc = internal.NewAESGCMContentCipher
	} else if strings.Contains(objectMetadata.CEKAlg, "AES/CBC") {
		if !m.client.Options.EnableLegacyUnauthenticatedModes {
			return out, metadata, fmt.Errorf("configure client with enable legacy unauthenticated modes set to true to decrypt with %s", objectMetadata.CEKAlg)
		}
		cekFunc = internal.NewAESCBCContentCipher
	} else {
		return out, metadata, fmt.Errorf("invalid content encryption algorithm found in metadata: %s", objectMetadata.CEKAlg)
	}

	cipherKey, err := objectMetadata.GetDecodedKey()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded key for materials: %w", err)
	}
	iv, err := objectMetadata.GetDecodedIV()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded IV for materials: %w", err)
	}
	matDesc, err := objectMetadata.GetMatDesc()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get Material Description for materials: %w", err)
	}

	// S3 server will encode metadata with non-US-ASCII characters
	// Decode it here to avoid parsing/decryption failure
	decodedMatDesc, err := customS3Decoder(matDesc)
	if err != nil {
		return out, metadata, fmt.Errorf("error while decoding Material Description: %w", err)
	}

	decryptMaterialsRequest := materials.DecryptMaterialsRequest{
		cipherKey,
		iv,
		decodedMatDesc,
		objectMetadata.KeyringAlg,
		objectMetadata.CEKAlg,
		objectMetadata.TagLen,
	}
	decryptMaterials, err := m.client.Options.CryptographicMaterialsManager.DecryptMaterials(ctx, decryptMaterialsRequest)
		return out, metadata, fmt.Errorf("error while decrypting materials: %w", err)
	}

	cipher, err := cekFunc(*decryptMaterials)
	reader, err := cipher.DecryptContents(result.Body)
	if err != nil {
		return out, metadata, err
	}

	result.Body = reader
	out.Result = result

	return out, metadata, err
}
	"encoding/hex"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"testing"
)

func TestDecryptionClientV3_GetObject(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	}
}

func TestDecryptionClientV3_GetObject_V1Interop_KMS_AESCBC(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "7ItX9CTGNWWegC62RlaNu6EJ3+J9yGO7yAqDNU4CdeA=", `"}`))
	}))

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDecryptionClientV3_GetObject_V1Interop_KMS_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "Hrjrkkt/vQwMYtqvK6+MiXh3xiMvviL1Ks7w2mgsJgU=", `"}`))
	}))
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	}
}

func TestDecryptionClientV3_GetObject_OnlyDecryptsRegisteredAlgorithms(t *testing.T) {
	httpClientFactory := func() *awstesting.MockHttpClient {
		b, err := hex.DecodeString("1bd0271b25951fdef3dbe51a9b7af85f66b311e091aa10a346655068f657b9da9acc0843ea0522b0d1ae4a25a31b13605dd1ac5d002db8965d9d4652fd602693")
		if err != nil {
				tConfig.HTTPClient = httpClientFactory()
				s3Client := s3.NewFromConfig(tConfig)

				client, err := New(s3Client, cmm)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
	}
}

func TestDecryptionClientV3_CheckValidCryptographicMaterialsManager(t *testing.T) {
	_, err := materials.NewCryptographicMaterialsManager(nil)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}
import (
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"io"
)

// DefaultMinFileSize is used to check whether we want to write to a temp file
// or store the data in memory.
const DefaultMinFileSize = 1024 * 512 * 5

// EncryptionContext is used to extract Encryption Context to use on a per-request basis
const EncryptionContext = "EncryptionContext"

	if err != nil {
		return out, metadata, err
	}
	cipher, err := internal.NewAESGCMContentCipher(*cryptoMaterials)
	if err != nil {
		return out, metadata, err
	}

	stream := reqCopy.GetStream()
import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	return cryptoMaterials, err
}

func TestEncryptionClientV3_PutObject_KMSCONTEXT_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm)

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
	}
}

func TestEncryptionClientV3_PutObject_KMSCONTEXT_AESGCM_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm)

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:   aws.String("test-bucket"),
		t.Errorf("expected no error, got %v", err)
	}
}

import (
	"context"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// S3EncryptionClientV3 provides client-side encryption for S3.
// The client embeds a default client to provide support for control plane operations
// which do not involve encryption.
type S3EncryptionClientV3 struct {
	*s3.Client                         // promoted anonymous field, it allows this type to call s3 Client methods
	Options    EncryptionClientOptions // options for encrypt/decrypt
}

// EncryptionClientOptions is the configuration options for the S3 Encryption Client.
type EncryptionClientOptions struct {
	// TempFolderPath is used to store temp files when calling PutObject
	// Temporary files are needed to compute the X-Amz-Content-Sha256 header
	// temporary file instead of using memory
	MinFileSize int64

	// The logger to write logging messages to
	Logger *log.Logger

	// The CryptographicMaterialsManager to use to manage encryption and decryption materials
	CryptographicMaterialsManager materials.CryptographicMaterialsManager

	// EnableLegacyUnauthenticatedModes MUST be set to true in order to decrypt objects encrypted
	//using legacy (unauthenticated) modes such as AES/CBC
	EnableLegacyUnauthenticatedModes bool
}

// New creates a new S3 Encryption Client v3 with the given CryptographicMaterialsManager
func New(s3Client *s3.Client, CryptographicMaterialsManager materials.CryptographicMaterialsManager, optFns ...func(options *EncryptionClientOptions)) (*S3EncryptionClientV3, error) {
	wrappedClient := s3Client
	// default options
	options := EncryptionClientOptions{
		MinFileSize:                      DefaultMinFileSize,
		Logger:                           log.Default(),
		CryptographicMaterialsManager:    CryptographicMaterialsManager,
		EnableLegacyUnauthenticatedModes: false,
	}
	for _, fn := range optFns {
		fn(&options)
	}

	// use the given wrappedClient for the promoted anon fields
	s3ec := &S3EncryptionClientV3{wrappedClient, options}
	return s3ec, nil
}

// GetObject will make a request to s3 and retrieve the object. In this process
// decryption will be done. The SDK only supports region reads of KMS and GCM.
func (c *S3EncryptionClientV3) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	m := &decryptMiddleware{
		client: c,
		input:  input,
	return c.Client.GetObject(ctx, input, opts...)
}

// PutObject will make encrypt the contents before sending the data to S3. Depending on the MinFileSize
// a temporary file may be used to buffer the encrypted contents to.
func (c *S3EncryptionClientV3) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	em := &encryptMiddleware{
		ec: c,
	}
	opts := append(optFns, encryptOpts...)
	return c.Client.PutObject(ctx, input, opts...)
}

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"io"
)

type aesGCM struct {
	aead  cipher.AEAD
	nonce []byte
}

// newAESGCM creates a new AES GCM cipher. Expects keys to be of
		return nil, err
	}

	return &aesGCM{aesgcm, materials.IV}, nil
}

// Encrypt will encrypt the data using AES GCM
		encrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
	}
	return reader
}
	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
}

func (reader *gcmEncryptReader) Read(data []byte) (int, error) {
		if err != nil {
			return 0, err
		}
		b = reader.encrypter.Seal(b[:0], reader.nonce, b, nil)
		reader.buf = bytes.NewBuffer(b)
	}

		decrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
	}
}

	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
}

func (reader *gcmDecryptReader) Read(data []byte) (int, error) {
	if reader.buf == nil {
		b, err := io.ReadAll(reader.src)
		if err != nil {
			return 0, err
		}
		b, err = reader.decrypter.Open(b[:0], reader.nonce, b, nil)
		if err != nil {
			return 0, err
		}
package internal

import (
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"io"
)

const (
	GcmTagSizeBits  = "128"
	AESGCMNoPadding = "AES/GCM/NoPadding"
)

// NewAESGCMContentCipher returns a new encryption only AES/GCM mode structure with a specific cipher data generator
// will be fully loaded into memory before encryption or decryption can occur. Caution must be taken to avoid memory
// allocation failures.
func NewAESGCMContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.CEKAlgorithm = AESGCMNoPadding
	materials.TagLength = GcmTagSizeBits

	cipher, err := newAESGCM(materials)
	"encoding/hex"
	"encoding/json"
	"fmt"
	materials2 "github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"io"
	"os"
	}
}

func aesgcmTest(t *testing.T, iv, key, plaintext, expected, tag []byte) {
	t.Helper()
	const gcmTagSize = 16
	materials := materials2.CryptographicMaterials{
		Key: key,
		IV:  iv,
	}
package internal

import (
	"errors"
	"io"
	"os"
	ws.i = abs
	return abs, nil
}
	"encoding/json"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"strconv"
)

const DefaultInstructionKeySuffix = ".instruction"

const (
	metaHeader                     = "x-amz-meta"
	keyV1Header                    = "x-amz-key"
	keyV2Header                    = keyV1Header + "-v2"
	ivHeader                       = "x-amz-iv"
	matDescHeader                  = "x-amz-matdesc"
	CekAlgorithmHeader             = "x-amz-cek-alg"
	KeyringAlgorithmHeader         = "x-amz-wrap-alg"
	tagLengthHeader                = "x-amz-tag-len"
	unencryptedContentLengthHeader = "x-amz-unencrypted-content-length"
)

// ObjectMetadata encryption starts off by generating a random symmetric key using
// AES GCM. The SDK generates a random IV based off the encryption cipher
// chosen. The master key that was provided, whether by the user or KMS, will be used
// to encrypt the randomly generated symmetric key and base64 encode the iv. This will
// allow for decryption of that same data later.
type ObjectMetadata struct {
	// IV is the randomly generated IV base64 encoded.
	IV string `json:"x-amz-iv"`
	// CipherKey is the randomly generated cipher key.
	CipherKey string `json:"x-amz-key-v2"`
	// MaterialDesc is a description to distinguish from other envelopes.
	MatDesc               string `json:"x-amz-matdesc"`
	KeyringAlg            string `json:"x-amz-wrap-alg"`
	CEKAlg                string `json:"x-amz-cek-alg"`
	TagLen                string `json:"x-amz-tag-len"`
	UnencryptedContentLen string `json:"x-amz-unencrypted-content-length"`
}

func (e *ObjectMetadata) GetDecodedKey() ([]byte, error) {
	key, err := base64.StdEncoding.DecodeString(e.CipherKey)
	if err != nil {
		return nil, err
	}
	return key, err
}

func (e *ObjectMetadata) GetDecodedIV() ([]byte, error) {
	iv, err := base64.StdEncoding.DecodeString(e.IV)
	if err != nil {
		return nil, err
	}
	return iv, err
}

func (e *ObjectMetadata) GetMatDesc() (string, error) {
	return e.MatDesc, nil
}

// UnmarshalJSON unmarshalls the given JSON bytes into ObjectMetadata
func (e *ObjectMetadata) UnmarshalJSON(value []byte) error {
	type StrictEnvelope ObjectMetadata
}

func EncodeMeta(reader lengthReader, cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	iv := base64.StdEncoding.EncodeToString(cryptographicMaterials.IV)
	key := base64.StdEncoding.EncodeToString(cryptographicMaterials.EncryptedKey)

		UnencryptedContentLen: strconv.FormatInt(contentLength, 10),
	}, nil
}
	"encoding/json"
	"reflect"
	"testing"
)

func TestEnvelope_UnmarshalJSON(t *testing.T) {
		})
	}
}
		input.Metadata = map[string]string{}
	}

	env := saveReq.Envelope
	input.Metadata[http.CanonicalHeaderKey(keyV2Header)] = env.CipherKey
	input.Metadata[http.CanonicalHeaderKey(ivHeader)] = env.IV
	input.Metadata[http.CanonicalHeaderKey(matDescHeader)] = env.MatDesc
	input.Metadata[http.CanonicalHeaderKey(KeyringAlgorithmHeader)] = env.KeyringAlg
	input.Metadata[http.CanonicalHeaderKey(CekAlgorithmHeader)] = env.CEKAlg
	input.Metadata[http.CanonicalHeaderKey(unencryptedContentLengthHeader)] = env.UnencryptedContentLen

	if len(env.TagLen) > 0 {
		input.Metadata[http.CanonicalHeaderKey(tagLengthHeader)] = env.TagLen
	}
	return nil
}

// LoadStrategyRequest represents a request sent to a LoadStrategy to load the contents of an ObjectMetadata
	return env, nil
}

// DefaultLoadStrategy This is the only exported LoadStrategy since cx are no longer able to configure their client
// with a specific load strategy. Instead, we figure out which strategy to use based on the response header on decrypt.
type DefaultLoadStrategy struct {
}

func (load DefaultLoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	if value := req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, keyV2Header}, "-")); value != "" {
		strat := headerV2LoadStrategy{}
		return strat.Load(ctx, req)
	} else if value = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, keyV1Header}, "-")); value != "" {
		// In other S3EC implementations, decryption of v1 objects is supported.
		// Go, however, does not support this.
		return ObjectMetadata{}, &smithy.GenericAPIError{
			Code:    "V1NotSupportedError",
			Message: "The AWS SDK for Go does not support version 1",
		}
	}

	var client GetObjectAPIClient
	if load.client == nil {
		cfg, err := config.LoadDefaultConfig(context.Background())
		if err != nil {
			return ObjectMetadata{}, fmt.Errorf("unable to create S3 client to load instruction file: ")
		}
		client = s3.NewFromConfig(cfg)
	} else {
		client = load.client
	}

	strat := s3LoadStrategy{
		APIClient:             client,
		InstructionFileSuffix: load.suffix,
	}
	return strat.Load(ctx, req)
}
// When EnableLegacyWrappingAlgorithms is set to true, the Keyring MAY decrypt objects encrypted
// using legacy wrapping algorithms such as KMS v1.
type KeyringOptions struct {
	EnableLegacyWrappingAlgorithms bool
}

// object. The KmsKeyring will always use the kmsKeyId provided to encrypt and decrypt messages.
func NewKmsKeyring(apiClient KmsAPIClient, kmsKeyId string, optFns ...func(options *KeyringOptions)) *KmsKeyring {
	options := KeyringOptions{
		EnableLegacyWrappingAlgorithms: false,
	}
	for _, fn := range optFns {
// for use with content decryption, or an error if the object cannot be decrypted
// by the Keyring as its configured.
func (k *KmsKeyring) OnDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey) (*CryptographicMaterials, error) {
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && !k.legacyWrappingAlgorithms {
		return nil, fmt.Errorf("to decrypt x-amz-cek-alg value `%s` you must enable legacyWrappingAlgorithms on the keyring", materials.DataKey.DataKeyAlgorithm)
	}

	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && k.legacyWrappingAlgorithms {
		return commonDecrypt(ctx, materials, encryptedDataKey, &k.KmsKeyId, nil, k.kmsClient)
	} else if materials.DataKey.DataKeyAlgorithm == KMSContextKeyring {
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	encryptionMaterials := NewEncryptionMaterials()

	_, err := keyring.OnEncrypt(ctx, encryptionMaterials)
		GrantTokens: grantTokens,
		KeySpec:     types.DataKeySpecAes256,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: kmsDefaultEncryptionContextKey,
		},
	}


	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	decryptionMaterials, err := NewDecryptionMaterials(DecryptMaterialsRequest{
		CipherKey:  []byte("test-cipher-key"),
		Iv:         []byte("test-iv"),
		MatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
		KeyringAlg: "kms+context",
		CekAlg:     kmsDefaultEncryptionContextKey,
	})
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
		GrantTokens:    grantTokens,
		CiphertextBlob: dataKey.EncryptedDataKey,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: kmsDefaultEncryptionContextKey,
		},
	}

	MaterialDescription MaterialDescription
	ContentAlgorithm    string
	TagLength           string
}

func NewDecryptionMaterials(req DecryptMaterialsRequest) (*DecryptionMaterials, error) {
	MaterialDescription MaterialDescription
	// EncryptedKey should be populated when calling GenerateCipherData
	EncryptedKey []byte
}
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
// This is meant to be a utility function, not a test function,
// but for simplicity and easy invocation it is a test function.
// To avoid running it each test run, it is left commented out.
//func TestGenerateCBCIntegTests(t *testing.T) {
//	arn := "arn:aws:kms:us-west-2:370957321024:alias/S3EC-Go-Github-KMS-Key"
//	bucket := "s3ec-go-github-test-bucket"
//	region := "us-west-2"
//	ctx := context.Background()
//	cfg, _ := config.LoadDefaultConfig(ctx,
//		config.WithRegion(region),
//	)
//
//	s3Client := s3.NewFromConfig(cfg)
//	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
//	// V2 client
//	var handler s3cryptoV2.CipherDataGenerator
//	sessKms, _ := sessionV1.NewSession(&awsV1.Config{
//		Region: aws.String(region),
//	})
//
//	// KMS v1
//	kmsSvc := kmsV1.New(sessKms)
//	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, arn)
//	// AES-CBC content cipher
//	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
//	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)
//
//	for caseKey, plaintext := range fixtures.Plaintexts {
//		_, err := encClient.PutObject(&s3V1.PutObjectInput{
//			Bucket: aws.String(bucket),
//			Key: aws.String(
//				fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
//					fixtures.BaseFolder, version, caseKey),
//			),
//			Body: bytes.NewReader(plaintext),
//		})
//		if err != nil {
//			t.Fatalf("failed to upload encrypted fixture, %v", err)
//		}
//	}
//
//}

func TestKmsV1toV3_CBC(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/V1toV3_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	result, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsV1toV3_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/V1toV3_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	result, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsContextV2toV3_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/V2toV3_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm)

	result, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsContextV3toV2_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/V3toV2_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm)

	_, err = s3ecV3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	}
}

func TestInstructionFileV2toV3(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/inst_file_test.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	result, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v3/language_Go/NegativeV1toV3_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	_, err = s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlgCbc := "aes_cbc"
	keyCbc := "crypto_tests/" + cekAlgCbc + "/v3/language_Go/BothFormats_CBC.txt"
	cekAlgGcm := "aes_gcm"
	keyGcm := "crypto_tests/" + cekAlgGcm + "/v3/language_Go/BothFormats_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	_, err = s3ecV3.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
		Body:   bytes.NewReader([]byte(plaintext)),
		t.Fatalf("error while calling PutObject: %v", err)
	}

	getResponseCbc, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyCbc),
	})
		t.Errorf("expect %v text, got %v", e, a)
	}

	getResponseGcm, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
	})
	}
}

func TestUnicodeEncryptionContextV3(t *testing.T) {
	rune128 := string(rune(128))
	rune200 := string(rune(200))
	rune256 := string(rune(256))

	unicodeStrings := []string{rune128, rune200, rune256, runeMaxInt, shorter, medium, longer, mix, mixTwo}
	for _, s := range unicodeStrings {
		UnicodeEncryptionContextV3(t, s)
	}
}

func UnicodeEncryptionContextV3(t *testing.T, metadataString string) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV3, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": metadataString})
	_, err = s3ecV3.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),

	time.Sleep(1 * time.Second)

	result, err := s3ecV3.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
		t.Errorf("expect %v text, got %v", e, a)
	}

	s3ecV3.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	"bytes"
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"strings"
	"testing"
	return fmt.Sprintf(arnFormat, region, accountId, shortAlias)
}

func TestInteg_EncryptFixtures(t *testing.T) {
	var region = LoadRegion()
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		KEK, bucket, region, CEK string
	}{
		{
			CEKAlg: "aes_gcm",
			KEK:    "kms", bucket: bucket, region: region, CEK: "aes_gcm",
		},
			if err != nil {
				t.Fatalf("failed to create new CMM")
			}
			encClient, _ := client.New(s3Client, cmm)

			for caseKey, plaintext := range fixtures.Plaintexts {
				_, err := encClient.PutObject(ctx, &s3.PutObjectInput{
					),
					Body: bytes.NewReader(plaintext),
				})
				if err != nil {
					t.Fatalf("failed to upload encrypted fixture, %v", err)
				}
		Lang    string
		Version string
	}{
		{CEKAlg: "aes_cbc", Lang: "Go", Version: "v3"}, // v3 doesn't support CBC but that's where the files are
		{CEKAlg: "aes_gcm", Lang: "Go", Version: "v3"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v2"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v3"},
	}
				cmmCbc, err := materials.NewCryptographicMaterialsManager(keyring)
				decClient, err = client.New(s3Client, cmmCbc, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = true
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			} else if c.CEKAlg == "aes_gcm" {
				decClient, err = client.New(s3Client, cmm)
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			}

			fixtures := getFixtures(t, s3Client, c.CEKAlg, bucket)
			ciphertexts := decryptFixtures(t, decClient, fixtures, bucket, c.Lang, version)

			if len(ciphertexts) == 0 {
				t.Fatalf("expected more than 0 ciphertexts to decrypt!")
	}

	switch cek {
	case "aes_gcm":
		return kmsKeyring
	case "aes_cbc":
	})
}

func TestIntegKmsContext(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
		Key:    &key,
	})
}
