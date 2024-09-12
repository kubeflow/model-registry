package api

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/config"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
)

const (
	Version              = "1.0.0"
	PathPrefix           = "/api/v1"
	ModelRegistryId      = "model_registry_id"
	RegisteredModelId    = "registered_model_id"
	HealthCheckPath      = PathPrefix + "/healthcheck"
	ModelRegistry        = PathPrefix + "/model_registry"
	RegisteredModelsPath = ModelRegistry + "/:" + ModelRegistryId + "/registered_models"
	RegisteredModelPath  = RegisteredModelsPath + "/:" + RegisteredModelId
)

type App struct {
	config              config.EnvConfig
	logger              *slog.Logger
	models              data.Models
	kubernetesClient    integrations.KubernetesClientInterface
	modelRegistryClient data.ModelRegistryClientInterface
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	var k8sClient integrations.KubernetesClientInterface
	var err error
	if cfg.MockK8Client {
		//mock all k8s calls
		k8sClient, err = mocks.NewKubernetesClient(logger)
	} else {
		k8sClient, err = integrations.NewKubernetesClient(logger)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	var mrClient data.ModelRegistryClientInterface

	if cfg.MockMRClient {
		mrClient, err = mocks.NewModelRegistryClient(logger)
	} else {
		mrClient, err = data.NewModelRegistryClient(logger)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create ModelRegistry client: %w", err)
	}

	app := &App{
		config:              cfg,
		logger:              logger,
		kubernetesClient:    k8sClient,
		modelRegistryClient: mrClient,
	}
	return app, nil
}

func (app *App) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// HTTP client routes
	router.GET(HealthCheckPath, app.HealthcheckHandler)
	router.GET(RegisteredModelsPath, app.AttachRESTClient(app.GetAllRegisteredModelsHandler))
	router.GET(RegisteredModelPath, app.AttachRESTClient(app.GetRegisteredModelHandler))
	router.POST(RegisteredModelsPath, app.AttachRESTClient(app.CreateRegisteredModelHandler))

	// Kubernetes client routes
	router.GET(ModelRegistry, app.ModelRegistryHandler)

	return app.RecoverPanic(app.enableCORS(router))
}
