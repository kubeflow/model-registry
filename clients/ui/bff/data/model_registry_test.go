package data

import (
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchAllModelRegistry(t *testing.T) {
	mockK8sClient := new(mocks.KubernetesClientMock)

	mockK8sClient.On("GetServiceNames").Return(mockK8sClient.MockServiceNames(), nil)

	model := ModelRegistryModel{}

	registries, err := model.FetchAllModelRegistry(mockK8sClient)

	assert.NoError(t, err)

	expectedRegistries := []ModelRegistryModel{
		{Name: mockK8sClient.MockServiceNames()[0]},
		{Name: mockK8sClient.MockServiceNames()[1]},
	}
	assert.Equal(t, expectedRegistries, registries)

	mockK8sClient.AssertExpectations(t)
}
