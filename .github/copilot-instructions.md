# Postie - API Client Project - Copilot Instructions

This is a Go-based native desktop API client project similar to Postman that demonstrates:
- HTTP request building and execution (GET, POST, PUT, DELETE, etc.)
- Authentication handling (API keys, Bearer tokens, Basic auth, custom headers)
- Response processing and formatting with JSON support
- Request/response middleware system
- Command-line interface for automation and scripting
- Comprehensive testing and examples

## Project Structure
- `main.go` - CLI application entry point
- `pkg/client/` - Core HTTP client implementation with fluent API
- `pkg/auth/` - Authentication handlers for multiple auth methods
- `pkg/middleware/` - Request/response middleware for logging, rate limiting, etc.
- `examples/` - Usage examples and comprehensive demos
- `tests/` - Unit and integration tests
- `.vscode/tasks.json` - VS Code build tasks

## Key Features
- Native desktop performance with single binary distribution
- Fluent API for building complex HTTP requests
- Multiple authentication methods with extensible interface
- Middleware system for cross-cutting concerns
- JSON response parsing and formatting
- Command-line interface for automation
- Comprehensive error handling and status checking
- Request timing and performance metrics

## Development Guidelines
- Follow Go conventions and best practices
- Use interfaces for extensibility and testability
- Include comprehensive error handling with proper error types
- Provide clear documentation and working examples
- Write tests for all public APIs
- Support both CLI and programmatic usage
- Design for single binary distribution

## Usage Examples
- CLI: `./postie get https://api.github.com/users/octocat`
- Programmatic: Fluent API with method chaining
- Authentication: Multiple auth methods with simple configuration
- Middleware: Extensible system for request/response processing