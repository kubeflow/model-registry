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
	// Note: source3 is disabled and should not appear in results
	sources := map[string]Source{
		"source1": {
			CatalogSource: model.CatalogSource{
				Id:      "source1",
				Name:    "Source 1",
				Enabled: apiutils.Of(true),
				Labels:  []string{"frontend", "production"},
			},
		},
		"source2": {
			CatalogSource: model.CatalogSource{
				Id:      "source2",
				Name:    "Source 2",
				Enabled: apiutils.Of(true),
				Labels:  []string{"Backend", "Development"}, // Mixed case to test case insensitivity
			},
		},
		"source3": {
			CatalogSource: model.CatalogSource{
				Id:      "source3",
				Name:    "Source 3",
				Enabled: apiutils.Of(false),                  // Disabled - should not appear in results
				Labels:  []string{"analytics", "PRODUCTION"}, // Mixed case
			},
		},
		"source4": {
			CatalogSource: model.CatalogSource{
				Id:      "source4",
				Name:    "Source 4",
				Enabled: apiutils.Of(true),
				Labels:  []string{"testing", "staging", "analytics"}, // Added analytics to test this label with enabled source
			},
		},
		"source5": {
			CatalogSource: model.CatalogSource{
				Id:      "source5",
				Name:    "Source 5",
				Enabled: apiutils.Of(true),
				Labels:  []string{}, // No labels
			},
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
			name:            "production label case insensitive - disabled source excluded",
			labels:          []string{"production"},
			expectedSources: []string{"source1"}, // source3 is disabled
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
			name:            "multiple different labels - disabled source excluded",
			labels:          []string{"analytics", "testing"},
			expectedSources: []string{"source4"}, // source3 is disabled, source4 has both labels
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
	sources := map[string]Source{
		"source1": {
			CatalogSource: model.CatalogSource{
				Id:      "source1",
				Name:    "Source 1",
				Enabled: apiutils.Of(true),
				Labels:  nil, // nil labels
			},
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

func TestSourceCollection_MergeOverride(t *testing.T) {
	tests := []struct {
		name          string
		originOrder   []string
		mergeSequence []struct {
			origin  string
			sources map[string]Source
		}
		expectedSources map[string]model.CatalogSource
	}{
		{
			name:        "later origin overrides earlier origin",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:             "hf",
							Name:           "Hugging Face",
							Enabled:        apiutils.Of(true),
							Labels:         []string{"default"},
							ExcludedModels: []string{"model-a"},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:             "hf",
							Name:           "Hugging Face Custom",
							Enabled:        apiutils.Of(true),
							Labels:         []string{"custom"},
							ExcludedModels: []string{"model-a", "DeepSeek"},
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:             "hf",
					Name:           "Hugging Face Custom",
					Enabled:        apiutils.Of(true),
					Labels:         []string{"custom"},
					ExcludedModels: []string{"model-a", "DeepSeek"},
				},
			},
		},
		{
			name:        "source from single origin behaves as before",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Hugging Face",
							Enabled: apiutils.Of(true),
							Labels:  []string{"default"},
						}},
					},
				},
				{
					origin:  "user.yaml",
					sources: map[string]Source{},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:      "hf",
					Name:    "Hugging Face",
					Enabled: apiutils.Of(true),
					Labels:  []string{"default"},
				},
			},
		},
		{
			name:        "multiple sources with partial override",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Hugging Face",
							Enabled: apiutils.Of(true),
							Labels:  []string{},
						}},
						"local": {CatalogSource: model.CatalogSource{
							Id:      "local",
							Name:    "Local Files",
							Enabled: apiutils.Of(true),
							Labels:  []string{},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:             "hf",
							Name:           "Hugging Face",
							Enabled:        apiutils.Of(true),
							Labels:         []string{},
							ExcludedModels: []string{"DeepSeek"},
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:             "hf",
					Name:           "Hugging Face",
					Enabled:        apiutils.Of(true),
					Labels:         []string{},
					ExcludedModels: []string{"DeepSeek"},
				},
				"local": {
					Id:      "local",
					Name:    "Local Files",
					Enabled: apiutils.Of(true),
					Labels:  []string{},
				},
			},
		},
		{
			name:        "three origins with cascading override",
			originOrder: []string{"base.yaml", "team.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "base.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Base HF",
							Enabled: apiutils.Of(true),
							Labels:  []string{"base"},
						}},
					},
				},
				{
					origin: "team.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Team HF",
							Enabled: apiutils.Of(true),
							Labels:  []string{"team"},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "User HF",
							Enabled: apiutils.Of(true),
							Labels:  []string{"user"},
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:      "hf",
					Name:    "User HF",
					Enabled: apiutils.Of(true),
					Labels:  []string{"user"},
				},
			},
		},
		{
			name:        "user can disable a source from default - disabled sources not returned",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Hugging Face",
							Enabled: apiutils.Of(true),
							Labels:  []string{},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Hugging Face",
							Enabled: apiutils.Of(false),
							Labels:  []string{},
						}},
					},
				},
			},
			// Disabled sources are not returned by All()
			expectedSources: map[string]model.CatalogSource{},
		},
		{
			name:        "sparse override: user enables disabled source with just id and enabled",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"models_x": {CatalogSource: model.CatalogSource{
							Id:             "models_x",
							Name:           "Models X Catalog",
							Enabled:        apiutils.Of(false), // Disabled in default
							Labels:         []string{"enterprise"},
							IncludedModels: []string{"model-a", "model-b"},
							ExcludedModels: []string{"model-c"},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						// Sparse override: only id and enabled
						"models_x": {CatalogSource: model.CatalogSource{
							Id:      "models_x",
							Enabled: apiutils.Of(true), // Enable it
							// Name, Labels, IncludedModels, ExcludedModels are nil/empty
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"models_x": {
					Id:             "models_x",
					Name:           "Models X Catalog",             // Inherited from default
					Enabled:        apiutils.Of(true),              // Overridden by user
					Labels:         []string{"enterprise"},         // Inherited from default
					IncludedModels: []string{"model-a", "model-b"}, // Inherited from default
					ExcludedModels: []string{"model-c"},            // Inherited from default
				},
			},
		},
		{
			name:        "sparse override: user changes only excluded models",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:             "hf",
							Name:           "Hugging Face",
							Enabled:        apiutils.Of(true),
							Labels:         []string{"public"},
							ExcludedModels: []string{"model-a"},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id: "hf",
							// Only override ExcludedModels
							ExcludedModels: []string{"model-a", "DeepSeek", "banned-model"},
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:             "hf",
					Name:           "Hugging Face",                                  // Inherited
					Enabled:        apiutils.Of(true),                               // Inherited
					Labels:         []string{"public"},                              // Inherited
					ExcludedModels: []string{"model-a", "DeepSeek", "banned-model"}, // Overridden
				},
			},
		},
		{
			name:        "sparse override: user clears labels with empty slice",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:      "hf",
							Name:    "Hugging Face",
							Enabled: apiutils.Of(true),
							Labels:  []string{"public", "ai"},
						}},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:     "hf",
							Labels: []string{}, // Explicitly clear labels
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:      "hf",
					Name:    "Hugging Face",    // Inherited
					Enabled: apiutils.Of(true), // Inherited
					Labels:  []string{},        // Overridden to empty
				},
			},
		},
		{
			name:        "defaults applied: enabled defaults to true, labels defaults to empty",
			originOrder: []string{"default.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"hf": {CatalogSource: model.CatalogSource{
							Id:   "hf",
							Name: "Hugging Face",
							// Enabled and Labels are nil
						}},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"hf": {
					Id:      "hf",
					Name:    "Hugging Face",
					Enabled: apiutils.Of(true), // Default applied
					Labels:  []string{},        // Default applied
				},
			},
		},
		{
			name:        "sparse override: type and properties are inherited",
			originOrder: []string{"default.yaml", "user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "default.yaml",
					sources: map[string]Source{
						"models_x": {
							CatalogSource: model.CatalogSource{
								Id:      "models_x",
								Name:    "Models X Catalog",
								Enabled: apiutils.Of(false), // Disabled in default
								Labels:  []string{"enterprise"},
							},
							Type: "yaml",
							Properties: map[string]any{
								"yamlCatalogPath": "models-x.yaml",
							},
						},
					},
				},
				{
					origin: "user.yaml",
					sources: map[string]Source{
						// Sparse override: only id and enabled
						"models_x": {
							CatalogSource: model.CatalogSource{
								Id:      "models_x",
								Enabled: apiutils.Of(true), // Enable it
							},
							// Type and Properties are empty/nil - should be inherited
						},
					},
				},
			},
			expectedSources: map[string]model.CatalogSource{
				"models_x": {
					Id:      "models_x",
					Name:    "Models X Catalog",     // Inherited from default
					Enabled: apiutils.Of(true),      // Overridden by user
					Labels:  []string{"enterprise"}, // Inherited from default
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewSourceCollection(tt.originOrder...)

			for _, merge := range tt.mergeSequence {
				err := sc.Merge(merge.origin, merge.sources)
				if err != nil {
					t.Fatalf("Merge(%s) failed: %v", merge.origin, err)
				}
			}

			result := sc.All()

			if len(result) != len(tt.expectedSources) {
				t.Errorf("All() returned %d sources, want %d", len(result), len(tt.expectedSources))
			}

			for id, expected := range tt.expectedSources {
				got, ok := result[id]
				if !ok {
					t.Errorf("source %s not found in result", id)
					continue
				}
				if got.Id != expected.Id {
					t.Errorf("source %s: Id = %s, want %s", id, got.Id, expected.Id)
				}
				if got.Name != expected.Name {
					t.Errorf("source %s: Name = %s, want %s", id, got.Name, expected.Name)
				}
				if *got.Enabled != *expected.Enabled {
					t.Errorf("source %s: Enabled = %v, want %v", id, *got.Enabled, *expected.Enabled)
				}
				if !reflect.DeepEqual(got.Labels, expected.Labels) {
					t.Errorf("source %s: Labels = %v, want %v", id, got.Labels, expected.Labels)
				}
				if !reflect.DeepEqual(got.ExcludedModels, expected.ExcludedModels) {
					t.Errorf("source %s: ExcludedModels = %v, want %v", id, got.ExcludedModels, expected.ExcludedModels)
				}
			}
		})
	}
}

func TestSourceCollection_MergeOverride_Get(t *testing.T) {
	sc := NewSourceCollection("default.yaml", "user.yaml")

	// Merge default config
	err := sc.Merge("default.yaml", map[string]Source{
		"hf": {CatalogSource: model.CatalogSource{
			Id:      "hf",
			Name:    "Hugging Face Default",
			Enabled: apiutils.Of(true),
			Labels:  []string{},
		}},
	})
	if err != nil {
		t.Fatalf("Merge(default.yaml) failed: %v", err)
	}

	// Merge user config that overrides
	err = sc.Merge("user.yaml", map[string]Source{
		"hf": {CatalogSource: model.CatalogSource{
			Id:             "hf",
			Name:           "Hugging Face User",
			Enabled:        apiutils.Of(true),
			Labels:         []string{"user-managed"},
			ExcludedModels: []string{"DeepSeek"},
		}},
	})
	if err != nil {
		t.Fatalf("Merge(user.yaml) failed: %v", err)
	}

	// Get should return the overridden source
	source, ok := sc.Get("hf")
	if !ok {
		t.Fatal("Get(hf) returned false, want true")
	}
	if source.Name != "Hugging Face User" {
		t.Errorf("Get(hf).Name = %s, want 'Hugging Face User'", source.Name)
	}
	if len(source.ExcludedModels) != 1 || source.ExcludedModels[0] != "DeepSeek" {
		t.Errorf("Get(hf).ExcludedModels = %v, want [DeepSeek]", source.ExcludedModels)
	}
}

func TestSourceCollection_MergeOverride_DynamicOrigin(t *testing.T) {
	// Test that origins not in the initial order are appended
	sc := NewSourceCollection("default.yaml")

	err := sc.Merge("default.yaml", map[string]Source{
		"hf": {CatalogSource: model.CatalogSource{Id: "hf", Name: "Default", Enabled: apiutils.Of(true), Labels: []string{}}},
	})
	if err != nil {
		t.Fatalf("Merge(default.yaml) failed: %v", err)
	}

	// Dynamic origin not in initial order
	err = sc.Merge("extra.yaml", map[string]Source{
		"hf": {CatalogSource: model.CatalogSource{Id: "hf", Name: "Extra", Enabled: apiutils.Of(true), Labels: []string{}}},
	})
	if err != nil {
		t.Fatalf("Merge(extra.yaml) failed: %v", err)
	}

	source, _ := sc.Get("hf")
	if source.Name != "Extra" {
		t.Errorf("dynamically added origin should override earlier origins, got Name = %s", source.Name)
	}
}

func TestSourceCollection_MergeOverride_TypeAndProperties(t *testing.T) {
	// Test that Type and Properties are properly inherited with sparse overrides
	sc := NewSourceCollection("default.yaml", "user.yaml")

	// Merge default config with full source definition
	err := sc.Merge("default.yaml", map[string]Source{
		"models_x": {
			CatalogSource: model.CatalogSource{
				Id:      "models_x",
				Name:    "Models X Catalog",
				Enabled: apiutils.Of(false), // Disabled
				Labels:  []string{"enterprise"},
			},
			Type: "yaml",
			Properties: map[string]any{
				"yamlCatalogPath": "models-x.yaml",
				"otherProp":       123,
			},
		},
	})
	if err != nil {
		t.Fatalf("Merge(default.yaml) failed: %v", err)
	}

	// Merge sparse user config that only enables the source
	err = sc.Merge("user.yaml", map[string]Source{
		"models_x": {
			CatalogSource: model.CatalogSource{
				Id:      "models_x",
				Enabled: apiutils.Of(true), // Enable it
			},
			// Type and Properties are empty - should be inherited
		},
	})
	if err != nil {
		t.Fatalf("Merge(user.yaml) failed: %v", err)
	}

	// AllSources should return the merged source with Type and Properties
	sources := sc.AllSources()
	source, ok := sources["models_x"]
	if !ok {
		t.Fatal("AllSources() should return models_x")
	}

	if source.Type != "yaml" {
		t.Errorf("Type = %s, want 'yaml' (inherited from default)", source.Type)
	}

	if source.Properties == nil {
		t.Fatal("Properties should be inherited from default")
	}

	if source.Properties["yamlCatalogPath"] != "models-x.yaml" {
		t.Errorf("Properties[yamlCatalogPath] = %v, want 'models-x.yaml'", source.Properties["yamlCatalogPath"])
	}

	if source.Properties["otherProp"] != 123 {
		t.Errorf("Properties[otherProp] = %v, want 123", source.Properties["otherProp"])
	}

	// CatalogSource fields should also be merged
	if source.Name != "Models X Catalog" {
		t.Errorf("Name = %s, want 'Models X Catalog'", source.Name)
	}

	if *source.Enabled != true {
		t.Errorf("Enabled = %v, want true", *source.Enabled)
	}
}

func TestSourceCollection_MergeOverride_Origin(t *testing.T) {
	// Test that Origin is correctly tracked through merge operations
	// This is important for resolving relative paths in source properties

	tests := []struct {
		name           string
		originOrder    []string
		mergeSequence  []struct {
			origin  string
			sources map[string]Source
		}
		expectedOrigins map[string]string // sourceId -> expected origin
	}{
		{
			name:        "origin is preserved from first definition",
			originOrder: []string{"/config/default.yaml", "/user-config/user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "/config/default.yaml",
					sources: map[string]Source{
						"hf": {
							CatalogSource: model.CatalogSource{
								Id:      "hf",
								Name:    "Hugging Face",
								Enabled: apiutils.Of(true),
							},
							Type: "yaml",
							Properties: map[string]any{
								"yamlCatalogPath": "models.yaml",
							},
							Origin: "/config/default.yaml",
						},
					},
				},
				{
					origin: "/user-config/user.yaml",
					sources: map[string]Source{
						// Sparse override: only enable and change name
						"hf": {
							CatalogSource: model.CatalogSource{
								Id:      "hf",
								Name:    "Hugging Face Custom",
								Enabled: apiutils.Of(true),
							},
							Origin: "/user-config/user.yaml",
							// Properties is nil, so Origin should stay with base
						},
					},
				},
			},
			// Origin should be from default.yaml since Properties weren't overridden
			expectedOrigins: map[string]string{
				"hf": "/config/default.yaml",
			},
		},
		{
			name:        "origin changes when properties are overridden",
			originOrder: []string{"/config/default.yaml", "/user-config/user.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "/config/default.yaml",
					sources: map[string]Source{
						"local": {
							CatalogSource: model.CatalogSource{
								Id:      "local",
								Name:    "Local Catalog",
								Enabled: apiutils.Of(true),
							},
							Type: "yaml",
							Properties: map[string]any{
								"yamlCatalogPath": "default-models.yaml",
							},
							Origin: "/config/default.yaml",
						},
					},
				},
				{
					origin: "/user-config/user.yaml",
					sources: map[string]Source{
						"local": {
							CatalogSource: model.CatalogSource{
								Id:      "local",
								Enabled: apiutils.Of(true),
							},
							// Override Properties - this changes where relative paths resolve from
							Properties: map[string]any{
								"yamlCatalogPath": "user-models.yaml",
							},
							Origin: "/user-config/user.yaml",
						},
					},
				},
			},
			// Origin should be from user.yaml since Properties were overridden
			expectedOrigins: map[string]string{
				"local": "/user-config/user.yaml",
			},
		},
		{
			name:        "multiple sources from different origins",
			originOrder: []string{"/admin/sources.yaml", "/user/sources.yaml"},
			mergeSequence: []struct {
				origin  string
				sources map[string]Source
			}{
				{
					origin: "/admin/sources.yaml",
					sources: map[string]Source{
						"admin-catalog": {
							CatalogSource: model.CatalogSource{
								Id:      "admin-catalog",
								Name:    "Admin Catalog",
								Enabled: apiutils.Of(true),
							},
							Type: "yaml",
							Properties: map[string]any{
								"yamlCatalogPath": "admin-models.yaml",
							},
							Origin: "/admin/sources.yaml",
						},
					},
				},
				{
					origin: "/user/sources.yaml",
					sources: map[string]Source{
						"user-catalog": {
							CatalogSource: model.CatalogSource{
								Id:      "user-catalog",
								Name:    "User Catalog",
								Enabled: apiutils.Of(true),
							},
							Type: "yaml",
							Properties: map[string]any{
								"yamlCatalogPath": "user-models.yaml",
							},
							Origin: "/user/sources.yaml",
						},
					},
				},
			},
			// Each source should keep its own origin
			expectedOrigins: map[string]string{
				"admin-catalog": "/admin/sources.yaml",
				"user-catalog":  "/user/sources.yaml",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sc := NewSourceCollection(tt.originOrder...)

			for _, merge := range tt.mergeSequence {
				err := sc.Merge(merge.origin, merge.sources)
				if err != nil {
					t.Fatalf("Merge(%s) failed: %v", merge.origin, err)
				}
			}

			sources := sc.AllSources()

			for id, expectedOrigin := range tt.expectedOrigins {
				source, ok := sources[id]
				if !ok {
					t.Errorf("source %s not found in AllSources()", id)
					continue
				}
				if source.Origin != expectedOrigin {
					t.Errorf("source %s: Origin = %s, want %s", id, source.Origin, expectedOrigin)
				}
			}
		})
	}
}

func TestSourceCollection_ByLabel_NullBehavior(t *testing.T) {
	// Create test sources with various label configurations to test "null" behavior
	// Note: disabled sources are filtered out and not returned
	sources := map[string]Source{
		"source_with_labels": {
			CatalogSource: model.CatalogSource{
				Id:      "source_with_labels",
				Name:    "Source With Labels",
				Enabled: apiutils.Of(true),
				Labels:  []string{"frontend", "production"},
			},
		},
		"source_empty_labels": {
			CatalogSource: model.CatalogSource{
				Id:      "source_empty_labels",
				Name:    "Source Empty Labels",
				Enabled: apiutils.Of(true),
				Labels:  []string{}, // Empty labels slice
			},
		},
		"source_nil_labels": {
			CatalogSource: model.CatalogSource{
				Id:      "source_nil_labels",
				Name:    "Source Nil Labels",
				Enabled: apiutils.Of(true),
				Labels:  nil, // Nil labels
			},
		},
		"source_another_with_labels": {
			CatalogSource: model.CatalogSource{
				Id:      "source_another_with_labels",
				Name:    "Another Source With Labels",
				Enabled: apiutils.Of(true), // Changed to enabled
				Labels:  []string{"backend", "testing"},
			},
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
