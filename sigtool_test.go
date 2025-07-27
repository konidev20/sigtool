package sigtool

import (
	"encoding/binary"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// createMockPEFile creates a minimal PE file with optional security directory
func createMockPEFile(t *testing.T, withSignature bool, signatureData []byte) string {
	t.Helper()

	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "test.exe")

	// Calculate total file size upfront
	headerSize := 64 + 4 + 20 + 224 // DOS + PE sig + COFF + Optional
	var totalSize int64 = int64(headerSize)
	if withSignature {
		totalSize += 8 + int64(len(signatureData)) // Security header + signature data
	}

	// Create file content in memory first
	content := make([]byte, totalSize)

	// Write DOS header
	copy(content[0:2], "MZ")                        // DOS signature
	binary.LittleEndian.PutUint32(content[60:], 64) // PE header offset

	// Write PE signature
	copy(content[64:68], "PE\x00\x00")

	// Write COFF header
	binary.LittleEndian.PutUint16(content[68:], 0x014c) // Machine (i386)
	binary.LittleEndian.PutUint16(content[70:], 0)      // NumberOfSections
	binary.LittleEndian.PutUint16(content[84:], 224)    // SizeOfOptionalHeader

	// Write Optional Header
	optHeaderStart := 88
	binary.LittleEndian.PutUint16(content[optHeaderStart:], 0x010b) // Magic (PE32)
	binary.LittleEndian.PutUint32(content[optHeaderStart+92:], 16)  // NumberOfRvaAndSizes

	if withSignature {
		// Security directory is at index 4, each entry is 8 bytes
		secDirOffset := optHeaderStart + 96 + 4*8
		signatureOffset := uint32(headerSize)

		// Write security directory entry
		binary.LittleEndian.PutUint32(content[secDirOffset:], signatureOffset)
		binary.LittleEndian.PutUint32(content[secDirOffset+4:], uint32(len(signatureData)+8))

		// Write security directory header at end
		secHeaderStart := headerSize
		binary.LittleEndian.PutUint32(content[secHeaderStart:], uint32(len(signatureData)+8)) // Length
		binary.LittleEndian.PutUint16(content[secHeaderStart+4:], 0x0200)                     // Revision
		binary.LittleEndian.PutUint16(content[secHeaderStart+6:], 0x0002)                     // Type

		// Write signature data
		copy(content[secHeaderStart+8:], signatureData)
	}

	// Write entire content to file
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	return filePath
}

func TestExtractDigitalSignature_ValidSignedPE(t *testing.T) {
	signatureData := []byte("mock-pkcs7-signature-data")
	filePath := createMockPEFile(t, true, signatureData)

	result, err := ExtractDigitalSignature(filePath)
	if err != nil {
		t.Fatalf("Expected no error, got: %v", err)
	}

	if string(result) != string(signatureData) {
		t.Errorf("Expected signature data %q, got %q", signatureData, result)
	}
}

func TestExtractDigitalSignature_UnsignedPE(t *testing.T) {
	filePath := createMockPEFile(t, false, nil)

	_, err := ExtractDigitalSignature(filePath)
	if err == nil {
		t.Fatal("Expected error for unsigned PE file, got nil")
	}

	if !strings.Contains(err.Error(), "not digitally signed") {
		t.Errorf("Expected 'not digitally signed' error, got: %v", err)
	}
}

func TestExtractDigitalSignature_EmptyFilePath(t *testing.T) {
	testCases := []string{"", "   ", "\t\n"}

	for _, filePath := range testCases {
		_, err := ExtractDigitalSignature(filePath)
		if err == nil {
			t.Errorf("Expected error for empty file path %q, got nil", filePath)
		}

		if !strings.Contains(err.Error(), "cannot be empty") {
			t.Errorf("Expected 'cannot be empty' error, got: %v", err)
		}
	}
}

func TestExtractDigitalSignature_NonExistentFile(t *testing.T) {
	_, err := ExtractDigitalSignature("/nonexistent/file.exe")
	if err == nil {
		t.Fatal("Expected error for non-existent file, got nil")
	}

	if !strings.Contains(err.Error(), "failed to open file") {
		t.Errorf("Expected 'failed to open file' error, got: %v", err)
	}
}

func TestExtractDigitalSignature_NonPEFile(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "notpe.txt")

	if err := os.WriteFile(filePath, []byte("not a PE file"), 0600); err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	_, err := ExtractDigitalSignature(filePath)
	if err == nil {
		t.Fatal("Expected error for non-PE file, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse PE file") {
		t.Errorf("Expected 'failed to parse PE file' error, got: %v", err)
	}
}

func TestExtractDigitalSignature_LargeSignature(t *testing.T) {
	// Create signature larger than MaxSignatureSize
	largeSignature := make([]byte, MaxSignatureSize+1)
	filePath := createMockPEFile(t, true, largeSignature)

	_, err := ExtractDigitalSignature(filePath)
	if err == nil {
		t.Fatal("Expected error for oversized signature, got nil")
	}

	if !strings.Contains(err.Error(), "exceeds maximum allowed size") {
		t.Errorf("Expected 'exceeds maximum allowed size' error, got: %v", err)
	}
}

func TestIsValidDigitalSignature_EmptyFilePath(t *testing.T) {
	err := IsValidDigitalSignature("")
	if err == nil {
		t.Fatal("Expected error for empty file path, got nil")
	}

	if !strings.Contains(err.Error(), "cannot be empty") {
		t.Errorf("Expected 'cannot be empty' error, got: %v", err)
	}
}

func TestIsValidDigitalSignature_InvalidSignature(t *testing.T) {
	// Create a PE file with invalid PKCS#7 data
	invalidSignature := []byte("invalid-pkcs7-data")
	filePath := createMockPEFile(t, true, invalidSignature)

	err := IsValidDigitalSignature(filePath)
	if err == nil {
		t.Fatal("Expected error for invalid signature, got nil")
	}

	if !strings.Contains(err.Error(), "failed to parse PKCS#7") {
		t.Errorf("Expected 'failed to parse PKCS#7' error, got: %v", err)
	}
}

func TestIsValidDigitalSignature_ExtractionFailure(t *testing.T) {
	// Test with unsigned PE file
	filePath := createMockPEFile(t, false, nil)

	err := IsValidDigitalSignature(filePath)
	if err == nil {
		t.Fatal("Expected error for unsigned PE file, got nil")
	}

	if !strings.Contains(err.Error(), "failed to extract signature") {
		t.Errorf("Expected 'failed to extract signature' error, got: %v", err)
	}
}

// Benchmark tests
func BenchmarkExtractDigitalSignature(b *testing.B) {
	signatureData := make([]byte, 1024) // 1KB signature
	filePath := createMockPEFileForBench(b, true, signatureData)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ExtractDigitalSignature(filePath)
		if err != nil {
			b.Fatalf("Unexpected error: %v", err)
		}
	}
}

// Helper function for benchmark
func createMockPEFileForBench(b *testing.B, withSignature bool, signatureData []byte) string {
	b.Helper()

	tmpDir := b.TempDir()
	filePath := filepath.Join(tmpDir, "test.exe")

	// Calculate total file size upfront
	headerSize := 64 + 4 + 20 + 224 // DOS + PE sig + COFF + Optional
	var totalSize int64 = int64(headerSize)
	if withSignature {
		totalSize += 8 + int64(len(signatureData)) // Security header + signature data
	}

	// Create file content in memory first
	content := make([]byte, totalSize)

	// Write DOS header
	copy(content[0:2], "MZ")                        // DOS signature
	binary.LittleEndian.PutUint32(content[60:], 64) // PE header offset

	// Write PE signature
	copy(content[64:68], "PE\x00\x00")

	// Write COFF header
	binary.LittleEndian.PutUint16(content[68:], 0x014c) // Machine (i386)
	binary.LittleEndian.PutUint16(content[70:], 0)      // NumberOfSections
	binary.LittleEndian.PutUint16(content[84:], 224)    // SizeOfOptionalHeader

	// Write Optional Header
	optHeaderStart := 88
	binary.LittleEndian.PutUint16(content[optHeaderStart:], 0x010b) // Magic (PE32)
	binary.LittleEndian.PutUint32(content[optHeaderStart+92:], 16)  // NumberOfRvaAndSizes

	if withSignature {
		// Security directory is at index 4, each entry is 8 bytes
		secDirOffset := optHeaderStart + 96 + 4*8
		signatureOffset := uint32(headerSize)

		// Write security directory entry
		binary.LittleEndian.PutUint32(content[secDirOffset:], signatureOffset)
		binary.LittleEndian.PutUint32(content[secDirOffset+4:], uint32(len(signatureData)+8))

		// Write security directory header at end
		secHeaderStart := headerSize
		binary.LittleEndian.PutUint32(content[secHeaderStart:], uint32(len(signatureData)+8)) // Length
		binary.LittleEndian.PutUint16(content[secHeaderStart+4:], 0x0200)                     // Revision
		binary.LittleEndian.PutUint16(content[secHeaderStart+6:], 0x0002)                     // Type

		// Write signature data
		copy(content[secHeaderStart+8:], signatureData)
	}

	// Write entire content to file
	if err := os.WriteFile(filePath, content, 0600); err != nil {
		b.Fatalf("Failed to write test file: %v", err)
	}

	return filePath
}
