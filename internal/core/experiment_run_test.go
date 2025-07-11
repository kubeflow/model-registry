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

func TestUpsertExperimentRun(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create a parent experiment first
	experiment := &openapi.Experiment{
		Name: "parent-experiment",
	}
	parentExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	t.Run("successful create", func(t *testing.T) {
		experimentRun := &openapi.ExperimentRun{
			Name:                apiutils.Of("test-experiment-run"),
			Description:         apiutils.Of("Test experiment run description"),
			Owner:               apiutils.Of("test-owner"),
			ExternalId:          apiutils.Of("exp-run-ext-123"),
			State:               apiutils.Of(openapi.EXPERIMENTRUNSTATE_LIVE),
			Status:              apiutils.Of(openapi.EXPERIMENTRUNSTATUS_RUNNING),
			StartTimeSinceEpoch: apiutils.Of("1234567890"),
			EndTimeSinceEpoch:   apiutils.Of("1234567999"),
			CustomProperties: &map[string]openapi.MetadataValue{
				"run_type": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "training",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}

		result, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.NotNil(t, result.Id)
		assert.Contains(t, *result.Name, "test-experiment-run") // Name might be prefixed with experiment ID
		assert.Equal(t, "exp-run-ext-123", *result.ExternalId)
		assert.Equal(t, "Test experiment run description", *result.Description)
		assert.Equal(t, "test-owner", *result.Owner)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATE_LIVE, *result.State)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATUS_RUNNING, *result.Status)
		assert.Equal(t, "1234567890", *result.StartTimeSinceEpoch)
		assert.Equal(t, "1234567999", *result.EndTimeSinceEpoch)
		assert.Equal(t, *parentExperiment.Id, result.ExperimentId)
		assert.NotNil(t, result.CreateTimeSinceEpoch)
		assert.NotNil(t, result.LastUpdateTimeSinceEpoch)
		assert.NotNil(t, result.CustomProperties)
		assert.Contains(t, *result.CustomProperties, "run_type")
	})

	t.Run("successful update", func(t *testing.T) {
		// Create initial experiment run
		experimentRun := &openapi.ExperimentRun{
			Name:        apiutils.Of("update-test-run"),
			Description: apiutils.Of("Original description"),
			Owner:       apiutils.Of("original-owner"),
		}

		created, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)
		require.NoError(t, err)

		// Update the experiment run
		created.Description = apiutils.Of("Updated description")
		created.Owner = apiutils.Of("updated-owner")
		created.State = apiutils.Of(openapi.EXPERIMENTRUNSTATE_ARCHIVED)
		created.Status = apiutils.Of(openapi.EXPERIMENTRUNSTATUS_FINISHED)

		updated, err := service.UpsertExperimentRun(created, parentExperiment.Id)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *updated.Id)
		assert.Equal(t, "Updated description", *updated.Description)
		assert.Equal(t, "updated-owner", *updated.Owner)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATE_ARCHIVED, *updated.State)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATUS_FINISHED, *updated.Status)
	})

	t.Run("error on nil experiment run", func(t *testing.T) {
		_, err := service.UpsertExperimentRun(nil, parentExperiment.Id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment run pointer")
	})

	t.Run("error on nil experiment id", func(t *testing.T) {
		experimentRun := &openapi.ExperimentRun{
			Name: apiutils.Of("test-run"),
		}
		_, err := service.UpsertExperimentRun(experimentRun, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experiment ID is required")
	})

	t.Run("error on invalid experiment id", func(t *testing.T) {
		experimentRun := &openapi.ExperimentRun{
			Name: apiutils.Of("test-run"),
		}
		invalidId := "invalid-id"
		_, err := service.UpsertExperimentRun(experimentRun, &invalidId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment ID")
	})

	t.Run("error on non-existent experiment id", func(t *testing.T) {
		experimentRun := &openapi.ExperimentRun{
			Name: apiutils.Of("test-run"),
		}
		nonExistentId := "999999"
		_, err := service.UpsertExperimentRun(experimentRun, &nonExistentId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "experiment not found")
	})
}

func TestGetExperimentRunById(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create a parent experiment
	experiment := &openapi.Experiment{
		Name: "parent-experiment-get",
	}
	parentExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	t.Run("successful get", func(t *testing.T) {
		// Create an experiment run
		experimentRun := &openapi.ExperimentRun{
			Name:        apiutils.Of("get-test-experiment-run"),
			Description: apiutils.Of("Get test description"),
			Owner:       apiutils.Of("get-test-owner"),
			State:       apiutils.Of(openapi.EXPERIMENTRUNSTATE_LIVE),
			Status:      apiutils.Of(openapi.EXPERIMENTRUNSTATUS_RUNNING),
		}

		created, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)
		require.NoError(t, err)

		// Get the experiment run by ID
		retrieved, err := service.GetExperimentRunById(*created.Id)

		require.NoError(t, err)
		assert.Equal(t, *created.Id, *retrieved.Id)
		assert.Equal(t, *created.Name, *retrieved.Name)
		assert.Equal(t, "Get test description", *retrieved.Description)
		assert.Equal(t, "get-test-owner", *retrieved.Owner)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATE_LIVE, *retrieved.State)
		assert.Equal(t, openapi.EXPERIMENTRUNSTATUS_RUNNING, *retrieved.Status)
		assert.Equal(t, *parentExperiment.Id, retrieved.ExperimentId)
	})

	t.Run("error on non-existent id", func(t *testing.T) {
		_, err := service.GetExperimentRunById("999999")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("error on invalid id", func(t *testing.T) {
		_, err := service.GetExperimentRunById("invalid-id")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment run ID")
	})
}

func TestGetExperimentRunByParams(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create parent experiments
	experiment1 := &openapi.Experiment{
		Name: "parent-experiment-params-1",
	}
	parentExperiment1, err := service.UpsertExperiment(experiment1)
	require.NoError(t, err)

	experiment2 := &openapi.Experiment{
		Name: "parent-experiment-params-2",
	}
	parentExperiment2, err := service.UpsertExperiment(experiment2)
	require.NoError(t, err)

	// Create test experiment runs
	experimentRun1 := &openapi.ExperimentRun{
		Name:       apiutils.Of("params-test-run-1"),
		ExternalId: apiutils.Of("params-run-ext-1"),
	}
	created1, err := service.UpsertExperimentRun(experimentRun1, parentExperiment1.Id)
	require.NoError(t, err)

	experimentRun2 := &openapi.ExperimentRun{
		Name:       apiutils.Of("params-test-run-2"),
		ExternalId: apiutils.Of("params-run-ext-2"),
	}
	created2, err := service.UpsertExperimentRun(experimentRun2, parentExperiment2.Id)
	require.NoError(t, err)

	t.Run("get by name and experiment id", func(t *testing.T) {
		result, err := service.GetExperimentRunByParams(apiutils.Of("params-test-run-1"), parentExperiment1.Id, nil)
		require.NoError(t, err)
		assert.Equal(t, *created1.Id, *result.Id)
		assert.Contains(t, *result.Name, "params-test-run-1")
	})

	t.Run("get by external id", func(t *testing.T) {
		result, err := service.GetExperimentRunByParams(nil, nil, apiutils.Of("params-run-ext-2"))
		require.NoError(t, err)
		assert.Equal(t, *created2.Id, *result.Id)
		assert.Equal(t, "params-run-ext-2", *result.ExternalId)
	})

	t.Run("error on no params", func(t *testing.T) {
		_, err := service.GetExperimentRunByParams(nil, nil, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "supply either (name and experimentId), or externalId")
	})

	t.Run("error on not found", func(t *testing.T) {
		_, err := service.GetExperimentRunByParams(apiutils.Of("non-existent"), parentExperiment1.Id, nil)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})
}

func TestGetExperimentRuns(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create parent experiments
	experiment1 := &openapi.Experiment{
		Name: "parent-experiment-list-1",
	}
	parentExperiment1, err := service.UpsertExperiment(experiment1)
	require.NoError(t, err)

	experiment2 := &openapi.Experiment{
		Name: "parent-experiment-list-2",
	}
	parentExperiment2, err := service.UpsertExperiment(experiment2)
	require.NoError(t, err)

	// Create multiple test experiment runs for first experiment
	for i := 0; i < 3; i++ {
		experimentRun := &openapi.ExperimentRun{
			Name:        apiutils.Of(fmt.Sprintf("list-test-run-exp1-%d", i)),
			Description: apiutils.Of(fmt.Sprintf("List test description exp1-%d", i)),
			Owner:       apiutils.Of("list-test-owner"),
		}
		_, err := service.UpsertExperimentRun(experimentRun, parentExperiment1.Id)
		require.NoError(t, err)
	}

	// Create experiment runs for second experiment
	for i := 0; i < 2; i++ {
		experimentRun := &openapi.ExperimentRun{
			Name:        apiutils.Of(fmt.Sprintf("list-test-run-exp2-%d", i)),
			Description: apiutils.Of(fmt.Sprintf("List test description exp2-%d", i)),
			Owner:       apiutils.Of("list-test-owner"),
		}
		_, err := service.UpsertExperimentRun(experimentRun, parentExperiment2.Id)
		require.NoError(t, err)
	}

	t.Run("get all experiment runs", func(t *testing.T) {
		result, err := service.GetExperimentRuns(api.ListOptions{}, nil)
		require.NoError(t, err)
		assert.GreaterOrEqual(t, int(result.Size), 5)
		assert.GreaterOrEqual(t, len(result.Items), 5)
	})

	t.Run("get experiment runs for specific experiment", func(t *testing.T) {
		result, err := service.GetExperimentRuns(api.ListOptions{}, parentExperiment1.Id)
		require.NoError(t, err)
		assert.Equal(t, int32(3), result.Size)
		assert.Equal(t, 3, len(result.Items))

		// Verify all returned runs belong to the correct experiment
		for _, run := range result.Items {
			assert.Equal(t, *parentExperiment1.Id, run.ExperimentId)
		}
	})

	t.Run("get experiment runs with pagination", func(t *testing.T) {
		pageSize := int32(2)
		result, err := service.GetExperimentRuns(api.ListOptions{
			PageSize: &pageSize,
		}, parentExperiment1.Id)
		require.NoError(t, err)
		assert.Equal(t, pageSize, result.PageSize)
		assert.Equal(t, pageSize, result.Size)
		assert.Equal(t, int(pageSize), len(result.Items))
	})

	t.Run("get experiment runs with ordering", func(t *testing.T) {
		orderBy := "CREATE_TIME"
		sortOrder := "DESC"
		result, err := service.GetExperimentRuns(api.ListOptions{
			OrderBy:   &orderBy,
			SortOrder: &sortOrder,
		}, parentExperiment1.Id)
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

	t.Run("error on invalid experiment id", func(t *testing.T) {
		invalidId := "invalid-id"
		_, err := service.GetExperimentRuns(api.ListOptions{}, &invalidId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid experiment ID")
	})
}
