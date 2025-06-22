package core

import (
	"context"
	"fmt"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/converter"
	"github.com/kubeflow/model-registry/internal/defaults"
	"github.com/kubeflow/model-registry/internal/ml_metadata/proto"
	"github.com/kubeflow/model-registry/pkg/api"
	"github.com/kubeflow/model-registry/pkg/openapi"
)

// METRIC HISTORY TESTS

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistory() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create metrics for the experiment run
	metric1 := &openapi.Metric{
		Name:        apiutils.StrPtr("accuracy"),
		Value:       apiutils.Of(0.95),
		Timestamp:   apiutils.StrPtr("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.StrPtr("Test accuracy metric"),
	}

	metric2 := &openapi.Metric{
		Name:        apiutils.StrPtr("loss"),
		Value:       apiutils.Of(0.05),
		Timestamp:   apiutils.StrPtr("1234567891"),
		Step:        apiutils.Of(int64(2)),
		Description: apiutils.StrPtr("Test loss metric"),
	}

	// Insert metric history records
	err := service.InsertMetricHistory(metric1, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric1: %v", err)

	err = service.InsertMetricHistory(metric2, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric2: %v", err)

	// Test getting all metric history
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting metric history: %v", err)
	suite.Equal(int32(2), result.Size, "should return 2 metric history records")
	suite.Equal(2, len(result.Items), "should have 2 items in the result")

	// Verify the metrics are returned correctly
	foundAccuracy := false
	foundLoss := false
	for _, item := range result.Items {
		if item.Metric != nil {
			if *item.Metric.Name == "accuracy" {
				foundAccuracy = true
				suite.Equal(0.95, *item.Metric.Value)
				suite.Equal("1234567890", *item.Metric.Timestamp)
				suite.Equal(int64(1), *item.Metric.Step)
			} else if *item.Metric.Name == "loss" {
				foundLoss = true
				suite.Equal(0.05, *item.Metric.Value)
				suite.Equal("1234567891", *item.Metric.Timestamp)
				suite.Equal(int64(2), *item.Metric.Step)
			}
		}
	}
	suite.True(foundAccuracy, "should find accuracy metric")
	suite.True(foundLoss, "should find loss metric")
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryWithNameFilter() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create metrics for the experiment run
	metric1 := &openapi.Metric{
		Name:      apiutils.StrPtr("accuracy"),
		Value:     apiutils.Of(0.95),
		Timestamp: apiutils.StrPtr("1234567890"),
		Step:      apiutils.Of(int64(1)),
	}

	metric2 := &openapi.Metric{
		Name:      apiutils.StrPtr("loss"),
		Value:     apiutils.Of(0.05),
		Timestamp: apiutils.StrPtr("1234567891"),
		Step:      apiutils.Of(int64(2)),
	}

	// Insert metric history records
	err := service.InsertMetricHistory(metric1, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric1: %v", err)

	err = service.InsertMetricHistory(metric2, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric2: %v", err)

	// Test filtering by name
	accuracyName := "accuracy"
	result, err := service.GetExperimentRunMetricHistory(&accuracyName, nil, api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting metric history with name filter: %v", err)
	suite.Equal(int32(1), result.Size, "should return 1 metric history record for accuracy")
	suite.Equal(1, len(result.Items), "should have 1 item in the result")
	suite.Equal("accuracy", *result.Items[0].Metric.Name)
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryWithStepFilter() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create metrics for the experiment run
	metric1 := &openapi.Metric{
		Name:      apiutils.StrPtr("accuracy"),
		Value:     apiutils.Of(0.95),
		Timestamp: apiutils.StrPtr("1234567890"),
		Step:      apiutils.Of(int64(1)),
	}

	metric2 := &openapi.Metric{
		Name:      apiutils.StrPtr("accuracy"),
		Value:     apiutils.Of(0.97),
		Timestamp: apiutils.StrPtr("1234567892"),
		Step:      apiutils.Of(int64(3)),
	}

	// Insert metric history records
	err := service.InsertMetricHistory(metric1, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric1: %v", err)

	err = service.InsertMetricHistory(metric2, experimentRunId)
	suite.Nilf(err, "error inserting metric history for metric2: %v", err)

	// Test filtering by step IDs
	stepIds := "1"
	result, err := service.GetExperimentRunMetricHistory(nil, &stepIds, api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting metric history with step filter: %v", err)
	suite.Equal(int32(1), result.Size, "should return 1 metric history record for step 1")
	suite.Equal(1, len(result.Items), "should have 1 item in the result")
	suite.Equal(int64(1), *result.Items[0].Metric.Step)
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryWithPagination() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create multiple metrics for the experiment run
	for i := 1; i <= 5; i++ {
		metric := &openapi.Metric{
			Name:      apiutils.StrPtr(fmt.Sprintf("metric_%d", i)),
			Value:     apiutils.Of(float64(i)),
			Timestamp: apiutils.StrPtr(fmt.Sprintf("123456789%d", i)),
			Step:      apiutils.Of(int64(i)),
		}

		err := service.InsertMetricHistory(metric, experimentRunId)
		suite.Nilf(err, "error inserting metric history for metric %d: %v", i, err)
	}

	// Test pagination with page size 2
	pageSize := int32(2)
	listOptions := api.ListOptions{
		PageSize: &pageSize,
	}

	result, err := service.GetExperimentRunMetricHistory(nil, nil, listOptions, &experimentRunId)
	suite.Nilf(err, "error getting metric history with pagination: %v", err)
	suite.Equal(int32(2), result.Size, "should return 2 metric history records")
	suite.Equal(2, len(result.Items), "should have 2 items in the result")
	suite.NotNil(result.NextPageToken, "should have next page token")
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryWithInvalidExperimentRunId() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Test with nil experiment run ID
	_, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, nil)
	suite.NotNil(err, "should return error for nil experiment run ID")
	suite.Contains(err.Error(), "experiment run ID is required")

	// Test with non-existent experiment run ID
	nonExistentId := "999999"
	_, err = service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, &nonExistentId)
	suite.NotNil(err, "should return error for non-existent experiment run ID")
	suite.Contains(err.Error(), "experiment run not found")
}

func (suite *CoreTestSuite) TestInsertMetricHistory() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Test 1: Basic metric history insertion
	metric := &openapi.Metric{
		Name:        apiutils.StrPtr("test_metric"),
		Value:       apiutils.Of(42.5),
		Timestamp:   apiutils.StrPtr("1234567890"),
		Step:        apiutils.Of(int64(1)),
		Description: apiutils.StrPtr("Test metric description"),
		CustomProperties: &map[string]openapi.MetadataValue{
			"custom_prop": {
				MetadataStringValue: converter.NewMetadataStringValue("custom_value"),
			},
		},
	}

	// Insert metric history
	err := service.InsertMetricHistory(metric, experimentRunId)
	suite.Nilf(err, "error inserting metric history: %v", err)

	// Verify the metric history was created in MLMD
	metricHistoryTypeName := defaults.MetricHistoryTypeName
	artifacts, err := suite.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &metricHistoryTypeName,
	})
	suite.Nilf(err, "error getting metric history artifacts: %v", err)
	suite.Equal(1, len(artifacts.Artifacts), "should have 1 metric history artifact")

	// Verify the artifact properties
	artifact := artifacts.Artifacts[0]
	suite.Equal(defaults.MetricHistoryTypeName, *artifact.Type, "artifact should be of MetricHistory type")
	suite.Contains(*artifact.Name, "test_metric", "artifact name should contain metric name")
	suite.Contains(*artifact.Name, experimentRunId, "artifact name should contain experiment run ID")
	suite.Equal(42.5, artifact.Properties["value"].GetDoubleValue(), "metric value should match")
	suite.Equal("1234567890", artifact.Properties["timestamp"].GetStringValue(), "metric timestamp should match")
	suite.Equal(int64(1), artifact.Properties["step"].GetIntValue(), "metric step should match")
	suite.Equal("custom_value", artifact.CustomProperties["custom_prop"].GetStringValue(), "custom property should match")

	// Test 2: Insertion with last update time but no timestamp
	metricWithLastUpdate := &openapi.Metric{
		Name:                     apiutils.StrPtr("test_metric_last_update"),
		Value:                    apiutils.Of(42.5),
		LastUpdateTimeSinceEpoch: apiutils.StrPtr("9876543210"),
		Step:                     apiutils.Of(int64(1)),
	}

	// Insert metric history with last update time
	err = service.InsertMetricHistory(metricWithLastUpdate, experimentRunId)
	suite.Nilf(err, "error inserting metric history with last update time: %v", err)

	// Verify the metric history was created with last update time
	artifacts, err = suite.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &metricHistoryTypeName,
	})
	suite.Nilf(err, "error getting metric history artifacts: %v", err)
	suite.Equal(2, len(artifacts.Artifacts), "should have 2 metric history artifacts")

	// Verify the artifact name contains the last update time
	foundLastUpdateArtifact := false
	for _, a := range artifacts.Artifacts {
		if *a.Name == fmt.Sprintf("%s:test_metric_last_update__9876543210", experimentRunId) {
			foundLastUpdateArtifact = true
			break
		}
	}
	suite.True(foundLastUpdateArtifact, "should find artifact with last update time in name")

	// Test 3: Error handling - nil metric
	err = service.InsertMetricHistory(nil, experimentRunId)
	suite.NotNil(err, "should return error for nil metric")
	suite.Contains(err.Error(), "metric cannot be nil")

	// Test 4: Error handling - empty experiment run ID
	err = service.InsertMetricHistory(metric, "")
	suite.NotNil(err, "should return error for empty experiment run ID")
	suite.Contains(err.Error(), "experiment run ID is required")

	// Test 5: Error handling - non-existent experiment run ID
	err = service.InsertMetricHistory(metric, "999999")
	suite.NotNil(err, "should return error for non-existent experiment run ID")
	suite.Contains(err.Error(), "experiment run not found")
}

func (suite *CoreTestSuite) TestUpsertExperimentRunArtifactTriggersMetricHistory() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create a metric artifact
	metricArtifact := &openapi.Artifact{
		Metric: &openapi.Metric{
			Name:      apiutils.StrPtr("test_metric"),
			Value:     apiutils.Of(42.5),
			Timestamp: apiutils.StrPtr("1234567890"),
			Step:      apiutils.Of(int64(1)),
		},
	}

	// Upsert the metric artifact (this should trigger InsertMetricHistory)
	createdArtifact, err := service.UpsertExperimentRunArtifact(metricArtifact, experimentRunId)
	suite.Nilf(err, "error upserting metric artifact: %v", err)
	suite.NotNil(createdArtifact.Metric, "should have created metric artifact")

	// Verify that metric history was also created
	metricHistoryTypeName := defaults.MetricHistoryTypeName
	artifacts, err := suite.mlmdClient.GetArtifactsByType(context.Background(), &proto.GetArtifactsByTypeRequest{
		TypeName: &metricHistoryTypeName,
	})
	suite.Nilf(err, "error getting metric history artifacts: %v", err)
	suite.Equal(1, len(artifacts.Artifacts), "should have 1 metric history artifact created automatically")

	// Verify the metric history has the correct name format
	artifact := artifacts.Artifacts[0]
	suite.Contains(*artifact.Name, experimentRunId, "metric history name should contain experiment run ID")
	suite.Contains(*artifact.Name, "test_metric", "metric history name should contain metric name")
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryEmptyResult() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Test getting metric history for experiment run with no metrics
	result, err := service.GetExperimentRunMetricHistory(nil, nil, api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting metric history: %v", err)
	suite.Equal(int32(0), result.Size, "should return 0 metric history records")
	suite.Equal(0, len(result.Items), "should have 0 items in the result")
}

func (suite *CoreTestSuite) TestGetExperimentRunMetricHistoryWithMultipleSteps() {
	// create model registry service
	service := suite.setupModelRegistryService()

	// Create experiment run (which also creates the experiment)
	experimentRunId := suite.registerExperimentRun(service, nil, nil, nil, nil)

	// Create metrics for different steps
	for step := 1; step <= 3; step++ {
		metric := &openapi.Metric{
			Name:      apiutils.StrPtr("accuracy"),
			Value:     apiutils.Of(float64(0.9 + float64(step)*0.01)),
			Timestamp: apiutils.StrPtr(fmt.Sprintf("123456789%d", step)),
			Step:      apiutils.Of(int64(step)),
		}

		err := service.InsertMetricHistory(metric, experimentRunId)
		suite.Nilf(err, "error inserting metric history for step %d: %v", step, err)
	}

	// Test filtering by multiple step IDs
	stepIds := "1,3"
	result, err := service.GetExperimentRunMetricHistory(nil, &stepIds, api.ListOptions{}, &experimentRunId)
	suite.Nilf(err, "error getting metric history with multiple step filter: %v", err)
	suite.Equal(int32(2), result.Size, "should return 2 metric history records for steps 1 and 3")
	suite.Equal(2, len(result.Items), "should have 2 items in the result")

	// Verify we got the correct steps
	steps := make(map[int64]bool)
	for _, item := range result.Items {
		if item.Metric != nil {
			steps[*item.Metric.Step] = true
		}
	}
	suite.True(steps[1], "should have step 1")
	suite.True(steps[3], "should have step 3")
	suite.False(steps[2], "should not have step 2")
}
