package catalog

import (
	"maps"
	"slices"
	"strings"
	"sync"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
)

// originEntry holds sources from a single origin (config file).
type originEntry struct {
	origin  string
	sources map[string]Source
}

// SourceCollection manages catalog sources from multiple origins with priority-based merging.
// Later entries in the slice take precedence over earlier ones.
type SourceCollection struct {
	mu           sync.RWMutex
	entries      []originEntry
	namedQueries map[string]map[string]FieldFilter
}

// NewSourceCollection creates a new SourceCollection with the given origin order.
// Origins listed later in the order take precedence over earlier ones.
// For example, if originOrder is ["default.yaml", "user.yaml"], sources from
// "user.yaml" will override sources with the same ID from "default.yaml".
func NewSourceCollection(originOrder ...string) *SourceCollection {
	entries := make([]originEntry, len(originOrder))
	for i, origin := range originOrder {
		entries[i] = originEntry{origin: origin, sources: nil}
	}
	return &SourceCollection{
		entries:      entries,
		namedQueries: make(map[string]map[string]FieldFilter),
	}
}

// Merge adds sources from one origin (ordinarily, a file path--but any unique
// string will do), completely replacing anything that was previously from that
// origin.
//
// If a source with the same ID exists in multiple origins, fields from
// higher-priority origins (listed later in entries) override fields from
// lower-priority origins. Fields that are not set (zero value for strings,
// nil for pointers/slices/maps) in the override are inherited from the base.
func (sc *SourceCollection) Merge(origin string, sources map[string]Source) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	return sc.mergeSourcesInternal(origin, sources)
}

// MergeWithNamedQueries adds sources and named queries from one origin.
func (sc *SourceCollection) MergeWithNamedQueries(origin string, sources map[string]Source, namedQueries map[string]map[string]FieldFilter) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Merge sources using existing logic
	if err := sc.mergeSourcesInternal(origin, sources); err != nil {
		return err
	}

	// Merge named queries (later origins override earlier ones)
	for queryName, fieldFilters := range namedQueries {
		if sc.namedQueries[queryName] == nil {
			sc.namedQueries[queryName] = make(map[string]FieldFilter)
		}
		for fieldName, filter := range fieldFilters {
			sc.namedQueries[queryName][fieldName] = filter
		}
	}

	return nil
}

// mergeSourcesInternal extracts the internal logic from Merge
func (sc *SourceCollection) mergeSourcesInternal(origin string, sources map[string]Source) error {
	// Find existing entry for this origin
	for i := range sc.entries {
		if sc.entries[i].origin == origin {
			sc.entries[i].sources = sources
			return nil
		}
	}

	// Origin not found, append it (dynamic registration)
	sc.entries = append(sc.entries, originEntry{origin: origin, sources: sources})
	return nil
}

// GetNamedQueries returns all merged named queries
func (sc *SourceCollection) GetNamedQueries() map[string]map[string]FieldFilter {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Return a copy to prevent external modification
	result := make(map[string]map[string]FieldFilter, len(sc.namedQueries))
	for queryName, fieldFilters := range sc.namedQueries {
		result[queryName] = make(map[string]FieldFilter, len(fieldFilters))
		for fieldName, filter := range fieldFilters {
			result[queryName][fieldName] = filter
		}
	}
	return result
}

// mergeSources performs field-level merging of two Source structs.
// Fields from 'override' take precedence over 'base' when they are explicitly set.
// A field is considered "set" if:
// - For strings: non-empty
// - For pointers: non-nil
// - For slices: non-nil (empty slice is considered explicitly set to "no items")
// - For maps: non-nil (empty map is considered explicitly set)
func mergeSources(base, override Source) Source {
	result := base

	// Id is always taken from override (it's the key)
	result.Id = override.Id

	// Name: override if non-empty
	if override.Name != "" {
		result.Name = override.Name
	}

	// Enabled: override if non-nil
	if override.Enabled != nil {
		result.Enabled = override.Enabled
	}

	// Labels: override if non-nil (empty slice means "explicitly no labels")
	if override.Labels != nil {
		result.Labels = override.Labels
	}

	// IncludedModels: override if non-nil
	if override.IncludedModels != nil {
		result.IncludedModels = override.IncludedModels
	}

	// ExcludedModels: override if non-nil
	if override.ExcludedModels != nil {
		result.ExcludedModels = override.ExcludedModels
	}

	// Type: override if non-empty
	if override.Type != "" {
		result.Type = override.Type
	}

	// Properties: override if non-nil (complete replacement, not deep merge)
	if override.Properties != nil {
		result.Properties = override.Properties
	}

	// Origin: use override's origin if Properties are overridden (since relative
	// paths in Properties should resolve relative to where they were defined).
	// Otherwise, keep base origin (where Type and original Properties came from).
	if override.Properties != nil && override.Origin != "" {
		result.Origin = override.Origin
	}

	return result
}

// applyDefaults applies default values to an Source for fields that are not set.
func applyDefaults(source Source) Source {
	// Default Enabled to true if not set
	if source.Enabled == nil {
		source.Enabled = apiutils.Of(true)
	}

	// Default Labels to empty slice if not set
	if source.Labels == nil {
		source.Labels = []string{}
	}

	return source
}

// merged computes the merged view of all sources with field-level merging.
// Must be called with lock held.
func (sc *SourceCollection) merged() map[string]Source {
	result := map[string]Source{}

	for _, entry := range sc.entries {
		for id, source := range entry.sources {
			if existing, ok := result[id]; ok {
				// Field-level merge: existing is base, source is override
				result[id] = mergeSources(existing, source)
			} else {
				result[id] = source
			}
		}
	}

	// Apply defaults to all merged sources
	for id, source := range result {
		result[id] = applyDefaults(source)
	}

	return result
}

// AllSources returns all merged sources including Type and Properties.
// This is used by the loader to get complete source information.
// All sources are returned regardless of enabled status.
func (sc *SourceCollection) AllSources() map[string]Source {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := map[string]Source{}
	for id, source := range sc.merged() {
		result[id] = source
	}
	return result
}

// All returns all sources as CatalogSource (for the API).
// This excludes internal fields like Type and Properties.
func (sc *SourceCollection) All() map[string]model.CatalogSource {
	result := map[string]model.CatalogSource{}
	for id, source := range sc.AllSources() {
		result[id] = source.CatalogSource
	}
	return result
}

// Get returns a source by name if it exists and is enabled.
func (sc *SourceCollection) Get(name string) (src model.CatalogSource, ok bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Get from merged view (which includes field-level merging and defaults)
	merged := sc.merged()
	source, exists := merged[name]
	if !exists {
		return model.CatalogSource{}, false
	}

	// Only return if enabled
	if source.Enabled != nil && *source.Enabled {
		return source.CatalogSource, true
	}
	return model.CatalogSource{}, false
}

// ByLabel returns enabled sources that have any of the labels provided. The matching
// is case insensitive.
//
// If a label is "null", every source without a label is returned.
func (sc *SourceCollection) ByLabel(labels []string) []model.CatalogSource {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	labelMap := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		labelMap[strings.ToLower(label)] = struct{}{}
	}

	matches := map[string]model.CatalogSource{}
	sources := sc.merged()

	if _, hasNull := labelMap["null"]; hasNull {
		for _, source := range sources {
			// Skip disabled sources
			if source.Enabled == nil || !*source.Enabled {
				continue
			}
			if len(source.Labels) == 0 {
				matches[source.Id] = source.CatalogSource
			}
		}
	}

OUTER:
	for _, source := range sources {
		// Skip disabled sources
		if source.Enabled == nil || !*source.Enabled {
			continue
		}
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[source.Id] = source.CatalogSource
				continue OUTER
			}
		}
	}

	return slices.Collect(maps.Values(matches))
}
