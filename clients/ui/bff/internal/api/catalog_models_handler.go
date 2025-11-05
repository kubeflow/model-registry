package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	catalogOpenapi "github.com/kubeflow/model-registry/catalog/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
)

type CatalogModelListEnvelope Envelope[*catalogOpenapi.CatalogModelList, None]

func (app *App) GetAllCatalogModelsAcrossSourcesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	catalogModels, err := app.repositories.ModelCatalogClient.GetAllCatalogModelsAcrossSources(client, r.URL.Query())

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	sourcesList := CatalogModelListEnvelope{
		Data: catalogModels,
	}

	err = app.WriteJSON(w, http.StatusOK, sourcesList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
