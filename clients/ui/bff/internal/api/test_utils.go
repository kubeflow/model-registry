package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
)

func setupApiTest[T any](method string, url string, body interface{}, k8sClient k8s.KubernetesClientInterface, kubeflowUserIDHeaderValue string, namespace string) (T, *http.Response, error) {
	mockMRClient, err := mocks.NewModelRegistryClient(nil)
	if err != nil {
		return *new(T), nil, err
	}

	mockClient := new(mocks.MockHTTPClient)

	testApp := App{
		repositories:     repositories.NewRepositories(mockMRClient),
		kubernetesClient: k8sClient,
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

	// Set the kubeflow-userid header
	req.Header.Set(constants.KubeflowUserIDHeader, kubeflowUserIDHeaderValue)

	ctx := mocks.NewMockSessionContext(req.Context())
	ctx = context.WithValue(ctx, constants.ModelRegistryHttpClientKey, mockClient)
	ctx = context.WithValue(ctx, constants.KubeflowUserIdKey, kubeflowUserIDHeaderValue)
	ctx = context.WithValue(ctx, constants.NamespaceHeaderParameterKey, namespace)
	mrHttpClient := k8s.HTTPClient{
		ModelRegistryID: "model-registry",
	}
	ctx = context.WithValue(ctx, constants.ModelRegistryHttpClientKey, mrHttpClient)
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

func resolveStaticAssetsDirOnTests() string {
	// Fall back to finding project root for testing
	projectRoot, err := findProjectRootOnTests()
	if err != nil {
		panic("Failed to find project root: ")
	}

	return filepath.Join(projectRoot, "static")
}

// on tests findProjectRoot searches for the project root by locating go.mod
func findProjectRootOnTests() (string, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Traverse up until go.mod is found
	for currentDir != "/" {
		if _, err := os.Stat(filepath.Join(currentDir, "go.mod")); err == nil {
			return currentDir, nil
		}
		currentDir = filepath.Dir(currentDir)
	}

	return "", os.ErrNotExist
}
