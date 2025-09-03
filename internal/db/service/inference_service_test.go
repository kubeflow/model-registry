package service_test

import (
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInferenceServiceRepository(t *testing.T) {
	sharedDB, cleanup := testutils.SetupMySQLWithMigrations(t)
	defer cleanup()

	// Get the actual InferenceService type ID from the database
	typeID := getInferenceServiceTypeID(t, sharedDB)
	repo := service.NewInferenceServiceRepository(sharedDB, typeID)

	// Also get other type IDs for creating parent and related entities
	servingEnvironmentTypeID := getServingEnvironmentTypeID(t, sharedDB)
	servingEnvironmentRepo := service.NewServingEnvironmentRepository(sharedDB, servingEnvironmentTypeID)

	registeredModelTypeID := getRegisteredModelTypeID(t, sharedDB)
	registeredModelRepo := service.NewRegisteredModelRepository(sharedDB, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, sharedDB)
	modelVersionRepo := service.NewModelVersionRepository(sharedDB, modelVersionTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// First create a parent serving environment
		parentServingEnv := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("parent-serving-env-for-inference"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(parentServingEnv)
		require.NoError(t, err)

		// Create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// Test creating a new inference service
		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name:       apiutils.Of("test-inference-service"),
				ExternalID: apiutils.Of("inference-ext-123"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test inference service description"),
				},
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
				{
					Name:        "runtime",
					StringValue: apiutils.Of("tensorflow"),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:        "custom-inference-prop",
					StringValue: apiutils.Of("custom-inference-value"),
				},
			},
		}

		saved, err := repo.Save(inferenceService)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-inference-service", *saved.GetAttributes().Name)
		assert.Equal(t, "inference-ext-123", *saved.GetAttributes().ExternalID)

		// Test updating the same inference service
		inferenceService.ID = saved.GetID()
		inferenceService.GetAttributes().Name = apiutils.Of("updated-inference-service")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		inferenceService.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(inferenceService)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-inference-service", *updated.GetAttributes().Name)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a parent serving environment
		parentServingEnv := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("parent-serving-env-for-getbyid"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(parentServingEnv)
		require.NoError(t, err)

		// Create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-getbyid"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// First create an inference service to retrieve
		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name:       apiutils.Of("get-test-inference-service"),
				ExternalID: apiutils.Of("get-inference-ext-123"),
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

		saved, err := repo.Save(inferenceService)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-inference-service", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-inference-ext-123", *retrieved.GetAttributes().ExternalID)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
	})

	t.Run("TestList", func(t *testing.T) {
		// Create a parent serving environment for the inference services
		parentServingEnv := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("parent-serving-env-for-list"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(parentServingEnv)
		require.NoError(t, err)

		// Create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-list"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// Create multiple inference services for listing
		testInferenceServices := []*models.InferenceServiceImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.InferenceServiceAttributes{
					Name:       apiutils.Of("list-inference-service-1"),
					ExternalID: apiutils.Of("list-inference-ext-1"),
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
					{
						Name:        "runtime",
						StringValue: apiutils.Of("tensorflow"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.InferenceServiceAttributes{
					Name:       apiutils.Of("list-inference-service-2"),
					ExternalID: apiutils.Of("list-inference-ext-2"),
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
					{
						Name:        "runtime",
						StringValue: apiutils.Of("pytorch"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.InferenceServiceAttributes{
					Name:       apiutils.Of("list-inference-service-3"),
					ExternalID: apiutils.Of("list-inference-ext-3"),
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
					{
						Name:        "runtime",
						StringValue: apiutils.Of("tensorflow"),
					},
				},
			},
		}

		for _, infSvc := range testInferenceServices {
			_, err := repo.Save(infSvc)
			require.NoError(t, err)
		}

		// Test listing all inference services with basic pagination
		pageSize := int32(10)
		listOptions := models.InferenceServiceListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test inference services

		// Test listing by name
		listOptions = models.InferenceServiceListOptions{
			Name: apiutils.Of("list-inference-service-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-inference-service-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.InferenceServiceListOptions{
			ExternalID: apiutils.Of("list-inference-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-inference-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test listing by parent resource ID (serving environment)
		listOptions = models.InferenceServiceListOptions{
			ParentResourceID: savedServingEnv.GetID(),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should find our 3 test inference services

		// Test listing by runtime
		listOptions = models.InferenceServiceListOptions{
			Runtime: apiutils.Of("tensorflow"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should find 2 tensorflow services

		// Test ordering by ID (deterministic)
		listOptions = models.InferenceServiceListOptions{
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
		// First create a parent serving environment
		parentServingEnv := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("parent-serving-env-for-ordering"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(parentServingEnv)
		require.NoError(t, err)

		// Create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-ordering"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// Create inference services sequentially with time delays to ensure deterministic ordering
		infSvc1 := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("time-test-inference-service-1"),
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
		saved1, err := repo.Save(infSvc1)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		infSvc2 := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("time-test-inference-service-2"),
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
		saved2, err := repo.Save(infSvc2)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.InferenceServiceListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test inference services in the results
		var foundInfSvc1, foundInfSvc2 models.InferenceService
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundInfSvc1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundInfSvc2 = item
				index2 = i
			}
		}

		// Verify both inference services were found and infSvc1 comes before infSvc2 (ascending order)
		require.NotEqual(t, -1, index1, "Inference Service 1 should be found in results")
		require.NotEqual(t, -1, index2, "Inference Service 2 should be found in results")
		assert.Less(t, index1, index2, "Inference Service 1 should come before Inference Service 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundInfSvc1.GetAttributes().CreateTimeSinceEpoch, *foundInfSvc2.GetAttributes().CreateTimeSinceEpoch, "Inference Service 1 should have earlier create time")
	})

	t.Run("TestSaveWithModelVersion", func(t *testing.T) {
		// First create a parent serving environment
		parentServingEnv := &models.ServingEnvironmentImpl{
			TypeID: apiutils.Of(int32(servingEnvironmentTypeID)),
			Attributes: &models.ServingEnvironmentAttributes{
				Name: apiutils.Of("parent-serving-env-for-model-version"),
			},
		}
		savedServingEnv, err := servingEnvironmentRepo.Save(parentServingEnv)
		require.NoError(t, err)

		// Create a registered model
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-with-version"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		// Create a model version
		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version"),
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

		// Create inference service with both registered model and model version
		inferenceService := &models.InferenceServiceImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.InferenceServiceAttributes{
				Name: apiutils.Of("inference-service-with-model-version"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Inference service with model version"),
				},
				{
					Name:     "serving_environment_id",
					IntValue: savedServingEnv.GetID(),
				},
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
				{
					Name:     "model_version_id",
					IntValue: savedModelVersion.GetID(),
				},
				{
					Name:        "runtime",
					StringValue: apiutils.Of("onnx"),
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

		saved, err := repo.Save(inferenceService)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 5) // description, serving_environment_id, registered_model_id, model_version_id, runtime

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify the specific properties exist
		properties := *retrieved.GetProperties()
		var foundModelVersionID, foundRegisteredModelID, foundServingEnvID bool
		for _, prop := range properties {
			if prop.Name == "model_version_id" && prop.IntValue != nil && *prop.IntValue == *savedModelVersion.GetID() {
				foundModelVersionID = true
			}
			if prop.Name == "registered_model_id" && prop.IntValue != nil && *prop.IntValue == *savedRegisteredModel.GetID() {
				foundRegisteredModelID = true
			}
			if prop.Name == "serving_environment_id" && prop.IntValue != nil && *prop.IntValue == *savedServingEnv.GetID() {
				foundServingEnvID = true
			}
		}
		assert.True(t, foundModelVersionID, "Should find model_version_id property")
		assert.True(t, foundRegisteredModelID, "Should find registered_model_id property")
		assert.True(t, foundServingEnvID, "Should find serving_environment_id property")
	})
}
