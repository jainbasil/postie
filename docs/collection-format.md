# API Collection JSON Format Specification

This document defines the JSON file format for storing API collections in the Postie tool. Collections allow you to group related API requests together with shared configuration and reusable properties.

## Collection Format Overview

```json
{
  "collection": {
    "info": {
      "name": "My API Collection",
      "description": "Collection description",
      "version": "1.0.0",
      "schema": "https://postie.dev/collection/v1.0.0/collection.json"
    },
    "variable": [
      {
        "key": "baseUrl",
        "value": "https://api.example.com",
        "type": "string"
      }
    ],
    "auth": {
      "type": "bearer",
      "bearer": [
        {
          "key": "token",
          "value": "{{apiToken}}",
          "type": "string"
        }
      ]
    },
    "event": [
      {
        "listen": "prerequest",
        "script": {
          "type": "text/javascript",
          "exec": [
            "console.log('Pre-request script');"
          ]
        }
      }
    ],
    "apiGroup": [
      {
        "id": "users",
        "name": "Get Users",
        "request": {
          "method": "GET",
          "header": [],
          "url": {
            "raw": "{{baseUrl}}/users",
            "host": ["{{baseUrl}}"],
            "path": ["users"],
            "query": []
          }
        },
        "response": []
      }
    ]
  }
}
```

## Detailed Specification

### 1. Root Object

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `collection` | Object | Yes | Root collection object |

### 2. Collection Info

The `info` object contains metadata about the collection.

```json
{
  "info": {
    "name": "API Collection Name",
    "description": "Detailed description of the collection",
    "version": "1.0.0",
    "schema": "https://postie.dev/collection/v1.0.0/collection.json",
    "author": "Developer Name",
    "license": "MIT",
    "documentation": "https://docs.example.com"
  }
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | String | Yes | Collection name |
| `description` | String | No | Collection description |
| `version` | String | No | Collection version (semver) |
| `schema` | String | No | Schema URL for validation |
| `author` | String | No | Collection author |
| `license` | String | No | License type |
| `documentation` | String | No | Documentation URL |

### 3. Variables

Variables allow you to define reusable values across the collection.

```json
{
  "variable": [
    {
      "key": "baseUrl",
      "value": "https://api.example.com",
      "type": "string",
      "description": "Base API URL"
    },
    {
      "key": "apiVersion",
      "value": "v1",
      "type": "string",
      "description": "API version"
    },
    {
      "key": "timeout",
      "value": 30000,
      "type": "number",
      "description": "Request timeout in milliseconds"
    }
  ]
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | String | Yes | Variable name |
| `value` | Any | Yes | Variable value |
| `type` | String | No | Variable type (string, number, boolean) |
| `description` | String | No | Variable description |

### 4. Environments

Environments allow you to define multiple sets of variable values for different deployment targets (development, staging, production, etc.). Each environment can override collection-level variables.

```json
{
  "environment": [
    {
      "name": "Development",
      "description": "Development environment configuration",
      "values": [
        {
          "key": "baseUrl",
          "value": "https://api-dev.example.com",
          "type": "string",
          "enabled": true
        },
        {
          "key": "apiKey",
          "value": "dev-api-key-12345",
          "type": "string",
          "enabled": true
        },
        {
          "key": "timeout",
          "value": 60000,
          "type": "number",
          "enabled": true
        },
        {
          "key": "debugMode",
          "value": true,
          "type": "boolean",
          "enabled": true
        }
      ]
    },
    {
      "name": "Production",
      "description": "Production environment configuration",
      "values": [
        {
          "key": "baseUrl",
          "value": "https://api.example.com",
          "type": "string",
          "enabled": true
        },
        {
          "key": "apiKey",
          "value": "{{PROD_API_KEY}}",
          "type": "string",
          "enabled": true
        },
        {
          "key": "timeout",
          "value": 30000,
          "type": "number",
          "enabled": true
        },
        {
          "key": "debugMode",
          "value": false,
          "type": "boolean",
          "enabled": true
        }
      ]
    },
    {
      "name": "Testing",
      "description": "Testing environment with mock data",
      "values": [
        {
          "key": "baseUrl",
          "value": "https://api-test.example.com",
          "type": "string",
          "enabled": true
        },
        {
          "key": "apiKey",
          "value": "test-api-key",
          "type": "string",
          "enabled": true
        },
        {
          "key": "useMockData",
          "value": true,
          "type": "boolean",
          "enabled": true
        }
      ]
    }
  ]
}
```

#### Environment Object Structure

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | String | Yes | Environment name |
| `description` | String | No | Environment description |
| `values` | Array | Yes | Array of environment variables |

#### Environment Variable Structure

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | String | Yes | Variable name |
| `value` | Any | Yes | Variable value |
| `type` | String | No | Variable type (string, number, boolean) |
| `enabled` | Boolean | No | Whether the variable is active (default: true) |
| `description` | String | No | Variable description |

### 5. Authentication

Collection-level authentication that applies to all requests unless overridden. Authentication can also be environment-specific.

#### Collection-Level Authentication

```json
{
  "auth": {
    "type": "bearer",
    "bearer": [
      {
        "key": "token",
        "value": "{{apiToken}}",
        "type": "string"
      }
    ]
  }
}
```

#### Environment-Specific Authentication

Environments can override collection-level authentication:

```json
{
  "environment": [
    {
      "name": "Development",
      "auth": {
        "type": "basic",
        "basic": [
          {
            "key": "username",
            "value": "dev-user",
            "type": "string"
          },
          {
            "key": "password",
            "value": "dev-password",
            "type": "string"
          }
        ]
      },
      "values": [...]
    },
    {
      "name": "Production",
      "auth": {
        "type": "bearer",
        "bearer": [
          {
            "key": "token",
            "value": "{{PROD_TOKEN}}",
            "type": "string"
          }
        ]
      },
      "values": [...]
    }
  ]
}
```

#### Bearer Token Authentication

```json
{
  "auth": {
    "type": "bearer",
    "bearer": [
      {
        "key": "token",
        "value": "{{apiToken}}",
        "type": "string"
      }
    ]
  }
}
```

#### API Key Authentication

```json
{
  "auth": {
    "type": "apikey",
    "apikey": [
      {
        "key": "key",
        "value": "X-API-Key",
        "type": "string"
      },
      {
        "key": "value",
        "value": "{{apiKey}}",
        "type": "string"
      },
      {
        "key": "in",
        "value": "header",
        "type": "string"
      }
    ]
  }
}
```

#### Basic Authentication

```json
{
  "auth": {
    "type": "basic",
    "basic": [
      {
        "key": "username",
        "value": "{{username}}",
        "type": "string"
      },
      {
        "key": "password",
        "value": "{{password}}",
        "type": "string"
      }
    ]
  }
}
```

#### OAuth 2.0 Authentication

```json
{
  "auth": {
    "type": "oauth2",
    "oauth2": [
      {
        "key": "accessToken",
        "value": "{{accessToken}}",
        "type": "string"
      },
      {
        "key": "tokenType",
        "value": "Bearer",
        "type": "string"
      }
    ]
  }
}
```

### 5. Events (Scripts)

Collection-level scripts that run before or after requests.

```json
### 6. Events (Scripts)

Collection-level scripts that run before or after requests. Scripts can also be environment-specific.

#### Collection-Level Events

```json
{
  "event": [
    {
      "listen": "prerequest",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// Set timestamp variable",
          "pm.collectionVariables.set('timestamp', Date.now());"
        ]
      }
    },
    {
      "listen": "test",
      "script": {
        "type": "text/javascript",
        "exec": [
          "// Common test for all requests",
          "pm.test('Status code is not 5xx', function () {",
          "    pm.expect(pm.response.code).to.be.below(500);",
          "});"
        ]
      }
    }
  ]
}
```

#### Environment-Specific Events

Environments can have their own scripts that run in addition to or instead of collection-level scripts:

```json
{
  "environment": [
    {
      "name": "Development",
      "event": [
        {
          "listen": "prerequest",
          "script": {
            "type": "text/javascript",
            "exec": [
              "// Development-specific setup",
              "console.log('Running in development mode');",
              "pm.environment.set('debugLevel', 'verbose');"
            ]
          }
        },
        {
          "listen": "test",
          "script": {
            "type": "text/javascript",
            "exec": [
              "// More lenient tests for development",
              "pm.test('Response time is acceptable for dev', function () {",
              "    pm.expect(pm.response.responseTime).to.be.below(10000);",
              "});"
            ]
          }
        }
      ],
      "values": [...]
    },
    {
      "name": "Production",
      "event": [
        {
          "listen": "test",
          "script": {
            "type": "text/javascript",
            "exec": [
              "// Strict tests for production",
              "pm.test('Response time is fast in production', function () {",
              "    pm.expect(pm.response.responseTime).to.be.below(2000);",
              "});",
              "",
              "pm.test('No debug information in headers', function () {",
              "    pm.expect(pm.response.headers.get('X-Debug')).to.be.undefined;",
              "});"
            ]
          }
        }
      ],
      "values": [...]
    }
  ]
}
```

### 7. API Groups
```

### 6. API Groups

The main content of the collection - individual API requests and groups.

```json
{
  "apiGroup": [
    {
      "name": "Get All Users",
      "description": "Retrieve all users from the system",
      "request": {
        "method": "GET",
        "header": [
          {
            "key": "Accept",
            "value": "application/json",
            "type": "text"
          }
        ],
        "url": {
          "raw": "{{baseUrl}}/api/{{apiVersion}}/users?page=1&limit=10",
          "protocol": "https",
          "host": ["api", "example", "com"],
          "path": ["api", "{{apiVersion}}", "users"],
          "query": [
            {
              "key": "page",
              "value": "1",
              "description": "Page number"
            },
            {
              "key": "limit",
              "value": "10",
              "description": "Items per page"
            }
          ]
        },
        "auth": {
          "type": "inherit"
        }
      },
      "response": [
        {
          "name": "Successful Response",
          "originalRequest": {
            "method": "GET",
            "header": [],
            "url": {
              "raw": "{{baseUrl}}/api/v1/users"
            }
          },
          "status": "OK",
          "code": 200,
          "header": [
            {
              "key": "Content-Type",
              "value": "application/json"
            }
          ],
          "body": "{\n  \"users\": [\n    {\n      \"id\": 1,\n      \"name\": \"John Doe\"\n    }\n  ]\n}",
          "responseTime": 245
        }
      ],
      "event": [
        {
          "listen": "test",
          "script": {
            "type": "text/javascript",
            "exec": [
              "pm.test('Response contains users array', function () {",
              "    const jsonData = pm.response.json();",
              "    pm.expect(jsonData).to.have.property('users');",
              "});"
            ]
          }
        }
      ]
    }
  ]
}
```

### 8. Folder Structure

Collections can contain folders to organize requests hierarchically. Folders can also have environment-specific configurations.

```json
{
  "apiGroup": [
    {
      "name": "User Management",
      "description": "All user-related endpoints",
      "environment": [
        {
          "name": "Development",
          "values": [
            {
              "key": "userApiEndpoint",
              "value": "/dev/users",
              "type": "string",
              "enabled": true
            }
          ]
        },
        {
          "name": "Production", 
          "values": [
            {
              "key": "userApiEndpoint",
              "value": "/api/v2/users",
              "type": "string",
              "enabled": true
            }
          ]
        }
      ],
      "apis": [
        {
          "name": "Get Users",
          "request": {
            "method": "GET",
            "url": "{{baseUrl}}{{userApiEndpoint}}"
          }
        },
        {
          "name": "Create User",
          "request": {
            "method": "POST",
            "url": "{{baseUrl}}{{userApiEndpoint}}",
            "body": {
              "mode": "raw",
              "raw": "{\n  \"name\": \"John Doe\",\n  \"email\": \"john@example.com\"\n}",
              "options": {
                "raw": {
                  "language": "json"
                }
              }
            }
          }
        }
      ],
      "auth": {
        "type": "inherit"
      },
      "event": [
        {
          "listen": "prerequest",
          "script": {
            "type": "text/javascript",
            "exec": [
              "console.log('User management folder pre-request');"
            ]
          }
        }
      ]
    }
  ]
}
```

## Complete Example

```json
{
  "collection": {
    "info": {
      "name": "GitHub API Collection",
      "description": "Collection for testing GitHub REST API endpoints",
      "version": "1.0.0",
      "schema": "https://postie.dev/collection/v1.0.0/collection.json",
      "author": "API Team",
      "documentation": "https://docs.github.com/en/rest"
    },
    "variable": [
      {
        "key": "username",
        "value": "octocat",
        "type": "string",
        "description": "Default GitHub username"
      },
      {
        "key": "timeout",
        "value": 30000,
        "type": "number",
        "description": "Default request timeout"
      }
    ],
    "environment": [
      {
        "name": "GitHub Public API",
        "description": "Public GitHub API configuration",
        "values": [
          {
            "key": "baseUrl",
            "value": "https://api.github.com",
            "type": "string",
            "enabled": true
          },
          {
            "key": "token",
            "value": "",
            "type": "string",
            "enabled": false,
            "description": "Optional personal access token"
          },
          {
            "key": "rateLimit",
            "value": 60,
            "type": "number",
            "enabled": true,
            "description": "Requests per hour for unauthenticated users"
          }
        ],
        "auth": {
          "type": "noauth"
        },
        "event": [
          {
            "listen": "prerequest",
            "script": {
              "type": "text/javascript",
              "exec": [
                "// Public API - no authentication needed",
                "console.log('Using public GitHub API');"
              ]
            }
          }
        ]
      },
      {
        "name": "GitHub Enterprise",
        "description": "GitHub Enterprise Server configuration",
        "values": [
          {
            "key": "baseUrl",
            "value": "https://github.company.com/api/v3",
            "type": "string",
            "enabled": true
          },
          {
            "key": "token",
            "value": "{{GITHUB_ENTERPRISE_TOKEN}}",
            "type": "string",
            "enabled": true,
            "description": "Enterprise access token"
          },
          {
            "key": "rateLimit",
            "value": 5000,
            "type": "number",
            "enabled": true,
            "description": "Higher rate limit for enterprise"
          }
        ],
        "auth": {
          "type": "bearer",
          "bearer": [
            {
              "key": "token",
              "value": "{{token}}",
              "type": "string"
            }
          ]
        },
        "event": [
          {
            "listen": "prerequest",
            "script": {
              "type": "text/javascript",
              "exec": [
                "// Enterprise API - add custom headers",
                "pm.request.headers.add({",
                "    key: 'X-Enterprise-Client',",
                "    value: 'Postie/1.0.0'",
                "});"
              ]
            }
          },
          {
            "listen": "test",
            "script": {
              "type": "text/javascript",
              "exec": [
                "// Enterprise-specific tests",
                "pm.test('Enterprise rate limit header present', function () {",
                "    pm.expect(pm.response.headers.get('X-RateLimit-Limit')).to.not.be.undefined;",
                "});"
              ]
            }
          }
        ]
      }
    ],
    "auth": {
      "type": "bearer",
      "bearer": [
        {
          "key": "token",
          "value": "{{token}}",
          "type": "string"
        }
      ]
    },
    "event": [
      {
        "listen": "prerequest",
        "script": {
          "type": "text/javascript",
          "exec": [
            "// Set User-Agent header for all requests",
            "pm.request.headers.add({",
            "    key: 'User-Agent',",
            "    value: 'Postie/1.0.0'",
            "});"
          ]
        }
      }
    ],
    "apiGroup": [
      {
        "name": "User Operations",
        "description": "GitHub user-related operations",
        "environment": [
          {
            "name": "GitHub Public API",
            "values": [
              {
                "key": "userEndpoint",
                "value": "/users",
                "type": "string",
                "enabled": true
              }
            ]
          },
          {
            "name": "GitHub Enterprise",
            "values": [
              {
                "key": "userEndpoint", 
                "value": "/enterprise/users",
                "type": "string",
                "enabled": true
              }
            ]
          }
        ],
        "apis": [
          {
            "name": "Get User",
            "description": "Get a user by username",
            "request": {
              "method": "GET",
              "header": [
                {
                  "key": "Accept",
                  "value": "application/vnd.github.v3+json",
                  "type": "text"
                }
              ],
              "url": {
                "raw": "{{baseUrl}}{{userEndpoint}}/{{username}}",
                "host": ["{{baseUrl}}"],
                "path": ["{{userEndpoint}}", "{{username}}"]
              }
            },
            "response": [
              {
                "name": "User Found",
                "originalRequest": {
                  "method": "GET",
                  "header": [],
                  "url": {
                    "raw": "{{baseUrl}}{{userEndpoint}}/{{username}}"
                  }
                },
                "status": "OK",
                "code": 200,
                "header": [
                  {
                    "key": "Content-Type",
                    "value": "application/json; charset=utf-8"
                  }
                ],
                "body": "{\n  \"login\": \"octocat\",\n  \"id\": 1,\n  \"name\": \"The Octocat\"\n}",
                "responseTime": 312
              }
            ],
            "event": [
              {
                "listen": "test",
                "script": {
                  "type": "text/javascript",
                  "exec": [
                    "pm.test('Status code is 200', function () {",
                    "    pm.response.to.have.status(200);",
                    "});",
                    "",
                    "pm.test('Response has login field', function () {",
                    "    const jsonData = pm.response.json();",
                    "    pm.expect(jsonData).to.have.property('login');",
                    "});"
                  ]
                }
              }
            ]
          },
          {
            "name": "List User Repositories",
            "description": "List repositories for a user",
            "request": {
              "method": "GET",
              "header": [],
              "url": {
                "raw": "{{baseUrl}}{{userEndpoint}}/{{username}}/repos?type=owner&sort=updated",
                "host": ["{{baseUrl}}"],
                "path": ["{{userEndpoint}}", "{{username}}", "repos"],
                "query": [
                  {
                    "key": "type",
                    "value": "owner",
                    "description": "Repository type filter"
                  },
                  {
                    "key": "sort",
                    "value": "updated",
                    "description": "Sort order"
                  }
                ]
              }
            }
          }
        ]
      }
    ]
  }
}
```

## Environment Files

Separate environment files can be used to manage different configurations:

### development.json
```json
{
  "name": "Development",
  "values": [
    {
      "key": "baseUrl",
      "value": "https://api-dev.example.com",
      "enabled": true
    },
    {
      "key": "token",
      "value": "dev-token-12345",
      "enabled": true
    }
  ]
}
```

### production.json
```json
{
  "name": "Production",
  "values": [
    {
      "key": "baseUrl",
      "value": "https://api.example.com",
      "enabled": true
    },
    {
      "key": "token",
      "value": "prod-token-67890",
      "enabled": true
    }
  ]
}
```

## Implementation Notes

1. **Variable Resolution**: Variables are resolved in the following order (highest to lowest priority):
   - Environment-specific variables (selected environment)
   - Collection-level variables
   - Global variables (external environment files)
   - System environment variables

2. **Environment Selection**: 
   - Collections can contain multiple environments internally
   - One environment is active at a time
   - Environment selection can be done via CLI: `postie run collection.json --env "Production"`
   - Default environment is the first one defined in the collection

3. **Authentication Inheritance**: 
   - Environment-specific auth overrides collection-level auth
   - Request-level auth overrides environment and collection auth
   - Folder-level auth overrides collection but not environment auth

4. **Script Execution**: Scripts execute in this order:
   - Collection pre-request scripts
   - Environment pre-request scripts (if active environment has them)
   - Folder pre-request scripts
   - Request pre-request scripts
   - Request execution
   - Request test scripts
   - Folder test scripts
   - Environment test scripts (if active environment has them)
   - Collection test scripts

5. **Environment-Specific Configuration Merging**:
   - Environment variables override collection variables with the same key
   - Environment auth completely replaces collection auth (no merging)
   - Environment scripts run in addition to collection scripts
   - Folder environments override collection environments for that folder's scope

6. **File Extensions**: 
   - Collections: `.json`
   - External environments (optional): `.json`
   - Suggested naming: `collection-name.postman_collection.json`

7. **Validation**: Collections should be validated against the JSON schema for consistency

8. **Backward Compatibility**: 
   - Collections without internal environments work as before
   - External environment files can still be used alongside internal environments
   - External environments take precedence over internal ones when both are present

This format provides a comprehensive way to organize and share API collections while maintaining compatibility with existing tools and workflows.