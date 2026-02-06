// Package yaml provides a base YAML file provider for catalog data.
// It handles file reading, parsing, and hot-reloading, while delegating
// entity-specific conversion to user-provided functions.
package yaml

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/kubeflow/model-registry/pkg/catalog"
	k8syaml "k8s.io/apimachinery/pkg/util/yaml"
)

// Config configures a YAML provider.
type Config[E any, A any] struct {
	// PathKey is the property key in source.Properties that contains the YAML file path.
	// Defaults to "yamlCatalogPath" if empty.
	PathKey string

	// Parse parses raw YAML bytes into a slice of entity records.
	Parse func(data []byte) ([]catalog.Record[E, A], error)

	// Filter optionally filters records before emitting them.
	// Return true to include the record, false to exclude it.
	// If nil, all records are included.
	Filter func(record catalog.Record[E, A]) bool

	// Logger for logging messages (optional).
	Logger Logger

	// WatchInterval is the interval for checking file changes.
	// Defaults to 5 seconds if zero.
	WatchInterval time.Duration
}

// Logger is an interface for logging.
type Logger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

type noopLogger struct{}

func (noopLogger) Infof(format string, args ...any)  {}
func (noopLogger) Errorf(format string, args ...any) {}

// Provider is a YAML file-based data provider.
type Provider[E any, A any] struct {
	config Config[E, A]
	path   string
	filter *catalog.ItemFilter
	logger Logger
}

// NewProvider creates a new YAML provider with the given configuration.
// It reads the file path from source.Properties using the configured PathKey.
func NewProvider[E any, A any](config Config[E, A], source *catalog.Source, reldir string) (*Provider[E, A], error) {
	pathKey := config.PathKey
	if pathKey == "" {
		pathKey = "yamlCatalogPath"
	}

	path, ok := source.Properties[pathKey].(string)
	if !ok || path == "" {
		return nil, fmt.Errorf("missing %s string property", pathKey)
	}

	// Resolve relative paths
	if !filepath.IsAbs(path) {
		path = filepath.Join(reldir, path)
	}

	// Build filter from source configuration
	filter, err := catalog.NewItemFilterFromSource(source, nil, nil)
	if err != nil {
		return nil, err
	}

	logger := config.Logger
	if logger == nil {
		logger = noopLogger{}
	}

	return &Provider[E, A]{
		config: config,
		path:   path,
		filter: filter,
		logger: logger,
	}, nil
}

// Records starts reading the YAML file and returns a channel of records.
// The channel is closed when the context is canceled.
// The provider watches for file changes and re-emits records when the file changes.
func (p *Provider[E, A]) Records(ctx context.Context) (<-chan catalog.Record[E, A], error) {
	// Read initial data to catch errors early
	records, err := p.read()
	if err != nil {
		return nil, err
	}

	ch := make(chan catalog.Record[E, A])
	go func() {
		defer close(ch)

		// Send initial records
		p.emit(ctx, records, ch)

		// Watch for changes
		p.watchAndReload(ctx, ch)
	}()

	return ch, nil
}

func (p *Provider[E, A]) read() ([]catalog.Record[E, A], error) {
	data, err := os.ReadFile(p.path)
	if err != nil {
		return nil, fmt.Errorf("failed to read YAML file %s: %w", p.path, err)
	}

	records, err := p.config.Parse(data)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML file %s: %w", p.path, err)
	}

	return records, nil
}

func (p *Provider[E, A]) emit(ctx context.Context, records []catalog.Record[E, A], out chan<- catalog.Record[E, A]) {
	done := ctx.Done()
	for _, record := range records {
		// Apply custom filter if provided
		if p.config.Filter != nil && !p.config.Filter(record) {
			continue
		}

		select {
		case out <- record:
		case <-done:
			return
		}
	}

	// Send an empty record to indicate batch completion
	var zero catalog.Record[E, A]
	select {
	case out <- zero:
	case <-done:
	}
}

func (p *Provider[E, A]) watchAndReload(ctx context.Context, ch chan<- catalog.Record[E, A]) {
	interval := p.config.WatchInterval
	if interval == 0 {
		interval = 5 * time.Second
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	var lastModTime time.Time
	if info, err := os.Stat(p.path); err == nil {
		lastModTime = info.ModTime()
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			info, err := os.Stat(p.path)
			if err != nil {
				continue
			}

			if info.ModTime().After(lastModTime) {
				lastModTime = info.ModTime()
				p.logger.Infof("Reloading YAML file %s", p.path)

				records, err := p.read()
				if err != nil {
					p.logger.Errorf("Failed to reload YAML file: %v", err)
					continue
				}

				p.emit(ctx, records, ch)
			}
		}
	}
}

// NewProviderFunc creates a ProviderFunc that can be registered with a ProviderRegistry.
// This is a convenience function for creating providers using the standard pattern.
func NewProviderFunc[E any, A any](config Config[E, A]) catalog.ProviderFunc[E, A] {
	return func(ctx context.Context, source *catalog.Source, reldir string) (<-chan catalog.Record[E, A], error) {
		provider, err := NewProvider(config, source, reldir)
		if err != nil {
			return nil, err
		}
		return provider.Records(ctx)
	}
}

// SimpleCatalog is a generic structure for simple YAML catalogs.
type SimpleCatalog[E any] struct {
	Source   string `json:"source" yaml:"source"`
	Entities []E    `json:"entities" yaml:"entities"`
}

// ParseSimpleCatalog creates a parser for SimpleCatalog format.
// The toRecord function converts each entity to a Record.
func ParseSimpleCatalog[E any, A any](toRecord func(E) catalog.Record[E, A]) func([]byte) ([]catalog.Record[E, A], error) {
	return func(data []byte) ([]catalog.Record[E, A], error) {
		var cat SimpleCatalog[E]
		if err := k8syaml.UnmarshalStrict(data, &cat); err != nil {
			return nil, err
		}

		records := make([]catalog.Record[E, A], 0, len(cat.Entities))
		for _, entity := range cat.Entities {
			records = append(records, toRecord(entity))
		}

		return records, nil
	}
}

// FileWatcher watches a file for changes using polling.
// This is a simple implementation that can be used if more sophisticated
// file watching (like fsnotify) is not available.
type FileWatcher struct {
	mu          sync.Mutex
	path        string
	lastModTime time.Time
	interval    time.Duration
}

// NewFileWatcher creates a new file watcher.
func NewFileWatcher(path string, interval time.Duration) *FileWatcher {
	if interval == 0 {
		interval = 5 * time.Second
	}

	w := &FileWatcher{
		path:     path,
		interval: interval,
	}

	if info, err := os.Stat(path); err == nil {
		w.lastModTime = info.ModTime()
	}

	return w
}

// Watch returns a channel that receives a value whenever the file changes.
// The channel is closed when the context is canceled.
func (w *FileWatcher) Watch(ctx context.Context) <-chan struct{} {
	ch := make(chan struct{})

	go func() {
		defer close(ch)

		ticker := time.NewTicker(w.interval)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				if w.hasChanged() {
					select {
					case ch <- struct{}{}:
					case <-ctx.Done():
						return
					}
				}
			}
		}
	}()

	return ch
}

func (w *FileWatcher) hasChanged() bool {
	w.mu.Lock()
	defer w.mu.Unlock()

	info, err := os.Stat(w.path)
	if err != nil {
		return false
	}

	if info.ModTime().After(w.lastModTime) {
		w.lastModTime = info.ModTime()
		return true
	}

	return false
}
