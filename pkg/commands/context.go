package commands

import (
	"fmt"
	"path/filepath"
	"strings"

	"postie/pkg/cli"
	"postie/pkg/collection"
	"postie/pkg/context"
)

// ContextCommands returns the context command with all subcommands
func ContextCommands() *cli.Command {
	return &cli.Command{
		Name:        "context",
		Description: "Manage CLI context",
		Subcommands: map[string]*cli.Command{
			"set":   contextSetCommand(),
			"show":  contextShowCommand(),
			"clear": contextClearCommand(),
		},
	}
}

func contextSetCommand() *cli.Command {
	return &cli.Command{
		Name:        "set",
		Description: "Set the current context",
		Action: func(args []string) error {
			var collectionPath, environment string

			collFlag := &cli.StringFlag{Name: "collection", ShortName: "c", Value: collectionPath, Usage: "Collection file path"}
			envFlag := &cli.StringFlag{Name: "environment", ShortName: "e", Value: environment, Usage: "Environment name"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{collFlag, envFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			collectionPath = collFlag.Value
			environment = envFlag.Value

			// Load current context
			ctx, err := context.Load()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			// Update collection if provided
			if collectionPath != "" {
				absPath, _ := filepath.Abs(collectionPath)
				if err := ctx.SetCollection(absPath); err != nil {
					return fmt.Errorf("error setting collection: %w", err)
				}
				fmt.Printf("✅ Collection set to: %s\n", absPath)
			}

			// Update environment if provided
			if environment != "" {
				// Validate environment exists in the collection
				if ctx.HasCollection() {
					coll, err := collection.LoadCollection(ctx.GetCollection())
					if err != nil {
						return fmt.Errorf("error loading collection: %w", err)
					}
					if _, err := coll.GetEnvironment(environment); err != nil {
						return fmt.Errorf("error: %w\nAvailable environments: %s", err, strings.Join(coll.GetEnvironmentNames(), ", "))
					}
				}
				ctx.SetEnvironment(environment)
				fmt.Printf("✅ Environment set to: %s\n", environment)
			}

			// Save context
			if err := ctx.Save(); err != nil {
				return fmt.Errorf("error saving context: %w", err)
			}

			if collectionPath == "" && environment == "" {
				fmt.Println("No changes made. Use --collection or --environment to set values.")
			}

			return nil
		},
	}
}

func contextShowCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show the current context",
		Action: func(args []string) error {
			var format string

			formatFlag := &cli.StringFlag{Name: "format", Value: format, Usage: "Output format (table, json, yaml)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{formatFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			format = formatFlag.Value

			ctx, err := context.Load()
			if err != nil {
				return fmt.Errorf("error loading context: %w", err)
			}

			if ctx.IsEmpty() {
				fmt.Println("No context is currently set.")
				fmt.Println()
				fmt.Println("Use 'postie context set' to configure a default collection and environment:")
				fmt.Println("  postie context set --collection <file> --environment <environment>")
				return nil
			}

			fmt.Println("Current Context:")
			fmt.Println(strings.Repeat("=", 50))

			if ctx.HasCollection() {
				fmt.Printf("Collection:  %s\n", ctx.GetCollection())

				// Try to load and show collection info
				if coll, err := collection.LoadCollection(ctx.GetCollection()); err == nil {
					fmt.Printf("Name:        %s\n", coll.Collection.Info.Name)
					if coll.Collection.Info.Description != "" {
						fmt.Printf("Description: %s\n", coll.Collection.Info.Description)
					}
				}
			} else {
				fmt.Println("Collection:  (not set)")
			}

			if ctx.HasEnvironment() {
				fmt.Printf("Environment: %s\n", ctx.GetEnvironment())
			} else {
				fmt.Println("Environment: (not set)")
			}

			fmt.Println(strings.Repeat("=", 50))

			return nil
		},
	}
}

func contextClearCommand() *cli.Command {
	return &cli.Command{
		Name:        "clear",
		Description: "Clear the current context",
		Action: func(args []string) error {
			var force bool

			forceFlag := &cli.BoolFlag{Name: "force", Value: force, Usage: "Skip confirmation prompt"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{}, []*cli.BoolFlag{forceFlag})
			if err != nil {
				return err
			}

			force = forceFlag.Value

			// Confirm if not forced
			if !force {
				fmt.Print("Are you sure you want to clear the current context? (y/N): ")
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "yes" && response != "Y" && response != "YES" {
					fmt.Println("Clear cancelled")
					return nil
				}
			}

			if err := context.Clear(); err != nil {
				return fmt.Errorf("error clearing context: %w", err)
			}

			fmt.Println("✅ Context cleared successfully")

			return nil
		},
	}
}
