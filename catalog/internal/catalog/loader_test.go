package catalog

import (
	"testing"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	apimodels "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/yaml"
)

func TestRemoveModelsFromMissingSources(t *testing.T) {
	// This test verifies the current behavior of removeModelsFromMissingSources.
	// The method removes models from sources that are either:
	// 1. Not present in the current configuration
	// 2. Explicitly disabled (Enabled = false)
	// Note: Sources with nil Enabled are filtered out by AllSources() - see bug comment below

	tests := []struct {
		name                   string
		enabledSources         map[string]*bool // source ID -> enabled status (nil means default true)
		existingSourceIDs      []string         // source IDs currently in database
		expectedDeletedSources []string         // source IDs that should be deleted
		repositoryError        string           // if set, repository returns this error
		expectError            bool
	}{
		{
			name: "removes models from sources not in config",
			enabledSources: map[string]*bool{
				"source1": apiutils.Of(true),
				"source2": apiutils.Of(true),
			},
			existingSourceIDs:      []string{"source1", "source2", "source3", "source4"},
			expectedDeletedSources: []string{"source3", "source4"}, // source3 and source4 not in config
		},
		{
			name: "no deletion when all database sources are in config",
			enabledSources: map[string]*bool{
				"source1": apiutils.Of(true),
				"source2": apiutils.Of(true),
				"source3": apiutils.Of(true),
			},
			existingSourceIDs:      []string{"source1", "source2"},
			expectedDeletedSources: []string{}, // no deletions - all database sources are in config
		},
		{
			name: "handles empty existing sources",
			enabledSources: map[string]*bool{
				"source1": apiutils.Of(true),
			},
			existingSourceIDs:      []string{},
			expectedDeletedSources: []string{}, // no deletions needed - no existing sources
		},
		{
			name:                   "handles empty config sources",
			enabledSources:         map[string]*bool{}, // no sources in config
			existingSourceIDs:      []string{"source1", "source2"},
			expectedDeletedSources: []string{"source1", "source2"}, // all existing sources deleted
		},
		{
			name: "correctly handles default enabled sources",
			enabledSources: map[string]*bool{
				"source1": nil,               // default enabled - converted to true by applyDefaults
				"source2": apiutils.Of(true), // explicitly enabled
			},
			existingSourceIDs:      []string{"source1", "source2", "source3"},
			expectedDeletedSources: []string{"source3"}, // only source3 (not in config) gets deleted
		},
		{
			name: "handles repository error on GetDistinctSourceIDs",
			enabledSources: map[string]*bool{
				"source1": apiutils.Of(true),
			},
			existingSourceIDs: []string{"source1"},
			repositoryError:   "get_distinct_source_ids_error",
			expectError:       true,
		},
		{
			name: "handles repository error on DeleteBySource",
			enabledSources: map[string]*bool{
				"source1": apiutils.Of(true),
			},
			existingSourceIDs:      []string{"source1", "source2"},
			repositoryError:        "delete_by_source_error",
			expectedDeletedSources: []string{"source2"}, // source2 should be attempted for deletion
			expectError:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock repository with test data
			mockModelRepo := &MockCatalogModelRepositoryWithSourceTracking{
				ExistingSourceIDs: tt.existingSourceIDs,
				DeletedSources:    []string{},
				ErrorType:         tt.repositoryError,
			}

			services := service.NewServices(
				mockModelRepo,
				&MockCatalogArtifactRepository{},
				&MockCatalogModelArtifactRepository{},
				&MockCatalogMetricsArtifactRepository{},
				&MockCatalogSourceRepository{},
				&MockPropertyOptionsRepository{},
			)

			// Create loader and populate sources
			loader := NewLoader(services, []string{})

			// Add all test sources to the loader's source collection in a single Merge call
			sourcesMap := make(map[string]Source)
			for sourceID, enabled := range tt.enabledSources {
				source := apimodels.CatalogSource{
					Id:      sourceID,
					Name:    "Test " + sourceID,
					Enabled: enabled,
				}

				sourcesMap[sourceID] = Source{
					CatalogSource: source,
					Type:          "test",
				}
			}

			if len(sourcesMap) > 0 {
				err := loader.Sources.Merge("test-path", sourcesMap)
				if err != nil {
					t.Fatalf("Failed to add sources: %v", err)
				}
			}

			// Call the method under test
			err := loader.removeModelsFromMissingSources()

			// Verify error expectation
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify which sources were deleted
			expectedSet := mapset.NewSet(tt.expectedDeletedSources...)
			actualSet := mapset.NewSet(mockModelRepo.DeletedSources...)

			if !expectedSet.Equal(actualSet) {
				t.Errorf("Expected deleted sources %v, got %v",
					tt.expectedDeletedSources, mockModelRepo.DeletedSources)
			}
		})
	}
}

// MockCatalogModelRepositoryWithSourceTracking extends the existing mock to add source tracking
type MockCatalogModelRepositoryWithSourceTracking struct {
	MockCatalogModelRepository
	ExistingSourceIDs []string
	DeletedSources    []string
	ErrorType         string // "get_distinct_source_ids_error" or "delete_by_source_error"
}

func (m *MockCatalogModelRepositoryWithSourceTracking) GetDistinctSourceIDs() ([]string, error) {
	if m.ErrorType == "get_distinct_source_ids_error" {
		return nil, NewMockError("failed to get distinct source IDs")
	}
	return m.ExistingSourceIDs, nil
}

func (m *MockCatalogModelRepositoryWithSourceTracking) DeleteBySource(sourceID string) error {
	if m.ErrorType == "delete_by_source_error" {
		return NewMockError("failed to delete models from source: " + sourceID)
	}
	m.DeletedSources = append(m.DeletedSources, sourceID)
	return nil
}

// Helper function to create mock errors
func NewMockError(message string) error {
	return &MockRepositoryError{Message: message}
}

type MockRepositoryError struct {
	Message string
}

func (e *MockRepositoryError) Error() string {
	return e.Message

}

func TestSourceConfigNamedQueries(t *testing.T) {
	yamlContent := `
catalogs: []
namedQueries:
  validation-default:
    ttft_p90:
      operator: '<'
      value: 70
    workload_type:
      operator: '='
      value: "Chat"
  high-performance:
    performance_score:
      operator: '>'
      value: 0.95
`
	var config sourceConfig
	err := yaml.UnmarshalStrict([]byte(yamlContent), &config)
	assert.NoError(t, err)
	assert.NotNil(t, config.NamedQueries)
	assert.Len(t, config.NamedQueries, 2)

	validationQuery := config.NamedQueries["validation-default"]
	assert.NotNil(t, validationQuery)
	assert.Equal(t, "<", validationQuery["ttft_p90"].Operator)
	assert.Equal(t, float64(70), validationQuery["ttft_p90"].Value)
	assert.Equal(t, "=", validationQuery["workload_type"].Operator)
	assert.Equal(t, "Chat", validationQuery["workload_type"].Value)
}
