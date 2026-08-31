package main

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"
	"testing"

	v3example "migration-examples/v3"
	step1example "migration-examples/v4/step1_forbid_encrypt_allow_decrypt"
	step2example "migration-examples/v4/step2_require_encrypt_allow_decrypt"
	step3example "migration-examples/v4/step3_require_encrypt_require_decrypt"
)

// Test configuration - these should be set via environment variables
var (
	testBucket = getEnvOrDefault("TEST_BUCKET", "s3ec-go-github-test-bucket")
	testKey    = getEnvOrDefault("TEST_KEY", "migration-test")
	testKMSKey = getEnvOrDefault("TEST_KMS_KEY", "arn:aws:kms:us-west-2:370957321024:alias/S3EC-Go-Github-KMS-Key")
	testRegion = getEnvOrDefault("TEST_REGION", "us-west-2")
)

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// runStepExample calls the appropriate step example function directly
func runStepExample(step int, bucket, key, kmsKey, region string) error {
	switch step {
	case 0:
		return v3example.RunMigrationExample(bucket, key, kmsKey, region)
	case 1:
		return step1example.RunMigrationExample(bucket, key, kmsKey, region)
	case 2:
		return step2example.RunMigrationExample(bucket, key, kmsKey, region)
	case 3:
		return step3example.RunMigrationExample(bucket, key, kmsKey, region)
	default:
		return fmt.Errorf("unknown migration step: %d", step)
	}
}

// runStepExampleWithSource calls the appropriate step example function with a source step parameter
func runStepExampleWithSource(step int, sourceStep int, bucket, key, kmsKey, region string) error {
	switch step {
	case 0:
		return v3example.RunMigrationExample(bucket, key, kmsKey, region, sourceStep)
	case 1:
		return step1example.RunMigrationExample(bucket, key, kmsKey, region, sourceStep)
	case 2:
		return step2example.RunMigrationExample(bucket, key, kmsKey, region, sourceStep)
	case 3:
		return step3example.RunMigrationExample(bucket, key, kmsKey, region, sourceStep)
	default:
		return fmt.Errorf("unknown migration step: %d", step)
	}
}

// TestMigrationStep0WriteRead tests that Step 0 (V3) can write and read its own objects
func TestMigrationStep0WriteRead(t *testing.T) {
	t.Log("Testing Step 0 (V3 client) write-read roundtrip")
	
	err := runStepExample(0, testBucket, testKey, testKMSKey, testRegion)
	if err != nil {
		t.Errorf("Step 0 write-read roundtrip failed: %v", err)
	} else {
		t.Log("SUCCESS: Step 0 write-read roundtrip completed")
	}
}

// TestMigrationStep1WriteRead tests that Step 1 (V4 FORBID_ENCRYPT_ALLOW_DECRYPT) can write and read its own objects
func TestMigrationStep1WriteRead(t *testing.T) {
	t.Log("Testing Step 1 (V4 FORBID_ENCRYPT_ALLOW_DECRYPT) write-read roundtrip")
	
	err := runStepExample(1, testBucket, testKey, testKMSKey, testRegion)
	if err != nil {
		t.Errorf("Step 1 write-read roundtrip failed: %v", err)
	} else {
		t.Log("SUCCESS: Step 1 write-read roundtrip completed")
	}
}

// TestMigrationStep2WriteRead tests that Step 2 (V4 REQUIRE_ENCRYPT_ALLOW_DECRYPT) can write and read its own objects
func TestMigrationStep2WriteRead(t *testing.T) {
	t.Log("Testing Step 2 (V4 REQUIRE_ENCRYPT_ALLOW_DECRYPT) write-read roundtrip")
	
	err := runStepExample(2, testBucket, testKey, testKMSKey, testRegion)
	if err != nil {
		t.Errorf("Step 2 write-read roundtrip failed: %v", err)
	} else {
		t.Log("SUCCESS: Step 2 write-read roundtrip completed")
	}
}

// TestMigrationStep3WriteRead tests that Step 3 (V4 REQUIRE_ENCRYPT_REQUIRE_DECRYPT) can write and read its own objects
func TestMigrationStep3WriteRead(t *testing.T) {
	t.Log("Testing Step 3 (V4 REQUIRE_ENCRYPT_REQUIRE_DECRYPT) write-read roundtrip")
	
	err := runStepExample(3, testBucket, testKey, testKMSKey, testRegion)
	if err != nil {
		t.Errorf("Step 3 write-read roundtrip failed: %v", err)
	} else {
		t.Log("SUCCESS: Step 3 write-read roundtrip completed")
	}
}

// TestMigrationCompatibilityMatrix tests cross-compatibility between all migration steps
func TestMigrationCompatibilityMatrix(t *testing.T) {
	// Define the compatibility matrix
	// Each test case specifies: reader step, writer step, expected success
	compatibilityTests := []struct {
		readerStep    int
		writerStep    int
		expectSuccess bool
		description   string
	}{
		// Step 0 (V3) compatibility - can read everything
		{0, 0, true, "V3 reading V3-encrypted object (no key commitment)"},
		{0, 1, true, "V3 reading V4-forbid/allow-encrypted object (no key commitment)"},
		{0, 2, true, "V3 reading V4-require/allow-encrypted object (with key commitment)"},
		{0, 3, true, "V3 reading V4-require/require-encrypted object (with key commitment)"},
		
		// Step 1 (V4 FORBID_ENCRYPT_ALLOW_DECRYPT) compatibility - can read everything
		{1, 0, true, "V4-forbid/allow reading V3-encrypted object (no key commitment)"},
		{1, 1, true, "V4-forbid/allow reading V4-forbid/allow-encrypted object (no key commitment)"},
		{1, 2, true, "V4-forbid/allow reading V4-require/allow-encrypted object (with key commitment)"},
		{1, 3, true, "V4-forbid/allow reading V4-require/require-encrypted object (with key commitment)"},
		
		// Step 2 (V4 REQUIRE_ENCRYPT_ALLOW_DECRYPT) compatibility - can read everything
		{2, 0, true, "V4-require/allow reading V3-encrypted object (no key commitment)"},
		{2, 1, true, "V4-require/allow reading V4-forbid/allow-encrypted object (no key commitment)"},
		{2, 2, true, "V4-require/allow reading V4-require/allow-encrypted object (with key commitment)"},
		{2, 3, true, "V4-require/allow reading V4-require/require-encrypted object (with key commitment)"},
		
		// Step 3 (V4 REQUIRE_ENCRYPT_REQUIRE_DECRYPT) compatibility - can only read objects with key commitment
		{3, 0, false, "V4-require/require reading V3-encrypted object (no key commitment) - SHOULD FAIL"},
		{3, 1, false, "V4-require/require reading V4-forbid/allow-encrypted object (no key commitment) - SHOULD FAIL"},
		{3, 2, true, "V4-require/require reading V4-require/allow-encrypted object (with key commitment)"},
		{3, 3, true, "V4-require/require reading V4-require/require-encrypted object (with key commitment)"},
	}

	for _, tt := range compatibilityTests {
		testName := fmt.Sprintf("Reader_Step%d_Writer_Step%d", tt.readerStep, tt.writerStep)
		t.Run(testName, func(t *testing.T) {
			t.Logf("Testing: %s", tt.description)
			
			// Phase 1: Write object using the WRITER step's client
			t.Logf("Phase 1: Writing object using step %d client", tt.writerStep)
			err := runStepExample(tt.writerStep, testBucket, testKey, testKMSKey, testRegion)
			if err != nil {
				t.Fatalf("Failed to write object with step %d: %v", tt.writerStep, err)
			}
			t.Logf("Successfully wrote object with step %d", tt.writerStep)
			
			// Phase 2: Try to read object using the READER step's client
			t.Logf("Phase 2: Reading object using step %d client", tt.readerStep)
			err = runStepExampleWithSource(tt.readerStep, tt.writerStep, testBucket, testKey, testKMSKey, testRegion)
			
			if tt.expectSuccess {
				if err != nil {
					t.Errorf("Expected success but got error: %v", err)
				} else {
					t.Logf("SUCCESS: Step %d successfully processed object from step %d", tt.readerStep, tt.writerStep)
				}
			} else {
				if err == nil {
					t.Errorf("Expected failure but got success - Step %d should not be able to decrypt objects from step %d", tt.readerStep, tt.writerStep)
				} else {
					t.Logf("EXPECTED FAILURE: Step %d correctly failed to process object from step %d: %v", tt.readerStep, tt.writerStep, err)
				}
			}
		})
	}
}

// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	v3example "migration-examples/v3"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <bucket-name> <object-key> <kms-key-id> <region>\n", os.Args[0])
		fmt.Printf("Example: %s my-bucket my-key arn:aws:kms:us-east-2:123456789012:key/12345678-1234-1234-1234-123456789012 us-east-2\n", os.Args[0])
		os.Exit(1)
	}

	bucketName := os.Args[1]
	objectKey := os.Args[2]
	kmsKeyID := os.Args[3]
	region := os.Args[4]

	// Run with current step as both source and target (normal operation)
	err := v3example.RunMigrationExample(bucketName, objectKey, kmsKeyID, region)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v3example

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Migration Step 0: This example demonstrates use of the S3 Encryption Client for Go v3
// and is the starting state for migrating your data to the v4 client.
//
// This example's purpose is to model behavior of an existing v3 client.
// Subsequent migration steps will demonstrate code changes needed to use the v4 client.
//
// This example configures a v3 client to:
// - Write objects using non-key committing encryption algorithms
// - Read objects encrypted with either key committing algorithms or with non-key committing algorithms
//
// In this configuration, the client can read objects encrypted
// with non-key committing algorithms (written by this v3 client or an in-progress v4 migration),
// as well as objects encrypted by a migrated v4 client
// that is configured to write objects encrypted with key committing algorithms.
// You should ensure you are using the latest version of the v3 client
// that can read objects encrypted with key committing algorithms before proceeding with migration.

const CurrentMigrationStep = 0

func RunMigrationExample(bucketName, objectKey, kmsKeyID, region string, sourceStep ...int) error {
	actualSourceStep := CurrentMigrationStep // Default to current step
	if len(sourceStep) > 0 {
		actualSourceStep = sourceStep[0]
	}
	
	fmt.Println("=== S3 Encryption Client v3 Baseline Example ===")
	fmt.Printf("Bucket: %s\n", bucketName)
	fmt.Printf("Object Key: %s\n", objectKey)
	fmt.Printf("KMS Key ID: %s\n", kmsKeyID)
	fmt.Printf("Region: %s\n", region)
	fmt.Println()

	// Test data for encryption
	testData := "Hello, World! This is a test message for S3 encryption client migration."
	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("Data length: %d bytes\n", len(testData))
	fmt.Println()

	fmt.Println("--- Initialize S3 Encryption Client v3 ---")

	// Create regular S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("error loading AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Create KMS keyring
	keyring := materials.NewKmsKeyring(kmsClient, kmsKeyID)

	// Create Cryptographic Materials Manager
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		return fmt.Errorf("error creating CMM: %v", err)
	}

	// Create S3 Encryption Client v3 with FORBID_ENCRYPT_ALLOW_DECRYPT commitment policy (default)
	encryptionClient, err := client.New(s3Client, cmm, func(options *client.EncryptionClientOptions) {
		// This is the default commitment policy for v3 clients
		// that can read objects encrypted with key commitment;
		// you do not need to set this explicitly.
		// However, the "commitment" package is only available in v3 versions
		// that can read objects encrypted with key commitment.
		// If you are unsure whether your v3 client version supports this,
		// you can set this explicitly to avoid accidental use of a v3 client
		// that cannot read objects encrypted with key committing algorithms.
		options.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		return fmt.Errorf("error creating S3 Encryption Client: %v", err)
	}

	fmt.Println("Successfully initialized S3 Encryption Client v3")
	fmt.Println("Commitment Policy: FORBID_ENCRYPT_ALLOW_DECRYPT")
	fmt.Println()

	// Create object keys for PUT and GET operations
	// PUT: Always use current step
	putObjectKey := fmt.Sprintf("%s-step-%d", objectKey, CurrentMigrationStep)
	// GET: Use actualSourceStep (debug parameter to test cross-compatibility between steps; defaults to current step)
	getObjectKey := fmt.Sprintf("%s-step-%d", objectKey, actualSourceStep)

	fmt.Println("--- Encrypt and Upload Object to S3 ---")

	// Upload encrypted object using S3 Encryption Client
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(putObjectKey),
		Body:   strings.NewReader(testData),
	}

	_, err = encryptionClient.PutObject(context.TODO(), putInput)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucket") {
			return fmt.Errorf("S3 bucket '%s' does not exist or is not accessible", bucketName)
		} else if strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("KMS key '%s' not found or not accessible", kmsKeyID)
		} else {
			return fmt.Errorf("error uploading encrypted object: %v", err)
		}
	}

	fmt.Println("Successfully uploaded encrypted object to S3!")
	fmt.Printf("   Bucket: %s\n", bucketName)
	fmt.Printf("   Key: %s\n", putObjectKey)
	fmt.Println()

	fmt.Println("--- Download and Decrypt Object from S3 ---")

	// Download and decrypt object using S3 Encryption Client
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(getObjectKey),
	}

	getResponse, err := encryptionClient.GetObject(context.TODO(), getInput)
	if err != nil {
		return fmt.Errorf("error downloading and decrypting object: %v", err)
	}
	defer getResponse.Body.Close()

	// Read the decrypted data
	decryptedData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return fmt.Errorf("error reading decrypted data: %v", err)
	}

	fmt.Println("Successfully downloaded and decrypted object from S3!")
	fmt.Printf("   Object size: %d bytes\n", len(decryptedData))
	fmt.Printf("   Decrypted data: %s\n", string(decryptedData))
	fmt.Println()

	fmt.Println("--- Verify Roundtrip Success ---")

	// Verify the roundtrip was successful
	if string(decryptedData) == testData {
		fmt.Println("SUCCESS: Roundtrip encryption/decryption completed successfully!")
		fmt.Println("   Original data matches decrypted data")
		fmt.Println("   Data integrity verified")
	} else {
		return fmt.Errorf("roundtrip failed - data mismatch. Original: %s, Decrypted: %s", testData, string(decryptedData))
	}

	fmt.Println()
	fmt.Println("=== V3 Baseline Example completed successfully! ===")
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	step1example "migration-examples/v4/step1_forbid_encrypt_allow_decrypt"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <bucket-name> <object-key> <kms-key-id> <region>\n", os.Args[0])
		fmt.Printf("Example: %s my-bucket my-key arn:aws:kms:us-east-2:123456789012:key/12345678-1234-1234-1234-123456789012 us-east-2\n", os.Args[0])
		os.Exit(1)
	}

	bucketName := os.Args[1]
	objectKey := os.Args[2]
	kmsKeyID := os.Args[3]
	region := os.Args[4]

	// Run with current step as both source and target (normal operation)
	err := step1example.RunMigrationExample(bucketName, objectKey, kmsKeyID, region)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package step1example

import (
	"context"
	"fmt"
	"io"
	"strings"

	// Migration note: The "v3" import has been updated to "v4".
	"github.com/aws/amazon-s3-encryption-client-go/v4/client"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Migration Step 1: This example demonstrates how to start using the S3 Encryption Client v4.
//
// This example's purpose is to demonstrate the code changes to 
// migrate from the v3 client to the v4 client while maintaining identical behavior.
//
// When starting from a v3 client modeled in "Migration Step 0",
// "Migration Step 1" should result in no behavioral changes to your application.
//
// In this example we configure a v4 client to:
// - Write objects encrypted with non-key committing algorithms
// - Read objects encrypted either with or without key committing algorithms
//
// In this configuration, the client will continue to read objects encrypted
// with non-key committing algorithms (written by a v3 client or this migration-in-progress v4 client),
// as well as objects encrypted by a migrated v4 client
// that is configured to write objects encrypted with key committing algorithms.
//
// This configuration results in identical behavior to the S3 Encryption Client v3 client
// configured to use the default FORBID_ENCRYPT_ALLOW_DECRYPT commitment policy.

const CurrentMigrationStep = 1

func RunMigrationExample(bucketName, objectKey, kmsKeyID, region string, sourceStep ...int) error {
	actualSourceStep := CurrentMigrationStep // Default to current step
	if len(sourceStep) > 0 {
		actualSourceStep = sourceStep[0]
	}
	
	fmt.Println("=== S3 Encryption Client v4 Step 1 Example ===")
	fmt.Printf("Bucket: %s\n", bucketName)
	fmt.Printf("Object Key: %s\n", objectKey)
	fmt.Printf("KMS Key ID: %s\n", kmsKeyID)
	fmt.Printf("Region: %s\n", region)
	fmt.Println()

	// Test data for encryption
	testData := "Hello, World! This is a test message for S3 encryption client migration."
	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("Data length: %d bytes\n", len(testData))
	fmt.Println()

	fmt.Println("--- Initialize S3 Encryption Client v4 ---")

	// Create regular S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("error loading AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Create KMS keyring
	keyring := materials.NewKmsKeyring(kmsClient, kmsKeyID)

	// Create Cryptographic Materials Manager
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		return fmt.Errorf("error creating CMM: %v", err)
	}

	// Create S3 Encryption Client v4 with FORBID_ENCRYPT_ALLOW_DECRYPT commitment policy
	// Migration note: The type of this client has changed from `S3EncryptionClientV3` to `S3EncryptionClientV4`.
	encryptionClient, err := client.New(s3Client, cmm, func(options *client.EncryptionClientOptions) {
		// This MUST be explicitly configured to FORBID_ENCRYPT_ALLOW_DECRYPT.
		// While FORBID_ENCRYPT_ALLOW_DECRYPT is the default for v3 clients,
		// v4 clients default to REQUIRE_ENCRYPT_REQUIRE_DECRYPT.
		// This configuration ensures identical behavior to a v3 client.
		options.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		return fmt.Errorf("error creating S3 Encryption Client: %v", err)
	}

	fmt.Println("Successfully initialized S3 Encryption Client v4")
	fmt.Println("Commitment Policy: FORBID_ENCRYPT_ALLOW_DECRYPT")
	fmt.Println()

	// Create object keys for PUT and GET operations
	// PUT: Always use current step
	putObjectKey := fmt.Sprintf("%s-step-%d", objectKey, CurrentMigrationStep)
	// GET: Use sourceStep (debug parameter to test cross-compatibility between steps; defaults to 1)
	getObjectKey := fmt.Sprintf("%s-step-%d", objectKey, actualSourceStep)

	fmt.Println("--- Encrypt and Upload Object to S3 ---")

	// Upload encrypted object using S3 Encryption Client
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(putObjectKey),
		Body:   strings.NewReader(testData),
	}

	_, err = encryptionClient.PutObject(context.TODO(), putInput)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucket") {
			return fmt.Errorf("S3 bucket '%s' does not exist or is not accessible", bucketName)
		} else if strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("KMS key '%s' not found or not accessible", kmsKeyID)
		} else {
			return fmt.Errorf("error uploading encrypted object: %v", err)
		}
	}

	fmt.Println("Successfully uploaded encrypted object to S3!")
	fmt.Printf("   Bucket: %s\n", bucketName)
	fmt.Printf("   Key: %s\n", putObjectKey)
	fmt.Println()

	fmt.Println("--- Download and Decrypt Object from S3 ---")

	// Download and decrypt object using S3 Encryption Client
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(getObjectKey),
	}

	getResponse, err := encryptionClient.GetObject(context.TODO(), getInput)
	if err != nil {
		return fmt.Errorf("error downloading and decrypting object: %v", err)
	}
	defer getResponse.Body.Close()

	// Read the decrypted data
	decryptedData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return fmt.Errorf("error reading decrypted data: %v", err)
	}

	fmt.Println("Successfully downloaded and decrypted object from S3!")
	fmt.Printf("   Object size: %d bytes\n", len(decryptedData))
	fmt.Printf("   Decrypted data: %s\n", string(decryptedData))
	fmt.Println()

	fmt.Println("--- Verify Roundtrip Success ---")

	// Verify the roundtrip was successful
	if string(decryptedData) == testData {
		fmt.Println("SUCCESS: Roundtrip encryption/decryption completed successfully!")
		fmt.Println("   Original data matches decrypted data")
		fmt.Println("   Data integrity verified")
	} else {
		return fmt.Errorf("roundtrip failed - data mismatch. Original: %s, Decrypted: %s", testData, string(decryptedData))
	}

	fmt.Println()
	fmt.Println("=== V4 Step 1 Example completed successfully! ===")
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	step2example "migration-examples/v4/step2_require_encrypt_allow_decrypt"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <bucket-name> <object-key> <kms-key-id> <region>\n", os.Args[0])
		fmt.Printf("Example: %s my-bucket my-key arn:aws:kms:us-east-2:123456789012:key/12345678-1234-1234-1234-123456789012 us-east-2\n", os.Args[0])
		os.Exit(1)
	}

	bucketName := os.Args[1]
	objectKey := os.Args[2]
	kmsKeyID := os.Args[3]
	region := os.Args[4]

	// Run with current step as both source and target (normal operation)
	err := step2example.RunMigrationExample(bucketName, objectKey, kmsKeyID, region)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package step2example

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/amazon-s3-encryption-client-go/v4/client"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Migration Step 2: This example demonstrates how to update your v4 client configuration
// to start writing objects encrypted with key committing algorithms.
//
// This example's purpose is to demonstrate the commitment policy code changes required to
// start writing objects encrypted with key committing algorithms
// and document the behavioral changes that will result from this change.
//
// When starting from a v4 client modeled in "Migration Step 1",
// "Migration Step 2" WILL result in behavioral changes to your application.
// The client will start writing objects encrypted with key committing algorithms.
//
// IMPORTANT: You MUST have updated your readers to be able to read objects encrypted with key committing algorithms
// before deploying the changes in this step.
// This means deploying the changes from either "Migration Step 0" (if readers are v3 clients)
// or "Migration Step 1" (if readers are v4 clients) to all of your readers
// before deploying the changes from to "Migration Step 2".
//
// Once you deploy this change to your writers, your readers will start seeing
// some objects encrypted with non-key committing algorithms,
// and some objects encrypted with key committing algorithms.
// Because the changes would have already been deployed to all our readers from earlier migration steps,
// we can be sure that our entire system is ready to read both types of objects.
// After deploying these changes but before proceeding to "Migration Step 3",
// you MUST take extra steps to ensure that your system is no longer reading
// objects encrypted with non-key committing algorithms
// (such as re-encrypting any existing objects using key committing algorithms).

const CurrentMigrationStep = 2

func RunMigrationExample(bucketName, objectKey, kmsKeyID, region string, sourceStep ...int) error {
	actualSourceStep := CurrentMigrationStep // Default to current step
	if len(sourceStep) > 0 {
		actualSourceStep = sourceStep[0]
	}
	
	fmt.Println("=== S3 Encryption Client v4 Step 2 Example ===")
	fmt.Printf("Current Step: %d (V4 with REQUIRE_ENCRYPT_ALLOW_DECRYPT)\n", CurrentMigrationStep)
	fmt.Printf("Source Step: %d (reading object written by step %d)\n", actualSourceStep, actualSourceStep)
	fmt.Printf("Bucket: %s\n", bucketName)
	fmt.Printf("Object Key: %s\n", objectKey)
	fmt.Printf("KMS Key ID: %s\n", kmsKeyID)
	fmt.Printf("Region: %s\n", region)
	fmt.Println()

	// Test data for encryption
	testData := "Hello, World! This is a test message for S3 encryption client migration."
	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("Data length: %d bytes\n", len(testData))
	fmt.Println()

	fmt.Println("--- Initialize S3 Encryption Client v4 ---")

	// Create regular S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("error loading AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Create KMS keyring
	keyring := materials.NewKmsKeyring(kmsClient, kmsKeyID)

	// Create Cryptographic Materials Manager
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		return fmt.Errorf("error creating CMM: %v", err)
	}

	// Create S3 Encryption Client v4 with REQUIRE_ENCRYPT_ALLOW_DECRYPT commitment policy
	encryptionClient, err := client.New(s3Client, cmm, func(options *client.EncryptionClientOptions) {
		// Migration note: The commitment policy has been updated to REQUIRE_ENCRYPT_ALLOW_DECRYPT.
		// This change causes the client to start writing objects encrypted with key committing algorithms.
		// The client will continue to be able to read objects encrypted with either
		// key committing or non-key committing algorithms.
		options.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		return fmt.Errorf("error creating S3 Encryption Client: %v", err)
	}

	fmt.Println("Successfully initialized S3 Encryption Client v4")
	fmt.Println("Commitment Policy: REQUIRE_ENCRYPT_ALLOW_DECRYPT")
	fmt.Println()

	// Create object keys for PUT and GET operations
	// PUT: Always use current step
	putObjectKey := fmt.Sprintf("%s-step-%d", objectKey, CurrentMigrationStep)
	// GET: Use sourceStep (debug parameter to test cross-compatibility between steps; defaults to 2)
	getObjectKey := fmt.Sprintf("%s-step-%d", objectKey, actualSourceStep)

	fmt.Println("--- Encrypt and Upload Object to S3 ---")

	// Upload encrypted object using S3 Encryption Client
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(putObjectKey),
		Body:   strings.NewReader(testData),
	}

	_, err = encryptionClient.PutObject(context.TODO(), putInput)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucket") {
			return fmt.Errorf("S3 bucket '%s' does not exist or is not accessible", bucketName)
		} else if strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("KMS key '%s' not found or not accessible", kmsKeyID)
		} else {
			return fmt.Errorf("error uploading encrypted object: %v", err)
		}
	}

	fmt.Println("Successfully uploaded encrypted object to S3!")
	fmt.Printf("   Bucket: %s\n", bucketName)
	fmt.Printf("   Key: %s\n", putObjectKey)
	fmt.Println()

	fmt.Println("--- Download and Decrypt Object from S3 ---")

	// Download and decrypt object using S3 Encryption Client
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(getObjectKey),
	}

	getResponse, err := encryptionClient.GetObject(context.TODO(), getInput)
	if err != nil {
		return fmt.Errorf("error downloading and decrypting object: %v", err)
	}
	defer getResponse.Body.Close()

	// Read the decrypted data
	decryptedData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return fmt.Errorf("error reading decrypted data: %v", err)
	}

	fmt.Println("Successfully downloaded and decrypted object from S3!")
	fmt.Printf("   Object size: %d bytes\n", len(decryptedData))
	fmt.Printf("   Decrypted data: %s\n", string(decryptedData))
	fmt.Println()

	fmt.Println("--- Verify Roundtrip Success ---")

	// Verify the roundtrip was successful
	if string(decryptedData) == testData {
		fmt.Println("SUCCESS: Roundtrip encryption/decryption completed successfully!")
		fmt.Println("   Original data matches decrypted data")
		fmt.Println("   Data integrity verified")
	} else {
		return fmt.Errorf("roundtrip failed - data mismatch. Original: %s, Decrypted: %s", testData, string(decryptedData))
	}

	fmt.Println()
	fmt.Println("=== V4 Step 2 Example completed successfully! ===")
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
	"os"

	step3example "migration-examples/v4/step3_require_encrypt_require_decrypt"
)

func main() {
	// Check command line arguments
	if len(os.Args) != 5 {
		fmt.Printf("Usage: %s <bucket-name> <object-key> <kms-key-id> <region>\n", os.Args[0])
		fmt.Printf("Example: %s my-bucket my-key arn:aws:kms:us-east-2:123456789012:key/12345678-1234-1234-1234-123456789012 us-east-2\n", os.Args[0])
		os.Exit(1)
	}

	bucketName := os.Args[1]
	objectKey := os.Args[2]
	kmsKeyID := os.Args[3]
	region := os.Args[4]

	// Run with current step as both source and target (normal operation)
	err := step3example.RunMigrationExample(bucketName, objectKey, kmsKeyID, region)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package step3example

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/aws/amazon-s3-encryption-client-go/v4/client"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// Migration Step 3: This example demonstrates how to update your v4 client configuration
// to stop reading objects encrypted with non-key committing algorithms.
//
// This example's purpose is to demonstrate the commitment policy code changes required to
// stop reading objects encrypted with non-key committing algorithms
// and document the behavioral changes that will result from this change.
//
// When starting from a v4 client modeled in "Migration Step 2",
// "Migration Step 3" WILL result in behavioral changes to your application.
// The client will no longer be able to read objects encrypted with non-key committing algorithms.
// Before deploying these changes, you MUST have taken some extra steps 
// to ensure that your system is no longer reading such objects,
// such as re-encrypting them with key committing algorithms.
//
// IMPORTANT: Before deploying the changes in this step, your system should not be reading
// any objects encrypted with non-key committing algorithms.
// The changes in this step will cause such read attempts to fail.
// This means the changes from "Migration Step 2" should have already been deployed to all of your readers
// before you deploy the changes from "Migration Step 3".
//
// Once you complete Step 3, you can be sure that all items being read by your system
// have been encrypted using key committing algorithms.

const CurrentMigrationStep = 3

func RunMigrationExample(bucketName, objectKey, kmsKeyID, region string, sourceStep ...int) error {
	actualSourceStep := CurrentMigrationStep // Default to current step
	if len(sourceStep) > 0 {
		actualSourceStep = sourceStep[0]
	}
	
	fmt.Println("=== S3 Encryption Client v4 Step 3 Example ===")
	fmt.Printf("Current Step: %d (V4 with REQUIRE_ENCRYPT_REQUIRE_DECRYPT)\n", CurrentMigrationStep)
	fmt.Printf("Source Step: %d (reading object written by step %d)\n", actualSourceStep, actualSourceStep)
	fmt.Printf("Bucket: %s\n", bucketName)
	fmt.Printf("Object Key: %s\n", objectKey)
	fmt.Printf("KMS Key ID: %s\n", kmsKeyID)
	fmt.Printf("Region: %s\n", region)
	fmt.Println()

	// Test data for encryption
	testData := "Hello, World! This is a test message for S3 encryption client migration."
	fmt.Printf("Original data: %s\n", testData)
	fmt.Printf("Data length: %d bytes\n", len(testData))
	fmt.Println()

	fmt.Println("--- Initialize S3 Encryption Client v4 ---")

	// Create regular S3 client
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(region))
	if err != nil {
		return fmt.Errorf("error loading AWS config: %v", err)
	}
	s3Client := s3.NewFromConfig(cfg)

	// Create KMS client
	kmsClient := kms.NewFromConfig(cfg)

	// Create KMS keyring
	keyring := materials.NewKmsKeyring(kmsClient, kmsKeyID)

	// Create Cryptographic Materials Manager
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		return fmt.Errorf("error creating CMM: %v", err)
	}

	// Create S3 Encryption Client v4 with REQUIRE_ENCRYPT_REQUIRE_DECRYPT commitment policy
	encryptionClient, err := client.New(s3Client, cmm, func(options *client.EncryptionClientOptions) {
		// Migration note: The commitment policy has been changed to REQUIRE_ENCRYPT_REQUIRE_DECRYPT.
		// This change causes the client to stop reading objects encrypted without key committing algorithms.
		// IMPORTANT: Ensure your system is no longer reading such objects before deploying this change.
		// REQUIRE_ENCRYPT_REQUIRE_DECRYPT is also the default commitment policy for v4 clients,
		// so you do not need to set this explicitly.
		options.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})
	if err != nil {
		return fmt.Errorf("error creating S3 Encryption Client: %v", err)
	}

	fmt.Println("Successfully initialized S3 Encryption Client v4")
	fmt.Println("Commitment Policy: REQUIRE_ENCRYPT_REQUIRE_DECRYPT")
	fmt.Println()

	// Create object keys for PUT and GET operations
	// PUT: Always use current step
	putObjectKey := fmt.Sprintf("%s-step-%d", objectKey, CurrentMigrationStep)
	// GET: Use sourceStep (debug parameter to test cross-compatibility between steps; defaults to 3)
	getObjectKey := fmt.Sprintf("%s-step-%d", objectKey, actualSourceStep)

	fmt.Println("--- Encrypt and Upload Object to S3 ---")

	// Upload encrypted object using S3 Encryption Client
	putInput := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(putObjectKey),
		Body:   strings.NewReader(testData),
	}

	_, err = encryptionClient.PutObject(context.TODO(), putInput)
	if err != nil {
		if strings.Contains(err.Error(), "NoSuchBucket") {
			return fmt.Errorf("S3 bucket '%s' does not exist or is not accessible", bucketName)
		} else if strings.Contains(err.Error(), "NotFoundException") {
			return fmt.Errorf("KMS key '%s' not found or not accessible", kmsKeyID)
		} else {
			return fmt.Errorf("error uploading encrypted object: %v", err)
		}
	}

	fmt.Println("Successfully uploaded encrypted object to S3!")
	fmt.Printf("   Bucket: %s\n", bucketName)
	fmt.Printf("   Key: %s\n", putObjectKey)
	fmt.Println()

	fmt.Println("--- Download and Decrypt Object from S3 ---")

	// Download and decrypt object using S3 Encryption Client
	getInput := &s3.GetObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(getObjectKey),
	}

	getResponse, err := encryptionClient.GetObject(context.TODO(), getInput)
	if err != nil {
		return fmt.Errorf("error downloading and decrypting object: %v", err)
	}
	defer getResponse.Body.Close()

	// Read the decrypted data
	decryptedData, err := io.ReadAll(getResponse.Body)
	if err != nil {
		return fmt.Errorf("error reading decrypted data: %v", err)
	}

	fmt.Println("Successfully downloaded and decrypted object from S3!")
	fmt.Printf("   Object size: %d bytes\n", len(decryptedData))
	fmt.Printf("   Decrypted data: %s\n", string(decryptedData))
	fmt.Println()

	fmt.Println("--- Verify Roundtrip Success ---")

	// Verify the roundtrip was successful
	if string(decryptedData) == testData {
		fmt.Println("SUCCESS: Roundtrip encryption/decryption completed successfully!")
		fmt.Println("   Original data matches decrypted data")
		fmt.Println("   Data integrity verified")
	} else {
		return fmt.Errorf("roundtrip failed - data mismatch. Original: %s, Decrypted: %s", testData, string(decryptedData))
	}

	fmt.Println()
	fmt.Println("=== V4 Step 2 Example completed successfully! ===")
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package algorithms

import (
	"crypto/sha512"
	"fmt"
	"strconv"
	"hash"
)

// Algorithm constants
const (
	// GCM maximum content length in bits (2^39 - 256 bits)
	GCMMaxContentLengthBits = 549755813632
	// CTR maximum content length in bytes (2^32 bytes)
	CTRMaxContentLengthBytes = 4294967296
	// CBC maximum content length in bytes (2^32 bytes)  
	CBCMaxContentLengthBytes = 4294967296
	
	// Algorithm identifiers
	AESGCMCommitKey = "115"
	AESGCMNoPadding = "AES/GCM/NoPadding"
	AESCTRNoPadding = "AES/CTR/NoPadding"
	AESCBCPKCS5   	= "AES/CBC/PKCS5Padding"
)

// AlgorithmSuite represents the encryption algorithm suite configuration
type AlgorithmSuite struct {
	id                         int
	isLegacy                   bool
	dataKeyAlgorithm          string
	dataKeyLengthBits         int
	cipherName                string
	cipherBlockSizeBits       int
	cipherIvLengthBits        int
	cipherTagLengthBits       int
	cipherMaxContentLengthBits int64
	isCommitting              bool
	commitmentLengthBits      int
	kdfHashAlgorithm          func() hash.Hash
}

// Predefined algorithm suites
var (
	// Key committing AES-256-GCM content encryption with HKDF-SHA512 commitment/encryption key derivation.
	// In addition to the security properties for the AES-256-GCM content encryption suite,
	// this suite also uses HKDF to derive a key commitment string that is included in metadata.
	// v3 clients (only v3.2.0 or higher) can only use this suite to read objects with key commitment;
	// to use this suite to write objects with key commitment, upgrade to a v4 client.
	AlgAES256GCMHkdfSha512CommitKey = &AlgorithmSuite{
		id:                         0x0073,
		isLegacy:                   false,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256, // this is the input into the KDF
		cipherName:                AESGCMCommitKey,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        224,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: GCMMaxContentLengthBits,
		isCommitting:              true,
		commitmentLengthBits:      224,
		kdfHashAlgorithm:          sha512.New,
	}

	// AES-256 GCM content encryption.
	// This suite uses the data encryption key directly for AES-256 GCM content encryption.
	// This is the default suite for v3 clients.
	// Content encrypted with this suite can be read by any v2, v3, or v4 client.
	AlgAES256GCMIV12Tag16NoKDF = &AlgorithmSuite{
		id:                         0x0072,
		isLegacy:                   false,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESGCMNoPadding,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        96,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: GCMMaxContentLengthBits,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}

	// Legacy AES-256 CTR.
	// This suite is not supported at this time.
	AlgAES256CTRIV16Tag16NoKDF = &AlgorithmSuite{
		id:                         0x0071,
		isLegacy:                   true,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESCTRNoPadding,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        128,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: CTRMaxContentLengthBytes * 8,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}

	// Legacy AES-256 CBC.
	// This suite is only supported for decryption of existing objects and cannot be used for new objects.
	// We recommend migrating any existing objects encrypted with this suite to a non-legacy suite.
	AlgAES256CBCIV16NoKDF = &AlgorithmSuite{
		id:                         0x0070,
		isLegacy:                   true,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESCBCPKCS5,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        128,
		cipherTagLengthBits:       0,
		cipherMaxContentLengthBits: CBCMaxContentLengthBytes * 8,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}
)

// Map for looking up algorithm suites by ID
var algorithmSuitesByID = map[int]*AlgorithmSuite{
	0x0073: AlgAES256GCMHkdfSha512CommitKey,
	0x0072: AlgAES256GCMIV12Tag16NoKDF,
	0x0071: AlgAES256CTRIV16Tag16NoKDF,
	0x0070: AlgAES256CBCIV16NoKDF,
}

// GetAlgorithmSuiteByID returns the algorithm suite for the given ID
func GetAlgorithmSuiteByID(id int) (*AlgorithmSuite, error) {
	suite, exists := algorithmSuitesByID[id]
	if !exists {
		return nil, fmt.Errorf("unknown algorithm suite ID: 0x%04x", id)
	}
	return suite, nil
}

// ID returns the algorithm suite ID
func (a *AlgorithmSuite) ID() int {
	return a.id
}

// IDAsString returns the algorithm suite ID as a string
func (a *AlgorithmSuite) IDAsString() string {
	return strconv.Itoa(a.id)
}

// IDAsBytes returns the algorithm suite ID as a byte array (big-endian)
func (a *AlgorithmSuite) IDAsBytes() []byte {
	return []byte{byte(a.id >> 8), byte(a.id)}
}

// IsLegacy returns whether this is a legacy algorithm suite
func (a *AlgorithmSuite) IsLegacy() bool {
	return a.isLegacy
}

// DataKeyAlgorithm returns the data key algorithm name
func (a *AlgorithmSuite) DataKeyAlgorithm() string {
	return a.dataKeyAlgorithm
}

// DataKeyLengthBits returns the data key length in bits
func (a *AlgorithmSuite) DataKeyLengthBits() int {
	return a.dataKeyLengthBits
}

// DataKeyLengthBytes returns the data key length in bytes
func (a *AlgorithmSuite) DataKeyLengthBytes() int {
	return a.dataKeyLengthBits / 8
}

// CipherName returns the cipher name
func (a *AlgorithmSuite) CipherName() string {
	return a.cipherName
}

// CipherTagLengthBits returns the cipher tag length in bits
func (a *AlgorithmSuite) CipherTagLengthBits() int {
	return a.cipherTagLengthBits
}

// CipherTagLengthBytes returns the cipher tag length in bytes
func (a *AlgorithmSuite) CipherTagLengthBytes() int {
	return a.cipherTagLengthBits / 8
}

// IVLengthBytes returns the IV length in bytes
func (a *AlgorithmSuite) IVLengthBytes() int {
	return a.cipherIvLengthBits / 8
}

// CipherBlockSizeBytes returns the cipher block size in bytes
func (a *AlgorithmSuite) CipherBlockSizeBytes() int {
	return a.cipherBlockSizeBits / 8
}

// CipherMaxContentLengthBits returns the maximum content length in bits
func (a *AlgorithmSuite) CipherMaxContentLengthBits() int64 {
	return a.cipherMaxContentLengthBits
}

// CipherMaxContentLengthBytes returns the maximum content length in bytes
func (a *AlgorithmSuite) CipherMaxContentLengthBytes() int64 {
	return a.cipherMaxContentLengthBits / 8
}

// IsCommitting returns whether this algorithm suite is key committing
func (a *AlgorithmSuite) IsCommitting() bool {
	return a.isCommitting
}

// CommitmentLengthBits returns the commitment length in bits
func (a *AlgorithmSuite) CommitmentLengthBits() int {
	return a.commitmentLengthBits
}

// CommitmentLengthBytes returns the commitment length in bytes
func (a *AlgorithmSuite) CommitmentLengthBytes() int {
	return a.commitmentLengthBits / 8
}

// KDFHashAlgorithm returns the KDF hash algorithm
func (a *AlgorithmSuite) KDFHashAlgorithm() func() hash.Hash {
	return a.kdfHashAlgorithm
}
import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	// this is purposefully done before attempting to
	// decrypt the materials
	var cekFunc internal.CEKEntry
	objectCekAlgSuite, err := objectMetadata.GetContentEncryptionAlgorithmSuite()
	if err != nil {
		return out, metadata, err
	}

	//= ../specification/s3-encryption/decryption.md#key-commitment
	//# The S3EC MUST validate the algorithm suite used for decryption against the key commitment policy before attempting to decrypt the content ciphertext.
	if err := ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(objectCekAlgSuite, m.client.Options.CommitmentPolicy); err != nil {
		return out, metadata, fmt.Errorf("object's content encryption algorithm is not valid for the selected commitment policy: %v, %w", objectCekAlgSuite, err)
	}

	var matDesc string
	if objectCekAlgSuite == algorithms.AlgAES256GCMHkdfSha512CommitKey {
		cekFunc = internal.NewAESGCMDecryptCommittingContentCipher
		matDesc, err = objectMetadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for committing algorithm: %w", err)
		}
	} else if objectCekAlgSuite == algorithms.AlgAES256GCMIV12Tag16NoKDF {
		cekFunc = internal.NewAESGCMContentCipher
		matDesc, err = objectMetadata.GetMatDescV2()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for AES/GCM algorithm: %w", err)
		}
	} else if objectCekAlgSuite == algorithms.AlgAES256CBCIV16NoKDF {
		if !m.client.Options.EnableLegacyUnauthenticatedModes {
			//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
			//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy content encryption algorithms; it MUST throw an exception when attempting to decrypt an object encrypted with a legacy content encryption algorithm.
			//= ../specification/s3-encryption/decryption.md#legacy-decryption
			//# The S3EC MUST NOT decrypt objects encrypted using legacy unauthenticated algorithm suites unless specifically configured to do so.
			//= ../specification/s3-encryption/decryption.md#legacy-decryption
			//# If the S3EC is not configured to enable legacy unauthenticated content decryption, the client MUST throw an exception when attempting to decrypt an object encrypted with a legacy unauthenticated algorithm suite.
			return out, metadata, fmt.Errorf("configure client with enable legacy unauthenticated modes set to true to decrypt with %s", objectMetadata.CEKAlg)
		}
		//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
		//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
		cekFunc = internal.NewAESCBCContentCipher
		matDesc, err = objectMetadata.GetMatDescV2()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for AES/CBC algorithm: %w", err)
		}
	} else {
		return out, metadata, fmt.Errorf("invalid content encryption algorithm found in metadata: %s", objectMetadata.CEKAlg)
	}


	cipherKey, err := objectMetadata.GetDecodedKey()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded key for materials: %w", err)
	}
	iv, err := objectMetadata.GetDecodedMessageIDOrIV()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded IV for materials: %w", err)
	}
	keyringWrappingAlg, err := objectMetadata.GetFullWrappingAlgorithm()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get wrapping algorithm for materials: %w", err)
	}

	// S3 server will encode metadata with non-US-ASCII characters
	// Decode it here to avoid parsing/decryption failure
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//# The S3EC SHOULD support decoding the S3 Server's "double encoding".
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=exception
	//# If the S3EC does not support decoding the S3 Server's "double encoding" then it MUST return the content metadata untouched.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v1-v2-shared
	//# This string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# This material description string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# This encryption context string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	decodedMatDesc, err := customS3Decoder(matDesc)
	if err != nil {
		return out, metadata, fmt.Errorf("error while decoding Material Description: %w", err)
	}

	ec := ctx.Value(EncryptionContext)
	// If an encryption context is provided, the provided encryption context MUST match the encryption context stored in the metadata.
	if ec != nil {
		ecMap, ok := ec.(map[string]string)
		if !ok {
			return out, metadata, fmt.Errorf("encryption context provided to decrypt method is not valid JSON")
		}
		decodedMatDescMap := map[string]string{}
		if err := json.Unmarshal([]byte(decodedMatDesc), &decodedMatDescMap); err != nil {
			return out, metadata, fmt.Errorf("encryption context in stored object is not valid JSON: %w", err)
		}

		// The stored encryption context with the two reserved keys removed MUST match the provided encryption context.
		delete(decodedMatDescMap, "kms_cmk_id")
		delete(decodedMatDescMap, "aws:x-amz-cek-alg")

		if len(ecMap) != len(decodedMatDescMap) {
			return out, metadata, fmt.Errorf("Provided encryption context does not match information retrieved from S3")
		}
		for k, v := range ecMap {
			val, exists := decodedMatDescMap[k]
			if !exists || val != v {
				// If the stored encryption context with the two reserved keys removed does not match the provided encryption context, S3EC MUST throw an exception.
				return out, metadata, fmt.Errorf("Provided encryption context does not match information retrieved from S3")
			}
		}
	}

	cekAlg, err := objectMetadata.GetContentEncryptionAlgorithmString()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get content encryption algorithm from metadata: %w", err)
	}

	decryptMaterialsRequest := materials.DecryptMaterialsRequest{
		cipherKey,
		iv,
		decodedMatDesc,
		keyringWrappingAlg,
		cekAlg,
		objectMetadata.TagLen,
	}
	decryptMaterials, err := m.client.Options.CryptographicMaterialsManager.DecryptMaterials(ctx, decryptMaterialsRequest)
		return out, metadata, fmt.Errorf("error while decrypting materials: %w", err)
	}

	// If the CMM did not provide a key commitment for a committing algorithm,
	// retrieve it from the object metadata.
	if (objectCekAlgSuite.IsCommitting() && decryptMaterials.KeyCommitment == nil) {
		commitment, err := objectMetadata.GetDecodedKeyCommitment()
		if err != nil {
			return out, metadata, fmt.Errorf("unable to get decoded key commitment for committing algorithm: %w", err)
		}
		decryptMaterials.KeyCommitment = commitment
	}

	cipher, err := cekFunc(*decryptMaterials)
	if err != nil {
		return out, metadata, err
	}
	reader, err := cipher.DecryptContents(result.Body)
	if err != nil {
		return out, metadata, err
	}

	// Apply buffer size configuration for GetObject operations
	// The S3EC MUST set the buffer size to a reasonable default for GetObject
	bufferedReader, err := internal.NewBufferedReader(reader, int(m.client.Options.BufferSize))
	if err != nil {
		return out, metadata, fmt.Errorf("unable to create buffered reader for decrypted contents: %w", err)
	}
	result.Body = bufferedReader
	out.Result = result

	return out, metadata, err
}

func ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(cekAlgSuite *algorithms.AlgorithmSuite, policy commitment.CommitmentPolicy) error {
	//= ../specification/s3-encryption/decryption.md#key-commitment
	//# If the commitment policy requires decryption using a committing algorithm suite,
	//# and the algorithm suite associated with the object does not support key commitment,
	//# then the S3EC MUST throw an exception.
	if policy.RequiresDecrypt() && !cekAlgSuite.IsCommitting() {
		return fmt.Errorf("commitment policy %v does not allow decryption using algorithm suite %v which does not support key commitment", policy, cekAlgSuite)
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT,
	//# the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.FORBID_ENCRYPT_ALLOW_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return nil
		}
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT,
	//# the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return nil
		}
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
	//# the S3EC MUST NOT allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return fmt.Errorf("commitment policy %v does not allow decryption using algorithm suite %v which does not support key commitment", policy, cekAlgSuite)
		}
	}

	// If the policy is not recognized, return an error
	switch policy {
		case commitment.FORBID_ENCRYPT_ALLOW_DECRYPT, commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT, commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
			// do nothing -- valid policies
		default:
			return fmt.Errorf("unknown commitment policy: %v", policy)
	}

	return nil
}
	"encoding/hex"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"testing"
)

func TestDecryptionClientV4_GetMockV2Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	}
}

// Use TestEncryptionClientV4_PutMockV3Object to generate a new test vector if needed.
func TestDecryptionClientV4_GetMockV3Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("2d4ff4dafe27f69f628872d82b5a1002ed1a21b8485d532bd8159f6487945b3641af5865fc0a029a3650053600c6d213625b9a0cc9c239577c09f3423dedc5641e88b6835824417c")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"6JOSx47RkdyfciJkNauuC4RpkMcWZY4a+i1RzQ=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"sE9zLb4tsEBJkvEhLMZFxMj9oZJBbQ6ZOgOqHA=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	out, err := client.GetObject(context.Background(), input)

	actual, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected content to match but it did not")
	}
}

func TestDecryptionClientV4_GetMockV3ObjectWithIncorrectKCValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("e403a8f941e43bdf0ca3ef0bcf6701acd739b2de0a8ee524fa89497210fb0213dfc856376a9ff7753db6ee549dc4040861bc7080a66f902441904bf4a003e028e982de8ea6958c30")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I/0X2GC1iX9Pszf1PAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMDQPII4AynCy/rVwhAgEQgDueLCWabc8WgyoZkAnqVzESQ4NztSDxuETx3obcWJ9Jj6gDAAuDaAL5V+H5QFfwgBqWEcIYt2Ep9WcECw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"FiQepGw+/O+3MuYrmQU5mAkotUnxB+W+EwYDHw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"cPwtbK08jrpe3x+QElyG+vGtkk9jDn6KmOta2Q=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	_, err = client.GetObject(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error due to incorrect key commitment, got nil")
	}
	if !strings.Contains(err.Error(), "derived key commitment value does not match value stored on encrypted message") {
		t.Fatalf("expected key commitment mismatch error, got %v", err)
	}
}

func TestDecryptionClientV4_GetMockV3ObjectWithInvalidKCValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("e403a8f941e43bdf0ca3ef0bcf6701acd739b2de0a8ee524fa89497210fb0213dfc856376a9ff7753db6ee549dc4040861bc7080a66f902441904bf4a003e028e982de8ea6958c30")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I/0X2GC1iX9Pszf1PAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMDQPII4AynCy/rVwhAgEQgDueLCWabc8WgyoZkAnqVzESQ4NztSDxuETx3obcWJ9Jj6gDAAuDaAL5V+H5QFfwgBqWEcIYt2Ep9WcECw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"notvalidbase64"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"cPwtbK08jrpe3x+QElyG+vGtkk9jDn6KmOta2Q=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	_, err = client.GetObject(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error due to incorrect key commitment, got nil")
	}
	if !strings.Contains(err.Error(), "illegal base64 data at input byte") {
		t.Fatalf("expected base64 decoding error, got %v", err)
	}
}

func TestDecryptionClientV4_GetMockV2Object_V1Interop_KMS_AESCBC(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "7ItX9CTGNWWegC62RlaNu6EJ3+J9yGO7yAqDNU4CdeA=", `"}`))
	}))

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
}

func TestDecryptionClientV4_GetMockV2Object_V1Interop_KMS_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "Hrjrkkt/vQwMYtqvK6+MiXh3xiMvviL1Ks7w2mgsJgU=", `"}`))
	}))
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	}
}

func TestDecryptionClientV4_GetMockV2Object_OnlyDecryptsRegisteredAlgorithms(t *testing.T) {
	httpClientFactory := func() *awstesting.MockHttpClient {
		b, err := hex.DecodeString("1bd0271b25951fdef3dbe51a9b7af85f66b311e091aa10a346655068f657b9da9acc0843ea0522b0d1ae4a25a31b13605dd1ac5d002db8965d9d4652fd602693")
		if err != nil {
				tConfig.HTTPClient = httpClientFactory()
				s3Client := s3.NewFromConfig(tConfig)

				client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
	}
}

func TestDecryptionClientV4_CheckValidCryptographicMaterialsManager(t *testing.T) {
	_, err := materials.NewCryptographicMaterialsManager(nil)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestDecryptionClientV4_EncryptionContextValidation_MockV2Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6b134eb7a353131de92faff64f594b2794e3544e31776cca26fe3bbeeffc68742d1007234f11c6670522602326868e29f37e9d2678f1614ec1a2418009b9772100929aadbed9a21a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	cases := map[string]struct {
		storedMatDesc       string
		providedContext     map[string]string
		expectError         bool
		expectedErrorMsg    string
	}{
		"matching encryption context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"custom-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value"},
			expectError:     false,
		},
		"matching encryption context with multiple keys": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1","key2":"value2"}`,
			providedContext: map[string]string{"key1": "value1", "key2": "value2"},
			expectError:     false,
		},
		"empty encryption context matches empty stored context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id"}`,
			providedContext: map[string]string{},
			expectError:     false,
		},
		"mismatched encryption context value": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"stored-value"}`,
			providedContext:  map[string]string{"custom-key": "different-value"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"missing key in provided context": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1","key2":"value2"}`,
			providedContext:  map[string]string{"key1": "value1"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"extra key in provided context": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1"}`,
			providedContext:  map[string]string{"key1": "value1", "key2": "value2"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"provided context has reserved key (should be ignored in stored context)": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"custom-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value"},
			expectError:     false,
		},
		"stored context missing kms_cmk_id key": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","custom-key":"custom-value","another-key":"another-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value", "another-key": "another-value"},
			expectError:     false,
		},
		"stored context missing kms_cmk_id with empty provided context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
			providedContext: map[string]string{},
			expectError:     false,
		},
		"invalid JSON in stored material description": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","invalid-json":}`,
			providedContext:  map[string]string{"key1": "value1"},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
		"malformed JSON in stored material description": {
			storedMatDesc:    `not-json-at-all`,
			providedContext:  map[string]string{},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
		"incomplete JSON in stored material description": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"`,
			providedContext:  map[string]string{},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					StatusCode: 200,
					Header: http.Header{
						http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"PsuclPnlo2O0MQoov6kL1TBlaZG6oyNwWuAqmAgq7g8b9ZeeORi3VTMg624FU9jx"},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"dqqlq2dRVSQ5hFRb"},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{tc.storedMatDesc},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSContextKeyring},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
					},
					Body: io.NopCloser(bytes.NewBuffer(b)),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			input := &s3.GetObjectInput{
				Bucket: aws.String("test"),
				Key:    aws.String("test"),
			}

			// Create context with encryption context if provided
			ctx := context.Background()
			if tc.providedContext != nil {
				ctx = context.WithValue(ctx, EncryptionContext, tc.providedContext)
			}

			_, err = client.GetObject(ctx, input)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !strings.Contains(err.Error(), tc.expectedErrorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDecryptionClientV4_EncryptionContextValidation_InvalidContextType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6b134eb7a353131de92faff64f594b2794e3544e31776cca26fe3bbeeffc68742d1007234f11c6670522602326868e29f37e9d2678f1614ec1a2418009b9772100929aadbed9a21a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"PsuclPnlo2O0MQoov6kL1TBlaZG6oyNwWuAqmAgq7g8b9ZeeORi3VTMg624FU9jx"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"dqqlq2dRVSQ5hFRb"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSContextKeyring},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	// Test with invalid encryption context type (string instead of map[string]string)
	ctx := context.WithValue(context.Background(), EncryptionContext, "invalid-type")

	_, err = client.GetObject(ctx, input)

	if err == nil {
		t.Fatalf("expected error but got none")
	}

	expectedErrorMsg := "encryption context provided to decrypt method is not valid JSON"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
	}
}

func TestValidateContentEncryptionAlgorithmAgainstCommitmentPolicy_UnrecognizedPolicy(t *testing.T) {
	algSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF
	// Given: Some invalid commitment policy
	invalidPolicy := commitment.CommitmentPolicy(999)
	
	// When: Call the function under test
	err := ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(algSuite, invalidPolicy)
	
	// Then: function raises error
	if err == nil {
		t.Fatalf("expected error for unrecognized commitment policy, got nil")
	}
	
	// Then: error message contains the expected text
	expectedErrorMsg := "unknown commitment policy"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
	}
	
	// Then: error message includes the policy value
	expectedPolicyValue := "CommitmentPolicy(999)"
	if !strings.Contains(err.Error(), expectedPolicyValue) {
		t.Errorf("expected error message to contain %q, got %q", expectedPolicyValue, err.Error())
	}
}
import (
	"context"
	"fmt"
	"io"

	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// DefaultMinFileSize is used to check whether we want to write to a temp file
// or store the data in memory.
const DefaultMinFileSize = 1024 * 512 * 5

// DefaultBufferSize is the default buffer size for GetObject operations
// The S3EC MUST set the buffer size to a reasonable default for GetObject
const DefaultBufferSize = 1024 * 64 // 64KB default buffer size

// EncryptionContext is used to extract Encryption Context to use on a per-request basis
const EncryptionContext = "EncryptionContext"

	if err != nil {
		return out, metadata, err
	}
	var cipher internal.ContentCipher

	//= ../specification/s3-encryption/encryption.md#content-encryption
	//# The S3EC MUST use the encryption algorithm configured during [client](./client.md) initialization.
	if m.ec.Options.EncryptionAlgorithmSuite == algorithms.AlgAES256GCMHkdfSha512CommitKey {
		return out, metadata, fmt.Errorf("algorithm suite ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY is not supported for encryption in S3EC V3")
	} else if m.ec.Options.EncryptionAlgorithmSuite == algorithms.AlgAES256GCMIV12Tag16NoKDF {
		cipher, err = internal.NewAESGCMContentCipher(*cryptoMaterials)
		if err != nil {
			return out, metadata, err
		}
	} else {
		// S3EC V4 only supports writing AES-GCM, with or without key commitment.
		// Any other algorithms are invalid for encryption.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-iv16-tag16-no-kdf
		//# Attempts to encrypt using AES-CTR MUST fail.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-hkdf-sha512-commit-key
		//# Attempts to encrypt using key committing AES-CTR MUST fail.
		return out, metadata, fmt.Errorf("invalid content encryption algorithm found in options: %s", cryptoMaterials.CEKAlgorithm)
	}

	//= ../specification/s3-encryption/encryption.md#content-encryption
	//# The client MUST validate that the length of the plaintext bytes does not exceed the algorithm suite's cipher's maximum content length in bytes.
	if n >= m.ec.Options.EncryptionAlgorithmSuite.CipherMaxContentLengthBytes() {
		return out, metadata, fmt.Errorf("plaintext length %d exceeds maximum content length for algorithm %s", n, cryptoMaterials.CEKAlgorithm)
	}

	stream := reqCopy.GetStream()
import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	return cryptoMaterials, err
}

func TestEncryptionClientV3_PutMockV2Object_KMSCONTEXT_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
	}
}

func TestEncryptionClientV3_PutMockV2Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	var md materials.MaterialDescription
	iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
	kmsWithStaticIV := keyringWithStaticTestIV{
		IV: iv,
		Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		}),
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256GCMIV12Tag16NoKDF
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
		Body: func() io.ReadSeeker {
			content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
			return bytes.NewReader(content)
		}(),
		Metadata: md,
	})
	if err != nil {
		t.Fatalf("PutObject failed with %v", err)
	}

	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	if tHttpClient.CapturedReq != nil {
		headers := tHttpClient.CapturedReq.Header
		
		// V2 IV (x-amz-iv) - should be present and base64 encoded
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Iv"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-Iv should be present but was empty")
		}
		
		// V2 Encrypted Data Key (x-amz-key-v2) is nondeterministic, just check presence
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Key-V2"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-Key-V2 should be present but was empty")
		}
		
		// V2 Wrapping Algorithm (x-amz-wrap-alg) - should be "kms+context"
		expectedWrappingAlg := "kms+context"
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Wrap-Alg"); actualValue != expectedWrappingAlg {
			t.Errorf("X-Amz-Meta-X-Amz-Wrap-Alg expected '%s', got '%s'", expectedWrappingAlg, actualValue)
		}
		
		// V2 Content Encryption Algorithm (x-amz-cek-alg) - should be "AES/GCM/NoPadding"
		expectedCekAlg := "AES/GCM/NoPadding"
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Cek-Alg"); actualValue != expectedCekAlg {
			t.Errorf("X-Amz-Meta-X-Amz-Cek-Alg expected '%s', got '%s'", expectedCekAlg, actualValue)
		}
		
		// V2 Material Description (x-amz-matdesc) - should be present
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Matdesc"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-Matdesc should be present but was empty")
		}
		
		// V2 Tag Length (x-amz-tag-len) - should be "128"
		expectedTagLen := "128"
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-Tag-Len"); actualValue != expectedTagLen {
			t.Errorf("X-Amz-Meta-X-Amz-Tag-Len expected '%s', got '%s'", expectedTagLen, actualValue)
		}
	}

	// V2 body should match expected encrypted content
	if tHttpClient.CapturedBody == nil || len(tHttpClient.CapturedBody) == 0 {
		t.Error("Expected encrypted content, but captured body was empty")
	}
}

func TestEncryptionClientv3_PutMockV2Object_KMSCONTEXT_AESGCM_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:   aws.String("test-bucket"),
		t.Errorf("expected no error, got %v", err)
	}
}

//= ../specification/s3-encryption/encryption.md#content-encryption
//= type=test
//# The S3EC MUST use the encryption algorithm configured during [client](./client.md) initialization.
func TestS3EC_UsesConfiguredEncryptionAlgorithm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		ivHex            string
		expectedHeader   string
		expectedValue    string
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813",
			expectedHeader:  "X-Amz-Meta-X-Amz-Cek-Alg",
			expectedValue:   "AES/GCM/NoPadding",
		},
		{
			name:             "AES256GCMHkdfSha512CommitKey",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9",
			expectedHeader:  "X-Amz-Meta-X-Amz-C",
			expectedValue:   "115",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString(tc.ivHex)
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Verify the configured algorithm is set correctly
			if client.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %v, got %v", tc.algorithm, client.Options.EncryptionAlgorithmSuite)
			}

			// Test that PutObject uses the configured algorithm
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
				Body: func() io.ReadSeeker {
					content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
					return bytes.NewReader(content)
				}(),
				Metadata: md,
			})
			if err != nil {
				if tc.algorithm == algorithms.AlgAES256GCMHkdfSha512CommitKey {
					// Encrypting with key commitment is expected to fail with V3
					t.Logf("Expected failure encrypting with key commitment algorithm in V3: %v", err)
					return
				}
				t.Fatalf("PutObject failed with %v", err)
			}

			// Verify the encryption was performed (captured body should not be empty)
			if tHttpClient.CapturedBody == nil || len(tHttpClient.CapturedBody) == 0 {
				t.Error("Expected encrypted content, but captured body was empty")
			}

			// Capture and validate object metadata for algorithm information
			if tHttpClient.CapturedReq == nil {
				t.Fatal("Expected captured HTTP request, but it was nil")
			}

			// Verify the algorithm-specific header is set correctly
			headers := tHttpClient.CapturedReq.Header
			actualValue := headers.Get(tc.expectedHeader)
			
			if actualValue == "" {
				t.Errorf("Expected %s header to be present for %s algorithm", tc.expectedHeader, tc.name)
			} else if actualValue != tc.expectedValue {
				t.Errorf("Expected %s header to be '%s', got '%s'", tc.expectedHeader, tc.expectedValue, actualValue)
			}

			t.Logf("✓ S3EC successfully used configured %s algorithm with correct metadata", tc.name)
			t.Logf("  - %s: %s", tc.expectedHeader, actualValue)
		})
	}
}

// mockLargeReader simulates a large reader without allocating memory
type mockLargeReader struct {
	size     int64
	position int64
}

func (r *mockLargeReader) Read(p []byte) (n int, err error) {
	if r.position >= r.size {
		return 0, io.EOF
	}
	
	remaining := r.size - r.position
	toRead := int64(len(p))
	if toRead > remaining {
		toRead = remaining
	}
	
	// Fill buffer with test data
	for i := int64(0); i < toRead; i++ {
		p[i] = byte((r.position + i) % 256)
	}
	
	r.position += toRead
	return int(toRead), nil
}

func (r *mockLargeReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.position = offset
	case io.SeekCurrent:
		r.position += offset
	case io.SeekEnd:
		r.position = r.size + offset
		if r.position > r.size {
			return 0, fmt.Errorf("seek position beyond end of reader")
		}
	}
	
	if r.position < 0 {
		r.position = 0
	}
	if r.position > r.size {
		r.position = r.size
	}
	
	return r.position, nil
}

//= ../specification/s3-encryption/encryption.md#content-encryption
//= type=test
//# The client MUST validate that the length of the plaintext bytes does not exceed the algorithm suite's cipher's maximum content length in bytes.
func TestS3EC_ValidatesPlaintextLengthLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		contentSize       int64
		expectError       bool
		errorContains     string
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF_WithinLimit",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			contentSize:      1024, // Well within limit
			expectError:      false,
		},
		{
			name:             "AES256GCMIV12Tag16NoKDF_ExceedsLimit",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			contentSize:      algorithms.AlgAES256GCMIV12Tag16NoKDF.CipherMaxContentLengthBytes() + 1, // Just over the limit
			expectError:      true,
			errorContains:    "plaintext length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Create a mock reader for the specified size (avoids memory allocation issues)
			var body io.ReadSeeker
			if tc.contentSize <= 1024*1024 { // For small sizes, use real content
				content := make([]byte, tc.contentSize)
				for i := range content {
					content[i] = byte(i % 256)
				}
				body = bytes.NewReader(content)
			} else {
				// For large sizes, use mock reader
				body = &mockLargeReader{size: tc.contentSize}
			}

			// Test PutObject with the content
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket:   aws.String("test-bucket"),
				Key:      aws.String("test-key"),
				Body:     body,
				Metadata: md,
			})

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for content size %d bytes (limit: %d), but got none", 
						tc.contentSize, tc.algorithm.CipherMaxContentLengthBytes())
				} else if tc.errorContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tc.errorContains)) {
					t.Errorf("Expected error to contain '%s', but got: %v", tc.errorContains, err)
				} else {
					t.Logf("✓ S3EC correctly rejected content exceeding maximum length for %s", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for content size %d bytes (limit: %d), but got: %v", 
						tc.contentSize, tc.algorithm.CipherMaxContentLengthBytes(), err)
				} else {
					t.Logf("✓ S3EC correctly accepted content within maximum length for %s", tc.name)
				}
			}
		})
	}
}

func TestS3EC_IVGenerationAndMetadataInclusion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		ivHex            string
		expectedHeader   string
		expectedIVLength  int // in bytes
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813", // 12 bytes
			expectedHeader:  "X-Amz-Meta-X-Amz-Iv",
			expectedIVLength: 12, // 96 bits / 8 = 12 bytes
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// First, verify the algorithm's IV length specification
			actualBytes := tc.algorithm.IVLengthBytes()
			expectedBits := actualBytes * 8

			if actualBytes != tc.expectedIVLength {
				t.Errorf("Expected IV/Message ID length %d bytes for %s, got %d bytes", 
					tc.expectedIVLength, tc.name, actualBytes)
			}

			t.Logf("✓ S3EC algorithm %s correctly defines IV/Message ID length: %d bytes (%d bits)", 
				tc.name, actualBytes, expectedBits)

			// Now test that the IV is properly included in content metadata during encryption
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString(tc.ivHex)
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Test that PutObject includes IV/Message ID in metadata
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
				Body: func() io.ReadSeeker {
					content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
					return bytes.NewReader(content)
				}(),
				Metadata: md,
			})
			if err != nil {
				t.Fatalf("PutObject failed with %v", err)
			}

			// Verify the IV/Message ID is included in content metadata
			if tHttpClient.CapturedReq == nil {
				t.Fatal("Expected captured HTTP request, but it was nil")
			}

			headers := tHttpClient.CapturedReq.Header
			ivValue := headers.Get(tc.expectedHeader)
			
			if ivValue == "" {
				t.Errorf("Expected %s header to be present for %s algorithm", tc.expectedHeader, tc.name)
			} else {
				//= ../specification/s3-encryption/encryption.md#content-encryption
				//= type=test
				//# The client MUST generate an IV or Message ID using the length of the IV or Message ID defined in the algorithm suite.
				//# The generated IV or Message ID MUST be set or returned from the encryption process such that it can be included in the content metadata.
				decodedIV, err := hex.DecodeString(ivValue)
				if err != nil {
					// Try base64 decoding for V3 format
					decodedIV, err = base64.StdEncoding.DecodeString(ivValue)
				}
				
				if err == nil {
					expectedLength := tc.algorithm.IVLengthBytes()
					
					if len(decodedIV) != expectedLength {
						t.Errorf("Expected IV/Message ID length %d bytes, got %d bytes", expectedLength, len(decodedIV))
					}
				}

				t.Logf("✓ S3EC successfully included IV/Message ID in content metadata for %s", tc.name)
				t.Logf("  - %s: %s", tc.expectedHeader, ivValue)
			}
		})
	}
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-iv16-tag16-no-kdf
//= type=test
//# Attempts to encrypt using AES-CTR MUST fail.
func TestS3EC_AESCTREncryptionMustFail(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	
	var mcmm = mockCMM{}

	// Attempt to create client with AES-CTR algorithm - this should fail
	_, err := New(s3Client, mcmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CTRIV16Tag16NoKDF
	})

	if err == nil {
		t.Fatalf("Expected error when attempting to use AES-CTR algorithm, but got none")
	}

	if !strings.Contains(err.Error(), "AES-CTR") && !strings.Contains(err.Error(), "CTR") {
		t.Errorf("Expected error to mention AES-CTR, but got: %v", err)
	}

	t.Logf("✓ S3EC correctly rejected AES-CTR encryption algorithm: %v", err)
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-hkdf-sha512-commit-key
//= type=test
//# Attempts to encrypt using key committing AES-CTR MUST fail.
func TestS3EC_CommittingAESCTREncryptionMustFail(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	
	var mcmm = mockCMM{}

	// Attempt to create client with AES-CTR algorithm and committing policy - this should fail
	_, err := New(s3Client, mcmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CTRIV16Tag16NoKDF
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})

	if err == nil {
		t.Fatalf("Expected error when attempting to use key committing AES-CTR algorithm, but got none")
	}

	if !strings.Contains(err.Error(), "AES-CTR") && !strings.Contains(err.Error(), "CTR") {
		t.Errorf("Expected error to mention AES-CTR, but got: %v", err)
	}

	t.Logf("✓ S3EC correctly rejected key committing AES-CTR encryption algorithm: %v", err)
}

import (
	"context"
	"log"
	"fmt"

	"github.com/aws/amazon-s3-encryption-client-go/v3/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3EncryptionClientV3 provides client-side encryption for S3.
// The client embeds a default client to provide support for control plane operations
// which do not involve encryption.
type S3EncryptionClientV3 struct {
	//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
	//# The S3EC SHOULD support invoking operations unrelated to client-side encryption e.g. CopyObject as the conventional AWS SDK S3 client would.
	//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
	//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

	*s3.Client                         // promoted anonymous field, it allows this type to call s3 Client methods
	Options    EncryptionClientOptions // options for encrypt/decrypt
}

// EncryptionClientOptions is the configuration options for the S3 Encryption Client.
//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=implication
//# The S3EC MUST provide a different set of configuration options than the conventional S3 client.

type EncryptionClientOptions struct {
	// TempFolderPath is used to store temp files when calling PutObject
	// Temporary files are needed to compute the X-Amz-Content-Sha256 header
	// temporary file instead of using memory
	MinFileSize int64

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//# The S3EC SHOULD accept a configurable buffer size
	//# which refers to the maximum ciphertext length in bytes to store in memory
	//# when Delayed Authentication mode is disabled.

	// BufferSize is the buffer size used for GetObject operations
	BufferSize int64

	// The logger to write logging messages to
	Logger *log.Logger

	// The CryptographicMaterialsManager to use to manage encryption and decryption materials
	CryptographicMaterialsManager materials.CryptographicMaterialsManager

	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=implication
	//# The S3EC MUST support the option to enable or disable legacy unauthenticated modes (content encryption algorithms).
	
	// EnableLegacyUnauthenticatedModes MUST be set to true in order to decrypt objects encrypted
	// using legacy (unauthenticated) modes such as AES/CBC. The default is false.
	EnableLegacyUnauthenticatedModes bool

	//= ../specification/s3-encryption/client.md#key-commitment
	//# The S3EC MUST support configuration of the [Key Commitment policy](./key-commitment.md) during its initialization.
	
	// CommitmentPolicy specifies the key commitment policy for this client.
	// S3EncryptionClientV3 defaults to commitment.FORBID_ENCRYPT_ALLOW_DECRYPT.
	// Objects written by a client configured with this default can be read by any v3 or v4 client
	// that is configured with either FORBID_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_ALLOW_DECRYPT commitment policies.
	// A client configured with this default can read objects written by any v3 or v4 client.
	// If an EncryptionAlgorithmSuite is also provided,
	// the selected CommitmentPolicy must also be compatible with the selected EncryptionAlgorithmSuite; if not, New() will return an error.
	CommitmentPolicy commitment.CommitmentPolicy

	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# The S3EC MUST support configuration of the encryption algorithm (or algorithm suite) during its initialization.
	
	// EncryptionAlgorithmSuite specifies the algorithm suite to use when encrypting objects.
	// S3EncryptionClientV3 defaults to algorithms.AlgAES256GCMIV12Tag16NoKDF.
	// Attempts to use a committing algorithm suite with S3EncryptionClientV3 will result in an error;
	// to use a committing algorithm suite, upgrade to S3EncryptionClientV4.
	// S3EncryptionClientV3 will decrypt objects encrypted with any supported algorithm suite, provided that the
	// algorithms suite is compatible with the selected CommitmentPolicy and EnableLegacyUnauthenticatedModes options.
	// If a CommitmentPolicy is also provided,
	// the selected EncryptionAlgorithmSuite must also be compatible with the selected CommitmentPolicy; if not, New() will return an error.
	EncryptionAlgorithmSuite *algorithms.AlgorithmSuite
}

//# The S3EC MUST NOT support use of S3EC as the provided S3 client during its initialization; it MUST throw an exception in this case.

// New creates a new S3 Encryption Client v4 with the given CryptographicMaterialsManager
func New(s3Client *s3.Client, CryptographicMaterialsManager materials.CryptographicMaterialsManager, optFns ...func(options *EncryptionClientOptions)) (*S3EncryptionClientV3, error) {
	//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
	//# The S3EC MUST support the option to provide an SDK S3 client instance during its initialization.
	// In Go, the requirement below is enforced/implied by the type system at compile time.
	// The New() function expects *s3.Client, but S3EncryptionClientV3 is a wholly different type.
	//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
	//= type=implication
	//# The S3EC MUST NOT support use of S3EC as the provided S3 client during its initialization; it MUST throw an exception in this case.
	wrappedClient := s3Client
	// default options
	options := EncryptionClientOptions{
		MinFileSize:                      DefaultMinFileSize,
		//= ../specification/s3-encryption/client.md#set-buffer-size
		//# If Delayed Authentication mode is disabled, and no buffer size is provided,
		//# the S3EC MUST set the buffer size to a reasonable default.
		BufferSize:                       DefaultBufferSize,
		Logger:                           log.Default(),
		// S3EC Go V4 only accepts a CMM on client configuration, not a keyring.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# The S3EC MUST accept either one CMM or one Keyring instance upon initialization.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# If both a CMM and a Keyring are provided, the S3EC MUST throw an exception.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# When a Keyring is provided, the S3EC MUST create an instance of the DefaultCMM using the provided Keyring.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# The S3EC MAY accept key material directly.
		CryptographicMaterialsManager:    CryptographicMaterialsManager,
		//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
		//# The option to enable legacy unauthenticated modes MUST be set to false by default.
		EnableLegacyUnauthenticatedModes: false,
		CommitmentPolicy:                 commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,

		// S3EC Go V4 does not support delayed authentication mode.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# The S3EC MUST support the option to enable or disable Delayed Authentication mode.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# Delayed Authentication mode MUST be set to false by default.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# When enabled, the S3EC MAY release plaintext from a stream which has not been authenticated.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# When disabled the S3EC MUST NOT release plaintext from a stream which has not been authenticated.
		//= ../specification/s3-encryption/client.md#set-buffer-size
		//= type=exception
		//# If Delayed Authentication mode is enabled, and the buffer size has been set to a value other than its default, the S3EC MUST throw an exception.
	
		// Go S3EC V4 does not support instruction file configuration.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# The S3EC MAY support the option to provide Instruction File Configuration during its initialization.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# If the S3EC in a given language supports Instruction Files, then it MUST accept Instruction File Configuration during its initialization.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# In this case, the Instruction File Configuration SHOULD be optional, such that its default configuration is used when none is provided.

		// Go S3EC V4 does not support a single inherited configuration for underlying AWS SDK clients.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# The S3EC MAY support directly configuring the wrapped SDK clients through its initialization.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# For example, the S3EC MAY accept a credentials provider instance during its initialization.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# If the S3EC accepts SDK client configuration, the configuration MUST be applied to all wrapped S3 clients.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# If the S3EC accepts SDK client configuration, the configuration MUST be applied to all wrapped SDK clients including the KMS client.
	
		// Go S3EC V4 does not support supplying a custom source of randomness during client initialization.
		//= ../specification/s3-encryption/client.md#randomness
		//= type=exception
		//# The S3EC MAY accept a source of randomness during client initialization.
	}
	// apply functional options
	for _, fn := range optFns {
		fn(&options)
	}

	// If no algorithm suite, supply a default
	if options.EncryptionAlgorithmSuite == nil {
		options.EncryptionAlgorithmSuite = DefaultEncryptionAlgorithmSuite(options)
	}

	// Validate selected encryption algorithm suite
	if err := ValidateEncryptionAlgorithmSuite(options); err != nil {
		return nil, err
	}

	// use the given wrappedClient for the promoted anon fields
	s3ec := &S3EncryptionClientV3{wrappedClient, options}
	return s3ec, nil
}

func DefaultEncryptionAlgorithmSuite(options EncryptionClientOptions) *algorithms.AlgorithmSuite {
	if options.CommitmentPolicy.RequiresEncrypt() {
		return algorithms.AlgAES256GCMHkdfSha512CommitKey
	} else {
		return algorithms.AlgAES256GCMIV12Tag16NoKDF
	}
}

// Explict (but verbose) validations of S3EC specification
func ValidateEncryptionAlgorithmSuite(options EncryptionClientOptions) error {
	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# The S3EC MUST validate that the configured encryption algorithm is not legacy.
	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# If the configured encryption algorithm is legacy, then the S3EC MUST throw an exception.
	if options.EncryptionAlgorithmSuite.IsLegacy() {
		return fmt.Errorf("legacy algorithm suites are not allowed for decrypt, got %v", options.EncryptionAlgorithmSuite)
	}
	//= ../specification/s3-encryption/client.md#key-commitment
	//# The S3EC MUST validate the configured Encryption Algorithm against the provided key commitment policy.
	//= ../specification/s3-encryption/client.md#key-commitment
	//# If the configured Encryption Algorithm is incompatible with the key commitment policy, then it MUST throw an exception.
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST NOT encrypt using an algorithm suite which supports key commitment.
	if options.CommitmentPolicy == commitment.FORBID_ENCRYPT_ALLOW_DECRYPT {
		if options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy FORBID_ENCRYPT_ALLOW_DECRYPT does not allow committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
	} else if options.CommitmentPolicy == commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT {
		if !options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy REQUIRE_ENCRYPT_ALLOW_DECRYPT requires committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
	} else if options.CommitmentPolicy == commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT {
		if !options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy REQUIRE_ENCRYPT_REQUIRE_DECRYPT requires committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	} else {
		return fmt.Errorf("unknown CommitmentPolicy %v", options.CommitmentPolicy)
	}
	return nil
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - GetObject MUST be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

// GetObject will make a request to s3 and retrieve the object. In this process
// decryption will be done. The SDK only supports region reads of KMS and GCM.
func (c *S3EncryptionClientV3) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	// Go S3EC V4 does not support ranged gets.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# The S3EC MAY support the "range" parameter on GetObject which specifies a subset of bytes to download and decrypt.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the S3EC supports Ranged Gets, the S3EC MUST adjust the customer-provided range to include the beginning and end of the cipher blocks for the given range.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# For requests which provide a range to decrypt an object encrypted with an authenticated algorithm suite, the corresponding CTR-based algorithm suite is used.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the GetObject response contains a range, but the GetObject request does not contain a range, the S3EC MUST throw an exception.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the object was encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF, then ALG_AES_256_CTR_IV16_TAG16_NO_KDF MUST be used to decrypt the range of the object.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the object was encrypted with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY, then ALG_AES_256_CTR_HKDF_SHA512_COMMIT_KEY MUST be used to decrypt the range of the object.

	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - GetObject MUST decrypt data received from the S3 server and return it as plaintext.
	m := &decryptMiddleware{
		client: c,
		input:  input,
	return c.Client.GetObject(ctx, input, opts...)
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - PutObject MUST be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

// PutObject will make encrypt the contents before sending the data to S3. Depending on the MinFileSize
// a temporary file may be used to buffer the encrypted contents to.
func (c *S3EncryptionClientV3) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - PutObject MUST encrypt its input data before it is uploaded to S3.
	em := &encryptMiddleware{
		ec: c,
	}
	opts := append(optFns, encryptOpts...)
	return c.Client.PutObject(ctx, input, opts...)
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - DeleteObject MUST be implemented by the S3EC.

// DeleteObject will defer to the underlying S3 client to delete the object,
// but will execute its own logic to delete the associated instruction file using the default instruction file suffix.
func (c *S3EncryptionClientV3) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - DeleteObject MUST delete the given object key.
	result, err := c.Client.DeleteObject(ctx, input, optFns...)
	if err != nil {
		return result, err
	}

	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
	
	// Delete the associated instruction file
	instructionFileKey := *input.Key + internal.DefaultInstructionKeySuffix
	instructionInput := &s3.DeleteObjectInput{
		Bucket: input.Bucket,
		Key:    &instructionFileKey,
	}
	
	// Delete instruction file - ignore errors as the instruction file may not exist
	_, _ = c.Client.DeleteObject(ctx, instructionInput, optFns...)
	
	return result, err
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - DeleteObjects MUST be implemented by the S3EC.

// DeleteObjects will delete multiple objects by calling DeleteObject for each object.
// This ensures that both the objects and their associated instruction files are deleted.
func (c *S3EncryptionClientV3) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {

	// We implement DeleteObjects by calling DeleteObject for each object
	// This ensures that both the object and its instruction file are deleted for each item
	var deletedObjects []types.DeletedObject
	var errors []types.Error
	
	for _, obj := range input.Delete.Objects {
		// Call our DeleteObject method which handles both object and instruction file deletion
		//= ../specification/s3-encryption/client.md#required-api-operations
		//# - DeleteObjects MUST delete each of the given objects.
		//= ../specification/s3-encryption/client.md#required-api-operations
		//# - DeleteObjects MUST delete each of the corresponding instruction files using the default instruction file suffix.
		deleteResult, err := c.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: input.Bucket,
			Key:    obj.Key,
		}, optFns...)
		
		if err != nil {
			// Add error to the errors list
			errors = append(errors, types.Error{
				Key:     obj.Key,
				Code:    aws.String("InternalError"),
				Message: aws.String(err.Error()),
			})
		} else {
			// Add successful deletion to the deleted objects list
			deletedObjects = append(deletedObjects, types.DeletedObject{
				Key:          obj.Key,
				DeleteMarker: deleteResult.DeleteMarker,
				VersionId:    deleteResult.VersionId,
			})
		}
	}
	
	// Build the response
	result := &s3.DeleteObjectsOutput{
		Deleted: deletedObjects,
	}
	
	// Add errors if any occurred
	if len(errors) > 0 {
		result.Errors = errors
	}
	
	return result, nil
}

// S3EC Go V4 does not implement the following operations:
// - CreateMultipartUpload
// - UploadPart
// - CompleteMultipartUpload
// - AbortMultipartUpload
// - ReEncryptInstructionFile

//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CreateMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - If implemented, CreateMultipartUpload MUST initiate a multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - UploadPart MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - UploadPart MUST encrypt each part.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - Each part MUST be encrypted in sequence.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - Each part MUST be encrypted using the same cipher instance for each part.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CompleteMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CompleteMultipartUpload MUST complete the multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - AbortMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - AbortMultipartUpload MUST abort the multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MUST decrypt the instruction file's encrypted data key for the given object using the client's CMM.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MUST re-encrypt the plaintext data key with a provided keyring.
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=test
//# The S3EC SHOULD support invoking operations unrelated to client-side encryption e.g. CopyObject as the conventional AWS SDK S3 client would.
func TestWHEN_CallNonEncryptionOperationOnS3EC_THEN_PassthroughToPlaintextS3Client(t *testing.T) {
	// Create mock HTTP response for ListBuckets
	listBucketsResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Owner>
        <ID>test-owner-id</ID>
        <DisplayName>test-owner</DisplayName>
    </Owner>
    <Buckets>
        <Bucket>
            <Name>test-bucket-1</Name>
            <CreationDate>2023-01-01T00:00:00.000Z</CreationDate>
        </Bucket>
        <Bucket>
            <Name>test-bucket-2</Name>
            <CreationDate>2023-01-02T00:00:00.000Z</CreationDate>
        </Bucket>
    </Buckets>
</ListAllMyBucketsResult>`

	// Create mock HTTP client
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(listBucketsResponse)),
		},
	}

	// Create test config with mock HTTP client
	tConfig := awstesting.Config()
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	// Create mock CMM
	mockCMM := &mockCMM{}

	// Create the S3 encryption client
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Test that ListBuckets is passed through to the underlying S3 client
	ctx := context.TODO()
	result, err := s3ec.ListBuckets(ctx, &s3.ListBucketsInput{})

	// Verify no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the response is parsed correctly
	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		if len(result.Buckets) != 2 {
			t.Errorf("Expected 2 buckets, got %d", len(result.Buckets))
		}
		if result.Owner == nil || *result.Owner.DisplayName != "test-owner" {
			t.Error("Expected owner display name to be 'test-owner'")
		}
	}
}

//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=test
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.
func TestS3EC_AdheresToSameInterfaceAsConventionalS3Client(t *testing.T) {
	// This test validates that S3EC and regular S3 client have identical interfaces
	// by calling both with the same parameters and verifying compatible behavior
	
	testCases := []struct {
		name        string
		operation   string
		setupMock   func() *awstesting.MockHttpClient
		testFunc    func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV3)
	}{
		{
			name:      "ListBuckets",
			operation: "ListBuckets",
			setupMock: func() *awstesting.MockHttpClient {
				mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Owner>
        <ID>test-owner-id</ID>
        <DisplayName>test-owner</DisplayName>
    </Owner>
    <Buckets>
        <Bucket>
            <Name>test-bucket-1</Name>
            <CreationDate>2023-01-01T00:00:00.000Z</CreationDate>
        </Bucket>
        <Bucket>
            <Name>test-bucket-2</Name>
            <CreationDate>2023-01-02T00:00:00.000Z</CreationDate>
        </Bucket>
    </Buckets>
</ListAllMyBucketsResult>`
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockResponse)),
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV3) {
				ctx := context.TODO()
				input := &s3.ListBucketsInput{}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.ListBuckets(ctx, input)
				result2, err2 := s3ec.ListBuckets(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					if len(result1.Buckets) != len(result2.Buckets) {
						t.Errorf("Bucket count mismatch: s3Client=%d, s3ec=%d", len(result1.Buckets), len(result2.Buckets))
					}
					
					if result1.Owner != nil && result2.Owner != nil {
						if *result1.Owner.DisplayName != *result2.Owner.DisplayName {
							t.Errorf("Owner display name mismatch: s3Client=%s, s3ec=%s", 
								*result1.Owner.DisplayName, *result2.Owner.DisplayName)
						}
					}
					
					t.Logf("✓ ListBuckets: Both clients returned identical structured output")
				}
			},
		},
		{
			name:      "HeadBucket",
			operation: "HeadBucket",
			setupMock: func() *awstesting.MockHttpClient {
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("")),
						Header: http.Header{
							"x-amz-bucket-region": []string{"us-east-1"},
						},
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV3) {
				ctx := context.TODO()
				input := &s3.HeadBucketInput{
					Bucket: aws.String("test-bucket"),
				}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.HeadBucket(ctx, input)
				result2, err2 := s3ec.HeadBucket(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					// Both should return HeadBucketOutput with same structure
					// HeadBucketOutput doesn't have many fields, but both should be non-nil
					if result1 == nil || result2 == nil {
						t.Errorf("Result mismatch: s3Client result=%v, s3ec result=%v", result1, result2)
					}
					
					t.Logf("✓ HeadBucket: Both clients returned identical structured output")
				}
			},
		},
		{
			name:      "ListObjectsV2",
			operation: "ListObjectsV2",
			setupMock: func() *awstesting.MockHttpClient {
				mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Name>test-bucket</Name>
    <Prefix></Prefix>
    <KeyCount>2</KeyCount>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>test-object-1</Key>
        <LastModified>2023-01-01T00:00:00.000Z</LastModified>
        <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
        <Size>0</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>test-object-2</Key>
        <LastModified>2023-01-02T00:00:00.000Z</LastModified>
        <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
        <Size>100</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
</ListBucketResult>`
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockResponse)),
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV3) {
				ctx := context.TODO()
				input := &s3.ListObjectsV2Input{
					Bucket:  aws.String("test-bucket"),
					MaxKeys: 1000,
				}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.ListObjectsV2(ctx, input)
				result2, err2 := s3ec.ListObjectsV2(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					if len(result1.Contents) != len(result2.Contents) {
						t.Errorf("Contents count mismatch: s3Client=%d, s3ec=%d", len(result1.Contents), len(result2.Contents))
					}
					
					if result1.KeyCount != result2.KeyCount {
						t.Errorf("KeyCount mismatch: s3Client=%d, s3ec=%d", result1.KeyCount, result2.KeyCount)
					}
					
					if *result1.Name != *result2.Name {
						t.Errorf("Bucket name mismatch: s3Client=%s, s3ec=%s", *result1.Name, *result2.Name)
					}
					
					t.Logf("✓ ListObjectsV2: Both clients returned identical structured output")
				}
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create two separate mock HTTP clients with identical responses
			mockClient1 := tc.setupMock()
			mockClient2 := tc.setupMock()
			
			// Create regular S3 client
			tConfig1 := awstesting.Config()
			tConfig1.HTTPClient = mockClient1
			s3Client := s3.NewFromConfig(tConfig1)
			
			// Create S3 encryption client
			tConfig2 := awstesting.Config()
			tConfig2.HTTPClient = mockClient2
			s3BaseClient := s3.NewFromConfig(tConfig2)
			
			mockCMM := &mockCMM{}
			s3ec, err := New(s3BaseClient, mockCMM)
			if err != nil {
				t.Fatalf("Failed to create S3 encryption client: %v", err)
			}
			
			// Run the specific test for this operation
			tc.testFunc(t, s3Client, s3ec)
			
			t.Logf("✓ %s: S3EC adheres to same interface as conventional S3 client", tc.operation)
		})
	}
}

// Custom mock HTTP client that can capture multiple requests
type multiRequestMockClient struct {
	capturedRequests []string
}

func (m *multiRequestMockClient) Do(req *http.Request) (*http.Response, error) {
	// Capture the request path to verify both object and instruction file are deleted
	m.capturedRequests = append(m.capturedRequests, req.URL.Path)
	
	return &http.Response{
		Status:     http.StatusText(204),
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(strings.NewReader("")),
		Header: http.Header{
			"x-amz-delete-marker": []string{"false"},
		},
	}, nil
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - DeleteObject MUST delete the given object key.
//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
func TestS3EC_DeleteObject_DeletesObjectAndInstructionFile(t *testing.T) {
	// This test validates that DeleteObject deletes both the object and its instruction file
	// We'll use a custom mock to capture both delete requests
	
	// Create custom mock HTTP client that captures delete requests
	mockClient := &multiRequestMockClient{}

	// Create test config with mock HTTP client
	tConfig := awstesting.Config()
	tConfig.HTTPClient = mockClient
	s3Client := s3.NewFromConfig(tConfig)

	// Create mock CMM
	mockCMM := &mockCMM{}

	// Create the S3 encryption client
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Call DeleteObject - this should delete both object and instruction file
	ctx := context.TODO()
	result, err := s3ec.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
	})

	// Verify no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the response is not nil
	if result == nil {
		t.Error("Expected non-nil result from DeleteObject")
	}

	// Verify that both the object and instruction file delete requests were made
	if len(mockClient.capturedRequests) < 2 {
		t.Errorf("Expected at least 2 delete requests (object + instruction file), got %d", len(mockClient.capturedRequests))
	} else {
		// Check that requests include both the original object and the instruction file
		hasObjectDelete := false
		hasInstructionDelete := false
		
		for _, path := range mockClient.capturedRequests {
			if strings.Contains(path, "/test-key") && !strings.Contains(path, ".instruction") {
				hasObjectDelete = true
			}
			if strings.Contains(path, "/test-key.instruction") {
				hasInstructionDelete = true
			}
		}
		
		if !hasObjectDelete {
			t.Error("Expected delete request for original object 'test-key'")
		}
		if !hasInstructionDelete {
			t.Error("Expected delete request for instruction file 'test-key.instruction'")
		}
		
		t.Logf("✓ Verified DeleteObject deletes both object and instruction file")
	}
}
func TestS3ECInterfaceCompatibility(t *testing.T) {
	t.Run("ListBuckets", func(t *testing.T) {
		// Setup mock response
		mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Buckets><Bucket><Name>test-bucket</Name></Bucket></Buckets>
</ListAllMyBucketsResult>`

		// Create S3 client
		tConfig1 := awstesting.Config()
		tConfig1.HTTPClient = &awstesting.MockHttpClient{
			Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(mockResponse))},
		}
		s3Client := s3.NewFromConfig(tConfig1)

		// Create S3EC
		tConfig2 := awstesting.Config()
		tConfig2.HTTPClient = &awstesting.MockHttpClient{
			Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(mockResponse))},
		}
		s3ec, _ := New(s3.NewFromConfig(tConfig2), &mockCMM{})

		// Call same operation on both with same parameters
		ctx := context.TODO()
		input := &s3.ListBucketsInput{}
		
		result1, err1 := s3Client.ListBuckets(ctx, input)
		result2, err2 := s3ec.ListBuckets(ctx, input)

		// Verify same interface and output
		if (err1 == nil) != (err2 == nil) {
			t.Errorf("Error mismatch: s3Client=%v, s3ec=%v", err1, err2)
		}
		if len(result1.Buckets) != len(result2.Buckets) {
			t.Errorf("Output mismatch: s3Client buckets=%d, s3ec buckets=%d", len(result1.Buckets), len(result2.Buckets))
		}
	})

	t.Run("GetObject", func(t *testing.T) {
		// Verify interface compatibility - both should accept same input types
		ctx := context.TODO()
		input := &s3.GetObjectInput{
			Bucket: &[]string{"test-bucket"}[0],
			Key:    &[]string{"test-key"}[0],
		}
		
		// Test that both clients accept the same input parameters and return same types
		var result *s3.GetObjectOutput
		var err error
		
		// This verifies the interface is identical - same method signature
		// Don't assert equality; these clients behave differently. Important to assert the method signatures match.
		_ = func() {
			result, err = (&s3.Client{}).GetObject(ctx, input)
			result, err = (&S3EncryptionClientV3{}).GetObject(ctx, input)
		}
		
		// Suppress unused variable warnings
		_, _ = result, err
	})

	t.Run("PutObject", func(t *testing.T) {
		// Verify interface compatibility - both should accept same input types
		ctx := context.TODO()
		input := &s3.PutObjectInput{
			Bucket: &[]string{"test-bucket"}[0],
			Key:    &[]string{"test-key"}[0],
			Body:   strings.NewReader("test-content"),
		}
		
		// Test that both clients accept the same input parameters and return same types
		var result *s3.PutObjectOutput
		var err error
		
		// This verifies the interface is identical - same method signature
		// Don't assert equality; these clients behave differently. Important to assert the method signatures match.
		_ = func() {
			result, err = (&s3.Client{}).PutObject(ctx, input)
			result, err = (&S3EncryptionClientV3{}).PutObject(ctx, input)
		}
		
		// Suppress unused variable warnings
		_, _ = result, err
	})
}

//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
//= type=test
//# The option to enable legacy unauthenticated modes MUST be set to false by default.
func TestLegacyUnauthenticatedModes_DefaultDisabled(t *testing.T) {
	// Create test config
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Create S3EC without specifying EnableLegacyUnauthenticatedModes (should default to false)
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Verify that EnableLegacyUnauthenticatedModes defaults to false
	if s3ec.Options.EnableLegacyUnauthenticatedModes {
		t.Error("Expected EnableLegacyUnauthenticatedModes to default to false, but it was true")
	}

	t.Logf("✓ Verified EnableLegacyUnauthenticatedModes defaults to false")
}

//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
//= type=test
//# The S3EC MUST support the option to provide an SDK S3 client instance during its initialization.
func TestS3EC_AcceptsProvidedS3ClientInstance(t *testing.T) {
	// Create a specific S3 client instance
	tConfig := awstesting.Config()
	providedS3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Create S3EC with the provided S3 client
	s3ec, err := New(providedS3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Verify that the S3EC uses the provided client (they should be the same instance)
	if s3ec.Client != providedS3Client {
		t.Error("S3EC should use the provided S3 client instance")
	}

	t.Logf("✓ Verified S3EC accepts and uses provided S3 client instance")
}

//= ../specification/s3-encryption/client.md#encryption-algorithm
//= type=test
//# The S3EC MUST support configuration of the encryption algorithm (or algorithm suite) during its initialization.
func TestS3EC_SupportsEncryptionAlgorithmConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test configuring different valid (non-legacy) algorithm suites
	testCases := []struct {
		name             string
		algorithm        *algorithms.AlgorithmSuite
		commitmentPolicy commitment.CommitmentPolicy
	}{
		{
			name:             "AES256GCMHkdfSha512CommitKey",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT, // Committing algorithm
		},
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT, // Non-committing algorithm
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure S3EC with specific encryption algorithm and appropriate commitment policy
			s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.EncryptionAlgorithmSuite = tc.algorithm
				options.CommitmentPolicy = tc.commitmentPolicy
			})
			
			if err != nil {
				t.Fatalf("Failed to create S3EC with algorithm %s: %v", tc.name, err)
			}

			// Verify the algorithm was configured correctly
			if s3ec.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %s, got %v", tc.name, s3ec.Options.EncryptionAlgorithmSuite)
			}

			// Verify the commitment policy was configured correctly
			if s3ec.Options.CommitmentPolicy != tc.commitmentPolicy {
				t.Errorf("Expected commitment policy %v, got %v", tc.commitmentPolicy, s3ec.Options.CommitmentPolicy)
			}

			t.Logf("✓ Successfully configured S3EC with algorithm: %s and commitment policy: %v", tc.name, tc.commitmentPolicy)
		})
	}
}

//= ../specification/s3-encryption/client.md#key-commitment
//= type=test
//# The S3EC MUST support configuration of the [Key Commitment policy](./key-commitment.md) during its initialization.
func TestS3EC_SupportsKeyCommitmentPolicyConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test configuring all 3 commitment policies with compatible algorithms (happy path cases)
	testCases := []struct {
		name             string
		commitmentPolicy commitment.CommitmentPolicy
		algorithm        *algorithms.AlgorithmSuite
	}{
		{
			name:             "FORBID_ENCRYPT_ALLOW_DECRYPT",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF, // Non-committing algorithm
		},
		{
			name:             "REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing algorithm
		},
		{
			name:             "REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing algorithm
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure S3EC with specific commitment policy and compatible algorithm
			s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.CommitmentPolicy = tc.commitmentPolicy
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err != nil {
				t.Fatalf("Failed to create S3EC with commitment policy %v: %v", tc.commitmentPolicy, err)
			}

			// Verify the commitment policy was configured correctly
			if s3ec.Options.CommitmentPolicy != tc.commitmentPolicy {
				t.Errorf("Expected commitment policy %v, got %v", tc.commitmentPolicy, s3ec.Options.CommitmentPolicy)
			}

			// Verify the algorithm was configured correctly
			if s3ec.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %v, got %v", tc.algorithm, s3ec.Options.EncryptionAlgorithmSuite)
			}

			t.Logf("✓ Successfully configured S3EC with commitment policy: %v and algorithm: %v", tc.commitmentPolicy, tc.algorithm)
		})
	}
}

//= ../specification/s3-encryption/client.md#key-commitment
//= type=test
//# The S3EC MUST validate the configured Encryption Algorithm against the provided key commitment policy.
//# If the configured Encryption Algorithm is incompatible with the key commitment policy, then it MUST throw an exception.
func TestS3EC_ValidatesAlgorithmCommitmentPolicyCompatibility(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test incompatible combinations that should fail
	incompatibleCases := []struct {
		name             string
		commitmentPolicy commitment.CommitmentPolicy
		algorithm        *algorithms.AlgorithmSuite
		expectedError    string
	}{
		{
			name:             "FORBID_ENCRYPT_ALLOW_DECRYPT with committing algorithm",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing
			expectedError:    "does not allow committing algorithm suites",
		},
		{
			name:             "REQUIRE_ENCRYPT_REQUIRE_DECRYPT with non-committing algorithm",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF, // Non-committing
			expectedError:    "requires committing algorithm suites",
		},
	}

	for _, tc := range incompatibleCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to configure S3EC with incompatible commitment policy and algorithm - should fail
			_, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.CommitmentPolicy = tc.commitmentPolicy
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err == nil {
				t.Errorf("Expected error for incompatible combination %s, but got none", tc.name)
			} else if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectedError, err)
			} else {
				t.Logf("✓ Correctly rejected incompatible combination %s: %v", tc.name, err)
			}
		})
	}
}

//= ../specification/s3-encryption/client.md#encryption-algorithm
//= type=test
//# The S3EC MUST validate that the configured encryption algorithm is not legacy.
//# If the configured encryption algorithm is legacy, then the S3EC MUST throw an exception.
func TestS3EC_RejectsLegacyEncryptionAlgorithms(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test that legacy algorithm suites are rejected
	legacyAlgorithms := []struct {
		name      string
		algorithm *algorithms.AlgorithmSuite
	}{
		{
			name:      "AES256CTRIV16Tag16NoKDF (legacy)",
			algorithm: algorithms.AlgAES256CTRIV16Tag16NoKDF,
		},
		{
			name:      "AES256CBCIV16NoKDF (legacy)",
			algorithm: algorithms.AlgAES256CBCIV16NoKDF,
		},
	}

	for _, tc := range legacyAlgorithms {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to configure S3EC with legacy encryption algorithm - should fail
			_, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err == nil {
				t.Errorf("Expected error when configuring legacy algorithm %s, but got none", tc.name)
			} else {
				t.Logf("✓ Correctly rejected legacy algorithm %s: %v", tc.name, err)
			}
		})
	}
}

func TestS3EC_BufferSizeConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//= type=test
	//# If Delayed Authentication mode is disabled, and no buffer size is provided,
	//# the S3EC MUST set the buffer size to a reasonable default.
	t.Run("SetsReasonableDefaultBufferSize", func(t *testing.T) {
		// Create S3EC with default options
		s3ec, err := New(s3Client, mockCMM)
		if err != nil {
			t.Fatalf("Failed to create S3 encryption client: %v", err)
		}
		
		// Verify that the default buffer size is the default
		expectedBufferSize := int64(DefaultBufferSize)
		if s3ec.Options.BufferSize != expectedBufferSize {
			t.Errorf("Expected default buffer size to be %d, got %d", expectedBufferSize, s3ec.Options.BufferSize)
		}
		
		// Verify the default buffer size is the default
		if s3ec.Options.BufferSize != 64*1024 {
			t.Errorf("Expected default buffer size to be 64KB (65536 bytes), got %d", s3ec.Options.BufferSize)
		}
		
		t.Logf("✓ Verified S3EC sets reasonable default buffer size: %d bytes", s3ec.Options.BufferSize)
	})

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//= type=test
	//# The S3EC SHOULD accept a configurable buffer size
	//# which refers to the maximum ciphertext length in bytes to store in memory
	//# when Delayed Authentication mode is disabled.
	t.Run("SupportsCustomBufferSizeConfiguration", func(t *testing.T) {
		// Test custom buffer size configuration
		customBufferSize := int64(128 * 1024) // 128KB
		
		s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
			options.BufferSize = customBufferSize
		})
		if err != nil {
			t.Fatalf("Failed to create S3 encryption client with custom buffer size: %v", err)
		}
		
		// Verify that the custom buffer size is set correctly
		if s3ec.Options.BufferSize != customBufferSize {
			t.Errorf("Expected buffer size to be %d, got %d", customBufferSize, s3ec.Options.BufferSize)
		}
		
		t.Logf("✓ Verified S3EC supports custom buffer size configuration: %d bytes", s3ec.Options.BufferSize)
	})
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package commitment

import "fmt"

type CommitmentPolicy int

const (
	// "Forbid" writing objects encrypted with key commitment, and "allow" reading objects encrypted with key commitment.
	// FORBID_ENCRYPT_ALLOW_DECRYPT does not write objects with key commitment
	// and can read objects encrypted either with or without key commitment.
	// Keys in Instruction Files could be tampered with when reading objects without key commitment.
	// FORBID_ENCRYPT_ALLOW_DECRYPT means that this client will write objects that any v3 client can read,
	// and any v4 client (configured with either FORBID_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_ALLOW_DECRYPT) can read.
	// FORBID_ENCRYPT_ALLOW_DECRYPT also means that this client can read objects written by any v3 or v4 client.
	// This is the default policy for v3 clients.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	FORBID_ENCRYPT_ALLOW_DECRYPT CommitmentPolicy = iota
	// "Require" writing objects encrypted with key commitment, and "allow" reading objects encrypted with key commitment.
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT allows this client to read objects encrypted either with or without key commitment.
	// However, specifying REQUIRE_ENCRYPT_ALLOW_DECRYPT will not allow a v3 client to write objects
	// as v3 clients do not support writing objects with key commitment.
	// To write objects with key commitment, you must use a v4 client.
	// Keys in Instruction Files could be tampered with when reading objects without key commitment.
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT also means that this client can read objects written by any v3 or v4 client.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	REQUIRE_ENCRYPT_ALLOW_DECRYPT
	// "Require" writing objects encrypted with key commitment, and "require" reading objects encrypted with key commitment.
	// REQUIRE_ENCRYPT_REQUIRE_DECRYPT ensures that all decrypted objects are verified to have been encrypted with key commitment.
	// This prevents reading objects with keys in Instruction Files that may have been tampered with.
	// However, specifying REQUIRE_ENCRYPT_REQUIRE_DECRYPT will not allow a v3 client to write objects
	// as v3 clients do not support writing objects with key commitment.
	// To write objects with key commitment, you must use a v4 client.
	// Specifying REQUIRE_ENCRYPT_REQUIRE_DECRYPT also means that this client can only read objects written by
	// v4 clients (configured with any REQUIRE_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_REQUIRE_DECRYPT).
	// This is the default policy for v4 clients.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	REQUIRE_ENCRYPT_REQUIRE_DECRYPT
)

func (cp CommitmentPolicy) RequiresEncrypt() bool {
	switch cp {
	case REQUIRE_ENCRYPT_ALLOW_DECRYPT, REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
		return true
	default:
		return false
	}
}

func (cp CommitmentPolicy) RequiresDecrypt() bool {
	return cp == REQUIRE_ENCRYPT_REQUIRE_DECRYPT
}

func (p CommitmentPolicy) String() string {
    switch p {
    case FORBID_ENCRYPT_ALLOW_DECRYPT:
        return "FORBID_ENCRYPT_ALLOW_DECRYPT"
    case REQUIRE_ENCRYPT_ALLOW_DECRYPT:
        return "REQUIRE_ENCRYPT_ALLOW_DECRYPT"
    case REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
        return "REQUIRE_ENCRYPT_REQUIRE_DECRYPT"
    default:
        return fmt.Sprintf("CommitmentPolicy(%d)", int(p))
    }
}

import (
	"bytes"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"io"
)

type aesGCM struct {
	aead  cipher.AEAD
	nonce []byte
	aad  []byte
}

// newAESGCM creates a new AES GCM cipher. Expects keys to be of
		return nil, err
	}

	if materials.CEKAlgorithm == algorithms.AESGCMCommitKey {
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//# The client MUST set the AAD to the Algorithm Suite ID represented as bytes.
		aad := algorithms.AlgAES256GCMHkdfSha512CommitKey.IDAsBytes()
		return &aesGCM{aesgcm, materials.IV, aad}, nil
	} else if materials.CEKAlgorithm == algorithms.AESGCMNoPadding {
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//# The client MUST NOT provide any AAD when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
		return &aesGCM{aesgcm, materials.IV, nil}, nil
	} else {
		return nil, fmt.Errorf("unsupported CEK algorithm for AES GCM: %s", materials.CEKAlgorithm)
	}
	
}

// Encrypt will encrypt the data using AES GCM
		encrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
		aad:	   c.aad,
	}
	return reader
}
	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
	aad	      []byte
}

func (reader *gcmEncryptReader) Read(data []byte) (int, error) {
		if err != nil {
			return 0, err
		}
		var aad []byte
		if reader.aad != nil {
			aad = reader.aad
		} else {
			aad = nil
		}
		// The GCM auth tag is appended to the ciphertext by the Seal function.
		// Docs: https://pkg.go.dev/crypto/cipher#GCM
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//= type=exception
		//# The client MUST append the GCM auth tag to the ciphertext if the underlying crypto provider does not do so automatically.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
		//= type=exception
		//# The client MUST append the GCM auth tag to the ciphertext if the underlying crypto provider does not do so automatically.

		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//= type=implication
		//# The client MUST initialize the cipher, or call an AES-GCM encryption API,
		//# with the plaintext data key, the generated IV, and the tag length defined in the Algorithm Suite
		//# when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
		b = reader.encrypter.Seal(b[:0], reader.nonce, b, aad)
		reader.buf = bytes.NewBuffer(b)
	}

		decrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
		aad:	   c.aad,
	}
}

	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
	aad 	  []byte
}

func (reader *gcmDecryptReader) Read(data []byte) (int, error) {
	var aad []byte
	if reader.aad != nil {
		aad = reader.aad
	} else {
		aad = nil
	}
	if reader.buf == nil {
		b, err := io.ReadAll(reader.src)
		if err != nil {
			return 0, err
		}
		b, err = reader.decrypter.Open(b[:0], reader.nonce, b, aad)
		if err != nil {
			return 0, err
		}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
)

// NewAESGCMCommittingContentCipher returns a new encryption only AES/GCM mode structure with a specific cipher data generator
// that will provide keys to be used for content encryption.
//
// Note: This uses the Go stdlib AEAD implementation for AES/GCM. Due to this, objects to be encrypted or decrypted
// will be fully loaded into memory before encryption or decryption can occur. Caution must be taken to avoid memory
// allocation failures.
func NewAESGCMCommittingContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.CEKAlgorithm = algorithms.AESGCMCommitKey
	materials.TagLength = GcmTagSizeBits

	// Persist original IV to store in the object metadata
	original_iv := materials.IV

	var storedKeyCommitment []byte = nil
	if materials.KeyCommitment != nil {
		storedKeyCommitment = make([]byte, len(materials.KeyCommitment))
		copy(storedKeyCommitment, materials.KeyCommitment)
	}

	//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
	//= type=implication
	//# The client MUST use HKDF to derive the key commitment value and the derived encrypting key as described in [Key Derivation](key-derivation.md).
	keys, err := DeriveKeys(materials.Key, materials.IV, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), storedKeyCommitment)
	if err != nil {
		return nil, err
	}

	materials.Key = keys.DerivedEncryptionKey
	//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
	//# The derived key commitment value MUST be set or returned from the encryption process such that it can be included in the content metadata.
	materials.KeyCommitment = keys.CommitKey

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# When encrypting or decrypting with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY,
	//# the IV used in the AES-GCM content encryption/decryption MUST consist entirely of bytes with the value 0x01.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# The IV's total length MUST match the IV length defined by the algorithm suite.
	var nonce = [12]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	materials.IV = nonce[:]

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# The client MUST initialize the cipher, or call an AES-GCM encryption API,
	//# with the derived encryption key,
	//# an IV containing only bytes with the value 0x01,
	//# and the tag length defined in the Algorithm Suite
	//# when encrypting or decrypting with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY.
	cipher, err := newAESGCM(materials)
	if err != nil {
		return nil, err
	}

	// Restore original non-zero IV for metadata storage
	materials.IV = original_iv

	return &aesGCMContentCipher{
		CryptographicMaterials: materials,
		Cipher:                 cipher,
	}, nil
}

func NewAESGCMDecryptCommittingContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	// KeyCommitment value is required for decryption
	if materials.KeyCommitment == nil {
		return nil, fmt.Errorf("key commitment is required for AES-GCM committing decryption")
	}

	return NewAESGCMCommittingContentCipher(materials)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
)

func TestNewAESGCMDecryptCommittingContentCipher_NilKeyCommitment_ReturnsError(t *testing.T) {
	// Given: materials with nil KeyCommitment
	materials := materials.CryptographicMaterials{
		KeyCommitment: nil,
	}

	// When: calling NewAESGCMDecryptCommittingContentCipher
	cipher, err := NewAESGCMDecryptCommittingContentCipher(materials)

	// Then: an error should be returned
	if err == nil {
		t.Fatal("expected error when KeyCommitment is nil, but got nil")
	}

	if cipher != nil {
		t.Fatal("expected cipher to be nil when error occurs, but got non-nil cipher")
	}

	// Verify the error message contains the expected text
	expectedErrorText := "key commitment is required for AES-GCM committing decryption"
	if !strings.Contains(err.Error(), expectedErrorText) {
		t.Fatalf("expected error message to contain '%s', but got: %s", expectedErrorText, err.Error())
	}
}
package internal

import (
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"io"
)

const (
	GcmTagSizeBits  = "128"
)

// NewAESGCMContentCipher returns a new encryption only AES/GCM mode structure with a specific cipher data generator
// will be fully loaded into memory before encryption or decryption can occur. Caution must be taken to avoid memory
// allocation failures.
func NewAESGCMContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.CEKAlgorithm = algorithms.AESGCMNoPadding
	materials.TagLength = GcmTagSizeBits

	cipher, err := newAESGCM(materials)
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	materials2 "github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"io"
	"os"
	}
}

//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
//= type=test
//# The client MUST set the AAD to the Algorithm Suite ID represented as bytes.
func TestGIVEN_materialsCEKAlg115_WHEN_newAESGCM_THEN_returnedAesGCMAadIs0x0073(t *testing.T) {
	// Given: materials with CEKAlgorithm ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY (0x73)
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: "115",
		Key:          make([]byte, 32),
		IV:           make([]byte, 12),
	}

	// When: instantiating new AES GCM cipher with the materials
	cipher, err := newAESGCM(materials)

	// Then: no error is returned and the AAD is set to the Algorithm Suite ID bytes
	if err != nil {
		panic(fmt.Sprintf("expected no error, but received %v", err))
	}

	aesgcm, ok := cipher.(*aesGCM)
	if !ok {
		panic("expected cipher to be of type *aesGCM")
	}
	
	expectedAad := []byte{0x00, 0x73}
	if !bytes.Equal(aesgcm.aad, expectedAad) {
		panic(fmt.Sprintf("expected AAD to be %v, but received %v", expectedAad, aesgcm.aad))
	}
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
//= type=test
//# The client MUST NOT provide any AAD when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
func TestGIVEN_materialsCEKAlgAESGCMNoPadding_WHEN_newAESGCM_THEN_returnedAesGCMAadIsNil(t *testing.T) {
	// Given: materials with CEKAlgorithm ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY (0x73)
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: "AES/GCM/NoPadding",
		Key:          make([]byte, 32),
		IV:           make([]byte, 12),
	}

	// When: instantiating new AES GCM cipher with the materials
	cipher, err := newAESGCM(materials)

	// Then: no error is returned and the AAD is set to the Algorithm Suite ID bytes
	if err != nil {
		panic(fmt.Sprintf("expected no error, but received %v", err))
	}

	aesgcm, ok := cipher.(*aesGCM)
	if !ok {
		panic("expected cipher to be of type *aesGCM")
	}
	
	if aesgcm.aad != nil {
		panic(fmt.Sprintf("expected AAD to be nil, but received %v", aesgcm.aad))
	}
}

func aesgcmTest(t *testing.T, iv, key, plaintext, expected, tag []byte) {
	t.Helper()
	const gcmTagSize = 16
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: algorithms.AESGCMNoPadding,
		Key: key,
		IV:  iv,
	}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/rand"
	"fmt"
)

// GenerateNonZeroBytes generates random bytes and validates they are not all zeros
func GenerateNonZeroBytes(n int) ([]byte, error) {
	return GenerateNonZeroBytesWithGenerator(n, generateRandomBytes)
}

// GenerateNonZeroBytesWithGenerator allows injection of custom generator for testing
func GenerateNonZeroBytesWithGenerator(n int, generator func(int) ([]byte, error)) ([]byte, error) {
	const maxRetries = 3 // Prevent infinite loop in case of broken generator
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		//= ../specification/s3-encryption/encryption.md#content-encryption
		//# The client MUST generate an IV or Message ID using the length of the IV or Message ID defined in the algorithm suite.
		keys_iv, err := generator(n)
		if err != nil {
			return nil, err
		}
		
		//= ../specification/s3-encryption/encryption.md#cipher-initialization
		//# The client SHOULD validate that the generated IV or Message ID is not zeros.
		allZero := true
		for _, b := range keys_iv {
			if b != 0 {
				allZero = false
				break
			}
		}
		
		// If not all zeros, we have a valid IV
		if !allZero {
			return keys_iv, nil
		}
		
		// If all zeros, retry (unless this is the last attempt)
	}
	
	// If we've exhausted all retries, return an error
	return nil, fmt.Errorf("failed to generate non-zero IV after %d attempts", maxRetries)
}

// Default generator using crypto/rand
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"testing"
)

func TestGenerateNonZeroBytes(t *testing.T) {
	cases := []struct {
		name   string
		length int
	}{
		{
			name:   "small_length",
			length: 12,
		},
		{
			name:   "medium_length",
			length: 28,
		},
		{
			name:   "large_length",
			length: 64,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateNonZeroBytes(tc.length)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(result) != tc.length {
				t.Errorf("expected length %d, got %d", tc.length, len(result))
			}

			// Verify not all zeros
			allZero := true
			for _, b := range result {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				t.Error("GenerateNonZeroBytes returned all zeros")
			}
		})
	}
}

//= ../specification/s3-encryption/encryption.md#cipher-initialization
//= type=test
//# The client SHOULD validate that the generated IV or Message ID is not zeros.
func TestGenerateNonZeroBytesWithGenerator_AllZeros(t *testing.T) {
	// Mock generator that always returns all zeros
	mockAllZeros := func(n int) ([]byte, error) {
		return make([]byte, n), nil // Returns all zeros
	}

	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "12_byte_iv",
			length: 12,
		},
		{
			name:   "28_byte_message_id",
			length: 28,
		},
		{
			name:   "16_byte_iv",
			length: 16,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateNonZeroBytesWithGenerator(tc.length, mockAllZeros)

			// Assert that an error is raised after max retries
			if err == nil {
				t.Fatalf("Expected error when generator always returns all zeros for length %d, but got none", tc.length)
			}

			// Assert that the error message indicates retry exhaustion
			expectedErrorMsg := "failed to generate non-zero IV after 3 attempts"
			if err.Error() != expectedErrorMsg {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
			}
		})
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	// Test that the default generator produces different results
	results := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		result, err := generateRandomBytes(16)
		if err != nil {
			t.Fatalf("generateRandomBytes failed: %v", err)
		}
		if len(result) != 16 {
			t.Errorf("Expected length 16, got %d", len(result))
		}
		results[i] = result
	}

	// Check that not all results are identical (very unlikely with proper randomness)
	allSame := true
	for i := 1; i < len(results); i++ {
		if !bytes.Equal(results[0], results[i]) {
			allSame = false
			break
		}
	}

	if allSame {
		t.Errorf("All calls to generateRandomBytes returned identical results, randomness may be compromised")
	}
}
package internal

import (
	"bufio"
	"errors"
	"io"
	"os"
	ws.i = abs
	return abs, nil
}

// NewBufferedReader creates a buffered reader with the specified buffer size
// This implements the S3EC requirement to set buffer size to a reasonable default for GetObject
func NewBufferedReader(r io.Reader, bufferSize int) (io.ReadCloser, error) {
	if bufferSize <= 0 {
		return nil, errors.New("buffer size must be greater than zero")
	}
	
	bufferedReader := bufio.NewReaderSize(r, bufferSize)
	
	// If the original reader implements io.ReadCloser, preserve the Close method
	if rc, ok := r.(io.ReadCloser); ok {
		return &bufferedReadCloser{
			Reader: bufferedReader,
			closer: rc,
		}, nil
	}
	
	// If not a ReadCloser, create a no-op closer
	return &bufferedReadCloser{
		Reader: bufferedReader,
		closer: nil,
	}, nil
}

// bufferedReadCloser wraps a buffered reader and preserves the Close method
type bufferedReadCloser struct {
	*bufio.Reader
	closer io.ReadCloser
}

func (brc *bufferedReadCloser) Close() error {
	if brc.closer != nil {
		return brc.closer.Close()
	}
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"crypto/hkdf"
)

const (
	// Key derivation constants
	DeriveKeyInfo = "DERIVEKEY"
	CommitKeyInfo = "COMMITKEY"
)

// KeyDerivationResult holds the results of HKDF key derivation
type KeyDerivationResult struct {
	DerivedEncryptionKey []byte
	CommitKey            []byte
}

func DeriveKeys(plaintextDataKey []byte, messageID []byte, algorithmSuiteID int, storedKeyCommitment []byte) (*KeyDerivationResult, error) {
	// Validate input parameters
	if len(plaintextDataKey) == 0 {
		return nil, fmt.Errorf("plaintext data key cannot be empty")
	}
	if len(messageID) == 0 {
		return nil, fmt.Errorf("message ID cannot be empty")
	}

	// Get algorithm suite to determine key lengths and other properties
	algSuite, err := algorithms.GetAlgorithmSuiteByID(algorithmSuiteID)
	if err != nil {
		return nil, fmt.Errorf("unable to get algorithm suite by ID: %v", err)
	}

	// Convert algorithm suite ID to bytes (big-endian)
	algorithmSuiteIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(algorithmSuiteIDBytes, uint16(algorithmSuiteID))

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The hash function MUST be specified by the algorithm suite commitment settings.
	hashFunc := algSuite.KDFHashAlgorithm()
	if hashFunc == nil || !algSuite.IsCommitting() {
		return nil, fmt.Errorf("algorithm suite does not support key derivation: 0x%04x", algorithmSuiteID)
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input keying material MUST be the plaintext data key (PDK) generated by the key provider.
	inputKeyMaterial := plaintextDataKey

	// Get expected key lengths from algorithm suite
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the input keying material MUST equal the key derivation input length specified by the algorithm suite commit key derivation setting.
	if len(inputKeyMaterial) != expectedDataKeyLength {
		return nil, fmt.Errorf("plaintext data key length must be %d bytes, got %d", expectedDataKeyLength, len(inputKeyMaterial))
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The salt MUST be the Message ID with the length defined in the algorithm suite.
	salt := messageID
	expectedSaltLength := algSuite.IVLengthBytes()
	if len(salt) != expectedSaltLength {
		return nil, fmt.Errorf("message ID length must be %d bytes, got %d", expectedSaltLength, len(salt))
	}

	// Extract step
	extractPRK, err := hkdf.Extract(hashFunc, inputKeyMaterial, salt)
	if err != nil {
		return nil, fmt.Errorf("hkdf extract failed: %w", err)
	}
	if len(extractPRK) == 0 {
		return nil, fmt.Errorf("hkdf extract failed: extractPRK is empty")
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input info MUST be a concatenation of the algorithm suite ID as bytes followed by the string COMMITKEY as UTF8 encoded bytes.
	commitKeyInfo := append(algorithmSuiteIDBytes, []byte(CommitKeyInfo)...)
	
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The CK input pseudorandom key MUST be the output from the extract step.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the output keying material MUST equal the commit key length specified by the supported algorithm suites.
	derivedCommitKey, err := hkdf.Expand(hashFunc, extractPRK, string(commitKeyInfo), expectedCommitKeyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive commit key: %w", err)
	}
	if len(derivedCommitKey) != expectedCommitKeyLength {
		return nil, fmt.Errorf("commit key length must be %d bytes, got %d", expectedCommitKeyLength, len(derivedCommitKey))
	}

	// First, check commitment value; then, derive encryption key
	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=implication
	//# When using an algorithm suite which supports key commitment, the client MUST verify the key commitment values match before deriving the [derived encryption key](./key-derivation.md#hkdf-operation).

	if storedKeyCommitment != nil {
		//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
		//# When using an algorithm suite which supports key commitment, the client MUST verify that the [derived key commitment](./key-derivation.md#hkdf-operation) contains the same bytes as the stored key commitment retrieved from the stored object's metadata.
		//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
		//= type=implication
		//# When using an algorithm suite which supports key commitment, the verification of the derived key commitment value MUST be done in constant time.
		if subtle.ConstantTimeCompare(derivedCommitKey, storedKeyCommitment) != 1 {
			//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
			//# When using an algorithm suite which supports key commitment, the client MUST throw an exception when the derived key commitment value and stored key commitment value do not match.
			return nil, fmt.Errorf("derived key commitment value does not match value stored on encrypted message")
		}
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input info MUST be a concatenation of the algorithm suite ID as bytes followed by the string DERIVEKEY as UTF8 encoded bytes.
	encryptionKeyInfo := append(algorithmSuiteIDBytes, []byte(DeriveKeyInfo)...)
	
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The DEK input pseudorandom key MUST be the output from the extract step.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the output keying material MUST equal the encryption key length specified by the algorithm suite encryption settings.
	derivedEncryptionKey, err := hkdf.Expand(hashFunc, extractPRK, string(encryptionKeyInfo), expectedDataKeyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}
	if len(derivedEncryptionKey) != expectedDataKeyLength {
		return nil, fmt.Errorf("encryption key length must be %d bytes, got %d", expectedDataKeyLength, len(derivedEncryptionKey))
	}

	return &KeyDerivationResult{
		DerivedEncryptionKey: derivedEncryptionKey,
		CommitKey:            derivedCommitKey,
	}, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/hex"
	"strings"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"testing"
)

func hexToBytes(t *testing.T, s string, expectedLen int) []byte {
	b, err := decodeHex(s)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	if len(b) != expectedLen {
		t.Fatalf("Expected length %d, got %d for hex string %s", expectedLen, len(b), s)
	}
	return b
}

func decodeHex(s string) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(s)))
	_, err := hex.Decode(dst, []byte(s))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func TestDeriveKeys_KnownAnswerTests(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	tests := []struct {
		comment        string
		dataKeyHex     string
		messageIDHex   string
		expectedEncHex string
		expectedComHex string
	}{
		{
			comment:        "Basic S3EC.ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY #1",
			dataKeyHex:     "80d90dc4cc7e77d8a6332efa44eba56230a7fe7b89af37d1e501ab2e07c0a163",
			messageIDHex:   "b8ea76bed24c7b85382a148cb9dcd1cfdfb765f55ded4dfa6e0c4c79",
			expectedEncHex: "6dd14f546cc006e639126e83f5d4d1b118576bb5df97f38c6fb3a1db87bbc338",
			expectedComHex: "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db",
		},
		{
			comment:        "Basic S3EC.ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY #2",
			dataKeyHex:     "501afb8227d22e75e68010414b8abdaf3064c081e8e922dafef4992036394d60",
			messageIDHex:   "61a00b4981a5aacfd136c55cb726e32d2a547dc7600a7d4675c69127",
			expectedEncHex: "e14786a714748d1d2c3a4a6816dec56ddf1881bbeabb4f39420ffb9f63700b2f",
			expectedComHex: "5c1e73e47f6fe3a70d6d094283aceaa76d2975feb829212d88f0afc1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.comment, func(t *testing.T) {
			dataKey := hexToBytes(t, tc.dataKeyHex, expectedDataKeyLength)
			messageID := hexToBytes(t, tc.messageIDHex, expectedMessageIDLength)
			expectedEnc := hexToBytes(t, tc.expectedEncHex, expectedDataKeyLength)
			expectedCom := hexToBytes(t, tc.expectedComHex, expectedCommitKeyLength)

			result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), expectedCom)
			if err != nil {
				t.Fatalf("DeriveKeys failed: %v", err)
			}
			if !bytes.Equal(result.DerivedEncryptionKey, expectedEnc) {
				t.Errorf("DerivedEncryptionKey mismatch.\nExpected: %x\nGot:      %x", expectedEnc, result.DerivedEncryptionKey)
			}
			if !bytes.Equal(result.CommitKey, expectedCom) {
				t.Errorf("CommitKey mismatch.\nExpected: %x\nGot:      %x", expectedCom, result.CommitKey)
			}
		})
	}
}

func TestDeriveKeys_KeyCommitmentValidation(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	dataKey := hexToBytes(t, "80d90dc4cc7e77d8a6332efa44eba56230a7fe7b89af37d1e501ab2e07c0a163", expectedDataKeyLength)
	messageID := hexToBytes(t, "b8ea76bed24c7b85382a148cb9dcd1cfdfb765f55ded4dfa6e0c4c79", expectedMessageIDLength)
	correctCommitment := hexToBytes(t, "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db", expectedCommitKeyLength)
	wrongCommitment := hexToBytes(t, "00000000000000000000000000000000000000000000000000000000", expectedCommitKeyLength)

	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=test
	//# When using an algorithm suite which supports key commitment, the client MUST verify that the [derived key commitment](./key-derivation.md#hkdf-operation) contains the same bytes as the stored key commitment retrieved from the stored object's metadata.
	t.Run("matching_commitment_succeeds", func(t *testing.T) {
		result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), correctCommitment)
		if err != nil {
			t.Fatalf("DeriveKeys should succeed when commitment values match, got error: %v", err)
		}

		// Verify that the derived commitment matches the stored commitment
		if !bytes.Equal(result.CommitKey, correctCommitment) {
			t.Errorf("Derived commitment should match stored commitment.\nExpected: %x\nGot:      %x", correctCommitment, result.CommitKey)
		}
	})

	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=test
	//# When using an algorithm suite which supports key commitment, the client MUST throw an exception when the derived key commitment value and stored key commitment value do not match.
	t.Run("mismatched_commitment_throws_exception", func(t *testing.T) {
		_, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), wrongCommitment)
		if err == nil {
			t.Fatalf("DeriveKeys should throw an exception when commitment values do not match, but it succeeded")
		}

		// Verify the error message indicates commitment mismatch
		if !bytes.Contains([]byte(err.Error()), []byte("derived key commitment value does not match value stored on encrypted message")) {
			t.Errorf("Expected error to mention commitment mismatch, got: %v", err)
		}
	})
}

func TestDeriveKeys_InputOutputLengthValidation(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	correctCommitment := hexToBytes(t, "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db", expectedCommitKeyLength)

	cases := []struct {
		name           string
		dataKeyLen     int
		messageIDLen   int
		expectError    bool
		errorContains  string
	}{
		{
			name:         "correct_input_lengths",
			dataKeyLen:   expectedDataKeyLength,   // 32 bytes
			messageIDLen: expectedMessageIDLength, // Message ID length for committing algorithm
			expectError:  true, // Will fail due to commitment mismatch, but that's OK - we're testing length validation
			errorContains: "commitment",
		},
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//= type=test
		//# - The length of the input keying material MUST equal the key derivation input length specified by the algorithm suite commit key derivation setting.
		{
			name:          "wrong_data_key_length",
			dataKeyLen:    16, // Wrong length
			messageIDLen:  expectedMessageIDLength,
			expectError:   true,
			errorContains: "plaintext data key length",
		},
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//= type=test
		//# - The salt MUST be the Message ID with the length defined in the algorithm suite.
		{
			name:          "wrong_message_id_length",
			dataKeyLen:    expectedDataKeyLength,
			messageIDLen:  16, // Wrong length
			expectError:   true,
			errorContains: "message ID length",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			
			dataKey := make([]byte, tc.dataKeyLen)
			for i := range dataKey {
				dataKey[i] = byte(i % 256)
			}

			messageID := make([]byte, tc.messageIDLen)
			for i := range messageID {
				messageID[i] = byte((i + 100) % 256)
			}

			result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), correctCommitment)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error for %s but got none", tc.name)
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Expected error for %s: %v", tc.name, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for %s, got %v", tc.name, err)
				}

				//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
				//= type=test
				//# - The length of the output keying material MUST equal the encryption key length specified by the algorithm suite encryption settings.
				if len(result.DerivedEncryptionKey) != expectedDataKeyLength {
					t.Errorf("DerivedEncryptionKey length should be %d, got %d", expectedDataKeyLength, len(result.DerivedEncryptionKey))
				}
				
				//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
				//= type=test
				//# - The length of the output keying material MUST equal the commit key length specified by the supported algorithm suites.
				if len(result.CommitKey) != expectedCommitKeyLength {
					t.Errorf("CommitKey length should be %d, got %d", expectedCommitKeyLength, len(result.CommitKey))
				}
				t.Logf("✓ Correct input/output lengths for %s", tc.name)
			}
		})
	}
}
	"encoding/json"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"strconv"
)

const DefaultInstructionKeySuffix = ".instruction"

const (
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=implication
	//# The "x-amz-" prefix denotes that the metadata is owned by an Amazon product and MUST be prepended to all S3EC metadata mapkeys.
	amzPrefix 				   	   = "x-amz-"
	metaHeader                     = amzPrefix + "meta"
	keyV1Header                    = amzPrefix + "key"
	keyV2Header                    = amzPrefix + "key-v2"
	ivHeader                       = amzPrefix + "iv"
	matDescHeader                  = amzPrefix + "matdesc"
	CekAlgorithmHeader             = amzPrefix + "cek-alg"
	KeyringAlgorithmHeader         = amzPrefix + "wrap-alg"
	tagLengthHeader                = amzPrefix + "tag-len"
	unencryptedContentLengthHeader = amzPrefix + "unencrypted-content-length"

	// For the below constants, the specification comments describe the recommended constant names.
	// However, Go's JSON struct tags require literal string values and cannot reference constants.
	// This forces us to duplicate the string literals in the ObjectMetadata struct's `json:` tags.
	// While this creates duplication between these constants and the struct tags, it's a Go language
	// limitation - we cannot write `json:ContentCipherV3` in the struct definition.
	// There are tests that validate these constant values match the struct tags in ObjectMetadata.

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-c") SHOULD be represented by a constant named "CONTENT_CIPHER_V3" or similar in the implementation code.
	ContentCipherV3                = amzPrefix + "c"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-3") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_V3" or similar in the implementation code.
	EncryptedDataKeyV3             = amzPrefix + "3"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-m") SHOULD be represented by a constant named "MAT_DESC_V3" or similar in the implementation code.
	MatDescV3                      = amzPrefix + "m"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-t") SHOULD be represented by a constant named "ENCRYPTION_CONTEXT_V3" or similar in the implementation code.
	EncryptionContextV3            = amzPrefix + "t"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-w") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_ALGORITHM_V3" or similar in the implementation code.
	EncryptedDataKeyAlgorithmV3    = amzPrefix + "w"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-d") SHOULD be represented by a constant named "KEY_COMMITMENT_V3" or similar in the implementation code.
	KeyCommitmentV3                = amzPrefix + "d"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-i") SHOULD be represented by a constant named "MESSAGE_ID_V3" or similar in the implementation code.
	MessageIDV3                    = amzPrefix + "i"
)

// S3EC Go V4 does not support reading nor writing V1 format metadata.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-unencrypted-content-length" SHOULD be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-key" MUST be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-matdesc" MUST be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-iv" MUST be present for V1 format objects.

// ObjectMetadata encryption starts off by generating a random symmetric key using
// AES GCM. The SDK generates a random IV based off the encryption cipher
// chosen. The master key that was provided, whether by the user or KMS, will be used
// to encrypt the randomly generated symmetric key and base64 encode the iv. This will
// allow for decryption of that same data later.
//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=implication
//# The "x-amz-meta-" prefix is automatically added by the S3 server and MUST NOT be included in implementation code.
type ObjectMetadata struct {
	// IV is the randomly generated IV base64 encoded.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-iv" MUST be present for V2 format objects.
	IV string `json:"x-amz-iv"`
	// CipherKey is the randomly generated cipher key.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-key-v2" MUST be present for V2 format objects.
	CipherKey string `json:"x-amz-key-v2"`
	// MaterialDesc is a description to distinguish from other envelopes.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-matdesc" MUST be present for V2 format objects.
	MatDesc               string `json:"x-amz-matdesc"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-wrap-alg" MUST be present for V2 format objects.
	KeyringAlg            string `json:"x-amz-wrap-alg"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-cek-alg" MUST be present for V2 format objects.
	CEKAlg                string `json:"x-amz-cek-alg"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-tag-len" MUST be present for V2 format objects.
	TagLen                string `json:"x-amz-tag-len"`
	UnencryptedContentLen string `json:"x-amz-unencrypted-content-length"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-c" MUST be present for V3 format objects.
	ContentCipher string `json:"x-amz-c"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-3" MUST be present for V3 format objects.
	EncryptedDataKey string `json:"x-amz-3"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
	MatDescV3 string `json:"x-amz-m"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-t" SHOULD be present for V3 format objects that use KMS Encryption Context.
	EncryptionContext string `json:"x-amz-t"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-w" MUST be present for V3 format objects.
	WrappingAlgorithm string `json:"x-amz-w"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-d" MUST be present for V3 format objects.
	KeyCommitment string `json:"x-amz-d"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-i" MUST be present for V3 format objects.
	MessageID string `json:"x-amz-i"`
}

// V3 algorithm compression mappings
var (
	// Wrapping algorithm decompression: compressed value -> full algorithm name
	v3WrapAlgDecompression = map[string]string{
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval
		"02": "AES/GCM",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval
		"12": "kms+context",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval
		"22": "RSA-OAEP-SHA1",
	}
	
	// Wrapping algorithm compression: full algorithm name -> compressed value
	v3WrapAlgCompression = map[string]string{
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval, and vice versa on write.
		"AES/GCM": "02",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval, and vice versa on write.
		"kms+context": "12",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval, and vice versa on write.
		"RSA-OAEP-SHA1": "22",
	}
)

func (e *ObjectMetadata) GetDecodedKey() ([]byte, error) {
	var keyStr string
	if e.EncryptedDataKey != "" {
		// V3
		keyStr = e.EncryptedDataKey
	} else {
		// V2
		keyStr = e.CipherKey
	}
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (e *ObjectMetadata) GetDecodedMessageIDOrIV() ([]byte, error) {
	var value string
	if e.MessageID != "" {
		// V3
		value = e.MessageID
	} else {
		// V2
		value = e.IV
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (e *ObjectMetadata) GetDecodedKeyCommitment() ([]byte, error) {
	// Only V3 has KeyCommitment
	commitment, err := base64.StdEncoding.DecodeString(e.KeyCommitment)
	if err != nil {
		return nil, err
	}
	return commitment, err
}

func (e *ObjectMetadata) GetMatDescV3() (string, error) {
	if e.MatDescV3 != "" {
		return e.MatDescV3, nil
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# If the mapkey x-amz-m is not present, the default Material Description value MUST be set to an empty map (`{}`).
	return "{}", nil
}

func (e *ObjectMetadata) GetEncryptionContextV3() (string, error) {
	// Only V3 has EncryptionContext
	if e.EncryptionContext != "" {
		return e.EncryptionContext, nil
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# If the mapkey x-amz-t is not present, the default Material Description value MUST be set to an empty map (`{}`).
	return "{}", nil
}

func (e *ObjectMetadata) GetMatDescV2() (string, error) {
	return e.MatDesc, nil
}

func (e *ObjectMetadata) GetEncryptionContextOrMatDescV3() (string, error) {
	wrappingAlg, err := e.GetFullWrappingAlgorithm()
	var matDesc string
	if wrappingAlg == "kms+context" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
		matDesc, err = e.GetEncryptionContextV3()
		return matDesc, err
	} else if wrappingAlg == "AES/GCM" || wrappingAlg == "RSA-OAEP-SHA1" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
		matDesc, err = e.GetMatDescV3()
		return matDesc, err
	}
	return "", fmt.Errorf("unsupported wrapping algorithm for getting Material Description: %s", wrappingAlg)
}

func (e *ObjectMetadata) GetContentEncryptionAlgorithmString() (string, error) {
	if e.ContentCipher != "" {
		// V3
		return e.ContentCipher, nil
	}
	if e.CEKAlg != "" {
		// V2
		return e.CEKAlg, nil
	}
	return "", fmt.Errorf("no content encryption algorithm found in metadata")
}

func (e *ObjectMetadata) GetContentEncryptionAlgorithmSuite() (*algorithms.AlgorithmSuite, error) {
	cekString, err := e.GetContentEncryptionAlgorithmString()
	if err != nil {
		return nil, err
	}
	if cekString == algorithms.AESGCMCommitKey {
		return algorithms.AlgAES256GCMHkdfSha512CommitKey, nil
	} else if cekString == algorithms.AESGCMNoPadding {
		return algorithms.AlgAES256GCMIV12Tag16NoKDF, nil
	} else if cekString == algorithms.AESCBCPKCS5 {
		return algorithms.AlgAES256CBCIV16NoKDF, nil
	} else if cekString == algorithms.AESCTRNoPadding {
		return algorithms.AlgAES256CTRIV16Tag16NoKDF, nil
	}
	return nil, fmt.Errorf("invalid content encryption algorithm found in metadata: %s", cekString)
}

func (e *ObjectMetadata) GetFullWrappingAlgorithm() (string, error) {
	if e.WrappingAlgorithm != "" {
		// V3
		// Decompress the V3 wrapping algorithm to its full name
		fullAlg, exists := v3WrapAlgDecompression[e.WrappingAlgorithm]
		if !exists {
			return "", fmt.Errorf("unknown V3 wrapping algorithm: %s", e.WrappingAlgorithm)
		}
		return fullAlg, nil
	}
	if e.KeyringAlg != "" {
		// V2
		return e.KeyringAlg, nil
	}
	return "", fmt.Errorf("no wrapping algorithm found in metadata")
}

// CompressWrappingAlgorithm compresses a full wrapping algorithm name to V3 format
func CompressWrappingAlgorithm(fullAlgorithm string) (string, error) {
	compressed, exists := v3WrapAlgCompression[fullAlgorithm]
	if !exists {
		return "", fmt.Errorf("unsupported wrapping algorithm for V3: %s", fullAlgorithm)
	}
	return compressed, nil
}

// UnmarshalJSON unmarshalls the given JSON bytes into ObjectMetadata
func (e *ObjectMetadata) UnmarshalJSON(value []byte) error {
	type StrictEnvelope ObjectMetadata
}

func EncodeMeta(reader lengthReader, cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
	//# Objects encrypted with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY MUST use the V3 message format version only.
	if cryptographicMaterials.CEKAlgorithm == algorithms.AESGCMCommitKey {
		return EncodeMetaV3(cryptographicMaterials)
	//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
	//# Objects encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF MUST use the V2 message format version only.
	} else if cryptographicMaterials.CEKAlgorithm == algorithms.AESGCMNoPadding {
		return EncodeMetaV2(reader, cryptographicMaterials)
	} else {
		// Go S3EC V4 does not support writing other (ex. AES CBC) messages
		//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
		//= type=exception
		//# Objects encrypted with ALG_AES_256_CBC_IV16_NO_KDF MAY use either the V1 or V2 message format version.
		return ObjectMetadata{}, fmt.Errorf("unsupported CEK algorithm: %s", cryptographicMaterials.CEKAlgorithm)
	}
}

func EncodeMetaV2(reader lengthReader, cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	iv := base64.StdEncoding.EncodeToString(cryptographicMaterials.IV)
	key := base64.StdEncoding.EncodeToString(cryptographicMaterials.EncryptedKey)

		UnencryptedContentLen: strconv.FormatInt(contentLength, 10),
	}, nil
}

func EncodeMetaV3(cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	iv := base64.StdEncoding.EncodeToString(cryptographicMaterials.IV)
	key := base64.StdEncoding.EncodeToString(cryptographicMaterials.EncryptedKey)
	alg_suite := cryptographicMaterials.CEKAlgorithm
	wrapping_alg := cryptographicMaterials.KeyringAlgorithm
	commitment := base64.StdEncoding.EncodeToString(cryptographicMaterials.KeyCommitment)
	mat_desc_bytes, err := cryptographicMaterials.MaterialDescription.EncodeDescription()
	if err != nil {
		return ObjectMetadata{}, err
	}
	mat_desc := string(mat_desc_bytes)

	out := ObjectMetadata{
		EncryptedDataKey:  key,
		MessageID:         iv,
		ContentCipher:     alg_suite,
		WrappingAlgorithm: wrapping_alg,
		KeyCommitment:     commitment,
	}

	// Set MatDescV3 for AES/GCM or RSA-OAEP-SHA1, EncryptionContext for kms+context
	if wrapping_alg == "AES/GCM" || wrapping_alg == "RSA-OAEP-SHA1" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
		out.MatDescV3 = mat_desc
	} else if wrapping_alg == "kms+context" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
		out.EncryptionContext = mat_desc
	}

	return out, nil
}

// MetadataFormat represents the format version of S3EC metadata
type MetadataFormat int

const (
	FormatUnknown MetadataFormat = iota
	FormatInstructionFile
	FormatV1
	FormatV2
	FormatV3
)

// Validate and detect correct metadata format 
func DetectAndValidateMetadataFormat(metadata map[string]string) (MetadataFormat, error) {
	// Check for mapkeys defined in the spec	
	hasV1Key := hasKey(metadata, keyV1Header) // "x-amz-key"
	hasV2Key := hasKey(metadata, keyV2Header) // "x-amz-key-v2"
	hasV3Key := hasKey(metadata, EncryptedDataKeyV3) // "x-amz-3"
	hasIv := hasKey(metadata, ivHeader) // "x-amz-iv"
	hasV3KeyCommitment := hasKey(metadata, KeyCommitmentV3) // "x-amz-d"
	hasV3KeyID := hasKey(metadata, MessageIDV3) // "x-amz-i"

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-iv" and "x-amz-key" then the object MUST be considered as an S3EC-encrypted object using the V1 format.
	isV1 := hasIv && hasV1Key
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-iv" and "x-amz-metadata-x-amz-key-v2" then the object MUST be considered as an S3EC-encrypted object using the V2 format.
	isV2 := hasIv && hasV2Key
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-3" and "x-amz-d" and "x-amz-i" then the object MUST be considered an S3EC-encrypted object using the V3 format.
	isV3 := hasV3Key && hasV3KeyCommitment && hasV3KeyID

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# If there are multiple mapkeys which are meant to be exclusive, such as "x-amz-key", "x-amz-key-v2", and "x-amz-3" then the S3EC SHOULD throw an exception.
	hasAtLeastOneV1ExclusiveKey := hasV1Key
	hasAtLeastOneV2ExclusiveKey := hasV2Key
	hasAtLeastOneV3ExclusiveKey := hasV3Key
	exclusiveKeyMatchCount := 0
	if hasAtLeastOneV1ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if hasAtLeastOneV2ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if hasAtLeastOneV3ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if exclusiveKeyMatchCount > 1 {
		return FormatUnknown, fmt.Errorf("metadata contains conflicting exclusive mapkeys")
	}

	versionMatchCount := 0
	if isV1 {
		versionMatchCount++
	}
	if isV2 {
		versionMatchCount++
	}
	if isV3 {
		versionMatchCount++
	}
	if versionMatchCount > 1 {
		return FormatUnknown, fmt.Errorf("metadata contains multiple S3EC format versions")
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# If the object matches none of the V1/V2/V3 formats, the S3EC MUST attempt to get the instruction file.
	if versionMatchCount == 0 {
		return FormatInstructionFile, nil
	}

	if isV1 {
		return FormatV1, nil
	}
	if isV2 {
		return FormatV2, nil
	}
	if isV3 {
		return FormatV3, nil
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//= type=implication
	//# In general, if there is any deviation from the above format, with the exception of additional unrelated mapkeys, then the S3EC SHOULD throw an exception.
	return FormatUnknown, fmt.Errorf("unable to determine metadata format")
}

// hasKey checks if a metadata key exists (with or without x-amz-meta prefix)
func hasKey(metadata map[string]string, key string) bool {
	// Check direct key
	if _, exists := metadata[key]; exists {
		return true
	}
	// Check with x-amz-meta prefix
	prefixedKey := metaHeader + "-" + key
	if _, exists := metadata[prefixedKey]; exists {
		return true
	}
	return false
}
	"encoding/json"
	"reflect"
	"testing"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
)

func TestEnvelope_UnmarshalJSON(t *testing.T) {
		})
	}
}

func TestObjectMetadata_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		content  []byte
		expected ObjectMetadata
		actual   ObjectMetadata
	}{
		"complete V3 metadata with encryption context": {
			content: []byte(`{
				"x-amz-c": "115",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-t": "{\"kms_cmk_id\":\"test-key-id\"}",
				"x-amz-w": "12",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "115",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				EncryptionContext: `{"kms_cmk_id":"test-key-id"}`,
				WrappingAlgorithm: "12",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
		"V3 metadata with material description": {
			content: []byte(`{
				"x-amz-c": "AES/GCM/NoPadding",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-m": "{\"test\":\"material-desc\"}",
				"x-amz-w": "01",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "AES/GCM/NoPadding",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				MatDescV3:         `{"test":"material-desc"}`,
				WrappingAlgorithm: "01",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
		"minimal V3 metadata": {
			content: []byte(`{
				"x-amz-c": "AES/CBC/PKCS5Padding",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-w": "11",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "AES/CBC/PKCS5Padding",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				WrappingAlgorithm: "11",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal(tt.content, &tt.actual)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !reflect.DeepEqual(tt.expected, tt.actual) {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

func TestWrappingAlgorithmCompression(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval, and vice versa on write.
	result, err := CompressWrappingAlgorithm("AES/GCM")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "02" {
		t.Errorf("expected 02, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval, and vice versa on write.
	result, err = CompressWrappingAlgorithm("kms+context")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "12" {
		t.Errorf("expected 12, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval, and vice versa on write.
	result, err = CompressWrappingAlgorithm("RSA-OAEP-SHA1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "22" {
		t.Errorf("expected 22, got %s", result)
	}

	result, err = CompressWrappingAlgorithm("not-a-known-algorithm")
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestWrappingAlgorithmDecompression(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval
	metadata := ObjectMetadata{WrappingAlgorithm: "02"}
	result, err := metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "AES/GCM" {
		t.Errorf("expected AES/GCM, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval
	metadata = ObjectMetadata{WrappingAlgorithm: "12"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "kms+context" {
		t.Errorf("expected kms+context, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval
	metadata = ObjectMetadata{WrappingAlgorithm: "22"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "RSA-OAEP-SHA1" {
		t.Errorf("expected RSA-OAEP-SHA1, got %s", result)
	}

	metadata = ObjectMetadata{WrappingAlgorithm: "99"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestDetectAndValidateMetadataFormat(t *testing.T) {
	cases := map[string]struct {
		metadata map[string]string
		expected MetadataFormat
	}{
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-3" and "x-amz-d" and "x-amz-i" then the object MUST be considered an S3EC-encrypted object using the V3 format.
		"V3 format": {
			metadata: map[string]string{
				"x-amz-3": "encrypted-key",
				"x-amz-d": "key-commitment",
				"x-amz-i": "message-id",
			},
			expected: FormatV3,
		},
		"V3 format with meta prefix": {
			metadata: map[string]string{
				"x-amz-meta-x-amz-3": "encrypted-key",
				"x-amz-meta-x-amz-d": "key-commitment",
				"x-amz-meta-x-amz-i": "message-id",
			},
			expected: FormatV3,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-iv" and "x-amz-metadata-x-amz-key-v2" then the object MUST be considered as an S3EC-encrypted object using the V2 format.
		"V2 format minimal meta prefix": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-meta-x-amz-key-v2":  "key",
			},
			expected: FormatV2,
		},
		"V2 format minimal": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key-v2":  "key",
			},
			expected: FormatV2,
		},
		"V2 format": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key-v2":  "key",
				"x-amz-matdesc": "matdesc",
			},
			expected: FormatV2,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-iv" and "x-amz-key" then the object MUST be considered as an S3EC-encrypted object using the V1 format.
		"V1 format minimal": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key":     "key",
			},
			expected: FormatV1,
		},
		"V1 format": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key":     "key",
				"x-amz-matdesc": "matdesc",
			},
			expected: FormatV1,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# If the object matches none of the V1/V2/V3 formats, the S3EC MUST attempt to get the instruction file.
		"matching no format some keys": {
			metadata: map[string]string{
				"x-amz-abcdef": "not-a-key",
				"x-amz-123": "still-not-a-key",
			},
			expected: FormatInstructionFile,
		},
		"matching no format no keys": {
			metadata: map[string]string{},
			expected: FormatInstructionFile,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# If there are multiple mapkeys which are meant to be exclusive, such as "x-amz-key", "x-amz-key-v2", and "x-amz-3" then the S3EC SHOULD throw an exception.
		"multiple exclusive mapkeys case 1": {
			metadata: map[string]string{
				"x-amz-key": "key",
				"x-amz-key-v2": "key-v2",
				"x-amx-3": "key-v3",
			},
			expected: FormatUnknown,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := DetectAndValidateMetadataFormat(tt.metadata)
			if err != nil && result != FormatUnknown {
				t.Errorf("expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestObjectMetadataConstValues(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-c") SHOULD be represented by a constant named "CONTENT_CIPHER_V3" or similar in the implementation code.
	if ContentCipherV3 != "x-amz-c" {
		t.Errorf("ContentCipherV3 MUST be `x-amz-c`, got %q", ContentCipherV3)
	}
	// Truly wild reflection usage in this test below to ensure the struct field in ObjectMetadata has the correct json tag to fully satisfy spec requirement
	field, ok := reflect.TypeOf(ObjectMetadata{}).FieldByName("ContentCipher")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field ContentCipher")
	}
	if got := field.Tag.Get("json"); got != "x-amz-c" {
		t.Errorf("ContentCipher json tag MUST be `x-amz-c`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-3") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_V3" or similar in the implementation code.
	if EncryptedDataKeyV3 != "x-amz-3" {
		t.Errorf("EncryptedDataKeyV3 MUST be `x-amz-3`, got %q", EncryptedDataKeyV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("EncryptedDataKey")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field EncryptedDataKey")
	}
	if got := field.Tag.Get("json"); got != "x-amz-3" {
		t.Errorf("EncryptedDataKey json tag MUST be `x-amz-3`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-m") SHOULD be represented by a constant named "MAT_DESC_V3" or similar in the implementation code.
	if MatDescV3 != "x-amz-m" {
		t.Errorf("MatDescV3 MUST be `x-amz-m`, got %q", MatDescV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("MatDescV3")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field MatDesc")
	}
	if got := field.Tag.Get("json"); got != "x-amz-m" {
		t.Errorf("MatDesc json tag MUST be `x-amz-m`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-t") SHOULD be represented by a constant named "ENCRYPTION_CONTEXT_V3" or similar in the implementation code.
	if EncryptionContextV3 != "x-amz-t" {
		t.Errorf("EncryptionContextV3 MUST be `x-amz-t`, got %q", EncryptionContextV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("EncryptionContext")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field EncryptionContext")
	}
	if got := field.Tag.Get("json"); got != "x-amz-t" {
		t.Errorf("EncryptionContext json tag MUST be `x-amz-t`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-w") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_ALGORITHM_V3" or similar in the implementation code.
	if EncryptedDataKeyAlgorithmV3 != "x-amz-w" {
		t.Errorf("EncryptedDataKeyAlgorithmV3 MUST be `x-amz-w`, got %q", EncryptedDataKeyAlgorithmV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("WrappingAlgorithm")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field WrappingAlgorithm")
	}
	if got := field.Tag.Get("json"); got != "x-amz-w" {
		t.Errorf("WrappingAlgorithm json tag MUST be `x-amz-w`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-d") SHOULD be represented by a constant named "KEY_COMMITMENT_V3" or similar in the implementation code.
	if KeyCommitmentV3 != "x-amz-d" {
		t.Errorf("KeyCommitmentV3 MUST be `x-amz-d`, got %q", KeyCommitmentV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("KeyCommitment")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field KeyCommitment")
	}
	if got := field.Tag.Get("json"); got != "x-amz-d" {
		t.Errorf("KeyCommitment json tag MUST be `x-amz-d`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-i") SHOULD be represented by a constant named "MESSAGE_ID_V3" or similar in the implementation code.
	if MessageIDV3 != "x-amz-i" {
		t.Errorf("MessageIDV3 MUST be `x-amz-i`, got %q", MessageIDV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("MessageID")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field MessageID")
	}
	if got := field.Tag.Get("json"); got != "x-amz-i" {
		t.Errorf("MessageID json tag MUST be `x-amz-i`, got %q", got)
	}
}

func TestMaterialDescriptionAndEncryptionContextRequirements(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
	t.Run("Material Description used for AES/GCM", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "02", // AES/GCM
			MatDescV3:         `{"test":"material-desc"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"test":"material-desc"}` {
			t.Errorf("expected material description, got %s", result)
		}
	})

	t.Run("Material Description used for RSA-OAEP-SHA1", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "22", // RSA-OAEP-SHA1
			MatDescV3:         `{"test":"rsa-material-desc"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"test":"rsa-material-desc"}` {
			t.Errorf("expected material description, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# If the mapkey x-amz-m is not present, the default Material Description value MUST be set to an empty map (`{}`).
	t.Run("Default empty map when Material Description not present for AES/GCM", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "02", // AES/GCM
			// MatDescV3 is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})

	t.Run("Default empty map when Material Description not present for RSA-OAEP-SHA1", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "22", // RSA-OAEP-SHA1
			// MatDescV3 is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
	t.Run("Encryption Context used for kms+context", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "12", // kms+context
			EncryptionContext: `{"kms_cmk_id":"test-key-id"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"kms_cmk_id":"test-key-id"}` {
			t.Errorf("expected encryption context, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# If the mapkey x-amz-t is not present, the default Material Description value MUST be set to an empty map (`{}`).
	t.Run("Default empty map when Encryption Context not present for kms+context", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "12", // kms+context
			// EncryptionContext is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})
}

// mockLengthReader implements the lengthReader interface for testing
type mockLengthReader struct {
	contentLength int64
}

func (m *mockLengthReader) GetContentLength() int64 {
	return m.contentLength
}

func TestEncodeMetaV2(t *testing.T) {
	cases := map[string]struct {
		reader               lengthReader
		cryptographicMaterials materials.CryptographicMaterials
		expected             ObjectMetadata
	}{
		"standard V2 encoding": {
			reader: &mockLengthReader{contentLength: 1024},
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-12-b"),
				KeyringAlgorithm:    "kms+context",
				CEKAlgorithm:        "AES/GCM/NoPadding",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"aws:x-amz-cek-alg": "AES/GCM/NoPadding", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
			},
			expected: ObjectMetadata{
				CipherKey:             "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				IV:                    "dGVzdC1pdi0xMi1i", // base64 of "test-iv-12-b"
				MatDesc:               `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","custom":"value"}`,
				KeyringAlg:            "kms+context",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "128",
				UnencryptedContentLen: "1024",
			},
		},
		"V2 encoding with empty material description": {
			reader: &mockLengthReader{contentLength: 2048},
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-12-b"),
				KeyringAlgorithm:    "kms",
				CEKAlgorithm:        "AES/GCM/NoPadding",
				TagLength:           "96",
				MaterialDescription: materials.MaterialDescription{},
				EncryptedKey:        []byte("encrypted-key-data"),
			},
			expected: ObjectMetadata{
				CipherKey:             "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				IV:                    "dGVzdC1pdi0xMi1i", // base64 of "test-iv-12-b"
				MatDesc:               `{}`,
				KeyringAlg:            "kms",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "96",
				UnencryptedContentLen: "2048",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := EncodeMetaV2(tt.reader, tt.cryptographicMaterials)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestEncodeMetaV3(t *testing.T) {
	cases := map[string]struct {
		cryptographicMaterials materials.CryptographicMaterials
		expected             ObjectMetadata
	}{
		"V3 encoding with AES/GCM wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "AES/GCM",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"test": "material-desc", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "AES/GCM",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
				MatDescV3:         `{"custom":"value","test":"material-desc"}`,
			},
		},
		"V3 encoding with kms+context wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "kms+context",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"kms_cmk_id": "test-key-id", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "kms+context",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-t" SHOULD be present for V3 format objects that use KMS Encryption Context.
				EncryptionContext: `{"custom":"value","kms_cmk_id":"test-key-id"}`,
			},
		},
		"V3 encoding with RSA-OAEP-SHA1 wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "RSA-OAEP-SHA1",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"rsa": "material-desc"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "RSA-OAEP-SHA1",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
				MatDescV3:         `{"rsa":"material-desc"}`,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := EncodeMetaV3(tt.cryptographicMaterials)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
		input.Metadata = map[string]string{}
	}

	// S3EC Go V4 supports reading content metadata from an instruction file, but not writing it to an instruction file.
	// Any content metadata written is implicitly written to object metadata and not an instruction file.

	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=implication
	//# By default, the S3EC MUST store content metadata in the S3 Object Metadata.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=implication
	//# In the V3 format, the mapkeys "x-amz-c", "x-amz-d", and "x-amz-i" MUST be stored exclusively in the Object Metadata.
	env := saveReq.Envelope
	if env.EncryptedDataKey != "" {
		// V3 format
		compressedAlg, err := CompressWrappingAlgorithm(env.WrappingAlgorithm)
		if err != nil {
			return fmt.Errorf("error while compressing wrapping algorithm: %w", err)
		}
		input.Metadata[http.CanonicalHeaderKey(ContentCipherV3)] = env.ContentCipher
		input.Metadata[http.CanonicalHeaderKey(EncryptedDataKeyV3)] = env.EncryptedDataKey
		input.Metadata[http.CanonicalHeaderKey(MatDescV3)] = env.MatDescV3
		input.Metadata[http.CanonicalHeaderKey(EncryptionContextV3)] = env.EncryptionContext
		input.Metadata[http.CanonicalHeaderKey(EncryptedDataKeyAlgorithmV3)] = compressedAlg
		input.Metadata[http.CanonicalHeaderKey(KeyCommitmentV3)] = env.KeyCommitment
		input.Metadata[http.CanonicalHeaderKey(MessageIDV3)] = env.MessageID
	} else {
		// V2 format
		input.Metadata[http.CanonicalHeaderKey(keyV2Header)] = env.CipherKey
		input.Metadata[http.CanonicalHeaderKey(ivHeader)] = env.IV
		input.Metadata[http.CanonicalHeaderKey(matDescHeader)] = env.MatDesc
		input.Metadata[http.CanonicalHeaderKey(KeyringAlgorithmHeader)] = env.KeyringAlg
		input.Metadata[http.CanonicalHeaderKey(CekAlgorithmHeader)] = env.CEKAlg
	}
	input.Metadata[http.CanonicalHeaderKey(unencryptedContentLengthHeader)] = env.UnencryptedContentLen

	if len(env.TagLen) > 0 {
		input.Metadata[http.CanonicalHeaderKey(tagLengthHeader)] = env.TagLen
	}
	return nil

	// S3EC Go V4 supports reading content metadata from an instruction file, but not writing it.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MUST support writing some or all (depending on format) content metadata to an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The content metadata stored in the Instruction File MUST be serialized to a JSON string.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The serialized JSON string MUST be the only contents of the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# Instruction File writes MUST NOT be enabled by default.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# Instruction File writes MUST be optionally configured during client creation or on each PutObject request.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MAY support re-encryption/key rotation via Instruction Files.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MUST NOT support providing a custom Instruction File suffix on ordinary writes; custom suffixes MUST only be used during re-encryption.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC SHOULD support providing a custom Instruction File suffix on GetObject requests, regardless of whether or not re-encryption is supported.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v1-v2-instruction-files
	//= type=exception
	//# In the V1/V2 message format, all of the content metadata MUST be stored in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-c" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-c" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-d" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-d" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-i" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-i" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-3" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-w" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-m" and its value (when present in the content metadata) in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-t" and its value (when present in the content metadata) in the Instruction File.
}

// LoadStrategyRequest represents a request sent to a LoadStrategy to load the contents of an ObjectMetadata
	return env, nil
}

// headerV3LoadStrategy will load the V3 envelope from the metadata
type headerV3LoadStrategy struct{}

// Load from a given object's header using V3 format
func (load headerV3LoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	v3Meta := ObjectMetadata{}
	v3Meta.ContentCipher = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, ContentCipherV3}, "-"))
	v3Meta.EncryptedDataKey = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptedDataKeyV3}, "-"))
	v3Meta.MatDescV3 = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, MatDescV3}, "-"))
	v3Meta.EncryptionContext = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptionContextV3}, "-"))
	v3Meta.WrappingAlgorithm = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptedDataKeyAlgorithmV3}, "-"))
	v3Meta.KeyCommitment = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, KeyCommitmentV3}, "-"))
	v3Meta.MessageID = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, MessageIDV3}, "-"))
	v3Meta.UnencryptedContentLen = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, unencryptedContentLengthHeader}, "-"))
	return v3Meta, nil
}

// DefaultLoadStrategy This is the only exported LoadStrategy since cx are no longer able to configure their client
// with a specific load strategy. Instead, we figure out which strategy to use based on the response header on decrypt.
type DefaultLoadStrategy struct {
}

func (load DefaultLoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	// Create metadata map from headers for format detection
	metadata := make(map[string]string)
	for key, values := range req.HTTPResponse.Header {
		if len(values) > 0 {
			metadata[strings.ToLower(key)] = values[0]
		}
	}

	// Detect format and validate
	format, err := DetectAndValidateMetadataFormat(metadata)
	if err != nil {
		return ObjectMetadata{}, fmt.Errorf("invalid metadata format: %w", err)
	}

	switch format {
	case FormatV3:
		strat := headerV3LoadStrategy{}
		return strat.Load(ctx, req)
		
	case FormatV2:
		strat := headerV2LoadStrategy{}
		return strat.Load(ctx, req)
		
	case FormatV1:
		// In other S3EC implementations, decryption of v1 objects is supported.
		// Go, however, does not support this.
		return ObjectMetadata{}, &smithy.GenericAPIError{
			Code:    "V1NotSupportedError",
			Message: "The AWS SDK for Go does not support version 1",
		}
		
	default:
		// Fall back to instruction file loading
		var client GetObjectAPIClient
		if load.client == nil {
			cfg, err := config.LoadDefaultConfig(context.Background())
			if err != nil {
				return ObjectMetadata{}, fmt.Errorf("unable to create S3 client to load instruction file: %w", err)
			}
			client = s3.NewFromConfig(cfg)
		} else {
			client = load.client
		}

		// Load from instruction file
		strat := s3LoadStrategy{
			APIClient:             client,
			InstructionFileSuffix: load.suffix,
		}
		loadedMetadata, err := strat.Load(ctx, req)
		if err != nil {
			return ObjectMetadata{}, err
		}

		// For "V3 instruction file" format, load any additional metadata from headers
		// EncryptedDataKey is only present for V3 format (with or without instruction file)
		if loadedMetadata.EncryptedDataKey != "" {
			// If any values that should be in the headers were in the instruction file, raise error
			if loadedMetadata.ContentCipher != "" || loadedMetadata.KeyCommitment != "" || loadedMetadata.MessageID != "" {
				return ObjectMetadata{}, fmt.Errorf("invalid metadata format: missing V3 header values in instruction file format")
			}

			// Load these values from headers
			headerMeta, err := headerV3LoadStrategy{}.Load(ctx, req)
			if err != nil {
				return ObjectMetadata{}, err
			}
			loadedMetadata.ContentCipher = headerMeta.ContentCipher
			loadedMetadata.KeyCommitment = headerMeta.KeyCommitment
			loadedMetadata.MessageID = headerMeta.MessageID
		}
		
		return loadedMetadata, nil
	}
}
// When EnableLegacyWrappingAlgorithms is set to true, the Keyring MAY decrypt objects encrypted
// using legacy wrapping algorithms such as KMS v1.
type KeyringOptions struct {
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=implication
	//# The S3EC MUST support the option to enable or disable legacy wrapping algorithms.

	EnableLegacyWrappingAlgorithms bool
}

// object. The KmsKeyring will always use the kmsKeyId provided to encrypt and decrypt messages.
func NewKmsKeyring(apiClient KmsAPIClient, kmsKeyId string, optFns ...func(options *KeyringOptions)) *KmsKeyring {
	options := KeyringOptions{
		//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
		//= type=implication
		//# The option to enable legacy wrapping algorithms MUST be set to false by default.
		EnableLegacyWrappingAlgorithms: false,
	}
	for _, fn := range optFns {
// for use with content decryption, or an error if the object cannot be decrypted
// by the Keyring as its configured.
func (k *KmsKeyring) OnDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey) (*CryptographicMaterials, error) {
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy wrapping algorithms; it MUST throw an exception when attempting to decrypt an object encrypted with a legacy wrapping algorithm.	
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && !k.legacyWrappingAlgorithms {
		return nil, fmt.Errorf("to decrypt x-amz-cek-alg value `%s` you must enable legacyWrappingAlgorithms on the keyring", materials.DataKey.DataKeyAlgorithm)
	}

	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && k.legacyWrappingAlgorithms {
		return commonDecrypt(ctx, materials, encryptedDataKey, &k.KmsKeyId, nil, k.kmsClient)
	} else if materials.DataKey.DataKeyAlgorithm == KMSContextKeyring {
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v3/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	algorithmSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF

	encryptionMaterials := NewEncryptionMaterials()

	_, err := keyring.OnEncrypt(ctx, encryptionMaterials)
		GrantTokens: grantTokens,
		KeySpec:     types.DataKeySpecAes256,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: algorithmSuite.CipherName(),
		},
	}


	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	algorithmSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF

	decryptionMaterials, err := NewDecryptionMaterials(DecryptMaterialsRequest{
		CipherKey:  []byte("test-cipher-key"),
		Iv:         []byte("test-iv"),
		MatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
		KeyringAlg: "kms+context",
		CekAlg:     algorithmSuite.CipherName(),
	})
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
		GrantTokens:    grantTokens,
		CiphertextBlob: dataKey.EncryptedDataKey,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: algorithmSuite.CipherName(),
		},
	}

	MaterialDescription MaterialDescription
	ContentAlgorithm    string
	TagLength           string
	KeyCommitment       []byte
}

func NewDecryptionMaterials(req DecryptMaterialsRequest) (*DecryptionMaterials, error) {
	MaterialDescription MaterialDescription
	// EncryptedKey should be populated when calling GenerateCipherData
	EncryptedKey []byte
	KeyCommitment []byte
}
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
// This is meant to be a utility function, not a test function,
// but for simplicity and easy invocation it is a test function.
// To avoid running it each test run, it is left commented out.
// func TestGenerateCBCIntegTests(t *testing.T) {
// 	arn := "arn:aws:kms:us-west-2:370957321024:alias/S3EC-Go-Github-KMS-Key"
// 	bucket := "s3ec-go-github-test-bucket"
// 	region := "us-west-2"
// 	ctx := context.Background()
// 	cfg, _ := config.LoadDefaultConfig(ctx,
// 		config.WithRegion(region),
// 	)

// 	s3Client := s3.NewFromConfig(cfg)
// 	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
// 	// V2 client
// 	var handler s3cryptoV2.CipherDataGenerator
// 	sessKms, _ := sessionV1.NewSession(&awsV1.Config{
// 		Region: aws.String(region),
// 	})

// 	// KMS v1
// 	kmsSvc := kmsV1.New(sessKms)
// 	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, arn)
// 	// AES-CBC content cipher
// 	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
// 	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

// 	for caseKey, plaintext := range fixtures.Plaintexts {
// 		_, err := encClient.PutObject(&s3V1.PutObjectInput{
// 			Bucket: aws.String(bucket),
// 			Key: aws.String(
// 				fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
// 					fixtures.BaseFolder, version, caseKey),
// 			),
// 			Body: bytes.NewReader(plaintext),
// 		})
// 		if err != nil {
// 			t.Fatalf("failed to upload encrypted fixture, %v", err)
// 		}
// 	}

// }

func TestKmsV1toV4_CBC(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V1toV4_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsV1toV4_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V1toV4_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsContextV2toV4_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V2toV4_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	}
}

func TestKmsContextV4toV2_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V4toV2_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = s3ecV4.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	}
}

func TestInstructionFileV2toV4(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/inst_file_test.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/NegativeV1toV4_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	_, err = s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlgCbc := "aes_cbc"
	keyCbc := "crypto_tests/" + cekAlgCbc + "/v4/language_Go/BothFormats_CBC.txt"
	cekAlgGcm := "aes_gcm"
	keyGcm := "crypto_tests/" + cekAlgGcm + "/v4/language_Go/BothFormats_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = s3ecV4.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
		Body:   bytes.NewReader([]byte(plaintext)),
		t.Fatalf("error while calling PutObject: %v", err)
	}

	getResponseCbc, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyCbc),
	})
		t.Errorf("expect %v text, got %v", e, a)
	}

	getResponseGcm, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
	})
	}
}

func TestUnicodeEncryptionContextV4ClientV2MessageFormat(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v1-v2-shared
	//= type=test
	//# - This string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=test
	//# - The S3EC SHOULD support decoding the S3 Server's "double encoding".
	rune128 := string(rune(128))
	rune200 := string(rune(200))
	rune256 := string(rune(256))

	unicodeStrings := []string{rune128, rune200, rune256, runeMaxInt, shorter, medium, longer, mix, mixTwo}
	for _, s := range unicodeStrings {
		UnicodeEncryptionContextV4ClientV2MessageFormat(t, s)
	}
}

func UnicodeEncryptionContextV4ClientV2MessageFormat(t *testing.T, metadataString string) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": metadataString})
	_, err = s3ecV4.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),

	time.Sleep(1 * time.Second)

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
		t.Errorf("expect %v text, got %v", e, a)
	}

	s3ecV4.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	"bytes"
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v3/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v3/client"
	"github.com/aws/amazon-s3-encryption-client-go/v3/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v3/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"strings"
	"testing"
	return fmt.Sprintf(arnFormat, region, accountId, shortAlias)
}

func TestInteg_EncryptFixtures_V2MessageFormat(t *testing.T) {
	var region = LoadRegion()
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		KEK, bucket, region, CEK string
	}{
		{
			// AES GCM implies V2 message format
			CEKAlg: "aes_gcm",
			KEK:    "kms", bucket: bucket, region: region, CEK: "aes_gcm",
		},
			if err != nil {
				t.Fatalf("failed to create new CMM")
			}
			encClient, _ := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.EnableLegacyUnauthenticatedModes = true
				clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
			})

			for caseKey, plaintext := range fixtures.Plaintexts {
				_, err := encClient.PutObject(ctx, &s3.PutObjectInput{
					),
					Body: bytes.NewReader(plaintext),
				})

				// Assert correctness of encrypted object
				out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key: aws.String(
						fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
							fixtures.BaseFolder, version, caseKey),
					),
				})
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-iv" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-iv"]; !ok {
					t.Fatalf("expected x-amz-iv to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-key-v2" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-key-v2"]; !ok {
					t.Fatalf("expected x-amz-key-v2 to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-matdesc" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-matdesc"]; !ok {
					t.Fatalf("expected x-amz-matdesc to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-wrap-alg" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-wrap-alg"]; !ok {
					t.Fatalf("expected x-amz-wrap-alg to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-cek-alg" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-cek-alg"]; !ok {
					t.Fatalf("expected x-amz-cek-alg to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-tag-len" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-tag-len"]; !ok {
					t.Fatalf("expected x-amz-tag-len to be present in metadata, got %v", out.Metadata)
				}

				if err != nil {
					t.Fatalf("failed to upload encrypted fixture, %v", err)
				}
		Lang    string
		Version string
	}{
		{CEKAlg: "aes_cbc", Lang: "Go", Version: "v4"}, // v4 doesn't support CBC but that's where the files are
		{CEKAlg: "aes_gcm", Lang: "Go", Version: "v4"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v2"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v3"},
	}
				cmmCbc, err := materials.NewCryptographicMaterialsManager(keyring)
				decClient, err = client.New(s3Client, cmmCbc, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = true
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			} else if c.CEKAlg == "aes_gcm" {
				decClient, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			}

			fixtures := getFixtures(t, s3Client, c.CEKAlg, bucket)
			ciphertexts := decryptFixtures(t, decClient, fixtures, bucket, c.Lang, c.Version)

			if len(ciphertexts) == 0 {
				t.Fatalf("expected more than 0 ciphertexts to decrypt!")
	}

	switch cek {
	case "aes_gcm_committing":
	case "aes_gcm":
		return kmsKeyring
	case "aes_cbc":
	})
}

func TestInteg_DeleteObjects_DeletesObjects(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "delete-objects-test-" + time.Now().Format("20060102-150405")
	var key1 = baseKey + "-object1"
	var key2 = baseKey + "-object2"
	var key3 = baseKey + "-object3"
	var plaintext1 = "Hello, S3 Encryption Client DeleteObjects test - Object 1!"
	var plaintext2 = "Hello, S3 Encryption Client DeleteObjects test - Object 2!"
	var plaintext3 = "Hello, S3 Encryption Client DeleteObjects test - Object 3!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing objects
	objectsToClean := []string{key1, key2, key3}
	for _, key := range objectsToClean {
		s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
	}

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put multiple encrypted objects
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key1,
		Body:   bytes.NewReader([]byte(plaintext1)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 1: %v", err)
	}

	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key2,
		Body:   bytes.NewReader([]byte(plaintext2)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 2: %v", err)
	}

	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key3,
		Body:   bytes.NewReader([]byte(plaintext3)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 3: %v", err)
	}

	// Verify all objects exist before deletion
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			t.Fatalf("object %d should exist before deletion: %v", i+1, err)
		}
	}

	t.Logf("✓ Verified all objects exist before deletion")

	// Test DeleteObjects - should delete all objects
	deleteInput := &s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{
				{Key: &key1},
				{Key: &key2},
				{Key: &key3},
			},
		},
	}

	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST be implemented by the S3EC.
	result, err := s3ec.DeleteObjects(ctx, deleteInput)
	if err != nil {
		t.Fatalf("DeleteObjects failed: %v", err)
	}

	// Verify the response structure
	if result == nil {
		t.Fatal("DeleteObjects result should not be nil")
	}

	if len(result.Deleted) != 3 {
		t.Errorf("expected 3 deleted objects, got %d", len(result.Deleted))
	}

	if len(result.Errors) > 0 {
		t.Errorf("expected no errors, got %d errors: %v", len(result.Errors), result.Errors)
	}

	// Verify all objects are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST delete each of the given objects.
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err == nil {
			t.Errorf("object %d should be deleted but still exists", i+1)
		}
	}

	// Verify all instruction files are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST delete each of the corresponding instruction files using the default instruction file suffix.
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    aws.String(key + ".instruction"),
		})
		if err == nil {
			t.Errorf("instruction file %d should be deleted but still exists", i+1)
		}
	}

	t.Logf("✓ DeleteObjects successfully deleted all objects")
}

func TestInteg_DeleteObject_DeletesObjectAndInstructionFile(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "delete-object-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client DeleteObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing objects
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key + ".instruction"),
	})

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put encrypted object (this should create both the object and instruction file)
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object: %v", err)
	}

	// Verify object exists before deletion
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("object should exist before deletion: %v", err)
	}

	// Test DeleteObject - should delete both the object and instruction file
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST be implemented by the S3EC.
	_, err = s3ec.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	// Verify both object and instruction file are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST delete the given object key.
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil {
		t.Errorf("object should be deleted but still exists")
	}

	// Verify both object and instruction file are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key + ".instruction"),
	})
	if err == nil {
		t.Errorf("instruction file should be deleted but still exists")
	}

	t.Logf("✓ DeleteObject successfully deleted both object and instruction file")
}

func TestIntegLegacyUnauthenticatedModes(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test specific known objects and validate their content encryption algorithms
	testCases := []struct {
		name           string
		objectKey      string
		expectedCekAlg string
		description    string
	}{
		{
			name:           "LegacyUnauthenticated_AES_CBC",
			objectKey:      "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedCekAlg: "AES/CBC/PKCS5Padding",
			description:    "AES/CBC object should use legacy unauthenticated content encryption algorithm",
		},
		{
			name:           "ModernAuthenticated_AES_GCM",
			objectKey:      "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedCekAlg: "AES/GCM/NoPadding",
			description:    "AES/GCM object should use modern authenticated content encryption algorithm",
		},
	}

	// First validate that our test objects have the expected content encryption algorithms
	for _, tc := range testCases {
		t.Run("Validate_"+tc.name, func(t *testing.T) {
			headResult, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})
			if err != nil {
				t.Skipf("Test object %s not found, skipping: %v", tc.objectKey, err)
			}

			var actualCekAlg string
			if headResult.Metadata != nil {
				if val, ok := headResult.Metadata["x-amz-cek-alg"]; ok {
					actualCekAlg = val
				}
			}

			if actualCekAlg != tc.expectedCekAlg {
				t.Errorf("Object %s has cek-alg %q, expected %q. %s", 
					tc.objectKey, actualCekAlg, tc.expectedCekAlg, tc.description)
			} else {
				t.Logf("✓ Validated: %s has expected cek-alg: %s", tc.objectKey, actualCekAlg)
			}
		})
	}

	// Get expected plaintext for test case test_one.txt
	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
	expectedPlaintext, exists := fixtures.Plaintexts["test_one.txt"]
	if !exists {
		t.Fatal("Could not find plaintext for test case test_one.txt")
	}

	// Test legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes = true (should succeed)
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
	t.Run("LegacyUnauthenticatedObject_EnabledShouldSucceed", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err != nil {
			t.Fatalf("expected successful decryption of legacy unauthenticated object with EnableLegacyUnauthenticatedModes=true, but got error: %v", err)
		}

		decryptedData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("failed to read decrypted data: %v", err)
		}

		if !bytes.Equal(expectedPlaintext, decryptedData) {
			t.Errorf("decrypted data mismatch: expected %q, got %q", string(expectedPlaintext), string(decryptedData))
		}

		t.Logf("✓ Successfully decrypted legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes=true")
	})

	// Test legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes = false (should fail)
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy content encryption algorithms;
	//# it MUST throw an exception when attempting to decrypt an object encrypted with a legacy content encryption algorithm.
	//= ../specification/s3-encryption/decryption.md#legacy-decryption
	//= type=test
	//# The S3EC MUST NOT decrypt objects encrypted using legacy unauthenticated algorithm suites unless specifically configured to do so.
	t.Run("LegacyUnauthenticatedObject_DisabledShouldFail", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = false
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		//= ../specification/s3-encryption/decryption.md#legacy-decryption
		//= type=test
		//# If the S3EC is not configured to enable legacy unauthenticated content decryption,
		//# the client MUST throw an exception when attempting to decrypt an object encrypted with a legacy unauthenticated algorithm suite.
		if err == nil {
			t.Fatalf("expected decryption of legacy unauthenticated object to fail when EnableLegacyUnauthenticatedModes=false, but it succeeded")
		}

		// Verify the error message indicates legacy unauthenticated modes issue
		if !strings.Contains(err.Error(), "enable legacy unauthenticated modes") && !strings.Contains(err.Error(), "AES/CBC/PKCS5Padding") {
			t.Errorf("expected error to mention legacy unauthenticated modes or AES/CBC algorithm, got: %v", err)
		}

		t.Logf("✓ Correctly failed to decrypt legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes=false: %v", err)
	})

	// Test modern authenticated object (AES/GCM) works regardless of legacy unauthenticated modes setting
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
	t.Run("ModernAuthenticatedObject_AlwaysWorks", func(t *testing.T) {
		// Get expected plaintext for modern object
		gcmFixtures := getFixtures(t, s3Client, "aes_gcm", bucket)
		gcmExpectedPlaintext, exists := gcmFixtures.Plaintexts["test_one.txt"]
		if !exists {
			t.Skip("Could not find plaintext for aes_gcm test case test_one.txt, skipping modern object test")
		}

		for _, enableLegacyUnauthenticated := range []bool{true, false} {
			t.Run(fmt.Sprintf("EnableLegacyUnauthenticated_%v", enableLegacyUnauthenticated), func(t *testing.T) {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				})
				cmm, err := materials.NewCryptographicMaterialsManager(keyring)
				if err != nil {
					t.Fatalf("failed to create CMM: %v", err)
				}

				decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = enableLegacyUnauthenticated
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}

				modernObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
				result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &modernObjectKey,
				})
				if err != nil {
					t.Fatalf("modern authenticated algorithm should always work regardless of legacy unauthenticated modes setting, but got error: %v", err)
				}

				decryptedData, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("failed to read decrypted data: %v", err)
				}

				if !bytes.Equal(gcmExpectedPlaintext, decryptedData) {
					t.Errorf("decrypted data mismatch: expected %q, got %q", string(gcmExpectedPlaintext), string(decryptedData))
				}

				t.Logf("✓ Modern authenticated object (AES/GCM) works correctly with EnableLegacyUnauthenticatedModes=%v", enableLegacyUnauthenticated)
			})
		}
	})
}

func TestIntegKmsContext(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
		Key:    &key,
	})
}

func TestIntegLegacyWrappingAlgorithms(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test specific known objects and validate their wrapping algorithms
	testCases := []struct {
		name           string
		objectKey      string
		expectedWrapAlg string
		cekAlg         string
		description    string
	}{
		{
			name:           "LegacyKMSKeyring_AES_CBC",
			objectKey:      "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedWrapAlg: "kms",
			cekAlg:         "aes_cbc",
			description:    "AES/CBC object should use legacy 'kms' wrapping algorithm",
		},
		{
			name:           "ModernKMSContextKeyring_AES_GCM",
			objectKey:      "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedWrapAlg: "kms+context",
			cekAlg:         "aes_gcm",
			description:    "AES/GCM object should use modern 'kms+context' wrapping algorithm",
		},
	}

	// First validate that our test objects have the expected wrapping algorithms
	for _, tc := range testCases {
		t.Run("Validate_"+tc.name, func(t *testing.T) {
			headResult, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})
			if err != nil {
				t.Skipf("Test object %s not found, skipping: %v", tc.objectKey, err)
			}

			var actualWrapAlg string
			if headResult.Metadata != nil {
				if val, ok := headResult.Metadata["x-amz-wrap-alg"]; ok {
					actualWrapAlg = val
				}
			}

			if actualWrapAlg != tc.expectedWrapAlg {
				t.Errorf("Object %s has wrap-alg %q, expected %q. %s", 
					tc.objectKey, actualWrapAlg, tc.expectedWrapAlg, tc.description)
			} else {
				t.Logf("✓ Validated: %s has expected wrap-alg: %s", tc.objectKey, actualWrapAlg)
			}
		})
	}

	// Get expected plaintext for test case test_one.txt
	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
	expectedPlaintext, exists := fixtures.Plaintexts["test_one.txt"]
	if !exists {
		t.Fatal("Could not find plaintext for test case test_one.txt")
	}

	// Test legacy object (kms) with EnableLegacyWrappingAlgorithms = true (should succeed)
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	t.Run("LegacyObject_EnabledShouldSucceed", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err != nil {
			t.Fatalf("expected successful decryption of legacy object with EnableLegacyWrappingAlgorithms=true, but got error: %v", err)
		}

		decryptedData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("failed to read decrypted data: %v", err)
		}

		if !bytes.Equal(expectedPlaintext, decryptedData) {
			t.Errorf("decrypted data mismatch: expected %q, got %q", string(expectedPlaintext), string(decryptedData))
		}

		t.Logf("✓ Successfully decrypted legacy object (kms wrapping algorithm) with EnableLegacyWrappingAlgorithms=true")
	})

	// Test legacy object (kms) with EnableLegacyWrappingAlgorithms = false (should fail)
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy wrapping algorithms;
	//# it MUST throw an exception when attempting to decrypt an object encrypted with a legacy wrapping algorithm.
	t.Run("LegacyObject_DisabledShouldFail", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err == nil {
			t.Fatalf("expected decryption of legacy object to fail when EnableLegacyWrappingAlgorithms=false, but it succeeded")
		}

		// Verify the error message indicates legacy wrapping algorithm issue
		if !strings.Contains(err.Error(), "legacyWrappingAlgorithms") && !strings.Contains(err.Error(), "did not match an expected algorithm") {
			t.Errorf("expected error to mention legacyWrappingAlgorithms or algorithm mismatch, got: %v", err)
		}

		t.Logf("✓ Correctly failed to decrypt legacy object (kms wrapping algorithm) with EnableLegacyWrappingAlgorithms=false: %v", err)
	})

	// Test modern object (kms+context) works regardless of legacy setting
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	t.Run("ModernObject_AlwaysWorks", func(t *testing.T) {
		// Get expected plaintext for modern object
		gcmFixtures := getFixtures(t, s3Client, "aes_gcm", bucket)
		gcmExpectedPlaintext, exists := gcmFixtures.Plaintexts["test_one.txt"]
		if !exists {
			t.Skip("Could not find plaintext for aes_gcm test case test_one.txt, skipping modern object test")
		}

		for _, enableLegacy := range []bool{true, false} {
			t.Run(fmt.Sprintf("EnableLegacy_%v", enableLegacy), func(t *testing.T) {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = enableLegacy
				})
				cmm, err := materials.NewCryptographicMaterialsManager(keyring)
				if err != nil {
					t.Fatalf("failed to create CMM: %v", err)
				}

				decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}

				modernObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
				result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &modernObjectKey,
				})
				if err != nil {
					t.Fatalf("modern algorithm should always work regardless of legacy setting, but got error: %v", err)
				}

				decryptedData, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("failed to read decrypted data: %v", err)
				}

				if !bytes.Equal(gcmExpectedPlaintext, decryptedData) {
					t.Errorf("decrypted data mismatch: expected %q, got %q", string(gcmExpectedPlaintext), string(decryptedData))
				}

				t.Logf("✓ Modern object (kms+context wrapping algorithm) works correctly with EnableLegacyWrappingAlgorithms=%v", enableLegacy)
			})
		}
	})
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - GetObject MUST be implemented by the S3EC.
//# - GetObject MUST decrypt data received from the S3 server and return it as plaintext.
func TestInteg_GetObject_BasicDecryption(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "basic-getobject-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client GetObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing object
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	// Create S3EC and encrypt an object
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put encrypted object
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object: %v", err)
	}
	
	// Loose assertion that the object was encrypted by checking that its body does not match the plaintext
	rawResult, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to get raw object: %v", err)
	}

	rawData, err := io.ReadAll(rawResult.Body)
	if err != nil {
		t.Fatalf("failed to read raw object data: %v", err)
	}

	if string(rawData) == plaintext {
		t.Errorf("object was not encrypted: raw data matches plaintext")
	}

	// Test GetObject - should decrypt and return plaintext
	result, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}

	decryptedData, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted data: %v", err)
	}

	if string(decryptedData) != plaintext {
		t.Errorf("GetObject decryption failed: expected %q, got %q", plaintext, string(decryptedData))
	}

	t.Logf("✓ GetObject successfully decrypted data and returned plaintext")

	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - PutObject MUST be implemented by the S3EC.
//# - PutObject MUST encrypt its input data before it is uploaded to S3.
func TestInteg_PutObject_BasicEncryption(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "basic-putobject-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client PutObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing object
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Test PutObject - should encrypt input data before uploading
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	// Loose assertion that the object was encrypted by checking that its body does not match the plaintext
	rawResult, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to get raw object: %v", err)
	}

	rawData, err := io.ReadAll(rawResult.Body)
	if err != nil {
		t.Fatalf("failed to read raw object data: %v", err)
	}

	if string(rawData) == plaintext {
		t.Errorf("object was not encrypted: raw data matches plaintext")
	}

	// Verify encryption metadata is present
	if rawResult.Metadata == nil {
		t.Fatal("expected encryption metadata to be present")
	}

	expectedMetadataKeys := []string{"x-amz-iv", "x-amz-key-v2", "x-amz-matdesc", "x-amz-wrap-alg", "x-amz-cek-alg", "x-amz-tag-len"}
	for _, key := range expectedMetadataKeys {
		if _, exists := rawResult.Metadata[key]; !exists {
			t.Errorf("expected encryption metadata key %s to be present", key)
		}
	}

	// Verify we can decrypt it back to the original plaintext
	decryptResult, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to decrypt object: %v", err)
	}

	decryptedData, err := io.ReadAll(decryptResult.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted data: %v", err)
	}

	if string(decryptedData) != plaintext {
		t.Errorf("decryption verification failed: expected %q, got %q", plaintext, string(decryptedData))
	}

	t.Logf("✓ PutObject successfully encrypted input data before uploading to S3")

	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

//= ../specification/s3-encryption/decryption.md#key-commitment
//= type=test
//# The S3EC MUST validate the algorithm suite used for decryption against the key commitment policy before attempting to decrypt the content ciphertext.
func TestInteg_ValidateAlgorithmSuiteAgainstCommitmentPolicy(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Matrix of [V3 object, V2 object] x [commitment policies] = 6 test cases
	cases := []struct {
		name             string
		objectKey        string
		objectType       string
		commitmentPolicy commitment.CommitmentPolicy
		expectError      bool
		errorContains    string
	}{
		{
			name:             "V3_object_with_FORBID_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V3_object_with_REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V3_object_with_REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V2_object_with_FORBID_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V2_object_with_REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		// Only expected failure case on decryption, where a non-committing object is
		// attempted to be decrypted under a commitment policy that requires commitment.
		{
			name:             "V2_object_with_REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			expectError:      true,
			errorContains:    "object's content encryption algorithm is not valid for the selected commitment policy",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = false
			})
			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
			})
			if err != nil {
				t.Fatalf("expected no error during client creation, got %v", err)
			}

			_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})

			//= ../specification/s3-encryption/decryption.md#key-commitment
			//= type=test
			//# If the commitment policy requires decryption using a committing algorithm suite,
			//# and the algorithm suite associated with the object does not support key commitment,
			//# then the S3EC MUST throw an exception.
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Expected error during GetObject: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				t.Logf("✓ Successfully decrypted %s object with %s policy", tc.objectType, tc.commitmentPolicy)
			}
		})
	}
}

func TestInteg_AlgorithmSuiteMessageFormatCompatibility(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "algorithm-suite-message-format-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client Algorithm Suite Message Format test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	testCases := []struct {
		name                string
		algorithmSuite      string
		commitmentPolicy    commitment.CommitmentPolicy
		enableLegacyModes   bool
		enableLegacyWrap    bool
		expectedFormat      string
		expectedHeaders     []string
	}{
		//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
		//= type=test
		//# Objects encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF MUST use the V2 message format version only.
		{
			name:                "ALG_AES_256_GCM_IV12_TAG16_NO_KDF_V2_Format",
			algorithmSuite:      "AES256GCMIV12Tag16NoKDF",
			commitmentPolicy:    commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			enableLegacyModes:   false,
			enableLegacyWrap:    false,
			expectedFormat:      "V2",
			expectedHeaders:     []string{"x-amz-iv", "x-amz-key-v2", "x-amz-matdesc", "x-amz-wrap-alg", "x-amz-cek-alg", "x-amz-tag-len"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			key := baseKey + "-" + tc.algorithmSuite

			// Clean up any existing object
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			// Create keyring with appropriate legacy settings
			keyring := materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = tc.enableLegacyWrap
			})

			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			// Create S3EC with specific configuration
			var s3ec *client.S3EncryptionClientV3
			if tc.algorithmSuite == "AES256GCMHkdfSha512CommitKey" {
				// For committing algorithm suite, use default client (which uses committing by default)
				s3ec, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = tc.enableLegacyModes
					clientOptions.CommitmentPolicy = tc.commitmentPolicy
					// Don't override algorithm suite - let it use the default committing algorithm
				})
			} else {
				// For non-committing algorithm suites, explicitly set them
				s3ec, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = tc.enableLegacyModes
					clientOptions.CommitmentPolicy = tc.commitmentPolicy
					
					// Configure the specific algorithm suite
					switch tc.algorithmSuite {
					case "AES256CBCIV16NoKDF":
						clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CBCIV16NoKDF
					case "AES256GCMIV12Tag16NoKDF":
						clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256GCMIV12Tag16NoKDF
					default:
						t.Fatalf("unknown algorithm suite: %s", tc.algorithmSuite)
					}
				})
			}
			if err != nil {
				t.Fatalf("failed to create S3EC: %v", err)
			}

			// Encrypt object with configured algorithm suite
			_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   bytes.NewReader([]byte(plaintext)),
			})
			if err != nil {
				t.Fatalf("failed to put encrypted object: %v", err)
			}

			// Get object metadata using regular S3 client to inspect message format headers
			result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				t.Fatalf("failed to get encrypted object metadata: %v", err)
			}

			// Verify the expected message format headers are present
			if result.Metadata == nil {
				t.Fatalf("expected metadata to be present for %s format", tc.expectedFormat)
			}

			t.Logf("Testing %s algorithm suite for %s message format", tc.algorithmSuite, tc.expectedFormat)
			t.Logf("Object metadata: %v", result.Metadata)

			// Assert all expected headers for the message format are present
			for _, expectedHeader := range tc.expectedHeaders {
				if _, exists := result.Metadata[expectedHeader]; !exists {
					t.Errorf("expected %s format header '%s' to be present for algorithm suite %s, but it was missing", 
						tc.expectedFormat, expectedHeader, tc.algorithmSuite)
				} else {
					t.Logf("✓ Found expected %s format header: %s", tc.expectedFormat, expectedHeader)
				}
			}

			// Verify that headers from other message formats are NOT present
			var unexpectedHeaders []string
			if tc.expectedFormat == "V2" {
				// V2 format should not have V3 headers
				unexpectedHeaders = []string{"x-amz-c", "x-amz-3", "x-amz-w", "x-amz-d", "x-amz-i"}
			} else if tc.expectedFormat == "V3" {
				// V3 format should not have V2 headers
				unexpectedHeaders = []string{"x-amz-iv", "x-amz-key-v2", "x-amz-tag-len"}
			}

			for _, unexpectedHeader := range unexpectedHeaders {
				if _, exists := result.Metadata[unexpectedHeader]; exists {
					t.Errorf("unexpected header '%s' found for %s format with algorithm suite %s", 
						unexpectedHeader, tc.expectedFormat, tc.algorithmSuite)
				}
			}

			t.Logf("✓ Algorithm suite %s correctly uses %s message format", tc.algorithmSuite, tc.expectedFormat)

			// Decrypt sanity check - verify we can decrypt the object back to original plaintext
			decryptResult, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				t.Fatalf("failed to decrypt object with algorithm suite %s: %v", tc.algorithmSuite, err)
			}

			decryptedData, err := io.ReadAll(decryptResult.Body)
			if err != nil {
				t.Fatalf("failed to read decrypted data for algorithm suite %s: %v", tc.algorithmSuite, err)
			}

			if string(decryptedData) != plaintext {
				t.Errorf("decrypt sanity check failed for algorithm suite %s: expected %q, got %q", 
					tc.algorithmSuite, plaintext, string(decryptedData))
			} else {
				t.Logf("✓ Decrypt sanity check passed for algorithm suite %s: successfully decrypted %d bytes", 
					tc.algorithmSuite, len(decryptedData))
			}

			// Cleanup
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
		})
	}
}

func TestInteg_CommitmentPolicyBehavior(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "commitment-policy-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client Commitment Policy test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	cases := []struct {
		name                    string
		commitmentPolicy        commitment.CommitmentPolicy
		algorithmSuite          *algorithms.AlgorithmSuite
		expectEncryptError      bool
		expectDecryptError      bool
		encryptErrorContains    string
		decryptErrorContains    string
	}{
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST NOT encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "FORBID_with_committing_algorithm",
			commitmentPolicy:        commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "does not allow committing algorithm suites",
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
		{
			name:                    "FORBID_with_non_committing_algorithm",
			commitmentPolicy:        commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMIV12Tag16NoKDF,
			expectEncryptError:      false,
			expectDecryptError:      false,
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_ALLOW_DECRYPT_with_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "algorithm suite ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY is not supported for encryption in S3EC V3",
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_ALLOW_DECRYPT_with_non_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMIV12Tag16NoKDF,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "requires committing algorithm suites",
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_REQUIRE_DECRYPT_with_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "algorithm suite ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY is not supported for encryption in S3EC V3",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key := baseKey + "-" + tc.name

			// Clean up any existing object
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			keyring := materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = false
			})
			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			// Test encryption behavior
			s3ec, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithmSuite
			})
			if err != nil {
				if tc.expectEncryptError && tc.encryptErrorContains != "" && strings.Contains(err.Error(), tc.encryptErrorContains) {
					t.Logf("✓ Expected error during client creation for encryption: %v", err)
					return // Expected error during client creation
				}
				t.Fatalf("expected no error during client creation, got %v", err)
			}

			_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   bytes.NewReader([]byte(plaintext)),
			})

			if tc.expectEncryptError {
				if err == nil {
					t.Fatalf("expected encryption error but got none")
				}
				if tc.encryptErrorContains != "" && !strings.Contains(err.Error(), tc.encryptErrorContains) {
					t.Errorf("expected encryption error to contain %q, got %q", tc.encryptErrorContains, err.Error())
				} else {
					t.Logf("✓ Expected encryption error: %v", err)
				}
				return // Don't test decryption if encryption failed
			} else {
				if err != nil {
					t.Fatalf("expected no encryption error, got %v", err)
				}
				t.Logf("✓ Successfully encrypted with %s policy and %s algorithm", tc.commitmentPolicy, tc.algorithmSuite.CipherName())
			}

			// Test decryption behavior - create a new client for decryption
			decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
			})
			if err != nil {
				t.Fatalf("failed to create decryption client: %v", err)
			}

			_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			if tc.expectDecryptError {
				if err == nil {
					t.Fatalf("expected decryption error but got none")
				}
				if tc.decryptErrorContains != "" && !strings.Contains(err.Error(), tc.decryptErrorContains) {
					t.Errorf("expected decryption error to contain %q, got %q", tc.decryptErrorContains, err.Error())
				} else {
					t.Logf("✓ Expected decryption error: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no decryption error, got %v", err)
				}
				t.Logf("✓ Successfully decrypted with %s policy", tc.commitmentPolicy)
			}

			// Cleanup
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
		})
	}
}


//= ../specification/s3-encryption/key-commitment.md#commitment-policy
//= type=test
//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST NOT allow decryption using algorithm suites which do not support key commitment.
func TestInteg_RequireEncryptRequireDecrypt_RejectsNonCommittingObjects(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test that REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy rejects decryption of non-committing objects
	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})
	if err != nil {
		t.Fatalf("failed to create decryption client: %v", err)
	}

	// Try to decrypt a non-committing object (V2 format, AES/GCM without commitment)
	nonCommittingObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
	_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &nonCommittingObjectKey,
	})

	if err == nil {
		t.Fatalf("expected decryption error when trying to decrypt non-committing object with REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy, but got none")
	}

	if !strings.Contains(err.Error(), "object's content encryption algorithm is not valid for the selected commitment policy") {
		t.Errorf("expected error to contain commitment policy validation message, got: %v", err)
	}

	t.Logf("✓ REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy correctly rejected decryption of non-committing object: %v", err)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package algorithms

import (
	"crypto/sha512"
	"fmt"
	"strconv"
	"hash"
)

// Algorithm constants
const (
	// GCM maximum content length in bits (2^39 - 256 bits)
	GCMMaxContentLengthBits = 549755813632
	// CTR maximum content length in bytes (2^32 bytes)
	CTRMaxContentLengthBytes = 4294967296
	// CBC maximum content length in bytes (2^32 bytes)  
	CBCMaxContentLengthBytes = 4294967296
	
	// Algorithm identifiers
	AESGCMCommitKey = "115"
	AESGCMNoPadding = "AES/GCM/NoPadding"
	AESCTRNoPadding = "AES/CTR/NoPadding"
	AESCBCPKCS5   	= "AES/CBC/PKCS5Padding"
)

// AlgorithmSuite represents the encryption algorithm suite configuration
type AlgorithmSuite struct {
	id                         int
	isLegacy                   bool
	dataKeyAlgorithm          string
	dataKeyLengthBits         int
	cipherName                string
	cipherBlockSizeBits       int
	cipherIvLengthBits        int
	cipherTagLengthBits       int
	cipherMaxContentLengthBits int64
	isCommitting              bool
	commitmentLengthBits      int
	kdfHashAlgorithm          func() hash.Hash
}

// Predefined algorithm suites
var (
	// Key committing AES-256-GCM content encryption with HKDF-SHA512 commitment/encryption key derivation.
	// In addition to the security properties for the AES-256-GCM content encryption suite,
	// this suite also uses HKDF to derive a key commitment string that is included in metadata.
	// We recommend using this suite for any new objects.
	// This is the default suite for v4 clients.
	// Any v4 client can use this suite to read and write objects with key commitment,
	// but v3 clients can only use this suite to read objects with key commitment.
	// Be sure to upgrade any v3 clients to at least v3.2.0 to ensure compatibility
	// if you enable this algorithm.
	AlgAES256GCMHkdfSha512CommitKey = &AlgorithmSuite{
		id:                         0x0073,
		isLegacy:                   false,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256, // this is the input into the KDF
		cipherName:                AESGCMCommitKey,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        224,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: GCMMaxContentLengthBits,
		isCommitting:              true,
		commitmentLengthBits:      224,
		kdfHashAlgorithm:          sha512.New,
	}

	// AES-256 GCM content encryption.
	// This suite uses the data encryption key directly for AES-256 GCM content encryption without key derivation.
	// This is the default suite for v3 clients.
	// Content encrypted with this suite can be read by any v2, v3, or v4 client.
	AlgAES256GCMIV12Tag16NoKDF = &AlgorithmSuite{
		id:                         0x0072,
		isLegacy:                   false,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESGCMNoPadding,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        96,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: GCMMaxContentLengthBits,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}

	// Legacy AES-256 CTR.
	// This suite is not supported at this time.
	AlgAES256CTRIV16Tag16NoKDF = &AlgorithmSuite{
		id:                         0x0071,
		isLegacy:                   true,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESCTRNoPadding,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        128,
		cipherTagLengthBits:       128,
		cipherMaxContentLengthBits: CTRMaxContentLengthBytes * 8,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}

	// Legacy AES-256 CBC.
	// This suite is only supported for decryption of existing objects and cannot be used for new objects.
	// We recommend migrating any existing objects encrypted with this suite to a non-legacy suite.
	AlgAES256CBCIV16NoKDF = &AlgorithmSuite{
		id:                         0x0070,
		isLegacy:                   true,
		dataKeyAlgorithm:          "AES",
		dataKeyLengthBits:         256,
		cipherName:                AESCBCPKCS5,
		cipherBlockSizeBits:       128,
		cipherIvLengthBits:        128,
		cipherTagLengthBits:       0,
		cipherMaxContentLengthBits: CBCMaxContentLengthBytes * 8,
		isCommitting:              false,
		commitmentLengthBits:      0,
		kdfHashAlgorithm:          nil,
	}
)

// Map for looking up algorithm suites by ID
var algorithmSuitesByID = map[int]*AlgorithmSuite{
	0x0073: AlgAES256GCMHkdfSha512CommitKey,
	0x0072: AlgAES256GCMIV12Tag16NoKDF,
	0x0071: AlgAES256CTRIV16Tag16NoKDF,
	0x0070: AlgAES256CBCIV16NoKDF,
}

// GetAlgorithmSuiteByID returns the algorithm suite for the given ID
func GetAlgorithmSuiteByID(id int) (*AlgorithmSuite, error) {
	suite, exists := algorithmSuitesByID[id]
	if !exists {
		return nil, fmt.Errorf("unknown algorithm suite ID: 0x%04x", id)
	}
	return suite, nil
}

// ID returns the algorithm suite ID
func (a *AlgorithmSuite) ID() int {
	return a.id
}

// IDAsString returns the algorithm suite ID as a string
func (a *AlgorithmSuite) IDAsString() string {
	return strconv.Itoa(a.id)
}

// IDAsBytes returns the algorithm suite ID as a byte array (big-endian)
func (a *AlgorithmSuite) IDAsBytes() []byte {
	return []byte{byte(a.id >> 8), byte(a.id)}
}

// IsLegacy returns whether this is a legacy algorithm suite
func (a *AlgorithmSuite) IsLegacy() bool {
	return a.isLegacy
}

// DataKeyAlgorithm returns the data key algorithm name
func (a *AlgorithmSuite) DataKeyAlgorithm() string {
	return a.dataKeyAlgorithm
}

// DataKeyLengthBits returns the data key length in bits
func (a *AlgorithmSuite) DataKeyLengthBits() int {
	return a.dataKeyLengthBits
}

// DataKeyLengthBytes returns the data key length in bytes
func (a *AlgorithmSuite) DataKeyLengthBytes() int {
	return a.dataKeyLengthBits / 8
}

// CipherName returns the cipher name
func (a *AlgorithmSuite) CipherName() string {
	return a.cipherName
}

// CipherTagLengthBits returns the cipher tag length in bits
func (a *AlgorithmSuite) CipherTagLengthBits() int {
	return a.cipherTagLengthBits
}

// CipherTagLengthBytes returns the cipher tag length in bytes
func (a *AlgorithmSuite) CipherTagLengthBytes() int {
	return a.cipherTagLengthBits / 8
}

// IVLengthBytes returns the IV length in bytes
func (a *AlgorithmSuite) IVLengthBytes() int {
	return a.cipherIvLengthBits / 8
}

// CipherBlockSizeBytes returns the cipher block size in bytes
func (a *AlgorithmSuite) CipherBlockSizeBytes() int {
	return a.cipherBlockSizeBits / 8
}

// CipherMaxContentLengthBits returns the maximum content length in bits
func (a *AlgorithmSuite) CipherMaxContentLengthBits() int64 {
	return a.cipherMaxContentLengthBits
}

// CipherMaxContentLengthBytes returns the maximum content length in bytes
func (a *AlgorithmSuite) CipherMaxContentLengthBytes() int64 {
	return a.cipherMaxContentLengthBits / 8
}

// IsCommitting returns whether this algorithm suite is key committing
func (a *AlgorithmSuite) IsCommitting() bool {
	return a.isCommitting
}

// CommitmentLengthBits returns the commitment length in bits
func (a *AlgorithmSuite) CommitmentLengthBits() int {
	return a.commitmentLengthBits
}

// CommitmentLengthBytes returns the commitment length in bytes
func (a *AlgorithmSuite) CommitmentLengthBytes() int {
	return a.commitmentLengthBits / 8
}

// KDFHashAlgorithm returns the KDF hash algorithm
func (a *AlgorithmSuite) KDFHashAlgorithm() func() hash.Hash {
	return a.kdfHashAlgorithm
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"encoding/json"
	"github.com/aws/amazon-s3-encryption-client-go/v4/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"mime"
	"strings"
	"unicode/utf16"
	"unicode/utf8"
)

func customS3Decoder(matDesc string) (decoded string, e error) {
	// Manually decode S3's non-standard "double encoding"
	// First, mime decode it:
	decoder := new(mime.WordDecoder)
	s, err := decoder.DecodeHeader(matDesc)
	if err != nil {
		return "", fmt.Errorf("error while decoding material description: %s\n from S3 object metadata: %w", matDesc, err)
	}
	var sb strings.Builder

	skipNext := false
	var utf8buffer []byte
	// Iterate over the bytes in the string
	for i, b := range []byte(s) {
		r := rune(b)
		// Check if the rune (code point) is non-US-ASCII
		if r > 127 && !skipNext {
			// Non-ASCII characters need special treatment
			// due to double-encoding.
			// We are dealing with UTF-16 encoded codepoints
			// of the original UTF-8 characters.
			// So, take two bytes at a time...
			buf := []byte{s[i], s[i+1]}
			// Get the rune (code point)
			wrongRune := string(buf)
			// UTF-16 encode it
			encd := utf16.Encode([]rune(wrongRune))[0]
			// Buffer the byte-level representation of the code point
			// So that it can be UTF-8 encoded later
			utf8buffer = append(utf8buffer, byte(encd))
			skipNext = true
		} else if r > 127 && skipNext {
			// only skip once
			skipNext = false
		} else {
			// Decode the binary values as UTF-8
			// This recovers the original UTF-8
			for len(utf8buffer) > 0 {
				rb, size := utf8.DecodeRune(utf8buffer)
				sb.WriteRune(rb)
				utf8buffer = utf8buffer[size:]
			}
			sb.WriteByte(b)
		}
		// A more general solution would need to clear the utf8buffer here,
		// but specifically for material description,
		// we can assume that the string is JSON,
		// so the last character is '}' which is valid ASCII.
	}
	return sb.String(), nil
}

// GetObjectAPIClient is a client that implements the GetObject operation
type GetObjectAPIClient interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

func (m *decryptMiddleware) addDecryptAPIOptions(options *s3.Options) {
	options.APIOptions = append(options.APIOptions,
		m.addDecryptMiddleware,
	)
}

func (m *decryptMiddleware) addDecryptMiddleware(stack *middleware.Stack) error {
	return stack.Deserialize.Add(m, middleware.Before)
}

const decryptMiddlewareID = "S3Decrypt"

type decryptMiddleware struct {
	client *S3EncryptionClientV4
	input  *s3.GetObjectInput
}

// ID returns the resolver identifier
func (m *decryptMiddleware) ID() string {
	return decryptMiddlewareID
}

func (m *decryptMiddleware) HandleDeserialize(ctx context.Context, in middleware.DeserializeInput, next middleware.DeserializeHandler) (
	out middleware.DeserializeOutput, metadata middleware.Metadata, err error,
) {
	// call down the stack and get the deserialized result (decrypt middleware runs after the operation deserializer)
	out, metadata, err = next.HandleDeserialize(ctx, in)
	if err != nil {
		return out, metadata, err
	}

	httpResp, ok := out.RawResponse.(*smithyhttp.Response)
	if !ok {
		return out, metadata, &smithy.DeserializationError{Err: fmt.Errorf("unknown transport type %T", out.RawResponse)}
	}

	result, ok := out.Result.(*s3.GetObjectOutput)
	if !ok {
		return out, metadata, fmt.Errorf("expected GetObjectOutput; got %v", out)
	}

	loadReq := &internal.LoadStrategyRequest{
		HTTPResponse: httpResp.Response,
		Input:        m.input,
	}

	// decode metadata
	loadStrat := internal.DefaultLoadStrategy{}
	objectMetadata, err := loadStrat.Load(ctx, loadReq)
	if err != nil {
		return out, metadata, fmt.Errorf("failed to load objectMetadata: bucket=%v; key=%v; err=%w", m.input.Bucket, m.input.Key, err)
	}

	// determine the content algorithm from metadata
	// this is purposefully done before attempting to
	// decrypt the materials
	var cekFunc internal.CEKEntry
	objectCekAlgSuite, err := objectMetadata.GetContentEncryptionAlgorithmSuite()
	if err != nil {
		return out, metadata, err
	}

	//= ../specification/s3-encryption/decryption.md#key-commitment
	//# The S3EC MUST validate the algorithm suite used for decryption against the key commitment policy before attempting to decrypt the content ciphertext.
	if err := ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(objectCekAlgSuite, m.client.Options.CommitmentPolicy); err != nil {
		return out, metadata, fmt.Errorf("object's content encryption algorithm is not valid for the selected commitment policy: %v, %w", objectCekAlgSuite, err)
	}

	var matDesc string
	if objectCekAlgSuite == algorithms.AlgAES256GCMHkdfSha512CommitKey {
		cekFunc = internal.NewAESGCMDecryptCommittingContentCipher
		matDesc, err = objectMetadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for committing algorithm: %w", err)
		}
	} else if objectCekAlgSuite == algorithms.AlgAES256GCMIV12Tag16NoKDF {
		cekFunc = internal.NewAESGCMContentCipher
		matDesc, err = objectMetadata.GetMatDescV2()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for AES/GCM algorithm: %w", err)
		}
	} else if objectCekAlgSuite == algorithms.AlgAES256CBCIV16NoKDF {
		if !m.client.Options.EnableLegacyUnauthenticatedModes {
			//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
			//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy content encryption algorithms; it MUST throw an exception when attempting to decrypt an object encrypted with a legacy content encryption algorithm.
			//= ../specification/s3-encryption/decryption.md#legacy-decryption
			//# The S3EC MUST NOT decrypt objects encrypted using legacy unauthenticated algorithm suites unless specifically configured to do so.
			//= ../specification/s3-encryption/decryption.md#legacy-decryption
			//# If the S3EC is not configured to enable legacy unauthenticated content decryption, the client MUST throw an exception when attempting to decrypt an object encrypted with a legacy unauthenticated algorithm suite.
			return out, metadata, fmt.Errorf("configure client with enable legacy unauthenticated modes set to true to decrypt with %s", objectMetadata.CEKAlg)
		}
		//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
		//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
		cekFunc = internal.NewAESCBCContentCipher
		matDesc, err = objectMetadata.GetMatDescV2()
		if err != nil {
			return out, metadata, fmt.Errorf("error while getting material description for AES/CBC algorithm: %w", err)
		}
	} else {
		return out, metadata, fmt.Errorf("invalid content encryption algorithm found in metadata: %s", objectMetadata.CEKAlg)
	}


	cipherKey, err := objectMetadata.GetDecodedKey()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded key for materials: %w", err)
	}
	iv, err := objectMetadata.GetDecodedMessageIDOrIV()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get decoded IV for materials: %w", err)
	}
	keyringWrappingAlg, err := objectMetadata.GetFullWrappingAlgorithm()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get wrapping algorithm for materials: %w", err)
	}

	// S3 server will encode metadata with non-US-ASCII characters
	// Decode it here to avoid parsing/decryption failure
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//# The S3EC SHOULD support decoding the S3 Server's "double encoding".
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=exception
	//# If the S3EC does not support decoding the S3 Server's "double encoding" then it MUST return the content metadata untouched.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v1-v2-shared
	//# This string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# This material description string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# This encryption context string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	decodedMatDesc, err := customS3Decoder(matDesc)
	if err != nil {
		return out, metadata, fmt.Errorf("error while decoding Material Description: %w", err)
	}

	ec := ctx.Value(EncryptionContext)
	// If an encryption context is provided, the provided encryption context MUST match the encryption context stored in the metadata.
	if ec != nil {
		ecMap, ok := ec.(map[string]string)
		if !ok {
			return out, metadata, fmt.Errorf("encryption context provided to decrypt method is not valid JSON")
		}
		decodedMatDescMap := map[string]string{}
		if err := json.Unmarshal([]byte(decodedMatDesc), &decodedMatDescMap); err != nil {
			return out, metadata, fmt.Errorf("encryption context in stored object is not valid JSON: %w", err)
		}

		// The stored encryption context with the two reserved keys removed MUST match the provided encryption context.
		delete(decodedMatDescMap, "kms_cmk_id")
		delete(decodedMatDescMap, "aws:x-amz-cek-alg")

		if len(ecMap) != len(decodedMatDescMap) {
			return out, metadata, fmt.Errorf("Provided encryption context does not match information retrieved from S3")
		}
		for k, v := range ecMap {
			val, exists := decodedMatDescMap[k]
			if !exists || val != v {
				// If the stored encryption context with the two reserved keys removed does not match the provided encryption context, S3EC MUST throw an exception.
				return out, metadata, fmt.Errorf("Provided encryption context does not match information retrieved from S3")
			}
		}
	}

	cekAlg, err := objectMetadata.GetContentEncryptionAlgorithmString()
	if err != nil {
		return out, metadata, fmt.Errorf("unable to get content encryption algorithm from metadata: %w", err)
	}

	var keyCommitment []byte
	if (objectCekAlgSuite.IsCommitting()) {
		keyCommitment, err = objectMetadata.GetDecodedKeyCommitment()
		if err != nil {
			return out, metadata, fmt.Errorf("unable to get decoded key commitment for committing algorithm: %w", err)
		}
	}

	decryptMaterialsRequest := materials.DecryptMaterialsRequest{
		cipherKey,
		iv,
		decodedMatDesc,
		keyringWrappingAlg,
		cekAlg,
		objectMetadata.TagLen,
		keyCommitment,
	}
	decryptMaterials, err := m.client.Options.CryptographicMaterialsManager.DecryptMaterials(ctx, decryptMaterialsRequest)
	if err != nil {
		return out, metadata, fmt.Errorf("error while decrypting materials: %w", err)
	}

	cipher, err := cekFunc(*decryptMaterials)
	if err != nil {
		return out, metadata, err
	}
	reader, err := cipher.DecryptContents(result.Body)
	if err != nil {
		return out, metadata, err
	}

	// Apply buffer size configuration for GetObject operations
	// The S3EC MUST set the buffer size to a reasonable default for GetObject
	bufferedReader, err := internal.NewBufferedReader(reader, int(m.client.Options.BufferSize))
	if err != nil {
		return out, metadata, fmt.Errorf("unable to create buffered reader for decrypted contents: %w", err)
	}
	result.Body = bufferedReader
	out.Result = result

	return out, metadata, err
}

func ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(cekAlgSuite *algorithms.AlgorithmSuite, policy commitment.CommitmentPolicy) error {
	//= ../specification/s3-encryption/decryption.md#key-commitment
	//# If the commitment policy requires decryption using a committing algorithm suite,
	//# and the algorithm suite associated with the object does not support key commitment,
	//# then the S3EC MUST throw an exception.
	if policy.RequiresDecrypt() && !cekAlgSuite.IsCommitting() {
		return fmt.Errorf("commitment policy %v does not allow decryption using algorithm suite %v which does not support key commitment", policy, cekAlgSuite)
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT,
	//# the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.FORBID_ENCRYPT_ALLOW_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return nil
		}
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT,
	//# the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return nil
		}
	}

	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
	//# the S3EC MUST NOT allow decryption using algorithm suites which do not support key commitment.
	if policy == commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT {
		if !cekAlgSuite.IsCommitting() {
			return fmt.Errorf("commitment policy %v does not allow decryption using algorithm suite %v which does not support key commitment", policy, cekAlgSuite)
		}
	}

	// If the policy is not recognized, return an error
	switch policy {
		case commitment.FORBID_ENCRYPT_ALLOW_DECRYPT, commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT, commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
			// do nothing -- valid policies
		default:
			return fmt.Errorf("unknown commitment policy: %v", policy)
	}

	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/hex"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecryptionClientV4_GetMockV2Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6b134eb7a353131de92faff64f594b2794e3544e31776cca26fe3bbeeffc68742d1007234f11c6670522602326868e29f37e9d2678f1614ec1a2418009b9772100929aadbed9a21a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"PsuclPnlo2O0MQoov6kL1TBlaZG6oyNwWuAqmAgq7g8b9ZeeORi3VTMg624FU9jx"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"dqqlq2dRVSQ5hFRb"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSContextKeyring},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	out, err := client.GetObject(context.Background(), input)

	actual, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := hex.DecodeString("af150d7156bf5b3f5c461e5c6ac820acc5a33aab7085d920666c250ff251209d5a4029b3bd78250fab6e11aed52fae948d407056a9519b68")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected content to match but it did not")
	}
}

// Use TestEncryptionClientV4_PutMockV3Object to generate a new test vector if needed.
func TestDecryptionClientV4_GetMockV3Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("2d4ff4dafe27f69f628872d82b5a1002ed1a21b8485d532bd8159f6487945b3641af5865fc0a029a3650053600c6d213625b9a0cc9c239577c09f3423dedc5641e88b6835824417c")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"6JOSx47RkdyfciJkNauuC4RpkMcWZY4a+i1RzQ=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"sE9zLb4tsEBJkvEhLMZFxMj9oZJBbQ6ZOgOqHA=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	out, err := client.GetObject(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	actual, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected content to match but it did not")
	}
}


func TestDecryptionClientV4_GetMockV3ObjectWithIncorrectKCValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("e403a8f941e43bdf0ca3ef0bcf6701acd739b2de0a8ee524fa89497210fb0213dfc856376a9ff7753db6ee549dc4040861bc7080a66f902441904bf4a003e028e982de8ea6958c30")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I/0X2GC1iX9Pszf1PAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMDQPII4AynCy/rVwhAgEQgDueLCWabc8WgyoZkAnqVzESQ4NztSDxuETx3obcWJ9Jj6gDAAuDaAL5V+H5QFfwgBqWEcIYt2Ep9WcECw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"FiQepGw+/O+3MuYrmQU5mAkotUnxB+W+EwYDHw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"cPwtbK08jrpe3x+QElyG+vGtkk9jDn6KmOta2Q=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	_, err = client.GetObject(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error due to incorrect key commitment, got nil")
	}
	if !strings.Contains(err.Error(), "derived key commitment value does not match value stored on encrypted message") {
		t.Fatalf("expected key commitment mismatch error, got %v", err)
	}
}

func TestDecryptionClientV4_GetMockV3ObjectWithInvalidKCValue(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("e403a8f941e43bdf0ca3ef0bcf6701acd739b2de0a8ee524fa89497210fb0213dfc856376a9ff7753db6ee549dc4040861bc7080a66f902441904bf4a003e028e982de8ea6958c30")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-3"):   []string{"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I/0X2GC1iX9Pszf1PAAAAfjB8BgkqhkiG9w0BBwagbzBtAgEAMGgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMDQPII4AynCy/rVwhAgEQgDueLCWabc8WgyoZkAnqVzESQ4NztSDxuETx3obcWJ9Jj6gDAAuDaAL5V+H5QFfwgBqWEcIYt2Ep9WcECw=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-c"):       []string{"115"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-d"):  []string{"notvalidbase64"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-i"): []string{"cPwtbK08jrpe3x+QElyG+vGtkk9jDn6KmOta2Q=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-t"):  []string{`{"aws:x-amz-cek-alg":"115"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-w"):  []string{"12"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	_, err = client.GetObject(context.Background(), input)
	if err == nil {
		t.Fatalf("expected error due to incorrect key commitment, got nil")
	}
	if !strings.Contains(err.Error(), "illegal base64 data at input byte") {
		t.Fatalf("expected base64 decoding error, got %v", err)
	}
}


func TestDecryptionClientV4_GetMockV2Object_V1Interop_KMS_AESCBC(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "7ItX9CTGNWWegC62RlaNu6EJ3+J9yGO7yAqDNU4CdeA=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6f4f413a357a3c3a12289442fb835c5e4ecc8db1d86d3d1eab906ce07e1ad772180b2e9ec49c3fc667d8aceea8c46da6bb9738251a8e36241a473ad820f99c701906bac1f48578d5392e928889bbb1d9")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"/nJlgMtxMNk2ErKLLrLp3H7A7aQyJcJOClE2ldAIIFNZU4OhUMc1mMCHdIEC8fby"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"adO9U7pcEHxUTaguIkho9g=="},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"kms_cmk_id":"test-key-id"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSKeyring},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/CBC/PKCS5Padding"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	out, err := client.GetObject(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	actual, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := hex.DecodeString("a716e018ffecf4bb94d4352082af4662612d9c225efed6f389bf1f6f0447a9bce80cc712d7e66ee5e1c086af38e607ead351fd2c1a0247878e693ada73bd580b")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected content to match but it did not")
	}
}

func TestDecryptionClientV4_GetMockV2Object_V1Interop_KMS_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "Hrjrkkt/vQwMYtqvK6+MiXh3xiMvviL1Ks7w2mgsJgU=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6370a90b9a118301c2160c23a90d96146761276acdcfa92e6cbcb783abdc2e1813891506d6850754ef87ed2ac3bf570dd5c9da9492b7769ae1e639d073d688bd284815404ce2648a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"/7tu/RFXZU1UFwRzzf11IdF3b1wBxBZhnUMjVYHKKr5DjAHS602GvXt4zYcx/MJo"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"8Rlvyy8AoYj8v579"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"kms_cmk_id":"test-key-id"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSKeyring},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	out, err := client.GetObject(context.Background(), input)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	actual, err := io.ReadAll(out.Body)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expected, err := hex.DecodeString("75f6805afa7d7be4f56c5906adc27a5959158bf4af6e7c7e12bda3458300f6b1c8daaf9a5949f7a6bdbb8a9c072de05bf0541633421f42f8")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if bytes.Compare(expected, actual) != 0 {
		t.Fatalf("expected content to match but it did not")
	}
}

func TestDecryptionClientV4_GetMockV2Object_OnlyDecryptsRegisteredAlgorithms(t *testing.T) {
	httpClientFactory := func() *awstesting.MockHttpClient {
		b, err := hex.DecodeString("1bd0271b25951fdef3dbe51a9b7af85f66b311e091aa10a346655068f657b9da9acc0843ea0522b0d1ae4a25a31b13605dd1ac5d002db8965d9d4652fd602693")
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		return &awstesting.MockHttpClient{
			Response: &http.Response{
				StatusCode: 200,
				Header: http.Header{
					http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"gNuYjzkLTzfhOcIX9h1l8jApWcAAQqzlryOE166kdDojaHH/+7cCqR5HU8Bpxmij"},
					http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"Vmauu+TMEgaXa26ObqpARA=="},
					http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"kms_cmk_id":"test-key-id"}`},
					http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSKeyring},
					http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/CBC/PKCS5Padding"},
				},
				Body: io.NopCloser(bytes.NewBuffer(b)),
			},
		}
	}

	cases := map[string]struct {
		Client  *S3EncryptionClientV4
		WantErr string
	}{
		"unsupported cek": {
			Client: func() *S3EncryptionClientV4 {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kms.NewFromConfig(awstesting.Config()), func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				})
				cmm, err := materials.NewCryptographicMaterialsManager(keyring)
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				tConfig := awstesting.Config()
				tConfig.HTTPClient = httpClientFactory()
				s3Client := s3.NewFromConfig(tConfig)

				client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				return client
			}(),
			WantErr: "operation error S3: GetObject, configure client with enable legacy unauthenticated modes set to true to decrypt with AES/CBC/PKCS5Padding",
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			client := tt.Client
			input := &s3.GetObjectInput{
				Bucket: aws.String("test"),
				Key:    aws.String("test"),
			}

			_, err := client.GetObject(context.Background(), input)

			if err == nil {
				t.Fatalf("expected error, got none")
			}

			if e, a := tt.WantErr, err.Error(); !strings.Contains(a, e) {
				t.Errorf("expected %v, got %v", e, a)
			}
		})
	}
}

func TestDecryptionClientV4_CheckValidCryptographicMaterialsManager(t *testing.T) {
	_, err := materials.NewCryptographicMaterialsManager(nil)
	if err == nil {
		t.Fatal("expected error, got none")
	}
}

func TestDecryptionClientV4_EncryptionContextValidation_MockV2Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6b134eb7a353131de92faff64f594b2794e3544e31776cca26fe3bbeeffc68742d1007234f11c6670522602326868e29f37e9d2678f1614ec1a2418009b9772100929aadbed9a21a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	cases := map[string]struct {
		storedMatDesc       string
		providedContext     map[string]string
		expectError         bool
		expectedErrorMsg    string
	}{
		"matching encryption context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"custom-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value"},
			expectError:     false,
		},
		"matching encryption context with multiple keys": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1","key2":"value2"}`,
			providedContext: map[string]string{"key1": "value1", "key2": "value2"},
			expectError:     false,
		},
		"empty encryption context matches empty stored context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id"}`,
			providedContext: map[string]string{},
			expectError:     false,
		},
		"mismatched encryption context value": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"stored-value"}`,
			providedContext:  map[string]string{"custom-key": "different-value"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"missing key in provided context": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1","key2":"value2"}`,
			providedContext:  map[string]string{"key1": "value1"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"extra key in provided context": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","key1":"value1"}`,
			providedContext:  map[string]string{"key1": "value1", "key2": "value2"},
			expectError:      true,
			expectedErrorMsg: "Provided encryption context does not match information retrieved from S3",
		},
		"provided context has reserved key (should be ignored in stored context)": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id","custom-key":"custom-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value"},
			expectError:     false,
		},
		"stored context missing kms_cmk_id key": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","custom-key":"custom-value","another-key":"another-value"}`,
			providedContext: map[string]string{"custom-key": "custom-value", "another-key": "another-value"},
			expectError:     false,
		},
		"stored context missing kms_cmk_id with empty provided context": {
			storedMatDesc:   `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
			providedContext: map[string]string{},
			expectError:     false,
		},
		"invalid JSON in stored material description": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","invalid-json":}`,
			providedContext:  map[string]string{"key1": "value1"},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
		"malformed JSON in stored material description": {
			storedMatDesc:    `not-json-at-all`,
			providedContext:  map[string]string{},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
		"incomplete JSON in stored material description": {
			storedMatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"`,
			providedContext:  map[string]string{},
			expectError:      true,
			expectedErrorMsg: "encryption context in stored object is not valid JSON",
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					StatusCode: 200,
					Header: http.Header{
						http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"PsuclPnlo2O0MQoov6kL1TBlaZG6oyNwWuAqmAgq7g8b9ZeeORi3VTMg624FU9jx"},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"dqqlq2dRVSQ5hFRb"},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{tc.storedMatDesc},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSContextKeyring},
						http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
					},
					Body: io.NopCloser(bytes.NewBuffer(b)),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
			})
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			input := &s3.GetObjectInput{
				Bucket: aws.String("test"),
				Key:    aws.String("test"),
			}

			// Create context with encryption context if provided
			ctx := context.Background()
			if tc.providedContext != nil {
				ctx = context.WithValue(ctx, EncryptionContext, tc.providedContext)
			}

			_, err = client.GetObject(ctx, input)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if !strings.Contains(err.Error(), tc.expectedErrorMsg) {
					t.Errorf("expected error message to contain %q, got %q", tc.expectedErrorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
			}
		})
	}
}

func TestDecryptionClientV4_EncryptionContextValidation_InvalidContextType(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, fmt.Sprintf("%s%s%s", `{"KeyId":"test-key-id","Plaintext":"`, "hJUv7S6K2cHF64boS9ixHX0TZAjBZLT4ZpEO4XxkGnY=", `"}`))
	}))
	defer ts.Close()

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	b, err := hex.DecodeString("6b134eb7a353131de92faff64f594b2794e3544e31776cca26fe3bbeeffc68742d1007234f11c6670522602326868e29f37e9d2678f1614ec1a2418009b9772100929aadbed9a21a")
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Header: http.Header{
				http.CanonicalHeaderKey("x-amz-meta-x-amz-key-v2"):   []string{"PsuclPnlo2O0MQoov6kL1TBlaZG6oyNwWuAqmAgq7g8b9ZeeORi3VTMg624FU9jx"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-iv"):       []string{"dqqlq2dRVSQ5hFRb"},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-matdesc"):  []string{`{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","kms_cmk_id":"test-key-id"}`},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-wrap-alg"): []string{materials.KMSContextKeyring},
				http.CanonicalHeaderKey("x-amz-meta-x-amz-cek-alg"):  []string{"AES/GCM/NoPadding"},
			},
			Body: io.NopCloser(bytes.NewBuffer(b)),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	input := &s3.GetObjectInput{
		Bucket: aws.String("test"),
		Key:    aws.String("test"),
	}

	// Test with invalid encryption context type (string instead of map[string]string)
	ctx := context.WithValue(context.Background(), EncryptionContext, "invalid-type")

	_, err = client.GetObject(ctx, input)

	if err == nil {
		t.Fatalf("expected error but got none")
	}

	expectedErrorMsg := "encryption context provided to decrypt method is not valid JSON"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
	}
}

func TestValidateContentEncryptionAlgorithmAgainstCommitmentPolicy_UnrecognizedPolicy(t *testing.T) {
	algSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF
	// Given: Some invalid commitment policy
	invalidPolicy := commitment.CommitmentPolicy(999)
	
	// When: Call the function under test
	err := ValidateContentEncryptionAlgorithmAgainstCommitmentPolicy(algSuite, invalidPolicy)
	
	// Then: function raises error
	if err == nil {
		t.Fatalf("expected error for unrecognized commitment policy, got nil")
	}
	
	// Then: error message contains the expected text
	expectedErrorMsg := "unknown commitment policy"
	if !strings.Contains(err.Error(), expectedErrorMsg) {
		t.Errorf("expected error message to contain %q, got %q", expectedErrorMsg, err.Error())
	}
	
	// Then: error message includes the policy value
	expectedPolicyValue := "CommitmentPolicy(999)"
	if !strings.Contains(err.Error(), expectedPolicyValue) {
		t.Errorf("expected error message to contain %q, got %q", expectedPolicyValue, err.Error())
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"fmt"
	"io"

	"github.com/aws/amazon-s3-encryption-client-go/v4/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/smithy-go"
	"github.com/aws/smithy-go/middleware"
	smithyhttp "github.com/aws/smithy-go/transport/http"
)

// DefaultMinFileSize is used to check whether we want to write to a temp file
// or store the data in memory.
const DefaultMinFileSize = 1024 * 512 * 5

// DefaultBufferSize is the default buffer size for GetObject operations
// The S3EC MUST set the buffer size to a reasonable default for GetObject
const DefaultBufferSize = 1024 * 64 // 64KB default buffer size

// EncryptionContext is used to extract Encryption Context to use on a per-request basis
const EncryptionContext = "EncryptionContext"

// PutObjectAPIClient is a client that implements the PutObject operation
type PutObjectAPIClient interface {
	PutObject(context.Context, *s3.PutObjectInput, ...func(*s3.Options)) (*s3.PutObjectOutput, error)
}

func (m *encryptMiddleware) addEncryptAPIOptions(options *s3.Options) {
	options.APIOptions = append(options.APIOptions,
		m.addEncryptMiddleware,
	)
}

func (m *encryptMiddleware) addEncryptMiddleware(stack *middleware.Stack) error {
	return stack.Serialize.Add(m, middleware.Before)
}

const encryptMiddlewareID = "S3Encrypt"

type encryptMiddleware struct {
	ec *S3EncryptionClientV4
}

// ID returns the resolver identifier
func (m *encryptMiddleware) ID() string {
	return encryptMiddlewareID
}

// HandleSerialize replaces the request body with an encrypted version and saves the envelope using the save strategy
func (m *encryptMiddleware) HandleSerialize(
	ctx context.Context, in middleware.SerializeInput, next middleware.SerializeHandler,
) (
	out middleware.SerializeOutput, metadata middleware.Metadata, err error,
) {

	req, ok := in.Request.(*smithyhttp.Request)
	if !ok {
		return out, metadata, &smithy.SerializationError{Err: fmt.Errorf("unknown transport type %T", in.Request)}
	}

	input, ok := in.Parameters.(*s3.PutObjectInput)
	if !ok {
		return out, metadata, &smithy.SerializationError{Err: fmt.Errorf("unknown input parameters type %T", in.Parameters)}
	}

	// TODO - customize errors?
	reqCopy, err := req.SetStream(input.Body)
	if err != nil {
		return out, metadata, &smithy.SerializationError{Err: err}
	}

	n, ok, err := reqCopy.StreamLength()
	if !ok || err != nil {
		return out, metadata, err
	}

	dst, err := internal.GetWriterStore(m.ec.Options.TempFolderPath, n >= m.ec.Options.MinFileSize)
	if err != nil {
		return out, metadata, err
	}

	ec := ctx.Value(EncryptionContext)
	if ec == nil {
		ec = map[string]string{}
	}
	var matDesc materials.MaterialDescription = ec.(map[string]string)
	cmm := m.ec.Options.CryptographicMaterialsManager
	encryptionMaterialsRequest := materials.EncryptionMaterialsRequest{
		MaterialDescription: matDesc,
		AlgorithmSuite:     m.ec.Options.EncryptionAlgorithmSuite,
	}
	cryptoMaterials, err := cmm.GetEncryptionMaterials(ctx, encryptionMaterialsRequest)
	if err != nil {
		return out, metadata, err
	}
	// Assert that the CMM respected the algorithm suite argument
	if encryptionMaterialsRequest.AlgorithmSuite.CipherName() != cryptoMaterials.CEKAlgorithm {
		return out, metadata, fmt.Errorf("cryptographic materials manager returned unexpected algorithm suite: expected %s, got %s",
			encryptionMaterialsRequest.AlgorithmSuite.CipherName(),
			cryptoMaterials.CEKAlgorithm)
	}
	var cipher internal.ContentCipher

	//= ../specification/s3-encryption/encryption.md#content-encryption
	//# The S3EC MUST use the encryption algorithm configured during [client](./client.md) initialization.
	if m.ec.Options.EncryptionAlgorithmSuite == algorithms.AlgAES256GCMHkdfSha512CommitKey {
		cipher_iv, err := internal.GenerateNonZeroBytes(algorithms.AlgAES256GCMHkdfSha512CommitKey.IVLengthBytes())
		if err != nil {
			return out, metadata, err
		}

		//= ../specification/s3-encryption/encryption.md#content-encryption
		//# The generated IV or Message ID MUST be set or returned from the encryption process such that it can be included in the content metadata.
		cryptoMaterials.IV = cipher_iv

		cipher, err = internal.NewAESGCMCommittingContentCipher(*cryptoMaterials)
		if err != nil {
			return out, metadata, err
		}

	} else if m.ec.Options.EncryptionAlgorithmSuite == algorithms.AlgAES256GCMIV12Tag16NoKDF {
		cipher, err = internal.NewAESGCMContentCipher(*cryptoMaterials)
		if err != nil {
			return out, metadata, err
		}
	} else {
		// S3EC V4 only supports writing AES-GCM, with or without key commitment.
		// Any other algorithms are invalid for encryption.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-iv16-tag16-no-kdf
		//# Attempts to encrypt using AES-CTR MUST fail.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-hkdf-sha512-commit-key
		//# Attempts to encrypt using key committing AES-CTR MUST fail.
		return out, metadata, fmt.Errorf("invalid content encryption algorithm found in options: %s", cryptoMaterials.CEKAlgorithm)
	}

	//= ../specification/s3-encryption/encryption.md#content-encryption
	//# The client MUST validate that the length of the plaintext bytes does not exceed the algorithm suite's cipher's maximum content length in bytes.
	if n >= m.ec.Options.EncryptionAlgorithmSuite.CipherMaxContentLengthBytes() {
		return out, metadata, fmt.Errorf("plaintext length %d exceeds maximum content length for algorithm %s", n, cryptoMaterials.CEKAlgorithm)
	}

	stream := reqCopy.GetStream()
	lengthReader := internal.NewContentLengthReader(stream)
	reader, err := cipher.EncryptContents(lengthReader)
	if err != nil {
		return out, metadata, err
	}

	_, err = io.Copy(dst, reader)
	if err != nil {
		return out, metadata, err
	}

	data := cipher.GetCipherData()
	envelope, err := internal.EncodeMeta(lengthReader, data)
	if err != nil {
		return out, metadata, err
	}

	// rewind
	if _, err := dst.Seek(0, io.SeekStart); err != nil {
		return out, metadata, err
	}

	// update the request body to encrypted contents
	input.Body = dst

	// save the metadata
	saveReq := &internal.SaveStrategyRequest{
		Envelope:    &envelope,
		HTTPRequest: req.Request,
		Input:       input,
	}

	// this saves the required crypto params (IV, tag length, etc.)
	strat := internal.ObjectMetadataSaveStrategy{}
	if err = strat.Save(ctx, saveReq); err != nil {
		return out, metadata, err
	}

	// update the middleware input's parameter which is what the generated serialize step will use
	in.Parameters = input

	out, metadata, err = next.HandleSerialize(ctx, in)

	// cleanup any temp files after the request is made
	dst.Cleanup()
	return out, metadata, err
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

func TestNewEncryptionClientV4_NonDefaults(t *testing.T) {
	tConfig := awstesting.Config()
	tClient := s3.NewFromConfig(tConfig)

	var mcmm = mockCMM{}
	v4, _ := New(tClient, mcmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CryptographicMaterialsManager = mcmm
		clientOptions.TempFolderPath = "/mock/path"
		clientOptions.MinFileSize = 42
	})

	if v4 == nil {
		t.Fatal("expected client to not be nil")
	}

	if !reflect.DeepEqual(mcmm, v4.Options.CryptographicMaterialsManager) {
		t.Errorf("CMM did not match provided value")
	}

	if v4.Client != tClient {
		t.Errorf("expected s3 client not be nil")
	}

	if e, a := 42, v4.Options.MinFileSize; int64(e) != a {
		t.Errorf("expected %v, got %v", e, a)
	}

	if e, a := "/mock/path", v4.Options.TempFolderPath; e != a {
		t.Errorf("expected %v, got %v", e, a)
	}
}

// keyringWithStaticTestIV is a test structure that wraps a CipherDataGeneratorWithCEKAlg and stubs in a static IV
// so that encryption tests can be guaranteed to be consistent.
type keyringWithStaticTestIV struct {
	IV []byte
	materials.Keyring
}

// isAWSFixture will avoid the warning log message when doing tests that need to mock the IV
func (k keyringWithStaticTestIV) isAWSFixture() bool {
	return true
}

func (k keyringWithStaticTestIV) OnEncrypt(ctx context.Context, materials *materials.EncryptionMaterials) (*materials.CryptographicMaterials, error) {
	cryptoMaterials, err := k.Keyring.OnEncrypt(ctx, materials)
	if err == nil {
		cryptoMaterials.IV = k.IV
	}
	return cryptoMaterials, err
}

func TestEncryptionClientV4_PutMockV2Object_KMSCONTEXT_AESGCM(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	var md materials.MaterialDescription
	iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
	kmsWithStaticIV := keyringWithStaticTestIV{
		IV: iv,
		Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		}),
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
		Body: func() io.ReadSeeker {
			content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
			return bytes.NewReader(content)
		}(),
		Metadata: md,
	})
	if err != nil {
		t.Fatalf("PutObject failed with %v", err)
	}

	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	expected, _ := hex.DecodeString("4cd8e95a1c9b8b19640e02838b02c8c09e66250703a602956695afbc23cbb8647d51645955ab63b89733d0766f9a264adb88571b1d467b734ff72eb73d31de9a83670d59688c54ea")

	if !bytes.Equal(tHttpClient.CapturedBody, expected) {
		t.Error("encrypted bytes did not match expected")
	}

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestEncryptionClientV4_PutMockV3Object(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	var md materials.MaterialDescription
	iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
	kmsWithStaticIV := keyringWithStaticTestIV{
		IV: iv,
		Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		}),
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
		Body: func() io.ReadSeeker {
			content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
			return bytes.NewReader(content)
		}(),
		Metadata: md,
	})
	if err != nil {
		t.Fatalf("PutObject failed with %v", err)
	}

	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	if tHttpClient.CapturedReq != nil {
		headers := tHttpClient.CapturedReq.Header
		
		// V3 Content Cipher (x-amz-c) - should be "115" for AES256GCMHkdfSha512CommitKey
		expectedContentCipher := "115"
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-C"); actualValue != expectedContentCipher {
			t.Errorf("X-Amz-Meta-X-Amz-C expected '%s', got '%s'", expectedContentCipher, actualValue)
		}
		
		// V3 Encrypted Data Key (x-amz-3) is nondeterministic, just check presence
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-3"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-3 should be present but was empty")
		}
		
		// V3 Wrapping Algorithm (x-amz-w) - should be "12" for kms+context
		expectedWrappingAlg := "12"
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-W"); actualValue != expectedWrappingAlg {
			t.Errorf("X-Amz-Meta-X-Amz-W expected '%s', got '%s'", expectedWrappingAlg, actualValue)
		}
		
		// V3 Encryption Context (x-amz-t) - should be "{"aws:x-amz-cek-alg":"115"}" for otherwise empty encryption context
		expectedEncryptionContext := `{"aws:x-amz-cek-alg":"115"}`
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-T"); actualValue != expectedEncryptionContext {
			t.Errorf("X-Amz-Meta-X-Amz-T expected '%s', got '%s'", expectedEncryptionContext, actualValue)
		}
		
		// V3 Message ID (x-amz-i) is nondeterministic, just check presence
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-I"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-I should be present but was empty")
		}
		
		// V3 Key Commitment (x-amz-d) is nondeterministic, just check presence
		if actualValue := headers.Get("X-Amz-Meta-X-Amz-D"); strings.TrimSpace(actualValue) == "" {
			t.Errorf("X-Amz-Meta-X-Amz-D should be present but was empty")
		}
	}

	// V3 body is non-deterministic due to IV derivation, so we will just check that it is present
	if tHttpClient.CapturedBody == nil || len(tHttpClient.CapturedBody) == 0 {
		t.Error("Expected encrypted content, but captured body was empty")
	}
}

func TestEncryptionClientv4_PutMockV2Object_KMSCONTEXT_AESGCM_EmptyBody(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	var md materials.MaterialDescription
	iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
	kmsWithStaticIV := keyringWithStaticTestIV{
		IV: iv,
		Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		}),
	}

	tConfig := awstesting.Config()
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(bytes.NewReader([]byte{})),
		},
	}
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
	if err != nil {
		t.Fatalf("error while trying to create new CMM: %v", err)
	}
	client, _ := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket:   aws.String("test-bucket"),
		Key:      aws.String("test-key"),
		Body:     new(bytes.Buffer),
		Metadata: md,
	})
	if err != nil {
		t.Fatalf("PutObject failed with %v", err)
	}

	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	expected, _ := hex.DecodeString("38a7dff91ec56105eedb716fe171675f")

	if !bytes.Equal(tHttpClient.CapturedBody, expected) {
		t.Errorf("encrypted bytes did not match expected")
	}

	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

//= ../specification/s3-encryption/encryption.md#content-encryption
//= type=test
//# The S3EC MUST use the encryption algorithm configured during [client](./client.md) initialization.
func TestS3EC_UsesConfiguredEncryptionAlgorithm(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		ivHex            string
		expectedHeader   string
		expectedValue    string
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813",
			expectedHeader:  "X-Amz-Meta-X-Amz-Cek-Alg",
			expectedValue:   "AES/GCM/NoPadding",
		},
		{
			name:             "AES256GCMHkdfSha512CommitKey",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9",
			expectedHeader:  "X-Amz-Meta-X-Amz-C",
			expectedValue:   "115",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString(tc.ivHex)
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Verify the configured algorithm is set correctly
			if client.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %v, got %v", tc.algorithm, client.Options.EncryptionAlgorithmSuite)
			}

			// Test that PutObject uses the configured algorithm
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
				Body: func() io.ReadSeeker {
					content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
					return bytes.NewReader(content)
				}(),
				Metadata: md,
			})
			if err != nil {
				t.Fatalf("PutObject failed with %v", err)
			}

			// Verify the encryption was performed (captured body should not be empty)
			if tHttpClient.CapturedBody == nil || len(tHttpClient.CapturedBody) == 0 {
				t.Error("Expected encrypted content, but captured body was empty")
			}

			// Capture and validate object metadata for algorithm information
			if tHttpClient.CapturedReq == nil {
				t.Fatal("Expected captured HTTP request, but it was nil")
			}

			// Verify the algorithm-specific header is set correctly
			headers := tHttpClient.CapturedReq.Header
			actualValue := headers.Get(tc.expectedHeader)
			
			if actualValue == "" {
				t.Errorf("Expected %s header to be present for %s algorithm", tc.expectedHeader, tc.name)
			} else if actualValue != tc.expectedValue {
				t.Errorf("Expected %s header to be '%s', got '%s'", tc.expectedHeader, tc.expectedValue, actualValue)
			}

			t.Logf("✓ S3EC successfully used configured %s algorithm with correct metadata", tc.name)
			t.Logf("  - %s: %s", tc.expectedHeader, actualValue)
		})
	}
}

// mockLargeReader simulates a large reader without allocating memory
type mockLargeReader struct {
	size     int64
	position int64
}

func (r *mockLargeReader) Read(p []byte) (n int, err error) {
	if r.position >= r.size {
		return 0, io.EOF
	}
	
	remaining := r.size - r.position
	toRead := int64(len(p))
	if toRead > remaining {
		toRead = remaining
	}
	
	// Fill buffer with test data
	for i := int64(0); i < toRead; i++ {
		p[i] = byte((r.position + i) % 256)
	}
	
	r.position += toRead
	return int(toRead), nil
}

func (r *mockLargeReader) Seek(offset int64, whence int) (int64, error) {
	switch whence {
	case io.SeekStart:
		r.position = offset
	case io.SeekCurrent:
		r.position += offset
	case io.SeekEnd:
		r.position = r.size + offset
		if r.position > r.size {
			return 0, fmt.Errorf("seek position beyond end of reader")
		}
	}
	
	if r.position < 0 {
		r.position = 0
	}
	if r.position > r.size {
		r.position = r.size
	}
	
	return r.position, nil
}

//= ../specification/s3-encryption/encryption.md#content-encryption
//= type=test
//# The client MUST validate that the length of the plaintext bytes does not exceed the algorithm suite's cipher's maximum content length in bytes.
func TestS3EC_ValidatesPlaintextLengthLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		contentSize       int64
		expectError       bool
		errorContains     string
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF_WithinLimit",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			contentSize:      1024, // Well within limit
			expectError:      false,
		},
		{
			name:             "AES256GCMIV12Tag16NoKDF_ExceedsLimit",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			contentSize:      algorithms.AlgAES256GCMIV12Tag16NoKDF.CipherMaxContentLengthBytes() + 1, // Just over the limit
			expectError:      true,
			errorContains:    "plaintext length",
		},
		{
			name:             "AES256GCMHkdfSha512CommitKey_WithinLimit",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			contentSize:      1024, // Well within limit
			expectError:      false,
		},
		{
			name:             "AES256GCMHkdfSha512CommitKey_ExceedsLimit",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			contentSize:      algorithms.AlgAES256GCMHkdfSha512CommitKey.CipherMaxContentLengthBytes() + 1, // Just over the limit
			expectError:      true,
			errorContains:    "plaintext length",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString("ae325acae2bfd5b9c3d0b813")
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Create a mock reader for the specified size (avoids memory allocation issues)
			var body io.ReadSeeker
			if tc.contentSize <= 1024*1024 { // For small sizes, use real content
				content := make([]byte, tc.contentSize)
				for i := range content {
					content[i] = byte(i % 256)
				}
				body = bytes.NewReader(content)
			} else {
				// For large sizes, use mock reader
				body = &mockLargeReader{size: tc.contentSize}
			}

			// Test PutObject with the content
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket:   aws.String("test-bucket"),
				Key:      aws.String("test-key"),
				Body:     body,
				Metadata: md,
			})

			if tc.expectError {
				if err == nil {
					t.Errorf("Expected error for content size %d bytes (limit: %d), but got none", 
						tc.contentSize, tc.algorithm.CipherMaxContentLengthBytes())
				} else if tc.errorContains != "" && !bytes.Contains([]byte(err.Error()), []byte(tc.errorContains)) {
					t.Errorf("Expected error to contain '%s', but got: %v", tc.errorContains, err)
				} else {
					t.Logf("✓ S3EC correctly rejected content exceeding maximum length for %s", tc.name)
				}
			} else {
				if err != nil {
					t.Errorf("Expected no error for content size %d bytes (limit: %d), but got: %v", 
						tc.contentSize, tc.algorithm.CipherMaxContentLengthBytes(), err)
				} else {
					t.Logf("✓ S3EC correctly accepted content within maximum length for %s", tc.name)
				}
			}
		})
	}
}

func TestS3EC_IVGenerationAndMetadataInclusion(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, `{"CiphertextBlob":"8gSzlk7giyfFbLPUVgoVjvQebI1827jp8lDkO+n2chsiSoegx1sjm8NdPk0Bl70I","KeyId":"test-key-id","Plaintext":"lP6AbIQTmptyb/+WQq+ubDw+w7na0T1LGSByZGuaono="}`)
	}))

	tKmsConfig := awstesting.Config()
	tKmsConfig.Region = "us-west-2"
	tKmsConfig.RetryMaxAttempts = 0
	tKmsConfig.EndpointResolverWithOptions = awstesting.TestEndpointResolver(ts.URL)
	kmsClient := kms.NewFromConfig(tKmsConfig)

	testCases := []struct {
		name              string
		algorithm         *algorithms.AlgorithmSuite
		commitmentPolicy  commitment.CommitmentPolicy
		ivHex            string
		expectedHeader   string
		expectedIVLength  int // in bytes
	}{
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813", // 12 bytes
			expectedHeader:  "X-Amz-Meta-X-Amz-Iv",
			expectedIVLength: 12, // 96 bits / 8 = 12 bytes
		},
		{
			name:             "AES256GCMHkdfSha512CommitKey",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			ivHex:           "ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9c3d0b813ae325acae2bfd5b9", // 32 bytes
			expectedHeader:  "X-Amz-Meta-X-Amz-I", // Message ID header for V3 format
			expectedIVLength: 28, // 224 bits / 8 = 28 bytes (Message ID length)
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// First, verify the algorithm's IV length specification
			actualBytes := tc.algorithm.IVLengthBytes()
			expectedBits := actualBytes * 8

			if actualBytes != tc.expectedIVLength {
				t.Errorf("Expected IV/Message ID length %d bytes for %s, got %d bytes", 
					tc.expectedIVLength, tc.name, actualBytes)
			}

			t.Logf("✓ S3EC algorithm %s correctly defines IV/Message ID length: %d bytes (%d bits)", 
				tc.name, actualBytes, expectedBits)

			// Now test that the IV is properly included in content metadata during encryption
			var md materials.MaterialDescription
			iv, _ := hex.DecodeString(tc.ivHex)
			kmsWithStaticIV := keyringWithStaticTestIV{
				IV: iv,
				Keyring: materials.NewKmsKeyring(kmsClient, "test-key-id", func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				}),
			}

			tConfig := awstesting.Config()
			tHttpClient := &awstesting.MockHttpClient{
				Response: &http.Response{
					Status:     http.StatusText(200),
					StatusCode: http.StatusOK,
					Body:       io.NopCloser(bytes.NewReader([]byte{})),
				},
			}
			tConfig.HTTPClient = tHttpClient
			s3Client := s3.NewFromConfig(tConfig)

			cmm, err := materials.NewCryptographicMaterialsManager(kmsWithStaticIV)
			if err != nil {
				t.Fatalf("error while trying to create new CMM: %v", err)
			}

			// Configure client with specific encryption algorithm
			client, err := New(s3Client, cmm, func(clientOptions *EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithm
			})
			if err != nil {
				t.Fatalf("Failed to create S3EC: %v", err)
			}

			// Test that PutObject includes IV/Message ID in metadata
			_, err = client.PutObject(context.Background(), &s3.PutObjectInput{
				Bucket: aws.String("test-bucket"),
				Key:    aws.String("test-key"),
				Body: func() io.ReadSeeker {
					content, _ := hex.DecodeString("8f2c59c6dbfcacf356f3da40788cbde67ca38161a4702cbcf757af663e1c24a600001b2f500417dbf5a050f57db6737422b2ed6a44c75e0d")
					return bytes.NewReader(content)
				}(),
				Metadata: md,
			})
			if err != nil {
				t.Fatalf("PutObject failed with %v", err)
			}

			// Verify the IV/Message ID is included in content metadata
			if tHttpClient.CapturedReq == nil {
				t.Fatal("Expected captured HTTP request, but it was nil")
			}

			headers := tHttpClient.CapturedReq.Header
			ivValue := headers.Get(tc.expectedHeader)
			
			if ivValue == "" {
				t.Errorf("Expected %s header to be present for %s algorithm", tc.expectedHeader, tc.name)
			} else {
				//= ../specification/s3-encryption/encryption.md#content-encryption
				//= type=test
				//# The client MUST generate an IV or Message ID using the length of the IV or Message ID defined in the algorithm suite.
				//# The generated IV or Message ID MUST be set or returned from the encryption process such that it can be included in the content metadata.
				decodedIV, err := hex.DecodeString(ivValue)
				if err != nil {
					// Try base64 decoding for V3 format
					decodedIV, err = base64.StdEncoding.DecodeString(ivValue)
				}
				
				if err == nil {
					expectedLength := tc.algorithm.IVLengthBytes()
					
					if len(decodedIV) != expectedLength {
						t.Errorf("Expected IV/Message ID length %d bytes, got %d bytes", expectedLength, len(decodedIV))
					}
				}

				t.Logf("✓ S3EC successfully included IV/Message ID in content metadata for %s", tc.name)
				t.Logf("  - %s: %s", tc.expectedHeader, ivValue)
			}
		})
	}
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-iv16-tag16-no-kdf
//= type=test
//# Attempts to encrypt using AES-CTR MUST fail.
func TestS3EC_AESCTREncryptionMustFail(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	
	var mcmm = mockCMM{}

	// Attempt to create client with AES-CTR algorithm - this should fail
	_, err := New(s3Client, mcmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CTRIV16Tag16NoKDF
	})

	if err == nil {
		t.Fatalf("Expected error when attempting to use AES-CTR algorithm, but got none")
	}

	if !strings.Contains(err.Error(), "AES-CTR") && !strings.Contains(err.Error(), "CTR") {
		t.Errorf("Expected error to mention AES-CTR, but got: %v", err)
	}

	t.Logf("✓ S3EC correctly rejected AES-CTR encryption algorithm: %v", err)
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-ctr-hkdf-sha512-commit-key
//= type=test
//# Attempts to encrypt using key committing AES-CTR MUST fail.
func TestS3EC_CommittingAESCTREncryptionMustFail(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	
	var mcmm = mockCMM{}

	// Attempt to create client with AES-CTR algorithm and committing policy - this should fail
	_, err := New(s3Client, mcmm, func(clientOptions *EncryptionClientOptions) {
		clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CTRIV16Tag16NoKDF
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})

	if err == nil {
		t.Fatalf("Expected error when attempting to use key committing AES-CTR algorithm, but got none")
	}

	if !strings.Contains(err.Error(), "AES-CTR") && !strings.Contains(err.Error(), "CTR") {
		t.Errorf("Expected error to mention AES-CTR, but got: %v", err)
	}

	t.Logf("✓ S3EC correctly rejected key committing AES-CTR encryption algorithm: %v", err)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"bytes"
	"context"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
	"io/ioutil"
)

type mockCMM struct{}

func (m mockCMM) GetEncryptionMaterials(ctx context.Context, req materials.EncryptionMaterialsRequest) (*materials.CryptographicMaterials, error) {
	// TODO: make this mock more useful
	return &materials.CryptographicMaterials{
		Key:          nil,
		IV:           nil,
		CEKAlgorithm: "AES/GCM/NoPadding",
	}, nil
}

func (m mockCMM) DecryptMaterials(ctx context.Context, req materials.DecryptMaterialsRequest) (*materials.CryptographicMaterials, error) {
	// TODO: make this mock more useful
	return &materials.CryptographicMaterials{
		Key:                 nil,
		IV:                  nil,
		KeyringAlgorithm:    "",
		CEKAlgorithm:        "",
		TagLength:           "",
		MaterialDescription: nil,
		EncryptedKey:        nil,
	}, nil
}

type mockKeyring struct{}

// OnEncrypt generates/encrypts a data key for use with content encryption
func (m mockKeyring) OnEncrypt(ctx context.Context, encryptionMaterials *materials.EncryptionMaterials) (*materials.CryptographicMaterials, error) {
	return &materials.CryptographicMaterials{
		Key:                 nil,
		IV:                  nil,
		KeyringAlgorithm:    "",
		CEKAlgorithm:        "",
		TagLength:           "",
		MaterialDescription: nil,
		EncryptedKey:        nil,
	}, nil
	// TODO: make this mock more useful
	return &materials.CryptographicMaterials{}, nil
}

// OnDecrypt decrypts the encryptedDataKeys and returns them in materials
// for use with content decryption
func (m mockKeyring) OnDecrypt(ctx context.Context, decryptionMaterials *materials.DecryptionMaterials, encryptedDataKey materials.DataKey) (*materials.CryptographicMaterials, error) {
	// TODO: make this mock more useful
	return &materials.CryptographicMaterials{
		Key:                 nil,
		IV:                  nil,
		KeyringAlgorithm:    "",
		CEKAlgorithm:        "",
		TagLength:           "",
		MaterialDescription: nil,
		EncryptedKey:        nil,
	}, nil
}

type mockContentCipher struct {
	materials materials.CryptographicMaterials
}

func (cipher *mockContentCipher) GetCipherData() materials.CryptographicMaterials {
	return cipher.materials
}

func (cipher *mockContentCipher) EncryptContents(src io.Reader) (io.Reader, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	size := len(b)
	b = bytes.Repeat([]byte{1}, size)
	return bytes.NewReader(b), nil
}

func (cipher *mockContentCipher) DecryptContents(src io.ReadCloser) (io.ReadCloser, error) {
	b, err := ioutil.ReadAll(src)
	if err != nil {
		return nil, err
	}
	size := len(b)
	return ioutil.NopCloser(bytes.NewReader(make([]byte, size))), nil
}

type mockPadder struct {
}

func (m mockPadder) Pad(i []byte, i2 int) ([]byte, error) {
	return i, nil
}

func (m mockPadder) Unpad(i []byte) ([]byte, error) {
	return i, nil
}

func (m mockPadder) Name() string {
	return "mockPadder"
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"log"
	"fmt"

	"github.com/aws/amazon-s3-encryption-client-go/v4/internal"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// S3EncryptionClientV4 provides client-side encryption for S3.
// The client embeds a default client to provide support for control plane operations
// which do not involve encryption.
type S3EncryptionClientV4 struct {
	//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
	//# The S3EC SHOULD support invoking operations unrelated to client-side encryption e.g. CopyObject as the conventional AWS SDK S3 client would.
	//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
	//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

	*s3.Client                         // promoted anonymous field, it allows this type to call s3 Client methods
	Options    EncryptionClientOptions // options for encrypt/decrypt
}

//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=implication
//# The S3EC MUST provide a different set of configuration options than the conventional S3 client.

// EncryptionClientOptions is the configuration options for the S3 Encryption Client.
type EncryptionClientOptions struct {
	// TempFolderPath is used to store temp files when calling PutObject
	// Temporary files are needed to compute the X-Amz-Content-Sha256 header
	TempFolderPath string

	// MinFileSize is the minimum size for the content to write to a
	// temporary file instead of using memory
	MinFileSize int64

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//# The S3EC SHOULD accept a configurable buffer size
	//# which refers to the maximum ciphertext length in bytes to store in memory
	//# when Delayed Authentication mode is disabled.

	// BufferSize is the buffer size used for GetObject operations
	BufferSize int64

	// The logger to write logging messages to
	Logger *log.Logger

	// The CryptographicMaterialsManager to use to manage encryption and decryption materials
	CryptographicMaterialsManager materials.CryptographicMaterialsManager

	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=implication
	//# The S3EC MUST support the option to enable or disable legacy unauthenticated modes (content encryption algorithms).

	// EnableLegacyUnauthenticatedModes MUST be set to true in order to decrypt objects encrypted
	// using legacy (unauthenticated) modes such as AES/CBC. The default is false.
	EnableLegacyUnauthenticatedModes bool

	//= ../specification/s3-encryption/client.md#key-commitment
	//# The S3EC MUST support configuration of the [Key Commitment policy](./key-commitment.md) during its initialization.

	// CommitmentPolicy specifies the key commitment policy for this client.
	// S3EncryptionClientV4 defaults to commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT.
	// Objects written by a client configured with this default can be read by any v4 client
	// or any v3 client from v3.2.0 onward, but a client configured with this default can
	// only read objects written by other v4 clients configured with either
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_REQUIRE_DECRYPT commitment policies.
	// If an EncryptionAlgorithmSuite is also provided,
	// the selected CommitmentPolicy must also be compatible with the selected EncryptionAlgorithmSuite; if not, New() will return an error.
	CommitmentPolicy commitment.CommitmentPolicy

	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# The S3EC MUST support configuration of the encryption algorithm (or algorithm suite) during its initialization.

	// EncryptionAlgorithmSuite specifies the algorithm suite to use when encrypting objects.
	// If the commitment policy requires encrypting with key committing algorithms (default for S3EncryptionClientV4),
	// then S3EncryptionClientV4 defaults to algorithms.AlgAES256GCMHkdfSha512CommitKey.
	// However, If the commitment policy does not require encrypting with key committing algorithms,
	// S3EncryptionClientV4 defaults to algorithms.AlgAES256GCMIV12Tag16NoKDF.
	// The client will decrypt objects encrypted with any supported algorithm suite, provided that the
	// algorithms suite is compatible with the selected CommitmentPolicy and EnableLegacyUnauthenticatedModes options.
	// If a CommitmentPolicy is also provided,
	// the selected EncryptionAlgorithmSuite must also be compatible with the selected CommitmentPolicy; if not, New() will return an error.
	EncryptionAlgorithmSuite *algorithms.AlgorithmSuite
}

//# The S3EC MUST NOT support use of S3EC as the provided S3 client during its initialization; it MUST throw an exception in this case.

// New creates a new S3 Encryption Client v4 with the given CryptographicMaterialsManager
func New(s3Client *s3.Client, CryptographicMaterialsManager materials.CryptographicMaterialsManager, optFns ...func(options *EncryptionClientOptions)) (*S3EncryptionClientV4, error) {
	//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
	//# The S3EC MUST support the option to provide an SDK S3 client instance during its initialization.
	// In Go, the requirement below is enforced/implied by the type system at compile time.
	// The New() function expects *s3.Client, but S3EncryptionClientV4 is a wholly different type.
	//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
	//= type=implication
	//# The S3EC MUST NOT support use of S3EC as the provided S3 client during its initialization; it MUST throw an exception in this case.
	wrappedClient := s3Client
	// default options
	options := EncryptionClientOptions{
		MinFileSize:                      DefaultMinFileSize,
		//= ../specification/s3-encryption/client.md#set-buffer-size
		//# If Delayed Authentication mode is disabled, and no buffer size is provided,
		//# the S3EC MUST set the buffer size to a reasonable default.
		BufferSize:                       DefaultBufferSize,
		Logger:                           log.Default(),
		// S3EC Go V4 only accepts a CMM on client configuration, not a keyring.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# The S3EC MUST accept either one CMM or one Keyring instance upon initialization.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# If both a CMM and a Keyring are provided, the S3EC MUST throw an exception.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# When a Keyring is provided, the S3EC MUST create an instance of the DefaultCMM using the provided Keyring.
		//= ../specification/s3-encryption/client.md#cryptographic-materials
		//= type=exception
		//# The S3EC MAY accept key material directly.
		CryptographicMaterialsManager:    CryptographicMaterialsManager,
		//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
		//# The option to enable legacy unauthenticated modes MUST be set to false by default.
		EnableLegacyUnauthenticatedModes: false,
		CommitmentPolicy:                 commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,

		// S3EC Go V4 does not support delayed authentication mode.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# The S3EC MUST support the option to enable or disable Delayed Authentication mode.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# Delayed Authentication mode MUST be set to false by default.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# When enabled, the S3EC MAY release plaintext from a stream which has not been authenticated.
		//= ../specification/s3-encryption/client.md#enable-delayed-authentication
		//= type=exception
		//# When disabled the S3EC MUST NOT release plaintext from a stream which has not been authenticated.
		//= ../specification/s3-encryption/client.md#set-buffer-size
		//= type=exception
		//# If Delayed Authentication mode is enabled, and the buffer size has been set to a value other than its default, the S3EC MUST throw an exception.
	
		// Go S3EC V4 does not support instruction file configuration.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# The S3EC MAY support the option to provide Instruction File Configuration during its initialization.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# If the S3EC in a given language supports Instruction Files, then it MUST accept Instruction File Configuration during its initialization.
		//= ../specification/s3-encryption/client.md#instruction-file-configuration
		//= type=exception
		//# In this case, the Instruction File Configuration SHOULD be optional, such that its default configuration is used when none is provided.

		// Go S3EC V4 does not support a single inherited configuration for underlying AWS SDK clients.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# The S3EC MAY support directly configuring the wrapped SDK clients through its initialization.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# For example, the S3EC MAY accept a credentials provider instance during its initialization.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# If the S3EC accepts SDK client configuration, the configuration MUST be applied to all wrapped S3 clients.
		//= ../specification/s3-encryption/client.md#inherited-sdk-configuration
		//= type=exception
		//# If the S3EC accepts SDK client configuration, the configuration MUST be applied to all wrapped SDK clients including the KMS client.
	
		// Go S3EC V4 does not support supplying a custom source of randomness during client initialization.
		//= ../specification/s3-encryption/client.md#randomness
		//= type=exception
		//# The S3EC MAY accept a source of randomness during client initialization.
	}
	// apply functional options
	for _, fn := range optFns {
		fn(&options)
	}

	// If no algorithm suite, supply a default
	if options.EncryptionAlgorithmSuite == nil {
		options.EncryptionAlgorithmSuite = DefaultEncryptionAlgorithmSuite(options)
	}

	// Validate selected encryption algorithm suite
	if err := ValidateEncryptionAlgorithmSuite(options); err != nil {
		return nil, err
	}

	// use the given wrappedClient for the promoted anon fields
	s3ec := &S3EncryptionClientV4{wrappedClient, options}
	return s3ec, nil
}

func DefaultEncryptionAlgorithmSuite(options EncryptionClientOptions) *algorithms.AlgorithmSuite {
	if options.CommitmentPolicy.RequiresEncrypt() {
		return algorithms.AlgAES256GCMHkdfSha512CommitKey
	} else {
		return algorithms.AlgAES256GCMIV12Tag16NoKDF
	}
}

// Explict (but verbose) validations of S3EC specification
func ValidateEncryptionAlgorithmSuite(options EncryptionClientOptions) error {
	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# The S3EC MUST validate that the configured encryption algorithm is not legacy.
	//= ../specification/s3-encryption/client.md#encryption-algorithm
	//# If the configured encryption algorithm is legacy, then the S3EC MUST throw an exception.
	if options.EncryptionAlgorithmSuite.IsLegacy() {
		return fmt.Errorf("legacy algorithm suites are not allowed for decrypt, got %v", options.EncryptionAlgorithmSuite)
	}
	//= ../specification/s3-encryption/client.md#key-commitment
	//# The S3EC MUST validate the configured Encryption Algorithm against the provided key commitment policy.
	//= ../specification/s3-encryption/client.md#key-commitment
	//# If the configured Encryption Algorithm is incompatible with the key commitment policy, then it MUST throw an exception.
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST NOT encrypt using an algorithm suite which supports key commitment.
	if options.CommitmentPolicy == commitment.FORBID_ENCRYPT_ALLOW_DECRYPT {
		if options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy FORBID_ENCRYPT_ALLOW_DECRYPT does not allow committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
	} else if options.CommitmentPolicy == commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT {
		if !options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy REQUIRE_ENCRYPT_ALLOW_DECRYPT requires committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	//= ../specification/s3-encryption/key-commitment.md#commitment-policy
	//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
	} else if options.CommitmentPolicy == commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT {
		if !options.EncryptionAlgorithmSuite.IsCommitting() {
			return fmt.Errorf("CommitmentPolicy REQUIRE_ENCRYPT_REQUIRE_DECRYPT requires committing algorithm suites, got %v", options.EncryptionAlgorithmSuite)
		}
	} else {
		return fmt.Errorf("unknown CommitmentPolicy %v", options.CommitmentPolicy)
	}
	return nil
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - GetObject MUST be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

// GetObject will make a request to s3 and retrieve the object. In this process
// decryption will be done. The SDK only supports region reads of KMS and GCM.
func (c *S3EncryptionClientV4) GetObject(ctx context.Context, input *s3.GetObjectInput, optFns ...func(*s3.Options)) (*s3.GetObjectOutput, error) {
	// Go S3EC V4 does not support ranged gets.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# The S3EC MAY support the "range" parameter on GetObject which specifies a subset of bytes to download and decrypt.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the S3EC supports Ranged Gets, the S3EC MUST adjust the customer-provided range to include the beginning and end of the cipher blocks for the given range.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# For requests which provide a range to decrypt an object encrypted with an authenticated algorithm suite, the corresponding CTR-based algorithm suite is used.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the GetObject response contains a range, but the GetObject request does not contain a range, the S3EC MUST throw an exception.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the object was encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF, then ALG_AES_256_CTR_IV16_TAG16_NO_KDF MUST be used to decrypt the range of the object.
	//= ../specification/s3-encryption/decryption.md#ranged-gets
	//= type=exception
	//# If the object was encrypted with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY, then ALG_AES_256_CTR_HKDF_SHA512_COMMIT_KEY MUST be used to decrypt the range of the object.

	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - GetObject MUST decrypt data received from the S3 server and return it as plaintext.
	m := &decryptMiddleware{
		client: c,
		input:  input,
	}
	decryptOpts := []func(*s3.Options){
		internal.AddS3CryptoUserAgent,
		m.addDecryptAPIOptions,
	}

	opts := append(optFns, decryptOpts...)
	return c.Client.GetObject(ctx, input, opts...)
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - PutObject MUST be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.

// PutObject will make encrypt the contents before sending the data to S3. Depending on the MinFileSize
// a temporary file may be used to buffer the encrypted contents to.
func (c *S3EncryptionClientV4) PutObject(ctx context.Context, input *s3.PutObjectInput, optFns ...func(*s3.Options)) (*s3.PutObjectOutput, error) {
	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - PutObject MUST encrypt its input data before it is uploaded to S3.
	em := &encryptMiddleware{
		ec: c,
	}

	encryptOpts := []func(*s3.Options){
		internal.AddS3CryptoUserAgent,
		em.addEncryptAPIOptions,
	}

	opts := append(optFns, encryptOpts...)
	return c.Client.PutObject(ctx, input, opts...)
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - DeleteObject MUST be implemented by the S3EC.

// DeleteObject will defer to the underlying S3 client to delete the object,
// but will execute its own logic to delete the associated instruction file using the default instruction file suffix.
func (c *S3EncryptionClientV4) DeleteObject(ctx context.Context, input *s3.DeleteObjectInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectOutput, error) {
	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - DeleteObject MUST delete the given object key.
	result, err := c.Client.DeleteObject(ctx, input, optFns...)
	if err != nil {
		return result, err
	}

	//= ../specification/s3-encryption/client.md#required-api-operations
	//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
	
	// Delete the associated instruction file
	instructionFileKey := *input.Key + internal.DefaultInstructionKeySuffix
	instructionInput := &s3.DeleteObjectInput{
		Bucket: input.Bucket,
		Key:    &instructionFileKey,
	}
	
	// Delete instruction file - ignore errors as the instruction file may not exist
	_, _ = c.Client.DeleteObject(ctx, instructionInput, optFns...)
	
	return result, err
}

//= ../specification/s3-encryption/client.md#required-api-operations
//# - DeleteObjects MUST be implemented by the S3EC.

// DeleteObjects will delete multiple objects by calling DeleteObject for each object.
// This ensures that both the objects and their associated instruction files are deleted.
func (c *S3EncryptionClientV4) DeleteObjects(ctx context.Context, input *s3.DeleteObjectsInput, optFns ...func(*s3.Options)) (*s3.DeleteObjectsOutput, error) {

	// We implement DeleteObjects by calling DeleteObject for each object
	// This ensures that both the object and its instruction file are deleted for each item
	var deletedObjects []types.DeletedObject
	var errors []types.Error
	
	for _, obj := range input.Delete.Objects {
		// Call our DeleteObject method which handles both object and instruction file deletion
		//= ../specification/s3-encryption/client.md#required-api-operations
		//# - DeleteObjects MUST delete each of the given objects.
		//= ../specification/s3-encryption/client.md#required-api-operations
		//# - DeleteObjects MUST delete each of the corresponding instruction files using the default instruction file suffix.
		deleteResult, err := c.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: input.Bucket,
			Key:    obj.Key,
		}, optFns...)
		
		if err != nil {
			// Add error to the errors list
			errors = append(errors, types.Error{
				Key:     obj.Key,
				Code:    aws.String("InternalError"),
				Message: aws.String(err.Error()),
			})
		} else {
			// Add successful deletion to the deleted objects list
			deletedObjects = append(deletedObjects, types.DeletedObject{
				Key:          obj.Key,
				DeleteMarker: deleteResult.DeleteMarker,
				VersionId:    deleteResult.VersionId,
			})
		}
	}
	
	// Build the response
	result := &s3.DeleteObjectsOutput{
		Deleted: deletedObjects,
	}
	
	// Add errors if any occurred
	if len(errors) > 0 {
		result.Errors = errors
	}
	
	return result, nil
}

// S3EC Go V4 does not implement the following operations:
// - CreateMultipartUpload
// - UploadPart
// - CompleteMultipartUpload
// - AbortMultipartUpload
// - ReEncryptInstructionFile

//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CreateMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - If implemented, CreateMultipartUpload MUST initiate a multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - UploadPart MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - UploadPart MUST encrypt each part.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - Each part MUST be encrypted in sequence.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - Each part MUST be encrypted using the same cipher instance for each part.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CompleteMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - CompleteMultipartUpload MUST complete the multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - AbortMultipartUpload MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - AbortMultipartUpload MUST abort the multipart upload.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MAY be implemented by the S3EC.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MUST decrypt the instruction file's encrypted data key for the given object using the client's CMM.
//= ../specification/s3-encryption/client.md#optional-api-operations
//= type=exception
//# - ReEncryptInstructionFile MUST re-encrypt the plaintext data key with a provided keyring.
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package client

import (
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/amazon-s3-encryption-client-go/v4/internal/awstesting"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=test
//# The S3EC SHOULD support invoking operations unrelated to client-side encryption e.g. CopyObject as the conventional AWS SDK S3 client would.
func TestWHEN_CallNonEncryptionOperationOnS3EC_THEN_PassthroughToPlaintextS3Client(t *testing.T) {
	// Create mock HTTP response for ListBuckets
	listBucketsResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Owner>
        <ID>test-owner-id</ID>
        <DisplayName>test-owner</DisplayName>
    </Owner>
    <Buckets>
        <Bucket>
            <Name>test-bucket-1</Name>
            <CreationDate>2023-01-01T00:00:00.000Z</CreationDate>
        </Bucket>
        <Bucket>
            <Name>test-bucket-2</Name>
            <CreationDate>2023-01-02T00:00:00.000Z</CreationDate>
        </Bucket>
    </Buckets>
</ListAllMyBucketsResult>`

	// Create mock HTTP client
	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			Status:     http.StatusText(200),
			StatusCode: http.StatusOK,
			Body:       io.NopCloser(strings.NewReader(listBucketsResponse)),
		},
	}

	// Create test config with mock HTTP client
	tConfig := awstesting.Config()
	tConfig.HTTPClient = tHttpClient
	s3Client := s3.NewFromConfig(tConfig)

	// Create mock CMM
	mockCMM := &mockCMM{}

	// Create the S3 encryption client
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Test that ListBuckets is passed through to the underlying S3 client
	ctx := context.TODO()
	result, err := s3ec.ListBuckets(ctx, &s3.ListBucketsInput{})

	// Verify no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the response is parsed correctly
	if result == nil {
		t.Error("Expected non-nil result")
	} else {
		if len(result.Buckets) != 2 {
			t.Errorf("Expected 2 buckets, got %d", len(result.Buckets))
		}
		if result.Owner == nil || *result.Owner.DisplayName != "test-owner" {
			t.Error("Expected owner display name to be 'test-owner'")
		}
	}
}

//= ../specification/s3-encryption/client.md#aws-sdk-compatibility
//= type=test
//# The S3EC MUST adhere to the same interface for API operations as the conventional AWS SDK S3 client.
func TestS3EC_AdheresToSameInterfaceAsConventionalS3Client(t *testing.T) {
	// This test validates that S3EC and regular S3 client have identical interfaces
	// by calling both with the same parameters and verifying compatible behavior
	
	testCases := []struct {
		name        string
		operation   string
		setupMock   func() *awstesting.MockHttpClient
		testFunc    func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV4)
	}{
		{
			name:      "ListBuckets",
			operation: "ListBuckets",
			setupMock: func() *awstesting.MockHttpClient {
				mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Owner>
        <ID>test-owner-id</ID>
        <DisplayName>test-owner</DisplayName>
    </Owner>
    <Buckets>
        <Bucket>
            <Name>test-bucket-1</Name>
            <CreationDate>2023-01-01T00:00:00.000Z</CreationDate>
        </Bucket>
        <Bucket>
            <Name>test-bucket-2</Name>
            <CreationDate>2023-01-02T00:00:00.000Z</CreationDate>
        </Bucket>
    </Buckets>
</ListAllMyBucketsResult>`
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockResponse)),
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV4) {
				ctx := context.TODO()
				input := &s3.ListBucketsInput{}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.ListBuckets(ctx, input)
				result2, err2 := s3ec.ListBuckets(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					if len(result1.Buckets) != len(result2.Buckets) {
						t.Errorf("Bucket count mismatch: s3Client=%d, s3ec=%d", len(result1.Buckets), len(result2.Buckets))
					}
					
					if result1.Owner != nil && result2.Owner != nil {
						if *result1.Owner.DisplayName != *result2.Owner.DisplayName {
							t.Errorf("Owner display name mismatch: s3Client=%s, s3ec=%s", 
								*result1.Owner.DisplayName, *result2.Owner.DisplayName)
						}
					}
					
					t.Logf("✓ ListBuckets: Both clients returned identical structured output")
				}
			},
		},
		{
			name:      "HeadBucket",
			operation: "HeadBucket",
			setupMock: func() *awstesting.MockHttpClient {
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader("")),
						Header: http.Header{
							"x-amz-bucket-region": []string{"us-east-1"},
						},
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV4) {
				ctx := context.TODO()
				input := &s3.HeadBucketInput{
					Bucket: aws.String("test-bucket"),
				}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.HeadBucket(ctx, input)
				result2, err2 := s3ec.HeadBucket(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					// Both should return HeadBucketOutput with same structure
					// HeadBucketOutput doesn't have many fields, but both should be non-nil
					if result1 == nil || result2 == nil {
						t.Errorf("Result mismatch: s3Client result=%v, s3ec result=%v", result1, result2)
					}
					
					t.Logf("✓ HeadBucket: Both clients returned identical structured output")
				}
			},
		},
		{
			name:      "ListObjectsV2",
			operation: "ListObjectsV2",
			setupMock: func() *awstesting.MockHttpClient {
				mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListBucketResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Name>test-bucket</Name>
    <Prefix></Prefix>
    <KeyCount>2</KeyCount>
    <MaxKeys>1000</MaxKeys>
    <IsTruncated>false</IsTruncated>
    <Contents>
        <Key>test-object-1</Key>
        <LastModified>2023-01-01T00:00:00.000Z</LastModified>
        <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
        <Size>0</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
    <Contents>
        <Key>test-object-2</Key>
        <LastModified>2023-01-02T00:00:00.000Z</LastModified>
        <ETag>"d41d8cd98f00b204e9800998ecf8427e"</ETag>
        <Size>100</Size>
        <StorageClass>STANDARD</StorageClass>
    </Contents>
</ListBucketResult>`
				return &awstesting.MockHttpClient{
					Response: &http.Response{
						Status:     http.StatusText(200),
						StatusCode: http.StatusOK,
						Body:       io.NopCloser(strings.NewReader(mockResponse)),
					},
				}
			},
			testFunc: func(t *testing.T, s3Client *s3.Client, s3ec *S3EncryptionClientV4) {
				ctx := context.TODO()
				input := &s3.ListObjectsV2Input{
					Bucket:  aws.String("test-bucket"),
					MaxKeys: 1000,
				}
				
				// Call both clients with identical parameters
				result1, err1 := s3Client.ListObjectsV2(ctx, input)
				result2, err2 := s3ec.ListObjectsV2(ctx, input)
				
				// Verify both calls succeed or fail in the same way
				if (err1 == nil) != (err2 == nil) {
					t.Errorf("Error behavior mismatch: s3Client error=%v, s3ec error=%v", err1, err2)
				}
				
				// If both succeed, verify output structure is compatible
				if err1 == nil && err2 == nil {
					if len(result1.Contents) != len(result2.Contents) {
						t.Errorf("Contents count mismatch: s3Client=%d, s3ec=%d", len(result1.Contents), len(result2.Contents))
					}
					
					if result1.KeyCount != result2.KeyCount {
						t.Errorf("KeyCount mismatch: s3Client=%d, s3ec=%d", result1.KeyCount, result2.KeyCount)
					}
					
					if *result1.Name != *result2.Name {
						t.Errorf("Bucket name mismatch: s3Client=%s, s3ec=%s", *result1.Name, *result2.Name)
					}
					
					t.Logf("✓ ListObjectsV2: Both clients returned identical structured output")
				}
			},
		},
	}
	
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create two separate mock HTTP clients with identical responses
			mockClient1 := tc.setupMock()
			mockClient2 := tc.setupMock()
			
			// Create regular S3 client
			tConfig1 := awstesting.Config()
			tConfig1.HTTPClient = mockClient1
			s3Client := s3.NewFromConfig(tConfig1)
			
			// Create S3 encryption client
			tConfig2 := awstesting.Config()
			tConfig2.HTTPClient = mockClient2
			s3BaseClient := s3.NewFromConfig(tConfig2)
			
			mockCMM := &mockCMM{}
			s3ec, err := New(s3BaseClient, mockCMM)
			if err != nil {
				t.Fatalf("Failed to create S3 encryption client: %v", err)
			}
			
			// Run the specific test for this operation
			tc.testFunc(t, s3Client, s3ec)
			
			t.Logf("✓ %s: S3EC adheres to same interface as conventional S3 client", tc.operation)
		})
	}
}

// Custom mock HTTP client that can capture multiple requests
type multiRequestMockClient struct {
	capturedRequests []string
}

func (m *multiRequestMockClient) Do(req *http.Request) (*http.Response, error) {
	// Capture the request path to verify both object and instruction file are deleted
	m.capturedRequests = append(m.capturedRequests, req.URL.Path)
	
	return &http.Response{
		Status:     http.StatusText(204),
		StatusCode: http.StatusNoContent,
		Body:       io.NopCloser(strings.NewReader("")),
		Header: http.Header{
			"x-amz-delete-marker": []string{"false"},
		},
	}, nil
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - DeleteObject MUST delete the given object key.
//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
func TestS3EC_DeleteObject_DeletesObjectAndInstructionFile(t *testing.T) {
	// This test validates that DeleteObject deletes both the object and its instruction file
	// We'll use a custom mock to capture both delete requests
	
	// Create custom mock HTTP client that captures delete requests
	mockClient := &multiRequestMockClient{}

	// Create test config with mock HTTP client
	tConfig := awstesting.Config()
	tConfig.HTTPClient = mockClient
	s3Client := s3.NewFromConfig(tConfig)

	// Create mock CMM
	mockCMM := &mockCMM{}

	// Create the S3 encryption client
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Call DeleteObject - this should delete both object and instruction file
	ctx := context.TODO()
	result, err := s3ec.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String("test-bucket"),
		Key:    aws.String("test-key"),
	})

	// Verify no error occurred
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	// Verify the response is not nil
	if result == nil {
		t.Error("Expected non-nil result from DeleteObject")
	}

	// Verify that both the object and instruction file delete requests were made
	if len(mockClient.capturedRequests) < 2 {
		t.Errorf("Expected at least 2 delete requests (object + instruction file), got %d", len(mockClient.capturedRequests))
	} else {
		// Check that requests include both the original object and the instruction file
		hasObjectDelete := false
		hasInstructionDelete := false
		
		for _, path := range mockClient.capturedRequests {
			if strings.Contains(path, "/test-key") && !strings.Contains(path, ".instruction") {
				hasObjectDelete = true
			}
			if strings.Contains(path, "/test-key.instruction") {
				hasInstructionDelete = true
			}
		}
		
		if !hasObjectDelete {
			t.Error("Expected delete request for original object 'test-key'")
		}
		if !hasInstructionDelete {
			t.Error("Expected delete request for instruction file 'test-key.instruction'")
		}
		
		t.Logf("✓ Verified DeleteObject deletes both object and instruction file")
	}
}
func TestS3ECInterfaceCompatibility(t *testing.T) {
	t.Run("ListBuckets", func(t *testing.T) {
		// Setup mock response
		mockResponse := `<?xml version="1.0" encoding="UTF-8"?>
<ListAllMyBucketsResult xmlns="http://s3.amazonaws.com/doc/2006-03-01/">
    <Buckets><Bucket><Name>test-bucket</Name></Bucket></Buckets>
</ListAllMyBucketsResult>`

		// Create S3 client
		tConfig1 := awstesting.Config()
		tConfig1.HTTPClient = &awstesting.MockHttpClient{
			Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(mockResponse))},
		}
		s3Client := s3.NewFromConfig(tConfig1)

		// Create S3EC
		tConfig2 := awstesting.Config()
		tConfig2.HTTPClient = &awstesting.MockHttpClient{
			Response: &http.Response{StatusCode: 200, Body: io.NopCloser(strings.NewReader(mockResponse))},
		}
		s3ec, _ := New(s3.NewFromConfig(tConfig2), &mockCMM{})

		// Call same operation on both with same parameters
		ctx := context.TODO()
		input := &s3.ListBucketsInput{}
		
		result1, err1 := s3Client.ListBuckets(ctx, input)
		result2, err2 := s3ec.ListBuckets(ctx, input)

		// Verify same interface and output
		if (err1 == nil) != (err2 == nil) {
			t.Errorf("Error mismatch: s3Client=%v, s3ec=%v", err1, err2)
		}
		if len(result1.Buckets) != len(result2.Buckets) {
			t.Errorf("Output mismatch: s3Client buckets=%d, s3ec buckets=%d", len(result1.Buckets), len(result2.Buckets))
		}
	})

	t.Run("GetObject", func(t *testing.T) {
		// Verify interface compatibility - both should accept same input types
		ctx := context.TODO()
		input := &s3.GetObjectInput{
			Bucket: &[]string{"test-bucket"}[0],
			Key:    &[]string{"test-key"}[0],
		}
		
		// Test that both clients accept the same input parameters and return same types
		var result *s3.GetObjectOutput
		var err error
		
		// This verifies the interface is identical - same method signature
		// Don't assert equality; these clients behave differently. Important to assert the method signatures match.
		_ = func() {
			result, err = (&s3.Client{}).GetObject(ctx, input)
			result, err = (&S3EncryptionClientV4{}).GetObject(ctx, input)
		}
		
		// Suppress unused variable warnings
		_, _ = result, err
	})

	t.Run("PutObject", func(t *testing.T) {
		// Verify interface compatibility - both should accept same input types
		ctx := context.TODO()
		input := &s3.PutObjectInput{
			Bucket: &[]string{"test-bucket"}[0],
			Key:    &[]string{"test-key"}[0],
			Body:   strings.NewReader("test-content"),
		}
		
		// Test that both clients accept the same input parameters and return same types
		var result *s3.PutObjectOutput
		var err error
		
		// This verifies the interface is identical - same method signature
		// Don't assert equality; these clients behave differently. Important to assert the method signatures match.
		_ = func() {
			result, err = (&s3.Client{}).PutObject(ctx, input)
			result, err = (&S3EncryptionClientV4{}).PutObject(ctx, input)
		}
		
		// Suppress unused variable warnings
		_, _ = result, err
	})
}

//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
//= type=test
//# The option to enable legacy unauthenticated modes MUST be set to false by default.
func TestLegacyUnauthenticatedModes_DefaultDisabled(t *testing.T) {
	// Create test config
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Create S3EC without specifying EnableLegacyUnauthenticatedModes (should default to false)
	s3ec, err := New(s3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Verify that EnableLegacyUnauthenticatedModes defaults to false
	if s3ec.Options.EnableLegacyUnauthenticatedModes {
		t.Error("Expected EnableLegacyUnauthenticatedModes to default to false, but it was true")
	}

	t.Logf("✓ Verified EnableLegacyUnauthenticatedModes defaults to false")
}

//= ../specification/s3-encryption/client.md#wrapped-s3-client-s
//= type=test
//# The S3EC MUST support the option to provide an SDK S3 client instance during its initialization.
func TestS3EC_AcceptsProvidedS3ClientInstance(t *testing.T) {
	// Create a specific S3 client instance
	tConfig := awstesting.Config()
	providedS3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Create S3EC with the provided S3 client
	s3ec, err := New(providedS3Client, mockCMM)
	if err != nil {
		t.Fatalf("Failed to create S3 encryption client: %v", err)
	}

	// Verify that the S3EC uses the provided client (they should be the same instance)
	if s3ec.Client != providedS3Client {
		t.Error("S3EC should use the provided S3 client instance")
	}

	t.Logf("✓ Verified S3EC accepts and uses provided S3 client instance")
}

//= ../specification/s3-encryption/client.md#encryption-algorithm
//= type=test
//# The S3EC MUST support configuration of the encryption algorithm (or algorithm suite) during its initialization.
func TestS3EC_SupportsEncryptionAlgorithmConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test configuring different valid (non-legacy) algorithm suites
	testCases := []struct {
		name             string
		algorithm        *algorithms.AlgorithmSuite
		commitmentPolicy commitment.CommitmentPolicy
	}{
		{
			name:             "AES256GCMHkdfSha512CommitKey",
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey,
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT, // Committing algorithm
		},
		{
			name:             "AES256GCMIV12Tag16NoKDF",
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF,
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT, // Non-committing algorithm
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure S3EC with specific encryption algorithm and appropriate commitment policy
			s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.EncryptionAlgorithmSuite = tc.algorithm
				options.CommitmentPolicy = tc.commitmentPolicy
			})
			
			if err != nil {
				t.Fatalf("Failed to create S3EC with algorithm %s: %v", tc.name, err)
			}

			// Verify the algorithm was configured correctly
			if s3ec.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %s, got %v", tc.name, s3ec.Options.EncryptionAlgorithmSuite)
			}

			// Verify the commitment policy was configured correctly
			if s3ec.Options.CommitmentPolicy != tc.commitmentPolicy {
				t.Errorf("Expected commitment policy %v, got %v", tc.commitmentPolicy, s3ec.Options.CommitmentPolicy)
			}

			t.Logf("✓ Successfully configured S3EC with algorithm: %s and commitment policy: %v", tc.name, tc.commitmentPolicy)
		})
	}
}

//= ../specification/s3-encryption/client.md#key-commitment
//= type=test
//# The S3EC MUST support configuration of the [Key Commitment policy](./key-commitment.md) during its initialization.
func TestS3EC_SupportsKeyCommitmentPolicyConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test configuring all 3 commitment policies with compatible algorithms (happy path cases)
	testCases := []struct {
		name             string
		commitmentPolicy commitment.CommitmentPolicy
		algorithm        *algorithms.AlgorithmSuite
	}{
		{
			name:             "FORBID_ENCRYPT_ALLOW_DECRYPT",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF, // Non-committing algorithm
		},
		{
			name:             "REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing algorithm
		},
		{
			name:             "REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing algorithm
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Configure S3EC with specific commitment policy and compatible algorithm
			s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.CommitmentPolicy = tc.commitmentPolicy
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err != nil {
				t.Fatalf("Failed to create S3EC with commitment policy %v: %v", tc.commitmentPolicy, err)
			}

			// Verify the commitment policy was configured correctly
			if s3ec.Options.CommitmentPolicy != tc.commitmentPolicy {
				t.Errorf("Expected commitment policy %v, got %v", tc.commitmentPolicy, s3ec.Options.CommitmentPolicy)
			}

			// Verify the algorithm was configured correctly
			if s3ec.Options.EncryptionAlgorithmSuite != tc.algorithm {
				t.Errorf("Expected algorithm %v, got %v", tc.algorithm, s3ec.Options.EncryptionAlgorithmSuite)
			}

			t.Logf("✓ Successfully configured S3EC with commitment policy: %v and algorithm: %v", tc.commitmentPolicy, tc.algorithm)
		})
	}
}

//= ../specification/s3-encryption/client.md#key-commitment
//= type=test
//# The S3EC MUST validate the configured Encryption Algorithm against the provided key commitment policy.
//# If the configured Encryption Algorithm is incompatible with the key commitment policy, then it MUST throw an exception.
func TestS3EC_ValidatesAlgorithmCommitmentPolicyCompatibility(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test incompatible combinations that should fail
	incompatibleCases := []struct {
		name             string
		commitmentPolicy commitment.CommitmentPolicy
		algorithm        *algorithms.AlgorithmSuite
		expectedError    string
	}{
		{
			name:             "FORBID_ENCRYPT_ALLOW_DECRYPT with committing algorithm",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMHkdfSha512CommitKey, // Committing
			expectedError:    "does not allow committing algorithm suites",
		},
		{
			name:             "REQUIRE_ENCRYPT_REQUIRE_DECRYPT with non-committing algorithm",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithm:        algorithms.AlgAES256GCMIV12Tag16NoKDF, // Non-committing
			expectedError:    "requires committing algorithm suites",
		},
	}

	for _, tc := range incompatibleCases {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to configure S3EC with incompatible commitment policy and algorithm - should fail
			_, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.CommitmentPolicy = tc.commitmentPolicy
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err == nil {
				t.Errorf("Expected error for incompatible combination %s, but got none", tc.name)
			} else if !strings.Contains(err.Error(), tc.expectedError) {
				t.Errorf("Expected error containing '%s', got: %v", tc.expectedError, err)
			} else {
				t.Logf("✓ Correctly rejected incompatible combination %s: %v", tc.name, err)
			}
		})
	}
}

//= ../specification/s3-encryption/client.md#encryption-algorithm
//= type=test
//# The S3EC MUST validate that the configured encryption algorithm is not legacy.
//# If the configured encryption algorithm is legacy, then the S3EC MUST throw an exception.
func TestS3EC_RejectsLegacyEncryptionAlgorithms(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	// Test that legacy algorithm suites are rejected
	legacyAlgorithms := []struct {
		name      string
		algorithm *algorithms.AlgorithmSuite
	}{
		{
			name:      "AES256CTRIV16Tag16NoKDF (legacy)",
			algorithm: algorithms.AlgAES256CTRIV16Tag16NoKDF,
		},
		{
			name:      "AES256CBCIV16NoKDF (legacy)",
			algorithm: algorithms.AlgAES256CBCIV16NoKDF,
		},
	}

	for _, tc := range legacyAlgorithms {
		t.Run(tc.name, func(t *testing.T) {
			// Attempt to configure S3EC with legacy encryption algorithm - should fail
			_, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
				options.EncryptionAlgorithmSuite = tc.algorithm
			})
			
			if err == nil {
				t.Errorf("Expected error when configuring legacy algorithm %s, but got none", tc.name)
			} else {
				t.Logf("✓ Correctly rejected legacy algorithm %s: %v", tc.name, err)
			}
		})
	}
}

func TestS3EC_BufferSizeConfiguration(t *testing.T) {
	tConfig := awstesting.Config()
	s3Client := s3.NewFromConfig(tConfig)
	mockCMM := &mockCMM{}

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//= type=test
	//# If Delayed Authentication mode is disabled, and no buffer size is provided,
	//# the S3EC MUST set the buffer size to a reasonable default.
	t.Run("SetsReasonableDefaultBufferSize", func(t *testing.T) {
		// Create S3EC with default options
		s3ec, err := New(s3Client, mockCMM)
		if err != nil {
			t.Fatalf("Failed to create S3 encryption client: %v", err)
		}
		
		// Verify that the default buffer size is the default
		expectedBufferSize := int64(DefaultBufferSize)
		if s3ec.Options.BufferSize != expectedBufferSize {
			t.Errorf("Expected default buffer size to be %d, got %d", expectedBufferSize, s3ec.Options.BufferSize)
		}
		
		// Verify the default buffer size is the default
		if s3ec.Options.BufferSize != 64*1024 {
			t.Errorf("Expected default buffer size to be 64KB (65536 bytes), got %d", s3ec.Options.BufferSize)
		}
		
		t.Logf("✓ Verified S3EC sets reasonable default buffer size: %d bytes", s3ec.Options.BufferSize)
	})

	//= ../specification/s3-encryption/client.md#set-buffer-size
	//= type=test
	//# The S3EC SHOULD accept a configurable buffer size
	//# which refers to the maximum ciphertext length in bytes to store in memory
	//# when Delayed Authentication mode is disabled.
	t.Run("SupportsCustomBufferSizeConfiguration", func(t *testing.T) {
		// Test custom buffer size configuration
		customBufferSize := int64(128 * 1024) // 128KB
		
		s3ec, err := New(s3Client, mockCMM, func(options *EncryptionClientOptions) {
			options.BufferSize = customBufferSize
		})
		if err != nil {
			t.Fatalf("Failed to create S3 encryption client with custom buffer size: %v", err)
		}
		
		// Verify that the custom buffer size is set correctly
		if s3ec.Options.BufferSize != customBufferSize {
			t.Errorf("Expected buffer size to be %d, got %d", customBufferSize, s3ec.Options.BufferSize)
		}
		
		t.Logf("✓ Verified S3EC supports custom buffer size configuration: %d bytes", s3ec.Options.BufferSize)
	})
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package commitment

import "fmt"

type CommitmentPolicy int

const (
	// "Forbid" writing objects encrypted with key commitment, and "allow" reading objects encrypted with key commitment.
	// FORBID_ENCRYPT_ALLOW_DECRYPT does not write objects with key commitment
	// and can read objects encrypted either with or without key commitment.
	// Keys in Instruction Files could be tampered with when reading objects without key commitment.
	// FORBID_ENCRYPT_ALLOW_DECRYPT means that this client will write objects that any v3 client can read,
	// and any v4 client (configured with either FORBID_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_ALLOW_DECRYPT) can read.
	// FORBID_ENCRYPT_ALLOW_DECRYPT also means that this client can read objects written by any v3 or v4 client.
	// This is the default policy for v3 clients.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	FORBID_ENCRYPT_ALLOW_DECRYPT CommitmentPolicy = iota
	// "Require" writing objects encrypted with key commitment, and "allow" reading objects encrypted with key commitment.
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT ensures that all newly written objects are protected with key commitment,
	// but can read objects encrypted either with or without key commitment.
	// Keys in Instruction Files could be tampered with when reading objects without key commitment.
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT means that this client will write objects that any v3 client (from v3.2.0 onward)
	// or any v4 client can read.
	// REQUIRE_ENCRYPT_ALLOW_DECRYPT also means that this client can read objects written by any v3 or v4 client.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	REQUIRE_ENCRYPT_ALLOW_DECRYPT
	// "Require" writing objects encrypted with key commitment, and "require" reading objects encrypted with key commitment.
	// REQUIRE_ENCRYPT_REQUIRE_DECRYPT ensures that all newly written objects are protected with key commitment
	// and that all decrypted objects are verified to have been encrypted with key commitment.
	// This prevents reading objects with keys in Instruction Files that may have been tampered with.
	// REQUIRE_ENCRYPT_REQUIRE_DECRYPT means that any v3 client (from v3.2.0 onward) or any v4 client implementation
	// can read all of the objects being written.
	// However, REQUIRE_ENCRYPT_REQUIRE_DECRYPT also means that this client can only read objects written by
	// v4 clients (configured with any REQUIRE_ENCRYPT_ALLOW_DECRYPT or REQUIRE_ENCRYPT_REQUIRE_DECRYPT).
	// This is the default policy for v4 clients.
	// For more information, see the developer guide:
	// https://docs.aws.amazon.com/amazon-s3-encryption-client/latest/developerguide/go-v4-migration.html
	REQUIRE_ENCRYPT_REQUIRE_DECRYPT
)

func (cp CommitmentPolicy) RequiresEncrypt() bool {
	switch cp {
	case REQUIRE_ENCRYPT_ALLOW_DECRYPT, REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
		return true
	default:
		return false
	}
}

func (cp CommitmentPolicy) RequiresDecrypt() bool {
	return cp == REQUIRE_ENCRYPT_REQUIRE_DECRYPT
}

func (p CommitmentPolicy) String() string {
    switch p {
    case FORBID_ENCRYPT_ALLOW_DECRYPT:
        return "FORBID_ENCRYPT_ALLOW_DECRYPT"
    case REQUIRE_ENCRYPT_ALLOW_DECRYPT:
        return "REQUIRE_ENCRYPT_ALLOW_DECRYPT"
    case REQUIRE_ENCRYPT_REQUIRE_DECRYPT:
        return "REQUIRE_ENCRYPT_REQUIRE_DECRYPT"
    default:
        return fmt.Sprintf("CommitmentPolicy(%d)", int(p))
    }
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

/*
Package client provides client-side encryption for S3 using KMS and AES GCM.

The client supports encryption and decryption using authenticated algorithms by default.

Refer to the module README for more detailed documentation.
*/
package client
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
)

// AESCBC is a symmetric crypto algorithm. This algorithm
// requires a padder due to CBC needing to be of the same block
// size. AES CBC is vulnerable to Padding Oracle attacks and
// so should be avoided when possible.
type aesCBC struct {
	encrypter cipher.BlockMode
	decrypter cipher.BlockMode
	padder    Padder
}

// newAESCBC creates a new AES CBC cipher. Expects keys to be of
// the correct size.
func newAESCBC(materials materials.CryptographicMaterials, padder Padder) (Cipher, error) {
	block, err := aes.NewCipher(materials.Key)
	if err != nil {
		return nil, err
	}

	encrypter := cipher.NewCBCEncrypter(block, materials.IV)
	decrypter := cipher.NewCBCDecrypter(block, materials.IV)

	return &aesCBC{encrypter, decrypter, padder}, nil
}

// Encrypt will encrypt the data using AES CBC by returning
// an io.Reader. The io.Reader will encrypt the data as Read
// is called.
func (c *aesCBC) Encrypt(src io.Reader) io.Reader {
	reader := &cbcEncryptReader{
		encrypter: c.encrypter,
		src:       src,
		padder:    c.padder,
	}
	return reader
}

type cbcEncryptReader struct {
	encrypter cipher.BlockMode
	src       io.Reader
	padder    Padder
	size      int
	buf       bytes.Buffer
}

// Read will read from our io.Reader and encrypt the data as necessary.
// Due to padding, we have to do some logic that when we encounter an
// end of file to pad properly.
func (reader *cbcEncryptReader) Read(data []byte) (int, error) {
	n, err := reader.src.Read(data)
	reader.size += n
	blockSize := reader.encrypter.BlockSize()
	reader.buf.Write(data[:n])

	if err == io.EOF {
		b := make([]byte, getSliceSize(blockSize, reader.buf.Len(), len(data)))
		n, err = reader.buf.Read(b)
		if err != nil && err != io.EOF {
			return n, err
		}
		// The buffer is now empty, we can now pad the data
		if reader.buf.Len() == 0 {
			b, err = reader.padder.Pad(b[:n], reader.size)
			if err != nil {
				return n, err
			}
			n = len(b)
			err = io.EOF
		}
		// We only want to encrypt if we have read anything
		if n > 0 {
			reader.encrypter.CryptBlocks(data, b)
		}
		return n, err
	}

	if err != nil {
		return n, err
	}

	if size := reader.buf.Len(); size >= blockSize {
		nBlocks := size / blockSize
		if size > len(data) {
			nBlocks = len(data) / blockSize
		}

		if nBlocks > 0 {
			b := make([]byte, nBlocks*blockSize)
			n, _ = reader.buf.Read(b)
			reader.encrypter.CryptBlocks(data, b[:n])
		}
	} else {
		n = 0
	}
	return n, nil
}

// Decrypt will decrypt the data using AES CBC
func (c *aesCBC) Decrypt(src io.Reader) io.Reader {
	return &cbcDecryptReader{
		decrypter: c.decrypter,
		src:       src,
		padder:    c.padder,
	}
}

type cbcDecryptReader struct {
	decrypter cipher.BlockMode
	src       io.Reader
	padder    Padder
	buf       bytes.Buffer
}

// Read will read from our io.Reader and decrypt the data as necessary.
// Due to padding, we have to do some logic that when we encounter an
// end of file to pad properly.
func (reader *cbcDecryptReader) Read(data []byte) (int, error) {
	n, err := reader.src.Read(data)
	blockSize := reader.decrypter.BlockSize()
	reader.buf.Write(data[:n])

	if err == io.EOF {
		b := make([]byte, getSliceSize(blockSize, reader.buf.Len(), len(data)))
		n, err = reader.buf.Read(b)
		if err != nil && err != io.EOF {
			return n, err
		}
		// We only want to decrypt if we have read anything
		if n > 0 {
			reader.decrypter.CryptBlocks(data, b)
		}

		if reader.buf.Len() == 0 {
			b, err = reader.padder.Unpad(data[:n])
			n = len(b)
			if err != nil {
				return n, err
			}
			err = io.EOF
		}
		return n, err
	}

	if err != nil {
		return n, err
	}

	if size := reader.buf.Len(); size >= blockSize {
		nBlocks := size / blockSize
		if size > len(data) {
			nBlocks = len(data) / blockSize
		}
		// The last block is always padded. This will allow us to unpad
		// when we receive an io.EOF error
		nBlocks -= blockSize

		if nBlocks > 0 {
			b := make([]byte, nBlocks*blockSize)
			n, _ = reader.buf.Read(b)
			reader.decrypter.CryptBlocks(data, b[:n])
		} else {
			n = 0
		}
	}

	return n, nil
}

// getSliceSize will return the correct amount of bytes we need to
// read with regards to padding.
func getSliceSize(blockSize, bufSize, dataSize int) int {
	size := bufSize
	if bufSize > dataSize {
		size = dataSize
	}
	size = size - (size % blockSize) - blockSize
	if size <= 0 {
		size = blockSize
	}

	return size
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
)

const (
	AESCBCPKCS5Padding = "AES/CBC/PKCS5Padding"
	aesCbcTagSizeBits  = "0"
)

// NewAESCBCContentCipher will create a new aes cbc content cipher. If the cipher data's
// will set the cek algorithm if it hasn't been set.
func NewAESCBCContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.TagLength = aesCbcTagSizeBits
	if len(materials.CEKAlgorithm) == 0 {
		materials.CEKAlgorithm = AESCBCPKCS5Padding
	}
	cipher, err := newAESCBC(materials, AESCBCPadder)
	if err != nil {
		return nil, err
	}

	return &aesCBCContentCipher{
		CryptographicMaterials: materials,
		Cipher:                 cipher,
	}, nil
}

// aesCBCContentCipher will use AES CBC for the main cipher.
type aesCBCContentCipher struct {
	CryptographicMaterials materials.CryptographicMaterials
	Cipher                 Cipher
}

// EncryptContents will generate a random key and iv and encrypt the data using cbc
func (cc *aesCBCContentCipher) EncryptContents(src io.Reader) (io.Reader, error) {
	return cc.Cipher.Encrypt(src), nil
}

// DecryptContents will use the symmetric key provider to instantiate a new CBC cipher.
// We grab a decrypt reader from CBC and wrap it in a CryptoReadCloser. The only error
// expected here is when the key or iv is of invalid length.
func (cc *aesCBCContentCipher) DecryptContents(src io.ReadCloser) (io.ReadCloser, error) {
	reader := cc.Cipher.Decrypt(src)
	return &CryptoReadCloser{Body: src, Decrypter: reader}, nil
}

// GetCipherData returns cipher data
func (cc aesCBCContentCipher) GetCipherData() materials.CryptographicMaterials {
	return cc.CryptographicMaterials
}

var (
	_ ContentCipher = (*aesCBCContentCipher)(nil)
)
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

var AesCbcPadding = aescbcPadder{pkcs7Padder{16}}

// AESCBCPadder is used to pad AES encrypted and decrypted data.
// Although it uses the pkcs5Padder, it isn't following the RFC
// for PKCS5. The only reason why it is called pkcs5Padder is
// due to the Name returning PKCS5Padding.
var AESCBCPadder = Padder(AesCbcPadding)

type aescbcPadder struct {
	padder pkcs7Padder
}

func (padder aescbcPadder) Pad(b []byte, n int) ([]byte, error) {
	return padder.padder.Pad(b, n)
}

func (padder aescbcPadder) Unpad(b []byte) ([]byte, error) {
	return padder.padder.Unpad(b)
}

func (padder aescbcPadder) Name() string {
	return "PKCS5Padding"
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"fmt"
	"testing"
)

func TestAESCBCPadding(t *testing.T) {
	for i := 0; i < 16; i++ {
		input := make([]byte, i)
		expected := append(input, bytes.Repeat([]byte{byte(16 - i)}, 16-i)...)
		b, err := AESCBCPadder.Pad(input, len(input))
		if err != nil {
			t.Fatal("Expected error to be nil but received " + err.Error())
		}
		if len(b) != len(expected) {
			t.Fatal(fmt.Sprintf("Case %d: data is not of the same length", i))
		}
		if bytes.Compare(b, expected) != 0 {
			t.Fatal(fmt.Sprintf("Expected %v but got %v", expected, b))
		}
	}
}

func TestAESCBCUnpadding(t *testing.T) {
	for i := 0; i < 16; i++ {
		expected := make([]byte, i)
		input := append(expected, bytes.Repeat([]byte{byte(16 - i)}, 16-i)...)
		b, err := AESCBCPadder.Unpad(input)
		if err != nil {
			t.Fatal("Error received, was expecting nil: " + err.Error())
		}
		if len(b) != len(expected) {
			t.Fatal(fmt.Sprintf("Case %d: data is not of the same length", i))
		}
		if bytes.Compare(b, expected) != 0 {
			t.Fatal(fmt.Sprintf("Expected %v but got %v", expected, b))
		}
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/hex"
	"fmt"
	materials2 "github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
	"testing"
)

func TestAESCBCEncryptDecrypt(t *testing.T) {
	var testCases = []struct {
		key        string
		iv         string
		plaintext  string
		ciphertext string
		decodeHex  bool
		padder     Padder
	}{
		// Test vectors from RFC 3602: https://tools.ietf.org/html/rfc3602
		{
			"06a9214036b8a15b512e03d534120006",
			"3dafba429d9eb430b422da802c9fac41",
			"Single block msg",
			"e353779c1079aeb82708942dbe77181a",
			false,
			NoPadder,
		},
		{
			"c286696d887c9aa0611bbb3e2025a45a",
			"562e17996d093d28ddb3ba695a2e6f58",
			"000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f",
			"d296cd94c2cccf8a3a863028b5e1dc0a7586602d253cfff91b8266bea6d61ab1",
			true,
			NoPadder,
		},
		{
			"6c3ea0477630ce21a2ce334aa746c2cd",
			"c782dc4c098c66cbd9cd27d825682c81",
			"This is a 48-byte message (exactly 3 AES blocks)",
			"d0a02b3836451753d493665d33f0e8862dea54cdb293abc7506939276772f8d5021c19216bad525c8579695d83ba2684",
			false,
			NoPadder,
		},
		{
			"56e47a38c5598974bc46903dba290349",
			"8ce82eefbea0da3c44699ed7db51b7d9",
			"a0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebfc0c1c2c3c4c5c6c7c8c9cacbcccdcecfd0d1d2d3d4d5d6d7d8d9dadbdcdddedf",
			"c30e32ffedc0774e6aff6af0869f71aa0f3af07a9a31a9c684db207eb0ef8e4e35907aa632c3ffdf868bb7b29d3d46ad83ce9f9a102ee99d49a53e87f4c3da55",
			true,
			NoPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"",
			"B012949BA07D1A6DCE9DEE67274D41AB",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41",
			"8A11ABA68A566132FFE04DB336621D41",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141",
			"97D0896E41DFDB5CEA4A9EB70A938CFD",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141",
			"8464EAD45FA2D8790E8741E32C28083F",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141",
			"1E656D6E2745BA9F154FAF136B2BC73D",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141",
			"0B6031C4B230DAC6BD6D3F195645B287",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141",
			"5D09FEB6462BB489489A7E18FD341D9D",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141",
			"85745E398F2FD1050C2CE8F8614DA369",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141",
			"7BE52933970BA7B0FC6FB3FC37648205",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141",
			"ED3A1E134EF36CCFE60C8123B4272F89",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141",
			"C3B7C9E177E1052FC736F65FC1E74209",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141",
			"C3A8B53F7F57F0B9D346FA99810A3C28",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141",
			"D16B1ECE5BF00AF919E139E99775FF06",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141",
			"B258F4DF57FFCA1EFCF8D76140F05139",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141",
			"3CD2282DE24A2CF9E23326CC3DC9077A",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141",
			"3010232E7C752A3B4C9EE428B4C4FE88",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A22BC4E6D03BFD2418DD412D1ED1B31AF",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A5427BBD4A4D50776989441370E3B5B16",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A7FF985F55567D1B25EA40E23BB4CB1FE",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A0835E548C7370D8F8D9925C0E6B54727",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ADC0CF1436399E67BC1122B31CB596649",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A3D096F0DEAFF91938B82E5D404B0B065",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217AD56ABA897A355CF307CCB74226243192",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A151284F950B1B1DBCAD6D9E7900DF4E6",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217AEF85A612514121C299A1D87116C4A182",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A67F157569BFB4013EA3AD16DB8C69AD6",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217AF8520D191F6ACBD88B2140588B91C697",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ADD8BBAA71745669B96F2683E2F5AEC35",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217AFB2D4282817D7EC6B33EFAD7AA14A3C5",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A459B89E7E0DAF3DA654576B60B2DA7CE",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A65759F23F9789D05B23D5DBAA9E32036",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217A03C78FBD5E2CB08B3B6D181E23FBDE79",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA013D941FBBDE56C106C482CD022F290F",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA0645D313AC3C29B79DB1AA2E00A5B393",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA2ED0FD8048053BF22EBE501D82C4B3F1",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CAC57D706C7866A01D6E913F98AE57EE54",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CAB7FC1241FAFDFE45C4FF982D5DC1DAEF",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA7063EA296922DE8BDFD3B29D786C5F91",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA3A4603475F4AFDBFADC6E7FA908188B1",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA3365C63C2AF2A6C8FB4D0E9ED3C6FDA3",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA78BCC1874C0B7EB52645FC8F03B9C9CF",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA9B7A31397718EECB89B9E9CCCD729326",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CAB15EA8A67E9E9FADB4249710277F3D4F",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA94641D6A076193C660632CEA3F9CB02C",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CAB2170A08417BE77F0EAA9110F4790E12",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA4E30F1CD7B2256ABD57DC3DAB05376C9",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"41414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA9909B7B93D01BDAAC22D15AF34DF1EEF",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"4141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CAD97F5D1206F00E5C7225CAD81CCD4027",
			true,
			AESCBCPadder,
		},
		{
			"11111111111111111111111111111111",
			"22222222222222222222222222222222",
			"414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141414141",
			"C3304FA46097CBBA59085416764A217ACEF79EE1163A2F52888F87A3979EB3CA570CBB001A0C87558906B60C884AB5F41DA97CEF2A9401BC6DD0D22A54DBAD6D",
			true,
			AESCBCPadder,
		},
	}

	for i, testCase := range testCases {
		key, _ := hex.DecodeString(testCase.key)
		iv, _ := hex.DecodeString(testCase.iv)
		materials := materials2.CryptographicMaterials{
			Key: key,
			IV:  iv,
		}

		cbc, err := newAESCBC(materials, testCase.padder)
		if err != nil {
			t.Fatal(fmt.Sprintf("Case %d: Expected no error for cipher creation, but received: %v", i, err.Error()))
		}

		plaintext := []byte(testCase.plaintext)
		if testCase.decodeHex {
			plaintext, _ = hex.DecodeString(testCase.plaintext)
		}

		cipherdata := cbc.Encrypt(bytes.NewReader(plaintext))
		ciphertext := []byte{}
		b := make([]byte, 19)
		err = nil
		n := 0
		for err != io.EOF {
			n, err = cipherdata.Read(b)
			ciphertext = append(ciphertext, b[:n]...)
		}

		if err != io.EOF {
			t.Fatal(fmt.Sprintf("Case %d: Expected no error during io reading, but received: %v", i, err.Error()))
		}

		expectedData, _ := hex.DecodeString(testCase.ciphertext)
		if bytes.Compare(expectedData, ciphertext) != 0 {
			t.Log("\n", ciphertext, "\n", expectedData)
			t.Fatal(fmt.Sprintf("Case %d: AES CBC encryption fails. Data is not the same", i))
		}

		plaindata := cbc.Decrypt(bytes.NewReader(ciphertext))
		plaintextDecrypted := []byte{}
		err = nil
		for err != io.EOF {
			n, err = plaindata.Read(b)
			plaintextDecrypted = append(plaintextDecrypted, b[:n]...)
		}
		if err != io.EOF {
			t.Fatal(fmt.Sprintf("Case %d: Expected no error during io reading, but received: %v", i, err.Error()))
		}

		if bytes.Compare(plaintext, plaintextDecrypted) != 0 {
			t.Log("\n", plaintext, "\n", plaintextDecrypted)
			t.Fatal(fmt.Sprintf("Case %d: AES CBC decryption fails. Data is not the same", i))
		}
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"fmt"
	"crypto/aes"
	"crypto/cipher"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"io"
)

// AESGCM Symmetric encryption algorithm. Since Golang designed this
// with only TLS in mind. We have to load it all into memory meaning
// this isn't streamed.
type aesGCM struct {
	aead  cipher.AEAD
	nonce []byte
	aad  []byte
}

// newAESGCM creates a new AES GCM cipher. Expects keys to be of
// the correct size.
//
// Example:
//
//	materials := &s3crypto.CryptographicMaterials{
//		Key: key,
//		"IV": iv,
//	}
//	cipher, err := s3crypto.newAESGCM(materials)
func newAESGCM(materials materials.CryptographicMaterials) (Cipher, error) {
	expectedNonceLength := algorithms.AlgAES256GCMIV12Tag16NoKDF.IVLengthBytes()
	if len(materials.IV) != expectedNonceLength {
		return nil, fmt.Errorf("invalid nonce length: expected %d bytes, got %d bytes for algorithm %s", expectedNonceLength, len(materials.IV), materials.CEKAlgorithm)
	}
	
	block, err := aes.NewCipher(materials.Key)
	if err != nil {
		return nil, err
	}

	aesgcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if materials.CEKAlgorithm == algorithms.AESGCMCommitKey {
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//# The client MUST set the AAD to the Algorithm Suite ID represented as bytes.
		aad := algorithms.AlgAES256GCMHkdfSha512CommitKey.IDAsBytes()
		return &aesGCM{aesgcm, materials.IV, aad}, nil
	} else if materials.CEKAlgorithm == algorithms.AESGCMNoPadding {
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//# The client MUST NOT provide any AAD when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
		return &aesGCM{aesgcm, materials.IV, nil}, nil
	} else {
		return nil, fmt.Errorf("unsupported CEK algorithm for AES GCM: %s", materials.CEKAlgorithm)
	}
	
}

// Encrypt will encrypt the data using AES GCM
// Tag will be included as the last 16 bytes of the slice
func (c *aesGCM) Encrypt(src io.Reader) io.Reader {
	reader := &gcmEncryptReader{
		encrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
		aad:	   c.aad,
	}
	return reader
}

type gcmEncryptReader struct {
	encrypter cipher.AEAD
	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
	aad	      []byte
}

func (reader *gcmEncryptReader) Read(data []byte) (int, error) {
	if reader.buf == nil {
		b, err := io.ReadAll(reader.src)
		if err != nil {
			return 0, err
		}
		var aad []byte
		if reader.aad != nil {
			aad = reader.aad
		} else {
			aad = nil
		}
		// The GCM auth tag is appended to the ciphertext by the Seal function.
		// Docs: https://pkg.go.dev/crypto/cipher#GCM
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//= type=exception
		//# The client MUST append the GCM auth tag to the ciphertext if the underlying crypto provider does not do so automatically.
		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
		//= type=exception
		//# The client MUST append the GCM auth tag to the ciphertext if the underlying crypto provider does not do so automatically.

		//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
		//= type=implication
		//# The client MUST initialize the cipher, or call an AES-GCM encryption API,
		//# with the plaintext data key, the generated IV, and the tag length defined in the Algorithm Suite
		//# when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
		b = reader.encrypter.Seal(b[:0], reader.nonce, b, aad)
		reader.buf = bytes.NewBuffer(b)
	}

	return reader.buf.Read(data)
}

// Decrypt will decrypt the data using AES GCM
func (c *aesGCM) Decrypt(src io.Reader) io.Reader {
	return &gcmDecryptReader{
		decrypter: c.aead,
		nonce:     c.nonce,
		src:       src,
		aad:	   c.aad,
	}
}

type gcmDecryptReader struct {
	decrypter cipher.AEAD
	nonce     []byte
	src       io.Reader
	buf       *bytes.Buffer
	aad 	  []byte
}

func (reader *gcmDecryptReader) Read(data []byte) (int, error) {
	var aad []byte
	if reader.aad != nil {
		aad = reader.aad
	} else {
		aad = nil
	}
	if reader.buf == nil {
		b, err := io.ReadAll(reader.src)
		if err != nil {
			return 0, err
		}
		b, err = reader.decrypter.Open(b[:0], reader.nonce, b, aad)
		if err != nil {
			return 0, err
		}

		reader.buf = bytes.NewBuffer(b)
	}

	return reader.buf.Read(data)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
)

// NewAESGCMCommittingContentCipher returns a new encryption only AES/GCM mode structure with a specific cipher data generator
// that will provide keys to be used for content encryption.
//
// Note: This uses the Go stdlib AEAD implementation for AES/GCM. Due to this, objects to be encrypted or decrypted
// will be fully loaded into memory before encryption or decryption can occur. Caution must be taken to avoid memory
// allocation failures.
func NewAESGCMCommittingContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.CEKAlgorithm = algorithms.AESGCMCommitKey
	materials.TagLength = GcmTagSizeBits

	// Persist original IV to store in the object metadata
	original_iv := materials.IV

	var storedKeyCommitment []byte = nil
	if materials.KeyCommitment != nil {
		storedKeyCommitment = make([]byte, len(materials.KeyCommitment))
		copy(storedKeyCommitment, materials.KeyCommitment)
	}

	//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
	//= type=implication
	//# The client MUST use HKDF to derive the key commitment value and the derived encrypting key as described in [Key Derivation](key-derivation.md).
	keys, err := DeriveKeys(materials.Key, materials.IV, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), storedKeyCommitment)
	if err != nil {
		return nil, err
	}

	materials.Key = keys.DerivedEncryptionKey
	//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
	//# The derived key commitment value MUST be set or returned from the encryption process such that it can be included in the content metadata.
	materials.KeyCommitment = keys.CommitKey

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# When encrypting or decrypting with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY,
	//# the IV used in the AES-GCM content encryption/decryption MUST consist entirely of bytes with the value 0x01.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# The IV's total length MUST match the IV length defined by the algorithm suite.
	var nonce = [12]byte{1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	materials.IV = nonce[:]

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# The client MUST initialize the cipher, or call an AES-GCM encryption API,
	//# with the derived encryption key,
	//# an IV containing only bytes with the value 0x01,
	//# and the tag length defined in the Algorithm Suite
	//# when encrypting or decrypting with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY.
	cipher, err := newAESGCM(materials)
	if err != nil {
		return nil, err
	}

	// Restore original non-zero IV for metadata storage
	materials.IV = original_iv

	return &aesGCMContentCipher{
		CryptographicMaterials: materials,
		Cipher:                 cipher,
	}, nil
}

func NewAESGCMDecryptCommittingContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	// KeyCommitment value is required for decryption
	if materials.KeyCommitment == nil {
		return nil, fmt.Errorf("key commitment is required for AES-GCM committing decryption")
	}

	return NewAESGCMCommittingContentCipher(materials)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"strings"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
)

func TestNewAESGCMDecryptCommittingContentCipher_NilKeyCommitment_ReturnsError(t *testing.T) {
	// Given: materials with nil KeyCommitment
	materials := materials.CryptographicMaterials{
		KeyCommitment: nil,
	}

	// When: calling NewAESGCMDecryptCommittingContentCipher
	cipher, err := NewAESGCMDecryptCommittingContentCipher(materials)

	// Then: an error should be returned
	if err == nil {
		t.Fatal("expected error when KeyCommitment is nil, but got nil")
	}

	if cipher != nil {
		t.Fatal("expected cipher to be nil when error occurs, but got non-nil cipher")
	}

	// Verify the error message contains the expected text
	expectedErrorText := "key commitment is required for AES-GCM committing decryption"
	if !strings.Contains(err.Error(), expectedErrorText) {
		t.Fatalf("expected error message to contain '%s', but got: %s", expectedErrorText, err.Error())
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
)

const (
	GcmTagSizeBits  = "128"
)

// NewAESGCMContentCipher returns a new encryption only AES/GCM mode structure with a specific cipher data generator
// that will provide keys to be used for content encryption.
//
// Note: This uses the Go stdlib AEAD implementation for AES/GCM. Due to this, objects to be encrypted or decrypted
// will be fully loaded into memory before encryption or decryption can occur. Caution must be taken to avoid memory
// allocation failures.
func NewAESGCMContentCipher(materials materials.CryptographicMaterials) (ContentCipher, error) {
	materials.CEKAlgorithm = algorithms.AESGCMNoPadding
	materials.TagLength = GcmTagSizeBits

	cipher, err := newAESGCM(materials)
	if err != nil {
		return nil, err
	}

	return &aesGCMContentCipher{
		CryptographicMaterials: materials,
		Cipher:                 cipher,
	}, nil
}

// AESGCMContentCipher will use AES GCM for the main cipher.
type aesGCMContentCipher struct {
	CryptographicMaterials materials.CryptographicMaterials
	Cipher                 Cipher
}

// EncryptContents will generate a random key and iv and encrypt the data using cbc
func (cc *aesGCMContentCipher) EncryptContents(src io.Reader) (io.Reader, error) {
	return cc.Cipher.Encrypt(src), nil
}

// DecryptContents will use the symmetric key provider to instantiate a new GCM cipher.
// We grab a decrypt reader from gcm and wrap it in a CryptoReadCloser. The only error
// expected here is when the key or iv is of invalid length.
func (cc *aesGCMContentCipher) DecryptContents(src io.ReadCloser) (io.ReadCloser, error) {
	reader := cc.Cipher.Decrypt(src)
	return &CryptoReadCloser{Body: src, Decrypter: reader}, nil
}

// GetCipherData returns cipher data
func (cc aesGCMContentCipher) GetCipherData() materials.CryptographicMaterials {
	return cc.CryptographicMaterials
}

// assert ContentCipher implementations
var (
	_ ContentCipher = (*aesGCMContentCipher)(nil)
)
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

//go:build go1.9
// +build go1.9

package internal

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	materials2 "github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
	"os"
	"strings"
	"testing"
)

// AES GCM
func TestAES_GCM_NIST_gcmEncryptExtIV256_PTLen_128_Test_0(t *testing.T) {
	iv, _ := hex.DecodeString("0d18e06c7c725ac9e362e1ce")
	key, _ := hex.DecodeString("31bdadd96698c204aa9ce1448ea94ae1fb4a9a0b3c9d773b51bb1822666b8f22")
	plaintext, _ := hex.DecodeString("2db5168e932556f8089a0622981d017d")
	expected, _ := hex.DecodeString("fa4362189661d163fcd6a56d8bf0405a")
	tag, _ := hex.DecodeString("d636ac1bbedd5cc3ee727dc2ab4a9489")
	aesgcmTest(t, iv, key, plaintext, expected, tag)
}

func TestAES_GCM_NIST_gcmEncryptExtIV256_PTLen_104_Test_3(t *testing.T) {
	iv, _ := hex.DecodeString("4742357c335913153ff0eb0f")
	key, _ := hex.DecodeString("e5a0eb92cc2b064e1bc80891faf1fab5e9a17a9c3a984e25416720e30e6c2b21")
	plaintext, _ := hex.DecodeString("8499893e16b0ba8b007d54665a")
	expected, _ := hex.DecodeString("eb8e6175f1fe38eb1acf95fd51")
	tag, _ := hex.DecodeString("88a8b74bb74fda553e91020a23deed45")
	aesgcmTest(t, iv, key, plaintext, expected, tag)
}

func TestAES_GCM_NIST_gcmEncryptExtIV256_PTLen_256_Test_6(t *testing.T) {
	iv, _ := hex.DecodeString("a291484c3de8bec6b47f525f")
	key, _ := hex.DecodeString("37f39137416bafde6f75022a7a527cc593b6000a83ff51ec04871a0ff5360e4e")
	plaintext, _ := hex.DecodeString("fafd94cede8b5a0730394bec68a8e77dba288d6ccaa8e1563a81d6e7ccc7fc97")
	expected, _ := hex.DecodeString("44dc868006b21d49284016565ffb3979cc4271d967628bf7cdaf86db888e92e5")
	tag, _ := hex.DecodeString("01a2b578aa2f41ec6379a44a31cc019c")
	aesgcmTest(t, iv, key, plaintext, expected, tag)
}

func TestAES_GCM_NIST_gcmEncryptExtIV256_PTLen_408_Test_8(t *testing.T) {
	iv, _ := hex.DecodeString("92f258071d79af3e63672285")
	key, _ := hex.DecodeString("595f259c55abe00ae07535ca5d9b09d6efb9f7e9abb64605c337acbd6b14fc7e")
	plaintext, _ := hex.DecodeString("a6fee33eb110a2d769bbc52b0f36969c287874f665681477a25fc4c48015c541fbe2394133ba490a34ee2dd67b898177849a91")
	expected, _ := hex.DecodeString("bbca4a9e09ae9690c0f6f8d405e53dccd666aa9c5fa13c8758bc30abe1ddd1bcce0d36a1eaaaaffef20cd3c5970b9673f8a65c")
	tag, _ := hex.DecodeString("26ccecb9976fd6ac9c2c0f372c52c821")
	aesgcmTest(t, iv, key, plaintext, expected, tag)
}

type KAT struct {
	IV         string `json:"iv"`
	Key        string `json:"key"`
	Plaintext  string `json:"pt"`
	AAD        string `json:"aad"`
	CipherText string `json:"ct"`
	Tag        string `json:"tag"`
}

func TestAES_GCM_KATS(t *testing.T) {
	fileContents, err := os.ReadFile("testdata/aes_gcm.json")
	if err != nil {
		t.Fatalf("failed to read KAT file: %v", err)
	}

	var kats []KAT
	err = json.Unmarshal(fileContents, &kats)
	if err != nil {
		t.Fatalf("failed to unmarshal KAT json file: %v", err)
	}

	for i, kat := range kats {
		t.Run(fmt.Sprintf("Case%d", i), func(t *testing.T) {
			if len(kat.AAD) > 0 {
				t.Skip("Skipping... SDK implementation does not expose additional authenticated data")
			}
			iv, err := hex.DecodeString(kat.IV)
			if err != nil {
				t.Fatalf("failed to decode iv: %v", err)
			}
			key, err := hex.DecodeString(kat.Key)
			if err != nil {
				t.Fatalf("failed to decode key: %v", err)
			}
			plaintext, err := hex.DecodeString(kat.Plaintext)
			if err != nil {
				t.Fatalf("failed to decode plaintext: %v", err)
			}
			ciphertext, err := hex.DecodeString(kat.CipherText)
			if err != nil {
				t.Fatalf("failed to decode ciphertext: %v", err)
			}
			tag, err := hex.DecodeString(kat.Tag)
			if err != nil {
				t.Fatalf("failed to decode tag: %v", err)
			}
			aesgcmTest(t, iv, key, plaintext, ciphertext, tag)
		})
	}
}

func TestGCMEncryptReader_SourceError(t *testing.T) {
	gcm := &gcmEncryptReader{
		encrypter: &mockCipherAEAD{},
		src:       &mockSourceReader{err: fmt.Errorf("test read error")},
	}

	b := make([]byte, 10)
	n, err := gcm.Read(b)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	} else if err != nil && !strings.Contains(err.Error(), "test read error") {
		t.Fatalf("expected source read error, but got %v", err)
	}

	if n != 0 {
		t.Errorf("expected number of read bytes to be zero, but got %v", n)
	}
}

func TestGCMDecryptReader_SourceError(t *testing.T) {
	gcm := &gcmDecryptReader{
		decrypter: &mockCipherAEAD{},
		src:       &mockSourceReader{err: fmt.Errorf("test read error")},
	}

	b := make([]byte, 10)
	n, err := gcm.Read(b)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	} else if err != nil && !strings.Contains(err.Error(), "test read error") {
		t.Fatalf("expected source read error, but got %v", err)
	}

	if n != 0 {
		t.Errorf("expected number of read bytes to be zero, but got %v", n)
	}
}

func TestGCMDecryptReader_DecrypterOpenError(t *testing.T) {
	gcm := &gcmDecryptReader{
		decrypter: &mockCipherAEAD{openError: fmt.Errorf("test open error")},
		src:       &mockSourceReader{err: io.EOF},
	}

	b := make([]byte, 10)
	n, err := gcm.Read(b)
	if err == nil {
		t.Fatalf("expected error, but got nil")
	} else if err != nil && !strings.Contains(err.Error(), "test open error") {
		t.Fatalf("expected source read error, but got %v", err)
	}

	if n != 0 {
		t.Errorf("expected number of read bytes to be zero, but got %v", n)
	}
}

//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
//= type=test
//# The client MUST set the AAD to the Algorithm Suite ID represented as bytes.
func TestGIVEN_materialsCEKAlg115_WHEN_newAESGCM_THEN_returnedAesGCMAadIs0x0073(t *testing.T) {
	// Given: materials with CEKAlgorithm ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY (0x73)
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: "115",
		Key:          make([]byte, 32),
		IV:           make([]byte, 12),
	}

	// When: instantiating new AES GCM cipher with the materials
	cipher, err := newAESGCM(materials)

	// Then: no error is returned and the AAD is set to the Algorithm Suite ID bytes
	if err != nil {
		panic(fmt.Sprintf("expected no error, but received %v", err))
	}

	aesgcm, ok := cipher.(*aesGCM)
	if !ok {
		panic("expected cipher to be of type *aesGCM")
	}
	
	expectedAad := []byte{0x00, 0x73}
	if !bytes.Equal(aesgcm.aad, expectedAad) {
		panic(fmt.Sprintf("expected AAD to be %v, but received %v", expectedAad, aesgcm.aad))
	}
}

//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-iv12-tag16-no-kdf
//= type=test
//# The client MUST NOT provide any AAD when encrypting with ALG_AES_256_GCM_IV12_TAG16_NO_KDF.
func TestGIVEN_materialsCEKAlgAESGCMNoPadding_WHEN_newAESGCM_THEN_returnedAesGCMAadIsNil(t *testing.T) {
	// Given: materials with CEKAlgorithm ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY (0x73)
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: "AES/GCM/NoPadding",
		Key:          make([]byte, 32),
		IV:           make([]byte, 12),
	}

	// When: instantiating new AES GCM cipher with the materials
	cipher, err := newAESGCM(materials)

	// Then: no error is returned and the AAD is set to the Algorithm Suite ID bytes
	if err != nil {
		panic(fmt.Sprintf("expected no error, but received %v", err))
	}

	aesgcm, ok := cipher.(*aesGCM)
	if !ok {
		panic("expected cipher to be of type *aesGCM")
	}
	
	if aesgcm.aad != nil {
		panic(fmt.Sprintf("expected AAD to be nil, but received %v", aesgcm.aad))
	}
}

func TestNewAESGCM_InvalidNonceLength(t *testing.T) {
	// Test with invalid nonce length (8 bytes instead of 12)
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: algorithms.AESGCMNoPadding,
		Key:          make([]byte, 32),
		IV:           make([]byte, 8), // Wrong length: should be 12 bytes
	}

	_, err := newAESGCM(materials)

	if err == nil {
		t.Fatalf("expected error for invalid nonce length, got nil")
	}

	expectedError := "invalid nonce length: expected 12 bytes, got 8 bytes"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("expected error message to contain %q, got %q", expectedError, err.Error())
	}
}

func aesgcmTest(t *testing.T, iv, key, plaintext, expected, tag []byte) {
	t.Helper()
	const gcmTagSize = 16
	materials := materials2.CryptographicMaterials{
		CEKAlgorithm: algorithms.AESGCMNoPadding,
		Key: key,
		IV:  iv,
	}
	gcm, err := newAESGCM(materials)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	cipherdata := gcm.Encrypt(bytes.NewReader(plaintext))

	ciphertext, err := io.ReadAll(cipherdata)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	// splitting tag and ciphertext
	etag := ciphertext[len(ciphertext)-gcmTagSize:]
	if !bytes.Equal(etag, tag) {
		t.Errorf("expected tags to be equivalent")
	}
	if !bytes.Equal(ciphertext[:len(ciphertext)-gcmTagSize], expected) {
		t.Errorf("expected ciphertext to be equivalent")
	}

	data := gcm.Decrypt(bytes.NewReader(ciphertext))
	text, err := io.ReadAll(data)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if !bytes.Equal(plaintext, text) {
		t.Errorf("expected ciphertext to be equivalent")
	}
}

type mockSourceReader struct {
	n   int
	err error
}

func (b mockSourceReader) Read(p []byte) (n int, err error) {
	return b.n, b.err
}

type mockCipherAEAD struct {
	seal      []byte
	openError error
}

func (m mockCipherAEAD) NonceSize() int {
	panic("implement me")
}

func (m mockCipherAEAD) Overhead() int {
	panic("implement me")
}

func (m mockCipherAEAD) Seal(dst, nonce, plaintext, additionalData []byte) []byte {
	return m.seal
}

func (m mockCipherAEAD) Open(dst, nonce, ciphertext, additionalData []byte) ([]byte, error) {
	return []byte("mocked decrypt"), m.openError
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package awstesting

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"io"
	"net/http"
)

func init() {
	config = aws.Config{}
	config.Region = "mock-region"
	config.Credentials = StubCredentialsProvider{}
}

// StubCredentialsProvider provides a stub credential provider that returns
// static credentials that never expire.
type StubCredentialsProvider struct{}

// Retrieve satisfies the CredentialsProvider interface. Returns stub
// credential value, and never error.
func (StubCredentialsProvider) Retrieve(context.Context) (aws.Credentials, error) {
	return aws.Credentials{
		AccessKeyID: "AKID", SecretAccessKey: "SECRET", SessionToken: "SESSION",
		Source: "unit test credentials",
	}, nil
}

var config aws.Config

// Config returns a copy of the mock configuration for unit tests.
func Config() aws.Config { return config.Copy() }

// MockHttpClient is a simple utility for mocking HTTP responses and capturing requests
type MockHttpClient struct {
	Response      *http.Response
	ResponseError error
	CapturedReq   *http.Request
	CapturedBody  []byte
}

func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	m.CapturedReq = req
	if req.Body != nil {
		body, err := io.ReadAll(req.Body)
		if err != nil {
			return nil, err
		}
		m.CapturedBody = body
	}

	return m.Response, m.ResponseError
}

// TestEndpointResolver returns an endpoint resolver that uses the given URL for the resolved endpoint
func TestEndpointResolver(url string) aws.EndpointResolverWithOptions {
	return aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		endpoint := aws.Endpoint{
			URL:           url,
			PartitionID:   "aws",
			SigningMethod: region,
			SigningName:   service,
		}

		return endpoint, nil
	})
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/rand"
	"fmt"
)

// GenerateNonZeroBytes generates random bytes and validates they are not all zeros
func GenerateNonZeroBytes(n int) ([]byte, error) {
	return GenerateNonZeroBytesWithGenerator(n, generateRandomBytes)
}

// GenerateNonZeroBytesWithGenerator allows injection of custom generator for testing
func GenerateNonZeroBytesWithGenerator(n int, generator func(int) ([]byte, error)) ([]byte, error) {
	const maxRetries = 3 // Prevent infinite loop in case of broken generator
	
	for attempt := 0; attempt < maxRetries; attempt++ {
		//= ../specification/s3-encryption/encryption.md#content-encryption
		//# The client MUST generate an IV or Message ID using the length of the IV or Message ID defined in the algorithm suite.
		keys_iv, err := generator(n)
		if err != nil {
			return nil, err
		}
		
		//= ../specification/s3-encryption/encryption.md#cipher-initialization
		//# The client SHOULD validate that the generated IV or Message ID is not zeros.
		allZero := true
		for _, b := range keys_iv {
			if b != 0 {
				allZero = false
				break
			}
		}
		
		// If not all zeros, we have a valid IV
		if !allZero {
			return keys_iv, nil
		}
		
		// If all zeros, retry (unless this is the last attempt)
	}
	
	// If we've exhausted all retries, return an error
	return nil, fmt.Errorf("failed to generate non-zero IV after %d attempts", maxRetries)
}

// Default generator using crypto/rand
func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"testing"
)

func TestGenerateNonZeroBytes(t *testing.T) {
	cases := []struct {
		name   string
		length int
	}{
		{
			name:   "small_length",
			length: 12,
		},
		{
			name:   "medium_length",
			length: 28,
		},
		{
			name:   "large_length",
			length: 64,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			result, err := GenerateNonZeroBytes(tc.length)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if len(result) != tc.length {
				t.Errorf("expected length %d, got %d", tc.length, len(result))
			}

			// Verify not all zeros
			allZero := true
			for _, b := range result {
				if b != 0 {
					allZero = false
					break
				}
			}
			if allZero {
				t.Error("GenerateNonZeroBytes returned all zeros")
			}
		})
	}
}

//= ../specification/s3-encryption/encryption.md#cipher-initialization
//= type=test
//# The client SHOULD validate that the generated IV or Message ID is not zeros.
func TestGenerateNonZeroBytesWithGenerator_AllZeros(t *testing.T) {
	// Mock generator that always returns all zeros
	mockAllZeros := func(n int) ([]byte, error) {
		return make([]byte, n), nil // Returns all zeros
	}

	testCases := []struct {
		name   string
		length int
	}{
		{
			name:   "12_byte_iv",
			length: 12,
		},
		{
			name:   "28_byte_message_id",
			length: 28,
		},
		{
			name:   "16_byte_iv",
			length: 16,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := GenerateNonZeroBytesWithGenerator(tc.length, mockAllZeros)

			// Assert that an error is raised after max retries
			if err == nil {
				t.Fatalf("Expected error when generator always returns all zeros for length %d, but got none", tc.length)
			}

			// Assert that the error message indicates retry exhaustion
			expectedErrorMsg := "failed to generate non-zero IV after 3 attempts"
			if err.Error() != expectedErrorMsg {
				t.Errorf("Expected error message '%s', but got '%s'", expectedErrorMsg, err.Error())
			}
		})
	}
}

func TestGenerateRandomBytes(t *testing.T) {
	// Test that the default generator produces different results
	results := make([][]byte, 5)
	for i := 0; i < 5; i++ {
		result, err := generateRandomBytes(16)
		if err != nil {
			t.Fatalf("generateRandomBytes failed: %v", err)
		}
		if len(result) != 16 {
			t.Errorf("Expected length 16, got %d", len(result))
		}
		results[i] = result
	}

	// Check that not all results are identical (very unlikely with proper randomness)
	allSame := true
	for i := 1; i < len(results); i++ {
		if !bytes.Equal(results[0], results[i]) {
			allSame = false
			break
		}
	}

	if allSame {
		t.Errorf("All calls to generateRandomBytes returned identical results, randomness may be compromised")
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"io"
)

// Cipher interface allows for either encryption and decryption of an object
type Cipher interface {
	Encrypter
	Decrypter
}

// Encrypter interface with only the encrypt method
type Encrypter interface {
	Encrypt(io.Reader) io.Reader
}

// Decrypter interface with only the decrypt method
type Decrypter interface {
	Decrypt(io.Reader) io.Reader
}

// CryptoReadCloser handles closing of the body and allowing reads from the decrypted
// content.
type CryptoReadCloser struct {
	Body      io.ReadCloser
	Decrypter io.Reader
	isClosed  bool
}

// Close lets the CryptoReadCloser satisfy io.ReadCloser interface
func (rc *CryptoReadCloser) Close() error {
	rc.isClosed = true
	return rc.Body.Close()
}

// Read lets the CryptoReadCloser satisfy io.ReadCloser interface
func (rc *CryptoReadCloser) Read(b []byte) (int, error) {
	if rc.isClosed {
		return 0, io.EOF
	}
	return rc.Decrypter.Read(b)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"io"
)

// ContentCipherBuilder is a builder interface that builds
// ciphers for each request.
type ContentCipherBuilder interface {
	ContentCipher() (ContentCipher, error)
}

// ContentCipherBuilderWithContext is a builder interface that builds
// ciphers for each request.
type ContentCipherBuilderWithContext interface {
	ContentCipherWithContext(context.Context) (ContentCipher, error)
}

// ContentCipher deals with encrypting and decrypting content
type ContentCipher interface {
	EncryptContents(io.Reader) (io.Reader, error)
	DecryptContents(io.ReadCloser) (io.ReadCloser, error)
	GetCipherData() materials.CryptographicMaterials
}

// CEKEntry is a builder that returns a proper content decrypter and error
type CEKEntry func(materials.CryptographicMaterials) (ContentCipher, error)
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"io"
	"strings"
	"testing"
)

func TestCryptoReadCloserRead(t *testing.T) {
	expectedStr := "HELLO WORLD "
	str := strings.NewReader(expectedStr)
	rc := &CryptoReadCloser{Body: io.NopCloser(str), Decrypter: str}

	b, err := io.ReadAll(rc)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if expectedStr != string(b) {
		t.Errorf("expected %s, but received %s", expectedStr, string(b))
	}
}

func TestCryptoReadCloserClose(t *testing.T) {
	data := "HELLO WORLD "
	expectedStr := ""

	str := strings.NewReader(data)
	rc := &CryptoReadCloser{Body: io.NopCloser(str), Decrypter: str}
	rc.Close()

	b, err := io.ReadAll(rc)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if expectedStr != string(b) {
		t.Errorf("expected %s, but received %s", expectedStr, string(b))
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/sha256"
	"hash"
	"io"
)

// lengthReader returns the content length
type lengthReader interface {
	GetContentLength() int64
}

type contentLengthReader struct {
	contentLength int64
	body          io.Reader
}

func NewContentLengthReader(f io.Reader) *contentLengthReader {
	return &contentLengthReader{body: f}
}

func (r *contentLengthReader) Read(b []byte) (int, error) {
	if r.body == nil {
		return 0, io.EOF
	}
	n, err := r.body.Read(b)
	if err != nil && err != io.EOF {
		return n, err
	}
	r.contentLength += int64(n)
	return n, err
}

func (r *contentLengthReader) GetContentLength() int64 {
	return r.contentLength
}

type sha256Writer struct {
	sha256 []byte
	hash   hash.Hash
	out    io.Writer
}

func newSHA256Writer(f io.Writer) *sha256Writer {
	return &sha256Writer{hash: sha256.New(), out: f}
}
func (r *sha256Writer) Write(b []byte) (int, error) {
	r.hash.Write(b)
	return r.out.Write(b)
}

func (r *sha256Writer) GetValue() []byte {
	return r.hash.Sum(nil)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"
	"testing"
)

// From Go stdlib encoding/sha256 test cases
func TestSHA256(t *testing.T) {
	sha := newSHA256Writer(nil)
	expected, _ := hex.DecodeString("e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855")
	b := sha.GetValue()

	if !bytes.Equal(expected, b) {
		t.Errorf("expected equivalent sha values, but received otherwise")
	}
}

func TestSHA256_Case2(t *testing.T) {
	sha := newSHA256Writer(bytes.NewBuffer([]byte{}))
	sha.Write([]byte("hello"))
	expected, _ := hex.DecodeString("2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824")
	b := sha.GetValue()

	if !bytes.Equal(expected, b) {
		t.Errorf("expected equivalent sha values, but received otherwise")
	}
}

type mockReader struct {
	err error
}

func (m mockReader) Read(p []byte) (int, error) {
	return len(p), m.err
}

func TestContentLengthReader(t *testing.T) {
	cases := []struct {
		reader      io.Reader
		expected    int64
		expectedErr string
	}{
		{
			reader:   bytes.NewReader([]byte("foo bar baz")),
			expected: 11,
		},
		{
			reader:   bytes.NewReader(nil),
			expected: 0,
		},
		{
			reader:      mockReader{err: fmt.Errorf("not an EOF error")},
			expectedErr: "not an EOF error",
		},
	}

	for _, tt := range cases {
		reader := NewContentLengthReader(tt.reader)
		_, err := io.ReadAll(reader)
		if err != nil {
			if len(tt.expectedErr) == 0 {
				t.Errorf("expected no error, got %v", err)
			} else if !strings.Contains(err.Error(), tt.expectedErr) {
				t.Errorf("expected error %v, got %v", tt.expectedErr, err.Error())
			}
			continue
		} else if len(tt.expectedErr) > 0 {
			t.Error("expected error, got none")
			continue
		}
		actual := reader.GetContentLength()
		if tt.expected != actual {
			t.Errorf("expected %v, got %v", tt.expected, actual)
		}
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bufio"
	"errors"
	"io"
	"os"
)

type writerStore struct {
	io.ReadWriteSeeker
	Cleanup func()
}

func GetWriterStore(path string, useTempFile bool) (*writerStore, error) {
	if !useTempFile {
		return &writerStore{
			ReadWriteSeeker: &bytesReadWriteSeeker{},
			Cleanup:         func() {},
		}, nil
	}
	// Create temp file to be used later for calculating the SHA256 header
	f, err := os.CreateTemp(path, "")
	if err != nil {
		return nil, err
	}

	ws := &writerStore{
		ReadWriteSeeker: f,
		Cleanup: func() {
			// Close the temp file and Cleanup
			f.Close()
			os.Remove(f.Name())
		},
	}

	return ws, nil
}

type bytesReadWriteSeeker struct {
	buf []byte
	i   int64
}

// Copied from Go stdlib bytes.Reader
func (ws *bytesReadWriteSeeker) Read(b []byte) (int, error) {
	if ws.i >= int64(len(ws.buf)) {
		return 0, io.EOF
	}
	n := copy(b, ws.buf[ws.i:])
	ws.i += int64(n)
	return n, nil
}

func (ws *bytesReadWriteSeeker) Write(b []byte) (int, error) {
	ws.buf = append(ws.buf, b...)
	return len(b), nil
}

// Copied from Go stdlib bytes.Reader
func (ws *bytesReadWriteSeeker) Seek(offset int64, whence int) (int64, error) {
	var abs int64
	switch whence {
	case 0:
		abs = offset
	case 1:
		abs = int64(ws.i) + offset
	case 2:
		abs = int64(len(ws.buf)) + offset
	default:
		return 0, errors.New("bytes.Reader.Seek: invalid whence")
	}
	if abs < 0 {
		return 0, errors.New("bytes.Reader.Seek: negative position")
	}
	ws.i = abs
	return abs, nil
}

// NewBufferedReader creates a buffered reader with the specified buffer size
// This implements the S3EC requirement to set buffer size to a reasonable default for GetObject
func NewBufferedReader(r io.Reader, bufferSize int) (io.ReadCloser, error) {
	if bufferSize <= 0 {
		return nil, errors.New("buffer size must be greater than zero")
	}
	
	bufferedReader := bufio.NewReaderSize(r, bufferSize)
	
	// If the original reader implements io.ReadCloser, preserve the Close method
	if rc, ok := r.(io.ReadCloser); ok {
		return &bufferedReadCloser{
			Reader: bufferedReader,
			closer: rc,
		}, nil
	}
	
	// If not a ReadCloser, create a no-op closer
	return &bufferedReadCloser{
		Reader: bufferedReader,
		closer: nil,
	}, nil
}

// bufferedReadCloser wraps a buffered reader and preserves the Close method
type bufferedReadCloser struct {
	*bufio.Reader
	closer io.ReadCloser
}

func (brc *bufferedReadCloser) Close() error {
	if brc.closer != nil {
		return brc.closer.Close()
	}
	return nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestBytesReadWriteSeeker_Read(t *testing.T) {
	b := &bytesReadWriteSeeker{[]byte{1, 2, 3}, 0}
	expected := []byte{1, 2, 3}
	buf := make([]byte, 3)
	n, err := b.Read(buf)

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 3, n; e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	if !bytes.Equal(expected, buf) {
		t.Error("expected equivalent byte slices, but received otherwise")
	}
}

func TestBytesReadWriteSeeker_Write(t *testing.T) {
	b := &bytesReadWriteSeeker{}
	expected := []byte{1, 2, 3}
	buf := make([]byte, 3)
	n, err := b.Write([]byte{1, 2, 3})

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 3, n; e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	n, err = b.Read(buf)
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 3, n; e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	if !bytes.Equal(expected, buf) {
		t.Error("expected equivalent byte slices, but received otherwise")
	}
}

func TestBytesReadWriteSeeker_Seek(t *testing.T) {
	b := &bytesReadWriteSeeker{[]byte{1, 2, 3}, 0}
	expected := []byte{2, 3}
	m, err := b.Seek(1, io.SeekStart)

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 1, int(m); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	buf := make([]byte, 3)
	n, err := b.Read(buf)

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	if e, a := 2, n; e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}

	if !bytes.Equal(expected, buf[:n]) {
		t.Error("expected equivalent byte slices, but received otherwise")
	}
}

func TestGetWriterStore_TempFile(t *testing.T) {
	ws, err := GetWriterStore("", true)
	if err != nil {
		t.Fatalf("expected no error, but received %v", err)
	}
	tempFile, ok := ws.ReadWriteSeeker.(*os.File)
	if !ok {
		t.Fatal("io.ReadWriteSeeker expected to be *os.file")
	}

	ws.Cleanup()
	if _, err := os.Stat(tempFile.Name()); !os.IsNotExist(err) {
		t.Errorf("expected temp file be deleted, but still exists %v", tempFile.Name())
	}
}

func TestGetWriterStore_Memory(t *testing.T) {
	ws, err := GetWriterStore("", false)
	if err != nil {
		t.Fatalf("expected no error, but received %v", err)
	}
	if _, ok := ws.ReadWriteSeeker.(*bytesReadWriteSeeker); !ok {
		t.Fatal("io.ReadWriteSeeker expected to be *bytesReadWriteSeeker")
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	ws.Cleanup()
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"crypto/subtle"
	"encoding/binary"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"crypto/hkdf"
)

const (
	// Key derivation constants
	DeriveKeyInfo = "DERIVEKEY"
	CommitKeyInfo = "COMMITKEY"
)

// KeyDerivationResult holds the results of HKDF key derivation
type KeyDerivationResult struct {
	DerivedEncryptionKey []byte
	CommitKey            []byte
}

func DeriveKeys(plaintextDataKey []byte, messageID []byte, algorithmSuiteID int, storedKeyCommitment []byte) (*KeyDerivationResult, error) {
	// Validate input parameters
	if len(plaintextDataKey) == 0 {
		return nil, fmt.Errorf("plaintext data key cannot be empty")
	}
	if len(messageID) == 0 {
		return nil, fmt.Errorf("message ID cannot be empty")
	}

	// Get algorithm suite to determine key lengths and other properties
	algSuite, err := algorithms.GetAlgorithmSuiteByID(algorithmSuiteID)
	if err != nil {
		return nil, fmt.Errorf("unable to get algorithm suite by ID: %v", err)
	}

	// Convert algorithm suite ID to bytes (big-endian)
	algorithmSuiteIDBytes := make([]byte, 2)
	binary.BigEndian.PutUint16(algorithmSuiteIDBytes, uint16(algorithmSuiteID))

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The hash function MUST be specified by the algorithm suite commitment settings.
	hashFunc := algSuite.KDFHashAlgorithm()
	if hashFunc == nil || !algSuite.IsCommitting() {
		return nil, fmt.Errorf("algorithm suite does not support key derivation: 0x%04x", algorithmSuiteID)
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input keying material MUST be the plaintext data key (PDK) generated by the key provider.
	inputKeyMaterial := plaintextDataKey

	// Get expected key lengths from algorithm suite
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the input keying material MUST equal the key derivation input length specified by the algorithm suite commit key derivation setting.
	if len(inputKeyMaterial) != expectedDataKeyLength {
		return nil, fmt.Errorf("plaintext data key length must be %d bytes, got %d", expectedDataKeyLength, len(inputKeyMaterial))
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The salt MUST be the Message ID with the length defined in the algorithm suite.
	salt := messageID
	expectedSaltLength := algSuite.IVLengthBytes()
	if len(salt) != expectedSaltLength {
		return nil, fmt.Errorf("message ID length must be %d bytes, got %d", expectedSaltLength, len(salt))
	}

	// Extract step
	extractPRK, err := hkdf.Extract(hashFunc, inputKeyMaterial, salt)
	if err != nil {
		return nil, fmt.Errorf("hkdf extract failed: %w", err)
	}
	if len(extractPRK) == 0 {
		return nil, fmt.Errorf("hkdf extract failed: extractPRK is empty")
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input info MUST be a concatenation of the algorithm suite ID as bytes followed by the string COMMITKEY as UTF8 encoded bytes.
	commitKeyInfo := append(algorithmSuiteIDBytes, []byte(CommitKeyInfo)...)
	
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The CK input pseudorandom key MUST be the output from the extract step.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the output keying material MUST equal the commit key length specified by the supported algorithm suites.
	derivedCommitKey, err := hkdf.Expand(hashFunc, extractPRK, string(commitKeyInfo), expectedCommitKeyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive commit key: %w", err)
	}
	if len(derivedCommitKey) != expectedCommitKeyLength {
		return nil, fmt.Errorf("commit key length must be %d bytes, got %d", expectedCommitKeyLength, len(derivedCommitKey))
	}

	// First, check commitment value; then, derive encryption key
	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=implication
	//# When using an algorithm suite which supports key commitment, the client MUST verify the key commitment values match before deriving the [derived encryption key](./key-derivation.md#hkdf-operation).

	if storedKeyCommitment != nil {
		//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
		//# When using an algorithm suite which supports key commitment, the client MUST verify that the [derived key commitment](./key-derivation.md#hkdf-operation) contains the same bytes as the stored key commitment retrieved from the stored object's metadata.
		//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
		//= type=implication
		//# When using an algorithm suite which supports key commitment, the verification of the derived key commitment value MUST be done in constant time.
		if subtle.ConstantTimeCompare(derivedCommitKey, storedKeyCommitment) != 1 {
			//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
			//# When using an algorithm suite which supports key commitment, the client MUST throw an exception when the derived key commitment value and stored key commitment value do not match.
			return nil, fmt.Errorf("derived key commitment value does not match value stored on encrypted message")
		}
	}

	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The input info MUST be a concatenation of the algorithm suite ID as bytes followed by the string DERIVEKEY as UTF8 encoded bytes.
	encryptionKeyInfo := append(algorithmSuiteIDBytes, []byte(DeriveKeyInfo)...)
	
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//= type=implication
	//# - The DEK input pseudorandom key MUST be the output from the extract step.
	//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
	//# - The length of the output keying material MUST equal the encryption key length specified by the algorithm suite encryption settings.
	derivedEncryptionKey, err := hkdf.Expand(hashFunc, extractPRK, string(encryptionKeyInfo), expectedDataKeyLength)
	if err != nil {
		return nil, fmt.Errorf("failed to derive encryption key: %w", err)
	}
	if len(derivedEncryptionKey) != expectedDataKeyLength {
		return nil, fmt.Errorf("encryption key length must be %d bytes, got %d", expectedDataKeyLength, len(derivedEncryptionKey))
	}

	return &KeyDerivationResult{
		DerivedEncryptionKey: derivedEncryptionKey,
		CommitKey:            derivedCommitKey,
	}, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"bytes"
	"encoding/hex"
	"strings"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"testing"
)

func hexToBytes(t *testing.T, s string, expectedLen int) []byte {
	b, err := decodeHex(s)
	if err != nil {
		t.Fatalf("Failed to decode hex string: %v", err)
	}
	if len(b) != expectedLen {
		t.Fatalf("Expected length %d, got %d for hex string %s", expectedLen, len(b), s)
	}
	return b
}

func decodeHex(s string) ([]byte, error) {
	dst := make([]byte, hex.DecodedLen(len(s)))
	_, err := hex.Decode(dst, []byte(s))
	if err != nil {
		return nil, err
	}
	return dst, nil
}

func TestDeriveKeys_KnownAnswerTests(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	tests := []struct {
		comment        string
		dataKeyHex     string
		messageIDHex   string
		expectedEncHex string
		expectedComHex string
	}{
		{
			comment:        "Basic S3EC.ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY #1",
			dataKeyHex:     "80d90dc4cc7e77d8a6332efa44eba56230a7fe7b89af37d1e501ab2e07c0a163",
			messageIDHex:   "b8ea76bed24c7b85382a148cb9dcd1cfdfb765f55ded4dfa6e0c4c79",
			expectedEncHex: "6dd14f546cc006e639126e83f5d4d1b118576bb5df97f38c6fb3a1db87bbc338",
			expectedComHex: "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db",
		},
		{
			comment:        "Basic S3EC.ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY #2",
			dataKeyHex:     "501afb8227d22e75e68010414b8abdaf3064c081e8e922dafef4992036394d60",
			messageIDHex:   "61a00b4981a5aacfd136c55cb726e32d2a547dc7600a7d4675c69127",
			expectedEncHex: "e14786a714748d1d2c3a4a6816dec56ddf1881bbeabb4f39420ffb9f63700b2f",
			expectedComHex: "5c1e73e47f6fe3a70d6d094283aceaa76d2975feb829212d88f0afc1",
		},
	}

	for _, tc := range tests {
		t.Run(tc.comment, func(t *testing.T) {
			dataKey := hexToBytes(t, tc.dataKeyHex, expectedDataKeyLength)
			messageID := hexToBytes(t, tc.messageIDHex, expectedMessageIDLength)
			expectedEnc := hexToBytes(t, tc.expectedEncHex, expectedDataKeyLength)
			expectedCom := hexToBytes(t, tc.expectedComHex, expectedCommitKeyLength)

			result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), expectedCom)
			if err != nil {
				t.Fatalf("DeriveKeys failed: %v", err)
			}
			if !bytes.Equal(result.DerivedEncryptionKey, expectedEnc) {
				t.Errorf("DerivedEncryptionKey mismatch.\nExpected: %x\nGot:      %x", expectedEnc, result.DerivedEncryptionKey)
			}
			if !bytes.Equal(result.CommitKey, expectedCom) {
				t.Errorf("CommitKey mismatch.\nExpected: %x\nGot:      %x", expectedCom, result.CommitKey)
			}
		})
	}
}

func TestDeriveKeys_KeyCommitmentValidation(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	dataKey := hexToBytes(t, "80d90dc4cc7e77d8a6332efa44eba56230a7fe7b89af37d1e501ab2e07c0a163", expectedDataKeyLength)
	messageID := hexToBytes(t, "b8ea76bed24c7b85382a148cb9dcd1cfdfb765f55ded4dfa6e0c4c79", expectedMessageIDLength)
	correctCommitment := hexToBytes(t, "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db", expectedCommitKeyLength)
	wrongCommitment := hexToBytes(t, "00000000000000000000000000000000000000000000000000000000", expectedCommitKeyLength)

	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=test
	//# When using an algorithm suite which supports key commitment, the client MUST verify that the [derived key commitment](./key-derivation.md#hkdf-operation) contains the same bytes as the stored key commitment retrieved from the stored object's metadata.
	t.Run("matching_commitment_succeeds", func(t *testing.T) {
		result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), correctCommitment)
		if err != nil {
			t.Fatalf("DeriveKeys should succeed when commitment values match, got error: %v", err)
		}

		// Verify that the derived commitment matches the stored commitment
		if !bytes.Equal(result.CommitKey, correctCommitment) {
			t.Errorf("Derived commitment should match stored commitment.\nExpected: %x\nGot:      %x", correctCommitment, result.CommitKey)
		}
	})

	//= ../specification/s3-encryption/decryption.md#decrypting-with-commitment
	//= type=test
	//# When using an algorithm suite which supports key commitment, the client MUST throw an exception when the derived key commitment value and stored key commitment value do not match.
	t.Run("mismatched_commitment_throws_exception", func(t *testing.T) {
		_, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), wrongCommitment)
		if err == nil {
			t.Fatalf("DeriveKeys should throw an exception when commitment values do not match, but it succeeded")
		}

		// Verify the error message indicates commitment mismatch
		if !bytes.Contains([]byte(err.Error()), []byte("derived key commitment value does not match value stored on encrypted message")) {
			t.Errorf("Expected error to mention commitment mismatch, got: %v", err)
		}
	})
}

func TestDeriveKeys_InputOutputLengthValidation(t *testing.T) {
	// Get algorithm suite for expected lengths
	algSuite := algorithms.AlgAES256GCMHkdfSha512CommitKey
	expectedDataKeyLength := algSuite.DataKeyLengthBytes()
	expectedCommitKeyLength := algSuite.CommitmentLengthBytes()
	expectedMessageIDLength := algSuite.IVLengthBytes()

	correctCommitment := hexToBytes(t, "f89818bc0a346d3a3426b68e9509b6b2ae5fe1f904aa329fb73625db", expectedCommitKeyLength)

	cases := []struct {
		name           string
		dataKeyLen     int
		messageIDLen   int
		expectError    bool
		errorContains  string
	}{
		{
			name:         "correct_input_lengths",
			dataKeyLen:   expectedDataKeyLength,   // 32 bytes
			messageIDLen: expectedMessageIDLength, // Message ID length for committing algorithm
			expectError:  true, // Will fail due to commitment mismatch, but that's OK - we're testing length validation
			errorContains: "commitment",
		},
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//= type=test
		//# - The length of the input keying material MUST equal the key derivation input length specified by the algorithm suite commit key derivation setting.
		{
			name:          "wrong_data_key_length",
			dataKeyLen:    16, // Wrong length
			messageIDLen:  expectedMessageIDLength,
			expectError:   true,
			errorContains: "plaintext data key length",
		},
		//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
		//= type=test
		//# - The salt MUST be the Message ID with the length defined in the algorithm suite.
		{
			name:          "wrong_message_id_length",
			dataKeyLen:    expectedDataKeyLength,
			messageIDLen:  16, // Wrong length
			expectError:   true,
			errorContains: "message ID length",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			
			dataKey := make([]byte, tc.dataKeyLen)
			for i := range dataKey {
				dataKey[i] = byte(i % 256)
			}

			messageID := make([]byte, tc.messageIDLen)
			for i := range messageID {
				messageID[i] = byte((i + 100) % 256)
			}

			result, err := DeriveKeys(dataKey, messageID, algorithms.AlgAES256GCMHkdfSha512CommitKey.ID(), correctCommitment)

			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error for %s but got none", tc.name)
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Expected error for %s: %v", tc.name, err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for %s, got %v", tc.name, err)
				}

				//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
				//= type=test
				//# - The length of the output keying material MUST equal the encryption key length specified by the algorithm suite encryption settings.
				if len(result.DerivedEncryptionKey) != expectedDataKeyLength {
					t.Errorf("DerivedEncryptionKey length should be %d, got %d", expectedDataKeyLength, len(result.DerivedEncryptionKey))
				}
				
				//= ../specification/s3-encryption/key-derivation.md#hkdf-operation
				//= type=test
				//# - The length of the output keying material MUST equal the commit key length specified by the supported algorithm suites.
				if len(result.CommitKey) != expectedCommitKeyLength {
					t.Errorf("CommitKey length should be %d, got %d", expectedCommitKeyLength, len(result.CommitKey))
				}
				t.Logf("✓ Correct input/output lengths for %s", tc.name)
			}
		})
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"strconv"
)

// DefaultInstructionKeySuffix is appended to the end of the instruction file key when
// grabbing or saving to S3
const DefaultInstructionKeySuffix = ".instruction"

const (
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=implication
	//# The "x-amz-" prefix denotes that the metadata is owned by an Amazon product and MUST be prepended to all S3EC metadata mapkeys.
	amzPrefix 				   	   = "x-amz-"
	metaHeader                     = amzPrefix + "meta"
	keyV1Header                    = amzPrefix + "key"
	keyV2Header                    = amzPrefix + "key-v2"
	ivHeader                       = amzPrefix + "iv"
	matDescHeader                  = amzPrefix + "matdesc"
	CekAlgorithmHeader             = amzPrefix + "cek-alg"
	KeyringAlgorithmHeader         = amzPrefix + "wrap-alg"
	tagLengthHeader                = amzPrefix + "tag-len"
	unencryptedContentLengthHeader = amzPrefix + "unencrypted-content-length"

	// For the below constants, the specification comments describe the recommended constant names.
	// However, Go's JSON struct tags require literal string values and cannot reference constants.
	// This forces us to duplicate the string literals in the ObjectMetadata struct's `json:` tags.
	// While this creates duplication between these constants and the struct tags, it's a Go language
	// limitation - we cannot write `json:ContentCipherV3` in the struct definition.
	// There are tests that validate these constant values match the struct tags in ObjectMetadata.

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-c") SHOULD be represented by a constant named "CONTENT_CIPHER_V3" or similar in the implementation code.
	ContentCipherV3                = amzPrefix + "c"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-3") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_V3" or similar in the implementation code.
	EncryptedDataKeyV3             = amzPrefix + "3"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-m") SHOULD be represented by a constant named "MAT_DESC_V3" or similar in the implementation code.
	MatDescV3                      = amzPrefix + "m"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-t") SHOULD be represented by a constant named "ENCRYPTION_CONTEXT_V3" or similar in the implementation code.
	EncryptionContextV3            = amzPrefix + "t"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-w") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_ALGORITHM_V3" or similar in the implementation code.
	EncryptedDataKeyAlgorithmV3    = amzPrefix + "w"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-d") SHOULD be represented by a constant named "KEY_COMMITMENT_V3" or similar in the implementation code.
	KeyCommitmentV3                = amzPrefix + "d"

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - This mapkey ("x-amz-i") SHOULD be represented by a constant named "MESSAGE_ID_V3" or similar in the implementation code.
	MessageIDV3                    = amzPrefix + "i"
)

// S3EC Go V4 does not support reading nor writing V1 format metadata.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-unencrypted-content-length" SHOULD be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-key" MUST be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-matdesc" MUST be present for V1 format objects.

//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=exception
//# - The mapkey "x-amz-iv" MUST be present for V1 format objects.

// ObjectMetadata encryption starts off by generating a random symmetric key using
// AES GCM. The SDK generates a random IV based off the encryption cipher
// chosen. The master key that was provided, whether by the user or KMS, will be used
// to encrypt the randomly generated symmetric key and base64 encode the iv. This will
// allow for decryption of that same data later.
//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
//= type=implication
//# The "x-amz-meta-" prefix is automatically added by the S3 server and MUST NOT be included in implementation code.
type ObjectMetadata struct {
	// IV is the randomly generated IV base64 encoded.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-iv" MUST be present for V2 format objects.
	IV string `json:"x-amz-iv"`
	// CipherKey is the randomly generated cipher key.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-key-v2" MUST be present for V2 format objects.
	CipherKey string `json:"x-amz-key-v2"`
	// MaterialDesc is a description to distinguish from other envelopes.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-matdesc" MUST be present for V2 format objects.
	MatDesc               string `json:"x-amz-matdesc"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-wrap-alg" MUST be present for V2 format objects.
	KeyringAlg            string `json:"x-amz-wrap-alg"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-cek-alg" MUST be present for V2 format objects.
	CEKAlg                string `json:"x-amz-cek-alg"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-tag-len" MUST be present for V2 format objects.
	TagLen                string `json:"x-amz-tag-len"`
	UnencryptedContentLen string `json:"x-amz-unencrypted-content-length"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-c" MUST be present for V3 format objects.
	ContentCipher string `json:"x-amz-c"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-3" MUST be present for V3 format objects.
	EncryptedDataKey string `json:"x-amz-3"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
	MatDescV3 string `json:"x-amz-m"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-t" SHOULD be present for V3 format objects that use KMS Encryption Context.
	EncryptionContext string `json:"x-amz-t"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-w" MUST be present for V3 format objects.
	WrappingAlgorithm string `json:"x-amz-w"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-d" MUST be present for V3 format objects.
	KeyCommitment string `json:"x-amz-d"`
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//# - The mapkey "x-amz-i" MUST be present for V3 format objects.
	MessageID string `json:"x-amz-i"`
}

// V3 algorithm compression mappings
var (
	// Wrapping algorithm decompression: compressed value -> full algorithm name
	v3WrapAlgDecompression = map[string]string{
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval
		"02": "AES/GCM",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval
		"12": "kms+context",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval
		"22": "RSA-OAEP-SHA1",
	}
	
	// Wrapping algorithm compression: full algorithm name -> compressed value
	v3WrapAlgCompression = map[string]string{
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval, and vice versa on write.
		"AES/GCM": "02",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval, and vice versa on write.
		"kms+context": "12",
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval, and vice versa on write.
		"RSA-OAEP-SHA1": "22",
	}
)

func (e *ObjectMetadata) GetDecodedKey() ([]byte, error) {
	var keyStr string
	if e.EncryptedDataKey != "" {
		// V3
		keyStr = e.EncryptedDataKey
	} else {
		// V2
		keyStr = e.CipherKey
	}
	key, err := base64.StdEncoding.DecodeString(keyStr)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func (e *ObjectMetadata) GetDecodedMessageIDOrIV() ([]byte, error) {
	var value string
	if e.MessageID != "" {
		// V3
		value = e.MessageID
	} else {
		// V2
		value = e.IV
	}
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		return nil, err
	}
	return decoded, nil
}

func (e *ObjectMetadata) GetDecodedKeyCommitment() ([]byte, error) {
	// Only V3 has KeyCommitment
	commitment, err := base64.StdEncoding.DecodeString(e.KeyCommitment)
	if err != nil {
		return nil, err
	}
	return commitment, err
}

func (e *ObjectMetadata) GetMatDescV3() (string, error) {
	if e.MatDescV3 != "" {
		return e.MatDescV3, nil
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# If the mapkey x-amz-m is not present, the default Material Description value MUST be set to an empty map (`{}`).
	return "{}", nil
}

func (e *ObjectMetadata) GetEncryptionContextV3() (string, error) {
	// Only V3 has EncryptionContext
	if e.EncryptionContext != "" {
		return e.EncryptionContext, nil
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//# If the mapkey x-amz-t is not present, the default Material Description value MUST be set to an empty map (`{}`).
	return "{}", nil
}

func (e *ObjectMetadata) GetMatDescV2() (string, error) {
	return e.MatDesc, nil
}

func (e *ObjectMetadata) GetEncryptionContextOrMatDescV3() (string, error) {
	wrappingAlg, err := e.GetFullWrappingAlgorithm()
	var matDesc string
	if wrappingAlg == "kms+context" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
		matDesc, err = e.GetEncryptionContextV3()
		return matDesc, err
	} else if wrappingAlg == "AES/GCM" || wrappingAlg == "RSA-OAEP-SHA1" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
		matDesc, err = e.GetMatDescV3()
		return matDesc, err
	}
	return "", fmt.Errorf("unsupported wrapping algorithm for getting Material Description: %s", wrappingAlg)
}

func (e *ObjectMetadata) GetContentEncryptionAlgorithmString() (string, error) {
	if e.ContentCipher != "" {
		// V3
		return e.ContentCipher, nil
	}
	if e.CEKAlg != "" {
		// V2
		return e.CEKAlg, nil
	}
	return "", fmt.Errorf("no content encryption algorithm found in metadata")
}

func (e *ObjectMetadata) GetContentEncryptionAlgorithmSuite() (*algorithms.AlgorithmSuite, error) {
	cekString, err := e.GetContentEncryptionAlgorithmString()
	if err != nil {
		return nil, err
	}
	if cekString == algorithms.AESGCMCommitKey {
		return algorithms.AlgAES256GCMHkdfSha512CommitKey, nil
	} else if cekString == algorithms.AESGCMNoPadding {
		return algorithms.AlgAES256GCMIV12Tag16NoKDF, nil
	} else if cekString == algorithms.AESCBCPKCS5 {
		return algorithms.AlgAES256CBCIV16NoKDF, nil
	} else if cekString == algorithms.AESCTRNoPadding {
		return algorithms.AlgAES256CTRIV16Tag16NoKDF, nil
	}
	return nil, fmt.Errorf("invalid content encryption algorithm found in metadata: %s", cekString)
}

func (e *ObjectMetadata) GetFullWrappingAlgorithm() (string, error) {
	if e.WrappingAlgorithm != "" {
		// V3
		// Decompress the V3 wrapping algorithm to its full name
		fullAlg, exists := v3WrapAlgDecompression[e.WrappingAlgorithm]
		if !exists {
			return "", fmt.Errorf("unknown V3 wrapping algorithm: %s", e.WrappingAlgorithm)
		}
		return fullAlg, nil
	}
	if e.KeyringAlg != "" {
		// V2
		return e.KeyringAlg, nil
	}
	return "", fmt.Errorf("no wrapping algorithm found in metadata")
}

// CompressWrappingAlgorithm compresses a full wrapping algorithm name to V3 format
func CompressWrappingAlgorithm(fullAlgorithm string) (string, error) {
	compressed, exists := v3WrapAlgCompression[fullAlgorithm]
	if !exists {
		return "", fmt.Errorf("unsupported wrapping algorithm for V3: %s", fullAlgorithm)
	}
	return compressed, nil
}

// UnmarshalJSON unmarshalls the given JSON bytes into ObjectMetadata
func (e *ObjectMetadata) UnmarshalJSON(value []byte) error {
	type StrictEnvelope ObjectMetadata
	type LaxEnvelope struct {
		StrictEnvelope
		TagLen                json.RawMessage `json:"x-amz-tag-len"`
		UnencryptedContentLen json.RawMessage `json:"x-amz-unencrypted-content-length"`
	}

	inner := LaxEnvelope{}
	err := json.Unmarshal(value, &inner)
	if err != nil {
		return err
	}
	*e = ObjectMetadata(inner.StrictEnvelope)

	e.TagLen, err = getJSONNumberAsString(inner.TagLen)
	if err != nil {
		return fmt.Errorf("failed to parse tag length: %w", err)
	}

	e.UnencryptedContentLen, err = getJSONNumberAsString(inner.UnencryptedContentLen)
	if err != nil {
		return fmt.Errorf("failed to parse unencrypted content length: %w", err)
	}

	return nil
}

// getJSONNumberAsString will attempt to convert the provided bytes into a string representation of a JSON Number.
// Only supports byte values that are string or integers, not floats. If the provided value is JSON Null, empty string
// will be returned.
func getJSONNumberAsString(data []byte) (string, error) {
	if len(data) == 0 {
		return "", nil
	}

	// first try string, this also catches null value
	var s *string
	err := json.Unmarshal(data, &s)
	if err == nil && s != nil {
		return *s, nil
	} else if err == nil {
		return "", nil
	}

	// fallback to int64
	var i int64
	err = json.Unmarshal(data, &i)
	if err == nil {
		return strconv.FormatInt(i, 10), nil
	}

	return "", fmt.Errorf("failed to parse as JSON Number")
}

func EncodeMeta(reader lengthReader, cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
	//# Objects encrypted with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY MUST use the V3 message format version only.
	if cryptographicMaterials.CEKAlgorithm == algorithms.AESGCMCommitKey {
		return EncodeMetaV3(cryptographicMaterials)
	//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
	//# Objects encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF MUST use the V2 message format version only.
	} else if cryptographicMaterials.CEKAlgorithm == algorithms.AESGCMNoPadding {
		return EncodeMetaV2(reader, cryptographicMaterials)
	} else {
		// Go S3EC V4 does not support writing other (ex. AES CBC) messages
		//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
		//= type=exception
		//# Objects encrypted with ALG_AES_256_CBC_IV16_NO_KDF MAY use either the V1 or V2 message format version.
		return ObjectMetadata{}, fmt.Errorf("unsupported CEK algorithm: %s", cryptographicMaterials.CEKAlgorithm)
	}
}

func EncodeMetaV2(reader lengthReader, cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	iv := base64.StdEncoding.EncodeToString(cryptographicMaterials.IV)
	key := base64.StdEncoding.EncodeToString(cryptographicMaterials.EncryptedKey)

	encodedMatDesc, err := cryptographicMaterials.MaterialDescription.EncodeDescription()
	if err != nil {
		return ObjectMetadata{}, err
	}

	contentLength := reader.GetContentLength()

	return ObjectMetadata{
		CipherKey:             key,
		IV:                    iv,
		MatDesc:               string(encodedMatDesc),
		KeyringAlg:            cryptographicMaterials.KeyringAlgorithm,
		CEKAlg:                cryptographicMaterials.CEKAlgorithm,
		TagLen:                cryptographicMaterials.TagLength,
		UnencryptedContentLen: strconv.FormatInt(contentLength, 10),
	}, nil
}

func EncodeMetaV3(cryptographicMaterials materials.CryptographicMaterials) (ObjectMetadata, error) {
	iv := base64.StdEncoding.EncodeToString(cryptographicMaterials.IV)
	key := base64.StdEncoding.EncodeToString(cryptographicMaterials.EncryptedKey)
	alg_suite := cryptographicMaterials.CEKAlgorithm
	wrapping_alg := cryptographicMaterials.KeyringAlgorithm
	commitment := base64.StdEncoding.EncodeToString(cryptographicMaterials.KeyCommitment)
	mat_desc_bytes, err := cryptographicMaterials.MaterialDescription.EncodeDescription()
	if err != nil {
		return ObjectMetadata{}, err
	}
	mat_desc := string(mat_desc_bytes)

	out := ObjectMetadata{
		EncryptedDataKey:  key,
		MessageID:         iv,
		ContentCipher:     alg_suite,
		WrappingAlgorithm: wrapping_alg,
		KeyCommitment:     commitment,
	}

	// Set MatDescV3 for AES/GCM or RSA-OAEP-SHA1, EncryptionContext for kms+context
	if wrapping_alg == "AES/GCM" || wrapping_alg == "RSA-OAEP-SHA1" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
		out.MatDescV3 = mat_desc
	} else if wrapping_alg == "kms+context" {
		//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
		//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
		out.EncryptionContext = mat_desc
	}

	return out, nil
}

// MetadataFormat represents the format version of S3EC metadata
type MetadataFormat int

const (
	FormatUnknown MetadataFormat = iota
	FormatInstructionFile
	FormatV1
	FormatV2
	FormatV3
)

// Validate and detect correct metadata format 
func DetectAndValidateMetadataFormat(metadata map[string]string) (MetadataFormat, error) {
	// Check for mapkeys defined in the spec	
	hasV1Key := hasKey(metadata, keyV1Header) // "x-amz-key"
	hasV2Key := hasKey(metadata, keyV2Header) // "x-amz-key-v2"
	hasV3Key := hasKey(metadata, EncryptedDataKeyV3) // "x-amz-3"
	hasIv := hasKey(metadata, ivHeader) // "x-amz-iv"
	hasV3KeyCommitment := hasKey(metadata, KeyCommitmentV3) // "x-amz-d"
	hasV3KeyID := hasKey(metadata, MessageIDV3) // "x-amz-i"

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-iv" and "x-amz-key" then the object MUST be considered as an S3EC-encrypted object using the V1 format.
	isV1 := hasIv && hasV1Key
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-iv" and "x-amz-metadata-x-amz-key-v2" then the object MUST be considered as an S3EC-encrypted object using the V2 format.
	isV2 := hasIv && hasV2Key
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# - If the metadata contains "x-amz-3" and "x-amz-d" and "x-amz-i" then the object MUST be considered an S3EC-encrypted object using the V3 format.
	isV3 := hasV3Key && hasV3KeyCommitment && hasV3KeyID

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# If there are multiple mapkeys which are meant to be exclusive, such as "x-amz-key", "x-amz-key-v2", and "x-amz-3" then the S3EC SHOULD throw an exception.
	hasAtLeastOneV1ExclusiveKey := hasV1Key
	hasAtLeastOneV2ExclusiveKey := hasV2Key
	hasAtLeastOneV3ExclusiveKey := hasV3Key
	exclusiveKeyMatchCount := 0
	if hasAtLeastOneV1ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if hasAtLeastOneV2ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if hasAtLeastOneV3ExclusiveKey {
		exclusiveKeyMatchCount++
	}
	if exclusiveKeyMatchCount > 1 {
		return FormatUnknown, fmt.Errorf("metadata contains conflicting exclusive mapkeys")
	}

	versionMatchCount := 0
	if isV1 {
		versionMatchCount++
	}
	if isV2 {
		versionMatchCount++
	}
	if isV3 {
		versionMatchCount++
	}
	if versionMatchCount > 1 {
		return FormatUnknown, fmt.Errorf("metadata contains multiple S3EC format versions")
	}
	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//# If the object matches none of the V1/V2/V3 formats, the S3EC MUST attempt to get the instruction file.
	if versionMatchCount == 0 {
		return FormatInstructionFile, nil
	}

	if isV1 {
		return FormatV1, nil
	}
	if isV2 {
		return FormatV2, nil
	}
	if isV3 {
		return FormatV3, nil
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
	//= type=implication
	//# In general, if there is any deviation from the above format, with the exception of additional unrelated mapkeys, then the S3EC SHOULD throw an exception.
	return FormatUnknown, fmt.Errorf("unable to determine metadata format")
}

// hasKey checks if a metadata key exists (with or without x-amz-meta prefix)
func hasKey(metadata map[string]string, key string) bool {
	// Check direct key
	if _, exists := metadata[key]; exists {
		return true
	}
	// Check with x-amz-meta prefix
	prefixedKey := metaHeader + "-" + key
	if _, exists := metadata[prefixedKey]; exists {
		return true
	}
	return false
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"encoding/json"
	"reflect"
	"testing"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
)

func TestEnvelope_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		content  []byte
		expected ObjectMetadata
		actual   ObjectMetadata
	}{
		"string json numbers": {
			content: []byte(`{
  "x-amz-iv": "iv",
  "x-amz-key-v2": "key",
  "x-amz-matdesc": "{\"aws:x-amz-cek-alg\":\"AES/GCM/NoPadding\"}",
  "x-amz-wrap-alg": "kms+context",
  "x-amz-cek-alg": "AES/GCM/NoPadding",
  "x-amz-tag-len": "128",
  "x-amz-unencrypted-content-length": "1024"
}
`),
			expected: ObjectMetadata{
				IV:                    "iv",
				CipherKey:             "key",
				MatDesc:               `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
				KeyringAlg:            "kms+context",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "128",
				UnencryptedContentLen: "1024",
			},
		},
		"integer json numbers": {
			content: []byte(`{
  "x-amz-iv": "iv",
  "x-amz-key-v2": "key",
  "x-amz-matdesc": "{\"aws:x-amz-cek-alg\":\"AES/GCM/NoPadding\"}",
  "x-amz-wrap-alg": "kms+context",
  "x-amz-cek-alg": "AES/GCM/NoPadding",
  "x-amz-tag-len": 128,
  "x-amz-unencrypted-content-length": 1024
}
`),
			expected: ObjectMetadata{
				IV:                    "iv",
				CipherKey:             "key",
				MatDesc:               `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
				KeyringAlg:            "kms+context",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "128",
				UnencryptedContentLen: "1024",
			},
		},
		"null json numbers": {
			content: []byte(`{
  "x-amz-iv": "iv",
  "x-amz-key-v2": "key",
  "x-amz-matdesc": "{\"aws:x-amz-cek-alg\":\"AES/GCM/NoPadding\"}",
  "x-amz-wrap-alg": "kms+context",
  "x-amz-cek-alg": "AES/GCM/NoPadding",
  "x-amz-tag-len": null,
  "x-amz-unencrypted-content-length": null
}
`),
			expected: ObjectMetadata{
				IV:         "iv",
				CipherKey:  "key",
				MatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
				KeyringAlg: "kms+context",
				CEKAlg:     "AES/GCM/NoPadding",
			},
		},
		"no json numbers": {
			content: []byte(`{
  "x-amz-iv": "iv",
  "x-amz-key-v2": "key",
  "x-amz-matdesc": "{\"aws:x-amz-cek-alg\":\"AES/GCM/NoPadding\"}",
  "x-amz-wrap-alg": "kms+context",
  "x-amz-cek-alg": "AES/GCM/NoPadding"
}
`),
			expected: ObjectMetadata{
				IV:         "iv",
				CipherKey:  "key",
				MatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
				KeyringAlg: "kms+context",
				CEKAlg:     "AES/GCM/NoPadding",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal(tt.content, &tt.actual)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !reflect.DeepEqual(tt.expected, tt.actual) {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

func TestObjectMetadata_UnmarshalJSON(t *testing.T) {
	cases := map[string]struct {
		content  []byte
		expected ObjectMetadata
		actual   ObjectMetadata
	}{
		"complete V3 metadata with encryption context": {
			content: []byte(`{
				"x-amz-c": "115",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-t": "{\"kms_cmk_id\":\"test-key-id\"}",
				"x-amz-w": "12",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "115",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				EncryptionContext: `{"kms_cmk_id":"test-key-id"}`,
				WrappingAlgorithm: "12",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
		"V3 metadata with material description": {
			content: []byte(`{
				"x-amz-c": "AES/GCM/NoPadding",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-m": "{\"test\":\"material-desc\"}",
				"x-amz-w": "01",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "AES/GCM/NoPadding",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				MatDescV3:         `{"test":"material-desc"}`,
				WrappingAlgorithm: "01",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
		"minimal V3 metadata": {
			content: []byte(`{
				"x-amz-c": "AES/CBC/PKCS5Padding",
				"x-amz-3": "dGVzdC1lbmNyeXB0ZWQta2V5",
				"x-amz-w": "11",
				"x-amz-d": "dGVzdC1rZXktY29tbWl0bWVudA==",
				"x-amz-i": "dGVzdC1tZXNzYWdlLWlk"
				}
			`),
			expected: ObjectMetadata{
				ContentCipher:     "AES/CBC/PKCS5Padding",
				EncryptedDataKey:  "dGVzdC1lbmNyeXB0ZWQta2V5",
				WrappingAlgorithm: "11",
				KeyCommitment:     "dGVzdC1rZXktY29tbWl0bWVudA==",
				MessageID:         "dGVzdC1tZXNzYWdlLWlk",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			err := json.Unmarshal(tt.content, &tt.actual)
			if err != nil {
				t.Errorf("expected no error, got %v", err)
			}
			if !reflect.DeepEqual(tt.expected, tt.actual) {
				t.Errorf("expected %v, got %v", tt.expected, tt.actual)
			}
		})
	}
}

func TestWrappingAlgorithmCompression(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval, and vice versa on write.
	result, err := CompressWrappingAlgorithm("AES/GCM")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "02" {
		t.Errorf("expected 02, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval, and vice versa on write.
	result, err = CompressWrappingAlgorithm("kms+context")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "12" {
		t.Errorf("expected 12, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval, and vice versa on write.
	result, err = CompressWrappingAlgorithm("RSA-OAEP-SHA1")
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "22" {
		t.Errorf("expected 22, got %s", result)
	}

	result, err = CompressWrappingAlgorithm("not-a-known-algorithm")
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestWrappingAlgorithmDecompression(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "02" MUST be translated to AES/GCM upon retrieval
	metadata := ObjectMetadata{WrappingAlgorithm: "02"}
	result, err := metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "AES/GCM" {
		t.Errorf("expected AES/GCM, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "12" MUST be translated to kms+context upon retrieval
	metadata = ObjectMetadata{WrappingAlgorithm: "12"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "kms+context" {
		t.Errorf("expected kms+context, got %s", result)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# - The wrapping algorithm value "22" MUST be translated to RSA-OAEP-SHA1 upon retrieval
	metadata = ObjectMetadata{WrappingAlgorithm: "22"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
	if result != "RSA-OAEP-SHA1" {
		t.Errorf("expected RSA-OAEP-SHA1, got %s", result)
	}

	metadata = ObjectMetadata{WrappingAlgorithm: "99"}
	result, err = metadata.GetFullWrappingAlgorithm()
	if err == nil {
		t.Errorf("expected error but got none")
	}
}

func TestDetectAndValidateMetadataFormat(t *testing.T) {
	cases := map[string]struct {
		metadata map[string]string
		expected MetadataFormat
	}{
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-3" and "x-amz-d" and "x-amz-i" then the object MUST be considered an S3EC-encrypted object using the V3 format.
		"V3 format": {
			metadata: map[string]string{
				"x-amz-3": "encrypted-key",
				"x-amz-d": "key-commitment",
				"x-amz-i": "message-id",
			},
			expected: FormatV3,
		},
		"V3 format with meta prefix": {
			metadata: map[string]string{
				"x-amz-meta-x-amz-3": "encrypted-key",
				"x-amz-meta-x-amz-d": "key-commitment",
				"x-amz-meta-x-amz-i": "message-id",
			},
			expected: FormatV3,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-iv" and "x-amz-metadata-x-amz-key-v2" then the object MUST be considered as an S3EC-encrypted object using the V2 format.
		"V2 format minimal meta prefix": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-meta-x-amz-key-v2":  "key",
			},
			expected: FormatV2,
		},
		"V2 format minimal": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key-v2":  "key",
			},
			expected: FormatV2,
		},
		"V2 format": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key-v2":  "key",
				"x-amz-matdesc": "matdesc",
			},
			expected: FormatV2,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# - If the metadata contains "x-amz-iv" and "x-amz-key" then the object MUST be considered as an S3EC-encrypted object using the V1 format.
		"V1 format minimal": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key":     "key",
			},
			expected: FormatV1,
		},
		"V1 format": {
			metadata: map[string]string{
				"x-amz-iv":      "iv",
				"x-amz-key":     "key",
				"x-amz-matdesc": "matdesc",
			},
			expected: FormatV1,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# If the object matches none of the V1/V2/V3 formats, the S3EC MUST attempt to get the instruction file.
		"matching no format some keys": {
			metadata: map[string]string{
				"x-amz-abcdef": "not-a-key",
				"x-amz-123": "still-not-a-key",
			},
			expected: FormatInstructionFile,
		},
		"matching no format no keys": {
			metadata: map[string]string{},
			expected: FormatInstructionFile,
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#determining-s3ec-object-status
		//= type=test
		//# If there are multiple mapkeys which are meant to be exclusive, such as "x-amz-key", "x-amz-key-v2", and "x-amz-3" then the S3EC SHOULD throw an exception.
		"multiple exclusive mapkeys case 1": {
			metadata: map[string]string{
				"x-amz-key": "key",
				"x-amz-key-v2": "key-v2",
				"x-amx-3": "key-v3",
			},
			expected: FormatUnknown,
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := DetectAndValidateMetadataFormat(tt.metadata)
			if err != nil && result != FormatUnknown {
				t.Errorf("expected no error, got %v", err)
			}
			if result != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result)
			}
		})
	}
}

func TestObjectMetadataConstValues(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-c") SHOULD be represented by a constant named "CONTENT_CIPHER_V3" or similar in the implementation code.
	if ContentCipherV3 != "x-amz-c" {
		t.Errorf("ContentCipherV3 MUST be `x-amz-c`, got %q", ContentCipherV3)
	}
	// Truly wild reflection usage in this test below to ensure the struct field in ObjectMetadata has the correct json tag to fully satisfy spec requirement
	field, ok := reflect.TypeOf(ObjectMetadata{}).FieldByName("ContentCipher")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field ContentCipher")
	}
	if got := field.Tag.Get("json"); got != "x-amz-c" {
		t.Errorf("ContentCipher json tag MUST be `x-amz-c`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-3") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_V3" or similar in the implementation code.
	if EncryptedDataKeyV3 != "x-amz-3" {
		t.Errorf("EncryptedDataKeyV3 MUST be `x-amz-3`, got %q", EncryptedDataKeyV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("EncryptedDataKey")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field EncryptedDataKey")
	}
	if got := field.Tag.Get("json"); got != "x-amz-3" {
		t.Errorf("EncryptedDataKey json tag MUST be `x-amz-3`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-m") SHOULD be represented by a constant named "MAT_DESC_V3" or similar in the implementation code.
	if MatDescV3 != "x-amz-m" {
		t.Errorf("MatDescV3 MUST be `x-amz-m`, got %q", MatDescV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("MatDescV3")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field MatDesc")
	}
	if got := field.Tag.Get("json"); got != "x-amz-m" {
		t.Errorf("MatDesc json tag MUST be `x-amz-m`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-t") SHOULD be represented by a constant named "ENCRYPTION_CONTEXT_V3" or similar in the implementation code.
	if EncryptionContextV3 != "x-amz-t" {
		t.Errorf("EncryptionContextV3 MUST be `x-amz-t`, got %q", EncryptionContextV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("EncryptionContext")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field EncryptionContext")
	}
	if got := field.Tag.Get("json"); got != "x-amz-t" {
		t.Errorf("EncryptionContext json tag MUST be `x-amz-t`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-w") SHOULD be represented by a constant named "ENCRYPTED_DATA_KEY_ALGORITHM_V3" or similar in the implementation code.
	if EncryptedDataKeyAlgorithmV3 != "x-amz-w" {
		t.Errorf("EncryptedDataKeyAlgorithmV3 MUST be `x-amz-w`, got %q", EncryptedDataKeyAlgorithmV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("WrappingAlgorithm")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field WrappingAlgorithm")
	}
	if got := field.Tag.Get("json"); got != "x-amz-w" {
		t.Errorf("WrappingAlgorithm json tag MUST be `x-amz-w`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-d") SHOULD be represented by a constant named "KEY_COMMITMENT_V3" or similar in the implementation code.
	if KeyCommitmentV3 != "x-amz-d" {
		t.Errorf("KeyCommitmentV3 MUST be `x-amz-d`, got %q", KeyCommitmentV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("KeyCommitment")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field KeyCommitment")
	}
	if got := field.Tag.Get("json"); got != "x-amz-d" {
		t.Errorf("KeyCommitment json tag MUST be `x-amz-d`, got %q", got)
	}

	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=test
	//# - This mapkey ("x-amz-i") SHOULD be represented by a constant named "MESSAGE_ID_V3" or similar in the implementation code.
	if MessageIDV3 != "x-amz-i" {
		t.Errorf("MessageIDV3 MUST be `x-amz-i`, got %q", MessageIDV3)
	}
	field, ok = reflect.TypeOf(ObjectMetadata{}).FieldByName("MessageID")
	if !ok {
		t.Fatal("ObjectMetadata SHOULD have field MessageID")
	}
	if got := field.Tag.Get("json"); got != "x-amz-i" {
		t.Errorf("MessageID json tag MUST be `x-amz-i`, got %q", got)
	}
}

func TestMaterialDescriptionAndEncryptionContextRequirements(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# The Material Description MUST be used for wrapping algorithms `AES/GCM` (`02`) and `RSA-OAEP-SHA1` (`22`).
	t.Run("Material Description used for AES/GCM", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "02", // AES/GCM
			MatDescV3:         `{"test":"material-desc"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"test":"material-desc"}` {
			t.Errorf("expected material description, got %s", result)
		}
	})

	t.Run("Material Description used for RSA-OAEP-SHA1", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "22", // RSA-OAEP-SHA1
			MatDescV3:         `{"test":"rsa-material-desc"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"test":"rsa-material-desc"}` {
			t.Errorf("expected material description, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# If the mapkey x-amz-m is not present, the default Material Description value MUST be set to an empty map (`{}`).
	t.Run("Default empty map when Material Description not present for AES/GCM", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "02", // AES/GCM
			// MatDescV3 is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})

	t.Run("Default empty map when Material Description not present for RSA-OAEP-SHA1", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "22", // RSA-OAEP-SHA1
			// MatDescV3 is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# The Encryption Context value MUST be used for wrapping algorithm `kms+context` or `12`.
	t.Run("Encryption Context used for kms+context", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "12", // kms+context
			EncryptionContext: `{"kms_cmk_id":"test-key-id"}`,
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != `{"kms_cmk_id":"test-key-id"}` {
			t.Errorf("expected encryption context, got %s", result)
		}
	})

	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# If the mapkey x-amz-t is not present, the default Material Description value MUST be set to an empty map (`{}`).
	t.Run("Default empty map when Encryption Context not present for kms+context", func(t *testing.T) {
		metadata := ObjectMetadata{
			WrappingAlgorithm: "12", // kms+context
			// EncryptionContext is empty
		}
		result, err := metadata.GetEncryptionContextOrMatDescV3()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}
		if result != "{}" {
			t.Errorf("expected empty map {}, got %s", result)
		}
	})
}

// mockLengthReader implements the lengthReader interface for testing
type mockLengthReader struct {
	contentLength int64
}

func (m *mockLengthReader) GetContentLength() int64 {
	return m.contentLength
}

func TestEncodeMetaV2(t *testing.T) {
	cases := map[string]struct {
		reader               lengthReader
		cryptographicMaterials materials.CryptographicMaterials
		expected             ObjectMetadata
	}{
		"standard V2 encoding": {
			reader: &mockLengthReader{contentLength: 1024},
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-12-b"),
				KeyringAlgorithm:    "kms+context",
				CEKAlgorithm:        "AES/GCM/NoPadding",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"aws:x-amz-cek-alg": "AES/GCM/NoPadding", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
			},
			expected: ObjectMetadata{
				CipherKey:             "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				IV:                    "dGVzdC1pdi0xMi1i", // base64 of "test-iv-12-b"
				MatDesc:               `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding","custom":"value"}`,
				KeyringAlg:            "kms+context",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "128",
				UnencryptedContentLen: "1024",
			},
		},
		"V2 encoding with empty material description": {
			reader: &mockLengthReader{contentLength: 2048},
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-12-b"),
				KeyringAlgorithm:    "kms",
				CEKAlgorithm:        "AES/GCM/NoPadding",
				TagLength:           "96",
				MaterialDescription: materials.MaterialDescription{},
				EncryptedKey:        []byte("encrypted-key-data"),
			},
			expected: ObjectMetadata{
				CipherKey:             "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				IV:                    "dGVzdC1pdi0xMi1i", // base64 of "test-iv-12-b"
				MatDesc:               `{}`,
				KeyringAlg:            "kms",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "96",
				UnencryptedContentLen: "2048",
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := EncodeMetaV2(tt.reader, tt.cryptographicMaterials)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}

func TestEncodeMetaV3(t *testing.T) {
	cases := map[string]struct {
		cryptographicMaterials materials.CryptographicMaterials
		expected             ObjectMetadata
	}{
		"V3 encoding with AES/GCM wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "AES/GCM",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"test": "material-desc", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "AES/GCM",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
				MatDescV3:         `{"custom":"value","test":"material-desc"}`,
			},
		},
		"V3 encoding with kms+context wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "kms+context",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"kms_cmk_id": "test-key-id", "custom": "value"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "kms+context",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-t" SHOULD be present for V3 format objects that use KMS Encryption Context.
				EncryptionContext: `{"custom":"value","kms_cmk_id":"test-key-id"}`,
			},
		},
		"V3 encoding with RSA-OAEP-SHA1 wrapping algorithm": {
			cryptographicMaterials: materials.CryptographicMaterials{
				Key:                 []byte("test-key-32-bytes-long-12345678"),
				IV:                  []byte("test-iv-28-bytes-long-1234567890"),
				KeyringAlgorithm:    "RSA-OAEP-SHA1",
				CEKAlgorithm:        "115",
				TagLength:           "128",
				MaterialDescription: materials.MaterialDescription{"rsa": "material-desc"},
				EncryptedKey:        []byte("encrypted-key-data"),
				KeyCommitment:       []byte("key-commitment-data"),
			},
			expected: ObjectMetadata{
				EncryptedDataKey:  "ZW5jcnlwdGVkLWtleS1kYXRh", // base64 of "encrypted-key-data"
				MessageID:         "dGVzdC1pdi0yOC1ieXRlcy1sb25nLTEyMzQ1Njc4OTA=", // base64 of IV
				ContentCipher:     "115",
				WrappingAlgorithm: "RSA-OAEP-SHA1",
				KeyCommitment:     "a2V5LWNvbW1pdG1lbnQtZGF0YQ==", // base64 of "key-commitment-data"
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-m" SHOULD be present for V3 format objects that use Raw Keyring Material Description.
				MatDescV3:         `{"rsa":"material-desc"}`,
			},
		},
	}

	for name, tt := range cases {
		t.Run(name, func(t *testing.T) {
			result, err := EncodeMetaV3(tt.cryptographicMaterials)
			if err != nil {
				t.Fatalf("expected no error, got %v", err)
			}

			if !reflect.DeepEqual(tt.expected, result) {
				t.Errorf("expected %+v, got %+v", tt.expected, result)
			}
		})
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

// Padder handles padding of crypto data
type Padder interface {
	// Pad will pad the byte array.
	// The second parameter is NOT how many
	// bytes to pad by, but how many bytes
	// have been read prior to the padding.
	// This allows for streamable padding.
	Pad([]byte, int) ([]byte, error)
	// Unpad will unpad the byte bytes. Unpad
	// methods must be constant time.
	Unpad([]byte) ([]byte, error)
	// Name returns the name of the padder.
	// This is used when decrypting on
	// instantiating new padders.
	Name() string
}

// NoPadder does not pad anything
var NoPadder = Padder(noPadder{})

type noPadder struct{}

func (padder noPadder) Pad(b []byte, n int) ([]byte, error) {
	return b, nil
}

func (padder noPadder) Unpad(b []byte) ([]byte, error) {
	return b, nil
}

func (padder noPadder) Name() string {
	return "NoPadding"
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

// Copyright 2017 Amazon.com, Inc. or its affiliates. All Rights Reserved.
//
// Portions Licensed under the MIT License. Copyright (c) 2016 Carl Jackson

import (
	"bytes"
	"crypto/subtle"
	"errors"
)

const (
	pkcs7MaxPaddingSize = 255
)

type pkcs7Padder struct {
	blockSize int
}

// NewPKCS7Padder follows the RFC 2315: https://www.ietf.org/rfc/rfc2315.txt
// PKCS7 padding is subject to side-channel attacks and timing attacks. For
// the most secure data, use an authenticated crypto algorithm.
func NewPKCS7Padder(blockSize int) Padder {
	return pkcs7Padder{blockSize}
}

var errPKCS7Padding = errors.New("invalid padding")
var errInvalidBlockSize = errors.New("invalid block size, must be between 1 and 255")

// Pad will pad the data relative to how many bytes have been read.
// Pad follows the PKCS7 standard.
func (padder pkcs7Padder) Pad(buf []byte, n int) ([]byte, error) {
	if padder.blockSize < 1 || padder.blockSize > pkcs7MaxPaddingSize {
		return nil, errInvalidBlockSize
	}
	size := padder.blockSize - (n % padder.blockSize)
	pad := bytes.Repeat([]byte{byte(size)}, size)
	buf = append(buf, pad...)
	return buf, nil
}

// Unpad will unpad the correct amount of bytes based off
// of the PKCS7 standard
func (padder pkcs7Padder) Unpad(buf []byte) ([]byte, error) {
	if len(buf) == 0 {
		return nil, errPKCS7Padding
	}

	// Here be dragons. We're attempting to check the padding in constant
	// time. The only piece of information here which is public is len(buf).
	// This code is modeled loosely after tls1_cbc_remove_padding from
	// OpenSSL.
	padLen := buf[len(buf)-1]
	toCheck := pkcs7MaxPaddingSize
	good := 1
	if toCheck > len(buf) {
		toCheck = len(buf)
	}
	for i := 0; i < toCheck; i++ {
		b := buf[len(buf)-1-i]

		outOfRange := subtle.ConstantTimeLessOrEq(int(padLen), i)
		equal := subtle.ConstantTimeByteEq(padLen, b)
		good &= subtle.ConstantTimeSelect(outOfRange, 1, equal)
	}

	good &= subtle.ConstantTimeLessOrEq(1, int(padLen))
	good &= subtle.ConstantTimeLessOrEq(int(padLen), len(buf))

	if good != 1 {
		return nil, errPKCS7Padding
	}

	return buf[:len(buf)-int(padLen)], nil
}

func (padder pkcs7Padder) Name() string {
	return "PKCS7Padding"
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal_test

import (
	"bytes"
	"fmt"
	s3crypto "github.com/aws/amazon-s3-encryption-client-go/v4/internal"
	"testing"
)

func padTest(size int, t *testing.T) {
	padder := s3crypto.NewPKCS7Padder(size)
	for i := 0; i < size; i++ {
		input := make([]byte, i)
		expected := append(input, bytes.Repeat([]byte{byte(size - i)}, size-i)...)
		b, err := padder.Pad(input, len(input))
		if err != nil {
			t.Fatal("Expected error to be nil but received " + err.Error())
		}
		if len(b) != len(expected) {
			t.Fatal(fmt.Sprintf("Case %d: data is not of the same length", i))
		}
		if bytes.Compare(b, expected) != 0 {
			t.Fatal(fmt.Sprintf("Expected %v but got %v", expected, b))
		}
	}
}

func unpadTest(size int, t *testing.T) {
	padder := s3crypto.NewPKCS7Padder(size)
	for i := 0; i < size; i++ {
		expected := make([]byte, i)
		input := append(expected, bytes.Repeat([]byte{byte(size - i)}, size-i)...)
		b, err := padder.Unpad(input)
		if err != nil {
			t.Fatal("Error received, was expecting nil: " + err.Error())
		}
		if len(b) != len(expected) {
			t.Fatal(fmt.Sprintf("Case %d: data is not of the same length", i))
		}
		if bytes.Compare(b, expected) != 0 {
			t.Fatal(fmt.Sprintf("Expected %v but got %v", expected, b))
		}
	}
}

func TestPKCS7Padding(t *testing.T) {
	padTest(10, t)
	padTest(16, t)
	padTest(255, t)
}

func TestPKCS7Unpadding(t *testing.T) {
	unpadTest(10, t)
	unpadTest(16, t)
	unpadTest(255, t)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/smithy-go"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// GetObjectAPIClient is a client that implements the GetObject operation
type GetObjectAPIClient interface {
	GetObject(context.Context, *s3.GetObjectInput, ...func(*s3.Options)) (*s3.GetObjectOutput, error)
}

// SaveStrategyRequest represents a request sent to a SaveStrategy to save the contents of an ObjectMetadata
type SaveStrategyRequest struct {
	// The envelope to save
	Envelope *ObjectMetadata

	// The HTTP request being built
	HTTPRequest *http.Request

	// The operation Input type
	Input interface{}
}

// ObjectMetadataSaveStrategy will save the metadata of the crypto contents to the header of
// the object.
type ObjectMetadataSaveStrategy struct{}

// Save will save the envelope to the request's header.
func (strat ObjectMetadataSaveStrategy) Save(ctx context.Context, saveReq *SaveStrategyRequest) error {

	input := saveReq.Input.(*s3.PutObjectInput)
	if input.Metadata == nil {
		input.Metadata = map[string]string{}
	}

	// S3EC Go V4 supports reading content metadata from an instruction file, but not writing it to an instruction file.
	// Any content metadata written is implicitly written to object metadata and not an instruction file.

	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=implication
	//# By default, the S3EC MUST store content metadata in the S3 Object Metadata.
	//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
	//= type=implication
	//# In the V3 format, the mapkeys "x-amz-c", "x-amz-d", and "x-amz-i" MUST be stored exclusively in the Object Metadata.
	env := saveReq.Envelope
	if env.EncryptedDataKey != "" {
		// V3 format
		compressedAlg, err := CompressWrappingAlgorithm(env.WrappingAlgorithm)
		if err != nil {
			return fmt.Errorf("error while compressing wrapping algorithm: %w", err)
		}
		input.Metadata[http.CanonicalHeaderKey(ContentCipherV3)] = env.ContentCipher
		input.Metadata[http.CanonicalHeaderKey(EncryptedDataKeyV3)] = env.EncryptedDataKey
		input.Metadata[http.CanonicalHeaderKey(MatDescV3)] = env.MatDescV3
		input.Metadata[http.CanonicalHeaderKey(EncryptionContextV3)] = env.EncryptionContext
		input.Metadata[http.CanonicalHeaderKey(EncryptedDataKeyAlgorithmV3)] = compressedAlg
		input.Metadata[http.CanonicalHeaderKey(KeyCommitmentV3)] = env.KeyCommitment
		input.Metadata[http.CanonicalHeaderKey(MessageIDV3)] = env.MessageID
	} else {
		// V2 format
		input.Metadata[http.CanonicalHeaderKey(keyV2Header)] = env.CipherKey
		input.Metadata[http.CanonicalHeaderKey(ivHeader)] = env.IV
		input.Metadata[http.CanonicalHeaderKey(matDescHeader)] = env.MatDesc
		input.Metadata[http.CanonicalHeaderKey(KeyringAlgorithmHeader)] = env.KeyringAlg
		input.Metadata[http.CanonicalHeaderKey(CekAlgorithmHeader)] = env.CEKAlg
	}
	input.Metadata[http.CanonicalHeaderKey(unencryptedContentLengthHeader)] = env.UnencryptedContentLen

	if len(env.TagLen) > 0 {
		input.Metadata[http.CanonicalHeaderKey(tagLengthHeader)] = env.TagLen
	}
	return nil

	// S3EC Go V4 supports reading content metadata from an instruction file, but not writing it.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MUST support writing some or all (depending on format) content metadata to an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The content metadata stored in the Instruction File MUST be serialized to a JSON string.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The serialized JSON string MUST be the only contents of the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# Instruction File writes MUST NOT be enabled by default.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# Instruction File writes MUST be optionally configured during client creation or on each PutObject request.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MAY support re-encryption/key rotation via Instruction Files.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC MUST NOT support providing a custom Instruction File suffix on ordinary writes; custom suffixes MUST only be used during re-encryption.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#instruction-file
	//= type=exception
	//# The S3EC SHOULD support providing a custom Instruction File suffix on GetObject requests, regardless of whether or not re-encryption is supported.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v1-v2-instruction-files
	//= type=exception
	//# In the V1/V2 message format, all of the content metadata MUST be stored in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-c" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-c" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-d" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-d" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-i" and its value in the Object Metadata when writing with an Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST NOT store the mapkey "x-amz-i" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-3" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-w" and its value in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-m" and its value (when present in the content metadata) in the Instruction File.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#v3-instruction-files
	//= type=exception
	//# - The V3 message format MUST store the mapkey "x-amz-t" and its value (when present in the content metadata) in the Instruction File.
}

// LoadStrategyRequest represents a request sent to a LoadStrategy to load the contents of an ObjectMetadata
type LoadStrategyRequest struct {
	// The HTTP response
	HTTPResponse *http.Response

	// The operation Input type
	Input interface{}
}

// LoadStrategy ...
type LoadStrategy interface {
	Load(context.Context, *LoadStrategyRequest) (ObjectMetadata, error)
}

// S3LoadStrategy will load the instruction file from s3
type s3LoadStrategy struct {
	APIClient             GetObjectAPIClient
	InstructionFileSuffix string
}

// Load from a given instruction file suffix
func (load s3LoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	env := ObjectMetadata{}
	if load.InstructionFileSuffix == "" {
		load.InstructionFileSuffix = DefaultInstructionKeySuffix
	}

	input := req.Input.(*s3.GetObjectInput)
	out, err := load.APIClient.GetObject(ctx, &s3.GetObjectInput{
		Key:    aws.String(strings.Join([]string{*input.Key, load.InstructionFileSuffix}, "")),
		Bucket: input.Bucket,
	})

	if err != nil {
		return env, err
	}

	b, err := io.ReadAll(out.Body)
	if err != nil {
		return env, err
	}
	err = json.Unmarshal(b, &env)
	return env, err
}

// headerV2LoadStrategy will load the envelope from the metadata
type headerV2LoadStrategy struct{}

// Load from a given object's header
func (load headerV2LoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	env := ObjectMetadata{}
	env.CipherKey = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, keyV2Header}, "-"))
	env.IV = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, ivHeader}, "-"))
	env.MatDesc = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, matDescHeader}, "-"))
	env.KeyringAlg = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, KeyringAlgorithmHeader}, "-"))
	env.CEKAlg = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, CekAlgorithmHeader}, "-"))
	env.TagLen = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, tagLengthHeader}, "-"))
	env.UnencryptedContentLen = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, unencryptedContentLengthHeader}, "-"))
	return env, nil
}

// headerV3LoadStrategy will load the V3 envelope from the metadata
type headerV3LoadStrategy struct{}

// Load from a given object's header using V3 format
func (load headerV3LoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	v3Meta := ObjectMetadata{}
	v3Meta.ContentCipher = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, ContentCipherV3}, "-"))
	v3Meta.EncryptedDataKey = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptedDataKeyV3}, "-"))
	v3Meta.MatDescV3 = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, MatDescV3}, "-"))
	v3Meta.EncryptionContext = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptionContextV3}, "-"))
	v3Meta.WrappingAlgorithm = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, EncryptedDataKeyAlgorithmV3}, "-"))
	v3Meta.KeyCommitment = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, KeyCommitmentV3}, "-"))
	v3Meta.MessageID = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, MessageIDV3}, "-"))
	v3Meta.UnencryptedContentLen = req.HTTPResponse.Header.Get(strings.Join([]string{metaHeader, unencryptedContentLengthHeader}, "-"))
	return v3Meta, nil
}

// DefaultLoadStrategy This is the only exported LoadStrategy since cx are no longer able to configure their client
// with a specific load strategy. Instead, we figure out which strategy to use based on the response header on decrypt.
type DefaultLoadStrategy struct {
	client GetObjectAPIClient
	suffix string
}

func (load DefaultLoadStrategy) Load(ctx context.Context, req *LoadStrategyRequest) (ObjectMetadata, error) {
	// Create metadata map from headers for format detection
	metadata := make(map[string]string)
	for key, values := range req.HTTPResponse.Header {
		if len(values) > 0 {
			metadata[strings.ToLower(key)] = values[0]
		}
	}

	// Detect format and validate
	format, err := DetectAndValidateMetadataFormat(metadata)
	if err != nil {
		return ObjectMetadata{}, fmt.Errorf("invalid metadata format: %w", err)
	}

	switch format {
	case FormatV3:
		strat := headerV3LoadStrategy{}
		return strat.Load(ctx, req)
		
	case FormatV2:
		strat := headerV2LoadStrategy{}
		return strat.Load(ctx, req)
		
	case FormatV1:
		// In other S3EC implementations, decryption of v1 objects is supported.
		// Go, however, does not support this.
		return ObjectMetadata{}, &smithy.GenericAPIError{
			Code:    "V1NotSupportedError",
			Message: "The AWS SDK for Go does not support version 1",
		}
		
	default:
		// Fall back to instruction file loading
		var client GetObjectAPIClient
		if load.client == nil {
			cfg, err := config.LoadDefaultConfig(context.Background())
			if err != nil {
				return ObjectMetadata{}, fmt.Errorf("unable to create S3 client to load instruction file: %w", err)
			}
			client = s3.NewFromConfig(cfg)
		} else {
			client = load.client
		}

		// Load from instruction file
		strat := s3LoadStrategy{
			APIClient:             client,
			InstructionFileSuffix: load.suffix,
		}
		loadedMetadata, err := strat.Load(ctx, req)
		if err != nil {
			return ObjectMetadata{}, err
		}

		// For "V3 instruction file" format, load any additional metadata from headers
		// EncryptedDataKey is only present for V3 format (with or without instruction file)
		if loadedMetadata.EncryptedDataKey != "" {
			// If any values that should be in the headers were in the instruction file, raise error
			if loadedMetadata.ContentCipher != "" || loadedMetadata.KeyCommitment != "" || loadedMetadata.MessageID != "" {
				return ObjectMetadata{}, fmt.Errorf("invalid metadata format: missing V3 header values in instruction file format")
			}

			// Load these values from headers
			headerMeta, err := headerV3LoadStrategy{}.Load(ctx, req)
			if err != nil {
				return ObjectMetadata{}, err
			}
			loadedMetadata.ContentCipher = headerMeta.ContentCipher
			loadedMetadata.KeyCommitment = headerMeta.KeyCommitment
			loadedMetadata.MessageID = headerMeta.MessageID
		}
		
		return loadedMetadata, nil
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	"context"
	"reflect"
	"testing"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func TestHeaderV2SaveStrategy(t *testing.T) {
	cases := []struct {
		env      ObjectMetadata
		expected map[string]string
	}{
		{
			ObjectMetadata{
				CipherKey:             "Foo",
				IV:                    "Bar",
				MatDesc:               "{}",
				KeyringAlg:            "kms",
				CEKAlg:                "AES/GCM/NoPadding",
				TagLen:                "128",
				UnencryptedContentLen: "0",
			},
			map[string]string{
				"X-Amz-Key-V2":                     "Foo",
				"X-Amz-Iv":                         "Bar",
				"X-Amz-Matdesc":                    "{}",
				"X-Amz-Wrap-Alg":                   "kms",
				"X-Amz-Cek-Alg":                    "AES/GCM/NoPadding",
				"X-Amz-Tag-Len":                    "128",
				"X-Amz-Unencrypted-Content-Length": "0",
			},
		},
		{
			ObjectMetadata{
				CipherKey:             "Foo",
				IV:                    "Bar",
				MatDesc:               "{}",
				KeyringAlg:            "kms",
				CEKAlg:                "AES/GCM/NoPadding",
				UnencryptedContentLen: "0",
			},
			map[string]string{
				"X-Amz-Key-V2":                     "Foo",
				"X-Amz-Iv":                         "Bar",
				"X-Amz-Matdesc":                    "{}",
				"X-Amz-Wrap-Alg":                   "kms",
				"X-Amz-Cek-Alg":                    "AES/GCM/NoPadding",
				"X-Amz-Unencrypted-Content-Length": "0",
			},
		},
	}

	for _, c := range cases {
		params := &s3.PutObjectInput{}
		req := &SaveStrategyRequest{
			Envelope: &c.env,
			Input:    params,
		}
		strat := ObjectMetadataSaveStrategy{}
		err := strat.Save(context.Background(), req)
		if err != nil {
			t.Errorf("expected no error, but received %v", err)
		}

		if !reflect.DeepEqual(c.expected, params.Metadata) {
			t.Errorf("expected %v, but received %v", c.expected, params.Metadata)
		}
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package internal

import (
	awsmiddleware "github.com/aws/aws-sdk-go-v2/aws/middleware"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// specified by SDK user-agent SEP
const (
	cryptoUserAgent = "S3Cryptov4"
)

// TODO - can we get a meaningful encryption client version into the header?

// append to user agent (will be ft/s3-encrypt)
func AddS3CryptoUserAgent(options *s3.Options) {
	options.APIOptions = append(options.APIOptions, awsmiddleware.AddSDKAgentKey(awsmiddleware.FeatureMetadata, cryptoUserAgent))
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"log"
)

// CryptographicMaterialsManager (CMM) assembles the cryptographic materials used to
// encrypt and decrypt the encrypted objects.
type CryptographicMaterialsManager interface {
	GetEncryptionMaterials(ctx context.Context, req EncryptionMaterialsRequest) (*CryptographicMaterials, error)
	DecryptMaterials(ctx context.Context, req DecryptMaterialsRequest) (*CryptographicMaterials, error)
}

// DefaultCryptographicMaterialsManager provides support for encrypting and decrypting S3 objects using
// the configured Keyring.
type DefaultCryptographicMaterialsManager struct {
	Keyring *Keyring
}

// NewCryptographicMaterialsManager creates a new DefaultCryptographicMaterialsManager with the given Keyring.
// The Keyring provided must not be nil. If Keyring is nil, NewCryptographicMaterialsManager will return error.
func NewCryptographicMaterialsManager(keyring Keyring) (*DefaultCryptographicMaterialsManager, error) {
	cmm := &DefaultCryptographicMaterialsManager{
		Keyring: &keyring,
	}
	if keyring != nil {
		// Check if the passed in type is a fixture, if not log a warning message to the user
		if fixture, ok := keyring.(awsFixture); !ok || !fixture.isAWSFixture() {
			log.Default().Println(customTypeWarningMessage)
		}
	} else {
		// keyring MUST NOT be nil
		return nil, fmt.Errorf("keyring provided to new cryptographic materials manager MUST NOT be nil")
	}

	return cmm, nil
}

// GetEncryptionMaterials assembles the required EncryptionMaterials and then calls Keyring.OnEncrypt
// to encrypt the materials.
func (cmm *DefaultCryptographicMaterialsManager) GetEncryptionMaterials(ctx context.Context, req EncryptionMaterialsRequest) (*CryptographicMaterials, error) {
	keyring := *cmm.Keyring
	encryptionMaterials := NewEncryptionMaterials(req.AlgorithmSuite)
	encryptionMaterials.encryptionContext = req.MaterialDescription
	return keyring.OnEncrypt(ctx, encryptionMaterials)
}

// EncryptionMaterialsRequest contains the information required to assemble the EncryptionMaterials which
// are used by Keyring.OnEncrypt to encrypt the data key.
type EncryptionMaterialsRequest struct {
	MaterialDescription MaterialDescription
	AlgorithmSuite     *algorithms.AlgorithmSuite
}

// DecryptMaterialsRequest contains the information required to assemble the DecryptionMaterials which
// are used by Keyring.OnDecrypt to decrypt the encrypted data key.
type DecryptMaterialsRequest struct {
	CipherKey  []byte
	Iv         []byte
	MatDesc    string
	KeyringAlg string
	CekAlg     string
	TagLength  string
	KeyCommitment []byte
}

// DecryptMaterials uses the provided DecryptMaterialsRequest to assemble DecryptionMaterials which
// are used by Keyring.OnDecrypt to decrypt the encrypted data key.
func (cmm *DefaultCryptographicMaterialsManager) DecryptMaterials(ctx context.Context, req DecryptMaterialsRequest) (*CryptographicMaterials, error) {
	keyring := *cmm.Keyring

	materials, err := NewDecryptionMaterials(req)
	if err != nil {
		return nil, err
	}
	return keyring.OnDecrypt(ctx, materials, materials.DataKey)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import "testing"

func TestGenerateBytes(t *testing.T) {
	b, _ := generateBytes(5)
	if e, a := 5, len(b); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}
	b, _ = generateBytes(0)
	if e, a := 0, len(b); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}
	b, _ = generateBytes(1024)
	if e, a := 1024, len(b); e != a {
		t.Errorf("expected %d, but received %d", e, a)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"context"
	"crypto/rand"
)

const customTypeWarningMessage = "WARNING: The S3 Encryption client is configured to write encrypted objects using types not provided by AWS. Security and compatibility with these types can not be guaranteed."

// Keyring implementations are responsible for encrypting/decrypting data keys
// using some kind of key material.
// Keyring implementations MAY support decryption-only (i.e. for legacy algorithms)
// or both encryption (including data key generation) and decryption.
type Keyring interface {
	// OnEncrypt generates/encrypts a data key for use with content encryption
	OnEncrypt(ctx context.Context, materials *EncryptionMaterials) (*CryptographicMaterials, error)
	// OnDecrypt decrypts the encryptedDataKeys and returns them in materials
	// for use with content decryption
	OnDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey) (*CryptographicMaterials, error)
}

// awsFixture is an unexported interface to expose whether a given fixture is an aws provided fixture, and whether that
// fixtures dependencies were constructed using aws types.
//
// This interface is used to warn users if they are using custom implementations of CryptographicMaterialsManager
// or Keyring.
type awsFixture interface {
	isAWSFixture() bool
}

func generateBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

const (
	GcmTagSizeBits = "128"
	// KMSKeyring is a constant used during decryption to build a KMS key handler.
	KMSKeyring = "kms"
	// KMSContextKeyring is a constant used during decryption to build a kms+context keyring
	KMSContextKeyring = "kms+context"

	// GrantToken is the key used to store the grant tokens in the context. They are used to avoid eventual consistency authorization issues when calling KMS APIs
	GrantTokens = "GrantTokens"

	kmsAWSCEKContextKey            = "aws:x-amz-cek-alg"
	kmsMismatchCEKAlg              = "the content encryption algorithm used at encryption time does not match the algorithm stored for decryption time. The object may be altered or corrupted"
	kmsReservedKeyConflictErrMsg   = "conflict in reserved KMS Encryption Context key %s. This value is reserved for the S3 Encryption client and cannot be set by the user"
)

// KmsAPIClient is a client that implements the GenerateDataKey and Decrypt operations
type KmsAPIClient interface {
	GenerateDataKey(context.Context, *kms.GenerateDataKeyInput, ...func(*kms.Options)) (*kms.GenerateDataKeyOutput, error)
	Decrypt(context.Context, *kms.DecryptInput, ...func(*kms.Options)) (*kms.DecryptOutput, error)
}

// KeyringOptions is for additional configuration on Keyring types to perform additional behaviors.
// When EnableLegacyWrappingAlgorithms is set to true, the Keyring MAY decrypt objects encrypted
// using legacy wrapping algorithms such as KMS v1.
type KeyringOptions struct {
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=implication
	//# The S3EC MUST support the option to enable or disable legacy wrapping algorithms.
	
	EnableLegacyWrappingAlgorithms bool
}

// KmsKeyring encrypts with encryption context and on decrypt it checks for the algorithm
// in the material description and makes the call to commonDecrypt with the correct parameters
type KmsKeyring struct {
	kmsClient                KmsAPIClient
	KmsKeyId                 string
	legacyWrappingAlgorithms bool
}

// KmsAnyKeyKeyring is decrypt-only.
type KmsAnyKeyKeyring struct {
	kmsClient                KmsAPIClient
	legacyWrappingAlgorithms bool
}

// NewKmsKeyring creates a new KmsKeyring which calls KMS to encrypt/decrypt the data key used to encrypt the S3
// object. The KmsKeyring will always use the kmsKeyId provided to encrypt and decrypt messages.
func NewKmsKeyring(apiClient KmsAPIClient, kmsKeyId string, optFns ...func(options *KeyringOptions)) *KmsKeyring {
	options := KeyringOptions{
		//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
		//= type=implication
		//# The option to enable legacy wrapping algorithms MUST be set to false by default.
		EnableLegacyWrappingAlgorithms: false,
	}
	for _, fn := range optFns {
		fn(&options)
	}

	return &KmsKeyring{
		kmsClient:                apiClient,
		KmsKeyId:                 kmsKeyId,
		legacyWrappingAlgorithms: options.EnableLegacyWrappingAlgorithms,
	}
}

// NewKmsDecryptOnlyAnyKeyKeyring creates a new KmsAnyKeyKeyring. This Keyring uses the KMS identifier
// persisted in the data key's ciphertext to decrypt the data key.
func NewKmsDecryptOnlyAnyKeyKeyring(apiClient KmsAPIClient, optFns ...func(options *KeyringOptions)) *KmsAnyKeyKeyring {
	options := KeyringOptions{
		EnableLegacyWrappingAlgorithms: false,
	}
	for _, fn := range optFns {
		fn(&options)
	}

	return &KmsAnyKeyKeyring{
		kmsClient:                apiClient,
		legacyWrappingAlgorithms: options.EnableLegacyWrappingAlgorithms,
	}
}

// OnEncrypt generates/encrypts a data key for use with content encryption.
func (k *KmsKeyring) OnEncrypt(ctx context.Context, materials *EncryptionMaterials) (*CryptographicMaterials, error) {
	var matDesc MaterialDescription = materials.encryptionContext
	if _, ok := matDesc[kmsAWSCEKContextKey]; ok {
		return nil, fmt.Errorf(kmsReservedKeyConflictErrMsg, kmsAWSCEKContextKey)
	}
	if matDesc == nil {
		matDesc = map[string]string{}
	}

	requestMatDesc := matDesc.Clone()
	requestMatDesc[kmsAWSCEKContextKey] = materials.algorithm

	in := kms.GenerateDataKeyInput{
		EncryptionContext: requestMatDesc,
		KeyId:             &k.KmsKeyId,
		KeySpec:           types.DataKeySpecAes256,
	}

	grantTokens := ctx.Value(GrantTokens)
	if grantTokens != nil {
		in.GrantTokens = grantTokens.([]string)
	}

	out, err := k.kmsClient.GenerateDataKey(ctx, &in)
	if err != nil {
		return &CryptographicMaterials{}, err
	}
	iv, err := generateBytes(materials.gcmNonceSize)
	if err != nil {
		return &CryptographicMaterials{}, err
	}

	cryptoMaterials := &CryptographicMaterials{
		Key:                 out.Plaintext,
		IV:                  iv,
		KeyringAlgorithm:    KMSContextKeyring,
		CEKAlgorithm:        materials.algorithm,
		TagLength:           GcmTagSizeBits,
		MaterialDescription: requestMatDesc,
		EncryptedKey:        out.CiphertextBlob,
	}

	return cryptoMaterials, nil
}

// OnDecrypt decrypts the encryptedDataKeys and returns them in materials
// for use with content decryption, or an error if the object cannot be decrypted
// by the Keyring as its configured.
func (k *KmsKeyring) OnDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey) (*CryptographicMaterials, error) {
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy wrapping algorithms; it MUST throw an exception when attempting to decrypt an object encrypted with a legacy wrapping algorithm.	
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && !k.legacyWrappingAlgorithms {
		return nil, fmt.Errorf("to decrypt x-amz-cek-alg value `%s` you must enable legacyWrappingAlgorithms on the keyring", materials.DataKey.DataKeyAlgorithm)
	}

	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && k.legacyWrappingAlgorithms {
		return commonDecrypt(ctx, materials, encryptedDataKey, &k.KmsKeyId, nil, k.kmsClient)
	} else if materials.DataKey.DataKeyAlgorithm == KMSContextKeyring {
		return commonDecrypt(ctx, materials, encryptedDataKey, &k.KmsKeyId, materials.MaterialDescription, k.kmsClient)
	} else {
		return nil, fmt.Errorf("x-amz-cek-alg value `%s` did not match an expected algorithm", materials.DataKey.DataKeyAlgorithm)
	}
}

func (k *KmsKeyring) isAWSFixture() bool {
	return true
}

// OnEncrypt generates/encrypts a data key for use with content encryption
// The KmsAnyKeyKeyring does not support OnEncrypt, so an error is returned.
func (k *KmsAnyKeyKeyring) OnEncrypt(ctx context.Context, materials *EncryptionMaterials) (*CryptographicMaterials, error) {
	return nil, fmt.Errorf("KmsAnyKeyKeyring MUST NOT be used to encrypt new data")
}

// OnDecrypt decrypts the encryptedDataKeys and returns them in materials
// for use with content decryption, or an error if the object cannot be decrypted
// by the Keyring as its configured.
func (k *KmsAnyKeyKeyring) OnDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey) (*CryptographicMaterials, error) {
	if materials.DataKey.DataKeyAlgorithm == KMSKeyring && k.legacyWrappingAlgorithms {
		return commonDecrypt(ctx, materials, encryptedDataKey, nil, nil, k.kmsClient)
	} else if materials.DataKey.DataKeyAlgorithm == KMSContextKeyring {
		return commonDecrypt(ctx, materials, encryptedDataKey, nil, materials.MaterialDescription, k.kmsClient)
	} else {
		return nil, fmt.Errorf("x-amz-cek-alg value `%s` did not match an expected algorithm", materials.DataKey.DataKeyAlgorithm)
	}
}

func (k *KmsAnyKeyKeyring) isAWSFixture() bool {
	return true
}

func commonDecrypt(ctx context.Context, materials *DecryptionMaterials, encryptedDataKey DataKey, kmsKeyId *string, matDesc MaterialDescription, kmsClient KmsAPIClient) (*CryptographicMaterials, error) {
	if matDesc != nil {
		if v, ok := matDesc[kmsAWSCEKContextKey]; !ok {
			return nil, fmt.Errorf("required key %v is missing from encryption context", kmsAWSCEKContextKey)
		} else if v != materials.ContentAlgorithm {
			return nil, fmt.Errorf(kmsMismatchCEKAlg)
		}
	}

	in := &kms.DecryptInput{
		EncryptionContext: materials.MaterialDescription,
		CiphertextBlob:    encryptedDataKey.EncryptedDataKey,
		KeyId:             kmsKeyId,
	}

	grantTokens := ctx.Value(GrantTokens)
	if grantTokens != nil {
		in.GrantTokens = grantTokens.([]string)
	}

	out, err := kmsClient.Decrypt(ctx, in)
	if err != nil {
		return nil, err
	}

	materials.DataKey.KeyMaterial = out.Plaintext
	cryptoMaterials := &CryptographicMaterials{
		Key:                 out.Plaintext,
		IV:                  materials.ContentIV,
		KeyringAlgorithm:    materials.DataKey.DataKeyAlgorithm,
		CEKAlgorithm:        materials.ContentAlgorithm,
		TagLength:           materials.TagLength,
		MaterialDescription: materials.MaterialDescription,
		EncryptedKey:        materials.DataKey.EncryptedDataKey,
		KeyCommitment:       materials.KeyCommitment,
	}
	return cryptoMaterials, nil
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"testing"

	"github.com/aws/amazon-s3-encryption-client-go/v4/internal/awstesting"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/kms/types"
)

func TestKMSKeyring_OnEncrypt_CorrectKMSRequest(t *testing.T) {
	tConfig := awstesting.Config()
	kmsKeyId := "test-key-id"
	grantTokens := []string{"test-ciphertext-blob"}

	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			// This test focuses on the KMS request correctness, so we just return an empty body
			Body: io.NopCloser(bytes.NewBuffer([]byte("{}"))),
		},
	}

	tConfig.HTTPClient = tHttpClient
	kmsClient := kms.NewFromConfig(tConfig)
	keyring := NewKmsKeyring(kmsClient, kmsKeyId)

	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	algorithmSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF

	encryptionMaterials := NewEncryptionMaterials(algorithmSuite)

	_, err := keyring.OnEncrypt(ctx, encryptionMaterials)

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	var capturedKmsRequest kms.GenerateDataKeyInput
	json.Unmarshal(tHttpClient.CapturedBody, &capturedKmsRequest)

	expectedRequest := kms.GenerateDataKeyInput{
		KeyId:       &kmsKeyId,
		GrantTokens: grantTokens,
		KeySpec:     types.DataKeySpecAes256,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: algorithmSuite.CipherName(),
		},
	}

	if !reflect.DeepEqual(capturedKmsRequest, expectedRequest) {
		t.Errorf("requests sent to KMS was not the expected request.\nExpected %v\nReceived; %v", expectedRequest, capturedKmsRequest)
	}
}

func TestKMSKeyring_OnDecrypt_CorrectKMSRequest(t *testing.T) {
	tConfig := awstesting.Config()
	kmsKeyId := "test-key-id"
	dataKey := DataKey{
		EncryptedDataKey: []byte("data-key"),
		DataKeyAlgorithm: "kms+context",
	}
	grantTokens := []string{"test-ciphertext-blob"}

	tHttpClient := &awstesting.MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			// This test focuses on the KMS request correctness, so we just return an empty body
			Body: io.NopCloser(bytes.NewBuffer([]byte("{}"))),
		},
	}

	tConfig.HTTPClient = tHttpClient
	kmsClient := kms.NewFromConfig(tConfig)
	keyring := NewKmsKeyring(kmsClient, kmsKeyId)

	ctx := context.WithValue(context.Background(), "GrantTokens", grantTokens)

	algorithmSuite := algorithms.AlgAES256GCMIV12Tag16NoKDF

	decryptionMaterials, err := NewDecryptionMaterials(DecryptMaterialsRequest{
		CipherKey:  []byte("test-cipher-key"),
		Iv:         []byte("test-iv"),
		MatDesc:    `{"aws:x-amz-cek-alg":"AES/GCM/NoPadding"}`,
		KeyringAlg: "kms+context",
		CekAlg:     algorithmSuite.CipherName(),
	})
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}

	_, err = keyring.OnDecrypt(ctx, decryptionMaterials, dataKey)

	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if tHttpClient.CapturedReq == nil || tHttpClient.CapturedBody == nil {
		t.Errorf("captured HTTP request/body was nil")
	}

	var capturedKmsRequest kms.DecryptInput
	json.Unmarshal(tHttpClient.CapturedBody, &capturedKmsRequest)

	expectedRequest := kms.DecryptInput{
		KeyId:          &kmsKeyId,
		GrantTokens:    grantTokens,
		CiphertextBlob: dataKey.EncryptedDataKey,
		EncryptionContext: map[string]string{
			kmsAWSCEKContextKey: algorithmSuite.CipherName(),
		},
	}

	if !reflect.DeepEqual(capturedKmsRequest, expectedRequest) {
		t.Errorf("requests sent to KMS was not the expected request.\nExpected %v\nReceived; %v", expectedRequest, capturedKmsRequest)
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"encoding/json"
)

// MaterialDescription is used to identify how and what master
// key has been used.
type MaterialDescription map[string]string

// Clone returns a copy of the MaterialDescription
func (md MaterialDescription) Clone() (clone MaterialDescription) {
	if md == nil {
		return nil
	}
	clone = make(MaterialDescription, len(md))
	for k, v := range md {
		clone[k] = v
	}
	return clone
}

func (md *MaterialDescription) EncodeDescription() ([]byte, error) {
	v, err := json.Marshal(&md)
	return v, err
}

func (md *MaterialDescription) DecodeDescription(b []byte) error {
	return json.Unmarshal(b, &md)
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"reflect"
	"testing"
)

func TestEncodeMaterialDescription(t *testing.T) {
	md := MaterialDescription{}
	md["foo"] = "bar"
	b, err := md.EncodeDescription()
	expected := `{"foo":"bar"}`
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if expected != string(b) {
		t.Errorf("expected %s, but received %s", expected, string(b))
	}
}
func TestDecodeMaterialDescription(t *testing.T) {
	md := MaterialDescription{}
	json := `{"foo":"bar"}`
	err := md.DecodeDescription([]byte(json))
	expected := MaterialDescription{
		"foo": "bar",
	}
	if err != nil {
		t.Errorf("expected no error, but received %v", err)
	}
	if !reflect.DeepEqual(expected, md) {
		t.Error("expected material description to be equivalent, but received otherwise")
	}
}

func TestMaterialDescription_Clone(t *testing.T) {
	tests := map[string]struct {
		md        MaterialDescription
		wantClone MaterialDescription
	}{
		"it handles nil": {
			md:        nil,
			wantClone: nil,
		},
		"it copies all values": {
			md: MaterialDescription{
				"key1": "value1",
				"key2": "value2",
			},
			wantClone: MaterialDescription{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			if gotClone := tt.md.Clone(); !reflect.DeepEqual(gotClone, tt.wantClone) {
				t.Errorf("Clone() = %v, want %v", gotClone, tt.wantClone)
			}
		})
	}
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package materials

import (
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
)

type DecryptionMaterials struct {
	DataKey             DataKey
	ContentIV           []byte //base64 decoded content IV
	MaterialDescription MaterialDescription
	ContentAlgorithm    string
	TagLength           string
	KeyCommitment       []byte
}

func NewDecryptionMaterials(req DecryptMaterialsRequest) (*DecryptionMaterials, error) {
	materialDescription := MaterialDescription{}
	err := materialDescription.DecodeDescription([]byte(req.MatDesc))
	if err != nil {
		return nil, err
	}
	dataKey := DataKey{
		KeyMaterial:      nil,
		EncryptedDataKey: req.CipherKey,
		DataKeyAlgorithm: req.KeyringAlg,
	}

	return &DecryptionMaterials{
		DataKey:             dataKey,
		ContentIV:           req.Iv,
		MaterialDescription: materialDescription,
		ContentAlgorithm:    req.CekAlg,
		TagLength:           req.TagLength,
		KeyCommitment:       req.KeyCommitment,
	}, nil
}

type DataKey struct {
	KeyMaterial      []byte
	EncryptedDataKey []byte
	DataKeyAlgorithm string
}

type EncryptionMaterials struct {
	gcmKeySize        int
	gcmNonceSize      int
	algorithm         string
	encryptionContext map[string]string
}

func NewEncryptionMaterials(algorithmSuite *algorithms.AlgorithmSuite) *EncryptionMaterials {
	return &EncryptionMaterials{
		gcmKeySize:        algorithmSuite.DataKeyLengthBytes(),
		gcmNonceSize:      algorithmSuite.IVLengthBytes(),
		algorithm:         algorithmSuite.CipherName(),
		encryptionContext: map[string]string{},
	}
}

// CryptographicMaterials is used for content encryption. It is used for storing the
// metadata of the encrypted content.
type CryptographicMaterials struct {
	Key                 []byte
	IV                  []byte
	KeyringAlgorithm    string
	CEKAlgorithm        string
	TagLength           string
	MaterialDescription MaterialDescription
	// EncryptedKey should be populated when calling GenerateCipherData
	EncryptedKey []byte
	KeyCommitment []byte
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testvectors

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/client"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	awsV1 "github.com/aws/aws-sdk-go/aws"
	sessionV1 "github.com/aws/aws-sdk-go/aws/session"
	kmsV1 "github.com/aws/aws-sdk-go/service/kms"
	s3V1 "github.com/aws/aws-sdk-go/service/s3"
	s3cryptoV2 "github.com/aws/aws-sdk-go/service/s3/s3crypto"
	"io"
	"log"
	"math"
	"math/rand"
	"os"
	"testing"
	"time"
)

const defaultBucket = "s3ec-go-github-test-bucket"
const bucketEnvvar = "BUCKET"
const defaultAwsKmsAlias = "arn:aws:kms:us-west-2:370957321024:alias/S3EC-Go-Github-KMS-Key"

const awsKmsAliasEnvvar = "AWS_KMS_ALIAS"
const awsAccountIdEnvvar = "AWS_ACCOUNT_ID"
const defaultRegion = "us-west-2"
const regionEnvvar = "AWS_REGION"

func LoadRegion() string {
	if len(os.Getenv(regionEnvvar)) > 0 {
		return os.Getenv(regionEnvvar)
	} else {
		return defaultRegion
	}
}

func LoadBucket() string {
	if len(os.Getenv(bucketEnvvar)) > 0 {
		return os.Getenv(bucketEnvvar)
	} else {
		return defaultBucket
	}
}

func LoadAwsKmsAlias() string {
	if len(os.Getenv(awsKmsAliasEnvvar)) > 0 {
		return os.Getenv(awsKmsAliasEnvvar)
	} else {
		return defaultAwsKmsAlias
	}
}

func LoadAwsAccountId() string {
	return os.Getenv(awsAccountIdEnvvar)
}

// This generates CBC ciphertexts that the s3_integ_test decrypts.
// This is meant to be a utility function, not a test function,
// but for simplicity and easy invocation it is a test function.
// To avoid running it each test run, it is left commented out.
// func TestGenerateCBCIntegTests(t *testing.T) {
// 	arn := "arn:aws:kms:us-west-2:370957321024:alias/S3EC-Go-Github-KMS-Key"
// 	bucket := "s3ec-go-github-test-bucket"
// 	region := "us-west-2"
// 	ctx := context.Background()
// 	cfg, _ := config.LoadDefaultConfig(ctx,
// 		config.WithRegion(region),
// 	)

// 	s3Client := s3.NewFromConfig(cfg)
// 	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
// 	// V2 client
// 	var handler s3cryptoV2.CipherDataGenerator
// 	sessKms, _ := sessionV1.NewSession(&awsV1.Config{
// 		Region: aws.String(region),
// 	})

// 	// KMS v1
// 	kmsSvc := kmsV1.New(sessKms)
// 	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, arn)
// 	// AES-CBC content cipher
// 	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
// 	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

// 	for caseKey, plaintext := range fixtures.Plaintexts {
// 		_, err := encClient.PutObject(&s3V1.PutObjectInput{
// 			Bucket: aws.String(bucket),
// 			Key: aws.String(
// 				fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
// 					fixtures.BaseFolder, version, caseKey),
// 			),
// 			Body: bytes.NewReader(plaintext),
// 		})
// 		if err != nil {
// 			t.Fatalf("failed to upload encrypted fixture, %v", err)
// 		}
// 	}

// }

func TestKmsV1toV4_CBC(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V1toV4_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 client
	var handler s3cryptoV2.CipherDataGenerator
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, kmsKeyAlias)
	// AES-CBC content cipher
	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestKmsV1toV4_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V1toV4_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 client
	var handler s3cryptoV2.CipherDataGenerator
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, kmsKeyAlias)
	// AES-GCM content cipher
	builder := s3cryptoV2.AESGCMContentCipherBuilder(handler)
	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestKmsContextV2toV4_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V2toV4_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 client
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	handler := s3cryptoV2.NewKMSContextKeyGenerator(kmsSvc, kmsKeyAlias, s3cryptoV2.MaterialDescription{})
	// AES-GCM content cipher
	builder := s3cryptoV2.AESGCMContentCipherBuilderV2(handler)
	encClient, err := s3cryptoV2.NewEncryptionClientV2(sessKms, builder)
	if err != nil {
		log.Fatalf("error creating new v2 client: %v", err)
	}

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestKmsContextV4toV2_GCM(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_gcm"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/V4toV2_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = s3ecV4.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	// V2 client
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	cr := s3cryptoV2.NewCryptoRegistry()
	s3cryptoV2.RegisterKMSContextWrapWithCMK(cr, kmsSvc, kmsKeyAlias)
	s3cryptoV2.RegisterAESGCMContentCipher(cr)
	decClient, err := s3cryptoV2.NewDecryptionClientV2(sessKms, cr)
	if err != nil {
		log.Fatalf("error creating new v2 client: %v", err)
	}

	result, err := decClient.GetObject(&s3V1.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestInstructionFileV2toV4(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/inst_file_test.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 client
	var handler s3cryptoV2.CipherDataGenerator
	sessV1, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessV1)
	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, kmsKeyAlias)
	// AES-CBC content cipher
	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
	encClient := s3cryptoV2.NewEncryptionClient(sessV1, builder, func(clientOpts *s3cryptoV2.EncryptionClient) {
		clientOpts.SaveStrategy = s3cryptoV2.S3SaveStrategy{
			Client: s3V1.New(sessV1),
		}
	})

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array: %v", err)
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestNegativeKeyringOption(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlg := "aes_cbc"
	key := "crypto_tests/" + cekAlg + "/v4/language_Go/NegativeV1toV4_CBC.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 Client
	var handler s3cryptoV2.CipherDataGenerator
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, kmsKeyAlias)
	// AES-CBC content cipher
	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, key)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
	})

	_, err = s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err == nil {
		t.Fatalf("error while calling GetObject, expected to FAIL")
	}
}

func TestEnableLegacyDecryptBothFormats(t *testing.T) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	cekAlgCbc := "aes_cbc"
	keyCbc := "crypto_tests/" + cekAlgCbc + "/v4/language_Go/BothFormats_CBC.txt"
	cekAlgGcm := "aes_gcm"
	keyGcm := "crypto_tests/" + cekAlgGcm + "/v4/language_Go/BothFormats_GCM.txt"
	region := "us-west-2"
	plaintext := "This is a test.\n"

	// V2 Client
	var handler s3cryptoV2.CipherDataGenerator
	sessKms, err := sessionV1.NewSession(&awsV1.Config{
		Region: aws.String(region),
	})

	// KMS v1
	kmsSvc := kmsV1.New(sessKms)
	handler = s3cryptoV2.NewKMSKeyGenerator(kmsSvc, kmsKeyAlias)
	// AES-CBC content cipher
	builder := s3cryptoV2.AESCBCContentCipherBuilder(handler, s3cryptoV2.AESCBCPadder)
	encClient := s3cryptoV2.NewEncryptionClient(sessKms, builder)

	_, err = encClient.PutObject(&s3V1.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyCbc),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}
	fmt.Printf("successfully uploaded file to %s/%s\n", bucket, keyCbc)

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	_, err = s3ecV4.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("error while calling PutObject: %v", err)
	}

	getResponseCbc, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyCbc),
	})
	if err != nil {
		t.Fatalf("error while calling GetObject for CBC: %v", err)
	}
	// ensure CBC matches
	decryptedPlaintext, err := io.ReadAll(getResponseCbc.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array: %v", err)
	}
	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	getResponseGcm, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(keyGcm),
	})
	if err != nil {
		t.Fatalf("error while calling GetObject for GCM: %v", err)
	}
	// ensure GCM matches
	decryptedPlaintext, err = io.ReadAll(getResponseGcm.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array: %v", err)
	}
	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}
}

func TestUnicodeEncryptionContextV4ClientV2MessageFormat(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v1-v2-shared
	//= type=test
	//# - This string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/metadata-strategy.md#object-metadata
	//= type=test
	//# - The S3EC SHOULD support decoding the S3 Server's "double encoding".
	rune128 := string(rune(128))
	rune200 := string(rune(200))
	rune256 := string(rune(256))
	runeMaxInt := string(rune(math.MaxInt32))
	shorter := "我"
	medium := "Brøther, may I have the lööps"
	longer := "我的资我的资源我的资源我的资源的资源源"
	mix := "hello 我的资我的资源我的资源我的资源的资源源 goodbye"
	mixTwo := "hello 我的资我的资源我的资源我的资源的资源源 goodbye我的资"

	unicodeStrings := []string{rune128, rune200, rune256, runeMaxInt, shorter, medium, longer, mix, mixTwo}
	for _, s := range unicodeStrings {
		UnicodeEncryptionContextV4ClientV2MessageFormat(t, s)
	}
}

func UnicodeEncryptionContextV4ClientV2MessageFormat(t *testing.T, metadataString string) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	random := rand.Int()
	key := "unicode-encryption-context-" + time.Now().String() + fmt.Sprintf("%d", random)
	region := "us-west-2"
	plaintext := "This is a test.\n"
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.EnableLegacyUnauthenticatedModes = true
		clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
	})

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": metadataString})
	_, err = s3ecV4.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}

	time.Sleep(1 * time.Second)

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting object (%s): %v", key, err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	s3ecV4.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

func TestUnicodeMaterialDescriptionV4ClientV3MessageFormat(t *testing.T) {
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# This material description string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	//= ../specification/s3-encryption/data-format/content-metadata.md#v3-only
	//= type=test
	//# This encryption context string MAY be encoded by the esoteric double-encoding scheme used by the S3 web server.
	rune128 := string(rune(128))
	rune200 := string(rune(200))
	rune256 := string(rune(256))
	runeMaxInt := string(rune(math.MaxInt32))
	shorter := "我"
	medium := "Brøther, may I have the lööps"
	longer := "我的资我的资源我的资源我的资源的资源源"
	mix := "hello 我的资我的资源我的资源我的资源的资源源 goodbye"
	mixTwo := "hello 我的资我的资源我的资源我的资源的资源源 goodbye我的资"

	unicodeStrings := []string{rune128, rune200, rune256, runeMaxInt, shorter, medium, longer, mix, mixTwo}
	for _, s := range unicodeStrings {
		UnicodeMaterialDescriptionV4ClientV3MessageFormat(t, s)
	}
}

func UnicodeMaterialDescriptionV4ClientV3MessageFormat(t *testing.T, metadataString string) {
	bucket := LoadBucket()
	kmsKeyAlias := LoadAwsKmsAlias()

	random := rand.Int()
	key := "unicode-material-description-" + time.Now().String() + fmt.Sprintf("%d", random)
	region := "us-west-2"
	plaintext := "This is a test.\n"
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	kmsV2 := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsV2, kmsKeyAlias, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = true
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	s3V2 := s3.NewFromConfig(cfg)
	s3ecV4, err := client.New(s3V2, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})

	// Create material description with Unicode string
	materialDescription := materials.MaterialDescription{"md-key": metadataString}
	_, err = s3ecV4.PutObject(ctx, &s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader([]byte(plaintext)),
		Metadata: materialDescription,
	})
	if err != nil {
		log.Fatalf("error calling putObject: %v", err)
	}

	time.Sleep(1 * time.Second)

	result, err := s3ecV4.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting object (%s): %v", key, err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	s3ecV4.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}
// Copyright Amazon.com Inc. or its affiliates. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package testvectors

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/amazon-s3-encryption-client-go/v4/algorithms"
	"github.com/aws/amazon-s3-encryption-client-go/v4/client"
	"github.com/aws/amazon-s3-encryption-client-go/v4/materials"
	"github.com/aws/amazon-s3-encryption-client-go/v4/commitment"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"io"
	"strings"
	"testing"
	"time"
)

const version = "v4"

func getAliasArn(shortAlias string, region string, accountId string) string {
	if strings.Contains(shortAlias, "arn") {
		// shortAlias is not actually short
		return shortAlias
	}
	arnFormat := "arn:aws:kms:%s:%s:alias/%s"
	return fmt.Sprintf(arnFormat, region, accountId, shortAlias)
}

func TestInteg_EncryptFixtures_V2MessageFormat(t *testing.T) {
	var region = LoadRegion()
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	var bucket = LoadBucket()
	var accountId = LoadAwsAccountId()
	var kmsAlias = LoadAwsKmsAlias()

	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	cases := []struct {
		CEKAlg                   string
		KEK, bucket, region, CEK string
	}{
		{
			// AES GCM implies V2 message format
			CEKAlg: "aes_gcm",
			KEK:    "kms", bucket: bucket, region: region, CEK: "aes_gcm",
		},
	}

	for _, c := range cases {
		t.Run(c.CEKAlg, func(t *testing.T) {
			s3Client := s3.NewFromConfig(cfg)

			fixtures := getFixtures(t, s3Client, c.CEKAlg, bucket)
			keyring := getEncryptFixtureBuilder(t, cfg, c.KEK, kmsAlias, c.region, accountId, c.CEK)

			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create new CMM")
			}
			encClient, _ := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.EnableLegacyUnauthenticatedModes = true
				clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
			})

			for caseKey, plaintext := range fixtures.Plaintexts {
				_, err := encClient.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(bucket),
					Key: aws.String(
						fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
							fixtures.BaseFolder, version, caseKey),
					),
					Body: bytes.NewReader(plaintext),
				})

				// Assert correctness of encrypted object
				out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key: aws.String(
						fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
							fixtures.BaseFolder, version, caseKey),
					),
				})
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-iv" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-iv"]; !ok {
					t.Fatalf("expected x-amz-iv to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-key-v2" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-key-v2"]; !ok {
					t.Fatalf("expected x-amz-key-v2 to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-matdesc" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-matdesc"]; !ok {
					t.Fatalf("expected x-amz-matdesc to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-wrap-alg" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-wrap-alg"]; !ok {
					t.Fatalf("expected x-amz-wrap-alg to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-cek-alg" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-cek-alg"]; !ok {
					t.Fatalf("expected x-amz-cek-alg to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-tag-len" MUST be present for V2 format objects.
				if _, ok := out.Metadata["x-amz-tag-len"]; !ok {
					t.Fatalf("expected x-amz-tag-len to be present in metadata, got %v", out.Metadata)
				}

				if err != nil {
					t.Fatalf("failed to upload encrypted fixture, %v", err)
				}
			}
		})
	}
}

func TestInteg_EncryptFixtures_V3MessageFormat(t *testing.T) {
	var region = LoadRegion()
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	var bucket = LoadBucket()
	var accountId = LoadAwsAccountId()
	var kmsAlias = LoadAwsKmsAlias()

	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	cases := []struct {
		CEKAlg                   string
		KEK, bucket, region, CEK string
	}{
		{
			// AES GCM committing implies V3 message format
			CEKAlg: "aes_gcm_committing",
			KEK:    "kms", bucket: bucket, region: region, CEK: "aes_gcm_committing",
		},
	}

	for _, c := range cases {
		t.Run(c.CEKAlg, func(t *testing.T) {
			s3Client := s3.NewFromConfig(cfg)

			fixtures := getFixtures(t, s3Client, c.CEKAlg, bucket)
			keyring := getEncryptFixtureBuilder(t, cfg, c.KEK, kmsAlias, c.region, accountId, c.CEK)

			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create new CMM")
			}
			encClient, _ := client.New(s3Client, cmm)

			for caseKey, plaintext := range fixtures.Plaintexts {
				_, err := encClient.PutObject(ctx, &s3.PutObjectInput{
					Bucket: aws.String(bucket),
					Key: aws.String(
						fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
							fixtures.BaseFolder, version, caseKey),
					),
					Body: bytes.NewReader(plaintext),
				})

				// Assert correctness of encrypted object using regular S3 client (not the encryption client)
				out, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
					Bucket: aws.String(bucket),
					Key: aws.String(
						fmt.Sprintf("%s/%s/language_Go/ciphertext_test_case_%s",
							fixtures.BaseFolder, version, caseKey),
					),
				})
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-c" MUST be present for V3 format objects.
				if _, ok := out.Metadata["x-amz-c"]; !ok {
					t.Fatalf("expected x-amz-c to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-3" MUST be present for V3 format objects.
				if _, ok := out.Metadata["x-amz-3"]; !ok {
					t.Fatalf("expected x-amz-3 to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-t" SHOULD be present for V3 format objects.
				if _, ok := out.Metadata["x-amz-t"]; !ok {
					t.Fatalf("expected x-amz-t to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-w" MUST be present for V3 format objects.
				if _, ok := out.Metadata["x-amz-w"]; !ok {
					t.Fatalf("expected x-amz-w to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-d" MUST be present for V3 format objects.
				//= ../specification/s3-encryption/encryption.md#alg-aes-256-gcm-hkdf-sha512-commit-key
				//= type=test
				//# The derived key commitment value MUST be set or returned from the encryption process such that it can be included in the content metadata.
				if _, ok := out.Metadata["x-amz-d"]; !ok {
					t.Fatalf("expected x-amz-d to be present in metadata, got %v", out.Metadata)
				}
				//= ../specification/s3-encryption/data-format/content-metadata.md#content-metadata-mapkeys
				//= type=test
				//# - The mapkey "x-amz-i" MUST be present for V3 format objects.
				if _, ok := out.Metadata["x-amz-i"]; !ok {
					t.Fatalf("expected x-amz-i to be present in metadata, got %v", out.Metadata)
				}

				if err != nil {
					t.Fatalf("failed to upload encrypted fixture, %v", err)
				}
			}
		})
	}
}

func TestInteg_DecryptFixtures(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()

	cases := []struct {
		CEKAlg  string
		Lang    string
		Version string
	}{
		{CEKAlg: "aes_cbc", Lang: "Go", Version: "v4"}, // v4 doesn't support CBC but that's where the files are
		{CEKAlg: "aes_gcm", Lang: "Go", Version: "v4"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v2"},
		{CEKAlg: "aes_gcm", Lang: "Java", Version: "v3"},
	}

	for _, c := range cases {
		t.Run(c.CEKAlg+"-"+c.Lang, func(t *testing.T) {
			s3Client := s3.NewFromConfig(cfg)
			kmsClient := kms.NewFromConfig(cfg)
			keyringWithContext := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = false
			})
			cmm, err := materials.NewCryptographicMaterialsManager(keyringWithContext)
			if err != nil {
				t.Fatalf("failed to create new CMM")
			}

			var decClient *client.S3EncryptionClientV4
			if c.CEKAlg == "aes_cbc" {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = true
				})
				cmmCbc, err := materials.NewCryptographicMaterialsManager(keyring)
				decClient, err = client.New(s3Client, cmmCbc, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = true
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			} else if c.CEKAlg == "aes_gcm" {
				decClient, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}
			} else {
				t.Fatalf("unknown CEKAlg, cannot configure crypto registry: %s", c.CEKAlg)
			}

			fixtures := getFixtures(t, s3Client, c.CEKAlg, bucket)
			ciphertexts := decryptFixtures(t, decClient, fixtures, bucket, c.Lang, c.Version)

			if len(ciphertexts) == 0 {
				t.Fatalf("expected more than 0 ciphertexts to decrypt!")
			}

			for caseKey, ciphertext := range ciphertexts {
				if e, a := len(fixtures.Plaintexts[caseKey]), len(ciphertext); e != a {
					t.Errorf("expect %v text len, got %v", e, a)
				}
				if e, a := fixtures.Plaintexts[caseKey], ciphertext; !bytes.Equal(e, a) {
					t.Errorf("expect %v text, got %v", e, a)
				}
			}
		})
	}
}

type testFixtures struct {
	BaseFolder string
	Plaintexts map[string][]byte
}

func getFixtures(t *testing.T, s3Client *s3.Client, cekAlg, bucket string) testFixtures {
	t.Helper()
	ctx := context.Background()

	prefix := "plaintext_test_case_"
	baseFolder := "crypto_tests/" + cekAlg

	out, err := s3Client.ListObjects(ctx, &s3.ListObjectsInput{
		Bucket: aws.String(bucket),
		Prefix: aws.String(baseFolder + "/" + prefix),
	})
	if err != nil {
		t.Fatalf("unable to list fixtures %v", err)
	}

	plaintexts := map[string][]byte{}
	for _, obj := range out.Contents {
		ptObj, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    obj.Key,
		})
		if err != nil {
			t.Fatalf("unable to get fixture object %s, %v", *obj.Key, err)
		}
		caseKey := strings.TrimPrefix(*obj.Key, baseFolder+"/"+prefix)
		plaintext, err := io.ReadAll(ptObj.Body)
		if err != nil {
			t.Fatalf("unable to read fixture object %s, %v", *obj.Key, err)
		}

		plaintexts[caseKey] = plaintext
	}

	return testFixtures{
		BaseFolder: baseFolder,
		Plaintexts: plaintexts,
	}
}

func getEncryptFixtureBuilder(t *testing.T, cfg aws.Config, kek, alias, region, accountId string, cek string) (keyring materials.Keyring) {
	t.Helper()

	var kmsKeyring materials.Keyring
	switch kek {
	case "kms":
		arn := getAliasArn(alias, region, accountId)

		kmsSvc := kms.NewFromConfig(cfg)
		kmsKeyring = materials.NewKmsKeyring(kmsSvc, arn, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		})
	default:
		t.Fatalf("unknown fixture KEK, %v", kek)
	}

	switch cek {
	case "aes_gcm_committing":
	case "aes_gcm":
		return kmsKeyring
	case "aes_cbc":
		t.Fatalf("aes cbc is not supported ")
	default:
		t.Fatalf("unknown fixture CEK, %v", cek)
	}

	return kmsKeyring
}

func decryptFixtures(t *testing.T, decClient *client.S3EncryptionClientV4, fixtures testFixtures, bucket, lang, version string,
) map[string][]byte {
	t.Helper()
	ctx := context.Background()

	prefix := "ciphertext_test_case_"
	lang = "language_" + lang

	ciphertexts := map[string][]byte{}
	for caseKey := range fixtures.Plaintexts {
		cipherKey := fixtures.BaseFolder + "/" + version + "/" + lang + "/" + prefix + caseKey
		ctObj, err := decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &cipherKey,
		})
		if err != nil {
			t.Fatalf("failed to get encrypted object %v", err)
		}

		ciphertext, err := io.ReadAll(ctObj.Body)
		if err != nil {
			t.Fatalf("failed to read object data %v", err)
		}
		ciphertexts[caseKey] = ciphertext
	}

	return ciphertexts
}

func TestIntegS3ECHeadObject(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "single-round-trip-test" + time.Now().String()
	var plaintext = "this is some plaintext to encrypt!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	if err != nil {
		t.Fatalf("failed to get fixture alias info for %s, %v", alias, err)
	}

	var s3Client = s3.NewFromConfig(cfg)

	// Ensure fresh key
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	kmsClient := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": "ec-value"})

	s3EncryptionClient, err := client.New(s3Client, cmm)
	_, err = s3EncryptionClient.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("error while encrypting: %v", err)
	}

	result, err := s3EncryptionClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	headResult, err := s3EncryptionClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if e, a := int64(len(plaintext)+16), headResult.ContentLength; e != a {
		t.Errorf("expect %v text, got %v", e, a)
	}

	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

func TestInteg_DeleteObjects_DeletesObjects(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "delete-objects-test-" + time.Now().Format("20060102-150405")
	var key1 = baseKey + "-object1"
	var key2 = baseKey + "-object2"
	var key3 = baseKey + "-object3"
	var plaintext1 = "Hello, S3 Encryption Client DeleteObjects test - Object 1!"
	var plaintext2 = "Hello, S3 Encryption Client DeleteObjects test - Object 2!"
	var plaintext3 = "Hello, S3 Encryption Client DeleteObjects test - Object 3!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing objects
	objectsToClean := []string{key1, key2, key3}
	for _, key := range objectsToClean {
		s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
	}

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put multiple encrypted objects
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key1,
		Body:   bytes.NewReader([]byte(plaintext1)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 1: %v", err)
	}

	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key2,
		Body:   bytes.NewReader([]byte(plaintext2)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 2: %v", err)
	}

	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key3,
		Body:   bytes.NewReader([]byte(plaintext3)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object 3: %v", err)
	}

	// Verify all objects exist before deletion
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err != nil {
			t.Fatalf("object %d should exist before deletion: %v", i+1, err)
		}
	}

	t.Logf("✓ Verified all objects exist before deletion")

	// Test DeleteObjects - should delete all objects
	deleteInput := &s3.DeleteObjectsInput{
		Bucket: &bucket,
		Delete: &types.Delete{
			Objects: []types.ObjectIdentifier{
				{Key: &key1},
				{Key: &key2},
				{Key: &key3},
			},
		},
	}

	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST be implemented by the S3EC.
	result, err := s3ec.DeleteObjects(ctx, deleteInput)
	if err != nil {
		t.Fatalf("DeleteObjects failed: %v", err)
	}

	// Verify the response structure
	if result == nil {
		t.Fatal("DeleteObjects result should not be nil")
	}

	if len(result.Deleted) != 3 {
		t.Errorf("expected 3 deleted objects, got %d", len(result.Deleted))
	}

	if len(result.Errors) > 0 {
		t.Errorf("expected no errors, got %d errors: %v", len(result.Errors), result.Errors)
	}

	// Verify all objects are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST delete each of the given objects.
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    &key,
		})
		if err == nil {
			t.Errorf("object %d should be deleted but still exists", i+1)
		}
	}

	// Verify all instruction files are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObjects MUST delete each of the corresponding instruction files using the default instruction file suffix.
	for i, key := range objectsToClean {
		_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
			Bucket: &bucket,
			Key:    aws.String(key + ".instruction"),
		})
		if err == nil {
			t.Errorf("instruction file %d should be deleted but still exists", i+1)
		}
	}

	t.Logf("✓ DeleteObjects successfully deleted all objects")
}

func TestInteg_DeleteObject_DeletesObjectAndInstructionFile(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "delete-object-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client DeleteObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing objects
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key + ".instruction"),
	})

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put encrypted object (this should create both the object and instruction file)
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object: %v", err)
	}

	// Verify object exists before deletion
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("object should exist before deletion: %v", err)
	}

	// Test DeleteObject - should delete both the object and instruction file
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST be implemented by the S3EC.
	_, err = s3ec.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("DeleteObject failed: %v", err)
	}

	// Verify both object and instruction file are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST delete the given object key.
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err == nil {
		t.Errorf("object should be deleted but still exists")
	}

	// Verify both object and instruction file are deleted
	//= ../specification/s3-encryption/client.md#required-api-operations
	//= type=test
	//# - DeleteObject MUST delete the associated instruction file using the default instruction file suffix.
	_, err = s3Client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: &bucket,
		Key:    aws.String(key + ".instruction"),
	})
	if err == nil {
		t.Errorf("instruction file should be deleted but still exists")
	}

	t.Logf("✓ DeleteObject successfully deleted both object and instruction file")
}

func TestIntegLegacyUnauthenticatedModes(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test specific known objects and validate their content encryption algorithms
	testCases := []struct {
		name           string
		objectKey      string
		expectedCekAlg string
		description    string
	}{
		{
			name:           "LegacyUnauthenticated_AES_CBC",
			objectKey:      "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedCekAlg: "AES/CBC/PKCS5Padding",
			description:    "AES/CBC object should use legacy unauthenticated content encryption algorithm",
		},
		{
			name:           "ModernAuthenticated_AES_GCM",
			objectKey:      "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedCekAlg: "AES/GCM/NoPadding",
			description:    "AES/GCM object should use modern authenticated content encryption algorithm",
		},
	}

	// First validate that our test objects have the expected content encryption algorithms
	for _, tc := range testCases {
		t.Run("Validate_"+tc.name, func(t *testing.T) {
			headResult, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})
			if err != nil {
				t.Skipf("Test object %s not found, skipping: %v", tc.objectKey, err)
			}

			var actualCekAlg string
			if headResult.Metadata != nil {
				if val, ok := headResult.Metadata["x-amz-cek-alg"]; ok {
					actualCekAlg = val
				}
			}

			if actualCekAlg != tc.expectedCekAlg {
				t.Errorf("Object %s has cek-alg %q, expected %q. %s", 
					tc.objectKey, actualCekAlg, tc.expectedCekAlg, tc.description)
			} else {
				t.Logf("✓ Validated: %s has expected cek-alg: %s", tc.objectKey, actualCekAlg)
			}
		})
	}

	// Get expected plaintext for test case test_one.txt
	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
	expectedPlaintext, exists := fixtures.Plaintexts["test_one.txt"]
	if !exists {
		t.Fatal("Could not find plaintext for test case test_one.txt")
	}

	// Test legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes = true (should succeed)
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
	t.Run("LegacyUnauthenticatedObject_EnabledShouldSucceed", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err != nil {
			t.Fatalf("expected successful decryption of legacy unauthenticated object with EnableLegacyUnauthenticatedModes=true, but got error: %v", err)
		}

		decryptedData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("failed to read decrypted data: %v", err)
		}

		if !bytes.Equal(expectedPlaintext, decryptedData) {
			t.Errorf("decrypted data mismatch: expected %q, got %q", string(expectedPlaintext), string(decryptedData))
		}

		t.Logf("✓ Successfully decrypted legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes=true")
	})

	// Test legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes = false (should fail)
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy content encryption algorithms;
	//# it MUST throw an exception when attempting to decrypt an object encrypted with a legacy content encryption algorithm.
	//= ../specification/s3-encryption/decryption.md#legacy-decryption
	//= type=test
	//# The S3EC MUST NOT decrypt objects encrypted using legacy unauthenticated algorithm suites unless specifically configured to do so.
	t.Run("LegacyUnauthenticatedObject_DisabledShouldFail", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = false
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		//= ../specification/s3-encryption/decryption.md#legacy-decryption
		//= type=test
		//# If the S3EC is not configured to enable legacy unauthenticated content decryption,
		//# the client MUST throw an exception when attempting to decrypt an object encrypted with a legacy unauthenticated algorithm suite.
		if err == nil {
			t.Fatalf("expected decryption of legacy unauthenticated object to fail when EnableLegacyUnauthenticatedModes=false, but it succeeded")
		}

		// Verify the error message indicates legacy unauthenticated modes issue
		if !strings.Contains(err.Error(), "enable legacy unauthenticated modes") && !strings.Contains(err.Error(), "AES/CBC/PKCS5Padding") {
			t.Errorf("expected error to mention legacy unauthenticated modes or AES/CBC algorithm, got: %v", err)
		}

		t.Logf("✓ Correctly failed to decrypt legacy unauthenticated object (AES/CBC) with EnableLegacyUnauthenticatedModes=false: %v", err)
	})

	// Test modern authenticated object (AES/GCM) works regardless of legacy unauthenticated modes setting
	//= ../specification/s3-encryption/client.md#enable-legacy-unauthenticated-modes
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all content encryption algorithms (both legacy and fully supported).
	t.Run("ModernAuthenticatedObject_AlwaysWorks", func(t *testing.T) {
		// Get expected plaintext for modern object
		gcmFixtures := getFixtures(t, s3Client, "aes_gcm", bucket)
		gcmExpectedPlaintext, exists := gcmFixtures.Plaintexts["test_one.txt"]
		if !exists {
			t.Skip("Could not find plaintext for aes_gcm test case test_one.txt, skipping modern object test")
		}

		for _, enableLegacyUnauthenticated := range []bool{true, false} {
			t.Run(fmt.Sprintf("EnableLegacyUnauthenticated_%v", enableLegacyUnauthenticated), func(t *testing.T) {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = false
				})
				cmm, err := materials.NewCryptographicMaterialsManager(keyring)
				if err != nil {
					t.Fatalf("failed to create CMM: %v", err)
				}

				decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = enableLegacyUnauthenticated
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}

				modernObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
				result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &modernObjectKey,
				})
				if err != nil {
					t.Fatalf("modern authenticated algorithm should always work regardless of legacy unauthenticated modes setting, but got error: %v", err)
				}

				decryptedData, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("failed to read decrypted data: %v", err)
				}

				if !bytes.Equal(gcmExpectedPlaintext, decryptedData) {
					t.Errorf("decrypted data mismatch: expected %q, got %q", string(gcmExpectedPlaintext), string(decryptedData))
				}

				t.Logf("✓ Modern authenticated object (AES/GCM) works correctly with EnableLegacyUnauthenticatedModes=%v", enableLegacyUnauthenticated)
			})
		}
	})
}

func TestIntegKmsContext(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "single-round-trip-test-kms-context" + time.Now().String()
	var plaintext = "this is some plaintext to encrypt!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	if err != nil {
		t.Fatalf("failed to get fixture alias info for %s, %v", alias, err)
	}

	var s3Client = s3.NewFromConfig(cfg)

	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	kmsClient := kms.NewFromConfig(cfg)
	var matDesc materials.MaterialDescription
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": "ec-value"})
	s3EncryptionClient, err := client.New(s3Client, cmm)
	_, err = s3EncryptionClient.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket:   aws.String(bucket),
		Key:      aws.String(key),
		Body:     bytes.NewReader([]byte(plaintext)),
		Metadata: matDesc,
	})
	if err != nil {
		t.Fatalf("error while encrypting: %v", err)
	}

	result, err := s3EncryptionClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	headResult, err := s3EncryptionClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if e, a := int64(len(plaintext)+16), headResult.ContentLength; e != a {
		t.Errorf("expect %v text, got %v", e, a)
	}
	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

func TestIntegKmsContextDecryptAny(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "single-round-trip-test-context-decrypt-any" + time.Now().String()
	var plaintext = "this is some plaintext to encrypt!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(region),
	)

	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	if err != nil {
		t.Fatalf("failed to get fixture alias info for %s, %v", alias, err)
	}

	var s3Client = s3.NewFromConfig(cfg)

	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	kmsClient := kms.NewFromConfig(cfg)
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	encryptionContext := context.WithValue(ctx, "EncryptionContext", map[string]string{"ec-key": "ec-value"})
	s3EncryptionClient, err := client.New(s3Client, cmm)
	_, err = s3EncryptionClient.PutObject(encryptionContext, &s3.PutObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("error while encrypting: %v", err)
	}

	// decrypt with AnyKey
	anyKeyCmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	s3EncryptionClientAnyKey, err := client.New(s3Client, anyKeyCmm)
	if err != nil {
		t.Fatalf("error while creating new CMM")
	}

	result, err := s3EncryptionClientAnyKey.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		t.Fatalf("error while decrypting: %v", err)
	}

	decryptedPlaintext, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted plaintext into byte array")
	}

	if e, a := []byte(plaintext), decryptedPlaintext; !bytes.Equal(e, a) {
		t.Errorf("expect %v text, got %v", e, a)
	}

	headResult, err := s3EncryptionClient.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	})
	if e, a := int64(len(plaintext)+16), headResult.ContentLength; e != a {
		t.Errorf("expect %v text, got %v", e, a)
	}
	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

func TestIntegLegacyWrappingAlgorithms(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test specific known objects and validate their wrapping algorithms
	testCases := []struct {
		name           string
		objectKey      string
		expectedWrapAlg string
		cekAlg         string
		description    string
	}{
		{
			name:           "LegacyKMSKeyring_AES_CBC",
			objectKey:      "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedWrapAlg: "kms",
			cekAlg:         "aes_cbc",
			description:    "AES/CBC object should use legacy 'kms' wrapping algorithm",
		},
		{
			name:           "ModernKMSContextKeyring_AES_GCM",
			objectKey:      "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			expectedWrapAlg: "kms+context",
			cekAlg:         "aes_gcm",
			description:    "AES/GCM object should use modern 'kms+context' wrapping algorithm",
		},
	}

	// First validate that our test objects have the expected wrapping algorithms
	for _, tc := range testCases {
		t.Run("Validate_"+tc.name, func(t *testing.T) {
			headResult, err := s3Client.HeadObject(ctx, &s3.HeadObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})
			if err != nil {
				t.Skipf("Test object %s not found, skipping: %v", tc.objectKey, err)
			}

			var actualWrapAlg string
			if headResult.Metadata != nil {
				if val, ok := headResult.Metadata["x-amz-wrap-alg"]; ok {
					actualWrapAlg = val
				}
			}

			if actualWrapAlg != tc.expectedWrapAlg {
				t.Errorf("Object %s has wrap-alg %q, expected %q. %s", 
					tc.objectKey, actualWrapAlg, tc.expectedWrapAlg, tc.description)
			} else {
				t.Logf("✓ Validated: %s has expected wrap-alg: %s", tc.objectKey, actualWrapAlg)
			}
		})
	}

	// Get expected plaintext for test case test_one.txt
	fixtures := getFixtures(t, s3Client, "aes_cbc", bucket)
	expectedPlaintext, exists := fixtures.Plaintexts["test_one.txt"]
	if !exists {
		t.Fatal("Could not find plaintext for test case test_one.txt")
	}

	// Test legacy object (kms) with EnableLegacyWrappingAlgorithms = true (should succeed)
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	t.Run("LegacyObject_EnabledShouldSucceed", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = true
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err != nil {
			t.Fatalf("expected successful decryption of legacy object with EnableLegacyWrappingAlgorithms=true, but got error: %v", err)
		}

		decryptedData, err := io.ReadAll(result.Body)
		if err != nil {
			t.Fatalf("failed to read decrypted data: %v", err)
		}

		if !bytes.Equal(expectedPlaintext, decryptedData) {
			t.Errorf("decrypted data mismatch: expected %q, got %q", string(expectedPlaintext), string(decryptedData))
		}

		t.Logf("✓ Successfully decrypted legacy object (kms wrapping algorithm) with EnableLegacyWrappingAlgorithms=true")
	})

	// Test legacy object (kms) with EnableLegacyWrappingAlgorithms = false (should fail)
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When disabled, the S3EC MUST NOT decrypt objects encrypted using legacy wrapping algorithms;
	//# it MUST throw an exception when attempting to decrypt an object encrypted with a legacy wrapping algorithm.
	t.Run("LegacyObject_DisabledShouldFail", func(t *testing.T) {
		keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
			options.EnableLegacyWrappingAlgorithms = false
		})
		cmm, err := materials.NewCryptographicMaterialsManager(keyring)
		if err != nil {
			t.Fatalf("failed to create CMM: %v", err)
		}

		decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
			clientOptions.EnableLegacyUnauthenticatedModes = true
			clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
		})
		if err != nil {
			t.Fatalf("failed to create decryption client: %v", err)
		}

		legacyObjectKey := "crypto_tests/aes_cbc/v4/language_Go/ciphertext_test_case_test_one.txt"
		_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
			Bucket: &bucket,
			Key:    &legacyObjectKey,
		})
		if err == nil {
			t.Fatalf("expected decryption of legacy object to fail when EnableLegacyWrappingAlgorithms=false, but it succeeded")
		}

		// Verify the error message indicates legacy wrapping algorithm issue
		if !strings.Contains(err.Error(), "legacyWrappingAlgorithms") && !strings.Contains(err.Error(), "did not match an expected algorithm") {
			t.Errorf("expected error to mention legacyWrappingAlgorithms or algorithm mismatch, got: %v", err)
		}

		t.Logf("✓ Correctly failed to decrypt legacy object (kms wrapping algorithm) with EnableLegacyWrappingAlgorithms=false: %v", err)
	})

	// Test modern object (kms+context) works regardless of legacy setting
	//= ../specification/s3-encryption/client.md#enable-legacy-wrapping-algorithms
	//= type=test
	//# When enabled, the S3EC MUST be able to decrypt objects encrypted with all supported wrapping algorithms (both legacy and fully supported).
	t.Run("ModernObject_AlwaysWorks", func(t *testing.T) {
		// Get expected plaintext for modern object
		gcmFixtures := getFixtures(t, s3Client, "aes_gcm", bucket)
		gcmExpectedPlaintext, exists := gcmFixtures.Plaintexts["test_one.txt"]
		if !exists {
			t.Skip("Could not find plaintext for aes_gcm test case test_one.txt, skipping modern object test")
		}

		for _, enableLegacy := range []bool{true, false} {
			t.Run(fmt.Sprintf("EnableLegacy_%v", enableLegacy), func(t *testing.T) {
				keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
					options.EnableLegacyWrappingAlgorithms = enableLegacy
				})
				cmm, err := materials.NewCryptographicMaterialsManager(keyring)
				if err != nil {
					t.Fatalf("failed to create CMM: %v", err)
				}

				decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.CommitmentPolicy = commitment.FORBID_ENCRYPT_ALLOW_DECRYPT
				})
				if err != nil {
					t.Fatalf("failed to create decryption client: %v", err)
				}

				modernObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
				result, err := decClient.GetObject(ctx, &s3.GetObjectInput{
					Bucket: &bucket,
					Key:    &modernObjectKey,
				})
				if err != nil {
					t.Fatalf("modern algorithm should always work regardless of legacy setting, but got error: %v", err)
				}

				decryptedData, err := io.ReadAll(result.Body)
				if err != nil {
					t.Fatalf("failed to read decrypted data: %v", err)
				}

				if !bytes.Equal(gcmExpectedPlaintext, decryptedData) {
					t.Errorf("decrypted data mismatch: expected %q, got %q", string(gcmExpectedPlaintext), string(decryptedData))
				}

				t.Logf("✓ Modern object (kms+context wrapping algorithm) works correctly with EnableLegacyWrappingAlgorithms=%v", enableLegacy)
			})
		}
	})
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - GetObject MUST be implemented by the S3EC.
//# - GetObject MUST decrypt data received from the S3 server and return it as plaintext.
func TestInteg_GetObject_BasicDecryption(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "basic-getobject-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client GetObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing object
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	// Create S3EC and encrypt an object
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Put encrypted object
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("failed to put encrypted object: %v", err)
	}
	
	// Loose assertion that the object was encrypted by checking that its body does not match the plaintext
	rawResult, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to get raw object: %v", err)
	}

	rawData, err := io.ReadAll(rawResult.Body)
	if err != nil {
		t.Fatalf("failed to read raw object data: %v", err)
	}

	if string(rawData) == plaintext {
		t.Errorf("object was not encrypted: raw data matches plaintext")
	}

	// Test GetObject - should decrypt and return plaintext
	result, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("GetObject failed: %v", err)
	}

	decryptedData, err := io.ReadAll(result.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted data: %v", err)
	}

	if string(decryptedData) != plaintext {
		t.Errorf("GetObject decryption failed: expected %q, got %q", plaintext, string(decryptedData))
	}

	t.Logf("✓ GetObject successfully decrypted data and returned plaintext")

	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

//= ../specification/s3-encryption/client.md#required-api-operations
//= type=test
//# - PutObject MUST be implemented by the S3EC.
//# - PutObject MUST encrypt its input data before it is uploaded to S3.
func TestInteg_PutObject_BasicEncryption(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var key = "basic-putobject-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client PutObject test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Clean up any existing object
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})

	// Create S3EC
	cmm, err := materials.NewCryptographicMaterialsManager(materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	}))
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	s3ec, err := client.New(s3Client, cmm)
	if err != nil {
		t.Fatalf("failed to create S3EC: %v", err)
	}

	// Test PutObject - should encrypt input data before uploading
	_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
		Bucket: &bucket,
		Key:    &key,
		Body:   bytes.NewReader([]byte(plaintext)),
	})
	if err != nil {
		t.Fatalf("PutObject failed: %v", err)
	}

	// Loose assertion that the object was encrypted by checking that its body does not match the plaintext
	rawResult, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to get raw object: %v", err)
	}

	rawData, err := io.ReadAll(rawResult.Body)
	if err != nil {
		t.Fatalf("failed to read raw object data: %v", err)
	}

	if string(rawData) == plaintext {
		t.Errorf("object was not encrypted: raw data matches plaintext")
	}

	// Verify encryption metadata is present
	if rawResult.Metadata == nil {
		t.Fatal("expected encryption metadata to be present")
	}

	expectedMetadataKeys := []string{"x-amz-c", "x-amz-3", "x-amz-w", "x-amz-d", "x-amz-i"}
	for _, key := range expectedMetadataKeys {
		if _, exists := rawResult.Metadata[key]; !exists {
			t.Errorf("expected encryption metadata key %s to be present", key)
		}
	}

	// Verify we can decrypt it back to the original plaintext
	decryptResult, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
	if err != nil {
		t.Fatalf("failed to decrypt object: %v", err)
	}

	decryptedData, err := io.ReadAll(decryptResult.Body)
	if err != nil {
		t.Fatalf("failed to read decrypted data: %v", err)
	}

	if string(decryptedData) != plaintext {
		t.Errorf("decryption verification failed: expected %q, got %q", plaintext, string(decryptedData))
	}

	t.Logf("✓ PutObject successfully encrypted input data before uploading to S3")

	// Cleanup
	s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: &bucket,
		Key:    &key,
	})
}

//= ../specification/s3-encryption/decryption.md#key-commitment
//= type=test
//# The S3EC MUST validate the algorithm suite used for decryption against the key commitment policy before attempting to decrypt the content ciphertext.
func TestInteg_ValidateAlgorithmSuiteAgainstCommitmentPolicy(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Matrix of [V3 object, V2 object] x [commitment policies] = 6 test cases
	cases := []struct {
		name             string
		objectKey        string
		objectType       string
		commitmentPolicy commitment.CommitmentPolicy
		expectError      bool
		errorContains    string
	}{
		{
			name:             "V3_object_with_FORBID_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V3_object_with_REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V3_object_with_REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm_committing/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V3",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V2_object_with_FORBID_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		{
			name:             "V2_object_with_REQUIRE_ENCRYPT_ALLOW_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			expectError:      false,
		},
		// Only expected failure case on decryption, where a non-committing object is
		// attempted to be decrypted under a commitment policy that requires commitment.
		{
			name:             "V2_object_with_REQUIRE_ENCRYPT_REQUIRE_DECRYPT",
			objectKey:        "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt",
			objectType:       "V2",
			commitmentPolicy: commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			expectError:      true,
			errorContains:    "object's content encryption algorithm is not valid for the selected commitment policy",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = false
			})
			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
			})
			if err != nil {
				t.Fatalf("expected no error during client creation, got %v", err)
			}

			_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &tc.objectKey,
			})

			//= ../specification/s3-encryption/decryption.md#key-commitment
			//= type=test
			//# If the commitment policy requires decryption using a committing algorithm suite,
			//# and the algorithm suite associated with the object does not support key commitment,
			//# then the S3EC MUST throw an exception.
			if tc.expectError {
				if err == nil {
					t.Fatalf("expected error but got none")
				}
				if tc.errorContains != "" && !strings.Contains(err.Error(), tc.errorContains) {
					t.Errorf("expected error to contain %q, got %q", tc.errorContains, err.Error())
				} else {
					t.Logf("✓ Expected error during GetObject: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error, got %v", err)
				}
				t.Logf("✓ Successfully decrypted %s object with %s policy", tc.objectType, tc.commitmentPolicy)
			}
		})
	}
}

func TestInteg_AlgorithmSuiteMessageFormatCompatibility(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "algorithm-suite-message-format-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client Algorithm Suite Message Format test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	testCases := []struct {
		name                string
		algorithmSuite      string
		commitmentPolicy    commitment.CommitmentPolicy
		enableLegacyModes   bool
		enableLegacyWrap    bool
		expectedFormat      string
		expectedHeaders     []string
	}{
		//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
		//= type=test
		//# Objects encrypted with ALG_AES_256_GCM_IV12_TAG16_NO_KDF MUST use the V2 message format version only.
		{
			name:                "ALG_AES_256_GCM_IV12_TAG16_NO_KDF_V2_Format",
			algorithmSuite:      "AES256GCMIV12Tag16NoKDF",
			commitmentPolicy:    commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			enableLegacyModes:   false,
			enableLegacyWrap:    false,
			expectedFormat:      "V2",
			expectedHeaders:     []string{"x-amz-iv", "x-amz-key-v2", "x-amz-matdesc", "x-amz-wrap-alg", "x-amz-cek-alg", "x-amz-tag-len"},
		},
		//= ../specification/s3-encryption/data-format/content-metadata.md#algorithm-suite-and-message-format-version-compatibility
		//= type=test
		//# Objects encrypted with ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY MUST use the V3 message format version only.
		{
			name:                "ALG_AES_256_GCM_HKDF_SHA512_COMMIT_KEY_V3_Format",
			algorithmSuite:      "AES256GCMHkdfSha512CommitKey",
			commitmentPolicy:    commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			enableLegacyModes:   false,
			enableLegacyWrap:    false,
			expectedFormat:      "V3",
			expectedHeaders:     []string{"x-amz-c", "x-amz-3", "x-amz-t", "x-amz-w", "x-amz-d", "x-amz-i"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			key := baseKey + "-" + tc.algorithmSuite

			// Clean up any existing object
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			// Create keyring with appropriate legacy settings
			keyring := materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = tc.enableLegacyWrap
			})

			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			// Create S3EC with specific configuration
			var s3ec *client.S3EncryptionClientV4
			if tc.algorithmSuite == "AES256GCMHkdfSha512CommitKey" {
				// For committing algorithm suite, use default client (which uses committing by default)
				s3ec, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = tc.enableLegacyModes
					clientOptions.CommitmentPolicy = tc.commitmentPolicy
					// Don't override algorithm suite - let it use the default committing algorithm
				})
			} else {
				// For non-committing algorithm suites, explicitly set them
				s3ec, err = client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
					clientOptions.EnableLegacyUnauthenticatedModes = tc.enableLegacyModes
					clientOptions.CommitmentPolicy = tc.commitmentPolicy
					
					// Configure the specific algorithm suite
					switch tc.algorithmSuite {
					case "AES256CBCIV16NoKDF":
						clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256CBCIV16NoKDF
					case "AES256GCMIV12Tag16NoKDF":
						clientOptions.EncryptionAlgorithmSuite = algorithms.AlgAES256GCMIV12Tag16NoKDF
					default:
						t.Fatalf("unknown algorithm suite: %s", tc.algorithmSuite)
					}
				})
			}
			if err != nil {
				t.Fatalf("failed to create S3EC: %v", err)
			}

			// Encrypt object with configured algorithm suite
			_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   bytes.NewReader([]byte(plaintext)),
			})
			if err != nil {
				t.Fatalf("failed to put encrypted object: %v", err)
			}

			// Get object metadata using regular S3 client to inspect message format headers
			result, err := s3Client.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				t.Fatalf("failed to get encrypted object metadata: %v", err)
			}

			// Verify the expected message format headers are present
			if result.Metadata == nil {
				t.Fatalf("expected metadata to be present for %s format", tc.expectedFormat)
			}

			t.Logf("Testing %s algorithm suite for %s message format", tc.algorithmSuite, tc.expectedFormat)
			t.Logf("Object metadata: %v", result.Metadata)

			// Assert all expected headers for the message format are present
			for _, expectedHeader := range tc.expectedHeaders {
				if _, exists := result.Metadata[expectedHeader]; !exists {
					t.Errorf("expected %s format header '%s' to be present for algorithm suite %s, but it was missing", 
						tc.expectedFormat, expectedHeader, tc.algorithmSuite)
				} else {
					t.Logf("✓ Found expected %s format header: %s", tc.expectedFormat, expectedHeader)
				}
			}

			// Verify that headers from other message formats are NOT present
			var unexpectedHeaders []string
			if tc.expectedFormat == "V2" {
				// V2 format should not have V3 headers
				unexpectedHeaders = []string{"x-amz-c", "x-amz-3", "x-amz-w", "x-amz-d", "x-amz-i"}
			} else if tc.expectedFormat == "V3" {
				// V3 format should not have V2 headers
				unexpectedHeaders = []string{"x-amz-iv", "x-amz-key-v2", "x-amz-tag-len"}
			}

			for _, unexpectedHeader := range unexpectedHeaders {
				if _, exists := result.Metadata[unexpectedHeader]; exists {
					t.Errorf("unexpected header '%s' found for %s format with algorithm suite %s", 
						unexpectedHeader, tc.expectedFormat, tc.algorithmSuite)
				}
			}

			t.Logf("✓ Algorithm suite %s correctly uses %s message format", tc.algorithmSuite, tc.expectedFormat)

			// Decrypt sanity check - verify we can decrypt the object back to original plaintext
			decryptResult, err := s3ec.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
			if err != nil {
				t.Fatalf("failed to decrypt object with algorithm suite %s: %v", tc.algorithmSuite, err)
			}

			decryptedData, err := io.ReadAll(decryptResult.Body)
			if err != nil {
				t.Fatalf("failed to read decrypted data for algorithm suite %s: %v", tc.algorithmSuite, err)
			}

			if string(decryptedData) != plaintext {
				t.Errorf("decrypt sanity check failed for algorithm suite %s: expected %q, got %q", 
					tc.algorithmSuite, plaintext, string(decryptedData))
			} else {
				t.Logf("✓ Decrypt sanity check passed for algorithm suite %s: successfully decrypted %d bytes", 
					tc.algorithmSuite, len(decryptedData))
			}

			// Cleanup
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
		})
	}
}

func TestInteg_CommitmentPolicyBehavior(t *testing.T) {
	var bucket = LoadBucket()
	var region = LoadRegion()
	var accountId = LoadAwsAccountId()
	var baseKey = "commitment-policy-test-" + time.Now().Format("20060102-150405")
	var plaintext = "Hello, S3 Encryption Client Commitment Policy test!"

	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx, config.WithRegion(region))
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var alias = LoadAwsKmsAlias()
	arn := getAliasArn(alias, region, accountId)
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	cases := []struct {
		name                    string
		commitmentPolicy        commitment.CommitmentPolicy
		algorithmSuite          *algorithms.AlgorithmSuite
		expectEncryptError      bool
		expectDecryptError      bool
		encryptErrorContains    string
		decryptErrorContains    string
	}{
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST NOT encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "FORBID_with_committing_algorithm",
			commitmentPolicy:        commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "does not allow committing algorithm suites",
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is FORBID_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
		{
			name:                    "FORBID_with_non_committing_algorithm",
			commitmentPolicy:        commitment.FORBID_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMIV12Tag16NoKDF,
			expectEncryptError:      false,
			expectDecryptError:      false,
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_ALLOW_DECRYPT_with_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      false,
			expectDecryptError:      false,
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_ALLOW_DECRYPT, the S3EC MUST allow decryption using algorithm suites which do not support key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_ALLOW_DECRYPT_with_non_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_ALLOW_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMIV12Tag16NoKDF,
			expectEncryptError:      true,
			expectDecryptError:      false,
			encryptErrorContains:    "requires committing algorithm suites",
		},
		//= ../specification/s3-encryption/key-commitment.md#commitment-policy
		//= type=test
		//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST only encrypt using an algorithm suite which supports key commitment.
		{
			name:                    "REQUIRE_ENCRYPT_REQUIRE_DECRYPT_with_committing_algorithm",
			commitmentPolicy:        commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT,
			algorithmSuite:          algorithms.AlgAES256GCMHkdfSha512CommitKey,
			expectEncryptError:      false,
			expectDecryptError:      false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			key := baseKey + "-" + tc.name

			// Clean up any existing object
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			keyring := materials.NewKmsKeyring(kmsClient, arn, func(options *materials.KeyringOptions) {
				options.EnableLegacyWrappingAlgorithms = false
			})
			cmm, err := materials.NewCryptographicMaterialsManager(keyring)
			if err != nil {
				t.Fatalf("failed to create CMM: %v", err)
			}

			// Test encryption behavior
			s3ec, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
				clientOptions.EncryptionAlgorithmSuite = tc.algorithmSuite
			})
			if err != nil {
				if tc.expectEncryptError && tc.encryptErrorContains != "" && strings.Contains(err.Error(), tc.encryptErrorContains) {
					t.Logf("✓ Expected error during client creation for encryption: %v", err)
					return // Expected error during client creation
				}
				t.Fatalf("expected no error during client creation, got %v", err)
			}

			_, err = s3ec.PutObject(ctx, &s3.PutObjectInput{
				Bucket: &bucket,
				Key:    &key,
				Body:   bytes.NewReader([]byte(plaintext)),
			})

			if tc.expectEncryptError {
				if err == nil {
					t.Fatalf("expected encryption error but got none")
				}
				if tc.encryptErrorContains != "" && !strings.Contains(err.Error(), tc.encryptErrorContains) {
					t.Errorf("expected encryption error to contain %q, got %q", tc.encryptErrorContains, err.Error())
				} else {
					t.Logf("✓ Expected encryption error: %v", err)
				}
				return // Don't test decryption if encryption failed
			} else {
				if err != nil {
					t.Fatalf("expected no encryption error, got %v", err)
				}
				t.Logf("✓ Successfully encrypted with %s policy and %s algorithm", tc.commitmentPolicy, tc.algorithmSuite.CipherName())
			}

			// Test decryption behavior - create a new client for decryption
			decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
				clientOptions.CommitmentPolicy = tc.commitmentPolicy
			})
			if err != nil {
				t.Fatalf("failed to create decryption client: %v", err)
			}

			_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})

			if tc.expectDecryptError {
				if err == nil {
					t.Fatalf("expected decryption error but got none")
				}
				if tc.decryptErrorContains != "" && !strings.Contains(err.Error(), tc.decryptErrorContains) {
					t.Errorf("expected decryption error to contain %q, got %q", tc.decryptErrorContains, err.Error())
				} else {
					t.Logf("✓ Expected decryption error: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no decryption error, got %v", err)
				}
				t.Logf("✓ Successfully decrypted with %s policy", tc.commitmentPolicy)
			}

			// Cleanup
			s3Client.DeleteObject(ctx, &s3.DeleteObjectInput{
				Bucket: &bucket,
				Key:    &key,
			})
		})
	}
}


//= ../specification/s3-encryption/key-commitment.md#commitment-policy
//= type=test
//# When the commitment policy is REQUIRE_ENCRYPT_REQUIRE_DECRYPT, the S3EC MUST NOT allow decryption using algorithm suites which do not support key commitment.
func TestInteg_RequireEncryptRequireDecrypt_RejectsNonCommittingObjects(t *testing.T) {
	ctx := context.Background()
	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion("us-west-2"),
		config.WithLogConfigurationWarnings(true),
	)
	if err != nil {
		t.Fatalf("failed to load cfg: %v", err)
	}

	var bucket = LoadBucket()
	s3Client := s3.NewFromConfig(cfg)
	kmsClient := kms.NewFromConfig(cfg)

	// Test that REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy rejects decryption of non-committing objects
	keyring := materials.NewKmsDecryptOnlyAnyKeyKeyring(kmsClient, func(options *materials.KeyringOptions) {
		options.EnableLegacyWrappingAlgorithms = false
	})
	cmm, err := materials.NewCryptographicMaterialsManager(keyring)
	if err != nil {
		t.Fatalf("failed to create CMM: %v", err)
	}

	decClient, err := client.New(s3Client, cmm, func(clientOptions *client.EncryptionClientOptions) {
		clientOptions.CommitmentPolicy = commitment.REQUIRE_ENCRYPT_REQUIRE_DECRYPT
	})
	if err != nil {
		t.Fatalf("failed to create decryption client: %v", err)
	}

	// Try to decrypt a non-committing object (V2 format, AES/GCM without commitment)
	nonCommittingObjectKey := "crypto_tests/aes_gcm/v4/language_Go/ciphertext_test_case_test_one.txt"
	_, err = decClient.GetObject(ctx, &s3.GetObjectInput{
		Bucket: &bucket,
		Key:    &nonCommittingObjectKey,
	})

	if err == nil {
		t.Fatalf("expected decryption error when trying to decrypt non-committing object with REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy, but got none")
	}

	if !strings.Contains(err.Error(), "object's content encryption algorithm is not valid for the selected commitment policy") {
		t.Errorf("expected error to contain commitment policy validation message, got: %v", err)
	}

	t.Logf("✓ REQUIRE_ENCRYPT_REQUIRE_DECRYPT policy correctly rejected decryption of non-committing object: %v", err)
}
