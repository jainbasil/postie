package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	configDirName  = ".postie"
	configFileName = "context.json"
)

// Context represents the current CLI context
type Context struct {
	Collection  string `json:"collection,omitempty"`
	Environment string `json:"environment,omitempty"`
}

// GetConfigDir returns the path to the Postie configuration directory
func GetConfigDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get user home directory: %w", err)
	}
	return filepath.Join(home, configDirName), nil
}

// GetConfigFilePath returns the path to the context configuration file
func GetConfigFilePath() (string, error) {
	configDir, err := GetConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(configDir, configFileName), nil
}

// ensureConfigDir creates the config directory if it doesn't exist
func ensureConfigDir() error {
	configDir, err := GetConfigDir()
	if err != nil {
		return err
	}

	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	return nil
}

// Load loads the current context from the configuration file
func Load() (*Context, error) {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return nil, err
	}

	// If config file doesn't exist, return empty context
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return &Context{}, nil
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, fmt.Errorf("failed to parse context JSON: %w", err)
	}

	return &ctx, nil
}

// Save saves the context to the configuration file
func (c *Context) Save() error {
	if err := ensureConfigDir(); err != nil {
		return err
	}

	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// Clear removes the context configuration file
func Clear() error {
	configPath, err := GetConfigFilePath()
	if err != nil {
		return err
	}

	// If file doesn't exist, nothing to clear
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil
	}

	if err := os.Remove(configPath); err != nil {
		return fmt.Errorf("failed to clear context: %w", err)
	}

	return nil
}

// SetCollection sets the collection in the context
func (c *Context) SetCollection(collectionPath string) error {
	// Convert to absolute path
	absPath, err := filepath.Abs(collectionPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %w", err)
	}

	// Verify the collection file exists
	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return fmt.Errorf("collection file not found: %s", absPath)
	}

	c.Collection = absPath
	return nil
}

// SetEnvironment sets the environment in the context
func (c *Context) SetEnvironment(environment string) {
	c.Environment = environment
}

// GetCollection returns the current collection path
func (c *Context) GetCollection() string {
	return c.Collection
}

// GetEnvironment returns the current environment
func (c *Context) GetEnvironment() string {
	return c.Environment
}

// IsEmpty returns true if the context has no values set
func (c *Context) IsEmpty() bool {
	return c.Collection == "" && c.Environment == ""
}

// HasCollection returns true if a collection is set
func (c *Context) HasCollection() bool {
	return c.Collection != ""
}

// HasEnvironment returns true if an environment is set
func (c *Context) HasEnvironment() bool {
	return c.Environment != ""
}
