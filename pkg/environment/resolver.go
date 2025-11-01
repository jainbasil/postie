package environment

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

// Resolver handles variable resolution and environment merging
type Resolver struct {
	systemEnvPrefix string
}

// NewResolver creates a new environment resolver
func NewResolver() *Resolver {
	return &Resolver{
		systemEnvPrefix: "", // Allow all system environment variables
	}
}

// NewResolverWithPrefix creates a resolver that only uses system env vars with a prefix
func NewResolverWithPrefix(prefix string) *Resolver {
	return &Resolver{
		systemEnvPrefix: prefix,
	}
}

// Resolve merges public and private environments and resolves variables
func (r *Resolver) Resolve(publicEnv, privateEnv EnvironmentFile, envName string) (*ResolvedEnvironment, error) {
	// Check if environment exists
	publicVars, publicExists := publicEnv[envName]
	privateVars, privateExists := privateEnv[envName]

	if !publicExists && !privateExists {
		return nil, fmt.Errorf("environment '%s' not found", envName)
	}

	// Merge environments with precedence: private > public > system
	merged := make(map[string]interface{})
	sources := make(map[string]string)

	// Start with public environment
	if publicExists {
		for key, value := range publicVars {
			merged[key] = value
			sources[key] = "public"
		}
	}

	// Override with private environment
	if privateExists {
		for key, value := range privateVars {
			merged[key] = value
			sources[key] = "private"
		}
	}

	// Resolve variables (expand {{var}} references and system env vars)
	resolved, err := r.resolveVariables(merged, sources)
	if err != nil {
		return nil, fmt.Errorf("variable resolution failed: %w", err)
	}

	return &ResolvedEnvironment{
		Name:      envName,
		Variables: resolved,
		Source:    sources,
	}, nil
}

// resolveVariables resolves variable references and system environment variables
func (r *Resolver) resolveVariables(variables map[string]interface{}, sources map[string]string) (map[string]interface{}, error) {
	resolved := make(map[string]interface{})

	// Copy all variables for resolution
	for key, value := range variables {
		resolved[key] = value
	}

	// Resolve variables in multiple passes to handle nested references
	maxPasses := 10
	for pass := 0; pass < maxPasses; pass++ {
		changed := false

		for key, value := range resolved {
			if strValue, ok := value.(string); ok {
				newValue, hasChanges, err := r.resolveStringValue(strValue, resolved, sources)
				if err != nil {
					return nil, &VariableResolutionError{
						Variable: key,
						Message:  err.Error(),
					}
				}
				if hasChanges {
					resolved[key] = newValue
					changed = true
				}
			}
		}

		// If no changes in this pass, we're done
		if !changed {
			break
		}

		// If we've hit max passes, we might have circular references
		if pass == maxPasses-1 {
			return nil, fmt.Errorf("maximum resolution passes exceeded, possible circular references")
		}
	}

	return resolved, nil
}

// resolveStringValue resolves variable references in a string value
func (r *Resolver) resolveStringValue(value string, variables map[string]interface{}, sources map[string]string) (interface{}, bool, error) {
	// Pattern for {{variable}} references
	varPattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	original := value
	hasChanges := false

	// Replace variable references
	value = varPattern.ReplaceAllStringFunc(value, func(match string) string {
		// Extract variable name (remove {{ and }})
		varName := match[2 : len(match)-2]
		varName = strings.TrimSpace(varName)

		// First try to resolve from local variables
		if varValue, exists := variables[varName]; exists {
			hasChanges = true
			return fmt.Sprintf("%v", varValue)
		}

		// Then try system environment variables
		if r.isSystemEnvVar(varName) {
			if sysValue := os.Getenv(varName); sysValue != "" {
				hasChanges = true
				// Track as system variable
				sources[varName] = "system"
				return sysValue
			}
		}

		// Return unchanged if not found (for multi-pass resolution)
		return match
	})

	// If the entire value was a single variable reference, try to preserve type
	if singleVarPattern := regexp.MustCompile(`^\{\{([^}]+)\}\}$`); singleVarPattern.MatchString(original) {
		varName := strings.TrimSpace(original[2 : len(original)-2])
		if varValue, exists := variables[varName]; exists && fmt.Sprintf("%v", varValue) == value {
			// Return the actual value with its original type
			return varValue, hasChanges, nil
		}
	}

	return value, hasChanges, nil
}

// isSystemEnvVar checks if a variable name should be resolved from system environment
func (r *Resolver) isSystemEnvVar(varName string) bool {
	// If no prefix is set, allow all system env vars
	if r.systemEnvPrefix == "" {
		return true
	}

	// Check if variable name starts with the required prefix
	return strings.HasPrefix(varName, r.systemEnvPrefix)
}

// MergeEnvironments merges multiple environments with precedence
func (r *Resolver) MergeEnvironments(environments ...*ResolvedEnvironment) *ResolvedEnvironment {
	if len(environments) == 0 {
		return &ResolvedEnvironment{
			Name:      "empty",
			Variables: make(map[string]interface{}),
			Source:    make(map[string]string),
		}
	}

	if len(environments) == 1 {
		return environments[0]
	}

	// Use first environment as base
	merged := &ResolvedEnvironment{
		Name:      environments[0].Name,
		Variables: make(map[string]interface{}),
		Source:    make(map[string]string),
	}

	// Copy variables from first environment
	for key, value := range environments[0].Variables {
		merged.Variables[key] = value
		merged.Source[key] = environments[0].Source[key]
	}

	// Override with subsequent environments (later ones have higher precedence)
	for i := 1; i < len(environments); i++ {
		env := environments[i]
		for key, value := range env.Variables {
			merged.Variables[key] = value
			merged.Source[key] = env.Source[key]
		}
		// Update name to reflect the last merged environment
		merged.Name = fmt.Sprintf("%s+%s", merged.Name, env.Name)
	}

	return merged
}

// ValidateResolution checks if all variable references can be resolved
func (r *Resolver) ValidateResolution(resolved *ResolvedEnvironment) []error {
	var errors []error

	// Pattern for unresolved variable references
	unresolvedPattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	for varName, value := range resolved.Variables {
		if strValue, ok := value.(string); ok {
			// Find unresolved variable references
			matches := unresolvedPattern.FindAllStringSubmatch(strValue, -1)
			for _, match := range matches {
				unresolvedVar := strings.TrimSpace(match[1])
				errors = append(errors, fmt.Errorf("unresolved variable reference '{{%s}}' in variable '%s'", unresolvedVar, varName))
			}
		}
	}

	return errors
}

// ExpandString expands variable references in a string using the resolved environment
func (r *Resolver) ExpandString(input string, resolved *ResolvedEnvironment) string {
	varPattern := regexp.MustCompile(`\{\{([^}]+)\}\}`)

	return varPattern.ReplaceAllStringFunc(input, func(match string) string {
		varName := strings.TrimSpace(match[2 : len(match)-2])

		if variable, exists := resolved.GetVariable(varName); exists {
			return variable.GetString()
		}

		// Return unchanged if variable not found
		return match
	})
}
