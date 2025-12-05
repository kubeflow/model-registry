package api

import (
	"fmt"
	"net/http"

	"log/slog"

	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
)

// BadRequest exposes the internal bad request helper for extensions.
func (app *App) BadRequest(w http.ResponseWriter, r *http.Request, err error) {
	if app == nil {
		return
	}
	app.badRequestResponse(w, r, err)
}

// ServerError exposes the internal server error helper for extensions.
func (app *App) ServerError(w http.ResponseWriter, r *http.Request, err error) {
	if app == nil {
		return
	}
	app.serverErrorResponse(w, r, err)
}

// NotImplemented writes a standard placeholder response for unimplemented endpoints.
func (app *App) NotImplemented(w http.ResponseWriter, r *http.Request, feature string) {
	app.serverErrorResponse(w, r, fmt.Errorf("%s is not implemented", feature))
}

// Config exposes the application configuration for extensions.
func (app *App) Config() config.EnvConfig {
	return app.config
}

// Logger exposes the application logger for extensions.
func (app *App) Logger() *slog.Logger {
	return app.logger
}

// KubernetesClientFactory exposes the k8s factory for extensions.
func (app *App) KubernetesClientFactory() k8s.KubernetesClientFactory {
	return app.kubernetesClientFactory
}

// Repositories exposes the repositories container for extensions.
func (app *App) Repositories() *repositories.Repositories {
	return app.repositories
}
