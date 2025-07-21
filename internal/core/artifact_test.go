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

func TestUpsertArtifact(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful create model artifact", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of("test-model-artifact"),
			Description:        apiutils.Of("Test model artifact description"),
			ExternalId:         apiutils.Of("model-ext-123"),
			Uri:                apiutils.Of("s3://bucket/model.pkl"),
			State:              apiutils.Of(openapi.ARTIFACTSTATE_LIVE),
			ModelFormatName:    apiutils.Of("pickle"),
			ModelFormatVersion: apiutils.Of("1.0"),
			StorageKey:         apiutils.Of("model-storage-key"),
			StoragePath:        apiutils.Of("/models/test"),
			ServiceAccountName: apiutils.Of("model-sa"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.NotNil(t, result.ModelArtifact.Id)
		assert.Equal(t, "test-model-artifact", *result.ModelArtifact.Name)
		assert.Equal(t, "model-ext-123", *result.ModelArtifact.ExternalId)
		assert.Equal(t, "s3://bucket/model.pkl", *result.ModelArtifact.Uri)
		assert.Equal(t, openapi.ARTIFACTSTATE_LIVE, *result.ModelArtifact.State)
		assert.Equal(t, "pickle", *result.ModelArtifact.ModelFormatName)
		assert.Equal(t, "1.0", *result.ModelArtifact.ModelFormatVersion)
		assert.Equal(t, "model-storage-key", *result.ModelArtifact.StorageKey)
		assert.Equal(t, "/models/test", *result.ModelArtifact.StoragePath)
		assert.Equal(t, "model-sa", *result.ModelArtifact.ServiceAccountName)
		assert.NotNil(t, result.ModelArtifact.CreateTimeSinceEpoch)
		assert.NotNil(t, result.ModelArtifact.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful create doc artifact", func(t *testing.T) {
		docArtifact := &openapi.DocArtifact{
			Name:        apiutils.Of("test-doc-artifact"),
			Description: apiutils.Of("Test doc artifact description"),
			ExternalId:  apiutils.Of("doc-ext-123"),
			Uri:         apiutils.Of("s3://bucket/doc.pdf"),
			State:       apiutils.Of(openapi.ARTIFACTSTATE_LIVE),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.DocArtifact)
		assert.NotNil(t, result.DocArtifact.Id)
		assert.Equal(t, "test-doc-artifact", *result.DocArtifact.Name)
		assert.Equal(t, "doc-ext-123", *result.DocArtifact.ExternalId)
		assert.Equal(t, "s3://bucket/doc.pdf", *result.DocArtifact.Uri)
		assert.Equal(t, openapi.ARTIFACTSTATE_LIVE, *result.DocArtifact.State)
		assert.NotNil(t, result.DocArtifact.CreateTimeSinceEpoch)
		assert.NotNil(t, result.DocArtifact.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful update model artifact", func(t *testing.T) {
		// Create first
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("update-model-artifact"),
			Uri:  apiutils.Of("s3://bucket/original.pkl"),
		}

		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update by modifying the created artifact
		created.Uri = apiutils.Of("s3://bucket/updated.pkl")
		created.Description = apiutils.Of("Updated description")

		updated, err := _service.UpsertModelArtifact(created)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "s3://bucket/updated.pkl", *updated.Uri)
		assert.Equal(t, "Updated description", *updated.Description)
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
		}

		modelArtifact := &openapi.ModelArtifact{
			Name:             apiutils.Of("custom-props-artifact"),
			CustomProperties: &customProps,
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result.ModelArtifact)
		assert.NotNil(t, result.ModelArtifact.CustomProperties)

		resultProps := *result.ModelArtifact.CustomProperties
		assert.Contains(t, resultProps, "accuracy")
		assert.Contains(t, resultProps, "framework")
		assert.Equal(t, 0.95, resultProps["accuracy"].MetadataDoubleValue.DoubleValue)
		assert.Equal(t, "tensorflow", resultProps["framework"].MetadataStringValue.StringValue)
	})

	t.Run("nil artifact error", func(t *testing.T) {
		result, err := _service.UpsertArtifact(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid artifact pointer")
	})

	t.Run("invalid artifact type", func(t *testing.T) {
		artifact := &openapi.Artifact{}

		result, err := _service.UpsertArtifact(artifact)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid artifact type, must be either ModelArtifact or DocArtifact")
	})

	t.Run("unicode characters in model artifact name", func(t *testing.T) {
		// Test with unicode characters: Chinese, Russian, Japanese, and emoji
		unicodeName := "模型工件-тест-モデルアーティファクト-🚀"
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of(unicodeName),
			Description: apiutils.Of("Test model artifact with unicode characters"),
			Uri:         apiutils.Of("s3://bucket/unicode-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Equal(t, unicodeName, *result.ModelArtifact.Name)
		assert.Equal(t, "Test model artifact with unicode characters", *result.ModelArtifact.Description)
		assert.Equal(t, "s3://bucket/unicode-model.pkl", *result.ModelArtifact.Uri)
		assert.NotNil(t, result.ModelArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Equal(t, unicodeName, *retrieved.ModelArtifact.Name)
	})

	t.Run("special characters in model artifact name", func(t *testing.T) {
		// Test with various special characters
		specialName := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of(specialName),
			Description: apiutils.Of("Test model artifact with special characters"),
			Uri:         apiutils.Of("s3://bucket/special-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Equal(t, specialName, *result.ModelArtifact.Name)
		assert.Equal(t, "Test model artifact with special characters", *result.ModelArtifact.Description)
		assert.NotNil(t, result.ModelArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Equal(t, specialName, *retrieved.ModelArtifact.Name)
	})

	t.Run("mixed unicode and special characters in doc artifact", func(t *testing.T) {
		// Test with mixed unicode and special characters
		mixedName := "文档@#$%工件-тест!@#-ドキュメント()アーティファクト-🚀[]"
		docArtifact := &openapi.DocArtifact{
			Name:        apiutils.Of(mixedName),
			Description: apiutils.Of("Test doc artifact with mixed unicode and special characters"),
			Uri:         apiutils.Of("s3://bucket/mixed-doc.pdf"),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.DocArtifact)
		assert.Equal(t, mixedName, *result.DocArtifact.Name)
		assert.Equal(t, "Test doc artifact with mixed unicode and special characters", *result.DocArtifact.Description)
		assert.NotNil(t, result.DocArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.DocArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.DocArtifact)
		assert.Equal(t, mixedName, *retrieved.DocArtifact.Name)
	})

	t.Run("pagination with 10+ artifacts", func(t *testing.T) {
		// Create 15 artifacts for pagination testing
		var createdArtifacts []string
		for i := 0; i < 15; i++ {
			artifactName := "paging-test-artifact-" + fmt.Sprintf("%02d", i)
			modelArtifact := &openapi.ModelArtifact{
				Name:        apiutils.Of(artifactName),
				Description: apiutils.Of("Pagination test artifact " + fmt.Sprintf("%02d", i)),
				Uri:         apiutils.Of("s3://bucket/paging-test-" + fmt.Sprintf("%02d", i) + ".pkl"),
			}

			artifact := &openapi.Artifact{
				ModelArtifact: modelArtifact,
			}

			result, err := _service.UpsertArtifact(artifact)
			require.NoError(t, err)
			createdArtifacts = append(createdArtifacts, *result.ModelArtifact.Id)
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
		firstPage, err := _service.GetArtifacts(listOptions, nil)
		require.NoError(t, err)
		require.NotNil(t, firstPage)
		assert.LessOrEqual(t, len(firstPage.Items), 5, "First page should have at most 5 items")
		assert.Equal(t, int32(5), firstPage.PageSize)

		// Filter to only our test artifacts in first page
		var firstPageTestArtifacts []openapi.Artifact
		firstPageIds := make(map[string]bool)
		for _, item := range firstPage.Items {
			// Only include our test artifacts (those with the specific prefix)
			var artifactName string
			if item.ModelArtifact != nil {
				artifactName = *item.ModelArtifact.Name
			} else if item.DocArtifact != nil {
				artifactName = *item.DocArtifact.Name
			}

			if strings.HasPrefix(artifactName, "paging-test-artifact-") {
				var artifactId string
				if item.ModelArtifact != nil {
					artifactId = *item.ModelArtifact.Id
				} else if item.DocArtifact != nil {
					artifactId = *item.DocArtifact.Id
				}
				assert.False(t, firstPageIds[artifactId], "Should not have duplicate IDs in first page")
				firstPageIds[artifactId] = true
				firstPageTestArtifacts = append(firstPageTestArtifacts, item)
			}
		}

		// Only proceed with second page test if we have a next page token and found test artifacts
		if firstPage.NextPageToken != "" && len(firstPageTestArtifacts) > 0 {
			// Get second page using next page token
			listOptions.NextPageToken = &firstPage.NextPageToken
			secondPage, err := _service.GetArtifacts(listOptions, nil)
			require.NoError(t, err)
			require.NotNil(t, secondPage)
			assert.LessOrEqual(t, len(secondPage.Items), 5, "Second page should have at most 5 items")

			// Verify no duplicates between pages (only check our test artifacts)
			for _, item := range secondPage.Items {
				var artifactName, artifactId string
				if item.ModelArtifact != nil {
					artifactName = *item.ModelArtifact.Name
					artifactId = *item.ModelArtifact.Id
				} else if item.DocArtifact != nil {
					artifactName = *item.DocArtifact.Name
					artifactId = *item.DocArtifact.Id
				}

				if strings.HasPrefix(artifactName, "paging-test-artifact-") {
					assert.False(t, firstPageIds[artifactId], "Should not have duplicate IDs between pages")
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

		allItems, err := _service.GetArtifacts(listOptions, nil)
		require.NoError(t, err)
		require.NotNil(t, allItems)
		assert.GreaterOrEqual(t, len(allItems.Items), 15, "Should have at least our 15 test artifacts")

		// Count our test artifacts in the results
		foundCount := 0
		for _, item := range allItems.Items {
			var artifactId string
			if item.ModelArtifact != nil {
				artifactId = *item.ModelArtifact.Id
			} else if item.DocArtifact != nil {
				artifactId = *item.DocArtifact.Id
			}

			for _, createdId := range createdArtifacts {
				if artifactId == createdId {
					foundCount++
					break
				}
			}
		}
		assert.Equal(t, 15, foundCount, "Should find all 15 created artifacts")

		// Test descending order
		descOrder := "DESC"
		listOptions = api.ListOptions{
			PageSize:  &pageSize,
			OrderBy:   &orderBy,
			SortOrder: &descOrder,
		}

		descPage, err := _service.GetArtifacts(listOptions, nil)
		require.NoError(t, err)
		require.NotNil(t, descPage)
		assert.LessOrEqual(t, len(descPage.Items), 5, "Desc page should have at most 5 items")

		// Verify ordering (names should be in descending order)
		if len(descPage.Items) > 1 {
			for i := 1; i < len(descPage.Items); i++ {
				var prevName, currName string
				if descPage.Items[i-1].ModelArtifact != nil {
					prevName = *descPage.Items[i-1].ModelArtifact.Name
				} else if descPage.Items[i-1].DocArtifact != nil {
					prevName = *descPage.Items[i-1].DocArtifact.Name
				}
				if descPage.Items[i].ModelArtifact != nil {
					currName = *descPage.Items[i].ModelArtifact.Name
				} else if descPage.Items[i].DocArtifact != nil {
					currName = *descPage.Items[i].DocArtifact.Name
				}
				assert.GreaterOrEqual(t, prevName, currName,
					"Items should be in descending order by name")
			}
		}
	})
}

func TestUpsertModelVersionArtifact(t *testing.T) {
	cleanupTestData(t, sharedDB)
	t.Run("successful create with model version", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "test-model-for-artifact",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:        "v1.0",
			Description: apiutils.Of("Version 1.0"),
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create artifact associated with model version
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("version-artifact"),
			Uri:  apiutils.Of("s3://bucket/version-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.NotNil(t, result.ModelArtifact.Id)
		// Name should be prefixed with model version ID
		assert.Contains(t, *result.ModelArtifact.Name, "version-artifact")
		assert.Equal(t, "s3://bucket/version-model.pkl", *result.ModelArtifact.Uri)
	})

	t.Run("invalid model version id", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("test-artifact"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertModelVersionArtifact(artifact, "invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})

	t.Run("unicode characters in model version artifact name", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "unicode-test-model-for-artifact",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0-unicode",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Test with unicode characters: Chinese, Russian, Japanese, and emoji
		unicodeName := "版本工件-тест-バージョンアーティファクト-🚀"
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of(unicodeName),
			Description: apiutils.Of("Test model version artifact with unicode characters"),
			Uri:         apiutils.Of("s3://bucket/unicode-version-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Contains(t, *result.ModelArtifact.Name, unicodeName)
		assert.Equal(t, "Test model version artifact with unicode characters", *result.ModelArtifact.Description)
		assert.Equal(t, "s3://bucket/unicode-version-model.pkl", *result.ModelArtifact.Uri)
		assert.NotNil(t, result.ModelArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Contains(t, *retrieved.ModelArtifact.Name, unicodeName)
	})

	t.Run("special characters in model version artifact name", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "special-chars-test-model-for-artifact",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0-special",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Test with various special characters
		specialName := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of(specialName),
			Description: apiutils.Of("Test model version artifact with special characters"),
			Uri:         apiutils.Of("s3://bucket/special-version-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Contains(t, *result.ModelArtifact.Name, specialName)
		assert.Equal(t, "Test model version artifact with special characters", *result.ModelArtifact.Description)
		assert.NotNil(t, result.ModelArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Contains(t, *retrieved.ModelArtifact.Name, specialName)
	})

	t.Run("mixed unicode and special characters in model version artifact", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "mixed-chars-test-model-for-artifact",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0-mixed",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Test with mixed unicode and special characters
		mixedName := "版本@#$%工件-тест!@#-バージョン()アーティファクト-🚀[]"
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of(mixedName),
			Description: apiutils.Of("Test model version artifact with mixed unicode and special characters"),
			Uri:         apiutils.Of("s3://bucket/mixed-version-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Contains(t, *result.ModelArtifact.Name, mixedName)
		assert.Equal(t, "Test model version artifact with mixed unicode and special characters", *result.ModelArtifact.Description)
		assert.NotNil(t, result.ModelArtifact.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetArtifactById(*result.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Contains(t, *retrieved.ModelArtifact.Name, mixedName)
	})

	t.Run("pagination with 10+ model version artifacts", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "paging-test-model-for-artifacts",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0-paging",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create 15 model version artifacts for pagination testing
		var createdArtifacts []string
		for i := 0; i < 15; i++ {
			artifactName := "paging-test-version-artifact-" + fmt.Sprintf("%02d", i)
			modelArtifact := &openapi.ModelArtifact{
				Name:        apiutils.Of(artifactName),
				Description: apiutils.Of("Pagination test model version artifact " + fmt.Sprintf("%02d", i)),
				Uri:         apiutils.Of("s3://bucket/paging-version-test-" + fmt.Sprintf("%02d", i) + ".pkl"),
			}

			artifact := &openapi.Artifact{
				ModelArtifact: modelArtifact,
			}

			result, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
			require.NoError(t, err)
			createdArtifacts = append(createdArtifacts, *result.ModelArtifact.Id)
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
		firstPage, err := _service.GetArtifacts(listOptions, createdVersion.Id)
		require.NoError(t, err)
		require.NotNil(t, firstPage)
		assert.LessOrEqual(t, len(firstPage.Items), 5, "First page should have at most 5 items")
		assert.Equal(t, int32(5), firstPage.PageSize)

		// Filter to only our test artifacts in first page
		var firstPageTestArtifacts []openapi.Artifact
		firstPageIds := make(map[string]bool)
		for _, item := range firstPage.Items {
			// Only include our test artifacts (those with the specific prefix)
			var artifactName string
			if item.ModelArtifact != nil {
				artifactName = *item.ModelArtifact.Name
			}

			if strings.Contains(artifactName, "paging-test-version-artifact-") {
				artifactId := *item.ModelArtifact.Id
				assert.False(t, firstPageIds[artifactId], "Should not have duplicate IDs in first page")
				firstPageIds[artifactId] = true
				firstPageTestArtifacts = append(firstPageTestArtifacts, item)
			}
		}

		// Only proceed with second page test if we have a next page token and found test artifacts
		if firstPage.NextPageToken != "" && len(firstPageTestArtifacts) > 0 {
			// Get second page using next page token
			listOptions.NextPageToken = &firstPage.NextPageToken
			secondPage, err := _service.GetArtifacts(listOptions, createdVersion.Id)
			require.NoError(t, err)
			require.NotNil(t, secondPage)
			assert.LessOrEqual(t, len(secondPage.Items), 5, "Second page should have at most 5 items")

			// Verify no duplicates between pages (only check our test artifacts)
			for _, item := range secondPage.Items {
				if item.ModelArtifact != nil && strings.Contains(*item.ModelArtifact.Name, "paging-test-version-artifact-") {
					assert.False(t, firstPageIds[*item.ModelArtifact.Id], "Should not have duplicate IDs between pages")
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

		allItems, err := _service.GetArtifacts(listOptions, createdVersion.Id)
		require.NoError(t, err)
		require.NotNil(t, allItems)
		assert.GreaterOrEqual(t, len(allItems.Items), 15, "Should have at least our 15 test artifacts")

		// Count our test artifacts in the results
		foundCount := 0
		for _, item := range allItems.Items {
			if item.ModelArtifact != nil {
				for _, createdId := range createdArtifacts {
					if *item.ModelArtifact.Id == createdId {
						foundCount++
						break
					}
				}
			}
		}
		assert.Equal(t, 15, foundCount, "Should find all 15 created model version artifacts")

		// Test descending order
		descOrder := "DESC"
		listOptions = api.ListOptions{
			PageSize:  &pageSize,
			OrderBy:   &orderBy,
			SortOrder: &descOrder,
		}

		descPage, err := _service.GetArtifacts(listOptions, createdVersion.Id)
		require.NoError(t, err)
		require.NotNil(t, descPage)
		assert.LessOrEqual(t, len(descPage.Items), 5, "Desc page should have at most 5 items")

		// Verify ordering (names should be in descending order)
		if len(descPage.Items) > 1 {
			for i := 1; i < len(descPage.Items); i++ {
				var prevName, currName string
				if descPage.Items[i-1].ModelArtifact != nil {
					prevName = *descPage.Items[i-1].ModelArtifact.Name
				}
				if descPage.Items[i].ModelArtifact != nil {
					currName = *descPage.Items[i].ModelArtifact.Name
				}
				if prevName != "" && currName != "" {
					assert.GreaterOrEqual(t, prevName, currName,
						"Items should be in descending order by name")
				}
			}
		}
	})
}

func TestGetArtifactById(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful get model artifact", func(t *testing.T) {
		// Create a model artifact first
		modelArtifact := &openapi.ModelArtifact{
			Name:        apiutils.Of("get-test-model-artifact"),
			Description: apiutils.Of("Test description"),
			Uri:         apiutils.Of("s3://bucket/test.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		created, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)
		require.NotNil(t, created.ModelArtifact.Id)

		// Get the artifact by ID
		result, err := _service.GetArtifactById(*created.ModelArtifact.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Equal(t, *created.ModelArtifact.Id, *result.ModelArtifact.Id)
		assert.Equal(t, "get-test-model-artifact", *result.ModelArtifact.Name)
		assert.Equal(t, "Test description", *result.ModelArtifact.Description)
		assert.Equal(t, "s3://bucket/test.pkl", *result.ModelArtifact.Uri)
	})

	t.Run("successful get doc artifact", func(t *testing.T) {
		// Create a doc artifact first
		docArtifact := &openapi.DocArtifact{
			Name:        apiutils.Of("get-test-doc-artifact"),
			Description: apiutils.Of("Test doc description"),
			Uri:         apiutils.Of("s3://bucket/test.pdf"),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		created, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)
		require.NotNil(t, created.DocArtifact.Id)

		// Get the artifact by ID
		result, err := _service.GetArtifactById(*created.DocArtifact.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.DocArtifact)
		assert.Equal(t, *created.DocArtifact.Id, *result.DocArtifact.Id)
		assert.Equal(t, "get-test-doc-artifact", *result.DocArtifact.Name)
		assert.Equal(t, "Test doc description", *result.DocArtifact.Description)
		assert.Equal(t, "s3://bucket/test.pdf", *result.DocArtifact.Uri)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := _service.GetArtifactById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := _service.GetArtifactById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetArtifactByParams(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful get by name and model version", func(t *testing.T) {
		// Create registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "test-model-for-params",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create artifact with model version
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("params-test-artifact"),
			Uri:  apiutils.Of("s3://bucket/params-test.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		created, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
		require.NoError(t, err)

		// Get by name and model version ID
		result, err := _service.GetArtifactByParams(apiutils.Of("params-test-artifact"), createdVersion.Id, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Equal(t, *created.ModelArtifact.Id, *result.ModelArtifact.Id)
	})

	t.Run("successful get by external id", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			Name:       apiutils.Of("external-id-artifact"),
			ExternalId: apiutils.Of("ext-params-123"),
			Uri:        apiutils.Of("s3://bucket/external.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		created, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)

		// Get by external ID
		result, err := _service.GetArtifactByParams(nil, nil, apiutils.Of("ext-params-123"))

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.Equal(t, *created.ModelArtifact.Id, *result.ModelArtifact.Id)
		assert.Equal(t, "ext-params-123", *result.ModelArtifact.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := _service.GetArtifactByParams(nil, nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters")
	})

	t.Run("artifact not found", func(t *testing.T) {
		result, err := _service.GetArtifactByParams(nil, nil, apiutils.Of("non-existent"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no artifacts found")
	})
}

func TestGetArtifacts(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful list all artifacts", func(t *testing.T) {
		// Create multiple artifacts
		artifacts := []*openapi.Artifact{
			{
				ModelArtifact: &openapi.ModelArtifact{
					Name: apiutils.Of("list-artifact-1"),
					Uri:  apiutils.Of("s3://bucket/artifact1.pkl"),
				},
			},
			{
				ModelArtifact: &openapi.ModelArtifact{
					Name: apiutils.Of("list-artifact-2"),
					Uri:  apiutils.Of("s3://bucket/artifact2.pkl"),
				},
			},
			{
				DocArtifact: &openapi.DocArtifact{
					Name: apiutils.Of("list-doc-artifact"),
					Uri:  apiutils.Of("s3://bucket/doc.pdf"),
				},
			},
		}

		for _, artifact := range artifacts {
			_, err := _service.UpsertArtifact(artifact)
			require.NoError(t, err)
		}

		// List all artifacts
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(10)),
		}

		result, err := _service.GetArtifacts(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3)
		assert.NotNil(t, result.Size)
		assert.Equal(t, int32(10), result.PageSize)
	})

	t.Run("successful list artifacts by model version", func(t *testing.T) {
		// Create registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "test-model-for-list",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create artifacts for this model version
		for i := 0; i < 3; i++ {
			artifact := &openapi.Artifact{
				ModelArtifact: &openapi.ModelArtifact{
					Name: apiutils.Of("version-artifact-" + string(rune('1'+i))),
					Uri:  apiutils.Of("s3://bucket/version" + string(rune('1'+i)) + ".pkl"),
				},
			}
			_, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
			require.NoError(t, err)
		}

		// List artifacts for this model version
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(10)),
		}

		result, err := _service.GetArtifacts(listOptions, createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, len(result.Items))
	})

	t.Run("invalid model version id", func(t *testing.T) {
		listOptions := api.ListOptions{}

		result, err := _service.GetArtifacts(listOptions, apiutils.Of("invalid"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})
}

func TestUpsertModelArtifact(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful create", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of("direct-model-artifact"),
			Description:        apiutils.Of("Direct model artifact"),
			Uri:                apiutils.Of("s3://bucket/direct.pkl"),
			ModelFormatName:    apiutils.Of("tensorflow"),
			ModelFormatVersion: apiutils.Of("2.8"),
		}

		result, err := _service.UpsertModelArtifact(modelArtifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "direct-model-artifact", *result.Name)
		assert.Equal(t, "Direct model artifact", *result.Description)
		assert.Equal(t, "s3://bucket/direct.pkl", *result.Uri)
		assert.Equal(t, "tensorflow", *result.ModelFormatName)
		assert.Equal(t, "2.8", *result.ModelFormatVersion)
	})

	t.Run("nil model artifact error", func(t *testing.T) {
		result, err := _service.UpsertModelArtifact(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid model artifact pointer")
	})

	t.Run("unicode characters in model artifact name", func(t *testing.T) {
		// Test with unicode characters: Chinese, Russian, Japanese, and emoji
		unicodeName := "直接模型工件-тест-ダイレクトモデルアーティファクト-🚀"
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of(unicodeName),
			Description:        apiutils.Of("Direct model artifact with unicode characters"),
			Uri:                apiutils.Of("s3://bucket/unicode-direct.pkl"),
			ModelFormatName:    apiutils.Of("tensorflow-unicode"),
			ModelFormatVersion: apiutils.Of("2.8-测试"),
		}

		result, err := _service.UpsertModelArtifact(modelArtifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, unicodeName, *result.Name)
		assert.Equal(t, "Direct model artifact with unicode characters", *result.Description)
		assert.Equal(t, "s3://bucket/unicode-direct.pkl", *result.Uri)
		assert.Equal(t, "tensorflow-unicode", *result.ModelFormatName)
		assert.Equal(t, "2.8-测试", *result.ModelFormatVersion)
		assert.NotNil(t, result.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelArtifactById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, unicodeName, *retrieved.Name)
		assert.Equal(t, "2.8-测试", *retrieved.ModelFormatVersion)
	})

	t.Run("special characters in model artifact name", func(t *testing.T) {
		// Test with various special characters
		specialName := "!@#$%^&*()_+-=[]{}|;':\",./<>?"
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of(specialName),
			Description:        apiutils.Of("Direct model artifact with special characters"),
			Uri:                apiutils.Of("s3://bucket/special-direct.pkl"),
			ModelFormatName:    apiutils.Of("format@#$%"),
			ModelFormatVersion: apiutils.Of("1.0!@#"),
		}

		result, err := _service.UpsertModelArtifact(modelArtifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, specialName, *result.Name)
		assert.Equal(t, "Direct model artifact with special characters", *result.Description)
		assert.Equal(t, "s3://bucket/special-direct.pkl", *result.Uri)
		assert.Equal(t, "format@#$%", *result.ModelFormatName)
		assert.Equal(t, "1.0!@#", *result.ModelFormatVersion)
		assert.NotNil(t, result.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelArtifactById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, specialName, *retrieved.Name)
		assert.Equal(t, "format@#$%", *retrieved.ModelFormatName)
	})

	t.Run("mixed unicode and special characters in model artifact", func(t *testing.T) {
		// Test with mixed unicode and special characters
		mixedName := "直接@#$%模型-тест!@#-ダイレクト()モデル-🚀[]"
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of(mixedName),
			Description:        apiutils.Of("Direct model artifact with mixed unicode and special characters"),
			Uri:                apiutils.Of("s3://bucket/mixed-direct.pkl"),
			ModelFormatName:    apiutils.Of("tensorflow@#$%-测试"),
			ModelFormatVersion: apiutils.Of("2.8!@#-тест"),
		}

		result, err := _service.UpsertModelArtifact(modelArtifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, mixedName, *result.Name)
		assert.Equal(t, "Direct model artifact with mixed unicode and special characters", *result.Description)
		assert.Equal(t, "s3://bucket/mixed-direct.pkl", *result.Uri)
		assert.Equal(t, "tensorflow@#$%-测试", *result.ModelFormatName)
		assert.Equal(t, "2.8!@#-тест", *result.ModelFormatVersion)
		assert.NotNil(t, result.Id)

		// Verify we can retrieve it by ID
		retrieved, err := _service.GetModelArtifactById(*result.Id)
		require.NoError(t, err)
		assert.Equal(t, mixedName, *retrieved.Name)
		assert.Equal(t, "tensorflow@#$%-测试", *retrieved.ModelFormatName)
	})

	t.Run("pagination test", func(t *testing.T) {
		// Create multiple model artifacts for pagination testing
		for i := 0; i < 15; i++ {
			modelArtifact := &openapi.ModelArtifact{
				Name: apiutils.Of(fmt.Sprintf("paging-test-direct-model-artifact-%d", i+1)),
				Uri:  apiutils.Of(fmt.Sprintf("s3://bucket/paging-direct-model-%d.pkl", i+1)),
			}

			result, err := _service.UpsertModelArtifact(modelArtifact)
			require.NoError(t, err)
			require.NotNil(t, result.Id)
		}

		// Test pagination with page size 5
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(5)),
		}

		// Get first page
		firstPage, err := _service.GetModelArtifacts(listOptions, nil)
		require.NoError(t, err)
		require.NotNil(t, firstPage)
		assert.Equal(t, 5, len(firstPage.Items))
		assert.NotNil(t, firstPage.NextPageToken)

		// Get second page
		listOptions.NextPageToken = apiutils.Of(firstPage.NextPageToken)
		secondPage, err := _service.GetModelArtifacts(listOptions, nil)
		require.NoError(t, err)
		require.NotNil(t, secondPage)
		assert.GreaterOrEqual(t, len(secondPage.Items), 5)

		// Verify no duplicate IDs between pages
		firstPageIds := make(map[string]bool)
		for _, item := range firstPage.Items {
			firstPageIds[*item.Id] = true
		}

		for _, item := range secondPage.Items {
			if firstPageIds[*item.Id] {
				t.Errorf("Found duplicate ID %s between pages", *item.Id)
			}
		}
	})
}

func TestGetModelArtifactById(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful get", func(t *testing.T) {
		// Create a model artifact
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("get-model-artifact"),
			Uri:  apiutils.Of("s3://bucket/get-model.pkl"),
		}

		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		result, err := _service.GetModelArtifactById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-model-artifact", *result.Name)
		assert.Equal(t, "s3://bucket/get-model.pkl", *result.Uri)
	})

	t.Run("artifact is not model artifact", func(t *testing.T) {
		// Create a doc artifact
		docArtifact := &openapi.DocArtifact{
			Name: apiutils.Of("doc-not-model"),
			Uri:  apiutils.Of("s3://bucket/doc.pdf"),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		created, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)
		require.NotNil(t, created.DocArtifact.Id)

		// Try to get as model artifact
		result, err := _service.GetModelArtifactById(*created.DocArtifact.Id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "is not a model artifact")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := _service.GetModelArtifactById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestGetModelArtifactByInferenceService(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful get", func(t *testing.T) {
		// Create the full chain: RegisteredModel -> ModelVersion -> InferenceService -> ModelArtifact
		registeredModel := &openapi.RegisteredModel{
			Name: "inference-artifact-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "inference-artifact-env",
		}
		createdEnv, err := _service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 apiutils.Of("inference-artifact-service"),
			RegisteredModelId:    *createdModel.Id,
			ServingEnvironmentId: *createdEnv.Id,
			ModelVersionId:       createdVersion.Id,
		}
		createdInference, err := _service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Create model artifact for the model version
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("inference-model-artifact"),
			Uri:  apiutils.Of("s3://bucket/inference-model.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		_, err = _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
		require.NoError(t, err)

		// Get model artifact by inference service
		result, err := _service.GetModelArtifactByInferenceService(*createdInference.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Contains(t, *result.Name, "inference-model-artifact")
		assert.Equal(t, "s3://bucket/inference-model.pkl", *result.Uri)
	})

	t.Run("no artifacts found", func(t *testing.T) {
		// Create inference service without artifacts
		registeredModel := &openapi.RegisteredModel{
			Name: "no-artifact-model",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		servingEnv := &openapi.ServingEnvironment{
			Name: "no-artifact-env",
		}
		createdEnv, err := _service.UpsertServingEnvironment(servingEnv)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		inferenceService := &openapi.InferenceService{
			Name:                 apiutils.Of("no-artifact-service"),
			RegisteredModelId:    *createdModel.Id,
			ServingEnvironmentId: *createdEnv.Id,
			ModelVersionId:       createdVersion.Id,
		}
		createdInference, err := _service.UpsertInferenceService(inferenceService)
		require.NoError(t, err)

		// Try to get model artifact
		result, err := _service.GetModelArtifactByInferenceService(*createdInference.Id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no artifacts found")
	})
}

func TestGetModelArtifactByParams(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful get by external id", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			Name:       apiutils.Of("params-model-artifact"),
			ExternalId: apiutils.Of("model-params-ext-123"),
			Uri:        apiutils.Of("s3://bucket/params-model.pkl"),
		}

		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)

		// Get by external ID
		result, err := _service.GetModelArtifactByParams(nil, nil, apiutils.Of("model-params-ext-123"))

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "model-params-ext-123", *result.ExternalId)
	})

	t.Run("artifact is not model artifact", func(t *testing.T) {
		// Create a doc artifact
		docArtifact := &openapi.DocArtifact{
			Name:       apiutils.Of("doc-params-artifact"),
			ExternalId: apiutils.Of("doc-params-ext-123"),
			Uri:        apiutils.Of("s3://bucket/doc-params.pdf"),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		_, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)

		// Try to get as model artifact
		result, err := _service.GetModelArtifactByParams(nil, nil, apiutils.Of("doc-params-ext-123"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "is not a model artifact")
	})
}

func TestGetModelArtifacts(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("successful list all model artifacts", func(t *testing.T) {
		// Create multiple model artifacts
		for i := 0; i < 3; i++ {
			modelArtifact := &openapi.ModelArtifact{
				Name: apiutils.Of("list-model-artifact-" + string(rune('1'+i))),
				Uri:  apiutils.Of("s3://bucket/model" + string(rune('1'+i)) + ".pkl"),
			}
			_, err := _service.UpsertModelArtifact(modelArtifact)
			require.NoError(t, err)
		}

		// List all model artifacts
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(10)),
		}

		result, err := _service.GetModelArtifacts(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3)
		assert.NotNil(t, result.Size)
		assert.Equal(t, int32(10), result.PageSize)
	})

	t.Run("successful list model artifacts by model version", func(t *testing.T) {
		// Create registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name: "test-model-for-model-artifacts",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name: "v1.0",
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create model artifacts for this model version
		for i := 0; i < 2; i++ {
			modelArtifact := &openapi.ModelArtifact{
				Name: apiutils.Of("version-model-artifact-" + string(rune('1'+i))),
				Uri:  apiutils.Of("s3://bucket/version-model" + string(rune('1'+i)) + ".pkl"),
			}

			artifact := &openapi.Artifact{
				ModelArtifact: modelArtifact,
			}

			_, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
			require.NoError(t, err)
		}

		// List model artifacts for this model version
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(10)),
		}

		result, err := _service.GetModelArtifacts(listOptions, createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))
	})

	t.Run("invalid model version id", func(t *testing.T) {
		listOptions := api.ListOptions{}

		result, err := _service.GetModelArtifacts(listOptions, apiutils.Of("invalid"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})
}

func TestArtifactRoundTrip(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create registered model and model version
		registeredModel := &openapi.RegisteredModel{
			Name:        "roundtrip-model",
			Description: apiutils.Of("Model for roundtrip test"),
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		modelVersion := &openapi.ModelVersion{
			Name:        "v1.0",
			Description: apiutils.Of("Version 1.0"),
		}
		createdVersion, err := _service.UpsertModelVersion(modelVersion, createdModel.Id)
		require.NoError(t, err)

		// Create model artifact
		modelArtifact := &openapi.ModelArtifact{
			Name:               apiutils.Of("roundtrip-artifact"),
			Description:        apiutils.Of("Roundtrip test artifact"),
			Uri:                apiutils.Of("s3://bucket/roundtrip.pkl"),
			ModelFormatName:    apiutils.Of("sklearn"),
			ModelFormatVersion: apiutils.Of("1.0"),
			StorageKey:         apiutils.Of("roundtrip-key"),
			StoragePath:        apiutils.Of("/models/roundtrip"),
			ServiceAccountName: apiutils.Of("roundtrip-sa"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		// Create
		created, err := _service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
		require.NoError(t, err)
		require.NotNil(t, created.ModelArtifact.Id)

		// Get by ID
		retrieved, err := _service.GetArtifactById(*created.ModelArtifact.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.ModelArtifact)
		assert.Equal(t, *created.ModelArtifact.Id, *retrieved.ModelArtifact.Id)
		assert.Contains(t, *retrieved.ModelArtifact.Name, "roundtrip-artifact")
		assert.Equal(t, "Roundtrip test artifact", *retrieved.ModelArtifact.Description)
		assert.Equal(t, "s3://bucket/roundtrip.pkl", *retrieved.ModelArtifact.Uri)
		assert.Equal(t, "sklearn", *retrieved.ModelArtifact.ModelFormatName)
		assert.Equal(t, "1.0", *retrieved.ModelArtifact.ModelFormatVersion)
		assert.Equal(t, "roundtrip-key", *retrieved.ModelArtifact.StorageKey)
		assert.Equal(t, "/models/roundtrip", *retrieved.ModelArtifact.StoragePath)
		assert.Equal(t, "roundtrip-sa", *retrieved.ModelArtifact.ServiceAccountName)

		// Update
		retrieved.ModelArtifact.Description = apiutils.Of("Updated description")
		retrieved.ModelArtifact.Uri = apiutils.Of("s3://bucket/updated-roundtrip.pkl")
		retrieved.ModelArtifact.State = apiutils.Of(openapi.ARTIFACTSTATE_DELETED)

		updated, err := _service.UpsertArtifact(retrieved)
		require.NoError(t, err)
		require.NotNil(t, updated.ModelArtifact)
		assert.Equal(t, *created.ModelArtifact.Id, *updated.ModelArtifact.Id)
		assert.Equal(t, "Updated description", *updated.ModelArtifact.Description)
		assert.Equal(t, "s3://bucket/updated-roundtrip.pkl", *updated.ModelArtifact.Uri)
		assert.Equal(t, openapi.ARTIFACTSTATE_DELETED, *updated.ModelArtifact.State)

		// List artifacts for model version
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(10)),
		}

		artifacts, err := _service.GetArtifacts(listOptions, createdVersion.Id)
		require.NoError(t, err)
		require.NotNil(t, artifacts)
		assert.Equal(t, 1, len(artifacts.Items))
		assert.Equal(t, *updated.ModelArtifact.Id, *artifacts.Items[0].ModelArtifact.Id)
	})

	t.Run("roundtrip with custom properties", func(t *testing.T) {
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

		modelArtifact := &openapi.ModelArtifact{
			Name:             apiutils.Of("custom-props-roundtrip"),
			Uri:              apiutils.Of("s3://bucket/custom-props.pkl"),
			CustomProperties: &customProps,
		}

		// Create
		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Verify custom properties
		retrieved, err := _service.GetModelArtifactById(*created.Id)
		require.NoError(t, err)
		require.NotNil(t, retrieved.CustomProperties)

		resultProps := *retrieved.CustomProperties
		assert.Contains(t, resultProps, "accuracy")
		assert.Contains(t, resultProps, "framework")
		assert.Contains(t, resultProps, "epochs")
		assert.Contains(t, resultProps, "is_production")

		assert.Equal(t, 0.95, resultProps["accuracy"].MetadataDoubleValue.DoubleValue)
		assert.Equal(t, "tensorflow", resultProps["framework"].MetadataStringValue.StringValue)
		assert.Equal(t, "100", resultProps["epochs"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["is_production"].MetadataBoolValue.BoolValue)

		// Update custom properties
		newProps := map[string]openapi.MetadataValue{
			"accuracy": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue: 0.97,
				},
			},
			"new_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "new_value",
				},
			},
		}

		retrieved.CustomProperties = &newProps

		updated, err := _service.UpsertModelArtifact(retrieved)
		require.NoError(t, err)
		require.NotNil(t, updated.CustomProperties)

		updatedProps := *updated.CustomProperties
		assert.Contains(t, updatedProps, "accuracy")
		assert.Contains(t, updatedProps, "new_prop")
		assert.Equal(t, 0.97, updatedProps["accuracy"].MetadataDoubleValue.DoubleValue)
		assert.Equal(t, "new_value", updatedProps["new_prop"].MetadataStringValue.StringValue)
	})
}

func TestModelArtifactNilFieldsPreservation(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("nil fields preserved during model artifact upsert", func(t *testing.T) {
		// Create model artifact with only required fields, leaving optional fields as nil
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of("nil-fields-test"),
			Uri:  apiutils.Of("s3://bucket/test.pkl"),
			// Explicitly leaving these fields as nil:
			// Description: nil,
			// ExternalId: nil,
			// ModelFormatName: nil,
			// ModelFormatVersion: nil,
			// StorageKey: nil,
			// StoragePath: nil,
			// ServiceAccountName: nil,
			// ModelSourceKind: nil,
			// ModelSourceClass: nil,
			// ModelSourceGroup: nil,
			// ModelSourceId: nil,
			// ModelSourceName: nil,
			// State: nil (will get default),
		}

		// Create the artifact
		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Verify nil fields are preserved (not set to default values)
		assert.Nil(t, created.Description)
		assert.Nil(t, created.ExternalId)
		assert.Nil(t, created.ModelFormatName)
		assert.Nil(t, created.ModelFormatVersion)
		assert.Nil(t, created.StorageKey)
		assert.Nil(t, created.StoragePath)
		assert.Nil(t, created.ServiceAccountName)
		assert.Nil(t, created.ModelSourceKind)
		assert.Nil(t, created.ModelSourceClass)
		assert.Nil(t, created.ModelSourceGroup)
		assert.Nil(t, created.ModelSourceId)
		assert.Nil(t, created.ModelSourceName)

		// Update the artifact while keeping nil fields as nil
		created.Uri = apiutils.Of("s3://bucket/updated.pkl")
		// Keep all other optional fields as nil

		updated, err := _service.UpsertModelArtifact(created)
		require.NoError(t, err)

		// Verify nil fields are still preserved after update
		assert.Equal(t, "s3://bucket/updated.pkl", *updated.Uri)
		assert.Nil(t, updated.Description)
		assert.Nil(t, updated.ExternalId)
		assert.Nil(t, updated.ModelFormatName)
		assert.Nil(t, updated.ModelFormatVersion)
		assert.Nil(t, updated.StorageKey)
		assert.Nil(t, updated.StoragePath)
		assert.Nil(t, updated.ServiceAccountName)
		assert.Nil(t, updated.ModelSourceKind)
		assert.Nil(t, updated.ModelSourceClass)
		assert.Nil(t, updated.ModelSourceGroup)
		assert.Nil(t, updated.ModelSourceId)
		assert.Nil(t, updated.ModelSourceName)
	})
}

func TestDocArtifactNilFieldsPreservation(t *testing.T) {
	cleanupTestData(t, sharedDB)

	t.Run("nil fields preserved during doc artifact upsert", func(t *testing.T) {
		// Create doc artifact with only required fields, leaving optional fields as nil
		docArtifact := &openapi.DocArtifact{
			Name: apiutils.Of("nil-fields-doc-test"),
			Uri:  apiutils.Of("s3://bucket/doc.pdf"),
			// Explicitly leaving these fields as nil:
			// Description: nil,
			// ExternalId: nil,
			// State: nil (will get default),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		// Create the artifact
		created, err := _service.UpsertArtifact(artifact)
		require.NoError(t, err)
		require.NotNil(t, created.DocArtifact.Id)

		// Verify nil fields are preserved (not set to default values)
		assert.Nil(t, created.DocArtifact.Description)
		assert.Nil(t, created.DocArtifact.ExternalId)

		// Update the artifact while keeping nil fields as nil
		created.DocArtifact.Uri = apiutils.Of("s3://bucket/updated-doc.pdf")
		// Keep all other optional fields as nil

		updated, err := _service.UpsertArtifact(created)
		require.NoError(t, err)

		// Verify nil fields are still preserved after update
		assert.Equal(t, "s3://bucket/updated-doc.pdf", *updated.DocArtifact.Uri)
		assert.Nil(t, updated.DocArtifact.Description)
		assert.Nil(t, updated.DocArtifact.ExternalId)
	})
}
