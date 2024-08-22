package mocks

import (
	k8s "github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

type KubernetesClientMock struct {
	mock.Mock
}

func NewKubernetesClient(logger *slog.Logger) (k8s.KubernetesClientInterface, error) {
	return &KubernetesClientMock{}, nil
}

func (m *KubernetesClientMock) GetServiceNames() ([]string, error) {
	return []string{"model-registry", "model-registry-dora", "model-registry-bella"}, nil
}

func (m *KubernetesClientMock) GetServiceDetails() ([]k8s.ServiceDetails, error) {
	return []k8s.ServiceDetails{
		{Name: "model-registry", Description: "Model registry description", DisplayName: "Model Registry"},
		{Name: "model-registry-dora", Description: "Model registry dora description", DisplayName: "Model Registry Dora"},
		{Name: "model-registry-bella", Description: "Model registry bella description", DisplayName: "Model Registry Bella"},
	}, nil
}

func (m *KubernetesClientMock) BearerToken() (string, error) {
	return "FAKE BEARER TOKEN", nil
}

func (m *KubernetesClientMock) GetServiceDetailsByName(serviceName string) (k8s.ServiceDetails, error) {
	//expected forward to docker compose -f docker-compose.yaml up
	return k8s.ServiceDetails{
		Name:      serviceName,
		ClusterIP: "127.0.0.1",
		HTTPPort:  8080,
	}, nil
}
