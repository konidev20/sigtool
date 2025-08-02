# sigtool

A Go library and CLI tool for extracting and validating digital signatures from Windows PE (Portable Executable) files.

## Features

- **Extract digital signatures** from signed PE files (PKCS#7 format)
- **Validate digital signatures** using certificate chain verification
- **Cross-platform support** (Windows, Linux, macOS)
- **CLI tool** for command-line operations
- **Go library** for programmatic integration
- **Comprehensive testing** with unit and integration tests

## Installation

### As a Go Module

```bash
go get github.com/konidev20/sigtool
```

### CLI Tool

```bash
go install github.com/konidev20/sigtool/cmd/gosigtool@latest
```

## Usage

### CLI Tool

Extract a digital signature from a PE file:

```bash
gosigtool path/to/signed.exe output_signature.pkcs7
```

The tool will:
- Extract the PKCS#7 signature from the PE file
- Save it to the specified output file
- Validate the signature and report the result

### Go Library

```go
package main

import (
    "fmt"
    "log"
    
    "github.com/konidev20/sigtool"
)

func main() {
    // Extract digital signature
    signature, err := sigtool.ExtractDigitalSignature("path/to/signed.exe")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Extracted signature: %d bytes\n", len(signature))
    
    // Validate the signature
    err = sigtool.IsValidDigitalSignature("path/to/signed.exe")
    if err != nil {
        fmt.Printf("Signature validation failed: %v\n", err)
    } else {
        fmt.Println("Signature is valid")
    }
}
```

## API Reference

### Functions

#### `ExtractDigitalSignature(filePath string) ([]byte, error)`

Extracts the PKCS#7 digital signature from a signed PE file.

**Parameters:**
- `filePath`: Path to the PE file

**Returns:**
- `[]byte`: Raw PKCS#7 signature data
- `error`: Error if extraction fails

**Errors:**
- File not found or cannot be opened
- Invalid PE file format
- File is not digitally signed
- Signature data is corrupted or invalid

#### `IsValidDigitalSignature(filePath string) error`

Validates the digital signature of a PE file using PKCS#7 verification.

**Parameters:**
- `filePath`: Path to the PE file

**Returns:**
- `error`: `nil` if signature is valid, error describing validation failure otherwise

**Note:** Validation may fail due to expired certificates, missing root certificates, or untrusted certificate chains, even if the signature format is correct.

## Requirements

- Go 1.21 or higher
- No external dependencies beyond the Go standard library and `go.mozilla.org/pkcs7`

## Development

### Building

```bash
# Build the library
go build ./...

# Build the CLI tool
go build -o gosigtool ./cmd/gosigtool
```

### Testing

```bash
# Run unit tests
go test -short ./...

# Run all tests including integration tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Integration Tests

Integration tests require signed PE files. The test suite will automatically search for:
- Windows system files (if available)
- Files in `testdata/` directory
- Custom files via `SIGTOOL_TEST_PE_FILE` environment variable

```bash
# Run with custom test file
SIGTOOL_TEST_PE_FILE=/path/to/signed.exe go test ./...
```

## CI/CD

The project includes comprehensive GitHub Actions workflows:
- **Unit tests** on multiple Go versions and platforms
- **Integration tests** with real PE files
- **Code quality** checks (golangci-lint, gosec, gofmt)
- **Security scanning** with Gosec
- **Build verification** for all supported platforms
- **Coverage reporting** via Codecov

## Security Considerations

- The tool is designed for defensive security analysis only
- File access is intentionally limited to user-specified files
- Input validation prevents buffer overflows and path traversal
- Maximum signature size limits prevent memory exhaustion
- All file operations include proper bounds checking

## Contributing

1. Fork the repository
2. Create a feature branch
3. Add tests for new functionality
4. Ensure all tests pass and linting is clean
5. Submit a pull request

## License

This project is licensed under the MIT License - see the LICENSE file for details.

## Acknowledgments

- Uses the [go.mozilla.org/pkcs7](https://github.com/mozilla-services/pkcs7) library for PKCS#7 parsing and verification
- Built with Go's excellent `debug/pe` package for PE file parsing

## Support

For bugs, feature requests, or questions:
- Open an issue on GitHub
- Ensure you provide sample files (if safe to share) and full error messages
- Include your Go version and operating system details
