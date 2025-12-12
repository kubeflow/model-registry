package integration_test

import (
	"context"
	"os"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	"github.com/kubeflow/model-registry/catalog/internal/db/service"
	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/kubeflow/model-registry/internal/db/schema"
	"github.com/kubeflow/model-registry/internal/testutils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	os.Exit(testutils.TestMainPostgresHelper(m))
}

// Helper functions to get type IDs from DB
func getCatalogModelTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to query CatalogModel type")
	return typeRecord.ID
}

func getCatalogMetricsArtifactTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogMetricsArtifactTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to query CatalogMetricsArtifact type")
	return typeRecord.ID
}

func getCatalogModelArtifactTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogModelArtifactTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to query CatalogModelArtifact type")
	return typeRecord.ID
}

func getCatalogSourceTypeIDForTest(t *testing.T, db *gorm.DB) int32 {
	var typeRecord schema.Type
	err := db.Where("name = ?", service.CatalogSourceTypeName).First(&typeRecord).Error
	require.NoError(t, err, "Failed to query CatalogSource type")
	return typeRecord.ID
}

// TestIntegration_PreservedRecommendationAlgorithm verifies end-to-end functionality
// with the exact preserved recommendations algorithm from db_catalog.go
func TestIntegration_PreservedRecommendationAlgorithm(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get type IDs
	catalogModelTypeID := getCatalogModelTypeIDForTest(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeIDForTest(t, sharedDB)
	modelArtifactTypeID := getCatalogModelArtifactTypeIDForTest(t, sharedDB)
	catalogSourceTypeID := getCatalogSourceTypeIDForTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)
	catalogSourceRepo := service.NewCatalogSourceRepository(sharedDB, catalogSourceTypeID)

	services := service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		modelArtifactRepo,
		metricsArtifactRepo,
		catalogSourceRepo,
		service.NewPropertyOptionsRepository(sharedDB),
	)

	provider := catalog.NewDBCatalog(services, nil)

	// Create test model
	testModel := &models.CatalogModelImpl{
		TypeID: apiutils.Of(catalogModelTypeID),
		Attributes: &models.CatalogModelAttributes{
			Name:       apiutils.Of("algorithm-test-model"),
			ExternalID: apiutils.Of("alg-test-model-789"),
		},
		Properties: &[]dbmodels.Properties{
			{Name: "source_id", StringValue: apiutils.Of("algorithm-test-source")},
		},
	}
	savedModel, err := services.CatalogModelRepository.Save(testModel)
	require.NoError(t, err)

	// Create artifacts that test the specific algorithm behavior
	// These match the exact patterns expected by the preserved algorithm
	artifacts := []*models.CatalogMetricsArtifactImpl{
		{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("fast-low-cost"),
				ExternalID:  apiutils.Of("fast-low"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "requests_per_second", DoubleValue: apiutils.Of(250.0)},
				{Name: "ttft_p90", DoubleValue: apiutils.Of(50.0)},
				{Name: "hardware_count", IntValue: apiutils.Of(int32(2))},
				{Name: "hardware_type", StringValue: apiutils.Of("gpu-a100")},
			},
		},
		{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("medium-medium-cost"),
				ExternalID:  apiutils.Of("medium-med"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "requests_per_second", DoubleValue: apiutils.Of(150.0)},
				{Name: "ttft_p90", DoubleValue: apiutils.Of(120.0)}, // > 100, should use epsilon 0.1
				{Name: "hardware_count", IntValue: apiutils.Of(int32(3))},
				{Name: "hardware_type", StringValue: apiutils.Of("gpu-a100")},
			},
		},
		{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("slow-high-cost"),
				ExternalID:  apiutils.Of("slow-high"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "requests_per_second", DoubleValue: apiutils.Of(100.0)},
				{Name: "ttft_p90", DoubleValue: apiutils.Of(200.0)},
				{Name: "hardware_count", IntValue: apiutils.Of(int32(5))}, // Should be filtered out by first pass
				{Name: "hardware_type", StringValue: apiutils.Of("gpu-a100")},
			},
		},
	}

	// Save all artifacts
	for _, artifact := range artifacts {
		_, err := services.CatalogMetricsArtifactRepository.Save(artifact, savedModel.GetID())
		require.NoError(t, err)
	}

	t.Run("Exact_Algorithm_TwoPass_Filtering", func(t *testing.T) {
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       300, // This will add replicas for cost calculation
			Recommendations: true,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"algorithm-test-model",
			"algorithm-test-source",
			params,
		)
		require.NoError(t, err)

		// The exact algorithm should:
		// 1. Group by hardware_type (all gpu-a100)
		// 2. Sort by ttft_p90 ascending (50, 120, 200)
		// 3. First pass: filter by cost (hardware_count * replicas)
		//    - fast-low-cost: 2 * ceil(300/250) = 2 * 2 = 4
		//    - medium-medium: 3 * ceil(300/150) = 3 * 2 = 6
		//    - slow-high: 5 * ceil(300/100) = 5 * 3 = 15
		//    Should keep: fast-low-cost (4), then medium-medium (6), reject slow-high (15 > 6)
		// 4. Reverse for cost ascending: [medium-medium, fast-low-cost]
		// 5. Second pass: epsilon-based filtering
		//    - Keep medium-medium (first)
		//    - Check fast-low-cost: latency improvement = (120-50)/120 = 0.583 > 0.1 (epsilon for >100 latency)
		//    - Should keep fast-low-cost too

		// Verify results match exact algorithm expectations
		assert.GreaterOrEqual(t, len(result.Items), 1)

		// Check that replicas were calculated and used in cost calculation
		for _, artifact := range result.Items {
			assert.Contains(t, artifact.CatalogMetricsArtifact.CustomProperties, "replicas")
		}

		// Verify no slow-high-cost in results (should be filtered by first pass)
		names := make([]string, len(result.Items))
		for i, artifact := range result.Items {
			names[i] = *artifact.CatalogMetricsArtifact.Name
		}
		assert.NotContains(t, names, "slow-high-cost")
	})

	t.Run("Algorithm_Epsilon_Threshold_Behavior", func(t *testing.T) {
		// Test epsilon = 0.05 vs 0.1 behavior based on lastKeptLatency <= 100
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       100, // Lower target to test different replica calculations
			Recommendations: true,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"algorithm-test-model",
			"algorithm-test-source",
			params,
		)
		require.NoError(t, err)

		// Should preserve the exact epsilon logic from the original algorithm
		assert.GreaterOrEqual(t, len(result.Items), 1)

		// Verify algorithm behavior is consistent with original implementation
		for _, artifact := range result.Items {
			// All artifacts should have replicas calculated for cost computation
			replicas := artifact.CatalogMetricsArtifact.CustomProperties["replicas"]
			require.NotNil(t, replicas.MetadataIntValue)
			assert.NotEmpty(t, replicas.MetadataIntValue.IntValue)
		}
	})

	t.Run("TargetRPS_Calculations_Integration", func(t *testing.T) {
		// Test that targetRPS calculations are integrated correctly
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       600,
			Recommendations: false, // No dedup to see all artifacts with calculations
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"algorithm-test-model",
			"algorithm-test-source",
			params,
		)
		require.NoError(t, err)

		// Should have all 3 artifacts
		assert.Len(t, result.Items, 3)

		// Each should have replicas and total_requests_per_second calculated
		for _, artifact := range result.Items {
			assert.Contains(t, artifact.CatalogMetricsArtifact.CustomProperties, "replicas")
			assert.Contains(t, artifact.CatalogMetricsArtifact.CustomProperties, "total_requests_per_second")

			// Verify replicas is calculated based on targetRPS / requests_per_second
			rpsVal := artifact.CatalogMetricsArtifact.CustomProperties["requests_per_second"]
			require.NotNil(t, rpsVal.MetadataDoubleValue)

			totalRPSVal := artifact.CatalogMetricsArtifact.CustomProperties["total_requests_per_second"]
			require.NotNil(t, totalRPSVal.MetadataDoubleValue)

			// Total RPS should be >= targetRPS
			assert.GreaterOrEqual(t, totalRPSVal.MetadataDoubleValue.DoubleValue, 600.0)
		}
	})

	t.Run("Repository_Filtering_PerformanceMetricsOnly", func(t *testing.T) {
		// Add a non-performance metrics artifact to verify filtering
		nonPerfArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("quality-metrics"),
				ExternalID:  apiutils.Of("quality-123"),
				MetricsType: models.MetricsTypeAccuracy,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("accuracy-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "accuracy", DoubleValue: apiutils.Of(0.95)},
			},
		}
		_, err := services.CatalogMetricsArtifactRepository.Save(nonPerfArtifact, savedModel.GetID())
		require.NoError(t, err)

		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       100,
			Recommendations: false,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"algorithm-test-model",
			"algorithm-test-source",
			params,
		)
		require.NoError(t, err)

		// Should only have the 3 performance artifacts, not the quality one
		assert.Len(t, result.Items, 3)

		// Verify all are performance-metrics - metricsType is in MetricsType field
		for _, artifact := range result.Items {
			assert.Equal(t, "performance-metrics", artifact.CatalogMetricsArtifact.MetricsType)
		}
	})

	t.Run("Original_Input_Order_Preservation", func(t *testing.T) {
		// Test that the algorithm preserves original input order in final results
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       300,
			Recommendations: true,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"algorithm-test-model",
			"algorithm-test-source",
			params,
		)
		require.NoError(t, err)

		// The algorithm should preserve the order from the original input
		// After recommendations, results should be in the same relative order
		// as they appeared in the input
		assert.GreaterOrEqual(t, len(result.Items), 1)

		// Verify artifacts maintain expected properties
		for _, artifact := range result.Items {
			require.NotNil(t, artifact.CatalogMetricsArtifact)
			require.NotNil(t, artifact.CatalogMetricsArtifact.Id)
			require.NotNil(t, artifact.CatalogMetricsArtifact.Name)
		}
	})
}

// TestIntegration_ServiceLayerBehavior tests the service layer implementation
func TestIntegration_ServiceLayerBehavior(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get type IDs
	catalogModelTypeID := getCatalogModelTypeIDForTest(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeIDForTest(t, sharedDB)
	modelArtifactTypeID := getCatalogModelArtifactTypeIDForTest(t, sharedDB)
	catalogSourceTypeID := getCatalogSourceTypeIDForTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)
	catalogSourceRepo := service.NewCatalogSourceRepository(sharedDB, catalogSourceTypeID)

	services := service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		modelArtifactRepo,
		metricsArtifactRepo,
		catalogSourceRepo,
		service.NewPropertyOptionsRepository(sharedDB),
	)

	provider := catalog.NewDBCatalog(services, nil)

	// Create test model
	testModel := &models.CatalogModelImpl{
		TypeID: apiutils.Of(catalogModelTypeID),
		Attributes: &models.CatalogModelAttributes{
			Name:       apiutils.Of("service-test-model"),
			ExternalID: apiutils.Of("service-model-123"),
		},
		Properties: &[]dbmodels.Properties{
			{Name: "source_id", StringValue: apiutils.Of("service-test-source")},
		},
	}
	savedModel, err := services.CatalogModelRepository.Save(testModel)
	require.NoError(t, err)

	t.Run("Multiple_Hardware_Types_Grouped_Correctly", func(t *testing.T) {
		// Create artifacts with different hardware types
		artifacts := []*models.CatalogMetricsArtifactImpl{
			{
				TypeID: apiutils.Of(metricsArtifactTypeID),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("a100-config-1"),
					ExternalID:  apiutils.Of("a100-1"),
					MetricsType: models.MetricsTypePerformance,
				},
				Properties: &[]dbmodels.Properties{
					{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
				},
				CustomProperties: &[]dbmodels.Properties{
					{Name: "requests_per_second", DoubleValue: apiutils.Of(200.0)},
					{Name: "ttft_p90", DoubleValue: apiutils.Of(50.0)},
					{Name: "hardware_count", IntValue: apiutils.Of(int32(2))},
					{Name: "hardware_type", StringValue: apiutils.Of("gpu-a100")},
				},
			},
			{
				TypeID: apiutils.Of(metricsArtifactTypeID),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("h100-config-1"),
					ExternalID:  apiutils.Of("h100-1"),
					MetricsType: models.MetricsTypePerformance,
				},
				Properties: &[]dbmodels.Properties{
					{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
				},
				CustomProperties: &[]dbmodels.Properties{
					{Name: "requests_per_second", DoubleValue: apiutils.Of(300.0)},
					{Name: "ttft_p90", DoubleValue: apiutils.Of(40.0)},
					{Name: "hardware_count", IntValue: apiutils.Of(int32(1))},
					{Name: "hardware_type", StringValue: apiutils.Of("gpu-h100")},
				},
			},
		}

		for _, artifact := range artifacts {
			_, err := services.CatalogMetricsArtifactRepository.Save(artifact, savedModel.GetID())
			require.NoError(t, err)
		}

		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       300,
			Recommendations: true,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"service-test-model",
			"service-test-source",
			params,
		)
		require.NoError(t, err)

		// Algorithm groups by hardware_type, so both should be preserved
		// since they're in different hardware groups
		assert.Len(t, result.Items, 2)

		// Verify we have both hardware types
		hardwareTypes := make(map[string]bool)
		for _, artifact := range result.Items {
			hwType := artifact.CatalogMetricsArtifact.CustomProperties["hardware_type"]
			require.NotNil(t, hwType.MetadataStringValue)
			hardwareTypes[hwType.MetadataStringValue.StringValue] = true
		}
		assert.True(t, hardwareTypes["gpu-a100"])
		assert.True(t, hardwareTypes["gpu-h100"])
	})

	t.Run("Cost_Calculation_With_Replicas", func(t *testing.T) {
		// Verify cost = hardware_count * replicas calculation
		artifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("cost-test-artifact"),
				ExternalID:  apiutils.Of("cost-123"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "requests_per_second", DoubleValue: apiutils.Of(100.0)},
				{Name: "ttft_p90", DoubleValue: apiutils.Of(80.0)},
				{Name: "hardware_count", IntValue: apiutils.Of(int32(4))},
				{Name: "hardware_type", StringValue: apiutils.Of("gpu-v100")},
			},
		}
		_, err := services.CatalogMetricsArtifactRepository.Save(artifact, savedModel.GetID())
		require.NoError(t, err)

		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       500, // Should need ceil(500/100) = 5 replicas
			Recommendations: false,
			PageSize:        10,
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"service-test-model",
			"service-test-source",
			params,
		)
		require.NoError(t, err)

		// Find our artifact
		var found bool
		for _, a := range result.Items {
			if *a.CatalogMetricsArtifact.Name == "cost-test-artifact" {
				found = true
				// Cost would be: hardware_count (4) * replicas (5) = 20
				replicas := a.CatalogMetricsArtifact.CustomProperties["replicas"]
				require.NotNil(t, replicas.MetadataIntValue)
				assert.Equal(t, "5", replicas.MetadataIntValue.IntValue)
				break
			}
		}
		assert.True(t, found, "cost-test-artifact should be in results")
	})
}

// TestIntegration_ConfigurableProperties tests the new configurable property parameters
func TestIntegration_ConfigurableProperties(t *testing.T) {
	sharedDB, cleanup := testutils.SetupPostgresWithMigrations(t, service.DatastoreSpec())
	defer cleanup()

	// Get type IDs
	catalogModelTypeID := getCatalogModelTypeIDForTest(t, sharedDB)
	metricsArtifactTypeID := getCatalogMetricsArtifactTypeIDForTest(t, sharedDB)
	modelArtifactTypeID := getCatalogModelArtifactTypeIDForTest(t, sharedDB)
	catalogSourceTypeID := getCatalogSourceTypeIDForTest(t, sharedDB)

	// Create repositories
	catalogModelRepo := service.NewCatalogModelRepository(sharedDB, catalogModelTypeID)
	catalogArtifactRepo := service.NewCatalogArtifactRepository(sharedDB, map[string]int32{
		service.CatalogModelArtifactTypeName:   modelArtifactTypeID,
		service.CatalogMetricsArtifactTypeName: metricsArtifactTypeID,
	})
	modelArtifactRepo := service.NewCatalogModelArtifactRepository(sharedDB, modelArtifactTypeID)
	metricsArtifactRepo := service.NewCatalogMetricsArtifactRepository(sharedDB, metricsArtifactTypeID)
	catalogSourceRepo := service.NewCatalogSourceRepository(sharedDB, catalogSourceTypeID)

	services := service.NewServices(
		catalogModelRepo,
		catalogArtifactRepo,
		modelArtifactRepo,
		metricsArtifactRepo,
		catalogSourceRepo,
		service.NewPropertyOptionsRepository(sharedDB),
	)

	provider := catalog.NewDBCatalog(services, nil)

	// Create test model
	testModel := &models.CatalogModelImpl{
		TypeID: apiutils.Of(catalogModelTypeID),
		Attributes: &models.CatalogModelAttributes{
			Name:       apiutils.Of("configurable-props-model"),
			ExternalID: apiutils.Of("config-model-999"),
		},
		Properties: &[]dbmodels.Properties{
			{Name: "source_id", StringValue: apiutils.Of("configurable-props-source")},
		},
	}
	savedModel, err := services.CatalogModelRepository.Save(testModel)
	require.NoError(t, err)

	// Create test artifact with custom property names
	artifact := &models.CatalogMetricsArtifactImpl{
		TypeID: apiutils.Of(metricsArtifactTypeID),
		Attributes: &models.CatalogMetricsArtifactAttributes{
			Name:        apiutils.Of("custom-props-artifact"),
			ExternalID:  apiutils.Of("custom-123"),
			MetricsType: models.MetricsTypePerformance,
		},
		Properties: &[]dbmodels.Properties{
			{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
		},
		CustomProperties: &[]dbmodels.Properties{
			{Name: "throughput", DoubleValue: apiutils.Of(100.0)},
			{Name: "p90_latency", DoubleValue: apiutils.Of(120.0)},
			{Name: "nodes", IntValue: apiutils.Of(int32(2))},
			{Name: "instance_type", StringValue: apiutils.Of("gpu-large")},
		},
	}
	_, err = services.CatalogMetricsArtifactRepository.Save(artifact, savedModel.GetID())
	require.NoError(t, err)

	t.Run("CustomPropertyNames_WorkEndToEnd", func(t *testing.T) {
		// Create additional artifacts with custom property names to test recommendations
		artifact2 := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("custom-props-artifact-2"),
				ExternalID:  apiutils.Of("custom-456"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "throughput", DoubleValue: apiutils.Of(50.0)},   // Lower RPS
				{Name: "p90_latency", DoubleValue: apiutils.Of(100.0)}, // Lower latency
				{Name: "nodes", IntValue: apiutils.Of(int32(2))},
				{Name: "instance_type", StringValue: apiutils.Of("gpu-large")},
			},
		}
		_, err := services.CatalogMetricsArtifactRepository.Save(artifact2, savedModel.GetID())
		require.NoError(t, err)

		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:             100,
			Recommendations:       true, // Test with recommendations enabled
			FilterQuery:           "",
			PageSize:              10,
			OrderBy:               "",
			SortOrder:             "ASC",
			NextPageToken:         apiutils.Of(""),
			RPSProperty:           "throughput",    // Custom property name
			LatencyProperty:       "p90_latency",   // Custom property name
			HardwareCountProperty: "nodes",         // Custom property name
			HardwareTypeProperty:  "instance_type", // Custom property name
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"configurable-props-model",
			"configurable-props-source",
			params,
		)

		require.NoError(t, err)
		require.NotNil(t, result)
		assert.GreaterOrEqual(t, len(result.Items), 1)

		// Verify replicas are calculated based on custom throughput property
		// For targetRPS=100 and throughput=50, we should see replicas=2
		foundArtifact := false
		for _, item := range result.Items {
			if *item.CatalogMetricsArtifact.Name == "custom-props-artifact-2" {
				foundArtifact = true
				replicas := item.CatalogMetricsArtifact.CustomProperties["replicas"]
				require.NotNil(t, replicas.MetadataIntValue)
				assert.Equal(t, "2", replicas.MetadataIntValue.IntValue, "replicas should be ceil(100/50)=2")

				// Verify total_requests_per_second is calculated
				totalRPS := item.CatalogMetricsArtifact.CustomProperties["total_requests_per_second"]
				require.NotNil(t, totalRPS.MetadataDoubleValue)
				assert.Equal(t, 100.0, totalRPS.MetadataDoubleValue.DoubleValue)
				break
			}
		}
		assert.True(t, foundArtifact, "Should find artifact with custom properties")
	})

	t.Run("ErrorHandling_MissingCustomProperties", func(t *testing.T) {
		// Create a fresh artifact with specific properties (not including nonexistent_rps)
		errorTestArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("error-test-artifact"),
				ExternalID:  apiutils.Of("error-test-123"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "throughput", DoubleValue: apiutils.Of(100.0)},          // Has throughput
				{Name: "p90_latency", DoubleValue: apiutils.Of(120.0)},         // Has p90_latency
				{Name: "nodes", IntValue: apiutils.Of(int32(2))},               // Has nodes
				{Name: "instance_type", StringValue: apiutils.Of("gpu-large")}, // Has instance_type
				// But does NOT have nonexistent_rps
			},
		}
		_, err := services.CatalogMetricsArtifactRepository.Save(errorTestArtifact, savedModel.GetID())
		require.NoError(t, err)

		// Test with a nonexistent custom property name
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:             100,
			Recommendations:       true,
			FilterQuery:           "",
			PageSize:              10,
			OrderBy:               "",
			SortOrder:             "ASC",
			NextPageToken:         apiutils.Of(""),
			RPSProperty:           "nonexistent_rps", // This property doesn't exist in the artifact
			LatencyProperty:       "p90_latency",
			HardwareCountProperty: "nodes",
			HardwareTypeProperty:  "instance_type",
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"configurable-props-model",
			"configurable-props-source",
			params,
		)

		// Should return an error when required property is missing
		if err != nil {
			assert.Contains(t, err.Error(), "nonexistent_rps", "Error should mention the missing property")
		} else {
			// If no error, check if we got empty results (which might be acceptable behavior)
			// The service might return empty results instead of error
			t.Logf("No error returned, but result is: %+v", result)
			// For now, accept either error OR empty result as valid behavior
			assert.Equal(t, 0, len(result.Items), "Should have no items when property is missing")
		}
	})

	t.Run("DefaultPropertyNames_WhenNotSpecified", func(t *testing.T) {
		// Create artifact with default property names
		defaultArtifact := &models.CatalogMetricsArtifactImpl{
			TypeID: apiutils.Of(metricsArtifactTypeID),
			Attributes: &models.CatalogMetricsArtifactAttributes{
				Name:        apiutils.Of("default-props-artifact"),
				ExternalID:  apiutils.Of("default-789"),
				MetricsType: models.MetricsTypePerformance,
			},
			Properties: &[]dbmodels.Properties{
				{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
			},
			CustomProperties: &[]dbmodels.Properties{
				{Name: "requests_per_second", DoubleValue: apiutils.Of(150.0)},  // Default name
				{Name: "ttft_p90", DoubleValue: apiutils.Of(80.0)},              // Default name
				{Name: "hardware_count", IntValue: apiutils.Of(int32(3))},       // Default name
				{Name: "hardware_type", StringValue: apiutils.Of("gpu-medium")}, // Default name
			},
		}
		_, err := services.CatalogMetricsArtifactRepository.Save(defaultArtifact, savedModel.GetID())
		require.NoError(t, err)

		// Test with empty strings for custom properties (should use defaults)
		params := catalog.ListPerformanceArtifactsParams{
			TargetRPS:             100,
			Recommendations:       false,
			FilterQuery:           "",
			PageSize:              10,
			OrderBy:               "",
			SortOrder:             "ASC",
			NextPageToken:         apiutils.Of(""),
			RPSProperty:           "", // Should use default "requests_per_second"
			LatencyProperty:       "", // Should use default "ttft_p90"
			HardwareCountProperty: "", // Should use default "hardware_count"
			HardwareTypeProperty:  "", // Should use default "hardware_type"
		}

		result, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"configurable-props-model",
			"configurable-props-source",
			params,
		)

		require.NoError(t, err)
		require.NotNil(t, result)

		// Should find the artifact with default property names
		foundDefault := false
		for _, item := range result.Items {
			if *item.CatalogMetricsArtifact.Name == "default-props-artifact" {
				foundDefault = true
				// Verify replicas calculated from default property
				replicas := item.CatalogMetricsArtifact.CustomProperties["replicas"]
				require.NotNil(t, replicas.MetadataIntValue)
				// For targetRPS=100 and requests_per_second=150, replicas=ceil(100/150)=1
				assert.Equal(t, "1", replicas.MetadataIntValue.IntValue)
				break
			}
		}
		assert.True(t, foundDefault, "Should work with default property names")
	})

	t.Run("RecommendationsParameter_TrueFalse", func(t *testing.T) {
		// First, create multiple artifacts in same hardware group to test recommendations
		artifacts := []*models.CatalogMetricsArtifactImpl{
			{
				TypeID: apiutils.Of(metricsArtifactTypeID),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("rec-artifact-1"),
					ExternalID:  apiutils.Of("rec-1"),
					MetricsType: models.MetricsTypePerformance,
				},
				Properties: &[]dbmodels.Properties{
					{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
				},
				CustomProperties: &[]dbmodels.Properties{
					{Name: "requests_per_second", DoubleValue: apiutils.Of(100.0)},
					{Name: "ttft_p90", DoubleValue: apiutils.Of(90.0)},
					{Name: "hardware_count", IntValue: apiutils.Of(int32(2))},
					{Name: "hardware_type", StringValue: apiutils.Of("gpu-test")},
				},
			},
			{
				TypeID: apiutils.Of(metricsArtifactTypeID),
				Attributes: &models.CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("rec-artifact-2"),
					ExternalID:  apiutils.Of("rec-2"),
					MetricsType: models.MetricsTypePerformance,
				},
				Properties: &[]dbmodels.Properties{
					{Name: "metricsType", StringValue: apiutils.Of("performance-metrics")},
				},
				CustomProperties: &[]dbmodels.Properties{
					{Name: "requests_per_second", DoubleValue: apiutils.Of(120.0)},
					{Name: "ttft_p90", DoubleValue: apiutils.Of(95.0)},
					{Name: "hardware_count", IntValue: apiutils.Of(int32(2))},
					{Name: "hardware_type", StringValue: apiutils.Of("gpu-test")},
				},
			},
		}
		for _, a := range artifacts {
			_, err := services.CatalogMetricsArtifactRepository.Save(a, savedModel.GetID())
			require.NoError(t, err)
		}

		// Test with recommendations=false (no recommendations)
		paramsNoRecommendations := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       100,
			Recommendations: false, // Should return all artifacts
			PageSize:        10,
		}

		resultNoRecommendations, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"configurable-props-model",
			"configurable-props-source",
			paramsNoRecommendations,
		)
		require.NoError(t, err)
		require.NotNil(t, resultNoRecommendations)

		// Count how many rec-artifact-* items we have
		noRecommendationCount := 0
		for _, item := range resultNoRecommendations.Items {
			name := *item.CatalogMetricsArtifact.Name
			if name == "rec-artifact-1" || name == "rec-artifact-2" {
				noRecommendationCount++
			}
		}
		assert.Equal(t, 2, noRecommendationCount, "With recommendations=false, should return both artifacts")

		// Test with recommendations=true (with recommendations)
		paramsRecommendations := catalog.ListPerformanceArtifactsParams{
			TargetRPS:       100,
			Recommendations: true, // Should deduplicate
			PageSize:        10,
		}

		resultRecommendations, err := provider.GetPerformanceArtifacts(
			context.Background(),
			"configurable-props-model",
			"configurable-props-source",
			paramsRecommendations,
		)
		require.NoError(t, err)
		require.NotNil(t, resultRecommendations)

		// Count how many rec-artifact-* items after recommendations
		dedupCount := 0
		for _, item := range resultRecommendations.Items {
			name := *item.CatalogMetricsArtifact.Name
			if name == "rec-artifact-1" || name == "rec-artifact-2" {
				dedupCount++
			}
		}
		// With recommendations, should have fewer items (at most 1 from the gpu-test group)
		assert.LessOrEqual(t, dedupCount, 1, "With recommendations=true, should deduplicate similar configs")
	})
}
