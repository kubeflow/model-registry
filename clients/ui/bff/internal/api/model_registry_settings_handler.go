package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
)

type ModelRegistrySettingsListEnvelope Envelope[[]models.ModelRegistryKind, None]
type ModelRegistrySettingsEnvelope Envelope[models.ModelRegistryKind, None]
type ModelRegistrySettingsPayloadEnvelope Envelope[models.ModelRegistrySettingsPayload, None]
type ModelRegistrySettingsUpdateEnvelope Envelope[models.ModelRegistryKind, None]

func (app *App) GetAllModelRegistriesSettingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	labelSelector, ok := r.Context().Value(constants.LabelSelectorHeaderParameterKey).(string)
	if !ok {
		labelSelector = ""
	}

	client, err := app.kubernetesClientFactory.GetClient(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return
	}

	registries, err := app.repositories.ModelRegistrySettings.GetAllModelRegistriesSettings(r.Context(), client, namespace, labelSelector)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := ModelRegistrySettingsListEnvelope{
		Data: registries,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *App) GetModelRegistrySettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	modelRegistryName := ps.ByName(ModelRegistryId)

	client, err := app.kubernetesClientFactory.GetClient(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return
	}

	modelRegistry, err := app.repositories.ModelRegistrySettings.GetModelRegistrySettings(r.Context(), client, namespace, modelRegistryName)
	if err != nil {
		ctxLogger.Error("Failed to fetch model registry settings", "name", modelRegistryName, "error", err)
		app.serverErrorResponse(w, r, err)
		return
	}

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: modelRegistry,
	}

	err = app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *App) CreateModelRegistrySettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	var envelope ModelRegistrySettingsPayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}
	model := envelope.Data.ModelRegistry
	dbPassword := envelope.Data.DatabasePassword

	dryRun := false
	if dr := r.URL.Query().Get("dryRun"); dr == "true" {
		dryRun = true
	}

	client, err := app.kubernetesClientFactory.GetClient(r.Context())
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
		return
	}
	created, err := app.repositories.ModelRegistrySettings.CreateModelRegistryKindWithSecret(r.Context(), client, namespace, model, dbPassword, dryRun)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	w.Header().Set("Location", r.URL.JoinPath(created.Metadata.Name).String())
	resp := ModelRegistrySettingsEnvelope{Data: created}
	writeErr := app.WriteJSON(w, http.StatusCreated, resp, nil)
	if writeErr != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) UpdateModelRegistrySettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	// Read request body for update data (following standard UPDATE pattern)
	var envelope ModelRegistrySettingsUpdateEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	modelId := ps.ByName(ModelRegistryId)

	// TODO: Implement actual update logic here
	// For now, return mock data but following proper pattern
	// data := envelope.Data
	// client, err := app.kubernetesClientFactory.GetClient(r.Context())
	// updatedRegistry, err := app.repositories.ModelRegistrySettings.UpdateModelRegistrySettings(r.Context(), client, namespace, modelId, data)

	// STUB: Return sample data for now
	registry := createSampleModelRegistry(modelId, namespace)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)
	if err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error writing JSON"))
		return
	}
}

func (app *App) DeleteModelRegistrySettingsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		return
	}

	modelId := ps.ByName(ModelRegistryId)

	// TODO: Implement actual delete logic here
	// client, err := app.kubernetesClientFactory.GetClient(r.Context())
	// err := app.repositories.ModelRegistrySettings.DeleteModelRegistrySettings(r.Context(), client, namespace, modelId)
	// if err != nil {
	//     app.serverErrorResponse(w, r, err)
	//     return
	// }

	ctxLogger.Info("STUB: Deleting Model Registry Settings", "name", modelId, "namespace", namespace)

	// Standard response for successful DELETE (no response body)
	w.WriteHeader(http.StatusNoContent)
}

// TODO: delete this (move to shared client on mocking for now)
func createSampleModelRegistry(name string, namespace string) models.ModelRegistryKind {

	creationTime, _ := time.Parse(time.RFC3339, "2024-03-14T08:01:42Z")
	lastTransitionTime, _ := time.Parse(time.RFC3339, "2024-03-22T09:30:02Z")

	return models.ModelRegistryKind{
		APIVersion: "modelregistry.io/v1alpha1",
		Kind:       "ModelRegistry",
		Metadata: models.Metadata{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: creationTime,
			Annotations:       map[string]string{},
		},
		Spec: models.ModelRegistrySpec{
			GRPC: models.EmptyObject{},
			REST: models.EmptyObject{},
			Istio: models.IstioConfig{
				Gateway: models.GatewayConfig{
					GRPC: models.GRPCConfig{
						TLS: models.EmptyObject{},
					},
					REST: models.RESTConfig{
						TLS: models.EmptyObject{},
					},
				},
			},
			DatabaseConfig: models.DatabaseConfig{
				DatabaseType: models.MySQL,
				Database:     "model-registry",
				Host:         "model-registry-db",
				//intentionally not set
				// PasswordSecret: models.PasswordSecret{
				// 	Key:  "database-password",
				// 	Name: "model-registry-db",
				// },
				Port:                        5432,
				SkipDBCreation:              false,
				Username:                    "mlmduser",
				SSLRootCertificateConfigMap: "ssl-config-map",
				SSLRootCertificateSecret:    "ssl-secret",
			},
		},
		Status: models.Status{
			Conditions: []models.Condition{
				{
					LastTransitionTime: lastTransitionTime,
					Message:            "Deployment for custom resource " + name + " was successfully created",
					Reason:             "CreatedDeployment",
					Status:             "True",
					Type:               "Progressing",
				},
			},
		},
	}
}
