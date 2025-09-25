package openapi

import (
	"context"
	"errors"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// ModelCatalogServiceAPIService is a service that implements the logic for the ModelCatalogServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelCatalogServiceAPI s.coreApi.
// Include any external packages or services that will be required by this service.
type ModelCatalogServiceAPIService struct {
	provider catalog.CatalogSourceProvider
}

// GetAllModelArtifacts retrieves all model artifacts for a given model from the specified source.
func (m *ModelCatalogServiceAPIService) GetAllModelArtifacts(ctx context.Context, sourceID string, name string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	if newName, err := url.PathUnescape(name); err == nil {
		name = newName
	}

	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return Response(http.StatusBadRequest, err), err
	}

	artifacts, err := m.provider.GetArtifacts(ctx, name, catalog.ListArtifactsParams{
		PageSize:      int32(pageSizeInt),
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	})
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, artifacts), nil
}

func (m *ModelCatalogServiceAPIService) FindModels(ctx context.Context, sourceIDs []string, q string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	pageSizeInt, err := strconv.ParseInt(pageSize, 10, 32)
	if err != nil {
		return Response(http.StatusBadRequest, err), err
	}

	listModelsParams := catalog.ListModelsParams{
		Query:         q,
		SourceIDs:     sourceIDs,
		PageSize:      int32(pageSizeInt),
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	}

	models, err := m.provider.ListModels(ctx, listModelsParams)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, models), nil
}

func (m *ModelCatalogServiceAPIService) GetModel(ctx context.Context, sourceID string, name string) (ImplResponse, error) {
	if name, ok := strings.CutSuffix(name, "/artifacts"); ok {
		return m.GetAllModelArtifacts(ctx, sourceID, name, "10", model.OrderByField(model.ORDERBYFIELD_CREATE_TIME), model.SortOrder(model.SORTORDER_ASC), "")
	}

	if newName, err := url.PathUnescape(name); err == nil {
		name = newName
	}

	model, err := m.provider.GetModel(ctx, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}
	if model == nil {
		return notFound("Unknown model or version"), nil
	}

	return Response(http.StatusOK, model), nil
}

func (m *ModelCatalogServiceAPIService) FindSources(ctx context.Context, name string, strPageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	// sources := m.sources.All()
	// if len(sources) > math.MaxInt32 {
	// 	err := errors.New("too many registered models")
	// 	return ErrorResponse(http.StatusInternalServerError, err), err
	// }

	// paginator, err := newPaginator[model.CatalogSource](strPageSize, orderBy, sortOrder, nextPageToken)
	// if err != nil {
	// 	return ErrorResponse(http.StatusBadRequest, err), err
	// }

	// items := make([]model.CatalogSource, 0, len(sources))

	// name = strings.ToLower(name)

	// for _, v := range sources {
	// 	if !strings.Contains(strings.ToLower(v.Metadata.Name), name) {
	// 		continue
	// 	}

	// 	items = append(items, v.Metadata)
	// }

	// cmpFunc, err := genCatalogCmpFunc(orderBy, sortOrder)
	// if err != nil {
	// 	return ErrorResponse(http.StatusBadRequest, err), err
	// }
	// slices.SortStableFunc(items, cmpFunc)

	// total := int32(len(items))

	// pagedItems, next := paginator.Paginate(items)

	// res := model.CatalogSourceList{
	// 	PageSize:      paginator.PageSize,
	// 	Items:         pagedItems,
	// 	Size:          total,
	// 	NextPageToken: next.Token(),
	// }
	return Response(http.StatusOK, model.CatalogSourceList{}), nil
}

var _ ModelCatalogServiceAPIServicer = &ModelCatalogServiceAPIService{}

// NewModelCatalogServiceAPIService creates a default api service
func NewModelCatalogServiceAPIService(provider catalog.CatalogSourceProvider) ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{
		provider: provider,
	}
}

func notFound(msg string) ImplResponse {
	if msg == "" {
		msg = "Resource not found"
	}
	return ErrorResponse(http.StatusNotFound, errors.New(msg))
}
