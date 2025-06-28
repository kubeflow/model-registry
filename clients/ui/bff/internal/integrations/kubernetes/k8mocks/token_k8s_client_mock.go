package k8mocks

import (
	"context"
	"fmt"
	"log/slog"

	k8s "github.com/kubeflow/model-registry/ui/bff/internal/integrations/kubernetes"
	"k8s.io/client-go/kubernetes"
)

// ⚠️ WHY THIS FILE EXISTS:
// envtest does NOT support real authentication or token evaluation.
// It allows you to simulate Kubernetes behavior, but all requests use the test client's identity (usually cluster-admin).
// So, we simulate token-based behavior by mapping FAKE tokens to preconfigured test users.

type TokenKubernetesClientMock struct {
	*k8s.TokenKubernetesClient
}

func newMockedTokenKubernetesClientFromClientset(clientset kubernetes.Interface, logger *slog.Logger) k8s.KubernetesClientInterface {
	return &TokenKubernetesClientMock{
		TokenKubernetesClient: &k8s.TokenKubernetesClient{
			SharedClientLogic: k8s.SharedClientLogic{
				Client: clientset,
				Logger: logger,
				Token:  k8s.NewBearerToken(""), // Unused because impersonation is already handled in the client config
			},
		},
	}
}

// GetServiceDetails overrides to simulate ClusterIP for localhost access
func (m *TokenKubernetesClientMock) GetServiceDetails(sessionCtx context.Context, namespace string) ([]k8s.ServiceDetails, error) {
	originalServices, err := m.TokenKubernetesClient.GetServiceDetails(sessionCtx, namespace)
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
func (m *TokenKubernetesClientMock) GetServiceDetailsByName(sessionCtx context.Context, namespace, serviceName string) (k8s.ServiceDetails, error) {
	originalService, err := m.TokenKubernetesClient.GetServiceDetailsByName(sessionCtx, namespace, serviceName)
	if err != nil {
		return k8s.ServiceDetails{}, fmt.Errorf("failed to get service details: %w", err)
	}
	originalService.ClusterIP = "127.0.0.1"
	originalService.HTTPPort = 8080
	originalService.IsHTTPS = false
	return originalService, nil
}

// BearerToken always returns a fake token for tests
func (m *TokenKubernetesClientMock) BearerToken() (string, error) {
	return "FAKE-BEARER-TOKEN", nil
}

func (kc *TokenKubernetesClientMock) GetGroups(ctx context.Context) ([]string, error) {
	return []string{"dora-group-mock", "bella-group-mock"}, nil
}
