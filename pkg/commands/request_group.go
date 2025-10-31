package commands

import (
	"encoding/json"
	"fmt"
	"os"

	"postie/pkg/cli"
	"postie/pkg/collection"
	"postie/pkg/context"
)

// RequestGroupCommands returns the request-group command with all subcommands
func RequestGroupCommands() *cli.Command {
	return &cli.Command{
		Name:        "request-group",
		Description: "Manage API request groups",
		Subcommands: map[string]*cli.Command{
			"create": requestGroupCreateCommand(),
			"update": requestGroupUpdateCommand(),
			"list":   requestGroupListCommand(),
			"delete": requestGroupDeleteCommand(),
		},
	}
}

func requestGroupCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Description: "Create a new request group",
		Action: func(args []string) error {
			var name, file, id, description string

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Group name (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Custom group ID (auto-generated from name if not provided)"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "Group description"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, fileFlag, idFlag, descFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			file = fileFlag.Value
			id = idFlag.Value
			description = descFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
			}

			// Generate ID if not provided
			if id == "" {
				id = collection.GenerateSlug(name)
			}

			// Set description if not provided
			if description == "" {
				description = fmt.Sprintf("API group for %s operations", name)
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Create new API group
			newGroup := collection.Item{
				ID:          id,
				Name:        name,
				Description: description,
				Apis:        []collection.Item{},
				Environment: []collection.Environment{},
			}

			// Add to collection
			coll.Collection.ApiGroup = append(coll.Collection.ApiGroup, newGroup)

			// Save back to file
			data, err := json.MarshalIndent(coll, "", "  ")
			if err != nil {
				return fmt.Errorf("error marshaling collection: %w", err)
			}

			err = os.WriteFile(file, data, 0644)
			if err != nil {
				return fmt.Errorf("error writing file: %w", err)
			}

			fmt.Printf("✅ Request Group '%s' created successfully\n", name)
			fmt.Printf("🆔 ID: %s\n", id)
			fmt.Printf("📁 Collection: %s\n", file)

			return nil
		},
	}
}

func requestGroupUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:        "update",
		Description: "Update an existing request group",
		Action: func(args []string) error {
			var id, file, name, description string

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Group ID (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "New group name"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "New description"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{idFlag, fileFlag, nameFlag, descFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			id = idFlag.Value
			file = fileFlag.Value
			name = nameFlag.Value
			description = descFlag.Value

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

			// Find and update API group
			found := false
			for i := range coll.Collection.ApiGroup {
				if coll.Collection.ApiGroup[i].ID == id {
					if name != "" {
						coll.Collection.ApiGroup[i].Name = name
					}
					if description != "" {
						coll.Collection.ApiGroup[i].Description = description
					}
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("request group with ID '%s' not found", id)
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

			fmt.Printf("✅ Request Group updated successfully\n")
			fmt.Printf("🆔 ID: %s\n", id)
			fmt.Printf("📁 Collection: %s\n", file)

			return nil
		},
	}
}

func requestGroupListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List all request groups in a collection",
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

			// Show groups
			fmt.Printf("Collection: %s\n", coll.Collection.Info.Name)
			fmt.Printf("Request Groups (%d):\n\n", len(coll.Collection.ApiGroup))

			for i, group := range coll.Collection.ApiGroup {
				fmt.Printf("%d. %s\n", i+1, group.Name)
				fmt.Printf("   ID: %s\n", group.ID)
				if group.Description != "" {
					fmt.Printf("   Description: %s\n", group.Description)
				}
				fmt.Printf("   Requests: %d\n", len(group.Apis))
				fmt.Println()
			}

			return nil
		},
	}
}

func requestGroupDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Delete a request group",
		Action: func(args []string) error {
			var id, file string
			var force bool

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Group ID (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			forceFlag := &cli.BoolFlag{Name: "force", Value: force, Usage: "Skip confirmation prompt"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{idFlag, fileFlag}, []*cli.BoolFlag{forceFlag})
			if err != nil {
				return err
			}

			id = idFlag.Value
			file = fileFlag.Value
			force = forceFlag.Value

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

			// Find group for confirmation
			var groupName string
			for _, group := range coll.Collection.ApiGroup {
				if group.ID == id {
					groupName = group.Name
					break
				}
			}

			if groupName == "" {
				return fmt.Errorf("request group with ID '%s' not found", id)
			}

			// Confirm deletion if not forced
			if !force {
				fmt.Printf("Are you sure you want to delete group '%s' (ID: %s)? (y/N): ", groupName, id)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "yes" && response != "Y" && response != "YES" {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			// Find and remove API group
			found := false
			for i := range coll.Collection.ApiGroup {
				if coll.Collection.ApiGroup[i].ID == id {
					coll.Collection.ApiGroup = append(coll.Collection.ApiGroup[:i], coll.Collection.ApiGroup[i+1:]...)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("request group with ID '%s' not found", id)
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

			fmt.Printf("✅ Request Group '%s' deleted successfully\n", groupName)

			return nil
		},
	}
}
