package api

import (
	"context"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

type contextKey string

const httpClientKey contextKey = "httpClientKey"
const kubeflowUserId = "kubeflow-userid"

func (app *App) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) enableCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// TODO(ederign) restrict CORS to a much smaller set of trusted origins.
		// TODO(ederign) deal with preflight requests
		w.Header().Set("Access-Control-Allow-Origin", "*")

		next.ServeHTTP(w, r)
	})
}

func (app *App) AttachRESTClient(handler func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		modelRegistryID := ps.ByName(ModelRegistryId)

		modelRegistryBaseURL, err := resolveModelRegistryURL(modelRegistryID, app.kubernetesClient, app.config)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to resolve model registry base URL): %v", err))
			return
		}

		client, err := integrations.NewHTTPClient(modelRegistryBaseURL)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), httpClientKey, client)
		handler(w, r.WithContext(ctx), ps)
	}
}

func resolveModelRegistryURL(id string, client integrations.KubernetesClientInterface, config config.EnvConfig) (string, error) {
	serviceDetails, err := client.GetServiceDetailsByName(id)
	if err != nil {
		return "", err
	}

	if config.DevMode {
		serviceDetails.ClusterIP = "localhost"
		serviceDetails.HTTPPort = int32(config.DevModePort)
	}

	url := fmt.Sprintf("http://%s:%d/api/model_registry/v1alpha3", serviceDetails.ClusterIP, serviceDetails.HTTPPort)
	return url, nil
}

func (app *App) RequireAccessControl(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Skip SAR for health check
		if r.URL.Path == HealthCheckPath {
			next.ServeHTTP(w, r)
			return
		}

		user := r.Header.Get(kubeflowUserId)
		if user == "" {
			app.forbiddenResponse(w, r, "missing kubeflow-userid header")
			return
		}

		allowed, err := app.kubernetesClient.PerformSAR(user)
		if err != nil {
			app.forbiddenResponse(w, r, "failed to perform SAR: %v")
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r, "access denied")
			return
		}

		next.ServeHTTP(w, r)
	})
}
