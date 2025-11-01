package executor

import (
	"encoding/json"
	"fmt"
	"strings"

	"postie/pkg/scripting"
)

// Formatter handles formatting and display of execution results
type Formatter struct {
	verbose bool
	color   bool
}

// NewFormatter creates a new result formatter
func NewFormatter(verbose bool) *Formatter {
	return &Formatter{
		verbose: verbose,
		color:   true,
	}
}

// FormatResult formats an execution result for display
func (f *Formatter) FormatResult(result *ExecutionResult, index int) string {
	var output strings.Builder

	// Header
	output.WriteString(f.formatHeader(result, index))
	output.WriteString("\n")

	// Status
	output.WriteString(f.formatStatus(result))
	output.WriteString("\n")

	// Request details (if verbose)
	if f.verbose {
		output.WriteString(f.formatRequestDetails(result))
		output.WriteString("\n")
	}

	// Response body
	if result.Response != nil {
		output.WriteString(f.formatResponseBody(result))
	}

	// Error (if any)
	if result.Error != nil {
		output.WriteString(f.formatError(result))
	}

	// Script results (if any)
	if result.ScriptResult != nil {
		output.WriteString(f.formatScriptResults(result.ScriptResult))
	}

	// Response file path (if saved)
	if result.ResponseFilePath != "" {
		output.WriteString(fmt.Sprintf("\nResponse saved to: %s\n", result.ResponseFilePath))
	}

	return output.String()
}

// formatHeader formats the result header
func (f *Formatter) formatHeader(result *ExecutionResult, index int) string {
	var header strings.Builder

	header.WriteString(fmt.Sprintf("\n%s Request %d: %s %s %s\n",
		strings.Repeat("=", 10),
		index,
		result.Request.Method,
		result.Request.URL.Raw,
		strings.Repeat("=", 10)))

	if result.Request.Name != "" {
		header.WriteString(fmt.Sprintf("Name: %s\n", result.Request.Name))
	}

	return header.String()
}

// formatStatus formats the status information
func (f *Formatter) formatStatus(result *ExecutionResult) string {
	var status strings.Builder

	if result.Response != nil {
		statusIcon := "✓"
		if result.IsError() {
			statusIcon = "✗"
		}

		status.WriteString(fmt.Sprintf("%s Status: %s\n", statusIcon, result.Status))
		status.WriteString(fmt.Sprintf("  Duration: %v\n", result.Duration))
		status.WriteString(fmt.Sprintf("  Size: %d bytes\n", result.Response.Size()))

		contentType := result.Response.ContentType()
		if contentType != "" {
			status.WriteString(fmt.Sprintf("  Content-Type: %s\n", contentType))
		}
	}

	return status.String()
}

// formatRequestDetails formats detailed request information
func (f *Formatter) formatRequestDetails(result *ExecutionResult) string {
	var details strings.Builder

	details.WriteString("\nRequest Details:\n")
	details.WriteString(fmt.Sprintf("  Method: %s\n", result.Request.Method))
	details.WriteString(fmt.Sprintf("  URL: %s\n", result.Request.URL.Raw))

	// Headers
	if len(result.Request.Headers) > 0 {
		details.WriteString("  Headers:\n")
		for _, header := range result.Request.Headers {
			details.WriteString(fmt.Sprintf("    %s: %s\n", header.Name, header.Value))
		}
	}

	// Body
	if result.Request.Body != nil && result.Request.Body.Content != "" {
		details.WriteString("  Body:\n")
		// Indent body content
		bodyLines := strings.Split(result.Request.Body.Content, "\n")
		for _, line := range bodyLines {
			if line != "" {
				details.WriteString(fmt.Sprintf("    %s\n", line))
			}
		}
	}

	return details.String()
}

// formatResponseBody formats the response body
func (f *Formatter) formatResponseBody(result *ExecutionResult) string {
	var body strings.Builder

	text, err := result.Response.Text()
	if err != nil {
		body.WriteString(fmt.Sprintf("\nError reading response body: %v\n", err))
		return body.String()
	}

	if text == "" {
		body.WriteString("\nResponse: (empty)\n")
		return body.String()
	}

	body.WriteString("\nResponse Body:\n")

	// Try to format as JSON
	contentType := result.Response.ContentType()
	if strings.Contains(contentType, "json") || f.looksLikeJSON(text) {
		formatted := f.formatJSON(text)
		body.WriteString(formatted)
	} else {
		// Display as plain text
		if len(text) > 1000 && !f.verbose {
			body.WriteString(text[:1000])
			body.WriteString(fmt.Sprintf("\n... [Response truncated - %d total characters]\n", len(text)))
		} else {
			body.WriteString(text)
			body.WriteString("\n")
		}
	}

	return body.String()
}

// formatJSON tries to format text as pretty JSON
func (f *Formatter) formatJSON(text string) string {
	var jsonData interface{}
	if err := json.Unmarshal([]byte(text), &jsonData); err != nil {
		// Not valid JSON, return as-is
		return text + "\n"
	}

	prettyJSON, err := json.MarshalIndent(jsonData, "", "  ")
	if err != nil {
		return text + "\n"
	}

	return string(prettyJSON) + "\n"
}

// looksLikeJSON checks if text looks like JSON
func (f *Formatter) looksLikeJSON(text string) bool {
	trimmed := strings.TrimSpace(text)
	return (strings.HasPrefix(trimmed, "{") && strings.HasSuffix(trimmed, "}")) ||
		(strings.HasPrefix(trimmed, "[") && strings.HasSuffix(trimmed, "]"))
}

// formatError formats error information
func (f *Formatter) formatError(result *ExecutionResult) string {
	return fmt.Sprintf("\n✗ Error: %v\n", result.Error)
}

// formatScriptResults formats response handler script execution results
func (f *Formatter) formatScriptResults(scriptResult *scripting.ScriptExecutionResult) string {
	var output strings.Builder

	if scriptResult == nil {
		return ""
	}

	output.WriteString("\nResponse Handler Results:\n")

	// Format script execution error
	if scriptResult.Error != nil {
		output.WriteString(fmt.Sprintf("  ✗ Script Error: %v\n", scriptResult.Error))
		return output.String()
	}

	// Format test results
	if len(scriptResult.Tests) > 0 {
		output.WriteString("\n  Tests:\n")
		for _, test := range scriptResult.Tests {
			icon := "✓"
			if !test.Passed {
				icon = "✗"
			}
			output.WriteString(fmt.Sprintf("    %s %s", icon, test.Name))
			if !test.Passed && test.Error != "" {
				output.WriteString(fmt.Sprintf(" - %s", test.Error))
			}
			output.WriteString("\n")
		}
	}

	// Format failed assertions
	if len(scriptResult.Assertions) > 0 {
		output.WriteString("\n  Assertions:\n")
		for _, assertion := range scriptResult.Assertions {
			output.WriteString(fmt.Sprintf("    ✗ %s\n", assertion.Message))
		}
	}

	// Format logs
	if len(scriptResult.Logs) > 0 {
		output.WriteString("\n  Logs:\n")
		for _, log := range scriptResult.Logs {
			output.WriteString(fmt.Sprintf("    %s\n", log))
		}
	}

	// Format globals (if verbose)
	if f.verbose && len(scriptResult.Globals) > 0 {
		output.WriteString("\n  Global Variables Set:\n")
		for name, value := range scriptResult.Globals {
			output.WriteString(fmt.Sprintf("    %s = %v\n", name, value))
		}
	}

	return output.String()
}

// FormatSummary formats a summary of multiple results
func (f *Formatter) FormatSummary(results []*ExecutionResult) string {
	var summary strings.Builder

	successCount := 0
	errorCount := 0
	failureCount := 0

	for _, result := range results {
		if result.HasError() {
			errorCount++
		} else if result.IsSuccess() {
			successCount++
		} else if result.IsError() {
			failureCount++
		}
	}

	summary.WriteString(fmt.Sprintf("\n%s Execution Summary %s\n", strings.Repeat("=", 20), strings.Repeat("=", 20)))
	summary.WriteString(fmt.Sprintf("Total Requests: %d\n", len(results)))
	summary.WriteString(fmt.Sprintf("✓ Successful: %d\n", successCount))
	summary.WriteString(fmt.Sprintf("✗ Failed: %d\n", failureCount))
	if errorCount > 0 {
		summary.WriteString(fmt.Sprintf("⚠ Errors: %d\n", errorCount))
	}

	return summary.String()
}
