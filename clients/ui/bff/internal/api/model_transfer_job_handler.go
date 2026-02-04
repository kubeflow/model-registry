package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelTransferJobListEnvelope Envelope[*models.ModelTransferJobList, None]
type ModelTransferJobEnvelope Envelope[*models.ModelTransferJob, None]

func (app *App) GetAllModelTransferJobsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("kubernetes client not found"))
		return
	}

	transferJobs, err := app.repositories.ModelRegistry.GetAllModelTransferJobs(ctx, client, namespace)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := ModelTransferJobListEnvelope{
		Data: transferJobs,
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateModelTransferJobHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("kubernetes client not found"))
		return
	}

	var envelope ModelTransferJobEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON: %v", err.Error()))
		return
	}
	payload := *envelope.Data

	err = app.repositories.ModelRegistry.CreateModelTransferJob(ctx, client, namespace, payload)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	response := ModelTransferJobEnvelope{Data: &payload}

	// TODO: uncomment the following when we implement the actual logic
	// modelTransferJob := ModelTransferJobEnvelope{
	// 	Data: newModelTransferJob,
	// }

	// w.Header().Set("Location", r.URL.JoinPath(modelTransferJob.Data.Id).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, response, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) UpdateModelTransferJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("kubernetes client not found"))
		return
	}

	var updates map[string]string
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	jobId := ps.ByName(ModelTransferJobId)

	err = app.repositories.ModelRegistry.UpdateModelTransferJob(ctx, client, namespace, jobId, updates)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// TODO: uncomment the following when we implement the actual logic
	// modelTransferJob := ModelTransferJobEnvelope{
	// 	Data: updatedModelTransferJob,
	// }

	err = app.WriteJSON(w, http.StatusOK, map[string]string{"status": "updated"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) DeleteModelTransferJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctx := r.Context()

	namespace, ok := ctx.Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
		return
	}

	client, err := app.kubernetesClientFactory.GetClient(ctx)
	if err != nil {
		app.serverErrorResponse(w, r, errors.New("kubernetes client not found"))
		return
	}

	jobId := ps.ByName(ModelTransferJobId)
	err = app.repositories.ModelRegistry.DeleteModelTransferJob(ctx, client, namespace, jobId)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)

	// TODO: uncomment the following when we implement the actual logic
	// modelTransferJob := ModelTransferJobEnvelope{
	// 	Data: deletedModelTransferJob,
	// }

	// err = app.WriteJSON(w, http.StatusCreated, modelTransferJob, nil)
	// if err != nil {
	// 	app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
	// 	return
	// }
}
