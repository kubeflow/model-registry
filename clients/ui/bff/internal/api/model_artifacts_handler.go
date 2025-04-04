package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/pkg/openapi"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
)

type ModelArtifactUpdateEnvelope Envelope[*openapi.ModelArtifactUpdate, None]

func (app *App) UpdateModelArtifactHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	client, ok := r.Context().Value(constants.ModelRegistryHttpClientKey).(mrserver.HTTPClientInterface)
	if !ok {
		app.serverErrorResponse(w, r, errors.New("REST client not found"))
		return
	}

	var envelope ModelArtifactUpdateEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	data := *envelope.Data

	jsonData, err := json.Marshal(data)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error marshaling model to JSON: %w", err))
		return
	}

	patchedModelArtifact, err := app.repositories.ModelRegistryClient.UpdateModelArtifact(client, ps.ByName(ArtifactId), jsonData)
	if err != nil {
		var httpErr *mrserver.HTTPError
		if errors.As(err, &httpErr) {
			app.errorResponse(w, r, httpErr)
		} else {
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if patchedModelArtifact == nil {
		app.serverErrorResponse(w, r, fmt.Errorf("created artifact is nil or does not contain valid data"))
		return
	}

	responseBody := ModelArtifactEnvelope{
		Data: patchedModelArtifact,
	}

	err = app.WriteJSON(w, http.StatusOK, responseBody, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
