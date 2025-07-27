#!/bin/bash

# Script to download a known signed PE file for integration testing

set -e

TESTDATA_DIR="$(dirname "$0")/../testdata"
mkdir -p "$TESTDATA_DIR"

echo "Setting up integration test files..."

# Option 1: Try to copy from Windows system (if available)
if [[ "$OSTYPE" == "msys" || "$OSTYPE" == "cygwin" || "$OSTYPE" == "win32" ]]; then
    echo "Windows detected. Trying to copy system files..."
    
    if [[ -f "C:/Windows/System32/notepad.exe" ]]; then
        cp "C:/Windows/System32/notepad.exe" "$TESTDATA_DIR/signed.exe"
        echo "✅ Copied notepad.exe as test file"
        exit 0
    fi
fi

# Option 2: Download a small, reliable signed PE file
echo "Downloading a sample signed PE file..."

# We'll use a small utility from Microsoft - sigcheck from Sysinternals
# This is a legitimate, small, signed executable that's perfect for testing
DOWNLOAD_URL="https://download.sysinternals.com/files/Sigcheck.zip"
TEMP_ZIP="$TESTDATA_DIR/sigcheck.zip"

if command -v curl >/dev/null 2>&1; then
    curl -L -o "$TEMP_ZIP" "$DOWNLOAD_URL"
elif command -v wget >/dev/null 2>&1; then
    wget -O "$TEMP_ZIP" "$DOWNLOAD_URL"
else
    echo "❌ Neither curl nor wget available. Please install one of them or manually place a signed PE file in testdata/signed.exe"
    exit 1
fi

# Extract the executable
if command -v unzip >/dev/null 2>&1; then
    cd "$TESTDATA_DIR"
    unzip -q sigcheck.zip
    
    # Find the sigcheck executable (could be sigcheck.exe or sigcheck64.exe)
    if [[ -f "sigcheck.exe" ]]; then
        mv sigcheck.exe signed.exe
        echo "✅ Downloaded and set up sigcheck.exe as test file"
    elif [[ -f "sigcheck64.exe" ]]; then
        mv sigcheck64.exe signed.exe
        echo "✅ Downloaded and set up sigcheck64.exe as test file"
    else
        echo "❌ Could not find sigcheck executable in downloaded zip"
        exit 1
    fi
    
    # Clean up
    rm -f sigcheck.zip sigcheck64.exe sigcheck.exe Eula.txt 2>/dev/null || true
    
else
    echo "❌ unzip not available. Please install unzip or manually extract and place a signed PE file in testdata/signed.exe"
    exit 1
fi

echo ""
echo "🎉 Test file setup complete!"
echo "You can now run integration tests with:"
echo "  go test -v -run Integration"