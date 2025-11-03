package catalog

import (
	"sync"
)

type LabelCollection struct {
	mu sync.RWMutex

	// origins keeps track of which labels came from which origin (file path).
	// Each origin maps to the indices of labels it contributed.
	origins map[string][]int

	// labels stores all unique labels
	labels []map[string]string
}

func NewLabelCollection() *LabelCollection {
	return &LabelCollection{
		origins: map[string][]int{},
		labels:  []map[string]string{},
	}
}

// Merge adds labels from one origin (ordinarily, a file path), completely
// replacing anything that was previously from that origin.
func (lc *LabelCollection) Merge(origin string, newLabels []map[string]string) {
	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Remove labels that were previously set for this origin
	if oldIndices, exists := lc.origins[origin]; exists {
		// Mark old labels for removal by setting them to nil
		for _, idx := range oldIndices {
			if idx < len(lc.labels) {
				lc.labels[idx] = nil
			}
		}
	}

	// Compact the slice by removing nil entries
	compacted := make([]map[string]string, 0, len(lc.labels))
	for _, label := range lc.labels {
		if label != nil {
			compacted = append(compacted, label)
		}
	}
	lc.labels = compacted

	// Add new labels from this origin (only if they don't already exist)
	newIndices := make([]int, 0, len(newLabels))
	for _, newLabel := range newLabels {
		if !containsLabel(lc.labels, newLabel) {
			lc.labels = append(lc.labels, newLabel)
			newIndices = append(newIndices, len(lc.labels)-1)
		}
	}

	if len(newIndices) > 0 {
		lc.origins[origin] = newIndices
	} else {
		delete(lc.origins, origin)
	}
}

func (lc *LabelCollection) All() []map[string]string {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	result := make([]map[string]string, len(lc.labels))
	copy(result, lc.labels)
	return result
}

// containsLabel checks if a label map already exists in the labels slice.
func containsLabel(labels []map[string]string, target map[string]string) bool {
	for _, label := range labels {
		if mapsEqual(label, target) {
			return true
		}
	}
	return false
}

// mapsEqual compares two string maps for equality.
func mapsEqual(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		if bv, ok := b[k]; !ok || bv != v {
			return false
		}
	}
	return true
}
