package openapi

import (
	"context"
	"net/http"
	"testing"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// mockCatalogProvider implements catalog.APIProvider for testing
type mockCatalogProvider struct{}

func (m *mockCatalogProvider) GetModel(ctx context.Context, modelName string, sourceID string) (*openapi.CatalogModel, error) {
	return nil, nil
}

func (m *mockCatalogProvider) ListModels(ctx context.Context, params catalog.ListModelsParams) (openapi.CatalogModelList, error) {
	// Create test models
	model1ID := "1"
	model1Name := "test-model-1"
	model2ID := "2"
	model2Name := "test-model-2"

	model1 := openapi.CatalogModel{
		Id:   &model1ID,
		Name: model1Name,
	}

	model2 := openapi.CatalogModel{
		Id:   &model2ID,
		Name: model2Name,
	}

	// If artifact types are requested, add artifacts
	if len(params.ArtifactTypesFilter) > 0 {
		for _, artifactType := range params.ArtifactTypesFilter {
			if artifactType == "metrics-artifact" {
				// Add metrics artifact to model1
				metricsName := "performance-metrics"
				metricsType := "performance-metric"
				metricsArtifact := openapi.CatalogMetricsArtifact{
					Name:         &metricsName,
					ArtifactType: "metrics-artifact",
					MetricsType:  metricsType,
				}
				model1.Artifacts = append(model1.Artifacts, openapi.CatalogArtifact{
					CatalogMetricsArtifact: &metricsArtifact,
				})
			}
			if artifactType == "model-artifact" {
				// Add model artifact to model1
				modelName := "model-file"
				modelArtifact := openapi.CatalogModelArtifact{
					Name:         &modelName,
					ArtifactType: "model-artifact",
					Uri:          "oci://test",
				}
				model1.Artifacts = append(model1.Artifacts, openapi.CatalogArtifact{
					CatalogModelArtifact: &modelArtifact,
				})
			}
		}

		// Model 2 has no artifacts, so set empty array
		model2.Artifacts = []openapi.CatalogArtifact{}
	}

	return openapi.CatalogModelList{
		Items:    []openapi.CatalogModel{model1, model2},
		PageSize: 10,
		Size:     2,
	}, nil
}

func (m *mockCatalogProvider) GetArtifacts(ctx context.Context, modelName string, sourceID string, params catalog.ListArtifactsParams) (openapi.CatalogArtifactList, error) {
	return openapi.CatalogArtifactList{}, nil
}

func (m *mockCatalogProvider) GetFilterOptions(ctx context.Context) (*openapi.FilterOptionsList, error) {
	return nil, nil
}

// TestArtifactInclusion tests the artifact inclusion feature
func TestArtifactInclusion(t *testing.T) {
	provider := &mockCatalogProvider{}
	sources := &catalog.SourceCollection{}
	service := NewModelCatalogServiceAPIService(provider, sources)

	t.Run("without artifactType parameter - artifacts should NOT be included", func(t *testing.T) {
		resp, err := service.FindModels(
			context.Background(),
			[]string{}, // sourceIDs
			"",         // q
			[]string{}, // sourceLabels
			"",         // filterQuery
			"10",       // pageSize
			"",         // orderBy
			"",         // sortOrder
			"",         // nextPageToken
			[]string{}, // artifactTypes - EMPTY
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		modelList, ok := resp.Body.(openapi.CatalogModelList)
		require.True(t, ok, "Response should be CatalogModelList")
		require.Len(t, modelList.Items, 2)

		// Artifacts should NOT be present (nil or empty)
		assert.Nil(t, modelList.Items[0].Artifacts, "artifacts should be nil when not requested")
		assert.Nil(t, modelList.Items[1].Artifacts, "artifacts should be nil when not requested")
	})

	t.Run("with artifactType=metrics-artifact - artifacts SHOULD be included", func(t *testing.T) {
		resp, err := service.FindModels(
			context.Background(),
			[]string{},
			"",
			[]string{},
			"",
			"10",
			"",
			"",
			"",
			[]string{"metrics-artifact"}, // Request metrics artifacts
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		modelList, ok := resp.Body.(openapi.CatalogModelList)
		require.True(t, ok)
		require.Len(t, modelList.Items, 2)

		// Model 1 should have 1 metrics artifact
		require.NotNil(t, modelList.Items[0].Artifacts, "artifacts should be present")
		assert.Len(t, modelList.Items[0].Artifacts, 1, "model 1 should have 1 artifact")
		assert.NotNil(t, modelList.Items[0].Artifacts[0].CatalogMetricsArtifact)
		assert.Equal(t, "performance-metrics", *modelList.Items[0].Artifacts[0].CatalogMetricsArtifact.Name)

		// Model 2 should have empty artifacts array
		require.NotNil(t, modelList.Items[1].Artifacts, "artifacts should be present even if empty")
		assert.Len(t, modelList.Items[1].Artifacts, 0, "model 2 should have 0 artifacts")
	})

	t.Run("with multiple artifactType parameters", func(t *testing.T) {
		resp, err := service.FindModels(
			context.Background(),
			[]string{},
			"",
			[]string{},
			"",
			"10",
			"",
			"",
			"",
			[]string{"metrics-artifact", "model-artifact"}, // Request both types
		)

		require.NoError(t, err)
		assert.Equal(t, http.StatusOK, resp.Code)

		modelList, ok := resp.Body.(openapi.CatalogModelList)
		require.True(t, ok)
		require.Len(t, modelList.Items, 2)

		// Model 1 should have 2 artifacts (1 metrics + 1 model)
		require.NotNil(t, modelList.Items[0].Artifacts)
		assert.Len(t, modelList.Items[0].Artifacts, 2, "model 1 should have 2 artifacts")

		// Verify both artifact types are present
		hasMetrics := false
		hasModel := false
		for _, artifact := range modelList.Items[0].Artifacts {
			if artifact.CatalogMetricsArtifact != nil {
				hasMetrics = true
			}
			if artifact.CatalogModelArtifact != nil {
				hasModel = true
			}
		}
		assert.True(t, hasMetrics, "should have metrics artifact")
		assert.True(t, hasModel, "should have model artifact")
	})
}
