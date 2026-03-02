package catalog

import (
	"context"
	"fmt"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/glog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/mcpcatalog"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/modelcatalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
)

// Loader is the unified catalog loader that handles both model catalogs
// and MCP catalogs. It delegates to specialized loaders while sharing
// leader state, file watching, and write tracking.
type Loader struct {
	*basecatalog.BaseLoader

	modelLoader *modelcatalog.ModelLoader
	mcpLoader   *mcpcatalog.MCPLoader
}

// NewLoader creates a new unified catalog loader
func NewLoader(services service.Services, paths []string) *Loader {
	base := basecatalog.NewBaseLoader(paths)
	return &Loader{
		BaseLoader:  base,
		modelLoader: modelcatalog.NewModelLoader(services, base),
		mcpLoader:   mcpcatalog.NewMCPLoaderWithState(services, base),
	}
}

// Sources returns the model source collection (for API service)
func (l *Loader) Sources() *modelcatalog.SourceCollection {
	return l.modelLoader.Sources
}

// Labels returns the label collection (for API service)
func (l *Loader) Labels() *modelcatalog.LabelCollection {
	return l.modelLoader.Labels
}

// MCPSources returns the MCP source collection
func (l *Loader) MCPSources() *mcpcatalog.MCPSourceCollection {
	return l.mcpLoader.Sources
}

// RegisterEventHandler adds a model load event handler
func (l *Loader) RegisterEventHandler(fn modelcatalog.LoaderEventHandler) {
	l.modelLoader.RegisterEventHandler(fn)
}

// RegisterMCPEventHandler adds an MCP server load event handler
func (l *Loader) RegisterMCPEventHandler(fn mcpcatalog.MCPLoaderEventHandler) {
	l.mcpLoader.RegisterEventHandler(fn)
}

// StartReadOnly initializes the loader in read-only mode (standby pod).
// Parses all config files into in-memory collections and sets up file watchers.
func (l *Loader) StartReadOnly(ctx context.Context) error {
	watcherCtx, err := l.SetupFileWatchers(ctx)
	if err != nil {
		return err
	}

	glog.Info("Starting unified loader in read-only mode (standby)")

	// Parse configs for both loaders
	if err := l.modelLoader.ParseAllConfigs(); err != nil {
		return fmt.Errorf("model config: %w", err)
	}
	if err := l.mcpLoader.ParseAllConfigs(); err != nil {
		return fmt.Errorf("mcp config: %w", err)
	}

	// Start file watchers (single set, triggers both reloads)
	for _, path := range l.Paths() {
		go l.watchFile(watcherCtx, path)
	}

	glog.Info("Read-only mode initialized successfully")
	return nil
}

// collectAllSourceIDs returns the union of all model and MCP source IDs.
// This combined set is used to prevent cross-contamination when each loader
// cleans up shared CatalogSource records.
func (l *Loader) collectAllSourceIDs() mapset.Set[string] {
	combined := mapset.NewSet[string]()
	for id := range l.modelLoader.Sources.AllSources() {
		combined.Add(id)
	}
	for id := range l.mcpLoader.Sources.AllSources() {
		combined.Add(id)
	}
	return combined
}

// watchFile watches a single config file for changes and reloads both loaders
func (l *Loader) watchFile(ctx context.Context, path string) {
	changes, err := basecatalog.GetMonitor().Path(ctx, path)
	if err != nil {
		glog.Errorf("unable to watch file (%s): %v", path, err)
		return
	}

	for range changes {
		glog.Infof("Config file changed, reloading: %s", path)

		// Phase 1: Re-parse both loaders so source collections reflect current config.
		// Both must parse before we compute combined IDs to avoid missing each other's sources.
		l.modelLoader.ReloadParsing()
		l.mcpLoader.ReloadParsing()

		// Phase 2: Perform leader writes with combined source IDs to prevent cross-contamination.
		if l.ShouldWriteDatabase() {
			allKnownSourceIDs := l.collectAllSourceIDs()
			if err := l.modelLoader.PerformLeaderOperations(ctx, allKnownSourceIDs); err != nil {
				glog.Errorf("unable to perform model leader writes on reload: %v", err)
			}
			if err := l.mcpLoader.PerformLeaderOperations(ctx, allKnownSourceIDs); err != nil {
				glog.Errorf("unable to perform MCP leader writes on reload: %v", err)
			}
		}
	}
}

// StartLeader transitions the loader to leader mode and blocks until the
// context is cancelled (leadership lost or pod shutdown).
func (l *Loader) StartLeader(ctx context.Context) error {
	if l.IsLeader() {
		return fmt.Errorf("already in leader mode")
	}
	l.SetLeader(true)

	glog.Info("Transitioning to leader mode (read-write)")

	// Collect combined source IDs from both loaders to prevent cross-contamination
	// when each loader cleans up shared CatalogSource records.
	allKnownSourceIDs := l.collectAllSourceIDs()

	// Perform leader operations for both loaders
	if err := l.modelLoader.PerformLeaderOperations(ctx, allKnownSourceIDs); err != nil {
		l.SetLeader(false)
		return fmt.Errorf("model leader operations: %w", err)
	}

	if err := l.mcpLoader.PerformLeaderOperations(ctx, allKnownSourceIDs); err != nil {
		l.SetLeader(false)
		return fmt.Errorf("mcp leader operations: %w", err)
	}

	glog.Info("Leader mode active")

	// Wait for context cancellation
	<-ctx.Done()
	glog.Info("Leadership context cancelled, cleaning up...")

	l.WaitForInflightWrites(5 * time.Second)
	l.SetLeader(false)

	glog.Info("Leader mode stopped")
	return ctx.Err()
}

// Shutdown gracefully shuts down the loader
func (l *Loader) Shutdown() error {
	glog.Info("Shutting down unified loader...")
	l.StopFileWatchers()
	l.WaitForInflightWrites(10 * time.Second)
	glog.Info("Loader shutdown complete")
	return nil
}
