package api

import (
	"errors"
	"net/http"
	"net/url"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
)

type CatalogSourceListEnvelope Envelope[*openapi.CatalogSourceList, None]
type CatalogModelEnvelope Envelope[*openapi.CatalogModel, None]
type catalogModelArtifactsListEnvelope Envelope[*openapi.CatalogModelArtifactList, None]

func (app *App) GetAllCatalogSourcesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	catalogSources, err := app.repositories.ModelCatalogClient.GetAllCatalogSources(client, r.URL.Query())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sourcesList := CatalogSourceListEnvelope{
		Data: catalogSources,
	}

	err = app.WriteJSON(w, http.StatusOK, sourcesList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetCatalogSourceModelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	ps.ByName(CatalogSourceId)
	modelName := strings.TrimPrefix(ps.ByName(CatalogModelName), "/")

	newModelName := url.PathEscape(modelName)

	catalogModel, err := app.repositories.ModelCatalogClient.GetCatalogSourceModel(client, ps.ByName(CatalogSourceId), newModelName)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sourcesList := CatalogModelEnvelope{
		Data: catalogModel,
	}

	err = app.WriteJSON(w, http.StatusOK, sourcesList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetCatalogSourceModelArtifactHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	ps.ByName(CatalogSourceId)
	modelName := strings.TrimPrefix(ps.ByName(CatalogModelName), "/")

	newModelName := url.PathEscape(modelName)

	catalogModelArtifacts, err := app.repositories.ModelCatalogClient.GetCatalogModelArtifacts(client, ps.ByName(CatalogSourceId), newModelName)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	catalogModelArtifactList := catalogModelArtifactsListEnvelope{
		Data: catalogModelArtifacts,
	}

	err = app.WriteJSON(w, http.StatusOK, catalogModelArtifactList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
