package commands

import (
	"flag"
	"fmt"
	"path/filepath"

	"postie/pkg/cli"
	"postie/pkg/context"
)

// ContextCommands returns the context command with subcommands
func ContextCommands() *cli.Command {
	subcommands := make(map[string]*cli.Command)
	subcommands["set"] = contextSetCommand()
	subcommands["show"] = contextShowCommand()
	subcommands["clear"] = contextClearCommand()

	return &cli.Command{
		Name:        "context",
		Description: "Manage context settings for the current directory",
		Subcommands: subcommands,
	}
}

func contextSetCommand() *cli.Command {
	return &cli.Command{
		Name:        "set",
		Description: "Set context values for the current directory",
		Action:      executeContextSet,
	}
}

func contextShowCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show current context settings",
		Action:      executeContextShow,
	}
}

func contextClearCommand() *cli.Command {
	return &cli.Command{
		Name:        "clear",
		Description: "Clear context settings for the current directory",
		Action:      executeContextClear,
	}
}

func executeContextSet(args []string) error {
	fs := flag.NewFlagSet("context set", flag.ExitOnError)
	httpFile := fs.String("http-file", "", "Path to HTTP request file")
	env := fs.String("env", "", "Environment name")
	envFile := fs.String("env-file", "", "Path to environment file")
	privateEnvFile := fs.String("private-env-file", "", "Path to private environment file")
	saveResponses := fs.Bool("save-responses", false, "Save responses to files")
	responsesDir := fs.String("responses-dir", "", "Directory to save responses")

	if err := fs.Parse(args); err != nil {
		return err
	}

	mgr := context.NewManager()

	// Load existing context
	ctx, err := mgr.Load()
	if err != nil {
		return err
	}

	// Update context with provided values
	updated := false
	if *httpFile != "" {
		// Convert to absolute path if relative
		absPath, err := filepath.Abs(*httpFile)
		if err != nil {
			return fmt.Errorf("failed to resolve http file path: %w", err)
		}
		ctx.HTTPFile = absPath
		updated = true
	}
	if *env != "" {
		ctx.Environment = *env
		updated = true
	}
	if *envFile != "" {
		absPath, err := filepath.Abs(*envFile)
		if err != nil {
			return fmt.Errorf("failed to resolve env file path: %w", err)
		}
		ctx.EnvFile = absPath
		updated = true
	}
	if *privateEnvFile != "" {
		absPath, err := filepath.Abs(*privateEnvFile)
		if err != nil {
			return fmt.Errorf("failed to resolve private env file path: %w", err)
		}
		ctx.PrivateEnvFile = absPath
		updated = true
	}
	if *saveResponses {
		ctx.SaveResponses = true
		updated = true
	}
	if *responsesDir != "" {
		absPath, err := filepath.Abs(*responsesDir)
		if err != nil {
			return fmt.Errorf("failed to resolve responses dir path: %w", err)
		}
		ctx.ResponsesDir = absPath
		updated = true
	}

	if !updated {
		return fmt.Errorf("no context values provided. Use flags like --http-file, --env, --env-file, etc.")
	}

	// Save context
	if err := mgr.Save(ctx); err != nil {
		return err
	}

	fmt.Printf("Context saved to %s\n", mgr.GetPath())
	return executeContextShow([]string{})
}

func executeContextShow(args []string) error {
	mgr := context.NewManager()

	if !mgr.Exists() {
		fmt.Println("No context file found in current directory.")
		fmt.Printf("Use 'postie context set' to create one.\n")
		return nil
	}

	ctx, err := mgr.Load()
	if err != nil {
		return err
	}

	fmt.Printf("Context file: %s\n\n", mgr.GetPath())

	if ctx.HTTPFile != "" {
		fmt.Printf("HTTP File:         %s\n", ctx.HTTPFile)
	}
	if ctx.Environment != "" {
		fmt.Printf("Environment:       %s\n", ctx.Environment)
	}
	if ctx.EnvFile != "" {
		fmt.Printf("Env File:          %s\n", ctx.EnvFile)
	}
	if ctx.PrivateEnvFile != "" {
		fmt.Printf("Private Env File:  %s\n", ctx.PrivateEnvFile)
	}
	if ctx.SaveResponses {
		fmt.Printf("Save Responses:    %t\n", ctx.SaveResponses)
	}
	if ctx.ResponsesDir != "" {
		fmt.Printf("Responses Dir:     %s\n", ctx.ResponsesDir)
	}

	if ctx.HTTPFile == "" && ctx.Environment == "" && ctx.EnvFile == "" &&
		ctx.PrivateEnvFile == "" && !ctx.SaveResponses && ctx.ResponsesDir == "" {
		fmt.Println("Context is empty.")
	}

	return nil
}

func executeContextClear(args []string) error {
	mgr := context.NewManager()

	if !mgr.Exists() {
		fmt.Println("No context file found in current directory.")
		return nil
	}

	contextPath := mgr.GetPath()
	if err := mgr.Clear(); err != nil {
		return err
	}

	fmt.Printf("Context cleared: %s\n", contextPath)
	return nil
}
