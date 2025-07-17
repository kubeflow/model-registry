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

func getDocArtifactTypeID(t *testing.T, db *gorm.DB) int64 {
	var typeRecord schema.Type
	err := db.Where("name = ?", defaults.DocArtifactTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to find DocArtifact type")
	return int64(typeRecord.ID)
}

func TestDocArtifactRepository(t *testing.T) {
	cleanupTestData(t, sharedDB)

	// Get the actual DocArtifact type ID from the database
	typeID := getDocArtifactTypeID(t, sharedDB)
	repo := service.NewDocArtifactRepository(sharedDB, typeID)

	// Also get other type IDs for creating related entities
	registeredModelTypeID := getRegisteredModelTypeID(t, sharedDB)
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, sharedDB)
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, modelVersionTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a registered model and model version for attribution
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-doc-artifact"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-doc-artifact"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// Test creating a new doc artifact
		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("test-doc-artifact"),
				ExternalID:   apiutils.Of("doc-artifact-ext-123"),
				URI:          apiutils.Of("s3://bucket/documentation.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test doc artifact description"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-doc-prop",
					StringValue:      apiutils.Of("custom-doc-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-doc-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, "doc-artifact-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/documentation.pdf", *saved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *saved.GetAttributes().State)

		// Test updating the same doc artifact
		docArtifact.ID = saved.GetID()
		docArtifact.GetAttributes().Name = apiutils.Of("updated-doc-artifact")
		docArtifact.GetAttributes().State = apiutils.Of("PENDING")

		updated, err := repo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-doc-artifact", *updated.GetAttributes().Name)
		assert.Equal(t, "PENDING", *updated.GetAttributes().State)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-getbyid"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-getbyid"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// First create a doc artifact to retrieve
		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("get-test-doc-artifact"),
				ExternalID:   apiutils.Of("get-doc-artifact-ext-123"),
				URI:          apiutils.Of("s3://bucket/get-documentation.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}

		saved, err := repo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-doc-artifact", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-doc-artifact-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/get-documentation.pdf", *retrieved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *retrieved.GetAttributes().State)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a registered model and model version for the artifacts
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-list"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-list"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// Create multiple doc artifacts for listing
		testArtifacts := []*models.DocArtifactImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.DocArtifactAttributes{
					Name:         apiutils.Of("list-doc-artifact-1"),
					ExternalID:   apiutils.Of("list-doc-artifact-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-doc-1.pdf"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("doc-artifact"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.DocArtifactAttributes{
					Name:         apiutils.Of("list-doc-artifact-2"),
					ExternalID:   apiutils.Of("list-doc-artifact-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-doc-2.pdf"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("doc-artifact"),
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.DocArtifactAttributes{
					Name:         apiutils.Of("list-doc-artifact-3"),
					ExternalID:   apiutils.Of("list-doc-artifact-ext-3"),
					URI:          apiutils.Of("s3://bucket/list-doc-3.pdf"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("doc-artifact"),
				},
			},
		}

		for _, artifact := range testArtifacts {
			_, err := repo.Save(artifact, savedModelVersion.GetID())
			require.NoError(t, err)
		}

		// Test listing all artifacts with basic pagination
		pageSize := int32(10)
		listOptions := models.DocArtifactListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test artifacts

		// Test listing by name
		listOptions = models.DocArtifactListOptions{
			Name: apiutils.Of("list-doc-artifact-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-doc-artifact-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.DocArtifactListOptions{
			ExternalID: apiutils.Of("list-doc-artifact-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-doc-artifact-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test listing by model version ID
		listOptions = models.DocArtifactListOptions{
			ModelVersionID: savedModelVersion.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should find our 3 test artifacts

		// Test ordering by ID (deterministic)
		listOptions = models.DocArtifactListOptions{
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
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-ordering"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-ordering"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// Create artifacts sequentially with time delays to ensure deterministic ordering
		artifact1 := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("time-test-doc-artifact-1"),
				URI:          apiutils.Of("s3://bucket/time-doc-1.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		saved1, err := repo.Save(artifact1, savedModelVersion.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		artifact2 := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("time-test-doc-artifact-2"),
				URI:          apiutils.Of("s3://bucket/time-doc-2.pdf"),
				State:        apiutils.Of("PENDING"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		saved2, err := repo.Save(artifact2, savedModelVersion.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.DocArtifactListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test artifacts in the results
		var foundArtifact1, foundArtifact2 models.DocArtifact
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

	t.Run("TestSaveWithoutModelVersion", func(t *testing.T) {
		// Test creating a doc artifact without model version attribution
		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("standalone-doc-artifact"),
				URI:          apiutils.Of("s3://bucket/standalone-doc.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Standalone doc artifact without model version"),
				},
			},
		}

		saved, err := repo.Save(docArtifact, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "standalone-doc-artifact", *saved.GetAttributes().Name)
		assert.Equal(t, "s3://bucket/standalone-doc.pdf", *saved.GetAttributes().URI)

		// Verify it can be retrieved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		assert.Equal(t, "standalone-doc-artifact", *retrieved.GetAttributes().Name)
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-props"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-props"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("props-test-doc-artifact"),
				URI:          apiutils.Of("s3://bucket/props-doc.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Doc artifact with properties"),
				},
				{
					Name:        "document_type",
					StringValue: apiutils.Of("user_manual"),
				},
				{
					Name:     "page_count",
					IntValue: apiutils.Of(int32(42)),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "team",
					StringValue:      apiutils.Of("ml-team"),
					IsCustomProperty: true,
				},
				{
					Name:             "priority",
					IntValue:         apiutils.Of(int32(5)),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, document_type, page_count

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify specific properties exist
		properties := *retrieved.GetProperties()
		var foundDescription, foundDocumentType, foundPageCount bool
		for _, prop := range properties {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Doc artifact with properties", *prop.StringValue)
			case "document_type":
				foundDocumentType = true
				assert.Equal(t, "user_manual", *prop.StringValue)
			case "page_count":
				foundPageCount = true
				assert.Equal(t, int32(42), *prop.IntValue)
			}
		}
		assert.True(t, foundDescription, "Should find description property")
		assert.True(t, foundDocumentType, "Should find document_type property")
		assert.True(t, foundPageCount, "Should find page_count property")
	})
}
