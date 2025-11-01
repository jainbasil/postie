package environment

import (
	"encoding/json"
	"fmt"
)

// Environment represents a single environment (dev, prod, test, etc.)
type Environment map[string]interface{}

// EnvironmentFile represents the structure of an http-client.env.json file
type EnvironmentFile map[string]Environment

// EnvironmentConfig holds configuration for environment loading
type EnvironmentConfig struct {
	PublicFile  string // Path to http-client.env.json
	PrivateFile string // Path to http-client.private.env.json
	Environment string // Active environment name (dev, prod, etc.)
}

// ResolvedEnvironment contains the merged environment variables
type ResolvedEnvironment struct {
	Name      string                 // Environment name
	Variables map[string]interface{} // Merged variables
	Source    map[string]string      // Variable source tracking (public/private/system)
}

// Variable represents a resolved environment variable
type Variable struct {
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Source      string      `json:"source"`      // "public", "private", "system"
	Environment string      `json:"environment"` // Environment name where defined
}

// VariableResolutionError occurs when variable resolution fails
type VariableResolutionError struct {
	Variable string
	Message  string
}

func (e *VariableResolutionError) Error() string {
	return fmt.Sprintf("variable resolution error for '%s': %s", e.Variable, e.Message)
}

// EnvironmentLoadError occurs when environment files cannot be loaded
type EnvironmentLoadError struct {
	File    string
	Message string
	Cause   error
}

func (e *EnvironmentLoadError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("failed to load environment file '%s': %s (%v)", e.File, e.Message, e.Cause)
	}
	return fmt.Sprintf("failed to load environment file '%s': %s", e.File, e.Message)
}

// String returns the string representation of an environment variable
func (v *Variable) String() string {
	return fmt.Sprintf("%v", v.Value)
}

// GetString returns the variable value as a string
func (v *Variable) GetString() string {
	if v.Value == nil {
		return ""
	}
	return fmt.Sprintf("%v", v.Value)
}

// GetInt returns the variable value as an integer
func (v *Variable) GetInt() (int, error) {
	switch val := v.Value.(type) {
	case int:
		return val, nil
	case float64:
		return int(val), nil
	case string:
		var result int
		_, err := fmt.Sscanf(val, "%d", &result)
		return result, err
	default:
		return 0, fmt.Errorf("cannot convert %T to int", v.Value)
	}
}

// GetBool returns the variable value as a boolean
func (v *Variable) GetBool() (bool, error) {
	switch val := v.Value.(type) {
	case bool:
		return val, nil
	case string:
		switch val {
		case "true", "TRUE", "True", "1":
			return true, nil
		case "false", "FALSE", "False", "0":
			return false, nil
		default:
			return false, fmt.Errorf("invalid boolean value: %s", val)
		}
	case int:
		return val != 0, nil
	case float64:
		return val != 0, nil
	default:
		return false, fmt.Errorf("cannot convert %T to bool", v.Value)
	}
}

// MarshalJSON provides custom JSON marshaling for Variable
func (v *Variable) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"name":        v.Name,
		"value":       v.Value,
		"source":      v.Source,
		"environment": v.Environment,
	})
}

// GetVariable returns a specific variable from the resolved environment
func (re *ResolvedEnvironment) GetVariable(name string) (*Variable, bool) {
	value, exists := re.Variables[name]
	if !exists {
		return nil, false
	}

	source := re.Source[name]
	if source == "" {
		source = "unknown"
	}

	return &Variable{
		Name:        name,
		Value:       value,
		Source:      source,
		Environment: re.Name,
	}, true
}

// GetString returns a variable value as a string
func (re *ResolvedEnvironment) GetString(name string) string {
	if variable, exists := re.GetVariable(name); exists {
		return variable.GetString()
	}
	return ""
}

// HasVariable checks if a variable exists in the environment
func (re *ResolvedEnvironment) HasVariable(name string) bool {
	_, exists := re.Variables[name]
	return exists
}

// ListVariables returns all variables in the environment
func (re *ResolvedEnvironment) ListVariables() []*Variable {
	variables := make([]*Variable, 0, len(re.Variables))
	for name := range re.Variables {
		if variable, exists := re.GetVariable(name); exists {
			variables = append(variables, variable)
		}
	}
	return variables
}
