package catalog

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// Source status constants matching the OpenAPI enum values
const (
	SourceStatusAvailable          = "available"
	SourceStatusPartiallyAvailable = "partially-available"
	SourceStatusError              = "error"
	SourceStatusDisabled           = "disabled"
)

// PartiallyAvailableError indicates that a source loaded some models successfully
// but encountered errors with others.
type PartiallyAvailableError struct {
	FailedModels []string
}

func (e *PartiallyAvailableError) Error() string {
	return fmt.Sprintf("Failed to fetch some models, ensure models exist and are accessible with given credentials. Failed models: %v", e.FailedModels)
}

func (e *PartiallyAvailableError) Is(target error) bool {
	_, ok := target.(*PartiallyAvailableError)
	return ok
}

// ErrPartiallyAvailable is used with errors.Is() to check for this error type.
var ErrPartiallyAvailable error = &PartiallyAvailableError{}

// ModelProviderRecord contains one model and its associated artifacts.
type ModelProviderRecord struct {
	Model     dbmodels.CatalogModel
	Artifacts []dbmodels.CatalogArtifact
	// Error can be set here to emit successfully loaded models before updating source status err.
	Error error
}

// ModelProviderFunc emits models and related data in the channel it returns. It is
// expected to spawn a goroutine and return immediately. The returned channel must
// close when the goroutine ends. The goroutine should end when the context is
// canceled, but may end sooner.
//
// The function may emit a record with a nil Model to indicate that the
// complete set of models has been sent.
type ModelProviderFunc func(ctx context.Context, source *Source, reldir string) (<-chan ModelProviderRecord, error)

var registeredModelProviders = map[string]ModelProviderFunc{}

func RegisterModelProvider(name string, callback ModelProviderFunc) error {
	if _, exists := registeredModelProviders[name]; exists {
		return fmt.Errorf("provider type %s already exists", name)
	}
	registeredModelProviders[name] = callback
	return nil
}

// LoaderEventHandler is the definition of a function called after a model is loaded.
type LoaderEventHandler func(ctx context.Context, record ModelProviderRecord) error

// FieldFilter represents a single field filter within a named query
type FieldFilter struct {
	Operator string `json:"operator" yaml:"operator"`
	Value    any    `json:"value" yaml:"value"`
}

// sourceConfig is the structure for the catalog sources YAML file.
type sourceConfig struct {
	Catalogs     []Source                          `json:"catalogs"`
	Labels       []map[string]any                  `json:"labels,omitempty"`
	NamedQueries map[string]map[string]FieldFilter `json:"namedQueries,omitempty" yaml:"namedQueries,omitempty"`
}

// Source is a single entry from the catalog sources YAML file.
type Source struct {
	apimodels.CatalogSource `json:",inline"`

	// Catalog type to use, must match one of the registered types
	Type string `json:"type"`

	// Properties used for configuring the catalog connection based on catalog implementation
	Properties map[string]any `json:"properties,omitempty"`

	// Origin is the absolute path of the config file this source was loaded from.
	// This is set automatically during loading and used for resolving relative paths.
	// It is not read from YAML; it's set programmatically.
	Origin string `json:"-" yaml:"-"`
}

type Loader struct {
	// Sources contains current source information loaded from the configuration files.
	Sources *SourceCollection

	// Labels contains current labels loaded from the configuration files.
	Labels *LabelCollection

	paths         []string
	services      service.Services
	closersMu     sync.Mutex
	closer        func() // cancels the current model loading goroutines
	handlers      []LoaderEventHandler
	loadedSources map[string]bool // tracks which source IDs have been loaded

	// Leader state management
	leaderMu sync.RWMutex
	isLeader bool           // true when in leader mode
	writesWG sync.WaitGroup // tracks number of database write operations in progress

	// File watcher state
	watchersMu      sync.Mutex
	watchersStarted bool
	watchersCancel  context.CancelFunc // cancels file watchers on shutdown
}

func NewLoader(services service.Services, paths []string) *Loader {
	// Convert paths to absolute for consistent origin ordering.
	// This matches how loadOne converts paths before calling Merge.
	absPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			// Fall back to original path if conversion fails
			absPath = p
		}
		absPaths = append(absPaths, absPath)
	}

	return &Loader{
		Sources:       NewSourceCollection(absPaths...),
		Labels:        NewLabelCollection(),
		paths:         paths,
		services:      services,
		loadedSources: map[string]bool{},
	}
}

// RegisterEventHandler adds a function that will be called for every
// successfully processed record. This should be called before StartReadOnly.
//
// Handlers are called in the order they are registered.
func (l *Loader) RegisterEventHandler(fn LoaderEventHandler) {
	l.handlers = append(l.handlers, fn)
}

// StartReadOnly initializes the loader in read-only mode (standby pod).
// Call once at pod startup.
//
// Behavior:
// - Parses all config files into in-memory collections (Sources, Labels)
// - Sets up file watchers that reload config when files change
// - Skips database writes (leader-only operation)
// - File watchers run until ctx cancels (pod shutdown)
//
// Returns immediately after setup; non-blocking.
func (l *Loader) StartReadOnly(ctx context.Context) error {
	l.watchersMu.Lock()
	if l.watchersStarted {
		l.watchersMu.Unlock()
		return fmt.Errorf("StartReadOnly already called")
	}
	l.watchersStarted = true

	watcherCtx, cancel := context.WithCancel(ctx)
	l.watchersCancel = cancel
	l.watchersMu.Unlock()

	glog.Info("Starting loader in read-only mode (standby)")

	for _, path := range l.paths {
		if err := l.parseAndMerge(path); err != nil {
			return fmt.Errorf("%s: %w", path, err)
		}
	}

	for _, path := range l.paths {
		go func(path string) {
			changes, err := getMonitor().Path(watcherCtx, path)
			if err != nil {
				glog.Errorf("unable to watch sources file (%s): %v", path, err)
				return
			}

			for range changes {
				glog.Infof("Config file changed, reloading: %s", path)
				l.reloadConfig(watcherCtx)
			}
		}(path)
	}

	glog.Info("Read-only mode initialized successfully")
	return nil
}

// StartLeader transitions the loader to leader mode and blocks until the
// context is cancelled (leadership lost or pod shutdown).
//
// Performs database writes (fetches models, writes to database, cleans up orphans).
// Can be called multiple times as leadership changes.
//
// Pass the leadership context (canceled when leadership is lost).
func (l *Loader) StartLeader(ctx context.Context) error {
	l.leaderMu.Lock()
	if l.isLeader {
		l.leaderMu.Unlock()
		return fmt.Errorf("already in leader mode")
	}
	l.isLeader = true
	l.leaderMu.Unlock()

	glog.Info("Transitioning to leader mode (read-write)")

	if err := l.performLeaderOperations(ctx); err != nil {
		l.leaderMu.Lock()
		l.isLeader = false
		l.leaderMu.Unlock()
		return fmt.Errorf("failed to perform leader operations: %w", err)
	}

	glog.Info("Leader mode active")

	// Just wait for context cancellation
	<-ctx.Done()
	glog.Info("Leadership context cancelled, cleaning up...")

	l.waitForInflightWrites(5 * time.Second)

	l.leaderMu.Lock()
	l.isLeader = false
	l.leaderMu.Unlock()

	glog.Info("Leader mode stopped")
	return ctx.Err()
}

// Shutdown gracefully shuts down the loader, stopping file watchers and
// waiting for any inflight operations.
func (l *Loader) Shutdown() error {
	glog.Info("Shutting down loader...")

	// Note: Leader operations stop when leaderCtx is cancelled by the elector.
	// No need to explicitly stop here.

	// Stop file watchers
	l.watchersMu.Lock()
	if l.watchersCancel != nil {
		l.watchersCancel()
	}
	l.watchersMu.Unlock()

	// Wait for inflight writes
	l.waitForInflightWrites(10 * time.Second)

	glog.Info("Loader shutdown complete")
	return nil
}

// performLeaderWrites executes database write operations: removing orphaned
// models and loading all models from sources.
func (l *Loader) performLeaderWrites(ctx context.Context) error {
	if err := l.removeModelsFromMissingSources(); err != nil {
		return fmt.Errorf("failed to remove models from missing sources: %w", err)
	}
	return l.loadAllModels(ctx)
}

// performLeaderOperations executes the database operations that only the leader performs.
func (l *Loader) performLeaderOperations(ctx context.Context) error {
	return l.performLeaderWrites(ctx)
}

// shouldWriteDatabase returns true when in leader mode.
func (l *Loader) shouldWriteDatabase() bool {
	l.leaderMu.RLock()
	defer l.leaderMu.RUnlock()
	return l.isLeader
}

// reloadConfig handles config file changes.
// Re-parses all config files and reloads models when leader.
func (l *Loader) reloadConfig(ctx context.Context) {
	// Re-parse all config files (both leader and standby)
	for _, path := range l.paths {
		if err := l.parseAndMerge(path); err != nil {
			glog.Errorf("unable to reload sources from %s: %v", path, err)
		}
	}

	// Database writes (leader only)
	if l.shouldWriteDatabase() {
		if err := l.performLeaderWrites(ctx); err != nil {
			glog.Errorf("unable to perform leader writes: %v", err)
		}
	}
}

// waitForInflightWrites waits for all inflight database writes to complete
// with a timeout.
func (l *Loader) waitForInflightWrites(timeout time.Duration) {
	glog.Info("Waiting for inflight writes to complete...")

	done := make(chan struct{})
	go func() {
		l.writesWG.Wait()
		close(done)
	}()

	select {
	case <-done:
		glog.Info("All inflight writes completed")
	case <-time.After(timeout):
		glog.Warningf("Timeout waiting for inflight writes to complete")
	}
}

// parseAndMerge parses a config file and merges its sources/labels into the collections.
func (l *Loader) parseAndMerge(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("failed to get absolute path for %s: %v", path, err)
	}

	config, err := l.read(path)
	if err != nil {
		return err
	}

	if err = l.updateSources(path, config); err != nil {
		return err
	}

	return l.updateLabels(path, config)
}

func (l *Loader) loadAllModels(ctx context.Context) error {
	l.loadedSources = map[string]bool{}

	return l.updateDatabase(ctx)
}

func (l *Loader) read(path string) (*sourceConfig, error) {
	config := &sourceConfig{}
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	// Validate named queries if present
	if config.NamedQueries != nil {
		if err := ValidateNamedQueries(config.NamedQueries); err != nil {
			return nil, fmt.Errorf("invalid named queries in %s: %w", path, err)
		}
	}

	// Note: We intentionally do NOT filter disabled sources or apply defaults here.
	// This allows field-level merging in SourceCollection to work correctly:
	// - A base source with enabled=false can be enabled by a user override with just id + enabled=true
	// - Defaults are applied after merging in SourceCollection.merged()

	return config, nil
}

func (l *Loader) updateSources(path string, config *sourceConfig) error {
	sources := make(map[string]Source, len(config.Catalogs))

	for _, source := range config.Catalogs {
		glog.Infof("reading config type %s...", source.Type)
		id := source.GetId()
		if len(id) == 0 {
			return fmt.Errorf("invalid source: missing id")
		}
		if _, exists := sources[id]; exists {
			return fmt.Errorf("invalid source: duplicate id %s", id)
		}

		// Validate includedModels/excludedModels patterns early (only if set)
		if err := ValidateSourceFilters(source.IncludedModels, source.ExcludedModels); err != nil {
			return fmt.Errorf("invalid source %s: %w", id, err)
		}

		// Set the origin path so relative paths in properties can be resolved
		// relative to this config file's directory
		source.Origin = path
		sources[id] = source
		glog.Infof("loaded source %s of type %s", id, source.Type)
	}

	// Use MergeWithNamedQueries if named queries exist, otherwise use regular Merge
	if config.NamedQueries != nil {
		return l.Sources.MergeWithNamedQueries(path, sources, config.NamedQueries)
	}
	return l.Sources.Merge(path, sources)
}

func (l *Loader) updateLabels(path string, config *sourceConfig) error {
	// Merge labels from config into the label collection
	if config.Labels == nil {
		// No labels in config, but we still need to clear any previous labels from this origin
		return l.Labels.Merge(path, []map[string]any{})
	}

	// Validate that each label has a required "name" field
	for i, label := range config.Labels {
		if name, ok := label["name"]; !ok || name == "" {
			return fmt.Errorf("invalid label at index %d: missing required 'name' field", i)
		}
	}

	return l.Labels.Merge(path, config.Labels)
}

func (l *Loader) updateDatabase(ctx context.Context) error {
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
			// Check if we're still the leader before each write
			if !l.shouldWriteDatabase() {
				glog.Info("No longer leader, stopping database writes")
				return
			}

			// Check context cancellation
			if ctx.Err() != nil {
				glog.Info("Context cancelled, stopping database writes")
				return
			}

			if record.Model == nil {
				continue
			}
			attr := record.Model.GetAttributes()
			if attr == nil || attr.Name == nil {
				continue
			}

			// Track this write operation
			l.writesWG.Add(1)
			defer l.writesWG.Done()

			glog.Infof("Loading model %s with %d artifact(s)", *attr.Name, len(record.Artifacts))

			model, err := l.services.CatalogModelRepository.Save(record.Model)
			if err != nil {
				glog.Errorf("%s: unable to save: %v", *attr.Name, err)
				continue
			}

			modelID := model.GetID()
			if modelID == nil {
				glog.Errorf("%s: model has no ID after save", *attr.Name)
				continue
			}

			// Remove artifacts that existed before.
			err = l.services.CatalogArtifactRepository.DeleteByParentID(service.CatalogModelArtifactTypeName, *modelID)
			if err != nil {
				glog.Errorf("%s: unable to remove old catalog model artifacts: %v", *attr.Name, err)
			}
			err = l.services.CatalogArtifactRepository.DeleteByParentID(service.CatalogMetricsArtifactTypeName, *modelID)
			if err != nil {
				glog.Errorf("%s: unable to remove old catalog metrics artifacts: %v", *attr.Name, err)
			}

			for i, artifact := range record.Artifacts {
				switch {
				case artifact.CatalogModelArtifact != nil:
					_, err = l.services.CatalogModelArtifactRepository.Save(artifact.CatalogModelArtifact, modelID)
				case artifact.CatalogMetricsArtifact != nil:
					_, err = l.services.CatalogMetricsArtifactRepository.Save(artifact.CatalogMetricsArtifact, modelID)
				default:
					err = errors.New("unknown artifact type")
				}

				if err != nil {
					glog.Errorf("%s, artifact %d: %v", *attr.Name, i, err)
				}
			}

			for _, handler := range l.handlers {
				handler(ctx, record)
			}
		}
	}()

	return nil
}

// readProviderRecords calls the provider for every merged source that hasn't
// been loaded yet, and merges the returned channels together. The returned
// channel is closed when the last provider channel is closed.
func (l *Loader) readProviderRecords(ctx context.Context) <-chan ModelProviderRecord {
	ch := make(chan ModelProviderRecord)
	var wg sync.WaitGroup

	// Get all sources from the merged collection.
	// This allows sparse overrides to work: a user can enable a disabled source
	// with just "id" and "enabled: true", inheriting Type and Properties from the base.
	mergedSources := l.Sources.AllSources()

	for _, source := range mergedSources {
		// Skip disabled sources - only load catalog data from enabled sources
		// Per OpenAPI spec, enabled defaults to true, so nil is treated as enabled
		if source.Enabled != nil && !*source.Enabled {
			// Persist disabled status
			l.saveSourceStatus(source.Id, SourceStatusDisabled, "")
			continue
		}

		// Skip sources that have already been loaded
		if l.loadedSources[source.Id] {
			continue
		}

		if source.Type == "" {
			glog.Errorf("source %s has no type defined, skipping", source.Id)
			l.saveSourceStatus(source.Id, SourceStatusError, "source has no type defined")
			continue
		}

		// Mark this source as loaded
		l.loadedSources[source.Id] = true

		glog.Infof("Reading models from %s source %s", source.Type, source.Id)

		registerFunc, ok := registeredModelProviders[source.Type]
		if !ok {
			glog.Errorf("catalog type %s not registered", source.Type)
			l.saveSourceStatus(source.Id, SourceStatusError, fmt.Sprintf("catalog type %q not registered", source.Type))
			continue
		}

		// Use the source's origin directory for resolving relative paths.
		// This allows sources from different config files (e.g., mounted from
		// different configmaps) to use relative paths correctly.
		sourceDir := filepath.Dir(source.Origin)

		records, err := registerFunc(ctx, &source, sourceDir)
		if err != nil {
			glog.Errorf("error reading catalog type %s with id %s: %v", source.Type, source.Id, err)
			l.saveSourceStatus(source.Id, SourceStatusError, err.Error())
			continue
		}

		wg.Add(1)
		go func(ctx context.Context, sourceID string) {
			defer wg.Done()

			modelNames := []string{}
			statusSaved := false

			for r := range records {
				if r.Model == nil {
					glog.Infof("%s: loaded %d models", sourceID, len(modelNames))

					// Copy the list of model names, then clear it.
					modelNameSet := mapset.NewSet(modelNames...)
					modelNames = modelNames[:0]

					go func() {
						count, err := l.removeOrphanedModelsFromSource(sourceID, modelNameSet)
						if err != nil {
							glog.Errorf("error removing orphaned models: %v", err)
						}
						glog.Infof("%s: cleaned up %d models", sourceID, count)
					}()

					// Only save status if context is still valid (no reload in progress)
					if ctx.Err() == nil {
						// Check if there was a partial error (some models failed to load)
						if errors.Is(r.Error, ErrPartiallyAvailable) {
							glog.Warningf("%s: partial error after loading models: %v", sourceID, r.Error)
							l.saveSourceStatus(sourceID, SourceStatusPartiallyAvailable, r.Error.Error())
						} else {
							l.saveSourceStatus(sourceID, SourceStatusAvailable, "")
						}
						statusSaved = true
					}
					continue
				}

				if attr := r.Model.GetAttributes(); attr != nil && attr.Name != nil {
					modelNames = append(modelNames, *attr.Name)
				}

				// Set source_id on every returned model.
				l.setModelSourceID(r.Model, sourceID)

				ch <- r
			}

			// If the channel closed without a nil Model marker and status wasn't already saved,
			// save available status if context is still valid and we processed some models
			if !statusSaved && ctx.Err() == nil && len(modelNames) > 0 {
				l.saveSourceStatus(sourceID, SourceStatusAvailable, "")
			}
		}(ctx, source.Id)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

func (l *Loader) setModelSourceID(model dbmodels.CatalogModel, sourceID string) {
	if model == nil {
		return
	}

	// Add a source_id property to the model's properties list.. the hard
	// way, because we use pointers to slices for some reason.

	props := model.GetProperties()
	if props == nil {
		if modelImpl, ok := model.(*dbmodels.CatalogModelImpl); ok {
			newProps := make([]mrmodels.Properties, 0, 1)
			modelImpl.Properties = &newProps
			props = &newProps
		} else {
			// Can't do anything with this.
			return
		}
	}

	for i := range *props {
		if (*props)[i].Name == "source_id" {
			// Already has a source_id, just update it
			(*props)[i].StringValue = &sourceID
			return
		}
	}

	*props = append(*props, mrmodels.NewStringProperty("source_id", sourceID, false))
}

func (l *Loader) removeModelsFromMissingSources() error {
	enabledSourceIDs := mapset.NewSet[string]()
	allSourceIDs := mapset.NewSet[string]()
	for id, source := range l.Sources.AllSources() {
		allSourceIDs.Add(id)
		if source.Enabled == nil || *source.Enabled {
			enabledSourceIDs.Add(id)
		}
	}

	existingSourceIDs, err := l.services.CatalogModelRepository.GetDistinctSourceIDs()
	if err != nil {
		return fmt.Errorf("unable to retrieve existing source IDs: %w", err)
	}

	for oldSource := range mapset.NewSet(existingSourceIDs...).Difference(enabledSourceIDs).Iter() {
		glog.Infof("Removing models from source %s", oldSource)

		err = l.services.CatalogModelRepository.DeleteBySource(oldSource)
		if err != nil {
			return fmt.Errorf("unable to remove models from source %q: %w", oldSource, err)
		}

		// If the source is completely gone from config (not just disabled), remove its status too
		if !allSourceIDs.Contains(oldSource) {
			glog.Infof("Removing status for source %s (no longer in config)", oldSource)
			if delErr := l.services.CatalogSourceRepository.Delete(oldSource); delErr != nil {
				glog.Errorf("failed to delete status for source %s: %v", oldSource, delErr)
			}
		}
	}

	// Also clean up CatalogSource records for sources that are no longer in config
	// This handles sources that never loaded models (e.g., error sources, disabled sources)
	if err := l.cleanupOrphanedCatalogSources(allSourceIDs); err != nil {
		glog.Errorf("failed to cleanup orphaned catalog sources: %v", err)
	}

	return nil
}

// cleanupOrphanedCatalogSources removes CatalogSource records for sources that are no longer in the config.
func (l *Loader) cleanupOrphanedCatalogSources(currentSourceIDs mapset.Set[string]) error {
	existingSources, err := l.services.CatalogSourceRepository.GetAll()
	if err != nil {
		return fmt.Errorf("unable to get existing catalog sources: %w", err)
	}

	for _, source := range existingSources {
		attrs := source.GetAttributes()
		if attrs == nil || attrs.Name == nil {
			continue
		}

		sourceID := *attrs.Name
		if !currentSourceIDs.Contains(sourceID) {
			glog.Infof("Removing orphaned catalog source %s (no longer in config)", sourceID)
			if delErr := l.services.CatalogSourceRepository.Delete(sourceID); delErr != nil {
				glog.Errorf("failed to delete orphaned catalog source %s: %v", sourceID, delErr)
			}
		}
	}

	return nil
}

func (l *Loader) removeOrphanedModelsFromSource(sourceID string, valid mapset.Set[string]) (int, error) {
	list, err := l.services.CatalogModelRepository.List(dbmodels.CatalogModelListOptions{
		SourceIDs: &[]string{sourceID},
	})
	if err != nil {
		return 0, fmt.Errorf("unable to list models from source %q: %w", sourceID, err)
	}

	count := 0
	for _, model := range list.Items {
		attr := model.GetAttributes()
		if attr == nil || attr.Name == nil || model.GetID() == nil {
			continue
		}

		if valid.Contains(*attr.Name) {
			continue
		}

		glog.Infof("Removing %s model %s", sourceID, *attr.Name)

		err = l.services.CatalogModelRepository.DeleteByID(*model.GetID())
		if err != nil {
			return count, fmt.Errorf("unable to remove model %d (%s from source %s): %w", *model.GetID(), *attr.Name, sourceID, err)
		}
		count++
	}

	return count, nil
}

// saveSourceStatus persists the operational status of a source to the database.
// This allows status to be consistent across multiple pods.
func (l *Loader) saveSourceStatus(sourceID, status string, errorMsg string) {
	// Validate status is a valid enum value
	switch status {
	case SourceStatusAvailable, SourceStatusPartiallyAvailable, SourceStatusError, SourceStatusDisabled:
		// valid
	default:
		glog.Errorf("invalid status %q for source %s", status, sourceID)
		return
	}

	source := &dbmodels.CatalogSourceImpl{
		Attributes: &dbmodels.CatalogSourceAttributes{
			Name: &sourceID,
		},
	}

	// Set status property
	props := []mrmodels.Properties{
		mrmodels.NewStringProperty("status", status, false),
	}

	// Only store error property when non-empty
	if errorMsg != "" {
		props = append(props, mrmodels.NewStringProperty("error", errorMsg, false))
	}

	source.Properties = &props

	_, err := l.services.CatalogSourceRepository.Save(source)
	if err != nil {
		glog.Errorf("failed to save status for source %s: %v", sourceID, err)
	}
}
