// Package sigtool provides functionality for extracting and validating digital signatures
// from Windows PE (Portable Executable) files.
//
// This package supports reading PKCS#7 signatures embedded in PE files and performing
// cryptographic validation of those signatures. It is designed for defensive security
// analysis and forensic examination of signed executables.
//
// Key features:
//   - Extract PKCS#7 digital signatures from signed PE files
//   - Validate signatures using certificate chain verification
//   - Cross-platform support (Windows, Linux, macOS)
//   - Comprehensive error handling and input validation
//   - Security-focused design with bounds checking
//
// Example usage:
//
//	// Extract a signature
//	signature, err := sigtool.ExtractDigitalSignature("signed.exe")
//	if err != nil {
//	    log.Fatal(err)
//	}
//
//	// Validate the signature
//	err = sigtool.IsValidDigitalSignature("signed.exe")
//	if err != nil {
//	    fmt.Printf("Validation failed: %v\n", err)
//	}
package sigtool

import (
	"debug/pe"
	"errors"
	"fmt"
	"os"
	"strings"

	"go.mozilla.org/pkcs7"
)

const (
	// PKCS#7 signature data starts after 8-byte security directory header
	SecurityDirHeaderSize = 8
	// Maximum reasonable signature size (10MB)
	MaxSignatureSize = 10 * 1024 * 1024
)

// ExtractDigitalSignature extracts the PKCS#7 digital signature from a signed PE file.
//
// This function parses the PE file structure and locates the security directory
// containing the digital signature. It performs comprehensive validation including
// bounds checking and file integrity verification.
//
// Parameters:
//   - filePath: The path to the PE file to extract the signature from
//
// Returns:
//   - []byte: The raw PKCS#7 signature data
//   - error: An error if extraction fails for any reason
//
// Common errors include:
//   - File not found or inaccessible
//   - Invalid PE file format
//   - File is not digitally signed
//   - Signature data is corrupted or extends beyond file bounds
//
// Example usage:
//
//	signature, err := sigtool.ExtractDigitalSignature("signed.exe")
//	if err != nil {
//	    log.Fatal(err)
//	}
//	fmt.Printf("Extracted %d bytes of signature data\n", len(signature))
func ExtractDigitalSignature(filePath string) (buf []byte, err error) {
	// Input validation
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("file path cannot be empty")
	}

	// Open file once and use for both PE parsing and signature extraction
	// #nosec G304 - This tool is designed to read user-specified PE files
	f, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file %q: %w", filePath, err)
	}
	defer f.Close()

	// Get file info for bounds checking
	fileInfo, err := f.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file info: %w", err)
	}
	fileSize := fileInfo.Size()

	// Parse PE file
	pefile, err := pe.NewFile(f)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PE file: %w", err)
	}
	defer pefile.Close()

	var vAddr uint32
	var size uint32
	switch t := pefile.OptionalHeader.(type) {
	case *pe.OptionalHeader32:
		vAddr = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].VirtualAddress
		size = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].Size
	case *pe.OptionalHeader64:
		vAddr = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].VirtualAddress
		size = t.DataDirectory[pe.IMAGE_DIRECTORY_ENTRY_SECURITY].Size
	default:
		return nil, errors.New("unsupported PE optional header type")
	}

	// Validate security directory
	if vAddr == 0 || size == 0 {
		return nil, errors.New("PE file is not digitally signed")
	}

	// Bounds checking
	if size > MaxSignatureSize {
		return nil, fmt.Errorf("signature size %d exceeds maximum allowed size %d", size, MaxSignatureSize)
	}

	// Calculate actual signature data size (excluding 8-byte header)
	signatureDataSize := size - SecurityDirHeaderSize
	signatureOffset := int64(vAddr + SecurityDirHeaderSize)

	if signatureOffset < 0 || signatureOffset >= fileSize {
		return nil, fmt.Errorf("invalid signature offset %d in file of size %d", signatureOffset, fileSize)
	}

	if signatureOffset+int64(signatureDataSize) > fileSize {
		return nil, fmt.Errorf("signature extends beyond file bounds")
	}

	// Read signature data (excluding the 8-byte security directory header)
	buf = make([]byte, signatureDataSize)
	n, err := f.ReadAt(buf, signatureOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature data: %w", err)
	}
	if n != int(signatureDataSize) {
		return nil, fmt.Errorf("incomplete read: expected %d bytes, got %d", signatureDataSize, n)
	}

	return buf, nil
}

// IsValidDigitalSignature validates the digital signature of a PE file using PKCS#7 verification.
//
// This function extracts the signature from the PE file and performs cryptographic
// verification including certificate chain validation. Note that validation may fail
// even for properly formatted signatures due to certificate trust issues.
//
// Parameters:
//   - filePath: The path to the PE file to validate
//
// Returns:
//   - error: nil if the signature is valid, otherwise an error describing the validation failure
//
// Common validation failures:
//   - Expired certificates
//   - Missing or untrusted root certificates
//   - Revoked certificates
//   - Invalid signature format or corruption
//   - Certificate chain verification failures
//
// Example usage:
//
//	err := sigtool.IsValidDigitalSignature("signed.exe")
//	if err != nil {
//	    fmt.Printf("Signature validation failed: %v\n", err)
//	} else {
//	    fmt.Println("Signature is valid")
//	}
func IsValidDigitalSignature(filePath string) (err error) {
	// Input validation
	if strings.TrimSpace(filePath) == "" {
		return errors.New("file path cannot be empty")
	}

	peExtract, err := ExtractDigitalSignature(filePath)
	if err != nil {
		return fmt.Errorf("failed to extract signature: %w", err)
	}

	pc, err := pkcs7.Parse(peExtract)
	if err != nil {
		return fmt.Errorf("failed to parse PKCS#7 signature: %w", err)
	}

	if err := pc.Verify(); err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}
