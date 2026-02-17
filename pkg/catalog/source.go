package catalog

import (
	"maps"
	"slices"
	"strings"
	"sync"
)

// Source represents a catalog data source configuration.
// A source defines where to fetch entities from and how to configure the provider.
type Source struct {
	// ID is the unique identifier for this source.
	ID string

	// Name is the human-readable display name for this source.
	Name string

	// Type identifies the provider type to use (e.g., "yaml", "http").
	Type string

	// Enabled indicates whether this source should be loaded.
	// Defaults to true if nil.
	Enabled *bool

	// Labels are tags used for filtering and categorization.
	Labels []string

	// Properties contains provider-specific configuration.
	Properties map[string]any

	// IncludedItems are glob patterns for items to include.
	// If non-empty, only items matching at least one pattern are included.
	IncludedItems []string

	// ExcludedItems are glob patterns for items to exclude.
	// Items matching any pattern are excluded, even if they match an include pattern.
	ExcludedItems []string

	// Origin is the absolute path of the config file this source was loaded from.
	// This is used for resolving relative paths in Properties.
	Origin string
}

// IsEnabled returns true if this source is enabled (defaults to true if nil).
func (s Source) IsEnabled() bool {
	return s.Enabled == nil || *s.Enabled
}

// originEntry holds sources from a single origin (config file).
type originEntry struct {
	origin  string
	sources map[string]Source
}

// SourceCollection manages catalog sources from multiple origins with priority-based merging.
// Later entries in the origin order take precedence over earlier ones.
type SourceCollection struct {
	mu      sync.RWMutex
	entries []originEntry
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
		entries: entries,
	}
}

// Merge adds sources from one origin, completely replacing anything that was
// previously from that origin.
//
// If a source with the same ID exists in multiple origins, fields from
// higher-priority origins (listed later in entries) override fields from
// lower-priority origins. Fields that are not set (zero value for strings,
// nil for pointers/slices/maps) in the override are inherited from the base.
func (sc *SourceCollection) Merge(origin string, sources map[string]Source) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

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

// mergeSources performs field-level merging of two Source structs.
// Fields from 'override' take precedence over 'base' when they are explicitly set.
func mergeSources(base, override Source) Source {
	result := base

	// ID is always taken from override (it's the key)
	result.ID = override.ID

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

	// IncludedItems: override if non-nil
	if override.IncludedItems != nil {
		result.IncludedItems = override.IncludedItems
	}

	// ExcludedItems: override if non-nil
	if override.ExcludedItems != nil {
		result.ExcludedItems = override.ExcludedItems
	}

	// Type: override if non-empty
	if override.Type != "" {
		result.Type = override.Type
	}

	// Properties: override if non-nil (complete replacement, not deep merge)
	if override.Properties != nil {
		result.Properties = override.Properties
	}

	// Origin: use override's origin if Properties are overridden
	if override.Properties != nil && override.Origin != "" {
		result.Origin = override.Origin
	}

	return result
}

// applyDefaults applies default values to a Source for fields that are not set.
func applyDefaults(source Source) Source {
	// Default Enabled to true if not set
	if source.Enabled == nil {
		enabled := true
		source.Enabled = &enabled
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

// AllSources returns all merged sources.
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

// EnabledSources returns only sources that are enabled.
func (sc *SourceCollection) EnabledSources() map[string]Source {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	result := map[string]Source{}
	for id, source := range sc.merged() {
		if source.IsEnabled() {
			result[id] = source
		}
	}
	return result
}

// Get returns a source by ID if it exists and is enabled.
func (sc *SourceCollection) Get(id string) (Source, bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	merged := sc.merged()
	source, exists := merged[id]
	if !exists {
		return Source{}, false
	}

	// Only return if enabled
	if source.IsEnabled() {
		return source, true
	}
	return Source{}, false
}

// ByLabel returns enabled sources that have any of the labels provided.
// The matching is case insensitive.
//
// If a label is "null", every source without a label is returned.
func (sc *SourceCollection) ByLabel(labels []string) []Source {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	labelMap := make(map[string]struct{}, len(labels))
	for _, label := range labels {
		labelMap[strings.ToLower(label)] = struct{}{}
	}

	matches := map[string]Source{}
	sources := sc.merged()

	if _, hasNull := labelMap["null"]; hasNull {
		for _, source := range sources {
			if !source.IsEnabled() {
				continue
			}
			if len(source.Labels) == 0 {
				matches[source.ID] = source
			}
		}
	}

OUTER:
	for _, source := range sources {
		if !source.IsEnabled() {
			continue
		}
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[source.ID] = source
				continue OUTER
			}
		}
	}

	return slices.Collect(maps.Values(matches))
}

// IDs returns all source IDs (both enabled and disabled).
func (sc *SourceCollection) IDs() []string {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	merged := sc.merged()
	ids := make([]string, 0, len(merged))
	for id := range merged {
		ids = append(ids, id)
	}
	return ids
}
