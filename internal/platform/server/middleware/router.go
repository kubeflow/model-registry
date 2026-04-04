package middleware

import "net/http"

// WrapWithValidation wraps an http.Handler with custom validation middleware
func WrapWithValidation(handler http.Handler) http.Handler {
	return ValidationMiddleware(handler)
}
