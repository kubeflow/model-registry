package api

import (
	"context"
	"errors"
	"fmt"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
	"log/slog"
	"net/http"
	"runtime/debug"
	"strings"

	"github.com/google/uuid"
	"github.com/julienschmidt/httprouter"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	helper "github.com/kubeflow/model-registry/ui/bff/internal/helpers"
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

func (app *App) InjectRequestIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//skip use headers check if we are not on /api/v1 (i.e. we are on /healthcheck and / (static fe files) )
		if !strings.HasPrefix(r.URL.Path, ApiPathPrefix) && !strings.HasPrefix(r.URL.Path, PathPrefix+ApiPathPrefix) {
			next.ServeHTTP(w, r)
			return
		}

		var identity *kubernetes.RequestIdentity

		switch app.config.AuthMethod {

		case config.AuthMethodInternal:
			userID := r.Header.Get(constants.KubeflowUserIDHeader)
			//`kubeflow-userid`: Contains the user's email address.
			if userID == "" {
				app.badRequestResponse(w, r, errors.New("missing required header on AuthMethodInternal: kubeflow-userid"))
				return
			}

			userGroupsHeader := r.Header.Get(constants.KubeflowUserGroupsIdHeader)
			// Note: The functionality for `kubeflow-groups` is not fully operational at Kubeflow platform at this time
			// but it's supported on Model Registry BFF
			//`kubeflow-groups`: Holds a comma-separated list of user groups.
			groups := []string{}
			if userGroupsHeader != "" {
				for _, g := range strings.Split(userGroupsHeader, ",") {
					groups = append(groups, strings.TrimSpace(g))
				}
			}
			identity = &kubernetes.RequestIdentity{
				UserID: userID,
				Groups: groups,
			}
		case config.AuthMethodUser:
			token := r.Header.Get(constants.XForwardedAccessTokenHeader)
			if token == "" {
				app.badRequestResponse(w, r, errors.New("missing required header on AuthMethodUser: access token"))
				return
			}
			identity = &kubernetes.RequestIdentity{
				Token: token,
			}
		default:
			app.badRequestResponse(w, r, fmt.Errorf("invalid auth method: %s", app.config.AuthMethod))
			return
		}

		ctx := context.WithValue(r.Context(), constants.RequestIdentityKey, identity)
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

		client, err := app.kubernetesClientFactory.GetClient(r.Context())
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
			return
		}

		modelRegistry, err := app.repositories.ModelRegistry.GetModelRegistry(r.Context(), client, namespace, modelRegistryID)
		if err != nil {
			app.notFoundResponse(w, r)
			return
		}
		modelRegistryBaseURL := modelRegistry.ServerAddress

		// If we are in dev mode, we need to resolve the server address to the local host
		// to allow the client to connect to the model registry via port forwarded from the cluster to the local machine.
		if app.config.DevMode {
			modelRegistryBaseURL = app.repositories.ModelRegistry.ResolveServerAddress("localhost", int32(app.config.DevModePort))
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

		restHttpClient, err := mrserver.NewHTTPClient(restClientLogger, modelRegistryID, modelRegistryBaseURL)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to create Kubernetes client: %v", err))
			return
		}
		ctx := context.WithValue(r.Context(), constants.ModelRegistryHttpClientKey, restHttpClient)
		next(w, r.WithContext(ctx), ps)
	}
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

func (app *App) RequireListServiceAccessInNamespace(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := r.Context()
		identity, ok := ctx.Value(constants.RequestIdentityKey).(*kubernetes.RequestIdentity)
		if !ok || identity == nil {
			app.badRequestResponse(w, r, fmt.Errorf("missing RequestIdentity in context"))
			return
		}

		if err := validateRequestIdentity(identity, app.config.AuthMethod); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		client, err := app.kubernetesClientFactory.GetClient(ctx)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
			return
		}

		allowed, err := client.CanListServicesInNamespace(ctx, identity, namespace)
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

func (app *App) RequireAccessToService(next func(http.ResponseWriter, *http.Request, httprouter.Params)) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

		ctx := r.Context()
		identity, ok := ctx.Value(constants.RequestIdentityKey).(*kubernetes.RequestIdentity)
		if !ok || identity == nil {
			app.badRequestResponse(w, r, fmt.Errorf("missing RequestIdentity in context"))
			return
		}

		if err := validateRequestIdentity(identity, app.config.AuthMethod); err != nil {
			app.badRequestResponse(w, r, err)
			return
		}

		namespace, ok := r.Context().Value(constants.NamespaceHeaderParameterKey).(string)
		if !ok || namespace == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		serviceName := ps.ByName(ModelRegistryId)
		if !ok || serviceName == "" {
			app.badRequestResponse(w, r, fmt.Errorf("missing namespace in context"))
			return
		}

		client, err := app.kubernetesClientFactory.GetClient(ctx)
		if err != nil {
			app.serverErrorResponse(w, r, fmt.Errorf("failed to get Kubernetes client: %w", err))
			return
		}

		allowed, err := client.CanAccessServiceInNamespace(r.Context(), identity, namespace, serviceName)

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

// ValidateRequestIdentity ensures the identity contains required values based on auth method.
func validateRequestIdentity(identity *kubernetes.RequestIdentity, authMethod string) error {
	if identity == nil {
		return errors.New("missing identity")
	}

	switch authMethod {
	case config.AuthMethodInternal:
		if identity.UserID == "" {
			return errors.New("user ID (kubeflow-userid) required for internal authentication")
		}
	case config.AuthMethodUser:
		if identity.Token == "" {
			return errors.New("token is required for token-based authentication")
		}
	default:
		return fmt.Errorf("unsupported authentication method: %s", authMethod)
	}

	return nil
}
