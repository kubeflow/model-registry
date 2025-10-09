package catalog

import (
	"context"

	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

type ListModelsParams struct {
	Query         string
	FilterQuery   string
	SourceIDs     []string
	PageSize      int32
	OrderBy       model.OrderByField
	SortOrder     model.SortOrder
	NextPageToken *string
}

type ListArtifactsParams struct {
	PageSize      int32
	OrderBy       model.OrderByField
	SortOrder     model.SortOrder
	NextPageToken *string
	ArtifactType  *string
}

// APIProvider implements the API endpoints.
type APIProvider interface {
	// GetModel returns model metadata for a single model by its name. If
	// nothing is found with the name provided it returns nil, without an
	// error.
	GetModel(ctx context.Context, modelName string, sourceID string) (*model.CatalogModel, error)

	// ListModels returns all models according to the parameters. If
	// nothing suitable is found, it returns an empty list.
	// If sourceIDs is provided, filter models by source IDs. If not provided, return all models.
	ListModels(ctx context.Context, params ListModelsParams) (model.CatalogModelList, error)

	// GetArtifacts returns all artifacts for a particular model. If no
	// model is found with that name, it returns nil. If the model is
	// found, but has no artifacts, an empty list is returned.
	GetArtifacts(ctx context.Context, modelName string, sourceID string, params ListArtifactsParams) (model.CatalogArtifactList, error)
}
