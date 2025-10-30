package core_test

import (
	"fmt"
	"strconv"
	"strings"
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
			CustomProperties: map[string]openapi.MetadataValue{
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
		assert.Contains(t, result.CustomProperties, "run_type")
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

	t.Run("error on EndTimeSinceEpoch less than StartTimeSinceEpoch", func(t *testing.T) {
		experimentRun := &openapi.ExperimentRun{
			Name:                apiutils.Of("test-run-time-validation"),
			StartTimeSinceEpoch: apiutils.Of("1234567890"), // Start time is later
			EndTimeSinceEpoch:   apiutils.Of("1234567000"), // End time is earlier (invalid)
		}
		_, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EndTimeSinceEpoch (1234567000) cannot be less than StartTimeSinceEpoch (1234567890)")
	})

	t.Run("successful update with valid timestamps", func(t *testing.T) {
		// First create an experiment run
		experimentRun := &openapi.ExperimentRun{
			Name:                apiutils.Of("test-run-for-update"),
			StartTimeSinceEpoch: apiutils.Of("1234567890"),
		}
		created, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)
		require.NoError(t, err)

		// Update with valid end time
		created.EndTimeSinceEpoch = apiutils.Of("1234568000") // End time is after start time
		updated, err := service.UpsertExperimentRun(created, parentExperiment.Id)
		require.NoError(t, err)
		assert.Equal(t, "1234568000", *updated.EndTimeSinceEpoch)
	})

	t.Run("error on update with invalid EndTimeSinceEpoch", func(t *testing.T) {
		// First create an experiment run
		experimentRun := &openapi.ExperimentRun{
			Name:                apiutils.Of("test-run-for-invalid-update"),
			StartTimeSinceEpoch: apiutils.Of("1234567890"),
		}
		created, err := service.UpsertExperimentRun(experimentRun, parentExperiment.Id)
		require.NoError(t, err)

		// Try to update with invalid end time
		created.EndTimeSinceEpoch = apiutils.Of("1234567000") // End time is before start time
		_, err = service.UpsertExperimentRun(created, parentExperiment.Id)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "EndTimeSinceEpoch (1234567000) cannot be less than StartTimeSinceEpoch (1234567890)")
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

	t.Run("same experiment run name across different experiments", func(t *testing.T) {
		// This test catches the bug where ExperimentID was not being used to filter experiment runs

		// Create third experiment
		experiment3 := &openapi.Experiment{
			Name: "experiment-with-shared-run-1",
		}
		createdExperiment1, err := service.UpsertExperiment(experiment3)
		require.NoError(t, err)

		// Create fourth experiment
		experiment4 := &openapi.Experiment{
			Name: "experiment-with-shared-run-2",
		}
		createdExperiment2, err := service.UpsertExperiment(experiment4)
		require.NoError(t, err)

		// Create experiment run "shared-run-name-test" for the first experiment
		run1 := &openapi.ExperimentRun{
			Name:        apiutils.Of("shared-run-name-test"),
			Description: apiutils.Of("Run for experiment 1"),
		}
		createdRun1, err := service.UpsertExperimentRun(run1, createdExperiment1.Id)
		require.NoError(t, err)

		// Create experiment run "shared-run-name-test" for the second experiment
		run2 := &openapi.ExperimentRun{
			Name:        apiutils.Of("shared-run-name-test"),
			Description: apiutils.Of("Run for experiment 2"),
		}
		createdRun2, err := service.UpsertExperimentRun(run2, createdExperiment2.Id)
		require.NoError(t, err)

		// Query for run "shared-run-name-test" of the first experiment
		runName := "shared-run-name-test"
		result1, err := service.GetExperimentRunByParams(&runName, createdExperiment1.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result1)
		assert.Equal(t, *createdRun1.Id, *result1.Id)
		assert.Equal(t, *createdExperiment1.Id, result1.ExperimentId)
		assert.Equal(t, "Run for experiment 1", *result1.Description)

		// Query for run "shared-run-name-test" of the second experiment
		result2, err := service.GetExperimentRunByParams(&runName, createdExperiment2.Id, nil)
		require.NoError(t, err)
		require.NotNil(t, result2)
		assert.Equal(t, *createdRun2.Id, *result2.Id)
		assert.Equal(t, *createdExperiment2.Id, result2.ExperimentId)
		assert.Equal(t, "Run for experiment 2", *result2.Description)

		// Ensure we got different runs
		assert.NotEqual(t, *result1.Id, *result2.Id)
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

	t.Run("error on non-existent experiment id", func(t *testing.T) {
		// Use a valid ID format but for an experiment that doesn't exist
		nonExistentId := "999999"
		_, err := service.GetExperimentRuns(api.ListOptions{}, &nonExistentId)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "not found")
	})

	t.Run("error on empty experiment id", func(t *testing.T) {
		// Empty string should be invalid
		emptyId := ""
		_, err := service.GetExperimentRuns(api.ListOptions{}, &emptyId)
		assert.Error(t, err)
		// Should fail ID validation
		assert.Contains(t, err.Error(), "invalid")
	})

	t.Run("error on negative experiment id", func(t *testing.T) {
		// Negative number should either be invalid or not found
		negativeId := "-1"
		_, err := service.GetExperimentRuns(api.ListOptions{}, &negativeId)
		assert.Error(t, err)
		// Could be either invalid ID or not found depending on validation
		errorMsg := err.Error()
		assert.True(t, strings.Contains(errorMsg, "invalid") || strings.Contains(errorMsg, "not found"),
			"Error should contain 'invalid' or 'not found', but got: %s", errorMsg)
	})
}

func TestGetExperimentRunsWithFilterQuery(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create an experiment to associate runs with
	experiment := &openapi.Experiment{
		Name: "test-experiment-for-runs",
	}
	createdExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create test experiment runs with various properties for filtering
	testRuns := []struct {
		run *openapi.ExperimentRun
	}{
		{
			run: &openapi.ExperimentRun{
				Name:                apiutils.Of("pytorch-run-1"),
				Description:         apiutils.Of("PyTorch training run with hyperparameter tuning"),
				ExternalId:          apiutils.Of("ext-pytorch-run-001"),
				Status:              (*openapi.ExperimentRunStatus)(apiutils.Of("FINISHED")),
				StartTimeSinceEpoch: apiutils.Of("1700000000"),
				EndTimeSinceEpoch:   apiutils.Of("1700003600"),
				ExperimentId:        *createdExperiment.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"framework": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "pytorch",
							MetadataType: "MetadataStringValue",
						},
					},
					"learning_rate": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.001,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"epochs": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "100",
							MetadataType: "MetadataIntValue",
						},
					},
					"gpu_enabled": {
						MetadataBoolValue: &openapi.MetadataBoolValue{
							BoolValue:    true,
							MetadataType: "MetadataBoolValue",
						},
					},
				},
			},
		},
		{
			run: &openapi.ExperimentRun{
				Name:                apiutils.Of("tensorflow-run-2"),
				Description:         apiutils.Of("TensorFlow training run for NLP model"),
				ExternalId:          apiutils.Of("ext-tf-run-002"),
				Status:              (*openapi.ExperimentRunStatus)(apiutils.Of("FAILED")),
				StartTimeSinceEpoch: apiutils.Of("1700010000"),
				ExperimentId:        *createdExperiment.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"framework": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "tensorflow",
							MetadataType: "MetadataStringValue",
						},
					},
					"learning_rate": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.01,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"epochs": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "50",
							MetadataType: "MetadataIntValue",
						},
					},
					"gpu_enabled": {
						MetadataBoolValue: &openapi.MetadataBoolValue{
							BoolValue:    false,
							MetadataType: "MetadataBoolValue",
						},
					},
				},
			},
		},
		{
			run: &openapi.ExperimentRun{
				Name:                apiutils.Of("pytorch-run-2"),
				Description:         apiutils.Of("PyTorch run with distributed training"),
				ExternalId:          apiutils.Of("ext-pytorch-run-003"),
				Status:              (*openapi.ExperimentRunStatus)(apiutils.Of("RUNNING")),
				StartTimeSinceEpoch: apiutils.Of("1700020000"),
				ExperimentId:        *createdExperiment.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"framework": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "pytorch",
							MetadataType: "MetadataStringValue",
						},
					},
					"learning_rate": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.0001,
							MetadataType: "MetadataDoubleValue",
						},
					},
					"epochs": {
						MetadataIntValue: &openapi.MetadataIntValue{
							IntValue:     "200",
							MetadataType: "MetadataIntValue",
						},
					},
					"distributed": {
						MetadataBoolValue: &openapi.MetadataBoolValue{
							BoolValue:    true,
							MetadataType: "MetadataBoolValue",
						},
					},
				},
			},
		},
		{
			run: &openapi.ExperimentRun{
				Name:         apiutils.Of("sklearn-run"),
				Description:  apiutils.Of("Scikit-learn baseline model"),
				ExternalId:   apiutils.Of("ext-sklearn-run-004"),
				Status:       (*openapi.ExperimentRunStatus)(apiutils.Of("FINISHED")),
				ExperimentId: *createdExperiment.Id,
				CustomProperties: map[string]openapi.MetadataValue{
					"framework": {
						MetadataStringValue: &openapi.MetadataStringValue{
							StringValue:  "sklearn",
							MetadataType: "MetadataStringValue",
						},
					},
					"learning_rate": {
						MetadataDoubleValue: &openapi.MetadataDoubleValue{
							DoubleValue:  0.1,
							MetadataType: "MetadataDoubleValue",
						},
					},
				},
			},
		},
	}

	// Create all test runs
	createdRuns := make([]*openapi.ExperimentRun, 0)
	for _, tr := range testRuns {
		created, err := service.UpsertExperimentRun(tr.run, createdExperiment.Id)
		require.NoError(t, err)
		createdRuns = append(createdRuns, created)
	}

	// Debug: Check if names were saved correctly
	for i, created := range createdRuns {
		fetched, err := service.GetExperimentRunById(*created.Id)
		require.NoError(t, err)
		nameStr := "<nil>"
		if created.Name != nil {
			nameStr = *created.Name
		}
		fetchedNameStr := "<nil>"
		if fetched.Name != nil {
			fetchedNameStr = *fetched.Name
		}
		t.Logf("Created run %d: ID=%s, Name=%s", i, *created.Id, nameStr)
		t.Logf("Fetched run %d: ID=%s, Name=%s", i, *fetched.Id, fetchedNameStr)
		if testRuns[i].run.Name != nil {
			require.NotNil(t, fetched.Name, "Fetched run should have a name")
			require.Equal(t, *testRuns[i].run.Name, *fetched.Name, "Fetched run name should match")
		}
	}

	testCases := []struct {
		name          string
		filterQuery   string
		expectedCount int
		expectedNames []string
	}{
		{
			name:          "Filter by exact name",
			filterQuery:   "name = 'pytorch-run-1'",
			expectedCount: 1,
			expectedNames: []string{"pytorch-run-1"},
		},
		{
			name:          "Get all runs (no filter)",
			filterQuery:   "",
			expectedCount: 4,
			expectedNames: []string{"pytorch-run-1", "tensorflow-run-2", "pytorch-run-2", "sklearn-run"},
		},
		{
			name:          "Filter by name pattern",
			filterQuery:   "name LIKE 'pytorch-%'",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "pytorch-run-2"},
		},
		{
			name:          "Filter by description",
			filterQuery:   "description LIKE '%NLP%'",
			expectedCount: 1,
			expectedNames: []string{"tensorflow-run-2"},
		},
		{
			name:          "Filter by external ID",
			filterQuery:   "externalId = 'ext-tf-run-002'",
			expectedCount: 1,
			expectedNames: []string{"tensorflow-run-2"},
		},
		{
			name:          "Filter by status",
			filterQuery:   "status = 'FINISHED'",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "sklearn-run"},
		},
		{
			name:          "Filter by custom property - string",
			filterQuery:   "framework = 'pytorch'",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "pytorch-run-2"},
		},
		{
			name:          "Filter by custom property - numeric comparison",
			filterQuery:   "learning_rate < 0.01",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "pytorch-run-2"},
		},
		{
			name:          "Filter by custom property - integer",
			filterQuery:   "epochs >= 100",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "pytorch-run-2"},
		},
		{
			name:          "Filter by custom property - boolean",
			filterQuery:   "gpu_enabled = true",
			expectedCount: 1,
			expectedNames: []string{"pytorch-run-1"},
		},
		{
			name:          "Complex filter with AND",
			filterQuery:   "framework = 'pytorch' AND epochs > 150",
			expectedCount: 1,
			expectedNames: []string{"pytorch-run-2"},
		},
		{
			name:          "Complex filter with OR",
			filterQuery:   "status = 'FAILED' OR status = 'RUNNING'",
			expectedCount: 2,
			expectedNames: []string{"tensorflow-run-2", "pytorch-run-2"},
		},
		{
			name:          "Complex filter with parentheses",
			filterQuery:   "(framework = 'pytorch' OR framework = 'tensorflow') AND learning_rate <= 0.001",
			expectedCount: 2,
			expectedNames: []string{"pytorch-run-1", "pytorch-run-2"},
		},
		{
			name:          "Case insensitive pattern matching",
			filterQuery:   "name ILIKE '%RUN%'",
			expectedCount: 4,
			expectedNames: []string{"pytorch-run-1", "tensorflow-run-2", "pytorch-run-2", "sklearn-run"},
		},
		{
			name:          "Filter with NOT condition",
			filterQuery:   "status != 'FINISHED'",
			expectedCount: 2,
			expectedNames: []string{"tensorflow-run-2", "pytorch-run-2"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			pageSize := int32(10)
			listOptions := api.ListOptions{
				PageSize:    &pageSize,
				FilterQuery: &tc.filterQuery,
			}

			// Debug: Log the filter query
			t.Logf("Running filter query: %s", tc.filterQuery)
			result, err := service.GetExperimentRuns(listOptions, nil)

			require.NoError(t, err)
			require.NotNil(t, result)

			// Extract names from results
			var actualNames []string
			for _, item := range result.Items {
				for _, expectedName := range tc.expectedNames {
					if *item.Name == expectedName {
						actualNames = append(actualNames, *item.Name)
						break
					}
				}
			}

			assert.Equal(t, tc.expectedCount, len(actualNames),
				"Expected %d runs for filter '%s', but got %d",
				tc.expectedCount, tc.filterQuery, len(actualNames))

			// Verify the expected runs are present
			assert.ElementsMatch(t, tc.expectedNames, actualNames,
				"Expected runs %v for filter '%s', but got %v",
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

		result, err := service.GetExperimentRuns(listOptions, nil)

		if assert.Error(t, err) {
			assert.Nil(t, result)
			assert.Contains(t, err.Error(), "invalid filter query")
		}
	})

	// Test combining filterQuery with experimentId parameter
	t.Run("Filter with experimentId parameter", func(t *testing.T) {
		// Create another experiment with runs
		anotherExperiment := &openapi.Experiment{
			Name: "another-experiment",
		}
		anotherCreatedExperiment, err := service.UpsertExperiment(anotherExperiment)
		require.NoError(t, err)

		anotherRun := &openapi.ExperimentRun{
			Name:         apiutils.Of("another-pytorch-run"),
			ExperimentId: *anotherCreatedExperiment.Id,
			CustomProperties: map[string]openapi.MetadataValue{
				"framework": {
					MetadataStringValue: &openapi.MetadataStringValue{
						StringValue:  "pytorch",
						MetadataType: "MetadataStringValue",
					},
				},
			},
		}
		_, err = service.UpsertExperimentRun(anotherRun, anotherCreatedExperiment.Id)
		require.NoError(t, err)

		// Filter by framework=pytorch should return runs from both experiments
		pageSize := int32(10)
		filterQuery := "framework = 'pytorch'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		// Without experimentId - should get 3 (2 from first experiment + 1 from second)
		allResult, err := service.GetExperimentRuns(listOptions, nil)
		require.NoError(t, err)
		assert.Equal(t, 3, len(allResult.Items))

		// With experimentId - should only get 2 from first experiment
		filteredResult, err := service.GetExperimentRuns(listOptions, createdExperiment.Id)
		require.NoError(t, err)
		assert.Equal(t, 2, len(filteredResult.Items))
		for _, item := range filteredResult.Items {
			assert.Equal(t, *createdExperiment.Id, item.ExperimentId)
		}
	})

	// Test combining filterQuery with pagination
	t.Run("Filter with pagination", func(t *testing.T) {
		pageSize := int32(1)
		filterQuery := "framework = 'pytorch'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		// Get first page
		firstPage, err := service.GetExperimentRuns(listOptions, createdExperiment.Id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(firstPage.Items))
		assert.NotEmpty(t, firstPage.NextPageToken)

		// Get second page
		listOptions.NextPageToken = &firstPage.NextPageToken
		secondPage, err := service.GetExperimentRuns(listOptions, createdExperiment.Id)
		require.NoError(t, err)
		assert.Equal(t, 1, len(secondPage.Items))

		// Ensure different items on each page
		assert.NotEqual(t, firstPage.Items[0].Id, secondPage.Items[0].Id)
	})

	// Test empty results
	t.Run("Filter with no matches", func(t *testing.T) {
		pageSize := int32(10)
		filterQuery := "framework = 'nonexistent'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		result, err := service.GetExperimentRuns(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
	})

	// Test filtering with time range
	t.Run("Filter by time range", func(t *testing.T) {
		pageSize := int32(10)
		filterQuery := "startTimeSinceEpoch >= '1700010000'"
		listOptions := api.ListOptions{
			PageSize:    &pageSize,
			FilterQuery: &filterQuery,
		}

		result, err := service.GetExperimentRuns(listOptions, nil)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items))

		// Verify all results have start time >= 1700010000
		for _, item := range result.Items {
			startTime, err := strconv.ParseInt(*item.StartTimeSinceEpoch, 10, 64)
			require.NoError(t, err)
			assert.GreaterOrEqual(t, startTime, int64(1700010000))
		}
	})
}
