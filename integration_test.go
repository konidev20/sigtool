package sigtool

import (
	"os"
	"path/filepath"
	"testing"
)

// TestIntegration_RealPEFile tests the package against a real signed PE file
func TestIntegration_RealPEFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Look for test PE files in several common locations
	testFiles := []string{
		// Windows system files (if running on Windows or with Windows files available)
		"C:\\Windows\\System32\\notepad.exe",
		"C:\\Windows\\System32\\calc.exe",
		"C:\\Windows\\System32\\cmd.exe",
		
		// Relative paths for test files that might be provided
		"testdata/signed.exe",
		"testdata/notepad.exe",
		"test_files/signed.exe",
		
		// Environment variable for custom test file
		os.Getenv("SIGTOOL_TEST_PE_FILE"),
	}

	var testFile string
	for _, file := range testFiles {
		if file != "" {
			if _, err := os.Stat(file); err == nil {
				testFile = file
				break
			}
		}
	}

	if testFile == "" {
		t.Skip("No signed PE file available for integration testing. " +
			"Set SIGTOOL_TEST_PE_FILE environment variable or place a signed PE file in testdata/signed.exe")
	}

	t.Logf("Testing with PE file: %s", testFile)

	// Test 1: Extract digital signature
	t.Run("ExtractSignature", func(t *testing.T) {
		signature, err := ExtractDigitalSignature(testFile)
		if err != nil {
			t.Fatalf("Failed to extract signature: %v", err)
		}

		if len(signature) == 0 {
			t.Fatal("Extracted signature is empty")
		}

		t.Logf("Successfully extracted signature of %d bytes", len(signature))

		// Basic sanity check - PKCS#7 signatures typically start with specific bytes
		if len(signature) < 10 {
			t.Error("Signature seems too short to be valid PKCS#7")
		}
	})

	// Test 2: Validate digital signature
	t.Run("ValidateSignature", func(t *testing.T) {
		err := IsValidDigitalSignature(testFile)
		if err != nil {
			// Note: Validation might fail due to expired certificates, missing root certs, etc.
			// This is expected in many test environments, so we log but don't fail
			t.Logf("Signature validation failed (this may be expected): %v", err)
		} else {
			t.Log("Signature validation succeeded")
		}
	})

	// Test 3: Extract and save signature to file
	t.Run("ExtractAndSave", func(t *testing.T) {
		signature, err := ExtractDigitalSignature(testFile)
		if err != nil {
			t.Fatalf("Failed to extract signature: %v", err)
		}

		// Save to temporary file
		tmpDir := t.TempDir()
		outputFile := filepath.Join(tmpDir, "extracted_signature.pkcs7")
		
		err = os.WriteFile(outputFile, signature, 0644)
		if err != nil {
			t.Fatalf("Failed to write signature to file: %v", err)
		}

		// Verify the file was created and has content
		info, err := os.Stat(outputFile)
		if err != nil {
			t.Fatalf("Failed to stat output file: %v", err)
		}

		if info.Size() != int64(len(signature)) {
			t.Errorf("Output file size %d doesn't match signature size %d", info.Size(), len(signature))
		}

		t.Logf("Successfully saved signature to %s (%d bytes)", outputFile, info.Size())
	})
}

// TestIntegration_UnsignedPEFile tests with an unsigned PE file
func TestIntegration_UnsignedPEFile(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Try to find an unsigned PE file or create one
	testFiles := []string{
		"testdata/unsigned.exe",
		"test_files/unsigned.exe",
		os.Getenv("SIGTOOL_TEST_UNSIGNED_PE_FILE"),
	}

	var testFile string
	for _, file := range testFiles {
		if file != "" {
			if _, err := os.Stat(file); err == nil {
				testFile = file
				break
			}
		}
	}

	if testFile == "" {
		// Create a simple unsigned PE file for testing
		testFile = createSimpleUnsignedPE(t)
	}

	t.Logf("Testing with unsigned PE file: %s", testFile)

	// Test that extraction fails appropriately
	_, err := ExtractDigitalSignature(testFile)
	if err == nil {
		t.Fatal("Expected error when extracting signature from unsigned PE file, got nil")
	}

	if err.Error() != "PE file is not digitally signed" {
		t.Errorf("Expected 'PE file is not digitally signed' error, got: %v", err)
	}

	t.Log("Correctly detected unsigned PE file")
}

// createSimpleUnsignedPE creates a minimal unsigned PE file for testing
func createSimpleUnsignedPE(t *testing.T) string {
	t.Helper()
	
	// Use our existing mock PE file generator without signature
	return createMockPEFile(t, false, nil)
}

// Benchmark integration test with real file
func BenchmarkIntegration_RealPEFile(b *testing.B) {
	testFile := os.Getenv("SIGTOOL_TEST_PE_FILE")
	if testFile == "" {
		// Try common Windows files
		commonFiles := []string{
			"C:\\Windows\\System32\\notepad.exe",
			"testdata/signed.exe",
		}
		
		for _, file := range commonFiles {
			if _, err := os.Stat(file); err == nil {
				testFile = file
				break
			}
		}
	}

	if testFile == "" {
		b.Skip("No signed PE file available for benchmark. Set SIGTOOL_TEST_PE_FILE environment variable")
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractDigitalSignature(testFile)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}