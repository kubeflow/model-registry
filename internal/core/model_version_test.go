package core_test

import (
	"fmt"
	"strings"
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertModelVersion(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		// First create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name:        "test-registered-model",
			Description: apiutils.Of("Test registered model for version"),
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)
		require.NotNil(t, createdModel.Id)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "test-version",
			Description:       apiutils.Of("Test version description"),
			Author:            apiutils.Of("test-author"),
			ExternalId:        apiutils.Of("version-ext-123"),
			State:             apiutils.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

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
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "update-test-version",
			Description:       apiutils.Of("Original description"),
			RegisteredModelId: *createdModel.Id,
		}

		created, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update the version
		updateVersion := &openapi.ModelVersion{
			Id:                created.Id,
			Name:              "update-test-version", // Name should remain the same
			Description:       apiutils.Of("Updated description"),
			Author:            apiutils.Of("new-author"),
			State:             apiutils.Of(openapi.MODELVERSIONSTATE_ARCHIVED),
			RegisteredModelId: *createdModel.Id,
		}

		updated, err := _service.UpsertModelVersion(updateVersion, createdModel.Id)
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
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
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
			CustomProperties:  customProps,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-version", result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := result.CustomProperties
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
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		inputVersion := &openapi.ModelVersion{
			Name:              "minimal-version",
			RegisteredModelId: *createdModel.Id,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-version", result.Name)
		assert.NotNil(t, result.Id)
		assert.Equal(t, *createdModel.Id, result.RegisteredModelId)
	})

	t.Run("nil model version error", func(t *testing.T) {
		result, err := _service.UpsertModelVersion(nil, apiutils.Of("1"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid model version pointer")
	})

	t.Run("nil fields preserved", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "nil-fields-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
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

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

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

	t.Run("unicode characters in name", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "unicode-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Test with unicode characters: Chinese, Russian, Japanese, and emoji
		unicodeName := "模型版本-тест-モデルバージョン-🚀"
		inputVersion := &openapi.ModelVersion{
			Name:              unicodeName,
			Description:       apiutils.Of("Test model version with unicode characters"),
			Author:            apiutils.Of("测试作者-тестовый автор-テスト作者"),
			RegisteredModelId: *createdModel.Id,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, unicodeName, result.Name)
		assert.Equal(t, "Test model version with unicode characters", *result.Description)
		assert.Equal(t, "测试作者-тестовый автор-テスト作者", *result.Author)
		assert.NotNil(t, result.Id)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelVersionById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, unicodeName, retrieved.Name)
		assert.Equal(t, "测试作者-тестовый автор-テスト作者", *retrieved.Author)
	})

	t.Run("special characters in name", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "special-chars-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Test with various special characters
		specialName := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		inputVersion := &openapi.ModelVersion{
			Name:              specialName,
			Description:       apiutils.Of("Test model version with special characters"),
			Author:            apiutils.Of("author@#$%^&*()"),
			RegisteredModelId: *createdModel.Id,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, specialName, result.Name)
		assert.Equal(t, "Test model version with special characters", *result.Description)
		assert.Equal(t, "author@#$%^&*()", *result.Author)
		assert.NotNil(t, result.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelVersionById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, specialName, retrieved.Name)
		assert.Equal(t, "author@#$%^&*()", *retrieved.Author)
	})

	t.Run("mixed unicode and special characters", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "mixed-chars-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Test with mixed unicode and special characters
		mixedName := "模型@#$%版本-тест!@#-モデル()バージョン-🚀[]"
		inputVersion := &openapi.ModelVersion{
			Name:              mixedName,
			Description:       apiutils.Of("Test model version with mixed unicode and special characters"),
			Author:            apiutils.Of("作者@#$%-автор!@#-作者()🚀"),
			RegisteredModelId: *createdModel.Id,
		}

		result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, mixedName, result.Name)
		assert.Equal(t, "Test model version with mixed unicode and special characters", *result.Description)
		assert.Equal(t, "作者@#$%-автор!@#-作者()🚀", *result.Author)
		assert.NotNil(t, result.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelVersionById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, mixedName, retrieved.Name)
		assert.Equal(t, "作者@#$%-автор!@#-作者()🚀", *retrieved.Author)
	})

	t.Run("pagination with 10+ model versions", func(t *testing.T) {
		// Create a registered model first
		registeredModel := &openapi.RegisteredModel{
			Name: "paging-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create 15 model versions for pagination testing
		var createdVersions []string
		for i := 0; i < 15; i++ {
			versionName := "paging-test-model-version-" + fmt.Sprintf("%02d", i)
			inputVersion := &openapi.ModelVersion{
				Name:              versionName,
				Description:       apiutils.Of("Pagination test model version " + fmt.Sprintf("%02d", i)),
				Author:            apiutils.Of("test-author-" + fmt.Sprintf("%02d", i)),
				RegisteredModelId: *createdModel.Id,
			}

			result, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)
			require.NoError(t, err)
			createdVersions = append(createdVersions, *result.Id)
		}

		// Test pagination with page size 5
		pageSize := int32(5)
		orderBy := "name"
		sortOrder := "ASC"
		listOptions := api.ListOptions{
			PageSize:  &pageSize,
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		}

		// Get first page
		firstPage, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, firstPage)
		assert.LessOrEqual(t, len(firstPage.Items), 5, "First page should have at most 5 items")
		assert.Equal(t, int32(5), firstPage.PageSize)

		// Filter to only our test model versions in first page
		var firstPageTestVersions []openapi.ModelVersion
		firstPageIds := make(map[string]bool)
		for _, item := range firstPage.Items {
			// Only include our test versions (those with the specific prefix)
			if strings.HasPrefix(item.Name, "paging-test-model-version-") {
				assert.False(t, firstPageIds[*item.Id], "Should not have duplicate IDs in first page")
				firstPageIds[*item.Id] = true
				firstPageTestVersions = append(firstPageTestVersions, item)
			}
		}

		// Only proceed with second page test if we have a next page token and found test versions
		if firstPage.NextPageToken != "" && len(firstPageTestVersions) > 0 {
			// Get second page using next page token
			listOptions.NextPageToken = &firstPage.NextPageToken
			secondPage, err := _service.GetModelVersions(listOptions, createdModel.Id)
			require.NoError(t, err)
			require.NotNil(t, secondPage)
			assert.LessOrEqual(t, len(secondPage.Items), 5, "Second page should have at most 5 items")

			// Verify no duplicates between pages (only check our test versions)
			for _, item := range secondPage.Items {
				if strings.HasPrefix(item.Name, "paging-test-model-version-") {
					assert.False(t, firstPageIds[*item.Id], "Should not have duplicate IDs between pages")
				}
			}
		}

		// Test with larger page size
		largePage := int32(100)
		listOptions = api.ListOptions{
			PageSize:  &largePage,
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		}

		allItems, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, allItems)
		assert.GreaterOrEqual(t, len(allItems.Items), 15, "Should have at least our 15 test model versions")

		// Count our test versions in the results
		foundCount := 0
		for _, item := range allItems.Items {
			for _, createdId := range createdVersions {
				if *item.Id == createdId {
					foundCount++
					break
				}
			}
		}
		assert.Equal(t, 15, foundCount, "Should find all 15 created model versions")

		// Test descending order
		descOrder := "DESC"
		listOptions = api.ListOptions{
			PageSize:  &pageSize,
			OrderBy:   &orderBy,
			SortOrder: &descOrder,
		}

		descPage, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, descPage)
		assert.LessOrEqual(t, len(descPage.Items), 5, "Desc page should have at most 5 items")

		// Verify ordering (names should be in descending order)
		if len(descPage.Items) > 1 {
			for i := 1; i < len(descPage.Items); i++ {
				assert.GreaterOrEqual(t, descPage.Items[i-1].Name, descPage.Items[i].Name,
					"Items should be in descending order by name")
			}
		}
	})
}

func TestGetModelVersionById(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "get-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "get-test-version",
			Description:       apiutils.Of("Test description"),
			ExternalId:        apiutils.Of("get-ext-123"),
			State:             apiutils.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		created, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the version by ID
		result, err := _service.GetModelVersionById(*created.Id)

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
		result, err := _service.GetModelVersionById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := _service.GetModelVersionById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no model version found")
	})
}

func TestGetModelVersionByInferenceService(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get with specific model version", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "inference-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		modelVersion := &openapi.ModelVersion{
			Name:              "inference-test-version",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create a serving environment
		servingEnv := &openapi.ServingEnvironment{
			Name: "inference-test-env",
		}
		createdEnv, err := _service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create an inference service with specific model version
		inferenceService := &openapi.InferenceService{
			Name:                 apiutils.Of("test-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			ModelVersionId:       createdVersion.Id,
		}
		createdInference, err := _service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Get model version by inference service
		result, err := _service.GetModelVersionByInferenceService(*createdInference.Id)

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
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create multiple model versions
		version1 := &openapi.ModelVersion{
			Name:              "version-1",
			RegisteredModelId: *createdModel.Id,
		}
		_, err = _service.UpsertModelVersion(version1, createdModel.Id)
		require.NoError(t, err)

		version2 := &openapi.ModelVersion{
			Name:              "version-2",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion2, err := _service.UpsertModelVersion(version2, createdModel.Id)
		require.NoError(t, err)

		// Create a serving environment
		servingEnv := &openapi.ServingEnvironment{
			Name: "latest-version-test-env",
		}
		createdEnv, err := _service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		// Create an inference service without specific model version (should get latest)
		inferenceService := &openapi.InferenceService{
			Name:                 apiutils.Of("latest-version-inference-service"),
			ServingEnvironmentId: *createdEnv.Id,
			RegisteredModelId:    *createdModel.Id,
			// ModelVersionId is nil, should get latest
		}
		createdInference, err := _service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Get model version by inference service (should return latest)
		result, err := _service.GetModelVersionByInferenceService(*createdInference.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		// Should be the latest version (version2)
		assert.Equal(t, *createdVersion2.Id, *result.Id)
		assert.Equal(t, "version-2", result.Name)
	})

	t.Run("invalid inference service id", func(t *testing.T) {
		result, err := _service.GetModelVersionByInferenceService("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})

	t.Run("non-existent inference service", func(t *testing.T) {
		result, err := _service.GetModelVersionByInferenceService("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetModelVersionByParams(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get by name and registered model id", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "params-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "params-test-version",
			ExternalId:        apiutils.Of("params-ext-123"),
			RegisteredModelId: *createdModel.Id,
		}
		created, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)

		// Get by name and registered model ID
		versionName := "params-test-version"
		result, err := _service.GetModelVersionByParams(&versionName, createdModel.Id, nil)

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
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version
		inputVersion := &openapi.ModelVersion{
			Name:              "params-ext-test-version",
			ExternalId:        apiutils.Of("params-unique-ext-456"),
			RegisteredModelId: *createdModel.Id,
		}
		created, err := _service.UpsertModelVersion(inputVersion, createdModel.Id)
		require.NoError(t, err)

		// Get by external ID
		externalId := "params-unique-ext-456"
		result, err := _service.GetModelVersionByParams(nil, nil, &externalId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-ext-test-version", result.Name)
		assert.Equal(t, "params-unique-ext-456", *result.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := _service.GetModelVersionByParams(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters call")
	})

	t.Run("no version found", func(t *testing.T) {
		versionName := "nonexistent-version"
		registeredModelId := "999"
		result, err := _service.GetModelVersionByParams(&versionName, &registeredModelId, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no model versions found")
	})

	t.Run("same version name across different models", func(t *testing.T) {
		// This test catches the bug where ParentResourceID was not being used to filter model versions

		// Create first registered model
		registeredModel1 := &openapi.RegisteredModel{
			Name: "model-with-shared-version-1",
		}
		createdModel1, err := _service.UpsertRegisteredModel(registeredModel1)
		require.NoError(t, err)

		// Create second registered model
		registeredModel2 := &openapi.RegisteredModel{
			Name: "model-with-shared-version-2",
		}
		createdModel2, err := _service.UpsertRegisteredModel(registeredModel2)
		require.NoError(t, err)

		// Create version "shared-version-name-test" for the first model
		version1 := &openapi.ModelVersion{
			Name:              "shared-version-name-test",
			RegisteredModelId: *createdModel1.Id,
			Description:       apiutils.Of("Version for model 1"),
		}
		createdVersion1, err := _service.UpsertModelVersion(version1, createdModel1.Id)
		require.NoError(t, err)

		// Create version "shared-version-name-test" for the second model
		version2 := &openapi.ModelVersion{
			Name:              "shared-version-name-test",
			RegisteredModelId: *createdModel2.Id,
			Description:       apiutils.Of("Version for model 2"),
		}
		createdVersion2, err := _service.UpsertModelVersion(version2, createdModel2.Id)
		require.NoError(t, err)

		// Query for version "shared-version-name-test" of the first model
		versionName := "shared-version-name-test"
		result1, err := _service.GetModelVersionByParams(&versionName, createdModel1.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result1)
		assert.Equal(t, *createdVersion1.Id, *result1.Id)
		assert.Equal(t, *createdModel1.Id, result1.RegisteredModelId)
		assert.Equal(t, "Version for model 1", *result1.Description)

		// Query for version "shared-version-name-test" of the second model
		result2, err := _service.GetModelVersionByParams(&versionName, createdModel2.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result2)
		assert.Equal(t, *createdVersion2.Id, *result2.Id)
		assert.Equal(t, *createdModel2.Id, result2.RegisteredModelId)
		assert.Equal(t, "Version for model 2", *result2.Description)

		// Ensure we got different versions
		assert.NotEqual(t, *result1.Id, *result2.Id)
	})
}

func TestGetModelVersions(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "list-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create multiple model versions
		versions := []*openapi.ModelVersion{
			{Name: "list-version-1", ExternalId: apiutils.Of("list-ext-1"), RegisteredModelId: *createdModel.Id},
			{Name: "list-version-2", ExternalId: apiutils.Of("list-ext-2"), RegisteredModelId: *createdModel.Id},
			{Name: "list-version-3", ExternalId: apiutils.Of("list-ext-3"), RegisteredModelId: *createdModel.Id},
		}

		var createdIds []string
		for _, version := range versions {
			created, err := _service.UpsertModelVersion(version, createdModel.Id)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List versions for the registered model
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := _service.GetModelVersions(listOptions, createdModel.Id)

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
			createdModel, err := _service.UpsertRegisteredModel(registeredModel)
			require.NoError(t, err)

			version := &openapi.ModelVersion{
				Name:              "all-versions-version-" + string(rune('A'+i)),
				RegisteredModelId: *createdModel.Id,
			}
			_, err = _service.UpsertModelVersion(version, createdModel.Id)
			require.NoError(t, err)
		}

		// List all versions (no registered model filter)
		listOptions := api.ListOptions{}
		result, err := _service.GetModelVersions(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least our 2 versions
	})

	t.Run("pagination and ordering", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "pagination-test-registered-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create several versions for pagination testing
		for i := 0; i < 5; i++ {
			version := &openapi.ModelVersion{
				Name:              "pagination-version-" + string(rune('A'+i)),
				ExternalId:        apiutils.Of("pagination-ext-" + string(rune('A'+i))),
				RegisteredModelId: *createdModel.Id,
			}
			_, err := _service.UpsertModelVersion(version, createdModel.Id)
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

		result, err := _service.GetModelVersions(listOptions, createdModel.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})

	t.Run("invalid registered model id", func(t *testing.T) {
		invalidId := "invalid"
		listOptions := api.ListOptions{}

		result, err := _service.GetModelVersions(listOptions, &invalidId)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})
}

func TestModelVersionRoundTrip(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name:        "roundtrip-registered-model",
			Description: apiutils.Of("Roundtrip test registered model"),
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create a model version with all fields
		originalVersion := &openapi.ModelVersion{
			Name:              "roundtrip-version",
			Description:       apiutils.Of("Roundtrip test description"),
			Author:            apiutils.Of("roundtrip-author"),
			ExternalId:        apiutils.Of("roundtrip-ext-123"),
			State:             apiutils.Of(openapi.MODELVERSIONSTATE_LIVE),
			RegisteredModelId: *createdModel.Id,
		}

		// Create
		created, err := _service.UpsertModelVersion(originalVersion, createdModel.Id)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := _service.GetModelVersionById(*created.Id)
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
		retrieved.Description = apiutils.Of("Updated description")
		retrieved.State = apiutils.Of(openapi.MODELVERSIONSTATE_ARCHIVED)

		updated, err := _service.UpsertModelVersion(retrieved, createdModel.Id)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, openapi.MODELVERSIONSTATE_ARCHIVED, *updated.State)

		// Get again to verify persistence
		final, err := _service.GetModelVersionById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, openapi.MODELVERSIONSTATE_ARCHIVED, *final.State)
	})
}

func TestGetModelVersionsWithFilterQuery(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create a registered model to associate versions with
	registeredModel := &openapi.RegisteredModel{
		Name: "test-model-for-versions",
	}
	createdModel, err := _service.UpsertRegisteredModel(registeredModel)
	require.NoError(t, err)

	// Create test model versions with various properties for filtering
	testVersions := []struct {
		version *openapi.ModelVersion
	}{
		{
			version: &openapi.ModelVersion{
				Name:              "v1.0.0",
				Description:       apiutils.Of("Initial release version"),
				ExternalId:        apiutils.Of("ext-v1-001"),
				Author:            apiutils.Of("alice"),
				State:             (*openapi.ModelVersionState)(apiutils.Of("LIVE")),
				RegisteredModelId: *createdModel.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"stage": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "production",
							MetadataType: "MetadataStringValue",
						},
					},
					"accuracy": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.95,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"batch_size": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "32",
							MetadataType: "MetadataIntValue",
						},
					},
				},
			},
		},
		{
			version: &openapi.ModelVersion{
				Name:              "v2.0.0",
				Description:       apiutils.Of("Major update with improvements"),
				ExternalId:        apiutils.Of("ext-v2-002"),
				Author:            apiutils.Of("bob"),
				State:             (*openapi.ModelVersionState)(apiutils.Of("ARCHIVED")),
				RegisteredModelId: *createdModel.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"stage": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "staging",
							MetadataType: "MetadataStringValue",
						},
					},
					"accuracy": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.87,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"batch_size": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "64",
							MetadataType: "MetadataIntValue",
						},
					},
				},
			},
		},
		{
			version: &openapi.ModelVersion{
				Name:              "v1.1.0",
				Description:       apiutils.Of("Minor update with bug fixes"),
				ExternalId:        apiutils.Of("ext-v1-1-003"),
				Author:            apiutils.Of("alice"),
				RegisteredModelId: *createdModel.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"stage": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "production",
							MetadataType: "MetadataStringValue",
						},
					},
					"accuracy": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.96,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"batch_size": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "32",
							MetadataType: "MetadataIntValue",
						},
					},
					"experimental": {
						MetadataBoolValue: &openapi.MetadataBoolValue{
							BoolValue:    true,
							MetadataType: "MetadataBoolValue",
						},
					},
				},
			},
		},
		{
			version: &openapi.ModelVersion{
				Name:              "v3.0.0-beta",
				Description:       apiutils.Of("Beta version for testing"),
				ExternalId:        apiutils.Of("ext-v3-beta-004"),
				Author:            apiutils.Of("charlie"),
				RegisteredModelId: *createdModel.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"stage": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "development",
							MetadataType: "MetadataStringValue",
						},
					},
					"accuracy": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.92,
							MetadataType: "MetadataDoubleValue",
						},
					},
				},
			},
		},
	}

	// Create all test versions
	for _, tv := range testVersions {
		_, err := _service.UpsertModelVersion(tv.version, createdModel.Id)
		require.NoError(t, err)
	}

	testCases := []struct {
		name          string
		filterQuery   string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "Filter by exact name",
			filterQuery:   "name = 'v1.0.0'",
			expectedCount: 1,
			expectedNames: []string{"v1.0.0"},
		},
		{
			name:          "Filter by name pattern",
			filterQuery:   "name LIKE 'v1.%'",
			expectedCount: 2,
			expectedNames: []string{"v1.0.0", "v1.1.0"},
		},
		{
			name:          "Filter by description",
			filterQuery:   "description LIKE '%bug fixes%'",
			expectedCount: 1,
			expectedNames: []string{"v1.1.0"},
		},
		{
			name:          "Filter by author",
			filterQuery:   "author = 'alice'",
			expectedCount: 2,
			expectedNames: []string{"v1.0.0", "v1.1.0"},
		},
		{
			name:          "Filter by external ID",
			filterQuery:   "externalId = 'ext-v2-002'",
			expectedCount: 1,
			expectedNames: []string{"v2.0.0"},
		},
		{
			name:          "Filter by state",
			filterQuery:   "state = 'ARCHIVED'",
			expectedCount: 1,
			expectedNames: []string{"v2.0.0"},
		},
		{
			name:          "Filter by custom property - string",
			filterQuery:   "stage = 'production'",
			expectedCount: 2,
			expectedNames: []string{"v1.0.0", "v1.1.0"},
		},
		{
			name:          "OpenAPI example: Filter by name and state",
			filterQuery:   "name='v1.0.0' AND state='LIVE'",
			expectedCount: 1,
			expectedNames: []string{"v1.0.0"},
		},
		{
			name:          "Filter by custom property - numeric comparison",
			filterQuery:   "accuracy > 0.94",
			expectedCount: 2,
			expectedNames: []string{"v1.0.0", "v1.1.0"},
		},
		{
			name:          "Filter by custom property - integer",
			filterQuery:   "batch_size = 32",
			expectedCount: 2,
			expectedNames: []string{"v1.0.0", "v1.1.0"},
		},
		{
			name:          "Filter by custom property - boolean",
			filterQuery:   "experimental = true",
			expectedCount: 1,
			expectedNames: []string{"v1.1.0"},
		},
		{
			name:          "Complex filter with AND",
			filterQuery:   "stage = 'production' AND accuracy > 0.95",
			expectedCount: 1,
			expectedNames: []string{"v1.1.0"},
		},
		{
			name:          "Complex filter with OR",
			filterQuery:   "author = 'alice' OR author = 'charlie'",
			expectedCount: 3,
			expectedNames: []string{"v1.0.0", "v1.1.0", "v3.0.0-beta"},
		},
		{
			name:          "Complex filter with parentheses",
			filterQuery:   "(stage = 'production' OR stage = 'staging') AND accuracy < 0.95",
			expectedCount: 1,
			expectedNames: []string{"v2.0.0"},
		},
		{
			name:          "Case insensitive pattern matching",
			filterQuery:   "name ILIKE '%BETA%'",
			expectedCount: 1,
			expectedNames: []string{"v3.0.0-beta"},
		},
		{
			name:          "Filter with NOT condition",
			filterQuery:   "stage != 'development'",
			expectedCount: 3,
			expectedNames: []string{"v1.0.0", "v2.0.0", "v1.1.0"},
		},

		{
			name:          "Filter by name pattern with version suffix",
			filterQuery:   "name LIKE '%-beta'",
			expectedCount: 1,
			expectedNames: []string{"v3.0.0-beta"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pageSize := int32(10)
			listOptions := api.ListOptions{
				PageSize:    &pageSize,
				FilterQuery: &tc.filterQuery,
			}

			result, err := _service.GetModelVersions(listOptions, nil)

			require.NoError(t, err)
			require.NotNil(t, result)

			// Extract names from results
			var actualNames []string
			for _, item := range result.Items {
				for _, expectedName := range tc.expectedNames {
					if item.Name == expectedName {
						actualNames = append(actualNames, item.Name)
						break
					}
				}
			}

			assert.Equal(t, tc.expectedCount, len(actualNames),
				"Expected %d versions for filter '%s', but got %d",
				tc.expectedCount, tc.filterQuery, len(actualNames))

			// Verify the expected versions are present
			assert.ElementsMatch(t, tc.expectedNames, actualNames,
				"Expected versions %v for filter '%s', but got %v",
				tc.expectedNames, tc.filterQuery, actualNames)
		})
	}

	// Test error cases
	t.Run("Invalid filter syntax", func(t *testing.T) {
		pageSize := int32(10)
		invalidFilter := "invalid <<< syntax"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &invalidFilter,
		}

		result, err := _service.GetModelVersions(listOptions, nil)

		if assert.Error(t, err) {
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid filter query")
		}
	})

	// Test combining filterQuery with registeredModelId parameter
	t.Run("Filter with registeredModelId parameter", func(t *testing.T) {
		// Create another registered model with versions
		anotherModel := &openapi.RegisteredModel{
			Name: "another-model",
		}
		anotherCreatedModel, err := _service.UpsertRegisteredModel(anotherModel)
		require.NoError(t, err)

		anotherVersion := &openapi.ModelVersion{
			Name:              "v1.0.0",
			RegisteredModelId: *anotherCreatedModel.Id,
			CustomProperties: map[string]openapi.MetadataValue{
				"stage": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "production",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}
		_, err = _service.UpsertModelVersion(anotherVersion, anotherCreatedModel.Id)
		require.NoError(t, err)

		// Filter by stage=production should return versions from both models
		pageSize := int32(10)
		filterQuery := "stage = 'production'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		// Without registeredModelId - should get 3 (2 from first model + 1 from second)
		allResult, err := _service.GetModelVersions(listOptions, nil)
		require.NoError(t, err)
		assert.Equal(t, 3, len(allResult.Items))

		// With registeredModelId - should only get 2 from first model
		filteredResult, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		assert.Equal(t, 2, len(filteredResult.Items))
		for _, item := range filteredResult.Items {
			assert.Equal(t, *createdModel.Id, item.RegisteredModelId)
		}
	})

	// Test combining filterQuery with pagination
	t.Run("Filter with pagination", func(t *testing.T) {
		pageSize := int32(1)
		filterQuery := "stage = 'production'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		// Get first page
		firstPage, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(firstPage.Items))
		assert.NotEmpty(t, firstPage.NextPageToken)

		// Get second page
		listOptions.NextPageToken = &firstPage.NextPageToken
		secondPage, err := _service.GetModelVersions(listOptions, createdModel.Id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(secondPage.Items))

		// Ensure different items on each page
		assert.NotEqual(t, firstPage.Items[0].Id, secondPage.Items[0].Id)
	})

	// Test empty results
	t.Run("Filter with no matches", func(t *testing.T) {
		pageSize := int32(10)
		filterQuery := "stage = 'nonexistent'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		result, err := _service.GetModelVersions(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
	})
}
