package commands

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"postie/pkg/cli"
	"postie/pkg/client"
	"postie/pkg/middleware"
)

// HTTPCommands returns the HTTP command with all subcommands
func HTTPCommands() *cli.Command {
	return &cli.Command{
		Name:        "http",
		Description: "Make direct HTTP requests",
		Subcommands: map[string]*cli.Command{
			"get":     httpGetCommand(),
			"post":    httpPostCommand(),
			"put":     httpPutCommand(),
			"delete":  httpDeleteCommand(),
			"patch":   httpPatchCommand(),
			"head":    httpHeadCommand(),
			"options": httpOptionsCommand(),
		},
	}
}

func httpGetCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Description: "Send a GET request",
		Action: func(args []string) error {
			var url string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value

			return executeHTTPRequest("GET", url, "")
		},
	}
}

func httpPostCommand() *cli.Command {
	return &cli.Command{
		Name:        "post",
		Description: "Send a POST request",
		Action: func(args []string) error {
			var url, body string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}
			bodyFlag := &cli.StringFlag{Name: "body", ShortName: "b", Value: body, Usage: "Request body"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag, bodyFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value
			body = bodyFlag.Value

			return executeHTTPRequest("POST", url, body)
		},
	}
}

func httpPutCommand() *cli.Command {
	return &cli.Command{
		Name:        "put",
		Description: "Send a PUT request",
		Action: func(args []string) error {
			var url, body string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}
			bodyFlag := &cli.StringFlag{Name: "body", ShortName: "b", Value: body, Usage: "Request body"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag, bodyFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value
			body = bodyFlag.Value

			return executeHTTPRequest("PUT", url, body)
		},
	}
}

func httpDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Send a DELETE request",
		Action: func(args []string) error {
			var url string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value

			return executeHTTPRequest("DELETE", url, "")
		},
	}
}

func httpPatchCommand() *cli.Command {
	return &cli.Command{
		Name:        "patch",
		Description: "Send a PATCH request",
		Action: func(args []string) error {
			var url, body string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}
			bodyFlag := &cli.StringFlag{Name: "body", ShortName: "b", Value: body, Usage: "Request body"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag, bodyFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value
			body = bodyFlag.Value

			return executeHTTPRequest("PATCH", url, body)
		},
	}
}

func httpHeadCommand() *cli.Command {
	return &cli.Command{
		Name:        "head",
		Description: "Send a HEAD request",
		Action: func(args []string) error {
			var url string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value

			return executeHTTPRequest("HEAD", url, "")
		},
	}
}

func httpOptionsCommand() *cli.Command {
	return &cli.Command{
		Name:        "options",
		Description: "Send an OPTIONS request",
		Action: func(args []string) error {
			var url string

			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{urlFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			url = urlFlag.Value

			return executeHTTPRequest("OPTIONS", url, "")
		},
	}
}

func executeHTTPRequest(method, url, body string) error {
	fmt.Printf("Sending %s request to: %s\n", method, url)

	apiClient := client.NewClient(&client.Config{
		Timeout: 30 * time.Second,
		Middleware: []client.Middleware{
			middleware.LoggingMiddleware,
		},
	})

	var req *client.Request
	switch method {
	case "GET":
		req = apiClient.GET(url)
	case "POST":
		req = apiClient.POST(url)
		if body != "" {
			req.Text(body)
		}
	case "PUT":
		req = apiClient.PUT(url)
		if body != "" {
			req.Text(body)
		}
	case "DELETE":
		req = apiClient.DELETE(url)
	case "PATCH":
		req = apiClient.PATCH(url)
		if body != "" {
			req.Text(body)
		}
	case "HEAD":
		req = apiClient.HEAD(url)
	case "OPTIONS":
		req = apiClient.OPTIONS(url)
	default:
		return fmt.Errorf("unsupported HTTP method: %s", method)
	}

	resp, err := req.Execute()
	if err != nil {
		return fmt.Errorf("error executing request: %w", err)
	}
	defer resp.Response.Body.Close()

	printHTTPResponse(resp)
	return nil
}

func printHTTPResponse(resp *client.Response) {
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
