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
# Run tests
go test ./...

# Run tests with verbose output
go test -v ./...
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

## Dependencies

- `go.mozilla.org/pkcs7 v0.9.0`: For PKCS#7 signature parsing and verification
- Go standard library packages: `debug/pe`, `os`, `flag`, `fmt`, `log`, `path/filepath`