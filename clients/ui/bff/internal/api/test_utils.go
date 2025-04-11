package api

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/kubeflow/model-registry/ui/bff/internal/config"
	"github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/mrserver"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"

	"github.com/kubeflow/model-registry/ui/bff/internal/constants"
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/repositories"
)

func setupApiTest[T any](method string, url string, body interface{}, k8Factory kubernetes.KubernetesClientFactory, requestIdentity kubernetes.RequestIdentity, namespace string) (T, *http.Response, error) {
	mockMRClient, err := mocks.NewModelRegistryClient(nil)
	if err != nil {
		return *new(T), nil, err
	}

	mockClient := new(mocks.MockHTTPClient)

	cfg := config.EnvConfig{
		AuthMethod: config.AuthMethodInternal,
	}
	//if token is set, use token auth
	if requestIdentity.Token != "" {
		cfg.AuthMethod = config.AuthMethodUser
	}
	testApp := App{
		repositories:            repositories.NewRepositories(mockMRClient),
		kubernetesClientFactory: k8Factory,
		logger:                  slog.Default(),
		config:                  cfg,
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

	// Set the kubeflow-userid header (middleware work)
	if requestIdentity.UserID != "" {
		req.Header.Set(constants.KubeflowUserIDHeader, requestIdentity.UserID)
	}

	ctx := mocks.NewMockSessionContext(req.Context())

	ctx = context.WithValue(ctx, constants.ModelRegistryHttpClientKey, mockClient)
	ctx = context.WithValue(ctx, constants.RequestIdentityKey, requestIdentity)
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
