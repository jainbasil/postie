# Postie Samples

This directory contains example files demonstrating various features of Postie.

## Files

### Environment Files

- **`http-client.env.json`** - Example public environment file with non-sensitive configuration
- **`http-client.private.env.json`** - Example private environment file with sensitive data (should be in `.gitignore`)

### HTTP Request Files

- **`basic-requests.http`** - Basic CRUD operations (GET, POST, PUT, PATCH, DELETE)
- **`authentication.http`** - Authentication flows (login, logout, registration, OAuth, password reset)
- **`advanced-features.http`** - Advanced features (variables, chaining, testing, error handling)

## Usage

### Quick Start

1. **Copy environment files to your project:**
   ```bash
   cp samples/http-client.env.json ./
   cp samples/http-client.private.env.json ./
   ```

2. **Copy an HTTP request file:**
   ```bash
   cp samples/basic-requests.http ./requests.http
   ```

3. **Run the requests:**
   ```bash
   # Run all requests
   postie http run requests.http

   # Run with specific environment
   postie http run requests.http --env production

   # Run specific request
   postie http run requests.http --request 1
   ```

### Using with Context

Set up context to avoid repeating file and environment names:

```bash
# Set context
postie context set --http-file requests.http --env development

# Now run without specifying file
postie http run --request "Get all posts"

# View current context
postie context show
```

### Environment Management

Inspect environment variables:

```bash
# List all environments
postie env list

# Show variables for development
postie env show development

# Show including private variables
postie env show development --show-private
```

## File Descriptions

### basic-requests.http

Demonstrates basic HTTP operations:
- Simple GET requests
- POST requests with JSON bodies
- PUT and PATCH for updates
- DELETE operations
- Response handler scripts
- Global variable usage
- Query parameters
- Custom headers

**Example requests:**
- Get all posts
- Create new post
- Update post
- Delete post
- Get users and filter by criteria

### authentication.http

Complete authentication workflows:
- User login with token storage
- Token refresh mechanism
- Profile management (get, update)
- Password management (change, reset)
- User registration
- Email verification
- OAuth flows (Google example)
- Logout and token cleanup

**Key features:**
- Token management using `client.global.set()`
- Chaining authenticated requests
- Error handling for auth failures
- Multi-step authentication flows

### advanced-features.http

Advanced Postie capabilities:
- Inline variable definitions (`@variable = value`)
- Dynamic values (`{{$uuid}}`, `{{$timestamp}}`)
- Request chaining with global variables
- Complex response validation
- Nested object testing
- Array response handling
- Performance testing
- Error handling patterns
- Conditional logic in tests

**Testing features:**
- `client.test()` - Create test assertions
- `client.assert()` - Make assertions
- `client.log()` - Log information
- `client.global.set()` / `.get()` - Share data between requests

## Environment File Structure

### Public Environment (`http-client.env.json`)

Contains non-sensitive configuration that can be committed to version control:

```json
{
  "development": {
    "baseUrl": "https://api-dev.example.com",
    "apiVersion": "v1",
    "timeout": 30000,
    "enableDebug": true
  },
  "production": {
    "baseUrl": "https://api.example.com",
    "apiVersion": "v2",
    "timeout": 10000,
    "enableDebug": false
  }
}
```

### Private Environment (`http-client.private.env.json`)

Contains sensitive data that should **never** be committed (add to `.gitignore`):

```json
{
  "development": {
    "apiKey": "dev-secret-key",
    "authToken": "dev-bearer-token",
    "dbPassword": "dev-password"
  },
  "production": {
    "apiKey": "prod-secret-key",
    "authToken": "prod-bearer-token",
    "dbPassword": "prod-password"
  }
}
```

## Variable Usage

### Environment Variables

Use `{{variableName}}` to reference environment variables:

```http
GET {{baseUrl}}/{{apiVersion}}/users
Authorization: Bearer {{authToken}}
X-API-Key: {{apiKey}}
```

### Inline Variables

Define variables directly in the .http file:

```http
@hostname = api.example.com
@protocol = https

GET {{protocol}}://{{hostname}}/users
```

### Dynamic Variables

Built-in dynamic variables:
- `{{$uuid}}` - Generate a random UUID
- `{{$timestamp}}` - Current Unix timestamp
- `{{$randomInt}}` - Random integer

### Global Variables

Share data between requests using response handlers:

```http
### Login
POST {{baseUrl}}/auth/login
Content-Type: application/json

{"email": "user@example.com", "password": "pass123"}

> {%
  // Save token for subsequent requests
  client.global.set("authToken", response.body.token);
  client.global.set("userId", response.body.user.id);
%}

### Use saved data
GET {{baseUrl}}/users/{{userId}}
Authorization: Bearer {{authToken}}
```

## Response Handler Scripts

Add JavaScript code after requests to test responses and extract data:

```http
GET {{baseUrl}}/api/data

> {%
  // Test response
  client.test("Request successful", function() {
    client.assert(response.status === 200);
    client.assert(response.body.hasOwnProperty("data"));
  });
  
  // Extract and save data
  client.global.set("dataId", response.body.data.id);
  
  // Log information
  client.log("Data ID: " + response.body.data.id);
%}
```

### Available in Scripts

**`response` object:**
- `response.status` - HTTP status code
- `response.statusText` - Status text
- `response.headers` - Response headers
- `response.body` - Parsed response body (JSON)
- `response.contentType` - Content type info
- `response.responseTime` - Response time in ms

**`client` object:**
- `client.test(name, fn)` - Create a test
- `client.assert(condition, message)` - Assert a condition
- `client.log(message)` - Log a message
- `client.global.set(key, value)` - Save global variable
- `client.global.get(key)` - Retrieve global variable
- `client.global.clear(key)` - Clear global variable

## Tips and Best Practices

### 1. Environment Organization

- Keep public and private environments in sync (same keys)
- Use descriptive variable names
- Document what each variable is for
- Never commit `http-client.private.env.json`

### 2. Request Organization

- Group related requests in the same file
- Use descriptive request names
- Add comments to explain complex requests
- Use `@name` directive for important requests

### 3. Variable Management

- Use environment variables for configuration
- Use global variables for runtime data
- Clear sensitive data after use
- Validate required variables exist

### 4. Testing Best Practices

- Test status codes first
- Validate response structure
- Test edge cases
- Log important information
- Use meaningful test names

### 5. Security

- Never commit private environment files
- Use strong dummy values in examples
- Clear auth tokens when done
- Be careful with logging sensitive data

## .gitignore Recommendations

Add these to your `.gitignore`:

```gitignore
# Private environment variables with sensitive data
http-client.private.env.json

# HTTP response storage (if saving responses)
.http-responses/

# Context files (local development preferences)
.postie-context.json

# Any environment-specific configuration
*.local.env.json
*.secret.json
```

## Example Workflow

Here's a complete workflow using these samples:

```bash
# 1. Copy samples to your project
cp samples/*.json ./
cp samples/basic-requests.http ./requests.http

# 2. Edit environment files with your API details
vim http-client.env.json          # Update baseUrl, etc.
vim http-client.private.env.json  # Add your API keys

# 3. Set up context
postie context set --http-file requests.http --env development

# 4. Test environment setup
postie env show development --show-private

# 5. Run requests
postie http run                    # All requests
postie http run --request 1        # First request
postie http run --verbose          # With detailed output

# 6. Save responses for debugging
postie http run --save-responses

# 7. Switch to production
postie context set --env production
postie http run --request "Get all posts"
```

## Customization

Feel free to customize these samples for your specific API:

1. Update `baseUrl` in environment files
2. Modify request endpoints
3. Adjust request bodies for your data model
4. Update response handlers for your API responses
5. Add new authentication methods if needed
6. Create environment-specific test scenarios

## Additional Resources

- [User Guide](../docs/user-guide.md) - Comprehensive usage guide
- [Command Reference](../docs/command-reference.md) - Detailed command documentation
- [README](../README.md) - Project overview and quick start

## Support

For issues or questions:
- GitHub Issues: https://github.com/jainbasil/postie/issues
- Documentation: See `docs/` directory
