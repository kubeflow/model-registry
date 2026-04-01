package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	kubernetes "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
)

type ModelTransferJobListEnvelope Envelope[*models.ModelTransferJobList, None]
type ModelTransferJobEnvelope Envelope[*models.ModelTransferJob, None]
type ModelTransferJobEventsEnvelope Envelope[models.ModelTransferJobEventsResponse, None]

func (app *App) GetAllModelTransferJobsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	modelRegistryID := ps.ByName(ModelRegistryId)
	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	// If jobNamespace is provided, scope the query to that namespace.
	// Otherwise, check if the user can list jobs cluster-wide before attempting.
	jobNamespace := r.URL.Query().Get("jobNamespace")
	if jobNamespace == "" {
		identity, ok := ctx.Value(constants.RequestIdentityKey).(*kubernetes.RequestIdentity)
		if !ok {
			app.serverErrorResponse(w, r, fmt.Errorf("request identity not found in context"))
			return
		}
		canList, err := client.CanListJobsClusterWide(ctx, identity)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to check job list permission: %w", err))
			return
		}
		if !canList {
			app.forbiddenResponse(w, r, "user does not have permission to list jobs across all namespaces; provide a jobNamespace query parameter to scope the request")
			return
		}
	}

	transferJobs, err := app.repositories.ModelRegistry.GetAllModelTransferJobs(ctx, client, namespace, modelRegistryID, jobNamespace)
	if err != nil {
		var statusErr *apierrors.StatusError
		if errors.As(err, &statusErr) && apierrors.IsForbidden(statusErr) {
			app.forbiddenResponse(w, r, fmt.Sprintf("you do not have permission to list jobs in namespace %q", jobNamespace))
			return
		}
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

func (app *App) GetModelTransferJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	jobName := ps.ByName(ModelTransferJobName)
	if jobName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("job name is required"))
		return
	}

	jobNamespace, err := getRequiredJobNamespace(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	modelRegistryID := ps.ByName(ModelRegistryId)
	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	modelTransferJob, err := app.repositories.ModelRegistry.GetModelTransferJob(ctx, client, jobNamespace, jobName, modelRegistryID)

	if err != nil {
		if errors.Is(err, repositories.ErrJobNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		if errors.Is(err, repositories.ErrJobValidationFailed) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}
	response := ModelTransferJobEnvelope{Data: modelTransferJob}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) CreateModelTransferJobHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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
		app.badRequestResponse(w, r, fmt.Errorf("error decoding JSON: %w", err))
		return
	}

	if envelope.Data == nil {
		app.badRequestResponse(w, r, fmt.Errorf("data is required"))
		return
	}

	payload := *envelope.Data

	modelRegistryID := ps.ByName(ModelRegistryId)

	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	newJob, err := app.repositories.ModelRegistry.CreateModelTransferJob(ctx, client, namespace, payload, modelRegistryID, app.config.DeploymentMode.IsFederatedMode(), app.podNamespace)
	if err != nil {
		if errors.Is(err, repositories.ErrJobValidationFailed) {
			app.badRequestResponse(w, r, err)
			return
		}
		if errors.Is(err, repositories.ErrModelRegistryNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	modelTransferJob := ModelTransferJobEnvelope{Data: newJob}

	w.Header().Set("Location", r.URL.JoinPath(modelTransferJob.Data.Name).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, modelTransferJob, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON: %w", writeErr))
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

	var envelope ModelTransferJobEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.badRequestResponse(w, r, fmt.Errorf("error decoding JSON: %w", err))
		return
	}
	if envelope.Data == nil {
		app.badRequestResponse(w, r, fmt.Errorf("data is required"))
		return
	}
	payload := *envelope.Data

	jobName := ps.ByName(ModelTransferJobName)
	if jobName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("job name is required"))
		return
	}
	modelRegistryID := ps.ByName(ModelRegistryId)
	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}
	deleteOldJob := r.URL.Query().Get("deleteOldJob") == "true"

	updatedJob, err := app.repositories.ModelRegistry.UpdateModelTransferJob(
		ctx, client, namespace, jobName, payload, deleteOldJob, modelRegistryID, app.config.DeploymentMode.IsFederatedMode(), app.podNamespace)
	if err != nil {
		if errors.Is(err, repositories.ErrJobNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		if errors.Is(err, repositories.ErrJobValidationFailed) {
			app.badRequestResponse(w, r, err)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	modelTransferJob := ModelTransferJobEnvelope{
		Data: updatedJob,
	}

	err = app.WriteJSON(w, http.StatusOK, modelTransferJob, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
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

	jobName := ps.ByName(ModelTransferJobName)
	if jobName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("job name is required"))
		return
	}

	jobNamespace, err := getRequiredJobNamespace(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	modelRegistryID := ps.ByName(ModelRegistryId)
	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	deletedJob, err := app.repositories.ModelRegistry.DeleteModelTransferJob(ctx, client, jobNamespace, jobName, modelRegistryID)
	if err != nil {
		if errors.Is(err, repositories.ErrJobNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	response := ModelTransferJobEnvelope{Data: deletedJob}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) GetModelTransferJobEventsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
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

	jobName := ps.ByName(ModelTransferJobName)
	if jobName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("job name is required"))
		return
	}

	modelRegistryID := ps.ByName(ModelRegistryId)
	if modelRegistryID == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	jobNamespace, err := getRequiredJobNamespace(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	events, err := app.repositories.ModelRegistry.GetModelTransferJobEvents(ctx, client, jobNamespace, jobName, modelRegistryID)
	if err != nil {
		if errors.Is(err, repositories.ErrJobNotFound) {
			app.notFoundResponse(w, r)
			return
		}
		app.serverErrorResponse(w, r, err)
		return
	}

	response := ModelTransferJobEventsEnvelope{
		Data: models.ModelTransferJobEventsResponse{Events: events},
	}

	err = app.WriteJSON(w, http.StatusOK, response, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}

func getRequiredJobNamespace(r *http.Request) (string, error) {
	jobNamespace := r.URL.Query().Get("jobNamespace")
	if jobNamespace == "" {
		return "", fmt.Errorf("missing required query parameter: jobNamespace")
	}
	return jobNamespace, nil
}
