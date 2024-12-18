package api

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"net/http"
)

type ModelRegistryListEnvelope Envelope[[]models.ModelRegistryModel, None]

func (app *App) ModelRegistryHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	namespace, ok := r.Context().Value(NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	registries, err := app.repositories.ModelRegistry.GetAllModelRegistries(app.kubernetesClient, namespace)
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
