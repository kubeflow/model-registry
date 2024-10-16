package repositories

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestFetchAllModelRegistry(t *testing.T) {
	mockK8sClient, _ := mocks.NewKubernetesClient(nil)

	mrClient := NewModelRegistryRepository()

	registries, err := mrClient.FetchAllModelRegistries(mockK8sClient)

	assert.NoError(t, err)

	expectedRegistries := []models.ModelRegistryModel{
		{Name: "model-registry", Description: "Model registry description", DisplayName: "Model Registry"},
		{Name: "model-registry-dora", Description: "Model registry dora description", DisplayName: "Model Registry Dora"},
		{Name: "model-registry-bella", Description: "Model registry bella description", DisplayName: "Model Registry Bella"},
	}
	assert.Equal(t, expectedRegistries, registries)

}
