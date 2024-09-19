package api

import (
	"context"
	"fmt"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"k8s.io/client-go/rest"
	"net/http"
)

type contextKey string

const httpClientKey contextKey = "httpClientKey"

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

		modelRegistryBaseURL, err := resolveModelRegistryURL(modelRegistryID, app.kubernetesClient)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to resolve model registry base URL): %v", err))
			return
		}
		var bearerToken string
		bearerToken, err = resolveBearerToken(app.kubernetesClient)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to resolve BearerToken): %v", err))
			return
		}

		client, err := integrations.NewHTTPClient(modelRegistryBaseURL, bearerToken)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), httpClientKey, client)
		handler(w, r.WithContext(ctx), ps)
	}
}

func resolveBearerToken(k8s integrations.KubernetesClientInterface) (string, error) {
	var bearerToken string
	_, err := rest.InClusterConfig()
	if err == nil {
		//in cluster
		//TODO (eder) load bearerToken probably from x-forwarded-access-bearerToken
		return "", fmt.Errorf("failed to create Rest client (not implemented yet - inside cluster): %v", err)
	} else {
		//off cluster (development)
		bearerToken, err = k8s.BearerToken()
		if err != nil {
			return "", fmt.Errorf("failed to fetch BearerToken in development mode: %v", err)
		}
	}
	return bearerToken, err
}

func resolveModelRegistryURL(id string, client integrations.KubernetesClientInterface) (string, error) {
	serviceDetails, err := client.GetServiceDetailsByName(id)
	if err != nil {
		return "", err
	}
	url := fmt.Sprintf("http://%s:%d/api/model_registry/v1alpha3", serviceDetails.ClusterIP, serviceDetails.HTTPPort)
	return url, nil
}
