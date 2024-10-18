package mocks

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"log/slog"
	"os"
	"testing"
)

func TestGetServiceDetails(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	k8sClient, err := NewKubernetesClient(logger)
	require.NoError(t, err, "Failed to initialize KubernetesClientMock")

	services, err := k8sClient.GetServiceDetails()
	require.NoError(t, err, "Failed to get service details")

	// Check that all services have the modified ClusterIP and HTTPPort
	for _, service := range services {
		assert.Equal(t, "127.0.0.1", service.ClusterIP, "ClusterIP should be set to 127.0.0.1")
		assert.Equal(t, int32(8080), service.HTTPPort, "HTTPPort should be set to 8080")
	}

	//Check that a specific service exists
	foundService := false
	for _, service := range services {
		if service.Name == "model-registry" {
			foundService = true
			assert.Equal(t, "Model Registry", service.DisplayName)
			assert.Equal(t, "Model Registry Description", service.Description)
			break
		}
	}
	assert.True(t, foundService, "Expected to find service 'model-registry'")
}

func TestGetServiceDetailsByName(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	k8sClient, err := NewKubernetesClient(logger)
	require.NoError(t, err, "Failed to initialize KubernetesClientMock")

	service, err := k8sClient.GetServiceDetailsByName("model-registry-dora")
	require.NoError(t, err, "Failed to get service details")

	assert.Equal(t, "model-registry-dora", service.Name)
	assert.Equal(t, "Model Registry Dora description", service.Description)
	assert.Equal(t, "Model Registry Dora", service.DisplayName)

}

func TestGetService(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	k8sClient, err := NewKubernetesClient(logger)
	require.NoError(t, err, "Failed to initialize KubernetesClientMock")

	services, err := k8sClient.GetServiceNames()
	require.NoError(t, err, "Failed to get service details")

	assert.Equal(t, "model-registry", services[0])
	assert.Equal(t, "model-registry-bella", services[1])
	assert.Equal(t, "model-registry-dora", services[2])

}
