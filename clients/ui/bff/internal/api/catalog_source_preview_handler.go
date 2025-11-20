package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/httpclient"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type CatalogSourcePreviewEnvelope Envelope[*models.CatalogSourcePreviewResult, None]

func (app *App) CreateCatalogSourcePreviewHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelCatalogHttpClientKey).(httpclient.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("catalog REST client not found"))
		return
	}

	var sourcePreviewPayload models.CatalogSourcePreviewRequest
	if err := json.NewDecoder(r.Body).Decode(&sourcePreviewPayload); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	sourcePreview, err := app.repositories.ModelCatalogClient.CreateCatalogSourcePreview(client, sourcePreviewPayload)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

	catalogSourcePreview := CatalogSourcePreviewEnvelope{
		Data: sourcePreview,
	}

	err = app.WriteJSON(w, http.StatusOK, catalogSourcePreview, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
