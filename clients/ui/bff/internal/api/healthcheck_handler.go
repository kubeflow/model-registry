package api

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *App) HealthcheckHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	healthCheck, err := app.repositories.HealthCheck.HealthCheck(Version)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.WriteJSON(w, http.StatusOK, healthCheck, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
