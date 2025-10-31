package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Response wraps http.Response with additional functionality
type Response struct {
	*http.Response
	Duration time.Duration
	body     []byte
}

// GetBody returns the response body as bytes
func (r *Response) GetBody() ([]byte, error) {
	if r.body != nil {
		return r.body, nil
	}

	defer r.Response.Body.Close()
	body, err := io.ReadAll(r.Response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	r.body = body
	return body, nil
}

// Text returns the response body as a string
func (r *Response) Text() (string, error) {
	body, err := r.GetBody()
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// JSON unmarshals the response body into the provided interface
func (r *Response) JSON(v interface{}) error {
	body, err := r.GetBody()
	if err != nil {
		return err
	}

	if err := json.Unmarshal(body, v); err != nil {
		return fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return nil
}

// IsSuccess returns true if the status code is 2xx
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsClientError returns true if the status code is 4xx
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError returns true if the status code is 5xx
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// IsError returns true if the status code is 4xx or 5xx
func (r *Response) IsError() bool {
	return r.IsClientError() || r.IsServerError()
}

// Size returns the size of the response body in bytes
func (r *Response) Size() int64 {
	if r.body != nil {
		return int64(len(r.body))
	}
	return r.ContentLength
}

// ContentType returns the content type of the response
func (r *Response) ContentType() string {
	return r.Header.Get("Content-Type")
}

// String returns a string representation of the response
func (r *Response) String() string {
	body, _ := r.Text()
	return fmt.Sprintf("Status: %s\nDuration: %v\nSize: %d bytes\nBody: %s",
		r.Status, r.Duration, r.Size(), body)
}
