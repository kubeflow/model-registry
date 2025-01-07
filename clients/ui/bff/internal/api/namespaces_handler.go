package api

import (
	"errors"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type NamespacesEnvelope Envelope[[]models.NamespaceModel, None]

func (app *App) GetNamespacesHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	userId, ok := r.Context().Value(KubeflowUserIdKey).(string)
	if !ok || userId == "" {
		app.serverErrorResponse(w, r, errors.New("failed to retrieve kubeflow-userid from context"))
		return
	}

	var userGroups []string
	if groups, ok := r.Context().Value(KubeflowUserGroupsKey).([]string); ok {
		userGroups = groups
	} else {
		userGroups = []string{}
	}

	namespaces, err := app.repositories.Namespace.GetNamespaces(app.kubernetesClient, userId, userGroups)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	namespacesEnvelope := NamespacesEnvelope{
		Data: namespaces,
	}

	err = app.WriteJSON(w, http.StatusOK, namespacesEnvelope, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
