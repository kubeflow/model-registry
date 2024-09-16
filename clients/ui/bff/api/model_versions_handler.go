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

type ModelVersionEnvelope Envelope[*openapi.ModelVersion, None]
type ModelVersionListEnvelope Envelope[*openapi.ModelVersionList, None]

func (app *App) GetModelVersionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	model, err := app.modelRegistryClient.GetModelVersion(client, ps.ByName(ModelVersionId))
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if _, ok := model.GetIdOk(); !ok {
		app.notFoundResponse(w, r)
		return
	}

	result := ModelVersionEnvelope{
		Data: model,
	}

	err = app.WriteJSON(w, http.StatusOK, result, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateModelVersionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var envelope ModelVersionEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	if err := validation.ValidateModelVersion(data); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling ModelVersion to JSON: %w", err))
		return
	}

	createdVersion, err := app.modelRegistryClient.CreateModelVersion(client, jsonData)
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
		app.serverErrorResponse(w, r, fmt.Errorf("created ModelVersion is nil"))
		return
	}

	response := ModelVersionEnvelope{
		Data: createdVersion,
	}

	w.Header().Set("Location", r.URL.JoinPath(*createdVersion.Id).String())
	err = app.WriteJSON(w, http.StatusCreated, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) UpdateModelVersionHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var envelope ModelVersionEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	if err := validation.ValidateModelVersion(data); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("validation error:: %v", err.Error()))
		return
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling ModelVersion to JSON: %w", err))
		return
	}

	patchedModel, err := app.modelRegistryClient.UpdateModelVersion(client, ps.ByName(ModelVersionId), jsonData)
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
		app.serverErrorResponse(w, r, fmt.Errorf("patched ModelVersion is nil"))
		return
	}

	responseBody := ModelVersionEnvelope{
		Data: patchedModel,
	}

	err = app.WriteJSON(w, http.StatusOK, responseBody, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}
