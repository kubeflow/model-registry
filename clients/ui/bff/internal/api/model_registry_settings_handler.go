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
type ModelRegistryAndCredentialsSettingsEnvelope Envelope[models.ModelRegistryAndCredentials, None]
type ModelRegistrySettingsPayloadEnvelope Envelope[models.ModelRegistrySettingsPayload, None]

func (app *App) GetAllModelRegistriesSettingsHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctxLogger := helper.GetContextLoggerFromReq(r)
	ctxLogger.Info("This functionality is not implement yet. This is a STUB API to unblock frontend development")

	namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
	if !ok || namespace == "" {
		app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
	}

	sslRootCertificateConfigMap := models.Entry{
		Name: "ssl-config-map-name",
		Key:  "ssl-config-map-key",
	}

	sslRootCertificateSecret := models.Entry{
		Name: "ssl-secret-name",
		Key:  "ssl-secret-key",
	}

	lastTransitionTime, _ := time.Parse(time.RFC3339, "2024-03-22T09:30:02Z")
	deletionTime := lastTransitionTime
	registries := []models.ModelRegistryKind{
		createSampleModelRegistry("model-registry", namespace, &sslRootCertificateConfigMap, nil, nil, nil),
		createSampleModelRegistry("model-registry-dora", namespace, nil, &sslRootCertificateSecret, nil, nil),
		createSampleModelRegistry("model-registry-bella", namespace, nil, nil, nil, nil),
		createSampleModelRegistry("model-registry-ready", namespace, nil, nil, nil,
			[]models.Condition{
				{LastTransitionTime: lastTransitionTime, Message: "Deployment for custom resource model-registry-ready is available", Reason: "Available", Status: "True", Type: "Available"},
			}),
		createSampleModelRegistry("model-registry-starting", namespace, nil, nil, nil,
			[]models.Condition{
				{LastTransitionTime: lastTransitionTime, Message: "Deployment for custom resource model-registry-starting was successfully created", Reason: "CreatedDeployment", Status: "True", Type: "Progressing"},
			}),
		createSampleModelRegistry("model-registry-stopping", namespace, nil, nil, &deletionTime,
			[]models.Condition{
				{LastTransitionTime: lastTransitionTime, Message: "Deployment for custom resource model-registry-stopping is available", Reason: "Available", Status: "True", Type: "Available"},
			}),
		createSampleModelRegistry("model-registry-degrading", namespace, nil, nil, nil,
			[]models.Condition{
				{LastTransitionTime: lastTransitionTime, Message: "Service is degrading", Reason: "Degraded", Status: "True", Type: "Degraded"},
			}),
		createSampleModelRegistry("model-registry-unavailable", namespace, nil, nil, nil,
			[]models.Condition{
				{LastTransitionTime: lastTransitionTime, Message: "Service is unavailable", Reason: "Unavailable", Status: "False", Type: "Available"},
				{LastTransitionTime: lastTransitionTime, Message: "Istio resources are unavailable", Reason: "IstioUnavailable", Status: "False", Type: "IstioAvailable"},
			}),
	}

	modelRegistryRes := ModelRegistrySettingsListEnvelope{
		Data: registries,
	}

	if err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil); err != nil {
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
	registry := createSampleModelRegistry(modelId, namespace, nil, nil, nil, nil)

	modelRegistryWithCreds := models.ModelRegistryAndCredentials{
		ModelRegistry:    registry,
		DatabasePassword: "dbpassword",
	}

	modelRegistryRes := ModelRegistryAndCredentialsSettingsEnvelope{
		Data: modelRegistryWithCreds,
	}

	if err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil); err != nil {
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
		app.serverErrorResponse(w, r, fmt.Errorf("error decoding JSON: %v", err.Error()))
		return
	}

	modelRegistryName := envelope.Data.ModelRegistry.Metadata.Name

	if modelRegistryName == "" {
		app.badRequestResponse(w, r, fmt.Errorf("model registry name is required"))
		return
	}

	ctxLogger.Info("Creating model registry", "name", modelRegistryName)

	registry := createSampleModelRegistry(modelRegistryName, namespace, nil, nil, nil, nil)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	w.Header().Set("Location", r.URL.JoinPath(modelRegistryRes.Data.Metadata.Name).String())
	if err := app.WriteJSON(w, http.StatusCreated, modelRegistryRes, nil); err != nil {
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
	registry := createSampleModelRegistry(modelId, namespace, nil, nil, nil, nil)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	if err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil); err != nil {
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
	}

	// This is a temporary fix to handle frontend error (as it is expecting ModelRegistryKind response) until we have a real implementation
	modelId := ps.ByName(ModelRegistryId)
	registry := createSampleModelRegistry(modelId, namespace, nil, nil, nil, nil)

	modelRegistryRes := ModelRegistrySettingsEnvelope{
		Data: registry,
	}

	if err := app.WriteJSON(w, http.StatusOK, modelRegistryRes, nil); err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

// This function is a temporary function to create a sample model registry kind until we have a real implementation
func createSampleModelRegistry(name string, namespace string, SSLRootCertificateConfigMap *models.Entry, SSLRootCertificateSecret *models.Entry, deletionTimestamp *time.Time, conditions []models.Condition) models.ModelRegistryKind {

	creationTime, _ := time.Parse(time.RFC3339, "2024-03-14T08:01:42Z")

	return models.ModelRegistryKind{
		APIVersion: "modelregistry.io/v1alpha1",
		Kind:       "ModelRegistry",
		Metadata: models.Metadata{
			Name:              name,
			Namespace:         namespace,
			CreationTimestamp: creationTime,
			DeletionTimestamp: deletionTimestamp,
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
				SSLRootCertificateConfigMap: SSLRootCertificateConfigMap,
				SSLRootCertificateSecret:    SSLRootCertificateSecret,
			},
		},
		Status: models.Status{
			Conditions: conditions,
		},
	}
}
