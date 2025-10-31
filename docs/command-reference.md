# Postie Command Reference

Complete reference for all Postie CLI commands with examples and detailed explanations.

## Table of Contents

1. [Collection Management](#collection-management)
2. [Environment Management](#environment-management)
3. [Request Group Management](#request-group-management)
4. [Request Management](#request-management)
5. [Context Management](#context-management)
6. [HTTP Commands](#http-commands)
7. [Utility Commands](#utility-commands)

---

## Collection Management

Collections are the primary organizational unit in Postie, containing API groups, environments, and requests.

### `postie collection create`

Create a new API collection.

**Usage:**
```bash
postie collection create --name <name> [options]
```

**Options:**
- `--name, -n` (required): Collection name
- `--file, -f` (optional): Output file path (auto-generated from name if not provided)
- `--description, -d` (optional): Collection description
- `--set-context` (optional): Set this collection as the current context

**Examples:**
```bash
# Create a collection with auto-generated filename
postie collection create --name "Blog API"
# Creates: blog-api.collection.json

# Create a collection with custom filename
postie collection create --name "My API Collection" --file my-api.collection.json

# Create and set as context
postie collection create --name "E-commerce API" --file ecommerce.collection.json --set-context true

# Create with description
postie collection create --name "User Service" --description "User management API endpoints"
```

**Output:**
```
✅ Collection 'Blog API' created successfully
📁 File: blog-api.collection.json
📌 Collection set as current context
```

---

### `postie collection update`

Update an existing collection's metadata.

**Usage:**
```bash
postie collection update [--file <file>] [options]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--name, -n` (optional): New collection name
- `--description, -d` (optional): New collection description
- `--set-context` (optional): Set this collection as the current context

**Examples:**
```bash
# Update collection name
postie collection update --file my-api.collection.json --name "Updated API Name"

# Update description using context
postie collection update --description "New description for the API"

# Update and set as context
postie collection update --file api.collection.json --name "Production API" --set-context true
```

---

### `postie collection show`

Display detailed information about a collection.

**Usage:**
```bash
postie collection show [--file <file>] [options]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--format` (optional): Output format (json, table, yaml) - default is table

**Examples:**
```bash
# Show collection details
postie collection show --file my-api.collection.json

# Show using context
postie collection show

# Show in JSON format
postie collection show --file my-api.collection.json --format json
```

**Output:**
```
Collection: Blog API
Description: API collection for Blog API
Version: 1.0.0
File: blog-api.collection.json

API Groups: 2
Environments: 2
Variables: 0
```

---

### `postie collection list`

List all collections in a directory.

**Usage:**
```bash
postie collection list [options]
```

**Options:**
- `--directory, -d` (optional): Directory to search (default: current directory)
- `--recursive, -r` (optional): Search recursively in subdirectories

**Examples:**
```bash
# List collections in current directory
postie collection list

# List collections in specific directory
postie collection list --directory ./collections

# List collections recursively
postie collection list --directory ./projects --recursive
```

**Output:**
```
Found 2 collection(s):

1. Blog API
   File: blog-api.collection.json
   Description: API collection for Blog API
   Groups: 2, Environments: 2

2. E-commerce API
   File: ecommerce.collection.json
   Groups: 3, Environments: 1
```

---

### `postie collection delete`

Delete a collection file.

**Usage:**
```bash
postie collection delete --file <file> [options]
```

**Options:**
- `--file, -f` (required): Collection file path
- `--force` (optional): Skip confirmation prompt

**Examples:**
```bash
# Delete with confirmation
postie collection delete --file old-api.collection.json

# Delete without confirmation
postie collection delete --file old-api.collection.json --force
```

---

## Environment Management

Environments allow you to manage different configurations (development, staging, production) within a single collection.

### `postie environment create`

Create a new environment in a collection.

**Usage:**
```bash
postie environment create --name <name> [options]
```

**Options:**
- `--name, -n` (required): Environment name
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--description, -d` (optional): Environment description
- `--set-context` (optional): Set this environment as current in context

**Examples:**
```bash
# Create environment in context collection
postie environment create --name "Development"

# Create environment in specific collection
postie environment create --name "Production" --file my-api.collection.json

# Create and set as context
postie environment create --name "Staging" --set-context true

# Create with description
postie environment create --name "Testing" --description "QA testing environment"
```

---

### `postie environment update`

Update an existing environment.

**Usage:**
```bash
postie environment update --name <name> [options]
```

**Options:**
- `--name, -n` (required): Environment name
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--description, -d` (optional): New description
- `--set-context` (optional): Set as current environment in context

**Examples:**
```bash
# Update environment description
postie environment update --name "Development" --description "Local development environment"

# Update and set as context
postie environment update --name "Production" --set-context true
```

---

### `postie environment list`

List all environments in a collection.

**Usage:**
```bash
postie environment list [--file <file>]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)

**Examples:**
```bash
# List environments using context
postie environment list

# List environments in specific collection
postie environment list --file my-api.collection.json
```

**Output:**
```
Collection: Blog API
Environments (3):

1. Development (default)
   Description: Local development environment
   Variables: 3
   Authentication: bearer

2. Staging
   Variables: 3
   Authentication: bearer

3. Production
   Variables: 2
   Authentication: apikey
```

---

### `postie environment delete`

Delete an environment from a collection.

**Usage:**
```bash
postie environment delete --name <name> [options]
```

**Options:**
- `--name, -n` (required): Environment name
- `--file, -f` (optional): Collection file path (uses context if not provided)

**Examples:**
```bash
# Delete environment
postie environment delete --name "Staging"

# Delete from specific collection
postie environment delete --name "Testing" --file my-api.collection.json
```

---

### `postie environment variable set`

Set a variable in an environment.

**Usage:**
```bash
postie environment variable set --name <env-name> --key <key> --value <value> [options]
```

**Options:**
- `--name, -n` (required): Environment name
- `--key, -k` (required): Variable key
- `--value, -v` (required): Variable value
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--secret` (optional): Mark as secret/sensitive variable

**Examples:**
```bash
# Set a simple variable
postie environment variable set --name "Development" --key "baseUrl" --value "http://localhost:3000"

# Set an API key (marked as secret)
postie environment variable set --name "Production" --key "apiKey" --value "sk_live_abc123" --secret true

# Set multiple variables
postie environment variable set --name "Development" --key "apiVersion" --value "v1"
postie environment variable set --name "Development" --key "timeout" --value "30"
```

---

### `postie environment variable get`

Get a variable value from an environment.

**Usage:**
```bash
postie environment variable get --name <env-name> [--key <key>] [options]
```

**Options:**
- `--name, -n` (required): Environment name
- `--key, -k` (optional): Specific variable key (omit to list all)
- `--file, -f` (optional): Collection file path (uses context if not provided)

**Examples:**
```bash
# Get specific variable
postie environment variable get --name "Development" --key "baseUrl"

# List all variables in environment
postie environment variable get --name "Development"
```

**Output:**
```
baseUrl=http://localhost:3000
```

---

### `postie environment variable list`

List all variables in an environment (alias for `get` without `--key`).

**Usage:**
```bash
postie environment variable list --name <env-name> [options]
```

**Examples:**
```bash
postie environment variable list --name "Development"
```

**Output:**
```
Variables in environment 'Development':

✓ baseUrl = http://localhost:3000
✓ apiKey = dev-key-123
✓ timeout = 30
```

---

## Request Group Management

Request groups organize related API requests together.

### `postie request-group create`

Create a new request group.

**Usage:**
```bash
postie request-group create --name <name> [options]
```

**Options:**
- `--name, -n` (required): Group name
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--id` (optional): Custom group ID (auto-generated from name if not provided)
- `--description, -d` (optional): Group description

**Examples:**
```bash
# Create a group with auto-generated ID
postie request-group create --name "Users"

# Create a group with custom ID
postie request-group create --name "Posts" --id "blog-posts"

# Create with description
postie request-group create --name "Authentication" --description "Auth and authorization endpoints"

# Create in specific collection
postie request-group create --name "Products" --file ecommerce.collection.json
```

---

### `postie request-group update`

Update an existing request group.

**Usage:**
```bash
postie request-group update --id <id> [options]
```

**Options:**
- `--id` (required): Group ID
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--name, -n` (optional): New group name
- `--description, -d` (optional): New description

**Examples:**
```bash
# Update group name
postie request-group update --id "users" --name "User Management"

# Update description
postie request-group update --id "posts" --description "Blog post CRUD operations"

# Update both
postie request-group update --id "auth" --name "Security" --description "Authentication and security"
```

---

### `postie request-group list`

List all request groups in a collection.

**Usage:**
```bash
postie request-group list [--file <file>]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)

**Examples:**
```bash
# List groups using context
postie request-group list

# List groups in specific collection
postie request-group list --file my-api.collection.json
```

**Output:**
```
Collection: Blog API
Request Groups (3):

1. Users
   ID: users
   Description: User management endpoints
   Requests: 5

2. Posts
   ID: posts
   Description: Blog post operations
   Requests: 7

3. Authentication
   ID: authentication
   Description: Auth endpoints
   Requests: 3
```

---

### `postie request-group delete`

Delete a request group.

**Usage:**
```bash
postie request-group delete --id <id> [options]
```

**Options:**
- `--id` (required): Group ID
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--force` (optional): Skip confirmation prompt

**Examples:**
```bash
# Delete with confirmation
postie request-group delete --id "old-apis"

# Delete without confirmation
postie request-group delete --id "deprecated" --force
```

---

## Request Management

Requests are individual API endpoints within a collection.

### `postie request create`

Create a new API request.

**Usage:**
```bash
postie request create --name <name> --method <method> --url <url> --group <group-id> [options]
```

**Options:**
- `--name, -n` (required): Request name
- `--method, -m` (required): HTTP method (GET, POST, PUT, DELETE, PATCH, etc.)
- `--url, -u` (required): Request URL (supports variables like `{{baseUrl}}`)
- `--group, -g` (required): Group ID to add the request to
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--id` (optional): Custom request ID (auto-generated from name if not provided)
- `--description, -d` (optional): Request description
- `--body` (optional): Request body (JSON string)
- `--header` (optional, repeatable): Add header (format: "Key: Value")

**Examples:**
```bash
# Create a simple GET request
postie request create --name "Get All Users" --method GET --url "{{baseUrl}}/users" --group "users"

# Create a POST request with body
postie request create --name "Create User" --method POST --url "{{baseUrl}}/users" --group "users" \
  --body '{"name":"{{userName}}","email":"{{userEmail}}"}'

# Create with custom ID
postie request create --name "Update Post" --method PUT --url "{{baseUrl}}/posts/{{postId}}" \
  --group "posts" --id "update-blog-post"

# Create with description
postie request create --name "Delete User" --method DELETE --url "{{baseUrl}}/users/{{userId}}" \
  --group "users" --description "Permanently delete a user account"
```

---

### `postie request update`

Update an existing request.

**Usage:**
```bash
postie request update --id <id> [options]
```

**Options:**
- `--id` (required): Request ID
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--name, -n` (optional): New request name
- `--method, -m` (optional): New HTTP method
- `--url, -u` (optional): New URL
- `--description, -d` (optional): New description
- `--body` (optional): Update request body

**Examples:**
```bash
# Update request name
postie request update --id "get-users" --name "Fetch All Users"

# Update method and URL
postie request update --id "create-user" --method PATCH --url "{{baseUrl}}/api/v2/users"

# Update request body
postie request update --id "update-post" --body '{"title":"{{title}}","content":"{{content}}"}'
```

---

### `postie request list`

List all requests in a collection.

**Usage:**
```bash
postie request list [options]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--group, -g` (optional): Filter by group ID
- `--environment, -e` (optional): Show with environment variables resolved

**Examples:**
```bash
# List all requests
postie request list

# List requests in specific group
postie request list --group "users"

# List with environment variables resolved
postie request list --environment "Development"

# List from specific collection
postie request list --file my-api.collection.json
```

**Output:**
```
Collection: Blog API
Environment: Development
Requests (5):

1. Users / Get All Users
   ID: get-all-users
   GET http://localhost:3000/users

2. Users / Create User
   ID: create-user
   POST http://localhost:3000/users

3. Posts / Get All Posts
   ID: get-all-posts
   GET http://localhost:3000/posts
```

---

### `postie request show`

Show detailed information about a specific request.

**Usage:**
```bash
postie request show [--id <id> | --name <name>] [options]
```

**Options:**
- `--id` (optional): Request ID (use this OR --name)
- `--name, -n` (optional): Request name (use this OR --id)
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--environment, -e` (optional): Show with environment variables resolved

**Examples:**
```bash
# Show request by ID
postie request show --id "get-users"

# Show request by name
postie request show --name "Get All Users"

# Show with environment variables resolved
postie request show --id "create-user" --environment "Production"
```

**Output:**
```
Request: Get All Users
ID: get-all-users
Group: Users
Description: Retrieve all user accounts

Method: GET
URL: http://localhost:3000/users

Headers:
  Authorization: Bearer dev-token-123
  Content-Type: application/json
```

---

### `postie request delete`

Delete a request from a collection.

**Usage:**
```bash
postie request delete --id <id> [options]
```

**Options:**
- `--id` (required): Request ID
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--force` (optional): Skip confirmation prompt

**Examples:**
```bash
# Delete with confirmation
postie request delete --id "deprecated-api"

# Delete without confirmation
postie request delete --id "old-endpoint" --force
```

---

### `postie request run`

Execute a specific request.

**Usage:**
```bash
postie request run [--id <id> | --name <name>] [options]
```

**Options:**
- `--id` (optional): Request ID (use this OR --name)
- `--name, -n` (optional): Request name (use this OR --id)
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--environment, -e` (optional): Environment to use (uses context if not provided)

**Examples:**
```bash
# Run request by ID
postie request run --id "get-users"

# Run request by name
postie request run --name "Create User" --environment "Development"

# Run with specific environment
postie request run --id "update-post" --environment "Production"
```

**Output:**
```
Running request: Users / Get All Users

==================================================
Status: 200 OK
Duration: 234ms
Size: 1523 bytes
Content-Type: application/json
==================================================
✅ Request successful

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

### `postie request run-all`

Run all requests in a collection.

**Usage:**
```bash
postie request run-all [options]
```

**Options:**
- `--file, -f` (optional): Collection file path (uses context if not provided)
- `--environment, -e` (optional): Environment to use (uses context if not provided)
- `--group, -g` (optional): Run only requests in specific group

**Examples:**
```bash
# Run all requests with context
postie request run-all

# Run all requests in specific environment
postie request run-all --environment "Development"

# Run all requests in a group
postie request run-all --group "users"

# Run from specific collection
postie request run-all --file my-api.collection.json --environment "Staging"
```

**Output:**
```
Running collection: Blog API
Environment: Development
Found 5 requests

[1/5] Users / Get All Users
  ✅ 200 OK (234ms)
  📦 1523 bytes

[2/5] Users / Create User
  ✅ 201 Created (156ms)
  📦 245 bytes

[3/5] Posts / Get All Posts
  ✅ 200 OK (189ms)
  📦 3421 bytes
```

---

## Context Management

Context allows you to set default collection and environment to avoid typing them repeatedly.

### `postie context set`

Set the default collection and/or environment.

**Usage:**
```bash
postie context set [--collection <file>] [--environment <name>]
```

**Options:**
- `--collection, -c` (optional): Collection file path
- `--environment, -e` (optional): Environment name

**Examples:**
```bash
# Set collection and environment
postie context set --collection blog-api.collection.json --environment "Development"

# Set only collection
postie context set --collection my-api.collection.json

# Set only environment (requires collection to be set first)
postie context set --environment "Production"

# Update context
postie context set --collection new-api.collection.json --environment "Staging"
```

**Context File Location:** `~/.postie/context.json`

---

### `postie context show`

Display the current context.

**Usage:**
```bash
postie context show [--format <format>]
```

**Options:**
- `--format` (optional): Output format (table, json, yaml)

**Examples:**
```bash
# Show context in table format
postie context show

# Show context in JSON format
postie context show --format json
```

**Output:**
```
Current Context:
==================================================
Collection:  /home/user/blog-api.collection.json
Name:        Blog API
Description: API collection for Blog API
Environment: Development
==================================================
```

---

### `postie context clear`

Clear the current context.

**Usage:**
```bash
postie context clear [--force]
```

**Options:**
- `--force` (optional): Skip confirmation prompt

**Examples:**
```bash
# Clear with confirmation
postie context clear

# Clear without confirmation
postie context clear --force
```

---

## HTTP Commands

Direct HTTP requests without collections - useful for quick API testing.

### `postie http get`

Send a GET request.

**Usage:**
```bash
postie http get --url <url> [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--header, -H` (optional, repeatable): Add headers
- `--param, -p` (optional, repeatable): Add query parameters

**Examples:**
```bash
# Simple GET request
postie http get --url "https://api.github.com/users/octocat"

# GET request with headers
postie http get --url "https://api.example.com/users" \
  --header "Authorization: Bearer token123"

# GET with query parameters
postie http get --url "https://api.example.com/search" \
  --param "q=postie" \
  --param "limit=10"
```

---

### `postie http post`

Send a POST request.

**Usage:**
```bash
postie http post --url <url> [--body <body>] [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--body, -b` (optional): Request body (JSON string or @filename)
- `--header, -H` (optional, repeatable): Add headers
- `--param, -p` (optional, repeatable): Add query parameters

**Examples:**
```bash
# POST with JSON body
postie http post --url "https://api.example.com/users" \
  --body '{"name":"John","email":"john@example.com"}'

# POST with body from file
postie http post --url "https://httpbin.org/post" \
  --body @request.json \
  --header "Content-Type: application/json"

# POST with headers
postie http post --url "https://api.example.com/posts" \
  --body '{"title":"Hello World"}' \
  --header "Authorization: Bearer token" \
  --header "Content-Type: application/json"
```

---

### `postie http put`

Send a PUT request.

**Usage:**
```bash
postie http put --url <url> [--body <body>] [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--body, -b` (optional): Request body (JSON string or @filename)
- `--header, -H` (optional, repeatable): Add headers

**Examples:**
```bash
# PUT request
postie http put --url "https://api.example.com/users/123" \
  --body '{"name":"Jane Doe"}'

# PUT with headers
postie http put --url "https://api.example.com/posts/456" \
  --body '{"title":"Updated Title"}' \
  --header "Authorization: Bearer token"
```

---

### `postie http delete`

Send a DELETE request.

**Usage:**
```bash
postie http delete --url <url> [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--header, -H` (optional, repeatable): Add headers

**Examples:**
```bash
# Simple DELETE
postie http delete --url "https://api.example.com/users/123"

# DELETE with authorization
postie http delete --url "https://api.example.com/posts/456" \
  --header "Authorization: Bearer token"
```

---

### `postie http patch`

Send a PATCH request.

**Usage:**
```bash
postie http patch --url <url> [--body <body>] [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--body, -b` (optional): Request body (JSON string or @filename)
- `--header, -H` (optional, repeatable): Add headers

**Examples:**
```bash
# PATCH request
postie http patch --url "https://api.example.com/users/123" \
  --body '{"email":"newemail@example.com"}'
```

---

### `postie http head`

Send a HEAD request.

**Usage:**
```bash
postie http head --url <url> [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--header, -H` (optional, repeatable): Add headers

**Examples:**
```bash
# HEAD request to check resource existence
postie http head --url "https://api.example.com/users/123"
```

---

### `postie http options`

Send an OPTIONS request.

**Usage:**
```bash
postie http options --url <url> [options]
```

**Options:**
- `--url, -u` (required): Request URL
- `--header, -H` (optional, repeatable): Add headers

**Examples:**
```bash
# OPTIONS request to check allowed methods
postie http options --url "https://api.example.com/users"
```

---

## Utility Commands

### `postie demo`

Run interactive demonstration of Postie features.

**Usage:**
```bash
postie demo [--examples <type>]
```

**Options:**
- `--examples` (optional): Which demo examples to run (all, basic, authentication, collections)

**Examples:**
```bash
# Run all demos
postie demo

# Run specific demo
postie demo --examples authentication
```

---

### `postie version`

Display version information.

**Usage:**
```bash
postie version [--format <format>]
```

**Options:**
- `--format` (optional): Output format (text, json)

**Examples:**
```bash
# Show version
postie version

# Show version in JSON
postie version --format json
```

---

### `postie help`

Display help information.

**Usage:**
```bash
postie help [<resource>]
```

**Examples:**
```bash
# Show general help
postie help

# Show help for specific resource
postie collection help
postie request help
postie http help
```

---

## Complete Workflow Example

Here's a complete workflow showing how to use Postie from creating a collection to running requests:

```bash
# 1. Create a new collection and set as context
postie collection create --name "Blog API" --file blog-api.collection.json --set-context true

# 2. Create environments
postie environment create --name "Development"
postie environment create --name "Production"

# 3. Set variables for development
postie environment variable set --name "Development" --key "baseUrl" --value "http://localhost:3000"
postie environment variable set --name "Development" --key "apiKey" --value "dev-key-123"

# 4. Set variables for production
postie environment variable set --name "Production" --key "baseUrl" --value "https://api.myblog.com"
postie environment variable set --name "Production" --key "apiKey" --value "prod-key-456"

# 5. Create API groups
postie request-group create --name "Users" --description "User management endpoints"
postie request-group create --name "Posts" --description "Blog post endpoints"

# 6. Create requests for Users group
postie request create --name "Get All Users" --method GET --url "{{baseUrl}}/users" --group "users"
postie request create --name "Get User by ID" --method GET --url "{{baseUrl}}/users/{{userId}}" --group "users"
postie request create --name "Create User" --method POST --url "{{baseUrl}}/users" --group "users" \
  --body '{"name":"{{userName}}","email":"{{userEmail}}"}'

# 7. Create requests for Posts group
postie request create --name "Get All Posts" --method GET --url "{{baseUrl}}/posts" --group "posts"
postie request create --name "Create Post" --method POST --url "{{baseUrl}}/posts" --group "posts" \
  --body '{"title":"{{postTitle}}","content":"{{postContent}}"}'

# 8. Set context to development environment
postie context set --environment "Development"

# 9. Run all requests in development
postie request run-all

# 10. Run specific request
postie request run --name "Get All Users"

# 11. Switch to production and run
postie context set --environment "Production"
postie request run --id "get-all-posts"

# 12. List all requests
postie request list

# 13. Show current context
postie context show

# 14. Make a quick HTTP request outside collections
postie http get --url "https://api.github.com/users/octocat"
```

---

## Tips and Best Practices

### Using Context Effectively

Context is a powerful feature that saves you from typing collection and environment paths repeatedly:

```bash
# Set once
postie context set --collection my-api.collection.json --environment "Development"

# Use everywhere without specifying file
postie request list
postie request run --id "get-users"
postie environment variable set --name "Development" --key "newKey" --value "newValue"
```

### Variable Substitution

Use `{{variableName}}` in URLs and request bodies to substitute environment-specific values:

```bash
# In URL
--url "{{baseUrl}}/api/{{apiVersion}}/users"

# In request body
--body '{"apiKey":"{{apiKey}}","endpoint":"{{baseUrl}}"}'
```

### Organization Best Practices

1. **Use meaningful group names**: Organize requests by resource or feature
2. **Consistent naming**: Use clear, descriptive names for requests
3. **Environment separation**: Keep dev/staging/prod configurations separate
4. **Variable usage**: Use variables for all environment-specific values
5. **Set context**: Use context when working on a single project for extended periods

### File Organization

```
project/
├── collections/
│   ├── api.collection.json          # Main collection
│   ├── development.environment.json # Dev environment (optional)
│   └── production.environment.json  # Prod environment (optional)
└── .postie/
    └── context.json                 # User context (in home directory)
```

---

## Quick Reference

### Most Used Commands

```bash
# Collections
postie collection create --name "API" --file api.collection.json
postie collection list

# Groups
postie request-group create --name "Users"
postie request-group list

# Requests
postie request create --name "Get Users" --method GET --url "{{baseUrl}}/users" --group "users"
postie request run --id "get-users"
postie request run-all

# Context
postie context set --collection api.collection.json --environment "Development"
postie context show

# HTTP
postie http get --url "https://api.example.com/endpoint"
postie http post --url "https://api.example.com/endpoint" --body '{"key":"value"}'
```

### Common Patterns

```bash
# Create complete API collection
postie collection create --name "My API" --file my-api.collection.json --set-context true
postie environment create --name "Development"
postie request-group create --name "Users"
postie request create --name "List Users" --method GET --url "{{baseUrl}}/users" --group "users"

# Run with different environments
postie request run-all --environment "Development"
postie request run-all --environment "Production"

# Quick testing without collection
postie http get --url "https://api.example.com/test"
postie http post --url "https://api.example.com/test" --body '{"test":"data"}'
```

---

## Error Handling

Postie provides clear error messages to help you troubleshoot issues:

```bash
# Missing required flag
$ postie collection create
Error: required flag --name not provided

# Collection not found
$ postie request list --file nonexistent.json
Error: error loading collection: failed to read collection file: open nonexistent.json: no such file or directory

# No context set
$ postie request list
Error: no collection file specified and no context set

# Invalid environment
$ postie context set --environment "NonExistent"
Error: environment 'NonExistent' not found
Available environments: Development, Production, Staging
```

---

## Support and Resources

- **GitHub Repository**: https://github.com/jainbasil/postie
- **Documentation**: See README.md and docs/ directory
- **Issue Tracker**: Report bugs and request features on GitHub
- **Collection Format**: See `docs/collection-format.md` for JSON structure

---

*Last Updated: October 31, 2025*
