package httprequest

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Validator validates parsed HTTP requests according to the specification
type Validator struct {
	strict     bool   // Enable strict validation mode
	workingDir string // Working directory for file path resolution
	errors     []ValidationError
}

// NewValidator creates a new validator
func NewValidator(strict bool, workingDir string) *Validator {
	if workingDir == "" {
		workingDir, _ = os.Getwd()
	}

	return &Validator{
		strict:     strict,
		workingDir: workingDir,
		errors:     make([]ValidationError, 0),
	}
}

// Validate validates a RequestsFile
func (v *Validator) Validate(requestsFile *RequestsFile) []ValidationError {
	v.errors = make([]ValidationError, 0)

	if requestsFile == nil {
		v.addError("", "RequestsFile is nil", nil)
		return v.errors
	}

	if len(requestsFile.Requests) == 0 {
		v.addError("", "No requests found in file", nil)
		return v.errors
	}

	// Validate each request
	for i, request := range requestsFile.Requests {
		v.validateRequest(&request, i)
	}

	// Check for duplicate request names
	v.validateUniqueNames(requestsFile.Requests)

	return v.errors
}

// validateRequest validates a single request
func (v *Validator) validateRequest(request *Request, index int) {
	if request == nil {
		v.addError("", fmt.Sprintf("Request at index %d is nil", index), nil)
		return
	}

	// Validate request method
	v.validateMethod(request)

	// Validate URL
	v.validateURL(request)

	// Validate HTTP version
	v.validateHTTPVersion(request)

	// Validate headers
	v.validateHeaders(request)

	// Validate body
	v.validateBody(request)

	// Validate response handler
	v.validateResponseHandler(request)

	// Validate response reference
	v.validateResponseReference(request)

	// Validate variables
	v.validateVariables(request)
}

// validateMethod validates the HTTP method
func (v *Validator) validateMethod(request *Request) {
	if request.Method == "" {
		v.addError("Method", "HTTP method is required", request)
		return
	}

	method := strings.ToUpper(request.Method)
	if !ValidHTTPMethods[method] {
		v.addError("Method", fmt.Sprintf("Invalid HTTP method: %s", request.Method), request)
		return
	}

	// Normalize method to uppercase
	request.Method = method
}

// validateURL validates the request URL
func (v *Validator) validateURL(request *Request) {
	if request.URL == nil {
		v.addError("URL", "URL is required", request)
		return
	}

	if request.URL.Raw == "" {
		v.addError("URL", "URL cannot be empty", request)
		return
	}

	// Skip validation for template variables and asterisk form
	if request.URL.Raw == "*" || strings.Contains(request.URL.Raw, "{{") {
		return
	}

	// Validate URL format
	if strings.HasPrefix(request.URL.Raw, "http://") || strings.HasPrefix(request.URL.Raw, "https://") {
		// Absolute URL
		_, err := url.Parse(request.URL.Raw)
		if err != nil {
			v.addError("URL", fmt.Sprintf("Invalid URL format: %s", err.Error()), request)
		}
	} else if strings.HasPrefix(request.URL.Raw, "/") {
		// Origin form - path only
		if v.strict && !v.hasHostHeader(request) {
			v.addError("URL", "Origin-form URL requires Host header", request)
		}
	} else if request.URL.Raw != "*" {
		v.addError("URL", "URL must be absolute, origin-form, or asterisk-form", request)
	}
}

// validateHTTPVersion validates the HTTP version
func (v *Validator) validateHTTPVersion(request *Request) {
	if request.HTTPVersion == "" {
		return // HTTP version is optional
	}

	validVersionRegex := regexp.MustCompile(`^HTTP/\d+\.\d+$`)
	if !validVersionRegex.MatchString(request.HTTPVersion) {
		v.addError("HTTPVersion", fmt.Sprintf("Invalid HTTP version format: %s", request.HTTPVersion), request)
	}
}

// validateHeaders validates request headers
func (v *Validator) validateHeaders(request *Request) {
	headerNames := make(map[string]bool)

	for i, header := range request.Headers {
		// Check for empty header name
		if header.Name == "" {
			v.addError("Headers", fmt.Sprintf("Header at index %d has empty name", i), request)
			continue
		}

		// Check for duplicate headers (case-insensitive)
		lowerName := strings.ToLower(header.Name)
		if headerNames[lowerName] {
			if v.strict {
				v.addError("Headers", fmt.Sprintf("Duplicate header: %s", header.Name), request)
			}
		}
		headerNames[lowerName] = true

		// Validate header name format
		if !v.isValidHeaderName(header.Name) {
			v.addError("Headers", fmt.Sprintf("Invalid header name: %s", header.Name), request)
		}

		// Validate specific headers
		v.validateSpecificHeader(header, request)
	}
}

// validateBody validates the request body
func (v *Validator) validateBody(request *Request) {
	if request.Body == nil {
		return // Body is optional
	}

	switch request.Body.Type {
	case BodyTypeInline:
		// For inline body, content should not be empty
		if request.Body.Content == "" && v.strict {
			v.addError("Body", "Inline body content is empty", request)
		}

	case BodyTypeFile:
		// Validate file reference
		if request.Body.FilePath == "" {
			v.addError("Body", "File path is required for file body type", request)
		} else {
			v.validateFilePath(request.Body.FilePath, "Body.FilePath", request)
		}

	case BodyTypeMultipart:
		// Validate multipart fields
		v.validateMultipartBody(request.Body, request)

	default:
		v.addError("Body", fmt.Sprintf("Invalid body type: %s", request.Body.Type), request)
	}

	// Validate body for methods that shouldn't have body
	if v.strict && (request.Method == "GET" || request.Method == "HEAD" || request.Method == "DELETE") {
		if request.Body != nil && request.Body.Content != "" {
			v.addError("Body", fmt.Sprintf("%s requests should not have a body", request.Method), request)
		}
	}
}

// validateMultipartBody validates multipart form data
func (v *Validator) validateMultipartBody(body *RequestBody, request *Request) {
	if len(body.Multipart) == 0 {
		v.addError("Body", "Multipart body must have at least one field", request)
		return
	}

	// Check for Content-Type header with boundary
	var hasContentType bool

	for _, header := range request.Headers {
		if strings.ToLower(header.Name) == "content-type" {
			hasContentType = true
			if strings.Contains(strings.ToLower(header.Value), "multipart/form-data") {
				// Extract boundary for potential future validation
				parts := strings.Split(header.Value, ";")
				for _, part := range parts {
					part = strings.TrimSpace(part)
					if strings.HasPrefix(part, "boundary=") {
						_ = strings.TrimPrefix(part, "boundary=") // boundary extracted but not used yet
						break
					}
				}
			}
		}
	}

	if v.strict && !hasContentType {
		v.addError("Body", "Multipart body requires Content-Type header", request)
	}

	// Validate each multipart field
	fieldNames := make(map[string]bool)
	for i, field := range body.Multipart {
		if field.Name == "" {
			v.addError("Body", fmt.Sprintf("Multipart field at index %d has no name", i), request)
		} else if fieldNames[field.Name] && v.strict {
			v.addError("Body", fmt.Sprintf("Duplicate multipart field name: %s", field.Name), request)
		}
		fieldNames[field.Name] = true

		// Validate field content
		if field.Content == "" && field.FilePath == "" {
			v.addError("Body", fmt.Sprintf("Multipart field '%s' has no content or file reference", field.Name), request)
		}

		// Validate file reference
		if field.FilePath != "" {
			v.validateFilePath(field.FilePath, fmt.Sprintf("Body.Multipart[%d].FilePath", i), request)
		}
	}
}

// validateResponseHandler validates response handler scripts
func (v *Validator) validateResponseHandler(request *Request) {
	if request.ResponseHandler == nil {
		return // Response handler is optional
	}

	switch request.ResponseHandler.Type {
	case HandlerTypeInline:
		if request.ResponseHandler.Script == "" {
			v.addError("ResponseHandler", "Inline response handler script is empty", request)
		}

		// Basic JavaScript syntax validation (if strict)
		if v.strict {
			v.validateJavaScript(request.ResponseHandler.Script, request)
		}

	case HandlerTypeFile:
		if request.ResponseHandler.FilePath == "" {
			v.addError("ResponseHandler", "File path is required for file response handler", request)
		} else {
			v.validateFilePath(request.ResponseHandler.FilePath, "ResponseHandler.FilePath", request)
		}

	default:
		v.addError("ResponseHandler", fmt.Sprintf("Invalid response handler type: %s", request.ResponseHandler.Type), request)
	}
}

// validateResponseReference validates response references
func (v *Validator) validateResponseReference(request *Request) {
	if request.ResponseRef == nil {
		return // Response reference is optional
	}

	if request.ResponseRef.FilePath == "" {
		v.addError("ResponseRef", "Response reference file path is required", request)
	}

	// In strict mode, validate that referenced response file exists
	if v.strict {
		v.validateFilePath(request.ResponseRef.FilePath, "ResponseRef.FilePath", request)
	}
}

// validateVariables validates variable usage
func (v *Validator) validateVariables(request *Request) {
	variables := request.GetAllVariables()

	for _, varName := range variables {
		if !v.isValidVariableName(varName) {
			v.addError("Variables", fmt.Sprintf("Invalid variable name: %s", varName), request)
		}
	}
}

// validateUniqueNames checks for duplicate request names
func (v *Validator) validateUniqueNames(requests []Request) {
	names := make(map[string]bool)

	for _, request := range requests {
		if request.Name != "" {
			if names[request.Name] {
				v.addError("Name", fmt.Sprintf("Duplicate request name: %s", request.Name), &request)
			}
			names[request.Name] = true
		}
	}
}

// Helper validation methods

// hasHostHeader checks if request has a Host header
func (v *Validator) hasHostHeader(request *Request) bool {
	for _, header := range request.Headers {
		if strings.ToLower(header.Name) == "host" {
			return true
		}
	}
	return false
}

// isValidHeaderName checks if header name is valid
func (v *Validator) isValidHeaderName(name string) bool {
	// HTTP header names should contain only ASCII letters, digits, and hyphens
	headerNameRegex := regexp.MustCompile(`^[a-zA-Z0-9-]+$`)
	return headerNameRegex.MatchString(name)
}

// validateSpecificHeader validates specific header types
func (v *Validator) validateSpecificHeader(header Header, request *Request) {
	switch strings.ToLower(header.Name) {
	case "content-length":
		// Content-Length should be a number
		if v.strict && header.Value != "" && !regexp.MustCompile(`^\d+$`).MatchString(header.Value) {
			v.addError("Headers", "Content-Length must be a number", request)
		}

	case "content-type":
		// Basic content-type validation
		if v.strict && header.Value != "" {
			// Should contain at least a media type
			if !strings.Contains(header.Value, "/") {
				v.addError("Headers", "Invalid Content-Type format", request)
			}
		}
	}
}

// validateFilePath validates that a file path exists and is readable
func (v *Validator) validateFilePath(path, field string, request *Request) {
	if path == "" {
		return
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(v.workingDir, path)
	}

	// Check if file exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		v.addError(field, fmt.Sprintf("File not found: %s", path), request)
	} else if err != nil {
		v.addError(field, fmt.Sprintf("Cannot access file: %s", err.Error()), request)
	}
}

// validateJavaScript performs basic JavaScript syntax validation
func (v *Validator) validateJavaScript(script string, request *Request) {
	// Basic checks for common JavaScript syntax errors
	script = strings.TrimSpace(script)

	if script == "" {
		return
	}

	// Check for unmatched braces
	braceCount := 0
	for _, char := range script {
		if char == '{' {
			braceCount++
		} else if char == '}' {
			braceCount--
		}
	}

	if braceCount != 0 {
		v.addError("ResponseHandler", "Unmatched braces in JavaScript code", request)
	}

	// Check for common typos in API usage
	if strings.Contains(script, "client.") {
		// This is probably using the client API - basic validation
		if !strings.Contains(script, "client.test") && !strings.Contains(script, "client.assert") &&
			!strings.Contains(script, "client.global") {
			// Warning: might be valid code, but let's check
		}
	}
}

// isValidVariableName checks if variable name is valid
func (v *Validator) isValidVariableName(name string) bool {
	if name == "" {
		return false
	}

	// Variable names should start with letter or underscore, followed by letters, digits, underscore, or hyphen
	variableNameRegex := regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_-]*$`)
	return variableNameRegex.MatchString(name)
}

// addError adds a validation error
func (v *Validator) addError(field, message string, request *Request) {
	error := ValidationError{
		Field:   field,
		Message: message,
		Request: request,
	}
	v.errors = append(v.errors, error)
}

// ValidateFile validates an HTTP request file
func ValidateFile(filename string, content string, strict bool) ([]ValidationError, error) {
	requestsFile, err := ParseFile(filename, content)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	workingDir := filepath.Dir(filename)
	validator := NewValidator(strict, workingDir)

	return validator.Validate(requestsFile), nil
}

// IsValid returns true if there are no validation errors
func IsValid(errors []ValidationError) bool {
	return len(errors) == 0
}

// FormatErrors formats validation errors for display
func FormatErrors(errors []ValidationError) string {
	if len(errors) == 0 {
		return ""
	}

	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("Found %d validation error(s):\n", len(errors)))

	for i, err := range errors {
		builder.WriteString(fmt.Sprintf("%d. %s\n", i+1, err.Error()))
	}

	return builder.String()
}
