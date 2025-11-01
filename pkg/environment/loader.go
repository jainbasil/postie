package environment

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Loader handles loading and parsing environment files
type Loader struct {
	workingDir string
}

// NewLoader creates a new environment loader
func NewLoader(workingDir string) *Loader {
	return &Loader{
		workingDir: workingDir,
	}
}

// LoadEnvironments loads both public and private environment files
func (l *Loader) LoadEnvironments(config *EnvironmentConfig) (*EnvironmentFile, *EnvironmentFile, error) {
	publicEnv, err := l.loadEnvironmentFile(config.PublicFile)
	if err != nil {
		return nil, nil, &EnvironmentLoadError{
			File:    config.PublicFile,
			Message: "failed to load public environment file",
			Cause:   err,
		}
	}

	// Private environment file is optional
	privateEnv := make(EnvironmentFile)
	if config.PrivateFile != "" {
		private, err := l.loadEnvironmentFile(config.PrivateFile)
		if err != nil && !os.IsNotExist(err) {
			return nil, nil, &EnvironmentLoadError{
				File:    config.PrivateFile,
				Message: "failed to load private environment file",
				Cause:   err,
			}
		}
		if private != nil {
			privateEnv = private
		}
	}

	return &publicEnv, &privateEnv, nil
}

// loadEnvironmentFile loads a single environment file
func (l *Loader) loadEnvironmentFile(filename string) (EnvironmentFile, error) {
	if filename == "" {
		return make(EnvironmentFile), nil
	}

	// Make path relative to working directory if not absolute
	fullPath := filename
	if !filepath.IsAbs(filename) {
		fullPath = filepath.Join(l.workingDir, filename)
	}

	// Check if file exists
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return make(EnvironmentFile), nil
	}

	// Read file content
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	// Parse JSON with comments support (strip comments first)
	cleanContent := l.stripJSONComments(string(content))

	var envFile EnvironmentFile
	if err := json.Unmarshal([]byte(cleanContent), &envFile); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return envFile, nil
}

// stripJSONComments removes JSON comments from content
// This supports both // line comments and /* block comments */
func (l *Loader) stripJSONComments(content string) string {
	lines := strings.Split(content, "\n")
	var cleanLines []string

	inBlockComment := false
	for _, line := range lines {
		cleanLine := l.processLine(line, &inBlockComment)
		// Always add the line (even if empty) to preserve line structure
		cleanLines = append(cleanLines, cleanLine)
	}

	return strings.Join(cleanLines, "\n")
}

// processLine processes a single line to remove comments
// This version is aware of JSON string literals and won't treat // inside strings as comments
func (l *Loader) processLine(line string, inBlockComment *bool) string {
	if *inBlockComment {
		// Look for end of block comment
		if endIndex := strings.Index(line, "*/"); endIndex != -1 {
			*inBlockComment = false
			return l.processLine(line[endIndex+2:], inBlockComment)
		}
		return "" // Entire line is in block comment
	}

	result := make([]rune, 0, len(line))
	runes := []rune(line)
	inString := false
	escaped := false

	for i := 0; i < len(runes); i++ {
		char := runes[i]

		if inString {
			// Inside a string literal
			if escaped {
				escaped = false
			} else if char == '\\' {
				escaped = true
			} else if char == '"' {
				inString = false
			}
			result = append(result, char)
		} else {
			// Outside string literal
			if char == '"' {
				inString = true
				result = append(result, char)
			} else if char == '/' && i+1 < len(runes) {
				if runes[i+1] == '/' {
					// Line comment found - stop processing this line
					break
				} else if runes[i+1] == '*' {
					// Block comment starts
					*inBlockComment = true
					i++ // Skip the '*'
					// Look for end of block comment on same line
					for j := i + 1; j < len(runes)-1; j++ {
						if runes[j] == '*' && runes[j+1] == '/' {
							*inBlockComment = false
							i = j + 1 // Skip to after the */
							break
						}
					}
					// If block comment didn't end on this line, we're done with this line
					if *inBlockComment {
						break
					}
				} else {
					result = append(result, char)
				}
			} else {
				result = append(result, char)
			}
		}
	}

	return strings.TrimRight(string(result), " \t")
}

// DiscoverEnvironmentFiles discovers standard JetBrains environment files
func (l *Loader) DiscoverEnvironmentFiles() *EnvironmentConfig {
	publicFile := filepath.Join(l.workingDir, "http-client.env.json")
	privateFile := filepath.Join(l.workingDir, "http-client.private.env.json")

	// Check if public file exists
	if _, err := os.Stat(publicFile); os.IsNotExist(err) {
		publicFile = ""
	}

	// Check if private file exists
	if _, err := os.Stat(privateFile); os.IsNotExist(err) {
		privateFile = ""
	}

	return &EnvironmentConfig{
		PublicFile:  publicFile,
		PrivateFile: privateFile,
		Environment: "development", // Default environment
	}
}

// ValidateEnvironmentFile validates the structure of an environment file
func (l *Loader) ValidateEnvironmentFile(envFile EnvironmentFile) []error {
	var errors []error

	if len(envFile) == 0 {
		errors = append(errors, fmt.Errorf("environment file is empty"))
		return errors
	}

	for envName, env := range envFile {
		if envName == "" {
			errors = append(errors, fmt.Errorf("environment name cannot be empty"))
			continue
		}

		if env == nil {
			errors = append(errors, fmt.Errorf("environment '%s' is null", envName))
			continue
		}

		// Validate variable names
		for varName, value := range env {
			if varName == "" {
				errors = append(errors, fmt.Errorf("variable name cannot be empty in environment '%s'", envName))
				continue
			}

			// Check for invalid variable types
			if !l.isValidVariableType(value) {
				errors = append(errors, fmt.Errorf("invalid variable type for '%s' in environment '%s': %T", varName, envName, value))
			}
		}
	}

	return errors
}

// isValidVariableType checks if a variable value is of a valid type
func (l *Loader) isValidVariableType(value interface{}) bool {
	switch value.(type) {
	case string, int, float64, bool, nil:
		return true
	default:
		return false
	}
}

// GetAvailableEnvironments returns a list of available environment names
func (l *Loader) GetAvailableEnvironments(publicEnv, privateEnv EnvironmentFile) []string {
	envSet := make(map[string]bool)

	// Collect from public environments
	for envName := range publicEnv {
		envSet[envName] = true
	}

	// Collect from private environments
	for envName := range privateEnv {
		envSet[envName] = true
	}

	// Convert to sorted slice
	environments := make([]string, 0, len(envSet))
	for env := range envSet {
		environments = append(environments, env)
	}

	return environments
}
