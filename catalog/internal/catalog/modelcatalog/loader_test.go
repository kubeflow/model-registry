package modelcatalog

import (
	"context"
	"testing"
	"time"

	mapset "github.com/deckarep/golang-set/v2"
	"github.com/kubeflow/model-registry/catalog/internal/catalog/basecatalog"
	dbmodels "github.com/kubeflow/model-registry/catalog/internal/db/models"
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
				nil, // MCPServerRepository
				nil, // MCPServerToolRepository
			)

			// Create loader and populate sources
			loader := NewModelLoader(services, basecatalog.NewBaseLoader([]string{}))

			// Add all test sources to the loader's source collection in a single Merge call
			sourcesMap := make(map[string]basecatalog.ModelSource)
			for sourceID, enabled := range tt.enabledSources {
				source := apimodels.CatalogSource{
					Id:      sourceID,
					Name:    "Test " + sourceID,
					Enabled: enabled,
				}

				sourcesMap[sourceID] = basecatalog.ModelSource{
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

			// Call the method under test (empty combined IDs: model-only test)
			err := loader.removeModelsFromMissingSources(mapset.NewSet[string]())

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
	var config basecatalog.SourceConfig
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

func TestLoader_StartWithLeaderElection(t *testing.T) {
	// Create mock repositories with tracking capabilities
	mockModelRepo := &MockCatalogModelRepositoryWithSourceTracking{
		ExistingSourceIDs: []string{},
		DeletedSources:    []string{},
	}
	mockArtifactRepo := &MockCatalogArtifactRepository{}
	mockModelArtifactRepo := &MockCatalogModelArtifactRepository{}
	mockMetricsArtifactRepo := &MockCatalogMetricsArtifactRepository{}
	mockSourceRepo := &MockCatalogSourceRepository{}

	services := service.NewServices(
		mockModelRepo,
		mockArtifactRepo,
		mockModelArtifactRepo,
		mockMetricsArtifactRepo,
		mockSourceRepo,
		&MockPropertyOptionsRepository{},
		nil, // MCPServerRepository
		nil, // MCPServerToolRepository
	)

	// Register a test provider
	testProviderName := "test-leader-provider"
	RegisterModelProvider(testProviderName, func(ctx context.Context, source *basecatalog.ModelSource, reldir string) (<-chan ModelProviderRecord, error) {
		ch := make(chan ModelProviderRecord, 2)

		modelName := "test-model-1"
		model := &dbmodels.CatalogModelImpl{
			Attributes: &dbmodels.CatalogModelAttributes{
				Name: &modelName,
			},
		}

		ch <- ModelProviderRecord{
			Model:     model,
			Artifacts: []dbmodels.CatalogArtifact{},
		}

		// Send completion marker
		ch <- ModelProviderRecord{
			Model: nil,
		}

		close(ch)
		return ch, nil
	})

	testConfig := &basecatalog.SourceConfig{
		ModelCatalogs: []basecatalog.ModelSource{
			{
				CatalogSource: apimodels.CatalogSource{
					Id:      "test-catalog",
					Name:    "Test Catalog",
					Enabled: apiutils.Of(true),
				},
				Type: testProviderName,
			},
		},
	}

	t.Run("standby mode skips database writes", func(t *testing.T) {
		baseLoader := basecatalog.NewBaseLoader([]string{})
		loader := NewModelLoader(services, baseLoader)

		// Populate sources
		err := loader.updateSources("test-path", testConfig)
		assert.NoError(t, err)

		// In standby mode (read-only), just parse configs without becoming leader
		err = loader.ParseAllConfigs()
		assert.NoError(t, err)

		// Verify we're not in leader mode
		assert.False(t, baseLoader.ShouldWriteDatabase(), "Standby mode should not write to database")

		// Wait a bit for any goroutines to process
		time.Sleep(100 * time.Millisecond)

		// Verify no database writes occurred in standby mode
		assert.Empty(t, mockModelRepo.SavedModels, "Standby mode should not write models to database")
		assert.Empty(t, mockModelArtifactRepo.SavedArtifacts, "Standby mode should not write artifacts to database")
		assert.Empty(t, mockMetricsArtifactRepo.SavedMetrics, "Standby mode should not write metrics to database")
	})

	t.Run("leader mode performs database writes", func(t *testing.T) {
		// Reset mock repositories
		mockModelRepo.SavedModels = []dbmodels.CatalogModel{}
		mockModelArtifactRepo.SavedArtifacts = []dbmodels.CatalogModelArtifact{}
		mockMetricsArtifactRepo.SavedMetrics = []dbmodels.CatalogMetricsArtifact{}

		baseLoader := basecatalog.NewBaseLoader([]string{})
		loader := NewModelLoader(services, baseLoader)
		ctx := context.Background()

		// Populate sources
		err := loader.updateSources("test-path", testConfig)
		assert.NoError(t, err)

		// Parse configs first
		err = loader.ParseAllConfigs()
		assert.NoError(t, err)

		// Set leader mode on the base loader
		baseLoader.SetLeader(true)

		// Verify we're in leader mode
		assert.True(t, baseLoader.ShouldWriteDatabase(), "Leader mode should write to database")

		// Perform leader operations (empty combined IDs: model-only test)
		err = loader.PerformLeaderOperations(ctx, mapset.NewSet[string]())
		assert.NoError(t, err)

		// Wait for goroutines to process
		time.Sleep(200 * time.Millisecond)

		// Verify database writes occurred in leader mode
		assert.NotEmpty(t, mockModelRepo.SavedModels, "Leader mode should write models to database")
	})

	// Note: Lifecycle tests (StartReadOnly/StartLeader) are now tested
	// at the integration level with the unified catalog.Loader
}
