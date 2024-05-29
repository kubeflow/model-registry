package api

import (
	"github.com/kubeflow/model-registry/ui/bff/config"
	"github.com/kubeflow/model-registry/ui/bff/data"
	"log/slog"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

const (
	// TODO(ederign) discuss versioning with the team
	Version         = "1.0.0"
	HealthCheckPath = "/api/v1/healthcheck/"
)

type App struct {
	config config.EnvConfig
	logger *slog.Logger
	models data.Models
}

func NewApp(cfg config.EnvConfig, logger *slog.Logger) *App {
	app := &App{
		config: cfg,
		logger: logger,
	}
	return app
}

func (app *App) Routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.GET(HealthCheckPath, app.HealthcheckHandler)

	return app.RecoverPanic(app.enableCORS(router))
}
