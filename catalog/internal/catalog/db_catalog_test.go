package catalog

import (
	"context"
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/internal/apiutils"
	mr_models "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainHelper(m))
}

func TestDBCatalog(t *testing.T) {
	// Setup test database
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get type IDs
	catalogModelTypeID := getCatalogModelTypeIDForDBTest(t, sharedDB)
	modelArtifactTypeID := getCatalogModelArtifactTypeIDForDBTest(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeIDForDBTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int64{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)

	// Create DB catalog instance
	dbCatalog := NewDBCatalog(catalogModelRepo, catalogArtifactRepo)
	ctx := context.Background()

	t.Run("TestNewDBCatalog", func(t *testing.T) {
		catalog := NewDBCatalog(catalogModelRepo, catalogArtifactRepo)
		require.NotNil(t, catalog)

		// Verify it implements the interface
		var _ CatalogSourceProvider = catalog
	})

	t.Run("TestGetModel_Success", func(t *testing.T) {
		// Create test model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-get-model"),
				ExternalID: apiutils.Of("test-get-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("test-source-id")},
				{Name: "description", StringValue: apiutils.Of("Test model description")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Test GetModel
		retrievedModel, err := dbCatalog.GetModel(ctx, "test-get-model", "test-source-id")
		require.NoError(t, err)
		require.NotNil(t, retrievedModel)

		assert.Equal(t, "test-get-model", retrievedModel.Name)
		assert.Equal(t, strconv.FormatInt(int64(*savedModel.GetID()), 10), *retrievedModel.Id)
		assert.Equal(t, "test-get-model-ext", *retrievedModel.ExternalId)
		assert.Equal(t, "test-source-id", *retrievedModel.SourceId)
		assert.Equal(t, "Test model description", *retrievedModel.Description)
	})

	t.Run("TestGetModel_NotFound", func(t *testing.T) {
		// Test with non-existent model
		_, err := dbCatalog.GetModel(ctx, "non-existent-model", "test-source-id")
		require.Error(t, err)
		assert.Contains(t, err.Error(), "no models found")
		assert.ErrorIs(t, err, api.ErrNotFound)
	})

	t.Run("TestGetModel_DatabaseConstraints", func(t *testing.T) {
		// Test database constraint behavior - attempting to create duplicate models should fail
		// Using timestamp to ensure uniqueness across test runs
		timestamp := strconv.FormatInt(time.Now().UnixMilli(), 10)
		modelName := "constraint-test-model-" + timestamp
		sourceID := "constraint-test-source-" + timestamp

		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of(modelName),
				ExternalID: apiutils.Of("constraint-test-1-" + timestamp),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of(sourceID)},
			},
		}

		// First model should save successfully
		savedModel1, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		require.NotNil(t, savedModel1)

		// Now test that GetModel works correctly with the single saved model
		retrievedModel, err := dbCatalog.GetModel(ctx, modelName, sourceID)
		require.NoError(t, err)
		require.NotNil(t, retrievedModel)
		assert.Equal(t, modelName, retrievedModel.Name)

		// Test attempting to create a duplicate with same name but different external ID
		// This should fail due to database constraints (which is expected behavior)
		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of(modelName),                        // Same name
				ExternalID: apiutils.Of("constraint-test-2-" + timestamp), // Different external ID
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of(sourceID)},
			},
		}

		_, err = catalogModelRepo.Save(model2)
		// This should fail due to database constraints preventing duplicate names
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "duplicated key")
	})

	t.Run("TestListModels_Success", func(t *testing.T) {
		// Create test models
		sourceIDs := []string{"list-test-source"}

		model1 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("list-test-model-1"),
				ExternalID: apiutils.Of("list-test-1"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("list-test-source")},
				{Name: "description", StringValue: apiutils.Of("First test model")},
			},
		}

		model2 := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("list-test-model-2"),
				ExternalID: apiutils.Of("list-test-2"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("list-test-source")},
				{Name: "description", StringValue: apiutils.Of("Second test model")},
			},
		}

		_, err := catalogModelRepo.Save(model1)
		require.NoError(t, err)
		_, err = catalogModelRepo.Save(model2)
		require.NoError(t, err)

		// Test ListModels
		params := ListModelsParams{
			SourceIDs:     sourceIDs,
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(result.Items), 2, "Should return at least 2 models")
		assert.Equal(t, int32(10), result.PageSize)
		assert.GreaterOrEqual(t, result.Size, int32(2))

		// Verify models are properly mapped
		modelNames := make(map[string]bool)
		for _, model := range result.Items {
			modelNames[model.Name] = true
			// Verify required fields are present
			assert.NotEmpty(t, *model.Id)
			assert.NotEmpty(t, *model.SourceId)
		}

		// Should contain our test models
		foundCount := 0
		if modelNames["list-test-model-1"] {
			foundCount++
		}
		if modelNames["list-test-model-2"] {
			foundCount++
		}
		assert.GreaterOrEqual(t, foundCount, 2, "Should find our test models")
	})

	t.Run("TestListModels_WithPagination", func(t *testing.T) {
		// Test pagination
		sourceIDs := []string{"pagination-test-source"}

		// Create multiple models
		for i := 0; i < 5; i++ {
			model := &models.CatalogModelImpl{
				TypeID: apiutils.Of(int32(catalogModelTypeID)),
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of(fmt.Sprintf("pagination-test-model-%d", i)),
					ExternalID: apiutils.Of(fmt.Sprintf("pagination-test-%d", i)),
				},
				Properties: &[]mr_models.Properties{
					{Name: "source_id", StringValue: apiutils.Of("pagination-test-source")},
				},
			}
			_, err := catalogModelRepo.Save(model)
			require.NoError(t, err)
		}

		params := ListModelsParams{
			SourceIDs:     sourceIDs,
			PageSize:      3,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.ListModels(ctx, params)
		require.NoError(t, err)

		assert.LessOrEqual(t, len(result.Items), 3, "Should respect page size")
		assert.Equal(t, int32(3), result.PageSize)
	})

	t.Run("TestGetArtifacts_Success", func(t *testing.T) {
		// Create test model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("artifact-test-model"),
				ExternalID: apiutils.Of("artifact-test-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("artifact-test-source")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create test artifacts
		modelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact"),
				ExternalID:   apiutils.Of("test-model-artifact-ext"),
				URI:          apiutils.Of("s3://test/model.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact"),
				ExternalID:   apiutils.Of("test-metrics-artifact-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}

		savedModelArt, err := modelArtifactRepo.Save(modelArtifact, savedModel.GetID())
		require.NoError(t, err)
		savedMetricsArt, err := metricsArtifactRepo.Save(metricsArtifact, savedModel.GetID())
		require.NoError(t, err)

		// Test GetArtifacts
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.GetArtifacts(ctx, "artifact-test-model", "artifact-test-source", params)
		require.NoError(t, err)

		assert.GreaterOrEqual(t, len(result.Items), 2, "Should return at least 2 artifacts")
		assert.Equal(t, int32(10), result.PageSize)

		// Verify both types of artifacts are returned
		var modelArtifactFound, metricsArtifactFound bool
		artifactIDs := make(map[string]bool)

		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil {
				modelArtifactFound = true
				artifactIDs[*artifact.CatalogModelArtifact.Id] = true
				assert.Equal(t, "model-artifact", artifact.CatalogModelArtifact.ArtifactType)
			}
			if artifact.CatalogMetricsArtifact != nil {
				metricsArtifactFound = true
				artifactIDs[*artifact.CatalogMetricsArtifact.Id] = true
				assert.Equal(t, "metrics-artifact", artifact.CatalogMetricsArtifact.ArtifactType)
			}
		}

		assert.True(t, modelArtifactFound, "Should find model artifact")
		assert.True(t, metricsArtifactFound, "Should find metrics artifact")

		// Verify our specific artifacts are in the results
		modelArtifactIDStr := strconv.FormatInt(int64(*savedModelArt.GetID()), 10)
		metricsArtifactIDStr := strconv.FormatInt(int64(*savedMetricsArt.GetID()), 10)
		assert.True(t, artifactIDs[modelArtifactIDStr], "Should contain our model artifact")
		assert.True(t, artifactIDs[metricsArtifactIDStr], "Should contain our metrics artifact")
	})

	t.Run("TestGetArtifacts_ModelNotFound", func(t *testing.T) {
		// Test with non-existent model
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		_, err := dbCatalog.GetArtifacts(ctx, "non-existent-model", "test-source", params)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid model name")
	})

	t.Run("TestGetArtifacts_WithCustomProperties", func(t *testing.T) {
		// Create model
		testModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("custom-props-model"),
				ExternalID: apiutils.Of("custom-props-model-ext"),
			},
			Properties: &[]mr_models.Properties{
				{Name: "source_id", StringValue: apiutils.Of("custom-props-source")},
			},
		}

		savedModel, err := catalogModelRepo.Save(testModel)
		require.NoError(t, err)

		// Create artifact with custom properties
		customProps := []mr_models.Properties{
			{Name: "custom_prop_1", StringValue: apiutils.Of("value_1")},
			{Name: "custom_prop_2", StringValue: apiutils.Of("value_2")},
		}

		artifactWithProps := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("artifact-with-props"),
				ExternalID:   apiutils.Of("artifact-with-props-ext"),
				URI:          apiutils.Of("s3://test/props.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
			CustomProperties: &customProps,
		}

		_, err = modelArtifactRepo.Save(artifactWithProps, savedModel.GetID())
		require.NoError(t, err)

		// Get artifacts and verify custom properties
		params := ListArtifactsParams{
			PageSize:      10,
			OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
			SortOrder:     model.SORTORDER_ASC,
			NextPageToken: apiutils.Of(""),
		}

		result, err := dbCatalog.GetArtifacts(ctx, "custom-props-model", "custom-props-source", params)
		require.NoError(t, err)

		// Find our artifact and check custom properties
		found := false
		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil &&
				artifact.CatalogModelArtifact.Name != nil &&
				*artifact.CatalogModelArtifact.Name == "artifact-with-props" {

				found = true
				assert.NotNil(t, artifact.CatalogModelArtifact.CustomProperties)

				// Verify custom properties are present and properly converted
				customPropsMap := *artifact.CatalogModelArtifact.CustomProperties
				assert.Contains(t, customPropsMap, "custom_prop_1")
				assert.Contains(t, customPropsMap, "custom_prop_2")

				// Verify the values are properly converted to MetadataValue
				prop1 := customPropsMap["custom_prop_1"]
				assert.NotNil(t, prop1.MetadataStringValue)
				assert.Equal(t, "value_1", prop1.MetadataStringValue.StringValue)

				break
			}
		}
		assert.True(t, found, "Should find artifact with custom properties")
	})

	t.Run("TestMappingFunctions", func(t *testing.T) {
		t.Run("TestMapCatalogModelToCatalogModel", func(t *testing.T) {
			// Create a catalog model with various properties
			catalogModel := &models.CatalogModelImpl{
				ID:     apiutils.Of(int32(123)),
				TypeID: apiutils.Of(int32(catalogModelTypeID)),
				Attributes: &models.CatalogModelAttributes{
					Name:                     apiutils.Of("mapping-test-model"),
					ExternalID:               apiutils.Of("mapping-test-ext"),
					CreateTimeSinceEpoch:     apiutils.Of(int64(1234567890)),
					LastUpdateTimeSinceEpoch: apiutils.Of(int64(1234567891)),
				},
				Properties: &[]mr_models.Properties{
					{Name: "source_id", StringValue: apiutils.Of("test-source")},
					{Name: "description", StringValue: apiutils.Of("Test description")},
					{Name: "library_name", StringValue: apiutils.Of("pytorch")},
					{Name: "language", StringValue: apiutils.Of("[\"python\", \"go\"]")},
					{Name: "tasks", StringValue: apiutils.Of("[\"classification\", \"regression\"]")},
				},
			}

			result := mapCatalogModelToCatalogModel(catalogModel)

			assert.Equal(t, "123", *result.Id)
			assert.Equal(t, "mapping-test-model", result.Name)
			assert.Equal(t, "mapping-test-ext", *result.ExternalId)
			assert.Equal(t, "test-source", *result.SourceId)
			assert.Equal(t, "Test description", *result.Description)
			assert.Equal(t, "pytorch", *result.LibraryName)
			assert.Equal(t, "1234567890", *result.CreateTimeSinceEpoch)
			assert.Equal(t, "1234567891", *result.LastUpdateTimeSinceEpoch)

			// Verify JSON arrays are properly parsed
			assert.Equal(t, []string{"python", "go"}, result.Language)
			assert.Equal(t, []string{"classification", "regression"}, result.Tasks)
		})

		t.Run("TestMapCatalogArtifactToCatalogArtifact", func(t *testing.T) {
			// Test model artifact mapping
			var catalogModelArtifact models.CatalogModelArtifact = &models.CatalogModelArtifactImpl{
				ID:     apiutils.Of(int32(456)),
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of("test-model-artifact"),
					ExternalID: apiutils.Of("test-model-artifact-ext"),
					URI:        apiutils.Of("s3://test/model.bin"),
				},
			}

			catalogArtifact := models.CatalogArtifact{
				CatalogModelArtifact: &catalogModelArtifact,
			}

			result, err := mapCatalogArtifactToCatalogArtifact(catalogArtifact)
			require.NoError(t, err)

			assert.NotNil(t, result.CatalogModelArtifact)
			assert.Nil(t, result.CatalogMetricsArtifact)
			assert.Equal(t, "456", *result.CatalogModelArtifact.Id)
			assert.Equal(t, "test-model-artifact", *result.CatalogModelArtifact.Name)
			assert.Equal(t, "s3://test/model.bin", result.CatalogModelArtifact.Uri)

			// Test metrics artifact mapping
			var catalogMetricsArtifact models.CatalogMetricsArtifact = &models.CatalogMetricsArtifactImpl{
				ID:     apiutils.Of(int32(789)),
				TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("test-metrics-artifact"),
					ExternalID:  apiutils.Of("test-metrics-artifact-ext"),
					MetricsType: models.MetricsTypePerformance,
				},
			}

			catalogArtifact2 := models.CatalogArtifact{
				CatalogMetricsArtifact: &catalogMetricsArtifact,
			}

			result2, err := mapCatalogArtifactToCatalogArtifact(catalogArtifact2)
			require.NoError(t, err)

			assert.Nil(t, result2.CatalogModelArtifact)
			assert.NotNil(t, result2.CatalogMetricsArtifact)
			assert.Equal(t, "789", *result2.CatalogMetricsArtifact.Id)
			assert.Equal(t, "test-metrics-artifact", *result2.CatalogMetricsArtifact.Name)
			assert.Equal(t, "performance-metrics", result2.CatalogMetricsArtifact.MetricsType)
		})

		t.Run("TestMapCatalogArtifact_EmptyArtifact", func(t *testing.T) {
			// Test with empty catalog artifact
			emptyCatalogArtifact := models.CatalogArtifact{}

			_, err := mapCatalogArtifactToCatalogArtifact(emptyCatalogArtifact)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid catalog artifact type")
		})
	})

	t.Run("TestErrorHandling", func(t *testing.T) {
		t.Run("TestGetArtifacts_InvalidModelID", func(t *testing.T) {
			// Create a model with invalid ID format for testing
			// This would be an edge case where the ID isn't a valid integer

			// We can't easily test this directly since IDs are generated as integers
			// But we can test the error case by mocking a scenario

			// For now, let's test a scenario where the model exists but has some issue
			params := ListArtifactsParams{
				PageSize:      10,
				OrderBy:       model.ORDERBYFIELD_CREATE_TIME,
				SortOrder:     model.SORTORDER_ASC,
				NextPageToken: apiutils.Of(""),
			}

			_, err := dbCatalog.GetArtifacts(ctx, "non-existent-model", "test-source", params)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "invalid model name")
		})
	})
}

// Helper functions to get type IDs from database

func getCatalogModelTypeIDForDBTest(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}
	return int64(typeRecord.ID)
}

func getCatalogModelArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}
	return int64(typeRecord.ID)
}

func getCatalogMetricsArtifactTypeIDForDBTest(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}
	return int64(typeRecord.ID)
}
