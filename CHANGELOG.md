# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Added
- Initial project setup and core functionality

### Changed
- Nothing yet

### Deprecated
- Nothing yet

### Removed
- Nothing yet

### Fixed
- Nothing yet

### Security
- Nothing yet

## [1.0.0] - 2025-10-31

### Added
- 🚀 **Core HTTP Client**: Full-featured HTTP client with fluent API
- 🔐 **Authentication Support**: API keys, Bearer tokens, Basic auth, custom headers
- 🌐 **HTTP Methods**: GET, POST, PUT, DELETE, PATCH, HEAD, OPTIONS support
- 📊 **Response Handling**: JSON formatting, status checking, response timing
- 🔄 **Middleware System**: Extensible middleware for logging, rate limiting, retries
- 📁 **Collection Format**: JSON-based API collection format with environment support
- 📝 **Variable Resolution**: Dynamic variable replacement in requests
- ⚡ **CLI Interface**: Comprehensive command-line interface for all operations
- 🛠️ **Collection Management**: Create, read, update, delete operations for collections
- 📋 **CRUD Operations**: Full CRUD support for collections, API groups, and APIs
- 🏷️ **Slug-based IDs**: Human-readable identifiers for collection items
- 🎯 **Demo Mode**: Interactive demonstrations of all features
- 📖 **Comprehensive Documentation**: Full documentation with examples
- 🧪 **Test Suite**: Complete test coverage for core functionality
- 🔨 **Build System**: VS Code tasks and Go build configuration
- 📦 **Single Binary**: Self-contained executable with no dependencies

### Technical Details
- **Language**: Go 1.21+
- **Architecture**: Modular package-based design
- **Dependencies**: Minimal external dependencies
- **Performance**: Native Go performance with efficient HTTP handling
- **Compatibility**: Cross-platform support (Linux, macOS, Windows)

### CLI Commands
```bash
# HTTP Methods
postie get <url>           # GET request
postie post <url>          # POST request  
postie put <url>           # PUT request
postie delete <url>        # DELETE request

# Collection Operations
postie run <collection>    # Run collection
postie list <collection>   # List requests
postie env <collection>    # Show environments

# Collection Management
postie create collection <name>              # Create collection
postie create apigroup <file> <name>         # Create API group
postie create api <file> <group> <name>      # Create API
postie update collection <file>              # Update collection
postie remove apigroup <file> <id>           # Remove API group

# Utilities
postie demo                # Run demonstrations
postie help                # Show help
```

### Collection Format Features
- **Schema Validation**: JSON schema for collection validation
- **Environment Support**: Multiple environment configurations
- **Variable Interpolation**: Dynamic variable replacement
- **Nested Structure**: Hierarchical organization with API groups
- **Metadata Support**: Rich metadata for collections and requests
- **Authentication**: Per-request and collection-level auth
- **Headers Management**: Global and request-specific headers

### Known Issues
- External service dependencies for demo (httpbin.org, jsonplaceholder.typicode.com)
- Collection import from other tools not yet implemented

### Migration Notes
- This is the initial release, no migration needed

---

## Release Notes Template

### [Version] - YYYY-MM-DD

### Added
- New features

### Changed  
- Changes in existing functionality

### Deprecated
- Soon-to-be removed features

### Removed
- Removed features

### Fixed
- Bug fixes

### Security
- Security improvements

---

## Contributing

Please read [CONTRIBUTING.md](CONTRIBUTING.md) for details on our code of conduct and the process for submitting pull requests.

## Versioning

We use [SemVer](http://semver.org/) for versioning. For the versions available, see the [tags on this repository](https://github.com/yourusername/postie/tags).