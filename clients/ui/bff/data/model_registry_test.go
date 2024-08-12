package data

import (
	"github.com/kubeflow/model-registry/ui/bff/internals/mocks"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchAllModelRegistry(t *testing.T) {
	mockK8sClient, _ := mocks.NewKubernetesClient(nil)

	model := ModelRegistryModel{}

	registries, err := model.FetchAllModelRegistries(mockK8sClient)

	assert.NoError(t, err)

	expectedRegistries := []ModelRegistryModel{
		{Name: "model-registry", Description: "Model registry description", DisplayName: "Model Registry"},
		{Name: "model-registry-dora", Description: "Model registry dora description", DisplayName: "Model Registry Dora"},
		{Name: "model-registry-bella", Description: "Model registry bella description", DisplayName: "Model Registry Bella"},
	}
	assert.Equal(t, expectedRegistries, registries)

}
