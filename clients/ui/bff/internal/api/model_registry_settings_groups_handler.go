package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type GroupsEnvelope Envelope[[]models.Group, None]

// STUB IMPLEMENTATION (see kubernetes clients for more details)
func (app *App) GetGroupsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	ctx := r.Context()
	identity, ok := ctx.Value(constants.RequestIdentityKey).(*kubernetes.RequestIdentity)
	if !ok || identity == nil {
		app.badRequestResponse(w, r, fmt.Errorf("missing RequestIdentity in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return
	}

	groups, err := app.repositories.ModelRegistrySettings.GetGroups(r.Context(), client)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get groups: %w", err))
		return
	}

	resp := GroupsEnvelope{Data: groups}

	err = app.WriteJSON(w, http.StatusOK, resp, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
