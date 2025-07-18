package core_test

import (
	"fmt"
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertServingEnvironment(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		input := &openapi.ServingEnvironment{
			Name:        "test-serving-env",
			Description: apiutils.Of("Test serving environment description"),
			ExternalId:  apiutils.Of("serving-ext-123"),
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-serving-env", result.Name)
		assert.Equal(t, "serving-ext-123", *result.ExternalId)
		assert.Equal(t, "Test serving environment description", *result.Description)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
	})

	t.Run("successful update", func(t *testing.T) {
		// Create first
		input := &openapi.ServingEnvironment{
			Name:        "update-test-serving-env",
			Description: apiutils.Of("Original description"),
		}

		created, err := _service.UpsertServingEnvironment(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Update
		update := &openapi.ServingEnvironment{
			Id:          created.Id,
			Name:        "update-test-serving-env", // Name should remain the same
			Description: apiutils.Of("Updated description"),
			ExternalId:  apiutils.Of("updated-ext-456"),
		}

		updated, err := _service.UpsertServingEnvironment(update)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "update-test-serving-env", updated.Name)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-ext-456", *updated.ExternalId)
	})

	t.Run("create with custom properties", func(t *testing.T) {
		customProps := map[string]openapi.MetadataValue{
			"cpu_limit": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "2000m",
				},
			},
			"memory_limit": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "4Gi",
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
			"cost_per_hour": {
				MetadataDoubleValue: &openapi.MetadataDoubleValue{
					DoubleValue: 0.50,
				},
			},
		}

		input := &openapi.ServingEnvironment{
			Name:             "custom-props-serving-env",
			CustomProperties: &customProps,
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "custom-props-serving-env", result.Name)
		assert.NotNil(t, result.CustomProperties)

		resultProps := *result.CustomProperties
		assert.Contains(t, resultProps, "cpu_limit")
		assert.Contains(t, resultProps, "memory_limit")
		assert.Contains(t, resultProps, "replicas")
		assert.Contains(t, resultProps, "auto_scaling")
		assert.Contains(t, resultProps, "cost_per_hour")

		assert.Equal(t, "2000m", resultProps["cpu_limit"].MetadataStringValue.StringValue)
		assert.Equal(t, "4Gi", resultProps["memory_limit"].MetadataStringValue.StringValue)
		assert.Equal(t, "3", resultProps["replicas"].MetadataIntValue.IntValue)
		assert.Equal(t, true, resultProps["auto_scaling"].MetadataBoolValue.BoolValue)
		assert.Equal(t, 0.50, resultProps["cost_per_hour"].MetadataDoubleValue.DoubleValue)
	})

	t.Run("minimal serving environment", func(t *testing.T) {
		input := &openapi.ServingEnvironment{
			Name: "minimal-serving-env",
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, "minimal-serving-env", result.Name)
		assert.NotNil(t, result.Id)
	})

	t.Run("nil serving environment error", func(t *testing.T) {
		result, err := _service.UpsertServingEnvironment(nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid serving environment pointer")
	})

	t.Run("unicode characters in name", func(t *testing.T) {
		unicodeName := "测试环境-тест-環境-🚀"
		input := &openapi.ServingEnvironment{
			Name:        unicodeName,
			Description: apiutils.Of("Unicode test environment with 中文, русский, 日本語, and emoji 🎯"),
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, unicodeName, result.Name)
		assert.Equal(t, "Unicode test environment with 中文, русский, 日本語, and emoji 🎯", *result.Description)
		assert.NotNil(t, result.Id)
	})

	t.Run("special characters in name", func(t *testing.T) {
		specialName := "test-env!@#$%^&*()_+-=[]{}|;':\",./<>?"
		input := &openapi.ServingEnvironment{
			Name:        specialName,
			Description: apiutils.Of("Environment with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?"),
			ExternalId:  apiutils.Of("ext-id-with-special-chars_123!@#"),
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, specialName, result.Name)
		assert.Equal(t, "Environment with special chars: !@#$%^&*()_+-=[]{}|;':\",./<>?", *result.Description)
		assert.Equal(t, "ext-id-with-special-chars_123!@#", *result.ExternalId)
		assert.NotNil(t, result.Id)
	})

	t.Run("mixed unicode and special characters", func(t *testing.T) {
		mixedName := "环境-test!@#-тест_123-🚀"
		input := &openapi.ServingEnvironment{
			Name:        mixedName,
			Description: apiutils.Of("Mixed: 测试!@# русский_test 日本語-123 🎯"),
			ExternalId:  apiutils.Of("ext-混合_test!@#-123"),
		}

		result, err := _service.UpsertServingEnvironment(input)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, mixedName, result.Name)
		assert.Equal(t, "Mixed: 测试!@# русский_test 日本語-123 🎯", *result.Description)
		assert.Equal(t, "ext-混合_test!@#-123", *result.ExternalId)
		assert.NotNil(t, result.Id)
	})

	t.Run("pagination with 10+ environments", func(t *testing.T) {
		// Create 15 environments to test pagination
		var createdEnvironments []string
		for i := 0; i < 15; i++ {
			input := &openapi.ServingEnvironment{
				Name:        fmt.Sprintf("paging-test-env-%02d", i),
				Description: apiutils.Of(fmt.Sprintf("Test environment %d for pagination", i)),
				ExternalId:  apiutils.Of(fmt.Sprintf("paging-ext-%02d", i)),
			}

			result, err := _service.UpsertServingEnvironment(input)
			require.NoError(t, err)
			require.NotNil(t, result.Id)
			createdEnvironments = append(createdEnvironments, *result.Id)
		}

		// Test first page with page size 5
		pageSize := int32(5)
		firstPageResult, err := _service.GetServingEnvironments(api.ListOptions{
			PageSize: &pageSize,
		})

		require.NoError(t, err)
		require.NotNil(t, firstPageResult)
		assert.LessOrEqual(t, len(firstPageResult.Items), int(pageSize))
		assert.Equal(t, pageSize, firstPageResult.PageSize)

		// Test second page if there's a next page token
		if firstPageResult.NextPageToken != "" {
			secondPageResult, err := _service.GetServingEnvironments(api.ListOptions{
				PageSize:      &pageSize,
				NextPageToken: &firstPageResult.NextPageToken,
			})

			require.NoError(t, err)
			require.NotNil(t, secondPageResult)
			assert.LessOrEqual(t, len(secondPageResult.Items), int(pageSize))
			assert.Equal(t, pageSize, secondPageResult.PageSize)

			// Verify no duplicate environments between pages
			firstPageIds := make(map[string]bool)
			for _, env := range firstPageResult.Items {
				firstPageIds[*env.Id] = true
			}

			for _, env := range secondPageResult.Items {
				assert.False(t, firstPageIds[*env.Id], "Environment %s appears in both pages", *env.Id)
			}
		}

		// Test larger page size to get more environments
		largePageSize := int32(100)
		largePageResult, err := _service.GetServingEnvironments(api.ListOptions{
			PageSize: &largePageSize,
		})

		require.NoError(t, err)
		require.NotNil(t, largePageResult)
		assert.GreaterOrEqual(t, len(largePageResult.Items), 15) // Should include our 15 environments
		assert.Equal(t, largePageSize, largePageResult.PageSize)

		// Verify our created environments are in the results
		resultIds := make(map[string]bool)
		for _, env := range largePageResult.Items {
			resultIds[*env.Id] = true
		}

		foundCount := 0
		for _, createdId := range createdEnvironments {
			if resultIds[createdId] {
				foundCount++
			}
		}
		assert.Equal(t, 15, foundCount, "Should find all 15 created environments in the results")

		// Test ordering by name
		orderBy := "name"
		sortOrder := "ASC"
		orderedResult, err := _service.GetServingEnvironments(api.ListOptions{
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
		orderedDescResult, err := _service.GetServingEnvironments(api.ListOptions{
			PageSize:  &largePageSize,
			OrderBy:   &orderBy,
			SortOrder: &sortOrderDesc,
		})

		require.NoError(t, err)
		require.NotNil(t, orderedDescResult)
		assert.Greater(t, len(orderedDescResult.Items), 0)
	})
}

func TestGetServingEnvironmentById(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// First create a serving environment to retrieve
		input := &openapi.ServingEnvironment{
			Name:        "get-test-serving-env",
			Description: apiutils.Of("Test description"),
			ExternalId:  apiutils.Of("get-ext-123"),
		}

		created, err := _service.UpsertServingEnvironment(input)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get the serving environment by ID
		result, err := _service.GetServingEnvironmentById(*created.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "get-test-serving-env", result.Name)
		assert.Equal(t, "get-ext-123", *result.ExternalId)
		assert.Equal(t, "Test description", *result.Description)
	})

	t.Run("invalid id", func(t *testing.T) {
		result, err := _service.GetServingEnvironmentById("invalid")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid syntax: bad request")
	})

	t.Run("non-existent id", func(t *testing.T) {
		result, err := _service.GetServingEnvironmentById("99999")

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no serving environment found")
	})
}

func TestGetServingEnvironmentByParams(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get by name", func(t *testing.T) {
		input := &openapi.ServingEnvironment{
			Name:       "params-test-serving-env",
			ExternalId: apiutils.Of("params-ext-123"),
		}
		created, err := _service.UpsertServingEnvironment(input)
		require.NoError(t, err)

		// Get by name
		envName := "params-test-serving-env"
		result, err := _service.GetServingEnvironmentByParams(&envName, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-test-serving-env", result.Name)
	})

	t.Run("successful get by external id", func(t *testing.T) {
		input := &openapi.ServingEnvironment{
			Name:       "params-ext-test-serving-env",
			ExternalId: apiutils.Of("params-unique-ext-456"),
		}
		created, err := _service.UpsertServingEnvironment(input)
		require.NoError(t, err)

		// Get by external ID
		externalId := "params-unique-ext-456"
		result, err := _service.GetServingEnvironmentByParams(nil, &externalId)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, *created.Id, *result.Id)
		assert.Equal(t, "params-unique-ext-456", *result.ExternalId)
	})

	t.Run("invalid parameters", func(t *testing.T) {
		result, err := _service.GetServingEnvironmentByParams(nil, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "invalid parameters call")
	})

	t.Run("no environment found", func(t *testing.T) {
		envName := "nonexistent-serving-env"
		result, err := _service.GetServingEnvironmentByParams(&envName, nil)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Contains(t, err.Error(), "no serving environment found")
	})
}

func TestGetServingEnvironments(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful list", func(t *testing.T) {
		// Create multiple serving environments for listing
		testEnvironments := []*openapi.ServingEnvironment{
			{Name: "list-serving-env-1", ExternalId: apiutils.Of("list-ext-1")},
			{Name: "list-serving-env-2", ExternalId: apiutils.Of("list-ext-2")},
			{Name: "list-serving-env-3", ExternalId: apiutils.Of("list-ext-3")},
		}

		var createdIds []string
		for _, env := range testEnvironments {
			created, err := _service.UpsertServingEnvironment(env)
			require.NoError(t, err)
			createdIds = append(createdIds, *created.Id)
		}

		// List serving environments with basic pagination
		pageSize := int32(10)
		listOptions := api.ListOptions{
			PageSize: &pageSize,
		}

		result, err := _service.GetServingEnvironments(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // Should have at least our 3 test environments
		assert.Equal(t, int32(10), result.PageSize)

		// Verify our environments are in the result
		foundEnvironments := 0
		for _, item := range result.Items {
			for _, createdId := range createdIds {
				if *item.Id == createdId {
					foundEnvironments++
					break
				}
			}
		}
		assert.Equal(t, 3, foundEnvironments, "All created environments should be found in the list")
	})

	t.Run("pagination and ordering", func(t *testing.T) {
		// Create several environments for pagination testing
		for i := 0; i < 5; i++ {
			env := &openapi.ServingEnvironment{
				Name:       "pagination-serving-env-" + string(rune('A'+i)),
				ExternalId: apiutils.Of("pagination-ext-" + string(rune('A'+i))),
			}
			_, err := _service.UpsertServingEnvironment(env)
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

		result, err := _service.GetServingEnvironments(listOptions)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 2) // Should have at least 2 items
		assert.Equal(t, int32(2), result.PageSize)
	})
}

func TestServingEnvironmentRoundTrip(t *testing.T) {
	_service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("complete roundtrip", func(t *testing.T) {
		// Create a serving environment with all fields
		original := &openapi.ServingEnvironment{
			Name:        "roundtrip-serving-env",
			Description: apiutils.Of("Roundtrip test description"),
			ExternalId:  apiutils.Of("roundtrip-ext-123"),
		}

		// Create
		created, err := _service.UpsertServingEnvironment(original)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := _service.GetServingEnvironmentById(*created.Id)
		require.NoError(t, err)

		// Verify all fields match
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, original.Name, retrieved.Name)
		assert.Equal(t, *original.Description, *retrieved.Description)
		assert.Equal(t, *original.ExternalId, *retrieved.ExternalId)

		// Update
		retrieved.Description = apiutils.Of("Updated description")
		retrieved.ExternalId = apiutils.Of("updated-ext-456")

		updated, err := _service.UpsertServingEnvironment(retrieved)
		require.NoError(t, err)

		// Verify update
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-ext-456", *updated.ExternalId)

		// Get again to verify persistence
		final, err := _service.GetServingEnvironmentById(*created.Id)
		require.NoError(t, err)
		assert.Equal(t, "Updated description", *final.Description)
		assert.Equal(t, "updated-ext-456", *final.ExternalId)
	})

	t.Run("roundtrip with custom properties", func(t *testing.T) {
		customProps := map[string]openapi.MetadataValue{
			"environment": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "production",
				},
			},
			"max_replicas": {
				MetadataIntValue: &openapi.MetadataIntValue{
					IntValue: "10",
				},
			},
		}

		original := &openapi.ServingEnvironment{
			Name:             "roundtrip-custom-props-env",
			CustomProperties: &customProps,
		}

		// Create
		created, err := _service.UpsertServingEnvironment(original)
		require.NoError(t, err)
		require.NotNil(t, created.Id)

		// Get by ID
		retrieved, err := _service.GetServingEnvironmentById(*created.Id)
		require.NoError(t, err)

		// Verify custom properties
		assert.NotNil(t, retrieved.CustomProperties)
		retrievedProps := *retrieved.CustomProperties
		assert.Contains(t, retrievedProps, "environment")
		assert.Contains(t, retrievedProps, "max_replicas")
		assert.Equal(t, "production", retrievedProps["environment"].MetadataStringValue.StringValue)
		assert.Equal(t, "10", retrievedProps["max_replicas"].MetadataIntValue.IntValue)

		// Update custom properties
		updatedProps := map[string]openapi.MetadataValue{
			"environment": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue: "staging",
				},
			},
			"max_replicas": {
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

		updated, err := _service.UpsertServingEnvironment(retrieved)
		require.NoError(t, err)

		// Verify updated custom properties
		assert.NotNil(t, updated.CustomProperties)
		finalProps := *updated.CustomProperties
		assert.Equal(t, "staging", finalProps["environment"].MetadataStringValue.StringValue)
		assert.Equal(t, "5", finalProps["max_replicas"].MetadataIntValue.IntValue)
		assert.Equal(t, "new_value", finalProps["new_prop"].MetadataStringValue.StringValue)
	})
}
