package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"net/http"
)

///api/v1/model_registry/{modelRegistryName}/model_versions/{modelversionId} - GET

//func (app *App) GetAllRegisteredModelsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
//	//TODO (ederign) implement pagination
//	client, ok := r.Context().Value(httpClientKey).(integrations.HTTPClientInterface)
//	if !ok {
//		app.serverErrorResponse(w, r, errors.New("REST client not found"))
//		return
//	}
//
//	modelList, err := app.modelRegistryClient.GetAllRegisteredModels(client)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//		return
//	}
//
//	modelRegistryRes := RegisteredModelListEnvelope{
//		Data: modelList,
//	}
//
//	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)
//	if err != nil {
//		app.serverErrorResponse(w, r, err)
//	}
//}

/*
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
*/

type ModelVersionEnvelope Envelope[*openapi.ModelVersion, None]

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
