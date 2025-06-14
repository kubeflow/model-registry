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

func TestUpsertModelVersion(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		// First create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name:        "test-registered-model",
			Description: ptr.Of("Test registered model for version"),
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)
		require.NotNil(t, createdModel.Id)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "test-version",
			Description:       ptr.Of("Test version description"),
			Author:            ptr.Of("test-author"),
			ExternalId:        ptr.Of("version-ext-123"),
			State:             ptr.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		result, err := service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-version", result.Name)
		assert.Equal(t, "version-ext-123", *result.ExternalId)
		assert.Equal(t, "Test version description", *result.Description)
		assert.Equal(t, "test-author", *result.Author)
		assert.Equal(t, openapi.MODELVERSIONSTATE_LIVE, *result.State)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful update", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "update-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "update-test-version",
			Description:       ptr.Of("Original description"),
			RegisteredModelId: *createdModel.Id,
		}

		created, err := service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update the version
		updateVersion := &openapi.ModelVersion{
			Id:                created.Id,
			Name:              "update-test-version", // Name should remain the same
			Description:       ptr.Of("Updated description"),
			Author:            ptr.Of("new-author"),
			State:             ptr.Of(openapi.MODELVERSIONSTATE_ARCHIVED),
			RegisteredModelId: *createdModel.Id,
		}

		updated, err := service.UpsertModelVersion(updateVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "update-test-version", updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "new-author", *updated.Author)
		assert.Equal(t, openapi.MODELVERSIONSTATE_ARCHIVED, *updated.State)
	})

	t.Run("create with custom properties", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "custom-props-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

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
			"epochs": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "100",
				},
			},
			"is_production": {
				MetadataBoolValue: &openapi.MetadataBoolValue{
					BoolValue: true,
				},
			},
		}

		inputVersion := &openapi.ModelVersion{
			Name:              "custom-props-version",
			RegisteredModelId: *createdModel.Id,
			CustomProperties:  &customProps,
		}

		result, err := service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-version", result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := *result.CustomProperties
		assert.Contains(t, resultProps, "accuracy")
		assert.Contains(t, resultProps, "framework")
		assert.Contains(t, resultProps, "epochs")
		assert.Contains(t, resultProps, "is_production")

		assert.Equal(t, 0.95, resultProps["accuracy"].MetadataDoubleValue.DoubleValue)
		assert.Equal(t, "tensorflow", resultProps["framework"].MetadataStringValue.StringValue)
		assert.Equal(t, "100", resultProps["epochs"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["is_production"].MetadataBoolValue.BoolValue)
	})

	t.Run("minimal version", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "minimal-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		inputVersion := &openapi.ModelVersion{
			Name:              "minimal-version",
			RegisteredModelId: *createdModel.Id,
		}

		result, err := service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-version", result.Name)
		assert.NotNil(t, result.Id)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
	})

	t.Run("nil model version error", func(t *testing.T) {
		result, err := service.UpsertModelVersion(nil, ptr.Of("1"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid model version pointer")
	})

	t.Run("nil fields preserved", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "nil-fields-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create model version with nil optional fields
		inputVersion := &openapi.ModelVersion{
			Name:              "nil-fields-version",
			Description:       nil, // Explicitly set to nil
			Author:            nil, // Explicitly set to nil
			ExternalId:        nil, // Explicitly set to nil
			State:             nil, // Explicitly set to nil
			RegisteredModelId: *createdModel.Id,
		}

		result, err := service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "nil-fields-version", result.Name)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
		assert.Nil(t, result.Description) // Verify description remains nil
		assert.Nil(t, result.Author)      // Verify author remains nil
		assert.Nil(t, result.ExternalId)  // Verify external ID remains nil
		assert.Nil(t, result.State)       // Verify state remains nil
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})
}

func TestGetModelVersionById(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "get-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "get-test-version",
			Description:       ptr.Of("Test description"),
			ExternalId:        ptr.Of("get-ext-123"),
			State:             ptr.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		created, err := service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the version by ID
		result, err := service.GetModelVersionById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-test-version", result.Name)
		assert.Equal(t, "get-ext-123", *result.ExternalId)
		assert.Equal(t, "Test description", *result.Description)
		assert.Equal(t, openapi.MODELVERSIONSTATE_LIVE, *result.State)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := service.GetModelVersionById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := service.GetModelVersionById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no model version found")
	})
}

func TestGetModelVersionByInferenceService(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get with specific model version", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "inference-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		modelVersion := &openapi.ModelVersion{
			Name:              "inference-test-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create a serving environment
		servingEnv := &openapi.ServingEnvironment{
			Name: "inference-test-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create an inference service with specific model version
		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			ModelVersionId:       createdVersion.Id,
		}
		createdInference, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Get model version by inference service
		result, err := service.GetModelVersionByInferenceService(*createdInference.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *createdVersion.Id, *result.Id)
		assert.Equal(t, "inference-test-version", result.Name)
	})

	t.Run("successful get with latest version", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "latest-version-test-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create multiple model versions
		version1 := &openapi.ModelVersion{
			Name:              "version-1",
			RegisteredModelId: *createdModel.Id,
		}
		_, err = service.UpsertModelVersion(version1, createdModel.Id)
		require.NoError(t, err)

		version2 := &openapi.ModelVersion{
			Name:              "version-2",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion2, err := service.UpsertModelVersion(version2, createdModel.Id)
		require.NoError(t, err)

		// Create a serving environment
		servingEnv := &openapi.ServingEnvironment{
			Name: "latest-version-test-env",
		}
		createdEnv, err := service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create an inference service without specific model version (should get latest)
		inferenceService := &openapi.InferenceService{
			Name:                 ptr.Of("latest-version-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			// ModelVersionId is nil, should get latest
		}
		createdInference, err := service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Get model version by inference service (should return latest)
		result, err := service.GetModelVersionByInferenceService(*createdInference.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		// Should be the latest version (version2)
		assert.Equal(t, *createdVersion2.Id, *result.Id)
		assert.Equal(t, "version-2", result.Name)
	})

	t.Run("invalid inference service id", func(t *testing.T) {
		result, err := service.GetModelVersionByInferenceService("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid inference service id")
	})

	t.Run("non-existent inference service", func(t *testing.T) {
		result, err := service.GetModelVersionByInferenceService("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetModelVersionByParams(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get by name and registered model id", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "params-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "params-test-version",
			ExternalId:        ptr.Of("params-ext-123"),
			RegisteredModelId: *createdModel.Id,
		}
		created, err := service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)

		// Get by name and registered model ID
		versionName := "params-test-version"
		result, err := service.GetModelVersionByParams(&versionName, createdModel.Id, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-test-version", result.Name)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
	})

	t.Run("successful get by external id", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "params-ext-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "params-ext-test-version",
			ExternalId:        ptr.Of("params-unique-ext-456"),
			RegisteredModelId: *createdModel.Id,
		}
		created, err := service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)

		// Get by external ID
		externalId := "params-unique-ext-456"
		result, err := service.GetModelVersionByParams(nil, nil, &externalId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-ext-test-version", result.Name)
		assert.Equal(t, "params-unique-ext-456", *result.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := service.GetModelVersionByParams(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters call")
	})

	t.Run("no version found", func(t *testing.T) {
		versionName := "nonexistent-version"
		registeredModelId := "999"
		result, err := service.GetModelVersionByParams(&versionName, &registeredModelId, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no model versions found")
	})
}

func TestGetModelVersions(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "list-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create multiple model versions
		versions := []*openapi.ModelVersion{
			{Name: "list-version-1", ExternalId: ptr.Of("list-ext-1"), RegisteredModelId: *createdModel.Id},
			{Name: "list-version-2", ExternalId: ptr.Of("list-ext-2"), RegisteredModelId: *createdModel.Id},
			{Name: "list-version-3", ExternalId: ptr.Of("list-ext-3"), RegisteredModelId: *createdModel.Id},
		}

		var createdIds []string
		for _, version := range versions {
			created, err := service.UpsertModelVersion(version, createdModel.Id)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List versions for the registered model
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := service.GetModelVersions(listOptions, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Len(t, result.Items, 3) // Should have exactly our 3 versions
		assert.Equal(t, int32(10), result.PageSize)

		// Verify our versions are in the result
		foundVersions := 0
		for _, item := range result.Items {
			for _, createdId := range createdIds {
				if *item.Id == createdId {
					foundVersions++
					break
				}
			}
		}
		assert.Equal(t, 3, foundVersions, "All created versions should be found in the list")
	})

	t.Run("list all versions without registered model filter", func(t *testing.T) {
		// Create multiple registered models with versions
		for i := 0; i < 2; i++ {
			registeredModel := &openapi.RegisteredModel{
				Name: "all-versions-model-" + string(rune('A'+i)),
			}
			createdModel, err := service.UpsertRegisteredModel(registeredModel)
			require.NoError(t, err)

			version := &openapi.ModelVersion{
				Name:              "all-versions-version-" + string(rune('A'+i)),
				RegisteredModelId: *createdModel.Id,
			}
			_, err = service.UpsertModelVersion(version, createdModel.Id)
			require.NoError(t, err)
		}

		// List all versions (no registered model filter)
		listOptions := api.ListOptions{}
		result, err := service.GetModelVersions(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least our 2 versions
	})

	t.Run("pagination and ordering", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "pagination-test-registered-model",
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create several versions for pagination testing
		for i := 0; i < 5; i++ {
			version := &openapi.ModelVersion{
				Name:              "pagination-version-" + string(rune('A'+i)),
				ExternalId:        ptr.Of("pagination-ext-" + string(rune('A'+i))),
				RegisteredModelId: *createdModel.Id,
			}
			_, err := service.UpsertModelVersion(version, createdModel.Id)
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

		result, err := service.GetModelVersions(listOptions, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})

	t.Run("invalid registered model id", func(t *testing.T) {
		invalidId := "invalid"
		listOptions := api.ListOptions{}

		result, err := service.GetModelVersions(listOptions, &invalidId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid registered model id")
	})
}

func TestModelVersionRoundTrip(t *testing.T) {
	service, cleanup := core.SetupModelRegistryService(t)
	defer cleanup()

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name:        "roundtrip-registered-model",
			Description: ptr.Of("Roundtrip test registered model"),
		}
		createdModel, err := service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version with all fields
		originalVersion := &openapi.ModelVersion{
			Name:              "roundtrip-version",
			Description:       ptr.Of("Roundtrip test description"),
			Author:            ptr.Of("roundtrip-author"),
			ExternalId:        ptr.Of("roundtrip-ext-123"),
			State:             ptr.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		// Create
		created, err := service.UpsertModelVersion(originalVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := service.GetModelVersionById(*created.Id)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, "roundtrip-version", retrieved.Name)
		assert.Equal(t, *originalVersion.Description, *retrieved.Description)
		assert.Equal(t, *originalVersion.Author, *retrieved.Author)
		assert.Equal(t, *originalVersion.ExternalId, *retrieved.ExternalId)
		assert.Equal(t, *originalVersion.State, *retrieved.State)
		assert.Equal(t, originalVersion.RegisteredModelId, retrieved.RegisteredModelId)

		// Update
		retrieved.Description = ptr.Of("Updated description")
		retrieved.State = ptr.Of(openapi.MODELVERSIONSTATE_ARCHIVED)

		updated, err := service.UpsertModelVersion(retrieved, createdModel.Id)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, openapi.MODELVERSIONSTATE_ARCHIVED, *updated.State)

		// Get again to verify persistence
		final, err := service.GetModelVersionById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, openapi.MODELVERSIONSTATE_ARCHIVED, *final.State)
	})
}
