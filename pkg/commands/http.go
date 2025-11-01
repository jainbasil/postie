package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"postie/pkg/cli"
	"postie/pkg/environment"
	"postie/pkg/executor"
	"postie/pkg/httprequest"
)

// HTTPCommands returns the http command with subcommands for working with .http files
func HTTPCommands() *cli.Command {
	return &cli.Command{
		Name:        "http",
		Description: "Work with HTTP request files (.http)",
		Subcommands: map[string]*cli.Command{
			"run":   httpRunCommand(),
			"parse": httpParseCommand(),
			"list":  httpListCommand(),
		},
	}
}

func httpRunCommand() *cli.Command {
	return &cli.Command{
		Name:        "run",
		Description: "Execute HTTP requests from .http file",
		Action: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("HTTP request file required\nUsage: postie http run <file.http> [--env development] [--request name_or_number]")
			}

			var env, envFile, privateEnvFile, requestFilter string
			var verbose bool

			envFlag := &cli.StringFlag{Name: "env", ShortName: "e", Value: env, Usage: "Environment to use", Required: false}
			envFileFlag := &cli.StringFlag{Name: "env-file", Value: envFile, Usage: "Path to environment file", Required: false}
			privateEnvFileFlag := &cli.StringFlag{Name: "private-env-file", Value: privateEnvFile, Usage: "Path to private environment file", Required: false}
			requestFlag := &cli.StringFlag{Name: "request", ShortName: "r", Value: requestFilter, Usage: "Specific request name or number to run", Required: false}
			verboseFlag := &cli.BoolFlag{Name: "verbose", ShortName: "v", Value: verbose, Usage: "Verbose output"}

			_, err := cli.ParseFlags(args[1:], []*cli.StringFlag{envFlag, envFileFlag, privateEnvFileFlag, requestFlag}, []*cli.BoolFlag{verboseFlag})
			if err != nil {
				return err
			}

			env = envFlag.Value
			if env == "" {
				env = "development"
			}
			envFile = envFileFlag.Value
			if envFile == "" {
				envFile = "http-client.env.json"
			}
			privateEnvFile = privateEnvFileFlag.Value
			if privateEnvFile == "" {
				privateEnvFile = "http-client.private.env.json"
			}
			requestFilter = requestFlag.Value
			verbose = verboseFlag.Value

			return executeHttpFileRun(args[0], env, envFile, privateEnvFile, requestFilter, verbose)
		},
	}
}

func httpParseCommand() *cli.Command {
	return &cli.Command{
		Name:        "parse",
		Description: "Parse and validate HTTP request file",
		Action: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("HTTP request file required\nUsage: postie http parse <file.http> [--format summary]")
			}

			var format string
			var validate bool

			formatFlag := &cli.StringFlag{Name: "format", ShortName: "f", Value: format, Usage: "Output format (json, summary)", Required: false}
			validateFlag := &cli.BoolFlag{Name: "validate", Value: validate, Usage: "Perform validation"}

			_, err := cli.ParseFlags(args[1:], []*cli.StringFlag{formatFlag}, []*cli.BoolFlag{validateFlag})
			if err != nil {
				return err
			}

			format = formatFlag.Value
			if format == "" {
				format = "summary"
			}
			validate = validateFlag.Value

			return executeHttpFileParse(args[0], format, validate)
		},
	}
}

func httpListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List HTTP request files in directory",
		Action: func(args []string) error {
			var recursive bool
			recursiveFlag := &cli.BoolFlag{Name: "recursive", ShortName: "r", Value: recursive, Usage: "Search recursively"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{}, []*cli.BoolFlag{recursiveFlag})
			if err != nil {
				return err
			}

			dir := "."
			// Get the directory from non-flag arguments
			for _, arg := range args {
				if !strings.HasPrefix(arg, "-") {
					dir = arg
					break
				}
			}

			recursive = recursiveFlag.Value

			return executeHttpFileList(dir, recursive)
		},
	}
}

// Execute functions

func executeHttpFileRun(filePath string, envName string, envFile string, privateEnvFile string, requestName string, verbose bool) error {
	// Load environment files
	resolvedEnv, err := loadEnvironmentFiles(envName, envFile, privateEnvFile)
	if err != nil {
		return fmt.Errorf("failed to load environment: %w", err)
	}

	// Read HTTP file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read HTTP file: %w", err)
	}

	// Parse the HTTP file
	requestsFile, err := httprequest.ParseFile(filePath, string(content))
	if err != nil {
		return fmt.Errorf("failed to parse HTTP file: %w", err)
	}

	// Create executor
	exec := executor.NewExecutor(resolvedEnv, nil)
	formatter := executor.NewFormatter(verbose)

	// Execute requests from file
	results, err := exec.ExecuteFile(requestsFile, requestName)
	if err != nil {
		return fmt.Errorf("failed to execute requests: %w", err)
	}

	if len(results) == 0 {
		return fmt.Errorf("no requests executed")
	}

	// Display results
	for i, result := range results {
		fmt.Print(formatter.FormatResult(result, i+1))
	}

	// Display summary if multiple requests
	if len(results) > 1 {
		fmt.Print(formatter.FormatSummary(results))
	}

	return nil
}

// loadEnvironmentFiles loads and merges environment files
func loadEnvironmentFiles(envName string, envFile string, privateEnvFile string) (*environment.ResolvedEnvironment, error) {
	// Get working directory for loader
	workingDir := "."
	if abs, err := filepath.Abs("."); err == nil {
		workingDir = abs
	}

	loader := environment.NewLoader(workingDir)

	// Create environment config
	config := &environment.EnvironmentConfig{
		PublicFile:  envFile,
		PrivateFile: privateEnvFile,
		Environment: envName,
	}

	// Load both environment files
	publicEnv, privateEnv, err := loader.LoadEnvironments(config)
	if err != nil {
		// Check if it's just missing files
		if os.IsNotExist(err) {
			// Use empty environment if files don't exist
			emptyEnv := make(environment.EnvironmentFile)
			emptyPrivate := make(environment.EnvironmentFile)
			publicEnv = &emptyEnv
			privateEnv = &emptyPrivate
		} else {
			return nil, err
		}
	}

	// Resolve variables for the specified environment
	resolver := environment.NewResolver()
	resolvedEnv, err := resolver.Resolve(*publicEnv, *privateEnv, envName)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve environment variables: %w", err)
	}

	return resolvedEnv, nil
}

func executeHttpFileParse(httpFile, format string, validate bool) error {
	// Read the HTTP file content
	content, err := os.ReadFile(httpFile)
	if err != nil {
		return fmt.Errorf("failed to read HTTP file: %w", err)
	}

	// Parse HTTP request file
	requestsFile, err := httprequest.ParseFile(httpFile, string(content))
	if err != nil {
		return fmt.Errorf("failed to parse HTTP file: %w", err)
	}

	// Validate if requested
	if validate {
		validator := httprequest.NewValidator(true, "")
		if errors := validator.Validate(requestsFile); len(errors) > 0 {
			fmt.Printf("Validation errors:\n")
			for _, validationErr := range errors {
				fmt.Printf("  - %s\n", validationErr.Message)
			}
			return fmt.Errorf("validation failed")
		}
	}

	// Output results
	switch format {
	case "json":
		return outputJSON(requestsFile)
	case "summary":
		return outputSummary(requestsFile)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func executeHttpFileList(dir string, recursive bool) error {
	// Find .http files
	httpFiles, err := findHTTPFiles(dir, recursive)
	if err != nil {
		return fmt.Errorf("failed to find HTTP files: %w", err)
	}

	if len(httpFiles) == 0 {
		fmt.Printf("No .http files found in %s\n", dir)
		return nil
	}

	fmt.Printf("Found %d HTTP request file(s):\n", len(httpFiles))
	for _, file := range httpFiles {
		fmt.Printf("  %s\n", file)
	}

	return nil
}

// Helper functions

func filterRequests(requests []httprequest.Request, filter string) ([]httprequest.Request, error) {
	var filtered []httprequest.Request

	for i, request := range requests {
		// Check if filter matches request name
		if request.Name != "" && strings.Contains(strings.ToLower(request.Name), strings.ToLower(filter)) {
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

func expandRequestVariables(request *httprequest.Request, env *environment.ResolvedEnvironment) (*httprequest.Request, error) {
	// Create a copy of the request
	expanded := *request

	// Expand URL
	if request.URL != nil {
		resolver := environment.NewResolver()
		expanded.URL.Raw = resolver.ExpandString(request.URL.Raw, env)
		// TODO: Re-parse the expanded URL
	}

	// Expand headers
	for i, header := range request.Headers {
		resolver := environment.NewResolver()
		expanded.Headers[i].Value = resolver.ExpandString(header.Value, env)
	}

	// TODO: Expand body content

	return &expanded, nil
}

func printRequestDetails(request *httprequest.Request) {
	fmt.Printf("Method: %s\n", request.Method)
	if request.URL != nil {
		fmt.Printf("URL: %s\n", request.URL.Raw)
	}
	if len(request.Headers) > 0 {
		fmt.Printf("Headers:\n")
		for _, header := range request.Headers {
			fmt.Printf("  %s: %s\n", header.Name, header.Value)
		}
	}
	if request.Body != nil {
		fmt.Printf("Body: %s\n", request.Body.Type)
	}
}

func outputJSON(data interface{}) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

func outputSummary(requestsFile *httprequest.RequestsFile) error {
	fmt.Printf("Requests: %d\n\n", len(requestsFile.Requests))

	for i, request := range requestsFile.Requests {
		fmt.Printf("%d. %s %s", i+1, request.Method, request.URL.Raw)
		if request.Name != "" {
			fmt.Printf(" (%s)", request.Name)
		}
		fmt.Println()
	}

	return nil
}

func findHTTPFiles(dir string, recursive bool) ([]string, error) {
	var httpFiles []string

	if recursive {
		err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(strings.ToLower(path), ".http") {
				httpFiles = append(httpFiles, path)
			}
			return nil
		})
		return httpFiles, err
	} else {
		entries, err := os.ReadDir(dir)
		if err != nil {
			return nil, err
		}

		for _, entry := range entries {
			if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".http") {
				httpFiles = append(httpFiles, filepath.Join(dir, entry.Name()))
			}
		}
	}

	return httpFiles, nil
}
