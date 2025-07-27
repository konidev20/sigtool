# Test Data

This directory contains test files for integration testing.

## Getting Test Files

### Option 1: Use Environment Variable
Set the `SIGTOOL_TEST_PE_FILE` environment variable to point to any signed PE file:

```bash
export SIGTOOL_TEST_PE_FILE=/path/to/your/signed.exe
go test -v
```

### Option 2: Use Windows System Files (Windows only)
If running on Windows, the tests will automatically use common signed system files like:
- `C:\Windows\System32\notepad.exe`
- `C:\Windows\System32\calc.exe`
- `C:\Windows\System32\cmd.exe`

### Option 3: Download Sample Files
You can download known signed PE files for testing:

```bash
# Example: Download Git for Windows installer (signed PE file)
curl -L -o testdata/signed.exe "https://github.com/git-for-windows/git/releases/download/v2.43.0.windows.1/Git-2.43.0-64-bit.exe"

# Or use any other signed Windows executable
```

### Option 4: Use Your Own Files
Place any signed PE file in this directory as `signed.exe`:

```bash
cp /path/to/your/signed.exe testdata/signed.exe
```

## Running Integration Tests

```bash
# Run all tests including integration tests
go test -v

# Run only integration tests
go test -v -run Integration

# Skip integration tests (fast unit tests only)
go test -short

# Run with your own test file
SIGTOOL_TEST_PE_FILE=/path/to/file.exe go test -v -run Integration
```

## Good Test Files

Some reliable sources for signed PE files:

1. **Windows System Files** (if available):
   - `notepad.exe`, `calc.exe`, `cmd.exe` from `C:\Windows\System32\`

2. **Popular Software Installers**:
   - Git for Windows: https://git-scm.com/download/win
   - Chrome installer: https://www.google.com/chrome/
   - Firefox installer: https://www.mozilla.org/firefox/

3. **Microsoft Tools**:
   - Visual Studio Code: https://code.visualstudio.com/
   - .NET Runtime installers: https://dotnet.microsoft.com/

**Note**: Always verify downloaded files are from trusted sources before testing.