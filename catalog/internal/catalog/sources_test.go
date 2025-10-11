package catalog

import (
	"reflect"
	"sort"
	"testing"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
)

func TestSourceCollection_ByLabel(t *testing.T) {
	// Create test sources with various labels
	sources := map[string]model.CatalogSource{
		"source1": {
			Id:      "source1",
			Name:    "Source 1",
			Enabled: apiutils.Of(true),
			Labels:  []string{"frontend", "production"},
		},
		"source2": {
			Id:      "source2",
			Name:    "Source 2",
			Enabled: apiutils.Of(true),
			Labels:  []string{"Backend", "Development"}, // Mixed case to test case insensitivity
		},
		"source3": {
			Id:      "source3",
			Name:    "Source 3",
			Enabled: apiutils.Of(false),
			Labels:  []string{"analytics", "PRODUCTION"}, // Mixed case
		},
		"source4": {
			Id:      "source4",
			Name:    "Source 4",
			Enabled: apiutils.Of(true),
			Labels:  []string{"testing", "staging"},
		},
		"source5": {
			Id:      "source5",
			Name:    "Source 5",
			Enabled: apiutils.Of(true),
			Labels:  []string{}, // No labels
		},
	}

	tests := []struct {
		name            string
		labels          []string
		expectedSources []string // IDs of expected sources
	}{
		{
			name:            "single label match",
			labels:          []string{"frontend"},
			expectedSources: []string{"source1"},
		},
		{
			name:            "case insensitive match",
			labels:          []string{"FRONTEND"},
			expectedSources: []string{"source1"},
		},
		{
			name:            "multiple labels - any match",
			labels:          []string{"frontend", "backend"},
			expectedSources: []string{"source1", "source2"},
		},
		{
			name:            "case insensitive multiple labels",
			labels:          []string{"FRONTEND", "backend"},
			expectedSources: []string{"source1", "source2"},
		},
		{
			name:            "production label case insensitive",
			labels:          []string{"production"},
			expectedSources: []string{"source1", "source3"},
		},
		{
			name:            "no matching labels",
			labels:          []string{"nonexistent"},
			expectedSources: nil,
		},
		{
			name:            "empty labels input",
			labels:          []string{},
			expectedSources: nil,
		},
		{
			name:            "multiple different labels",
			labels:          []string{"analytics", "testing"},
			expectedSources: []string{"source3", "source4"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new SourceCollection and populate it
			sc := NewSourceCollection()
			err := sc.Merge("test-origin", sources)
			if err != nil {
				t.Fatalf("Failed to merge sources: %v", err)
			}

			// Call ByLabel
			result := sc.ByLabel(tt.labels)

			if tt.expectedSources == nil && result != nil {
				t.Errorf("ByLabel() = %v, want %v", result, tt.expectedSources)
				return
			}

			// Extract IDs from result for comparison
			var resultIDs []string
			for _, source := range result {
				resultIDs = append(resultIDs, source.Id)
			}

			// Sort both slices for comparison
			sort.Strings(resultIDs)
			sort.Strings(tt.expectedSources)

			if !reflect.DeepEqual(resultIDs, tt.expectedSources) {
				t.Errorf("ByLabel() = %v, want %v", resultIDs, tt.expectedSources)
			}

			// Verify that the returned sources are complete objects, not just IDs
			for _, source := range result {
				if source.Name == "" {
					t.Errorf("Returned source %s has empty name", source.Id)
				}
				if source.Labels == nil {
					t.Errorf("Returned source %s has nil labels", source.Id)
				}
			}
		})
	}
}

func TestSourceCollection_ByLabel_EmptyCollection(t *testing.T) {
	sc := NewSourceCollection()

	tests := []struct {
		name   string
		labels []string
	}{
		{
			name:   "empty collection with regular labels",
			labels: []string{"frontend"},
		},
		{
			name:   "empty collection with null label",
			labels: []string{"null"},
		},
		{
			name:   "empty collection with empty labels",
			labels: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sc.ByLabel(tt.labels)
			if len(result) != 0 {
				t.Errorf("ByLabel() on empty collection should return empty slice, got %d items", len(result))
			}
		})
	}
}

func TestSourceCollection_ByLabel_NilLabels(t *testing.T) {
	sc := NewSourceCollection()

	// Add a source with nil labels (edge case)
	sources := map[string]model.CatalogSource{
		"source1": {
			Id:      "source1",
			Name:    "Source 1",
			Enabled: apiutils.Of(true),
			Labels:  nil, // nil labels
		},
	}

	err := sc.Merge("test-origin", sources)
	if err != nil {
		t.Fatalf("Failed to merge sources: %v", err)
	}

	tests := []struct {
		name          string
		labels        []string
		expectedCount int
	}{
		{
			name:          "search with regular label on source with nil labels",
			labels:        []string{"frontend"},
			expectedCount: 0,
		},
		{
			name:          "search with null label on source with nil labels",
			labels:        []string{"null"},
			expectedCount: 1, // Should return all sources including those with nil labels
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sc.ByLabel(tt.labels)
			if len(result) != tt.expectedCount {
				t.Errorf("ByLabel() = %d sources, want %d", len(result), tt.expectedCount)
			}
		})
	}
}

func TestSourceCollection_ByLabel_NullBehavior(t *testing.T) {
	// Create test sources with various label configurations to test "null" behavior
	sources := map[string]model.CatalogSource{
		"source_with_labels": {
			Id:      "source_with_labels",
			Name:    "Source With Labels",
			Enabled: apiutils.Of(true),
			Labels:  []string{"frontend", "production"},
		},
		"source_empty_labels": {
			Id:      "source_empty_labels",
			Name:    "Source Empty Labels",
			Enabled: apiutils.Of(true),
			Labels:  []string{}, // Empty labels slice
		},
		"source_nil_labels": {
			Id:      "source_nil_labels",
			Name:    "Source Nil Labels",
			Enabled: apiutils.Of(true),
			Labels:  nil, // Nil labels
		},
		"source_another_with_labels": {
			Id:      "source_another_with_labels",
			Name:    "Another Source With Labels",
			Enabled: apiutils.Of(false),
			Labels:  []string{"backend", "testing"},
		},
	}

	tests := []struct {
		name            string
		labels          []string
		expectedSources []string // IDs of expected sources
		description     string
	}{
		{
			name:            "null label returns sources without labels",
			labels:          []string{"null"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels"},
			description:     "Should return sources with empty or nil labels when searching for 'null'",
		},
		{
			name:            "null label case insensitive",
			labels:          []string{"NULL"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels"},
			description:     "Should be case insensitive for 'null' label",
		},
		{
			name:            "null label mixed case",
			labels:          []string{"Null"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels"},
			description:     "Should handle mixed case 'null' label",
		},
		{
			name:            "null with other labels",
			labels:          []string{"null", "frontend"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels", "source_with_labels"},
			description:     "Should return sources without labels AND sources with matching labels",
		},
		{
			name:            "null with multiple other labels",
			labels:          []string{"null", "frontend", "backend"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels", "source_with_labels", "source_another_with_labels"},
			description:     "Should return sources without labels AND sources with any matching labels",
		},
		{
			name:            "multiple nulls should work same as single null",
			labels:          []string{"null", "NULL"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels"},
			description:     "Multiple 'null' variants should behave same as single 'null'",
		},
		{
			name:            "null only matches unlabeled sources",
			labels:          []string{"null", "nonexistent"},
			expectedSources: []string{"source_empty_labels", "source_nil_labels"},
			description:     "Should return sources without labels even when other labels don't match",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a new SourceCollection and populate it
			sc := NewSourceCollection()
			err := sc.Merge("test-origin", sources)
			if err != nil {
				t.Fatalf("Failed to merge sources: %v", err)
			}

			// Call ByLabel
			result := sc.ByLabel(tt.labels)

			// Extract IDs from result for comparison
			var resultIDs []string
			for _, source := range result {
				resultIDs = append(resultIDs, source.Id)
			}

			// Sort both slices for comparison
			sort.Strings(resultIDs)
			sort.Strings(tt.expectedSources)

			if !reflect.DeepEqual(resultIDs, tt.expectedSources) {
				t.Errorf("ByLabel(%v) = %v, want %v\nDescription: %s", tt.labels, resultIDs, tt.expectedSources, tt.description)
			}

			// Verify that each returned source is a complete object
			for _, source := range result {
				if source.Id == "" {
					t.Errorf("Returned source has empty ID")
				}
				if source.Name == "" {
					t.Errorf("Returned source %s has empty name", source.Id)
				}
			}
		})
	}
}
