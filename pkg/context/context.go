package context

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Context represents the saved context configuration for a directory
type Context struct {
	HTTPFile           string `json:"httpFile,omitempty"`
	Environment        string `json:"environment,omitempty"`
	EnvFile            string `json:"envFile,omitempty"`
	PrivateEnvFile     string `json:"privateEnvFile,omitempty"`
	SaveResponses      bool   `json:"saveResponses,omitempty"`
	ResponsesDir       string `json:"responsesDir,omitempty"`
}

// Manager handles reading and writing context files
type Manager struct {
	contextFile string
}

// NewManager creates a new context manager
// It looks for .postie-context.json in the current directory
func NewManager() *Manager {
	cwd, err := os.Getwd()
	if err != nil {
		cwd = "."
	}
	return &Manager{
		contextFile: filepath.Join(cwd, ".postie-context.json"),
	}
}

// NewManagerWithPath creates a context manager for a specific directory
func NewManagerWithPath(dir string) *Manager {
	return &Manager{
		contextFile: filepath.Join(dir, ".postie-context.json"),
	}
}

// Load reads the context from the context file
func (m *Manager) Load() (*Context, error) {
	data, err := os.ReadFile(m.contextFile)
	if err != nil {
		if os.IsNotExist(err) {
			return &Context{}, nil // Return empty context if file doesn't exist
		}
		return nil, fmt.Errorf("failed to read context file: %w", err)
	}

	var ctx Context
	if err := json.Unmarshal(data, &ctx); err != nil {
		return nil, fmt.Errorf("failed to parse context file: %w", err)
	}

	return &ctx, nil
}

// Save writes the context to the context file
func (m *Manager) Save(ctx *Context) error {
	data, err := json.MarshalIndent(ctx, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal context: %w", err)
	}

	if err := os.WriteFile(m.contextFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write context file: %w", err)
	}

	return nil
}

// Clear removes the context file
func (m *Manager) Clear() error {
	if err := os.Remove(m.contextFile); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove context file: %w", err)
	}
	return nil
}

// Exists checks if a context file exists
func (m *Manager) Exists() bool {
	_, err := os.Stat(m.contextFile)
	return err == nil
}

// GetPath returns the path to the context file
func (m *Manager) GetPath() string {
	return m.contextFile
}

// MergeWithFlags merges context values with command-line flags
// Flags take precedence over context values
func MergeWithFlags(ctx *Context, httpFile, env, envFile, privateEnvFile, responsesDir *string, saveResponses *bool) {
	if *httpFile == "" && ctx.HTTPFile != "" {
		*httpFile = ctx.HTTPFile
	}
	if *env == "" && ctx.Environment != "" {
		*env = ctx.Environment
	}
	if *envFile == "" && ctx.EnvFile != "" {
		*envFile = ctx.EnvFile
	}
	if *privateEnvFile == "" && ctx.PrivateEnvFile != "" {
		*privateEnvFile = ctx.PrivateEnvFile
	}
	if *responsesDir == "" && ctx.ResponsesDir != "" {
		*responsesDir = ctx.ResponsesDir
	}
	if !*saveResponses && ctx.SaveResponses {
		*saveResponses = ctx.SaveResponses
	}
}
