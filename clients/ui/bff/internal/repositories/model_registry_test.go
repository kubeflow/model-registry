package repositories

import (
	"github.com/kubeflow/model-registry/ui/bff/internal/mocks"
	"github.com/kubeflow/model-registry/ui/bff/internal/models"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"os"
	"testing"
)

func TestFetchAllModelRegistry(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	mockK8sClient, _ := mocks.NewKubernetesClient(logger)

	mrClient := NewModelRegistryRepository()

	registries, err := mrClient.FetchAllModelRegistries(mockK8sClient)

	assert.NoError(t, err)

	expectedRegistries := []models.ModelRegistryModel{
		{Name: "model-registry", Description: "Model Registry Description", DisplayName: "Model Registry"},
		{Name: "model-registry-bella", Description: "Model Registry Bella description", DisplayName: "Model Registry Bella"},
		{Name: "model-registry-dora", Description: "Model Registry Dora description", DisplayName: "Model Registry Dora"},
	}
	assert.Equal(t, expectedRegistries, registries)
}
