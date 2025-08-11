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

func TestGetExperimentRunMetricHistory(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment",
		Description: apiutils.Of("Test experiment for metric history"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run"),
		Description: apiutils.Of("Test experiment run for metric history"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics for the experiment run
	metric1 := &openapi.Metric{
		Name:        apiutils.Of("accuracy"),
		Value:       apiutils.Of(0.95),
		Timestamp:   apiutils.Of("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Test accuracy metric"),
	}

	metric2 := &openapi.Metric{
		Name:        apiutils.Of("loss"),
		Value:       apiutils.Of(0.05),
		Timestamp:   apiutils.Of("1234567891"),
		Step:        apiutils.Of(int64(2)),
		Description: apiutils.Of("Test loss metric"),
	}

	// Insert metric history records
	err = service.InsertMetricHistory(metric1, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history for metric1")

	err = service.InsertMetricHistory(metric2, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history for metric2")

	// Test getting all metric history
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history")
	assert.Equal(t, int32(2), result.Size, "should return 2 metric history records")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify the metrics are returned correctly
	foundAccuracy := false
	foundLoss := false
	for _, item := range result.Items {
		switch *item.Name {
		case "accuracy":
			foundAccuracy = true
			assert.Equal(t, 0.95, *item.Value)
			assert.Equal(t, "1234567890", *item.Timestamp)
			assert.Equal(t, int64(1), *item.Step)
		case "loss":
			foundLoss = true
			assert.Equal(t, 0.05, *item.Value)
			assert.Equal(t, "1234567891", *item.Timestamp)
			assert.Equal(t, int64(2), *item.Step)
		}
	}
	assert.True(t, foundAccuracy, "should find accuracy metric")
	assert.True(t, foundLoss, "should find loss metric")
}

func TestGetExperimentRunMetricHistoryWithNameFilter(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-filter",
		Description: apiutils.Of("Test experiment for metric history name filter"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-filter"),
		Description: apiutils.Of("Test experiment run for metric history name filter"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics for the experiment run
	metric1 := &openapi.Metric{
		Name:      apiutils.Of("accuracy"),
		Value:     apiutils.Of(0.95),
		Timestamp: apiutils.Of("1234567890"),
		Step:      apiutils.Of(int64(1)),
	}

	metric2 := &openapi.Metric{
		Name:      apiutils.Of("loss"),
		Value:     apiutils.Of(0.05),
		Timestamp: apiutils.Of("1234567891"),
		Step:      apiutils.Of(int64(2)),
	}

	// Insert metric history records
	err = service.InsertMetricHistory(metric1, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history for metric1")

	err = service.InsertMetricHistory(metric2, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history for metric2")

	// Test filtering by name
	accuracyName := "accuracy"
	result, err := service.GetExperimentRunMetricHistory(&accuracyName, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history with name filter")
	assert.Equal(t, int32(1), result.Size, "should return 1 metric history record for accuracy")
	assert.Equal(t, 1, len(result.Items), "should have 1 item in the result")
	assert.Equal(t, "accuracy", *result.Items[0].Name)
	assert.Equal(t, 0.95, *result.Items[0].Value)
}

func TestInsertMetricHistory(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-insert",
		Description: apiutils.Of("Test experiment for metric history insert"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-insert"),
		Description: apiutils.Of("Test experiment run for metric history insert"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Test 1: Basic metric history insertion
	metric := &openapi.Metric{
		Name:        apiutils.Of("test_metric"),
		Value:       apiutils.Of(42.5),
		Timestamp:   apiutils.Of("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Test metric description"),
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_prop": {
				MetadataStringValue: &openapi.MetadataStringValue{
					StringValue:  "custom_value",
					MetadataType: "MetadataStringValue",
				},
			},
		},
	}

	// Insert metric history
	err = service.InsertMetricHistory(metric, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history")

	// Verify the metric history was created
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history after insertion")
	assert.Equal(t, int32(1), result.Size, "should have 1 metric history record")
	assert.Equal(t, 1, len(result.Items), "should have 1 item in the result")

	// Verify the inserted metric properties
	insertedMetric := &result.Items[0]
	assert.Equal(t, "test_metric", *insertedMetric.Name)
	assert.Equal(t, 42.5, *insertedMetric.Value)
	assert.Equal(t, "1234567890", *insertedMetric.Timestamp)
	assert.Equal(t, int64(1), *insertedMetric.Step)
	assert.Equal(t, "Test metric description", *insertedMetric.Description)

	// Test 2: Error handling - nil metric
	err = service.InsertMetricHistory(nil, *savedExperimentRun.Id)
	assert.Error(t, err, "should return error for nil metric")
	assert.Contains(t, err.Error(), "metric cannot be nil")

	// Test 3: Error handling - empty experiment run ID
	err = service.InsertMetricHistory(metric, "")
	assert.Error(t, err, "should return error for empty experiment run ID")
	assert.Contains(t, err.Error(), "experiment run ID is required")

	// Test 4: Error handling - non-existent experiment run ID
	err = service.InsertMetricHistory(metric, "999999")
	assert.Error(t, err, "should return error for non-existent experiment run ID")
	assert.Contains(t, err.Error(), "experiment run not found")
}

func TestUpsertExperimentRunArtifactTriggersMetricHistory(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-upsert",
		Description: apiutils.Of("Test experiment for metric history upsert"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-upsert"),
		Description: apiutils.Of("Test experiment run for metric history upsert"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create a metric artifact
	metricArtifact := &openapi.Artifact{
		Metric: &openapi.Metric{
			Name:      apiutils.Of("test_metric"),
			Value:     apiutils.Of(42.5),
			Timestamp: apiutils.Of("1234567890"),
			Step:      apiutils.Of(int64(1)),
		},
	}

	// Upsert the metric artifact (this should trigger InsertMetricHistory)
	createdArtifact, err := service.UpsertExperimentRunArtifact(metricArtifact, *savedExperimentRun.Id)
	require.NoError(t, err, "error upserting metric artifact")
	assert.NotNil(t, createdArtifact.Metric, "should have created metric artifact")

	// Verify that metric history was also created
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history after upsert")
	assert.Equal(t, int32(1), result.Size, "should have 1 metric history record created automatically")

	// Verify the metric history has the correct properties
	insertedMetric := &result.Items[0]
	assert.Equal(t, "test_metric", *insertedMetric.Name)
	assert.Equal(t, 42.5, *insertedMetric.Value)
	assert.Equal(t, "1234567890", *insertedMetric.Timestamp)
	assert.Equal(t, int64(1), *insertedMetric.Step)
}

func TestGetExperimentRunMetricHistoryEmptyResult(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-empty",
		Description: apiutils.Of("Test experiment for empty metric history"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-empty"),
		Description: apiutils.Of("Test experiment run for empty metric history"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Test getting metric history for experiment run with no metrics
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting empty metric history")
	assert.Equal(t, int32(0), result.Size, "should return 0 metric history records")
	assert.Equal(t, 0, len(result.Items), "should have 0 items in the result")
}

func TestGetExperimentRunMetricHistoryWithPagination(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-pagination",
		Description: apiutils.Of("Test experiment for metric history pagination"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-pagination"),
		Description: apiutils.Of("Test experiment run for metric history pagination"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create multiple metrics for pagination test
	for i := 0; i < 5; i++ {
		metric := &openapi.Metric{
			Name:      apiutils.Of(fmt.Sprintf("metric_%d", i)),
			Value:     apiutils.Of(float64(i) * 0.1),
			Timestamp: apiutils.Of(fmt.Sprintf("123456789%d", i)),
			Step:      apiutils.Of(int64(i)),
		}
		err = service.InsertMetricHistory(metric, *savedExperimentRun.Id)
		require.NoError(t, err, "error inserting metric history %d", i)
	}

	// Test pagination with page size 2
	pageSize := int32(2)
	listOptions := api.ListOptions{
		PageSize: &pageSize,
	}

	result, err := service.GetExperimentRunMetricHistory(nil, nil, listOptions, savedExperimentRun.Id)
	require.NoError(t, err, "error getting paginated metric history")
	assert.Equal(t, int32(2), result.Size, "should return 2 metric history records on first page")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")
	assert.Equal(t, pageSize, result.PageSize, "should have correct page size")
}

func TestGetExperimentRunMetricHistoryWithInvalidExperimentRunId(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Test with nil experiment run ID
	_, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, nil)
	assert.Error(t, err, "should return error for nil experiment run ID")
	assert.Contains(t, err.Error(), "experiment run ID is required")

	// Test with non-existent experiment run ID
	nonExistentId := "999999"
	_, err = service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, &nonExistentId)
	assert.Error(t, err, "should return error for non-existent experiment run ID")
	assert.Contains(t, err.Error(), "experiment run not found")
}

func TestInsertMetricHistoryWithLastUpdateTime(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-last-update",
		Description: apiutils.Of("Test experiment for metric history with last update time"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-last-update"),
		Description: apiutils.Of("Test experiment run for metric history with last update time"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Test insertion with last update time but no timestamp
	metricWithLastUpdate := &openapi.Metric{
		Name:                     apiutils.Of("test_metric_last_update"),
		Value:                    apiutils.Of(42.5),
		LastUpdateTimeSinceEpoch: apiutils.Of("9876543210"),
		Step:                     apiutils.Of(int64(1)),
	}

	// Insert metric history with last update time
	err = service.InsertMetricHistory(metricWithLastUpdate, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history with last update time")

	// Verify the metric history was created
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history after insertion")
	assert.Equal(t, int32(1), result.Size, "should have 1 metric history record")

	// Verify the inserted metric has the correct name (should use last update time as timestamp)
	insertedMetric := &result.Items[0]
	assert.Equal(t, "test_metric_last_update", *insertedMetric.Name)
	assert.Equal(t, 42.5, *insertedMetric.Value)
}

func TestGetExperimentRunMetricHistoryReturnsCorrectArtifactType(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-artifact-type",
		Description: apiutils.Of("Test experiment for metric history artifact type validation"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-artifact-type"),
		Description: apiutils.Of("Test experiment run for metric history artifact type validation"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create a metric for the experiment run
	metric := &openapi.Metric{
		Name:        apiutils.Of("test_metric"),
		Value:       apiutils.Of(0.95),
		Timestamp:   apiutils.Of("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Test metric for artifact type validation"),
	}

	// Insert metric history (this stores as 'metric-history' in database)
	err = service.InsertMetricHistory(metric, *savedExperimentRun.Id)
	require.NoError(t, err, "error inserting metric history")

	// Get metric history via the REST API endpoint
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history")
	assert.Equal(t, int32(1), result.Size, "should return 1 metric history record")
	assert.Equal(t, 1, len(result.Items), "should have 1 item in the result")

	// Verify the artifact type is 'metric' (not 'metric-history') for REST API compliance
	retrievedMetric := &result.Items[0]

	// The key validation: ensure artifact type is "metric" for REST API, even though it's stored as "metric-history" internally
	assert.Equal(t, "metric", *retrievedMetric.ArtifactType, "artifact type should be 'metric' for REST API compliance")

	// Verify other properties are correctly mapped
	assert.Equal(t, "test_metric", *retrievedMetric.Name)
	assert.Equal(t, 0.95, *retrievedMetric.Value)
	assert.Equal(t, "1234567890", *retrievedMetric.Timestamp)
	assert.Equal(t, int64(1), *retrievedMetric.Step)
	assert.Equal(t, "Test metric for artifact type validation", *retrievedMetric.Description)
}

func TestGetExperimentRunMetricHistoryWithStepIdsFilter(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-stepids",
		Description: apiutils.Of("Test experiment for metric history stepIds filter"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-stepids"),
		Description: apiutils.Of("Test experiment run for metric history stepIds filter"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics with different step values
	metrics := []*openapi.Metric{
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.85),
			Timestamp: apiutils.Of("1234567890"),
			Step:      apiutils.Of(int64(1)),
		},
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.90),
			Timestamp: apiutils.Of("1234567891"),
			Step:      apiutils.Of(int64(2)),
		},
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.95),
			Timestamp: apiutils.Of("1234567892"),
			Step:      apiutils.Of(int64(3)),
		},
		{
			Name:      apiutils.Of("loss"),
			Value:     apiutils.Of(0.15),
			Timestamp: apiutils.Of("1234567893"),
			Step:      apiutils.Of(int64(1)),
		},
		{
			Name:      apiutils.Of("loss"),
			Value:     apiutils.Of(0.10),
			Timestamp: apiutils.Of("1234567894"),
			Step:      apiutils.Of(int64(2)),
		},
	}

	// Insert all metric history records
	for _, metric := range metrics {
		err = service.InsertMetricHistory(metric, *savedExperimentRun.Id)
		require.NoError(t, err, "error inserting metric history")
	}

	// Test filtering by single step ID
	stepIds := "1"
	result, err := service.GetExperimentRunMetricHistory(nil, &stepIds, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history with step filter")
	assert.Equal(t, int32(2), result.Size, "should return 2 metric history records for step 1")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify the returned metrics are from step 1
	for _, item := range result.Items {
		assert.Equal(t, int64(1), *item.Step, "all metrics should be from step 1")
	}

	// Test filtering by multiple step IDs
	stepIds = "1,3"
	result, err = service.GetExperimentRunMetricHistory(nil, &stepIds, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history with multiple step filter")
	assert.Equal(t, int32(3), result.Size, "should return 3 metric history records for steps 1 and 3")
	assert.Equal(t, 3, len(result.Items), "should have 3 items in the result")

	// Verify the returned metrics are from steps 1 and 3
	stepValues := make(map[int64]bool)
	for _, item := range result.Items {
		stepValues[*item.Step] = true
	}
	assert.True(t, stepValues[1], "should contain metrics from step 1")
	assert.True(t, stepValues[3], "should contain metrics from step 3")
	assert.False(t, stepValues[2], "should not contain metrics from step 2")

	// Test filtering by non-existent step ID
	stepIds = "999"
	result, err = service.GetExperimentRunMetricHistory(nil, &stepIds, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history with non-existent step filter")
	assert.Equal(t, int32(0), result.Size, "should return 0 metric history records for non-existent step")
	assert.Equal(t, 0, len(result.Items), "should have 0 items in the result")

	// Test combining stepIds with name filter
	accuracyName := "accuracy"
	stepIds = "1,2"
	result, err = service.GetExperimentRunMetricHistory(&accuracyName, &stepIds, api.ListOptions{}, savedExperimentRun.Id)
	require.NoError(t, err, "error getting metric history with name and step filter")
	assert.Equal(t, int32(2), result.Size, "should return 2 accuracy metric history records for steps 1 and 2")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify all returned metrics are accuracy metrics from steps 1 and 2
	for _, item := range result.Items {
		assert.Equal(t, "accuracy", *item.Name, "all metrics should be accuracy")
		assert.True(t, *item.Step == 1 || *item.Step == 2, "all metrics should be from step 1 or 2")
	}
}

// Tests for GetExperimentRunsMetricHistory (multiple experiment runs)

func TestGetExperimentRunsMetricHistory(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-multiple",
		Description: apiutils.Of("Test experiment for multiple experiment runs metric history"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create multiple experiment runs
	experimentRun1 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-1"),
		Description: apiutils.Of("First test experiment run"),
	}
	savedExperimentRun1, err := service.UpsertExperimentRun(experimentRun1, savedExperiment.Id)
	require.NoError(t, err)

	experimentRun2 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-2"),
		Description: apiutils.Of("Second test experiment run"),
	}
	savedExperimentRun2, err := service.UpsertExperimentRun(experimentRun2, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics for first experiment run
	metric1 := &openapi.Metric{
		Name:        apiutils.Of("accuracy"),
		Value:       apiutils.Of(0.95),
		Timestamp:   apiutils.Of("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Accuracy metric for run 1"),
	}

	metric2 := &openapi.Metric{
		Name:        apiutils.Of("loss"),
		Value:       apiutils.Of(0.05),
		Timestamp:   apiutils.Of("1234567891"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Loss metric for run 1"),
	}

	// Create metrics for second experiment run
	metric3 := &openapi.Metric{
		Name:        apiutils.Of("accuracy"),
		Value:       apiutils.Of(0.92),
		Timestamp:   apiutils.Of("1234567892"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.Of("Accuracy metric for run 2"),
	}

	metric4 := &openapi.Metric{
		Name:        apiutils.Of("precision"),
		Value:       apiutils.Of(0.88),
		Timestamp:   apiutils.Of("1234567893"),
		Step:        apiutils.Of(int64(2)),
		Description: apiutils.Of("Precision metric for run 2"),
	}

	// Insert metric history records
	err = service.InsertMetricHistory(metric1, *savedExperimentRun1.Id)
	require.NoError(t, err, "error inserting metric history for run 1 metric 1")

	err = service.InsertMetricHistory(metric2, *savedExperimentRun1.Id)
	require.NoError(t, err, "error inserting metric history for run 1 metric 2")

	err = service.InsertMetricHistory(metric3, *savedExperimentRun2.Id)
	require.NoError(t, err, "error inserting metric history for run 2 metric 1")

	err = service.InsertMetricHistory(metric4, *savedExperimentRun2.Id)
	require.NoError(t, err, "error inserting metric history for run 2 metric 2")

	// Test getting metric history for multiple experiment runs
	experimentRunIds := fmt.Sprintf("%s,%s", *savedExperimentRun1.Id, *savedExperimentRun2.Id)
	result, err := service.GetExperimentRunsMetricHistory(nil, nil, &experimentRunIds, api.ListOptions{})
	require.NoError(t, err, "error getting metric history for multiple runs")
	assert.Equal(t, int32(4), result.Size, "should return 4 metric history records")
	assert.Equal(t, 4, len(result.Items), "should have 4 items in the result")

	// Verify we got metrics from both experiment runs
	foundMetrics := make(map[string]bool)
	for _, item := range result.Items {
		key := fmt.Sprintf("%s-%.2f", *item.Name, *item.Value)
		foundMetrics[key] = true
	}

	assert.True(t, foundMetrics["accuracy-0.95"], "should find accuracy metric from run 1")
	assert.True(t, foundMetrics["loss-0.05"], "should find loss metric from run 1")
	assert.True(t, foundMetrics["accuracy-0.92"], "should find accuracy metric from run 2")
	assert.True(t, foundMetrics["precision-0.88"], "should find precision metric from run 2")
}

func TestGetExperimentRunsMetricHistoryWithNameFilter(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-multiple-name-filter",
		Description: apiutils.Of("Test experiment for multiple experiment runs metric history with name filter"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create multiple experiment runs
	experimentRun1 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-name-filter-1"),
		Description: apiutils.Of("First test experiment run for name filter"),
	}
	savedExperimentRun1, err := service.UpsertExperimentRun(experimentRun1, savedExperiment.Id)
	require.NoError(t, err)

	experimentRun2 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-name-filter-2"),
		Description: apiutils.Of("Second test experiment run for name filter"),
	}
	savedExperimentRun2, err := service.UpsertExperimentRun(experimentRun2, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics with same names across runs
	metrics := []*openapi.Metric{
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.95),
			Timestamp: apiutils.Of("1234567890"),
			Step:      apiutils.Of(int64(1)),
		},
		{
			Name:      apiutils.Of("loss"),
			Value:     apiutils.Of(0.05),
			Timestamp: apiutils.Of("1234567891"),
			Step:      apiutils.Of(int64(1)),
		},
	}

	// Insert metrics for both runs
	for _, metric := range metrics {
		err = service.InsertMetricHistory(metric, *savedExperimentRun1.Id)
		require.NoError(t, err, "error inserting metric history for run 1")

		err = service.InsertMetricHistory(metric, *savedExperimentRun2.Id)
		require.NoError(t, err, "error inserting metric history for run 2")
	}

	// Test filtering by name across multiple runs
	accuracyName := "accuracy"
	experimentRunIds := fmt.Sprintf("%s,%s", *savedExperimentRun1.Id, *savedExperimentRun2.Id)
	result, err := service.GetExperimentRunsMetricHistory(&accuracyName, nil, &experimentRunIds, api.ListOptions{})
	require.NoError(t, err, "error getting metric history with name filter")
	assert.Equal(t, int32(2), result.Size, "should return 2 accuracy metric history records")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify all returned metrics are accuracy metrics
	for _, item := range result.Items {
		assert.Equal(t, "accuracy", *item.Name, "all metrics should be accuracy")
		assert.Equal(t, 0.95, *item.Value, "all accuracy metrics should have value 0.95")
	}
}

func TestGetExperimentRunsMetricHistoryWithStepIdsFilter(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-multiple-step-filter",
		Description: apiutils.Of("Test experiment for multiple experiment runs metric history with step filter"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create multiple experiment runs
	experimentRun1 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-step-filter-1"),
		Description: apiutils.Of("First test experiment run for step filter"),
	}
	savedExperimentRun1, err := service.UpsertExperimentRun(experimentRun1, savedExperiment.Id)
	require.NoError(t, err)

	experimentRun2 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-step-filter-2"),
		Description: apiutils.Of("Second test experiment run for step filter"),
	}
	savedExperimentRun2, err := service.UpsertExperimentRun(experimentRun2, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics with different step values for both runs
	metricsRun1 := []*openapi.Metric{
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.85),
			Timestamp: apiutils.Of("1234567890"),
			Step:      apiutils.Of(int64(1)),
		},
		{
			Name:      apiutils.Of("accuracy"),
			Value:     apiutils.Of(0.90),
			Timestamp: apiutils.Of("1234567891"),
			Step:      apiutils.Of(int64(2)),
		},
		{
			Name:      apiutils.Of("loss"),
			Value:     apiutils.Of(0.15),
			Timestamp: apiutils.Of("1234567892"),
			Step:      apiutils.Of(int64(3)),
		},
	}

	metricsRun2 := []*openapi.Metric{
		{
			Name:      apiutils.Of("precision"),
			Value:     apiutils.Of(0.88),
			Timestamp: apiutils.Of("1234567893"),
			Step:      apiutils.Of(int64(1)),
		},
		{
			Name:      apiutils.Of("recall"),
			Value:     apiutils.Of(0.82),
			Timestamp: apiutils.Of("1234567894"),
			Step:      apiutils.Of(int64(2)),
		},
	}

	// Insert metrics for both runs
	for _, metric := range metricsRun1 {
		err = service.InsertMetricHistory(metric, *savedExperimentRun1.Id)
		require.NoError(t, err, "error inserting metric history for run 1")
	}

	for _, metric := range metricsRun2 {
		err = service.InsertMetricHistory(metric, *savedExperimentRun2.Id)
		require.NoError(t, err, "error inserting metric history for run 2")
	}

	// Test filtering by step ID across multiple runs
	stepIds := "1"
	experimentRunIds := fmt.Sprintf("%s,%s", *savedExperimentRun1.Id, *savedExperimentRun2.Id)
	result, err := service.GetExperimentRunsMetricHistory(nil, &stepIds, &experimentRunIds, api.ListOptions{})
	require.NoError(t, err, "error getting metric history with step filter")
	assert.Equal(t, int32(2), result.Size, "should return 2 metric history records for step 1")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify all returned metrics are from step 1
	for _, item := range result.Items {
		assert.Equal(t, int64(1), *item.Step, "all metrics should be from step 1")
	}

	// Test filtering by multiple step IDs
	stepIds = "1,2"
	result, err = service.GetExperimentRunsMetricHistory(nil, &stepIds, &experimentRunIds, api.ListOptions{})
	require.NoError(t, err, "error getting metric history with multiple step filter")
	assert.Equal(t, int32(4), result.Size, "should return 4 metric history records for steps 1 and 2")
	assert.Equal(t, 4, len(result.Items), "should have 4 items in the result")

	// Verify returned metrics are from steps 1 and 2 only
	for _, item := range result.Items {
		assert.True(t, *item.Step == 1 || *item.Step == 2, "all metrics should be from step 1 or 2")
	}
}

func TestGetExperimentRunsMetricHistoryEmptyExperimentRunIds(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Test with nil experiment run IDs
	result, err := service.GetExperimentRunsMetricHistory(nil, nil, nil, api.ListOptions{})
	require.NoError(t, err, "should not error with nil experiment run IDs")
	assert.Equal(t, int32(0), result.Size, "should return empty list")
	assert.Equal(t, 0, len(result.Items), "should have 0 items")

	// Test with empty experiment run IDs string
	emptyIds := ""
	result, err = service.GetExperimentRunsMetricHistory(nil, nil, &emptyIds, api.ListOptions{})
	require.NoError(t, err, "should not error with empty experiment run IDs")
	assert.Equal(t, int32(0), result.Size, "should return empty list")
	assert.Equal(t, 0, len(result.Items), "should have 0 items")

	// Test with whitespace-only experiment run IDs
	whitespaceIds := "   "
	result, err = service.GetExperimentRunsMetricHistory(nil, nil, &whitespaceIds, api.ListOptions{})
	require.NoError(t, err, "should not error with whitespace experiment run IDs")
	assert.Equal(t, int32(0), result.Size, "should return empty list")
	assert.Equal(t, 0, len(result.Items), "should have 0 items")
}

func TestGetExperimentRunsMetricHistoryNonExistentExperimentRuns(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Test with non-existent experiment run IDs
	nonExistentIds := "999999,888888"
	_, err := service.GetExperimentRunsMetricHistory(nil, nil, &nonExistentIds, api.ListOptions{})
	assert.Error(t, err, "should return error for non-existent experiment run IDs")
	assert.Contains(t, err.Error(), "experiment run with ID 999999 not found", "should specify which experiment run was not found")
}

func TestGetExperimentRunsMetricHistoryMixedExistentAndNonExistentRuns(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment and experiment run
	experiment := &openapi.Experiment{
		Name:        "test-experiment-mixed",
		Description: apiutils.Of("Test experiment for mixed existent/non-existent runs"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-mixed"),
		Description: apiutils.Of("Test experiment run for mixed test"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Test with mix of existent and non-existent experiment run IDs
	mixedIds := fmt.Sprintf("%s,999999", *savedExperimentRun.Id)
	_, err = service.GetExperimentRunsMetricHistory(nil, nil, &mixedIds, api.ListOptions{})
	assert.Error(t, err, "should return error when any experiment run ID is not found")
	assert.Contains(t, err.Error(), "experiment run with ID 999999 not found", "should specify which experiment run was not found")
}

func TestGetExperimentRunsMetricHistoryWithPagination(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-multiple-pagination",
		Description: apiutils.Of("Test experiment for multiple experiment runs metric history pagination"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create experiment run
	experimentRun := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-pagination"),
		Description: apiutils.Of("Test experiment run for pagination"),
	}
	savedExperimentRun, err := service.UpsertExperimentRun(experimentRun, savedExperiment.Id)
	require.NoError(t, err)

	// Create multiple metrics for pagination test
	for i := 0; i < 5; i++ {
		metric := &openapi.Metric{
			Name:      apiutils.Of(fmt.Sprintf("metric_%d", i)),
			Value:     apiutils.Of(float64(i) * 0.1),
			Timestamp: apiutils.Of(fmt.Sprintf("123456789%d", i)),
			Step:      apiutils.Of(int64(i)),
		}
		err = service.InsertMetricHistory(metric, *savedExperimentRun.Id)
		require.NoError(t, err, "error inserting metric history %d", i)
	}

	// Test pagination with page size 2
	pageSize := int32(2)
	listOptions := api.ListOptions{
		PageSize: &pageSize,
	}

	result, err := service.GetExperimentRunsMetricHistory(nil, nil, savedExperimentRun.Id, listOptions)
	require.NoError(t, err, "error getting paginated metric history for multiple runs")
	assert.Equal(t, int32(2), result.Size, "should return 2 metric history records on first page")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")
	assert.Equal(t, pageSize, result.PageSize, "should have correct page size")
}

func TestGetExperimentRunsMetricHistoryWithCombinedFilters(t *testing.T) {
	service, cleanup := SetupModelRegistryService(t)
	defer cleanup()

	// Create experiment
	experiment := &openapi.Experiment{
		Name:        "test-experiment-combined-filters",
		Description: apiutils.Of("Test experiment for combined filters on multiple runs"),
	}
	savedExperiment, err := service.UpsertExperiment(experiment)
	require.NoError(t, err)

	// Create multiple experiment runs
	experimentRun1 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-combined-1"),
		Description: apiutils.Of("First test experiment run for combined filters"),
	}
	savedExperimentRun1, err := service.UpsertExperimentRun(experimentRun1, savedExperiment.Id)
	require.NoError(t, err)

	experimentRun2 := &openapi.ExperimentRun{
		Name:        apiutils.Of("test-experiment-run-combined-2"),
		Description: apiutils.Of("Second test experiment run for combined filters"),
	}
	savedExperimentRun2, err := service.UpsertExperimentRun(experimentRun2, savedExperiment.Id)
	require.NoError(t, err)

	// Create metrics with varying names and steps for both runs
	metricsData := []struct {
		runId *string
		name  string
		value float64
		step  int64
	}{
		{savedExperimentRun1.Id, "accuracy", 0.85, 1},
		{savedExperimentRun1.Id, "accuracy", 0.90, 2},
		{savedExperimentRun1.Id, "loss", 0.15, 1},
		{savedExperimentRun1.Id, "loss", 0.10, 2},
		{savedExperimentRun2.Id, "accuracy", 0.88, 1},
		{savedExperimentRun2.Id, "accuracy", 0.92, 2},
		{savedExperimentRun2.Id, "precision", 0.75, 1},
		{savedExperimentRun2.Id, "precision", 0.80, 2},
	}

	for i, data := range metricsData {
		metric := &openapi.Metric{
			Name:      apiutils.Of(data.name),
			Value:     apiutils.Of(data.value),
			Timestamp: apiutils.Of(fmt.Sprintf("123456789%d", i)),
			Step:      apiutils.Of(data.step),
		}
		err = service.InsertMetricHistory(metric, *data.runId)
		require.NoError(t, err, "error inserting metric history")
	}

	// Test combined name and step filters
	accuracyName := "accuracy"
	stepIds := "1"
	experimentRunIds := fmt.Sprintf("%s,%s", *savedExperimentRun1.Id, *savedExperimentRun2.Id)

	result, err := service.GetExperimentRunsMetricHistory(&accuracyName, &stepIds, &experimentRunIds, api.ListOptions{})
	require.NoError(t, err, "error getting metric history with combined filters")
	assert.Equal(t, int32(2), result.Size, "should return 2 accuracy metrics from step 1")
	assert.Equal(t, 2, len(result.Items), "should have 2 items in the result")

	// Verify all returned metrics match the filters
	for _, item := range result.Items {
		assert.Equal(t, "accuracy", *item.Name, "all metrics should be accuracy")
		assert.Equal(t, int64(1), *item.Step, "all metrics should be from step 1")
	}

	// Verify we got metrics from both runs
	values := make(map[float64]bool)
	for _, item := range result.Items {
		values[*item.Value] = true
	}
	assert.True(t, values[0.85], "should have accuracy metric from run 1 (value 0.85)")
	assert.True(t, values[0.88], "should have accuracy metric from run 2 (value 0.88)")
}
