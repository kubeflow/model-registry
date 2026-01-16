package mocks

import (
	"fmt"
	"log/slog"
	"math"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/kubeflow/model-registry/ui/bff/internal/models"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/stretchr/testify/mock"
)

func isModelValidated(model models.CatalogModel) bool {
	if model.CustomProperties == nil {
		return false
	}
	_, hasValidated := (*model.CustomProperties)["validated"]
	return hasValidated
}

func parseFilterQuery(filterQuery string) []string {
	if filterQuery == "" {
		return nil
	}
	parts := regexp.MustCompile(`(?i)\s+AND\s+`).Split(filterQuery, -1)
	var conditions []string
	for _, part := range parts {
		trimmed := strings.TrimSpace(part)
		if trimmed != "" {
			conditions = append(conditions, trimmed)
		}
	}
	return conditions
}

func hasArtifactPrefix(condition string) bool {
	return strings.HasPrefix(condition, "artifacts.")
}

func modelMatchesFilterQuery(model models.CatalogModel, filterQuery string) bool {
	conditions := parseFilterQuery(filterQuery)
	if len(conditions) == 0 {
		return true // No filter, include all
	}

	for _, condition := range conditions {
		if hasArtifactPrefix(condition) {
			if !isModelValidated(model) {
				return false
			}
			continue
		}

		if !modelMatchesCondition(model, condition) {
			return false
		}
	}

	return true
}

func modelMatchesCondition(model models.CatalogModel, condition string) bool {
	eqRegex := regexp.MustCompile(`^(\w+)='([^']*)'$`)
	if matches := eqRegex.FindStringSubmatch(condition); matches != nil {
		field := strings.ToLower(matches[1])
		value := matches[2]
		return modelFieldEquals(model, field, value)
	}

	inRegex := regexp.MustCompile(`^(\w+)\s+IN\s+\(([^)]+)\)$`)
	if matches := inRegex.FindStringSubmatch(condition); matches != nil {
		field := strings.ToLower(matches[1])
		valuesStr := matches[2]
		// Parse values from 'value1','value2',...
		valueRegex := regexp.MustCompile(`'([^']*)'`)
		valueMatches := valueRegex.FindAllStringSubmatch(valuesStr, -1)
		var values []string
		for _, vm := range valueMatches {
			values = append(values, vm[1])
		}
		return modelFieldInValues(model, field, values)
	}
	return true
}

func modelFieldEquals(model models.CatalogModel, field string, value string) bool {
	switch field {
	case "provider":
		return model.Provider != nil && strings.EqualFold(*model.Provider, value)
	case "license":
		return model.License != nil && strings.EqualFold(*model.License, value)
	case "tasks":
		for _, task := range model.Tasks {
			if strings.EqualFold(task, value) {
				return true
			}
		}
		return false
	case "language":
		for _, lang := range model.Language {
			if strings.EqualFold(lang, value) {
				return true
			}
		}
		return false
	}
	return true
}

func modelFieldInValues(model models.CatalogModel, field string, values []string) bool {
	for _, value := range values {
		if modelFieldEquals(model, field, value) {
			return true
		}
	}
	return false
}

type ModelCatalogClientMock struct {
	mock.Mock
}

func NewModelCatalogClientMock(logger *slog.Logger) (*ModelCatalogClientMock, error) {
	return &ModelCatalogClientMock{}, nil
}

func (m *ModelCatalogClientMock) GetAllCatalogModelsAcrossSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogModelList, error) {
	allModels := GetCatalogModelMocks()
	var filteredModels []models.CatalogModel

	sourceId := pageValues.Get("source")
	sourceLabel := pageValues.Get("sourceLabel")
	query := pageValues.Get("q")
	filterQuery := pageValues.Get("filterQuery")

	if sourceId != "" {
		for _, model := range allModels {
			if model.SourceId != nil && *model.SourceId == sourceId {
				filteredModels = append(filteredModels, model)
			}
		}
	} else if sourceLabel != "" {
		allSources := GetCatalogSourceMocks()
		var matchingSourceIds []string

		if sourceLabel == "null" {
			for _, source := range allSources {
				if len(source.Labels) == 0 {
					matchingSourceIds = append(matchingSourceIds, source.Id)
				}
			}
		} else {
			for _, source := range allSources {
				for _, label := range source.Labels {
					if label == sourceLabel {
						matchingSourceIds = append(matchingSourceIds, source.Id)
						break
					}
				}
			}
		}

		for _, model := range allModels {
			if model.SourceId != nil {
				for _, sid := range matchingSourceIds {
					if *model.SourceId == sid {
						filteredModels = append(filteredModels, model)
						break
					}
				}
			}
		}
	} else {
		filteredModels = allModels
	}

	if query != "" {
		var queryFilteredModels []models.CatalogModel
		queryLower := strings.ToLower(query)

		for _, model := range filteredModels {
			matchFound := false

			// Check name
			if strings.Contains(strings.ToLower(model.Name), queryLower) {
				matchFound = true
			}

			// Check description
			if !matchFound && model.Description != nil && strings.Contains(strings.ToLower(*model.Description), queryLower) {
				matchFound = true
			}

			// Check provider
			if !matchFound && model.Provider != nil && strings.Contains(strings.ToLower(*model.Provider), queryLower) {
				matchFound = true
			}

			if matchFound {
				queryFilteredModels = append(queryFilteredModels, model)
			}
		}

		filteredModels = queryFilteredModels
	}

	if filterQuery != "" {
		var filterQueryFilteredModels []models.CatalogModel
		for _, model := range filteredModels {
			if modelMatchesFilterQuery(model, filterQuery) {
				filterQueryFilteredModels = append(filterQueryFilteredModels, model)
			}
		}
		filteredModels = filterQueryFilteredModels
	}

	pageSizeStr := pageValues.Get("pageSize")
	pageSize := 10 // default
	if pageSizeStr != "" {
		if parsed, err := strconv.Atoi(pageSizeStr); err == nil && parsed > 0 {
			pageSize = parsed
		}
	}

	pageTokenStr := pageValues.Get("nextPageToken")
	startIndex := 0
	if pageTokenStr != "" {
		if parsed, err := strconv.Atoi(pageTokenStr); err == nil && parsed > 0 {
			startIndex = parsed
		}
	}

	totalSize := len(filteredModels)
	endIndex := startIndex + pageSize
	if endIndex > totalSize {
		endIndex = totalSize
	}

	var pagedModels []models.CatalogModel
	if startIndex < totalSize {
		pagedModels = filteredModels[startIndex:endIndex]
	} else {
		pagedModels = []models.CatalogModel{}
	}

	var nextPageToken string
	if endIndex < totalSize {
		nextPageToken = strconv.Itoa(endIndex)
	}

	size := len(pagedModels)
	if size > math.MaxInt32 {
		size = math.MaxInt32
	}
	ps := pageSize
	if ps > math.MaxInt32 {
		ps = math.MaxInt32
	}

	catalogModelList := models.CatalogModelList{
		Items:         pagedModels,
		Size:          int32(size),
		PageSize:      int32(ps),
		NextPageToken: nextPageToken,
	}

	return &catalogModelList, nil

}

func (m *ModelCatalogClientMock) GetCatalogSourceModel(client httpclient.HTTPClientInterface, sourceId string, modelName string) (*models.CatalogModel, error) {
	allModels := GetCatalogModelMocks()

	decodedModelName, err := url.QueryUnescape(modelName)
	if err != nil {
		return nil, fmt.Errorf("failed to decode modelName: %w", err)
	}

	decodedModelName = strings.TrimPrefix(decodedModelName, "/")

	for _, model := range allModels {
		if model.SourceId != nil && *model.SourceId == sourceId && model.Name == decodedModelName {
			return &model, nil
		}
	}

	return nil, fmt.Errorf("catalog model not found for sourceId: %s, modelName: %s", sourceId, decodedModelName)
}

func (m *ModelCatalogClientMock) GetAllCatalogSources(client httpclient.HTTPClientInterface, pageValues url.Values) (*models.CatalogSourceList, error) {
	allMockSources := GetCatalogSourceListMock()
	var filteredMockSources []models.CatalogSource

	name := pageValues.Get("name")

	if name != "" {
		nameFilterLower := strings.ToLower(name)
		for _, source := range allMockSources.Items {
			if strings.ToLower(source.Id) == nameFilterLower || strings.ToLower(source.Name) == nameFilterLower {
				filteredMockSources = append(filteredMockSources, source)
			}
		}
	} else {
		filteredMockSources = allMockSources.Items
	}
	catalogSourceList := models.CatalogSourceList{
		Items:         filteredMockSources,
		PageSize:      int32(10),
		NextPageToken: "",
		Size:          int32(len(filteredMockSources)),
	}

	return &catalogSourceList, nil
}

func (m *ModelCatalogClientMock) GetCatalogSourceModelArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string, pageValues url.Values) (*models.CatalogModelArtifactList, error) {
	var allMockModelArtifacts models.CatalogModelArtifactList

	if sourceId == "sample-source" && (modelName == "repo1%2Fgranite-8b-code-instruct" || modelName == "repo1%2Fgranite-8b-code-instruct-quantized.w4a16") {
		performanceArtifacts := GetCatalogPerformanceMetricsArtifactListMock(4)
		accuracyArtifacts := GetCatalogAccuracyMetricsArtifactListMock()
		modelArtifacts := GetCatalogModelArtifactListMock()
		combinedItems := append(performanceArtifacts.Items, accuracyArtifacts.Items...)
		combinedItems = append(combinedItems, modelArtifacts.Items...)
		allMockModelArtifacts = models.CatalogModelArtifactList{
			Items:         combinedItems,
			Size:          int32(len(combinedItems)),
			PageSize:      performanceArtifacts.PageSize,
			NextPageToken: "",
		}
	} else if sourceId == "sample-source" && modelName == "repo1%2Fgranite-7b-instruct" {
		accuracyArtifacts := GetCatalogAccuracyMetricsArtifactListMock()
		modelArtifacts := GetCatalogModelArtifactListMock()
		combinedItems := append(accuracyArtifacts.Items, modelArtifacts.Items...)
		allMockModelArtifacts = models.CatalogModelArtifactList{
			Items:         combinedItems,
			Size:          int32(len(combinedItems)),
			PageSize:      accuracyArtifacts.PageSize,
			NextPageToken: "",
		}
	} else if sourceId == "sample-source" && (modelName == "repo1%2Fgranite-3b-code-base") {
		allMockModelArtifacts = GetCatalogModelArtifactListMock()
	} else {
		allMockModelArtifacts = GetCatalogModelArtifactListMock()
	}

	return &allMockModelArtifacts, nil
}

func (m *ModelCatalogClientMock) GetCatalogModelPerformanceArtifacts(client httpclient.HTTPClientInterface, sourceId string, modelName string, pageValues url.Values) (*models.CatalogModelArtifactList, error) {
	allMockModelPerformanceArtifacts := GetCatalogPerformanceMetricsArtifactListMock(4)
	return &allMockModelPerformanceArtifacts, nil

}

func (m *ModelCatalogClientMock) GetCatalogFilterOptions(client httpclient.HTTPClientInterface) (*models.FilterOptionsList, error) {
	filterOptions := GetFilterOptionsListMock()

	return &filterOptions, nil
}

func (m *ModelCatalogClientMock) CreateCatalogSourcePreview(client httpclient.HTTPClientInterface, sourcePreviewPayload models.CatalogSourcePreviewRequest, pageValues url.Values) (*models.CatalogSourcePreviewResult, error) {
	filterStatus := pageValues.Get("filterStatus")
	if filterStatus == "" {
		filterStatus = "all"
	}

	pageSize := 20
	if ps := pageValues.Get("pageSize"); ps != "" {
		_, _ = fmt.Sscanf(ps, "%d", &pageSize)
	}

	nextPageToken := pageValues.Get("nextPageToken")

	catalogSourcePreview := CreateCatalogSourcePreviewMockWithFilter(filterStatus, pageSize, nextPageToken)

	return &catalogSourcePreview, nil
}
