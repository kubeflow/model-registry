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
	Version = "1.0.0"

	PathPrefix                   = "/api/v1"
	ModelRegistryId              = "model_registry_id"
	RegisteredModelId            = "registered_model_id"
	ModelVersionId               = "model_version_id"
	ModelArtifactId              = "model_artifact_id"
	HealthCheckPath              = PathPrefix + "/healthcheck"
	UserPath                     = PathPrefix + "/user"
	ModelRegistryListPath        = PathPrefix + "/model_registry"
	NamespaceListPath            = PathPrefix + "/namespaces"
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

	// HTTP client routes (requests that we forward to Model Registry API)
	// on those, we perform SAR on Specific Service on a given namespace
	router.GET(HealthCheckPath, app.HealthcheckHandler)
	router.GET(RegisteredModelListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllRegisteredModelsHandler))))
	router.GET(RegisteredModelPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetRegisteredModelHandler))))
	router.POST(RegisteredModelListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateRegisteredModelHandler))))
	router.PATCH(RegisteredModelPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateRegisteredModelHandler))))
	router.GET(RegisteredModelVersionsPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllModelVersionsForRegisteredModelHandler))))
	router.POST(RegisteredModelVersionsPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelVersionForRegisteredModelHandler))))
	router.GET(ModelVersionPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient((app.GetModelVersionHandler)))))
	router.POST(ModelVersionListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelVersionHandler))))
	router.PATCH(ModelVersionPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateModelVersionHandler))))
	router.GET(ModelVersionArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllModelArtifactsByModelVersionHandler))))
	router.POST(ModelVersionArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelArtifactByModelVersionHandler))))
	router.PATCH(ModelRegistryPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateModelVersionHandler))))

	// Kubernetes routes
	router.GET(UserPath, app.UserHandler)
	// Perform SAR to Get List Services by Namspace
	router.GET(ModelRegistryListPath, app.AttachNamespace(app.PerformSARonGetListServicesByNamespace(app.ModelRegistryHandler)))
	if app.config.DevMode {
		router.GET(NamespaceListPath, app.AttachNamespace(app.GetNamespacesHandler))
	}

	return app.RecoverPanic(app.enableCORS(app.InjectUserHeaders(router)))
}
