package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/rs/cors"
)

func (app *App) RecoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
				app.logger.Error("Recovered from panic", slog.String("stack_trace", string(debug.Stack())))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *App) InjectUserHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//skip use headers check if we are not on /api/v1
		if !strings.HasPrefix(r.URL.Path, ApiPathPrefix) && !strings.HasPrefix(r.URL.Path, PathPrefix+ApiPathPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		userIdHeader := r.Header.Get(constants.KubeflowUserIDHeader)
		userGroupsHeader := r.Header.Get(constants.KubeflowUserGroupsIdHeader)
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
		ctx = context.WithValue(ctx, constants.KubeflowUserIdKey, userIdHeader)
		ctx = context.WithValue(ctx, constants.KubeflowUserGroupsKey, userGroups)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) EnableCORS(next http.Handler) http.Handler {
	if len(app.config.AllowedOrigins) == 0 {
		// CORS is disabled, this middleware becomes a noop.
		return next
	}

	c := cors.New(cors.Options{
		AllowedOrigins:     app.config.AllowedOrigins,
		AllowCredentials:   true,
		AllowedMethods:     []string{"GET", "PUT", "POST", "PATCH", "DELETE"},
		AllowedHeaders:     []string{constants.KubeflowUserIDHeader, constants.KubeflowUserGroupsIdHeader},
		Debug:              app.config.LogLevel == slog.LevelDebug,
		OptionsPassthrough: false,
	})

	return c.Handler(next)
}

func (app *App) EnableTelemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Adds a unique id to the context to allow tracing of requests
		traceId := uuid.NewString()
		ctx := context.WithValue(r.Context(), constants.TraceIdKey, traceId)

		// logger will only be nil in tests.
		if app.logger != nil {
			traceLogger := app.logger.With(slog.String("trace_id", traceId))
			ctx = context.WithValue(ctx, constants.TraceLoggerKey, traceLogger)

			traceLogger.Debug("Incoming HTTP request", slog.Any("request", helper.RequestLogValuer{Request: r}))
		}
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (app *App) AttachRESTClient(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		modelRegistryID := ps.ByName(ModelRegistryId)

		namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in the context"))
		}

		modelRegistryBaseURL, err := resolveModelRegistryURL(r.Context(), namespace, modelRegistryID, app.kubernetesClient, app.config)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}

		// Set up a child logger for the rest client that automatically adds the request id to all statements for
		// tracing.
		restClientLogger := app.logger
		traceId, ok := r.Context().Value(constants.TraceIdKey).(string)
		if app.logger != nil {
			if ok {
				restClientLogger = app.logger.With(slog.String("trace_id", traceId))
			} else {
				app.logger.Warn("Failed to set trace_id for tracing")
			}
		}

		client, err := integrations.NewHTTPClient(restClientLogger, modelRegistryID, modelRegistryBaseURL)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), constants.ModelRegistryHttpClientKey, client)
		next(w, r.WithContext(ctx), ps)
	}
}

func resolveModelRegistryURL(sessionCtx context.Context, namespace string, serviceName string, client integrations.KubernetesClientInterface, config config.EnvConfig) (string, error) {

	serviceDetails, err := client.GetServiceDetailsByName(sessionCtx, namespace, serviceName)
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
		namespace := r.URL.Query().Get(string(constants.NamespaceHeaderParameterKey))
		if namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing required query parameter: %s", constants.NamespaceHeaderParameterKey))
			return
		}

		ctx := context.WithValue(r.Context(), constants.NamespaceHeaderParameterKey, namespace)
		r = r.WithContext(ctx)

		next(w, r, ps)
	}
}

func (app *App) PerformSARonGetListServicesByNamespace(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		user, ok := r.Context().Value(constants.KubeflowUserIdKey).(string)
		if !ok || user == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing user in context"))
			return
		}
		namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		var userGroups []string
		if groups, ok := r.Context().Value(constants.KubeflowUserGroupsKey).([]string); ok {
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

		user, ok := r.Context().Value(constants.KubeflowUserIdKey).(string)
		if !ok || user == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing user in context"))
			return
		}

		namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
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
		if groups, ok := r.Context().Value(constants.KubeflowUserGroupsKey).([]string); ok {
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
