# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go tool for analyzing digital signatures in PE (Portable Executable) files. The project extracts PKCS#7 digital signatures from signed Windows executables and can verify their validity using Mozilla's PKCS#7 library.

## Commands

### Building and Running
```bash
# Build the CLI tool
go build -o gosigtool cmd/gosigtool/main.go

# Install from source
go install ./cmd/gosigtool

# Run the tool
./gosigtool -in <signed_pe_file> [-out <output_file>] [-validate]
```

### Testing
```bash
# Run all tests (unit + integration)
go test -v ./...

# Run only unit tests (fast)
go test -short -v ./...

# Run only integration tests
go test -v -run Integration

# Run tests with benchmarks
go test -bench=. -v

# Run tests with coverage
go test -cover ./...

# Setup integration test files
./scripts/setup-test-files.sh

# Run with custom PE file
SIGTOOL_TEST_PE_FILE=/path/to/signed.exe go test -v -run Integration
```

### Dependency Management
```bash
# Download dependencies
go mod download

# Clean up dependencies
go mod tidy
```

## Architecture

The project consists of two main components:

1. **Core Library (`sigtool.go`)**: Contains the main functionality
   - `ExtractDigitalSignature()`: Extracts PKCS#7 signature from PE files by reading the security directory
   - `IsValidDigitalSignature()`: Validates the extracted signature using PKCS#7 verification

2. **CLI Tool (`cmd/gosigtool/main.go`)**: Command-line interface
   - Handles flag parsing for input/output files and validation options
   - Outputs extracted signatures to `.pkcs7` files
   - Optional signature verification with `-validate` flag

## Key Implementation Details

- Uses Go's `debug/pe` package to parse PE file headers and locate the security directory
- Supports both 32-bit and 64-bit PE files through optional header type switching
- Reads signature data from the file offset specified in the security directory (skipping 8-byte header)
- Relies on `go.mozilla.org/pkcs7` for signature parsing and verification
- Default output naming: `<input_filename>.pkcs7`

## Testing

The project includes comprehensive testing at multiple levels:

### Unit Tests
- **Happy path scenarios**: Valid signed PE files with proper signature extraction
- **Error conditions**: Empty file paths, non-existent files, unsigned PE files, corrupted data
- **Edge cases**: Oversized signatures, malformed PE files, bounds checking
- **Security validation**: PKCS#7 signature parsing and verification testing
- Uses mock PE file generation to create minimal but valid PE structures

### Integration Tests
- **Real PE file testing**: Tests against actual signed executables
- **Automatic test file setup**: Downloads known good test files via `scripts/setup-test-files.sh`
- **Multiple file sources**: Supports Windows system files, environment variables, or custom files
- **CI/CD friendly**: Gracefully skips when no test files available
- **Benchmark testing**: Performance testing with real-world files

To set up integration tests:
```bash
./scripts/setup-test-files.sh  # Downloads a sample signed PE file
go test -v -run Integration     # Run integration tests
```

## CI/CD Pipeline

The project uses GitHub Actions for automated testing and quality assurance:

### Automatic Triggers
- **Push to master**: Runs full test suite on all supported platforms
- **Pull requests**: Validates changes before merging

### Manual Triggers (Repository Owners/Collaborators Only)
- **Workflow Dispatch**: Manual execution with options for test type and Go version
- **Access Control**: Only users with write/admin/maintain permissions can trigger manually

### Test Matrix
- **Platforms**: Ubuntu, Windows, macOS
- **Go Versions**: 1.21, 1.22
- **Test Types**: Unit tests (fast), Integration tests (with real PE files)

### Quality Gates
- **Unit Tests**: Fast tests with race detection and coverage
- **Integration Tests**: Real PE file testing with benchmarks
- **Linting**: golangci-lint, go vet, gofmt, security scanning
- **Build Verification**: Cross-platform binary builds

### Coverage and Artifacts
- **Code Coverage**: Uploaded to Codecov
- **Build Artifacts**: Binaries for all platforms (30-day retention)
- **Test Files**: Integration test files (7-day retention)

## Dependencies

- `go.mozilla.org/pkcs7 v0.9.0`: For PKCS#7 signature parsing and verification
- Go standard library packages: `debug/pe`, `os`, `flag`, `fmt`, `log`, `path/filepath`