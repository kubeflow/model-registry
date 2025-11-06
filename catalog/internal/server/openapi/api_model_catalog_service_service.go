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
	labels   *catalog.LabelCollection
}

// GetAllModelArtifacts retrieves all model artifacts for a given model from the specified source.
func (m *ModelCatalogServiceAPIService) GetAllModelArtifacts(ctx context.Context, sourceID string, modelName string, artifactType []model.ArtifactTypeQueryParam, artifactType2 []model.ArtifactTypeQueryParam, pageSize string, orderBy model.OrderByField, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	// Handle multiple artifact_type parameters (snake case - deprecated, will be removed in future)
	for _, v := range artifactType2 {
		if v != "" {
			artifactType = append(artifactType, v)
		}
	}

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

	// Handle multiple artifact types
	var artifactTypesFilter []string

	if len(artifactType) > 0 {
		// Convert slice of ArtifactTypeQueryParam to slice of strings
		artifactTypesFilter = make([]string, len(artifactType))
		for i, at := range artifactType {
			artifactTypesFilter[i] = string(at)
		}
	}

	artifacts, err := m.provider.GetArtifacts(ctx, modelName, sourceID, catalog.ListArtifactsParams{
		ArtifactTypesFilter: artifactTypesFilter,
		PageSize:            pageSizeInt,
		OrderBy:             orderBy,
		SortOrder:           sortOrder,
		NextPageToken:       &nextPageToken,
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

func (m *ModelCatalogServiceAPIService) FindLabels(ctx context.Context, pageSize string, orderBy string, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
	labels := m.labels.All()
	if len(labels) > math.MaxInt32 {
		err := errors.New("too many registered labels")
		return ErrorResponse(http.StatusInternalServerError, err), err
	}

	// Wrap labels to make them sortable
	sortableLabels := make([]sortableLabel, len(labels))
	for i, label := range labels {
		sortableLabels[i] = sortableLabel{
			data:  label,
			index: i, // Keep original index for stable sort
			id:    generateLabelID(i),
		}
	}

	// Create paginator - use empty OrderByField since we don't use it for labels
	paginator, err := newPaginator[sortableLabel](pageSize, model.OrderByField(""), sortOrder, nextPageToken)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	// Create comparison function for labels using the string key
	cmpFunc := genLabelCmpFunc(orderBy, sortOrder)
	slices.SortStableFunc(sortableLabels, cmpFunc)

	// Paginate the sorted labels
	pagedSortableLabels, next := paginator.Paginate(sortableLabels)

	// Convert map[string]string to model.CatalogLabel
	pagedLabels := make([]model.CatalogLabel, len(pagedSortableLabels))
	for i, sl := range pagedSortableLabels {
		// Extract the "name" field (required)
		name, ok := sl.data["name"]
		if !ok || name == "" {
			err := fmt.Errorf("internal error: label at index %d missing required name field", i)
			return ErrorResponse(http.StatusInternalServerError, err), err
		}

		// Create CatalogLabel with name
		label := model.NewCatalogLabel(name)

		// Add all other properties to AdditionalProperties
		label.AdditionalProperties = make(map[string]interface{})
		for key, value := range sl.data {
			if key != "name" {
				label.AdditionalProperties[key] = value
			}
		}

		pagedLabels[i] = *label
	}

	res := model.CatalogLabelList{
		PageSize:      paginator.PageSize,
		Items:         pagedLabels,
		Size:          int32(len(pagedLabels)), // Number of items in current page, not total
		NextPageToken: next.Token(),
	}
	return Response(http.StatusOK, res), nil
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

	pagedItems, next := paginator.Paginate(items)

	res := model.CatalogSourceList{
		PageSize:      paginator.PageSize,
		Items:         pagedItems,
		Size:          int32(len(pagedItems)), // Number of items in current page, not total
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

// generateLabelID creates a stable, unique ID for a label based on its index
func generateLabelID(index int) string {
	return strconv.Itoa(index)
}

// sortableLabel wraps a label map to make it sortable
type sortableLabel struct {
	data  map[string]string
	index int    // Original position for stable sort when key is missing
	id    string // Stable ID for pagination
}

// SortValue implements the Sortable interface for labels
func (sl sortableLabel) SortValue(field model.OrderByField) string {
	// Return ID for pagination purposes
	if field == model.ORDERBYFIELD_ID {
		return sl.id
	}
	// For other fields, labels use string keys directly in genLabelCmpFunc
	return ""
}

// genLabelCmpFunc generates a comparison function for sorting labels by a string key
func genLabelCmpFunc(orderByKey string, sortOrder model.SortOrder) func(sortableLabel, sortableLabel) int {
	multiplier := 1
	switch model.SortOrder(strings.ToUpper(string(sortOrder))) {
	case model.SORTORDER_DESC:
		multiplier = -1
	case model.SORTORDER_ASC, "":
		multiplier = 1
	}

	return func(a, b sortableLabel) int {
		// If no orderBy key specified, maintain original order
		if orderByKey == "" {
			if a.index < b.index {
				return -1
			}
			if a.index > b.index {
				return 1
			}
			return 0
		}

		// Get values for the orderBy key
		aVal, aHasKey := a.data[orderByKey]
		bVal, bHasKey := b.data[orderByKey]

		// If both have the key, compare their values
		if aHasKey && bHasKey {
			return multiplier * strings.Compare(aVal, bVal)
		}

		// If only one has the key, put it first
		if aHasKey && !bHasKey {
			return -1 // a comes first
		}
		if !aHasKey && bHasKey {
			return 1 // b comes first
		}

		// If neither has the key, maintain original order
		if a.index < b.index {
			return -1
		}
		if a.index > b.index {
			return 1
		}
		return 0
	}
}

var _ ModelCatalogServiceAPIServicer = &ModelCatalogServiceAPIService{}

// NewModelCatalogServiceAPIService creates a default api service
func NewModelCatalogServiceAPIService(provider catalog.APIProvider, sources *catalog.SourceCollection, labels *catalog.LabelCollection) ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{
		provider: provider,
		sources:  sources,
		labels:   labels,
	}
}

func notFound(msg string) ImplResponse {
	if msg == "" {
		msg = "Resource not found"
	}
	return ErrorResponse(http.StatusNotFound, errors.New(msg))
}
