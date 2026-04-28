package middleware

import (
	"net/http"

	platformmw "github.com/kubeflow/hub/internal/platform/server/middleware"
	"github.com/kubeflow/hub/internal/server/openapi"
)

// WrapWithValidation wraps the auto-generated router with custom validation middleware
func WrapWithValidation(routers ...openapi.Router) http.Handler {
	// Create the auto-generated router
	baseRouter := openapi.NewRouter(routers...)

	// Wrap it with our custom validation middleware
	return platformmw.ValidationMiddleware(baseRouter)
}
