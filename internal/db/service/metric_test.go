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

func TestMetricRepository(t *testing.T) {
	db, cleanup := setupTestDB(t)
	defer cleanup()

	// Get the actual Metric type ID from the database
	typeID := getMetricTypeID(t, db)
	repo := service.NewMetricRepository(db, typeID)

	// Also get model version type for creating related entities
	registeredModelTypeID := getRegisteredModelTypeID(t, db)
	registeredModelRepo := service.NewRegisteredModelRepository(db, registeredModelTypeID)

	modelVersionTypeID := getModelVersionTypeID(t, db)
	modelVersionRepo := service.NewModelVersionRepository(db, modelVersionTypeID)

	t.Run("TestSave", func(t *testing.T) {
		// Test creating a new metric
		metric := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of("test-metric"),
				ExternalID:   apiutils.Of("metric-ext-123"),
				URI:          apiutils.Of("s3://bucket/metric.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Test metric description"),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.95),
				},
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(100)),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "custom-metric-prop",
					StringValue:      apiutils.Of("custom-metric-value"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(metric, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, "test-metric", *saved.GetAttributes().Name)
		assert.Equal(t, "metric-ext-123", *saved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/metric.json", *saved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *saved.GetAttributes().State)
		assert.Equal(t, "metric", *saved.GetAttributes().ArtifactType)
		assert.NotNil(t, saved.GetAttributes().CreateTimeSinceEpoch)
		assert.NotNil(t, saved.GetAttributes().LastUpdateTimeSinceEpoch)

		// Verify properties were saved
		assert.NotNil(t, saved.GetProperties())
		assert.Len(t, *saved.GetProperties(), 3) // description, value, step

		// Verify specific properties
		var foundDescription, foundValue, foundStep bool
		for _, prop := range *saved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Test metric description", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, 0.95, *prop.DoubleValue)
			case "step":
				foundStep = true
				assert.Equal(t, int32(100), *prop.IntValue)
			}
		}
		assert.True(t, foundDescription, "description property should exist")
		assert.True(t, foundValue, "value property should exist")
		assert.True(t, foundStep, "step property should exist")

		// Verify custom properties were saved
		assert.NotNil(t, saved.GetCustomProperties())
		assert.Len(t, *saved.GetCustomProperties(), 1)

		var foundCustomProp bool
		for _, prop := range *saved.GetCustomProperties() {
			if prop.Name == "custom-metric-prop" {
				foundCustomProp = true
				assert.Equal(t, "custom-metric-value", *prop.StringValue)
				assert.True(t, prop.IsCustomProperty)
			}
		}
		assert.True(t, foundCustomProp, "custom-metric-prop should exist")

		// Test updating the same metric
		metric.ID = saved.GetID()
		metric.GetAttributes().Name = apiutils.Of("updated-metric")
		metric.GetAttributes().State = apiutils.Of("PENDING")
		// Preserve CreateTimeSinceEpoch from the saved entity (simulating what OpenAPI converter would do)
		metric.GetAttributes().CreateTimeSinceEpoch = saved.GetAttributes().CreateTimeSinceEpoch

		updated, err := repo.Save(metric, nil)
		require.NoError(t, err)
		require.NotNil(t, updated)
		assert.Equal(t, *saved.GetID(), *updated.GetID())
		assert.Equal(t, "updated-metric", *updated.GetAttributes().Name)
		assert.Equal(t, "PENDING", *updated.GetAttributes().State)
		assert.Equal(t, *saved.GetAttributes().CreateTimeSinceEpoch, *updated.GetAttributes().CreateTimeSinceEpoch)
		assert.Greater(t, *updated.GetAttributes().LastUpdateTimeSinceEpoch, *saved.GetAttributes().LastUpdateTimeSinceEpoch)
	})

	t.Run("TestGetByID", func(t *testing.T) {
		// First create a metric to retrieve
		metric := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of("get-test-metric"),
				ExternalID:   apiutils.Of("get-metric-ext-123"),
				URI:          apiutils.Of("s3://bucket/get-metric.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric for get test"),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.85),
				},
				{
					Name:     "step",
					IntValue: apiutils.Of(int32(50)),
				},
			},
			CustomProperties: &[]models.Properties{
				{
					Name:             "test-category",
					StringValue:      apiutils.Of("retrieval-test"),
					IsCustomProperty: true,
				},
			},
		}

		saved, err := repo.Save(metric, nil)
		require.NoError(t, err)
		require.NotNil(t, saved.GetID())

		// Test retrieving by ID
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)
		require.NotNil(t, retrieved)
		assert.Equal(t, *saved.GetID(), *retrieved.GetID())
		assert.Equal(t, "get-test-metric", *retrieved.GetAttributes().Name)
		assert.Equal(t, "get-metric-ext-123", *retrieved.GetAttributes().ExternalID)
		assert.Equal(t, "s3://bucket/get-metric.json", *retrieved.GetAttributes().URI)
		assert.Equal(t, "LIVE", *retrieved.GetAttributes().State)

		// Verify type-specific properties were retrieved
		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3)

		var foundDescription, foundValue, foundStep bool
		for _, prop := range *retrieved.GetProperties() {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Metric for get test", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, 0.85, *prop.DoubleValue)
			case "step":
				foundStep = true
				assert.Equal(t, int32(50), *prop.IntValue)
			}
		}
		assert.True(t, foundDescription, "description should be retrieved")
		assert.True(t, foundValue, "value should be retrieved")
		assert.True(t, foundStep, "step should be retrieved")

		// Verify custom properties were retrieved
		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 1)

		customProp := (*retrieved.GetCustomProperties())[0]
		assert.Equal(t, "test-category", customProp.Name)
		assert.Equal(t, "retrieval-test", *customProp.StringValue)
		assert.True(t, customProp.IsCustomProperty)

		// Test retrieving non-existent ID
		_, err = repo.GetByID(99999)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "metric by id not found")
	})

	t.Run("TestList", func(t *testing.T) {
		// Create multiple metrics for listing
		testMetrics := []*models.MetricImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of("list-metric-1"),
					ExternalID:   apiutils.Of("list-metric-ext-1"),
					URI:          apiutils.Of("s3://bucket/list-metric-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("First metric"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of("list-metric-2"),
					ExternalID:   apiutils.Of("list-metric-ext-2"),
					URI:          apiutils.Of("s3://bucket/list-metric-2.json"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("metric"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Second metric"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of("list-metric-3"),
					ExternalID:   apiutils.Of("list-metric-ext-3"),
					URI:          apiutils.Of("s3://bucket/list-metric-3.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Third metric"),
					},
				},
			},
		}

		for _, metric := range testMetrics {
			_, err := repo.Save(metric, nil)
			require.NoError(t, err)
		}

		// Test listing all metrics with basic pagination
		pageSize := int32(10)
		listOptions := models.MetricListOptions{}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 3) // At least our 3 test metrics

		// Test listing by name
		listOptions = models.MetricListOptions{
			Name: apiutils.Of("list-metric-1"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-metric-1", *result.Items[0].GetAttributes().Name)
		}

		// Test listing by external ID
		listOptions = models.MetricListOptions{
			ExternalID: apiutils.Of("list-metric-ext-2"),
		}
		listOptions.PageSize = &pageSize

		result, err = repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		if len(result.Items) > 0 {
			assert.Equal(t, 1, len(result.Items))
			assert.Equal(t, "list-metric-ext-2", *result.Items[0].GetAttributes().ExternalID)
		}

		// Test ordering by ID (deterministic)
		listOptions = models.MetricListOptions{
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
		// Create metrics sequentially with time delays to ensure deterministic ordering
		metric1 := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of("time-test-metric-1"),
				URI:          apiutils.Of("s3://bucket/time-metric-1.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric"),
			},
		}
		saved1, err := repo.Save(metric1, nil)
		require.NoError(t, err)

		// Small delay to ensure different timestamps
		time.Sleep(10 * time.Millisecond)

		metric2 := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of("time-test-metric-2"),
				URI:          apiutils.Of("s3://bucket/time-metric-2.json"),
				State:        apiutils.Of("PENDING"),
				ArtifactType: apiutils.Of("metric"),
			},
		}
		saved2, err := repo.Save(metric2, nil)
		require.NoError(t, err)

		// Test ordering by CREATE_TIME
		pageSize := int32(10)
		listOptions := models.MetricListOptions{
			Pagination: models.Pagination{
				OrderBy: apiutils.Of("CREATE_TIME"),
			},
		}
		listOptions.PageSize = &pageSize

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Find our test metrics in the results
		var foundMetric1, foundMetric2 models.Metric
		var index1, index2 = -1, -1

		for i, item := range result.Items {
			if *item.GetID() == *saved1.GetID() {
				foundMetric1 = item
				index1 = i
			}
			if *item.GetID() == *saved2.GetID() {
				foundMetric2 = item
				index2 = i
			}
		}

		// Verify both metrics were found and metric1 comes before metric2 (ascending order)
		require.NotEqual(t, -1, index1, "Metric 1 should be found in results")
		require.NotEqual(t, -1, index2, "Metric 2 should be found in results")
		assert.Less(t, index1, index2, "Metric 1 should come before Metric 2 when ordered by CREATE_TIME")
		assert.Less(t, *foundMetric1.GetAttributes().CreateTimeSinceEpoch, *foundMetric2.GetAttributes().CreateTimeSinceEpoch, "Metric 1 should have earlier create time")
	})

	t.Run("TestSaveWithModelVersion", func(t *testing.T) {
		// First create a registered model and model version
		registeredModel := &models.RegisteredModelImpl{
			TypeID: apiutils.Of(int32(registeredModelTypeID)),
			Attributes: &models.RegisteredModelAttributes{
				Name: apiutils.Of("test-registered-model-for-metric"),
			},
		}
		savedRegisteredModel, err := registeredModelRepo.Save(registeredModel)
		require.NoError(t, err)

		modelVersion := &models.ModelVersionImpl{
			TypeID: apiutils.Of(int32(modelVersionTypeID)),
			Attributes: &models.ModelVersionAttributes{
				Name: apiutils.Of("test-model-version-for-metric"),
			},
			Properties: &[]models.Properties{
				{
					Name:     "registered_model_id",
					IntValue: savedRegisteredModel.GetID(),
				},
			},
		}
		savedModelVersion, err := modelVersionRepo.Save(modelVersion)
		require.NoError(t, err)

		// Test creating a metric with model version attribution
		metric := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of(fmt.Sprintf("%d:model-version-metric", *savedModelVersion.GetID())),
				URI:          apiutils.Of("s3://bucket/model-version-metric.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric associated with model version"),
				},
			},
		}

		saved, err := repo.Save(metric, savedModelVersion.GetID())
		require.NoError(t, err)
		require.NotNil(t, saved)
		require.NotNil(t, saved.GetID())
		assert.Equal(t, fmt.Sprintf("%d:model-version-metric", *savedModelVersion.GetID()), *saved.GetAttributes().Name)
		assert.Equal(t, "s3://bucket/model-version-metric.json", *saved.GetAttributes().URI)

		// Test listing by model version ID
		listOptions := models.MetricListOptions{
			ParentResourceID: savedModelVersion.GetID(),
		}
		listOptions.PageSize = apiutils.Of(int32(10))

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1) // Should find our metric

		// Verify the found metric
		found := false
		for _, item := range result.Items {
			if *item.GetID() == *saved.GetID() {
				found = true
				break
			}
		}
		assert.True(t, found, "Should find the metric associated with the model version")
	})

	t.Run("TestSaveWithProperties", func(t *testing.T) {
		metric := &models.MetricImpl{
			TypeID: apiutils.Of(int32(typeID)),
			Attributes: &models.MetricAttributes{
				Name:         apiutils.Of("props-test-metric"),
				ExternalID:   apiutils.Of("props-metric-ext-123"),
				URI:          apiutils.Of("s3://bucket/props-metric.json"),
				State:        apiutils.Of("LIVE"),
				ArtifactType: apiutils.Of("metric"),
			},
			Properties: &[]models.Properties{
				{
					Name:        "description",
					StringValue: apiutils.Of("Metric with properties"),
				},
				{
					Name:        "value",
					DoubleValue: apiutils.Of(0.95),
				},
				{
					Name:        "threshold",
					DoubleValue: apiutils.Of(0.90),
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

		saved, err := repo.Save(metric, nil)
		require.NoError(t, err)
		require.NotNil(t, saved)

		// Verify properties were saved
		retrieved, err := repo.GetByID(*saved.GetID())
		require.NoError(t, err)

		assert.NotNil(t, retrieved.GetProperties())
		assert.Len(t, *retrieved.GetProperties(), 3) // description, value, threshold

		assert.NotNil(t, retrieved.GetCustomProperties())
		assert.Len(t, *retrieved.GetCustomProperties(), 2)

		// Verify specific properties exist
		properties := *retrieved.GetProperties()
		var foundDescription, foundValue, foundThreshold bool
		for _, prop := range properties {
			switch prop.Name {
			case "description":
				foundDescription = true
				assert.Equal(t, "Metric with properties", *prop.StringValue)
			case "value":
				foundValue = true
				assert.Equal(t, 0.95, *prop.DoubleValue)
			case "threshold":
				foundThreshold = true
				assert.Equal(t, 0.90, *prop.DoubleValue)
			}
		}
		assert.True(t, foundDescription, "Should find description property")
		assert.True(t, foundValue, "Should find value property")
		assert.True(t, foundThreshold, "Should find threshold property")
	})

	t.Run("TestPagination", func(t *testing.T) {
		// Create multiple metrics for pagination testing
		for i := 0; i < 5; i++ {
			metric := &models.MetricImpl{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of(fmt.Sprintf("page-metric-%d", i)),
					URI:          apiutils.Of(fmt.Sprintf("s3://bucket/page-metric-%d.json", i)),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric"),
				},
			}
			_, err := repo.Save(metric, nil)
			require.NoError(t, err)
		}

		// Test pagination with page size 2
		pageSize := int32(2)
		listOptions := models.MetricListOptions{
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
		}
	})

	t.Run("TestEmptyResults", func(t *testing.T) {
		// Test listing with filter that returns no results
		listOptions := models.MetricListOptions{
			Name: apiutils.Of("non-existent-metric"),
		}

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)
		assert.Equal(t, 0, len(result.Items))
		assert.Equal(t, int32(0), result.Size)
		assert.Empty(t, result.NextPageToken)
	})

	t.Run("TestListByState", func(t *testing.T) {
		// Create metrics with different states for testing
		metrics := []*models.MetricImpl{
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of("state-metric-1"),
					URI:          apiutils.Of("s3://bucket/state-metric-1.json"),
					State:        apiutils.Of("LIVE"),
					ArtifactType: apiutils.Of("metric"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Live metric"),
					},
				},
			},
			{
				TypeID: apiutils.Of(int32(typeID)),
				Attributes: &models.MetricAttributes{
					Name:         apiutils.Of("state-metric-2"),
					URI:          apiutils.Of("s3://bucket/state-metric-2.json"),
					State:        apiutils.Of("PENDING"),
					ArtifactType: apiutils.Of("metric"),
				},
				Properties: &[]models.Properties{
					{
						Name:        "description",
						StringValue: apiutils.Of("Pending metric"),
					},
				},
			},
		}

		for _, metric := range metrics {
			_, err := repo.Save(metric, nil)
			require.NoError(t, err)
		}

		// Test listing all metrics (should find at least our 2 test metrics)
		listOptions := models.MetricListOptions{}
		listOptions.PageSize = apiutils.Of(int32(10))

		result, err := repo.List(listOptions)
		require.NoError(t, err)
		require.NotNil(t, result)

		// Should find at least one metric with each state
		foundLive := false
		foundPending := false
		for _, item := range result.Items {
			if item.GetAttributes().State != nil {
				if *item.GetAttributes().State == "LIVE" {
					foundLive = true
				}
				if *item.GetAttributes().State == "PENDING" {
					foundPending = true
				}
			}
		}
		assert.True(t, foundLive, "Should find metric with LIVE state")
		assert.True(t, foundPending, "Should find metric with PENDING state")
	})
}
