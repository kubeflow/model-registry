package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"net/http"
)

type UserEnvelope Envelope[*models.User, None]

func (app *App) UserHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	userHeader := r.Header.Get(kubeflowUserId)
	if userHeader == "" {
		app.serverErrorResponse(w, r, errors.New("kubeflow-userid not present on header"))
		return
	}

	user, err := app.repositories.User.GetUser(app.kubernetesClient, userHeader)
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
