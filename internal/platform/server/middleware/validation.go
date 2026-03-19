package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/golang/glog"
)

func validateStringParameter(paramName, paramValue string) error {
	if paramValue == "" {
		return nil
	}

	if checkNullBytesInString(paramValue) {
		return fmt.Errorf("invalid %s parameter: contains null bytes which are not allowed", paramName)
	}

	return nil
}

func checkNullBytesInString(s string) bool {
	return strings.Contains(s, "\x00") || strings.Contains(s, "\\u0000")
}

// ValidationMiddleware validates all query parameters and request body for unsafe characters
func ValidationMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queryParams, err := url.ParseQuery(r.URL.RawQuery)
		if err != nil {
			// If we can't parse query parameters, let the request continue.
			// The parsing error will be handled elsewhere.
			next.ServeHTTP(w, r)
			return
		}

		for paramName, values := range queryParams {
			for _, value := range values {
				if err := validateStringParameter(paramName, value); err != nil {
					returnValidationError(w, fmt.Sprintf("Invalid %s query parameter: %v", paramName, err))
					return
				}
			}
		}

		if r.Body != nil && r.ContentLength != 0 {
			bodyBytes, err := io.ReadAll(r.Body)
			if err != nil {
				glog.Errorf("Error reading request body: %v", err)
				next.ServeHTTP(w, r)
				return
			}

			r.Body = io.NopCloser(bytes.NewReader(bodyBytes))

			bodyContent := string(bodyBytes)
			if checkNullBytesInString(bodyContent) {
				returnValidationError(w, "Request body contains null bytes which are not allowed")
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func returnValidationError(w http.ResponseWriter, message string) {
	glog.Errorf("Validation error: %s", message)

	errorResponse := struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	}{
		Code:    "Bad Request",
		Message: message,
	}

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusBadRequest)
	if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
		glog.Errorf("Error encoding JSON error response: %v", err)
	}
}
