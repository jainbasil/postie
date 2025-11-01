package scripting

import (
	"sync"
)

// GlobalStore manages global variables that persist across requests
type GlobalStore struct {
	mu        sync.RWMutex
	variables map[string]interface{}
}

// NewGlobalStore creates a new global variable store
func NewGlobalStore() *GlobalStore {
	return &GlobalStore{
		variables: make(map[string]interface{}),
	}
}

// Set sets a global variable
func (g *GlobalStore) Set(name string, value interface{}) {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.variables[name] = value
}

// Get retrieves a global variable
func (g *GlobalStore) Get(name string) (interface{}, bool) {
	g.mu.RLock()
	defer g.mu.RUnlock()
	value, exists := g.variables[name]
	return value, exists
}

// GetString retrieves a global variable as a string
func (g *GlobalStore) GetString(name string) string {
	value, exists := g.Get(name)
	if !exists {
		return ""
	}
	if str, ok := value.(string); ok {
		return str
	}
	return ""
}

// Clear removes a global variable
func (g *GlobalStore) Clear(name string) {
	g.mu.Lock()
	defer g.mu.Unlock()
	delete(g.variables, name)
}

// ClearAll removes all global variables
func (g *GlobalStore) ClearAll() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.variables = make(map[string]interface{})
}

// GetAll returns all global variables
func (g *GlobalStore) GetAll() map[string]interface{} {
	g.mu.RLock()
	defer g.mu.RUnlock()
	
	// Return a copy to avoid race conditions
	copy := make(map[string]interface{}, len(g.variables))
	for k, v := range g.variables {
		copy[k] = v
	}
	return copy
}

// Has checks if a global variable exists
func (g *GlobalStore) Has(name string) bool {
	g.mu.RLock()
	defer g.mu.RUnlock()
	_, exists := g.variables[name]
	return exists
}
