package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"postie/pkg/cli"
	"postie/pkg/environment"
)

// EnvCommands returns the env command with subcommands for environment management
func EnvCommands() *cli.Command {
	return &cli.Command{
		Name:        "env",
		Description: "Manage environment files and variables",
		Subcommands: map[string]*cli.Command{
			"list": envListCommand(),
			"show": envShowCommand(),
		},
	}
}

func envListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List available environments",
		Action: func(args []string) error {
			var envFile, privateEnvFile string

			envFileFlag := &cli.StringFlag{Name: "env-file", Value: envFile, Usage: "Path to environment file", Required: false}
			privateEnvFileFlag := &cli.StringFlag{Name: "private-env-file", Value: privateEnvFile, Usage: "Path to private environment file", Required: false}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{envFileFlag, privateEnvFileFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			envFile = envFileFlag.Value
			if envFile == "" {
				envFile = "http-client.env.json"
			}
			privateEnvFile = privateEnvFileFlag.Value
			if privateEnvFile == "" {
				privateEnvFile = "http-client.private.env.json"
			}

			return executeEnvList(envFile, privateEnvFile)
		},
	}
}

func envShowCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show variables for a specific environment",
		Action: func(args []string) error {
			if len(args) == 0 {
				return fmt.Errorf("environment name required\nUsage: postie env show <environment> [--env-file file.json]")
			}

			var envFile, privateEnvFile string
			var showPrivate bool

			envFileFlag := &cli.StringFlag{Name: "env-file", Value: envFile, Usage: "Path to environment file", Required: false}
			privateEnvFileFlag := &cli.StringFlag{Name: "private-env-file", Value: privateEnvFile, Usage: "Path to private environment file", Required: false}
			showPrivateFlag := &cli.BoolFlag{Name: "show-private", Value: showPrivate, Usage: "Show private variables"}

			_, err := cli.ParseFlags(args[1:], []*cli.StringFlag{envFileFlag, privateEnvFileFlag}, []*cli.BoolFlag{showPrivateFlag})
			if err != nil {
				return err
			}

			envFile = envFileFlag.Value
			if envFile == "" {
				envFile = "http-client.env.json"
			}
			privateEnvFile = privateEnvFileFlag.Value
			if privateEnvFile == "" {
				privateEnvFile = "http-client.private.env.json"
			}
			showPrivate = showPrivateFlag.Value

			return executeEnvShow(args[0], envFile, privateEnvFile, showPrivate)
		},
	}
}

func executeEnvList(envFile string, privateEnvFile string) error {
	// Get working directory
	workingDir := "."
	if abs, err := filepath.Abs("."); err == nil {
		workingDir = abs
	}

	loader := environment.NewLoader(workingDir)

	// Create environment config
	config := &environment.EnvironmentConfig{
		PublicFile:  envFile,
		PrivateFile: privateEnvFile,
	}

	// Load both environment files
	publicEnv, privateEnv, err := loader.LoadEnvironments(config)
	if err != nil {
		// Check if files just don't exist
		if os.IsNotExist(err) {
			fmt.Println("No environment files found.")
			fmt.Printf("Create %s to define environments.\n", envFile)
			return nil
		}
		return fmt.Errorf("failed to load environments: %w", err)
	}

	// Collect all unique environment names
	envNames := make(map[string]bool)
	if publicEnv != nil {
		for name := range *publicEnv {
			envNames[name] = true
		}
	}
	if privateEnv != nil {
		for name := range *privateEnv {
			envNames[name] = true
		}
	}

	if len(envNames) == 0 {
		fmt.Println("No environments defined.")
		return nil
	}

	// Sort and display
	names := make([]string, 0, len(envNames))
	for name := range envNames {
		names = append(names, name)
	}
	sort.Strings(names)

	fmt.Println("Available environments:")
	for _, name := range names {
		// Count variables
		publicCount := 0
		privateCount := 0

		if publicEnv != nil {
			if env, exists := (*publicEnv)[name]; exists {
				publicCount = len(env)
			}
		}
		if privateEnv != nil {
			if env, exists := (*privateEnv)[name]; exists {
				privateCount = len(env)
			}
		}

		fmt.Printf("  %s", name)
		if publicCount > 0 || privateCount > 0 {
			fmt.Printf(" (%d public", publicCount)
			if privateCount > 0 {
				fmt.Printf(", %d private", privateCount)
			}
			fmt.Printf(" variables)")
		}
		fmt.Println()
	}

	return nil
}

func executeEnvShow(envName string, envFile string, privateEnvFile string, showPrivate bool) error {
	// Get working directory
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
		return fmt.Errorf("failed to load environments: %w", err)
	}

	// Get environment variables
	var publicVars, privateVars environment.Environment
	publicExists := false
	privateExists := false

	if publicEnv != nil {
		if env, exists := (*publicEnv)[envName]; exists {
			publicVars = env
			publicExists = true
		}
	}
	if privateEnv != nil {
		if env, exists := (*privateEnv)[envName]; exists {
			privateVars = env
			privateExists = true
		}
	}

	if !publicExists && !privateExists {
		return fmt.Errorf("environment '%s' not found", envName)
	}

	fmt.Printf("Environment: %s\n\n", envName)

	// Display public variables
	if publicExists && len(publicVars) > 0 {
		fmt.Println("Public variables:")
		displayVariables(publicVars)
		fmt.Println()
	}

	// Display private variables if requested
	if showPrivate && privateExists && len(privateVars) > 0 {
		fmt.Println("Private variables:")
		displayVariables(privateVars)
		fmt.Println()
	} else if privateExists && len(privateVars) > 0 {
		fmt.Printf("Private variables: %d (use --show-private to display)\n\n", len(privateVars))
	}

	// Resolve and show merged result
	resolver := environment.NewResolver()
	resolved, err := resolver.Resolve(*publicEnv, *privateEnv, envName)
	if err != nil {
		return fmt.Errorf("failed to resolve variables: %w", err)
	}

	fmt.Printf("Total resolved variables: %d\n", len(resolved.Variables))

	return nil
}

func displayVariables(vars environment.Environment) {
	// Sort keys
	keys := make([]string, 0, len(vars))
	for k := range vars {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// Display each variable
	for _, key := range keys {
		value := vars[key]

		// Format value based on type
		var valueStr string
		switch v := value.(type) {
		case string:
			valueStr = fmt.Sprintf("\"%s\"", v)
		case bool, int, int64, float64:
			valueStr = fmt.Sprintf("%v", v)
		default:
			// Try to JSON encode complex types
			if jsonBytes, err := json.Marshal(v); err == nil {
				valueStr = string(jsonBytes)
			} else {
				valueStr = fmt.Sprintf("%v", v)
			}
		}

		fmt.Printf("  %s = %s\n", key, valueStr)
	}
}
