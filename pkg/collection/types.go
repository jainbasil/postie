package collection

// Collection represents the root collection structure
type Collection struct {
	Collection CollectionInfo `json:"collection"`
}

// CollectionInfo contains the main collection data
type CollectionInfo struct {
	Info        Info          `json:"info"`
	Variable    []Variable    `json:"variable,omitempty"`
	Environment []Environment `json:"environment,omitempty"`
	Auth        *Auth         `json:"auth,omitempty"`
	Event       []Event       `json:"event,omitempty"`
	ApiGroup    []ApiGroup    `json:"apiGroup"`
}

// Info contains collection metadata
type Info struct {
	Name          string `json:"name"`
	Description   string `json:"description,omitempty"`
	Version       string `json:"version,omitempty"`
	Schema        string `json:"schema,omitempty"`
	Author        string `json:"author,omitempty"`
	License       string `json:"license,omitempty"`
	Documentation string `json:"documentation,omitempty"`
}

// Variable represents a collection or environment variable
type Variable struct {
	Key         string      `json:"key"`
	Value       interface{} `json:"value"`
	Type        string      `json:"type,omitempty"`
	Description string      `json:"description,omitempty"`
	Enabled     bool        `json:"enabled,omitempty"`
}

// Environment represents an environment configuration
type Environment struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Values      []Variable `json:"values"`
	Auth        *Auth      `json:"auth,omitempty"`
	Event       []Event    `json:"event,omitempty"`
}

// Auth represents authentication configuration
type Auth struct {
	Type   string                 `json:"type"`
	NoAuth []AuthParam            `json:"noauth,omitempty"`
	Bearer []AuthParam            `json:"bearer,omitempty"`
	APIKey []AuthParam            `json:"apikey,omitempty"`
	Basic  []AuthParam            `json:"basic,omitempty"`
	OAuth2 []AuthParam            `json:"oauth2,omitempty"`
	Custom map[string]interface{} `json:",omitempty"`
}

// AuthParam represents an authentication parameter
type AuthParam struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// Event represents a script event
type Event struct {
	Listen string `json:"listen"`
	Script Script `json:"script"`
}

// Script represents a script configuration
type Script struct {
	Type string   `json:"type"`
	Exec []string `json:"exec"`
}

// Item represents a request or folder
type Item struct {
	ID          string        `json:"id"`
	Name        string        `json:"name"`
	Description string        `json:"description,omitempty"`
	Request     *Request      `json:"request,omitempty"`
	Response    []Response    `json:"response,omitempty"`
	Apis        []Item        `json:"apis,omitempty"`
	Environment []Environment `json:"environment,omitempty"`
	Auth        *Auth         `json:"auth,omitempty"`
	Event       []Event       `json:"event,omitempty"`
}

// ApiGroup is an alias for Item to support the new naming convention
type ApiGroup = Item

// Request represents an HTTP request
type Request struct {
	Method string      `json:"method"`
	Header []Header    `json:"header,omitempty"`
	Body   *Body       `json:"body,omitempty"`
	URL    interface{} `json:"url"` // Can be string or URL object
	Auth   *Auth       `json:"auth,omitempty"`
	Event  []Event     `json:"event,omitempty"`
}

// URL represents a structured URL
type URL struct {
	Raw      string       `json:"raw"`
	Protocol string       `json:"protocol,omitempty"`
	Host     []string     `json:"host,omitempty"`
	Path     []string     `json:"path,omitempty"`
	Query    []QueryParam `json:"query,omitempty"`
}

// Header represents an HTTP header
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type,omitempty"`
}

// QueryParam represents a URL query parameter
type QueryParam struct {
	Key         string `json:"key"`
	Value       string `json:"value"`
	Description string `json:"description,omitempty"`
	Disabled    bool   `json:"disabled,omitempty"`
}

// Body represents request body
type Body struct {
	Mode    string                 `json:"mode"`
	Raw     string                 `json:"raw,omitempty"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// Response represents a sample response
type Response struct {
	Name            string   `json:"name"`
	OriginalRequest *Request `json:"originalRequest,omitempty"`
	Status          string   `json:"status"`
	Code            int      `json:"code"`
	Header          []Header `json:"header,omitempty"`
	Body            string   `json:"body,omitempty"`
	ResponseTime    int      `json:"responseTime,omitempty"`
}

// RequestItem represents a request with its path context
type RequestItem struct {
	Name       string
	Path       string
	Request    *Request
	Item       *Item // The request item itself
	ParentItem *Item // The parent folder item (for variable resolution)
}
