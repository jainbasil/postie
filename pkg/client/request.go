package client

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Request represents an HTTP request builder
type Request struct {
	client *APIClient
	method string
	url    string
	header http.Header
	params url.Values
	body   io.Reader
	ctx    context.Context
}

// Header sets a request header
func (r *Request) Header(key, value string) *Request {
	r.header.Set(key, value)
	return r
}

// Headers sets multiple headers
func (r *Request) Headers(headers map[string]string) *Request {
	for key, value := range headers {
		r.header.Set(key, value)
	}
	return r
}

// Param sets a URL parameter
func (r *Request) Param(key, value string) *Request {
	r.params.Set(key, value)
	return r
}

// Params sets multiple URL parameters
func (r *Request) Params(params map[string]string) *Request {
	for key, value := range params {
		r.params.Set(key, value)
	}
	return r
}

// JSON sets the request body as JSON
func (r *Request) JSON(data interface{}) *Request {
	jsonData, err := json.Marshal(data)
	if err != nil {
		// Handle error appropriately in production
		return r
	}

	r.body = bytes.NewReader(jsonData)
	r.header.Set("Content-Type", "application/json")
	return r
}

// Body sets the request body
func (r *Request) Body(body io.Reader) *Request {
	r.body = body
	return r
}

// Text sets the request body as plain text
func (r *Request) Text(text string) *Request {
	r.body = strings.NewReader(text)
	r.header.Set("Content-Type", "text/plain")
	return r
}

// Form sets the request body as form data
func (r *Request) Form(data map[string]string) *Request {
	form := url.Values{}
	for key, value := range data {
		form.Set(key, value)
	}

	r.body = strings.NewReader(form.Encode())
	r.header.Set("Content-Type", "application/x-www-form-urlencoded")
	return r
}

// Context sets the request context
func (r *Request) Context(ctx context.Context) *Request {
	r.ctx = ctx
	return r
}

// Execute sends the HTTP request and returns the response
func (r *Request) Execute() (*Response, error) {
	// Build URL with parameters
	finalURL := r.url
	if len(r.params) > 0 {
		if strings.Contains(finalURL, "?") {
			finalURL += "&" + r.params.Encode()
		} else {
			finalURL += "?" + r.params.Encode()
		}
	}

	// Create HTTP request
	req, err := http.NewRequest(r.method, finalURL, r.body)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header = r.header

	// Set context if provided
	if r.ctx != nil {
		req = req.WithContext(r.ctx)
	}

	// Execute request
	start := time.Now()
	resp, err := r.client.httpClient.Do(req)
	duration := time.Since(start)

	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}

	// Create response wrapper
	response := &Response{
		Response: resp,
		Duration: duration,
	}

	// Apply middleware
	for _, middleware := range r.client.middleware {
		if err := middleware(req, resp); err != nil {
			return response, fmt.Errorf("middleware error: %w", err)
		}
	}

	return response, nil
}

// Send is an alias for Execute
func (r *Request) Send() (*Response, error) {
	return r.Execute()
}
