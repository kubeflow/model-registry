package api

import (
	"context"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
)

const (
	Version                      = "1.0.0"
	PathPrefix                   = "/api/v1"
	ModelRegistryId              = "model_registry_id"
	RegisteredModelId            = "registered_model_id"
	ModelVersionId               = "model_version_id"
	ModelArtifactId              = "model_artifact_id"
	HealthCheckPath              = PathPrefix + "/healthcheck"
	ModelRegistryListPath        = PathPrefix + "/model_registry"
	ModelRegistryPath            = ModelRegistryListPath + "/:" + ModelRegistryId
	RegisteredModelListPath      = ModelRegistryPath + "/registered_models"
	RegisteredModelPath          = RegisteredModelListPath + "/:" + RegisteredModelId
	RegisteredModelVersionsPath  = RegisteredModelPath + "/versions"
	ModelVersionListPath         = ModelRegistryPath + "/model_versions"
	ModelVersionPath             = ModelVersionListPath + "/:" + ModelVersionId
	ModelVersionArtifactListPath = ModelVersionPath + "/artifacts"
	ModelArtifactListPath        = ModelRegistryPath + "/model_artifacts"
	ModelArtifactPath            = ModelArtifactListPath + "/:" + ModelArtifactId
)

type App struct {
	config           config.EnvConfig
	logger           *slog.Logger
	kubernetesClient integrations.KubernetesClientInterface
	repositories     *repositories.Repositories
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	var k8sClient integrations.KubernetesClientInterface
	var err error
	if cfg.MockK8Client {
		//mock all k8s calls
		ctx, cancel := context.WithCancel(context.Background())
		k8sClient, err = mocks.NewKubernetesClient(logger, ctx, cancel)
	} else {
		k8sClient, err = integrations.NewKubernetesClient(logger)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	var mrClient repositories.ModelRegistryClientInterface

	if cfg.MockMRClient {
		//mock all model registry calls
		mrClient, err = mocks.NewModelRegistryClient(logger)
	} else {
		mrClient, err = repositories.NewModelRegistryClient(logger)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create ModelRegistryListPath client: %w", err)
	}

	app := &App{
		config:           cfg,
		logger:           logger,
		kubernetesClient: k8sClient,
		repositories:     repositories.NewRepositories(mrClient),
	}
	return app, nil
}

func (app *App) Shutdown(ctx context.Context, logger *slog.Logger) error {
	return app.kubernetesClient.Shutdown(ctx, logger)
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
	router.GET(ModelVersionArtifactListPath, app.AttachRESTClient(app.GetAllModelArtifactsByModelVersionHandler))
	router.POST(ModelVersionArtifactListPath, app.AttachRESTClient(app.CreateModelArtifactByModelVersionHandler))

	// Kubernetes client routes
	router.GET(ModelRegistryListPath, app.ModelRegistryHandler)
	router.PATCH(ModelRegistryPath, app.AttachRESTClient(app.UpdateModelVersionHandler))

	return app.RecoverPanic(app.enableCORS(router))
}
