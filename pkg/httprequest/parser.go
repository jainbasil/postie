package httprequest

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Parser parses HTTP request files into structured data
type Parser struct {
	tokens   []Token
	position int
	current  Token
	file     string
}

// NewParser creates a new parser for the given tokens
func NewParser(tokens []Token, filename string) *Parser {
	parser := &Parser{
		tokens: tokens,
		file:   filename,
	}

	if len(tokens) > 0 {
		parser.current = tokens[0]
	}

	return parser
}

// ParseFile parses an HTTP request file and returns the structured data
func ParseFile(filename string, content string) (*RequestsFile, error) {
	lexer := NewLexer(content)
	tokens, err := lexer.Tokenize()
	if err != nil {
		return nil, fmt.Errorf("lexing error: %w", err)
	}

	parser := NewParser(tokens, filename)
	return parser.Parse()
}

// Parse parses the tokens into a RequestsFile
func (p *Parser) Parse() (*RequestsFile, error) {
	var requests []Request

	// Skip initial request separators and whitespace
	p.skipIgnorable()

	var pendingRequestName string

	for !p.isAtEnd() {
		// Check for request separators
		if p.check(TokenRequestSeparator) {
			// Extract name from request separator if present
			separatorValue := p.current.Value
			p.advance() // consume the separator
			pendingRequestName = p.extractRequestName(separatorValue)
			p.skipIgnorable()
			continue
		}

		// Check if we have tokens that could start a request
		if !p.hasValidRequestStart() {
			// Skip tokens that don't start a request
			p.advance()
			continue
		}

		request, err := p.parseRequest()
		if err != nil {
			return nil, err
		}

		if request != nil {
			// Apply pending name to the request
			if pendingRequestName != "" {
				request.Name = pendingRequestName
				pendingRequestName = "" // Clear the pending name
			}

			requests = append(requests, *request)
		}

		p.skipIgnorable()
	}

	return &RequestsFile{
		Requests: requests,
	}, nil
}

// parseRequest parses a single HTTP request
func (p *Parser) parseRequest() (*Request, error) {
	if p.isAtEnd() {
		return nil, nil
	}

	request := &Request{
		LineNumber: p.current.Line,
	}

	// Parse optional name from preceding comment
	if p.check(TokenComment) {
		request.Name = p.extractRequestName(p.current.Value)
		p.advance()
		p.skipIgnorable()
	}

	// Parse request line
	if err := p.parseRequestLine(request); err != nil {
		return nil, err
	}

	p.skipNewlines()

	// Parse headers
	if err := p.parseHeaders(request); err != nil {
		return nil, err
	}

	// Check for empty line before body
	if p.check(TokenNewline) {
		p.advance() // consume the empty line that separates headers from body

		// Skip any additional newlines
		for p.check(TokenNewline) {
			p.advance()
		}

		// Parse body content
		if !p.isAtEnd() && !p.check(TokenRequestSeparator) &&
			!p.check(TokenResponseHandlerStart) && !p.check(TokenResponseRefStart) &&
			!p.check(TokenMethod) {
			if err := p.parseBody(request); err != nil {
				return nil, err
			}
		}
	}

	// Parse response handler
	if p.check(TokenResponseHandlerStart) ||
		(p.check(TokenText) && strings.HasPrefix(strings.TrimSpace(p.current.Value), ">")) {
		if err := p.parseResponseHandler(request); err != nil {
			return nil, err
		}
	}

	// Parse response reference
	if p.check(TokenResponseRefStart) {
		if err := p.parseResponseReference(request); err != nil {
			return nil, err
		}
	}

	return request, nil
}

// parseRequestLine parses the HTTP request line
func (p *Parser) parseRequestLine(request *Request) error {
	// Method is optional (defaults to GET)
	if p.check(TokenMethod) {
		request.Method = p.current.Value
		p.advance()
	} else {
		request.Method = "GET"
	}

	// URL is required
	if !p.check(TokenURL) && !p.check(TokenVariableStart) && !p.check(TokenText) {
		return p.error("expected URL after method")
	}

	// Handle different URL token types
	var urlStr string
	if p.check(TokenURL) {
		urlStr = p.current.Value
		p.advance()
	} else if p.check(TokenVariableStart) {
		// Handle URL starting with variable
		urlStr = p.collectURLTokens()
	} else if p.check(TokenText) && p.looksLikeURL(p.current.Value) {
		urlStr = p.collectURLTokens()
	} else {
		return p.error("expected valid URL")
	}

	url, err := p.parseURL(urlStr)
	if err != nil {
		return err
	}
	request.URL = url

	// HTTP version is optional
	if p.check(TokenHTTPVersion) {
		request.HTTPVersion = p.current.Value
		p.advance()
	}

	return nil
}

// parseURL parses a URL string into a URL struct
func (p *Parser) parseURL(urlStr string) (*URL, error) {
	url := &URL{
		Raw: urlStr,
	}

	// Extract variables
	url.GetVariables()

	// Parse URL components
	if strings.HasPrefix(urlStr, "http://") || strings.HasPrefix(urlStr, "https://") {
		// Absolute URL
		parts := strings.SplitN(urlStr, "://", 2)
		url.Scheme = parts[0]

		if len(parts) > 1 {
			remaining := parts[1]

			// Extract host and port
			pathStart := strings.Index(remaining, "/")
			queryStart := strings.Index(remaining, "?")
			fragmentStart := strings.Index(remaining, "#")

			hostEnd := len(remaining)
			if pathStart != -1 && pathStart < hostEnd {
				hostEnd = pathStart
			}
			if queryStart != -1 && queryStart < hostEnd {
				hostEnd = queryStart
			}
			if fragmentStart != -1 && fragmentStart < hostEnd {
				hostEnd = fragmentStart
			}

			hostPort := remaining[:hostEnd]
			if strings.Contains(hostPort, ":") {
				parts := strings.SplitN(hostPort, ":", 2)
				url.Host = parts[0]
				url.Port = parts[1]
			} else {
				url.Host = hostPort
			}

			// Extract path
			if pathStart != -1 {
				pathEnd := len(remaining)
				if queryStart != -1 && queryStart > pathStart {
					pathEnd = queryStart
				}
				if fragmentStart != -1 && fragmentStart > pathStart {
					pathEnd = fragmentStart
				}
				url.Path = remaining[pathStart:pathEnd]
			}

			// Extract query
			if queryStart != -1 {
				queryEnd := len(remaining)
				if fragmentStart != -1 && fragmentStart > queryStart {
					queryEnd = fragmentStart
				}

				queryStr := remaining[queryStart+1 : queryEnd]
				url.Query = p.parseQuery(queryStr)
			}

			// Extract fragment
			if fragmentStart != -1 {
				url.Fragment = remaining[fragmentStart+1:]
			}
		}
	} else if strings.HasPrefix(urlStr, "/") {
		// Path-only URL
		url.Path = urlStr
	} else if urlStr == "*" {
		// Asterisk form
		url.Path = "*"
	}

	return url, nil
}

// parseQuery parses query string into key-value pairs
func (p *Parser) parseQuery(queryStr string) map[string]string {
	query := make(map[string]string)

	if queryStr == "" {
		return query
	}

	pairs := strings.Split(queryStr, "&")
	for _, pair := range pairs {
		if strings.Contains(pair, "=") {
			parts := strings.SplitN(pair, "=", 2)
			query[parts[0]] = parts[1]
		} else {
			query[pair] = ""
		}
	}

	return query
}

// parseHeaders parses HTTP headers
func (p *Parser) parseHeaders(request *Request) error {
	for !p.isAtEnd() && !p.check(TokenRequestSeparator) &&
		!p.check(TokenResponseHandlerStart) && !p.check(TokenResponseRefStart) {

		// Check for empty line (end of headers)
		if p.check(TokenNewline) {
			// Look ahead to see if this is an empty line or just a newline between headers
			return nil // Stop parsing headers when we hit a newline (empty line separator)
		}

		// Check if this line looks like a header
		if p.check(TokenText) && strings.Contains(p.current.Value, ":") {
			header, err := p.parseHeader()
			if err != nil {
				return err
			}
			request.Headers = append(request.Headers, *header)
		} else if p.check(TokenHeaderName) {
			header, err := p.parseHeader()
			if err != nil {
				return err
			}
			request.Headers = append(request.Headers, *header)
		} else {
			// Not a header, stop parsing headers
			break
		}
	}

	return nil
}

// parseHeader parses a single HTTP header
func (p *Parser) parseHeader() (*Header, error) {
	var name, value string

	if p.check(TokenHeaderName) {
		name = p.current.Value
		p.advance()

		if !p.check(TokenColon) {
			return nil, p.error("expected ':' after header name")
		}
		p.advance()

		// Collect header value tokens (including variables) until newline
		var valueParts []string

		// Skip any immediate whitespace after the colon
		for p.check(TokenText) && strings.TrimSpace(p.current.Value) == "" {
			p.advance()
		}

		for !p.check(TokenNewline) && !p.isAtEnd() && !p.check(TokenRequestSeparator) {
			if p.check(TokenVariableStart) {
				// Collect complete variable
				valueParts = append(valueParts, p.current.Value) // {{
				p.advance()
				if p.check(TokenVariableName) {
					valueParts = append(valueParts, p.current.Value) // name
					p.advance()
				}
				if p.check(TokenVariableEnd) {
					valueParts = append(valueParts, p.current.Value) // }}
					p.advance()
				}
			} else if p.check(TokenText) {
				valueParts = append(valueParts, p.current.Value)
				p.advance()
			} else if p.check(TokenMethod) || p.check(TokenURL) {
				// Stop if we hit tokens that could start a new request
				break
			} else {
				// Skip whitespace and other tokens but consume them
				p.advance()
			}
		}

		value = strings.TrimSpace(strings.Join(valueParts, ""))

	} else if p.check(TokenText) && strings.Contains(p.current.Value, ":") {
		// Parse "Name: Value" format (most common case)
		text := p.current.Value
		colonIndex := strings.Index(text, ":")
		name = strings.TrimSpace(text[:colonIndex])

		// Get value part after colon
		if colonIndex+1 < len(text) {
			value = strings.TrimSpace(text[colonIndex+1:])
		}
		p.advance()

		// Continue collecting value tokens (like variables) until newline
		var valueParts []string
		if value != "" {
			valueParts = append(valueParts, value)
		}

		for !p.check(TokenNewline) && !p.isAtEnd() && !p.check(TokenRequestSeparator) {
			if p.check(TokenVariableStart) {
				// Collect complete variable
				valueParts = append(valueParts, p.current.Value) // {{
				p.advance()
				if p.check(TokenVariableName) {
					valueParts = append(valueParts, p.current.Value) // name
					p.advance()
				}
				if p.check(TokenVariableEnd) {
					valueParts = append(valueParts, p.current.Value) // }}
					p.advance()
				}
			} else if p.check(TokenText) {
				valueParts = append(valueParts, p.current.Value)
				p.advance()
			} else {
				p.advance()
			}
		}

		value = strings.TrimSpace(strings.Join(valueParts, ""))

	} else {
		return nil, p.error("expected header name")
	}

	header := &Header{
		Name:  name,
		Value: value,
	}

	// Extract variables from header value
	header.Variables = p.extractVariables(value)

	return header, nil
}

// parseBody parses the request body
func (p *Parser) parseBody(request *Request) error {
	if p.isAtEnd() {
		return nil
	}

	// Check for file reference (< ./file.json)
	if p.check(TokenFileReference) {
		request.Body = &RequestBody{
			Type:     BodyTypeFile,
			FilePath: p.current.Value,
		}
		p.advance()
		return nil
	}

	// Check for text that starts with < (file reference)
	if p.check(TokenText) && strings.HasPrefix(strings.TrimSpace(p.current.Value), "<") {
		text := strings.TrimSpace(p.current.Value)
		if strings.HasPrefix(text, "< ") {
			filePath := strings.TrimSpace(text[2:])
			request.Body = &RequestBody{
				Type:     BodyTypeFile,
				FilePath: filePath,
			}
			p.advance()
			return nil
		}
	}

	// Check for multipart boundary
	if p.check(TokenBoundary) {
		return p.parseMultipartBody(request)
	}

	// Parse inline body - collect all remaining content until next section
	var bodyLines []string
	for !p.isAtEnd() && !p.check(TokenRequestSeparator) &&
		!p.check(TokenResponseHandlerStart) && !p.check(TokenResponseRefStart) {

		if p.check(TokenText) {
			bodyLines = append(bodyLines, p.current.Value)
		} else if p.check(TokenNewline) {
			bodyLines = append(bodyLines, "\n")
		} else if p.check(TokenVariableStart) {
			// Handle variables in body
			bodyLines = append(bodyLines, p.current.Value) // {{
			p.advance()
			if p.check(TokenVariableName) {
				bodyLines = append(bodyLines, p.current.Value) // name
				p.advance()
			}
			if p.check(TokenVariableEnd) {
				bodyLines = append(bodyLines, p.current.Value) // }}
				p.advance()
			}
			continue
		}
		p.advance()
	}

	if len(bodyLines) > 0 {
		content := strings.Join(bodyLines, "")
		content = strings.TrimSpace(content)

		if content != "" {
			request.Body = &RequestBody{
				Type:    BodyTypeInline,
				Content: content,
			}

			// Extract variables from body
			request.Body.Variables = p.extractVariables(content)

			// Set content type
			request.Body.ContentType = request.Body.GetContentType()
		}
	}

	return nil
}

// parseMultipartBody parses multipart form data
func (p *Parser) parseMultipartBody(request *Request) error {
	var fields []MultipartField
	var boundary string

	// Extract boundary from first boundary token
	if p.check(TokenBoundary) {
		boundary = strings.TrimPrefix(p.current.Value, "--")
		p.advance()
	}

	for p.check(TokenBoundary) {
		field := MultipartField{
			Boundary: boundary,
		}

		p.advance() // skip boundary
		p.skipNewlines()

		// Parse field headers
		for p.check(TokenHeaderName) || (p.check(TokenText) && strings.Contains(p.current.Value, ":")) {
			header, err := p.parseHeader()
			if err != nil {
				return err
			}
			field.Headers = append(field.Headers, *header)

			// Extract field name from Content-Disposition header
			if strings.ToLower(header.Name) == "content-disposition" {
				field.Name = p.extractFormFieldName(header.Value)
			}

			p.skipNewlines()
		}

		p.skipNewlines()

		// Parse field content
		if p.check(TokenFileReference) {
			field.FilePath = p.current.Value
			p.advance()
		} else {
			var contentLines []string
			for !p.isAtEnd() && !p.check(TokenBoundary) && !p.check(TokenRequestSeparator) {
				if p.check(TokenText) || p.check(TokenNewline) {
					contentLines = append(contentLines, p.current.Value)
				}
				p.advance()
			}

			field.Content = strings.TrimSpace(strings.Join(contentLines, ""))
		}

		// Extract variables
		field.Variables = p.extractVariables(field.Content)

		fields = append(fields, field)
	}

	request.Body = &RequestBody{
		Type:      BodyTypeMultipart,
		Multipart: fields,
	}

	return nil
}

// parseResponseHandler parses response handler scripts
func (p *Parser) parseResponseHandler(request *Request) error {
	if p.check(TokenResponseHandlerStart) {
		// Inline handler
		p.advance() // consume start token

		if !p.check(TokenResponseHandlerCode) {
			return p.error("expected response handler code")
		}

		script := p.current.Value
		p.advance()

		if !p.check(TokenResponseHandlerEnd) {
			return p.error("expected response handler end")
		}
		p.advance()

		request.ResponseHandler = &ResponseHandler{
			Type:   HandlerTypeInline,
			Script: script,
		}
	} else if p.check(TokenText) && strings.HasPrefix(strings.TrimSpace(p.current.Value), ">") {
		// File handler
		text := strings.TrimSpace(p.current.Value)
		if strings.HasPrefix(text, "> ") {
			filePath := strings.TrimSpace(text[2:])
			request.ResponseHandler = &ResponseHandler{
				Type:     HandlerTypeFile,
				FilePath: filePath,
			}
		}
		p.advance()
	}

	return nil
}

// parseResponseReference parses response references
func (p *Parser) parseResponseReference(request *Request) error {
	if !p.check(TokenResponseRefStart) {
		return p.error("expected response reference start")
	}
	p.advance()

	if !p.check(TokenResponseRefPath) {
		return p.error("expected response reference path")
	}

	request.ResponseRef = &ResponseRef{
		FilePath: p.current.Value,
	}
	p.advance()

	return nil
}

// Helper methods

// advance moves to the next token
func (p *Parser) advance() {
	if !p.isAtEnd() {
		p.position++
		if p.position < len(p.tokens) {
			p.current = p.tokens[p.position]
		}
	}
}

// isAtEnd checks if we're at the end of tokens
func (p *Parser) isAtEnd() bool {
	return p.position >= len(p.tokens) || p.current.Type == TokenEOF
}

// check returns true if current token is of the given type
func (p *Parser) check(tokenType TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.current.Type == tokenType
}

// match consumes token if it matches the given type
func (p *Parser) match(tokenType TokenType) bool {
	if p.check(tokenType) {
		p.advance()
		return true
	}
	return false
}

// checkEmpty checks if current line is empty (only contains newline)
func (p *Parser) checkEmpty() bool {
	return p.check(TokenNewline)
}

// skipIgnorable skips whitespace, newlines, and comments
func (p *Parser) skipIgnorable() {
	for p.check(TokenWhitespace) || p.check(TokenNewline) || p.check(TokenComment) {
		p.advance()
	}
}

// skipNewlines skips newline tokens
func (p *Parser) skipNewlines() {
	for p.check(TokenNewline) {
		p.advance()
	}
}

// skipWhitespace skips whitespace tokens
func (p *Parser) skipWhitespace() {
	for p.check(TokenWhitespace) {
		p.advance()
	}
}

// error creates a parse error
func (p *Parser) error(message string) error {
	return &ParseError{
		Message:  message,
		Line:     p.current.Line,
		Column:   p.current.Column,
		Position: p.current.Position,
		Token:    &p.current,
	}
}

// extractRequestName extracts request name from comment
func (p *Parser) extractRequestName(comment string) string {
	// Remove comment markers and clean up
	name := strings.TrimSpace(comment)

	// Remove ### prefix (for request separators)
	if strings.HasPrefix(name, "###") {
		name = strings.TrimSpace(name[3:])
	} else {
		// Remove other comment prefixes
		name = strings.TrimPrefix(name, "//")
		name = strings.TrimPrefix(name, "#")
		name = strings.TrimSpace(name)
	}

	return name
}

// extractVariables extracts {{variable}} references from text
func (p *Parser) extractVariables(text string) []string {
	var variables []string
	content := text

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

	return variables
}

// extractFormFieldName extracts field name from Content-Disposition header
func (p *Parser) extractFormFieldName(value string) string {
	// Look for name="fieldname" pattern
	parts := strings.Split(value, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "name=") {
			name := strings.TrimPrefix(part, "name=")
			name = strings.Trim(name, `"'`)
			return name
		}
	}
	return ""
}

// looksLikeURL checks if text looks like a URL
func (p *Parser) looksLikeURL(text string) bool {
	text = strings.TrimSpace(text)
	return strings.HasPrefix(text, "http://") ||
		strings.HasPrefix(text, "https://") ||
		strings.HasPrefix(text, "/") ||
		strings.HasPrefix(text, "{{") ||
		text == "*"
}

// collectURLTokens collects tokens that form a complete URL
func (p *Parser) collectURLTokens() string {
	var parts []string

	// Collect tokens until we hit HTTP version, newline, or end
	for !p.isAtEnd() && !p.check(TokenNewline) && !p.check(TokenHTTPVersion) {
		if p.check(TokenURL) || p.check(TokenText) ||
			p.check(TokenVariableStart) || p.check(TokenVariableName) || p.check(TokenVariableEnd) {

			// For variables, collect the complete {{name}} sequence
			if p.check(TokenVariableStart) {
				parts = append(parts, p.current.Value) // {{
				p.advance()
				if p.check(TokenVariableName) {
					parts = append(parts, p.current.Value) // name
					p.advance()
				}
				if p.check(TokenVariableEnd) {
					parts = append(parts, p.current.Value) // }}
					p.advance()
				}
			} else {
				parts = append(parts, p.current.Value)
				p.advance()
			}
		} else {
			break
		}
	}

	return strings.Join(parts, "")
}

// hasValidRequestStart checks if current token can start a request
func (p *Parser) hasValidRequestStart() bool {
	// Check for HTTP method
	if p.check(TokenMethod) {
		return true
	}

	// Check for comment (### request name)
	if p.check(TokenComment) {
		return true
	}

	// Check for text that could be a method or URL
	if p.check(TokenText) {
		text := strings.TrimSpace(p.current.Value)

		// Check if it's a known HTTP method
		words := strings.Fields(text)
		if len(words) > 0 {
			method := strings.ToUpper(words[0])
			if ValidHTTPMethods[method] {
				return true
			}
		}

		// Check if it looks like a URL (for GET with default method)
		if p.looksLikeURL(text) {
			return true
		}
	}

	// Check for variable that could start a URL
	if p.check(TokenVariableStart) {
		return true
	}

	return false
}

// GetAbsolutePath resolves relative file paths relative to the request file
func (p *Parser) GetAbsolutePath(path string) string {
	if filepath.IsAbs(path) {
		return path
	}

	if p.file != "" {
		dir := filepath.Dir(p.file)
		return filepath.Join(dir, path)
	}

	return path
}
