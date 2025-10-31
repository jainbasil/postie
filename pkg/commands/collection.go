package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"postie/pkg/cli"
	"postie/pkg/collection"
	"postie/pkg/context"
)

// CollectionCommands returns the collection command with all subcommands
func CollectionCommands() *cli.Command {
	return &cli.Command{
		Name:        "collection",
		Description: "Manage API collections",
		Subcommands: map[string]*cli.Command{
			"create": collectionCreateCommand(),
			"update": collectionUpdateCommand(),
			"show":   collectionShowCommand(),
			"list":   collectionListCommand(),
			"delete": collectionDeleteCommand(),
		},
	}
}

func collectionCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Description: "Create a new collection",
		Action: func(args []string) error {
			var name, file, description string
			var setContext bool

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Collection name (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Output file path (auto-generated if not provided)"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "Collection description"}
			contextFlag := &cli.BoolFlag{Name: "set-context", Value: setContext, Usage: "Set this collection as the current context"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, fileFlag, descFlag}, []*cli.BoolFlag{contextFlag})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			file = fileFlag.Value
			description = descFlag.Value
			setContext = contextFlag.Value

			// Generate filename if not provided
			if file == "" {
				file = strings.ToLower(strings.ReplaceAll(name, " ", "-")) + ".collection.json"
			}

			// Set description if not provided
			if description == "" {
				description = fmt.Sprintf("API collection for %s", name)
			}

			// Create new collection
			newCollection := &collection.Collection{
				Collection: collection.CollectionInfo{
					Info: collection.Info{
						Name:        name,
						Description: description,
						Version:     "1.0.0",
						Schema:      "https://postie.dev/collection/v1.0.0/collection.json",
					},
					Variable:    []collection.Variable{},
					Environment: []collection.Environment{},
					ApiGroup:    []collection.Item{},
				},
			}

			// Save to file
			data, err := json.MarshalIndent(newCollection, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("Collection '%s' created successfully\n", name)
			fmt.Printf("File: %s\n", file)

			// Set context if requested
			if setContext {
				ctx, _ := context.Load()
				absPath, _ := filepath.Abs(file)
				ctx.SetCollection(absPath)
				if err := ctx.Save(); err != nil {
					fmt.Printf("Warning: Could not set context: %v\n", err)
				} else {
					fmt.Printf("Collection set as current context\n")
				}
			}

			return nil
		},
	}
}

func collectionUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:        "update",
		Description: "Update an existing collection",
		Action: func(args []string) error {
			var file, name, description string
			var setContext bool

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "New collection name"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "New collection description"}
			contextFlag := &cli.BoolFlag{Name: "set-context", Value: setContext, Usage: "Set this collection as the current context"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag, nameFlag, descFlag}, []*cli.BoolFlag{contextFlag})
			if err != nil {
				return err
			}

			file = fileFlag.Value
			name = nameFlag.Value
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

			// Update fields
			if name != "" {
				coll.Collection.Info.Name = name
			}
			if description != "" {
				coll.Collection.Info.Description = description
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

			fmt.Printf("Collection updated successfully\n")
			fmt.Printf("File: %s\n", file)
			fmt.Printf("Name: %s\n", coll.Collection.Info.Name)

			// Set context if requested
			if setContext {
				ctx, _ := context.Load()
				absPath, _ := filepath.Abs(file)
				ctx.SetCollection(absPath)
				if err := ctx.Save(); err != nil {
					fmt.Printf("Warning: Could not set context: %v\n", err)
				} else {
					fmt.Printf("Collection set as current context\n")
				}
			}

			return nil
		},
	}
}

func collectionShowCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show collection details",
		Action: func(args []string) error {
			var file, format string

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			formatFlag := &cli.StringFlag{Name: "format", Value: format, Usage: "Output format (json, table, yaml)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag, formatFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			file = fileFlag.Value
			format = formatFlag.Value

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

			// Display based on format
			if format == "json" {
				data, _ := json.MarshalIndent(coll, "", "  ")
				fmt.Println(string(data))
			} else {
				// Table format (default)
				fmt.Printf("Collection: %s\n", coll.Collection.Info.Name)
				fmt.Printf("Description: %s\n", coll.Collection.Info.Description)
				fmt.Printf("Version: %s\n", coll.Collection.Info.Version)
				fmt.Printf("File: %s\n", file)
				fmt.Printf("\nAPI Groups: %d\n", len(coll.Collection.ApiGroup))
				fmt.Printf("Environments: %d\n", len(coll.Collection.Environment))
				fmt.Printf("Variables: %d\n", len(coll.Collection.Variable))
			}

			return nil
		},
	}
}

func collectionListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List all collections in a directory",
		Action: func(args []string) error {
			var directory string
			var recursive bool

			dirFlag := &cli.StringFlag{Name: "directory", ShortName: "d", Value: directory, Usage: "Directory to search for collections (default: current directory)"}
			recursiveFlag := &cli.BoolFlag{Name: "recursive", ShortName: "r", Value: recursive, Usage: "Search recursively"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{dirFlag}, []*cli.BoolFlag{recursiveFlag})
			if err != nil {
				return err
			}

			directory = dirFlag.Value
			recursive = recursiveFlag.Value

			if directory == "" {
				directory = "."
			}

			// Find collection files
			var collections []string

			if recursive {
				err = filepath.Walk(directory, func(path string, info os.FileInfo, err error) error {
					if err != nil {
						return err
					}
					if !info.IsDir() && strings.HasSuffix(path, ".collection.json") {
						collections = append(collections, path)
					}
					return nil
				})
			} else {
				files, err := os.ReadDir(directory)
				if err != nil {
					return fmt.Errorf("error reading directory: %w", err)
				}
				for _, file := range files {
					if !file.IsDir() && strings.HasSuffix(file.Name(), ".collection.json") {
						collections = append(collections, filepath.Join(directory, file.Name()))
					}
				}
			}

			if err != nil {
				return fmt.Errorf("error finding collections: %w", err)
			}

			if len(collections) == 0 {
				fmt.Println("No collections found")
				return nil
			}

			fmt.Printf("Found %d collection(s):\n\n", len(collections))
			for i, collPath := range collections {
				coll, err := collection.LoadCollection(collPath)
				if err != nil {
					fmt.Printf("%d. %s (error loading: %v)\n", i+1, collPath, err)
					continue
				}
				fmt.Printf("%d. %s\n", i+1, coll.Collection.Info.Name)
				fmt.Printf("   File: %s\n", collPath)
				if coll.Collection.Info.Description != "" {
					fmt.Printf("   Description: %s\n", coll.Collection.Info.Description)
				}
				fmt.Printf("   Groups: %d, Environments: %d\n", len(coll.Collection.ApiGroup), len(coll.Collection.Environment))
				fmt.Println()
			}

			return nil
		},
	}
}

func collectionDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Delete a collection file",
		Action: func(args []string) error {
			var file string
			var force bool

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (required)", Required: true}
			forceFlag := &cli.BoolFlag{Name: "force", Value: force, Usage: "Skip confirmation prompt"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag}, []*cli.BoolFlag{forceFlag})
			if err != nil {
				return err
			}

			file = fileFlag.Value
			force = forceFlag.Value

			// Check if file exists
			if _, err := os.Stat(file); os.IsNotExist(err) {
				return fmt.Errorf("collection file not found: %s", file)
			}

			// Confirm deletion if not forced
			if !force {
				fmt.Printf("Are you sure you want to delete '%s'? (y/N): ", file)
				var response string
				fmt.Scanln(&response)
				if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			// Delete file
			if err := os.Remove(file); err != nil {
				return fmt.Errorf("error deleting file: %w", err)
			}

			fmt.Printf("Collection deleted: %s\n", file)

			return nil
		},
	}
}
