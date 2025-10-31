package collection

import (
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"postie/pkg/client"
)

// Runner executes collection requests
type Runner struct {
	client      *client.APIClient
	collection  *Collection
	environment string
	variables   map[string]interface{}
}

// NewRunner creates a new collection runner
func NewRunner(collection *Collection, environment string) *Runner {
	// Create HTTP client
	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
	})

	return &Runner{
		client:      apiClient,
		collection:  collection,
		environment: environment,
		variables:   make(map[string]interface{}), // Will be resolved per request
	}
}

// RunRequest executes a single request from the collection
func (r *Runner) RunRequest(requestItem RequestItem) (*client.Response, error) {
	req := requestItem.Request

	// Get the parent folder for variable resolution
	var folderItem *Item
	if requestItem.ParentItem != nil {
		folderItem = requestItem.ParentItem
	}

	// Resolve variables for this specific request context
	variables := r.collection.ResolveVariables(r.environment, folderItem)

	// Get URL with variable substitution
	url := r.collection.GetRequestURL(req, variables)
	if url == "" {
		return nil, fmt.Errorf("request URL is empty")
	}

	// Create request based on method
	var clientReq *client.Request
	switch strings.ToUpper(req.Method) {
	case "GET":
		clientReq = r.client.GET(url)
	case "POST":
		clientReq = r.client.POST(url)
	case "PUT":
		clientReq = r.client.PUT(url)
	case "DELETE":
		clientReq = r.client.DELETE(url)
	case "PATCH":
		clientReq = r.client.PATCH(url)
	case "HEAD":
		clientReq = r.client.HEAD(url)
	case "OPTIONS":
		clientReq = r.client.OPTIONS(url)
	default:
		return nil, fmt.Errorf("unsupported HTTP method: %s", req.Method)
	}

	// Add headers
	for _, header := range req.Header {
		value := ReplaceVariables(header.Value, variables)
		clientReq.Header(header.Key, value)
	}

	// Add body if present
	if req.Body != nil {
		switch req.Body.Mode {
		case "raw":
			content := ReplaceVariables(req.Body.Raw, variables)
			if options, ok := req.Body.Options["raw"].(map[string]interface{}); ok {
				if language, ok := options["language"].(string); ok && language == "json" {
					// Set as JSON
					clientReq.Header("Content-Type", "application/json")
				}
			}
			clientReq.Text(content)
		}
	}

	// Apply authentication
	if err := r.applyAuthentication(clientReq, requestItem, variables); err != nil {
		return nil, fmt.Errorf("failed to apply authentication: %w", err)
	}

	// Execute request
	return clientReq.Execute()
}

// applyAuthentication applies the appropriate authentication to the request
func (r *Runner) applyAuthentication(clientReq *client.Request, requestItem RequestItem, variables map[string]interface{}) error {
	auth := r.collection.GetAuth(r.environment, nil, requestItem.Item)
	if auth == nil || auth.Type == "noauth" {
		return nil
	}

	switch auth.Type {
	case "bearer":
		token := r.getAuthValue(auth.Bearer, "token")
		if token != "" {
			token = ReplaceVariables(token, variables)
			clientReq.Header("Authorization", fmt.Sprintf("Bearer %s", token))
		}

	case "apikey":
		key := r.getAuthValue(auth.APIKey, "key")
		value := r.getAuthValue(auth.APIKey, "value")
		in := r.getAuthValue(auth.APIKey, "in")

		if key != "" && value != "" {
			value = ReplaceVariables(value, variables)
			if in == "header" {
				clientReq.Header(key, value)
			} else if in == "query" {
				clientReq.Param(key, value)
			}
		}

	case "basic":
		username := r.getAuthValue(auth.Basic, "username")
		password := r.getAuthValue(auth.Basic, "password")

		if username != "" && password != "" {
			username = ReplaceVariables(username, variables)
			password = ReplaceVariables(password, variables)

			// Create basic auth header
			credentials := base64.StdEncoding.EncodeToString([]byte(username + ":" + password))
			clientReq.Header("Authorization", "Basic "+credentials)
		}
	}

	return nil
}

// getAuthValue extracts a value from auth parameters
func (r *Runner) getAuthValue(params []AuthParam, key string) string {
	for _, param := range params {
		if param.Key == key {
			return param.Value
		}
	}
	return ""
}

// RunAll executes all requests in the collection
func (r *Runner) RunAll() error {
	requests := r.collection.FindAllRequests()

	fmt.Printf("Running collection: %s\n", r.collection.Collection.Info.Name)
	fmt.Printf("Environment: %s\n", r.environment)
	fmt.Printf("Found %d requests\n\n", len(requests))

	for i, requestItem := range requests {
		fmt.Printf("[%d/%d] %s\n", i+1, len(requests), requestItem.Path)

		start := time.Now()
		resp, err := r.RunRequest(requestItem)
		duration := time.Since(start)

		if err != nil {
			fmt.Printf("  ❌ Error: %v\n", err)
			continue
		}

		defer resp.Response.Body.Close()

		// Print results
		status := "✅"
		if resp.IsError() {
			status = "❌"
		}

		fmt.Printf("  %s %s (%v)\n", status, resp.Status, duration)

		// Print response size
		if resp.Size() > 0 {
			fmt.Printf("  📦 %d bytes\n", resp.Size())
		}

		fmt.Println()
	}

	return nil
}

// RunByName executes a specific request by name
func (r *Runner) RunByName(name string) error {
	return r.RunByNameOrID(name, false)
}

// RunByID executes a specific request by ID
func (r *Runner) RunByID(id string) error {
	return r.RunByNameOrID(id, true)
}

// RunByNameOrID executes a specific request by name or ID
func (r *Runner) RunByNameOrID(identifier string, isID bool) error {
	requests := r.collection.FindAllRequests()

	for _, requestItem := range requests {
		var match bool
		if isID {
			match = requestItem.Item.ID == identifier
		} else {
			match = requestItem.Name == identifier
		}

		if match {
			fmt.Printf("Running request: %s\n", requestItem.Path)

			resp, err := r.RunRequest(requestItem)
			if err != nil {
				return fmt.Errorf("failed to execute request: %w", err)
			}
			defer resp.Response.Body.Close()

			// Print detailed response
			r.printDetailedResponse(resp)
			return nil
		}
	}

	if isID {
		return fmt.Errorf("request with ID '%s' not found in collection", identifier)
	}
	return fmt.Errorf("request '%s' not found in collection", identifier)
}

// printDetailedResponse prints a detailed response similar to the CLI
func (r *Runner) printDetailedResponse(resp *client.Response) {
	separator := strings.Repeat("=", 50)
	fmt.Println(separator)
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Duration: %v\n", resp.Duration)
	fmt.Printf("Size: %d bytes\n", resp.Size())
	fmt.Printf("Content-Type: %s\n", resp.ContentType())
	fmt.Println(separator)

	if resp.IsSuccess() {
		fmt.Println("✅ Request successful")
	} else if resp.IsError() {
		fmt.Printf("❌ Request failed: %s\n", resp.Status)
	}

	// Print response body
	text, err := resp.Text()
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	fmt.Println("\nResponse Body:")
	if len(text) > 1000 {
		fmt.Printf("%s...\n[Response truncated - %d total characters]\n", text[:1000], len(text))
	} else {
		fmt.Println(text)
	}
}

// ListRequests lists all requests in the collection
func (r *Runner) ListRequests() {
	requests := r.collection.FindAllRequests()

	fmt.Printf("Collection: %s\n", r.collection.Collection.Info.Name)
	fmt.Printf("Environment: %s\n", r.environment)
	fmt.Printf("Requests (%d):\n\n", len(requests))

	for i, requestItem := range requests {
		// Use the parent folder for variable resolution
		var folderItem *Item
		if requestItem.ParentItem != nil {
			folderItem = requestItem.ParentItem
		}
		variables := r.collection.ResolveVariables(r.environment, folderItem)
		fmt.Printf("%d. %s\n", i+1, requestItem.Path)
		fmt.Printf("   ID: %s\n", requestItem.Item.ID)
		fmt.Printf("   %s %s\n", requestItem.Request.Method,
			r.collection.GetRequestURL(requestItem.Request, variables))
		fmt.Println()
	}
}
