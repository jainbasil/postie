package httprequest

import (
	"strings"
	"testing"
)

func TestLexerBasicTokenization(t *testing.T) {
	input := "GET https://example.com"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()

	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	if len(tokens) < 3 {
		t.Fatalf("Expected at least 3 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TokenMethod || tokens[0].Value != "GET" {
		t.Errorf("Expected METHOD token 'GET', got %s '%s'", tokens[0].Type.String(), tokens[0].Value)
	}

	if tokens[1].Type != TokenURL || tokens[1].Value != "https://example.com" {
		t.Errorf("Expected URL token 'https://example.com', got %s '%s'", tokens[1].Type.String(), tokens[1].Value)
	}
}

func TestLexerRequestSeparator(t *testing.T) {
	input := "### Test Request"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()

	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	if len(tokens) < 1 {
		t.Fatalf("Expected at least 1 token, got %d", len(tokens))
	}

	if tokens[0].Type != TokenRequestSeparator {
		t.Errorf("Expected REQUEST_SEPARATOR token, got %s", tokens[0].Type.String())
	}
}

func TestLexerVariable(t *testing.T) {
	input := "{{baseUrl}}"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()

	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	if len(tokens) < 3 {
		t.Fatalf("Expected at least 3 tokens, got %d", len(tokens))
	}

	if tokens[0].Type != TokenVariableStart {
		t.Errorf("Expected VARIABLE_START token, got %s", tokens[0].Type.String())
	}

	if tokens[1].Type != TokenVariableName || tokens[1].Value != "baseUrl" {
		t.Errorf("Expected VARIABLE_NAME 'baseUrl', got %s '%s'", tokens[1].Type.String(), tokens[1].Value)
	}

	if tokens[2].Type != TokenVariableEnd {
		t.Errorf("Expected VARIABLE_END token, got %s", tokens[2].Type.String())
	}
}

func TestLexerResponseHandler(t *testing.T) {
	input := "> {% client.test() %}"
	lexer := NewLexer(input)
	tokens, err := lexer.Tokenize()

	if err != nil {
		t.Fatalf("Lexer error: %v", err)
	}

	if len(tokens) < 1 {
		t.Fatalf("Expected at least 1 token, got %d", len(tokens))
	}

	if tokens[0].Type != TokenResponseHandlerStart {
		t.Errorf("Expected RESPONSE_HANDLER_START token, got %s", tokens[0].Type.String())
	}
}

func TestParserSimpleRequest(t *testing.T) {
	input := "GET https://api.example.com/users"

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.Method != "GET" {
		t.Errorf("Expected method GET, got %s", req.Method)
	}

	if req.URL == nil {
		t.Fatal("Expected URL to be parsed")
	}

	if req.URL.Raw != "https://api.example.com/users" {
		t.Errorf("Expected URL 'https://api.example.com/users', got '%s'", req.URL.Raw)
	}
}

func TestParserRequestWithBody(t *testing.T) {
	input := `POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John"
}`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.Method != "POST" {
		t.Errorf("Expected method POST, got %s", req.Method)
	}

	if req.Body == nil {
		t.Fatal("Expected request body")
	}

	if req.Body.Type != BodyTypeInline {
		t.Errorf("Expected inline body, got %s", req.Body.Type)
	}

	if len(req.Headers) == 0 {
		t.Error("Expected headers")
	}
}

func TestParserFileReference(t *testing.T) {
	input := `POST https://api.example.com/upload
Content-Type: application/json

< ./data.json`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.Body == nil {
		t.Fatal("Expected request body, got nil")
	}

	if req.Body.Type != BodyTypeFile {
		t.Errorf("Expected file body, got %s", req.Body.Type)
	}

	if req.Body.FilePath != "./data.json" {
		t.Errorf("Expected file path './data.json', got '%s'", req.Body.FilePath)
	}
}

func TestParserVariables(t *testing.T) {
	input := "GET {{baseUrl}}/api/v1/users?page={{page}}"

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.URL == nil {
		t.Fatal("Expected URL to be parsed")
	}

	variables := req.GetAllVariables()
	expectedVars := []string{"baseUrl", "page"}

	for _, expectedVar := range expectedVars {
		found := false
		for _, variable := range variables {
			if variable == expectedVar {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected variable '%s' not found", expectedVar)
		}
	}
}

func TestParserMultipleRequests(t *testing.T) {
	input := `### Get Users
GET https://api.example.com/users

###
POST https://api.example.com/users
Content-Type: application/json

{
  "name": "John"
}

### Delete User
DELETE https://api.example.com/users/123`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 3 {
		t.Fatalf("Expected 3 requests, got %d", len(requestsFile.Requests))
	}

	// Check first request
	req1 := requestsFile.Requests[0]
	if req1.Name != "Get Users" {
		t.Errorf("Expected first request name 'Get Users', got '%s'", req1.Name)
	}
	if req1.Method != "GET" {
		t.Errorf("Expected first request method GET, got %s", req1.Method)
	}

	// Check second request (no name)
	req2 := requestsFile.Requests[1]
	if req2.Name != "" {
		t.Errorf("Expected second request name to be empty, got '%s'", req2.Name)
	}
	if req2.Method != "POST" {
		t.Errorf("Expected second request method POST, got %s", req2.Method)
	}

	// Check third request
	req3 := requestsFile.Requests[2]
	if req3.Name != "Delete User" {
		t.Errorf("Expected third request name 'Delete User', got '%s'", req3.Name)
	}
	if req3.Method != "DELETE" {
		t.Errorf("Expected third request method DELETE, got %s", req3.Method)
	}
}

func TestParserResponseHandler(t *testing.T) {
	input := `GET https://api.example.com/users

> {%
    client.test("Request successful", function() {
        client.assert(response.status === 200);
    });
%}`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.ResponseHandler == nil {
		t.Fatal("Expected response handler")
	}

	if !strings.Contains(req.ResponseHandler.Script, "client.test") {
		t.Error("Expected response handler to contain 'client.test'")
	}
}

func TestParserResponseReference(t *testing.T) {
	input := `GET https://api.example.com/users

<> response.json`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]
	if req.ResponseRef == nil {
		t.Fatal("Expected response reference")
	}

	if req.ResponseRef.FilePath != "response.json" {
		t.Errorf("Expected response file 'response.json', got '%s'", req.ResponseRef.FilePath)
	}
}

func TestValidatorBasicValidation(t *testing.T) {
	request := &Request{
		Method: "GET",
		URL: &URL{
			Raw: "https://api.example.com/users",
		},
	}

	requestsFile := &RequestsFile{
		Requests: []Request{*request},
	}

	validator := NewValidator(true, "")
	errors := validator.Validate(requestsFile)

	if len(errors) != 0 {
		t.Errorf("Expected no validation errors, got %d: %v", len(errors), errors)
	}
}

func TestValidatorInvalidMethod(t *testing.T) {
	request := &Request{
		Method: "INVALID",
		URL: &URL{
			Raw: "https://api.example.com/users",
		},
	}

	requestsFile := &RequestsFile{
		Requests: []Request{*request},
	}

	validator := NewValidator(true, "")
	errors := validator.Validate(requestsFile)

	if len(errors) == 0 {
		t.Error("Expected validation errors for invalid method")
	}
}

func TestValidatorMissingURL(t *testing.T) {
	request := &Request{
		Method: "GET",
		URL:    nil,
	}

	requestsFile := &RequestsFile{
		Requests: []Request{*request},
	}

	validator := NewValidator(true, "")
	errors := validator.Validate(requestsFile)

	if len(errors) == 0 {
		t.Error("Expected validation errors for missing URL")
	}
}

func TestValidatorInvalidHeaderName(t *testing.T) {
	request := &Request{
		Method: "GET",
		URL: &URL{
			Raw: "https://api.example.com/users",
		},
		Headers: []Header{
			{Name: "Invalid Header Name With Spaces"},
		},
	}

	requestsFile := &RequestsFile{
		Requests: []Request{*request},
	}

	validator := NewValidator(true, "")
	errors := validator.Validate(requestsFile)

	if len(errors) == 0 {
		t.Error("Expected validation errors for invalid header name")
	}

	found := false
	for _, err := range errors {
		if strings.Contains(err.Message, "Invalid header name") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected error about invalid header name")
	}
}

func TestEndToEndParsing(t *testing.T) {
	// Simpler test to debug the issue
	input := `GET https://example.com HTTP/1.1
Host: {{host}}

POST https://example.com/create`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 2 {
		t.Fatalf("Expected 2 requests, got %d", len(requestsFile.Requests))
	}

	// Check first request
	req1 := requestsFile.Requests[0]
	if req1.Method != "GET" {
		t.Errorf("Expected GET method, got %s", req1.Method)
	}
	if len(req1.Headers) == 0 {
		t.Error("Expected headers on first request")
	}

	// Check second request
	req2 := requestsFile.Requests[1]
	if req2.Method != "POST" {
		t.Errorf("Expected POST method, got %s", req2.Method)
	}
}

func TestParserMultipleHeadersWithVariables(t *testing.T) {
	// Test case for the reported issue: multiple headers with variables
	input := `### Q1
GET http://localhost:5001/plugin?userMessage=Hello%20World
Content-type: application/json
Accept: application/json
Authorization: Bearer {{token}}

> {%
  client.test("Login successful", function() {
    client.assert(response.status === 200, "Expected status 200");
    client.assert(response.body.token, "Token should be present");
  });
  
  // Save token for subsequent requests
  client.global.set("authToken", response.body.token);
%}`

	requestsFile, err := ParseFile("test.http", input)
	if err != nil {
		t.Fatalf("Parse error: %v", err)
	}

	if len(requestsFile.Requests) != 1 {
		t.Fatalf("Expected 1 request, got %d", len(requestsFile.Requests))
	}

	req := requestsFile.Requests[0]

	// Issue 1: Check that all headers are parsed (not just the first one)
	if len(req.Headers) != 3 {
		t.Errorf("Expected 3 headers, got %d", len(req.Headers))
	}

	// Verify each header
	expectedHeaders := map[string]string{
		"Content-type":  "application/json",
		"Accept":        "application/json",
		"Authorization": "Bearer {{token}}",
	}

	for _, header := range req.Headers {
		expectedValue, exists := expectedHeaders[header.Name]
		if !exists {
			t.Errorf("Unexpected header: %s", header.Name)
			continue
		}

		if header.Value != expectedValue {
			t.Errorf("Header %s: expected value %q, got %q", header.Name, expectedValue, header.Value)
		}
	}

	// Issue 2: Check that the Authorization header has correct spacing
	// (space between "Bearer" and "{{token}}")
	authHeader := ""
	for _, header := range req.Headers {
		if header.Name == "Authorization" {
			authHeader = header.Value
			break
		}
	}

	if authHeader != "Bearer {{token}}" {
		t.Errorf("Authorization header spacing issue: expected %q, got %q", "Bearer {{token}}", authHeader)
	}

	// Verify the variable was extracted
	if len(req.Headers[2].Variables) != 1 || req.Headers[2].Variables[0] != "token" {
		t.Error("Expected 'token' variable to be extracted from Authorization header")
	}

	// Verify response handler was parsed
	if req.ResponseHandler == nil {
		t.Error("Expected response handler to be parsed")
	}
}
