# Postie - Native Desktop API Testing Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/yourusername/postie?style=for-the-badge)](https://goreportcard.com/report/github.com/yourusername/postie)
[![Release](https://img.shields.io/github/v/release/yourusername/postie?style=for-the-badge)](https://github.com/yourusername/postie/releases)
[![Downloads](https://img.shields.io/github/downloads/yourusername/postie/total?style=for-the-badge)](https://github.com/yourusername/postie/releases)

A powerful, native desktop API client built in Go that provides both command-line and programmatic interfaces for testing and debugging APIs, similar to Postman but designed for developers who love the terminal.

## 🌟 Features

- 🚀 **Native Performance**: Built in Go for fast, native desktop performance
- 🔧 **Command-Line Interface**: Full-featured CLI for automation and scripting
- 🔐 **Multiple Authentication Methods**: API keys, Bearer tokens, Basic auth, and custom headers
- 🌐 **Full HTTP Support**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS
- 📊 **Response Analysis**: JSON formatting, status checking, response timing
- 🔄 **Middleware Support**: Extensible middleware for logging, rate limiting, retries
- 📝 **Request Building**: Fluent API for building complex requests
- 📁 **Collection Support**: Import/export API collections in standard JSON format
- ⚡ **Fast & Lightweight**: Single binary with no external dependencies

## 📦 Installation

### Binary Releases (Recommended)

Download the latest binary for your platform from the [releases page](https://github.com/yourusername/postie/releases):

```bash
# Linux (x64)
curl -L https://github.com/yourusername/postie/releases/latest/download/postie-linux-amd64 -o postie
chmod +x postie

# macOS (x64)
curl -L https://github.com/yourusername/postie/releases/latest/download/postie-darwin-amd64 -o postie
chmod +x postie

# Windows (x64)
# Download postie-windows-amd64.exe from releases page
```

### Using Go Install

```bash
go install github.com/yourusername/postie@latest
```

### Build from Source

```bash
git clone https://github.com/yourusername/postie.git
cd postie
go build -o postie .
```

### Package Managers

```bash
# Homebrew (macOS/Linux)
brew install yourusername/tap/postie

# Scoop (Windows)
scoop bucket add yourusername https://github.com/yourusername/scoop-bucket
scoop install postie

# Arch Linux (AUR)
yay -S postie-bin
```

## 🚀 Quick Start

### Basic Usage

```bash
# Basic GET request
./postie get https://api.github.com/users/octocat

# POST request with data
./postie post https://httpbin.org/post

# Run interactive demo
./postie demo

# Show help and all commands
./postie help
```

### Advanced Examples

```bash
# POST with JSON data and custom headers
./postie post https://api.example.com/users \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer token123" \
  --data '{"name": "John", "email": "john@example.com"}'

# Run API collection
./postie run collections/my-api.collection.json

# Run with specific environment
./postie run collections/my-api.collection.json --env "Production"

# List all requests in a collection
./postie list collections/my-api.collection.json
```

## 📋 Commands

### HTTP Methods
```bash
postie get <url>          # GET request
postie post <url>         # POST request  
postie put <url>          # PUT request
postie delete <url>       # DELETE request
postie patch <url>        # PATCH request
postie head <url>         # HEAD request
postie options <url>      # OPTIONS request
```

### Collection Management
```bash
postie create collection <name>              # Create new collection
postie create apigroup <file> <name>         # Create API group
postie create api <file> <group> <name>      # Create API endpoint
postie update collection <file>              # Update collection
postie remove apigroup <file> <id>           # Remove API group
postie list <collection>                     # List all requests
postie env <collection>                      # Show environments
postie run <collection>                      # Run collection
```

### Programmatic Usage

```go
package main

import (
    "fmt"
    "time"
    "postie/pkg/client"
    "postie/pkg/auth"
)

func main() {
    // Create API client
    apiClient := client.NewClient(&client.Config{
        BaseURL: "https://api.example.com",
        Timeout: 30 * time.Second,
        Headers: map[string]string{
            "User-Agent": "MyApp/1.0",
        },
    })

    // Make authenticated GET request
    resp, err := apiClient.GET("/users").
        Header("Authorization", "Bearer your-token").
        Param("page", "1").
        Execute()
    
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer resp.Response.Body.Close()

    // Parse JSON response
    var users []User
    if err := resp.JSON(&users); err != nil {
        fmt.Printf("JSON error: %v\n", err)
        return
    }

    fmt.Printf("Retrieved %d users\n", len(users))
}
```

## Architecture

### Project Structure

```
postie/
├── main.go                 # CLI application entry point
├── pkg/
│   ├── client/            # Core HTTP client
│   │   ├── client.go      # Main client implementation
│   │   ├── request.go     # Request builder
│   │   └── response.go    # Response wrapper
│   ├── auth/              # Authentication handlers
│   │   └── authenticator.go
│   └── middleware/        # Request/response middleware
│       └── common.go
├── examples/              # Usage examples
│   └── demo.go           # Comprehensive demo
├── collections/          # API collection files
│   ├── jsonplaceholder.collection.json
│   ├── development.environment.json
│   └── production.environment.json
├── docs/                 # Documentation
│   └── collection-format.md
└── tests/                # Test files
```

### Core Components

#### 1. HTTP Client (`pkg/client`)

The core client provides a fluent interface for building and executing HTTP requests:

```go
// Create client with configuration
client := client.NewClient(&client.Config{
    BaseURL: "https://api.example.com",
    Timeout: 30 * time.Second,
    Middleware: []client.Middleware{
        middleware.LoggingMiddleware,
    },
})

// Build and execute request
response, err := client.POST("/api/users").
    Header("Content-Type", "application/json").
    JSON(userData).
    Execute()
```

#### 2. Authentication (`pkg/auth`)

Multiple authentication methods are supported:

```go
// API Key in header
apiKey := auth.NewAPIKeyAuth("X-API-Key", "your-key", "header")

// Bearer token
bearer := auth.NewBearerTokenAuth("your-jwt-token")

// Basic authentication
basic := auth.NewBasicAuth("username", "password")

// Custom header
custom := auth.NewCustomHeaderAuth("X-Custom-Auth", "value")
```

#### 3. Middleware (`pkg/middleware`)

Extensible middleware system for cross-cutting concerns:

```go
// Built-in middleware
middleware.LoggingMiddleware          // Request/response logging
middleware.UserAgentMiddleware("...")  // Custom User-Agent
middleware.ErrorHandlingMiddleware    // Error handling
middleware.NewRateLimitMiddleware(5.0) // Rate limiting
```

## API Reference

### Client Methods

- `GET(url)` - Create GET request
- `POST(url)` - Create POST request  
- `PUT(url)` - Create PUT request
- `DELETE(url)` - Create DELETE request
- `PATCH(url)` - Create PATCH request
- `HEAD(url)` - Create HEAD request
- `OPTIONS(url)` - Create OPTIONS request

### Request Builder

- `Header(key, value)` - Set request header
- `Headers(map[string]string)` - Set multiple headers
- `Param(key, value)` - Set URL parameter
- `Params(map[string]string)` - Set multiple parameters
- `JSON(data)` - Set JSON body
- `Text(string)` - Set text body
- `Form(map[string]string)` - Set form data
- `Body(io.Reader)` - Set raw body
- `Context(context.Context)` - Set request context
- `Execute()` - Send the request

### Response Methods

- `GetBody()` - Get response body as bytes
- `Text()` - Get response body as string
- `JSON(interface{})` - Unmarshal JSON response
- `IsSuccess()` - Check if status 2xx
- `IsError()` - Check if status 4xx/5xx
- `IsClientError()` - Check if status 4xx
- `IsServerError()` - Check if status 5xx
- `Size()` - Get response size in bytes
- `ContentType()` - Get Content-Type header

## Command Line Interface

### Commands

```bash
postie get <url>          # GET request
postie post <url>         # POST request  
postie put <url>          # PUT request
postie delete <url>       # DELETE request

# Collection operations
postie demo               # Run demonstrations
postie help               # Show help
```

### Examples

```bash
# Simple GET request
# GET with custom headers
./postie get https://httpbin.org/get

# POST with JSON data  
./postie post https://httpbin.org/post

# Test different response status codes
./postie get https://httpbin.org/status/404
./postie get https://httpbin.org/status/200

# Test request timing
./postie get https://httpbin.org/delay/2
```

## Advanced Usage

### Custom Middleware

```go
// Create custom middleware
func CustomMiddleware(req *http.Request, resp *http.Response) error {
    // Add custom logic here
    fmt.Printf("Custom processing for %s\n", req.URL)
    return nil
}

// Use in client
client := client.NewClient(&client.Config{
    Middleware: []client.Middleware{
        CustomMiddleware,
        middleware.LoggingMiddleware,
    },
})
```

### Error Handling

```go
resp, err := client.GET("/api/data").Execute()
if err != nil {
    fmt.Printf("Request failed: %v\n", err)
    return
}

if resp.IsError() {
    fmt.Printf("HTTP Error: %s\n", resp.Status)
    if resp.IsClientError() {
        // Handle 4xx errors
    }
    if resp.IsServerError() {
        // Handle 5xx errors  
    }
}
```

### Timeouts and Context

```go
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()

resp, err := client.GET("/slow-endpoint").
    Context(ctx).
    Execute()
```

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

## 🤝 Contributing

We love contributions! Please read our [Contributing Guide](CONTRIBUTING.md) to get started.

### Quick Contributing Steps

1. **Fork the repository**
   ```bash
   git clone https://github.com/yourusername/postie.git
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
git clone https://github.com/yourusername/postie.git
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

### MIT License Summary

```
MIT License

Copyright (c) 2025 Postie Contributors

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all
copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
```

## 🗺️ Roadmap

### Current Version (v1.0.0)
- ✅ Full HTTP method support (GET, POST, PUT, DELETE, etc.)
- ✅ Multiple authentication methods  
- ✅ JSON collection format
- ✅ Environment variable support
- ✅ Command-line interface
- ✅ Middleware system

### Upcoming Features

#### v1.1.0 - Enhanced Collections
- [ ] Collection import/export from Postman
- [ ] Environment switching improvements
- [ ] Request templates
- [ ] Response caching

#### v1.2.0 - Advanced Features  
- [ ] WebSocket support
- [ ] GraphQL query support
- [ ] Request history
- [ ] Performance benchmarking
- [ ] Response assertions

#### v2.0.0 - Desktop GUI
- [ ] Cross-platform desktop GUI using Fyne
- [ ] Visual request builder
- [ ] Response visualization
- [ ] Interactive collection management

#### Future Considerations
- [ ] Plugin system for extensibility
- [ ] Team collaboration features
- [ ] API documentation generation
- [ ] Mock server capabilities
- [ ] CI/CD integrations

## 📊 Project Statistics

![Lines of Code](https://img.shields.io/tokei/lines/github/yourusername/postie?style=flat-square)
![Code Size](https://img.shields.io/github/languages/code-size/yourusername/postie?style=flat-square)
![Repository Size](https://img.shields.io/github/repo-size/yourusername/postie?style=flat-square)

## 🆘 Support & Community

### Getting Help

- � **Documentation**: [docs.postie.dev](https://docs.postie.dev)
- 🐛 **Bug Reports**: [GitHub Issues](https://github.com/yourusername/postie/issues)
- 💬 **Discussions**: [GitHub Discussions](https://github.com/yourusername/postie/discussions)
- 💭 **Feature Requests**: [GitHub Issues](https://github.com/yourusername/postie/issues/new?template=feature_request.md)

### Community

- 🌟 **Star us on GitHub** if you find Postie useful!
- 🐦 **Follow updates**: [@PostieAPI](https://twitter.com/PostieAPI)
- 📧 **Email**: support@postie.dev
- 💼 **LinkedIn**: [Postie API Tool](https://linkedin.com/company/postie-api)

### Security

If you discover a security vulnerability, please send an email to security@postie.dev. All security vulnerabilities will be promptly addressed.

## 🙏 Acknowledgments

- Inspired by [Postman](https://postman.com) for API testing workflows
- Built with the amazing [Go](https://golang.org) programming language
- Thanks to all [contributors](https://github.com/yourusername/postie/contributors) who make this project possible!

## 📈 Star History

[![Star History Chart](https://api.star-history.com/svg?repos=yourusername/postie&type=Date)](https://star-history.com/#yourusername/postie&Date)

---

Made with ❤️ in Go