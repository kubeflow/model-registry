package catalog

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
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

func TestRhecCatalogListModels(t *testing.T) {
	modelTime := time.Now()
	createTime := modelTime.Format(time.RFC3339)
	lastUpdateTime := modelTime.Add(5 * time.Minute).Format(time.RFC3339)
	sourceId := "rhec"
	provider := "redhat"
	artifactCreateTime := modelTime.Add(10 * time.Minute).Format(time.RFC3339)
	artifactLastUpdateTime := modelTime.Add(15 * time.Minute).Format(time.RFC3339)

	testModels := map[string]*rhecModel{
		"model3": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model3",
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
		"model1:v2": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model1:v2",
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
		"model2": {
			CatalogModel: openapi.CatalogModel{
				Name:                     "model2",
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
		params    ListModelsParams
		want      openapi.CatalogModelList
		wantErr   bool
	}{
		{
			name: "list models and sort order",
			want: openapi.CatalogModelList{
				Items: []openapi.CatalogModel{
					{
						Name:                     "model1",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
					{
						Name:                     "model1:v2",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
					{
						Name:                     "model2",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
					{
						Name:                     "model3",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
				},
				PageSize: 4,
				Size:     4,
			},
			wantErr: false,
		},
		{
			name:      "list models with query and sort order",
			modelName: "model1",
			params: ListModelsParams{
				Query:     "model1",
				OrderBy:   model.ORDERBYFIELD_NAME,
				SortOrder: model.SORTORDER_ASC,
			},
			want: openapi.CatalogModelList{
				Items: []openapi.CatalogModel{
					{
						Name:                     "model1",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
					{
						Name:                     "model1:v2",
						CreateTimeSinceEpoch:     &createTime,
						LastUpdateTimeSinceEpoch: &lastUpdateTime,
						Provider:                 &provider,
						SourceId:                 &sourceId,
					},
				},
				PageSize: 2,
				Size:     2,
			},
			wantErr: false,
		},
		{
			name:      "get non-existent model",
			modelName: "not-exist",
			params: ListModelsParams{
				Query:     "not-exist",
				OrderBy:   model.ORDERBYFIELD_NAME,
				SortOrder: model.SORTORDER_ASC,
			},
			want:    openapi.CatalogModelList{Items: []openapi.CatalogModel{}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := r.ListModels(context.Background(), tt.params)
			if (err != nil) != tt.wantErr {
				t.Errorf("ListModels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("ListModels() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestIsModelExcluded(t *testing.T) {
	tests := []struct {
		name      string
		modelName string
		patterns  []string
		want      bool
	}{
		{
			name:      "exact match",
			modelName: "model1:v1",
			patterns:  []string{"model1:v1"},
			want:      true,
		},
		{
			name:      "wildcard match",
			modelName: "model1:v2",
			patterns:  []string{"model1:*"},
			want:      true,
		},
		{
			name:      "no match",
			modelName: "model2:v1",
			patterns:  []string{"model1:*"},
			want:      false,
		},
		{
			name:      "multiple patterns with match",
			modelName: "model3:v1",
			patterns:  []string{"model2:*", "model3:v1"},
			want:      true,
		},
		{
			name:      "empty patterns",
			modelName: "model1:v1",
			patterns:  []string{},
			want:      false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isModelExcluded(tt.modelName, tt.patterns)
			if got != tt.want {
				t.Errorf("isModelExcluded() = %v, want %v", got, tt.want)
			}
		})
	}
}
