package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidationMiddleware(t *testing.T) {
	// Create a dummy handler that just returns 200 OK
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			// In tests, we don't need to handle this error, but satisfy linter
			return
		}
	})

	// Wrap it with our middleware
	handler := ValidationMiddleware(dummyHandler)

	testCases := []struct {
		name           string
		queryParams    string
		rawQuery       string // For cases where we need to manually set RawQuery
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid query parameters",
			queryParams:    "name=test&externalId=123",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "empty query parameters",
			queryParams:    "",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "URL encoded null byte",
			queryParams:    "name=test%00invalid",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request
			var req *http.Request
			if tc.rawQuery != "" {
				req = httptest.NewRequest("GET", "/test", nil)
				req.URL.RawQuery = tc.rawQuery
			} else {
				req = httptest.NewRequest("GET", "/test?"+tc.queryParams, nil)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check response body
			if tc.expectedStatus == http.StatusBadRequest {
				// For error responses, parse JSON and check structure
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Bad Request", response["code"])
				assert.Contains(t, response["message"], "contains null bytes which are not allowed")
			} else if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, rr.Body.String())
			}
		})
	}

	// Additional test for direct null byte injection
	t.Run("direct null byte in query", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/test", nil)
		// Manually set RawQuery with null byte
		req.URL.RawQuery = "name=test" + string([]byte{0}) + "invalid"

		rr := httptest.NewRecorder()
		handler.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusBadRequest, rr.Code)

		var response map[string]string
		err := json.Unmarshal(rr.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, "Bad Request", response["code"])
		assert.Contains(t, response["message"], "contains null bytes which are not allowed")
	})
}

func TestValidationMiddleware_RequestBody(t *testing.T) {
	// Create a dummy handler that just returns 200 OK
	dummyHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		if _, err := w.Write([]byte("OK")); err != nil {
			// In tests, we don't need to handle this error, but satisfy linter
			return
		}
	})

	// Wrap it with our middleware
	handler := ValidationMiddleware(dummyHandler)

	testCases := []struct {
		name           string
		body           string
		contentType    string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "valid JSON body",
			body:           `{"name": "test", "description": "valid description"}`,
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "empty body",
			body:           "",
			contentType:    "application/json",
			expectedStatus: http.StatusOK,
			expectedBody:   "OK",
		},
		{
			name:           "body with null byte",
			body:           "{\"name\": \"test\x00invalid\", \"description\": \"test\"}",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "body with embedded null byte in JSON string",
			body:           "{\"name\": \"test\x00invalid\"}",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:           "body with JSON unicode null byte escape sequence",
			body:           "{\"name\": \"\\u00c9B\\u0000b\\u0001\\u00cc\"}",
			contentType:    "application/json",
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create request with body
			var req *http.Request
			if tc.body != "" {
				req = httptest.NewRequest("POST", "/test", strings.NewReader(tc.body))
				req.Header.Set("Content-Type", tc.contentType)
			} else {
				req = httptest.NewRequest("POST", "/test", nil)
			}

			// Create response recorder
			rr := httptest.NewRecorder()

			// Call the handler
			handler.ServeHTTP(rr, req)

			// Check status code
			assert.Equal(t, tc.expectedStatus, rr.Code)

			// Check response body
			if tc.expectedStatus == http.StatusBadRequest {
				// For error responses, parse JSON and check structure
				var response map[string]string
				err := json.Unmarshal(rr.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Equal(t, "Bad Request", response["code"])
				// Check for either direct null bytes or JSON escape sequence errors
				message := response["message"]
				assert.True(t,
					strings.Contains(message, "null bytes which are not allowed") ||
						strings.Contains(message, "null byte escape sequences"),
					"Expected error message about null bytes, got: %s", message)
			} else if tc.expectedBody != "" {
				assert.Equal(t, tc.expectedBody, rr.Body.String())
			}
		})
	}
}
