package scripting

import (
	"postie/pkg/client"
	"postie/pkg/httprequest"
)

// ScriptContext contains the context for script execution
type ScriptContext struct {
	Request  *httprequest.Request
	Response *client.Response
	Env      map[string]interface{} // Environment variables
	Globals  *GlobalStore           // Global variables (persist across requests)
}

// TestResult represents the result of a client.test() call
type TestResult struct {
	Name   string
	Passed bool
	Error  string
	Line   int
	Column int
}

// AssertionError represents a failed assertion
type AssertionError struct {
	Message  string
	Expected interface{}
	Actual   interface{}
	Line     int
	Column   int
}

func (e *AssertionError) Error() string {
	return e.Message
}

// ScriptExecutionResult contains the results of script execution
type ScriptExecutionResult struct {
	Tests      []*TestResult
	Assertions []*AssertionError
	Logs       []string
	Globals    map[string]interface{}
	Error      error
}

// IsSuccess returns true if all tests passed and no errors occurred
func (r *ScriptExecutionResult) IsSuccess() bool {
	if r.Error != nil {
		return false
	}

	for _, test := range r.Tests {
		if !test.Passed {
			return false
		}
	}

	return len(r.Assertions) == 0
}

// HasTests returns true if any tests were executed
func (r *ScriptExecutionResult) HasTests() bool {
	return len(r.Tests) > 0
}

// HasAssertions returns true if any assertions failed
func (r *ScriptExecutionResult) HasAssertions() bool {
	return len(r.Assertions) > 0
}
