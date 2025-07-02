package openapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"slices"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
)

// ModelCatalogServiceAPIService is a service that implements the logic for the ModelCatalogServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelCatalogServiceAPI s.coreApi.
// Include any external packages or services that will be required by this service.
type ModelCatalogServiceAPIService struct {
	sources *catalog.SourceCollection
}

func (m *ModelCatalogServiceAPIService) GetAllModelArtifacts(context.Context, string, string) (ImplResponse, error) {
	return Response(http.StatusNotImplemented, "Not implemented"), nil
}

func (m *ModelCatalogServiceAPIService) FindModels(ctx context.Context, source string, q string, pageSize string, orderBy model.OrderByField, sortOder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	return Response(http.StatusNotImplemented, "Not implemented"), nil
}

func (m *ModelCatalogServiceAPIService) GetModel(ctx context.Context, sourceID string, name string) (ImplResponse, error) {
	if name, ok := strings.CutSuffix(name, "/artifacts"); ok {
		return m.GetAllModelArtifacts(ctx, sourceID, name)
	}

	source, ok := m.sources.Get(sourceID)
	if !ok {
		return notFound("Unknown source"), nil
	}

	model, err := source.Provider.GetModel(ctx, name)
	if err != nil {
		return Response(http.StatusInternalServerError, err), err
	}
	if model == nil {
		return notFound("Unknown model"), nil
	}

	return Response(http.StatusOK, model), nil
}

func (m *ModelCatalogServiceAPIService) FindSources(ctx context.Context, name string, strPageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	// TODO: Implement real pagination in here by reusing the nextPageToken
	// code from https://github.com/kubeflow/model-registry/pull/1205.

	sources := m.sources.All()
	if len(sources) > math.MaxInt32 {
		err := errors.New("too many registered models")
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	var pageSize int32 = 10
	if strPageSize != "" {
		pageSize64, err := strconv.ParseInt(strPageSize, 10, 32)
		if err != nil {
			return ErrorResponse(http.StatusBadRequest, err), err
		}
		pageSize = int32(pageSize64)
	}

	items := make([]model.CatalogSource, 0, len(sources))

	name = strings.ToLower(name)

	for _, v := range sources {
		if !strings.Contains(strings.ToLower(v.Metadata.Name), name) {
			continue
		}

		items = append(items, v.Metadata)
	}

	cmpFunc, err := genCatalogCmpFunc(orderBy, sortOrder)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	slices.SortStableFunc(items, cmpFunc)

	total := int32(len(items))
	if total > pageSize {
		items = items[:pageSize]
	}

	res := model.CatalogSourceList{
		PageSize:      pageSize,
		Items:         items,
		Size:          total,
		NextPageToken: "",
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
func NewModelCatalogServiceAPIService(sources *catalog.SourceCollection) ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{
		sources: sources,
	}
}

func notFound(msg string) ImplResponse {
	if msg == "" {
		msg = "Resource not found"
	}
	return ErrorResponse(http.StatusNotFound, errors.New(msg))
}
