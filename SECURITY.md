# Security Policy

## 🔒 Supported Versions

We actively support the following versions of Postie with security updates:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | ✅ Yes             |
| < 1.0   | ❌ No              |

## 🚨 Reporting a Vulnerability

We take the security of Postie seriously. If you discover a security vulnerability, please follow these steps:

### 📧 Private Disclosure

**DO NOT** create a public GitHub issue for security vulnerabilities.

Instead, please report security vulnerabilities to:
- **Email**: security@postie.dev
- **Subject**: `[SECURITY] Vulnerability Report`

### 📋 What to Include

When reporting a vulnerability, please include:

1. **Description**: Clear description of the vulnerability
2. **Impact**: Potential impact and attack scenarios
3. **Reproduction**: Step-by-step instructions to reproduce
4. **Environment**: Affected versions and configurations
5. **Proof of Concept**: Code or commands demonstrating the issue (if safe)
6. **Suggested Fix**: Any ideas for fixing the issue (optional)

### 📨 Report Template

```
Subject: [SECURITY] Vulnerability Report

Vulnerability Type: [e.g., Remote Code Execution, Information Disclosure]
Affected Component: [e.g., HTTP Client, Authentication, Collection Parser]
Affected Versions: [e.g., 1.0.0 - 1.0.5]

Description:
[Clear description of the vulnerability]

Impact:
[Describe the potential impact]

Steps to Reproduce:
1. [Step one]
2. [Step two]
3. [Step three]

Environment:
- OS: [e.g., Ubuntu 22.04]
- Go Version: [e.g., 1.21.3]
- Postie Version: [e.g., 1.0.0]

Proof of Concept:
[Safe demonstration of the vulnerability]

Additional Context:
[Any other relevant information]
```

## 🕐 Response Timeline

We aim to respond to security reports according to the following timeline:

- **Initial Response**: Within 24 hours
- **Triage**: Within 48 hours  
- **Status Update**: Weekly updates on progress
- **Resolution**: Varies based on complexity and severity

## 🔄 Disclosure Process

1. **Report Received**: We acknowledge receipt of your report
2. **Investigation**: We investigate and validate the vulnerability
3. **Fix Development**: We develop and test a fix
4. **Coordinated Disclosure**: We work with you on disclosure timing
5. **Public Release**: We release the fix and security advisory
6. **Recognition**: We acknowledge your contribution (if desired)

## 🏆 Security Researcher Recognition

We appreciate security researchers who help keep Postie secure:

- **Hall of Fame**: Recognition in our security hall of fame
- **Public Thanks**: Acknowledgment in release notes (with permission)
- **Direct Communication**: Direct line to our security team

## 🛡️ Security Best Practices

### For Users

- **Keep Updated**: Always use the latest version of Postie
- **Secure Collections**: Don't store sensitive data in collection files
- **Environment Variables**: Use environment variables for secrets
- **File Permissions**: Restrict access to collection and config files
- **Network Security**: Be cautious when making requests to untrusted endpoints

### For Developers

- **Dependency Security**: Regularly update dependencies
- **Input Validation**: Validate all inputs thoroughly
- **Error Handling**: Don't expose sensitive information in error messages
- **Authentication**: Implement secure authentication handling
- **File Operations**: Validate file paths and contents

## 🔍 Security Considerations

### Data Handling

- **Credentials**: Never log or store credentials in plain text
- **Request Data**: Be careful with sensitive request/response data
- **File Storage**: Collection files may contain sensitive information
- **Memory**: Clear sensitive data from memory when possible

### Network Security

- **TLS**: Always use HTTPS for sensitive communications
- **Certificate Validation**: Properly validate SSL/TLS certificates
- **Proxy Support**: Secure proxy authentication handling
- **Rate Limiting**: Implement proper rate limiting

### Code Security

- **Input Sanitization**: Sanitize all user inputs
- **Path Traversal**: Prevent directory traversal attacks
- **Code Injection**: Prevent code injection vulnerabilities
- **Dependency Management**: Keep dependencies updated and secure

## 📚 Security Resources

- **Go Security**: [Go Security Policy](https://golang.org/security)
- **OWASP**: [OWASP API Security Top 10](https://owasp.org/www-project-api-security/)
- **CVE Database**: [Common Vulnerabilities and Exposures](https://cve.mitre.org/)
- **Security Headers**: [Security Headers Reference](https://securityheaders.com/)

## 🔗 Contact Information

- **Security Email**: security@postie.dev
- **General Support**: support@postie.dev
- **GitHub Issues**: For non-security bugs only
- **Website**: https://postie.dev

## 📄 Legal

By reporting security vulnerabilities to us, you agree that:

1. You will not exploit the vulnerability for malicious purposes
2. You will not disclose the vulnerability publicly until we have had time to address it
3. You will not access more data than necessary to demonstrate the vulnerability
4. You will make good faith efforts to avoid privacy violations and disruption

We commit to:

1. Responding to your report in a timely manner
2. Working with you to understand and validate the issue
3. Acknowledging your contribution (if desired)
4. Not taking legal action against researchers who follow this policy

---

**Last Updated**: October 31, 2025
**Version**: 1.0