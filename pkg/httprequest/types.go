package httprequest

import (
	"fmt"
	"net/http"
	"strings"
)

// RequestsFile represents the top-level structure of an HTTP requests file
type RequestsFile struct {
	Requests []Request `json:"requests"`
}

// Request represents a complete HTTP request with all its components
type Request struct {
	Name            string           `json:"name,omitempty"`             // From ### comments
	Method          string           `json:"method"`                     // HTTP method (GET, POST, etc.)
	URL             *URL             `json:"url"`                        // Request target
	HTTPVersion     string           `json:"http_version,omitempty"`     // HTTP version (optional)
	Headers         []Header         `json:"headers,omitempty"`          // Request headers
	Body            *RequestBody     `json:"body,omitempty"`             // Request body
	ResponseHandler *ResponseHandler `json:"response_handler,omitempty"` // Response handler script
	ResponseRef     *ResponseRef     `json:"response_ref,omitempty"`     // Response reference
	Comments        []string         `json:"comments,omitempty"`         // Associated comments
	LineNumber      int              `json:"line_number,omitempty"`      // Line number in file
}

// URL represents the request target with all its components
type URL struct {
	Raw       string            `json:"raw"`                 // Original URL string
	Scheme    string            `json:"scheme,omitempty"`    // http, https
	Host      string            `json:"host,omitempty"`      // hostname or IP
	Port      string            `json:"port,omitempty"`      // port number
	Path      string            `json:"path,omitempty"`      // path segments
	Query     map[string]string `json:"query,omitempty"`     // query parameters
	Fragment  string            `json:"fragment,omitempty"`  // URL fragment
	Variables []string          `json:"variables,omitempty"` // Found variables
}

// Header represents an HTTP header field
type Header struct {
	Name      string   `json:"name"`                // Header name
	Value     string   `json:"value"`               // Header value
	Variables []string `json:"variables,omitempty"` // Found variables in value
}

// RequestBody represents the message body of a request
type RequestBody struct {
	Type        BodyType         `json:"type"`                   // Type of body content
	Content     string           `json:"content,omitempty"`      // Inline content
	FilePath    string           `json:"file_path,omitempty"`    // File reference path
	Multipart   []MultipartField `json:"multipart,omitempty"`    // Multipart fields
	Variables   []string         `json:"variables,omitempty"`    // Found variables
	ContentType string           `json:"content_type,omitempty"` // Detected content type
}

// BodyType represents the type of request body
type BodyType string

const (
	BodyTypeInline    BodyType = "inline"
	BodyTypeFile      BodyType = "file"
	BodyTypeMultipart BodyType = "multipart"
)

// MultipartField represents a field in multipart form data
type MultipartField struct {
	Name      string   `json:"name"`                // Field name
	Headers   []Header `json:"headers,omitempty"`   // Field headers
	Content   string   `json:"content,omitempty"`   // Inline content
	FilePath  string   `json:"file_path,omitempty"` // File reference
	Variables []string `json:"variables,omitempty"` // Found variables
	Boundary  string   `json:"boundary,omitempty"`  // Multipart boundary
}

// ResponseHandler represents a response handler script
type ResponseHandler struct {
	Type     HandlerType `json:"type"`                // inline or file
	Script   string      `json:"script,omitempty"`    // Inline script content
	FilePath string      `json:"file_path,omitempty"` // Script file path
}

// HandlerType represents the type of response handler
type HandlerType string

const (
	HandlerTypeInline HandlerType = "inline"
	HandlerTypeFile   HandlerType = "file"
)

// ResponseRef represents a reference to a previous response
type ResponseRef struct {
	FilePath string `json:"file_path"` // Path to response file
}

// Token represents a lexical token
type Token struct {
	Type     TokenType `json:"type"`
	Value    string    `json:"value"`
	Line     int       `json:"line"`
	Column   int       `json:"column"`
	Position int       `json:"position"` // Absolute position in file
}

// TokenType represents the type of lexical token
type TokenType int

const (
	// Special tokens
	TokenEOF TokenType = iota
	TokenError
	TokenComment
	TokenNewline
	TokenWhitespace

	// Separators
	TokenRequestSeparator // ###

	// Request line tokens
	TokenMethod      // GET, POST, etc.
	TokenURL         // URL components
	TokenHTTPVersion // HTTP/1.1

	// Header tokens
	TokenHeaderName  // header name
	TokenColon       // :
	TokenHeaderValue // header value

	// Body tokens
	TokenBodyContent   // inline body content
	TokenFileReference // < ./file.json

	// Multipart tokens
	TokenBoundary         // --boundary
	TokenMultipartHeader  // multipart header
	TokenMultipartContent // multipart content

	// Response handler tokens
	TokenResponseHandlerStart // > {%
	TokenResponseHandlerEnd   // %}
	TokenResponseHandlerCode  // JavaScript code

	// Response reference tokens
	TokenResponseRefStart // <>
	TokenResponseRefPath  // file path

	// Variable tokens
	TokenVariableStart // {{
	TokenVariableEnd   // }}
	TokenVariableName  // variable name

	// Content tokens
	TokenText       // general text content
	TokenIdentifier // identifiers
	TokenString     // quoted strings
)

// String returns the string representation of a token type
func (t TokenType) String() string {
	switch t {
	case TokenEOF:
		return "EOF"
	case TokenError:
		return "ERROR"
	case TokenComment:
		return "COMMENT"
	case TokenNewline:
		return "NEWLINE"
	case TokenWhitespace:
		return "WHITESPACE"
	case TokenRequestSeparator:
		return "REQUEST_SEPARATOR"
	case TokenMethod:
		return "METHOD"
	case TokenURL:
		return "URL"
	case TokenHTTPVersion:
		return "HTTP_VERSION"
	case TokenHeaderName:
		return "HEADER_NAME"
	case TokenColon:
		return "COLON"
	case TokenHeaderValue:
		return "HEADER_VALUE"
	case TokenBodyContent:
		return "BODY_CONTENT"
	case TokenFileReference:
		return "FILE_REFERENCE"
	case TokenBoundary:
		return "BOUNDARY"
	case TokenMultipartHeader:
		return "MULTIPART_HEADER"
	case TokenMultipartContent:
		return "MULTIPART_CONTENT"
	case TokenResponseHandlerStart:
		return "RESPONSE_HANDLER_START"
	case TokenResponseHandlerEnd:
		return "RESPONSE_HANDLER_END"
	case TokenResponseHandlerCode:
		return "RESPONSE_HANDLER_CODE"
	case TokenResponseRefStart:
		return "RESPONSE_REF_START"
	case TokenResponseRefPath:
		return "RESPONSE_REF_PATH"
	case TokenVariableStart:
		return "VARIABLE_START"
	case TokenVariableEnd:
		return "VARIABLE_END"
	case TokenVariableName:
		return "VARIABLE_NAME"
	case TokenText:
		return "TEXT"
	case TokenIdentifier:
		return "IDENTIFIER"
	case TokenString:
		return "STRING"
	default:
		return "UNKNOWN"
	}
}

// ParseError represents a parsing error
type ParseError struct {
	Message  string `json:"message"`
	Line     int    `json:"line"`
	Column   int    `json:"column"`
	Position int    `json:"position"`
	Token    *Token `json:"token,omitempty"`
}

// Error implements the error interface
func (e *ParseError) Error() string {
	if e.Token != nil {
		return fmt.Sprintf("parse error at line %d, column %d: %s (token: %s)",
			e.Line, e.Column, e.Message, e.Token.Type.String())
	}
	return fmt.Sprintf("parse error at line %d, column %d: %s", e.Line, e.Column, e.Message)
}

// ValidationError represents a validation error
type ValidationError struct {
	Field   string   `json:"field"`
	Message string   `json:"message"`
	Request *Request `json:"request,omitempty"`
}

// Error implements the error interface
func (e *ValidationError) Error() string {
	if e.Request != nil && e.Request.Name != "" {
		return fmt.Sprintf("validation error in request '%s', field '%s': %s",
			e.Request.Name, e.Field, e.Message)
	}
	return fmt.Sprintf("validation error in field '%s': %s", e.Field, e.Message)
}

// HTTPMethod represents valid HTTP methods according to the spec
type HTTPMethod string

const (
	MethodGET     HTTPMethod = "GET"
	MethodHEAD    HTTPMethod = "HEAD"
	MethodPOST    HTTPMethod = "POST"
	MethodPUT     HTTPMethod = "PUT"
	MethodDELETE  HTTPMethod = "DELETE"
	MethodCONNECT HTTPMethod = "CONNECT"
	MethodPATCH   HTTPMethod = "PATCH"
	MethodOPTIONS HTTPMethod = "OPTIONS"
	MethodTRACE   HTTPMethod = "TRACE"
)

// ValidHTTPMethods contains all valid HTTP methods from the spec
var ValidHTTPMethods = map[string]bool{
	"GET":     true,
	"HEAD":    true,
	"POST":    true,
	"PUT":     true,
	"DELETE":  true,
	"CONNECT": true,
	"PATCH":   true,
	"OPTIONS": true,
	"TRACE":   true,
}

// ExecutionContext holds context for request execution
type ExecutionContext struct {
	Variables       map[string]interface{} `json:"variables"`
	Globals         map[string]interface{} `json:"globals"`
	Environment     string                 `json:"environment"`
	ResponseHistory []ExecutedResponse     `json:"response_history"`
	WorkingDir      string                 `json:"working_dir"`
}

// ExecutedResponse represents the result of executing a request
type ExecutedResponse struct {
	Request    *Request       `json:"request"`
	Response   *http.Response `json:"-"` // Don't serialize the actual response
	StatusCode int            `json:"status_code"`
	Status     string         `json:"status"`
	Headers    http.Header    `json:"headers"`
	Body       []byte         `json:"body"`
	Duration   int64          `json:"duration_ms"` // Duration in milliseconds
	Timestamp  int64          `json:"timestamp"`   // Unix timestamp
	Error      string         `json:"error,omitempty"`
}

// String implements fmt.Stringer for pretty printing
func (r *Request) String() string {
	name := r.Name
	if name == "" {
		name = "Unnamed Request"
	}
	return fmt.Sprintf("%s: %s %s", name, r.Method, r.URL.Raw)
}

// GetContentType returns the content type from headers or detects it
func (rb *RequestBody) GetContentType() string {
	if rb.ContentType != "" {
		return rb.ContentType
	}

	// Try to detect from content for inline bodies
	if rb.Type == BodyTypeInline && rb.Content != "" {
		trimmed := strings.TrimSpace(rb.Content)
		if strings.HasPrefix(trimmed, "{") || strings.HasPrefix(trimmed, "[") {
			return "application/json"
		}
		if strings.HasPrefix(trimmed, "<") {
			return "application/xml"
		}
	}

	return "text/plain"
}

// HasVariables returns true if the URL contains template variables
func (u *URL) HasVariables() bool {
	return len(u.Variables) > 0
}

// GetVariables extracts all variable names from the URL
func (u *URL) GetVariables() []string {
	if u.Variables != nil {
		return u.Variables
	}

	// Extract variables from raw URL
	var variables []string
	content := u.Raw

	for {
		start := strings.Index(content, "{{")
		if start == -1 {
			break
		}
		end := strings.Index(content[start:], "}}")
		if end == -1 {
			break
		}

		varName := strings.TrimSpace(content[start+2 : start+end])
		if varName != "" {
			variables = append(variables, varName)
		}
		content = content[start+end+2:]
	}

	u.Variables = variables
	return variables
}

// IsValidMethod checks if the method is valid according to the spec
func (r *Request) IsValidMethod() bool {
	return ValidHTTPMethods[strings.ToUpper(r.Method)]
}

// GetAllVariables returns all variables used in the request
func (r *Request) GetAllVariables() []string {
	var variables []string

	// URL variables
	if r.URL != nil {
		variables = append(variables, r.URL.GetVariables()...)
	}

	// Header variables
	for _, header := range r.Headers {
		variables = append(variables, header.Variables...)
	}

	// Body variables
	if r.Body != nil {
		variables = append(variables, r.Body.Variables...)
	}

	// Remove duplicates
	seen := make(map[string]bool)
	var unique []string
	for _, v := range variables {
		if !seen[v] {
			seen[v] = true
			unique = append(unique, v)
		}
	}

	return unique
}
