# Security Policy

## Supported Versions

We release patches for security vulnerabilities for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 2.0.x   | :white_check_mark: |
| < 2.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability, please follow these steps:

### DO NOT

- Open a public GitHub issue
- Discuss the vulnerability publicly

### DO

1. Email details to: security@vex.dev (or create a private security advisory on GitHub)
2. Include:
   - Description of the vulnerability
   - Steps to reproduce
   - Potential impact
   - Suggested fix (if any)

### What to Expect

- **Acknowledgment**: Within 48 hours
- **Initial Assessment**: Within 7 days
- **Regular Updates**: At least weekly until resolved
- **Fix Timeline**: Varies by severity
  - Critical: 1-7 days
  - High: 7-30 days
  - Medium: 30-90 days
  - Low: 90+ days

## Security Best Practices for Users

1. **Keep Updated**: Always use the latest version
2. **Verify Downloads**: Check checksums for binary downloads
3. **File Permissions**: Be cautious when opening files from untrusted sources
4. **Sandboxing**: Run in a sandboxed environment when processing sensitive data

## Known Security Considerations

### File Parsing

- Excel files can contain macros (not executed by this tool)
- CSV files are parsed as text only
- Formula evaluation is display-only, not executed

### Data Privacy

- All processing happens locally
- No data is sent to external servers
- Clipboard access is used only when explicitly requested

### Dependencies

- We regularly update dependencies
- Security advisories are monitored
- Dependency scanning is automated

## Security Features

1. **No Code Execution**: Formulas are displayed but never executed
2. **Read-Only by Default**: Files are opened in read-only mode
3. **Safe Parsing**: Uses well-maintained libraries (Excelize)
4. **Input Validation**: All user inputs are validated and sanitized
5. **Error Handling**: Graceful error handling prevents crashes

## Disclosure Policy

When a security vulnerability is confirmed:

1. We will develop a fix
2. We will prepare a security advisory
3. We will notify affected users
4. We will publish the fix and advisory simultaneously
5. We will credit the reporter (unless anonymity is requested)

## Compliance

This project follows:

- OWASP guidelines for secure coding
- Go security best practices
- CWE/SANS Top 25 mitigation strategies

## Contact

For security concerns, contact:

- Email: security@vex-tui.dev
- GitHub Security Advisory: [Create Advisory](https://github.com/odesaur/vex-tui/security/advisories/new)

---

Thank you for helping keep Excel TUI secure!
