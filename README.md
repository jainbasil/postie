# Postie - API Testing Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jainbasil/postie?style=for-the-badge)](https://goreportcard.com/report/github.com/jainbasil/postie)
[![Release](https://img.shields.io/github/v/release/jainbasil/postie?style=for-the-badge)](https://github.com/jainbasil/postie/releases)
[![Downloads](https://img.shields.io/github/downloads/jainbasil/postie/total?style=for-the-badge)](https://github.com/jainbasil/postie/releases)

A powerful API test client built in Go that provides command-line interface for testing and debugging APIs, similar to Postman but designed for developers who love the terminal.

## Features

- **Native Performance**: Built in Go for fast, native desktop performance
- **Command- e Interface**: Full- tured CLI for automation and scripting
- **Multiple Authentication Methods**: API keys, Bearer tokens, Basic auth, and custom headers
- **Full HTTP Support**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- **Response Analysis**: JSON formatting, status checking, response timing
- **Middleware Support**: Extensible middleware for logging, rate limiting, retries
- **Request Building**: Fluent API for building complex requests
- **Collection Support**: Import/export API collections in standard JSON format
- **Fast & Lightweight**: Single binary with no external dependencies

## Installation

### Binary Releases (Recommended)

Download the latest binary for your platform from the [releases page](https://github.com/jainbasil/postie/releases):

```bash
# Linux (x64)
curl -L https://github.com/jainbasil/postie/releases/latest/download/postie-linux-amd64 -o postie
chmod +x postie

# macOS (x64)
curl -L https://github.com/jainbasil/postie/releases/latest/download/postie-darwin-amd64 -o postie
chmod +x postie

# Windows (x64)
# Download postie-windows-amd64.exe from releases page
```

### Using Go Install

```bash
go install github.com/jainbasil/postie@latest
```

### Build from Source

```bash
git clone https://github.com/jainbasil/postie.git
cd postie
go build -o postie .
```

## Quick Start

See wiki - [[Postie-Command-Reference]]

### Collection File Format

Collections are stored as JSON files. Here's the basic structure:

```bash
# View collection structure
./postie list collections/jsonplaceholder.collection.json

# The collection includes:
# - Collection metadata (name, description, version)
# - Environment variables
# - API groups and endpoints
# - Authentication settings
# - Pre/post request scripts
```

For detailed collection format documentation, see `docs/collection-format.md`.

## Development

### Building

```bash
# Build for current platform
go build -o postie .

# Cross-compile for different platforms
GOOS=windows GOARCH=amd64 go build -o postie.exe .
GOOS=linux GOARCH=amd64 go build -o postie-linux .
GOOS=darwin GOARCH=amd64 go build -o postie-mac .
```

### Testing

```bash
# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run demo
./postie demo
```

### VS Code Tasks

The project includes VS Code tasks for development:

- **Build Postie**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Build Postie"

## Contributing

We love contributions! Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

### Quick Contributing Steps

1. **Fork the repository**
   ```bash
   git clone https://github.com/jainbasil/postie.git
   cd postie
   ```

2. **Create a feature branch**
   ```bash
   git checkout -b feature/amazing-feature
   ```

3. **Make your changes and test**
   ```bash
   go test ./...
   go build -o postie .
   ./postie demo  # Test the changes
   ```

4. **Commit with conventional commits**
   ```bash
   git commit -m "feat: add amazing new feature"
   # or
   git commit -m "fix: resolve issue with authentication"
   # or  
   git commit -m "docs: update API documentation"
   ```

5. **Push and create a Pull Request**
   ```bash
   git push origin feature/amazing-feature
   ```

### Development Setup

```bash
# Clone the repository
git clone https://github.com/jainbasil/postie.git
cd postie

# Install dependencies
go mod download

# Build the project
go build -o postie .

# Run tests
go test ./...

# Run with coverage
go test -cover ./...

# Run demo
./postie demo
```

### Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Write tests for new features
- Update documentation for API changes
- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages

### VS Code Development

The project includes VS Code configuration:

- **Build Postie**: `Ctrl+Shift+P` → "Tasks: Run Task" → "Build Postie"
- Includes Go extension recommendations
- Pre-configured debug settings

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ in Go
