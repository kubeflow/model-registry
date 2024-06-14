package mocks

import (
	k8s "github.com/kubeflow/model-registry/ui/bff/integrations"
	"github.com/stretchr/testify/mock"
)

type KubernetesClientMock struct {
	mock.Mock
}

func (m *KubernetesClientMock) GetServiceNames() ([]string, error) {
	args := m.Called()
	return args.Get(0).([]string), args.Error(1)
}

func (m *KubernetesClientMock) BearerToken() (string, error) {
	args := m.Called()
	return args.String(0), args.Error(1)
}

func (m *KubernetesClientMock) GetServiceDetailsByName(serviceName string) (k8s.ServiceDetails, error) {
	args := m.Called(serviceName)
	return args.Get(0).(k8s.ServiceDetails), args.Error(1)
}
func (m *KubernetesClientMock) MockServiceNames() []string {
	return []string{"model-registry-dora", "model-registry-bella"}
}
