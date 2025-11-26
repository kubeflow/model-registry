package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelCatalogSettingsSourceConfigEnvelope Envelope[*models.CatalogSourceConfig, None]
type ModelCatalogSettingsSourceConfigListEnvelope Envelope[*models.CatalogSourceConfigList, None]
type ModelCatalogSourcePayloadEnvelope Envelope[*models.CatalogSourceConfigPayload, None]

func (app *App) GetAllCatalogSourceConfigsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("catalog client not found"))
		return
	}
	catalogSourceConfigs, err := app.repositories.ModelCatalogSettingsRepository.GetAllCatalogSourceConfigs(ctx, client, namespace)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelCatalogSource := ModelCatalogSettingsSourceConfigListEnvelope{
		Data: catalogSourceConfigs,
	}

	err = app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) GetCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	catalogSourceId := ps.ByName(CatalogSourceId)
	// TODO ppadti write the real implementation here
	catalogSourceConfig := mocks.CreateSampleCatalogSource(catalogSourceId, "catalog-source-1", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: &catalogSourceConfig,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) CreateCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
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
	// TODO ppadti write the real implementation here
	newCatalogSource := mocks.CreateSampleCatalogSource(sourceId, sourceName, sourceType)

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: &newCatalogSource,
	}

	w.Header().Set("Location", r.URL.JoinPath(modelCatalogSource.Data.Id).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, modelCatalogSource, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}

}

func (app *App) UpdateCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	var envelope ModelCatalogSourcePayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	catalogSourceId := envelope.Data.Id
	// TODO ppadti write the real implementation here
	newCatalogSource := mocks.CreateSampleCatalogSource(catalogSourceId, "Updated Catalog", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: &newCatalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) DeleteCatalogSourceConfigHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	// TODO ppadti write the real implementation here
	catalogSourceId := ps.ByName(CatalogSourceId)

	deletedCatalogSource := mocks.CreateSampleCatalogSource(catalogSourceId, "Updated Catalog", "yaml")

	modelCatalogSource := ModelCatalogSettingsSourceConfigEnvelope{
		Data: &deletedCatalogSource,
	}

	err := app.WriteJSON(w, http.StatusOK, modelCatalogSource, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
