package catalog

import (
	"maps"
	"slices"
	"strings"
	"sync"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// originEntry holds sources from a single origin (config file).
type originEntry struct {
	origin  string
	sources map[string]model.CatalogSource
}

// SourceCollection manages catalog sources from multiple origins with priority-based merging.
// Later entries in the slice take precedence over earlier ones.
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
	return &SourceCollection{entries: entries}
}

// Merge adds sources from one origin (ordinarily, a file path--but any unique
// string will do), completely replacing anything that was previously from that
// origin.
//
// If a source with the same ID exists in multiple origins, the source from
// the origin with higher priority (listed later in entries) takes precedence.
func (sc *SourceCollection) Merge(origin string, sources map[string]model.CatalogSource) error {
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

// merged computes the merged view of all sources. Must be called with lock held.
func (sc *SourceCollection) merged() map[string]model.CatalogSource {
	result := map[string]model.CatalogSource{}
	for _, entry := range sc.entries {
		for id, source := range entry.sources {
			result[id] = source
		}
	}
	return result
}

func (sc *SourceCollection) All() map[string]model.CatalogSource {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.merged()
}

func (sc *SourceCollection) Get(name string) (src model.CatalogSource, ok bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	// Iterate in reverse to find highest-priority match
	for i := len(sc.entries) - 1; i >= 0; i-- {
		if source, exists := sc.entries[i].sources[name]; exists {
			return source, true
		}
	}
	return model.CatalogSource{}, false
}

// ByLabel returns sources that have any of the labels provided. The matching
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
			if len(source.Labels) == 0 {
				matches[source.Id] = source
			}
		}
	}

OUTER:
	for _, source := range sources {
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[source.Id] = source
				continue OUTER
			}
		}
	}

	return slices.Collect(maps.Values(matches))
}
