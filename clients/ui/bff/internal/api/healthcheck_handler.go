package api

import (
	"errors"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"net/http"
)

func (app *App) HealthcheckHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	userId, ok := r.Context().Value(constants.KubeflowUserIdKey).(string)
	if !ok || userId == "" {
		app.serverErrorResponse(w, r, errors.New("failed to retrieve kubeflow-userid from context"))
		return
	}

	healthCheck, err := app.repositories.HealthCheck.HealthCheck(Version, userId)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, healthCheck, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
