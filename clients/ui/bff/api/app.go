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
	Version                     = "1.0.0"
	PathPrefix                  = "/api/v1"
	ModelRegistryId             = "model_registry_id"
	RegisteredModelId           = "registered_model_id"
	ModelVersionId              = "model_version_id"
	HealthCheckPath             = PathPrefix + "/healthcheck"
	ModelRegistryListPath       = PathPrefix + "/model_registry"
	ModelRegistryPath           = ModelRegistryListPath + "/:" + ModelRegistryId
	RegisteredModelListPath     = ModelRegistryPath + "/registered_models"
	RegisteredModelPath         = RegisteredModelListPath + "/:" + RegisteredModelId
	RegisteredModelVersionsPath = RegisteredModelPath + "/versions"
	ModelVersionListPath        = ModelRegistryPath + "/model_versions"
	ModelVersionPath            = ModelVersionListPath + "/:" + ModelVersionId
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
		return nil, fmt.Errorf("failed to create ModelRegistryListPath client: %w", err)
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
	router.GET(RegisteredModelListPath, app.AttachRESTClient(app.GetAllRegisteredModelsHandler))
	router.GET(RegisteredModelPath, app.AttachRESTClient(app.GetRegisteredModelHandler))
	router.POST(RegisteredModelListPath, app.AttachRESTClient(app.CreateRegisteredModelHandler))
	router.PATCH(RegisteredModelPath, app.AttachRESTClient(app.UpdateRegisteredModelHandler))
	router.GET(RegisteredModelVersionsPath, app.AttachRESTClient(app.GetAllModelVersionsForRegisteredModelHandler))
	router.POST(RegisteredModelVersionsPath, app.AttachRESTClient(app.CreateModelVersionForRegisteredModelHandler))

	router.GET(ModelVersionPath, app.AttachRESTClient(app.GetModelVersionHandler))
	router.POST(ModelVersionListPath, app.AttachRESTClient(app.CreateModelVersionHandler))
	router.PATCH(ModelVersionPath, app.AttachRESTClient(app.UpdateModelVersionHandler))

	// Kubernetes client routes
	router.GET(ModelRegistryListPath, app.ModelRegistryHandler)
	router.PATCH(ModelRegistryPath, app.AttachRESTClient(app.UpdateModelVersionHandler))

	return app.RecoverPanic(app.enableCORS(router))
}
