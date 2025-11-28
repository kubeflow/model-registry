package models

import (
	"testing"

	"github.com/kubeflow/model-registry/internal/apiutils"
	dbmodels "github.com/kubeflow/model-registry/internal/db/models"
	mrmodels "github.com/kubeflow/model-registry/internal/db/models"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// MockCatalogArtifactRepository is a mock implementation for testing
type MockCatalogArtifactRepository struct {
	mock.Mock
}

func (m *MockCatalogArtifactRepository) GetByID(id int32) (CatalogArtifact, error) {
	args := m.Called(id)
	return args.Get(0).(CatalogArtifact), args.Error(1)
}

func (m *MockCatalogArtifactRepository) List(listOptions CatalogArtifactListOptions) (*mrmodels.ListWrapper[CatalogArtifact], error) {
	args := m.Called(listOptions)
	return args.Get(0).(*mrmodels.ListWrapper[CatalogArtifact]), args.Error(1)
}

func (m *MockCatalogArtifactRepository) DeleteByParentID(artifactType string, parentResourceID int32) error {
	args := m.Called(artifactType, parentResourceID)
	return args.Error(0)
}

func TestPerformanceArtifactService_GetArtifacts(t *testing.T) {
	// Mock repository for testing
	mockRepo := &MockCatalogArtifactRepository{}
	service := NewPerformanceArtifactService(mockRepo)

	// Setup mock to return test artifacts with performance metrics
	id := int32(1)
	testArtifacts := []CatalogArtifact{
		{
			CatalogMetricsArtifact: &CatalogMetricsArtifactImpl{
				Attributes: &CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("test-perf-artifact"),
					MetricsType: MetricsTypePerformance,
				},
				Properties: &[]dbmodels.Properties{},
				CustomProperties: &[]dbmodels.Properties{
					{Name: "requests_per_second", DoubleValue: apiutils.Of(200.0)},
				},
			},
		},
	}
	// Set ID for the artifact
	testArtifacts[0].CatalogMetricsArtifact.SetID(id)

	mockRepo.On("List", mock.AnythingOfType("models.CatalogArtifactListOptions")).
		Return(&mrmodels.ListWrapper[CatalogArtifact]{
			Items: testArtifacts,
		}, nil)

	params := PerformanceArtifactParams{
		ModelID:         123,
		TargetRPS:       300,
		Recommendations: false,
		PageSize:        10,
	}

	result, err := service.GetArtifacts(params)

	require.NoError(t, err)
	require.Len(t, result.Items, 1)

	// Verify repository was called with correct performance filtering
	mockRepo.AssertCalled(t, "List", mock.MatchedBy(func(opts CatalogArtifactListOptions) bool {
		return opts.ParentResourceID != nil && *opts.ParentResourceID == 123
	}))
}

func TestPerformanceArtifactService_ProcessWithTargetRPSAndRecommendataion(t *testing.T) {
	mockRepo := &MockCatalogArtifactRepository{}
	service := NewPerformanceArtifactService(mockRepo)

	// Mock repository to return test artifacts with performance data
	id := int32(1)
	testDBMetrics := &CatalogMetricsArtifactImpl{
		Attributes: &CatalogMetricsArtifactAttributes{
			Name:        apiutils.Of("test-perf-artifact"),
			MetricsType: MetricsTypePerformance,
		},
		Properties: &[]dbmodels.Properties{},
		CustomProperties: &[]dbmodels.Properties{
			{Name: "requests_per_second", DoubleValue: apiutils.Of(200.0)},
			{Name: "ttft_p90", DoubleValue: apiutils.Of(50.0)},
			{Name: "estimated_cost_per_hour", DoubleValue: apiutils.Of(80.0)},
			{Name: "hardware_type", StringValue: apiutils.Of("gpu-a100")},
		},
	}
	testDBMetrics.SetID(id)

	testDBResult := &mrmodels.ListWrapper[CatalogArtifact]{
		Items: []CatalogArtifact{
			{CatalogMetricsArtifact: testDBMetrics},
		},
	}

	mockRepo.On("List", mock.Anything).Return(testDBResult, nil)

	params := PerformanceArtifactParams{
		ModelID:         123,
		TargetRPS:       600, // Should result in 3 replicas (600/200)
		Recommendations: true,
		PageSize:        10,
	}

	result, err := service.GetArtifacts(params)

	require.NoError(t, err)
	require.Len(t, result.Items, 1)

	artifact := result.Items[0]

	// Check targetRPS calculations were applied
	customProps := artifact.GetCustomProperties()
	require.NotNil(t, customProps)

	// Helper function to find property by name
	findProperty := func(name string) *dbmodels.Properties {
		for _, prop := range *customProps {
			if prop.Name == name {
				return &prop
			}
		}
		return nil
	}

	replicasProp := findProperty("replicas")
	require.NotNil(t, replicasProp)
	require.NotNil(t, replicasProp.IntValue)
	require.Equal(t, int32(3), *replicasProp.IntValue)

	totalRPSProp := findProperty("total_requests_per_second")
	require.NotNil(t, totalRPSProp)
	require.NotNil(t, totalRPSProp.DoubleValue)
	require.Equal(t, 600.0, *totalRPSProp.DoubleValue)
}

func TestPerformanceArtifactParamsStructure(t *testing.T) {
	params := PerformanceArtifactParams{
		ModelID:               1,
		TargetRPS:             100,
		Recommendations:       true,
		RPSProperty:           "throughput",
		LatencyProperty:       "p90_latency",
		HardwareCountProperty: "nodes",
		HardwareTypeProperty:  "instance_type",
	}

	require.Equal(t, "throughput", params.RPSProperty)
	require.Equal(t, "p90_latency", params.LatencyProperty)
	require.Equal(t, "nodes", params.HardwareCountProperty)
	require.Equal(t, "instance_type", params.HardwareTypeProperty)
}

func TestPerformanceArtifactService_Recommendataions(t *testing.T) {
	service := NewPerformanceArtifactService(nil) // No repo needed for this test

	id1 := int32(1)
	id2 := int32(2)
	artifacts := []CatalogMetricsArtifact{
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("fast-expensive"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewStringProperty("hardware_type", "gpu-a100", true),
				dbmodels.NewDoubleProperty("ttft_p90", 50.0, true),
				dbmodels.NewIntProperty("hardware_count", 2, true),
				dbmodels.NewIntProperty("replicas", 3, true),
			},
		},
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("slow-cheap"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewStringProperty("hardware_type", "gpu-a100", true),
				dbmodels.NewDoubleProperty("ttft_p90", 150.0, true),
				dbmodels.NewIntProperty("hardware_count", 1, true),
				dbmodels.NewIntProperty("replicas", 2, true),
			},
		},
	}

	// Set IDs for the artifacts
	artifacts[0].SetID(id1)
	artifacts[1].SetID(id2)

	result := service.generateRecommendations(artifacts, "ttft_p90", "hardware_count", "hardware_type")

	// Should apply the two-pass filtering with epsilon thresholds
	require.True(t, len(result) <= len(artifacts))
}

// TestPropertyValidation tests the validateCustomProperties method
func TestPropertyValidation(t *testing.T) {
	service := NewPerformanceArtifactService(nil)

	// Test case: property exists in some artifacts
	id1 := int32(1)
	id2 := int32(2)
	artifacts := []CatalogMetricsArtifact{
		// artifact with custom property
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("artifact-with-prop"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewDoubleProperty("custom_rps", 100.0, true),
				dbmodels.NewDoubleProperty("custom_latency", 50.0, true),
				dbmodels.NewIntProperty("custom_hw_count", 2, true),
				dbmodels.NewStringProperty("custom_hw_type", "gpu", true),
			},
		},
		// artifact without custom property
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("artifact-without-prop"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewDoubleProperty("other_prop", 200.0, true),
			},
		},
	}
	artifacts[0].SetID(id1)
	artifacts[1].SetID(id2)

	// Should not error when property exists in at least one artifact
	err := service.validateCustomProperties(artifacts, "custom_rps", "custom_latency", "custom_hw_count", "custom_hw_type")
	require.NoError(t, err)

	// Test case: property doesn't exist in any artifact
	err = service.validateCustomProperties(artifacts, "nonexistent_prop", "", "", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "nonexistent_prop")
	require.Contains(t, err.Error(), "not found")

	// Test case: multiple properties missing
	err = service.validateCustomProperties(artifacts, "missing1", "missing2", "", "")
	require.Error(t, err)
	require.Contains(t, err.Error(), "not found")

	// Test case: empty property names should be ignored
	err = service.validateCustomProperties(artifacts, "", "", "", "")
	require.NoError(t, err)
}

// TestGetArtifactsWithValidation tests that GetArtifacts validates custom properties
func TestGetArtifactsWithValidation(t *testing.T) {
	mockRepo := &MockCatalogArtifactRepository{}
	service := NewPerformanceArtifactService(mockRepo)

	id1 := int32(1)
	testArtifacts := []CatalogArtifact{
		{
			CatalogMetricsArtifact: &CatalogMetricsArtifactImpl{
				Attributes: &CatalogMetricsArtifactAttributes{
					Name:        apiutils.Of("test-artifact"),
					MetricsType: MetricsTypePerformance,
				},
				CustomProperties: &[]dbmodels.Properties{
					dbmodels.NewDoubleProperty("custom_rps", 100.0, true),
					dbmodels.NewDoubleProperty("custom_latency", 50.0, true),
				},
			},
		},
	}
	testArtifacts[0].CatalogMetricsArtifact.SetID(id1)

	mockRepo.On("List", mock.AnythingOfType("models.CatalogArtifactListOptions")).
		Return(&mrmodels.ListWrapper[CatalogArtifact]{Items: testArtifacts}, nil)

	// Test with valid custom properties
	params := PerformanceArtifactParams{
		ModelID:         123,
		TargetRPS:       100,
		Recommendations: false,
		PageSize:        10,
		RPSProperty:     "custom_rps",
		LatencyProperty: "custom_latency",
	}

	result, err := service.GetArtifacts(params)
	require.NoError(t, err)
	require.NotNil(t, result)

	// Test with missing custom property
	params.RPSProperty = "nonexistent_property"
	result, err = service.GetArtifacts(params)
	require.Error(t, err)
	require.Nil(t, result)
	require.Contains(t, err.Error(), "invalid custom properties")
	require.Contains(t, err.Error(), "nonexistent_property")
}

// TestConfigurablePropertyUsage tests that custom property names are used in calculations
func TestConfigurablePropertyUsage(t *testing.T) {
	service := NewPerformanceArtifactService(nil) // No repo needed for this test

	// Create test artifact with custom property names
	id := int32(1)
	artifact := &CatalogMetricsArtifactImpl{
		Attributes: &CatalogMetricsArtifactAttributes{
			Name: apiutils.Of("test-artifact"),
		},
		CustomProperties: &[]dbmodels.Properties{
			dbmodels.NewDoubleProperty("throughput", 50.0, true),
			dbmodels.NewDoubleProperty("p90_latency", 120.0, true),
			dbmodels.NewIntProperty("nodes", 2, true),
			dbmodels.NewStringProperty("instance_type", "gpu-large", true),
		},
	}
	artifact.SetID(id)

	artifacts := []CatalogMetricsArtifact{artifact}

	params := PerformanceArtifactParams{
		TargetRPS:             100,
		Recommendations:       true,
		RPSProperty:           "throughput",
		LatencyProperty:       "p90_latency",
		HardwareCountProperty: "nodes",
		HardwareTypeProperty:  "instance_type",
	}

	result := service.processArtifacts(artifacts, params)

	// Should calculate replicas based on custom throughput property
	// 100 targetRPS / 50 throughput = 2 replicas
	replicas := service.extractCustomPropertiesIntValue(result[0].GetCustomProperties(), "replicas", 0)
	require.Equal(t, int32(2), replicas)

	// Should calculate total RPS correctly
	totalRPS := service.extractCustomPropertiesDoubleValue(result[0].GetCustomProperties(), "total_requests_per_second", 0)
	require.Equal(t, 100.0, totalRPS)
}

// TestConfigurableRecommendataion tests deduplication with configurable property names
func TestConfigurableRecommendataion(t *testing.T) {
	service := NewPerformanceArtifactService(nil)

	id1 := int32(1)
	id2 := int32(2)
	id3 := int32(3)
	artifacts := []CatalogMetricsArtifact{
		// Dominated artifact: slower AND more expensive
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("artifact-slow-expensive"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewStringProperty("instance_type", "gpu-small", true),
				dbmodels.NewDoubleProperty("p90_latency", 150.0, true),
				dbmodels.NewIntProperty("nodes", 3, true),
				dbmodels.NewIntProperty("replicas", 1, true),
			},
		},
		// Fast but expensive
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("artifact-fast-expensive"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewStringProperty("instance_type", "gpu-small", true),
				dbmodels.NewDoubleProperty("p90_latency", 100.0, true),
				dbmodels.NewIntProperty("nodes", 2, true),
				dbmodels.NewIntProperty("replicas", 1, true),
			},
		},
		// Cheapest option (but slower than artifact-fast-expensive)
		&CatalogMetricsArtifactImpl{
			Attributes: &CatalogMetricsArtifactAttributes{
				Name: apiutils.Of("artifact-cheap"),
			},
			CustomProperties: &[]dbmodels.Properties{
				dbmodels.NewStringProperty("instance_type", "gpu-small", true),
				dbmodels.NewDoubleProperty("p90_latency", 120.0, true),
				dbmodels.NewIntProperty("nodes", 1, true),
				dbmodels.NewIntProperty("replicas", 1, true),
			},
		},
	}

	artifacts[0].SetID(id1)
	artifacts[1].SetID(id2)
	artifacts[2].SetID(id3)

	result := service.generateRecommendations(artifacts, "p90_latency", "nodes", "instance_type")

	// The dominated artifact (slow and expensive) should be filtered out
	require.Less(t, len(result), len(artifacts))

	// Verify that the function uses configurable property names by checking the results
	// The algorithm should filter based on the custom property names provided
	for _, artifact := range result {
		latency := service.extractCustomPropertiesDoubleValue(artifact.GetCustomProperties(), "p90_latency", 0)
		nodes := service.extractCustomPropertiesIntValue(artifact.GetCustomProperties(), "nodes", 0)
		instanceType := service.extractCustomPropertiesStringValue(artifact.GetCustomProperties(), "instance_type")

		// Verify we can extract values using custom property names
		require.Greater(t, latency, 0.0)
		require.Greater(t, nodes, int32(0))
		require.Equal(t, "gpu-small", instanceType)
	}
}
