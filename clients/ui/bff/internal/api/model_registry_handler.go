package api

import (
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"net/http"
)

type ModelRegistryListEnvelope Envelope[[]models.ModelRegistryModel, None]

func (app *App) ModelRegistryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	registries, err := app.repositories.ModelRegistry.FetchAllModelRegistries(app.kubernetesClient)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := ModelRegistryListEnvelope{
		Data: registries,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
