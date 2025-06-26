package core_test

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/core"
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
			Description: apiutils.Of("Test model description"),
			Owner:       apiutils.Of("test-owner"),
			ExternalId:  apiutils.Of("ext-123"),
			State:       apiutils.Of(openapi.REGISTEREDMODELSTATE_LIVE),
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
			Description: apiutils.Of("Original description"),
		}

		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update
		update := &openapi.RegisteredModel{
			Id:          created.Id,
			Name:        "update-test-model", // Name should remain the same
			Description: apiutils.Of("Updated description"),
			Owner:       apiutils.Of("new-owner"),
			State:       apiutils.Of(openapi.REGISTEREDMODELSTATE_ARCHIVED),
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

	t.Run("nil fields preserved", func(t *testing.T) {
		// Create registered model with nil optional fields
		input := &openapi.RegisteredModel{
			Name:        "nil-fields-model",
			Description: nil, // Explicitly set to nil
			Owner:       nil, // Explicitly set to nil
			ExternalId:  nil, // Explicitly set to nil
			State:       nil, // Explicitly set to nil
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "nil-fields-model", result.Name)
		assert.Nil(t, result.Description) // Verify description remains nil
		assert.Nil(t, result.Owner)       // Verify owner remains nil
		assert.Nil(t, result.ExternalId)  // Verify external ID remains nil
		assert.Nil(t, result.State)       // Verify state remains nil
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})

	t.Run("unicode characters in name", func(t *testing.T) {
		unicodeName := "测试模型-тест-モデル-🚀"
		input := &openapi.RegisteredModel{
			Name:        unicodeName,
			Description: apiutils.Of("Unicode test model with 中文, русский, 日本語, and emoji 🎯"),
			Owner:       apiutils.Of("用户-пользователь-ユーザー"),
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, unicodeName, result.Name)
		assert.Equal(t, "Unicode test model with 中文, русский, 日本語, and emoji 🎯", *result.Description)
		assert.Equal(t, "用户-пользователь-ユーザー", *result.Owner)
		assert.NotNil(t, result.Id)
	})

	t.Run("special characters in name", func(t *testing.T) {
		specialName := "test-model!@#$%^&*()_+-=[]{}|;':\",./<>?"
		input := &openapi.RegisteredModel{
			Name:        specialName,
			Description: apiutils.Of("Model with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?"),
			ExternalId:  apiutils.Of("ext-id-with-special-chars_123!@#"),
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, specialName, result.Name)
		assert.Equal(t, "Model with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?", *result.Description)
		assert.Equal(t, "ext-id-with-special-chars_123!@#", *result.ExternalId)
		assert.NotNil(t, result.Id)
	})

	t.Run("mixed unicode and special characters", func(t *testing.T) {
		mixedName := "模型-test!@#-тест_123-🚀"
		input := &openapi.RegisteredModel{
			Name:        mixedName,
			Description: apiutils.Of("Mixed: 测试!@# русский_test 日本語-123 🎯"),
			Owner:       apiutils.Of("owner@domain.com-用户_123"),
			ExternalId:  apiutils.Of("ext-混合_test!@#-123"),
		}

		result, err := service.UpsertRegisteredModel(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, mixedName, result.Name)
		assert.Equal(t, "Mixed: 测试!@# русский_test 日本語-123 🎯", *result.Description)
		assert.Equal(t, "owner@domain.com-用户_123", *result.Owner)
		assert.Equal(t, "ext-混合_test!@#-123", *result.ExternalId)
		assert.NotNil(t, result.Id)
	})

	t.Run("pagination with 10+ models", func(t *testing.T) {
		// Create 15 models to test pagination
		var createdModels []string
		for i := 0; i < 15; i++ {
			input := &openapi.RegisteredModel{
				Name:        fmt.Sprintf("paging-test-model-%02d", i),
				Description: apiutils.Of(fmt.Sprintf("Test model %d for pagination", i)),
				ExternalId:  apiutils.Of(fmt.Sprintf("paging-ext-%02d", i)),
			}

			result, err := service.UpsertRegisteredModel(input)
			require.NoError(t, err)
			require.NotNil(t, result.Id)
			createdModels = append(createdModels, *result.Id)
		}

		// Test first page with page size 5
		pageSize := int32(5)
		firstPageResult, err := service.GetRegisteredModels(api.ListOptions{
			PageSize: &pageSize,
		})

		require.NoError(t, err)
		require.NotNil(t, firstPageResult)
		assert.LessOrEqual(t, len(firstPageResult.Items), int(pageSize))
		assert.Equal(t, pageSize, firstPageResult.PageSize)

		// Test second page if there's a next page token
		if firstPageResult.NextPageToken != "" {
			secondPageResult, err := service.GetRegisteredModels(api.ListOptions{
				PageSize:      &pageSize,
				NextPageToken: &firstPageResult.NextPageToken,
			})

			require.NoError(t, err)
			require.NotNil(t, secondPageResult)
			assert.LessOrEqual(t, len(secondPageResult.Items), int(pageSize))
			assert.Equal(t, pageSize, secondPageResult.PageSize)

			// Verify no duplicate models between pages
			firstPageIds := make(map[string]bool)
			for _, model := range firstPageResult.Items {
				firstPageIds[*model.Id] = true
			}

			for _, model := range secondPageResult.Items {
				assert.False(t, firstPageIds[*model.Id], "Model %s appears in both pages", *model.Id)
			}
		}

		// Test larger page size to get more models
		largePageSize := int32(100)
		largePageResult, err := service.GetRegisteredModels(api.ListOptions{
			PageSize: &largePageSize,
		})

		require.NoError(t, err)
		require.NotNil(t, largePageResult)
		assert.GreaterOrEqual(t, len(largePageResult.Items), 15) // Should include our 15 models
		assert.Equal(t, largePageSize, largePageResult.PageSize)

		// Verify our created models are in the results
		resultIds := make(map[string]bool)
		for _, model := range largePageResult.Items {
			resultIds[*model.Id] = true
		}

		foundCount := 0
		for _, createdId := range createdModels {
			if resultIds[createdId] {
				foundCount++
			}
		}
		assert.Equal(t, 15, foundCount, "Should find all 15 created models in the results")

		// Test ordering by name
		orderBy := "name"
		sortOrder := "ASC"
		orderedResult, err := service.GetRegisteredModels(api.ListOptions{
			PageSize:  &largePageSize,
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		})

		require.NoError(t, err)
		require.NotNil(t, orderedResult)

		// Verify ordering (at least check that we have results)
		assert.Greater(t, len(orderedResult.Items), 0)

		// Test descending order
		sortOrderDesc := "DESC"
		orderedDescResult, err := service.GetRegisteredModels(api.ListOptions{
			PageSize:  &largePageSize,
			OrderBy:   &orderBy,
			SortOrder: &sortOrderDesc,
		})

		require.NoError(t, err)
		require.NotNil(t, orderedDescResult)
		assert.Greater(t, len(orderedDescResult.Items), 0)
	})
}

func TestGetRegisteredModelById(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// First create a model to retrieve
		input := &openapi.RegisteredModel{
			Name:        "get-test-model",
			Description: apiutils.Of("Test description"),
			ExternalId:  apiutils.Of("get-ext-123"),
			State:       apiutils.Of(openapi.REGISTEREDMODELSTATE_LIVE),
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
			Name:                 apiutils.Of("test-inference-service"),
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
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
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
			ExternalId: apiutils.Of("params-ext-123"),
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
			ExternalId: apiutils.Of("params-unique-ext-456"),
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

	t.Run("no external id", func(t *testing.T) {
		modelName := "params-test-model-no-external-id"
		input := &openapi.RegisteredModel{
			Name: modelName,
		}
		created, err := service.UpsertRegisteredModel(input)
		require.NoError(t, err)

		result, err := service.GetRegisteredModelByParams(&modelName, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, modelName, result.Name)
	})
}

func TestGetRegisteredModels(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create multiple models for listing
		testModels := []*openapi.RegisteredModel{
			{Name: "list-model-1", ExternalId: apiutils.Of("list-ext-1")},
			{Name: "list-model-2", ExternalId: apiutils.Of("list-ext-2")},
			{Name: "list-model-3", ExternalId: apiutils.Of("list-ext-3")},
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
				ExternalId: apiutils.Of("pagination-ext-" + string(rune('A'+i))),
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
			Description: apiutils.Of("Roundtrip test description"),
			Owner:       apiutils.Of("roundtrip-owner"),
			ExternalId:  apiutils.Of("roundtrip-ext-123"),
			State:       apiutils.Of(openapi.REGISTEREDMODELSTATE_LIVE),
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
		retrieved.Description = apiutils.Of("Updated description")
		retrieved.State = apiutils.Of(openapi.REGISTEREDMODELSTATE_ARCHIVED)

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
