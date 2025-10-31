# Contributing to Postie

Thank you for your interest in contributing to Postie! We welcome contributions from developers of all experience levels.

## 🚀 Quick Start for Contributors

1. Fork the repository on GitHub
2. Clone your fork locally
3. Create a feature branch
4. Make your changes
5. Test your changes
6. Submit a pull request

## 📋 Table of Contents

- [Code of Conduct](#code-of-conduct)
- [Development Setup](#development-setup)
- [How to Contribute](#how-to-contribute)
- [Coding Standards](#coding-standards)
- [Testing](#testing)
- [Pull Request Process](#pull-request-process)
- [Issue Guidelines](#issue-guidelines)
- [Release Process](#release-process)

## 🤝 Code of Conduct

This project adheres to a code of conduct. By participating, you are expected to uphold this code. Please report unacceptable behavior to conduct@postie.dev.

### Our Standards

- **Be respectful**: Treat everyone with respect and kindness
- **Be inclusive**: Welcome and support people of all backgrounds
- **Be collaborative**: Work together constructively
- **Be patient**: Help others learn and grow
- **Be constructive**: Provide helpful feedback

## 🛠️ Development Setup

### Prerequisites

- Go 1.21 or later
- Git
- A code editor (VS Code recommended)

### Setup Steps

```bash
# 1. Fork the repository on GitHub
# 2. Clone your fork
git clone https://github.com/yourusername/postie.git
cd postie

# 3. Add upstream remote
git remote add upstream https://github.com/originalowner/postie.git

# 4. Install dependencies
go mod download

# 5. Build the project
go build -o postie .

# 6. Run tests to ensure everything works
go test ./...

# 7. Run the demo to verify functionality
./postie demo
```

### VS Code Setup

1. Install the Go extension
2. Open the project folder
3. The project includes pre-configured tasks and debug settings

## 🔄 How to Contribute

### Types of Contributions

We welcome various types of contributions:

- 🐛 **Bug fixes**
- ✨ **New features**  
- 📝 **Documentation improvements**
- 🧪 **Tests**
- 🎨 **UI/UX improvements**
- 🔧 **Performance optimizations**
- 🌐 **Translations**

### Contribution Workflow

1. **Check existing issues** or create a new one to discuss your contribution
2. **Fork the repository** and create a branch from `main`
3. **Make your changes** following our coding standards
4. **Add tests** for new functionality
5. **Update documentation** if needed
6. **Test your changes** thoroughly
7. **Submit a pull request** with a clear description

### Branch Naming Convention

Use descriptive branch names:

```bash
feature/add-websocket-support
fix/authentication-bug
docs/update-api-examples
refactor/improve-error-handling
```

## 🎨 Coding Standards

### Go Code Style

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` to format your code
- Use `go vet` to check for common errors
- Follow Go naming conventions

### Code Organization

```
pkg/
├── client/         # HTTP client implementation
├── auth/          # Authentication handlers  
├── middleware/    # Request/response middleware
├── collection/    # Collection management
└── utils/         # Utility functions

examples/          # Usage examples
tests/            # Test files
docs/             # Documentation
```

### Comments and Documentation

- Document all public functions and types
- Use clear, concise comments
- Include examples in documentation
- Update README.md for new features

### Error Handling

- Use proper error wrapping
- Provide meaningful error messages
- Handle edge cases gracefully
- Include context in errors

```go
// Good
if err != nil {
    return fmt.Errorf("failed to parse URL %s: %w", url, err)
}

// Bad
if err != nil {
    return err
}
```

## 🧪 Testing

### Writing Tests

- Write tests for all new functionality
- Use table-driven tests where appropriate
- Include both positive and negative test cases
- Test error conditions

### Test Categories

```bash
# Unit tests
go test ./pkg/...

# Integration tests  
go test ./tests/integration/...

# End-to-end tests
go test ./tests/e2e/...

# Benchmark tests
go test -bench=. ./...

# Coverage report
go test -cover ./...
```

### Test Examples

```go
func TestClientGET(t *testing.T) {
    tests := []struct {
        name     string
        url      string
        expected int
        wantErr  bool
    }{
        {
            name:     "valid GET request",
            url:      "https://httpbin.org/get",
            expected: 200,
            wantErr:  false,
        },
        {
            name:     "invalid URL",
            url:      "invalid-url",
            expected: 0,
            wantErr:  true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client := NewClient(&Config{})
            resp, err := client.GET(tt.url).Execute()
            
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            
            assert.NoError(t, err)
            assert.Equal(t, tt.expected, resp.StatusCode)
        })
    }
}
```

## 📥 Pull Request Process

### Before Submitting

- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Tests added/updated
- [ ] Documentation updated
- [ ] No merge conflicts
- [ ] All tests pass

### PR Template

When creating a pull request, include:

```markdown
## Description
Brief description of the changes

## Type of Change
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Tests added/updated
- [ ] Manual testing completed
- [ ] All tests pass

## Checklist
- [ ] Code follows style guidelines
- [ ] Self-review completed
- [ ] Documentation updated
```

### Review Process

1. **Automated checks** must pass (tests, linting)
2. **Maintainer review** (usually within 2-3 days)
3. **Address feedback** if any
4. **Final approval** and merge

## 🐛 Issue Guidelines

### Bug Reports

Use the bug report template and include:

- **Steps to reproduce**
- **Expected behavior**
- **Actual behavior**
- **Environment details** (OS, Go version, etc.)
- **Code examples** or logs

### Feature Requests

Use the feature request template and include:

- **Problem description**
- **Proposed solution**
- **Alternative solutions considered**
- **Additional context**

### Issue Labels

- `bug` - Something isn't working
- `enhancement` - New feature or request
- `documentation` - Improvements to docs
- `good first issue` - Good for newcomers
- `help wanted` - Extra attention needed
- `priority:high` - High priority issue

## 🚀 Release Process

### Versioning

We use [Semantic Versioning](https://semver.org/):

- **MAJOR**: Breaking changes
- **MINOR**: New features (backward compatible)
- **PATCH**: Bug fixes (backward compatible)

### Release Checklist

1. Update version in relevant files
2. Update CHANGELOG.md
3. Create release notes
4. Tag the release
5. Build and upload binaries
6. Update package managers

## 🙋 Getting Help

- **Questions**: Use [GitHub Discussions](https://github.com/yourusername/postie/discussions)
- **Bugs**: Create an [issue](https://github.com/yourusername/postie/issues)
- **Security**: Email security@postie.dev
- **General**: Email support@postie.dev

## 🎉 Recognition

Contributors are recognized in:

- **README.md** contributors section
- **Release notes** for their contributions
- **Hall of Fame** for significant contributions

Thank you for contributing to Postie! 🚀