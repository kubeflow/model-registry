package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"path"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	k8mocks "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes/k8mocks"
	"k8s.io/client-go/kubernetes"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
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
	CertificatesPath              = SettingsPath + "/certificates"
	RoleBindingListPath           = SettingsPath + "/role_bindings"
	GroupsPath                    = SettingsPath + "/groups"
	SettingsNamespacePath         = SettingsPath + "/namespaces"
	RoleBindingPath               = RoleBindingListPath + "/:" + RoleBindingNameParam

	RegisteredModelListPath      = ModelRegistryPath + "/registered_models"
	RegisteredModelPath          = RegisteredModelListPath + "/:" + RegisteredModelId
	RegisteredModelVersionsPath  = RegisteredModelPath + "/versions"
	ModelVersionListPath         = ModelRegistryPath + "/model_versions"
	ModelVersionPath             = ModelVersionListPath + "/:" + ModelVersionId
	ModelVersionArtifactListPath = ModelVersionPath + "/artifacts"
	ModelArtifactListPath        = ModelRegistryPath + "/model_artifacts"
	ModelArtifactPath            = ModelArtifactListPath + "/:" + ModelArtifactId
	ArtifactListPath             = ModelRegistryPath + "/artifacts"
	ArtifactPath                 = ArtifactListPath + "/:" + ArtifactId

	// model catalog
	SourceId              = "source_id"
	CatalogModelName      = "model_name"
	CatalogPathPrefix     = ApiPathPrefix + "/model_catalog"
	CatalogModelListPath  = CatalogPathPrefix + "/models"
	CatalogSourceListPath = CatalogPathPrefix + "/sources"
	CatalogModelPath      = CatalogPathPrefix + "/sources" + "/:" + SourceId + "/models" + "/*model_name"
)

type App struct {
	config                  config.EnvConfig
	logger                  *slog.Logger
	kubernetesClientFactory k8s.KubernetesClientFactory
	repositories            *repositories.Repositories
	//used only on mocked k8s client
	testEnv *envtest.Environment
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	logger.Debug("Initializing app with config", slog.Any("config", cfg))
	var k8sFactory k8s.KubernetesClientFactory
	var err error
	// used only on mocked k8s client
	var testEnv *envtest.Environment

	if cfg.MockK8Client {
		//mock all k8s calls with 'env test'
		var clientset kubernetes.Interface
		ctx, cancel := context.WithCancel(context.Background())
		testEnv, clientset, err = k8mocks.SetupEnvTest(k8mocks.TestEnvInput{
			Logger: logger,
			Ctx:    ctx,
			Cancel: cancel,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to setup envtest: %w", err)
		}
		//create mocked kubernetes client factory
		k8sFactory, err = k8mocks.NewMockedKubernetesClientFactory(clientset, testEnv, cfg, logger)

	} else {
		//create kubernetes client factory
		k8sFactory, err = k8s.NewKubernetesClientFactory(cfg, logger)
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
		config:                  cfg,
		logger:                  logger,
		kubernetesClientFactory: k8sFactory,
		repositories:            repositories.NewRepositories(mrClient),
		testEnv:                 testEnv,
	}
	return app, nil
}

func (app *App) Shutdown() error {
	app.logger.Info("shutting down app...")
	if app.testEnv == nil {
		return nil
	}
	//shutdown the envtest control plane when we are in the mock mode.
	app.logger.Info("shutting env test...")
	return app.testEnv.Stop()
}

func (app *App) Routes() http.Handler {
	// Router for /api/v1/*
	apiRouter := httprouter.New()

	apiRouter.NotFound = http.HandlerFunc(app.notFoundResponse)
	apiRouter.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	// HTTP client routes (requests that we forward to Model Registry API)
	// on those, we perform SAR or SSAR on Specific Service on a given namespace
	apiRouter.GET(RegisteredModelListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetAllRegisteredModelsHandler))))
	apiRouter.GET(RegisteredModelPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetRegisteredModelHandler))))
	apiRouter.POST(RegisteredModelListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.CreateRegisteredModelHandler))))
	apiRouter.PATCH(RegisteredModelPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.UpdateRegisteredModelHandler))))
	apiRouter.GET(RegisteredModelVersionsPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetAllModelVersionsForRegisteredModelHandler))))
	apiRouter.POST(RegisteredModelVersionsPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.CreateModelVersionForRegisteredModelHandler))))
	apiRouter.POST(ModelVersionListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.CreateModelVersionHandler))))
	apiRouter.GET(ModelVersionListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetAllModelVersionHandler))))
	apiRouter.GET(ModelVersionPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetModelVersionHandler))))
	apiRouter.PATCH(ModelVersionPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.UpdateModelVersionHandler))))
	apiRouter.GET(ArtifactListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetAllArtifactsHandler))))
	apiRouter.GET(ArtifactPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetArtifactHandler))))
	apiRouter.POST(ArtifactListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.CreateArtifactHandler))))
	apiRouter.GET(ModelVersionArtifactListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.GetAllModelArtifactsByModelVersionHandler))))
	apiRouter.POST(ModelVersionArtifactListPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.CreateModelArtifactByModelVersionHandler))))
	apiRouter.PATCH(ModelRegistryPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.UpdateModelVersionHandler))))
	apiRouter.PATCH(ModelArtifactPath, app.AttachNamespace(app.RequireAccessToService(app.AttachRESTClient(app.UpdateModelArtifactHandler))))

	// Kubernetes routes
	apiRouter.GET(UserPath, app.UserHandler)
	apiRouter.GET(ModelRegistryListPath, app.AttachNamespace(app.RequireListServiceAccessInNamespace(app.GetAllModelRegistriesHandler)))

	// Enable these routes in all cases except Kubeflow integration mode
	// (Kubeflow integration mode is when DeploymentMode is kubeflow)
	isKubeflowIntegrationMode := app.config.DeploymentMode.IsKubeflowMode()
	if !isKubeflowIntegrationMode {
		// This namespace endpoint is used on standalone mode to simulate
		// Kubeflow Central Dashboard namespace selector dropdown on our standalone web app
		apiRouter.GET(NamespaceListPath, app.GetNamespacesHandler)

		// SettingsPath endpoints are used to manage the model registry settings and create new model registries
		// We are still discussing the best way to create model registries in the community
		// But in the meantime, those endpoints are STUBs endpoints used to unblock the frontend development
		apiRouter.GET(ModelRegistrySettingsListPath, app.AttachNamespace(app.GetAllModelRegistriesSettingsHandler))
		apiRouter.POST(ModelRegistrySettingsListPath, app.AttachNamespace(app.CreateModelRegistrySettingsHandler))
		apiRouter.GET(ModelRegistrySettingsPath, app.AttachNamespace(app.GetModelRegistrySettingsHandler))
		apiRouter.PATCH(ModelRegistrySettingsPath, app.AttachNamespace(app.UpdateModelRegistrySettingsHandler))
		apiRouter.DELETE(ModelRegistrySettingsPath, app.AttachNamespace(app.DeleteModelRegistrySettingsHandler))

		//SettingsPath: Certificate endpoints
		apiRouter.GET(CertificatesPath, app.AttachNamespace(app.GetCertificatesHandler))

		//SettingsPath: Role Binding endpoints
		apiRouter.GET(RoleBindingListPath, app.AttachNamespace(app.GetRoleBindingsHandler))
		apiRouter.POST(RoleBindingListPath, app.AttachNamespace(app.CreateRoleBindingHandler))
		apiRouter.PATCH(RoleBindingPath, app.AttachNamespace(app.PatchRoleBindingHandler))
		apiRouter.DELETE(RoleBindingPath, app.AttachNamespace(app.DeleteRoleBindingHandler))

		//SettingsPath Groups endpoints
		apiRouter.GET(GroupsPath, app.GetGroupsHandler)

		//SettingsPath Namespace endpoints
		//This namespace endpoint is used to get the namespaces for the current user inside the model registry settings
		apiRouter.GET(SettingsNamespacePath, app.GetNamespacesHandler)

		// Model catalog endpoints
		apiRouter.GET(CatalogModelListPath, app.AttachNamespace((app.GetAllCatalogModelsHandler)))
		apiRouter.GET(CatalogSourceListPath, app.AttachNamespace((app.GetAllCatalogSourcesHandler)))
		apiRouter.GET(CatalogModelPath, app.AttachNamespace((app.GetCatalogModelHandler)))

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
	combinedMux.Handle("/", app.RecoverPanic(app.EnableTelemetry(app.EnableCORS(app.InjectRequestIdentity(appMux)))))

	return combinedMux
}
