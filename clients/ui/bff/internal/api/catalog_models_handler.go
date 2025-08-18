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

	// TODO: this is to unblock frontend development, will be removed in actual implementation
	for _, model := range allMockModels {
		if model.SourceId != nil && *model.SourceId == source {
			filteredModels = append(filteredModels, model)
		}
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

	catalogModels := models.CatalogModelList{
		Items:         filteredModels,
		Size:          int32(len(filteredModels)),
		PageSize:      int32(10),
		NextPageToken: "",
	}

	catalogModelList := CatalogModelListEnvelope{
		Data: catalogModels,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogModelList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetCatalogModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	catalogSourceID := ps.ByName(SourceId)
	modelPath := ps.ByName(CatalogModelName)

	// This will handle tha wildcard route by parsing the model name
	if _, ok := strings.CutSuffix(modelPath, "/artifacts"); ok {
		app.GetAllCatalogModelArtifactsHandler(w, r, ps)
		return
	}

	modelPath = strings.TrimPrefix(modelPath, "/")

	// TODO: Implement actual catalog API call
	allMockModels := mocks.GetCatalogModelMocks()
	var catalogMockModels models.CatalogModel

	for _, model := range allMockModels {
		if model.Name == modelPath && model.SourceId != nil && *model.SourceId == catalogSourceID {
			catalogMockModels = model
		}
	}

	catalogModels := CatalogModelEnvelope{
		Data: catalogMockModels,
	}

	err := app.WriteJSON(w, http.StatusOK, catalogModels, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
