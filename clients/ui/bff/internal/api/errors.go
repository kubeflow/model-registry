package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
)

type HTTPError struct {
	StatusCode int `json:"-"`
	ErrorResponse
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorEnvelope struct {
	Error *mrserver.HTTPError `json:"error"`
}

func (app *App) LogError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
}

func (app *App) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusBadRequest,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusBadRequest),
			Message: err.Error(),
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) forbiddenResponse(w http.ResponseWriter, r *http.Request, message string) {
	// Log the detailed error message as a warning
	app.logger.Warn("Access forbidden", "message", message, "method", r.Method, "uri", r.URL.RequestURI())

	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusForbidden,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusForbidden),
			Message: "Access forbidden",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) errorResponse(w http.ResponseWriter, r *http.Request, error *mrserver.HTTPError) {

	env := ErrorEnvelope{Error: error}

	err := app.WriteJSON(w, error.StatusCode, env, nil)

	if err != nil {
		app.LogError(r, err)
		w.WriteHeader(error.StatusCode)
	}
}

func (app *App) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.LogError(r, err)

	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusInternalServerError,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusInternalServerError),
			Message: "the server encountered a problem and could not process your request",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) notFoundResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusNotFound,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusNotFound),
			Message: "the requested resource could not be found",
		},
	}
	app.errorResponse(w, r, httpError)
}

func (app *App) methodNotAllowedResponse(w http.ResponseWriter, r *http.Request) {

	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusMethodNotAllowed,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusMethodNotAllowed),
			Message: fmt.Sprintf("the %s method is not supported for this resource", r.Method),
		},
	}
	app.errorResponse(w, r, httpError)
}

// TODO remove nolint comment below when we use this method
func (app *App) failedValidationResponse(w http.ResponseWriter, r *http.Request, errors map[string]string) { //nolint:unused

	message, err := json.Marshal(errors)
	if err != nil {
		message = []byte("{}")
	}
	httpError := &mrserver.HTTPError{
		StatusCode: http.StatusUnprocessableEntity,
		ErrorResponse: mrserver.ErrorResponse{
			Code:    strconv.Itoa(http.StatusUnprocessableEntity),
			Message: string(message),
		},
	}
	app.errorResponse(w, r, httpError)
}
