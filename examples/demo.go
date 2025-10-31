package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"postie/pkg/auth"
	"postie/pkg/client"
	"postie/pkg/middleware"
)

func main() {
	// Example 1: Basic API client usage
	fmt.Println("=== Basic API Client Example ===")
	basicExample()

	// Example 2: API client with authentication
	fmt.Println("\n=== Authentication Example ===")
	authExample()

	// Example 3: API client with middleware
	fmt.Println("\n=== Middleware Example ===")
	middlewareExample()

	// Example 4: POST request with JSON body
	fmt.Println("\n=== POST with JSON Example ===")
	postJSONExample()

	// Example 5: Error handling
	fmt.Println("\n=== Error Handling Example ===")
	errorHandlingExample()
}

func basicExample() {
	// Create a simple client
	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 10 * time.Second,
	})

	// Make a GET request
	resp, err := apiClient.GET("/posts/1").Execute()
	if err != nil {
		log.Printf("Error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Print response details
	fmt.Printf("Status: %s\n", resp.Status)
	fmt.Printf("Duration: %v\n", resp.Duration)

	// Get response as text
	text, err := resp.Text()
	if err != nil {
		log.Printf("Error reading response: %v", err)
		return
	}
	fmt.Printf("Response: %s\n", text[:100]+"...")
}

func authExample() {
	// Create client with Bearer token authentication
	bearerAuth := auth.NewBearerTokenAuth("your-api-token")

	// Demonstrate usage of bearerAuth
	fmt.Printf("Bearer auth configured: %+v\n", bearerAuth)

	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://api.example.com",
		Timeout: 15 * time.Second,
		Headers: map[string]string{
			"Accept": "application/json",
		},
	})

	// Make an authenticated request
	resp, err := apiClient.GET("/user/profile").
		Header("Authorization", "Bearer your-token").
		Execute()

	if err != nil {
		log.Printf("Auth example error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Authenticated request status: %s\n", resp.Status)

	// Example with API key in header
	apiKeyAuth := auth.NewAPIKeyAuth("X-API-Key", "your-api-key", "header")
	_ = apiKeyAuth // Demonstrate the API

	// Example with Basic auth
	basicAuth := auth.NewBasicAuth("username", "password")
	_ = basicAuth // Demonstrate the API
}

func middlewareExample() {
	// Create client with middleware
	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://httpbin.org",
		Timeout: 10 * time.Second,
		Middleware: []client.Middleware{
			middleware.LoggingMiddleware,
			middleware.UserAgentMiddleware("Postie/1.0"),
		},
	})

	// Make request with middleware
	resp, err := apiClient.GET("/get").
		Param("param1", "value1").
		Param("param2", "value2").
		Execute()

	if err != nil {
		log.Printf("Middleware example error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("Request with middleware completed: %s\n", resp.Status)
}

func postJSONExample() {
	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://jsonplaceholder.typicode.com",
		Timeout: 10 * time.Second,
	})

	// Create data to send
	postData := map[string]interface{}{
		"title":  "My New Post",
		"body":   "This is the content of my post",
		"userId": 1,
	}

	// Make POST request with JSON body
	resp, err := apiClient.POST("/posts").
		JSON(postData).
		Execute()

	if err != nil {
		log.Printf("POST JSON example error: %v", err)
		return
	}
	defer resp.Body.Close()

	fmt.Printf("POST request status: %s\n", resp.Status)

	// Parse JSON response
	var result map[string]interface{}
	if err := resp.JSON(&result); err != nil {
		log.Printf("Error parsing JSON response: %v", err)
		return
	}

	fmt.Printf("Created post ID: %v\n", result["id"])
}

func errorHandlingExample() {
	apiClient := client.NewClient(&client.Config{
		BaseURL: "https://httpbin.org",
		Timeout: 5 * time.Second,
	})

	// Make request that will return 404
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	resp, err := apiClient.GET("/status/404").
		Context(ctx).
		Execute()

	if err != nil {
		log.Printf("Request error: %v", err)
		return
	}
	defer resp.Body.Close()

	// Check response status
	if resp.IsError() {
		fmt.Printf("HTTP Error: %s\n", resp.Status)
		if resp.IsClientError() {
			fmt.Println("This is a client error (4xx)")
		}
		if resp.IsServerError() {
			fmt.Println("This is a server error (5xx)")
		}
	}

	if resp.IsSuccess() {
		fmt.Println("Request was successful")
	}

	fmt.Printf("Response size: %d bytes\n", resp.Size())
	fmt.Printf("Content type: %s\n", resp.ContentType())
}
