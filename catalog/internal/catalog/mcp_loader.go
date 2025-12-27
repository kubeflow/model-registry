package catalog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/golang/glog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/catalog/internal/mcp"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
	"k8s.io/apimachinery/pkg/util/yaml"
)

// McpLoaderEventHandler is the definition of a function called after an MCP server is loaded.
type McpLoaderEventHandler func(ctx context.Context, record mcp.McpServerProviderRecord) error

// mcpSourceConfig is the structure for the MCP catalog sources YAML file.
type mcpSourceConfig struct {
	Catalogs     []mcp.McpSource                   `json:"catalogs" yaml:"catalogs"`
	NamedQueries map[string]map[string]FieldFilter `json:"namedQueries,omitempty" yaml:"namedQueries,omitempty"`
}

// McpLoader handles loading MCP servers from YAML sources into the database.
type McpLoader struct {
	paths         []string
	services      service.Services
	sources       *SourceCollection // shared source collection for unified sources API
	closersMu     sync.Mutex
	closer        func() // cancels the current MCP loading goroutines
	handlers      []McpLoaderEventHandler
	loadedSources map[string]bool                    // tracks which source IDs have been loaded
	namedQueries  map[string]map[string]FieldFilter  // merged named queries from all config files
}

// NewMcpLoader creates a new MCP loader.
func NewMcpLoader(services service.Services, paths []string, sources *SourceCollection) *McpLoader {
	// Convert paths to absolute for consistent origin ordering.
	absPaths := make([]string, 0, len(paths))
	for _, p := range paths {
		absPath, err := filepath.Abs(p)
		if err != nil {
			// Fall back to original path if conversion fails
			absPath = p
		}
		absPaths = append(absPaths, absPath)
	}

	return &McpLoader{
		paths:         absPaths,
		services:      services,
		sources:       sources,
		loadedSources: map[string]bool{},
		namedQueries:  make(map[string]map[string]FieldFilter),
	}
}

// RegisterEventHandler adds a function that will be called for every
// successfully processed record. This should be called before Start.
func (l *McpLoader) RegisterEventHandler(fn McpLoaderEventHandler) {
	l.handlers = append(l.handlers, fn)
}

// Start processes the MCP sources YAML files. Background goroutines will be
// stopped when the context is canceled.
func (l *McpLoader) Start(ctx context.Context) error {
	// Phase 1: Parse all config files and merge sources with field-level priority.
	// Sources from later config files override fields from earlier ones.
	allSources, err := l.readAndMergeSources()
	if err != nil {
		return fmt.Errorf("failed to read MCP sources: %w", err)
	}

	// Phase 1.5: Read and merge named queries from all config files
	l.readAndMergeNamedQueries()

	if len(allSources) == 0 {
		glog.Infof("No MCP catalog sources found")
		return nil
	}

	// Delete MCP servers from unknown or disabled sources
	err = l.removeMcpServersFromMissingSources(allSources)
	if err != nil {
		return fmt.Errorf("failed to remove MCP servers from missing sources: %w", err)
	}

	// Merge MCP sources and named queries into the shared SourceCollection for unified /sources API
	l.mergeMcpSourcesIntoCollection(allSources)

	// Phase 2: Load MCP servers from sources
	err = l.loadAllMcpServers(ctx, allSources)
	if err != nil {
		return err
	}

	// Phase 3: Set up file watchers for hot-reload
	for _, path := range l.paths {
		go func(path string) {
			changes, err := getMonitor().Path(ctx, path)
			if err != nil {
				glog.Errorf("unable to watch MCP sources file (%s): %v", path, err)
				return
			}

			for range changes {
				glog.Infof("Reloading MCP sources %s", path)
				l.reloadAll(ctx)
			}
		}(path)
	}

	return nil
}

// read reads an MCP source configuration file.
func (l *McpLoader) read(path string) ([]mcp.McpSource, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &mcpSourceConfig{}
	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	return config.Catalogs, nil
}

// readConfig reads an MCP source configuration file and returns both sources and named queries.
func (l *McpLoader) readConfig(path string) (*mcpSourceConfig, error) {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	config := &mcpSourceConfig{}
	if err = yaml.UnmarshalStrict(bytes, &config); err != nil {
		return nil, err
	}

	return config, nil
}

// loadAllMcpServers loads MCP servers from all sources.
func (l *McpLoader) loadAllMcpServers(ctx context.Context, sources map[string]mcp.McpSource) error {
	// Clear the loaded sources tracker for a fresh load
	l.loadedSources = map[string]bool{}

	return l.updateDatabase(ctx, sources)
}

// reloadAll re-parses all config files and reloads all MCP servers.
func (l *McpLoader) reloadAll(ctx context.Context) {
	// Use the same merge logic as Start()
	allSources, err := l.readAndMergeSources()
	if err != nil {
		glog.Errorf("unable to read and merge MCP sources: %v", err)
		return
	}

	// Re-read named queries on reload
	l.readAndMergeNamedQueries()

	// Merge into shared SourceCollection for unified /sources API
	l.mergeMcpSourcesIntoCollection(allSources)

	// Clean up MCP servers that are no longer in config
	if err := l.removeMcpServersFromMissingSources(allSources); err != nil {
		glog.Errorf("unable to remove MCP servers from missing sources: %v", err)
	}

	// Reload all MCP servers
	if err := l.loadAllMcpServers(ctx, allSources); err != nil {
		glog.Errorf("unable to reload MCP servers: %v", err)
	}
}

// readAndMergeSources reads MCP sources from all config paths and merges sources
// with the same ID using field-level priority. Sources from later paths override
// fields from earlier paths.
func (l *McpLoader) readAndMergeSources() (map[string]mcp.McpSource, error) {
	return MergeMcpSourcesFromPaths(l.paths, l.readWithWarning)
}

// readAndMergeNamedQueries reads named queries from all config files and merges them.
// Named queries from later paths override those from earlier paths.
func (l *McpLoader) readAndMergeNamedQueries() {
	// Clear existing named queries
	l.namedQueries = make(map[string]map[string]FieldFilter)

	// Read from all paths in priority order (later paths override earlier)
	for _, path := range l.paths {
		config, err := l.readConfig(path)
		if err != nil {
			glog.V(2).Infof("Skipping named queries from %s: %v", path, err)
			continue
		}

		if config.NamedQueries == nil {
			continue
		}

		// Merge named queries (later files override earlier ones)
		for queryName, fieldFilters := range config.NamedQueries {
			if l.namedQueries[queryName] == nil {
				l.namedQueries[queryName] = make(map[string]FieldFilter)
			}
			for fieldName, filter := range fieldFilters {
				l.namedQueries[queryName][fieldName] = filter
			}
		}
	}

	if len(l.namedQueries) > 0 {
		glog.Infof("Loaded %d MCP named queries", len(l.namedQueries))
	}
}

// GetNamedQueries returns all merged named queries for MCP servers.
func (l *McpLoader) GetNamedQueries() map[string]map[string]FieldFilter {
	// Return a copy to prevent external modification
	result := make(map[string]map[string]FieldFilter, len(l.namedQueries))
	for queryName, fieldFilters := range l.namedQueries {
		result[queryName] = make(map[string]FieldFilter, len(fieldFilters))
		for fieldName, filter := range fieldFilters {
			result[queryName][fieldName] = filter
		}
	}
	return result
}

// readWithWarning reads an MCP source configuration file and logs warnings on failure.
func (l *McpLoader) readWithWarning(path string) ([]mcp.McpSource, error) {
	sources, err := l.read(path)
	if err != nil {
		glog.Warningf("MCP catalog source file %s not found or invalid: %v", path, err)
		return nil, err
	}
	return sources, nil
}

// updateDatabase loads MCP servers into the database.
func (l *McpLoader) updateDatabase(ctx context.Context, sources map[string]mcp.McpSource) error {
	ctx, cancel := context.WithCancel(ctx)

	l.closersMu.Lock()
	if l.closer != nil {
		l.closer()
	}
	l.closer = cancel
	l.closersMu.Unlock()

	records := l.readProviderRecords(ctx, sources)

	go func() {
		for record := range records {
			if record.Server == nil {
				continue
			}
			attr := record.Server.GetAttributes()
			if attr == nil || attr.Name == nil {
				continue
			}

			glog.Infof("Loading MCP server %s with %d tool(s)", *attr.Name, len(record.Tools))

			_, err := l.services.McpServerRepository.Save(record.Server)
			if err != nil {
				glog.Errorf("%s: unable to save MCP server: %v", *attr.Name, err)
				continue
			}

			for _, handler := range l.handlers {
				handler(ctx, record)
			}
		}
	}()

	return nil
}

// readProviderRecords calls the provider for every source and merges the returned channels together.
func (l *McpLoader) readProviderRecords(ctx context.Context, sources map[string]mcp.McpSource) <-chan mcp.McpServerProviderRecord {
	ch := make(chan mcp.McpServerProviderRecord)
	var wg sync.WaitGroup

	for _, source := range sources {
		// Skip disabled sources
		if source.Enabled != nil && !*source.Enabled {
			continue
		}

		// Skip sources that have already been loaded
		if l.loadedSources[source.Id] {
			continue
		}

		if source.Type == "" {
			glog.Errorf("MCP source %s has no type defined, skipping", source.Id)
			continue
		}

		// Mark this source as loaded
		l.loadedSources[source.Id] = true

		// Build the server filter for this source
		serverFilter, err := NewMcpServerFilterFromSource(&source)
		if err != nil {
			glog.Errorf("MCP source %s has invalid filter configuration: %v", source.Id, err)
			continue
		}

		glog.Infof("Reading MCP servers from %s source %s", source.Type, source.Id)

		providerFunc, ok := mcp.RegisteredMcpProviders[source.Type]
		if !ok {
			glog.Errorf("MCP catalog type %s not registered", source.Type)
			continue
		}

		// Use the source's origin directory for resolving relative paths.
		sourceDir := filepath.Dir(source.Origin)
		sourceCopy := source

		records, err := providerFunc(ctx, &sourceCopy, sourceDir)
		if err != nil {
			glog.Errorf("error reading MCP catalog type %s with id %s: %v", source.Type, source.Id, err)
			continue
		}

		wg.Add(1)
		go func(ctx context.Context, sourceID string, filter *McpServerFilter) {
			defer wg.Done()

			serverNames := []string{}

			for r := range records {
				if r.Server == nil {
					glog.V(2).Infof("%s: MCP servers batch complete", sourceID)

					// Copy the list of server names, then clear it.
					serverNameSet := mapset.NewSet(serverNames...)
					serverNames = serverNames[:0]

					go func() {
						err := l.removeOrphanedMcpServersFromSource(sourceID, serverNameSet)
						if err != nil {
							glog.Errorf("error removing orphaned MCP servers: %v", err)
						}
					}()
					continue
				}

				// Apply server filter
				attr := r.Server.GetAttributes()
				if attr != nil && attr.Name != nil {
					serverName := *attr.Name
					if !filter.Allows(serverName) {
						glog.V(2).Infof("%s: MCP server %s excluded by filter", sourceID, serverName)
						continue
					}
					serverNames = append(serverNames, serverName)
				}

				// Set source_id on every returned server.
				l.setMcpServerSourceID(r.Server, sourceID)

				ch <- r
			}
		}(ctx, source.Id, serverFilter)
	}

	go func() {
		defer close(ch)
		wg.Wait()
	}()

	return ch
}

// setMcpServerSourceID adds the source_id property to an MCP server.
func (l *McpLoader) setMcpServerSourceID(server dbmodels.McpServer, sourceID string) {
	if server == nil {
		return
	}

	props := server.GetProperties()
	if props == nil {
		if serverImpl, ok := server.(*dbmodels.McpServerImpl); ok {
			newProps := make([]mrmodels.Properties, 0, 1)
			serverImpl.Properties = &newProps
			props = &newProps
		} else {
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

// removeMcpServersFromMissingSources removes MCP servers from sources that are no longer in config or disabled.
func (l *McpLoader) removeMcpServersFromMissingSources(sources map[string]mcp.McpSource) error {
	enabledSourceIDs := mapset.NewSet[string]()
	for id, source := range sources {
		if source.Enabled == nil || *source.Enabled {
			enabledSourceIDs.Add(id)
		}
	}

	existingSourceIDs, err := l.services.McpServerRepository.GetDistinctSourceIDs()
	if err != nil {
		return fmt.Errorf("unable to retrieve existing MCP source IDs: %w", err)
	}

	for oldSource := range mapset.NewSet(existingSourceIDs...).Difference(enabledSourceIDs).Iter() {
		glog.Infof("Removing MCP servers from source %s", oldSource)

		err = l.services.McpServerRepository.DeleteBySource(oldSource)
		if err != nil {
			return fmt.Errorf("unable to remove MCP servers from source %q: %w", oldSource, err)
		}
	}

	return nil
}

// removeOrphanedMcpServersFromSource removes MCP servers that are no longer in the source.
func (l *McpLoader) removeOrphanedMcpServersFromSource(sourceID string, valid mapset.Set[string]) error {
	list, err := l.services.McpServerRepository.List(dbmodels.McpServerListOptions{
		SourceIDs: &[]string{sourceID},
	})
	if err != nil {
		return fmt.Errorf("unable to list MCP servers from source %q: %w", sourceID, err)
	}

	for _, server := range list.Items {
		attr := server.GetAttributes()
		if attr == nil || attr.Name == nil || server.GetID() == nil {
			continue
		}

		if valid.Contains(*attr.Name) {
			continue
		}

		glog.Infof("Removing %s MCP server %s", sourceID, *attr.Name)

		err = l.services.McpServerRepository.DeleteByID(*server.GetID())
		if err != nil {
			return fmt.Errorf("unable to remove MCP server %d (%s from source %s): %w", *server.GetID(), *attr.Name, sourceID, err)
		}
	}

	return nil
}

// mergeMcpSourcesIntoCollection converts MCP sources to catalog Sources and merges them
// into the shared SourceCollection. This enables the unified /sources API to return
// MCP sources with assetType=mcp_servers. It also merges MCP named queries.
func (l *McpLoader) mergeMcpSourcesIntoCollection(mcpSources map[string]mcp.McpSource) {
	if l.sources == nil {
		return
	}

	// Convert MCP sources to catalog Sources
	catalogSources := make(map[string]Source)
	for id, mcpSource := range mcpSources {
		enabled := mcpSource.Enabled
		if enabled == nil {
			defaultEnabled := true
			enabled = &defaultEnabled
		}

		catalogSources[id] = Source{
			CatalogSource: model.CatalogSource{
				Id:      mcpSource.Id,
				Name:    mcpSource.Name,
				Labels:  mcpSource.Labels,
				Enabled: enabled,
			},
			Type:              mcpSource.Type,
			Properties:        mcpSource.Properties,
			Origin:            mcpSource.Origin,
			DetectedAssetType: AssetTypeMcpServers,
		}
	}

	// Merge sources and named queries into the shared source collection
	if len(l.paths) > 0 {
		if err := l.sources.MergeWithNamedQueries(l.paths[0], catalogSources, l.namedQueries); err != nil {
			glog.Errorf("unable to merge MCP sources into source collection: %v", err)
		}
	}
}
