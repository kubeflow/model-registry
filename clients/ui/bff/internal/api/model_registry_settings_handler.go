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

func (app *App) GetAllModelRegistriesSettingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	registries := []models.ModelRegistryKind{createSampleModelRegistry("model-registry", namespace),
		createSampleModelRegistry("model-registry-dora", namespace),
		createSampleModelRegistry("model-registry-bella", namespace)}

	modelRegistryRes := ModelRegistrySettingsListEnvelope{
		Data: registries,
	}

	err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

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
	}

	modelId := ps.ByName(ModelRegistryId)
	registry := createSampleModelRegistry(modelId, namespace)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil)

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
	}

	var envelope ModelRegistrySettingsPayloadEnvelope
	if err := json.NewDecoder(r.Body).Decode(&envelope); err != nil {
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON:: %v", err.Error()))
		return
	}

	var modelRegistryName = envelope.Data.ModelRegistry.Metadata.Name

	if modelRegistryName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	ctxLogger.Info("Creating model registry", "name", modelRegistryName)

	// For now, we're using the stub implementation, but we'd use envelope.Data.ModelRegistry
	// and other fields from the payload in a real implementation
	registry := createSampleModelRegistry(modelRegistryName, namespace)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	w.Header().Set("Location", r.URL.JoinPath(modelRegistryRes.Data.Metadata.Name).String())
	writeErr := app.WriteJSON(w, http.StatusCreated, modelRegistryRes, nil)
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
	}

	modelId := ps.ByName(ModelRegistryId)
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

func (app *App) DeleteModelRegistrySettingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	w.WriteHeader(200)
}

// This function is a temporary function to create a sample model registry kind until we have a real implementation
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
