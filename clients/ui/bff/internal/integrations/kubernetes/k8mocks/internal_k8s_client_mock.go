package k8mocks

import (
	"context"
	"fmt"
	"log/slog"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"k8s.io/client-go/kubernetes"
)

type InternalKubernetesClientMock struct {
	*k8s.InternalKubernetesClient
}

// newMockedInternalKubernetesClientFromClientset creates a mock from existing envtest clientset
func newMockedInternalKubernetesClientFromClientset(clientset kubernetes.Interface, logger *slog.Logger) k8s.KubernetesClientInterface {
	return &InternalKubernetesClientMock{
		InternalKubernetesClient: &k8s.InternalKubernetesClient{
			SharedClientLogic: k8s.SharedClientLogic{
				Client: clientset,
				Logger: logger,
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
	}

	return originalServices, nil
}

// GetServiceDetailsByName overrides to simulate local service access
func (m *InternalKubernetesClientMock) GetServiceDetailsByName(sessionCtx context.Context, namespace, serviceName string) (k8s.ServiceDetails, error) {
	originalService, err := m.InternalKubernetesClient.GetServiceDetailsByName(sessionCtx, namespace, serviceName)
	if err != nil {
		return k8s.ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}
	originalService.ClusterIP = "127.0.0.1"
	originalService.HTTPPort = 8080
	return originalService, nil
}

// BearerToken always returns a fake token for tests
func (m *InternalKubernetesClientMock) BearerToken() (string, error) {
	return "FAKE-BEARER-TOKEN", nil
}

func (kc *InternalKubernetesClientMock) GetGroups(ctx context.Context) ([]string, error) {
	return []string{"dora-group-mock", "bella-group-mock"}, nil
}
