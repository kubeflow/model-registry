package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/kubeflow/model-registry/ui/bff/validation"
	"net/http"
)

func (app *App) GetRegisteredModelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO (ederign) implement pagination
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	modelList, err := data.FetchAllRegisteredModels(client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := Envelope{
		"registered_models": modelList,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var model openapi.RegisteredModel
	if err := json.NewDecoder(r.Body).Decode(&model); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	if err := validation.ValidateRegisteredModel(model); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(model)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
		return
	}

	createdModel, err := data.CreateRegisteredModel(client, jsonData)
	if err != nil {
		var httpErr *integrations.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if createdModel == nil {
		app.serverErrorResponse(w, r, fmt.Errorf("created model is nil"))
		return
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", RegisteredModelsPath, *createdModel.Id))
	err = app.WriteJSON(w, http.StatusCreated, createdModel, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}
