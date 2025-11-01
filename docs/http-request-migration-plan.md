# HTTP Request in Editor Migration Plan

This document outlines the complete migration from Postie's current JSON-based collection format to the "HTTP Request in Editor" specification as defined by JetBrains.

## Overview

We are migrating from JSON collections to the HTTP Request in Editor format to provide:
- Better version control with text-based diffs
- Industry-standard format with IDE support
- Simpler, more intuitive syntax for HTTP requests
- Built-in scripting and response handling capabilities

## Key Differences Between Current and Target Formats

### Current Format (Postie JSON Collections)
- **Format**: JSON-based collection files
- **Structure**: Hierarchical with collections, folders (apiGroup), requests
- **Variables**: JSON objects with key-value pairs, environment support
- **Authentication**: JSON configuration with multiple auth types
- **Scripts**: JavaScript in JSON arrays
- **File Extension**: `.collection.json`
- **Multiple Requests**: Stored in single JSON file with hierarchical structure

### Target Format (HTTP Request in Editor)
- **Format**: Plain text files with HTTP-like syntax
- **Structure**: Sequential requests separated by `###`
- **Variables**: Template variables `{{variableName}}`
- **Authentication**: Headers-based (implicit)
- **Scripts**: Response handlers with `> {%  %}`
- **File Extension**: `.http` or `.rest`
- **Multiple Requests**: Text-based with separator syntax

## Migration Strategy

### Phase 1: Fresh Implementation - HTTP Request Parser

**Create New HTTP Request Parser (No Legacy Support)**
- Build lexer for tokenizing `.http` files
- Implement grammar-based parser following the JetBrains spec exactly
- Support all request elements per specification:
  - Request line (method, URL, HTTP version) with optional method and version
  - Headers with multi-line support
  - Message body (inline and file references with `< ./file.json`)
  - Environment variables `{{var}}`
  - Comments (`#` and `//`)
  - Request separators (`###`)
  - Multipart form data
  - Response handlers `> {% script %}`
  - Response references `<> response.json`
  - Unicode support for URLs, paths, queries

**New Package Structure:**
```
pkg/httprequest/
  ├── parser.go       # Main parser implementing the spec grammar
  ├── lexer.go        # Tokenizer for HTTP request format
  ├── types.go        # AST types matching spec grammar
  ├── validator.go    # Validation according to spec rules
  ├── executor.go     # Request execution engine
  └── parser_test.go  # Comprehensive tests covering all spec examples
```

### Phase 2: JetBrains Environment System

**Environment File Format (JetBrains Standard)**
Following JetBrains format exactly, create environment files:

```json
// http-client.env.json
{
  "development": {
    "baseUrl": "https://api-dev.example.com",
    "apiKey": "dev-key-123",
    "timeout": 30000,
    "debugMode": true
  },
  "production": {
    "baseUrl": "https://api.example.com", 
    "apiKey": "prod-key-456",
    "timeout": 10000,
    "debugMode": false
  },
  "testing": {
    "baseUrl": "https://api-test.example.com",
    "apiKey": "test-key-789",
    "useMockData": true
  }
}
```

**Private Environment File:**
```json
// http-client.private.env.json (for sensitive data, git-ignored)
{
  "development": {
    "secretToken": "dev-secret-xyz",
    "databasePassword": "dev-db-pass"
  },
  "production": {
    "secretToken": "{{PROD_SECRET_TOKEN}}",
    "databasePassword": "{{PROD_DB_PASSWORD}}"
  }
}
```

**Implementation:**
```
pkg/environment/
  ├── loader.go       # Load JetBrains .env.json files
  ├── resolver.go     # Variable resolution with precedence
  ├── types.go        # Environment types
  └── merger.go       # Merge public + private environments
```

### Phase 3: Request Execution Engine

**New Execution System**
- Replace existing collection runner with HTTP request executor
- Support all spec features including response handlers
- Implement JavaScript execution for response scripts
- File-based request body support

**Features:**
```go
// pkg/httprequest/executor.go
type Executor struct {
    client      *client.APIClient
    environment map[string]interface{}
    globals     map[string]interface{}
    responses   []Response // Response history
}

func (e *Executor) ExecuteFile(filename string, env string) error
func (e *Executor) ExecuteRequest(request *ParsedRequest) (*Response, error)
func (e *Executor) RunResponseHandler(handler *ResponseHandler, response *Response) error
```

### Phase 4: Response Handler System

**JavaScript Response Handlers**
Following the spec for response handler scripts:

```http
GET https://api.example.com/auth

> {% 
    client.test("Authentication successful", function() {
        client.assert(response.status === 200, "Expected 200 status");
        client.assert(response.body.token, "Token should be present");
    });
    
    // Set global variables
    client.global.set("authToken", response.body.token);
    client.global.set("userId", response.body.user.id);
%}
```

**Implementation:**
```
pkg/scripting/
  ├── engine.go       # JavaScript execution engine
  ├── client_api.go   # client.* API implementation
  ├── assertions.go   # Test assertion framework
  └── globals.go      # Global variable management
```

### Phase 5: Response History & References

**Response Storage**
```http
GET https://api.example.com/users/1

<> 2024-10-31T120000.200.json
```

- Store responses in timestamped files
- Support referencing previous responses
- Enable response comparison
- Organize by request name/path

**Implementation:**
```
pkg/responses/
  ├── storage.go      # Response file management
  ├── history.go      # Response history tracking
  └── comparison.go   # Response comparison tools
```

### Phase 6: Command Line Interface

**New CLI Commands (Replace existing collection commands):**

```bash
# Run .http file with default environment
postie run requests.http

# Run with specific environment
postie run requests.http --env production

# Run specific request by line number
postie run requests.http --request 3

# Run specific request by name (from ### comments)
postie run requests.http --name "Create User"

# List all requests in file
postie list requests.http

# Validate .http file syntax
postie validate requests.http

# Show available environments
postie env list

# Set active environment
postie env set production

# Run in watch mode (re-run on file changes)
postie run requests.http --watch
```

**Remove old collection commands:**
- `postie collection create`
- `postie collection update` 
- `postie collection show`
- `postie collection list`
- `postie collection delete`

## Implementation Timeline

### Week 1-2: Parser Foundation
- ✅ Create `pkg/httprequest` package
- ✅ Implement lexer for tokenization
- ✅ Build basic parser for simple requests
- ✅ Add comprehensive tests for core parsing

### Week 3-4: Advanced Parser Features  
- ✅ Multipart form data support
- ✅ Environment variable parsing
- ✅ Response handler parsing
- ✅ File reference support (`< ./file.json`)
- ✅ Unicode support for URLs and paths

### Week 5-6: Environment System
- ✅ JetBrains environment file loader
- ✅ Variable resolution with precedence rules
- ✅ Private environment file support
- ✅ Environment merging logic

### Week 7-8: Execution Engine
- ✅ Replace collection runner with HTTP executor
- ✅ Integrate with existing HTTP client
- ✅ Request execution with variable substitution
- ✅ Response collection and storage

### Week 9-10: Response Handlers & Scripting
- ✅ JavaScript execution engine
- ✅ `client.*` API implementation
- ✅ Test assertion framework
- ✅ Global variable management

### Week 11-12: CLI Updates & Polish
- ✅ Update all CLI commands
- ✅ Remove old collection commands
- ✅ Add environment management commands
- ✅ Comprehensive testing and documentation

## New File Structure

```
postie/
├── pkg/
│   ├── httprequest/         # NEW: HTTP Request format (replaces collection/)
│   │   ├── parser.go
│   │   ├── lexer.go
│   │   ├── types.go
│   │   ├── validator.go
│   │   ├── executor.go
│   │   └── parser_test.go
│   ├── environment/         # NEW: JetBrains environment support
│   │   ├── loader.go
│   │   ├── resolver.go
│   │   ├── types.go
│   │   └── merger.go
│   ├── scripting/           # NEW: Response handler scripting
│   │   ├── engine.go
│   │   ├── client_api.go
│   │   ├── assertions.go
│   │   └── globals.go
│   ├── responses/           # NEW: Response history & references
│   │   ├── storage.go
│   │   ├── history.go
│   │   └── comparison.go
│   └── client/              # Existing HTTP client (unchanged)
│       └── client.go
├── examples/
│   ├── requests/            # NEW: .http examples
│   │   ├── basic.http
│   │   ├── auth.http
│   │   ├── advanced.http
│   │   └── multipart.http
│   └── environments/        # NEW: Environment examples
│       ├── http-client.env.json
│       └── http-client.private.env.json
└── docs/
    ├── http-request-format.md
    └── migration-guide.md
```

## Example: Complete Migration

### Before (JSON Collection) - TO BE REMOVED
```json
{
  "collection": {
    "info": {
      "name": "API Tests"
    },
    "variable": [
      {"key": "baseUrl", "value": "https://api.example.com"},
      {"key": "userId", "value": "123"}
    ],
    "apiGroup": [{
      "name": "Get User",
      "request": {
        "method": "GET",
        "url": "{{baseUrl}}/users/{{userId}}",
        "header": [
          {"key": "Authorization", "value": "Bearer {{token}}"}
        ]
      }
    }]
  }
}
```

### After (HTTP Request Format) - NEW IMPLEMENTATION
```http
### Get User
GET {{baseUrl}}/users/{{userId}}
Authorization: Bearer {{token}}
Accept: application/json

> {%
    client.test("User retrieved successfully", function() {
        client.assert(response.status === 200, "Expected 200 status");
        client.assert(response.body.id === client.global.get("userId"), "User ID should match");
    });
    
    // Store user data for subsequent requests
    client.global.set("userName", response.body.name);
    client.global.set("userEmail", response.body.email);
%}

### Update User
PUT {{baseUrl}}/users/{{userId}}
Authorization: Bearer {{token}}
Content-Type: application/json

{
  "name": "{{userName}}",
  "email": "updated@example.com"
}

> {%
    client.test("User updated successfully", function() {
        client.assert(response.status === 200, "Expected 200 status");
    });
%}
```

### Environment File (JetBrains Format)
```json
{
  "development": {
    "baseUrl": "https://api-dev.example.com",
    "token": "dev-token-123",
    "userId": "dev-user-456"
  },
  "production": {
    "baseUrl": "https://api.example.com",
    "token": "{{PROD_TOKEN}}",
    "userId": "prod-user-789"
  }
}
```

## Breaking Changes

**Complete Format Change:**
- ❌ JSON collections no longer supported
- ❌ All existing `.collection.json` files must be manually converted
- ❌ Collection-based CLI commands removed
- ✅ New `.http` format required
- ✅ JetBrains environment format required
- ✅ New CLI commands for HTTP requests

**Migration Requirements:**
1. Convert all existing collections to `.http` format manually
2. Convert environment configurations to `http-client.env.json`
3. Update any automation scripts to use new CLI commands
4. Update documentation and examples

## Advantages of Fresh Implementation

1. **Standards Compliance**: 100% compatible with JetBrains HTTP Request in Editor spec
2. **IDE Integration**: Works seamlessly with IntelliJ IDEA, WebStorm, etc.
3. **Cleaner Codebase**: No legacy support means simpler, more maintainable code
4. **Better Performance**: Optimized for text-based format from the ground up
5. **Version Control**: Superior diff and merge capabilities
6. **Editor Agnostic**: Works in any text editor with proper syntax highlighting
7. **Response Handling**: Built-in scripting with full JavaScript support
8. **Industry Standard**: Follows established conventions used by major IDEs

## Risks & Mitigation

| Risk | Impact | Mitigation |
|------|--------|------------|
| User migration required | HIGH | Provide clear migration guide and examples |
| Breaking all existing workflows | HIGH | Comprehensive documentation and migration tools |
| Parser complexity | MEDIUM | Follow spec exactly, extensive testing |
| JavaScript security concerns | MEDIUM | Sandboxed execution environment |
| Loss of existing collections | HIGH | Document manual conversion process clearly |

## Recommendation

**Proceed with Fresh Implementation:**

1. ✅ **Week 1-4**: Build robust HTTP Request parser
2. ✅ **Week 5-6**: Implement JetBrains environment system  
3. ✅ **Week 7-8**: Create new execution engine
4. ✅ **Week 9-10**: Add response handlers and scripting
5. ✅ **Week 11-12**: Update CLI and comprehensive testing
6. 🚀 **Release**: Complete format migration with migration guide

**Key Benefits:**
- Clean, standards-compliant implementation
- Full JetBrains IDE compatibility
- Modern text-based format
- Superior version control support
- Built-in scripting capabilities
- Industry-standard environment management

The fresh implementation approach provides a cleaner codebase and ensures 100% compatibility with the HTTP Request in Editor specification, positioning Postie as a modern, standards-compliant API testing tool.