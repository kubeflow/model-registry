package catalog

import (
	"fmt"
	"maps"
	"slices"
	"strings"
	"sync"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type SourceCollection struct {
	mu sync.RWMutex

	// origins keeps track of where a source came from by some name (intended to be a file path).
	origins map[string][]string

	sources map[string]model.CatalogSource
}

func NewSourceCollection() *SourceCollection {
	return &SourceCollection{
		origins: map[string][]string{},
		sources: map[string]model.CatalogSource{},
	}
}

// Merge adds sources from one origin (ordinarily, a file path--but any unique
// string will do), completely replacing anything that was previously from that
// origin.
func (sc *SourceCollection) Merge(origin string, sources map[string]model.CatalogSource) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	// Remove everything that was set before for this origin so that
	// unreferenced sources are deleted.
	for _, id := range sc.origins[origin] {
		delete(sc.sources, id)
	}
	sc.origins[origin] = slices.Collect(maps.Keys(sources))

	for sourceID, source := range sources {
		// Everything was deleted above, so if there's a source that
		// already exists it must have come from another origin (file).
		if _, exists := sc.sources[sourceID]; exists {
			return fmt.Errorf("source %s exists from multiple origins", sourceID)
		}

		sc.sources[sourceID] = source
	}

	return nil
}

func (sc *SourceCollection) All() map[string]model.CatalogSource {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	return sc.sources
}

func (sc *SourceCollection) Get(name string) (src model.CatalogSource, ok bool) {
	sc.mu.RLock()
	defer sc.mu.RUnlock()

	src, ok = sc.sources[name]
	return
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

	if _, hasNull := labelMap["null"]; hasNull {
		for _, source := range sc.sources {
			if len(source.Labels) == 0 {
				matches[source.Id] = source
			}
		}
	}

OUTER:
	for _, source := range sc.sources {
		for _, label := range source.Labels {
			if _, match := labelMap[strings.ToLower(label)]; match {
				matches[source.Id] = source
				continue OUTER
			}
		}
	}

	return slices.Collect(maps.Values(matches))
}
