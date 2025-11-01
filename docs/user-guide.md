# Postie User Guide

## Overview

Postie is a powerful API testing tool that supports the HTTP Request in Editor format, making it easy to write, organize, and execute HTTP requests directly from `.http` files. With built-in support for environment variables, response handler scripts, and global variable sharing, Postie provides a complete solution for API testing and development.

## Table of Contents

- [Getting Started](#getting-started)
- [Writing HTTP Requests](#writing-http-requests)
- [Context Management](#context-management)
- [Environment Variables](#environment-variables)
- [Response Handler Scripts](#response-handler-scripts)
- [Global Variables](#global-variables)
- [Command Reference](#command-reference)
- [Examples](#examples)

## Getting Started

### Installation

```bash
# Build from source
go build -o postie .

# Run a request file
./postie http run api-tests/sample-requests.http
```

### Basic Usage

Create a file with `.http` extension and write HTTP requests:

```http
### Simple GET request
GET https://api.example.com/users

### POST request with JSON body
POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com"
}
```

Run the requests:

```bash
postie http run requests.http
```

## Writing HTTP Requests

### Request Format

Requests follow the HTTP Request in Editor format:

```http
### Request Name (optional comment)
# @name requestName
METHOD URL
Header-Name: Header-Value
Another-Header: Value

Request body (for POST/PUT/PATCH)
```

### Request Separators

Use `###` to separate multiple requests in a file:

```http
### First Request
GET https://api.example.com/posts

### Second Request
POST https://api.example.com/posts
Content-Type: application/json

{"title": "New Post"}
```

### Supported HTTP Methods

- GET
- POST
- PUT
- DELETE
- PATCH
- HEAD
- OPTIONS

### Headers

Add headers after the request line:

```http
GET https://api.example.com/data
Accept: application/json
Authorization: Bearer {{authToken}}
User-Agent: Postie/1.0
```

### Request Body

For POST, PUT, and PATCH requests, add the body after headers:

```http
POST https://api.example.com/users
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "age": 30
}
```

## Context Management

Context management allows you to set default values for HTTP files and environments in a specific directory, eliminating the need to specify them with every command.

### Setting Context

Use the `context set` command to configure defaults for the current directory:

```bash
# Set default HTTP file and environment
postie context set --http-file api-tests/requests.http --env development

# Include save responses option
postie context set --http-file requests.http --env production --save-responses

# Set custom environment file paths
postie context set \
  --http-file api.http \
  --env staging \
  --env-file custom-env.json \
  --private-env-file custom-private.json
```

Available options:
- `--http-file`: Default HTTP request file
- `--env`: Default environment name
- `--env-file`: Path to public environment file
- `--private-env-file`: Path to private environment file
- `--save-responses`: Enable automatic response saving
- `--responses-dir`: Custom directory for saved responses

### Using Context

Once context is set, you can omit common flags:

```bash
# Without context
postie http run api-tests/requests.http --env development --request getUserById

# With context (after setting http-file and env)
postie http run --request getUserById

# Context values can still be overridden
postie http run --env production
```

### Viewing Context

Check current context settings:

```bash
postie context show
```

Output example:
```
Context file: /path/to/project/.postie-context.json

HTTP File:         /path/to/project/api-tests/requests.http
Environment:       development
Save Responses:    true
```

### Clearing Context

Remove context settings:

```bash
postie context clear
```

### Context File

Context is stored in `.postie-context.json` in the current directory. This file should be added to `.gitignore` as it contains local development preferences.

Example `.postie-context.json`:
```json
{
  "httpFile": "/home/user/project/requests.http",
  "environment": "development",
  "saveResponses": true
}
```

## Environment Variables

### Environment Files

Postie uses JetBrains-style environment files for managing variables across different environments.

#### Public Environment File (`http-client.env.json`)

```json
{
  "development": {
    "baseUrl": "https://api-dev.example.com",
    "apiKey": "dev-key-123",
    "timeout": 30000
  },
  "production": {
    "baseUrl": "https://api.example.com",
    "apiKey": "prod-key-456",
    "timeout": 10000
  }
}
```

#### Private Environment File (`http-client.private.env.json`)

Store sensitive data in this file (should be git-ignored):

```json
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

### Using Variables

Use `{{variableName}}` syntax to reference variables:

```http
### Get User Profile
GET {{baseUrl}}/users/profile
Authorization: Bearer {{secretToken}}
Accept: application/json
```

### Variable Expansion

Variables are expanded in:
- URLs
- Headers
- Request bodies
- Response handler scripts

### Selecting Environments

Specify the environment when running requests:

```bash
# Use development environment (default)
postie http run requests.http --env development

# Use production environment
postie http run requests.http --env production
```

## Response Handler Scripts

Response handlers allow you to write JavaScript code that executes after receiving a response. Use them for testing, validation, and storing data for subsequent requests.

### Basic Syntax

```http
GET https://api.example.com/posts/1

> {%
    // JavaScript code runs after the response is received
    client.log("Response received");
%}
```

### Available APIs

#### `client.test(name, function)`

Define tests with assertions:

```http
GET https://api.example.com/posts/1

> {%
    client.test("Request was successful", function() {
        client.assert(response.status === 200, "Expected status 200");
    });
    
    client.test("Response contains data", function() {
        client.assert(response.body.id === 1, "Expected post ID 1");
        client.assert(response.body.title, "Post should have a title");
    });
%}
```

#### `client.assert(condition, message)`

Make inline assertions:

```http
POST https://api.example.com/users

> {%
    client.assert(response.status === 201, "Expected status 201 Created");
    client.assert(response.body.id, "Should return created user ID");
%}
```

#### `client.log(...messages)`

Log messages during script execution:

```http
GET https://api.example.com/users/1

> {%
    client.log("User ID:", response.body.id);
    client.log("User Name:", response.body.name);
    client.log("Response status:", response.status);
%}
```

#### `client.global.set(name, value)`

Store values in global variables for use in subsequent requests:

```http
### Login
POST https://api.example.com/auth/login
Content-Type: application/json

{
  "username": "user@example.com",
  "password": "password123"
}

> {%
    client.test("Login successful", function() {
        client.assert(response.status === 200, "Expected status 200");
    });
    
    // Store the authentication token
    client.global.set("authToken", response.body.token);
    client.global.set("userId", response.body.user.id);
%}

### Get User Profile (uses token from previous request)
GET https://api.example.com/users/{{userId}}
Authorization: Bearer {{authToken}}
```

#### `client.global.get(name)`

Retrieve previously stored global variables:

```http
> {%
    const token = client.global.get("authToken");
    client.log("Using token:", token);
%}
```

#### `client.global.clear(name)`

Remove a global variable:

```http
> {%
    client.global.clear("authToken");
%}
```

### Response Object

Access response data in scripts:

```javascript
// Response status
response.status          // 200
response.statusText      // "OK"

// Response headers
response.headers["content-type"]  // "application/json"

// Response body (automatically parsed if JSON)
response.body.id         // Access JSON properties
response.body.name

// Content type
response.contentType     // "application/json; charset=utf-8"
```

### Request Object

Access request data in scripts:

```javascript
// Request method and URL
request.method           // "GET"
request.url              // "https://api.example.com/users/1"

// Request headers
request.headers["authorization"]
```

### Environment Object

Access environment variables in scripts:

```javascript
// Access environment variables
env.baseUrl
env.apiKey
env.timeout
```

## Global Variables

Global variables persist across all requests in a session, making it easy to:
- Store authentication tokens
- Pass data between requests
- Build request chains
- Implement workflows

### Complete Example

```http
### Step 1: Create a post
# @name createPost
POST {{baseUrl}}/posts
Content-Type: application/json

{
  "title": "Test Post",
  "body": "This is a test post",
  "userId": 1
}

> {%
    client.test("Post created successfully", function() {
        client.assert(response.status === 201, "Expected 201 Created");
    });
    
    // Save the post ID for next request
    client.global.set("postId", response.body.id);
    client.log("Created post with ID:", response.body.id);
%}

###

### Step 2: Get the created post
GET {{baseUrl}}/posts/{{postId}}

> {%
    client.test("Retrieved post successfully", function() {
        client.assert(response.status === 200, "Expected 200 OK");
        client.assert(response.body.id === client.global.get("postId"));
    });
    
    client.log("Post title:", response.body.title);
%}

###

### Step 3: Update the post
PUT {{baseUrl}}/posts/{{postId}}
Content-Type: application/json

{
  "id": {{postId}},
  "title": "Updated Title",
  "body": "Updated content",
  "userId": 1
}

> {%
    client.test("Post updated successfully", function() {
        client.assert(response.status === 200, "Expected 200 OK");
    });
%}

###

### Step 4: Delete the post
DELETE {{baseUrl}}/posts/{{postId}}

> {%
    client.test("Post deleted successfully", function() {
        client.assert(response.status === 200, "Expected 200 OK");
    });
    
    // Clean up global variable
    client.global.clear("postId");
%}
```

## Command Reference

### Context Management

Set default HTTP file and environment for a directory to avoid specifying them every time:

```bash
# Set context with HTTP file and environment
postie context set --http-file requests.http --env development

# Set with save responses enabled
postie context set --http-file api.http --env production --save-responses

# Show current context
postie context show

# Clear context
postie context clear
```

Once context is set, you can run requests without specifying the file:

```bash
# Runs using context defaults
postie http run --request getUserById

# Override environment from context
postie http run --env staging
```

Context is saved to `.postie-context.json` in the current directory (should be added to `.gitignore`).

### Environment Management

Inspect and manage environment files:

```bash
# List all available environments
postie env list

# Show variables for a specific environment
postie env show development

# Show including private variables
postie env show development --show-private

# Use custom environment file path
postie env list --env-file custom-env.json
```

### Run Requests

Execute requests from a `.http` file:

```bash
# Run all requests
postie http run requests.http

# Run with specific environment
postie http run requests.http --env production

# Run with verbose output
postie http run requests.http --verbose

# Run specific request by name or number
postie http run requests.http --request createUser
postie http run requests.http --request 1

# Save responses to files
postie http run requests.http --save-responses

# Using context (no file needed if context is set)
postie http run --request getUserById
```

### Parse Requests

Parse and validate `.http` files without executing:

```bash
# Parse and show summary
postie http parse requests.http

# Parse with JSON output
postie http parse requests.http --format json

# Parse with validation
postie http parse requests.http --validate
```

### List Requests

List all requests in a file:

```bash
# List requests in current directory
postie http list

# List requests in specific directory
postie http list api-tests/

# List recursively
postie http list --recursive
```

## Examples

### Basic CRUD Operations

```http
### Create User
POST https://jsonplaceholder.typicode.com/users
Content-Type: application/json

{
  "name": "John Doe",
  "username": "johndoe",
  "email": "john@example.com"
}

###

### Get All Users
GET https://jsonplaceholder.typicode.com/users
Accept: application/json

###

### Get Single User
GET https://jsonplaceholder.typicode.com/users/1

###

### Update User
PUT https://jsonplaceholder.typicode.com/users/1
Content-Type: application/json

{
  "id": 1,
  "name": "Jane Doe",
  "username": "janedoe",
  "email": "jane@example.com"
}

###

### Delete User
DELETE https://jsonplaceholder.typicode.com/users/1
```

### Authentication Flow

```http
### Login
# @name login
POST {{baseUrl}}/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}

> {%
    client.test("Login successful", function() {
        client.assert(response.status === 200, "Login should succeed");
        client.assert(response.body.token, "Should return auth token");
    });
    
    client.global.set("authToken", response.body.token);
    client.log("Logged in successfully");
%}

###

### Get Protected Resource
GET {{baseUrl}}/api/protected
Authorization: Bearer {{authToken}}

> {%
    client.test("Access granted", function() {
        client.assert(response.status === 200, "Should access protected route");
    });
%}

###

### Logout
POST {{baseUrl}}/auth/logout
Authorization: Bearer {{authToken}}

> {%
    client.global.clear("authToken");
    client.log("Logged out");
%}
```

### API Testing with Validation

```http
### Test API Endpoint
GET https://jsonplaceholder.typicode.com/posts/1

> {%
    client.test("Status code is 200", function() {
        client.assert(response.status === 200, "Expected 200 OK");
    });
    
    client.test("Response is JSON", function() {
        client.assert(response.contentType.includes("application/json"));
    });
    
    client.test("Response has required fields", function() {
        client.assert(response.body.id, "Should have id field");
        client.assert(response.body.title, "Should have title field");
        client.assert(response.body.body, "Should have body field");
        client.assert(response.body.userId, "Should have userId field");
    });
    
    client.test("Data types are correct", function() {
        client.assert(typeof response.body.id === "number");
        client.assert(typeof response.body.title === "string");
        client.assert(typeof response.body.userId === "number");
    });
    
    client.log("All validations passed");
%}
```

### Chained Requests with Dependencies

```http
### Create Post
POST https://jsonplaceholder.typicode.com/posts
Content-Type: application/json

{
  "title": "My Post",
  "body": "Post content",
  "userId": 1
}

> {%
    client.global.set("createdPostId", response.body.id);
    client.log("Created post:", response.body.id);
%}

###

### Add Comment to Post
POST https://jsonplaceholder.typicode.com/comments
Content-Type: application/json

{
  "postId": {{createdPostId}},
  "name": "Great post!",
  "email": "commenter@example.com",
  "body": "This is a great post!"
}

> {%
    client.global.set("commentId", response.body.id);
    client.log("Added comment:", response.body.id);
%}

###

### Get Post with Comments
GET https://jsonplaceholder.typicode.com/posts/{{createdPostId}}/comments

> {%
    client.test("Comments retrieved", function() {
        client.assert(response.status === 200);
        client.assert(Array.isArray(response.body));
    });
    
    client.log("Found", response.body.length, "comments");
%}
```

## Best Practices

### 1. Organize Requests

Group related requests in separate files:

```
api-tests/
├── auth.http           # Authentication requests
├── users.http          # User management
├── posts.http          # Post operations
└── admin.http          # Admin operations
```

### 2. Use Environment Files

Keep sensitive data in `http-client.private.env.json` and add it to `.gitignore`:

```bash
# .gitignore
http-client.private.env.json
```

### 3. Name Your Requests

Use `# @name` to identify requests:

```http
# @name getUserProfile
GET {{baseUrl}}/users/me
```

### 4. Add Descriptive Comments

Document what each request does:

```http
### Get User Profile
# This endpoint retrieves the authenticated user's profile
# Requires: Authentication token
# @name getUserProfile
GET {{baseUrl}}/users/me
Authorization: Bearer {{authToken}}
```

### 5. Test Critical Paths

Use response handlers to validate important workflows:

```http
> {%
    client.test("User creation successful", function() {
        client.assert(response.status === 201);
        client.assert(response.body.id);
        client.assert(response.body.email === "test@example.com");
    });
%}
```

### 6. Clean Up Global Variables

Clear global variables when no longer needed:

```http
> {%
    client.global.clear("tempToken");
    client.global.clear("sessionId");
%}
```

## Troubleshooting

### Common Issues

**Variable not expanding:**
- Check that the variable is defined in your environment file
- Verify you're using the correct environment (`--env`)
- Make sure the variable name matches exactly (case-sensitive)

**Script errors:**
- Check JavaScript syntax in response handlers
- Verify response object structure before accessing properties
- Use `client.log()` to debug values

**Request fails:**
- Use `--verbose` flag to see detailed request information
- Check headers and body formatting
- Verify the URL is correct

### Debug Mode

Run with verbose output to see detailed information:

```bash
postie http run requests.http --verbose
```

This shows:
- Request method and URL
- All headers
- Request body
- Response status and timing
- Response headers and body
- Script execution results
- Global variables set

## Conclusion

Postie provides a powerful and flexible way to test APIs using the industry-standard HTTP Request in Editor format. With support for environment variables, response handler scripts, and global variable sharing, you can build complex testing workflows while keeping your requests organized and version-controlled.

For more examples, see the `api-tests/` directory in the repository.
