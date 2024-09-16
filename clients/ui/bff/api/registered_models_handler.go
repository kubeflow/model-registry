package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/kubeflow/model-registry/ui/bff/validation"
	"net/http"
)

type RegisteredModelEnvelope Envelope[*openapi.RegisteredModel, None]
type RegisteredModelListEnvelope Envelope[*openapi.RegisteredModelList, None]

func (app *App) GetAllRegisteredModelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO (ederign) implement pagination
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	modelList, err := app.modelRegistryClient.GetAllRegisteredModels(client)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := RegisteredModelListEnvelope{
		Data: modelList,
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

	var envelope RegisteredModelEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	if err := validation.ValidateRegisteredModel(data); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
		return
	}

	createdModel, err := app.modelRegistryClient.CreateRegisteredModel(client, jsonData)
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

	response := RegisteredModelEnvelope{
		Data: createdModel,
	}

	w.Header().Set("Location", r.URL.JoinPath(*createdModel.Id).String())
	err = app.WriteJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) GetRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	model, err := app.modelRegistryClient.GetRegisteredModel(client, ps.ByName(RegisteredModelId))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if _, ok := model.GetIdOk(); !ok {
		app.notFoundResponse(w, r)
		return
	}

	result := RegisteredModelEnvelope{
		Data: model,
	}

	err = app.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) UpdateRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var envelope RegisteredModelEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	if err := validation.ValidateRegisteredModel(data); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
		return
	}

	patchedModel, err := app.modelRegistryClient.UpdateRegisteredModel(client, ps.ByName("id"), jsonData)
	if err != nil {
		var httpErr *integrations.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if patchedModel == nil {
		app.serverErrorResponse(w, r, fmt.Errorf("patched model is nil"))
		return
	}

	responseBody := RegisteredModelEnvelope{
		Data: patchedModel,
	}

	err = app.WriteJSON(w, http.StatusOK, responseBody, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) GetAllModelVersionsForRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	//TODO (acreasy) implement pagination
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	versionList, err := app.modelRegistryClient.GetAllModelVersions(client, ps.ByName(RegisteredModelId))

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	responseBody := ModelVersionListEnvelope{
		Data: versionList,
	}

	err = app.WriteJSON(w, http.StatusOK, responseBody, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateModelVersionForRegisteredModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var envelope ModelVersionEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	if err := validation.ValidateModelVersion(data); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
	}

	createdVersion, err := app.modelRegistryClient.CreateModelVersionForRegisteredModel(client, ps.ByName(RegisteredModelId), jsonData)
	if err != nil {
		var httpErr *integrations.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if createdVersion == nil {
		app.serverErrorResponse(w, r, fmt.Errorf("created model version is nil"))
		return
	}

	responseBody := ModelVersionEnvelope{
		Data: createdVersion,
	}

	w.Header().Set("Location", ParseURLTemplate(ModelVersionPath, map[string]string{ModelRegistryId: ps.ByName(ModelRegistryId), ModelVersionId: createdVersion.GetId()}))
	err = app.WriteJSON(w, http.StatusCreated, responseBody, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
