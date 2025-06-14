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

func TestUpsertInferenceService(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		// Create prerequisites: registered model and serving environment
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

		// Create inference service
		input := &openapi.InferenceService{
			Name:                 ptr.Of("test-inference-service"),
			Description:          ptr.Of("Test inference service description"),
			ExternalId:           ptr.Of("inference-ext-123"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("tensorflow"),
		}

		result, err := service.UpsertInferenceService(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-inference-service", *result.Name)
		assert.Equal(t, "inference-ext-123", *result.ExternalId)
		assert.Equal(t, "Test inference service description", *result.Description)
		assert.Equal(t, *createdEnv.Id, result.ServingEnvironmentId)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
		assert.Equal(t, "tensorflow", *result.Runtime)
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

		// Create first
		input := &openapi.InferenceService{
			Name:                 ptr.Of("update-test-inference-service"),
			Description:          ptr.Of("Original description"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}

		created, err := service.UpsertInferenceService(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update
		update := &openapi.InferenceService{
			Id:                   created.Id,
			Name:                 ptr.Of("update-test-inference-service"), // Name should remain the same
			Description:          ptr.Of("Updated description"),
			ExternalId:           ptr.Of("updated-ext-456"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("pytorch"),
		}

		updated, err := service.UpsertInferenceService(update)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "update-test-inference-service", *updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-ext-456", *updated.ExternalId)
		assert.Equal(t, "pytorch", *updated.Runtime)
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

		customProps := map[string]openapi.MetadataValue{
			"model_uri": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "s3://bucket/model",
				},
			},
			"batch_size": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "32",
				},
			},
			"enable_logging": {
				MetadataBoolValue: &openapi.MetadataBoolValue{
					BoolValue: true,
				},
			},
			"confidence_threshold": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue: 0.85,
				},
			},
		}

		input := &openapi.InferenceService{
			Name:                 ptr.Of("custom-props-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			CustomProperties:     &customProps,
		}

		result, err := service.UpsertInferenceService(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-inference-service", *result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := *result.CustomProperties
		assert.Contains(t, resultProps, "model_uri")
		assert.Contains(t, resultProps, "batch_size")
		assert.Contains(t, resultProps, "enable_logging")
		assert.Contains(t, resultProps, "confidence_threshold")

		assert.Equal(t, "s3://bucket/model", resultProps["model_uri"].MetadataStringValue.StringValue)
		assert.Equal(t, "32", resultProps["batch_size"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["enable_logging"].MetadataBoolValue.BoolValue)
		assert.Equal(t, 0.85, resultProps["confidence_threshold"].MetadataDoubleValue.DoubleValue)
	})

	t.Run("minimal inference service", func(t *testing.T) {
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

		input := &openapi.InferenceService{
			Name:                 ptr.Of("minimal-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}

		result, err := service.UpsertInferenceService(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-inference-service", *result.Name)
		assert.NotNil(t, result.Id)
		assert.Equal(t, *createdEnv.Id, result.ServingEnvironmentId)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
	})

	t.Run("nil inference service error", func(t *testing.T) {
		result, err := service.UpsertInferenceService(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid inference service pointer")
	})

	t.Run("nil desired state preserved", func(t *testing.T) {
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

		// Create inference service with nil DesiredState and other optional fields
		input := &openapi.InferenceService{
			Name:                 ptr.Of("nil-state-inference-service"),
			Description:          ptr.Of("Test inference service with nil desired state"),
			ExternalId:           ptr.Of("nil-state-ext-123"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("tensorflow"),
			DesiredState:         nil, // Explicitly set to nil
		}

		result, err := service.UpsertInferenceService(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "nil-state-inference-service", *result.Name)
		assert.Equal(t, "nil-state-ext-123", *result.ExternalId)
		assert.Equal(t, "Test inference service with nil desired state", *result.Description)
		assert.Equal(t, *createdEnv.Id, result.ServingEnvironmentId)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
		assert.Equal(t, "tensorflow", *result.Runtime)
		assert.Nil(t, result.DesiredState) // Verify desired state remains nil
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})
}

func TestGetInferenceServiceById(t *testing.T) {
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

		// First create an inference service to retrieve
		input := &openapi.InferenceService{
			Name:                 ptr.Of("get-test-inference-service"),
			Description:          ptr.Of("Test description"),
			ExternalId:           ptr.Of("get-ext-123"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("tensorflow"),
		}

		created, err := service.UpsertInferenceService(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the inference service by ID
		result, err := service.GetInferenceServiceById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-test-inference-service", *result.Name)
		assert.Equal(t, "get-ext-123", *result.ExternalId)
		assert.Equal(t, "Test description", *result.Description)
		assert.Equal(t, "tensorflow", *result.Runtime)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := service.GetInferenceServiceById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := service.GetInferenceServiceById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no InferenceService found")
	})
}

func TestGetInferenceServiceByParams(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get by name and parent resource id", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "params-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "params-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		input := &openapi.InferenceService{
			Name:                 ptr.Of("params-test-inference-service"),
			ExternalId:           ptr.Of("params-ext-123"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		created, err := service.UpsertInferenceService(input)
		require.NoError(t, err)

		// Get by name and parent resource ID
		serviceName := "params-test-inference-service"
		result, err := service.GetInferenceServiceByParams(&serviceName, createdEnv.Id, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-test-inference-service", *result.Name)
	})

	t.Run("successful get by external id", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "params-ext-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "params-ext-test-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		input := &openapi.InferenceService{
			Name:                 ptr.Of("params-ext-test-inference-service"),
			ExternalId:           ptr.Of("params-unique-ext-456"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		created, err := service.UpsertInferenceService(input)
		require.NoError(t, err)

		// Get by external ID
		externalId := "params-unique-ext-456"
		result, err := service.GetInferenceServiceByParams(nil, nil, &externalId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-unique-ext-456", *result.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := service.GetInferenceServiceByParams(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters call")
	})

	t.Run("no inference service found", func(t *testing.T) {
		serviceName := "nonexistent-inference-service"
		parentResourceId := "999"
		result, err := service.GetInferenceServiceByParams(&serviceName, &parentResourceId, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no inference service found")
	})
}

func TestGetInferenceServices(t *testing.T) {
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

		// Create multiple inference services for listing
		testInferenceServices := []*openapi.InferenceService{
			{
				Name:                 ptr.Of("list-inference-service-1"),
				ExternalId:           ptr.Of("list-ext-1"),
				ServingEnvironmentId: *createdEnv.Id,
				RegisteredModelId:    *createdModel.Id,
				Runtime:              ptr.Of("tensorflow"),
			},
			{
				Name:                 ptr.Of("list-inference-service-2"),
				ExternalId:           ptr.Of("list-ext-2"),
				ServingEnvironmentId: *createdEnv.Id,
				RegisteredModelId:    *createdModel.Id,
				Runtime:              ptr.Of("pytorch"),
			},
			{
				Name:                 ptr.Of("list-inference-service-3"),
				ExternalId:           ptr.Of("list-ext-3"),
				ServingEnvironmentId: *createdEnv.Id,
				RegisteredModelId:    *createdModel.Id,
				Runtime:              ptr.Of("onnx"),
			},
		}

		var createdIds []string
		for _, infSvc := range testInferenceServices {
			created, err := service.UpsertInferenceService(infSvc)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List inference services with basic pagination
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetInferenceServices(listOptions, nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should have at least our 3 test inference services
		assert.Equal(t, int32(10), result.PageSize)

		// Verify our inference services are in the result
		foundServices := 0
		for _, item := range result.Items {
			for _, createdId := range createdIds {
				if *item.Id == createdId {
					foundServices++
					break
				}
			}
		}
		assert.Equal(t, 3, foundServices, "All created inference services should be found in the list")
	})

	t.Run("list with serving environment filter", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "filter-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv1 := &openapi.ServingEnvironment{
			Name: "filter-test-serving-env-1",
		}
		createdEnv1, err := service.UpsertServingEnvironment(servingEnv1)
		require.NoError(t, err)

		servingEnv2 := &openapi.ServingEnvironment{
			Name: "filter-test-serving-env-2",
		}
		createdEnv2, err := service.UpsertServingEnvironment(servingEnv2)
		require.NoError(t, err)

		// Create inference services in different serving environments
		infSvc1 := &openapi.InferenceService{
			Name:                 ptr.Of("filter-inference-service-1"),
			ServingEnvironmentId: *createdEnv1.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		created1, err := service.UpsertInferenceService(infSvc1)
		require.NoError(t, err)

		infSvc2 := &openapi.InferenceService{
			Name:                 ptr.Of("filter-inference-service-2"),
			ServingEnvironmentId: *createdEnv2.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		_, err = service.UpsertInferenceService(infSvc2)
		require.NoError(t, err)

		// List inference services filtered by serving environment
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetInferenceServices(listOptions, createdEnv1.Id, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should have at least our 1 service in env1

		// Verify that only services from the specified environment are returned
		found := false
		for _, item := range result.Items {
			if *item.Id == *created1.Id {
				found = true
				assert.Equal(t, *createdEnv1.Id, item.ServingEnvironmentId)
			}
		}
		assert.True(t, found, "Should find the inference service in the specified serving environment")
	})

	t.Run("list with runtime filter", func(t *testing.T) {
		// Create prerequisites
		registeredModel := &openapi.RegisteredModel{
			Name: "runtime-filter-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "runtime-filter-serving-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create inference services with different runtimes
		infSvcTensorflow := &openapi.InferenceService{
			Name:                 ptr.Of("runtime-tensorflow-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("tensorflow"),
		}
		createdTensorflow, err := service.UpsertInferenceService(infSvcTensorflow)
		require.NoError(t, err)

		infSvcPytorch := &openapi.InferenceService{
			Name:                 ptr.Of("runtime-pytorch-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("pytorch"),
		}
		_, err = service.UpsertInferenceService(infSvcPytorch)
		require.NoError(t, err)

		// List inference services filtered by runtime
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}
		runtime := "tensorflow"

		result, err := service.GetInferenceServices(listOptions, nil, &runtime)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should have at least our 1 tensorflow service

		// Verify that only services with the specified runtime are returned
		found := false
		for _, item := range result.Items {
			if *item.Id == *createdTensorflow.Id {
				found = true
				assert.Equal(t, "tensorflow", *item.Runtime)
			}
		}
		assert.True(t, found, "Should find the inference service with the specified runtime")
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

		// Create several inference services for pagination testing
		for i := 0; i < 5; i++ {
			infSvc := &openapi.InferenceService{
				Name:                 ptr.Of("pagination-inference-service-" + string(rune('A'+i))),
				ExternalId:           ptr.Of("pagination-ext-" + string(rune('A'+i))),
				ServingEnvironmentId: *createdEnv.Id,
				RegisteredModelId:    *createdModel.Id,
			}
			_, err := service.UpsertInferenceService(infSvc)
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

		result, err := service.GetInferenceServices(listOptions, nil, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})

	t.Run("invalid serving environment id", func(t *testing.T) {
		invalidId := "invalid"
		listOptions := api.ListOptions{}

		result, err := service.GetInferenceServices(listOptions, &invalidId, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid serving environment id")
	})
}

func TestInferenceServiceRoundTrip(t *testing.T) {
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

		// Create an inference service with all fields
		original := &openapi.InferenceService{
			Name:                 ptr.Of("roundtrip-inference-service"),
			Description:          ptr.Of("Roundtrip test description"),
			ExternalId:           ptr.Of("roundtrip-ext-123"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			Runtime:              ptr.Of("tensorflow"),
		}

		// Create
		created, err := service.UpsertInferenceService(original)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetInferenceServiceById(*created.Id)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, *original.Name, *retrieved.Name)
		assert.Equal(t, *original.Description, *retrieved.Description)
		assert.Equal(t, *original.ExternalId, *retrieved.ExternalId)
		assert.Equal(t, original.ServingEnvironmentId, retrieved.ServingEnvironmentId)
		assert.Equal(t, original.RegisteredModelId, retrieved.RegisteredModelId)
		assert.Equal(t, *original.Runtime, *retrieved.Runtime)

		// Update
		retrieved.Description = ptr.Of("Updated description")
		retrieved.Runtime = ptr.Of("pytorch")

		updated, err := service.UpsertInferenceService(retrieved)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "pytorch", *updated.Runtime)

		// Get again to verify persistence
		final, err := service.GetInferenceServiceById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, "pytorch", *final.Runtime)
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

		customProps := map[string]openapi.MetadataValue{
			"deployment_type": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "canary",
				},
			},
			"replicas": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "3",
				},
			},
		}

		original := &openapi.InferenceService{
			Name:                 ptr.Of("roundtrip-custom-props-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			CustomProperties:     &customProps,
		}

		// Create
		created, err := service.UpsertInferenceService(original)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetInferenceServiceById(*created.Id)
		require.NoError(t, err)

		// Verify custom properties
		assert.NotNil(t, retrieved.CustomProperties)
		retrievedProps := *retrieved.CustomProperties
		assert.Contains(t, retrievedProps, "deployment_type")
		assert.Contains(t, retrievedProps, "replicas")
		assert.Equal(t, "canary", retrievedProps["deployment_type"].MetadataStringValue.StringValue)
		assert.Equal(t, "3", retrievedProps["replicas"].MetadataIntValue.IntValue)

		// Update custom properties
		updatedProps := map[string]openapi.MetadataValue{
			"deployment_type": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "blue-green",
				},
			},
			"replicas": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "5",
				},
			},
			"new_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "new_value",
				},
			},
		}
		retrieved.CustomProperties = &updatedProps

		updated, err := service.UpsertInferenceService(retrieved)
		require.NoError(t, err)

		// Verify updated custom properties
		assert.NotNil(t, updated.CustomProperties)
		finalProps := *updated.CustomProperties
		assert.Equal(t, "blue-green", finalProps["deployment_type"].MetadataStringValue.StringValue)
		assert.Equal(t, "5", finalProps["replicas"].MetadataIntValue.IntValue)
		assert.Equal(t, "new_value", finalProps["new_prop"].MetadataStringValue.StringValue)
	})
}
