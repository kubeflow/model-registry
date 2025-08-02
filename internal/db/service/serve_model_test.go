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

func TestServeModelRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t)
	defer cleanup()

	// Get the actual ServeModel type ID from the database
	typeID := getServeModelTypeID(t, sharedDB)
	repo := service.NewServeModelRepository(sharedDB, typeID)

	// Also get other type IDs for creating related entities
	registeredModelTypeID := getRegisteredModelTypeID(t, sharedDB)
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, sharedDB)
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, modelVersionTypeID)

	inferenceServiceTypeID := getInferenceServiceTypeID(t, sharedDB)
	inferenceServiceRepo := service.NewInferenceServiceRepository(sharedDB, inferenceServiceTypeID)

	servingEnvironmentTypeID := getServingEnvironmentTypeID(t, sharedDB)
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(sharedDB, servingEnvironmentTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-serve"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// Create a model version
		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-serve"),
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

		// Create an inference service for the serve model
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-save"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-save"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// Test creating a new serve model
		serveModel := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:test-serve-model", *savedInferenceService.GetID())),
				ExternalID:     apiutils.Of("serve-ext-123"),
				LastKnownState: apiutils.Of("RUNNING"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test serve model description"),
				},
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-serve-prop",
					StringValue:      apiutils.Of("custom-serve-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(serveModel, savedInferenceService.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:test-serve-model", *savedInferenceService.GetID()), *saved.GetAttributes().Name)
		assert.Equal(t, "serve-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "RUNNING", *saved.GetAttributes().LastKnownState)

		// Test updating the same serve model
		serveModel.ID = saved.GetID()
		serveModel.GetAttributes().Name = apiutils.Of(fmt.Sprintf("%d:updated-serve-model", *savedInferenceService.GetID()))
		serveModel.GetAttributes().LastKnownState = apiutils.Of("COMPLETE")

		updated, err := repo.Save(serveModel, nil)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, fmt.Sprintf("%d:updated-serve-model", *savedInferenceService.GetID()), *updated.GetAttributes().Name)
		assert.Equal(t, "COMPLETE", *updated.GetAttributes().LastKnownState)
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

		// Create an inference service for the serve model
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-getbyid"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-getbyid"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// First create a serve model to retrieve
		serveModel := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:get-test-serve-model", *savedInferenceService.GetID())),
				ExternalID:     apiutils.Of("get-serve-ext-123"),
				LastKnownState: apiutils.Of("NEW"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}

		saved, err := repo.Save(serveModel, savedInferenceService.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:get-test-serve-model", *savedInferenceService.GetID()), *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-serve-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "NEW", *retrieved.GetAttributes().LastKnownState)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a registered model and model version for the serve models
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

		// Create an inference service for the serve models
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-list"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-list"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// Create multiple serve models for listing
		testServeModels := []*models.ServeModelImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServeModelAttributes{
					Name:           apiutils.Of(fmt.Sprintf("%d:list-serve-model-1", *savedInferenceService.GetID())),
					ExternalID:     apiutils.Of("list-serve-ext-1"),
					LastKnownState: apiutils.Of("RUNNING"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "model_version_id",
						IntValue: savedModelVersion.GetID(),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServeModelAttributes{
					Name:           apiutils.Of(fmt.Sprintf("%d:list-serve-model-2", *savedInferenceService.GetID())),
					ExternalID:     apiutils.Of("list-serve-ext-2"),
					LastKnownState: apiutils.Of("COMPLETE"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "model_version_id",
						IntValue: savedModelVersion.GetID(),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServeModelAttributes{
					Name:           apiutils.Of(fmt.Sprintf("%d:list-serve-model-3", *savedInferenceService.GetID())),
					ExternalID:     apiutils.Of("list-serve-ext-3"),
					LastKnownState: apiutils.Of("FAILED"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "model_version_id",
						IntValue: savedModelVersion.GetID(),
					},
				},
			},
		}

		for _, srvModel := range testServeModels {
			_, err := repo.Save(srvModel, savedInferenceService.GetID())
			require.NoError(t, err)
		}

		// Test listing all serve models with basic pagination
		pageSize := int32(10)
		listOptions := models.ServeModelListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test serve models

		// Test listing by name
		listOptions = models.ServeModelListOptions{
			Name: apiutils.Of("list-serve-model-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:list-serve-model-1", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.ServeModelListOptions{
			ExternalID: apiutils.Of("list-serve-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-serve-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.ServeModelListOptions{
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

		// Create an inference service for the serve models
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-ordering"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-ordering"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// Create serve models sequentially with time delays to ensure deterministic ordering
		serveModel1 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:time-test-serve-model-1", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("RUNNING"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}
		saved1, err := repo.Save(serveModel1, savedInferenceService.GetID())
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		serveModel2 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:time-test-serve-model-2", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("COMPLETE"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}
		saved2, err := repo.Save(serveModel2, savedInferenceService.GetID())
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.ServeModelListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test serve models in the results
		var foundServeModel1, foundServeModel2 models.ServeModel
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundServeModel1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundServeModel2 = item
				index2 = i
			}
		}

		// Verify both serve models were found and serveModel1 comes before serveModel2 (ascending order)
		require.NotEqual(t, -1, index1, "Serve Model 1 should be found in results")
		require.NotEqual(t, -1, index2, "Serve Model 2 should be found in results")
		assert.Less(t, index1, index2, "Serve Model 1 should come before Serve Model 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundServeModel1.GetAttributes().CreateTimeSinceEpoch, *foundServeModel2.GetAttributes().CreateTimeSinceEpoch, "Serve Model 1 should have earlier create time")
	})

	t.Run("TestListByInferenceService", func(t *testing.T) {
		// First create a serving environment
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-inference"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		// Create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-inference"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-inference"),
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

		// Create an inference service
		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-serve"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// Create another inference service for comparison
		inferenceService2 := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-2-for-serve"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService2, err := inferenceServiceRepo.Save(inferenceService2)
		require.NoError(t, err)

		// Create serve models - some associated with the first inference service, some with the second
		serveModel1 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:serve-model-with-inference-1", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("RUNNING"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}
		_, err = repo.Save(serveModel1, savedInferenceService.GetID())
		require.NoError(t, err)

		serveModel2 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:serve-model-with-inference-2", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("COMPLETE"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}
		_, err = repo.Save(serveModel2, savedInferenceService.GetID())
		require.NoError(t, err)

		// Create a serve model with the second inference service
		serveModel3 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:serve-model-with-inference-3", *savedInferenceService2.GetID())),
				LastKnownState: apiutils.Of("NEW"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
		}
		_, err = repo.Save(serveModel3, savedInferenceService2.GetID())
		require.NoError(t, err)

		// Test listing serve models by first inference service ID
		pageSize := int32(10)
		listOptions := models.ServeModelListOptions{
			InferenceServiceID: savedInferenceService.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items)) // Should find exactly 2 serve models associated with the first inference service

		// Verify the correct serve models are returned
		foundNames := make(map[string]bool)
		for _, item := range result.Items {
			foundNames[*item.GetAttributes().Name] = true
		}
		assert.True(t, foundNames[fmt.Sprintf("%d:serve-model-with-inference-1", *savedInferenceService.GetID())])
		assert.True(t, foundNames[fmt.Sprintf("%d:serve-model-with-inference-2", *savedInferenceService.GetID())])
		assert.False(t, foundNames[fmt.Sprintf("%d:serve-model-with-inference-3", *savedInferenceService2.GetID())])

		// Test listing serve models by second inference service ID
		listOptions = models.ServeModelListOptions{
			InferenceServiceID: savedInferenceService2.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items)) // Should find exactly 1 serve model associated with the second inference service

		// Verify the correct serve model is returned
		assert.Equal(t, fmt.Sprintf("%d:serve-model-with-inference-3", *savedInferenceService2.GetID()), *result.Items[0].GetAttributes().Name)
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

		// Create an inference service for the serve model
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-props"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-props"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		serveModel := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:props-test-serve-model", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("RUNNING"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Serve model with properties"),
				},
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
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

		saved, err := repo.Save(serveModel, savedInferenceService.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 2) // description, model_version_id

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify the required model_version_id property exists
		properties := *retrieved.GetProperties()
		var foundModelVersionID bool
		for _, prop := range properties {
			if prop.Name == "model_version_id" && prop.IntValue != nil && *prop.IntValue == *savedModelVersion.GetID() {
				foundModelVersionID = true
				break
			}
		}
		assert.True(t, foundModelVersionID, "Should find model_version_id property")
	})

	t.Run("TestFilterQuery", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-filter"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-filter"),
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

		// Create an inference service for the serve models
		servingEnvironment := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("test-serving-env-for-filter"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(servingEnvironment)
		require.NoError(t, err)

		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(inferenceServiceTypeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("test-inference-service-for-filter"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedInferenceService, err := inferenceServiceRepo.Save(inferenceService)
		require.NoError(t, err)

		// Create multiple serve models with different properties for filtering
		serveModel1 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("RUNNING"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "framework",
					StringValue:      apiutils.Of("pytorch"),
					IsCustomProperty: true,
				},
				{
					Name:             "replicas",
					IntValue:         apiutils.Of(int32(3)),
					IsCustomProperty: true,
				},
				{
					Name:             "accuracy",
					DoubleValue:      apiutils.Of(0.95),
					IsCustomProperty: true,
				},
			},
		}
		_, err = repo.Save(serveModel1, savedInferenceService.GetID())
		require.NoError(t, err)

		serveModel2 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:tensorflow-serve-model", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("COMPLETE"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "framework",
					StringValue:      apiutils.Of("tensorflow"),
					IsCustomProperty: true,
				},
				{
					Name:             "replicas",
					IntValue:         apiutils.Of(int32(5)),
					IsCustomProperty: true,
				},
				{
					Name:             "accuracy",
					DoubleValue:      apiutils.Of(0.89),
					IsCustomProperty: true,
				},
			},
		}
		_, err = repo.Save(serveModel2, savedInferenceService.GetID())
		require.NoError(t, err)

		serveModel3 := &models.ServeModelImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.ServeModelAttributes{
				Name:           apiutils.Of(fmt.Sprintf("%d:sklearn-serve-model", *savedInferenceService.GetID())),
				LastKnownState: apiutils.Of("NEW"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "framework",
					StringValue:      apiutils.Of("sklearn"),
					IsCustomProperty: true,
				},
				{
					Name:             "replicas",
					IntValue:         apiutils.Of(int32(2)),
					IsCustomProperty: true,
				},
				{
					Name:             "accuracy",
					DoubleValue:      apiutils.Of(0.92),
					IsCustomProperty: true,
				},
			},
		}
		_, err = repo.Save(serveModel3, savedInferenceService.GetID())
		require.NoError(t, err)

		// Test core property filtering
		t.Run("CorePropertyFilter", func(t *testing.T) {
			filterQuery := `name = "pytorch-serve-model"`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
		})

		// Test custom property filtering
		t.Run("CustomPropertyFilter", func(t *testing.T) {
			filterQuery := `framework = "tensorflow"`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:tensorflow-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
		})

		// Test numeric custom property filtering
		t.Run("NumericCustomPropertyFilter", func(t *testing.T) {
			filterQuery := `replicas >= 3`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 2, len(result.Items)) // pytorch (3) and tensorflow (5)
		})

		// Test double custom property filtering
		t.Run("DoubleCustomPropertyFilter", func(t *testing.T) {
			filterQuery := `accuracy > 0.9`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 2, len(result.Items)) // pytorch (0.95) and sklearn (0.92)
		})

		// Test complex AND filter
		t.Run("ComplexANDFilter", func(t *testing.T) {
			filterQuery := `framework = "pytorch" AND replicas = 3`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
		})

		// Test complex OR filter
		t.Run("ComplexORFilter", func(t *testing.T) {
			filterQuery := `framework = "pytorch" OR framework = "sklearn"`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 2, len(result.Items)) // pytorch and sklearn
		})

		// Test ILIKE operator
		t.Run("ILIKEFilter", func(t *testing.T) {
			filterQuery := `name ILIKE "%Serve-Model%"`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.GreaterOrEqual(t, len(result.Items), 3) // All serve models contain "serve-model"
		})

		// Test mixed core and custom property filter
		t.Run("MixedCoreAndCustomFilter", func(t *testing.T) {
			filterQuery := `name ILIKE "%pytorch%" AND accuracy > 0.9`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
		})

		// Test invalid filter query
		t.Run("InvalidFilterQuery", func(t *testing.T) {
			filterQuery := `invalid syntax =`
			pageSize := int32(10)
			listOptions := models.ServeModelListOptions{
				Pagination: models.Pagination{
					PageSize:    &pageSize,
					FilterQuery: &filterQuery,
				},
			}

			result, err := repo.List(listOptions)
			require.Error(t, err)
			require.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid filter query")
		})

		// Test combining old parameters with new filterQuery
		t.Run("CombinedOldAndNewFilters", func(t *testing.T) {
			// Setup additional test data with ExternalID
			serveModelWithExternalID := &models.ServeModelImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.ServeModelAttributes{
					Name:           apiutils.Of(fmt.Sprintf("%d:pytorch-serve-model-with-external-id", *savedInferenceService.GetID())),
					ExternalID:     apiutils.Of("ext-pytorch-123"),
					LastKnownState: apiutils.Of("RUNNING"),
				},
				CustomProperties: &[]models.Properties{
					{
						Name:             "framework",
						StringValue:      apiutils.Of("pytorch"),
						IsCustomProperty: true,
					},
					{
						Name:             "replicas",
						IntValue:         apiutils.Of(int32(3)),
						IsCustomProperty: true,
					},
				},
			}
			_, err := repo.Save(serveModelWithExternalID, savedInferenceService.GetID())
			require.NoError(t, err)

			// Test old Name parameter alone
			t.Run("OldNameParameterAlone", func(t *testing.T) {
				name := "pytorch-serve-model"
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					Name: &name,
					Pagination: models.Pagination{
						PageSize: &pageSize,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 1, len(result.Items))
				assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
			})

			// Test old Name parameter combined with filterQuery (should be AND)
			t.Run("OldNameParameterCombinedWithFilterQuery", func(t *testing.T) {
				name := "pytorch-serve-model"
				filterQuery := `framework = "pytorch"`
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					Name: &name,
					Pagination: models.Pagination{
						PageSize:    &pageSize,
						FilterQuery: &filterQuery,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 1, len(result.Items))
				assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
			})

			// Test old Name parameter combined with filterQuery (should return 0 results)
			t.Run("OldNameParameterCombinedWithFilterQueryNoMatch", func(t *testing.T) {
				name := "pytorch-serve-model"
				filterQuery := `framework = "tensorflow"` // This model has framework=pytorch, so should return 0
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					Name: &name,
					Pagination: models.Pagination{
						PageSize:    &pageSize,
						FilterQuery: &filterQuery,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 0, len(result.Items))
			})

			// Test old ExternalID parameter alone
			t.Run("OldExternalIDParameterAlone", func(t *testing.T) {
				externalID := "ext-pytorch-123"
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					ExternalID: &externalID,
					Pagination: models.Pagination{
						PageSize: &pageSize,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 1, len(result.Items))
				assert.Equal(t, "ext-pytorch-123", *result.Items[0].GetAttributes().ExternalID)
			})

			// Test old ExternalID parameter combined with filterQuery
			t.Run("OldExternalIDParameterCombinedWithFilterQuery", func(t *testing.T) {
				externalID := "ext-pytorch-123"
				filterQuery := `replicas = 3` // The model with this ExternalID has replicas=3
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					ExternalID: &externalID,
					Pagination: models.Pagination{
						PageSize:    &pageSize,
						FilterQuery: &filterQuery,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.Equal(t, 1, len(result.Items))
				assert.Equal(t, fmt.Sprintf("%d:pytorch-serve-model-with-external-id", *savedInferenceService.GetID()), *result.Items[0].GetAttributes().Name)
			})

			// Test old InferenceServiceID parameter combined with filterQuery
			t.Run("OldInferenceServiceIDParameterCombinedWithFilterQuery", func(t *testing.T) {
				inferenceServiceID := savedInferenceService.GetID()
				filterQuery := `framework = "pytorch"` // Should match existing pytorch models + new one
				pageSize := int32(10)
				listOptions := models.ServeModelListOptions{
					InferenceServiceID: inferenceServiceID,
					Pagination: models.Pagination{
						PageSize:    &pageSize,
						FilterQuery: &filterQuery,
					},
				}

				result, err := repo.List(listOptions)
				require.NoError(t, err)
				require.NotNil(t, result)
				assert.GreaterOrEqual(t, len(result.Items), 1) // At least one pytorch model should match

				// Check that all returned models have framework=pytorch
				for _, item := range result.Items {
					// Find the framework property
					found := false
					if item.GetCustomProperties() != nil {
						for _, prop := range *item.GetCustomProperties() {
							if prop.Name == "framework" && prop.StringValue != nil {
								assert.Equal(t, "pytorch", *prop.StringValue)
								found = true
								break
							}
						}
					}
					assert.True(t, found, "Framework property not found for model %s", *item.GetAttributes().Name)
				}
			})
		})
	})
}
