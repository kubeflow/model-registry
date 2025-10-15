package openapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/pkg/api"
)

// ModelCatalogServiceAPIService is a service that implements the logic for the ModelCatalogServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelCatalogServiceAPI s.coreApi.
// Include any external packages or services that will be required by this service.
type ModelCatalogServiceAPIService struct {
	provider catalog.APIProvider
	sources  *catalog.SourceCollection
}

// GetAllModelArtifacts retrieves all model artifacts for a given model from the specified source.
func (m *ModelCatalogServiceAPIService) GetAllModelArtifacts(ctx context.Context, sourceID string, modelName string, artifactType string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	if newName, err := url.PathUnescape(modelName); err == nil {
		modelName = newName
	}

	var err error
	pageSizeInt := int32(10)

	if pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 32)
		if err != nil {
			return Response(http.StatusBadRequest, err), err
		}
		pageSizeInt = int32(parsed)
	}

	artifacts, err := m.provider.GetArtifacts(ctx, modelName, sourceID, catalog.ListArtifactsParams{
		ArtifactType:  &artifactType,
		PageSize:      pageSizeInt,
		OrderBy:       orderBy,
		SortOrder:     sortOrder,
		NextPageToken: &nextPageToken,
	})
	if err != nil {
		statusCode := api.ErrToStatus(err)
		var errorMsg string
		if errors.Is(err, api.ErrBadRequest) {
			errorMsg = fmt.Sprintf("Invalid model name '%s' for source '%s'", modelName, sourceID)
		} else if errors.Is(err, api.ErrNotFound) {
			errorMsg = fmt.Sprintf("No model found '%s' in source '%s'", modelName, sourceID)
		} else {
			errorMsg = err.Error()
		}
		return ErrorResponse(statusCode, errors.New(errorMsg)), err
	}

	return Response(http.StatusOK, artifacts), nil
}

func (m *ModelCatalogServiceAPIService) FindModels(ctx context.Context, sourceIDs []string, q string, sourceLabels []string, filterQuery string, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	var err error
	pageSizeInt := int32(10)

	if pageSize != "" {
		parsed, err := strconv.ParseInt(pageSize, 10, 32)
		if err != nil {
			return Response(http.StatusBadRequest, err), err
		}
		pageSizeInt = int32(parsed)
	}

	if len(sourceIDs) == 1 && sourceIDs[0] == "" {
		sourceIDs = nil
	}
	if len(sourceLabels) == 1 && sourceLabels[0] == "" {
		sourceLabels = nil
	}

	if len(sourceIDs) > 0 && len(sourceLabels) > 0 {
		err := fmt.Errorf("source and sourceLabel cannot be used together")
		return Response(http.StatusBadRequest, err), err
	}

	listModelsParams := catalog.ListModelsParams{
		Query:         q,
		FilterQuery:   filterQuery,
		SourceIDs:     sourceIDs,
		SourceLabels:  sourceLabels,
		PageSize:      pageSizeInt,
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

func (m *ModelCatalogServiceAPIService) FindModelsFilterOptions(ctx context.Context) (ImplResponse, error) {
	filterOptions, err := m.provider.GetFilterOptions(ctx)
	if err != nil {
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	return Response(http.StatusOK, filterOptions), nil
}

func (m *ModelCatalogServiceAPIService) GetModel(ctx context.Context, sourceID, modelName string) (ImplResponse, error) {
	if newName, err := url.PathUnescape(modelName); err == nil {
		modelName = newName
	}

	model, err := m.provider.GetModel(ctx, modelName, sourceID)
	if err != nil {
		statusCode := api.ErrToStatus(err)
		var errorMsg string
		if errors.Is(err, api.ErrNotFound) {
			errorMsg = fmt.Sprintf("No model found '%s' in source '%s'", modelName, sourceID)
		} else {
			errorMsg = err.Error()
		}
		return ErrorResponse(statusCode, errors.New(errorMsg)), err
	}

	if model == nil {
		return notFound("Unknown model or version"), nil
	}

	return Response(http.StatusOK, model), nil
}

func (m *ModelCatalogServiceAPIService) FindSources(ctx context.Context, name string, strPageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	sources := m.sources.All()
	if len(sources) > math.MaxInt32 {
		err := errors.New("too many registered models")
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	paginator, err := newPaginator[model.CatalogSource](strPageSize, orderBy, sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	items := make([]model.CatalogSource, 0, len(sources))

	name = strings.ToLower(name)

	for _, v := range sources {
		if !strings.Contains(strings.ToLower(v.Name), name) {
			continue
		}

		items = append(items, v)
	}

	cmpFunc, err := genCatalogCmpFunc(orderBy, sortOrder)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	slices.SortStableFunc(items, cmpFunc)

	total := int32(len(items))

	pagedItems, next := paginator.Paginate(items)

	res := model.CatalogSourceList{
		PageSize:      paginator.PageSize,
		Items:         pagedItems,
		Size:          total,
		NextPageToken: next.Token(),
	}
	return Response(http.StatusOK, res), nil
}

func genCatalogCmpFunc(orderBy model.OrderByField, sortOrder model.SortOrder) (func(model.CatalogSource, model.CatalogSource) int, error) {
	multiplier := 1
	switch model.SortOrder(strings.ToUpper(string(sortOrder))) {
	case model.SORTORDER_DESC:
		multiplier = -1
	case model.SORTORDER_ASC, "":
		multiplier = 1
	default:
		return nil, fmt.Errorf("unsupported sort order field")
	}

	switch model.OrderByField(strings.ToUpper(string(orderBy))) {
	case model.ORDERBYFIELD_ID, "":
		return func(a, b model.CatalogSource) int {
			return multiplier * strings.Compare(a.Id, b.Id)
		}, nil
	case model.ORDERBYFIELD_NAME:
		return func(a, b model.CatalogSource) int {
			return multiplier * strings.Compare(a.Name, b.Name)
		}, nil
	default:
		return nil, fmt.Errorf("unsupported order by field")
	}
}

var _ ModelCatalogServiceAPIServicer = &ModelCatalogServiceAPIService{}

// NewModelCatalogServiceAPIService creates a default api service
func NewModelCatalogServiceAPIService(provider catalog.APIProvider, sources *catalog.SourceCollection) ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{
		provider: provider,
		sources:  sources,
	}
}

func notFound(msg string) ImplResponse {
	if msg == "" {
		msg = "Resource not found"
	}
	return ErrorResponse(http.StatusNotFound, errors.New(msg))
}
