package main

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"postie/pkg/auth"
	"postie/pkg/cli"
	"postie/pkg/client"
	"postie/pkg/commands"
	"postie/pkg/middleware"
)

func main() {
	// Create new CLI
	app := cli.NewCLI("postie", "1.0.0", "A powerful command-line API testing tool")

	// Add commands
	app.AddCommand(commands.HTTPCommands())
	app.AddCommand(commands.EnvCommands())
	app.AddCommand(commands.ContextCommands())
	app.AddCommand(demoCommand())

	// Run CLI
	if err := app.Run(os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// demoCommand returns the demo command
func demoCommand() *cli.Command {
	return &cli.Command{
		Name:        "demo",
		Description: "Run demonstration examples",
		Action: func(args []string) error {
			runDemo()
			return nil
		},
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

	// Demo 2: POST with custom headers
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
	fmt.Println("🎉 Demo completed! Try the new CLI commands:")
	fmt.Println("  postie http get --url https://httpbin.org/get")
	fmt.Println("  postie http post --url https://httpbin.org/post --body '{\"test\":\"data\"}'")
	fmt.Println("  postie collection create --name \"My API\" --file my-api.collection.json")
	fmt.Println("  postie request-group create --name \"Users\"")
	fmt.Println("  postie request create --name \"Get Users\" --method GET --url \"{{baseUrl}}/users\" --group \"users\"")
}

func printResponse(resp *client.Response) {
	separator := "=================================================="
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
