// Package catalog provides reusable abstractions for building catalog-style
// read-only aggregation services. It extracts common patterns from the Model Catalog
// to enable rapid creation of new catalog components (e.g., MCP Catalog).
package catalog

import (
	"context"
	"fmt"
	"sync"
)

// Record represents a single entity record with its associated artifacts.
// This is generic to support different entity types.
type Record[E any, A any] struct {
	// Entity is the main entity being loaded (e.g., a catalog model).
	// A nil Entity signals that a batch of records has been fully sent.
	Entity E

	// Artifacts are the associated artifacts for this entity.
	Artifacts []A
}

// ProviderFunc emits records in a channel and is expected to spawn a goroutine
// and return immediately. The returned channel must close when the goroutine ends.
// The goroutine should end when the context is canceled, but may end sooner.
//
// The function may emit a record with a nil Entity to indicate that the
// complete set of entities has been sent (batch completion marker).
//
// Parameters:
//   - ctx: Context for cancellation
//   - source: The source configuration for this provider
//   - reldir: The directory to resolve relative paths from (typically the config file's directory)
type ProviderFunc[E any, A any] func(ctx context.Context, source *Source, reldir string) (<-chan Record[E, A], error)

// ProviderRegistry manages provider type registrations.
// It allows registering provider functions by name (e.g., "yaml", "http")
// and retrieving them for use when loading from sources.
type ProviderRegistry[E any, A any] struct {
	mu        sync.RWMutex
	providers map[string]ProviderFunc[E, A]
}

// NewProviderRegistry creates a new empty provider registry.
func NewProviderRegistry[E any, A any]() *ProviderRegistry[E, A] {
	return &ProviderRegistry[E, A]{
		providers: make(map[string]ProviderFunc[E, A]),
	}
}

// Register adds a provider function with the given name.
// Returns an error if a provider with that name already exists.
func (r *ProviderRegistry[E, A]) Register(name string, fn ProviderFunc[E, A]) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.providers[name]; exists {
		return fmt.Errorf("provider type %q already exists", name)
	}
	r.providers[name] = fn
	return nil
}

// MustRegister is like Register but panics on error.
// This is useful for init() functions.
func (r *ProviderRegistry[E, A]) MustRegister(name string, fn ProviderFunc[E, A]) {
	if err := r.Register(name, fn); err != nil {
		panic(err)
	}
}

// Get retrieves a provider function by name.
// Returns the function and true if found, or nil and false if not.
func (r *ProviderRegistry[E, A]) Get(name string) (ProviderFunc[E, A], bool) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	fn, ok := r.providers[name]
	return fn, ok
}

// Names returns a list of all registered provider names.
func (r *ProviderRegistry[E, A]) Names() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()

	names := make([]string, 0, len(r.providers))
	for name := range r.providers {
		names = append(names, name)
	}
	return names
}

// Has returns true if a provider with the given name is registered.
func (r *ProviderRegistry[E, A]) Has(name string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()

	_, ok := r.providers[name]
	return ok
}
