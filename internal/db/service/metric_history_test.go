package service_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/kubeflow/model-registry/internal/apiutils"
	"github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMetricHistoryRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual MetricHistory type ID from the database
	typeID := getMetricHistoryTypeID(t, db)
	repo := service.NewMetricHistoryRepository(db, typeID)

	// Also get experiment and experiment run types for creating related entities
	experimentTypeID := getExperimentTypeID(t, db)
	experimentRepo := service.NewExperimentRepository(db, experimentTypeID)

	experimentRunTypeID := getExperimentRunTypeID(t, db)
	experimentRunRepo := service.NewExperimentRunRepository(db, experimentRunTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new metric history
		metricHistory := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("test-metric-history"),
				ExternalID:   apiutils.Of("metric-history-ext-123"),
				URI:          apiutils.Of("s3://bucket/metric-history.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test metric history description"),
				},
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(100)),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-metric-history-prop",
					StringValue:      apiutils.Of("custom-metric-history-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-metric-history", *saved.GetAttributes().Name)
		assert.Equal(t, "metric-history-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/metric-history.json", *saved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *saved.GetAttributes().State)
		assert.Equal(t, "metric-history", *saved.GetAttributes().ArtifactType)
		assert.NotNil(t, saved.GetAttributes().CreateTimeSinceEpoch)
		assert.NotNil(t, saved.GetAttributes().LastUpdateTimeSinceEpoch)

		// Test updating the same metric history
		metricHistory.ID = saved.GetID()
		metricHistory.GetAttributes().Name = apiutils.Of("updated-metric-history")
		metricHistory.GetAttributes().State = apiutils.Of("PENDING")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		metricHistory.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-metric-history", *updated.GetAttributes().Name)
		assert.Equal(t, "PENDING", *updated.GetAttributes().State)
		assert.Equal(t, *saved.GetAttributes().CreateTimeSinceEpoch, *updated.GetAttributes().CreateTimeSinceEpoch)
		assert.Greater(t, *updated.GetAttributes().LastUpdateTimeSinceEpoch, *saved.GetAttributes().LastUpdateTimeSinceEpoch)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a metric history to retrieve
		metricHistory := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("get-test-metric-history"),
				ExternalID:   apiutils.Of("get-metric-history-ext-123"),
				URI:          apiutils.Of("s3://bucket/get-metric-history.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric history for get test"),
				},
			},
		}

		saved, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-metric-history", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-metric-history-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/get-metric-history.json", *retrieved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *retrieved.GetAttributes().State)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "metric history by id not found")
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple metric histories for listing
		testMetricHistories := []*models.MetricHistoryImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("list-metric-history-1"),
					ExternalID:   apiutils.Of("list-metric-history-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-metric-history-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("First metric history"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("list-metric-history-2"),
					ExternalID:   apiutils.Of("list-metric-history-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-metric-history-2.json"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Second metric history"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("list-metric-history-3"),
					ExternalID:   apiutils.Of("list-metric-history-ext-3"),
					URI:          apiutils.Of("s3://bucket/list-metric-history-3.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Third metric history"),
					},
				},
			},
		}

		for _, metricHistory := range testMetricHistories {
			_, err := repo.Save(metricHistory, nil)
			require.NoError(t, err)
		}

		// Test listing all metric histories with basic pagination
		pageSize := int32(10)
		listOptions := models.MetricHistoryListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test metric histories

		// Test listing by name
		listOptions = models.MetricHistoryListOptions{
			Name: apiutils.Of("list-metric-history-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-metric-history-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.MetricHistoryListOptions{
			ExternalID: apiutils.Of("list-metric-history-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-metric-history-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.MetricHistoryListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("ID"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		// Verify we get results back and they are ordered by ID
		assert.GreaterOrEqual(t, len(result.Items), 1)
		if len(result.Items) > 1 {
			// Verify ascending ID order
			firstID := *result.Items[0].GetID()
			secondID := *result.Items[1].GetID()
			assert.Less(t, firstID, secondID, "Results should be ordered by ID ascending")
		}
	})

	t.Run("TestListOrdering", func(t *testing.T) {
		// Create metric histories sequentially with time delays to ensure deterministic ordering
		metricHistory1 := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("time-test-metric-history-1"),
				URI:          apiutils.Of("s3://bucket/time-metric-history-1.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
		}
		saved1, err := repo.Save(metricHistory1, nil)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		metricHistory2 := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("time-test-metric-history-2"),
				URI:          apiutils.Of("s3://bucket/time-metric-history-2.json"),
				State:        apiutils.Of("PENDING"),
				ArtifactType: apiutils.Of("metric-history"),
			},
		}
		saved2, err := repo.Save(metricHistory2, nil)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.MetricHistoryListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test metric histories in the results
		var foundMetricHistory1, foundMetricHistory2 models.MetricHistory
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundMetricHistory1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundMetricHistory2 = item
				index2 = i
			}
		}

		// Verify both metric histories were found and metricHistory1 comes before metricHistory2 (ascending order)
		require.NotEqual(t, -1, index1, "MetricHistory 1 should be found in results")
		require.NotEqual(t, -1, index2, "MetricHistory 2 should be found in results")
		assert.Less(t, index1, index2, "MetricHistory 1 should come before MetricHistory 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundMetricHistory1.GetAttributes().CreateTimeSinceEpoch, *foundMetricHistory2.GetAttributes().CreateTimeSinceEpoch, "MetricHistory 1 should have earlier create time")
	})

	t.Run("TestSaveWithExperimentRun", func(t *testing.T) {
		// First create an experiment and experiment run
		experiment := &models.ExperimentImpl{
			TypeID: apiutils.Of(int32(experimentTypeID)),
			Attributes: &models.ExperimentAttributes{
				Name: apiutils.Of("test-experiment-for-metric-history"),
			},
		}
		savedExperiment, err := experimentRepo.Save(experiment)
		require.NoError(t, err)

		experimentRun := &models.ExperimentRunImpl{
			TypeID: apiutils.Of(int32(experimentRunTypeID)),
			Attributes: &models.ExperimentRunAttributes{
				Name: apiutils.Of("test-experiment-run-for-metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "experiment_id",
					StringValue: apiutils.Of(fmt.Sprintf("%d", *savedExperiment.GetID())),
				},
			},
		}
		savedExperimentRun, err := experimentRunRepo.Save(experimentRun, savedExperiment.GetID())
		require.NoError(t, err)

		// Test creating a metric history with experiment run attribution
		metricHistory := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of(fmt.Sprintf("%d:experiment-run-metric-history", *savedExperimentRun.GetID())),
				URI:          apiutils.Of("s3://bucket/experiment-run-metric-history.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric history associated with experiment run"),
				},
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(50)),
				},
			},
		}

		saved, err := repo.Save(metricHistory, savedExperimentRun.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:experiment-run-metric-history", *savedExperimentRun.GetID()), *saved.GetAttributes().Name)
		assert.Equal(t, "s3://bucket/experiment-run-metric-history.json", *saved.GetAttributes().URI)

		// Test listing by experiment run ID
		listOptions := models.MetricHistoryListOptions{
			ExperimentRunID: savedExperimentRun.GetID(),
		}
		listOptions.PageSize = apiutils.Of(int32(10))

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should find our metric history

		// Verify the found metric history
		found := false
		for _, item := range result.Items {
			if *item.GetID() == *saved.GetID() {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the metric history associated with the experiment run")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		metricHistory := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("props-test-metric-history"),
				ExternalID:   apiutils.Of("props-metric-history-ext-123"),
				URI:          apiutils.Of("s3://bucket/props-metric-history.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric history with properties"),
				},
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(100)),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.95),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "team",
					StringValue:      apiutils.Of("ml-team"),
					IsCustomProperty: true,
				},
				{
					Name:             "priority",
					IntValue:         apiutils.Of(int32(5)),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, step, value

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2) // team, priority

		// Verify specific properties exist
		foundDescription := false
		foundStep := false
		foundValue := false
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Metric history with properties", *prop.StringValue)
			case "step":
				foundStep = true
				assert.Equal(t, int32(100), *prop.IntValue)
			case "value":
				foundValue = true
				assert.Equal(t, 0.95, *prop.DoubleValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundStep, "step property should exist")
		assert.True(t, foundValue, "value property should exist")

		// Verify custom properties
		foundTeam := false
		foundPriority := false
		for _, prop := range *retrieved.GetCustomProperties() {
			switch prop.Name {
			case "team":
				foundTeam = true
				assert.Equal(t, "ml-team", *prop.StringValue)
			case "priority":
				foundPriority = true
				assert.Equal(t, int32(5), *prop.IntValue)
			}
		}
		assert.True(t, foundTeam, "team custom property should exist")
		assert.True(t, foundPriority, "priority custom property should exist")
	})

	t.Run("TestPagination", func(t *testing.T) {
		// Create multiple metric histories for pagination testing
		for i := 0; i < 5; i++ {
			metricHistory := &models.MetricHistoryImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of(fmt.Sprintf("page-metric-history-%d", i)),
					URI:          apiutils.Of(fmt.Sprintf("s3://bucket/page-metric-history-%d.json", i)),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
			}
			_, err := repo.Save(metricHistory, nil)
			require.NoError(t, err)
		}

		// Test pagination with page size 2
		pageSize := int32(2)
		listOptions := models.MetricHistoryListOptions{
			Pagination: models.Pagination{
				PageSize: &pageSize,
			},
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should return at most 2 items
		assert.LessOrEqual(t, len(result.Items), 2)
		assert.Equal(t, pageSize, result.PageSize)
		assert.Equal(t, int32(len(result.Items)), result.Size)

		// If we have more items, there should be a next page token
		if len(result.Items) == 2 {
			assert.NotEmpty(t, result.NextPageToken)

			// Test getting the next page
			listOptions.NextPageToken = apiutils.Of(result.NextPageToken)
			nextResult, err := repo.List(listOptions)
			require.NoError(t, err)
			require.NotNil(t, nextResult)

			// Should get different items
			if len(nextResult.Items) > 0 {
				firstPageFirstID := *result.Items[0].GetID()
				nextPageFirstID := *nextResult.Items[0].GetID()
				assert.NotEqual(t, firstPageFirstID, nextPageFirstID, "Next page should have different items")
			}
		}
	})

	t.Run("TestUpdateExistingProperties", func(t *testing.T) {
		// Create metric history with initial properties
		metricHistory := &models.MetricHistoryImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricHistoryAttributes{
				Name:         apiutils.Of("update-props-metric-history"),
				URI:          apiutils.Of("s3://bucket/update-props-metric-history.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric-history"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(1)),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.5),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "team",
					StringValue:      apiutils.Of("team-a"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)

		// Update properties
		metricHistory.ID = saved.GetID()
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		metricHistory.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch
		metricHistory.Properties = &[]models.Properties{
			{
				Name:     "step",
				IntValue: apiutils.Of(int32(2)), // Updated value
			},
			{
				Name:        "value",
				DoubleValue: apiutils.Of(0.8), // Updated value
			},
			{
				Name:        "new_prop",
				StringValue: apiutils.Of("new_value"), // New property
			},
		}
		metricHistory.CustomProperties = &[]models.Properties{
			{
				Name:             "team",
				StringValue:      apiutils.Of("team-b"), // Updated custom property
				IsCustomProperty: true,
			},
		}

		updated, err := repo.Save(metricHistory, nil)
		require.NoError(t, err)

		// Verify properties were updated
		retrieved, err := repo.GetByID(*updated.GetID())
		require.NoError(t, err)

		assert.Len(t, *retrieved.GetProperties(), 3)       // step, value, new_prop
		assert.Len(t, *retrieved.GetCustomProperties(), 1) // team

		// Verify updated values
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "step":
				assert.Equal(t, int32(2), *prop.IntValue)
			case "value":
				assert.Equal(t, 0.8, *prop.DoubleValue)
			case "new_prop":
				assert.Equal(t, "new_value", *prop.StringValue)
			}
		}

		for _, prop := range *retrieved.GetCustomProperties() {
			if prop.Name == "team" {
				assert.Equal(t, "team-b", *prop.StringValue)
			}
		}
	})

	t.Run("TestListByStepIds", func(t *testing.T) {
		// Create metric histories with different step values
		testMetricHistories := []*models.MetricHistoryImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("step-metric-history-1"),
					ExternalID:   apiutils.Of("step-metric-history-ext-1"),
					URI:          apiutils.Of("s3://bucket/step-metric-history-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "step",
						IntValue: apiutils.Of(int32(1)),
					},
					{
						Name:        "description",
						StringValue: apiutils.Of("Metric history for step 1"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("step-metric-history-2"),
					ExternalID:   apiutils.Of("step-metric-history-ext-2"),
					URI:          apiutils.Of("s3://bucket/step-metric-history-2.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "step",
						IntValue: apiutils.Of(int32(2)),
					},
					{
						Name:        "description",
						StringValue: apiutils.Of("Metric history for step 2"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricHistoryAttributes{
					Name:         apiutils.Of("step-metric-history-3"),
					ExternalID:   apiutils.Of("step-metric-history-ext-3"),
					URI:          apiutils.Of("s3://bucket/step-metric-history-3.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric-history"),
				},
				Properties: &[]models.Properties{
					{
						Name:     "step",
						IntValue: apiutils.Of(int32(3)),
					},
					{
						Name:        "description",
						StringValue: apiutils.Of("Metric history for step 3"),
					},
				},
			},
		}

		// Save all metric histories
		for _, metricHistory := range testMetricHistories {
			_, err := repo.Save(metricHistory, nil)
			require.NoError(t, err)
		}

		// Test filtering by single step ID
		pageSize := int32(10)
		stepIds := "1"
		listOptions := models.MetricHistoryListOptions{
			StepIds: &stepIds,
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 1, len(result.Items), "should return 1 metric history for step 1")
		if len(result.Items) > 0 {
			assert.Equal(t, "step-metric-history-1", *result.Items[0].GetAttributes().Name)
		}

		// Test filtering by multiple step IDs
		stepIds = "1,3"
		listOptions = models.MetricHistoryListOptions{
			StepIds: &stepIds,
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items), "should return 2 metric histories for steps 1 and 3")

		// Verify the returned metric histories are from steps 1 and 3
		names := make(map[string]bool)
		for _, item := range result.Items {
			names[*item.GetAttributes().Name] = true
		}
		assert.True(t, names["step-metric-history-1"], "should contain metric history from step 1")
		assert.True(t, names["step-metric-history-3"], "should contain metric history from step 3")
		assert.False(t, names["step-metric-history-2"], "should not contain metric history from step 2")

		// Test filtering by non-existent step ID
		stepIds = "999"
		listOptions = models.MetricHistoryListOptions{
			StepIds: &stepIds,
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items), "should return 0 metric histories for non-existent step")

		// Test with empty string
		emptyStepIds := ""
		listOptions = models.MetricHistoryListOptions{
			StepIds: &emptyStepIds,
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		// Should return all metric histories since no filter is applied
		assert.GreaterOrEqual(t, len(result.Items), 3)

		// Test with whitespace-only values (should be ignored)
		whitespaceStepIds := "1, ,3"
		listOptions = models.MetricHistoryListOptions{
			StepIds: &whitespaceStepIds,
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items), "should return 2 metric histories for steps 1 and 3, ignoring whitespace")

		// Test with leading/trailing whitespace (should be trimmed)
		trimStepIds := " 1 , 3 "
		listOptions = models.MetricHistoryListOptions{
			StepIds: &trimStepIds,
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 2, len(result.Items), "should return 2 metric histories for steps 1 and 3, trimming whitespace")
	})
}
