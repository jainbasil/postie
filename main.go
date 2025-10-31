package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"postie/pkg/auth"
	"postie/pkg/client"
	"postie/pkg/collection"
	"postie/pkg/middleware"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		return
	}

	command := os.Args[1]

	switch command {
	case "get":
		handleGET()
	case "post":
		handlePOST()
	case "put":
		handlePUT()
	case "delete":
		handleDELETE()
	case "demo":
		runDemo()
	case "collection", "run":
		handleCollection()
	case "list":
		handleListCollection()
	case "env":
		handleEnvironment()
	case "create":
		handleCreate()
	case "update":
		handleUpdate()
	case "remove":
		handleRemove()
	case "help", "--help", "-h":
		printUsage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
	}
}

func printUsage() {
	fmt.Println("Postie - A powerful command-line API testing tool")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  postie <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  get <url>           Send a GET request")
	fmt.Println("  post <url>          Send a POST request")
	fmt.Println("  put <url>           Send a PUT request")
	fmt.Println("  delete <url>        Send a DELETE request")
	fmt.Println("  run <collection>    Run a collection file")
	fmt.Println("  list <collection>   List requests in a collection")
	fmt.Println("  env <collection>    Show environments in a collection")
	fmt.Println("  demo                Run demonstration examples")
	fmt.Println()
	fmt.Println("Collection Management:")
	fmt.Println("  create collection <name>        Create a new collection")
	fmt.Println("  create apigroup <file> <name>   Create a new API group")
	fmt.Println("  create api <file> <group-id> <name> <method> <url>  Create a new API")
	fmt.Println("  update collection <file>        Update collection metadata")
	fmt.Println("  update apigroup <file> <id>     Update API group")
	fmt.Println("  update api <file> <id>          Update API")
	fmt.Println("  remove apigroup <file> <id>     Remove API group")
	fmt.Println("  remove api <file> <id>          Remove API")
	fmt.Println("  help                Show this help message")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  postie get https://api.github.com/users/octocat")
	fmt.Println("  postie post https://httpbin.org/post")
	fmt.Println("  postie run collections/jsonplaceholder.collection.json")
	fmt.Println("  postie run collections/jsonplaceholder.collection.json --env \"Local Development\"")
	fmt.Println("  postie list collections/jsonplaceholder.collection.json")
	fmt.Println("  postie create collection \"My API\" --file my-api.collection.json")
	fmt.Println("  postie create apigroup my-api.collection.json \"Users\"")
	fmt.Println("  postie demo")
}

func handleGET() {
	if len(os.Args) < 3 {
		fmt.Println("Error: URL required for GET request")
		fmt.Println("Usage: postie get <url>")
		return
	}

	url := os.Args[2]
	fmt.Printf("Sending GET request to: %s\n", url)

	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
		Middleware: []client.Middleware{
			middleware.LoggingMiddleware,
		},
	})

	resp, err := apiClient.GET(url).Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Response.Body.Close()

	printResponse(resp)
}

func handlePOST() {
	if len(os.Args) < 3 {
		fmt.Println("Error: URL required for POST request")
		fmt.Println("Usage: postie post <url>")
		return
	}

	url := os.Args[2]
	fmt.Printf("Sending POST request to: %s\n", url)

	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
		Middleware: []client.Middleware{
			middleware.LoggingMiddleware,
		},
	})

	// Simple JSON payload for demo
	data := map[string]interface{}{
		"title":     "Test Post",
		"body":      "This is a test post from Postie",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	resp, err := apiClient.POST(url).JSON(data).Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Response.Body.Close()

	printResponse(resp)
}

func handlePUT() {
	if len(os.Args) < 3 {
		fmt.Println("Error: URL required for PUT request")
		fmt.Println("Usage: postie put <url>")
		return
	}

	url := os.Args[2]
	fmt.Printf("Sending PUT request to: %s\n", url)

	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
	})

	data := map[string]interface{}{
		"title":     "Updated Post",
		"body":      "This is an updated post from Postie",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	resp, err := apiClient.PUT(url).JSON(data).Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Response.Body.Close()

	printResponse(resp)
}

func handleDELETE() {
	if len(os.Args) < 3 {
		fmt.Println("Error: URL required for DELETE request")
		fmt.Println("Usage: postie delete <url>")
		return
	}

	url := os.Args[2]
	fmt.Printf("Sending DELETE request to: %s\n", url)

	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
	})

	resp, err := apiClient.DELETE(url).Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Response.Body.Close()

	printResponse(resp)
}

func printResponse(resp *client.Response) {
	separator := strings.Repeat("=", 50)
	fmt.Println("\n" + separator)
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

	// Try to format JSON response
	text, err := resp.Text()
	if err != nil {
		fmt.Printf("Error reading response: %v\n", err)
		return
	}

	// Try to pretty-print JSON
	var jsonData interface{}
	if err := json.Unmarshal([]byte(text), &jsonData); err == nil {
		prettyJSON, _ := json.MarshalIndent(jsonData, "", "  ")
		fmt.Println("\nResponse Body:")
		fmt.Println(string(prettyJSON))
	} else {
		// Not JSON, print as text
		fmt.Println("\nResponse Body:")
		if len(text) > 1000 {
			fmt.Printf("%s...\n[Response truncated - %d total characters]\n", text[:1000], len(text))
		} else {
			fmt.Println(text)
		}
	}
}

func runDemo() {
	fmt.Println("🚀 Running Postie Demo...")
	fmt.Println()

	// Demo 1: Basic GET request
	fmt.Println("Demo 1: Basic GET request to JSONPlaceholder")
	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 10 * time.Second,
		Middleware: []client.Middleware{
			middleware.LoggingMiddleware,
		},
	})

	resp, err := apiClient.GET("/posts/1").Execute()
	if err != nil {
		fmt.Printf("Demo 1 failed: %v\n", err)
	} else {
		defer resp.Response.Body.Close()
		fmt.Printf("✅ Status: %s, Duration: %v\n", resp.Status, resp.Duration)
	}

	fmt.Println()

	// Demo 2: POST with authentication simulation
	fmt.Println("Demo 2: POST request with custom headers")
	resp2, err := apiClient.POST("/posts").
		Header("X-API-Key", "demo-key").
		JSON(map[string]interface{}{
			"title":  "Demo Post",
			"body":   "This is a demo post",
			"userId": 1,
		}).Execute()

	if err != nil {
		fmt.Printf("Demo 2 failed: %v\n", err)
	} else {
		defer resp2.Response.Body.Close()
		fmt.Printf("✅ Status: %s, Duration: %v\n", resp2.Status, resp2.Duration)
	}

	fmt.Println()

	// Demo 3: Authentication examples
	fmt.Println("Demo 3: Authentication examples")

	// API Key auth
	apiKey := auth.NewAPIKeyAuth("X-API-Key", "your-api-key", "header")
	fmt.Printf("✅ API Key Auth configured: %+v\n", apiKey)

	// Bearer token auth
	bearer := auth.NewBearerTokenAuth("your-bearer-token")
	fmt.Printf("✅ Bearer Token Auth configured: %+v\n", bearer)

	// Basic auth
	basic := auth.NewBasicAuth("username", "password")
	fmt.Printf("✅ Basic Auth configured: %+v\n", basic)

	fmt.Println()
	fmt.Println("🎉 Demo completed! Try the CLI commands:")
	fmt.Println("  postie get https://httpbin.org/get")
	fmt.Println("  postie post https://httpbin.org/post")
	fmt.Println("  postie run collections/jsonplaceholder.collection.json")
	fmt.Println("  postie run collections/jsonplaceholder.collection.json --request \"Get All Posts\"")
	fmt.Println("  postie run collections/jsonplaceholder.collection.json --id get-all-posts")
}

func handleCollection() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Collection file required")
		fmt.Println("Usage: postie run <collection.json> [--env <environment>] [--request <request-name>] [--id <request-id>]")
		return
	}

	collectionFile := os.Args[2]
	environment := ""
	requestName := ""
	requestID := ""

	// Parse additional arguments
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--env", "-e":
			if i+1 < len(os.Args) {
				environment = os.Args[i+1]
				i++
			}
		case "--request", "-r":
			if i+1 < len(os.Args) {
				requestName = os.Args[i+1]
				i++
			}
		case "--id":
			if i+1 < len(os.Args) {
				requestID = os.Args[i+1]
				i++
			}
		}
	}

	// Load collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Use default environment if none specified
	if environment == "" {
		if defaultEnv := coll.GetDefaultEnvironment(); defaultEnv != nil {
			environment = defaultEnv.Name
		}
	}

	// Validate environment
	if environment != "" {
		if _, err := coll.GetEnvironment(environment); err != nil {
			fmt.Printf("Error: %v\n", err)
			fmt.Printf("Available environments: %s\n", strings.Join(coll.GetEnvironmentNames(), ", "))
			return
		}
	}

	// Create runner
	runner := collection.NewRunner(coll, environment)

	// Run specific request or all requests
	if requestID != "" {
		if err := runner.RunByID(requestID); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else if requestName != "" {
		if err := runner.RunByName(requestName); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	} else {
		if err := runner.RunAll(); err != nil {
			fmt.Printf("Error: %v\n", err)
		}
	}
}

func handleListCollection() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Collection file required")
		fmt.Println("Usage: postie list <collection.json> [--env <environment>]")
		return
	}

	collectionFile := os.Args[2]
	environment := ""

	// Parse additional arguments
	for i := 3; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--env", "-e":
			if i+1 < len(os.Args) {
				environment = os.Args[i+1]
				i++
			}
		}
	}

	// Load collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Use default environment if none specified
	if environment == "" {
		if defaultEnv := coll.GetDefaultEnvironment(); defaultEnv != nil {
			environment = defaultEnv.Name
		}
	}

	// Create runner and list requests
	runner := collection.NewRunner(coll, environment)
	runner.ListRequests()
}

func handleEnvironment() {
	if len(os.Args) < 3 {
		fmt.Println("Error: Collection file required")
		fmt.Println("Usage: postie env <collection.json>")
		return
	}

	collectionFile := os.Args[2]

	// Load collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Show environments
	fmt.Printf("Collection: %s\n", coll.Collection.Info.Name)
	fmt.Printf("Environments (%d):\n\n", len(coll.Collection.Environment))

	for i, env := range coll.Collection.Environment {
		defaultMarker := ""
		if i == 0 {
			defaultMarker = " (default)"
		}

		fmt.Printf("%d. %s%s\n", i+1, env.Name, defaultMarker)
		if env.Description != "" {
			fmt.Printf("   %s\n", env.Description)
		}

		fmt.Printf("   Variables: %d\n", len(env.Values))
		if env.Auth != nil {
			fmt.Printf("   Authentication: %s\n", env.Auth.Type)
		}
		fmt.Println()
	}
}

// handleCreate handles creating collections, API groups, and APIs
func handleCreate() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  postie create collection <name> [--file <filename>]")
		fmt.Println("  postie create apigroup <collection-file> <group-name> [--id <group-id>]")
		fmt.Println("  postie create api <collection-file> <group-id> <api-name> <method> <url> [--id <api-id>]")
		return
	}

	resourceType := os.Args[2]

	switch resourceType {
	case "collection":
		handleCreateCollection()
	case "apigroup", "group":
		handleCreateApiGroup()
	case "api":
		handleCreateApi()
	default:
		fmt.Printf("Unknown resource type: %s\n", resourceType)
		fmt.Println("Supported types: collection, apigroup, api")
	}
}

// handleUpdate handles updating collections, API groups, and APIs
func handleUpdate() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  postie update collection <collection-file> [--name <new-name>] [--description <new-desc>]")
		fmt.Println("  postie update apigroup <collection-file> <group-id> [--name <new-name>] [--description <new-desc>]")
		fmt.Println("  postie update api <collection-file> <api-id> [--name <new-name>] [--method <new-method>] [--url <new-url>]")
		return
	}

	resourceType := os.Args[2]

	switch resourceType {
	case "collection":
		handleUpdateCollection()
	case "apigroup", "group":
		handleUpdateApiGroup()
	case "api":
		handleUpdateApi()
	default:
		fmt.Printf("Unknown resource type: %s\n", resourceType)
		fmt.Println("Supported types: collection, apigroup, api")
	}
}

// handleRemove handles removing API groups and APIs (collections are just files)
func handleRemove() {
	if len(os.Args) < 3 {
		fmt.Println("Usage:")
		fmt.Println("  postie remove apigroup <collection-file> <group-id>")
		fmt.Println("  postie remove api <collection-file> <api-id>")
		fmt.Println("Note: To remove a collection, just delete the JSON file")
		return
	}

	resourceType := os.Args[2]

	switch resourceType {
	case "apigroup", "group":
		handleRemoveApiGroup()
	case "api":
		handleRemoveApi()
	default:
		fmt.Printf("Unknown resource type: %s\n", resourceType)
		fmt.Println("Supported types: apigroup, api")
	}
}

// Collection CRUD operations
func handleCreateCollection() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: postie create collection <name> [--file <filename>]")
		return
	}

	collectionName := os.Args[3]
	filename := ""

	// Parse additional arguments
	for i := 4; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--file", "-f":
			if i+1 < len(os.Args) {
				filename = os.Args[i+1]
				i++
			}
		}
	}

	if filename == "" {
		// Generate filename from collection name
		filename = strings.ToLower(strings.ReplaceAll(collectionName, " ", "-")) + ".collection.json"
	}

	// Create new collection
	newCollection := &collection.Collection{
		Collection: collection.CollectionInfo{
			Info: collection.Info{
				Name:        collectionName,
				Description: fmt.Sprintf("API collection for %s", collectionName),
				Version:     "1.0.0",
				Schema:      "https://postie.dev/collection/v1.0.0/collection.json",
			},
			Variable:    []collection.Variable{},
			Environment: []collection.Environment{},
			ApiGroup:    []collection.ApiGroup{},
		},
	}

	// Save to file
	data, err := json.MarshalIndent(newCollection, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ Collection '%s' created successfully\n", collectionName)
	fmt.Printf("📁 File: %s\n", filename)
}

func handleUpdateCollection() {
	if len(os.Args) < 4 {
		fmt.Println("Usage: postie update collection <collection-file> [--name <new-name>] [--description <new-desc>]")
		return
	}

	collectionFile := os.Args[3]
	newName := ""
	newDescription := ""

	// Parse additional arguments
	for i := 4; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--name", "-n":
			if i+1 < len(os.Args) {
				newName = os.Args[i+1]
				i++
			}
		case "--description", "--desc", "-d":
			if i+1 < len(os.Args) {
				newDescription = os.Args[i+1]
				i++
			}
		}
	}

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Update fields
	if newName != "" {
		coll.Collection.Info.Name = newName
	}
	if newDescription != "" {
		coll.Collection.Info.Description = newDescription
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ Collection updated successfully\n")
	fmt.Printf("📁 File: %s\n", collectionFile)
	fmt.Printf("📝 Name: %s\n", coll.Collection.Info.Name)
}

// API Group CRUD operations
func handleCreateApiGroup() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: postie create apigroup <collection-file> <group-name> [--id <group-id>]")
		return
	}

	collectionFile := os.Args[3]
	groupName := os.Args[4]
	groupID := ""

	// Parse additional arguments
	for i := 5; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--id":
			if i+1 < len(os.Args) {
				groupID = os.Args[i+1]
				i++
			}
		}
	}

	if groupID == "" {
		groupID = collection.GenerateSlug(groupName)
	}

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Create new API group
	newGroup := collection.ApiGroup{
		ID:          groupID,
		Name:        groupName,
		Description: fmt.Sprintf("API group for %s operations", groupName),
		Apis:        []collection.Item{},
		Environment: []collection.Environment{},
	}

	// Add to collection
	coll.Collection.ApiGroup = append(coll.Collection.ApiGroup, newGroup)

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API Group '%s' created successfully\n", groupName)
	fmt.Printf("🆔 ID: %s\n", groupID)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}

func handleUpdateApiGroup() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: postie update apigroup <collection-file> <group-id> [--name <new-name>] [--description <new-desc>]")
		return
	}

	collectionFile := os.Args[3]
	groupID := os.Args[4]
	newName := ""
	newDescription := ""

	// Parse additional arguments
	for i := 5; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--name", "-n":
			if i+1 < len(os.Args) {
				newName = os.Args[i+1]
				i++
			}
		case "--description", "--desc", "-d":
			if i+1 < len(os.Args) {
				newDescription = os.Args[i+1]
				i++
			}
		}
	}

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Find and update API group
	found := false
	for i := range coll.Collection.ApiGroup {
		if coll.Collection.ApiGroup[i].ID == groupID {
			if newName != "" {
				coll.Collection.ApiGroup[i].Name = newName
			}
			if newDescription != "" {
				coll.Collection.ApiGroup[i].Description = newDescription
			}
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: API Group with ID '%s' not found\n", groupID)
		return
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API Group updated successfully\n")
	fmt.Printf("🆔 ID: %s\n", groupID)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}

func handleRemoveApiGroup() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: postie remove apigroup <collection-file> <group-id>")
		return
	}

	collectionFile := os.Args[3]
	groupID := os.Args[4]

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Find and remove API group
	found := false
	for i := range coll.Collection.ApiGroup {
		if coll.Collection.ApiGroup[i].ID == groupID {
			// Remove the group
			coll.Collection.ApiGroup = append(coll.Collection.ApiGroup[:i], coll.Collection.ApiGroup[i+1:]...)
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: API Group with ID '%s' not found\n", groupID)
		return
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API Group removed successfully\n")
	fmt.Printf("🆔 ID: %s\n", groupID)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}

// API CRUD operations
func handleCreateApi() {
	if len(os.Args) < 8 {
		fmt.Println("Usage: postie create api <collection-file> <group-id> <api-name> <method> <url> [--id <api-id>]")
		return
	}

	collectionFile := os.Args[3]
	groupID := os.Args[4]
	apiName := os.Args[5]
	method := strings.ToUpper(os.Args[6])
	apiURL := os.Args[7]
	apiID := ""

	// Parse additional arguments
	for i := 8; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--id":
			if i+1 < len(os.Args) {
				apiID = os.Args[i+1]
				i++
			}
		}
	}

	if apiID == "" {
		apiID = collection.GenerateSlug(apiName)
	}

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Find the API group
	found := false
	for i := range coll.Collection.ApiGroup {
		if coll.Collection.ApiGroup[i].ID == groupID {
			// Create new API
			newAPI := collection.Item{
				ID:          apiID,
				Name:        apiName,
				Description: fmt.Sprintf("%s operation for %s", method, apiName),
				Request: &collection.Request{
					Method: method,
					Header: []collection.Header{},
					URL:    apiURL,
				},
			}

			// Add to API group
			coll.Collection.ApiGroup[i].Apis = append(coll.Collection.ApiGroup[i].Apis, newAPI)
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("Error: API Group with ID '%s' not found\n", groupID)
		return
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API '%s' created successfully\n", apiName)
	fmt.Printf("🆔 ID: %s\n", apiID)
	fmt.Printf("🔗 %s %s\n", method, apiURL)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}

func handleUpdateApi() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: postie update api <collection-file> <api-id> [--name <new-name>] [--method <new-method>] [--url <new-url>]")
		return
	}

	collectionFile := os.Args[3]
	apiID := os.Args[4]
	newName := ""
	newMethod := ""
	newURL := ""

	// Parse additional arguments
	for i := 5; i < len(os.Args); i++ {
		switch os.Args[i] {
		case "--name", "-n":
			if i+1 < len(os.Args) {
				newName = os.Args[i+1]
				i++
			}
		case "--method", "-m":
			if i+1 < len(os.Args) {
				newMethod = strings.ToUpper(os.Args[i+1])
				i++
			}
		case "--url", "-u":
			if i+1 < len(os.Args) {
				newURL = os.Args[i+1]
				i++
			}
		}
	}

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Find and update API
	found := false
	for gi := range coll.Collection.ApiGroup {
		for ai := range coll.Collection.ApiGroup[gi].Apis {
			if coll.Collection.ApiGroup[gi].Apis[ai].ID == apiID {
				if newName != "" {
					coll.Collection.ApiGroup[gi].Apis[ai].Name = newName
				}
				if newMethod != "" && coll.Collection.ApiGroup[gi].Apis[ai].Request != nil {
					coll.Collection.ApiGroup[gi].Apis[ai].Request.Method = newMethod
				}
				if newURL != "" && coll.Collection.ApiGroup[gi].Apis[ai].Request != nil {
					coll.Collection.ApiGroup[gi].Apis[ai].Request.URL = newURL
				}
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Printf("Error: API with ID '%s' not found\n", apiID)
		return
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API updated successfully\n")
	fmt.Printf("🆔 ID: %s\n", apiID)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}

func handleRemoveApi() {
	if len(os.Args) < 5 {
		fmt.Println("Usage: postie remove api <collection-file> <api-id>")
		return
	}

	collectionFile := os.Args[3]
	apiID := os.Args[4]

	// Load existing collection
	coll, err := collection.LoadCollection(collectionFile)
	if err != nil {
		fmt.Printf("Error loading collection: %v\n", err)
		return
	}

	// Find and remove API
	found := false
	for gi := range coll.Collection.ApiGroup {
		for ai := range coll.Collection.ApiGroup[gi].Apis {
			if coll.Collection.ApiGroup[gi].Apis[ai].ID == apiID {
				// Remove the API
				coll.Collection.ApiGroup[gi].Apis = append(
					coll.Collection.ApiGroup[gi].Apis[:ai],
					coll.Collection.ApiGroup[gi].Apis[ai+1:]...,
				)
				found = true
				break
			}
		}
		if found {
			break
		}
	}

	if !found {
		fmt.Printf("Error: API with ID '%s' not found\n", apiID)
		return
	}

	// Save back to file
	data, err := json.MarshalIndent(coll, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling collection: %v\n", err)
		return
	}

	err = os.WriteFile(collectionFile, data, 0644)
	if err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}

	fmt.Printf("✅ API removed successfully\n")
	fmt.Printf("🆔 ID: %s\n", apiID)
	fmt.Printf("📁 Collection: %s\n", collectionFile)
}
