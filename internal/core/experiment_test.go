package core_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUpsertExperiment(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful create", func(t *testing.T) {
		experiment := &openapi.Experiment{
			Name:        "test-experiment",
			Description: apiutils.Of("Test experiment description"),
			Owner:       apiutils.Of("test-owner"),
			ExternalId:  apiutils.Of("exp-ext-123"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"project": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "ml-project",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		result, err := service.UpsertExperiment(experiment)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Equal(t, "test-experiment", result.Name)
		assert.Equal(t, "exp-ext-123", *result.ExternalId)
		assert.Equal(t, "Test experiment description", *result.Description)
		assert.Equal(t, "test-owner", *result.Owner)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
		assert.NotNil(t, result.CustomProperties)
		assert.Contains(t, *result.CustomProperties, "project")
	})

	t.Run("successful update", func(t *testing.T) {
		// Create initial experiment
		experiment := &openapi.Experiment{
			Name:        "update-test-experiment",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Update the experiment
		created.Description = apiutils.Of("Updated description")
		created.Owner = apiutils.Of("updated-owner")

		updated, err := service.UpsertExperiment(created)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-owner", *updated.Owner)
	})

	t.Run("error on nil experiment", func(t *testing.T) {
		_, err := service.UpsertExperiment(nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment pointer")
	})

	t.Run("error on duplicate name", func(t *testing.T) {
		experiment := &openapi.Experiment{
			Name: "duplicate-name-test",
		}

		_, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Try to create another experiment with the same name
		duplicate := &openapi.Experiment{
			Name: "duplicate-name-test",
		}

		_, err = service.UpsertExperiment(duplicate)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "already exists")
	})
}

func TestGetExperimentById(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("successful get", func(t *testing.T) {
		// Create an experiment
		experiment := &openapi.Experiment{
			Name:        "get-test-experiment",
			Description: apiutils.Of("Get test description"),
			Owner:       apiutils.Of("get-test-owner"),
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Get the experiment by ID
		retrieved, err := service.GetExperimentById(*created.Id)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, "get-test-experiment", retrieved.Name)
		assert.Equal(t, "Get test description", *retrieved.Description)
		assert.Equal(t, "get-test-owner", *retrieved.Owner)
	})

	t.Run("error on non-existent id", func(t *testing.T) {
		_, err := service.GetExperimentById("999999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("error on invalid id", func(t *testing.T) {
		_, err := service.GetExperimentById("invalid-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid syntax")
	})
}

func TestGetExperimentByParams(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create test experiments
	experiment1 := &openapi.Experiment{
		Name:       "params-test-experiment-1",
		ExternalId: apiutils.Of("params-exp-ext-1"),
	}
	created1, err := service.UpsertExperiment(experiment1)
	require.NoError(t, err)

	experiment2 := &openapi.Experiment{
		Name:       "params-test-experiment-2",
		ExternalId: apiutils.Of("params-exp-ext-2"),
	}
	created2, err := service.UpsertExperiment(experiment2)
	require.NoError(t, err)

	t.Run("get by name", func(t *testing.T) {
		result, err := service.GetExperimentByParams(apiutils.Of("params-test-experiment-1"), nil)
		require.NoError(t, err)
		assert.Equal(t, *created1.Id, *result.Id)
		assert.Equal(t, "params-test-experiment-1", result.Name)
	})

	t.Run("get by external id", func(t *testing.T) {
		result, err := service.GetExperimentByParams(nil, apiutils.Of("params-exp-ext-2"))
		require.NoError(t, err)
		assert.Equal(t, *created2.Id, *result.Id)
		assert.Equal(t, "params-exp-ext-2", *result.ExternalId)
	})

	t.Run("error on no params", func(t *testing.T) {
		_, err := service.GetExperimentByParams(nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supply either name or externalId")
	})

	t.Run("error on not found", func(t *testing.T) {
		_, err := service.GetExperimentByParams(apiutils.Of("non-existent-experiment"), nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestGetExperiments(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create multiple test experiments
	for i := 0; i < 5; i++ {
		experiment := &openapi.Experiment{
			Name:        fmt.Sprintf("list-test-experiment-%d", i),
			Description: apiutils.Of(fmt.Sprintf("List test description %d", i)),
			Owner:       apiutils.Of("list-test-owner"),
		}
		_, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)
	}

	t.Run("get all experiments", func(t *testing.T) {
		result, err := service.GetExperiments(api.ListOptions{})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, int(result.Size), 5)
		assert.GreaterOrEqual(t, len(result.Items), 5)
	})

	t.Run("get experiments with pagination", func(t *testing.T) {
		pageSize := int32(3)
		result, err := service.GetExperiments(api.ListOptions{
			PageSize: &pageSize,
		})
		require.NoError(t, err)
		assert.Equal(t, pageSize, result.PageSize)
		assert.Equal(t, pageSize, result.Size)
		assert.Equal(t, int(pageSize), len(result.Items))
	})

	t.Run("get experiments with ordering", func(t *testing.T) {
		orderBy := "CREATE_TIME"
		sortOrder := "DESC"
		result, err := service.GetExperiments(api.ListOptions{
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		})
		require.NoError(t, err)
		assert.Greater(t, result.Size, int32(0))

		// Verify ordering (newest first) - now with proper constants
		if len(result.Items) > 1 {
			for i := 0; i < len(result.Items)-1; i++ {
				assert.GreaterOrEqual(t, *result.Items[i].CreateTimeSinceEpoch, *result.Items[i+1].CreateTimeSinceEpoch,
					"Items should be in descending order by create time, but %s is not >= %s",
					*result.Items[i].CreateTimeSinceEpoch, *result.Items[i+1].CreateTimeSinceEpoch)
			}
		}
	})
}

func TestExperimentNonEditableFieldsProtection(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	t.Run("non-editable fields are protected during update", func(t *testing.T) {
		// Create initial experiment
		experiment := &openapi.Experiment{
			Name:        "test-non-editable",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
			ExternalId:  apiutils.Of("original-ext-id"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"original_prop": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "original_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Store original values
		originalId := *created.Id
		originalName := created.Name
		originalCreateTime := *created.CreateTimeSinceEpoch
		originalUpdateTime := *created.LastUpdateTimeSinceEpoch

		// Wait a moment to ensure timestamp difference
		time.Sleep(10 * time.Millisecond)

		// Attempt to update non-editable fields along with editable fields
		updateRequest := &openapi.Experiment{
			Id:                       created.Id,                         // This should be preserved
			Name:                     "HACKED_NAME",                      // This should be ignored (non-editable)
			CreateTimeSinceEpoch:     apiutils.Of("9999999999"),          // This should be ignored (non-editable)
			LastUpdateTimeSinceEpoch: apiutils.Of("8888888888"),          // This should be ignored (non-editable)
			Description:              apiutils.Of("Updated description"), // This should be updated (editable)
			Owner:                    apiutils.Of("updated-owner"),       // This should be updated (editable)
			ExternalId:               apiutils.Of("updated-ext-id"),      // This should be updated (editable)
			CustomProperties: &map[string]openapi.MetadataValue{
				"updated_prop": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "updated_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		updated, err := service.UpsertExperiment(updateRequest)
		require.NoError(t, err)

		// Verify non-editable fields are preserved
		assert.Equal(t, originalId, *updated.Id, "ID should not be changeable")
		assert.Equal(t, originalName, updated.Name, "Name should not be changeable")
		assert.Equal(t, originalCreateTime, *updated.CreateTimeSinceEpoch, "CreateTimeSinceEpoch should not be changeable")

		// LastUpdateTimeSinceEpoch should be updated by the system, not by user input
		assert.NotEqual(t, originalUpdateTime, *updated.LastUpdateTimeSinceEpoch, "LastUpdateTimeSinceEpoch should be updated by system")
		assert.NotEqual(t, "8888888888", *updated.LastUpdateTimeSinceEpoch, "LastUpdateTimeSinceEpoch should not use user-provided value")

		// Verify editable fields are updated
		assert.Equal(t, "Updated description", *updated.Description, "Description should be updatable")
		assert.Equal(t, "updated-owner", *updated.Owner, "Owner should be updatable")
		assert.Equal(t, "updated-ext-id", *updated.ExternalId, "ExternalId should be updatable")
		assert.Contains(t, *updated.CustomProperties, "updated_prop", "CustomProperties should be updatable")
	})

	t.Run("partial update preserves existing editable fields", func(t *testing.T) {
		// Create initial experiment with multiple editable fields
		experiment := &openapi.Experiment{
			Name:        "test-partial-update",
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
			ExternalId:  apiutils.Of("original-ext-id"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"keep_this": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "keep_value",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		created, err := service.UpsertExperiment(experiment)
		require.NoError(t, err)

		// Partial update - only update description, should preserve other editable fields
		partialUpdate := &openapi.Experiment{
			Id:          created.Id,
			Description: apiutils.Of("Updated description only"),
		}

		updated, err := service.UpsertExperiment(partialUpdate)
		require.NoError(t, err)

		// Verify partial update worked correctly
		assert.Equal(t, "Updated description only", *updated.Description, "Description should be updated")
		assert.Equal(t, "original-owner", *updated.Owner, "Owner should be preserved from existing")
		assert.Equal(t, "original-ext-id", *updated.ExternalId, "ExternalId should be preserved from existing")
		assert.Contains(t, *updated.CustomProperties, "keep_this", "CustomProperties should be preserved from existing")
	})
}
