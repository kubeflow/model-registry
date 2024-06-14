package api

import (
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/config"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	Version              = "1.0.0"
	PathPrefix           = "/api/v1"
	ModelRegistryId      = "model_registry_id"
	HealthCheckPath      = PathPrefix + "/healthcheck/"
	ModelRegistry        = PathPrefix + "/model-registry/"
	RegisteredModelsPath = ModelRegistry + ":" + ModelRegistryId + "/registered_models"
)

type App struct {
	config           config.EnvConfig
	logger           *slog.Logger
	models           data.Models
	kubernetesClient integrations.KubernetesClientInterface
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	k8sClient, err := integrations.NewKubernetesClient(logger)
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	app := &App{
		config:           cfg,
		logger:           logger,
		kubernetesClient: k8sClient,
	}
	return app, nil
}

func (app *App) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// HTTP client routes
	router.GET(HealthCheckPath, app.HealthcheckHandler)
	router.GET(RegisteredModelsPath, app.AttachRESTClient(app.GetRegisteredModelsHandler))
	router.POST(RegisteredModelsPath, app.AttachRESTClient(app.CreateRegisteredModelHandler))

	// Kubernetes client routes
	router.GET(ModelRegistry, app.ModelRegistryHandler)

	return app.RecoverPanic(app.enableCORS(router))
}
