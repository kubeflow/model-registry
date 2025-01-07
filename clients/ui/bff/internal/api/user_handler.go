package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"net/http"
)

type UserEnvelope Envelope[*models.User, None]

func (app *App) UserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	userId, ok := r.Context().Value(KubeflowUserIdKey).(string)
	if !ok || userId == "" {
		app.serverErrorResponse(w, r, errors.New("failed to retrieve kubeflow-userid from context"))
		return
	}

	user, err := app.repositories.User.GetUser(app.kubernetesClient, userId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	userRes := UserEnvelope{
		Data: user,
	}

	err = app.WriteJSON(w, http.StatusOK, userRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
