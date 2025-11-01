package executor

import (
	"fmt"
	"time"

	"postie/pkg/client"
	"postie/pkg/environment"
	"postie/pkg/httprequest"
	"postie/pkg/responses"
	"postie/pkg/scripting"
)

// Executor executes HTTP requests with environment variable resolution
type Executor struct {
	client          *client.APIClient
	environment     *environment.ResolvedEnvironment
	verbose         bool
	globals         *scripting.GlobalStore     // Global variables for response handlers
	responseStorage *responses.Storage         // Response storage
	saveResponses   bool                       // Whether to save responses
}

// ExecutorConfig holds configuration for the executor
type ExecutorConfig struct {
	Timeout       time.Duration
	Verbose       bool
	SaveResponses bool                   // Enable response saving
	StorageConfig *responses.StorageConfig // Response storage configuration
}

// NewExecutor creates a new request executor
func NewExecutor(env *environment.ResolvedEnvironment, config *ExecutorConfig) *Executor {
	if config == nil {
		config = &ExecutorConfig{
			Timeout:       30 * time.Second,
			Verbose:       false,
			SaveResponses: false,
		}
	}

	var storage *responses.Storage
	if config.SaveResponses {
		storage = responses.NewStorage(config.StorageConfig)
	}

	return &Executor{
		client: client.NewClient(&client.Config{
			Timeout: config.Timeout,
		}),
		environment:     env,
		verbose:         config.Verbose,
		globals:         scripting.NewGlobalStore(),
		responseStorage: storage,
		saveResponses:   config.SaveResponses,
	}
}

// ExecuteRequest executes a single HTTP request
func (e *Executor) ExecuteRequest(request *httprequest.Request) (*ExecutionResult, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Expand variables in the request
	expandedRequest, err := e.expandRequestVariables(request)
	if err != nil {
		return nil, fmt.Errorf("failed to expand variables: %w", err)
	}

	// Build the HTTP request using the client
	req, err := e.buildClientRequest(expandedRequest)
	if err != nil {
		return nil, fmt.Errorf("failed to build request: %w", err)
	}

	// Execute the request
	startTime := time.Now()
	resp, err := req.Execute()
	duration := time.Since(startTime)

	if err != nil {
		return &ExecutionResult{
			Request:  expandedRequest,
			Error:    err,
			Duration: duration,
		}, err
	}

	// Build execution result
	result := &ExecutionResult{
		Request:    expandedRequest,
		Response:   resp,
		Duration:   duration,
		StatusCode: resp.Response.StatusCode,
		Status:     resp.Status,
	}

	// Execute response handler if present
	if expandedRequest.ResponseHandler != nil {
		envVars := make(map[string]interface{})
		if e.environment != nil {
			envVars = e.environment.Variables
		}

		scriptResult := scripting.ExecuteResponseHandler(
			expandedRequest.ResponseHandler,
			resp,
			expandedRequest,
			envVars,
			e.globals,
		)

		result.ScriptResult = scriptResult
	}

	// Save response if enabled
	if e.saveResponses && e.responseStorage != nil {
		storedResponse, err := responses.FromClientResponse(resp, expandedRequest, duration)
		if err == nil {
			filePath, err := e.responseStorage.Save(storedResponse)
			if err == nil {
				result.ResponseFilePath = filePath
			}
			// Don't fail the request if save fails, just skip
		}
	}

	return result, nil
}

// ExecuteFile executes all requests in an HTTP request file
func (e *Executor) ExecuteFile(requestsFile *httprequest.RequestsFile, filter string) ([]*ExecutionResult, error) {
	if requestsFile == nil {
		return nil, fmt.Errorf("requests file cannot be nil")
	}

	requestsToRun := requestsFile.Requests

	// Apply filter if specified
	if filter != "" {
		filtered, err := e.filterRequests(requestsFile.Requests, filter)
		if err != nil {
			return nil, fmt.Errorf("failed to filter requests: %w", err)
		}
		requestsToRun = filtered
	}

	// Execute each request
	results := make([]*ExecutionResult, 0, len(requestsToRun))
	for _, request := range requestsToRun {
		result, err := e.ExecuteRequest(&request)
		if err != nil && e.verbose {
			fmt.Printf("Error executing request: %v\n", err)
		}
		results = append(results, result)
	}

	return results, nil
}

// expandRequestVariables expands all variables in a request
func (e *Executor) expandRequestVariables(request *httprequest.Request) (*httprequest.Request, error) {
	// Create a copy of the request
	expanded := *request

	resolver := environment.NewResolver()

	// Create a combined environment with both env vars and globals
	combinedEnv := e.getCombinedEnvironment()

	// Expand URL
	if request.URL != nil {
		expanded.URL = &httprequest.URL{
			Raw:       resolver.ExpandString(request.URL.Raw, combinedEnv),
			Variables: request.URL.Variables,
		}
	}

	// Expand headers
	if len(request.Headers) > 0 {
		expanded.Headers = make([]httprequest.Header, len(request.Headers))
		for i, header := range request.Headers {
			expanded.Headers[i] = httprequest.Header{
				Name:  header.Name,
				Value: resolver.ExpandString(header.Value, combinedEnv),
			}
		}
	}

	// Expand body content
	if request.Body != nil {
		expanded.Body = &httprequest.RequestBody{
			Type:        request.Body.Type,
			ContentType: request.Body.ContentType,
			Content:     resolver.ExpandString(request.Body.Content, combinedEnv),
			Variables:   request.Body.Variables,
		}
	}

	return &expanded, nil
}

// getCombinedEnvironment merges environment variables and global variables
func (e *Executor) getCombinedEnvironment() *environment.ResolvedEnvironment {
	// Start with environment variables
	vars := make(map[string]interface{})
	if e.environment != nil {
		for k, v := range e.environment.Variables {
			vars[k] = v
		}
	}

	// Override/add global variables
	if e.globals != nil {
		globals := e.globals.GetAll()
		for k, v := range globals {
			vars[k] = v
		}
	}

	return &environment.ResolvedEnvironment{
		Name:      "combined",
		Variables: vars,
		Source:    make(map[string]string),
	}
}

// buildClientRequest converts a parsed request to a client request
func (e *Executor) buildClientRequest(request *httprequest.Request) (*client.Request, error) {
	if request.URL == nil {
		return nil, fmt.Errorf("request URL is required")
	}

	// Create request based on method
	var req *client.Request
	switch request.Method {
	case "GET":
		req = e.client.GET(request.URL.Raw)
	case "POST":
		req = e.client.POST(request.URL.Raw)
	case "PUT":
		req = e.client.PUT(request.URL.Raw)
	case "DELETE":
		req = e.client.DELETE(request.URL.Raw)
	case "PATCH":
		req = e.client.PATCH(request.URL.Raw)
	case "HEAD":
		req = e.client.HEAD(request.URL.Raw)
	case "OPTIONS":
		req = e.client.OPTIONS(request.URL.Raw)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", request.Method)
	}

	// Add headers
	for _, header := range request.Headers {
		req.Header(header.Name, header.Value)
	}

	// Add body if present
	if request.Body != nil && request.Body.Content != "" {
		// Determine content type
		contentType := request.Body.ContentType
		if contentType == "" {
			contentType = "text/plain"
		}

		// Set body based on content type
		if contentType == "application/json" || contentType == "text/json" {
			req.Text(request.Body.Content)
			req.Header("Content-Type", "application/json")
		} else {
			req.Text(request.Body.Content)
			req.Header("Content-Type", contentType)
		}
	}

	return req, nil
}

// filterRequests filters requests by name or number
func (e *Executor) filterRequests(requests []httprequest.Request, filter string) ([]httprequest.Request, error) {
	var filtered []httprequest.Request

	for i, request := range requests {
		// Check if filter matches request name
		if request.Name != "" && containsIgnoreCase(request.Name, filter) {
			filtered = append(filtered, request)
			continue
		}

		// Check if filter matches request number (1-based)
		if fmt.Sprintf("%d", i+1) == filter {
			filtered = append(filtered, request)
			continue
		}
	}

	if len(filtered) == 0 {
		return nil, fmt.Errorf("no requests match filter: %s", filter)
	}

	return filtered, nil
}

// containsIgnoreCase checks if s contains substr (case-insensitive)
func containsIgnoreCase(s, substr string) bool {
	s, substr = toLower(s), toLower(substr)
	return contains(s, substr)
}

func toLower(s string) string {
	result := make([]rune, len(s))
	for i, r := range s {
		if r >= 'A' && r <= 'Z' {
			result[i] = r + 32
		} else {
			result[i] = r
		}
	}
	return string(result)
}

func contains(s, substr string) bool {
	if len(substr) == 0 {
		return true
	}
	if len(substr) > len(s) {
		return false
	}
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
