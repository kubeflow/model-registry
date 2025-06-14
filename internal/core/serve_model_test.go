package core_test

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/core"
	"github.com/kubeflow/model-registry/internal/ptr"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertServeModel(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		// Create prerequisites: registered model, serving environment, model version, and inference service
		registeredModel := &openapi.RegisteredModel{
			Name: "test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create serve model
		input := &openapi.ServeModel{
			Name:           ptr.Of("test-serve-model"),
			Description:    ptr.Of("Test serve model description"),
			ExternalId:     ptr.Of("serve-ext-123"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_RUNNING),
		}

		result, err := service.UpsertServeModel(input, createdInfSvc.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-serve-model", *result.Name)
		assert.Equal(t, "serve-ext-123", *result.ExternalId)
		assert.Equal(t, "Test serve model description", *result.Description)
		assert.Equal(t, *createdVersion.Id, result.ModelVersionId)
		assert.Equal(t, openapi.EXECUTIONSTATE_RUNNING, *result.LastKnownState)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful update", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "update-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "update-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "update-test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("update-test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create first
		input := &openapi.ServeModel{
			Name:           ptr.Of("update-test-serve-model"),
			Description:    ptr.Of("Original description"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_NEW),
		}

		created, err := service.UpsertServeModel(input, createdInfSvc.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update
		update := &openapi.ServeModel{
			Id:             created.Id,
			Name:           ptr.Of("update-test-serve-model"), // Name should remain the same
			Description:    ptr.Of("Updated description"),
			ExternalId:     ptr.Of("updated-ext-456"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_COMPLETE),
		}

		updated, err := service.UpsertServeModel(update, createdInfSvc.Id)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "update-test-serve-model", *updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-ext-456", *updated.ExternalId)
		assert.Equal(t, openapi.EXECUTIONSTATE_COMPLETE, *updated.LastKnownState)
	})

	t.Run("create with custom properties", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "custom-props-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "custom-props-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "custom-props-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("custom-props-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		customProps := map[string]openapi.MetadataValue{
			"deployment_config": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "production",
				},
			},
			"replicas": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "3",
				},
			},
			"auto_scaling": {
				MetadataBoolValue: &openapi.MetadataBoolValue{
					BoolValue: true,
				},
			},
			"cpu_limit": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue: 2.5,
				},
			},
		}

		input := &openapi.ServeModel{
			Name:             ptr.Of("custom-props-serve-model"),
			ModelVersionId:   *createdVersion.Id,
			LastKnownState:   ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
			CustomProperties: &customProps,
		}

		result, err := service.UpsertServeModel(input, createdInfSvc.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-serve-model", *result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := *result.CustomProperties
		assert.Contains(t, resultProps, "deployment_config")
		assert.Contains(t, resultProps, "replicas")
		assert.Contains(t, resultProps, "auto_scaling")
		assert.Contains(t, resultProps, "cpu_limit")

		assert.Equal(t, "production", resultProps["deployment_config"].MetadataStringValue.StringValue)
		assert.Equal(t, "3", resultProps["replicas"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["auto_scaling"].MetadataBoolValue.BoolValue)
		assert.Equal(t, 2.5, resultProps["cpu_limit"].MetadataDoubleValue.DoubleValue)
	})

	t.Run("minimal serve model", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "minimal-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "minimal-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "minimal-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("minimal-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		input := &openapi.ServeModel{
			Name:           ptr.Of("minimal-serve-model"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
		}

		result, err := service.UpsertServeModel(input, createdInfSvc.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-serve-model", *result.Name)
		assert.NotNil(t, result.Id)
		assert.Equal(t, *createdVersion.Id, result.ModelVersionId)
	})

	t.Run("nil serve model error", func(t *testing.T) {
		inferenceServiceId := "1"
		result, err := service.UpsertServeModel(nil, &inferenceServiceId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid serve model pointer")
	})

	t.Run("nil state preserved", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "nil-state-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "nil-state-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "nil-state-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("nil-state-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create serve model with nil LastKnownState
		input := &openapi.ServeModel{
			Name:           ptr.Of("nil-state-serve-model"),
			Description:    ptr.Of("Test serve model with nil state"),
			ExternalId:     ptr.Of("nil-state-ext-123"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: nil, // Explicitly set to nil
		}

		result, err := service.UpsertServeModel(input, createdInfSvc.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "nil-state-serve-model", *result.Name)
		assert.Equal(t, "nil-state-ext-123", *result.ExternalId)
		assert.Equal(t, "Test serve model with nil state", *result.Description)
		assert.Equal(t, *createdVersion.Id, result.ModelVersionId)
		assert.Nil(t, result.LastKnownState) // Verify state remains nil
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})
}

func TestGetServeModelById(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "get-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "get-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "get-test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("get-test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// First create a serve model to retrieve
		input := &openapi.ServeModel{
			Name:           ptr.Of("get-test-serve-model"),
			Description:    ptr.Of("Test description"),
			ExternalId:     ptr.Of("get-ext-123"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_RUNNING),
		}

		created, err := service.UpsertServeModel(input, createdInfSvc.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the serve model by ID
		result, err := service.GetServeModelById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-test-serve-model", *result.Name)
		assert.Equal(t, "get-ext-123", *result.ExternalId)
		assert.Equal(t, "Test description", *result.Description)
		assert.Equal(t, openapi.EXECUTIONSTATE_RUNNING, *result.LastKnownState)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := service.GetServeModelById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := service.GetServeModelById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no ServeModel found")
	})
}

func TestGetServeModels(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "list-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "list-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "list-test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("list-test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create multiple serve models for listing
		testServeModels := []*openapi.ServeModel{
			{
				Name:           ptr.Of("list-serve-model-1"),
				ExternalId:     ptr.Of("list-ext-1"),
				ModelVersionId: *createdVersion.Id,
				LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_RUNNING),
			},
			{
				Name:           ptr.Of("list-serve-model-2"),
				ExternalId:     ptr.Of("list-ext-2"),
				ModelVersionId: *createdVersion.Id,
				LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_NEW),
			},
			{
				Name:           ptr.Of("list-serve-model-3"),
				ExternalId:     ptr.Of("list-ext-3"),
				ModelVersionId: *createdVersion.Id,
				LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_COMPLETE),
			},
		}

		var createdIds []string
		for _, srvModel := range testServeModels {
			created, err := service.UpsertServeModel(srvModel, createdInfSvc.Id)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List serve models with basic pagination
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetServeModels(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should have at least our 3 test serve models
		assert.Equal(t, int32(10), result.PageSize)

		// Verify our serve models are in the result
		foundModels := 0
		for _, item := range result.Items {
			for _, createdId := range createdIds {
				if *item.Id == createdId {
					foundModels++
					break
				}
			}
		}
		assert.Equal(t, 3, foundModels, "All created serve models should be found in the list")
	})

	t.Run("list with inference service filter", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "filter-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "filter-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "filter-test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService1 := &openapi.InferenceService{
			Name:                 ptr.Of("filter-test-inference-service-1"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc1, err := service.UpsertInferenceService(inferenceService1)
		require.NoError(t, err)

		inferenceService2 := &openapi.InferenceService{
			Name:                 ptr.Of("filter-test-inference-service-2"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc2, err := service.UpsertInferenceService(inferenceService2)
		require.NoError(t, err)

		// Create serve models in different inference services
		srvModel1 := &openapi.ServeModel{
			Name:           ptr.Of("filter-serve-model-1"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
		}
		created1, err := service.UpsertServeModel(srvModel1, createdInfSvc1.Id)
		require.NoError(t, err)

		srvModel2 := &openapi.ServeModel{
			Name:           ptr.Of("filter-serve-model-2"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
		}
		_, err = service.UpsertServeModel(srvModel2, createdInfSvc2.Id)
		require.NoError(t, err)

		// List serve models filtered by inference service
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetServeModels(listOptions, createdInfSvc1.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should have at least our 1 model in infSvc1

		// Verify that only models from the specified inference service are returned
		found := false
		for _, item := range result.Items {
			if *item.Id == *created1.Id {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the serve model in the specified inference service")
	})

	t.Run("pagination and ordering", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "pagination-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "pagination-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "pagination-test-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("pagination-test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create several serve models for pagination testing
		for i := 0; i < 5; i++ {
			srvModel := &openapi.ServeModel{
				Name:           ptr.Of("pagination-serve-model-" + string(rune('A'+i))),
				ExternalId:     ptr.Of("pagination-ext-" + string(rune('A'+i))),
				ModelVersionId: *createdVersion.Id,
				LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
			}
			_, err := service.UpsertServeModel(srvModel, createdInfSvc.Id)
			require.NoError(t, err)
		}

		// Test with small page size and ordering
		pageSize := int32(2)
		orderBy := "name"
		sortOrder := "asc"
		listOptions := api.ListOptions{
			PageSize:  &pageSize,
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		}

		result, err := service.GetServeModels(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})

	t.Run("invalid inference service id", func(t *testing.T) {
		invalidId := "invalid"
		listOptions := api.ListOptions{}

		result, err := service.GetServeModels(listOptions, &invalidId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid inference service id")
	})
}

func TestServeModelRoundTrip(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name:        "roundtrip-registered-model",
			Description: ptr.Of("Roundtrip test registered model"),
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name:        "roundtrip-serving-env",
			Description: ptr.Of("Roundtrip test serving environment"),
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "roundtrip-model-version",
			Description:       ptr.Of("Roundtrip test model version"),
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("roundtrip-inference-service"),
			Description:          ptr.Of("Roundtrip test inference service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create a serve model with all fields
		original := &openapi.ServeModel{
			Name:           ptr.Of("roundtrip-serve-model"),
			Description:    ptr.Of("Roundtrip test description"),
			ExternalId:     ptr.Of("roundtrip-ext-123"),
			ModelVersionId: *createdVersion.Id,
			LastKnownState: ptr.Of(openapi.EXECUTIONSTATE_RUNNING),
		}

		// Create
		created, err := service.UpsertServeModel(original, createdInfSvc.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetServeModelById(*created.Id)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, *original.Name, *retrieved.Name)
		assert.Equal(t, *original.Description, *retrieved.Description)
		assert.Equal(t, *original.ExternalId, *retrieved.ExternalId)
		assert.Equal(t, original.ModelVersionId, retrieved.ModelVersionId)
		assert.Equal(t, *original.LastKnownState, *retrieved.LastKnownState)

		// Update
		retrieved.Description = ptr.Of("Updated description")
		retrieved.LastKnownState = ptr.Of(openapi.EXECUTIONSTATE_COMPLETE)

		updated, err := service.UpsertServeModel(retrieved, createdInfSvc.Id)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, openapi.EXECUTIONSTATE_COMPLETE, *updated.LastKnownState)

		// Get again to verify persistence
		final, err := service.GetServeModelById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, openapi.EXECUTIONSTATE_COMPLETE, *final.LastKnownState)
	})

	t.Run("roundtrip with custom properties", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "roundtrip-custom-props-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "roundtrip-custom-props-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:              "roundtrip-custom-props-model-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("roundtrip-custom-props-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInfSvc, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		customProps := map[string]openapi.MetadataValue{
			"environment": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "staging",
				},
			},
			"max_requests": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "100",
				},
			},
		}

		original := &openapi.ServeModel{
			Name:             ptr.Of("roundtrip-custom-props-serve-model"),
			ModelVersionId:   *createdVersion.Id,
			LastKnownState:   ptr.Of(openapi.EXECUTIONSTATE_UNKNOWN),
			CustomProperties: &customProps,
		}

		// Create
		created, err := service.UpsertServeModel(original, createdInfSvc.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetServeModelById(*created.Id)
		require.NoError(t, err)

		// Verify custom properties
		assert.NotNil(t, retrieved.CustomProperties)
		retrievedProps := *retrieved.CustomProperties
		assert.Contains(t, retrievedProps, "environment")
		assert.Contains(t, retrievedProps, "max_requests")
		assert.Equal(t, "staging", retrievedProps["environment"].MetadataStringValue.StringValue)
		assert.Equal(t, "100", retrievedProps["max_requests"].MetadataIntValue.IntValue)

		// Update custom properties
		updatedProps := map[string]openapi.MetadataValue{
			"environment": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "production",
				},
			},
			"max_requests": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "500",
				},
			},
			"new_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "new_value",
				},
			},
		}
		retrieved.CustomProperties = &updatedProps

		updated, err := service.UpsertServeModel(retrieved, createdInfSvc.Id)
		require.NoError(t, err)

		// Verify updated custom properties
		assert.NotNil(t, updated.CustomProperties)
		finalProps := *updated.CustomProperties
		assert.Equal(t, "production", finalProps["environment"].MetadataStringValue.StringValue)
		assert.Equal(t, "500", finalProps["max_requests"].MetadataIntValue.IntValue)
		assert.Equal(t, "new_value", finalProps["new_prop"].MetadataStringValue.StringValue)
	})
}
