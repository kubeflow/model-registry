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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
		assert.Contains(t, err.Error(), "invalid artifact type, must be either ModelArtifact, DocArtifact, DataSet, Metric, or Parameter")
	})

	t.Run("metric without value error", func(t *testing.T) {
		artifact := &openapi.Artifact{
			Metric: &openapi.Metric{
				Name: apiutils.Of("test-metric-no-value"),
				// Value is intentionally omitted
			},
		}

		result, err := _service.UpsertArtifact(artifact)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "metric value is required")
	})

	// Tests for null name handling - should generate UUID for all artifact types
	t.Run("create model artifact with null name generates UUID", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			// Name is intentionally nil/not set
			Uri: apiutils.Of("s3://bucket/model-no-name.pkl"),
		}

		artifact := &openapi.Artifact{
			ModelArtifact: modelArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.ModelArtifact)
		assert.NotNil(t, result.ModelArtifact.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.ModelArtifact.Name, "Generated name should not be empty")
		// Check if it looks like a UUID (basic check for format)
		assert.Len(t, *result.ModelArtifact.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.ModelArtifact.Name, "-", "Generated name should have UUID format")
	})

	t.Run("create doc artifact with null name generates UUID", func(t *testing.T) {
		docArtifact := &openapi.DocArtifact{
			// Name is intentionally nil/not set
			Uri: apiutils.Of("s3://bucket/doc-no-name.pdf"),
		}

		artifact := &openapi.Artifact{
			DocArtifact: docArtifact,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.DocArtifact)
		assert.NotNil(t, result.DocArtifact.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.DocArtifact.Name, "Generated name should not be empty")
		assert.Len(t, *result.DocArtifact.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.DocArtifact.Name, "-", "Generated name should have UUID format")
	})

	t.Run("create dataset with null name generates UUID", func(t *testing.T) {
		dataSet := &openapi.DataSet{
			// Name is intentionally nil/not set
			Uri: apiutils.Of("s3://bucket/dataset-no-name.csv"),
		}

		artifact := &openapi.Artifact{
			DataSet: dataSet,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.DataSet)
		assert.NotNil(t, result.DataSet.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.DataSet.Name, "Generated name should not be empty")
		assert.Len(t, *result.DataSet.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.DataSet.Name, "-", "Generated name should have UUID format")
	})

	t.Run("create metric with null name generates UUID", func(t *testing.T) {
		metric := &openapi.Metric{
			// Name is intentionally nil/not set
			Value: apiutils.Of(0.99),
		}

		artifact := &openapi.Artifact{
			Metric: metric,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Metric)
		assert.NotNil(t, result.Metric.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.Metric.Name, "Generated name should not be empty")
		assert.Len(t, *result.Metric.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.Metric.Name, "-", "Generated name should have UUID format")
	})

	t.Run("create parameter with null name generates UUID", func(t *testing.T) {
		parameter := &openapi.Parameter{
			// Name is intentionally nil/not set
			Value: apiutils.Of("param-value"),
		}

		artifact := &openapi.Artifact{
			Parameter: parameter,
		}

		result, err := _service.UpsertArtifact(artifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		require.NotNil(t, result.Parameter)
		assert.NotNil(t, result.Parameter.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.Parameter.Name, "Generated name should not be empty")
		assert.Len(t, *result.Parameter.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.Parameter.Name, "-", "Generated name should have UUID format")
	})

	t.Run("update artifact with null name preserves existing name", func(t *testing.T) {
		// First create an artifact with a specific name
		originalName := "original-artifact-name"
		modelArtifact := &openapi.ModelArtifact{
			Name: apiutils.Of(originalName),
			Uri:  apiutils.Of("s3://bucket/original.pkl"),
		}

		created, err := _service.UpsertModelArtifact(modelArtifact)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update with nil name - should preserve existing name
		updateArtifact := &openapi.ModelArtifact{
			Id: created.Id,
			// Name is intentionally nil
			Uri: apiutils.Of("s3://bucket/updated.pkl"),
		}

		updated, err := _service.UpsertModelArtifact(updateArtifact)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, originalName, *updated.Name, "Name should be preserved during update")
		assert.Equal(t, "s3://bucket/updated.pkl", *updated.Uri, "Uri should be updated")
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
		firstPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, nil)
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
			secondPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, nil)
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

		allItems, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, nil)
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

		descPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, nil)
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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
		firstPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)
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
			secondPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)
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

		allItems, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)
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

		descPage, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)
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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

	t.Run("same artifact name across different model versions - all artifact types", func(t *testing.T) {
		// This test verifies that parentResourceId filtering works correctly for all artifact types

		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "model-for-all-artifact-types",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create two model versions
		version1 := &openapi.ModelVersion{
			Name: "version-with-all-artifacts-1",
		}
		createdVersion1, err := _service.UpsertModelVersion(version1, createdModel.Id)
		require.NoError(t, err)

		version2 := &openapi.ModelVersion{
			Name: "version-with-all-artifacts-2",
		}
		createdVersion2, err := _service.UpsertModelVersion(version2, createdModel.Id)
		require.NoError(t, err)

		// Test cases for each artifact type
		artifactTypes := []struct {
			name            string
			artifactName    string
			createArtifact1 *openapi.Artifact
			createArtifact2 *openapi.Artifact
			checkField      func(*openapi.Artifact) interface{}
			getDescription  func(*openapi.Artifact) string
		}{
			{
				name:         "ModelArtifact",
				artifactName: "shared-model-artifact-name",
				createArtifact1: &openapi.Artifact{
					ModelArtifact: &openapi.ModelArtifact{
						Name:        apiutils.Of("shared-model-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/model-v1.pkl"),
						Description: apiutils.Of("Model artifact for version 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					ModelArtifact: &openapi.ModelArtifact{
						Name:        apiutils.Of("shared-model-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/model-v2.pkl"),
						Description: apiutils.Of("Model artifact for version 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.ModelArtifact },
				getDescription: func(a *openapi.Artifact) string {
					if a.ModelArtifact != nil && a.ModelArtifact.Description != nil {
						return *a.ModelArtifact.Description
					}
					return ""
				},
			},
			{
				name:         "DocArtifact",
				artifactName: "shared-doc-artifact-name",
				createArtifact1: &openapi.Artifact{
					DocArtifact: &openapi.DocArtifact{
						Name:        apiutils.Of("shared-doc-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/doc-v1.pdf"),
						Description: apiutils.Of("Doc artifact for version 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					DocArtifact: &openapi.DocArtifact{
						Name:        apiutils.Of("shared-doc-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/doc-v2.pdf"),
						Description: apiutils.Of("Doc artifact for version 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.DocArtifact },
				getDescription: func(a *openapi.Artifact) string {
					if a.DocArtifact != nil && a.DocArtifact.Description != nil {
						return *a.DocArtifact.Description
					}
					return ""
				},
			},
			{
				name:         "DataSet",
				artifactName: "shared-dataset-artifact-name",
				createArtifact1: &openapi.Artifact{
					DataSet: &openapi.DataSet{
						Name:        apiutils.Of("shared-dataset-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/dataset-v1.csv"),
						Description: apiutils.Of("Dataset for version 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					DataSet: &openapi.DataSet{
						Name:        apiutils.Of("shared-dataset-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/dataset-v2.csv"),
						Description: apiutils.Of("Dataset for version 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.DataSet },
				getDescription: func(a *openapi.Artifact) string {
					if a.DataSet != nil && a.DataSet.Description != nil {
						return *a.DataSet.Description
					}
					return ""
				},
			},
			{
				name:         "Metric",
				artifactName: "shared-metric-artifact-name",
				createArtifact1: &openapi.Artifact{
					Metric: &openapi.Metric{
						Name:        apiutils.Of("shared-metric-artifact-name"),
						Value:       apiutils.Of(0.95),
						Description: apiutils.Of("Metric for version 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					Metric: &openapi.Metric{
						Name:        apiutils.Of("shared-metric-artifact-name"),
						Value:       apiutils.Of(0.97),
						Description: apiutils.Of("Metric for version 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.Metric },
				getDescription: func(a *openapi.Artifact) string {
					if a.Metric != nil && a.Metric.Description != nil {
						return *a.Metric.Description
					}
					return ""
				},
			},
			{
				name:         "Parameter",
				artifactName: "shared-parameter-artifact-name",
				createArtifact1: &openapi.Artifact{
					Parameter: &openapi.Parameter{
						Name:        apiutils.Of("shared-parameter-artifact-name"),
						Value:       apiutils.Of("0.001"),
						Description: apiutils.Of("Parameter for version 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					Parameter: &openapi.Parameter{
						Name:        apiutils.Of("shared-parameter-artifact-name"),
						Value:       apiutils.Of("0.002"),
						Description: apiutils.Of("Parameter for version 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.Parameter },
				getDescription: func(a *openapi.Artifact) string {
					if a.Parameter != nil && a.Parameter.Description != nil {
						return *a.Parameter.Description
					}
					return ""
				},
			},
		}

		for _, tc := range artifactTypes {
			t.Run(tc.name, func(t *testing.T) {
				// Create artifact with same name for version 1
				created1, err := _service.UpsertModelVersionArtifact(tc.createArtifact1, *createdVersion1.Id)
				require.NoError(t, err)
				require.NotNil(t, tc.checkField(created1))

				// Create artifact with same name for version 2
				created2, err := _service.UpsertModelVersionArtifact(tc.createArtifact2, *createdVersion2.Id)
				require.NoError(t, err)
				require.NotNil(t, tc.checkField(created2))

				// Query for artifact by name and version 1
				result1, err := _service.GetArtifactByParams(&tc.artifactName, createdVersion1.Id, nil)
				require.NoError(t, err)
				require.NotNil(t, result1)
				require.NotNil(t, tc.checkField(result1))
				assert.Contains(t, tc.getDescription(result1), "version 1")

				// Query for artifact by name and version 2
				result2, err := _service.GetArtifactByParams(&tc.artifactName, createdVersion2.Id, nil)
				require.NoError(t, err)
				require.NotNil(t, result2)
				require.NotNil(t, tc.checkField(result2))
				assert.Contains(t, tc.getDescription(result2), "version 2")

				// Ensure we got different artifacts
				assert.NotEqual(t, tc.getDescription(result1), tc.getDescription(result2))
			})
		}
	})

	t.Run("same artifact name across different experiment runs - all artifact types", func(t *testing.T) {
		// This test verifies that parentResourceId filtering works correctly for all artifact types in experiment runs

		// Create an experiment
		experiment := &openapi.Experiment{
			Name: "experiment-for-all-artifact-types",
		}
		createdExperiment, err := _service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Create two experiment runs
		run1 := &openapi.ExperimentRun{
			Name: apiutils.Of("run-with-all-artifacts-1"),
		}
		createdRun1, err := _service.UpsertExperimentRun(run1, createdExperiment.Id)
		require.NoError(t, err)

		run2 := &openapi.ExperimentRun{
			Name: apiutils.Of("run-with-all-artifacts-2"),
		}
		createdRun2, err := _service.UpsertExperimentRun(run2, createdExperiment.Id)
		require.NoError(t, err)

		// Test cases for each artifact type
		artifactTypes := []struct {
			name            string
			artifactName    string
			createArtifact1 *openapi.Artifact
			createArtifact2 *openapi.Artifact
			checkField      func(*openapi.Artifact) interface{}
			getDescription  func(*openapi.Artifact) string
		}{
			{
				name:         "ModelArtifact",
				artifactName: "shared-run-model-artifact-name",
				createArtifact1: &openapi.Artifact{
					ModelArtifact: &openapi.ModelArtifact{
						Name:        apiutils.Of("shared-run-model-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run1-model.pkl"),
						Description: apiutils.Of("Model artifact for run 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					ModelArtifact: &openapi.ModelArtifact{
						Name:        apiutils.Of("shared-run-model-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run2-model.pkl"),
						Description: apiutils.Of("Model artifact for run 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.ModelArtifact },
				getDescription: func(a *openapi.Artifact) string {
					if a.ModelArtifact != nil && a.ModelArtifact.Description != nil {
						return *a.ModelArtifact.Description
					}
					return ""
				},
			},
			{
				name:         "DocArtifact",
				artifactName: "shared-run-doc-artifact-name",
				createArtifact1: &openapi.Artifact{
					DocArtifact: &openapi.DocArtifact{
						Name:        apiutils.Of("shared-run-doc-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run1-doc.pdf"),
						Description: apiutils.Of("Doc artifact for run 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					DocArtifact: &openapi.DocArtifact{
						Name:        apiutils.Of("shared-run-doc-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run2-doc.pdf"),
						Description: apiutils.Of("Doc artifact for run 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.DocArtifact },
				getDescription: func(a *openapi.Artifact) string {
					if a.DocArtifact != nil && a.DocArtifact.Description != nil {
						return *a.DocArtifact.Description
					}
					return ""
				},
			},
			{
				name:         "DataSet",
				artifactName: "shared-run-dataset-artifact-name",
				createArtifact1: &openapi.Artifact{
					DataSet: &openapi.DataSet{
						Name:        apiutils.Of("shared-run-dataset-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run1-dataset.csv"),
						Description: apiutils.Of("Dataset for run 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					DataSet: &openapi.DataSet{
						Name:        apiutils.Of("shared-run-dataset-artifact-name"),
						Uri:         apiutils.Of("s3://bucket/run2-dataset.csv"),
						Description: apiutils.Of("Dataset for run 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.DataSet },
				getDescription: func(a *openapi.Artifact) string {
					if a.DataSet != nil && a.DataSet.Description != nil {
						return *a.DataSet.Description
					}
					return ""
				},
			},
			{
				name:         "Metric",
				artifactName: "shared-run-metric-artifact-name",
				createArtifact1: &openapi.Artifact{
					Metric: &openapi.Metric{
						Name:        apiutils.Of("shared-run-metric-artifact-name"),
						Value:       apiutils.Of(0.91),
						Description: apiutils.Of("Metric for run 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					Metric: &openapi.Metric{
						Name:        apiutils.Of("shared-run-metric-artifact-name"),
						Value:       apiutils.Of(0.93),
						Description: apiutils.Of("Metric for run 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.Metric },
				getDescription: func(a *openapi.Artifact) string {
					if a.Metric != nil && a.Metric.Description != nil {
						return *a.Metric.Description
					}
					return ""
				},
			},
			{
				name:         "Parameter",
				artifactName: "shared-run-parameter-artifact-name",
				createArtifact1: &openapi.Artifact{
					Parameter: &openapi.Parameter{
						Name:        apiutils.Of("shared-run-parameter-artifact-name"),
						Value:       apiutils.Of("0.01"),
						Description: apiutils.Of("Parameter for run 1"),
					},
				},
				createArtifact2: &openapi.Artifact{
					Parameter: &openapi.Parameter{
						Name:        apiutils.Of("shared-run-parameter-artifact-name"),
						Value:       apiutils.Of("0.02"),
						Description: apiutils.Of("Parameter for run 2"),
					},
				},
				checkField: func(a *openapi.Artifact) interface{} { return a.Parameter },
				getDescription: func(a *openapi.Artifact) string {
					if a.Parameter != nil && a.Parameter.Description != nil {
						return *a.Parameter.Description
					}
					return ""
				},
			},
		}

		for _, tc := range artifactTypes {
			t.Run(tc.name, func(t *testing.T) {
				// Create artifact with same name for run 1
				created1, err := _service.UpsertExperimentRunArtifact(tc.createArtifact1, *createdRun1.Id)
				require.NoError(t, err)
				require.NotNil(t, tc.checkField(created1))

				// Create artifact with same name for run 2
				created2, err := _service.UpsertExperimentRunArtifact(tc.createArtifact2, *createdRun2.Id)
				require.NoError(t, err)
				require.NotNil(t, tc.checkField(created2))

				// Query for artifact by name and run 1
				result1, err := _service.GetArtifactByParams(&tc.artifactName, createdRun1.Id, nil)
				require.NoError(t, err)
				require.NotNil(t, result1)
				require.NotNil(t, tc.checkField(result1))
				assert.Contains(t, tc.getDescription(result1), "run 1")

				// Query for artifact by name and run 2
				result2, err := _service.GetArtifactByParams(&tc.artifactName, createdRun2.Id, nil)
				require.NoError(t, err)
				require.NotNil(t, result2)
				require.NotNil(t, tc.checkField(result2))
				assert.Contains(t, tc.getDescription(result2), "run 2")

				// Ensure we got different artifacts
				assert.NotEqual(t, tc.getDescription(result1), tc.getDescription(result2))
			})
		}
	})
}

func TestGetArtifacts(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

		result, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2)
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

		result, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 3, len(result.Items))
	})

	t.Run("invalid model version id", func(t *testing.T) {
		listOptions := api.ListOptions{}

		result, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, apiutils.Of("invalid"))

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})
}

func TestUpsertModelArtifact(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

	t.Run("create with null name generates UUID", func(t *testing.T) {
		modelArtifact := &openapi.ModelArtifact{
			// Name is intentionally nil
			Uri:             apiutils.Of("s3://bucket/direct-no-name.pkl"),
			ModelFormatName: apiutils.Of("tensorflow"),
		}

		result, err := _service.UpsertModelArtifact(modelArtifact)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Name, "Name should be auto-generated")
		assert.NotEmpty(t, *result.Name, "Generated name should not be empty")
		assert.Len(t, *result.Name, 36, "Generated name should be UUID length")
		assert.Contains(t, *result.Name, "-", "Generated name should have UUID format")
		assert.Equal(t, "s3://bucket/direct-no-name.pkl", *result.Uri)
		assert.Equal(t, "tensorflow", *result.ModelFormatName)
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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

	t.Run("same model artifact name across different model versions", func(t *testing.T) {
		// This test catches the bug where ParentResourceID was not being used to filter artifacts

		// Create a registered model
		registeredModel := &openapi.RegisteredModel{
			Name: "model-with-shared-artifacts",
		}
		createdModel, err := _service.UpsertRegisteredModel(registeredModel)
		require.NoError(t, err)

		// Create first model version
		version1 := &openapi.ModelVersion{
			Name:              "version-with-shared-artifact-1",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion1, err := _service.UpsertModelVersion(version1, createdModel.Id)
		require.NoError(t, err)

		// Create second model version
		version2 := &openapi.ModelVersion{
			Name:              "version-with-shared-artifact-2",
			RegisteredModelId: *createdModel.Id,
		}
		createdVersion2, err := _service.UpsertModelVersion(version2, createdModel.Id)
		require.NoError(t, err)

		// Create model artifact "shared-artifact-name-test" for the first version
		artifact1 := &openapi.ModelArtifact{
			Name:            apiutils.Of("shared-artifact-name-test"),
			Uri:             apiutils.Of("s3://bucket/artifact-v1.pkl"),
			Description:     apiutils.Of("Artifact for version 1"),
			ModelFormatName: apiutils.Of("pickle"),
		}
		artifactWrapper1 := &openapi.Artifact{
			ModelArtifact: artifact1,
		}
		createdArtifact1, err := _service.UpsertModelVersionArtifact(artifactWrapper1, *createdVersion1.Id)
		require.NoError(t, err)

		// Create model artifact "shared-artifact-name-test" for the second version
		artifact2 := &openapi.ModelArtifact{
			Name:            apiutils.Of("shared-artifact-name-test"),
			Uri:             apiutils.Of("s3://bucket/artifact-v2.pkl"),
			Description:     apiutils.Of("Artifact for version 2"),
			ModelFormatName: apiutils.Of("pickle"),
		}
		artifactWrapper2 := &openapi.Artifact{
			ModelArtifact: artifact2,
		}
		createdArtifact2, err := _service.UpsertModelVersionArtifact(artifactWrapper2, *createdVersion2.Id)
		require.NoError(t, err)

		// Query for artifact "shared-artifact-name-test" of the first version
		artifactName := "shared-artifact-name-test"
		result1, err := _service.GetModelArtifactByParams(&artifactName, createdVersion1.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result1)
		assert.Equal(t, *createdArtifact1.ModelArtifact.Id, *result1.Id)
		assert.Equal(t, "Artifact for version 1", *result1.Description)
		assert.Equal(t, "s3://bucket/artifact-v1.pkl", *result1.Uri)

		// Query for artifact "shared-artifact-name-test" of the second version
		result2, err := _service.GetModelArtifactByParams(&artifactName, createdVersion2.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result2)
		assert.Equal(t, *createdArtifact2.ModelArtifact.Id, *result2.Id)
		assert.Equal(t, "Artifact for version 2", *result2.Description)
		assert.Equal(t, "s3://bucket/artifact-v2.pkl", *result2.Uri)

		// Ensure we got different artifacts
		assert.NotEqual(t, *result1.Id, *result2.Id)
	})
}

func TestGetModelArtifacts(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

		artifacts, err := _service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdVersion.Id)
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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

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

func TestArtifactTypeFiltering(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Setup: Create a registered model, model version, and experiment + experiment run
	registeredModel := &openapi.RegisteredModel{
		Name: "artifact-type-test-model",
	}
	createdModel, err := service.UpsertRegisteredModel(registeredModel)
	require.NoError(t, err)

	modelVersion := &openapi.ModelVersion{
		Name: "v1.0",
	}
	createdVersion, err := service.UpsertModelVersion(modelVersion, createdModel.Id)
	require.NoError(t, err)

	experiment := &openapi.Experiment{
		Name: "artifact-type-test-experiment",
	}
	createdExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	experimentRun := &openapi.ExperimentRun{
		Name: apiutils.Of("artifact-type-test-run"),
	}
	createdExperimentRun, err := service.UpsertExperimentRun(experimentRun, createdExperiment.Id)
	require.NoError(t, err)

	// Create one artifact of each type for general testing
	t.Run("setup artifacts", func(t *testing.T) {
		// Create ModelArtifact
		modelArtifact := &openapi.Artifact{
			ModelArtifact: &openapi.ModelArtifact{
				Name: apiutils.Of("test-model-artifact"),
				Uri:  apiutils.Of("s3://bucket/model.pkl"),
			},
		}
		_, err := service.UpsertArtifact(modelArtifact)
		require.NoError(t, err)

		// Create DocArtifact
		docArtifact := &openapi.Artifact{
			DocArtifact: &openapi.DocArtifact{
				Name: apiutils.Of("test-doc-artifact"),
				Uri:  apiutils.Of("s3://bucket/doc.pdf"),
			},
		}
		_, err = service.UpsertArtifact(docArtifact)
		require.NoError(t, err)

		// Create DataSet
		dataSet := &openapi.Artifact{
			DataSet: &openapi.DataSet{
				Name: apiutils.Of("test-dataset-artifact"),
				Uri:  apiutils.Of("s3://bucket/dataset.csv"),
			},
		}
		_, err = service.UpsertArtifact(dataSet)
		require.NoError(t, err)

		// Create Metric
		metric := &openapi.Artifact{
			Metric: &openapi.Metric{
				Name:  apiutils.Of("test-metric-artifact"),
				Value: apiutils.Of(0.95),
			},
		}
		_, err = service.UpsertArtifact(metric)
		require.NoError(t, err)

		// Create Parameter
		parameter := &openapi.Artifact{
			Parameter: &openapi.Parameter{
				Name:  apiutils.Of("test-parameter-artifact"),
				Value: apiutils.Of("param-value"),
			},
		}
		_, err = service.UpsertArtifact(parameter)
		require.NoError(t, err)
	})

	// Test all artifact types for GetArtifacts (general endpoint)
	t.Run("GetArtifacts endpoint filtering", func(t *testing.T) {
		testCases := []struct {
			name         string
			artifactType openapi.ArtifactTypeQueryParam
			expectField  string
		}{
			{"model-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, "ModelArtifact"},
			{"doc-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT, "DocArtifact"},
			{"dataset-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT, "DataSet"},
			{"metric filter", openapi.ARTIFACTTYPEQUERYPARAM_METRIC, "Metric"},
			{"parameter filter", openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER, "Parameter"},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				listOptions := api.ListOptions{
					PageSize: apiutils.Of(int32(100)),
				}

				result, err := service.GetArtifacts(tc.artifactType, listOptions, nil)
				require.NoError(t, err)
				require.NotNil(t, result)

				// Should have at least one artifact of the specified type
				assert.GreaterOrEqual(t, len(result.Items), 1, "Should find at least one artifact of type %s", tc.artifactType)

				// Verify all returned artifacts are of the correct type
				for i, artifact := range result.Items {
					switch tc.expectField {
					case "ModelArtifact":
						assert.NotNil(t, artifact.ModelArtifact, "Artifact %d should be ModelArtifact", i)
						assert.Equal(t, string(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT), *artifact.ModelArtifact.ArtifactType)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DocArtifact":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.NotNil(t, artifact.DocArtifact, "Artifact %d should be DocArtifact", i)
						assert.Equal(t, string(openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT), *artifact.DocArtifact.ArtifactType)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DataSet":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.NotNil(t, artifact.DataSet, "Artifact %d should be DataSet", i)
						assert.Equal(t, string(openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT), *artifact.DataSet.ArtifactType)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Metric":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.NotNil(t, artifact.Metric, "Artifact %d should be Metric", i)
						assert.Equal(t, string(openapi.ARTIFACTTYPEQUERYPARAM_METRIC), *artifact.Metric.ArtifactType)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Parameter":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.NotNil(t, artifact.Parameter, "Artifact %d should be Parameter", i)
						assert.Equal(t, string(openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER), *artifact.Parameter.ArtifactType)
					}
				}
			})
		}

		// Test empty filter returns all types
		t.Run("no filter returns all types", func(t *testing.T) {
			listOptions := api.ListOptions{
				PageSize: apiutils.Of(int32(100)),
			}

			result, err := service.GetArtifacts("", listOptions, nil)
			require.NoError(t, err)
			require.NotNil(t, result)

			// Should have at least 5 artifacts (ModelArtifact, DocArtifact, DataSet, Metric, Parameter)
			assert.GreaterOrEqual(t, len(result.Items), 5, "Should find artifacts of all types when no filter is applied")
		})
	})

	// Create artifacts specifically associated with model version
	t.Run("setup model version artifacts", func(t *testing.T) {
		// Create different types of artifacts for the model version
		artifacts := []*openapi.Artifact{
			{
				ModelArtifact: &openapi.ModelArtifact{
					Name: apiutils.Of("mv-model-artifact"),
					Uri:  apiutils.Of("s3://bucket/mv-model.pkl"),
				},
			},
			{
				DocArtifact: &openapi.DocArtifact{
					Name: apiutils.Of("mv-doc-artifact"),
					Uri:  apiutils.Of("s3://bucket/mv-doc.pdf"),
				},
			},
			{
				DataSet: &openapi.DataSet{
					Name: apiutils.Of("mv-dataset-artifact"),
					Uri:  apiutils.Of("s3://bucket/mv-dataset.csv"),
				},
			},
			{
				Metric: &openapi.Metric{
					Name:  apiutils.Of("mv-metric-artifact"),
					Value: apiutils.Of(0.95),
				},
			},
			{
				Parameter: &openapi.Parameter{
					Name:  apiutils.Of("mv-parameter-artifact"),
					Value: apiutils.Of("mv-param-value"),
				},
			},
		}

		for _, artifact := range artifacts {
			_, err := service.UpsertModelVersionArtifact(artifact, *createdVersion.Id)
			require.NoError(t, err)
		}
	})

	// Test all artifact types for GetArtifacts with model version (scoped endpoint)
	t.Run("GetArtifacts with model version filtering", func(t *testing.T) {
		testCases := []struct {
			name         string
			artifactType openapi.ArtifactTypeQueryParam
			expectField  string
			expectCount  int
		}{
			{"model-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, "ModelArtifact", 1},
			{"doc-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT, "DocArtifact", 1},
			{"dataset-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT, "DataSet", 1},
			{"metric filter", openapi.ARTIFACTTYPEQUERYPARAM_METRIC, "Metric", 1},
			{"parameter filter", openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER, "Parameter", 1},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				listOptions := api.ListOptions{
					PageSize: apiutils.Of(int32(100)),
				}

				result, err := service.GetArtifacts(tc.artifactType, listOptions, createdVersion.Id)
				require.NoError(t, err)
				require.NotNil(t, result)

				assert.Equal(t, tc.expectCount, len(result.Items), "Should find exactly %d artifacts of type %s for this model version", tc.expectCount, tc.artifactType)

				// Verify all returned artifacts are of the correct type (if any)
				for i, artifact := range result.Items {
					switch tc.expectField {
					case "ModelArtifact":
						assert.NotNil(t, artifact.ModelArtifact, "Artifact %d should be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DocArtifact":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.NotNil(t, artifact.DocArtifact, "Artifact %d should be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DataSet":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.NotNil(t, artifact.DataSet, "Artifact %d should be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Metric":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.NotNil(t, artifact.Metric, "Artifact %d should be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Parameter":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.NotNil(t, artifact.Parameter, "Artifact %d should be Parameter", i)
					}
				}
			})
		}
	})

	// Create artifacts specifically associated with experiment run
	t.Run("setup experiment run artifacts", func(t *testing.T) {
		// Create different types of artifacts for the experiment run
		artifacts := []*openapi.Artifact{
			{
				ModelArtifact: &openapi.ModelArtifact{
					Name: apiutils.Of("er-model-artifact"),
					Uri:  apiutils.Of("s3://bucket/er-model.pkl"),
				},
			},
			{
				DocArtifact: &openapi.DocArtifact{
					Name: apiutils.Of("er-doc-artifact"),
					Uri:  apiutils.Of("s3://bucket/er-doc.pdf"),
				},
			},
			{
				DataSet: &openapi.DataSet{
					Name: apiutils.Of("er-dataset-artifact"),
					Uri:  apiutils.Of("s3://bucket/er-dataset.csv"),
				},
			},
			{
				Metric: &openapi.Metric{
					Name:  apiutils.Of("er-metric-artifact"),
					Value: apiutils.Of(0.85),
				},
			},
			{
				Parameter: &openapi.Parameter{
					Name:  apiutils.Of("er-parameter-artifact"),
					Value: apiutils.Of("er-param-value"),
				},
			},
		}

		for _, artifact := range artifacts {
			_, err := service.UpsertExperimentRunArtifact(artifact, *createdExperimentRun.Id)
			require.NoError(t, err)
		}

		// Create multiple metric values to generate metric history records
		metricName := "er-accuracy-history"
		values := []float64{0.1, 0.5, 0.8, 0.95}
		for i, value := range values {
			metricArtifact := &openapi.Artifact{
				Metric: &openapi.Metric{
					Name:        apiutils.Of(metricName),
					Value:       apiutils.Of(value),
					Description: apiutils.Of(fmt.Sprintf("Accuracy step %d", i+1)),
				},
			}
			_, err := service.UpsertExperimentRunArtifact(metricArtifact, *createdExperimentRun.Id)
			require.NoError(t, err)
		}
	})

	// Test all artifact types for GetExperimentRunArtifacts (scoped endpoint)
	t.Run("GetExperimentRunArtifacts filtering", func(t *testing.T) {
		testCases := []struct {
			name         string
			artifactType openapi.ArtifactTypeQueryParam
			expectField  string
			expectCount  int
		}{
			{"model-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, "ModelArtifact", 1},
			{"doc-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DOC_ARTIFACT, "DocArtifact", 1},
			{"dataset-artifact filter", openapi.ARTIFACTTYPEQUERYPARAM_DATASET_ARTIFACT, "DataSet", 1},
			{"metric filter", openapi.ARTIFACTTYPEQUERYPARAM_METRIC, "Metric", 2},
			{"parameter filter", openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER, "Parameter", 1},
		}

		for _, tc := range testCases {
			t.Run(tc.name, func(t *testing.T) {
				listOptions := api.ListOptions{
					PageSize: apiutils.Of(int32(100)),
				}

				result, err := service.GetExperimentRunArtifacts(tc.artifactType, listOptions, createdExperimentRun.Id)
				require.NoError(t, err)
				require.NotNil(t, result)

				assert.Equal(t, tc.expectCount, len(result.Items), "Should find exactly %d artifacts of type %s for this experiment run", tc.expectCount, tc.artifactType)

				// Verify all returned artifacts are of the correct type (if any)
				for i, artifact := range result.Items {
					switch tc.expectField {
					case "ModelArtifact":
						assert.NotNil(t, artifact.ModelArtifact, "Artifact %d should be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DocArtifact":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.NotNil(t, artifact.DocArtifact, "Artifact %d should be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "DataSet":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.NotNil(t, artifact.DataSet, "Artifact %d should be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Metric":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.NotNil(t, artifact.Metric, "Artifact %d should be Metric", i)
						assert.Nil(t, artifact.Parameter, "Artifact %d should not be Parameter", i)
					case "Parameter":
						assert.Nil(t, artifact.ModelArtifact, "Artifact %d should not be ModelArtifact", i)
						assert.Nil(t, artifact.DocArtifact, "Artifact %d should not be DocArtifact", i)
						assert.Nil(t, artifact.DataSet, "Artifact %d should not be DataSet", i)
						assert.Nil(t, artifact.Metric, "Artifact %d should not be Metric", i)
						assert.NotNil(t, artifact.Parameter, "Artifact %d should be Parameter", i)
					}
				}
			})
		}
	})

	// Test edge cases
	t.Run("edge cases", func(t *testing.T) {
		t.Run("invalid artifact type", func(t *testing.T) {
			listOptions := api.ListOptions{
				PageSize: apiutils.Of(int32(100)),
			}

			result, err := service.GetArtifacts("invalid-artifact-type", listOptions, nil)
			assert.Error(t, err)
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid artifact type")
		})

		t.Run("empty result with valid filter", func(t *testing.T) {
			// Create a new model version with no artifacts
			emptyModel := &openapi.RegisteredModel{
				Name: "empty-test-model",
			}
			createdEmptyModel, err := service.UpsertRegisteredModel(emptyModel)
			require.NoError(t, err)

			emptyModelVersion := &openapi.ModelVersion{
				Name: "v1.0",
			}
			createdEmptyVersion, err := service.UpsertModelVersion(emptyModelVersion, createdEmptyModel.Id)
			require.NoError(t, err)

			listOptions := api.ListOptions{
				PageSize: apiutils.Of(int32(100)),
			}

			result, err := service.GetArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_MODEL_ARTIFACT, listOptions, createdEmptyVersion.Id)
			require.NoError(t, err)
			require.NotNil(t, result)
			assert.Equal(t, 0, len(result.Items), "Should find no artifacts for empty model version")
		})
	})

	// Test that metric history records are NOT returned as artifacts
	t.Run("metric history filtering", func(t *testing.T) {
		// Verify that GetExperimentRunArtifacts does NOT return metric history records
		listOptions := api.ListOptions{
			PageSize: apiutils.Of(int32(100)),
		}

		result, err := service.GetExperimentRunArtifacts("", listOptions, createdExperimentRun.Id)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Count artifacts by type - should have exactly 6 artifacts:
		// 1 ModelArtifact, 1 DocArtifact, 1 DataSet, 2 Metrics (er-metric-artifact + er-accuracy-history), 1 Parameter
		// NOTE: Should NOT have 4 additional metric history records
		var modelCount, docCount, datasetCount, metricCount, parameterCount int
		metricNames := make([]string, 0)

		for _, artifact := range result.Items {
			switch {
			case artifact.ModelArtifact != nil:
				modelCount++
			case artifact.DocArtifact != nil:
				docCount++
			case artifact.DataSet != nil:
				datasetCount++
			case artifact.Metric != nil:
				metricCount++
				metricNames = append(metricNames, *artifact.Metric.Name)
			case artifact.Parameter != nil:
				parameterCount++
			}
		}

		assert.Equal(t, 1, modelCount, "Should have exactly 1 ModelArtifact")
		assert.Equal(t, 1, docCount, "Should have exactly 1 DocArtifact")
		assert.Equal(t, 1, datasetCount, "Should have exactly 1 DataSet")
		assert.Equal(t, 2, metricCount, "Should have exactly 2 Metrics (not 6 with history records)")
		assert.Equal(t, 1, parameterCount, "Should have exactly 1 Parameter")

		// Verify the metric names are the expected ones (current metrics, not history)
		expectedMetricNames := []string{"er-metric-artifact", "er-accuracy-history"}
		assert.ElementsMatch(t, expectedMetricNames, metricNames, "Should only have current metric artifacts, not history records")

		// Total should be 6 artifacts, not 10 (6 + 4 history records)
		assert.Equal(t, 6, len(result.Items), "Should have exactly 6 artifacts total (no metric history records)")

		// Verify metric history is still accessible via dedicated endpoint
		metricName := "er-accuracy-history"
		metricHistory, err := service.GetExperimentRunMetricHistory(&metricName, nil, api.ListOptions{}, createdExperimentRun.Id)
		require.NoError(t, err)
		require.NotNil(t, metricHistory)

		// Should have all 4 history values
		assert.Equal(t, 4, len(metricHistory.Items), "Metric history endpoint should return all 4 history records")

		// Verify values are correct
		expectedValues := []float64{0.1, 0.5, 0.8, 0.95}
		for i, historyItem := range metricHistory.Items {
			assert.Equal(t, expectedValues[i], *historyItem.Value,
				fmt.Sprintf("History item %d should have value %f", i, expectedValues[i]))
		}
	})
}

func TestEmbedMDMetricDuplicateHandling(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-duplicate-metrics",
		Description: apiutils.Of("Test experiment for duplicate metric handling"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-duplicate-metrics"),
		Description: apiutils.Of("Test experiment run for duplicate metric handling"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create first metric
	firstMetric := &openapi.Artifact{
		Metric: &openapi.Metric{
			Name:        apiutils.Of("accuracy"),
			Value:       apiutils.Of(0.85),
			Timestamp:   apiutils.Of("1234567890"),
			Step:        apiutils.Of(int64(1)),
			Description: apiutils.Of("First accuracy measurement"),
		},
	}

	// Upsert the first metric
	firstResult, err := service.UpsertExperimentRunArtifact(firstMetric, *savedExperimentRun.Id)
	require.NoError(t, err, "error creating first metric")
	require.NotNil(t, firstResult.Metric)
	firstMetricId := firstResult.Metric.Id

	// Create second metric with same name but different value
	secondMetric := &openapi.Artifact{
		Metric: &openapi.Metric{
			Name:        apiutils.Of("accuracy"), // Same name as first metric
			Value:       apiutils.Of(0.92),       // Different value
			Timestamp:   apiutils.Of("1234567900"),
			Step:        apiutils.Of(int64(2)),
			Description: apiutils.Of("Updated accuracy measurement"),
		},
	}

	// Upsert the second metric - should update the existing one
	secondResult, err := service.UpsertExperimentRunArtifact(secondMetric, *savedExperimentRun.Id)
	require.NoError(t, err, "error creating/updating second metric")
	require.NotNil(t, secondResult.Metric)

	// Verify that it's the same metric ID (updated, not created new)
	assert.Equal(t, firstMetricId, secondResult.Metric.Id, "should update existing metric, not create new one")

	// Verify the value was updated
	assert.Equal(t, 0.92, *secondResult.Metric.Value, "metric value should be updated")
	assert.Equal(t, "Updated accuracy measurement", *secondResult.Metric.Description, "metric description should be updated")

	// Verify only one metric exists for this experiment run
	artifacts, err := service.GetExperimentRunArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_METRIC, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err)
	assert.Equal(t, int32(1), artifacts.Size, "should have only one metric artifact")
	assert.Equal(t, 1, len(artifacts.Items), "should have only one metric in results")

	// Verify it's the updated metric
	retrievedMetric := artifacts.Items[0].Metric
	assert.Equal(t, "accuracy", *retrievedMetric.Name)
	assert.Equal(t, 0.92, *retrievedMetric.Value)
}

func TestEmbedMDParameterDuplicateHandling(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-duplicate-parameters",
		Description: apiutils.Of("Test experiment for duplicate parameter handling"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-duplicate-parameters"),
		Description: apiutils.Of("Test experiment run for duplicate parameter handling"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create first parameter
	firstParameter := &openapi.Artifact{
		Parameter: &openapi.Parameter{
			Name:        apiutils.Of("learning_rate"),
			Value:       apiutils.Of("0.01"),
			Description: apiutils.Of("Initial learning rate"),
		},
	}

	// Upsert the first parameter
	firstResult, err := service.UpsertExperimentRunArtifact(firstParameter, *savedExperimentRun.Id)
	require.NoError(t, err, "error creating first parameter")
	require.NotNil(t, firstResult.Parameter)
	firstParameterId := firstResult.Parameter.Id

	// Create second parameter with same name but different value
	secondParameter := &openapi.Artifact{
		Parameter: &openapi.Parameter{
			Name:        apiutils.Of("learning_rate"), // Same name as first parameter
			Value:       apiutils.Of("0.001"),         // Different value
			Description: apiutils.Of("Updated learning rate"),
		},
	}

	// Upsert the second parameter - should update the existing one
	secondResult, err := service.UpsertExperimentRunArtifact(secondParameter, *savedExperimentRun.Id)
	require.NoError(t, err, "error creating/updating second parameter")
	require.NotNil(t, secondResult.Parameter)

	// Verify that it's the same parameter ID (updated, not created new)
	assert.Equal(t, firstParameterId, secondResult.Parameter.Id, "should update existing parameter, not create new one")

	// Verify the value was updated
	assert.Equal(t, "0.001", *secondResult.Parameter.Value, "parameter value should be updated")
	assert.Equal(t, "Updated learning rate", *secondResult.Parameter.Description, "parameter description should be updated")

	// Verify only one parameter exists for this experiment run
	artifacts, err := service.GetExperimentRunArtifacts(openapi.ARTIFACTTYPEQUERYPARAM_PARAMETER, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err)
	assert.Equal(t, int32(1), artifacts.Size, "should have only one parameter artifact")
	assert.Equal(t, 1, len(artifacts.Items), "should have only one parameter in results")

	// Verify it's the updated parameter
	retrievedParameter := artifacts.Items[0].Parameter
	assert.Equal(t, "learning_rate", *retrievedParameter.Name)
	assert.Equal(t, "0.001", *retrievedParameter.Value)
}
