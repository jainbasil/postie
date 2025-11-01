package responses

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Storage handles saving and loading responses
type Storage struct {
	config *StorageConfig
}

// NewStorage creates a new response storage
func NewStorage(config *StorageConfig) *Storage {
	if config == nil {
		config = DefaultStorageConfig()
	}
	return &Storage{
		config: config,
	}
}

// Save saves a response to disk
func (s *Storage) Save(response *StoredResponse) (string, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(s.config.BaseDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create base directory: %w", err)
	}

	// Generate file path
	filePath := s.generateFilePath(response)

	// Ensure directory exists
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("failed to create response directory: %w", err)
	}

	// Marshal response to JSON
	data, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal response: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write response file: %w", err)
	}

	return filePath, nil
}

// Load loads a response from disk
func (s *Storage) Load(filePath string) (*StoredResponse, error) {
	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read response file: %w", err)
	}

	// Unmarshal JSON
	var response StoredResponse
	if err := json.Unmarshal(data, &response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &response, nil
}

// generateFilePath generates the file path for storing a response
func (s *Storage) generateFilePath(response *StoredResponse) string {
	parts := []string{s.config.BaseDir}

	// Add request name directory if configured
	if s.config.UseRequestName && response.RequestName != "" {
		// Sanitize request name for filesystem
		safeName := sanitizeFilename(response.RequestName)
		parts = append(parts, safeName)
	} else {
		// Use method and URL-based path
		parts = append(parts, strings.ToLower(response.Method))
	}

	// Generate filename
	filename := s.generateFilename(response)
	parts = append(parts, filename)

	return filepath.Join(parts...)
}

// generateFilename generates the filename for a response
func (s *Storage) generateFilename(response *StoredResponse) string {
	var parts []string

	// Add timestamp if configured
	if s.config.UseTimestamp {
		timestamp := response.Timestamp.Format("2006-01-02T150405")
		parts = append(parts, timestamp)
	}

	// Add status code
	parts = append(parts, fmt.Sprintf("%d", response.StatusCode))

	// Join parts and add extension
	filename := strings.Join(parts, ".")
	return filename + ".json"
}

// GetHistory returns the response history for a request
func (s *Storage) GetHistory(requestName string) (*ResponseHistory, error) {
	if requestName == "" {
		return nil, fmt.Errorf("request name is required")
	}

	safeName := sanitizeFilename(requestName)
	dir := filepath.Join(s.config.BaseDir, safeName)

	// Check if directory exists
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		return &ResponseHistory{
			RequestName: requestName,
			Responses:   []HistoryEntry{},
		}, nil
	}

	// Read directory
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("failed to read response directory: %w", err)
	}

	// Build history
	history := &ResponseHistory{
		RequestName: requestName,
		Responses:   make([]HistoryEntry, 0, len(entries)),
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".json") {
			continue
		}

		filePath := filepath.Join(dir, entry.Name())

		// Load response to get metadata
		response, err := s.Load(filePath)
		if err != nil {
			continue // Skip invalid files
		}

		history.Responses = append(history.Responses, HistoryEntry{
			Timestamp: response.Timestamp,
			FilePath:  filePath,
			Status:    response.Status,
			Duration:  response.Duration,
		})
	}

	return history, nil
}

// CleanupHistory removes old responses beyond the configured limit
func (s *Storage) CleanupHistory(requestName string) error {
	if s.config.MaxHistoryPerReq <= 0 {
		return nil // Unlimited history
	}

	history, err := s.GetHistory(requestName)
	if err != nil {
		return err
	}

	if len(history.Responses) <= s.config.MaxHistoryPerReq {
		return nil // Within limit
	}

	// Sort by timestamp (oldest first) and remove excess
	// For simplicity, we'll remove the oldest files
	toRemove := len(history.Responses) - s.config.MaxHistoryPerReq

	for i := 0; i < toRemove; i++ {
		if err := os.Remove(history.Responses[i].FilePath); err != nil {
			return fmt.Errorf("failed to remove old response: %w", err)
		}
	}

	return nil
}

// sanitizeFilename removes unsafe characters from filenames
func sanitizeFilename(name string) string {
	// Replace unsafe characters with underscores
	unsafe := []string{"/", "\\", ":", "*", "?", "\"", "<", ">", "|", " "}
	result := name
	for _, char := range unsafe {
		result = strings.ReplaceAll(result, char, "_")
	}
	return result
}

// List returns all stored responses
func (s *Storage) List() ([]*StoredResponse, error) {
	var responses []*StoredResponse

	err := filepath.Walk(s.config.BaseDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() || !strings.HasSuffix(path, ".json") {
			return nil
		}

		response, err := s.Load(path)
		if err != nil {
			return nil // Skip invalid files
		}

		responses = append(responses, response)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to list responses: %w", err)
	}

	return responses, nil
}
