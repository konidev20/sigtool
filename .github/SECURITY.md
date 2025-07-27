# Security Policy

## Supported Versions

We actively support the following versions of sigtool with security updates:

| Version | Supported          |
| ------- | ------------------ |
| main    | :white_check_mark: |
| Latest release | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability in sigtool, please report it responsibly:

### For Critical Security Issues

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please:

1. **Email**: Send details to the repository maintainers privately
2. **Include**: 
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if available)

### Response Timeline

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 1 week
- **Status Updates**: Weekly until resolved
- **Fix Release**: Based on severity (critical issues within 7 days)

### Security Best Practices

When using sigtool:

1. **Input Validation**: Always validate PE file sources before processing
2. **File Permissions**: Use appropriate file permissions for extracted signatures
3. **Network Security**: When downloading test files, verify sources and checksums
4. **Environment**: Keep Go runtime and dependencies updated

### Scope

This security policy covers:

- The core sigtool library (`sigtool.go`)
- The CLI tool (`cmd/gosigtool/`)
- Build and CI/CD processes
- Dependencies and their security updates

### Out of Scope

- Security of the PE files being analyzed (sigtool analyzes but doesn't validate the security of the target files themselves)
- Third-party tools and utilities not part of this repository
- Security of development environments and systems