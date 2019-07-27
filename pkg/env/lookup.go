package env

import (
	"fmt"
	"os"
	"sync"
)

// Lookup holds environment variables.
type Lookup struct {
	mu   sync.Mutex
	vars map[string]string
}

// Set environment variable into container.
func (e *Lookup) Set(key, val string) {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.vars[key] = val
}

// Get environment variable from container.
func (e *Lookup) Get(key string) string {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.vars[key]
}

// NewLookup create a new container from input keys.
func NewLookup(keys ...string) (*Lookup, error) {
	e := &Lookup{vars: map[string]string{}}
	for _, key := range keys {
		val, ok := os.LookupEnv(key)
		if !ok {
			return nil, fmt.Errorf("$%s is not present", key)
		}
		e.Set(key, val)
	}
	return e, nil
}
