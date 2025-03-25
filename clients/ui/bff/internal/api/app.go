package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"

	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
)

const (
	Version = "1.0.0"

	PathPrefix                    = "/model-registry"
	ApiPathPrefix                 = "/api/v1"
	ModelRegistryId               = "model_registry_id"
	RegisteredModelId             = "registered_model_id"
	ModelVersionId                = "model_version_id"
	ModelArtifactId               = "model_artifact_id"
	ArtifactId                    = "artifact_id"
	HealthCheckPath               = "/healthcheck"
	UserPath                      = ApiPathPrefix + "/user"
	ModelRegistryListPath         = ApiPathPrefix + "/model_registry"
	ModelRegistryPath             = ModelRegistryListPath + "/:" + ModelRegistryId
	NamespaceListPath             = ApiPathPrefix + "/namespaces"
	SettingsPath                  = ApiPathPrefix + "/settings"
	ModelRegistrySettingsListPath = SettingsPath + "/model_registry"
	ModelRegistrySettingsPath     = ModelRegistrySettingsListPath + "/:" + ModelRegistryId
	RegisteredModelListPath       = ModelRegistryPath + "/registered_models"
	RegisteredModelPath           = RegisteredModelListPath + "/:" + RegisteredModelId
	RegisteredModelVersionsPath   = RegisteredModelPath + "/versions"
	ModelVersionListPath          = ModelRegistryPath + "/model_versions"
	ModelVersionPath              = ModelVersionListPath + "/:" + ModelVersionId
	ModelVersionArtifactListPath  = ModelVersionPath + "/artifacts"
	ModelArtifactListPath         = ModelRegistryPath + "/model_artifacts"
	ModelArtifactPath             = ModelArtifactListPath + "/:" + ModelArtifactId
	ArtifactListPath              = ModelRegistryPath + "/artifacts"
	ArtifactPath                  = ArtifactListPath + "/:" + ArtifactId
)

type App struct {
	config           config.EnvConfig
	logger           *slog.Logger
	kubernetesClient integrations.KubernetesClientInterface
	repositories     *repositories.Repositories
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	logger.Debug("Initializing app with config", slog.Any("config", cfg))
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
	// Router for /api/v1/*
	apiRouter := httprouter.New()

	apiRouter.NotFound = http.HandlerFunc(app.notFoundResponse)
	apiRouter.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// HTTP client routes (requests that we forward to Model Registry API)
	// on those, we perform SAR on Specific Service on a given namespace
	apiRouter.GET(RegisteredModelListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllRegisteredModelsHandler))))
	apiRouter.GET(RegisteredModelPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetRegisteredModelHandler))))
	apiRouter.POST(RegisteredModelListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateRegisteredModelHandler))))
	apiRouter.PATCH(RegisteredModelPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateRegisteredModelHandler))))
	apiRouter.GET(RegisteredModelVersionsPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllModelVersionsForRegisteredModelHandler))))
	apiRouter.POST(RegisteredModelVersionsPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelVersionForRegisteredModelHandler))))
	apiRouter.POST(ModelVersionListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelVersionHandler))))
	apiRouter.GET(ModelVersionListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllModelVersionHandler))))
	apiRouter.GET(ModelVersionPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetModelVersionHandler))))
	apiRouter.PATCH(ModelVersionPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateModelVersionHandler))))
	apiRouter.GET(ArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllArtifactsHandler))))
	apiRouter.GET(ArtifactPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetArtifactHandler))))
	apiRouter.POST(ArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateArtifactHandler))))
	apiRouter.GET(ModelVersionArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.GetAllModelArtifactsByModelVersionHandler))))
	apiRouter.POST(ModelVersionArtifactListPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.CreateModelArtifactByModelVersionHandler))))
	apiRouter.PATCH(ModelRegistryPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateModelVersionHandler))))
	apiRouter.PATCH(ModelArtifactPath, app.AttachNamespace(app.PerformSARonSpecificService(app.AttachRESTClient(app.UpdateModelArtifactHandler))))

	// Kubernetes routes
	apiRouter.GET(UserPath, app.UserHandler)
	apiRouter.GET(ModelRegistryListPath, app.AttachNamespace(app.PerformSARonGetListServicesByNamespace(app.GetAllModelRegistriesHandler)))
	if app.config.StandaloneMode {
		apiRouter.GET(NamespaceListPath, app.GetNamespacesHandler)
		//Those endpoints are not implement yet. This is a STUB API to unblock frontend development
		apiRouter.GET(ModelRegistrySettingsListPath, app.AttachNamespace(app.GetAllModelRegistriesSettingsHandler))
		apiRouter.POST(ModelRegistrySettingsListPath, app.AttachNamespace(app.CreateModelRegistrySettingsHandler))
		apiRouter.GET(ModelRegistrySettingsPath, app.AttachNamespace(app.GetModelRegistrySettingsHandler))
		apiRouter.PATCH(ModelRegistrySettingsPath, app.AttachNamespace(app.UpdateModelRegistrySettingsHandler))
		apiRouter.DELETE(ModelRegistrySettingsPath, app.AttachNamespace(app.DeleteModelRegistrySettingsHandler))
	}

	// App Router
	appMux := http.NewServeMux()

	// handler for api calls
	appMux.Handle(ApiPathPrefix+"/", apiRouter)
	appMux.Handle(PathPrefix+ApiPathPrefix+"/", http.StripPrefix(PathPrefix, apiRouter))

	// file server for the frontend file and SPA routes
	staticDir := http.Dir(app.config.StaticAssetsDir)
	fileServer := http.FileServer(staticDir)
	appMux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctxLogger := helper.GetContextLoggerFromReq(r)
		// Check if the requested file exists
		if _, err := staticDir.Open(r.URL.Path); err == nil {
			ctxLogger.Debug("Serving static file", slog.String("path", r.URL.Path))
			// Serve the file if it exists
			fileServer.ServeHTTP(w, r)
			return
		}

		// Fallback to index.html for SPA routes
		ctxLogger.Debug("Static asset not found, serving index.html", slog.String("path", r.URL.Path))
		http.ServeFile(w, r, path.Join(app.config.StaticAssetsDir, "index.html"))
	})

	// Create a mux for the healthcheck endpoint
	healthcheckMux := http.NewServeMux()
	healthcheckRouter := httprouter.New()
	healthcheckRouter.GET(HealthCheckPath, app.HealthcheckHandler)
	healthcheckMux.Handle(HealthCheckPath, app.RecoverPanic(app.EnableTelemetry(healthcheckRouter)))

	// Combines the healthcheck endpoint with the rest of the routes
	combinedMux := http.NewServeMux()
	combinedMux.Handle(HealthCheckPath, healthcheckMux)
	combinedMux.Handle("/", app.RecoverPanic(app.EnableTelemetry(app.EnableCORS(app.InjectUserHeaders(appMux)))))

	return combinedMux
}
