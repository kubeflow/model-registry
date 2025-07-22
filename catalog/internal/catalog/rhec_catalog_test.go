package catalog

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

func TestRhecCatalogImpl_GetModel(t *testing.T) {
	modelTime := time.Now()
	createTime := modelTime.Format(time.RFC3339)
	lastUpdateTime := modelTime.Add(5 * time.Minute).Format(time.RFC3339)
	sourceId := "rhec"
	provider := "redhat"

	testModels := map[string]*rhecModel{
		"model1": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model1",
				CreateTimeSinceEpoch:     &createTime,
				LastUpdateTimeSinceEpoch: &lastUpdateTime,
				Provider:                 &provider,
				SourceId:                 &sourceId,
			},
			Artifacts: []*openapi.CatalogModelArtifact{},
		},
	}

	r := &rhecCatalogImpl{
		models: testModels,
	}

	tests := []struct {
		name      string
		modelName string
		want      *openapi.CatalogModel
		wantErr   bool
	}{
		{
			name:      "get existing model",
			modelName: "model1",
			want: &openapi.CatalogModel{
				Name:                     "model1",
				CreateTimeSinceEpoch:     &createTime,
				LastUpdateTimeSinceEpoch: &lastUpdateTime,
				Provider:                 &provider,
				SourceId:                 &sourceId,
			},
			wantErr: false,
		},
		{
			name:      "get non-existent model",
			modelName: "not-exist",
			want:      nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetModel(context.Background(), tt.modelName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetModel() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetModel() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestRhecCatalogImpl_GetArtifacts(t *testing.T) {
	modelTime := time.Now()
	createTime := modelTime.Format(time.RFC3339)
	lastUpdateTime := modelTime.Add(5 * time.Minute).Format(time.RFC3339)
	sourceId := "rhec"
	provider := "redhat"
	artifactCreateTime := modelTime.Add(10 * time.Minute).Format(time.RFC3339)
	artifactLastUpdateTime := modelTime.Add(15 * time.Minute).Format(time.RFC3339)

	testModels := map[string]*rhecModel{
		"model1": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model1",
				CreateTimeSinceEpoch:     &createTime,
				LastUpdateTimeSinceEpoch: &lastUpdateTime,
				Provider:                 &provider,
				SourceId:                 &sourceId,
			},
			Artifacts: []*openapi.CatalogModelArtifact{
				{
					Uri:                      "test-uri",
					CreateTimeSinceEpoch:     &artifactCreateTime,
					LastUpdateTimeSinceEpoch: &artifactLastUpdateTime,
				},
			},
		},
		"model2-no-artifacts": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model2-no-artifacts",
				CreateTimeSinceEpoch:     &createTime,
				LastUpdateTimeSinceEpoch: &lastUpdateTime,
				Provider:                 &provider,
				SourceId:                 &sourceId,
			},
			Artifacts: []*openapi.CatalogModelArtifact{},
		},
	}

	r := &rhecCatalogImpl{
		models: testModels,
	}

	tests := []struct {
		name      string
		modelName string
		want      *openapi.CatalogModelArtifactList
		wantErr   bool
	}{
		{
			name:      "get artifacts for existing model",
			modelName: "model1",
			want: &openapi.CatalogModelArtifactList{
				Items: []openapi.CatalogModelArtifact{
					{
						Uri:                      "test-uri",
						CreateTimeSinceEpoch:     &artifactCreateTime,
						LastUpdateTimeSinceEpoch: &artifactLastUpdateTime,
					},
				},
				PageSize: 1,
				Size:     1,
			},
			wantErr: false,
		},
		{
			name:      "get artifacts for model with no artifacts",
			modelName: "model2-no-artifacts",
			want: &openapi.CatalogModelArtifactList{
				Items:    []openapi.CatalogModelArtifact{},
				PageSize: 0,
				Size:     0,
			},
			wantErr: false,
		},
		{
			name:      "get artifacts for non-existent model",
			modelName: "not-exist",
			want:      nil,
			wantErr:   false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.GetArtifacts(context.Background(), tt.modelName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetArtifacts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("GetArtifacts() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
