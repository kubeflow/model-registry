package service_test

import (
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func getModelVersionTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.ModelVersionTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find ModelVersion type")
	return int64(typeRecord.ID)
}

func TestModelVersionRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual ModelVersion type ID from the database
	typeID := getModelVersionTypeID(t, db)
	repo := service.NewModelVersionRepository(db, typeID)

	// Also get RegisteredModel type ID for creating parent models
	registeredModelTypeID := getRegisteredModelTypeID(t, db)
	registeredModelRepo := service.NewRegisteredModelRepository(db, registeredModelTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a parent registered model
		parentModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("parent-model-for-version"),
			},
		}
		savedParent, err := registeredModelRepo.Save(parentModel)
		require.NoError(t, err)

		// Test creating a new model version
		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ModelVersionAttributes{
				Name:       apiutils.Of("test-version"),
				ExternalID: apiutils.Of("version-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test version description"),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedParent.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "custom-version-prop",
					StringValue: apiutils.Of("custom-version-value"),
				},
			},
		}

		saved, err := repo.Save(modelVersion)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-version", *saved.GetAttributes().Name)
		assert.Equal(t, "version-ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same model version
		modelVersion.ID = saved.GetID()
		modelVersion.GetAttributes().Name = apiutils.Of("updated-version")

		updated, err := repo.Save(modelVersion)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-version", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a parent registered model
		parentModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("parent-model-for-getbyid"),
			},
		}
		savedParent, err := registeredModelRepo.Save(parentModel)
		require.NoError(t, err)

		// First create a model version to retrieve
		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ModelVersionAttributes{
				Name:       apiutils.Of("get-test-version"),
				ExternalID: apiutils.Of("get-version-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedParent.GetID(),
				},
			},
		}

		saved, err := repo.Save(modelVersion)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-version", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-version-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a parent registered model for the versions
		parentModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("parent-model-for-list"),
			},
		}
		savedParent, err := registeredModelRepo.Save(parentModel)
		require.NoError(t, err)

		// Create multiple model versions for listing
		testVersions := []*models.ModelVersionImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ModelVersionAttributes{
					Name:       apiutils.Of("list-version-1"),
					ExternalID: apiutils.Of("list-version-ext-1"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "registered_model_id",
						IntValue: savedParent.GetID(),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ModelVersionAttributes{
					Name:       apiutils.Of("list-version-2"),
					ExternalID: apiutils.Of("list-version-ext-2"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "registered_model_id",
						IntValue: savedParent.GetID(),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ModelVersionAttributes{
					Name:       apiutils.Of("list-version-3"),
					ExternalID: apiutils.Of("list-version-ext-3"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "registered_model_id",
						IntValue: savedParent.GetID(),
					},
				},
			},
		}

		for _, version := range testVersions {
			_, err := repo.Save(version)
			require.NoError(t, err)
		}

		// Test listing all versions with basic pagination
		pageSize := int32(10)
		listOptions := models.ModelVersionListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test versions

		// Test listing by name
		listOptions = models.ModelVersionListOptions{
			Name: apiutils.Of("list-version-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-version-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ModelVersionListOptions{
			ExternalID: apiutils.Of("list-version-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-version-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test listing by parent resource ID
		listOptions = models.ModelVersionListOptions{
			ParentResourceID: savedParent.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Test ordering by ID (deterministic)
		listOptions = models.ModelVersionListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("ID"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		// Verify we get results back and they are ordered by ID
		assert.GreaterOrEqual(t, len(result.Items), 1)
		if len(result.Items) > 1 {
			// Verify ascending ID order
			firstID := *result.Items[0].GetID()
			secondID := *result.Items[1].GetID()
			assert.Less(t, firstID, secondID, "Results should be ordered by ID ascending")
		}
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// First create a parent registered model
		parentModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("parent-model-for-ordering"),
			},
		}
		savedParent, err := registeredModelRepo.Save(parentModel)
		require.NoError(t, err)

		// Create versions sequentially with time delays to ensure deterministic ordering
		version1 := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("time-test-version-1"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedParent.GetID(),
				},
			},
		}
		saved1, err := repo.Save(version1)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		version2 := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("time-test-version-2"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedParent.GetID(),
				},
			},
		}
		saved2, err := repo.Save(version2)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ModelVersionListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test versions in the results
		var foundVersion1, foundVersion2 models.ModelVersion
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundVersion1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundVersion2 = item
				index2 = i
			}
		}

		// Verify both versions were found and version1 comes before version2 (ascending order)
		require.NotEqual(t, -1, index1, "Version 1 should be found in results")
		require.NotEqual(t, -1, index2, "Version 2 should be found in results")
		assert.Less(t, index1, index2, "Version 1 should come before Version 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundVersion1.GetAttributes().CreateTimeSinceEpoch, *foundVersion2.GetAttributes().CreateTimeSinceEpoch, "Version 1 should have earlier create time")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		// First create a parent registered model
		parentModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("parent-model-for-props"),
			},
		}
		savedParent, err := registeredModelRepo.Save(parentModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("props-test-version"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Version with properties"),
				},
				{
					Name:     "version_number",
					IntValue: apiutils.Of(int32(1)),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedParent.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "team",
					StringValue: apiutils.Of("ml-team"),
				},
				{
					Name:     "priority",
					IntValue: apiutils.Of(int32(5)),
				},
			},
		}

		saved, err := repo.Save(modelVersion)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, version_number, registered_model_id

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)
	})
}
