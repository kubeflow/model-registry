package catalog

import (
	"fmt"
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
// Returns an error if:
//   - duplicate label names exist within newLabels
//   - a label name conflicts with an existing label from a different origin
func (lc *LabelCollection) Merge(origin string, newLabels []map[string]string) error {
	newLabelNames := make(map[string]bool)
	for _, newLabel := range newLabels {
		if name, ok := newLabel["name"]; ok {
			if newLabelNames[name] {
				return fmt.Errorf("duplicate label name '%s' within the same origin", name)
			}
			newLabelNames[name] = true
		}
	}

	lc.mu.Lock()
	defer lc.mu.Unlock()

	// Build a map of existing label names from OTHER origins (excluding this origin)
	// This allows us to validate BEFORE mutating state
	oldIndices, originExists := lc.origins[origin]
	existingNamesFromOtherOrigins := make(map[string]bool)
	for i, label := range lc.labels {
		// Skip labels from this origin (they will be replaced)
		isFromThisOrigin := false
		if originExists {
			for _, idx := range oldIndices {
				if i == idx {
					isFromThisOrigin = true
					break
				}
			}
		}

		if !isFromThisOrigin {
			if name, ok := label["name"]; ok {
				existingNamesFromOtherOrigins[name] = true
			}
		}
	}

	// Validate conflicts and prepare labels to add in a single pass
	labelsToAdd := make([]map[string]string, 0, len(newLabels))
	for _, newLabel := range newLabels {
		// Check for conflicts with other origins
		if name, ok := newLabel["name"]; ok {
			if existingNamesFromOtherOrigins[name] {
				return fmt.Errorf("label with name '%s' already exists from another origin", name)
			}
		}

		labelsToAdd = append(labelsToAdd, newLabel)
	}

	// All validation passed, now proceed with mutation
	// Remove labels that were previously set for this origin
	if originExists {
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

	// Add the validated new labels
	newIndices := make([]int, 0, len(labelsToAdd))
	for _, newLabel := range labelsToAdd {
		lc.labels = append(lc.labels, newLabel)
		newIndices = append(newIndices, len(lc.labels)-1)
	}

	if len(newIndices) > 0 {
		lc.origins[origin] = newIndices
	} else {
		delete(lc.origins, origin)
	}

	return nil
}

func (lc *LabelCollection) All() []map[string]string {
	lc.mu.RLock()
	defer lc.mu.RUnlock()

	result := make([]map[string]string, len(lc.labels))
	copy(result, lc.labels)
	return result
}
