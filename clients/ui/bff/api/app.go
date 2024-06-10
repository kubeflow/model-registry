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
	Version         = "1.0.0"
	HealthCheckPath = "/api/v1/healthcheck/"
	ModelRegistry   = "/api/v1/model-registry/"
)

type App struct {
	config           config.EnvConfig
	logger           *slog.Logger
	models           data.Models
	kubernetesClient integrations.KubernetesClientInterface
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) (*App, error) {
	k8sClient, err := integrations.NewKubernetesClient()
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

	// Kubernetes client routes
	router.GET(ModelRegistry, app.ModelRegistryHandler)

	return app.RecoverPanic(app.enableCORS(router))
}
