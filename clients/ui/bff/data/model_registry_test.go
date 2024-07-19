package data

import (
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchAllModelRegistry(t *testing.T) {
	mockK8sClient, _ := mocks.NewKubernetesClient(nil)

	model := ModelRegistryModel{}

	registries, err := model.FetchAllModelRegistry(mockK8sClient)

	assert.NoError(t, err)

	expectedRegistries := []ModelRegistryModel{
		{Name: "model-registry"},
		{Name: "model-registry-dora"},
		{Name: "model-registry-bella"},
	}
	assert.Equal(t, expectedRegistries, registries)

}
