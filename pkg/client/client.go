package client

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// Client represents an HTTP API client
type Client interface {
	// GET makes a GET request
	GET(url string) *Request
	// POST makes a POST request
	POST(url string) *Request
	// PUT makes a PUT request
	PUT(url string) *Request
	// DELETE makes a DELETE request
	DELETE(url string) *Request
	// PATCH makes a PATCH request
	PATCH(url string) *Request
	// HEAD makes a HEAD request
	HEAD(url string) *Request
	// OPTIONS makes an OPTIONS request
	OPTIONS(url string) *Request
}

// APIClient implements the Client interface
type APIClient struct {
	httpClient *http.Client
	baseURL    string
	headers    http.Header
	middleware []Middleware
}

// Middleware represents request/response middleware
type Middleware func(*http.Request, *http.Response) error

// Config holds client configuration
type Config struct {
	BaseURL    string
	Timeout    time.Duration
	Headers    map[string]string
	Middleware []Middleware
}

// NewClient creates a new API client
func NewClient(config *Config) *APIClient {
	if config == nil {
		config = &Config{}
	}

	// Use the timeout from config (0 means no timeout)
	timeout := config.Timeout

	client := &APIClient{
		httpClient: &http.Client{
			Timeout: timeout,
		},
		baseURL:    config.BaseURL,
		headers:    make(http.Header),
		middleware: config.Middleware,
	}

	// Set default headers
	for key, value := range config.Headers {
		client.headers.Set(key, value)
	}

	return client
}

// GET creates a GET request
func (c *APIClient) GET(url string) *Request {
	return c.newRequest(http.MethodGet, url)
}

// POST creates a POST request
func (c *APIClient) POST(url string) *Request {
	return c.newRequest(http.MethodPost, url)
}

// PUT creates a PUT request
func (c *APIClient) PUT(url string) *Request {
	return c.newRequest(http.MethodPut, url)
}

// DELETE creates a DELETE request
func (c *APIClient) DELETE(url string) *Request {
	return c.newRequest(http.MethodDelete, url)
}

// PATCH creates a PATCH request
func (c *APIClient) PATCH(url string) *Request {
	return c.newRequest(http.MethodPatch, url)
}

// HEAD creates a HEAD request
func (c *APIClient) HEAD(url string) *Request {
	return c.newRequest(http.MethodHead, url)
}

// OPTIONS creates an OPTIONS request
func (c *APIClient) OPTIONS(url string) *Request {
	return c.newRequest(http.MethodOptions, url)
}

// newRequest creates a new request builder
func (c *APIClient) newRequest(method, requestURL string) *Request {
	// Build full URL
	fullURL := c.buildURL(requestURL)

	return &Request{
		client: c,
		method: method,
		url:    fullURL,
		header: c.headers.Clone(),
		params: make(url.Values),
	}
}

// buildURL builds the full URL from base URL and request URL
func (c *APIClient) buildURL(requestURL string) string {
	if c.baseURL == "" {
		return requestURL
	}

	if strings.HasPrefix(requestURL, "http://") || strings.HasPrefix(requestURL, "https://") {
		return requestURL
	}

	baseURL := strings.TrimSuffix(c.baseURL, "/")
	requestURL = strings.TrimPrefix(requestURL, "/")

	return fmt.Sprintf("%s/%s", baseURL, requestURL)
}
