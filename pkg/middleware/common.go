package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

// LoggingMiddleware logs request and response details
func LoggingMiddleware(req *http.Request, resp *http.Response) error {
	log.Printf("%s %s - %d %s",
		req.Method,
		req.URL.String(),
		resp.StatusCode,
		resp.Status)
	return nil
}

// UserAgentMiddleware sets a custom User-Agent header
func UserAgentMiddleware(userAgent string) func(*http.Request, *http.Response) error {
	return func(req *http.Request, resp *http.Response) error {
		req.Header.Set("User-Agent", userAgent)
		return nil
	}
}

// RetryMiddleware provides retry functionality
func RetryMiddleware(maxRetries int, retryDelay time.Duration) func(*http.Request, *http.Response) error {
	return func(req *http.Request, resp *http.Response) error {
		// This is a simplified implementation
		// In practice, you'd need to handle retries at the client level
		if resp.StatusCode >= 500 && maxRetries > 0 {
			log.Printf("Server error %d, retries remaining: %d", resp.StatusCode, maxRetries)
		}
		return nil
	}
}

// RateLimitMiddleware implements basic rate limiting
type RateLimiter struct {
	lastRequest time.Time
	minInterval time.Duration
}

func NewRateLimitMiddleware(requestsPerSecond float64) func(*http.Request, *http.Response) error {
	limiter := &RateLimiter{
		minInterval: time.Duration(float64(time.Second) / requestsPerSecond),
	}

	return func(req *http.Request, resp *http.Response) error {
		now := time.Now()
		if !limiter.lastRequest.IsZero() {
			elapsed := now.Sub(limiter.lastRequest)
			if elapsed < limiter.minInterval {
				time.Sleep(limiter.minInterval - elapsed)
			}
		}
		limiter.lastRequest = time.Now()
		return nil
	}
}

// ErrorHandlingMiddleware provides custom error handling
func ErrorHandlingMiddleware(req *http.Request, resp *http.Response) error {
	if resp.StatusCode >= 400 {
		return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}
	return nil
}

// TimeoutMiddleware sets request timeout
func TimeoutMiddleware(timeout time.Duration) func(*http.Request, *http.Response) error {
	return func(req *http.Request, resp *http.Response) error {
		// This should be handled at the client level with context
		// This is just for demonstration
		log.Printf("Request timeout set to %v", timeout)
		return nil
	}
}
