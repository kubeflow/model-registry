package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"k8s.io/apimachinery/pkg/util/yaml"
)

// LoaderEventHandler is called after each record is successfully processed.
// Use this to trigger side effects like cache invalidation or notifications.
type LoaderEventHandler[E any, A any] func(ctx context.Context, record Record[E, A]) error

// EntitySaver persists an entity and returns the saved entity (with ID populated).
type EntitySaver[E any] func(entity E) (E, error)

// ArtifactSaver persists an artifact associated with an entity.
type ArtifactSaver[A any] func(artifact A, entityID int32) error

// EntityIDGetter extracts the ID from an entity (returns nil if not set).
type EntityIDGetter[E any] func(entity E) *int32

// EntityNameGetter extracts the name from an entity for logging.
type EntityNameGetter[E any] func(entity E) string

// SourceStatusSaver persists the status of a source (available/error/disabled).
type SourceStatusSaver func(sourceID, status, errorMsg string)

// Source status constants
const (
	SourceStatusAvailable = "available"
	SourceStatusError     = "error"
	SourceStatusDisabled  = "disabled"
)

// LoaderConfig configures a Loader instance.
type LoaderConfig[E any, A any] struct {
	// Paths are the config file paths to load sources from.
	Paths []string

	// ProviderRegistry contains registered provider implementations.
	ProviderRegistry *ProviderRegistry[E, A]

	// SaveEntity persists an entity to the database.
	SaveEntity EntitySaver[E]

	// SaveArtifact persists an artifact associated with an entity.
	SaveArtifact ArtifactSaver[A]

	// GetEntityID extracts the ID from an entity.
	GetEntityID EntityIDGetter[E]

	// GetEntityName extracts a name from an entity for logging.
	GetEntityName EntityNameGetter[E]

	// SaveSourceStatus persists source status to the database.
	// Optional - if nil, status is not persisted.
	SaveSourceStatus SourceStatusSaver

	// DeleteArtifactsByEntity removes all artifacts for an entity before re-adding.
	// Optional - if nil, artifacts are not cleaned up before save.
	DeleteArtifactsByEntity func(entityID int32) error

	// DeleteEntitiesBySource removes all entities for a source ID.
	// Called when a source is removed or disabled.
	DeleteEntitiesBySource func(sourceID string) error

	// GetDistinctSourceIDs returns all source IDs that have entities in the DB.
	// Used for cleanup of orphaned sources.
	GetDistinctSourceIDs func() ([]string, error)

	// SetEntitySourceID sets the source ID on an entity.
	SetEntitySourceID func(entity E, sourceID string)

	// IsEntityNil checks if an entity is nil (for batch completion detection).
	IsEntityNil func(entity E) bool

	// Logger for logging messages (optional, defaults to no-op).
	Logger LoaderLogger
}

// LoaderLogger is an interface for logging within the loader.
type LoaderLogger interface {
	Infof(format string, args ...any)
	Errorf(format string, args ...any)
}

// noopLogger is a no-op logger implementation.
type noopLogger struct{}

func (noopLogger) Infof(format string, args ...any)  {}
func (noopLogger) Errorf(format string, args ...any) {}

// SourceConfig is the structure for catalog sources YAML files.
type SourceConfig struct {
	Catalogs []SourceConfigEntry `json:"catalogs" yaml:"catalogs"`
}

// SourceConfigEntry is a single entry in the sources YAML file.
type SourceConfigEntry struct {
	ID             string         `json:"id" yaml:"id"`
	Name           string         `json:"name,omitempty" yaml:"name,omitempty"`
	Type           string         `json:"type" yaml:"type"`
	Enabled        *bool          `json:"enabled,omitempty" yaml:"enabled,omitempty"`
	Labels         []string       `json:"labels,omitempty" yaml:"labels,omitempty"`
	Properties     map[string]any `json:"properties,omitempty" yaml:"properties,omitempty"`
	IncludedModels []string       `json:"includedModels,omitempty" yaml:"includedModels,omitempty"`
	ExcludedModels []string       `json:"excludedModels,omitempty" yaml:"excludedModels,omitempty"`
}

// ToSource converts a config entry to a Source.
func (e SourceConfigEntry) ToSource(origin string) Source {
	return Source{
		ID:            e.ID,
		Name:          e.Name,
		Type:          e.Type,
		Enabled:       e.Enabled,
		Labels:        e.Labels,
		Properties:    e.Properties,
		IncludedItems: e.IncludedModels,
		ExcludedItems: e.ExcludedModels,
		Origin:        origin,
	}
}

// Loader manages loading data from sources into repositories.
type Loader[E any, A any] struct {
	config LoaderConfig[E, A]

	// Sources contains current source information loaded from the configuration files.
	Sources *SourceCollection

	closersMu     sync.Mutex
	closer        func() // cancels the current loading goroutines
	handlers      []LoaderEventHandler[E, A]
	loadedSources map[string]bool // tracks which source IDs have been loaded
	logger        LoaderLogger
}

// NewLoader creates a new Loader with the given configuration.
func NewLoader[E any, A any](config LoaderConfig[E, A]) *Loader[E, A] {
	// Convert paths to absolute for consistent origin ordering.
	absPaths := make([]string, 0, len(config.Paths))
	for _, p := range config.Paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			absPath = p
		}
		absPaths = append(absPaths, absPath)
	}

	logger := config.Logger
	if logger == nil {
		logger = noopLogger{}
	}

	return &Loader[E, A]{
		config:        config,
		Sources:       NewSourceCollection(absPaths...),
		loadedSources: map[string]bool{},
		logger:        logger,
	}
}

// RegisterEventHandler adds a function that will be called for every
// successfully processed record. This should be called before Start.
func (l *Loader[E, A]) RegisterEventHandler(fn LoaderEventHandler[E, A]) {
	l.handlers = append(l.handlers, fn)
}

// Start processes the sources YAML files and loads data.
// Background goroutines will be stopped when the context is canceled.
func (l *Loader[E, A]) Start(ctx context.Context) error {
	// Phase 1: Parse all config files and merge sources
	for _, path := range l.config.Paths {
		if err := l.parseAndMerge(path); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}

	// Delete entities from unknown or disabled sources
	if err := l.removeEntitiesFromMissingSources(); err != nil {
		return fmt.Errorf("failed to remove entities from missing sources: %w", err)
	}

	// Phase 2: Load entities from merged sources
	if err := l.loadAllEntities(ctx); err != nil {
		return err
	}

	// Phase 3: Watch config files for hot-reload
	for _, path := range l.config.Paths {
		watcher := NewFileWatcher(path, 5*time.Second)
		changes := watcher.Watch(ctx)
		go func(p string) {
			for range changes {
				l.logger.Infof("Config file changed: %s, reloading...", p)
				if err := l.Reload(ctx); err != nil {
					l.logger.Errorf("Failed to reload after config change %s: %v", p, err)
				}
			}
		}(path)
	}

	return nil
}

// Reload re-parses config files, cleans up missing sources, and reloads all entities.
func (l *Loader[E, A]) Reload(ctx context.Context) error {
	for _, path := range l.config.Paths {
		if err := l.parseAndMerge(path); err != nil {
			l.logger.Errorf("failed to reload config %s: %v", path, err)
		}
	}
	_ = l.removeEntitiesFromMissingSources()
	return l.loadAllEntities(ctx)
}

// parseAndMerge parses a config file and merges its sources into the collection.
func (l *Loader[E, A]) parseAndMerge(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", path, err)
	}

	config, err := l.readConfig(path)
	if err != nil {
		return err
	}

	sources := make(map[string]Source, len(config.Catalogs))
	for _, entry := range config.Catalogs {
		l.logger.Infof("reading config type %s...", entry.Type)
		if entry.ID == "" {
			return fmt.Errorf("invalid source: missing id")
		}
		if _, exists := sources[entry.ID]; exists {
			return fmt.Errorf("invalid source: duplicate id %s", entry.ID)
		}

		// Validate include/exclude patterns early
		if err := ValidatePatterns(entry.IncludedModels, entry.ExcludedModels); err != nil {
			return fmt.Errorf("invalid source %s: %w", entry.ID, err)
		}

		sources[entry.ID] = entry.ToSource(path)
		l.logger.Infof("loaded source %s of type %s", entry.ID, entry.Type)
	}

	return l.Sources.Merge(path, sources)
}

func (l *Loader[E, A]) readConfig(path string) (*SourceConfig, error) {
	config := &SourceConfig{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(bytes, config); err != nil {
		return nil, err
	}

	return config, nil
}

// loadAllEntities loads entities from all merged sources.
func (l *Loader[E, A]) loadAllEntities(ctx context.Context) error {
	l.loadedSources = map[string]bool{}
	return l.updateDatabase(ctx)
}

func (l *Loader[E, A]) updateDatabase(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)

	l.closersMu.Lock()
	if l.closer != nil {
		l.closer()
	}
	l.closer = cancel
	l.closersMu.Unlock()

	records := l.readProviderRecords(ctx)

	go func() {
		for record := range records {
			if l.config.IsEntityNil != nil && l.config.IsEntityNil(record.Entity) {
				continue
			}

			name := ""
			if l.config.GetEntityName != nil {
				name = l.config.GetEntityName(record.Entity)
			}

			l.logger.Infof("Loading entity %s with %d artifact(s)", name, len(record.Artifacts))

			entity, err := l.config.SaveEntity(record.Entity)
			if err != nil {
				l.logger.Errorf("%s: unable to save: %v", name, err)
				continue
			}

			entityID := l.config.GetEntityID(entity)
			if entityID == nil {
				l.logger.Errorf("%s: entity has no ID after save", name)
				continue
			}

			// Remove old artifacts before adding new ones
			if l.config.DeleteArtifactsByEntity != nil {
				if err := l.config.DeleteArtifactsByEntity(*entityID); err != nil {
					l.logger.Errorf("%s: unable to remove old artifacts: %v", name, err)
				}
			}

			// Save new artifacts
			for i, artifact := range record.Artifacts {
				if err := l.config.SaveArtifact(artifact, *entityID); err != nil {
					l.logger.Errorf("%s, artifact %d: %v", name, i, err)
				}
			}

			// Call event handlers
			for _, handler := range l.handlers {
				if err := handler(ctx, record); err != nil {
					l.logger.Errorf("%s: event handler error: %v", name, err)
				}
			}
		}
	}()

	return nil
}

// readProviderRecords calls the provider for every merged source that hasn't
// been loaded yet, and merges the returned channels together.
func (l *Loader[E, A]) readProviderRecords(ctx context.Context) <-chan Record[E, A] {
	ch := make(chan Record[E, A])
	var wg sync.WaitGroup

	mergedSources := l.Sources.AllSources()

	for _, source := range mergedSources {
		// Skip disabled sources
		if !source.IsEnabled() {
			if l.config.SaveSourceStatus != nil {
				l.config.SaveSourceStatus(source.ID, SourceStatusDisabled, "")
			}
			continue
		}

		// Skip already loaded sources
		if l.loadedSources[source.ID] {
			continue
		}

		if source.Type == "" {
			l.logger.Errorf("source %s has no type defined, skipping", source.ID)
			if l.config.SaveSourceStatus != nil {
				l.config.SaveSourceStatus(source.ID, SourceStatusError, "source has no type defined")
			}
			continue
		}

		l.loadedSources[source.ID] = true

		l.logger.Infof("Reading entities from %s source %s", source.Type, source.ID)

		providerFunc, ok := l.config.ProviderRegistry.Get(source.Type)
		if !ok {
			l.logger.Errorf("provider type %s not registered", source.Type)
			if l.config.SaveSourceStatus != nil {
				l.config.SaveSourceStatus(source.ID, SourceStatusError, fmt.Sprintf("provider type %q not registered", source.Type))
			}
			continue
		}

		sourceDir := filepath.Dir(source.Origin)
		sourceCopy := source // capture for goroutine

		records, err := providerFunc(ctx, &sourceCopy, sourceDir)
		if err != nil {
			l.logger.Errorf("error reading provider type %s with id %s: %v", source.Type, source.ID, err)
			if l.config.SaveSourceStatus != nil {
				l.config.SaveSourceStatus(source.ID, SourceStatusError, err.Error())
			}
			continue
		}

		wg.Add(1)
		go func(sourceID string) {
			defer wg.Done()

			for record := range records {
				// Set source ID on entity
				if l.config.SetEntitySourceID != nil && !l.config.IsEntityNil(record.Entity) {
					l.config.SetEntitySourceID(record.Entity, sourceID)
				}
				ch <- record
			}

			// Mark source as available
			if l.config.SaveSourceStatus != nil && ctx.Err() == nil {
				l.config.SaveSourceStatus(sourceID, SourceStatusAvailable, "")
			}
		}(source.ID)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

func (l *Loader[E, A]) removeEntitiesFromMissingSources() error {
	if l.config.DeleteEntitiesBySource == nil || l.config.GetDistinctSourceIDs == nil {
		return nil
	}

	enabledSourceIDs := make(map[string]bool)
	for id, source := range l.Sources.AllSources() {
		if source.IsEnabled() {
			enabledSourceIDs[id] = true
		}
	}

	existingSourceIDs, err := l.config.GetDistinctSourceIDs()
	if err != nil {
		return fmt.Errorf("unable to retrieve existing source IDs: %w", err)
	}

	for _, oldSource := range existingSourceIDs {
		if !enabledSourceIDs[oldSource] {
			l.logger.Infof("Removing entities from source %s", oldSource)
			if err := l.config.DeleteEntitiesBySource(oldSource); err != nil {
				return fmt.Errorf("unable to remove entities from source %q: %w", oldSource, err)
			}
		}
	}

	return nil
}
