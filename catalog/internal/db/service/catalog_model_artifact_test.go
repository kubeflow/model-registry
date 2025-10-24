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

func TestCatalogModelArtifactRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get the CatalogModelArtifact type ID
	typeID := getCatalogModelArtifactTypeID(t, sharedDB)
	repo := service.NewCatalogModelArtifactRepository(sharedDB, typeID)

	// Also get CatalogModel type ID for creating parent entities
	catalogModelTypeID := getCatalogModelTypeID(t, sharedDB)
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a catalog model for attribution
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-artifact"),
				ExternalID: apiutils.Of("catalog-model-ext-123"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Test creating a new catalog model artifact
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("test-catalog-model-artifact"),
				ExternalID: apiutils.Of("catalog-artifact-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/model.pkl"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test catalog model artifact description"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "custom-catalog-prop",
					StringValue: apiutils.Of("custom-catalog-value"),
				},
			},
		}

		saved, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-catalog-model-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, "catalog-artifact-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://catalog-bucket/model.pkl", *saved.GetAttributes().URI)

		// Test updating the same catalog model artifact
		catalogModelArtifact.ID = saved.GetID()
		catalogModelArtifact.GetAttributes().Name = apiutils.Of("updated-catalog-model-artifact")
		catalogModelArtifact.GetAttributes().URI = apiutils.Of("s3://catalog-bucket/updated-model.pkl")
		// Preserve CreateTimeSinceEpoch from the saved entity
		catalogModelArtifact.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-catalog-model-artifact", *updated.GetAttributes().Name)
		assert.Equal(t, "s3://catalog-bucket/updated-model.pkl", *updated.GetAttributes().URI)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-getbyid"),
				ExternalID: apiutils.Of("catalog-model-getbyid-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create a catalog model artifact to retrieve
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("get-test-catalog-model-artifact"),
				ExternalID: apiutils.Of("get-catalog-artifact-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/get-model.pkl"),
			},
		}

		saved, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-catalog-model-artifact", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-catalog-artifact-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://catalog-bucket/get-model.pkl", *retrieved.GetAttributes().URI)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.ErrorIs(t, err, service.ErrCatalogModelArtifactNotFound)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a catalog model for the artifacts
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-list"),
				ExternalID: apiutils.Of("catalog-model-list-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create multiple catalog model artifacts for listing
		testArtifacts := []*models.CatalogModelArtifactImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of("list-catalog-artifact-1"),
					ExternalID: apiutils.Of("list-catalog-artifact-ext-1"),
					URI:        apiutils.Of("s3://catalog-bucket/list-model-1.pkl"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of("list-catalog-artifact-2"),
					ExternalID: apiutils.Of("list-catalog-artifact-ext-2"),
					URI:        apiutils.Of("s3://catalog-bucket/list-model-2.pkl"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of("list-catalog-artifact-3"),
					ExternalID: apiutils.Of("list-catalog-artifact-ext-3"),
					URI:        apiutils.Of("s3://catalog-bucket/list-model-3.pkl"),
				},
			},
		}

		// Save all test artifacts
		var savedArtifacts []models.CatalogModelArtifact
		for _, artifact := range testArtifacts {
			saved, err := repo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
			savedArtifacts = append(savedArtifacts, saved)
		}

		// Test listing all artifacts
		listOptions := models.CatalogModelArtifactListOptions{}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test artifacts

		// Test filtering by name
		nameFilter := "list-catalog-artifact-1"
		listOptions = models.CatalogModelArtifactListOptions{
			Name: &nameFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-catalog-artifact-1", *result.Items[0].GetAttributes().Name)
		}

		// Test filtering by external ID
		externalIDFilter := "list-catalog-artifact-ext-2"
		listOptions = models.CatalogModelArtifactListOptions{
			ExternalID: &externalIDFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-catalog-artifact-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test filtering by parent resource ID (catalog model)
		listOptions = models.CatalogModelArtifactListOptions{
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
				Name:       apiutils.Of("test-catalog-model-for-props"),
				ExternalID: apiutils.Of("catalog-model-props-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create a catalog model artifact with both properties and custom properties
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("props-test-catalog-artifact"),
				ExternalID: apiutils.Of("props-catalog-artifact-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/props-model.pkl"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:     "size_bytes",
					IntValue: apiutils.Of(int32(2048000)),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("catalog-ml-team"),
				},
				{
					Name:      "is_public",
					BoolValue: apiutils.Of(true),
				},
			},
		}

		saved, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Retrieve and verify properties
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)

		// Check regular properties
		require.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 2)

		// Check custom properties
		require.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify specific properties exist
		properties := *retrieved.GetProperties()
		var foundVersion, foundSizeBytes bool
		for _, prop := range properties {
			switch prop.Name {
			case "version":
				foundVersion = true
				assert.Equal(t, "1.0.0", *prop.StringValue)
			case "size_bytes":
				foundSizeBytes = true
				assert.Equal(t, int32(2048000), *prop.IntValue)
			}
		}
		assert.True(t, foundVersion, "Should find version property")
		assert.True(t, foundSizeBytes, "Should find size_bytes property")

		// Verify custom properties
		customProperties := *retrieved.GetCustomProperties()
		var foundTeam, foundIsPublic bool
		for _, prop := range customProperties {
			switch prop.Name {
			case "team":
				foundTeam = true
				assert.Equal(t, "catalog-ml-team", *prop.StringValue)
			case "is_public":
				foundIsPublic = true
				assert.Equal(t, true, *prop.BoolValue)
			}
		}
		assert.True(t, foundTeam, "Should find team custom property")
		assert.True(t, foundIsPublic, "Should find is_public custom property")
	})

	t.Run("TestSaveWithoutParentResource", func(t *testing.T) {
		// Test creating a catalog model artifact without parent resource attribution
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("standalone-catalog-artifact"),
				ExternalID: apiutils.Of("standalone-catalog-artifact-ext"),
				URI:        apiutils.Of("s3://catalog-bucket/standalone-model.pkl"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Standalone catalog artifact without parent"),
				},
			},
		}

		saved, err := repo.Save(catalogModelArtifact, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "standalone-catalog-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, "s3://catalog-bucket/standalone-model.pkl", *saved.GetAttributes().URI)

		// Verify it can be retrieved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		assert.Equal(t, "standalone-catalog-artifact", *retrieved.GetAttributes().Name)
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-ordering"),
				ExternalID: apiutils.Of("catalog-model-ordering-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create artifacts sequentially with time delays to ensure deterministic ordering
		artifact1 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("time-test-catalog-artifact-1"),
				ExternalID: apiutils.Of("time-catalog-artifact-ext-1"),
				URI:        apiutils.Of("s3://catalog-bucket/time-model-1.pkl"),
			},
		}
		saved1, err := repo.Save(artifact1, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		artifact2 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("time-test-catalog-artifact-2"),
				ExternalID: apiutils.Of("time-catalog-artifact-ext-2"),
				URI:        apiutils.Of("s3://catalog-bucket/time-model-2.pkl"),
			},
		}
		saved2, err := repo.Save(artifact2, savedCatalogModel.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		listOptions := models.CatalogModelArtifactListOptions{
			Pagination: dbmodels.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test artifacts in the results
		var foundArtifact1, foundArtifact2 models.CatalogModelArtifact
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

	t.Run("TestListPagination", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-pagination"),
				ExternalID: apiutils.Of("catalog-model-pagination-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create multiple artifacts for pagination testing
		for i := 0; i < 5; i++ {
			artifact := &models.CatalogModelArtifactImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.CatalogModelArtifactAttributes{
					Name:       apiutils.Of(fmt.Sprintf("pagination-artifact-%d", i)),
					ExternalID: apiutils.Of(fmt.Sprintf("pagination-artifact-ext-%d", i)),
					URI:        apiutils.Of(fmt.Sprintf("s3://catalog-bucket/pagination-model-%d.pkl", i)),
				},
			}
			_, err := repo.Save(artifact, savedCatalogModel.GetID())
			require.NoError(t, err)
		}

		// Test pagination with page size
		pageSize := int32(2)
		listOptions := models.CatalogModelArtifactListOptions{
			ParentResourceID: savedCatalogModel.GetID(),
			Pagination: dbmodels.Pagination{
				PageSize: &pageSize,
				OrderBy:  apiutils.Of("ID"),
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.LessOrEqual(t, len(result.Items), 2, "Should respect page size limit")
		assert.GreaterOrEqual(t, len(result.Items), 1, "Should return at least one item")
	})

	t.Run("TestSaveWithTypeIDSetting", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-typeid"),
				ExternalID: apiutils.Of("catalog-model-typeid-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Test creating artifact without explicit type_id (should be set automatically)
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			// Intentionally not setting TypeID to test auto-setting
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("typeid-test-artifact"),
				ExternalID: apiutils.Of("typeid-artifact-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/typeid-model.pkl"),
			},
		}

		saved, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetTypeID())
		assert.Equal(t, int32(typeID), *saved.GetTypeID())
		assert.Equal(t, "typeid-test-artifact", *saved.GetAttributes().Name)

		// Test with explicitly set type_id (should not be overridden)
		explicitTypeID := int32(typeID)
		catalogModelArtifact2 := &models.CatalogModelArtifactImpl{
			TypeID: &explicitTypeID,
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("explicit-typeid-artifact"),
				ExternalID: apiutils.Of("explicit-typeid-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/explicit-model.pkl"),
			},
		}

		saved2, err := repo.Save(catalogModelArtifact2, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved2)
		require.NotNil(t, saved2.GetTypeID())
		assert.Equal(t, explicitTypeID, *saved2.GetTypeID())
	})

	t.Run("TestSaveWithNameMatching", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-name-matching"),
				ExternalID: apiutils.Of("catalog-model-name-match-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Create initial artifact
		artifactName := "name-matching-artifact"
		catalogModelArtifact1 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of(artifactName),
				ExternalID: apiutils.Of("name-match-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/original.pkl"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
		}

		saved1, err := repo.Save(catalogModelArtifact1, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved1)
		originalID := *saved1.GetID()
		assert.Equal(t, artifactName, *saved1.GetAttributes().Name)
		assert.Equal(t, "s3://catalog-bucket/original.pkl", *saved1.GetAttributes().URI)

		// Create second artifact with same name (should update existing)
		catalogModelArtifact2 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of(artifactName), // Same name
				ExternalID: apiutils.Of("name-match-ext-456"),
				URI:        apiutils.Of("s3://catalog-bucket/updated.pkl"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("2.0.0"),
				},
			},
		}

		saved2, err := repo.Save(catalogModelArtifact2, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved2)

		// Should have same ID (updated existing)
		assert.Equal(t, originalID, *saved2.GetID())
		assert.Equal(t, artifactName, *saved2.GetAttributes().Name)
		assert.Equal(t, "s3://catalog-bucket/updated.pkl", *saved2.GetAttributes().URI)
		assert.Equal(t, "name-match-ext-456", *saved2.GetAttributes().ExternalID)

		// Verify by retrieving from database
		retrieved, err := repo.GetByID(originalID)
		require.NoError(t, err)
		assert.Equal(t, "s3://catalog-bucket/updated.pkl", *retrieved.GetAttributes().URI)
		assert.Equal(t, "name-match-ext-456", *retrieved.GetAttributes().ExternalID)

		// Verify properties were updated
		require.NotNil(t, retrieved.GetProperties())
		properties := *retrieved.GetProperties()
		var foundVersion bool
		for _, prop := range properties {
			if prop.Name == "version" {
				foundVersion = true
				assert.Equal(t, "2.0.0", *prop.StringValue)
				break
			}
		}
		assert.True(t, foundVersion, "Should find updated version property")

		// Test that artifact with different name creates new entity
		catalogModelArtifact3 := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("different-name-artifact"),
				ExternalID: apiutils.Of("different-name-ext-789"),
				URI:        apiutils.Of("s3://catalog-bucket/different.pkl"),
			},
		}

		saved3, err := repo.Save(catalogModelArtifact3, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved3)

		// Should have different ID (new entity)
		assert.NotEqual(t, originalID, *saved3.GetID())
		assert.Equal(t, "different-name-artifact", *saved3.GetAttributes().Name)
	})

	t.Run("TestSaveWithNameMatchingNoExistingName", func(t *testing.T) {
		// Create a catalog model
		catalogModel := &models.CatalogModelImpl{
			TypeID: apiutils.Of(int32(catalogModelTypeID)),
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model-for-no-match"),
				ExternalID: apiutils.Of("catalog-model-no-match-ext"),
			},
		}
		savedCatalogModel, err := catalogModelRepo.Save(catalogModel)
		require.NoError(t, err)

		// Test saving artifact when no existing artifact with same name exists
		catalogModelArtifact := &models.CatalogModelArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.CatalogModelArtifactAttributes{
				Name:       apiutils.Of("unique-artifact-name"),
				ExternalID: apiutils.Of("unique-ext-123"),
				URI:        apiutils.Of("s3://catalog-bucket/unique.pkl"),
			},
		}

		saved, err := repo.Save(catalogModelArtifact, savedCatalogModel.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "unique-artifact-name", *saved.GetAttributes().Name)
		assert.Equal(t, "s3://catalog-bucket/unique.pkl", *saved.GetAttributes().URI)
	})
}

// Helper function to get or create CatalogModelArtifact type ID
func getCatalogModelArtifactTypeID(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	}

	return typeRecord.ID
}
