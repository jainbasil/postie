package auth

import (
	"encoding/base64"
	"fmt"
	"net/http"
)

// Authenticator interface for different authentication methods
type Authenticator interface {
	Apply(req *http.Request) error
}

// NoAuth represents no authentication
type NoAuth struct{}

func (a *NoAuth) Apply(req *http.Request) error {
	return nil
}

// APIKeyAuth represents API key authentication
type APIKeyAuth struct {
	Key   string
	Value string
	In    string // "header", "query"
}

func NewAPIKeyAuth(key, value, in string) *APIKeyAuth {
	return &APIKeyAuth{
		Key:   key,
		Value: value,
		In:    in,
	}
}

func (a *APIKeyAuth) Apply(req *http.Request) error {
	switch a.In {
	case "header":
		req.Header.Set(a.Key, a.Value)
	case "query":
		q := req.URL.Query()
		q.Set(a.Key, a.Value)
		req.URL.RawQuery = q.Encode()
	default:
		return fmt.Errorf("unsupported API key location: %s", a.In)
	}
	return nil
}

// BearerTokenAuth represents Bearer token authentication
type BearerTokenAuth struct {
	Token string
}

func NewBearerTokenAuth(token string) *BearerTokenAuth {
	return &BearerTokenAuth{Token: token}
}

func (a *BearerTokenAuth) Apply(req *http.Request) error {
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", a.Token))
	return nil
}

// BasicAuth represents HTTP Basic authentication
type BasicAuth struct {
	Username string
	Password string
}

func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

func (a *BasicAuth) Apply(req *http.Request) error {
	auth := base64.StdEncoding.EncodeToString([]byte(a.Username + ":" + a.Password))
	req.Header.Set("Authorization", fmt.Sprintf("Basic %s", auth))
	return nil
}

// CustomHeaderAuth represents custom header authentication
type CustomHeaderAuth struct {
	Header string
	Value  string
}

func NewCustomHeaderAuth(header, value string) *CustomHeaderAuth {
	return &CustomHeaderAuth{
		Header: header,
		Value:  value,
	}
}

func (a *CustomHeaderAuth) Apply(req *http.Request) error {
	req.Header.Set(a.Header, a.Value)
	return nil
}
