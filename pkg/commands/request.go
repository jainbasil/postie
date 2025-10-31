package commands

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"postie/pkg/cli"
	"postie/pkg/collection"
	"postie/pkg/context"
)

// RequestCommands returns the request command with all subcommands
func RequestCommands() *cli.Command {
	return &cli.Command{
		Name:        "request",
		Description: "Manage API requests",
		Subcommands: map[string]*cli.Command{
			"create":  requestCreateCommand(),
			"update":  requestUpdateCommand(),
			"list":    requestListCommand(),
			"show":    requestShowCommand(),
			"delete":  requestDeleteCommand(),
			"run":     requestRunCommand(),
			"run-all": requestRunAllCommand(),
		},
	}
}

func requestCreateCommand() *cli.Command {
	return &cli.Command{
		Name:        "create",
		Description: "Create a new API request",
		Action: func(args []string) error {
			var name, method, url, group, file, id, description, body string

			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Request name (required)", Required: true}
			methodFlag := &cli.StringFlag{Name: "method", ShortName: "m", Value: method, Usage: "HTTP method (required)", Required: true}
			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "Request URL (required)", Required: true}
			groupFlag := &cli.StringFlag{Name: "group", ShortName: "g", Value: group, Usage: "Group ID to add request to (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Custom request ID (auto-generated from name if not provided)"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "Request description"}
			bodyFlag := &cli.StringFlag{Name: "body", Value: body, Usage: "Request body (JSON string)"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{nameFlag, methodFlag, urlFlag, groupFlag, fileFlag, idFlag, descFlag, bodyFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			name = nameFlag.Value
			method = strings.ToUpper(methodFlag.Value)
			url = urlFlag.Value
			group = groupFlag.Value
			file = fileFlag.Value
			id = idFlag.Value
			description = descFlag.Value
			body = bodyFlag.Value

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
				description = fmt.Sprintf("%s operation for %s", method, name)
			}

			// Load existing collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Create new request
			newRequest := collection.Item{
				ID:          id,
				Name:        name,
				Description: description,
				Request: &collection.Request{
					Method: method,
					Header: []collection.Header{},
					URL:    url,
				},
			}

			// Add body if provided
			if body != "" {
				newRequest.Request.Body = &collection.Body{
					Mode: "raw",
					Raw:  body,
					Options: map[string]interface{}{
						"raw": map[string]interface{}{
							"language": "json",
						},
					},
				}
			}

			// Find the API group and add request
			found := false
			for i := range coll.Collection.ApiGroup {
				if coll.Collection.ApiGroup[i].ID == group {
					coll.Collection.ApiGroup[i].Apis = append(coll.Collection.ApiGroup[i].Apis, newRequest)
					found = true
					break
				}
			}

			if !found {
				return fmt.Errorf("request group with ID '%s' not found", group)
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

			fmt.Printf("Request '%s' created successfully\n", name)
			fmt.Printf("ID: %s\n", id)
			fmt.Printf("%s %s\n", method, url)
			fmt.Printf("Collection: %s\n", file)

			return nil
		},
	}
}

func requestUpdateCommand() *cli.Command {
	return &cli.Command{
		Name:        "update",
		Description: "Update an existing request",
		Action: func(args []string) error {
			var id, file, name, method, url, description, body string

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Request ID (required)", Required: true}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "New request name"}
			methodFlag := &cli.StringFlag{Name: "method", ShortName: "m", Value: method, Usage: "New HTTP method"}
			urlFlag := &cli.StringFlag{Name: "url", ShortName: "u", Value: url, Usage: "New URL"}
			descFlag := &cli.StringFlag{Name: "description", ShortName: "d", Value: description, Usage: "New description"}
			bodyFlag := &cli.StringFlag{Name: "body", Value: body, Usage: "Update request body"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{idFlag, fileFlag, nameFlag, methodFlag, urlFlag, descFlag, bodyFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			id = idFlag.Value
			file = fileFlag.Value
			name = nameFlag.Value
			method = strings.ToUpper(methodFlag.Value)
			url = urlFlag.Value
			description = descFlag.Value
			body = bodyFlag.Value

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

			// Find and update request
			found := false
			for gi := range coll.Collection.ApiGroup {
				for ri := range coll.Collection.ApiGroup[gi].Apis {
					if coll.Collection.ApiGroup[gi].Apis[ri].ID == id {
						if name != "" {
							coll.Collection.ApiGroup[gi].Apis[ri].Name = name
						}
						if description != "" {
							coll.Collection.ApiGroup[gi].Apis[ri].Description = description
						}
						if coll.Collection.ApiGroup[gi].Apis[ri].Request != nil {
							if method != "" {
								coll.Collection.ApiGroup[gi].Apis[ri].Request.Method = method
							}
							if url != "" {
								coll.Collection.ApiGroup[gi].Apis[ri].Request.URL = url
							}
							if body != "" {
								if coll.Collection.ApiGroup[gi].Apis[ri].Request.Body == nil {
									coll.Collection.ApiGroup[gi].Apis[ri].Request.Body = &collection.Body{}
								}
								coll.Collection.ApiGroup[gi].Apis[ri].Request.Body.Mode = "raw"
								coll.Collection.ApiGroup[gi].Apis[ri].Request.Body.Raw = body
							}
						}
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			if !found {
				return fmt.Errorf("request with ID '%s' not found", id)
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

			fmt.Printf("Request updated successfully\n")
			fmt.Printf("ID: %s\n", id)
			fmt.Printf("Collection: %s\n", file)

			return nil
		},
	}
}

func requestListCommand() *cli.Command {
	return &cli.Command{
		Name:        "list",
		Description: "List all requests in a collection",
		Action: func(args []string) error {
			var file, group, environment string

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			groupFlag := &cli.StringFlag{Name: "group", ShortName: "g", Value: group, Usage: "Filter by group ID"}
			envFlag := &cli.StringFlag{Name: "environment", ShortName: "e", Value: environment, Usage: "Show with environment variables resolved"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag, groupFlag, envFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			file = fileFlag.Value
			group = groupFlag.Value
			environment = envFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
				if environment == "" && ctx.HasEnvironment() {
					environment = ctx.GetEnvironment()
				}
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Create runner for variable resolution
			runner := collection.NewRunner(coll, environment)
			runner.ListRequests()

			return nil
		},
	}
}

func requestShowCommand() *cli.Command {
	return &cli.Command{
		Name:        "show",
		Description: "Show details of a specific request",
		Action: func(args []string) error {
			var id, name, file, environment string

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Request ID (use this OR --name)"}
			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Request name (use this OR --id)"}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			envFlag := &cli.StringFlag{Name: "environment", ShortName: "e", Value: environment, Usage: "Show with environment variables resolved"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{idFlag, nameFlag, fileFlag, envFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			id = idFlag.Value
			name = nameFlag.Value
			file = fileFlag.Value
			environment = envFlag.Value

			if id == "" && name == "" {
				return fmt.Errorf("either --id or --name must be provided")
			}

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
				if environment == "" && ctx.HasEnvironment() {
					environment = ctx.GetEnvironment()
				}
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Find request
			var foundReq *collection.Item
			var foundGroup *collection.Item
			for gi := range coll.Collection.ApiGroup {
				for ri := range coll.Collection.ApiGroup[gi].Apis {
					if (id != "" && coll.Collection.ApiGroup[gi].Apis[ri].ID == id) ||
						(name != "" && coll.Collection.ApiGroup[gi].Apis[ri].Name == name) {
						foundReq = &coll.Collection.ApiGroup[gi].Apis[ri]
						foundGroup = &coll.Collection.ApiGroup[gi]
						break
					}
				}
				if foundReq != nil {
					break
				}
			}

			if foundReq == nil {
				return fmt.Errorf("request not found")
			}

			// Display request details
			fmt.Printf("Request: %s\n", foundReq.Name)
			fmt.Printf("ID: %s\n", foundReq.ID)
			fmt.Printf("Group: %s\n", foundGroup.Name)
			if foundReq.Description != "" {
				fmt.Printf("Description: %s\n", foundReq.Description)
			}
			if foundReq.Request != nil {
				fmt.Printf("\nMethod: %s\n", foundReq.Request.Method)

				// Resolve variables if environment is provided
				url := foundReq.Request.URL
				if urlStr, ok := url.(string); ok {
					if environment != "" {
						variables := coll.ResolveVariables(environment, foundGroup)
						url = collection.ReplaceVariables(urlStr, variables)
					}
				}
				fmt.Printf("URL: %v\n", url)

				if len(foundReq.Request.Header) > 0 {
					fmt.Printf("\nHeaders:\n")
					for _, h := range foundReq.Request.Header {
						fmt.Printf("  %s: %s\n", h.Key, h.Value)
					}
				}

				if foundReq.Request.Body != nil && foundReq.Request.Body.Raw != "" {
					fmt.Printf("\nBody:\n%s\n", foundReq.Request.Body.Raw)
				}
			}

			return nil
		},
	}
}

func requestDeleteCommand() *cli.Command {
	return &cli.Command{
		Name:        "delete",
		Description: "Delete a request",
		Action: func(args []string) error {
			var id, file string
			var force bool

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Request ID (required)", Required: true}
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

			// Find request for confirmation
			var requestName string
			for _, group := range coll.Collection.ApiGroup {
				for _, req := range group.Apis {
					if req.ID == id {
						requestName = req.Name
						break
					}
				}
				if requestName != "" {
					break
				}
			}

			if requestName == "" {
				return fmt.Errorf("request with ID '%s' not found", id)
			}

			// Confirm deletion if not forced
			if !force {
				fmt.Printf("Are you sure you want to delete request '%s' (ID: %s)? (y/N): ", requestName, id)
				var response string
				fmt.Scanln(&response)
				if response != "y" && response != "yes" && response != "Y" && response != "YES" {
					fmt.Println("Deletion cancelled")
					return nil
				}
			}

			// Find and remove request
			found := false
			for gi := range coll.Collection.ApiGroup {
				for ri := range coll.Collection.ApiGroup[gi].Apis {
					if coll.Collection.ApiGroup[gi].Apis[ri].ID == id {
						coll.Collection.ApiGroup[gi].Apis = append(
							coll.Collection.ApiGroup[gi].Apis[:ri],
							coll.Collection.ApiGroup[gi].Apis[ri+1:]...,
						)
						found = true
						break
					}
				}
				if found {
					break
				}
			}

			if !found {
				return fmt.Errorf("request with ID '%s' not found", id)
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

			fmt.Printf("Request '%s' deleted successfully\n", requestName)

			return nil
		},
	}
}

func requestRunCommand() *cli.Command {
	return &cli.Command{
		Name:        "run",
		Description: "Run a specific request",
		Action: func(args []string) error {
			var id, name, file, environment string

			idFlag := &cli.StringFlag{Name: "id", Value: id, Usage: "Request ID (use this OR --name)"}
			nameFlag := &cli.StringFlag{Name: "name", ShortName: "n", Value: name, Usage: "Request name (use this OR --id)"}
			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			envFlag := &cli.StringFlag{Name: "environment", ShortName: "e", Value: environment, Usage: "Environment to use"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{idFlag, nameFlag, fileFlag, envFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			id = idFlag.Value
			name = nameFlag.Value
			file = fileFlag.Value
			environment = envFlag.Value

			if id == "" && name == "" {
				return fmt.Errorf("either --id or --name must be provided")
			}

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
				if environment == "" && ctx.HasEnvironment() {
					environment = ctx.GetEnvironment()
				}
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Create runner
			runner := collection.NewRunner(coll, environment)

			// Run request
			if id != "" {
				return runner.RunByID(id)
			}
			return runner.RunByName(name)
		},
	}
}

func requestRunAllCommand() *cli.Command {
	return &cli.Command{
		Name:        "run-all",
		Description: "Run all requests in a collection",
		Action: func(args []string) error {
			var file, environment, group string

			fileFlag := &cli.StringFlag{Name: "file", ShortName: "f", Value: file, Usage: "Collection file path (uses context if not provided)"}
			envFlag := &cli.StringFlag{Name: "environment", ShortName: "e", Value: environment, Usage: "Environment to use"}
			groupFlag := &cli.StringFlag{Name: "group", ShortName: "g", Value: group, Usage: "Run only requests in specific group"}

			_, err := cli.ParseFlags(args, []*cli.StringFlag{fileFlag, envFlag, groupFlag}, []*cli.BoolFlag{})
			if err != nil {
				return err
			}

			file = fileFlag.Value
			environment = envFlag.Value
			group = groupFlag.Value

			// Use context if file not provided
			if file == "" {
				ctx, err := context.Load()
				if err != nil || !ctx.HasCollection() {
					return fmt.Errorf("no collection file specified and no context set")
				}
				file = ctx.GetCollection()
				if environment == "" && ctx.HasEnvironment() {
					environment = ctx.GetEnvironment()
				}
			}

			// Load collection
			coll, err := collection.LoadCollection(file)
			if err != nil {
				return fmt.Errorf("error loading collection: %w", err)
			}

			// Create runner
			runner := collection.NewRunner(coll, environment)

			// Run all requests (filtering by group is a TODO enhancement)
			return runner.RunAll()
		},
	}
}
