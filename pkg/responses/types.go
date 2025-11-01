package responses

import (
	"time"

	"postie/pkg/client"
	"postie/pkg/httprequest"
)

// StoredResponse represents a saved response with metadata
type StoredResponse struct {
	// Metadata
	RequestName string    `json:"request_name,omitempty"`
	RequestURL  string    `json:"request_url"`
	Method      string    `json:"method"`
	Timestamp   time.Time `json:"timestamp"`
	Duration    int64     `json:"duration_ms"` // Duration in milliseconds

	// Request details
	RequestHeaders map[string]string `json:"request_headers,omitempty"`
	RequestBody    string            `json:"request_body,omitempty"`

	// Response details
	StatusCode    int               `json:"status_code"`
	Status        string            `json:"status"`
	Headers       map[string]string `json:"headers"`
	Body          string            `json:"body"`
	ContentType   string            `json:"content_type"`
	ContentLength int64             `json:"content_length"`
}

// StorageConfig holds configuration for response storage
type StorageConfig struct {
	BaseDir          string // Base directory for storing responses
	UseRequestName   bool   // Organize by request name
	UseTimestamp     bool   // Include timestamp in filename
	MaxHistoryPerReq int    // Maximum number of responses to keep per request (0 = unlimited)
}

// DefaultStorageConfig returns the default storage configuration
func DefaultStorageConfig() *StorageConfig {
	return &StorageConfig{
		BaseDir:          ".http-responses",
		UseRequestName:   true,
		UseTimestamp:     true,
		MaxHistoryPerReq: 10,
	}
}

// ResponseHistory represents the history of responses for a request
type ResponseHistory struct {
	RequestName string           `json:"request_name"`
	RequestURL  string           `json:"request_url"`
	Responses   []HistoryEntry   `json:"responses"`
}

// HistoryEntry represents a single entry in response history
type HistoryEntry struct {
	Timestamp time.Time `json:"timestamp"`
	FilePath  string    `json:"file_path"`
	Status    string    `json:"status"`
	Duration  int64     `json:"duration_ms"`
}

// ResponseComparison represents a comparison between two responses
type ResponseComparison struct {
	Request1    *StoredResponse
	Request2    *StoredResponse
	StatusMatch bool
	BodyMatch   bool
	Differences []Difference
}

// Difference represents a difference between two responses
type Difference struct {
	Field    string      `json:"field"`
	Value1   interface{} `json:"value1"`
	Value2   interface{} `json:"value2"`
	DiffType string      `json:"diff_type"` // "added", "removed", "changed"
}

// FromClientResponse converts a client.Response to StoredResponse
func FromClientResponse(response *client.Response, request *httprequest.Request, duration time.Duration) (*StoredResponse, error) {
	body, err := response.Text()
	if err != nil {
		return nil, err
	}

	// Convert headers to map
	headers := make(map[string]string)
	for key, values := range response.Header {
		if len(values) > 0 {
			headers[key] = values[0]
		}
	}

	// Convert request headers
	reqHeaders := make(map[string]string)
	for _, h := range request.Headers {
		reqHeaders[h.Name] = h.Value
	}

	reqBody := ""
	if request.Body != nil {
		reqBody = request.Body.Content
	}

	return &StoredResponse{
		RequestName:    request.Name,
		RequestURL:     request.URL.Raw,
		Method:         request.Method,
		Timestamp:      time.Now(),
		Duration:       duration.Milliseconds(),
		RequestHeaders: reqHeaders,
		RequestBody:    reqBody,
		StatusCode:     response.StatusCode,
		Status:         response.Status,
		Headers:        headers,
		Body:           body,
		ContentType:    response.ContentType(),
		ContentLength:  response.ContentLength,
	}, nil
}
