package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *App) ModelRegistryHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	registries, err := app.models.ModelRegistry.FetchAllModelRegistries(app.kubernetesClient)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := Envelope{
		"model_registry": registries,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
