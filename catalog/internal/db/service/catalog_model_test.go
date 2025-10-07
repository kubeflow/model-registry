package service_test

import (
	"testing"

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

func TestCatalogModelRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Create or get the CatalogModel type ID
	typeID := getCatalogModelTypeID(t, sharedDB)
	repo := service.NewCatalogModelRepository(sharedDB, typeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new catalog model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("test-catalog-model"),
				ExternalID: apiutils.Of("catalog-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test catalog model description"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "custom-prop",
					StringValue: apiutils.Of("custom-value"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-catalog-model", *saved.GetAttributes().Name)
		assert.Equal(t, "catalog-ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same model
		catalogModel.ID = saved.GetID()
		catalogModel.GetAttributes().Name = apiutils.Of("updated-catalog-model")
		// Preserve CreateTimeSinceEpoch from the saved entity
		catalogModel.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-catalog-model", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a model to retrieve
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("get-test-catalog-model"),
				ExternalID: apiutils.Of("get-catalog-ext-123"),
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-catalog-model", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-catalog-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.ErrorIs(t, err, service.ErrCatalogModelNotFound)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple models for listing
		testModels := []*models.CatalogModelImpl{
			{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of("list-catalog-model-1"),
					ExternalID: apiutils.Of("list-catalog-ext-1"),
				},
			},
			{
				Attributes: &models.CatalogModelAttributes{
					Name:       apiutils.Of("list-catalog-model-2"),
					ExternalID: apiutils.Of("list-catalog-ext-2"),
				},
			},
		}

		// Save all test models
		var savedModels []models.CatalogModel
		for _, model := range testModels {
			saved, err := repo.Save(model)
			require.NoError(t, err)
			savedModels = append(savedModels, saved)
		}

		// Test listing all models
		listOptions := models.CatalogModelListOptions{}
		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // At least our 2 test models

		// Test filtering by name
		nameFilter := "list-catalog-model-1"
		listOptions = models.CatalogModelListOptions{
			Name: &nameFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, "list-catalog-model-1", *result.Items[0].GetAttributes().Name)

		// Test filtering by external ID
		externalIDFilter := "list-catalog-ext-2"
		listOptions = models.CatalogModelListOptions{
			ExternalID: &externalIDFilter,
		}
		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.Equal(t, "list-catalog-ext-2", *result.Items[0].GetAttributes().ExternalID)
	})

	t.Run("TestGetByName", func(t *testing.T) {
		// First create a model to retrieve by name
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("get-by-name-test-model"),
				ExternalID: apiutils.Of("get-by-name-ext-123"),
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by name
		retrieved, err := repo.GetByName("get-by-name-test-model")
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-by-name-test-model", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-by-name-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent name
		_, err = repo.GetByName("non-existent-model")
		assert.ErrorIs(t, err, service.ErrCatalogModelNotFound)
	})

	t.Run("TestUpdateWithID", func(t *testing.T) {
		// First create a model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-test-model"),
				ExternalID: apiutils.Of("update-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "environment",
					StringValue: apiutils.Of("dev"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Update the model with ID specified
		updateModel := &models.CatalogModelImpl{
			ID: saved.GetID(), // Specify the ID for update
			Attributes: &models.CatalogModelAttributes{
				Name:                     apiutils.Of("updated-test-model"),
				ExternalID:               apiutils.Of("updated-ext-456"),
				CreateTimeSinceEpoch:     saved.GetAttributes().CreateTimeSinceEpoch, // Preserve create time
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("2.0.0"), // Updated version
				},
				{
					Name:        "description",
					StringValue: apiutils.Of("Updated description"), // New property
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "environment",
					StringValue: apiutils.Of("prod"), // Updated environment
				},
			},
		}

		updated, err := repo.Save(updateModel)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify the update
		assert.Equal(t, *saved.GetID(), *updated.GetID()) // Same ID
		assert.Equal(t, "updated-test-model", *updated.GetAttributes().Name)
		assert.Equal(t, "updated-ext-456", *updated.GetAttributes().ExternalID)

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		assert.Len(t, *updated.GetProperties(), 2)

		// Verify custom properties were updated
		require.NotNil(t, updated.GetCustomProperties())
		assert.Len(t, *updated.GetCustomProperties(), 1)
	})

	t.Run("TestUpdateWithName", func(t *testing.T) {
		// First create a model
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-by-name-model"),
				ExternalID: apiutils.Of("update-by-name-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: apiutils.Of("draft"),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Update the model without specifying ID (should lookup by name)
		updateModel := &models.CatalogModelImpl{
			// No ID specified - should trigger name lookup in Save method
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("update-by-name-model"), // Same name to trigger lookup
				ExternalID: apiutils.Of("updated-by-name-ext-789"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "status",
					StringValue: apiutils.Of("published"), // Updated status
				},
				{
					Name:        "category",
					StringValue: apiutils.Of("ml-model"), // New property
				},
			},
		}

		updated, err := repo.Save(updateModel)
		require.NoError(t, err)
		require.NotNil(t, updated)

		// Verify the update happened (same ID, updated fields)
		assert.Equal(t, *saved.GetID(), *updated.GetID()) // Should have same ID from lookup
		assert.Equal(t, "update-by-name-model", *updated.GetAttributes().Name)
		assert.Equal(t, "updated-by-name-ext-789", *updated.GetAttributes().ExternalID)

		// Verify properties were updated
		require.NotNil(t, updated.GetProperties())
		assert.Len(t, *updated.GetProperties(), 2)
	})

	t.Run("TestListWithPropertiesAndCustomProperties", func(t *testing.T) {
		// Create a model with both properties and custom properties
		catalogModel := &models.CatalogModelImpl{
			Attributes: &models.CatalogModelAttributes{
				Name:       apiutils.Of("props-test-catalog-model"),
				ExternalID: apiutils.Of("props-catalog-ext-123"),
			},
			Properties: &[]dbmodels.Properties{
				{
					Name:        "version",
					StringValue: apiutils.Of("1.0.0"),
				},
				{
					Name:     "priority",
					IntValue: apiutils.Of(int32(5)),
				},
			},
			CustomProperties: &[]dbmodels.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("ml-team"),
				},
				{
					Name:      "active",
					BoolValue: apiutils.Of(true),
				},
			},
		}

		saved, err := repo.Save(catalogModel)
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
	})
}

// Helper function to get or create CatalogModel type ID
func getCatalogModelTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	if err != nil {
		require.NoError(t, err, "Failed to query CatalogModel type")
	}

	return int64(typeRecord.ID)
}
