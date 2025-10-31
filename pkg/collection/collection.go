package collection

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

// LoadCollection loads a collection from a JSON file
func LoadCollection(filename string) (*Collection, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read collection file: %w", err)
	}

	var collection Collection
	if err := json.Unmarshal(data, &collection); err != nil {
		return nil, fmt.Errorf("failed to parse collection JSON: %w", err)
	}

	return &collection, nil
}

// SaveCollection saves a collection to a JSON file
func (c *Collection) SaveCollection(filename string) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal collection: %w", err)
	}

	if err := os.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("failed to write collection file: %w", err)
	}

	return nil
}

// GetEnvironment returns an environment by name
func (c *Collection) GetEnvironment(name string) (*Environment, error) {
	for _, env := range c.Collection.Environment {
		if env.Name == name {
			return &env, nil
		}
	}
	return nil, fmt.Errorf("environment '%s' not found", name)
}

// GetEnvironmentNames returns all environment names
func (c *Collection) GetEnvironmentNames() []string {
	names := make([]string, len(c.Collection.Environment))
	for i, env := range c.Collection.Environment {
		names[i] = env.Name
	}
	return names
}

// GetDefaultEnvironment returns the first environment or nil
func (c *Collection) GetDefaultEnvironment() *Environment {
	if len(c.Collection.Environment) > 0 {
		return &c.Collection.Environment[0]
	}
	return nil
}

// ResolveVariables resolves variables for a given environment with folder overrides
func (c *Collection) ResolveVariables(envName string, folderItem *Item) map[string]interface{} {
	variables := make(map[string]interface{})

	// Start with collection variables
	for _, variable := range c.Collection.Variable {
		variables[variable.Key] = variable.Value
	}

	// Override with environment variables
	if env, err := c.GetEnvironment(envName); err == nil {
		for _, variable := range env.Values {
			if variable.Enabled {
				variables[variable.Key] = variable.Value
			}
		}
	}

	// Override with folder environment variables if folder has environments
	if folderItem != nil && len(folderItem.Environment) > 0 {
		for _, folderEnv := range folderItem.Environment {
			if folderEnv.Name == envName {
				for _, variable := range folderEnv.Values {
					if variable.Enabled {
						variables[variable.Key] = variable.Value
					}
				}
				break
			}
		}
	}

	return variables
}

// ReplaceVariables replaces {{variable}} placeholders in a string
func ReplaceVariables(text string, variables map[string]interface{}) string {
	result := text
	for key, value := range variables {
		placeholder := fmt.Sprintf("{{%s}}", key)
		result = strings.ReplaceAll(result, placeholder, fmt.Sprintf("%v", value))
	}
	return result
}

// FindAllRequests recursively finds all requests in the collection
func (c *Collection) FindAllRequests() []RequestItem {
	var requests []RequestItem
	c.findRequestsInItems(c.Collection.ApiGroup, "", nil, &requests)
	return requests
}

func (c *Collection) findRequestsInItems(items []Item, path string, parentItem *Item, requests *[]RequestItem) {
	for _, item := range items {
		currentPath := path
		if currentPath != "" {
			currentPath += " / "
		}
		currentPath += item.Name

		if item.Request != nil {
			*requests = append(*requests, RequestItem{
				Name:       item.Name,
				Path:       currentPath,
				Request:    item.Request,
				Item:       &item,
				ParentItem: parentItem,
			})
		}

		if len(item.Apis) > 0 {
			// Pass the current item as the parent for nested requests
			c.findRequestsInItems(item.Apis, currentPath, &item, requests)
		}
	}
}

// GetAuth resolves authentication for environment, folder, and request
func (c *Collection) GetAuth(envName string, folderItem *Item, requestItem *Item) *Auth {
	// Request level auth has highest priority
	if requestItem != nil && requestItem.Auth != nil {
		return requestItem.Auth
	}

	// Folder level auth
	if folderItem != nil && folderItem.Auth != nil && folderItem.Auth.Type != "inherit" {
		return folderItem.Auth
	}

	// Environment level auth
	if env, err := c.GetEnvironment(envName); err == nil && env.Auth != nil {
		return env.Auth
	}

	// Collection level auth
	if c.Collection.Auth != nil {
		return c.Collection.Auth
	}

	return nil
}

// GetRequestURL resolves the URL for a request with variable substitution
func (c *Collection) GetRequestURL(req *Request, variables map[string]interface{}) string {
	var urlStr string

	switch url := req.URL.(type) {
	case string:
		urlStr = url
	case map[string]interface{}:
		if raw, ok := url["raw"].(string); ok {
			urlStr = raw
		}
	}

	return ReplaceVariables(urlStr, variables)
}

// GenerateSlug creates a URL-friendly slug from a name
func GenerateSlug(name string) string {
	// Convert to lowercase
	slug := strings.ToLower(name)

	// Replace spaces and special characters with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")
	slug = strings.ReplaceAll(slug, "_", "-")
	slug = strings.ReplaceAll(slug, "/", "-")
	slug = strings.ReplaceAll(slug, "\\", "-")
	slug = strings.ReplaceAll(slug, ".", "-")
	slug = strings.ReplaceAll(slug, ",", "-")
	slug = strings.ReplaceAll(slug, ":", "-")
	slug = strings.ReplaceAll(slug, ";", "-")
	slug = strings.ReplaceAll(slug, "!", "-")
	slug = strings.ReplaceAll(slug, "?", "-")
	slug = strings.ReplaceAll(slug, "&", "-and-")
	slug = strings.ReplaceAll(slug, "+", "-plus-")
	slug = strings.ReplaceAll(slug, "=", "-equals-")
	slug = strings.ReplaceAll(slug, "(", "-")
	slug = strings.ReplaceAll(slug, ")", "-")
	slug = strings.ReplaceAll(slug, "[", "-")
	slug = strings.ReplaceAll(slug, "]", "-")
	slug = strings.ReplaceAll(slug, "{", "-")
	slug = strings.ReplaceAll(slug, "}", "-")
	slug = strings.ReplaceAll(slug, "<", "-")
	slug = strings.ReplaceAll(slug, ">", "-")
	slug = strings.ReplaceAll(slug, "@", "-at-")
	slug = strings.ReplaceAll(slug, "#", "-")
	slug = strings.ReplaceAll(slug, "$", "-")
	slug = strings.ReplaceAll(slug, "%", "-")
	slug = strings.ReplaceAll(slug, "^", "-")
	slug = strings.ReplaceAll(slug, "*", "-")
	slug = strings.ReplaceAll(slug, "|", "-")
	slug = strings.ReplaceAll(slug, "\"", "-")
	slug = strings.ReplaceAll(slug, "'", "-")
	slug = strings.ReplaceAll(slug, "`", "-")
	slug = strings.ReplaceAll(slug, "~", "-")

	// Remove multiple consecutive hyphens
	for strings.Contains(slug, "--") {
		slug = strings.ReplaceAll(slug, "--", "-")
	}

	// Remove leading and trailing hyphens
	slug = strings.Trim(slug, "-")

	// If empty, return a default
	if slug == "" {
		slug = "unnamed"
	}

	return slug
}
