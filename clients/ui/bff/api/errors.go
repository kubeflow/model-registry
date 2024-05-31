package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type HTTPError struct {
	StatusCode int `json:"-"`
	ErrorResponse
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (app *App) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	httpError := &HTTPError{
		StatusCode: http.StatusBadRequest,
		ErrorResponse: ErrorResponse{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: err.Error(),
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) errorResponse(w http.ResponseWriter, r *http.Request, error *HTTPError) {

	env := Envelope{"error": error}

	err := app.WriteJSON(w, error.StatusCode, env, nil)

	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(error.StatusCode)
	}
}

func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	httpError := &HTTPError{
		StatusCode: http.StatusInternalServerError,
		ErrorResponse: ErrorResponse{
			Code:    strconv.Itoa(http.StatusInternalServerError),
			Message: "the server encountered a problem and could not process your request",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &HTTPError{
		StatusCode: http.StatusNotFound,
		ErrorResponse: ErrorResponse{
			Code:    strconv.Itoa(http.StatusNotFound),
			Message: "the requested resource could not be found",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		ErrorResponse: ErrorResponse{
			Code:    strconv.Itoa(http.StatusMethodNotAllowed),
			Message: fmt.Sprintf("the %s method is not supported for this resource", r.Method),
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) {

	message, err := json.Marshal(errors)
	if err != nil {
		message = []byte("{}")
	}
	httpError := &HTTPError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorResponse: ErrorResponse{
			Code:    strconv.Itoa(http.StatusUnprocessableEntity),
			Message: string(message),
		},
	}
	app.errorResponse(w, r, httpError)
}
