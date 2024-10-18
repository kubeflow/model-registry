package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
)

func setupApiTest[T any](method string, url string, body interface{}) (T, *http.Response, error) {
	mockMRClient, err := mocks.NewModelRegistryClient(nil)
	if err != nil {
		return *new(T), nil, err
	}
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockK8sClient, err := mocks.NewKubernetesClient(logger)
	if err != nil {
		return *new(T), nil, err
	}

	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		repositories:     repositories.NewRepositories(mockMRClient),
		kubernetesClient: mockK8sClient,
	}

	var req *http.Request
	if body != nil {
		r, err := json.Marshal(body)
		if err != nil {
			return *new(T), nil, err
		}
		bytes.NewReader(r)
		req, err = http.NewRequest(method, url, bytes.NewReader(r))
		if err != nil {
			return *new(T), nil, err
		}
	} else {
		req, err = http.NewRequest(method, url, nil)
		if err != nil {
			return *new(T), nil, err
		}
	}

	ctx := context.WithValue(req.Context(), httpClientKey, mockClient)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	testApp.Routes().ServeHTTP(rr, req)

	rs := rr.Result()
	defer rs.Body.Close()
	respBody, err := io.ReadAll(rs.Body)
	if err != nil {
		return *new(T), nil, err
	}

	var entity T
	err = json.Unmarshal(respBody, &entity)
	if err != nil {
		if err == io.EOF {
			// There's no body to parse.
			return *new(T), rs, nil
		}
		return *new(T), nil, err
	}

	return entity, rs, nil
}
