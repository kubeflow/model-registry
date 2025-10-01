package service_test

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCatalogArtifactRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get the catalog artifact type IDs
	modelArtifactTypeID := getCatalogModelArtifactTypeID(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeID(t, sharedDB)

	// Create unified artifact repository with both types
	artifactTypeMap := map[string]int64{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	}
	repo := service.NewCatalogArtifactRepository(sharedDB, artifactTypeMap)

	// Also get CatalogModel type ID for creating parent entities
	catalogModelTypeID := getCatalogModelTypeID(t, sharedDB)
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)

	// Create shared test data
	catalogModel := &models.CatalogModelImpl{
		TypeID: apiutils.Of(int32(catalogModelTypeID)),
		Attributes: &models.CatalogModelAttributes{
			Name:       apiutils.Of("test-catalog-model-for-artifacts"),
			ExternalID: apiutils.Of("catalog-model-artifacts-ext-123"),
		},
	}
	savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
	require.NoError(t, err)

	t.Run("GetByID_ModelArtifact", func(t *testing.T) {
		// Create a model artifact using the specific repository
		modelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-getbyid"),
				ExternalID:   apiutils.Of("model-art-getbyid-ext-123"),
				URI:          apiutils.Of("s3://test-bucket/model.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}
		savedModelArtifact, err := modelArtifactRepo.Save(modelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedModelArtifact.GetID())
		require.NoError(t, err)

		// Verify it's a model artifact
		assert.NotNil(t, retrieved.CatalogModelArtifact)
		assert.Nil(t, retrieved.CatalogMetricsArtifact)
		assert.Equal(t, "test-model-artifact-getbyid", *retrieved.CatalogModelArtifact.GetAttributes().Name)
		assert.Equal(t, "model-art-getbyid-ext-123", *retrieved.CatalogModelArtifact.GetAttributes().ExternalID)
		assert.Equal(t, "s3://test-bucket/model.bin", *retrieved.CatalogModelArtifact.GetAttributes().URI)
	})

	t.Run("GetByID_MetricsArtifact", func(t *testing.T) {
		// Create a metrics artifact using the specific repository
		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact-getbyid"),
				ExternalID:   apiutils.Of("metrics-art-getbyid-ext-123"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedMetricsArtifact, err := metricsArtifactRepo.Save(metricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedMetricsArtifact.GetID())
		require.NoError(t, err)

		// Verify it's a metrics artifact
		assert.Nil(t, retrieved.CatalogModelArtifact)
		assert.NotNil(t, retrieved.CatalogMetricsArtifact)
		assert.Equal(t, "test-metrics-artifact-getbyid", *retrieved.CatalogMetricsArtifact.GetAttributes().Name)
		assert.Equal(t, "metrics-art-getbyid-ext-123", *retrieved.CatalogMetricsArtifact.GetAttributes().ExternalID)
		assert.Equal(t, models.MetricsTypeAccuracy, retrieved.CatalogMetricsArtifact.GetAttributes().MetricsType)
	})

	t.Run("GetByID_NotFound", func(t *testing.T) {
		nonExistentID := int32(99999)
		_, err := repo.GetByID(nonExistentID)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "catalog artifact by id not found")
	})

	t.Run("List_AllArtifacts", func(t *testing.T) {
		// Create test artifacts of both types
		modelArtifact1 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-list-1"),
				ExternalID:   apiutils.Of("model-list-1-ext"),
				URI:          apiutils.Of("s3://test/model1.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		modelArtifact2 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("test-model-artifact-list-2"),
				ExternalID:   apiutils.Of("model-list-2-ext"),
				URI:          apiutils.Of("s3://test/model2.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
		}

		metricsArtifact1 := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-metrics-artifact-list-1"),
				ExternalID:   apiutils.Of("metrics-list-1-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}

		// Save artifacts
		savedModelArt1, err := modelArtifactRepo.Save(modelArtifact1, savedCatalogModel.GetID())
		require.NoError(t, err)
		savedModelArt2, err := modelArtifactRepo.Save(modelArtifact2, savedCatalogModel.GetID())
		require.NoError(t, err)
		savedMetricsArt1, err := metricsArtifactRepo.Save(metricsArtifact1, savedCatalogModel.GetID())
		require.NoError(t, err)

		// List all artifacts for the parent resource
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return all 3 artifacts (2 model + 1 metrics)
		assert.GreaterOrEqual(t, len(result.Items), 3, "Should return at least the 3 artifacts we created")

		// Verify we got both types
		var modelArtifactCount, metricsArtifactCount int
		artifactIDs := make(map[int32]bool)

		for _, artifact := range result.Items {
			if artifact.CatalogModelArtifact != nil {
				modelArtifactCount++
				artifactIDs[*artifact.CatalogModelArtifact.GetID()] = true
			} else if artifact.CatalogMetricsArtifact != nil {
				metricsArtifactCount++
				artifactIDs[*artifact.CatalogMetricsArtifact.GetID()] = true
			}
		}

		assert.GreaterOrEqual(t, modelArtifactCount, 2, "Should have at least 2 model artifacts")
		assert.GreaterOrEqual(t, metricsArtifactCount, 1, "Should have at least 1 metrics artifact")

		// Verify our specific artifacts are in the results
		assert.True(t, artifactIDs[*savedModelArt1.GetID()], "Should contain first model artifact")
		assert.True(t, artifactIDs[*savedModelArt2.GetID()], "Should contain second model artifact")
		assert.True(t, artifactIDs[*savedMetricsArt1.GetID()], "Should contain metrics artifact")
	})

	t.Run("List_FilterByArtifactType_ModelArtifact", func(t *testing.T) {
		// Filter by model artifact type only
		artifactType := "model-artifact"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &artifactType,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// All results should be model artifacts
		for _, artifact := range result.Items {
			assert.NotNil(t, artifact.CatalogModelArtifact, "Should only return model artifacts")
			assert.Nil(t, artifact.CatalogMetricsArtifact, "Should not return metrics artifacts")
		}
	})

	t.Run("List_FilterByArtifactType_MetricsArtifact", func(t *testing.T) {
		// Filter by metrics artifact type only
		artifactType := "metrics-artifact"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &artifactType,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// All results should be metrics artifacts
		for _, artifact := range result.Items {
			assert.Nil(t, artifact.CatalogModelArtifact, "Should not return model artifacts")
			assert.NotNil(t, artifact.CatalogMetricsArtifact, "Should only return metrics artifacts")
		}
	})

	t.Run("List_FilterByExternalID", func(t *testing.T) {
		// Create artifact with specific external ID for filtering
		testArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("external-id-filter-test"),
				ExternalID:   apiutils.Of("unique-external-id-123"),
				MetricsType:  models.MetricsTypePerformance,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedArtifact, err := metricsArtifactRepo.Save(testArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Filter by external ID
		externalID := "unique-external-id-123"
		listOptions := models.CatalogArtifactListOptions{
			ExternalID: &externalID,
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 1, "Should find exactly one artifact with the external ID")

		// Verify it's the correct artifact
		artifact := result.Items[0]
		assert.NotNil(t, artifact.CatalogMetricsArtifact)
		assert.Equal(t, *savedArtifact.GetID(), *artifact.CatalogMetricsArtifact.GetID())
		assert.Equal(t, "unique-external-id-123", *artifact.CatalogMetricsArtifact.GetAttributes().ExternalID)
	})

	t.Run("List_WithPagination", func(t *testing.T) {
		// Create multiple artifacts for pagination testing
		for i := 0; i < 5; i++ {
			artifact := &models.CatalogModelArtifactImpl{
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("pagination-test-%d", i)),
					ExternalID:   apiutils.Of(fmt.Sprintf("pagination-ext-%d", i)),
					URI:          apiutils.Of(fmt.Sprintf("s3://test/pagination-%d.bin", i)),
					ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
				},
			}
			_, err := modelArtifactRepo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
		}

		// Test pagination
		pageSize := int32(3)
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			Pagination: dbmodels.Pagination{
				PageSize: &pageSize,
				OrderBy:  apiutils.Of("ID"),
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Items), 3, "Should respect page size limit")
		assert.GreaterOrEqual(t, len(result.Items), 1, "Should return at least one item")
	})

	t.Run("List_InvalidArtifactType", func(t *testing.T) {
		// Test with invalid artifact type
		invalidType := "invalid-artifact-type"
		listOptions := models.CatalogArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			ArtifactType:     &invalidType,
		}

		_, err := repo.List(listOptions)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid catalog artifact type")
		assert.Contains(t, err.Error(), "invalid-artifact-type")
	})

	t.Run("List_WithCustomProperties", func(t *testing.T) {
		// Create artifacts with custom properties
		customProps := []dbmodels.Properties{
			{
				Name:        "custom_prop_1",
				StringValue: apiutils.Of("custom_value_1"),
			},
			{
				Name:        "custom_prop_2",
				StringValue: apiutils.Of("custom_value_2"),
			},
		}

		artifactWithCustomProps := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:         apiutils.Of("artifact-with-custom-props"),
				ExternalID:   apiutils.Of("custom-props-ext"),
				URI:          apiutils.Of("s3://test/custom-props.bin"),
				ArtifactType: apiutils.Of(models.CatalogModelArtifactType),
			},
			CustomProperties: &customProps,
		}

		savedArtifact, err := modelArtifactRepo.Save(artifactWithCustomProps, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Retrieve using unified repository
		retrieved, err := repo.GetByID(*savedArtifact.GetID())
		require.NoError(t, err)

		// Verify custom properties are preserved
		require.NotNil(t, retrieved.CatalogModelArtifact)
		assert.NotNil(t, retrieved.CatalogModelArtifact.GetCustomProperties())

		customPropsMap := make(map[string]string)
		for _, prop := range *retrieved.CatalogModelArtifact.GetCustomProperties() {
			if prop.StringValue != nil {
				customPropsMap[prop.Name] = *prop.StringValue
			}
		}

		assert.Equal(t, "custom_value_1", customPropsMap["custom_prop_1"])
		assert.Equal(t, "custom_value_2", customPropsMap["custom_prop_2"])
	})

	t.Run("MappingErrors", func(t *testing.T) {
		// Test error handling for invalid type mapping
		// This would typically happen if there's data inconsistency in the database

		// We can't easily test this without directly manipulating the database
		// but we can test the GetByID with an artifact that has an unknown type
		// by temporarily modifying the repository's type mapping

		// Create a repository with incomplete type mapping
		incompleteTypeMap := map[string]int64{
			service.CatalogModelArtifactTypeName: modelArtifactTypeID,
			// Missing CatalogMetricsArtifactTypeName intentionally
		}
		incompleteRepo := service.NewCatalogArtifactRepository(sharedDB, incompleteTypeMap)

		// Create a metrics artifact first using the complete repo
		metricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(metricsArtifactTypeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:         apiutils.Of("test-mapping-error"),
				ExternalID:   apiutils.Of("mapping-error-ext"),
				MetricsType:  models.MetricsTypeAccuracy,
				ArtifactType: apiutils.Of("metrics-artifact"),
			},
		}
		savedMetricsArtifact, err := metricsArtifactRepo.Save(metricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Try to retrieve using incomplete repo - should get mapping error
		_, err = incompleteRepo.GetByID(*savedMetricsArtifact.GetID())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "invalid catalog artifact type")
	})
}
