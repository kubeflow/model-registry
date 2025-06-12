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

func TestUpsertRegisteredModel(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		input := &openapi.RegisteredModel{
			Name:        "test-model",
			Description: ptr.Of("Test model description"),
			Owner:       ptr.Of("test-owner"),
			ExternalId:  ptr.Of("ext-123"),
			State:       ptr.Of(openapi.REGISTEREDMODELSTATE_LIVE),
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-model", result.Name)
		assert.Equal(t, "ext-123", *result.ExternalId)
		assert.Equal(t, "Test model description", *result.Description)
		assert.Equal(t, "test-owner", *result.Owner)
		assert.Equal(t, openapi.REGISTEREDMODELSTATE_LIVE, *result.State)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful update", func(t *testing.T) {
		// Create first
		input := &openapi.RegisteredModel{
			Name:        "update-test-model",
			Description: ptr.Of("Original description"),
		}

		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update
		update := &openapi.RegisteredModel{
			Id:          created.Id,
			Name:        "update-test-model", // Name should remain the same
			Description: ptr.Of("Updated description"),
			Owner:       ptr.Of("new-owner"),
			State:       ptr.Of(openapi.REGISTEREDMODELSTATE_ARCHIVED),
		}

		updated, err := service.UpsertRegisteredModel(update)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "update-test-model", updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "new-owner", *updated.Owner)
		assert.Equal(t, openapi.REGISTEREDMODELSTATE_ARCHIVED, *updated.State)
	})

	t.Run("create with custom properties", func(t *testing.T) {
		customProps := map[string]openapi.MetadataValue{
			"accuracy": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue: 0.95,
				},
			},
			"framework": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "tensorflow",
				},
			},
			"version": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "2",
				},
			},
			"is_production": {
				MetadataBoolValue: &openapi.MetadataBoolValue{
					BoolValue: true,
				},
			},
		}

		input := &openapi.RegisteredModel{
			Name:             "custom-props-model",
			CustomProperties: &customProps,
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-model", result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := *result.CustomProperties
		assert.Contains(t, resultProps, "accuracy")
		assert.Contains(t, resultProps, "framework")
		assert.Contains(t, resultProps, "version")
		assert.Contains(t, resultProps, "is_production")

		assert.Equal(t, 0.95, resultProps["accuracy"].MetadataDoubleValue.DoubleValue)
		assert.Equal(t, "tensorflow", resultProps["framework"].MetadataStringValue.StringValue)
		assert.Equal(t, "2", resultProps["version"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["is_production"].MetadataBoolValue.BoolValue)
	})

	t.Run("minimal model", func(t *testing.T) {
		input := &openapi.RegisteredModel{
			Name: "minimal-model",
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-model", result.Name)
		assert.NotNil(t, result.Id)
	})

	t.Run("nil model error", func(t *testing.T) {
		result, err := service.UpsertRegisteredModel(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid registered model pointer")
	})
}

func TestGetRegisteredModelById(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// First create a model to retrieve
		input := &openapi.RegisteredModel{
			Name:        "get-test-model",
			Description: ptr.Of("Test description"),
			ExternalId:  ptr.Of("get-ext-123"),
			State:       ptr.Of(openapi.REGISTEREDMODELSTATE_LIVE),
		}

		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the model by ID
		result, err := service.GetRegisteredModelById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-test-model", result.Name)
		assert.Equal(t, "get-ext-123", *result.ExternalId)
		assert.Equal(t, "Test description", *result.Description)
		assert.Equal(t, openapi.REGISTEREDMODELSTATE_LIVE, *result.State)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := service.GetRegisteredModelById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := service.GetRegisteredModelById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no registered model found")
	})
}

func TestGetRegisteredModelByInferenceService(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "inference-test-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a serving environment
		servingEnv := &openapi.ServingEnvironment{
			Name: "test-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create an inference service
		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
		}
		createdInference, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Get registered model by inference service
		result, err := service.GetRegisteredModelByInferenceService(*createdInference.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *createdModel.Id, *result.Id)
		assert.Equal(t, "inference-test-model", result.Name)
	})

	t.Run("invalid inference service id", func(t *testing.T) {
		result, err := service.GetRegisteredModelByInferenceService("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid inference service id")
	})

	t.Run("non-existent inference service", func(t *testing.T) {
		result, err := service.GetRegisteredModelByInferenceService("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetRegisteredModelByParams(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get by name", func(t *testing.T) {
		input := &openapi.RegisteredModel{
			Name:       "params-test-model",
			ExternalId: ptr.Of("params-ext-123"),
		}
		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)

		// Get by name
		modelName := "params-test-model"
		result, err := service.GetRegisteredModelByParams(&modelName, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-test-model", result.Name)
	})

	t.Run("successful get by external id", func(t *testing.T) {
		input := &openapi.RegisteredModel{
			Name:       "params-ext-test-model",
			ExternalId: ptr.Of("params-unique-ext-456"),
		}
		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)

		// Get by external ID
		externalId := "params-unique-ext-456"
		result, err := service.GetRegisteredModelByParams(nil, &externalId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-unique-ext-456", *result.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := service.GetRegisteredModelByParams(nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters call")
	})

	t.Run("no model found", func(t *testing.T) {
		modelName := "nonexistent-model"
		result, err := service.GetRegisteredModelByParams(&modelName, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no registered models found")
	})
}

func TestGetRegisteredModels(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create multiple models for listing
		testModels := []*openapi.RegisteredModel{
			{Name: "list-model-1", ExternalId: ptr.Of("list-ext-1")},
			{Name: "list-model-2", ExternalId: ptr.Of("list-ext-2")},
			{Name: "list-model-3", ExternalId: ptr.Of("list-ext-3")},
		}

		var createdIds []string
		for _, model := range testModels {
			created, err := service.UpsertRegisteredModel(model)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List models with basic pagination
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetRegisteredModels(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should have at least our 3 test models
		assert.Equal(t, int32(10), result.PageSize)

		// Verify our models are in the result
		foundModels := 0
		for _, item := range result.Items {
			for _, createdId := range createdIds {
				if *item.Id == createdId {
					foundModels++
					break
				}
			}
		}
		assert.Equal(t, 3, foundModels, "All created models should be found in the list")
	})

	t.Run("pagination and ordering", func(t *testing.T) {
		// Create several models for pagination testing
		for i := 0; i < 5; i++ {
			model := &openapi.RegisteredModel{
				Name:       "pagination-model-" + string(rune('A'+i)),
				ExternalId: ptr.Of("pagination-ext-" + string(rune('A'+i))),
			}
			_, err := service.UpsertRegisteredModel(model)
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

		result, err := service.GetRegisteredModels(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})
}

func TestRegisteredModelRoundTrip(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create a model with all fields
		original := &openapi.RegisteredModel{
			Name:        "roundtrip-model",
			Description: ptr.Of("Roundtrip test description"),
			Owner:       ptr.Of("roundtrip-owner"),
			ExternalId:  ptr.Of("roundtrip-ext-123"),
			State:       ptr.Of(openapi.REGISTEREDMODELSTATE_LIVE),
		}

		// Create
		created, err := service.UpsertRegisteredModel(original)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetRegisteredModelById(*created.Id)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, original.Name, retrieved.Name)
		assert.Equal(t, *original.Description, *retrieved.Description)
		assert.Equal(t, *original.Owner, *retrieved.Owner)
		assert.Equal(t, *original.ExternalId, *retrieved.ExternalId)
		assert.Equal(t, *original.State, *retrieved.State)

		// Update
		retrieved.Description = ptr.Of("Updated description")
		retrieved.State = ptr.Of(openapi.REGISTEREDMODELSTATE_ARCHIVED)

		updated, err := service.UpsertRegisteredModel(retrieved)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, openapi.REGISTEREDMODELSTATE_ARCHIVED, *updated.State)

		// Get again to verify persistence
		final, err := service.GetRegisteredModelById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, openapi.REGISTEREDMODELSTATE_ARCHIVED, *final.State)
	})
}
