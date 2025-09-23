package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestCatalogMetricsArtifactRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get the CatalogMetricsArtifact type ID
	typeID := getCatalogMetricsArtifactTypeID(t, sharedDB)
	repo := service.NewCatalogMetricsArtifactRepository(sharedDB, typeID)

	// Also get CatalogModel type ID for creating parent entities
	catalogModelTypeID := getCatalogModelTypeID(t, sharedDB)
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a catalog model for attribution
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-metrics"),
				ExternalID: apiutils.Of("catalog-model-metrics-ext-123"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Test creating a new catalog metrics artifact
		catalogMetricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("test-catalog-metrics-artifact"),
				ExternalID:  apiutils.Of("catalog-metrics-ext-123"),
				MetricsType: models.MetricsTypeAccuracy,
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test catalog metrics artifact description"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "custom-metrics-prop",
					StringValue: apiutils.Of("custom-metrics-value"),
				},
			},
		}

		saved, err := repo.Save(catalogMetricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-catalog-metrics-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, "catalog-metrics-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, models.MetricsTypeAccuracy, saved.GetAttributes().MetricsType)

		// Test updating the same catalog metrics artifact
		catalogMetricsArtifact.ID = saved.GetID()
		catalogMetricsArtifact.GetAttributes().Name = apiutils.Of("updated-catalog-metrics-artifact")
		catalogMetricsArtifact.GetAttributes().MetricsType = models.MetricsTypePerformance
		// Preserve CreateTimeSinceEpoch from the saved entity
		catalogMetricsArtifact.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(catalogMetricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-catalog-metrics-artifact", *updated.GetAttributes().Name)
		assert.Equal(t, models.MetricsTypePerformance, updated.GetAttributes().MetricsType)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-getbyid-metrics"),
				ExternalID: apiutils.Of("catalog-model-getbyid-metrics-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create a catalog metrics artifact to retrieve
		catalogMetricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("get-test-catalog-metrics-artifact"),
				ExternalID:  apiutils.Of("get-catalog-metrics-ext-123"),
				MetricsType: models.MetricsTypeAccuracy,
			},
		}

		saved, err := repo.Save(catalogMetricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-catalog-metrics-artifact", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-catalog-metrics-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, models.MetricsTypeAccuracy, retrieved.GetAttributes().MetricsType)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.ErrorIs(t, err, service.ErrCatalogMetricsArtifactNotFound)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a catalog model for the artifacts
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-list-metrics"),
				ExternalID: apiutils.Of("catalog-model-list-metrics-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create multiple catalog metrics artifacts for listing
		testArtifacts := []*models.CatalogMetricsArtifactImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("list-catalog-metrics-artifact-1"),
					ExternalID:  apiutils.Of("list-catalog-metrics-ext-1"),
					MetricsType: models.MetricsTypeAccuracy,
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("list-catalog-metrics-artifact-2"),
					ExternalID:  apiutils.Of("list-catalog-metrics-ext-2"),
					MetricsType: models.MetricsTypePerformance,
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("list-catalog-metrics-artifact-3"),
					ExternalID:  apiutils.Of("list-catalog-metrics-ext-3"),
					MetricsType: models.MetricsTypePerformance,
				},
			},
		}

		// Save all test artifacts
		var savedArtifacts []models.CatalogMetricsArtifact
		for _, artifact := range testArtifacts {
			saved, err := repo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
			savedArtifacts = append(savedArtifacts, saved)
		}

		// Test listing all artifacts
		listOptions := models.CatalogMetricsArtifactListOptions{}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test artifacts

		// Test filtering by name
		nameFilter := "list-catalog-metrics-artifact-1"
		listOptions = models.CatalogMetricsArtifactListOptions{
			Name: &nameFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-catalog-metrics-artifact-1", *result.Items[0].GetAttributes().Name)
		}

		// Test filtering by external ID
		externalIDFilter := "list-catalog-metrics-ext-2"
		listOptions = models.CatalogMetricsArtifactListOptions{
			ExternalID: &externalIDFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-catalog-metrics-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test filtering by parent resource ID (catalog model)
		listOptions = models.CatalogMetricsArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should find our 3 test artifacts
	})

	t.Run("TestListWithPropertiesAndCustomProperties", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-props-metrics"),
				ExternalID: apiutils.Of("catalog-model-props-metrics-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create a catalog metrics artifact with both properties and custom properties
		catalogMetricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("props-test-catalog-metrics-artifact"),
				ExternalID:  apiutils.Of("props-catalog-metrics-ext-123"),
				MetricsType: models.MetricsTypeAccuracy,
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.95),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("catalog-metrics-team"),
				},
				{
					Name:      "is_validated",
					BoolValue: apiutils.Of(true),
				},
			},
		}

		saved, err := repo.Save(catalogMetricsArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Retrieve and verify properties
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Check that metricsType is properly set
		assert.Equal(t, models.MetricsTypeAccuracy, retrieved.GetAttributes().MetricsType)

		// Check regular properties
		require.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 2)

		// Check custom properties
		require.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify specific properties exist
		properties := *retrieved.GetProperties()
		var foundVersion, foundValue bool
		for _, prop := range properties {
			switch prop.Name {
			case "version":
				foundVersion = true
				assert.Equal(t, "1.0.0", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, 0.95, *prop.DoubleValue)
			}
		}
		assert.True(t, foundVersion, "Should find version property")
		assert.True(t, foundValue, "Should find value property")

		// Verify custom properties
		customProperties := *retrieved.GetCustomProperties()
		var foundTeam, foundIsValidated bool
		for _, prop := range customProperties {
			switch prop.Name {
			case "team":
				foundTeam = true
				assert.Equal(t, "catalog-metrics-team", *prop.StringValue)
			case "is_validated":
				foundIsValidated = true
				assert.Equal(t, true, *prop.BoolValue)
			}
		}
		assert.True(t, foundTeam, "Should find team custom property")
		assert.True(t, foundIsValidated, "Should find is_validated custom property")
	})

	t.Run("TestSaveWithoutParentResource", func(t *testing.T) {
		// Test creating a catalog metrics artifact without parent resource attribution
		catalogMetricsArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("standalone-catalog-metrics-artifact"),
				ExternalID:  apiutils.Of("standalone-catalog-metrics-ext"),
				MetricsType: models.MetricsTypeAccuracy,
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Standalone catalog metrics artifact without parent"),
				},
			},
		}

		saved, err := repo.Save(catalogMetricsArtifact, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "standalone-catalog-metrics-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, models.MetricsTypeAccuracy, saved.GetAttributes().MetricsType)

		// Verify it can be retrieved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		assert.Equal(t, "standalone-catalog-metrics-artifact", *retrieved.GetAttributes().Name)
		assert.Equal(t, models.MetricsTypeAccuracy, retrieved.GetAttributes().MetricsType)
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-ordering-metrics"),
				ExternalID: apiutils.Of("catalog-model-ordering-metrics-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create artifacts sequentially with time delays to ensure deterministic ordering
		artifact1 := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("time-test-catalog-metrics-artifact-1"),
				ExternalID:  apiutils.Of("time-catalog-metrics-ext-1"),
				MetricsType: models.MetricsTypeAccuracy,
			},
		}
		saved1, err := repo.Save(artifact1, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		artifact2 := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("time-test-catalog-metrics-artifact-2"),
				ExternalID:  apiutils.Of("time-catalog-metrics-ext-2"),
				MetricsType: models.MetricsTypePerformance,
			},
		}
		saved2, err := repo.Save(artifact2, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		listOptions := models.CatalogMetricsArtifactListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test artifacts in the results
		var foundArtifact1, foundArtifact2 models.CatalogMetricsArtifact
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundArtifact1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundArtifact2 = item
				index2 = i
			}
		}

		// Verify both artifacts were found and artifact1 comes before artifact2 (ascending order)
		require.NotEqual(t, -1, index1, "Artifact 1 should be found in results")
		require.NotEqual(t, -1, index2, "Artifact 2 should be found in results")
		assert.Less(t, index1, index2, "Artifact 1 should come before Artifact 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundArtifact1.GetAttributes().CreateTimeSinceEpoch, *foundArtifact2.GetAttributes().CreateTimeSinceEpoch, "Artifact 1 should have earlier create time")
	})

	t.Run("TestMetricsTypeField", func(t *testing.T) {
		// Test various metrics types
		metricsTypes := []models.MetricsType{models.MetricsTypeAccuracy, models.MetricsTypePerformance}

		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-metrics-types"),
				ExternalID: apiutils.Of("catalog-model-metrics-types-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		for i, metricsType := range metricsTypes {
			artifact := &models.CatalogMetricsArtifactImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of(fmt.Sprintf("metrics-type-test-%d", i)),
					ExternalID:  apiutils.Of(fmt.Sprintf("metrics-type-ext-%d", i)),
					MetricsType: metricsType,
				},
			}

			saved, err := repo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
			assert.Equal(t, metricsType, saved.GetAttributes().MetricsType)

			// Verify retrieval preserves metricsType
			retrieved, err := repo.GetByID(*saved.GetID())
			require.NoError(t, err)
			assert.Equal(t, metricsType, retrieved.GetAttributes().MetricsType)
		}
	})
}

// Helper function to get or create CatalogMetricsArtifact type ID
func getCatalogMetricsArtifactTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	}

	return int64(typeRecord.ID)
}
