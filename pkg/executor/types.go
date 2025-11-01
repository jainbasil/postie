package executor

import (
	"time"

	"postie/pkg/client"
	"postie/pkg/httprequest"
)

// ExecutionResult contains the result of executing an HTTP request
type ExecutionResult struct {
	// Request is the expanded request that was executed
	Request *httprequest.Request

	// Response is the HTTP response (nil if error occurred)
	Response *client.Response

	// Error is any error that occurred during execution
	Error error

	// Duration is how long the request took to execute
	Duration time.Duration

	// StatusCode is the HTTP status code
	StatusCode int

	// Status is the HTTP status text
	Status string
}

// IsSuccess returns true if the request was successful (2xx status code)
func (r *ExecutionResult) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError returns true if the request resulted in an error status (4xx or 5xx)
func (r *ExecutionResult) IsError() bool {
	return r.StatusCode >= 400
}

// HasError returns true if there was an execution error
func (r *ExecutionResult) HasError() bool {
	return r.Error != nil
}
