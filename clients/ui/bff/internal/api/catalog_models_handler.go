package api

import (
	"errors"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogModelEnvelope Envelope[models.CatalogModel, None]
type CatalogModelListEnvelope Envelope[models.CatalogModelList, None]
type CatalogGroupedResponseEnvelope Envelope[CatalogSourceGroup, None]

type CatalogSourceGroup struct {
	Sources    []SourceGroup  `json:"sources"`
	Pagination PaginationMeta `json:"pagination"`
}

type SourceGroup struct {
	Source string                `json:"source"`
	Models []models.CatalogModel `json:"models"`
}

type PaginationMeta struct {
	ReturnedModels int32  `json:"returnedModels"`
	PageSize       int32  `json:"pageSize"`
	NextPageToken  string `json:"nextPageToken"`
}

func (app *App) GetAllCatalogModelsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	source := r.URL.Query().Get("source")
	if source == "" {
		app.badRequestResponse(w, r, errors.New("source query parameter is required"))
		return
	}

	query := r.URL.Query().Get("q")

	// TODO: Implement actual catalog API call
	// For now, return mock response to unblock frontend development
	allMockModels := mocks.GetCatalogModelMocks()

	var filteredModels []models.CatalogModel

	if source == "ALL" || source == "all" {
		filteredModels = allMockModels
	} else {
		for _, model := range allMockModels {
			if model.SourceId != nil && *model.SourceId == source {
				filteredModels = append(filteredModels, model)
			}
		}
	}

	// TODO: this is to unblock frontend development, will be removed in actual implementation
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

			// Check tasks
			if !matchFound && model.Tasks != nil {
				for _, task := range model.Tasks {
					if strings.Contains(strings.ToLower(task), queryLower) {
						matchFound = true
						break
					}
				}
			}

			// Check license
			if !matchFound && model.License != nil && strings.Contains(strings.ToLower(*model.License), queryLower) {
				matchFound = true
			}

			if matchFound {
				queryFilteredModels = append(queryFilteredModels, model)
			}
		}

		filteredModels = queryFilteredModels
	}

	sourceGroups := groupModelsBySource(filteredModels)

	catalogResponse := CatalogSourceGroup{
		Sources: sourceGroups,
		Pagination: PaginationMeta{
			ReturnedModels: int32(len(filteredModels)),
			PageSize:       int32(10),
			NextPageToken:  "",
		},
	}

	catalogRes := CatalogGroupedResponseEnvelope{
		Data: catalogResponse,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogRes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func groupModelsBySource(catalogModels []models.CatalogModel) []SourceGroup {
	sourceMap := make(map[string][]models.CatalogModel)

	for _, model := range catalogModels {
		source := "unknown"
		if model.SourceId != nil {
			source = *model.SourceId
		}
		sourceMap[source] = append(sourceMap[source], model)
	}

	var sources []SourceGroup
	for source, models := range sourceMap {
		sources = append(sources, SourceGroup{
			Source: source,
			Models: models,
		})
	}

	return sources
}

func (app *App) GetCatalogModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	catalogSourceID := ps.ByName(SourceId)
	modelName := ps.ByName(CatalogModelName)

	allMockModels := mocks.GetCatalogModelMocks()
	var filteredMockModels []models.CatalogModel

	for _, model := range allMockModels {
		if model.Name == modelName && model.SourceId != nil && *model.SourceId == catalogSourceID {
			filteredMockModels = append(filteredMockModels, model)
		}
	}

	if len(filteredMockModels) == 0 {
		app.notFoundResponse(w, r)
		return
	}

	catalogModelList := models.CatalogModelList{
		Items:         filteredMockModels,
		Size:          int32(len(filteredMockModels)),
		PageSize:      int32(10),
		NextPageToken: "",
	}

	catalogModels := CatalogModelListEnvelope{
		Data: catalogModelList,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogModels, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
