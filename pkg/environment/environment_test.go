package environment

import (
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestLoaderBasic(t *testing.T) {
	// Create temporary directory for test files
	tmpDir := t.TempDir()

	// Create test environment file - use proper JSON formatting
	envContent := `{
	"development": {
		"baseUrl": "https://api-dev.example.com",
		"apiKey": "dev-key-123",
		"timeout": 30000
	},
	"production": {
		"baseUrl": "https://api.example.com",
		"apiKey": "prod-key-456"
	}
}`

	envFile := filepath.Join(tmpDir, "http-client.env.json")
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	// Test loading
	loader := NewLoader(tmpDir)
	config := &EnvironmentConfig{
		PublicFile:  envFile,
		Environment: "development",
	}

	publicEnv, privateEnv, err := loader.LoadEnvironments(config)
	if err != nil {
		t.Fatalf("Failed to load environments: %v", err)
	}

	// Verify public environment
	if len(*publicEnv) != 2 {
		t.Fatalf("Expected 2 environments, got %d", len(*publicEnv))
	}

	devEnv := (*publicEnv)["development"]
	if devEnv["baseUrl"] != "https://api-dev.example.com" {
		t.Errorf("Expected dev baseUrl, got %v", devEnv["baseUrl"])
	}

	// Verify private environment is empty
	if len(*privateEnv) != 0 {
		t.Errorf("Expected empty private environment, got %d variables", len(*privateEnv))
	}
}

func TestLoaderWithComments(t *testing.T) {
	tmpDir := t.TempDir()

	// Environment file with comments
	envContent := `{
	// Development environment
	"development": {
		"baseUrl": "https://api-dev.example.com", // Dev API URL
		"apiKey": "dev-key-123"
	},
	/* Production environment
	   with multi-line comment */
	"production": {
		"baseUrl": "https://api.example.com"
	}
}`

	envFile := filepath.Join(tmpDir, "http-client.env.json")
	err := os.WriteFile(envFile, []byte(envContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write test file: %v", err)
	}

	loader := NewLoader(tmpDir)
	config := &EnvironmentConfig{
		PublicFile:  envFile,
		Environment: "development",
	}

	publicEnv, _, err := loader.LoadEnvironments(config)
	if err != nil {
		t.Fatalf("Failed to load environments with comments: %v", err)
	}

	if len(*publicEnv) != 2 {
		t.Fatalf("Expected 2 environments, got %d", len(*publicEnv))
	}
}

func TestResolverBasic(t *testing.T) {
	publicEnv := EnvironmentFile{
		"development": Environment{
			"baseUrl": "https://api-dev.example.com",
			"apiKey":  "dev-key-123",
			"fullUrl": "{{baseUrl}}/v1",
		},
	}

	privateEnv := EnvironmentFile{
		"development": Environment{
			"secretKey": "secret-123",
		},
	}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(publicEnv, privateEnv, "development")
	if err != nil {
		t.Fatalf("Failed to resolve environment: %v", err)
	}

	// Check merged variables
	if resolved.GetString("baseUrl") != "https://api-dev.example.com" {
		t.Errorf("Expected baseUrl from public env")
	}

	if resolved.GetString("secretKey") != "secret-123" {
		t.Errorf("Expected secretKey from private env")
	}

	// Check variable resolution
	if resolved.GetString("fullUrl") != "https://api-dev.example.com/v1" {
		t.Errorf("Expected resolved fullUrl, got %s", resolved.GetString("fullUrl"))
	}
}

func TestResolverWithSystemEnv(t *testing.T) {
	// Set test system environment variable
	os.Setenv("TEST_VAR", "system-value")
	defer os.Unsetenv("TEST_VAR")

	publicEnv := EnvironmentFile{
		"test": Environment{
			"systemValue": "{{TEST_VAR}}",
			"combined":    "prefix-{{TEST_VAR}}-suffix",
		},
	}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(publicEnv, EnvironmentFile{}, "test")
	if err != nil {
		t.Fatalf("Failed to resolve environment: %v", err)
	}

	if resolved.GetString("systemValue") != "system-value" {
		t.Errorf("Expected system env var resolution, got %s", resolved.GetString("systemValue"))
	}

	if resolved.GetString("combined") != "prefix-system-value-suffix" {
		t.Errorf("Expected combined value, got %s", resolved.GetString("combined"))
	}
}

func TestResolverTypePreservation(t *testing.T) {
	publicEnv := EnvironmentFile{
		"test": Environment{
			"timeout":    30000,
			"debug":      true,
			"timeoutRef": "{{timeout}}",
			"debugRef":   "{{debug}}",
		},
	}

	resolver := NewResolver()
	resolved, err := resolver.Resolve(publicEnv, EnvironmentFile{}, "test")
	if err != nil {
		t.Fatalf("Failed to resolve environment: %v", err)
	}

	// Check that type is preserved for single variable references
	timeoutVar, exists := resolved.GetVariable("timeoutRef")
	if !exists {
		t.Fatal("timeoutRef variable not found")
	}

	if timeoutInt, err := timeoutVar.GetInt(); err != nil || timeoutInt != 30000 {
		t.Errorf("Expected int value 30000, got %v", timeoutVar.Value)
	}

	debugVar, exists := resolved.GetVariable("debugRef")
	if !exists {
		t.Fatal("debugRef variable not found")
	}

	if debugBool, err := debugVar.GetBool(); err != nil || !debugBool {
		t.Errorf("Expected bool value true, got %v", debugVar.Value)
	}
}

func TestMergerBasic(t *testing.T) {
	publicEnv := EnvironmentFile{
		"development": Environment{
			"baseUrl": "https://api-dev.example.com",
			"apiKey":  "dev-key-123",
		},
		"production": Environment{
			"baseUrl": "https://api.example.com",
			"apiKey":  "prod-key-456",
		},
	}

	privateEnv := EnvironmentFile{
		"development": Environment{
			"secretKey": "dev-secret",
		},
	}

	merger := NewMerger()
	config := DefaultMergeConfig("development")

	resolved, err := merger.MergeEnvironments(publicEnv, privateEnv, config)
	if err != nil {
		t.Fatalf("Failed to merge environments: %v", err)
	}

	if resolved.Name != "development" {
		t.Errorf("Expected environment name 'development', got %s", resolved.Name)
	}

	if len(resolved.Variables) != 3 {
		t.Errorf("Expected 3 variables, got %d", len(resolved.Variables))
	}
}

func TestMergerEnvironmentInfo(t *testing.T) {
	publicEnv := EnvironmentFile{
		"development": Environment{
			"baseUrl": "https://api-dev.example.com",
			"apiKey":  "dev-key-123",
		},
		"production": Environment{
			"baseUrl": "https://api.example.com",
		},
	}

	privateEnv := EnvironmentFile{
		"development": Environment{
			"secretKey": "dev-secret",
		},
		"staging": Environment{
			"stagingVar": "staging-value",
		},
	}

	merger := NewMerger()
	info := merger.GetEnvironmentInfo(publicEnv, privateEnv)

	envNames := info.GetEnvironmentNames()
	expectedNames := []string{"development", "production", "staging"}

	if !reflect.DeepEqual(envNames, expectedNames) {
		t.Errorf("Expected environments %v, got %v", expectedNames, envNames)
	}

	// Check development environment details
	devDetails, exists := info.GetEnvironmentDetails("development")
	if !exists {
		t.Fatal("Development environment details not found")
	}

	if devDetails.TotalVariables != 3 {
		t.Errorf("Expected 3 total variables in development, got %d", devDetails.TotalVariables)
	}

	if len(devDetails.PublicVariables) != 2 {
		t.Errorf("Expected 2 public variables, got %d", len(devDetails.PublicVariables))
	}

	if len(devDetails.PrivateVariables) != 1 {
		t.Errorf("Expected 1 private variable, got %d", len(devDetails.PrivateVariables))
	}
}

func TestValidation(t *testing.T) {
	tmpDir := t.TempDir()
	loader := NewLoader(tmpDir)

	// Test invalid environment file
	invalidEnv := EnvironmentFile{
		"": Environment{ // Empty environment name
			"var1": "value1",
		},
		"test": Environment{
			"": "value2", // Empty variable name
		},
	}

	errors := loader.ValidateEnvironmentFile(invalidEnv)
	if len(errors) == 0 {
		t.Error("Expected validation errors for invalid environment file")
	}

	// Should have errors for empty environment name and empty variable name
	found := make(map[string]bool)
	for _, err := range errors {
		errMsg := err.Error()
		if strings.Contains(errMsg, "environment name cannot be empty") {
			found["empty_env_name"] = true
		}
		if strings.Contains(errMsg, "variable name cannot be empty") {
			found["empty_var_name"] = true
		}
	}

	if !found["empty_env_name"] {
		t.Error("Expected error for empty environment name")
	}

	if !found["empty_var_name"] {
		t.Error("Expected error for empty variable name")
	}
}

func TestCircularReferenceDetection(t *testing.T) {
	publicEnv := EnvironmentFile{
		"test": Environment{
			"var1": "{{var2}}",
			"var2": "{{var1}}", // Circular reference
		},
	}

	resolver := NewResolver()
	_, err := resolver.Resolve(publicEnv, EnvironmentFile{}, "test")

	if err == nil {
		t.Error("Expected error for circular reference")
	}

	if !strings.Contains(err.Error(), "circular") {
		t.Errorf("Expected circular reference error, got: %v", err)
	}
}

func TestEnvironmentDiff(t *testing.T) {
	publicEnv := EnvironmentFile{
		"env1": Environment{
			"common":    "value1",
			"only1":     "unique1",
			"different": "value1",
		},
		"env2": Environment{
			"common":    "value1", // Same value
			"only2":     "unique2",
			"different": "value2", // Different value
		},
	}

	merger := NewMerger()
	diff, err := merger.DiffEnvironments(publicEnv, EnvironmentFile{}, "env1", "env2")
	if err != nil {
		t.Fatalf("Failed to diff environments: %v", err)
	}

	if len(diff.OnlyIn1) != 1 || diff.OnlyIn1[0] != "only1" {
		t.Errorf("Expected OnlyIn1 to contain 'only1', got %v", diff.OnlyIn1)
	}

	if len(diff.OnlyIn2) != 1 || diff.OnlyIn2[0] != "only2" {
		t.Errorf("Expected OnlyIn2 to contain 'only2', got %v", diff.OnlyIn2)
	}

	if len(diff.Same) != 1 || diff.Same[0] != "common" {
		t.Errorf("Expected Same to contain 'common', got %v", diff.Same)
	}

	if len(diff.Different) != 1 || diff.Different[0].Name != "different" {
		t.Errorf("Expected Different to contain 'different', got %v", diff.Different)
	}
}

func TestDiscoverEnvironmentFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create test files
	publicFile := filepath.Join(tmpDir, "http-client.env.json")
	privateFile := filepath.Join(tmpDir, "http-client.private.env.json")

	err := os.WriteFile(publicFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create public file: %v", err)
	}

	err = os.WriteFile(privateFile, []byte("{}"), 0644)
	if err != nil {
		t.Fatalf("Failed to create private file: %v", err)
	}

	loader := NewLoader(tmpDir)
	config := loader.DiscoverEnvironmentFiles()

	if config.PublicFile != publicFile {
		t.Errorf("Expected public file %s, got %s", publicFile, config.PublicFile)
	}

	if config.PrivateFile != privateFile {
		t.Errorf("Expected private file %s, got %s", privateFile, config.PrivateFile)
	}

	if config.Environment != "development" {
		t.Errorf("Expected default environment 'development', got %s", config.Environment)
	}
}

func TestVariableTypes(t *testing.T) {
	variable := &Variable{
		Name:  "testVar",
		Value: 42,
	}

	// Test GetInt
	intVal, err := variable.GetInt()
	if err != nil || intVal != 42 {
		t.Errorf("Expected int value 42, got %d, error: %v", intVal, err)
	}

	// Test GetString
	strVal := variable.GetString()
	if strVal != "42" {
		t.Errorf("Expected string value '42', got '%s'", strVal)
	}

	// Test boolean variable
	boolVar := &Variable{
		Name:  "boolVar",
		Value: true,
	}

	boolVal, err := boolVar.GetBool()
	if err != nil || !boolVal {
		t.Errorf("Expected bool value true, got %v, error: %v", boolVal, err)
	}
}

func TestJSONMarshaling(t *testing.T) {
	variable := &Variable{
		Name:        "testVar",
		Value:       "testValue",
		Source:      "public",
		Environment: "development",
	}

	data, err := json.Marshal(variable)
	if err != nil {
		t.Fatalf("Failed to marshal variable: %v", err)
	}

	var unmarshaled map[string]interface{}
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal variable: %v", err)
	}

	if unmarshaled["name"] != "testVar" {
		t.Errorf("Expected name 'testVar', got %v", unmarshaled["name"])
	}

	if unmarshaled["value"] != "testValue" {
		t.Errorf("Expected value 'testValue', got %v", unmarshaled["value"])
	}
}
