package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

type contextKey string

const (
	httpClientKey contextKey = "httpClientKey"

	//Kubeflow authorization operates using custom authentication headers:
	// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
	// But it will be soon implemented on Model Registry BFF
	KubeflowUserIdKey          contextKey = "kubeflowUserId" // kubeflow-userid :contains the user's email address
	KubeflowUserIDHeader                  = "kubeflow-userid"
	KubeflowUserGroupsKey      contextKey = "kubeflowUserGroups" // kubeflow-groups : Holds a comma-separated list of user groups
	KubeflowUserGroupsIdHeader            = "kubeflow-groups"
)

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

func (app *App) InjectUserHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		userId := r.Header.Get(KubeflowUserIDHeader)
		userGroups := r.Header.Get(KubeflowUserGroupsIdHeader)

		//Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
		if userId == "" {
			app.badRequestResponse(w, r, errors.New("missing required header: kubeflow-userid"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, KubeflowUserIdKey, userId)
		ctx = context.WithValue(ctx, KubeflowUserGroupsKey, userGroups)

		next.ServeHTTP(w, r.WithContext(ctx))
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

func (app *App) RequireAccessControl(next http.Handler, exemptPaths map[string]struct{}) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Skip SAR for exempt paths
		if _, exempt := exemptPaths[r.URL.Path]; exempt {
			next.ServeHTTP(w, r)
			return
		}

		user := r.Header.Get(KubeflowUserIDHeader)
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
