package api

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
)

type contextKey string

const (
	ModelRegistryHttpClientKey  contextKey = "ModelRegistryHttpClientKey"
	NamespaceHeaderParameterKey contextKey = "namespace"

	//Kubeflow authorization operates using custom authentication headers:
	// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
	// but it's supported on Model Registry BFF
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

		//skip use headers check if we are not on /api/v1
		if !strings.HasPrefix(r.URL.Path, PathPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		userIdHeader := r.Header.Get(KubeflowUserIDHeader)
		userGroupsHeader := r.Header.Get(KubeflowUserGroupsIdHeader)
		//`kubeflow-userid`: Contains the user's email address.
		if userIdHeader == "" {
			app.badRequestResponse(w, r, errors.New("missing required header: kubeflow-userid"))
			return
		}

		// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
		// but it's supported on Model Registry BFF
		//`kubeflow-groups`: Holds a comma-separated list of user groups.
		var userGroups []string
		if userGroupsHeader != "" {
			userGroups = strings.Split(userGroupsHeader, ",")
			// Trim spaces from each group name
			for i, group := range userGroups {
				userGroups[i] = strings.TrimSpace(group)
			}
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, KubeflowUserIdKey, userIdHeader)
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

func (app *App) AttachRESTClient(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		modelRegistryID := ps.ByName(ModelRegistryId)

		namespace, ok := r.Context().Value(NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		}

		modelRegistryBaseURL, err := resolveModelRegistryURL(namespace, modelRegistryID, app.kubernetesClient, app.config)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		client, err := integrations.NewHTTPClient(modelRegistryID, modelRegistryBaseURL)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), ModelRegistryHttpClientKey, client)
		next(w, r.WithContext(ctx), ps)
	}
}

func resolveModelRegistryURL(namespace string, serviceName string, client integrations.KubernetesClientInterface, config config.EnvConfig) (string, error) {

	serviceDetails, err := client.GetServiceDetailsByName(namespace, serviceName)
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

func (app *App) AttachNamespace(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		namespace := r.URL.Query().Get(string(NamespaceHeaderParameterKey))
		if namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing required query parameter: %s", NamespaceHeaderParameterKey))
			return
		}

		ctx := context.WithValue(r.Context(), NamespaceHeaderParameterKey, namespace)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

func (app *App) PerformSARonGetListServicesByNamespace(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, ok := r.Context().Value(KubeflowUserIdKey).(string)
		if !ok || user == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing user in context"))
			return
		}
		namespace, ok := r.Context().Value(NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		var userGroups []string
		if groups, ok := r.Context().Value(KubeflowUserGroupsKey).([]string); ok {
			userGroups = groups
		} else {
			userGroups = []string{}
		}

		allowed, err := app.kubernetesClient.PerformSARonGetListServicesByNamespace(user, userGroups, namespace)
		if err != nil {
			app.forbiddenResponse(w, r, fmt.Sprintf("failed to perform SAR: %v", err))
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r, "access denied")
			return
		}

		next(w, r, ps)
	}
}

func (app *App) PerformSARonSpecificService(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		user, ok := r.Context().Value(KubeflowUserIdKey).(string)
		if !ok || user == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing user in context"))
			return
		}

		namespace, ok := r.Context().Value(NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		modelRegistryID := ps.ByName(ModelRegistryId)
		if !ok || modelRegistryID == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		var userGroups []string
		if groups, ok := r.Context().Value(KubeflowUserGroupsKey).([]string); ok {
			userGroups = groups
		} else {
			userGroups = []string{}
		}

		allowed, err := app.kubernetesClient.PerformSARonSpecificService(user, userGroups, namespace, modelRegistryID)
		if err != nil {
			app.forbiddenResponse(w, r, "failed to perform SAR: %v")
			return
		}
		if !allowed {
			app.forbiddenResponse(w, r, "access denied")
			return
		}

		next(w, r, ps)
	}
}
