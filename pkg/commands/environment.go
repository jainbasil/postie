package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"postie/pkg/cli"
	"postie/pkg/collection"
	"postie/pkg/context"
)

// EnvironmentCommands returns the environment command with all subcommands
func EnvironmentCommands() *cli.Command {
	return &cli.Command{
		Name:        "environment",
		Description: "Manage environments",
		Subcommands: map[string]*cli.Command{
			"create":   environmentCreateCommand(),
			"update":   environmentUpdateCommand(),
			"list":     environmentListCommand(),
			"delete":   environmentDeleteCommand(),
			"variable": environmentVariableCommand(),
		},
	}
}

func environmentCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Description: "Create a new environment",
		Action: func(args []string) error {
			var name, file, description string
			var setContext bool

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Environment name (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "Environment description"}
			contextFlag := &cli.BoolFlag{Name: "set-context", Value: setContext, Usage: "Set this environment as current in context"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, fileFlag, descFlag}, []*cli.BoolFlag{contextFlag})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			file = fileFlag.Value
			description = descFlag.Value
			setContext = contextFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Check if environment already exists
			for _, env := range coll.Collection.Environment {
				if env.Name == name {
					return fmt.Errorf("environment '%s' already exists", name)
				}
			}

			// Create new environment
			newEnv := collection.Environment{
				Name:        name,
				Description: description,
				Values:      []collection.Variable{},
			}

			// Add to collection
			coll.Collection.Environment = append(coll.Collection.Environment, newEnv)

			// Save back to file
			data, err := json.MarshalIndent(coll, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("✅ Environment '%s' created successfully\n", name)
			fmt.Printf("📁 Collection: %s\n", file)

			// Set context if requested
			if setContext {
				ctx, _ := context.Load()
				ctx.SetEnvironment(name)
				if err := ctx.Save(); err != nil {
					fmt.Printf("⚠️  Warning: Could not set context: %v\n", err)
				} else {
					fmt.Printf("📌 Environment set as current context\n")
				}
			}

			return nil
		},
	}
}

func environmentUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:        "update",
		Description: "Update an existing environment",
		Action: func(args []string) error {
			var name, file, description string
			var setContext bool

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Environment name (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "New description"}
			contextFlag := &cli.BoolFlag{Name: "set-context", Value: setContext, Usage: "Set as current environment in context"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, fileFlag, descFlag}, []*cli.BoolFlag{contextFlag})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			file = fileFlag.Value
			description = descFlag.Value
			setContext = contextFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Find and update environment
			found := false
			for i := range coll.Collection.Environment {
				if coll.Collection.Environment[i].Name == name {
					if description != "" {
						coll.Collection.Environment[i].Description = description
					}
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("environment '%s' not found", name)
			}

			// Save back to file
			data, err := json.MarshalIndent(coll, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("✅ Environment '%s' updated successfully\n", name)
			fmt.Printf("📁 Collection: %s\n", file)

			// Set context if requested
			if setContext {
				ctx, _ := context.Load()
				ctx.SetEnvironment(name)
				if err := ctx.Save(); err != nil {
					fmt.Printf("⚠️  Warning: Could not set context: %v\n", err)
				} else {
					fmt.Printf("📌 Environment set as current context\n")
				}
			}

			return nil
		},
	}
}

func environmentListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List all environments in a collection",
		Action: func(args []string) error {
			var file string

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			file = fileFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
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
					fmt.Printf("   Description: %s\n", env.Description)
				}
				fmt.Printf("   Variables: %d\n", len(env.Values))
				if env.Auth != nil {
					fmt.Printf("   Authentication: %s\n", env.Auth.Type)
				}
				fmt.Println()
			}

			return nil
		},
	}
}

func environmentDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Delete an environment",
		Action: func(args []string) error {
			var name, file string

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Environment name (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, fileFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			file = fileFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Find and remove environment
			found := false
			for i := range coll.Collection.Environment {
				if coll.Collection.Environment[i].Name == name {
					coll.Collection.Environment = append(coll.Collection.Environment[:i], coll.Collection.Environment[i+1:]...)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("environment '%s' not found", name)
			}

			// Save back to file
			data, err := json.MarshalIndent(coll, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("✅ Environment '%s' deleted successfully\n", name)

			return nil
		},
	}
}

func environmentVariableCommand() *cli.Command {
	return &cli.Command{
		Name:        "variable",
		Description: "Manage environment variables",
		Subcommands: map[string]*cli.Command{
			"set":  envVariableSetCommand(),
			"get":  envVariableGetCommand(),
			"list": envVariableListCommand(),
		},
	}
}

func envVariableSetCommand() *cli.Command {
	return &cli.Command{
		Name:        "set",
		Description: "Set an environment variable",
		Action: func(args []string) error {
			var name, key, value, file string
			var secret bool

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Environment name (required)", Required: true}
			keyFlag := &cli.StringFlag{Name: "key", ShortName: "k", Value: key, Usage: "Variable key (required)", Required: true}
			valueFlag := &cli.StringFlag{Name: "value", ShortName: "v", Value: value, Usage: "Variable value (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			secretFlag := &cli.BoolFlag{Name: "secret", Value: secret, Usage: "Mark as secret/sensitive variable"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, keyFlag, valueFlag, fileFlag}, []*cli.BoolFlag{secretFlag})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			key = keyFlag.Value
			value = valueFlag.Value
			file = fileFlag.Value
			secret = secretFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Find environment
			found := false
			for i := range coll.Collection.Environment {
				if coll.Collection.Environment[i].Name == name {
					// Check if variable exists and update, or add new
					varFound := false
					for j := range coll.Collection.Environment[i].Values {
						if coll.Collection.Environment[i].Values[j].Key == key {
							coll.Collection.Environment[i].Values[j].Value = value
							varFound = true
							break
						}
					}

					if !varFound {
						newVar := collection.Variable{
							Key:     key,
							Value:   value,
							Enabled: true,
						}
						coll.Collection.Environment[i].Values = append(coll.Collection.Environment[i].Values, newVar)
					}

					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("environment '%s' not found", name)
			}

			// Save back to file
			data, err := json.MarshalIndent(coll, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("✅ Variable '%s' set in environment '%s'\n", key, name)
			if !secret {
				fmt.Printf("   Value: %s\n", value)
			}

			return nil
		},
	}
}

func envVariableGetCommand() *cli.Command {
	return &cli.Command{
		Name:        "get",
		Description: "Get an environment variable",
		Action: func(args []string) error {
			var name, key, file string

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Environment name (required)", Required: true}
			keyFlag := &cli.StringFlag{Name: "key", ShortName: "k", Value: key, Usage: "Variable key"}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, keyFlag, fileFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			key = keyFlag.Value
			file = fileFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Find environment
			for _, env := range coll.Collection.Environment {
				if env.Name == name {
					if key != "" {
						// Get specific variable
						for _, v := range env.Values {
							if v.Key == key {
								fmt.Printf("%s=%s\n", v.Key, v.Value)
								return nil
							}
						}
						return fmt.Errorf("variable '%s' not found in environment '%s'", key, name)
					} else {
						// List all variables (same as list command)
						fmt.Printf("Variables in environment '%s':\n\n", name)
						for _, v := range env.Values {
							enabledMark := "✓"
							if !v.Enabled {
								enabledMark = "✗"
							}
							fmt.Printf("%s %s = %s\n", enabledMark, v.Key, v.Value)
						}
						return nil
					}
				}
			}

			return fmt.Errorf("environment '%s' not found", name)
		},
	}
}

func envVariableListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List all variables in an environment",
		Action: func(args []string) error {
			// Same implementation as get without key
			return envVariableGetCommand().Action(args)
		},
	}
}
