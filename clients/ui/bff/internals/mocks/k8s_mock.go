package mocks

import (
	"github.com/stretchr/testify/mock"
)

type KubernetesClientMock struct {
	mock.Mock
}

func (m *KubernetesClientMock) FetchServiceNamesByComponent(componentValue string) ([]string, error) {
	args := m.Called(componentValue)
	return args.Get(0).([]string), args.Error(1)
}
