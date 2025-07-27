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

// ExtractDigitalSignature extracts a digital signature specified in a signed PE file.
// It returns a digital signature (pkcs#7) in bytes.
func ExtractDigitalSignature(filePath string) (buf []byte, err error) {
	// Input validation
	if strings.TrimSpace(filePath) == "" {
		return nil, errors.New("file path cannot be empty")
	}

	// Open file once and use for both PE parsing and signature extraction
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

	signatureOffset := int64(vAddr + SecurityDirHeaderSize)
	if signatureOffset < 0 || signatureOffset >= fileSize {
		return nil, fmt.Errorf("invalid signature offset %d in file of size %d", signatureOffset, fileSize)
	}

	if signatureOffset+int64(size) > fileSize {
		return nil, fmt.Errorf("signature extends beyond file bounds")
	}

	// Read signature data
	buf = make([]byte, size)
	n, err := f.ReadAt(buf, signatureOffset)
	if err != nil {
		return nil, fmt.Errorf("failed to read signature data: %w", err)
	}
	if n != int(size) {
		return nil, fmt.Errorf("incomplete read: expected %d bytes, got %d", size, n)
	}

	return buf, nil
}

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
