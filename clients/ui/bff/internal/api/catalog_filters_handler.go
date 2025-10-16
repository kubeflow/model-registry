package api

import (
	"errors"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogFilterOptionEnvelope Envelope[*models.FilterOption, None]
type CatalogFilterOptionsListEnvelope Envelope[*models.FilterOptionsList, None]

func (app *App) GetCatalogFilterListHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)

	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	catalogFilterOptions, err := app.repositories.ModelCatalogClient.GetCatalogFilterOptions(client)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	catalogFilterOptionList := CatalogFilterOptionsListEnvelope{
		Data: catalogFilterOptions,
	}

	err = app.WriteJSON(w, http.StatusOK, catalogFilterOptionList, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
