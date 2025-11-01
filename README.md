# Postie - API Testing Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jainbasil/postie?style=for-the-badge)](https://goreportcard.com/report/github.com/jainbasil/postie)
[![Release](https://img.shields.io/github/v/release/jainbasil/postie?style=for-the-badge)](https://github.com/jainbasil/postie/releases)
[![Downloads](https://img.shields.io/github/downloads/jainbasil/postie/total?style=for-the-badge)](https://github.com/jainbasil/postie/releases)

A powerful API testing tool built in Go that supports HTTP Request in Editor format (.http files). Execute HTTP requests, manage environments, and automate API testing workflows from the command line.

## ✨ Features

- **HTTP Request Files**: Write and execute requests in standard `.http` format (JetBrains HTTP Client compatible)
- **Environment Management**: Separate public and private environment files with variable substitution
- **Response Handler Scripts**: JavaScript-based response handlers for testing and assertions
- **Global Variables**: Share data between requests using global variable storage
- **Context Management**: Set default files and environments per directory for streamlined workflows
- **Response Storage**: Automatically save responses with timestamps for debugging
- **Native Performance**: Built in Go for fast, native desktop performance with single binary distribution
- **Command-Line Interface**: Full-featured CLI for automation and scripting
- **Multiple Authentication Methods**: API keys, Bearer tokens, Basic auth, and custom headers

## 📦 Installation

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

## 🚀 Quick Start

### Basic Usage

Create a `.http` file with your API requests:

```http
### Get all posts
GET https://jsonplaceholder.typicode.com/posts
Accept: application/json

### Create a new post
POST https://jsonplaceholder.typicode.com/posts
Content-Type: application/json

{
  "title": "My Post",
  "body": "This is the post content",
  "userId": 1
}
```

Run the requests:

```bash
# Run all requests in the file
postie http run requests.http

# Run a specific request by number
postie http run requests.http --request 1

# Run with verbose output
postie http run requests.http --verbose
```

### Using Environments

Create environment files for different configurations:

**http-client.env.json** (public variables):
```json
{
  "development": {
    "baseUrl": "https://api-dev.example.com",
    "apiVersion": "v1"
  },
  "production": {
    "baseUrl": "https://api.example.com",
    "apiVersion": "v1"
  }
}
```

**http-client.private.env.json** (sensitive data):
```json
{
  "development": {
    "apiKey": "dev-secret-key-12345",
    "authToken": "dev-bearer-token-xyz"
  },
  "production": {
    "apiKey": "prod-secret-key-67890",
    "authToken": "prod-bearer-token-abc"
  }
}
```

Use variables in your requests:

```http
### Get user data
GET {{baseUrl}}/{{apiVersion}}/users
Authorization: Bearer {{authToken}}
X-API-Key: {{apiKey}}
```

Run with a specific environment:

```bash
postie http run requests.http --env production
```

### Context Management

Set default values to avoid repetitive flags:

```bash
# Set context for current directory
postie context set --http-file requests.http --env development

# Now run without specifying file and environment
postie http run --request 1

# View current context
postie context show

# Clear context
postie context clear
```

### Response Handler Scripts

Add JavaScript code to test responses and extract data:

```http
### Login and save token
# @name login
POST {{baseUrl}}/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}

> {%
  client.test("Login successful", function() {
    client.assert(response.status === 200, "Expected status 200");
    client.assert(response.body.token, "Token should be present");
  });
  
  // Save token for subsequent requests
  client.global.set("authToken", response.body.token);
%}

### Use saved token
GET {{baseUrl}}/api/profile
Authorization: Bearer {{authToken}}
```

## 📖 Command Reference

### HTTP Commands

```bash
# Run requests from a file
postie http run <file.http> [options]
  --env <name>              Environment to use (default: development)
  --request <name|number>   Run specific request by name or number
  --verbose                 Show detailed output
  --save-responses          Save responses to .http-responses/ directory

# Parse and validate HTTP file
postie http parse <file.http> [options]
  --format <json|summary>   Output format (default: summary)
  --validate                Perform validation

# List requests in file or directory
postie http list [path] [options]
  --recursive               List recursively
```

### Environment Commands

```bash
# List all available environments
postie env list [options]
  --env-file <path>         Path to environment file
  --private-env-file <path> Path to private environment file

# Show variables for an environment
postie env show <environment> [options]
  --show-private            Display private variables
  --env-file <path>         Path to environment file
  --private-env-file <path> Path to private environment file
```

### Context Commands

```bash
# Set context for current directory
postie context set [options]
  --http-file <path>        Default HTTP file
  --env <name>              Default environment
  --env-file <path>         Default environment file
  --private-env-file <path> Default private environment file
  --save-responses          Enable response saving

# Show current context
postie context show

# Clear context settings
postie context clear
```

## 📚 Documentation

- [User Guide](docs/user-guide.md) - Comprehensive usage guide with examples
- [Command Reference](docs/command-reference.md) - Detailed command documentation
- [HTTP Request Format](docs/http-request-format.md) - .http file syntax and features

## 💡 Examples

### Complete API Testing Workflow

```http
### Variables can be used throughout
@hostname = api.example.com
@contentType = application/json

### Health check
GET https://{{hostname}}/health

### Login
# @name login
POST https://{{hostname}}/auth/login
Content-Type: {{contentType}}

{
  "email": "{{userEmail}}",
  "password": "{{userPassword}}"
}

> {%
  client.test("Login successful", function() {
    client.assert(response.status === 200);
    client.assert(response.body.hasOwnProperty("token"));
  });
  client.global.set("authToken", response.body.token);
  client.log("Token saved: " + response.body.token);
%}

### Get user profile (uses saved token)
GET https://{{hostname}}/api/user/profile
Authorization: Bearer {{authToken}}

> {%
  client.test("Profile retrieved", function() {
    client.assert(response.status === 200);
    client.assert(response.body.email === client.global.get("userEmail"));
  });
%}

### Create post
POST https://{{hostname}}/api/posts
Authorization: Bearer {{authToken}}
Content-Type: {{contentType}}

{
  "title": "Test Post",
  "content": "This is a test post"
}

> {%
  client.test("Post created", function() {
    client.assert(response.status === 201);
  });
  client.global.set("postId", response.body.id);
%}

### Update post
PUT https://{{hostname}}/api/posts/{{postId}}
Authorization: Bearer {{authToken}}
Content-Type: {{contentType}}

{
  "title": "Updated Post",
  "content": "This post has been updated"
}

### Delete post
DELETE https://{{hostname}}/api/posts/{{postId}}
Authorization: Bearer {{authToken}}
```

## 🛠️ Development

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

### Code Standards

- Follow [Go Code Review Comments](https://github.com/golang/go/wiki/CodeReviewComments)
- Use `gofmt` for formatting
- Write tests for new features
- Update documentation for API changes
- Follow [Conventional Commits](https://www.conventionalcommits.org/) for commit messages

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

Made with ❤️ in Go
