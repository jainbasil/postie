# Postie Command Reference

Complete reference for all Postie CLI commands with examples and detailed explanations.

## Table of Contents

1. [HTTP Commands](#http-commands)
2. [Environment Management](#environment-management)
3. [Context Management](#context-management)
4. [Utility Commands](#utility-commands)

---

## HTTP Commands

Execute HTTP requests from `.http` files using the JetBrains HTTP Client format.

### `postie http run`

Execute HTTP requests from a `.http` file.

**Usage:**
```bash
postie http run [<file.http>] [options]
```

**Options:**
- `--env, -e` (optional): Environment name (default: development)
- `--env-file` (optional): Path to environment file (default: http-client.env.json)
- `--private-env-file` (optional): Path to private environment file (default: http-client.private.env.json)
- `--request, -r` (optional): Run specific request by name or number
- `--verbose, -v` (optional): Show detailed output
- `--save-responses, -s` (optional): Save responses to `.http-responses/` directory

**Examples:**
```bash
# Run all requests in a file
postie http run requests.http

# Run with specific environment
postie http run requests.http --env production

# Run specific request by name
postie http run requests.http --request "Get Users"

# Run specific request by number
postie http run requests.http --request 1

# Run with verbose output
postie http run requests.http --verbose

# Save responses to files
postie http run requests.http --save-responses

# Using context (no file needed if context is set)
postie http run --request getUserById

# Override environment from context
postie http run --env staging
```

**Output:**
```
========== Request 1: GET https://api.example.com/users ====================
Name: Get all users

✓ Status: 200 OK
  Duration: 234.567ms
  Size: 1523 bytes
  Content-Type: application/json; charset=utf-8

Response Body:
[
  {
    "id": 1,
    "name": "John Doe",
    "email": "john@example.com"
  },
  ...
]
```

---

### `postie http parse`

Parse and validate HTTP request files without executing them.

**Usage:**
```bash
postie http parse <file.http> [options]
```

**Options:**
- `--format, -f` (optional): Output format (summary, json) - default: summary
- `--validate` (optional): Perform validation checks

**Examples:**
```bash
# Parse and show summary
postie http parse requests.http

# Parse with JSON output
postie http parse requests.http --format json

# Parse with validation
postie http parse requests.http --validate
```

**Output (summary format):**
```
HTTP Request File: requests.http
Total Requests: 5

Request 1: Get All Users
  Method: GET
  URL: https://api.example.com/users
  Headers: 2
  
Request 2: Create User
  Method: POST
  URL: https://api.example.com/users
  Headers: 1
  Body: Present
  Handler Script: Yes
```

---

### `postie http list`

List all `.http` files in a directory.

**Usage:**
```bash
postie http list [directory] [options]
```

**Options:**
- `--recursive, -r` (optional): Search recursively in subdirectories

**Examples:**
```bash
# List .http files in current directory
postie http list

# List in specific directory
postie http list api-tests/

# List recursively
postie http list --recursive
```

**Output:**
```
Found 3 .http file(s):

1. requests.http (5 requests)
2. api-tests/auth.http (3 requests)
3. api-tests/users.http (7 requests)
```

---

## Environment Management

Manage environment files and inspect environment variables.

### `postie env list`

List all available environments from environment files.

**Usage:**
```bash
postie env list [options]
```

**Options:**
- `--env-file` (optional): Path to environment file (default: http-client.env.json)
- `--private-env-file` (optional): Path to private environment file (default: http-client.private.env.json)

**Examples:**
```bash
# List environments from default files
postie env list

# List from custom environment file
postie env list --env-file custom-env.json

# List with custom private file
postie env list --env-file http-client.env.json --private-env-file http-client.private.env.json
```

**Output:**
```
Available environments:
  development (4 public, 3 private variables)
  production (4 public, 3 private variables)
  staging (4 public, 3 private variables)
```

---

### `postie env show`

Show variables for a specific environment.

**Usage:**
```bash
postie env show <environment> [options]
```

**Options:**
- `--env-file` (optional): Path to environment file (default: http-client.env.json)
- `--private-env-file` (optional): Path to private environment file (default: http-client.private.env.json)
- `--show-private` (optional): Display private/sensitive variables

**Examples:**
```bash
# Show public variables for development environment
postie env show development

# Show including private variables
postie env show development --show-private

# Show from custom environment file
postie env show production --env-file custom-env.json
```

**Output:**
```
Environment: development

Public variables:
  apiVersion = "v1"
  baseUrl = "https://api-dev.example.com"
  timeout = 30000
  userAgent = "Postie/1.0.0"

Private variables: 3 (use --show-private to display)

Total resolved variables: 7
```

**Output with `--show-private`:**
```
Environment: development

Public variables:
  apiVersion = "v1"
  baseUrl = "https://api-dev.example.com"
  timeout = 30000
  userAgent = "Postie/1.0.0"

Private variables:
  apiKey = "dev-secret-key-12345"
  authToken = "dev-bearer-token-xyz"
  dbPassword = "dev-password"

Total resolved variables: 7
```

---

## Context Management

Set default HTTP files and environments for a directory to streamline your workflow.

### `postie context set`

Set context values for the current directory.

**Usage:**
```bash
postie context set [options]
```

**Options:**
- `--http-file` (optional): Default HTTP request file path
- `--env` (optional): Default environment name
- `--env-file` (optional): Path to environment file
- `--private-env-file` (optional): Path to private environment file
- `--save-responses` (optional): Enable automatic response saving
- `--responses-dir` (optional): Custom directory for saved responses

**Examples:**
```bash
# Set default HTTP file and environment
postie context set --http-file requests.http --env development

# Set with save responses enabled
postie context set --http-file api.http --env production --save-responses

# Set custom environment file paths
postie context set \
  --http-file api.http \
  --env staging \
  --env-file custom-env.json \
  --private-env-file custom-private.json

# Update only the environment (keeps existing HTTP file)
postie context set --env production
```

**Output:**
```
Context saved to /path/to/project/.postie-context.json
Context file: /path/to/project/.postie-context.json

HTTP File:         /path/to/project/requests.http
Environment:       development
```

**Context File Location:** `.postie-context.json` in the current directory

---

### `postie context show`

Display current context settings.

**Usage:**
```bash
postie context show
```

**Examples:**
```bash
# Show current context
postie context show
```

**Output:**
```
Context file: /path/to/project/.postie-context.json

HTTP File:         /path/to/project/api-tests/requests.http
Environment:       development
Save Responses:    true
```

**Output (when no context exists):**
```
No context file found in current directory.
Use 'postie context set' to create one.
```

---

### `postie context clear`

Clear context settings for the current directory.

**Usage:**
```bash
postie context clear
```

**Examples:**
```bash
# Clear context
postie context clear
```

**Output:**
```
Context cleared: /path/to/project/.postie-context.json
```

---

## Utility Commands

### `postie demo`

Run interactive demonstration of Postie features.

**Usage:**
```bash
postie demo
```

**Examples:**
```bash
# Run demo
postie demo
```

---

### `postie version`

Display version information.

**Usage:**
```bash
postie version
```

**Examples:**
```bash
# Show version
postie version
postie --version
postie -v
```

**Output:**
```
postie version 1.0.0
```

---

### `postie help`

Display help information.

**Usage:**
```bash
postie help [<resource>]
postie <resource> help
```

**Examples:**
```bash
# Show general help
postie help
postie --help
postie -h

# Show help for specific resource
postie http help
postie env help
postie context help
```

---

## Complete Workflow Example

Here's a complete workflow showing how to use Postie for API testing:

```bash
# 1. Create environment files
# Create http-client.env.json with public variables
cat > http-client.env.json << 'EOF'
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
EOF

# Create http-client.private.env.json with sensitive data
cat > http-client.private.env.json << 'EOF'
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
EOF

# 2. Create HTTP request file (requests.http)
cat > requests.http << 'EOF'
### Login
# @name login
POST {{baseUrl}}/{{apiVersion}}/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

> {%
  client.test("Login successful", function() {
    client.assert(response.status === 200);
  });
  client.global.set("authToken", response.body.token);
%}

###

### Get Users
GET {{baseUrl}}/{{apiVersion}}/users
Authorization: Bearer {{authToken}}
X-API-Key: {{apiKey}}

> {%
  client.test("Users retrieved", function() {
    client.assert(response.status === 200);
    client.assert(Array.isArray(response.body));
  });
%}

###

### Create User
POST {{baseUrl}}/{{apiVersion}}/users
Authorization: Bearer {{authToken}}
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}

> {%
  client.test("User created", function() {
    client.assert(response.status === 201);
  });
  client.global.set("userId", response.body.id);
%}
EOF

# 3. Set context for the project
postie context set --http-file requests.http --env development

# 4. List available environments
postie env list

# 5. Show environment variables
postie env show development
postie env show development --show-private

# 6. Run all requests
postie http run

# 7. Run specific request
postie http run --request 1
postie http run --request "Get Users"

# 8. Run with verbose output
postie http run --verbose

# 9. Save responses
postie http run --save-responses

# 10. Switch to production environment
postie context set --env production

# 11. Run in production
postie http run --request login

# 12. Parse file to see structure
postie http parse requests.http

# 13. List all .http files
postie http list

# 14. Show current context
postie context show

# 15. Clear context when done
postie context clear
```

---

## Tips and Best Practices

### Using Context Effectively

Context eliminates repetitive command-line arguments:

```bash
# Without context - repetitive
postie http run api/requests.http --env development --request "Get Users"
postie http run api/requests.http --env development --request "Create User"
postie http run api/requests.http --env development --request "Update User"

# With context - streamlined
postie context set --http-file api/requests.http --env development
postie http run --request "Get Users"
postie http run --request "Create User"
postie http run --request "Update User"
```

### Environment File Organization

1. **Public file** (`http-client.env.json`): Non-sensitive configuration
   - Base URLs
   - API versions
   - Timeouts
   - Feature flags

2. **Private file** (`http-client.private.env.json`): Sensitive data (add to `.gitignore`)
   - API keys
   - Auth tokens
   - Passwords
   - Secrets

### Variable Substitution

Use `{{variableName}}` in requests:

```http
### Example with variables
GET {{baseUrl}}/{{apiVersion}}/users/{{userId}}
Authorization: Bearer {{authToken}}
X-API-Key: {{apiKey}}
```

### Response Handler Scripts

Add JavaScript for testing and data extraction:

```http
### Request with handler
POST {{baseUrl}}/api/login
Content-Type: application/json

{"email": "{{email}}", "password": "{{password}}"}

> {%
  // Test the response
  client.test("Login successful", function() {
    client.assert(response.status === 200);
    client.assert(response.body.hasOwnProperty("token"));
  });
  
  // Save data for subsequent requests
  client.global.set("authToken", response.body.token);
  client.global.set("userId", response.body.user.id);
  
  // Log information
  client.log("Logged in as: " + response.body.user.email);
%}
```

### File Organization

```
project/
├── api-tests/
│   ├── requests.http              # HTTP requests
│   ├── auth.http                  # Authentication requests
│   ├── users.http                 # User API requests
│   ├── http-client.env.json       # Public environment variables
│   ├── http-client.private.env.json  # Private variables (in .gitignore)
│   └── .postie-context.json       # Context (in .gitignore)
└── .http-responses/               # Saved responses (in .gitignore)
    └── Get_Users_2025-11-01_143052.json
```

### .gitignore Recommendations

Add these to your `.gitignore`:

```gitignore
# Private environment variables
http-client.private.env.json

# HTTP response storage
.http-responses/

# Context files (local development preferences)
.postie-context.json
```

---

## Quick Reference

### Most Used Commands

```bash
# Run requests
postie http run requests.http
postie http run requests.http --env production
postie http run --request "Get Users"  # with context

# Environment management
postie env list
postie env show development
postie env show production --show-private

# Context management
postie context set --http-file requests.http --env development
postie context show
postie context clear

# Parsing and listing
postie http parse requests.http
postie http list
```

### Common Patterns

```bash
# Setup project context
postie context set --http-file api/requests.http --env development --save-responses

# Run all tests
postie http run --verbose

# Test specific endpoint
postie http run --request "Login" --verbose

# Switch environments
postie context set --env staging
postie http run

# Debug with saved responses
postie http run --save-responses --verbose
# Check .http-responses/ directory for detailed response data
```

---

## Error Handling

Postie provides clear error messages:

```bash
# Missing HTTP file
$ postie http run
Error: HTTP request file required
Usage: postie http run <file.http> [--env development] [--request name_or_number]
Or use 'postie context set --http-file <file>' to set a default

# File not found
$ postie http run nonexistent.http
Error: failed to read HTTP file: open nonexistent.http: no such file or directory

# Invalid environment
$ postie env show nonexistent
Error: environment 'nonexistent' not found

# Parse error in .http file
$ postie http run bad-syntax.http
Error: failed to parse HTTP file: syntax error at line 5: invalid request format
```

---

## Support and Resources

- **GitHub Repository**: https://github.com/jainbasil/postie
- **User Guide**: See [docs/user-guide.md](user-guide.md) for comprehensive examples
- **HTTP Request Format**: Follows JetBrains HTTP Client format specification
- **Issue Tracker**: Report bugs and request features on GitHub

---

*Last Updated: November 1, 2025*

---

