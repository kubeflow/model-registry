package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestArtifactRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t)
	defer cleanup()

	// Get the actual type IDs from the database
	modelArtifactTypeID := getModelArtifactTypeID(t, sharedDB)
	docArtifactTypeID := getDocArtifactTypeID(t, sharedDB)
	dataSetTypeID := getDataSetTypeID(t, sharedDB)
	metricTypeID := getMetricTypeID(t, sharedDB)
	parameterTypeID := getParameterTypeID(t, sharedDB)
	metricHistoryTypeID := getMetricHistoryTypeID(t, sharedDB)
	repo := service.NewArtifactRepository(sharedDB, modelArtifactTypeID, docArtifactTypeID, dataSetTypeID, metricTypeID, parameterTypeID, metricHistoryTypeID)

	// Also get other type IDs for creating related entities
	registeredModelTypeID := getRegisteredModelTypeID(t, sharedDB)
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, sharedDB)
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, modelVersionTypeID)

	// Create shared test data
	registeredModel := &models.RegisteredModelImpl{
		TypeID: apiutils.Of(int32(registeredModelTypeID)),
		Attributes: &models.RegisteredModelAttributes{
			Name: apiutils.Of("test-registered-model-for-artifacts"),
		},
	}
	savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
	require.NoError(t, err)

	modelVersion := &models.ModelVersionImpl{
		TypeID: apiutils.Of(int32(modelVersionTypeID)),
		Attributes: &models.ModelVersionAttributes{
			Name: apiutils.Of("test-model-version-for-artifacts"),
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

	t.Run("TestGetByID", func(t *testing.T) {
		// Create a model artifact using the model artifact repository
		modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, modelArtifactTypeID)
		modelArtifact := &models.ModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.ModelArtifactAttributes{
				Name:         apiutils.Of(fmt.Sprintf("%d:test-model-artifact-for-getbyid", *savedModelVersion.GetID())),
				URI:          apiutils.Of("s3://bucket/model.pkl"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("model-artifact"),
			},
		}
		savedModelArtifact, err := modelArtifactRepo.Save(modelArtifact, savedModelVersion.GetID())
		require.NoError(t, err)

		// Create a doc artifact using the doc artifact repository
		docArtifactRepo := service.NewDocArtifactRepository(sharedDB, docArtifactTypeID)
		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(docArtifactTypeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of(fmt.Sprintf("%d:unified-test-doc-artifact-for-getbyid", *savedModelVersion.GetID())),
				URI:          apiutils.Of("s3://bucket/doc.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		savedDocArtifact, err := docArtifactRepo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)

		// Test retrieving model artifact by ID
		retrievedModelArtifact, err := repo.GetByID(*savedModelArtifact.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrievedModelArtifact.ModelArtifact)
		require.Nil(t, retrievedModelArtifact.DocArtifact)
		assert.Equal(t, *savedModelArtifact.GetID(), *apiutils.ZeroIfNil(retrievedModelArtifact.ModelArtifact).GetID())
		assert.Equal(t, fmt.Sprintf("%d:test-model-artifact-for-getbyid", *savedModelVersion.GetID()), *apiutils.ZeroIfNil(retrievedModelArtifact.ModelArtifact).GetAttributes().Name)
		assert.Equal(t, "s3://bucket/model.pkl", *apiutils.ZeroIfNil(retrievedModelArtifact.ModelArtifact).GetAttributes().URI)

		// Test retrieving doc artifact by ID
		retrievedDocArtifact, err := repo.GetByID(*savedDocArtifact.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrievedDocArtifact.DocArtifact)
		require.Nil(t, retrievedDocArtifact.ModelArtifact)
		assert.Equal(t, *savedDocArtifact.GetID(), *apiutils.ZeroIfNil(retrievedDocArtifact.DocArtifact).GetID())
		assert.Equal(t, fmt.Sprintf("%d:unified-test-doc-artifact-for-getbyid", *savedModelVersion.GetID()), *apiutils.ZeroIfNil(retrievedDocArtifact.DocArtifact).GetAttributes().Name)
		assert.Equal(t, "s3://bucket/doc.pdf", *apiutils.ZeroIfNil(retrievedDocArtifact.DocArtifact).GetAttributes().URI)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple artifacts of both types using their respective repositories
		modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, modelArtifactTypeID)
		docArtifactRepo := service.NewDocArtifactRepository(sharedDB, docArtifactTypeID)

		// Create model artifacts
		modelArtifacts := []*models.ModelArtifactImpl{
			{
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.ModelArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("%d:list-model-artifact-1", *savedModelVersion.GetID())),
					ExternalID:   apiutils.Of("list-model-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-model-1.pkl"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("model-artifact"),
				},
			},
			{
				TypeID: apiutils.Of(int32(modelArtifactTypeID)),
				Attributes: &models.ModelArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("%d:list-model-artifact-2", *savedModelVersion.GetID())),
					ExternalID:   apiutils.Of("list-model-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-model-2.pkl"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("model-artifact"),
				},
			},
		}

		for _, artifact := range modelArtifacts {
			_, err := modelArtifactRepo.Save(artifact, savedModelVersion.GetID())
			require.NoError(t, err)
		}

		// Create doc artifacts
		docArtifacts := []*models.DocArtifactImpl{
			{
				TypeID: apiutils.Of(int32(docArtifactTypeID)),
				Attributes: &models.DocArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("%d:unified-list-doc-artifact-1", *savedModelVersion.GetID())),
					ExternalID:   apiutils.Of("unified-list-doc-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-doc-1.pdf"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("doc-artifact"),
				},
			},
			{
				TypeID: apiutils.Of(int32(docArtifactTypeID)),
				Attributes: &models.DocArtifactAttributes{
					Name:         apiutils.Of(fmt.Sprintf("%d:unified-list-doc-artifact-2", *savedModelVersion.GetID())),
					ExternalID:   apiutils.Of("unified-list-doc-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-doc-2.pdf"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("doc-artifact"),
				},
			},
		}

		for _, artifact := range docArtifacts {
			_, err := docArtifactRepo.Save(artifact, savedModelVersion.GetID())
			require.NoError(t, err)
		}

		// Test listing all artifacts with basic pagination
		pageSize := int32(10)
		listOptions := models.ArtifactListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 4) // At least our 4 test artifacts

		// Verify we get both types of artifacts
		var foundModelArtifacts, foundDocArtifacts int
		for _, item := range result.Items {
			if item.ModelArtifact != nil {
				foundModelArtifacts++
			}
			if item.DocArtifact != nil {
				foundDocArtifacts++
			}
		}
		assert.GreaterOrEqual(t, foundModelArtifacts, 2, "Should find at least 2 model artifacts")
		assert.GreaterOrEqual(t, foundDocArtifacts, 2, "Should find at least 2 doc artifacts")

		// Test listing by name (model artifact)
		listOptions = models.ArtifactListOptions{
			Name: apiutils.Of("list-model-artifact-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.NotNil(t, result.Items[0].ModelArtifact)
		assert.Equal(t, fmt.Sprintf("%d:list-model-artifact-1", *savedModelVersion.GetID()), *apiutils.ZeroIfNil(result.Items[0].ModelArtifact).GetAttributes().Name)

		// Test listing by name (doc artifact)
		listOptions = models.ArtifactListOptions{
			Name: apiutils.Of("unified-list-doc-artifact-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.NotNil(t, result.Items[0].DocArtifact)
		assert.Equal(t, fmt.Sprintf("%d:unified-list-doc-artifact-1", *savedModelVersion.GetID()), *apiutils.ZeroIfNil(result.Items[0].DocArtifact).GetAttributes().Name)

		// Test listing by external ID
		listOptions = models.ArtifactListOptions{
			ExternalID: apiutils.Of("list-model-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items))
		assert.NotNil(t, result.Items[0].ModelArtifact)
		assert.Equal(t, "list-model-ext-2", *apiutils.ZeroIfNil(result.Items[0].ModelArtifact).GetAttributes().ExternalID)

		// Test listing by model version ID
		listOptions = models.ArtifactListOptions{
			ParentResourceID: savedModelVersion.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 4) // Should find our 4 test artifacts

		// Test ordering by ID (deterministic)
		listOptions = models.ArtifactListOptions{
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
			var firstID, secondID int32
			if result.Items[0].ModelArtifact != nil {
				firstID = *apiutils.ZeroIfNil(result.Items[0].ModelArtifact).GetID()
			} else {
				firstID = *apiutils.ZeroIfNil(result.Items[0].DocArtifact).GetID()
			}
			if result.Items[1].ModelArtifact != nil {
				secondID = *apiutils.ZeroIfNil(result.Items[1].ModelArtifact).GetID()
			} else {
				secondID = *apiutils.ZeroIfNil(result.Items[1].DocArtifact).GetID()
			}
			assert.Less(t, firstID, secondID, "Results should be ordered by ID ascending")
		}
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// Create artifacts sequentially with time delays to ensure deterministic ordering
		modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, modelArtifactTypeID)
		docArtifactRepo := service.NewDocArtifactRepository(sharedDB, docArtifactTypeID)

		// Create first artifact (model artifact)
		artifact1 := &models.ModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.ModelArtifactAttributes{
				Name:         apiutils.Of("time-test-model-artifact-1"),
				URI:          apiutils.Of("s3://bucket/time-model-1.pkl"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("model-artifact"),
			},
		}
		saved1, err := modelArtifactRepo.Save(artifact1, savedModelVersion.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		// Create second artifact (doc artifact)
		artifact2 := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(docArtifactTypeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("unified-time-test-doc-artifact-2"),
				URI:          apiutils.Of("s3://bucket/time-doc-2.pdf"),
				State:        apiutils.Of("PENDING"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		saved2, err := docArtifactRepo.Save(artifact2, savedModelVersion.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ArtifactListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test artifacts in the results
		var foundArtifact1, foundArtifact2 models.Artifact
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if item.ModelArtifact != nil && *apiutils.ZeroIfNil(item.ModelArtifact).GetID() == *saved1.GetID() {
				foundArtifact1 = item
				index1 = i
			}
			if item.DocArtifact != nil && *apiutils.ZeroIfNil(item.DocArtifact).GetID() == *saved2.GetID() {
				foundArtifact2 = item
				index2 = i
			}
		}

		// Verify both artifacts were found and artifact1 comes before artifact2 (ascending order)
		require.NotEqual(t, -1, index1, "Artifact 1 should be found in results")
		require.NotEqual(t, -1, index2, "Artifact 2 should be found in results")
		assert.Less(t, index1, index2, "Artifact 1 should come before Artifact 2 when ordered by CREATE_TIME")

		// Verify timestamps
		var createTime1, createTime2 int64
		if foundArtifact1.ModelArtifact != nil {
			createTime1 = *apiutils.ZeroIfNil(foundArtifact1.ModelArtifact).GetAttributes().CreateTimeSinceEpoch
		} else {
			createTime1 = *apiutils.ZeroIfNil(foundArtifact1.DocArtifact).GetAttributes().CreateTimeSinceEpoch
		}
		if foundArtifact2.ModelArtifact != nil {
			createTime2 = *apiutils.ZeroIfNil(foundArtifact2.ModelArtifact).GetAttributes().CreateTimeSinceEpoch
		} else {
			createTime2 = *apiutils.ZeroIfNil(foundArtifact2.DocArtifact).GetAttributes().CreateTimeSinceEpoch
		}
		assert.Less(t, createTime1, createTime2, "Artifact 1 should have earlier create time")
	})

	t.Run("TestListMixedTypes", func(t *testing.T) {
		// Test that the unified repository correctly handles mixed artifact types
		modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, modelArtifactTypeID)
		docArtifactRepo := service.NewDocArtifactRepository(sharedDB, docArtifactTypeID)

		// Create artifacts with similar names but different types
		modelArtifact := &models.ModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.ModelArtifactAttributes{
				Name:         apiutils.Of("mixed-test-artifact"),
				URI:          apiutils.Of("s3://bucket/mixed-model.pkl"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("model-artifact"),
			},
		}
		savedModelArtifact, err := modelArtifactRepo.Save(modelArtifact, savedModelVersion.GetID())
		require.NoError(t, err)

		docArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(docArtifactTypeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("unified-mixed-test-doc"),
				URI:          apiutils.Of("s3://bucket/mixed-doc.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		savedDocArtifact, err := docArtifactRepo.Save(docArtifact, savedModelVersion.GetID())
		require.NoError(t, err)

		// Test listing all artifacts
		pageSize := int32(10)
		listOptions := models.ArtifactListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our mixed artifacts
		var foundModelArtifact, foundDocArtifact bool
		for _, item := range result.Items {
			if item.ModelArtifact != nil && *apiutils.ZeroIfNil(item.ModelArtifact).GetID() == *savedModelArtifact.GetID() {
				foundModelArtifact = true
				assert.Equal(t, "mixed-test-artifact", *apiutils.ZeroIfNil(item.ModelArtifact).GetAttributes().Name)
				assert.Equal(t, "s3://bucket/mixed-model.pkl", *apiutils.ZeroIfNil(item.ModelArtifact).GetAttributes().URI)
			}
			if item.DocArtifact != nil && *apiutils.ZeroIfNil(item.DocArtifact).GetID() == *savedDocArtifact.GetID() {
				foundDocArtifact = true
				assert.Equal(t, "unified-mixed-test-doc", *apiutils.ZeroIfNil(item.DocArtifact).GetAttributes().Name)
				assert.Equal(t, "s3://bucket/mixed-doc.pdf", *apiutils.ZeroIfNil(item.DocArtifact).GetAttributes().URI)
			}
		}

		assert.True(t, foundModelArtifact, "Should find the model artifact in mixed results")
		assert.True(t, foundDocArtifact, "Should find the doc artifact in mixed results")
	})

	t.Run("TestListWithoutModelVersion", func(t *testing.T) {
		// Test listing artifacts that are not attributed to any model version
		modelArtifactRepo := service.NewModelArtifactRepository(sharedDB, modelArtifactTypeID)
		docArtifactRepo := service.NewDocArtifactRepository(sharedDB, docArtifactTypeID)

		// Create standalone artifacts (without model version attribution)
		standaloneModelArtifact := &models.ModelArtifactImpl{
			TypeID: apiutils.Of(int32(modelArtifactTypeID)),
			Attributes: &models.ModelArtifactAttributes{
				Name:         apiutils.Of("standalone-model-artifact"),
				URI:          apiutils.Of("s3://bucket/standalone-model.pkl"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("model-artifact"),
			},
		}
		savedStandaloneModel, err := modelArtifactRepo.Save(standaloneModelArtifact, nil)
		require.NoError(t, err)

		standaloneDocArtifact := &models.DocArtifactImpl{
			TypeID: apiutils.Of(int32(docArtifactTypeID)),
			Attributes: &models.DocArtifactAttributes{
				Name:         apiutils.Of("unified-standalone-doc-artifact"),
				URI:          apiutils.Of("s3://bucket/standalone-doc.pdf"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("doc-artifact"),
			},
		}
		savedStandaloneDoc, err := docArtifactRepo.Save(standaloneDocArtifact, nil)
		require.NoError(t, err)

		// Test listing all artifacts (should include standalone ones)
		pageSize := int32(20)
		listOptions := models.ArtifactListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our standalone artifacts
		var foundStandaloneModel, foundStandaloneDoc bool
		for _, item := range result.Items {
			if item.ModelArtifact != nil && *apiutils.ZeroIfNil(item.ModelArtifact).GetID() == *savedStandaloneModel.GetID() {
				foundStandaloneModel = true
				assert.Equal(t, "standalone-model-artifact", *apiutils.ZeroIfNil(item.ModelArtifact).GetAttributes().Name)
			}
			if item.DocArtifact != nil && *apiutils.ZeroIfNil(item.DocArtifact).GetID() == *savedStandaloneDoc.GetID() {
				foundStandaloneDoc = true
				assert.Equal(t, "unified-standalone-doc-artifact", *apiutils.ZeroIfNil(item.DocArtifact).GetAttributes().Name)
			}
		}

		assert.True(t, foundStandaloneModel, "Should find standalone model artifact")
		assert.True(t, foundStandaloneDoc, "Should find standalone doc artifact")

		// Test filtering by model version ID (should NOT include standalone artifacts)
		listOptions = models.ArtifactListOptions{
			ParentResourceID: savedModelVersion.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Verify standalone artifacts are NOT in the filtered results
		foundStandaloneModel = false
		foundStandaloneDoc = false
		for _, item := range result.Items {
			if item.ModelArtifact != nil && *apiutils.ZeroIfNil(item.ModelArtifact).GetID() == *savedStandaloneModel.GetID() {
				foundStandaloneModel = true
			}
			if item.DocArtifact != nil && *apiutils.ZeroIfNil(item.DocArtifact).GetID() == *savedStandaloneDoc.GetID() {
				foundStandaloneDoc = true
			}
		}

		assert.False(t, foundStandaloneModel, "Should NOT find standalone model artifact when filtering by model version")
		assert.False(t, foundStandaloneDoc, "Should NOT find standalone doc artifact when filtering by model version")
	})
}
