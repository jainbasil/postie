package environment

import (
	"fmt"
	"sort"
)

// Merger handles merging of environment files and configurations
type Merger struct {
	resolver *Resolver
}

// NewMerger creates a new environment merger
func NewMerger() *Merger {
	return &Merger{
		resolver: NewResolver(),
	}
}

// NewMergerWithResolver creates a merger with a custom resolver
func NewMergerWithResolver(resolver *Resolver) *Merger {
	return &Merger{
		resolver: resolver,
	}
}

// MergeConfig represents configuration for environment merging
type MergeConfig struct {
	Environment          string // Target environment name
	AllowSystemVariables bool   // Allow system environment variables
	SystemVariablePrefix string // Prefix for system variables (empty = allow all)
	FailOnMissing        bool   // Fail if environment doesn't exist
	FailOnUnresolved     bool   // Fail if variables can't be resolved
}

// DefaultMergeConfig returns default merge configuration
func DefaultMergeConfig(environment string) *MergeConfig {
	return &MergeConfig{
		Environment:          environment,
		AllowSystemVariables: true,
		SystemVariablePrefix: "",
		FailOnMissing:        true,
		FailOnUnresolved:     true,
	}
}

// MergeEnvironments merges public and private environment files for a specific environment
func (m *Merger) MergeEnvironments(publicEnv, privateEnv EnvironmentFile, config *MergeConfig) (*ResolvedEnvironment, error) {
	if config == nil {
		config = DefaultMergeConfig("development")
	}

	// Configure resolver based on merge config
	if config.SystemVariablePrefix != "" {
		m.resolver = NewResolverWithPrefix(config.SystemVariablePrefix)
	}

	// Resolve the target environment
	resolved, err := m.resolver.Resolve(publicEnv, privateEnv, config.Environment)
	if err != nil {
		if config.FailOnMissing {
			return nil, fmt.Errorf("failed to resolve environment '%s': %w", config.Environment, err)
		}
		// Return empty environment if not failing on missing
		return &ResolvedEnvironment{
			Name:      config.Environment,
			Variables: make(map[string]interface{}),
			Source:    make(map[string]string),
		}, nil
	}

	// Validate resolution if required
	if config.FailOnUnresolved {
		if errors := m.resolver.ValidateResolution(resolved); len(errors) > 0 {
			return nil, fmt.Errorf("unresolved variables: %v", errors)
		}
	}

	return resolved, nil
}

// MergeMultipleEnvironments merges variables from multiple environments
func (m *Merger) MergeMultipleEnvironments(publicEnv, privateEnv EnvironmentFile, environments []string) (*ResolvedEnvironment, error) {
	var resolvedEnvs []*ResolvedEnvironment

	for _, envName := range environments {
		resolved, err := m.resolver.Resolve(publicEnv, privateEnv, envName)
		if err != nil {
			return nil, fmt.Errorf("failed to resolve environment '%s': %w", envName, err)
		}
		resolvedEnvs = append(resolvedEnvs, resolved)
	}

	return m.resolver.MergeEnvironments(resolvedEnvs...), nil
}

// GetEnvironmentInfo returns information about available environments
func (m *Merger) GetEnvironmentInfo(publicEnv, privateEnv EnvironmentFile) *EnvironmentInfo {
	info := &EnvironmentInfo{
		Environments: make(map[string]*EnvironmentDetails),
	}

	// Collect all environment names
	allEnvs := make(map[string]bool)
	for env := range publicEnv {
		allEnvs[env] = true
	}
	for env := range privateEnv {
		allEnvs[env] = true
	}

	// Build details for each environment
	for envName := range allEnvs {
		details := &EnvironmentDetails{
			Name:             envName,
			PublicVariables:  make([]string, 0),
			PrivateVariables: make([]string, 0),
			TotalVariables:   0,
		}

		// Count public variables
		if pubEnv, exists := publicEnv[envName]; exists {
			for varName := range pubEnv {
				details.PublicVariables = append(details.PublicVariables, varName)
			}
		}

		// Count private variables
		if privEnv, exists := privateEnv[envName]; exists {
			for varName := range privEnv {
				details.PrivateVariables = append(details.PrivateVariables, varName)
			}
		}

		// Sort variable lists
		sort.Strings(details.PublicVariables)
		sort.Strings(details.PrivateVariables)

		// Calculate total unique variables
		allVars := make(map[string]bool)
		for _, varName := range details.PublicVariables {
			allVars[varName] = true
		}
		for _, varName := range details.PrivateVariables {
			allVars[varName] = true
		}
		details.TotalVariables = len(allVars)

		info.Environments[envName] = details
	}

	return info
}

// EnvironmentInfo provides information about available environments
type EnvironmentInfo struct {
	Environments map[string]*EnvironmentDetails `json:"environments"`
}

// EnvironmentDetails provides details about a specific environment
type EnvironmentDetails struct {
	Name             string   `json:"name"`
	PublicVariables  []string `json:"public_variables"`
	PrivateVariables []string `json:"private_variables"`
	TotalVariables   int      `json:"total_variables"`
}

// GetEnvironmentNames returns sorted list of environment names
func (info *EnvironmentInfo) GetEnvironmentNames() []string {
	names := make([]string, 0, len(info.Environments))
	for name := range info.Environments {
		names = append(names, name)
	}
	sort.Strings(names)
	return names
}

// GetEnvironmentDetails returns details for a specific environment
func (info *EnvironmentInfo) GetEnvironmentDetails(name string) (*EnvironmentDetails, bool) {
	details, exists := info.Environments[name]
	return details, exists
}

// HasEnvironment checks if an environment exists
func (info *EnvironmentInfo) HasEnvironment(name string) bool {
	_, exists := info.Environments[name]
	return exists
}

// ValidateEnvironmentReferences validates that all required environments exist
func (m *Merger) ValidateEnvironmentReferences(publicEnv, privateEnv EnvironmentFile, requiredEnvs []string) []error {
	var errors []error

	availableEnvs := make(map[string]bool)
	for env := range publicEnv {
		availableEnvs[env] = true
	}
	for env := range privateEnv {
		availableEnvs[env] = true
	}

	for _, required := range requiredEnvs {
		if !availableEnvs[required] {
			errors = append(errors, fmt.Errorf("required environment '%s' not found", required))
		}
	}

	return errors
}

// DiffEnvironments compares two environments and returns differences
func (m *Merger) DiffEnvironments(publicEnv, privateEnv EnvironmentFile, env1, env2 string) (*EnvironmentDiff, error) {
	resolved1, err := m.resolver.Resolve(publicEnv, privateEnv, env1)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve environment '%s': %w", env1, err)
	}

	resolved2, err := m.resolver.Resolve(publicEnv, privateEnv, env2)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve environment '%s': %w", env2, err)
	}

	diff := &EnvironmentDiff{
		Environment1: env1,
		Environment2: env2,
		OnlyIn1:      make([]string, 0),
		OnlyIn2:      make([]string, 0),
		Different:    make([]VariableDiff, 0),
		Same:         make([]string, 0),
	}

	// Find variables only in env1
	for varName := range resolved1.Variables {
		if _, exists := resolved2.Variables[varName]; !exists {
			diff.OnlyIn1 = append(diff.OnlyIn1, varName)
		}
	}

	// Find variables only in env2
	for varName := range resolved2.Variables {
		if _, exists := resolved1.Variables[varName]; !exists {
			diff.OnlyIn2 = append(diff.OnlyIn2, varName)
		}
	}

	// Find different and same variables
	for varName, value1 := range resolved1.Variables {
		if value2, exists := resolved2.Variables[varName]; exists {
			if fmt.Sprintf("%v", value1) != fmt.Sprintf("%v", value2) {
				diff.Different = append(diff.Different, VariableDiff{
					Name:   varName,
					Value1: value1,
					Value2: value2,
				})
			} else {
				diff.Same = append(diff.Same, varName)
			}
		}
	}

	// Sort all lists
	sort.Strings(diff.OnlyIn1)
	sort.Strings(diff.OnlyIn2)
	sort.Strings(diff.Same)
	sort.Slice(diff.Different, func(i, j int) bool {
		return diff.Different[i].Name < diff.Different[j].Name
	})

	return diff, nil
}

// EnvironmentDiff represents differences between two environments
type EnvironmentDiff struct {
	Environment1 string         `json:"environment1"`
	Environment2 string         `json:"environment2"`
	OnlyIn1      []string       `json:"only_in_1"`
	OnlyIn2      []string       `json:"only_in_2"`
	Different    []VariableDiff `json:"different"`
	Same         []string       `json:"same"`
}

// VariableDiff represents a difference in variable values
type VariableDiff struct {
	Name   string      `json:"name"`
	Value1 interface{} `json:"value1"`
	Value2 interface{} `json:"value2"`
}
