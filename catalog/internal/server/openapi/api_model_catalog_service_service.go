package openapi

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/catalog/internal/catalog"
	"github.com/kubeflow/model-registry/catalog/internal/db/models"
	model "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/pkg/api"
)

// ModelCatalogServiceAPIService is a service that implements the logic for the ModelCatalogServiceAPIServicer
// This service should implement the business logic for every endpoint for the ModelCatalogServiceAPI s.coreApi.
// Include any external packages or services that will be required by this service.
type ModelCatalogServiceAPIService struct {
	provider         catalog.APIProvider
	sources          *catalog.SourceCollection
	labels           *catalog.LabelCollection
	sourceRepository models.CatalogSourceRepository
}

// GetAllModelArtifacts retrieves all model artifacts for a given model from the specified source.
func (m *ModelCatalogServiceAPIService) GetAllModelArtifacts(ctx context.Context, sourceID string, modelName string, artifactType []model.ArtifactTypeQueryParam, artifactType2 []model.ArtifactTypeQueryParam, filterQuery string, pageSize string, orderBy string, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
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
		FilterQuery:         filterQuery,
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
			// Use the original error message which should be more specific
			errorMsg = err.Error()
		} else if errors.Is(err, api.ErrNotFound) {
			errorMsg = fmt.Sprintf("No model found '%s' in source '%s'", modelName, sourceID)
		} else {
			errorMsg = err.Error()
		}
		return ErrorResponse(statusCode, errors.New(errorMsg)), err
	}

	return Response(http.StatusOK, artifacts), nil
}

func (m *ModelCatalogServiceAPIService) GetAllModelPerformanceArtifacts(ctx context.Context, sourceID string, modelName string, targetRPS int32, recommendations bool, rpsProperty string, latencyProperty string, hardwareCountProperty string, hardwareTypeProperty string, filterQuery string, pageSize string, orderBy string, sortOrder model.SortOrder, nextPageToken string) (ImplResponse, error) {
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

	// Call the provider's GetPerformanceArtifacts method
	artifacts, err := m.provider.GetPerformanceArtifacts(ctx, modelName, sourceID, catalog.ListPerformanceArtifactsParams{
		FilterQuery:           filterQuery,
		PageSize:              pageSizeInt,
		OrderBy:               orderBy,
		SortOrder:             sortOrder,
		NextPageToken:         &nextPageToken,
		TargetRPS:             targetRPS,
		Recommendations:       recommendations,
		RPSProperty:           rpsProperty,
		LatencyProperty:       latencyProperty,
		HardwareCountProperty: hardwareCountProperty,
		HardwareTypeProperty:  hardwareTypeProperty,
	})
	if err != nil {
		statusCode := api.ErrToStatus(err)
		var errorMsg string
		if errors.Is(err, api.ErrBadRequest) {
			errorMsg = err.Error()
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

		// Create CatalogLabel with name (which may be null)
		var label *model.CatalogLabel
		if nameStr, ok := name.(string); ok {
			label = model.NewCatalogLabel(*model.NewNullableString(&nameStr))
		} else {
			label = model.NewCatalogLabel(*model.NewNullableString(nil))
		}

		// Add all other properties to AdditionalProperties
		label.AdditionalProperties = make(map[string]any)
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

	// Fetch status from database
	var statuses map[string]models.SourceStatus
	if m.sourceRepository != nil {
		var err error
		statuses, err = m.sourceRepository.GetAllStatuses()
		if err != nil {
			// Log error but continue - status is optional
			statuses = nil
		}
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

		// Merge status from database if available
		if statuses != nil {
			if status, ok := statuses[v.Id]; ok {
				if status.Status != "" {
					statusEnum := model.CatalogSourceStatus(status.Status)
					v.Status = &statusEnum
				}
				if status.Error != "" {
					v.Error = *model.NewNullableString(&status.Error)
				} else {
					v.Error = *model.NewNullableString(nil)
				}
			}
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

func (m *ModelCatalogServiceAPIService) PreviewCatalogSource(ctx context.Context, configParam *os.File, pageSizeParam string, nextPageTokenParam string, filterStatusParam string, catalogDataParam *os.File) (ImplResponse, error) {
	// Parse page size
	pageSize := int32(10)
	if pageSizeParam != "" {
		parsed, err := strconv.ParseInt(pageSizeParam, 10, 32)
		if err != nil {
			return ErrorResponse(http.StatusBadRequest, fmt.Errorf("invalid pageSize: %w", err)), err
		}
		pageSize = int32(parsed)
	}

	// Parse filterStatus (default: "all")
	filterStatus := "all"
	if filterStatusParam != "" {
		filterStatus = strings.ToLower(filterStatusParam)
		if filterStatus != "all" && filterStatus != "included" && filterStatus != "excluded" {
			err := fmt.Errorf("invalid filterStatus: must be 'all', 'included', or 'excluded'")
			return ErrorResponse(http.StatusBadRequest, err), err
		}
	}

	// Read and parse the uploaded config file
	if configParam == nil {
		err := errors.New("config file is required")
		return ErrorResponse(http.StatusBadRequest, err), err
	}
	defer configParam.Close()

	configBytes, err := os.ReadFile(configParam.Name())
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, fmt.Errorf("failed to read config file: %w", err)), err
	}

	// Read catalog data if provided (stateless mode)
	var catalogDataBytes []byte
	if catalogDataParam != nil {
		defer catalogDataParam.Close()
		catalogDataBytes, err = os.ReadFile(catalogDataParam.Name())
		if err != nil {
			return ErrorResponse(http.StatusBadRequest, fmt.Errorf("failed to read catalogData file: %w", err)), err
		}
	}

	// Parse the config as a preview request
	previewRequest, err := catalog.ParsePreviewConfig(configBytes)
	if err != nil {
		return ErrorResponse(http.StatusUnprocessableEntity, fmt.Errorf("invalid config: %w", err)), err
	}

	// Load models using the preview service
	previewResults, err := catalog.PreviewSourceModels(ctx, previewRequest, catalogDataBytes)
	if err != nil {
		return ErrorResponse(http.StatusUnprocessableEntity, fmt.Errorf("failed to load models: %w", err)), err
	}

	// Filter by status
	var filteredResults []model.ModelPreviewResult
	for _, result := range previewResults {
		switch filterStatus {
		case "included":
			if result.Included {
				filteredResults = append(filteredResults, result)
			}
		case "excluded":
			if !result.Included {
				filteredResults = append(filteredResults, result)
			}
		default: // "all"
			filteredResults = append(filteredResults, result)
		}
	}

	// Calculate summary from ALL results (not filtered)
	var includedCount, excludedCount int32
	for _, result := range previewResults {
		if result.Included {
			includedCount++
		} else {
			excludedCount++
		}
	}

	summary := model.CatalogSourcePreviewResponseAllOfSummary{
		TotalModels:    int32(len(previewResults)),
		IncludedModels: includedCount,
		ExcludedModels: excludedCount,
	}

	// Apply pagination
	paginator, err := newPaginator[model.ModelPreviewResult](pageSizeParam, model.OrderByField(""), model.SortOrder(""), nextPageTokenParam)
	if err != nil {
		return ErrorResponse(http.StatusBadRequest, err), err
	}

	pagedResults, next := paginator.Paginate(filteredResults)

	response := model.CatalogSourcePreviewResponse{
		PageSize:      pageSize,
		Size:          int32(len(pagedResults)),
		NextPageToken: next.Token(),
		Items:         pagedResults,
		Summary:       summary,
	}

	return Response(http.StatusOK, response), nil
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
	data  map[string]any
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
		aValRaw, aHasKey := a.data[orderByKey]
		bValRaw, bHasKey := b.data[orderByKey]

		var aVal string
		if aHasKey {
			aVal, aHasKey = aValRaw.(string)
		}
		var bVal string
		if bHasKey {
			bVal, bHasKey = bValRaw.(string)
		}

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
func NewModelCatalogServiceAPIService(provider catalog.APIProvider, sources *catalog.SourceCollection, labels *catalog.LabelCollection, sourceRepository models.CatalogSourceRepository) ModelCatalogServiceAPIServicer {
	return &ModelCatalogServiceAPIService{
		provider:         provider,
		sources:          sources,
		labels:           labels,
		sourceRepository: sourceRepository,
	}
}

func notFound(msg string) ImplResponse {
	if msg == "" {
		msg = "Resource not found"
	}
	return ErrorResponse(http.StatusNotFound, errors.New(msg))
}
