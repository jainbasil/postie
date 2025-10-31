# Postie - Native Desktop API Testing Tool

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg?style=for-the-badge)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/jainbasil/postie?style=for-the-badge)](https://goreportcard.com/report/github.com/jainbasil/postie)
[![Release](https://img.shields.io/github/v/release/jainbasil/postie?style=for-the-badge)](https://github.com/jainbasil/postie/releases)
[![Downloads](https://img.shields.io/github/downloads/jainbasil/postie/total?style=for-the-badge)](https://github.com/jainbasil/postie/releases)

A powerful, native desktop API client built in Go that provides both command-line and programmatic interfaces for testing and debugging APIs, similar to Postman but designed for developers who love the terminal.

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

This guide will walk you through the complete workflow of using Postie, from creating your first collection to running API tests.

### 1. Basic API Testing

Start with simple HTTP requests to test any API:

```bash
# Basic GET request
./postie get https://api.github.com/users/octocat

# POST request with JSON data
./postie post https://httpbin.org/post \
  --header "Content-Type: application/json" \
  --data '{"name": "John", "email": "john@example.com"}'

# Test different HTTP methods
./postie put https://httpbin.org/put --data '{"id": 1, "name": "Updated"}'
./postie delete https://httpbin.org/delete

# Run interactive demo to see all features
./postie demo
```

### 2. Create Your First Collection

Collections organize your API requests into reusable groups:

```bash
# Create a new collection
./postie create collection "My Blog API" --file blog-api.collection.json

# This creates a new file: blog-api.collection.json
```

### 3. Create API Groups

Organize related APIs into logical groups:

```bash
# Create a group for user-related APIs
./postie create apigroup blog-api.collection.json "Users"

# Create a group for blog post APIs  
./postie create apigroup blog-api.collection.json "Posts"

# View your collection structure
./postie list blog-api.collection.json
```

### 4. Add API Endpoints

Add specific API endpoints to your groups:

```bash
# Add a GET request to fetch all users
./postie create api blog-api.collection.json "users-group-id" "Get All Users" GET "{{baseUrl}}/users"

# Add a POST request to create a user
./postie create api blog-api.collection.json "users-group-id" "Create User" POST "{{baseUrl}}/users"

# Add a GET request to fetch posts
./postie create api blog-api.collection.json "posts-group-id" "Get All Posts" GET "{{baseUrl}}/posts"
```

### 5. Set Up Environments

Create different environments for development, staging, and production:

```bash
# View existing environments (if any)
./postie env blog-api.collection.json

# Environments are defined in the collection or separate files
# See collections/development.environment.json for examples
```

### 6. Run Your APIs

Execute your collection with different environments:

```bash
# Run entire collection
./postie run blog-api.collection.json

# Run with specific environment
./postie run blog-api.collection.json --env "Development"

# Run specific request by name
./postie run blog-api.collection.json --request "Get All Users"

# Run specific request by ID
./postie run blog-api.collection.json --id "user-123"
```

### 7. Set Default Collection Context

Save time by setting a default collection and environment - no need to type the collection path every time:

```bash
# Set your collection as the default context
./postie context set --collection blog-api.collection.json --env "Development"

# Now run commands without specifying the collection
./postie list              # Uses blog-api.collection.json
./postie env               # Shows environments from blog-api.collection.json
./postie run               # Runs the collection with Development environment

# Override the environment when needed
./postie run --env "Production"
./postie list --env "Staging"

# View current context
./postie context show

# Clear context when done
./postie context clear
```

**Context Benefits:**
- Set once, use everywhere - no need to type collection paths repeatedly
- Perfect for working on a single project
- Similar to `kubectl config` or `az account set` in other CLI tools
- Context is stored in `~/.postie/context.json`

### 8. Update and Manage APIs

Modify your APIs as they evolve:

```bash
# Update API endpoint
./postie update api blog-api.collection.json "user-get-id" --url "{{baseUrl}}/users/v2"

# Update API method
./postie update api blog-api.collection.json "user-post-id" --method PATCH

# Update group information
./postie update apigroup blog-api.collection.json "users-group-id" --name "User Management"

# Remove APIs or groups when no longer needed
./postie remove api blog-api.collection.json "old-api-id"
./postie remove apigroup blog-api.collection.json "deprecated-group-id"
```

### 9. Explore Examples

Try the included examples to learn advanced features:

```bash
# Run the comprehensive demo
./postie demo

# Explore the JSONPlaceholder collection
./postie list collections/jsonplaceholder.collection.json
./postie run collections/jsonplaceholder.collection.json

# Run with different environments
./postie run collections/jsonplaceholder.collection.json --env "JSONPlaceholder"
```

### Pro Tips

- Use `{{baseUrl}}` and other variables in your URLs for environment flexibility
- **Set context once** with `./postie context set` to avoid typing collection paths
- Context works like `kubectl config` or `az account set` - set it and forget it
- Start with the demo to understand all features: `./postie demo`
- Use `--help` with any command for detailed usage: `./postie create --help`
- Collections are JSON files - you can edit them directly if needed
- Override context environment anytime with `--env` flag

## Commands

### HTTP Request Commands

Execute HTTP requests directly from the command line:

```bash
# GET requests
./postie get https://api.github.com/users/octocat
./postie get https://httpbin.org/get --header "Authorization: Bearer token123"
./postie get "https://api.example.com/users?page=1&limit=10"

# POST requests
./postie post https://httpbin.org/post
./postie post https://api.example.com/users \
  --header "Content-Type: application/json" \
  --data '{"name": "John Doe", "email": "john@example.com"}'

# PUT requests
./postie put https://api.example.com/users/123 \
  --header "Content-Type: application/json" \
  --data '{"name": "Jane Doe", "email": "jane@example.com"}'

# DELETE requests
./postie delete https://api.example.com/users/123
./postie delete https://api.example.com/users/123 \
  --header "Authorization: Bearer token123"
```

### Collection Management Commands

#### Create Commands

```bash
# Create a new collection
./postie create collection "My API Collection"
./postie create collection "Blog API" --file blog-api.collection.json

# Create API groups within a collection
./postie create apigroup blog-api.collection.json "Users"
./postie create apigroup blog-api.collection.json "Posts" --id "posts-group"

# Create API endpoints
./postie create api blog-api.collection.json "users-group-id" "Get All Users" GET "{{baseUrl}}/users"
./postie create api blog-api.collection.json "users-group-id" "Create User" POST "{{baseUrl}}/users" --id "create-user"
./postie create api blog-api.collection.json "posts-group-id" "Get Post" GET "{{baseUrl}}/posts/{{postId}}"
```

#### Update Commands

```bash
# Update collection metadata
./postie update collection blog-api.collection.json --name "Updated Blog API"
./postie update collection blog-api.collection.json --description "API for blog management"

# Update API groups
./postie update apigroup blog-api.collection.json "users-group-id" --name "User Management"
./postie update apigroup blog-api.collection.json "posts-group-id" --description "Blog post operations"

# Update API endpoints
./postie update api blog-api.collection.json "get-users-id" --name "Fetch All Users"
./postie update api blog-api.collection.json "create-user-id" --method PATCH
./postie update api blog-api.collection.json "get-post-id" --url "{{baseUrl}}/api/v2/posts/{{postId}}"
```

#### Remove Commands

```bash
# Remove API endpoints
./postie remove api blog-api.collection.json "deprecated-api-id"

# Remove API groups (removes all APIs in the group)
./postie remove apigroup blog-api.collection.json "deprecated-group-id"
```

### Collection Execution Commands

```bash
# Run entire collection
./postie run collections/jsonplaceholder.collection.json
./postie run blog-api.collection.json

# Run with specific environment
./postie run blog-api.collection.json --env "Development"
./postie run blog-api.collection.json --env "Production"

# Run specific request by name
./postie run blog-api.collection.json --request "Get All Users"
./postie run blog-api.collection.json --request "Create User"

# Run specific request by ID
./postie run blog-api.collection.json --id "user-get-all"
./postie run blog-api.collection.json --id "post-create"
```

### Information Commands

```bash
# List all requests in a collection
./postie list collections/jsonplaceholder.collection.json
./postie list blog-api.collection.json --env "Development"

# List using context (no collection file needed)
./postie list
./postie list --env "Production"

# Show environments in a collection
./postie env collections/jsonplaceholder.collection.json
./postie env blog-api.collection.json

# Show environments using context
./postie env

# Show help for any command
./postie help
./postie --help
./postie -h
```

### Context Management Commands

Set default collection and environment to avoid typing collection paths repeatedly:

```bash
# Set default collection and environment
./postie context set --collection blog-api.collection.json --env "Development"
./postie context set --collection collections/jsonplaceholder.collection.json --env "Production"

# Set only collection (use default environment)
./postie context set --collection blog-api.collection.json

# Set only environment (keep existing collection)
./postie context set --env "Staging"

# View current context
./postie context show

# Clear all context settings
./postie context clear

# After setting context, run commands without specifying collection:
./postie run                    # Uses context collection and environment
./postie list                   # Uses context collection
./postie env                    # Uses context collection
./postie run --env "Testing"    # Override context environment
```

**Context is stored in:** `~/.postie/context.json`

### Demo and Learning Commands

```bash
# Run interactive demonstration
./postie demo

# This demo includes:
# - Basic HTTP requests (GET, POST, PUT, DELETE)
# - Authentication examples
# - Error handling demonstrations
# - Response parsing examples
# - Collection workflow examples
```

### Advanced Usage Examples

```bash
# Chain multiple operations
./postie create collection "E-commerce API" --file ecommerce.collection.json
./postie create apigroup ecommerce.collection.json "Products"
./postie create api ecommerce.collection.json "products-group-id" "List Products" GET "{{baseUrl}}/products"
./postie run ecommerce.collection.json --env "Development"

# Work with authentication
./postie post https://api.example.com/auth/login \
  --header "Content-Type: application/json" \
  --data '{"username": "admin", "password": "secret"}'

# Test API responses and status codes
./postie get https://httpbin.org/status/200  # Success
./postie get https://httpbin.org/status/404  # Not Found
./postie get https://httpbin.org/status/500  # Server Error

# Test with delays (for timeout testing)
./postie get https://httpbin.org/delay/2     # 2 second delay
./postie get https://httpbin.org/delay/5     # 5 second delay

# Complex POST with multiple headers
./postie post https://api.example.com/users \
  --header "Content-Type: application/json" \
  --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..." \
  --header "X-API-Version: v2" \
  --data '{"name": "Alice", "email": "alice@example.com", "role": "admin"}'
```

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

## Command Line Interface

Postie provides a comprehensive CLI for API testing and collection management. The tool supports:

- **HTTP Methods**: GET, POST, PUT, DELETE with full header and body support
- **Collection Management**: Create, update, and organize API collections
- **Environment Support**: Switch between development, staging, and production environments
- **Authentication**: Multiple auth methods including Bearer tokens, API keys, and Basic auth
- **Response Analysis**: Automatic JSON formatting, status checking, and timing metrics

### Common Flags and Options

```bash
# Global options available for HTTP requests:
--header, -H    Add custom headers (can be used multiple times)
--data, -d      Request body data (JSON, form data, etc.)
--env, -e       Specify environment for collection runs
--request, -r   Run specific request by name
--id           Run specific request by ID
--file, -f      Specify output file for collection creation

# Examples with flags:
./postie get https://api.example.com/users \
  --header "Authorization: Bearer token123" \
  --header "Content-Type: application/json"

./postie post https://api.example.com/users \
  --data '{"name": "John", "email": "john@example.com"}' \
  --header "Content-Type: application/json"
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
