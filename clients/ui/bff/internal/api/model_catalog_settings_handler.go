package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogSettingsSourceConfigEnvelope Envelope[models.CatalogSourceConfig, None]
type ModelCatalogSettingsSourceConfigListEnvelope Envelope[models.CatalogSourceConfigList, None]
type ModelCatalogSourcePayloadEnvelope Envelope[models.CatalogSourceConfigPayload, None]

func (app *App) GetAllCatalogSourceConfigsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	catalogConfigs := []models.CatalogSourceConfig{
		createSampleCatalogSource("catalog-1", "Default Catalog", "yaml"),
		createSampleCatalogSource("catalog-2", "HuggingFace Catalog", "huggingface"),
		createSampleCatalogSource("catalog-3", "Custom Catalog", "yaml"),
	}

	catalogSourceConfigs := models.CatalogSourceConfigList{
		Catalogs: catalogConfigs,
	}

	modelCatalogSource := ModelCatalogSettingsSourceConfigListEnvelope{
		Data: catalogSourceConfigs,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) GetCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	catalogSourceId := ps.ByName(CatalogSourceId)

	catalogSourceConfig := createSampleCatalogSource(catalogSourceId, "catalog-source-1", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: catalogSourceConfig,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) CreateCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	var envelope ModelCatalogSourcePayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	var sourceName = envelope.Data.Name
	var sourceId = envelope.Data.Id
	var sourceType = envelope.Data.Type

	if sourceName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source name is required"))
		return
	}
	if sourceId == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source ID is required"))
		return
	}
	if sourceType == "" {
		app.badRequestResponse(w, r, fmt.Errorf("source type is required"))
		return
	}

	ctxLogger.Info("Creating catalog source", "name", sourceName)

	newCatalogSource := createSampleCatalogSource(sourceId, sourceName, sourceType)

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: newCatalogSource,
	}

	w.Header().Set("Location", r.URL.JoinPath(modelCatalogSource.Data.Id).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, modelCatalogSource, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}

}

func (app *App) UpdateCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	// this is the temoprary fix to start fronetend development
	catalogSourceId := ps.ByName(CatalogSourceId)

	newCatalogSource := createSampleCatalogSource(catalogSourceId, "Updated Catalog", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: newCatalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	// this is the temoprary fix to start fronetend development
	catalogSourceId := ps.ByName(CatalogSourceId)

	newCatalogSource := createSampleCatalogSource(catalogSourceId, "Deleted Catalog", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: newCatalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func createSampleCatalogSource(id string, name string, catalogType string) models.CatalogSourceConfig {
	return models.CatalogSourceConfig{
		Name:           name,
		Id:             id,
		Type:           catalogType,
		Enabled:        BoolPtr(true),
		Labels:         []string{},
		IncludedModels: []string{"rhelai1/modelcar-granite-7b-starter"},
		ExcludedModels: []string{"model-a:1.0", "model-b:*"},
		Properties: &models.CatalogSourceProperties{
			YamlCatalogPath: stringToPointer("/path/to/catalog.yaml"),
		},
	}
}

func stringToPointer(s string) *string {
	return &s
}

func BoolPtr(b bool) *bool {
	return &b
}
