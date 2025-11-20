package k8mocks

import (
	"context"
	"fmt"
	"log/slog"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type InternalKubernetesClientMock struct {
	*k8s.InternalKubernetesClient
}

// newMockedInternalKubernetesClientFromClientset creates a mock from existing envtest clientset
func newMockedInternalKubernetesClientFromClientset(clientset kubernetes.Interface, restConfig *rest.Config, logger *slog.Logger) k8s.KubernetesClientInterface {
	// Create dynamic client from rest config
	dynamicClient, err := dynamic.NewForConfig(restConfig)
	if err != nil {
		logger.Error("failed to create dynamic client for mock", "error", err)
		// Return mock without dynamic client - tests may fail if they need it
		return &InternalKubernetesClientMock{
			InternalKubernetesClient: &k8s.InternalKubernetesClient{
				SharedClientLogic: k8s.SharedClientLogic{
					Client: clientset,
					Logger: logger,
				},
			},
		}
	}

	return &InternalKubernetesClientMock{
		InternalKubernetesClient: &k8s.InternalKubernetesClient{
			SharedClientLogic: k8s.SharedClientLogic{
				Client:        clientset,
				DynamicClient: dynamicClient,
				Logger:        logger,
			},
		},
	}
}

// GetServiceDetails overrides to simulate ClusterIP for localhost access
func (m *InternalKubernetesClientMock) GetServiceDetails(sessionCtx context.Context, namespace string) ([]k8s.ServiceDetails, error) {
	originalServices, err := m.InternalKubernetesClient.GetServiceDetails(sessionCtx, namespace)
	if err != nil {
		return nil, fmt.Errorf("failed to get service details: %w", err)
	}

	for i := range originalServices {
		originalServices[i].ClusterIP = "127.0.0.1"
		originalServices[i].HTTPPort = 8080
		originalServices[i].IsHTTPS = false
	}

	return originalServices, nil
}

// GetServiceDetailsByName overrides to simulate local service access
func (m *InternalKubernetesClientMock) GetServiceDetailsByName(sessionCtx context.Context, namespace, serviceName string, serviceType string) (k8s.ServiceDetails, error) {
	originalService, err := m.InternalKubernetesClient.GetServiceDetailsByName(sessionCtx, namespace, serviceName, serviceType)
	if err != nil {
		return k8s.ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}
	originalService.ClusterIP = "127.0.0.1"
	originalService.HTTPPort = 8080
	originalService.IsHTTPS = false
	return originalService, nil
}

// BearerToken always returns a fake token for tests
func (m *InternalKubernetesClientMock) BearerToken() (string, error) {
	return "FAKE-BEARER-TOKEN", nil
}

func (kc *InternalKubernetesClientMock) GetGroups(ctx context.Context) ([]string, error) {
	return []string{"dora-group-mock", "bella-group-mock"}, nil
}
